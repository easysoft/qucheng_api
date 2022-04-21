// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package helm

import (
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
)

type HelmAction struct {
	actionConfig *action.Configuration
	settings     *cli.EnvSettings
	namespace    string
}

func newAction(settings *cli.EnvSettings, config *action.Configuration) *HelmAction {
	return &HelmAction{
		settings: settings, actionConfig: config,
	}
}

func NamespaceScope(namespace string) (*HelmAction, error) {
	settings := cli.New()
	settings.SetNamespace(namespace)
	actionConfig := &action.Configuration{}
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}
	h := newAction(settings, actionConfig)
	h.namespace = namespace

	return h, nil
}

func (h *HelmAction) Install(name, chart string) (*release.Release, error) {
	client := action.NewInstall(h.actionConfig)

	cp, err := client.ChartPathOptions.LocateChart(chart, h.settings)
	if err != nil {
		return nil, err
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	valueOpts := &values.Options{}

	p := getter.All(h.settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	client.Namespace = h.namespace
	client.ReleaseName = name

	rel, err := client.Run(chartRequested, vals)
	if err != nil {
		return nil, err
	}
	return rel, nil
}

func (h *HelmAction) Uninstall(name string) error {
	client := action.NewUninstall(h.actionConfig)
	_, err := client.Run(name)
	return err
}

//func Push() error {
//	cli, err := push.New()
//	if err != nil {
//		return err
//	}
//	res, err := cli.UploadChartPackage("/tmp/redis-10.5.8.tgz", false)
//
//	fmt.Println(res.StatusCode, res.Body)
//	return err
//}
