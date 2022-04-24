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
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type AppListManager struct {
	clusterName string
	ks          *cluster.Cluster
	namespace   string
}

func NewApps(clusterName, namespace string) *AppListManager {
	return &AppListManager{
		clusterName: clusterName, namespace: namespace,
		ks: cluster.Get(clusterName),
	}
}

func (am *AppListManager) Install(name string, body model.AppCreateModel) error {
	h, err := helm.NamespaceScope(am.namespace)
	if err != nil {
		return err
	}

	_, err = h.Install(name, defaultChartRepo+"/"+body.Chart)
	return err
}

func (am *AppListManager) UnInstall(name string) error {
	h, err := helm.NamespaceScope(am.namespace)
	if err != nil {
		return err
	}

	err = h.Uninstall(name)
	return err
}

func (am *AppListManager) GetApp(name string) (*AppInstance, error) {
	app := newApp(am, name)
	app.componets = component.NewComponents()

	selector := labels.NewSelector()
	label1, _ := labels.NewRequirement("app.kubernetes.io/managed-by", selection.Equals, []string{"Helm"})
	//label2, _ := labels.NewRequirement("heritage", selection.Equals, []string{"Helm"})
	labelRelease, _ := labels.NewRequirement("release", selection.Equals, []string{name})
	selector = selector.Add(*label1, *labelRelease)

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

type AppInstance struct {
	clusterName string
	namespace   string
	name        string

	componets *component.Components

	ks *cluster.Cluster
}

func newApp(am *AppListManager, name string) *AppInstance {
	return &AppInstance{
		clusterName: am.clusterName, namespace: am.namespace, name: name,
		ks: am.ks,
	}
}

func (a *AppInstance) ParseStatus() *model.AppRespStatus {
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
