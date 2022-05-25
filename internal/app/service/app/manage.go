// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package app

import (
	"context"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service/app/component"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/tlog"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type Manager struct {
	ctx context.Context

	clusterName string
	ks          *cluster.Cluster
	namespace   string
}

func NewApps(ctx context.Context, clusterName, namespace string) *Manager {
	return &Manager{
		ctx:         ctx,
		clusterName: clusterName, namespace: namespace,
		ks: cluster.Get(clusterName),
	}
}

func (m *Manager) Install(name string, body model.AppCreateModel) error {
	h, err := helm.NamespaceScope(m.namespace)
	if err != nil {
		return err
	}

	var settings = make([]string, len(body.Settings))
	for _, s := range body.Settings {
		settings = append(settings, s.Key+"="+s.Val)
	}
	tlog.WithCtx(m.ctx).InfoS("build install settings", "namespace", m.namespace, "name", name, "settings", settings)
	_, err = h.Install(name, genChart(body.Channel, body.Chart), settings)
	return err
}

func (m *Manager) UnInstall(name string) error {
	h, err := helm.NamespaceScope(m.namespace)
	if err != nil {
		return err
	}

	err = h.Uninstall(name)
	return err
}

func (m *Manager) GetApp(name string) (*Instance, error) {
	app := newApp(m, name)
	app.components = component.NewComponents()

	selector := labels.NewSelector()
	label1, _ := labels.NewRequirement("app.kubernetes.io/managed-by", selection.Equals, []string{"Helm"})
	labelRelease, _ := labels.NewRequirement("release", selection.Equals, []string{name})
	selector = selector.Add(*label1, *labelRelease)

	app.selector = selector

	deployments, err := m.ks.Store.ListDeployments(m.namespace, selector)
	if err != nil {
		return nil, err
	}
	if len(deployments) >= 1 {
		for _, d := range deployments {
			app.components.Add(component.NewDeployComponent(d, app.ks))
		}
	}

	statefulsets, err := m.ks.Store.ListStatefulSets(m.namespace, selector)
	if err != nil {
		return nil, err
	}

	if len(statefulsets) >= 1 {
		for _, s := range statefulsets {
			app.components.Add(component.NewStatefulsetComponent(s, app.ks))
		}
	}

	if len(app.Components().Items()) == 0 {
		return nil, &ErrAppNotFound{Name: app.name}
	}

	return app, nil
}
