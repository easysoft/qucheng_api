package store

import (
	"fmt"
	"time"

	metav1 "k8s.io/api/core/v1"
	metanetworkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/core/v1"
	networkv1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

const (
	resyncPeriod = time.Minute * 10
)

type Informer struct {
	Namespaces  cache.SharedIndexInformer
	CloneSets   cache.SharedIndexInformer
	Pods        cache.SharedIndexInformer
	DivideRules cache.SharedIndexInformer
	Ingresses   cache.SharedIndexInformer
	Endpoints 	cache.SharedIndexInformer
}

func (i *Informer) Run(stopCh chan struct{}) {
	go i.Namespaces.Run(stopCh)
	go i.Pods.Run(stopCh)
	go i.Ingresses.Run(stopCh)
	go i.Endpoints.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh,
		i.Namespaces.HasSynced,
		i.Pods.HasSynced,
		i.Ingresses.HasSynced,
		i.Endpoints.HasSynced,
	) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
	}
}

type Lister struct {
	Namespaces  v1.NamespaceLister
	Pods        v1.PodLister
	Ingresses   networkv1.IngressLister
	Endpoints   v1.EndpointsLister
}

type Clients struct {
	Base   *kubernetes.Clientset
}

type Storer struct {
	informers *Informer
	listers   *Lister
	Clients   *Clients
}

func NewStorer(config rest.Config) *Storer {
	s := &Storer{
		informers: &Informer{},
		listers:   &Lister{},
		Clients:   &Clients{},
	}

	if cs, err := kubernetes.NewForConfig(&config); err != nil {
		klog.ErrorS(err, "failed to prepare kubeclient")
	} else {
		s.Clients.Base = cs
		factory := informers.NewSharedInformerFactoryWithOptions(cs, resyncPeriod)

		s.informers.Namespaces = factory.Core().V1().Namespaces().Informer()
		s.listers.Namespaces = factory.Core().V1().Namespaces().Lister()

		s.informers.Pods = factory.Core().V1().Pods().Informer()
		s.listers.Pods = factory.Core().V1().Pods().Lister()

		s.informers.Ingresses = factory.Networking().V1().Ingresses().Informer()
		s.listers.Ingresses = factory.Networking().V1().Ingresses().Lister()

		s.informers.Endpoints = factory.Core().V1().Endpoints().Informer()
		s.listers.Endpoints = factory.Core().V1().Endpoints().Lister()
	}

	return s
}

func (s *Storer) Run(stopCh chan struct{}) {
	s.informers.Run(stopCh)
}

func (s *Storer) GetNamespace(name string) (*metav1.Namespace, error) {
	return s.listers.Namespaces.Get(name)
}

func (s *Storer) ListNamespaces(selector labels.Selector) ([]*metav1.Namespace, error) {
	return s.listers.Namespaces.List(selector)
}

func (s *Storer) GetPod(namespace, name string) (*metav1.Pod, error) {
	return s.listers.Pods.Pods(namespace).Get(name)
}

func (s *Storer) ListPods(namespace string, selector labels.Selector) ([]*metav1.Pod, error) {
	return s.listers.Pods.Pods(namespace).List(selector)
}

func (s *Storer) GetIngress(namespace, name string) (*metanetworkv1.Ingress, error) {
	return s.listers.Ingresses.Ingresses(namespace).Get(name)
}

func (s *Storer) ListIngresses(namespace string, selector labels.Selector) ([]*metanetworkv1.Ingress, error) {
	return s.listers.Ingresses.Ingresses(namespace).List(selector)
}

func (s *Storer) GetEndpoint(namespace, name string) (*metav1.Endpoints, error) {
	return s.listers.Endpoints.Endpoints(namespace).Get(name)
}

func (s *Storer) ListEndpoints(namespace string, selector labels.Selector) ([]*metav1.Endpoints, error) {
	return s.listers.Endpoints.Endpoints(namespace).List(selector)
}
