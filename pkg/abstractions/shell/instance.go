package shell

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"

	abstractions "github.com/beam-cloud/beta9/pkg/abstractions/common"
	"github.com/beam-cloud/beta9/pkg/types"
)

type shellInstance struct {
	*abstractions.AutoscaledInstance
	buffer *RequestBuffer
	isASGI bool
}

func withAutoscaler(constructor func(i *shellInstance) *abstractions.Autoscaler[*shellInstance, *shellAutoscalerSample]) func(*shellInstance) {
	return func(i *shellInstance) {
		i.Autoscaler = constructor(i)
	}
}

func withEntryPoint(entryPoint func(instance *shellInstance) []string) func(*shellInstance) {
	return func(i *shellInstance) {
		i.EntryPoint = entryPoint(i)
	}
}

func (i *shellInstance) startContainers(containersToRun int) error {
	secrets, err := abstractions.ConfigureContainerRequestSecrets(i.Workspace, *i.buffer.stubConfig)
	if err != nil {
		return err
	}

	mounts, err := abstractions.ConfigureContainerRequestMounts(
		i.Stub.Object.ExternalId,
		i.Workspace,
		*i.buffer.stubConfig,
		i.Stub.ExternalId,
	)
	if err != nil {
		return err
	}

	env := []string{
		fmt.Sprintf("BETA9_TOKEN=%s", i.Token.Key),
		fmt.Sprintf("HANDLER=%s", i.StubConfig.Handler),
		fmt.Sprintf("ON_START=%s", i.StubConfig.OnStart),
		fmt.Sprintf("STUB_ID=%s", i.Stub.ExternalId),
		fmt.Sprintf("STUB_TYPE=%s", i.Stub.Type),
		fmt.Sprintf("WORKERS=%d", i.StubConfig.Workers),
		fmt.Sprintf("KEEP_WARM_SECONDS=%d", i.StubConfig.KeepWarmSeconds),
		fmt.Sprintf("PYTHON_VERSION=%s", i.StubConfig.PythonVersion),
		fmt.Sprintf("CALLBACK_URL=%s", i.StubConfig.CallbackUrl),
		fmt.Sprintf("TIMEOUT=%d", i.StubConfig.TaskPolicy.Timeout),
	}

	env = append(secrets, env...)

	gpuRequest := types.GpuTypesToStrings(i.StubConfig.Runtime.Gpus)
	if i.StubConfig.Runtime.Gpu != "" {
		gpuRequest = append(gpuRequest, i.StubConfig.Runtime.Gpu.String())
	}

	gpuCount := i.StubConfig.Runtime.GpuCount
	if i.StubConfig.RequiresGPU() && gpuCount == 0 {
		gpuCount = 1
	}

	checkpointEnabled := i.StubConfig.CheckpointEnabled
	if i.Stub.Type.IsServe() {
		checkpointEnabled = false
	}

	if gpuCount > 1 {
		checkpointEnabled = false
	}

	for c := 0; c < containersToRun; c++ {
		containerId := i.genContainerId()

		runRequest := &types.ContainerRequest{
			ContainerId:       containerId,
			Env:               env,
			Cpu:               i.StubConfig.Runtime.Cpu,
			Memory:            i.StubConfig.Runtime.Memory,
			GpuRequest:        gpuRequest,
			GpuCount:          uint32(gpuCount),
			ImageId:           i.StubConfig.Runtime.ImageId,
			StubId:            i.Stub.ExternalId,
			WorkspaceId:       i.Workspace.ExternalId,
			Workspace:         *i.Workspace,
			EntryPoint:        i.EntryPoint,
			Mounts:            mounts,
			Stub:              *i.Stub,
			CheckpointEnabled: checkpointEnabled,
		}

		err := i.Scheduler.Run(runRequest)
		if err != nil {
			log.Printf("<%s> unable to run container: %v", i.Name, err)
			return err
		}

		continue
	}

	return nil
}

func (i *shellInstance) stopContainers(containersToStop int) error {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	containerIds, err := i.stoppableContainers()
	if err != nil {
		return err
	}

	for c := 0; c < containersToStop && len(containerIds) > 0; c++ {
		idx := rnd.Intn(len(containerIds))
		containerId := containerIds[idx]

		err := i.Scheduler.Stop(&types.StopContainerArgs{ContainerId: containerId})
		if err != nil {
			log.Printf("<%s> unable to stop container: %v", i.Name, err)
			return err
		}

		// Remove the containerId from the containerIds slice to avoid
		// sending multiple stop requests to the same container
		containerIds = append(containerIds[:idx], containerIds[idx+1:]...)
	}

	return nil
}

func (i *shellInstance) stoppableContainers() ([]string, error) {
	containers, err := i.ContainerRepo.GetActiveContainersByStubId(i.Stub.ExternalId)
	if err != nil {
		return nil, err
	}

	// Create a slice to hold the keys
	keys := make([]string, 0, len(containers))
	for _, container := range containers {
		if container.Status == types.ContainerStatusStopping || container.Status == types.ContainerStatusPending {
			continue
		}

		// When deployment is stopped, all containers should be stopped even if they have keep warm
		if !i.IsActive {
			keys = append(keys, container.ContainerId)
			continue
		}
		keys = append(keys, container.ContainerId)
	}

	return keys, nil
}

func (i *shellInstance) genContainerId() string {
	return fmt.Sprintf("%s-%s-%s", shellContainerPrefix, i.Stub.ExternalId, uuid.New().String()[:8])
}