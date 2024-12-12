package shell

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	abstractions "github.com/beam-cloud/beta9/pkg/abstractions/common"
	"github.com/beam-cloud/beta9/pkg/common"
	"github.com/beam-cloud/beta9/pkg/network"
	"github.com/beam-cloud/beta9/pkg/repository"
	"github.com/beam-cloud/beta9/pkg/types"
)

const (
	requestProcessingInterval time.Duration = time.Millisecond * 100
)

type request struct {
	ctx       echo.Context
	payload   *types.TaskPayload
	task      *ShellTask
	done      chan bool
	processed bool
}

type container struct {
	id               string
	address          string
	inFlightRequests int
}

type RequestBuffer struct {
	ctx                     context.Context
	httpClient              *http.Client
	tailscale               *network.Tailscale
	tsConfig                types.TailscaleConfig
	stubId                  string
	stubConfig              *types.StubConfigV1
	workspace               *types.Workspace
	rdb                     *common.RedisClient
	containerRepo           repository.ContainerRepository
	buffer                  *abstractions.RingBuffer[*request]
	availableContainers     []container
	availableContainersLock sync.RWMutex
	maxTokens               int
	isASGI                  bool
}

func NewRequestBuffer(
	ctx context.Context,
	rdb *common.RedisClient,
	workspace *types.Workspace,
	stubId string,
	size int,
	containerRepo repository.ContainerRepository,
	stubConfig *types.StubConfigV1,
	tailscale *network.Tailscale,
	tsConfig types.TailscaleConfig,
	isASGI bool,
) *RequestBuffer {
	b := &RequestBuffer{
		ctx:                     ctx,
		rdb:                     rdb,
		workspace:               workspace,
		stubId:                  stubId,
		stubConfig:              stubConfig,
		buffer:                  abstractions.NewRingBuffer[*request](size),
		availableContainers:     []container{},
		availableContainersLock: sync.RWMutex{},
		containerRepo:           containerRepo,
		httpClient:              &http.Client{},
		tailscale:               tailscale,
		tsConfig:                tsConfig,
		maxTokens:               int(stubConfig.Workers),
		isASGI:                  isASGI,
	}

	if stubConfig.ConcurrentRequests > 1 && isASGI {
		// Floor is set to the number of workers
		b.maxTokens = max(int(stubConfig.ConcurrentRequests), b.maxTokens)
	}

	go b.discoverContainers()
	go b.processRequests()

	return b
}

func (rb *RequestBuffer) ForwardRequest(ctx echo.Context, task *ShellTask) error {
	done := make(chan bool)
	req := &request{
		ctx:  ctx,
		done: done,
		payload: &types.TaskPayload{
			Args:   task.msg.Args,
			Kwargs: task.msg.Kwargs,
		},
		task: task,
	}
	rb.buffer.Push(req, false)

	for {
		select {
		case <-rb.ctx.Done():
			return nil
		case <-ctx.Request().Context().Done():
			if !req.processed {
				rb.cancelInFlightTask(req.task)
			}
			return nil
		case <-done:
			return nil
		}
	}
}

func (rb *RequestBuffer) processRequests() {
	for {
		select {
		case <-rb.ctx.Done():
			return
		default:
			if len(rb.availableContainers) == 0 {
				time.Sleep(requestProcessingInterval)
				continue
			}

			req, ok := rb.buffer.Pop()
			if !ok {
				time.Sleep(requestProcessingInterval)
				continue
			}

			go rb.handleRequest(req)
		}
	}
}

