// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package app

import (
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service/app/component"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type Manager struct {
	clusterName string
	ks          *cluster.Cluster
	namespace   string
}

func NewApps(clusterName, namespace string) *Manager {
	return &Manager{
		clusterName: clusterName, namespace: namespace,
		ks: cluster.Get(clusterName),
	}
}

func (am *Manager) Install(name string, body model.AppCreateModel) error {
	h, err := helm.NamespaceScope(am.namespace)
	if err != nil {
		return err
	}

	_, err = h.Install(name, defaultChartRepo+"/"+body.Chart)
	return err
}

func (am *Manager) UnInstall(name string) error {
	h, err := helm.NamespaceScope(am.namespace)
	if err != nil {
		return err
	}

	err = h.Uninstall(name)
	return err
}

func (am *Manager) GetApp(name string) (*Instance, error) {
	app := newApp(am, name)
	app.componets = component.NewComponents()

	selector := labels.NewSelector()
	label1, _ := labels.NewRequirement("app.kubernetes.io/managed-by", selection.Equals, []string{"Helm"})
	labelRelease, _ := labels.NewRequirement("release", selection.Equals, []string{name})
	selector = selector.Add(*label1, *labelRelease)

	app.selector = selector

	deployments, err := am.ks.Store.ListDeployments(am.namespace, selector)
	if err != nil {
		return nil, err
	}
	if len(deployments) == 1 {
		deploy := deployments[0]
		app.componets.Add(component.NewDeployComponent(deploy, app.ks))
	}

	statefulsets, err := am.ks.Store.ListStatefulSets(am.namespace, selector)
	if err != nil {
		return nil, err
	}

	if len(statefulsets) == 1 {
		sts := statefulsets[0]
		app.componets.Add(component.NewStatefulsetComponent(sts, app.ks))
	}

	return app, nil
}

type Instance struct {
	clusterName string
	namespace   string
	name        string

	selector  labels.Selector
	componets *component.Components

	ks *cluster.Cluster
}

func newApp(am *Manager, name string) *Instance {
	return &Instance{
		clusterName: am.clusterName, namespace: am.namespace, name: name,
		ks: am.ks,
	}
}

func (a *Instance) getServices() ([]*v1.Service, error) {
	return a.ks.Store.ListServices(a.namespace, a.selector)
}

func (a *Instance) ParseStatus() *model.AppRespStatus {
	data := &model.AppRespStatus{
		Components: make([]model.AppRespStatusComponent, 0),
	}
	for _, c := range a.componets.Items() {
		resC := model.AppRespStatusComponent{
			Name:       c.Name(),
			Kind:       c.Kind(),
			Replicas:   c.Replicas(),
			StatusCode: c.Status(),
			Status:     constant.AppStatusMap[c.Status()],
		}
		data.Components = append(data.Components, resC)
	}

	minStatusCode := data.Components[0].StatusCode
	for _, comp := range data.Components {
		if comp.StatusCode < minStatusCode {
			minStatusCode = comp.StatusCode
		}
	}

	data.Status = constant.AppStatusMap[minStatusCode]
	return data
}

func (a *Instance) ParseNodePort() int32 {
	var nodePort int32 = 0
	services, err := a.getServices()
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
