package node

import (
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sort"
)

type Manager struct {
	clusterName string
	ks          *cluster.Cluster
}

func NewNodes(clusterName string) *Manager {
	return &Manager{
		clusterName: clusterName,
		ks:          cluster.Get(clusterName),
	}
}

func (m *Manager) filteNodes(selector labels.Selector) ([]*v1.Node, error) {
	return m.ks.Store.ListNodes(selector)
}

func (m *Manager) ListNodePortIPS() []string {
	ips := make([]string, 0)
	nodes, err := m.filteNodes(labels.NewSelector())
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
