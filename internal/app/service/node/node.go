// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package node

import (
	"context"
	"sort"

	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	v1 "k8s.io/api/core/v1"
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
