package cluster

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/store"
)

var kubeClusters = make(map[string]*Cluster)

type Cluster struct {
	Config  rest.Config
	Store   *store.Storer
	Clients *store.Clients
	inner   bool
	primary bool
}

func (c *Cluster) IsInner() bool {
	return c.inner
}

func (c *Cluster) IsPrimary() bool {
	return c.primary
}

func Exist(name string) bool {
	_, ok := kubeClusters[name]
	if ok {
		return true
	}
	return false
}

func Get(name string) (*Cluster, error) {
	if !Exist(name) {
		return nil, &NotFound{Name: name}
	}

	c, _ := kubeClusters[name]
	return c, nil
}

func add(name string, config rest.Config, inner, primary bool) error {
	if Exist(name) {
		return &AlreadyRegistered{Name: name}
	}

	s := store.NewStorer(config)
	cluster := &Cluster{
		Config:  config,
		Store:   s,
		Clients: s.Clients,
		inner:   inner,
		primary: primary,
	}

	kubeClusters[name] = cluster
	return nil
}

func Init(stopChan chan struct{}) error {
	restCfg, err := loadPrimaryCluster()
	if err != nil {
		return err
	}

	if err = add("primary", *restCfg, true, true); err != nil {
		return err
	}

	for _, c := range kubeClusters {
		go c.Store.Run(stopChan)
	}
	return nil
}

func loadPrimaryCluster() (*rest.Config, error) {
	restCfg, err := rest.InClusterConfig()
	if err == nil {
		return restCfg, nil
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(userHome, ".kube", "config")

	restCfg, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}

	return restCfg, nil
}