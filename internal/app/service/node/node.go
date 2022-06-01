// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package node

import (
	"context"
	"sort"

	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/metric"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/tlog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
)

type Manager struct {
	ctx context.Context

	clusterName string
	ks          *cluster.Cluster
}

func NewNodes(ctx context.Context, clusterName string) *Manager {
	return &Manager{
		ctx:         ctx,
		clusterName: clusterName,
		ks:          cluster.Get(clusterName),
	}
}

func (m *Manager) filterNodes(selector labels.Selector) ([]*v1.Node, error) {
	return m.ks.Store.ListNodes(selector)
}

func (m *Manager) ListNodePortIPS() []string {
	ips := make([]string, 0)
	nodes, err := m.filterNodes(labels.NewSelector())
	if err != nil {
		return ips
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

	for _, node := range nodes {
		for _, ad := range node.Status.Addresses {
			if ad.Type == v1.NodeInternalIP {
				ips = append(ips, ad.Address)
			}
		}
	}

	return ips
}

func (m *Manager) GetNodes() []*v1.Node {
	var nodes []*v1.Node
	var err error
	nodes, err = m.ks.Store.ListNodes(labels.Everything())
	if err != nil {
		tlog.WithCtx(m.ctx).ErrorS(err, "get nodes failed")
	}
	return nodes
}

func (m *Manager) GetNodesMetrics() []*metric.Res {
	return m.ks.Metric.ListNodeMetrics(m.ctx)
}

func (m *Manager) Statistic() (model.NodeMetric, error) {
	nodes := m.GetNodes()
	metrics := m.GetNodesMetrics()

	var usage metric.Res
	var allocatable metric.Res
	var capacity metric.Res

	sumNodeUsage(&usage, metrics)
	sumNodeCapacity(&capacity, &allocatable, nodes)

	memUsage, _ := usage.Memory.AsInt64()
	memCapacity, _ := capacity.Memory.AsInt64()
	memAllocatable, _ := allocatable.Memory.AsInt64()

	metricData := model.NodeMetric{
		Cpu: model.NodeResourceCpu{
			Usage: usage.Cpu.AsApproximateFloat64(), Capacity: capacity.Cpu.AsApproximateFloat64(),
			Allocatable: allocatable.Cpu.AsApproximateFloat64(),
		},
		Memory: model.NodeResourceMemory{
			Usage: memUsage, Capacity: memCapacity, Allocatable: memAllocatable,
		},
	}

	return metricData, nil
}

func sumNodeUsage(dst *metric.Res, metrics []*metric.Res) {
	dst.Cpu = resource.NewQuantity(0, resource.DecimalExponent)
	dst.Memory = resource.NewQuantity(0, resource.DecimalExponent)

	for _, m := range metrics {
		dst.Cpu.Add(m.Cpu.DeepCopy())
		dst.Memory.Add(m.Memory.DeepCopy())
	}
}

func sumNodeCapacity(cap *metric.Res, allocat *metric.Res, nodes []*v1.Node) {
	cap.Cpu = resource.NewQuantity(0, resource.DecimalExponent)
	cap.Memory = resource.NewQuantity(0, resource.DecimalExponent)
	allocat.Cpu = resource.NewQuantity(0, resource.DecimalExponent)
	allocat.Memory = resource.NewQuantity(0, resource.DecimalExponent)

	for _, node := range nodes {
		cap.Cpu.Add(*node.Status.Capacity.Cpu())
		cap.Memory.Add(*node.Status.Capacity.Memory())

		allocat.Cpu.Add(*node.Status.Allocatable.Cpu())
		allocat.Memory.Add(*node.Status.Allocatable.Memory())
	}
}
