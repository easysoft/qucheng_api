// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package helm

import (
	"context"
	"log"
	"os"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/strvals"

	"github.com/imdario/mergo"

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

func (h *HelmAction) Install(name, chart string, settings []string) (*release.Release, error) {
	client := action.NewInstall(h.actionConfig)

	cp, err := client.ChartPathOptions.LocateChart(chart, h.settings)
	if err != nil {
		return nil, err
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	valueOpts := &values.Options{StringValues: settings}

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

func (h *HelmAction) GetValues(name string) (map[string]interface{}, error) {
	client := action.NewGetValues(h.actionConfig)
	vals, err := client.Run(name)
	return vals, err
}

func (h *HelmAction) Upgrade(name string, chart string, chartValues map[string]interface{}) (interface{}, error) {
	client := action.NewUpgrade(h.actionConfig)
	valueOpts := &values.Options{}

	client.Namespace = h.namespace

	cp, err := client.ChartPathOptions.LocateChart(chart, h.settings)
	if err != nil {
		return nil, err
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	p := getter.All(h.settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	if err := mergo.Merge(&vals, chartValues, mergo.WithOverwriteWithEmptyValue); err != nil {
		return nil, err
	}

	ctx := context.Background()
	rel, err := client.RunWithContext(ctx, name, chartRequested, vals)
	return rel, err
}

func (h *HelmAction) PatchValues(dest map[string]interface{}, setvals []string) error {
	for _, value := range setvals {
		if err := strvals.ParseInto(value, dest); err != nil {
			return errors.Wrap(err, "failed parsing --set data")
		}
	}

	return nil
}
