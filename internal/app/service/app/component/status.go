package component

import (
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"

	v1 "k8s.io/api/core/v1"
)

func parseStatus(replicas, availableReplicas, updatedReplicas, readyReplicas int32,
	pods []*v1.Pod) (appStatus constant.AppStatusType) {

	appStatus = constant.AppStatusUnknown
	if replicas == 0 {
		appStatus = constant.AppStatusStop
		return
	}

	if replicas > 0 && readyReplicas < replicas {
		appStatus = constant.AppStatusStartup
		return
	}

	if updatedReplicas == replicas && readyReplicas == replicas {
		appStatus = constant.AppStatusRunning
		return
	}

	for _, pod := range pods {
		for _, ctnStatus := range pod.Status.ContainerStatuses {
			if !*ctnStatus.Started {
				if ctnStatus.State.Waiting != nil && ctnStatus.State.Waiting.Reason == "CrashLoopBackOff" {
					appStatus = constant.AppStatusAbnormal
					break
				}
			}
		}
	}
	return
}
