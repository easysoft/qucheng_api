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

	_, err = h.Install(name, body.Chart)
	return err
}
