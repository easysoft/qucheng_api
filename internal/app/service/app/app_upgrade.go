package app

import (
	"github.com/imdario/mergo"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm"
)

func (a *AppInstance) Stop(chart string) error {
	h, _ := helm.NamespaceScope(a.namespace)
	vals, err := h.GetValues(a.name)
	if err != nil {
		return err
	}

	updateMap := map[string]interface{}{
		"global": map[string]interface{}{
			"stoped": true,
		},
	}

	if err := mergo.Merge(&vals, updateMap, mergo.WithOverwriteWithEmptyValue); err != nil {
		return err
	}

	_, err = h.Upgrade(a.name, defaultChartRepo+"/"+chart, vals)
	return err
}

func (a *AppInstance) Start(chart string) error {
	h, _ := helm.NamespaceScope(a.namespace)
	vals, err := h.GetValues(a.name)
	if err != nil {
		return err
	}

	globalVals, ok := vals["global"]
	if ok {
		globalVals := globalVals.(map[string]interface{})
		delete(globalVals, "stoped")
		vals["global"] = globalVals
	}

	_, err = h.Upgrade(a.name, defaultChartRepo+"/"+chart, vals)
	return err
}