func (rb *RequestBuffer) checkAddressIsReady(address string) bool {
	httpClient, err := rb.getHttpClient(address)
	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(rb.ctx, 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/health", address), nil)
	if err != nil {
		return false
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (rb *RequestBuffer) discoverContainers() {
	for {
		select {
		case <-rb.ctx.Done():
			return
		default:
			containerStates, err := rb.containerRepo.GetActiveContainersByStubId(rb.stubId)
			if err != nil {
				continue
			}

			var wg sync.WaitGroup
			availableContainersChan := make(chan container, len(containerStates))

			for _, containerState := range containerStates {
				wg.Add(1)

				go func(cs types.ContainerState) {
					defer wg.Done()
					if cs.Status != types.ContainerStatusRunning {
						return
					}

					containerAddress, err := rb.containerRepo.GetContainerAddress(cs.ContainerId)
					if err != nil {
						return
					}

					availableTokens, err := rb.requestTokens(cs.ContainerId)
					if err != nil {
						return
					}

					// Let's say we have 5 workers available, and there are three tokens left in this bucket
					// that means we currently have 5-3 -> 2 requests in flight
					inFlightRequests := rb.maxTokens - availableTokens

					if rb.checkAddressIsReady(containerAddress) {
						availableContainersChan <- container{
							id:               cs.ContainerId,
							address:          containerAddress,
							inFlightRequests: inFlightRequests,
						}

						return
					}
				}(containerState)
			}

			wg.Wait()
			close(availableContainersChan)

			// Collect available containers
			availableContainers := make([]container, 0)
			for c := range availableContainersChan {
				availableContainers = append(availableContainers, c)
			}

			// Sort availableContainers by # of in-flight requests (ascending)
			sort.Slice(availableContainers, func(i, j int) bool {
				return availableContainers[i].inFlightRequests < availableContainers[j].inFlightRequests
			})

			rb.availableContainersLock.Lock()
			rb.availableContainers = availableContainers
			rb.availableContainersLock.Unlock()

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (rb *RequestBuffer) requestTokens(containerId string) (int, error) {
	tokenKey := Keys.shellRequestTokens(rb.workspace.Name, rb.stubId, containerId)

	val, err := rb.rdb.Get(rb.ctx, tokenKey).Int()
	if err != nil && err != redis.Nil {
		return 0, err
	} else if err == redis.Nil {
		created, err := rb.rdb.SetNX(rb.ctx, tokenKey, rb.maxTokens, 0).Result()
		if err != nil {
			return 0, err
		}

		if created {
			return rb.maxTokens, nil
		}

		tokens, err := rb.rdb.Get(rb.ctx, tokenKey).Int()
		if err != nil {
			return 0, err
		}

		return tokens, nil
	}

	return val, nil
}

func (rb *RequestBuffer) acquireRequestToken(containerId string) error {
	tokenKey := Keys.shellRequestTokens(rb.workspace.Name, rb.stubId, containerId)
	tokenCount, err := rb.rdb.Decr(rb.ctx, tokenKey).Result()
	if err != nil {
		return err
	}

	// If the token count is negative, we exceeded our threshold of
	// available request tokens, just reverse the operation
	if tokenCount < 0 {
		rb.rdb.Incr(rb.ctx, tokenKey)
		return errors.New("too many in-flight requests")
	}

	err = rb.rdb.Expire(rb.ctx, tokenKey, time.Duration(rb.stubConfig.TaskPolicy.Timeout)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rb *RequestBuffer) releaseRequestToken(containerId string) error {
	// TODO: if a gateway crashes before releasing the token, it could lead to a drift
	// in the count of available request tokens for a particular container. To handle this
	// we could move the release logic to the task implementation (e.g. task.Complete), so that
	// it handles the release of the token and is not tied to a specific gateway

	tokenKey := Keys.shellRequestTokens(rb.workspace.Name, rb.stubId, containerId)

	err := rb.rdb.Incr(rb.ctx, tokenKey).Err()
	if err != nil {
		return err
	}

	err = rb.rdb.Expire(rb.ctx, tokenKey, time.Duration(rb.stubConfig.TaskPolicy.Timeout)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rb *RequestBuffer) getHttpClient(address string) (*http.Client, error) {
	// If it isn't an tailnet address, just return the standard http client
	if !rb.tsConfig.Enabled || !strings.Contains(address, rb.tsConfig.HostName) {
		return rb.httpClient, nil
	}

	conn, err := network.ConnectToHost(rb.ctx, address, types.RequestTimeoutDurationS, rb.tailscale, rb.tsConfig)
	if err != nil {
		return nil, err
	}

	// Create a custom transport that uses the established connection
	// Either using tailscale or not
	transport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return conn, nil
		},
	}

	client := &http.Client{
		Transport: transport,
	}

	return client, nil
}

func (rb *RequestBuffer) handleRequest(req *request) {
	rb.availableContainersLock.RLock()

	if len(rb.availableContainers) == 0 {
		rb.buffer.Push(req, true)
		rb.availableContainersLock.RUnlock()
		return
	}

	// Select an available container to forward the request to (whichever one has the lowest # of inflight requests)
	// Basically least-connections load balancing
	c := rb.availableContainers[0]

	rb.availableContainersLock.RUnlock()

	err := rb.acquireRequestToken(c.id)
	if err != nil {
		rb.buffer.Push(req, true)
		return
	}
	defer rb.afterRequest(req, c.id)

	req.processed = true
	if req.ctx.IsWebSocket() {
		rb.handleWSRequest(req, c)
	} else {
		rb.handleHttpRequest(req, c)
	}
}

func (rb *RequestBuffer) handleWSRequest(req *request, c container) {
	dstDialer := websocket.Dialer{
		NetDialContext: network.GetDialer(c.address, rb.tailscale, rb.tsConfig),
	}

	err := rb.proxyWebsocketConnection(
		req,
		c,
		dstDialer,
		fmt.Sprintf("ws://%s/%s", c.address, req.ctx.Param("subPath")),
	)
	if err != nil {
		return
	}
}

func (rb *RequestBuffer) handleHttpRequest(req *request, c container) {
	request := req.ctx.Request()

	requestBody := request.Body
	if !rb.isASGI {
		b, err := json.Marshal(req.payload)
		if err != nil {
			return
		}
		requestBody = io.NopCloser(bytes.NewReader(b))
	}

	httpClient, err := rb.getHttpClient(c.address)
	if err != nil {
		return
	}

	containerUrl := fmt.Sprintf("http://%s/%s", c.address, req.ctx.Param("subPath"))

	// Forward query params to the container if ASGI
	if rb.isASGI {
		queryParams := req.ctx.QueryString()
		if queryParams != "" {
			containerUrl += "?" + queryParams
		}
	}

	httpReq, err := http.NewRequestWithContext(request.Context(), request.Method, containerUrl, requestBody)
	if err != nil {
		req.ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
		return
	}

	// Copy headers to new request
	for key, values := range request.Header {
		for _, val := range values {
			httpReq.Header.Add(key, val)
		}
	}

	httpReq.Header.Add("X-TASK-ID", req.task.msg.TaskId) // Add task ID to header
	go rb.heartBeat(req, c.id)                           // Send heartbeat via redis for duration of request

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		req.ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
		return
	}

	defer resp.Body.Close()

	// Write response headers
	for key, values := range resp.Header {
		for _, value := range values {
			req.ctx.Response().Writer.Header().Add(key, value)
		}
	}

	// Write status code header
	req.ctx.Response().Writer.WriteHeader(resp.StatusCode)

	// Check if we can stream the response
	streamingSupported := true
	flusher, ok := req.ctx.Response().Writer.(http.Flusher)
	if !ok {
		streamingSupported = false
	}

	// Send response to client in chunks
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			req.ctx.Response().Writer.Write(buf[:n])

			if streamingSupported {
				flusher.Flush()
			}
		}

		if err != nil {
			if err != io.EOF {
				req.ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "Internal server error",
				})
			}

			break
		}
	}
}

func (rb *RequestBuffer) cancelInFlightTask(task *ShellTask) {
	task.Cancel(context.Background(), types.TaskCancellationReason(types.TaskRequestCancelled))
}

func (rb *RequestBuffer) heartBeat(req *request, containerId string) {
	ctx := req.ctx.Request().Context()
	ticker := time.NewTicker(shellRequestHeartbeatInterval)
	defer ticker.Stop()

	rb.rdb.Set(rb.ctx, Keys.shellRequestHeartbeat(rb.workspace.Name, rb.stubId, req.task.msg.TaskId), containerId, shellRequestHeartbeatInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-rb.ctx.Done():
			return
		case <-ticker.C:
			rb.rdb.Set(rb.ctx, Keys.shellRequestHeartbeat(rb.workspace.Name, rb.stubId, req.task.msg.TaskId), containerId, shellRequestHeartbeatInterval)
		}
	}
}

