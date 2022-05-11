package component

import (
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	metaappsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type Statefulset struct {
	name   string
	kind   string
	object *metaappsv1.StatefulSet

	pods []*v1.Pod
	ks   *cluster.Cluster
}

func NewStatefulsetComponent(obj *metaappsv1.StatefulSet, ks *cluster.Cluster) Component {
	return &Statefulset{
		name: obj.Name, kind: KindStatefulSet,
		object: obj, ks: ks,
	}
}

func (s *Statefulset) Name() string {
	return s.name
}

func (s *Statefulset) Kind() string {
	return s.kind
}

func (s *Statefulset) Replicas() int32 {
	return s.object.Status.Replicas
}

func (s *Statefulset) Age() int64 {
	return parseOldestAge(s.getPods())
}

func (s *Statefulset) Status() constant.AppStatusType {
	status := s.object.Status
	return parseStatus(status.Replicas, status.AvailableReplicas, status.UpdatedReplicas, status.ReadyReplicas, s.getPods())
}

func (s *Statefulset) getPods() []*v1.Pod {
	matchLabels := s.object.Spec.Selector.MatchLabels

	pods, _ := s.ks.Store.ListPods(s.object.Namespace, labels.SelectorFromValidatedSet(matchLabels))
	return pods
}
