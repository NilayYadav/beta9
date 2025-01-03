package shell

import (
	"math"

	abstractions "github.com/beam-cloud/beta9/pkg/abstractions/common"
)

type shellAutoscalerSample struct {
	TotalRequests     int64
	CurrentContainers int64
}

func shellSampleFunc(i *shellInstance) (*shellAutoscalerSample, error) {
	totalRequests, err := i.TaskRepo.TasksInFlight(i.Ctx, i.Workspace.Name, i.Stub.ExternalId)
	if err != nil {
		totalRequests = -1
	}

	currentContainers := 0
	state, err := i.Sta4te()
	if err != nil {
		currentContainers = -1
	}

	currentContainers = state.PendingContainers + state.RunningContainers

	sample := &shellAutoscalerSample{
		TotalRequests:     int64(totalRequests),
		CurrentContainers: int64(currentContainers),
	}

	return sample, nil
}

func shellDeploymentScaleFunc(i *shellInstance, s *shellAutoscalerSample) *abstractions.AutoscalerResult {
	desiredContainers := 0

	if s.TotalRequests == 0 {
		desiredContainers = 0
	} else {
		if s.TotalRequests == -1 {
			return &abstractions.AutoscalerResult{
				ResultValid: false,
			}
		}

		desiredContainers = int(s.TotalRequests / int64(i.StubConfig.Autoscaler.TasksPerContainer))
		if s.TotalRequests%int64(i.StubConfig.Autoscaler.TasksPerContainer) > 0 {
			desiredContainers += 1
		}

		// Limit max replicas to either what was set in autoscaler config, or the limit specified on the gateway config (whichever is lower)
		maxReplicas := math.Min(float64(i.StubConfig.Autoscaler.MaxContainers), float64(i.AppConfig.GatewayService.StubLimits.MaxReplicas))
		desiredContainers = int(math.Min(maxReplicas, float64(desiredContainers)))
	}

	return &abstractions.AutoscalerResult{
		DesiredContainers: desiredContainers,
		ResultValid:       true,
	}
}

// func shellServeScaleFunc(i *shellInstance, sample *shellAutoscalerSample) *abstractions.AutoscalerResult {
// 	desiredContainers := 1

// 	timeoutKey := Keys.shellServeLock(i.Workspace.Name, i.Stub.ExternalId)
// 	exists, err := i.Rdb.Exists(i.Ctx, timeoutKey).Result()
// 	if err != nil {
// 		return &abstractions.AutoscalerResult{
// 			ResultValid: false,
// 		}
// 	}

// 	if exists == 0 {
// 		desiredContainers = 0
// 	}

// 	return &abstractions.AutoscalerResult{
// 		DesiredContainers: desiredContainers,
// 		ResultValid:       true,
// 	}
// }