func (rb *RequestBuffer) afterRequest(req *request, containerId string) {
	defer func() {
		req.done <- true
	}()

	defer rb.releaseRequestToken(containerId)

	// Set keep warm lock
	if rb.stubConfig.KeepWarmSeconds == 0 {
		return
	}

	rb.rdb.SetEx(
		context.Background(),
		Keys.shellKeepWarmLock(rb.workspace.Name, rb.stubId, containerId),
		1,
		time.Duration(rb.stubConfig.KeepWarmSeconds)*time.Second,
	)
}

func (rb *RequestBuffer) proxyWebsocketConnection(r *request, c container, dialer websocket.Dialer, dstAddress string) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins
			return true
		},
	}

	w := r.ctx.Response().Writer
	req := r.ctx.Request()

	wsSrc, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return err
	}

	headers := http.Header{}
	headers.Add("X-TASK-ID", r.task.msg.TaskId) // Add task ID to header

	wsDst, resp, err := dialer.Dial(dstAddress, headers)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	go rb.heartBeat(r, c.id) // Send heartbeat via redis for duration of request
	go forwardWSConn(wsSrc.NetConn(), wsDst.NetConn())
	forwardWSConn(wsDst.NetConn(), wsSrc.NetConn())
	return nil
}

func forwardWSConn(src net.Conn, dst net.Conn) {
	defer func() {
		src.Close()
		dst.Close()
	}()

	_, err := io.Copy(src, dst)
	if err != nil {
		return
	}
}