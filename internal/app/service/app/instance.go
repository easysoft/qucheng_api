package app

import (
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service/app/component"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type Instance struct {
	clusterName string
	namespace   string
	name        string

	selector   labels.Selector
	components *component.Components

	ks *cluster.Cluster
}

func newApp(am *Manager, name string) *Instance {
	return &Instance{
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
