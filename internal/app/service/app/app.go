// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package app

import (
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm"
)

type AppListManager struct {
	cluster   string
	namespace string
}

func NewApps(clusterName, namespace string) *AppListManager {
	return &AppListManager{
		cluster: clusterName, namespace: namespace,
	}
}

func (am *AppListManager) Install(name string, body model.AppCreateModel) error {
	h, err := helm.NamespaceScope(am.namespace)
	if err != nil {
		return err
	}

	_, err = h.Install(name, defaultChartRepo + "/" + body.Chart)
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
