// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package app

import (
	"github.com/imdario/mergo"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm"
)

func (a *Instance) Stop(chart string) error {
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

func (a *Instance) Start(chart string) error {
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

func (a *Instance) PatchSettings(chart string, body model.AppCreateModel) error {
	var (
		err  error
		vals map[string]interface{}
	)

	h, _ := helm.NamespaceScope(a.namespace)
	vals, err = h.GetValues(a.name)
	if err != nil {
		return err
	}

	if vals == nil {
		vals = make(map[string]interface{})
	}

	var settings = make([]string, len(body.Settings))
	for _, s := range body.Settings {
		settings = append(settings, s.Key+"="+s.Val)
	}

	if err = h.PatchValues(vals, settings); err != nil {
		return err
	}

	_, err = h.Upgrade(a.name, defaultChartRepo+"/"+chart, vals)
	return err
}
