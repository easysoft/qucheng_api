// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package app

import (
	"context"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm"
	"reflect"
	"strings"

	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service/app/component"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/metric"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/tlog"

	"helm.sh/helm/v3/pkg/release"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
)

type Instance struct {
	ctx context.Context

	clusterName string
	namespace   string
	name        string

	selector labels.Selector

	ks *cluster.Cluster

	release *release.Release

	ChartName           string
	CurrentChartVersion string
}

func newApp(ctx context.Context, am *Manager, name string) *Instance {

	i := &Instance{
		ctx:         ctx,
		clusterName: am.clusterName, namespace: am.namespace, name: name,
		ks: am.ks,
	}

	i.release = i.fetchRelease()
	return i
}

func (i *Instance) prepare() {
	i.ChartName = i.release.Chart.Metadata.Name
	i.CurrentChartVersion = i.release.Chart.Metadata.Version

	i.selector = labels.Set{"release": i.name}.AsSelector()
}

func (i *Instance) fetchRelease() *release.Release {
	getter := &ReleaseGetter{
		namespace: i.namespace,
		store:     i.ks.Store,
	}
	rel, err := getter.Last(i.name)
	if err != nil {
		tlog.WithCtx(i.ctx).ErrorS(err, "parse release failed")
	}
	return rel
}

func (i *Instance) getServices() ([]*v1.Service, error) {
	return i.ks.Store.ListServices(i.namespace, i.selector)
}

func (i *Instance) getComponents() *component.Components {
	components := component.NewComponents()

	deployments, _ := i.ks.Store.ListDeployments(i.namespace, i.selector)

	if len(deployments) >= 1 {
		for _, d := range deployments {
			components.Add(component.NewDeployComponent(d, i.ks))
		}
	}

	statefulsets, _ := i.ks.Store.ListStatefulSets(i.namespace, i.selector)

	if len(statefulsets) >= 1 {
		for _, s := range statefulsets {
			components.Add(component.NewStatefulsetComponent(s, i.ks))
		}
	}

	return components
}

func (i *Instance) ParseStatus() *model.AppRespStatus {
	var settingStopped = false
	components := i.getComponents()

	data := &model.AppRespStatus{
		Components: make([]model.AppRespStatusComponent, 0),
		Status:     constant.AppStatusMap[constant.AppStatusUnknown],
		Version:    i.CurrentChartVersion,
		Age:        0,
	}

	if len(components.Items()) == 0 {
		return data
	}

	stopped, err := i.settingParse("global.stopped")
	if err == nil {
		if stop, ok := stopped.(bool); ok {
			settingStopped = stop
		}
	}

	for _, c := range components.Items() {
		resC := model.AppRespStatusComponent{
			Name:       c.Name(),
			Kind:       c.Kind(),
			Replicas:   c.Replicas(),
			StatusCode: c.Status(settingStopped),
			Status:     constant.AppStatusMap[c.Status(settingStopped)],
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

func (i *Instance) ListIngressHosts() []string {
	var hosts []string
	ingresses, err := i.ks.Store.ListIngresses(i.namespace, i.selector)
	if err != nil {
		return hosts
	}
	for _, ing := range ingresses {
		for _, rule := range ing.Spec.Rules {
			hosts = append(hosts, rule.Host)
		}
	}
	return hosts
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
				if p.Name == constant.ServicePortWeb || p.Name == "http" {
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

func (i *Instance) settingParse(path string) (interface{}, error) {
	var err error
	var ok bool
	var node map[string]interface{}
	var data interface{}

	frames := strings.Split(path, ".")
	node = i.release.Config

	if len(frames) > 1 {
		for _, frame := range frames[0 : len(frames)-1] {
			n, ok := node[frame]
			if !ok {
				err = ErrPathParseFailed
				break
			}
			ntype := reflect.TypeOf(n)
			if ntype.Kind() != reflect.Map {
				err = ErrPathParseFailed
				break
			}
			node = n.(map[string]interface{})
		}
		if err != nil {
			return nil, err
		}
		data, ok = node[frames[len(frames)-1]]
	} else {
		data, ok = node[path]
	}

	if !ok {
		return nil, ErrPathParseFailed
	}

	return data, nil
}

func (i *Instance) GetComponents() []model.Component {
	var components []model.Component

	components = append(components, model.Component{Name: i.ChartName})
	deps, err := helm.ParseChartDependencies(i.release.Chart)
	if err != nil {
		return components
	}

	for _, dep := range deps {
		components = append(components, model.Component{
			Name: dep.Name,
		})
	}

	return components
}

func (i *Instance) GetSchemaCategories(component string) interface{} {
	var data []string

	if component == i.ChartName {
		data = helm.ParseChartCategories(i.release.Chart)
	} else {
		var exist = false
		for _, dep := range i.release.Chart.Lock.Dependencies {
			if dep.Name == component {
				exist = true
				break
			}
		}

		if exist {
			ch, err := helm.GetChart(genChart("test", i.ChartName), i.CurrentChartVersion)
			if err != nil {
				return data
			}

			for _, dep := range ch.Dependencies() {
				if dep.Name() == component {
					data = helm.ParseChartCategories(dep)
					break
				}
			}
		}
	}

	return data
}

func (i *Instance) GetSchema(component, category string) string {
	var data string

	if component == i.ChartName {
		jbody, err := helm.ReadSchemaFromChart(i.release.Chart, category, "test")
		tlog.WithCtx(i.ctx).InfoS("get schema content", "data", string(jbody))
		if err != nil {
			tlog.WithCtx(i.ctx).ErrorS(err, "get schema failed")
			return ""
		}
		data = string(jbody)
	} else {
		ch, err := helm.GetChart(genChart("test", i.ChartName), i.CurrentChartVersion)
		if err != nil {
			return ""
		}
		for _, dep := range ch.Dependencies() {
			if dep.Name() == component {
				jbody, err := helm.ReadSchemaFromChart(dep, category, "test")
				if err != nil {
					return ""
				}
				data = string(jbody)
			}
		}
	}
	return data
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
	dst.Cpu = resource.NewQuantity(0, resource.DecimalExponent)
	dst.Memory = resource.NewQuantity(0, resource.DecimalExponent)

	for _, m := range metrics {
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
