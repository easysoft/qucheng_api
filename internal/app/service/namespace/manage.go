package namespace

import (
	"context"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Manager struct {
	clusterName string
	ks          *cluster.Cluster
}

func NewNamespaces(clusterName string) *Manager {
	return &Manager{
		clusterName: clusterName,
		ks:          cluster.Get(clusterName),
	}
}

func (m *Manager) Create(name string) error {
	newNS := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{
		Name:        name,
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	}}

	if _, err := m.ks.Clients.Base.CoreV1().Namespaces().Create(context.TODO(), newNS, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func (m *Manager) Recycle(name string) error {
	return m.ks.Clients.Base.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (m *Manager) Has(name string) bool {
	_, err := m.ks.Clients.Base.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	return err == nil
}
