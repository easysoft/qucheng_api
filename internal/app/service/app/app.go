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
	"k8s.io/klog/v2"
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

	var settings = make([]string, len(body.Settings))
	for _, s := range body.Settings {
		settings = append(settings, s.Key+"="+s.Val)
	}
	klog.Infoln(settings)
	_, err = h.Install(name, genRepo(body.Channel)+"/"+body.Chart, settings)
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
	app.components = component.NewComponents()

	selector := labels.NewSelector()
	label1, _ := labels.NewRequirement("app.kubernetes.io/managed-by", selection.Equals, []string{"Helm"})
	labelRelease, _ := labels.NewRequirement("release", selection.Equals, []string{name})
	selector = selector.Add(*label1, *labelRelease)

	app.selector = selector

	deployments, err := am.ks.Store.ListDeployments(am.namespace, selector)
	if err != nil {
		return nil, err
	}
	if len(deployments) >= 1 {
		for _, d := range deployments {
			app.components.Add(component.NewDeployComponent(d, app.ks))
		}
	}

	statefulsets, err := am.ks.Store.ListStatefulSets(am.namespace, selector)
	if err != nil {
		return nil, err
	}

	if len(statefulsets) >= 1 {
		for _, s := range statefulsets {
			app.components.Add(component.NewStatefulsetComponent(s, app.ks))
		}
	}

	return app, nil
}

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

func (a *Instance) getServices() ([]*v1.Service, error) {
	return a.ks.Store.ListServices(a.namespace, a.selector)
}

func (a *Instance) Components() *component.Components {
	return a.components
}

func (a *Instance) ParseStatus() *model.AppRespStatus {
	data := &model.AppRespStatus{
		Components: make([]model.AppRespStatusComponent, 0),
		Status:     constant.AppStatusMap[constant.AppStatusUnknown],
		Age:        0,
	}

	if len(a.components.Items()) == 0 {
		return data
	}

	for _, c := range a.components.Items() {
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
