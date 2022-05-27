package app

import (
	"context"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service/app/component"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/metric"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
)

type Instance struct {
	ctx context.Context

	clusterName string
	namespace   string
	name        string

	selector   labels.Selector
	components *component.Components

	ks *cluster.Cluster
}

func newApp(ctx context.Context, am *Manager, name string) *Instance {
	return &Instance{
		ctx:         ctx,
		clusterName: am.clusterName, namespace: am.namespace, name: name,
		ks: am.ks,
	}
}

func (i *Instance) getServices() ([]*v1.Service, error) {
	return i.ks.Store.ListServices(i.namespace, i.selector)
}

func (i *Instance) Components() *component.Components {
	return i.components
}

func (i *Instance) ParseStatus() *model.AppRespStatus {
	data := &model.AppRespStatus{
		Components: make([]model.AppRespStatusComponent, 0),
		Status:     constant.AppStatusMap[constant.AppStatusUnknown],
		Age:        0,
	}

	if len(i.components.Items()) == 0 {
		return data
	}

	for _, c := range i.components.Items() {
		resC := model.AppRespStatusComponent{
			Name:       c.Name(),
			Kind:       c.Kind(),
			Replicas:   c.Replicas(),
			StatusCode: c.Status(),
			Status:     constant.AppStatusMap[c.Status()],
			Age:        c.Age(),
		}
		data.Components = append(data.Components, resC)
	}

	minStatusCode := data.Components[0].StatusCode
	maxAge := data.Components[0].Age
	for _, comp := range data.Components {
		if comp.StatusCode < minStatusCode {
			minStatusCode = comp.StatusCode
		}

		if comp.Age > maxAge {
			maxAge = comp.Age
		}
	}

	data.Status = constant.AppStatusMap[minStatusCode]
	data.Age = maxAge
	return data
}

func (i *Instance) ParseNodePort() int32 {
	var nodePort int32 = 0
	services, err := i.getServices()
	if err != nil {
		return nodePort
	}

	for _, s := range services {
		if s.Spec.Type == v1.ServiceTypeNodePort {
			for _, p := range s.Spec.Ports {
				if p.Name == constant.ServicePortWeb {
					nodePort = p.NodePort
					break
				}
			}
		}
	}

	return nodePort
}

func (i *Instance) Settings() *Settings {
	return newSettings(i)
}

func (i *Instance) GetMetrics() *model.AppMetric {
	metrics := i.ks.Metric.ListPodMetrics(i.ctx, i.namespace, i.selector)
	pods, _ := i.ks.Store.ListPods(i.namespace, i.selector)

	var usage metric.Res
	var limit metric.Res

	sumPodUsage(&usage, metrics)
	sumPodLimit(&limit, pods)

	memUsage, _ := usage.Memory.AsInt64()
	memLimit, _ := limit.Memory.AsInt64()

	data := model.AppMetric{
		Cpu: model.ResourceCpu{
			Usage: usage.Cpu.AsApproximateFloat64(), Limit: limit.Cpu.AsApproximateFloat64(),
		},
		Memory: model.ResourceMemory{
			Usage: memUsage, Limit: memLimit,
		},
	}
	return &data
}

func sumPodUsage(dst *metric.Res, metrics []*metric.Res) {
	count := len(metrics)

	if count == 0 {
		return
	}

	dst.Cpu = metrics[0].Cpu
	dst.Memory = metrics[0].Memory

	for _, m := range metrics[1:] {
		dst.Cpu.Add(m.Cpu.DeepCopy())
		dst.Memory.Add(m.Memory.DeepCopy())
	}
}

func sumPodLimit(dst *metric.Res, pods []*v1.Pod) {
	dst.Cpu = resource.NewQuantity(0, resource.DecimalExponent)
	dst.Memory = resource.NewQuantity(0, resource.DecimalExponent)

	for _, pod := range pods {
		for _, ctn := range pod.Spec.Containers {
			l := ctn.Resources.Limits
			dst.Cpu.Add(*l.Cpu())
			dst.Memory.Add(*l.Memory())
		}
	}
}
