// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package app

import (
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm"
	v1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/klog/v2"
)

type AppListManager struct {
	clusterName   string
	ks 				*cluster.Cluster
	namespace string
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

func (am *AppListManager) GetApp(name string) (*AppInstance, error) {
	return newApp(am, name), nil
}

func (am *AppListManager) GetAppWithOutHelm(name string) (bool, error) {
	selector := labels.NewSelector()
	reqName, _ := labels.NewRequirement("name", selection.Equals, []string{name})
	reqOwner, _ := labels.NewRequirement("owner", selection.Equals, []string{"helm"})
	selector = selector.Add(*reqName, *reqOwner)
	secrets, err := am.ks.Store.ListSecrets(am.namespace, selector)
	return len(secrets) > 0, err
}

type AppInstance struct {
	clusterName   string
	namespace string
	name		string

	object 	*appsv1.Deployment
	ks *cluster.Cluster
}

func newApp(am *AppListManager, name string) *AppInstance {
	return &AppInstance{
		clusterName: am.clusterName, namespace: am.namespace, name: name,
		ks: am.ks,
	}
}

func (a *AppInstance) podList() ([]*v1.Pod, error) {
	matchLabels := a.object.Spec.Selector.MatchLabels

	pods, err := a.ks.Store.ListPods(a.namespace, labels.SelectorFromValidatedSet(matchLabels))
	return pods, err
}

// need to support statefulset
func (a *AppInstance) ParseStatus(data *model.AppRespStatus) {
	//data.Replicas = a.object.Status.Replicas
	data.Status = a.parseStatus().String()
	//data.ReadyReplicas = a.object.Status.AvailableReplicas
}

func (a *AppInstance) parseStatus() (appStatus constant.AppStatusType) {
	appStatus = constant.AppStatusUnknown
	if a.object.Status.Replicas == 0 {
		appStatus = constant.AppStatusStop
		return
	}

	if a.object.Status.Replicas > 0 &&  a.object.Status.AvailableReplicas < a.object.Status.Replicas {
		appStatus = constant.AppStatusStop
		return
	}

	if a.object.Status.UpdatedReplicas == a.object.Status.Replicas &&
		a.object.Status.ReadyReplicas == a.object.Status.Replicas {
		appStatus = constant.AppStatusRunning
		return
	}

	pods, err := a.podList()
	if err != nil {
		appStatus = constant.AppStatusAbnormal
		klog.ErrorS(err, "get pod list failed", "app", a.name, "namespace", a.namespace)
		return
	}
	for _, pod := range pods {
		for _, ctnStatus := range pod.Status.ContainerStatuses {
			if !*ctnStatus.Started && ctnStatus.Image == a.object.Spec.Template.Spec.Containers[0].Image {
				if ctnStatus.State.Waiting != nil && ctnStatus.State.Waiting.Reason == "CrashLoopBackOff" {
					appStatus = constant.AppStatusAbnormal
				}
				break
			}
		}
	}
	return
}
