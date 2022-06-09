// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package helm

import (
	"context"
	"log"
	"os"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/strvals"

	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm/form"
)

type Action struct {
	actionConfig *action.Configuration
	settings     *cli.EnvSettings
	namespace    string
}

func newAction(settings *cli.EnvSettings, config *action.Configuration) *Action {
	return &Action{
		settings: settings, actionConfig: config,
	}
}

func NamespaceScope(namespace string) (*Action, error) {
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

func (h *Action) Install(name, chart, version string, settings []string) (*release.Release, error) {
	client := action.NewInstall(h.actionConfig)

	client.ChartPathOptions.Version = version
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

func (h *Action) Uninstall(name string) error {
	client := action.NewUninstall(h.actionConfig)
	_, err := client.Run(name)
	return err
}

func (h *Action) GetValues(name string) (map[string]interface{}, error) {
	client := action.NewGetValues(h.actionConfig)
	vars, err := client.Run(name)
	return vars, err
}

func (h *Action) GetRelease(name string) (*release.Release, error) {
	client := action.NewGet(h.actionConfig)
	rel, err := client.Run(name)
	return rel, err
}

func (h *Action) Upgrade(name, chart, version string, chartValues map[string]interface{}) (interface{}, error) {
	client := action.NewUpgrade(h.actionConfig)
	valueOpts := &values.Options{}

	client.Namespace = h.namespace
	client.ChartPathOptions.Version = version

	cp, err := client.ChartPathOptions.LocateChart(chart, h.settings)
	if err != nil {
		return nil, err
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	p := getter.All(h.settings)
	vars, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	if err := mergo.Merge(&vars, chartValues, mergo.WithOverwriteWithEmptyValue); err != nil {
		return nil, err
	}

	ctx := context.Background()
	rel, err := client.RunWithContext(ctx, name, chartRequested, vars)
	return rel, err
}

func (h *Action) PatchValues(dest map[string]interface{}, setvals []string) error {
	for _, value := range setvals {
		if err := strvals.ParseInto(value, dest); err != nil {
			return errors.Wrap(err, "failed parsing --set data")
		}
	}

	return nil
}

func GetChart(name string) (*chart.Chart, error) {
	var (
		err error
	)

	settings := cli.New()
	actionConfig := &action.Configuration{}

	h := newAction(settings, actionConfig)
	client := action.NewShowWithConfig(action.ShowAll, h.actionConfig)

	cp, err := client.ChartPathOptions.LocateChart(name, settings)
	if err != nil {
		return nil, err
	}

	return GetChartByFile(cp)
}

func GetChartByFile(fp string) (*chart.Chart, error) {
	c, err := loader.Load(fp)
	return c, err
}

func ParseForm(files []*chart.File) (*form.DynamicForm, error) {
	var (
		err     error
		dynForm form.DynamicForm
	)

	for _, f := range files {
		if f.Name == "form.yaml" {
			err = yaml.Unmarshal(f.Data, &dynForm)
			break
		}
	}

	if &dynForm == nil {
		err = errors.New("no dynamic form found")
	}
	return &dynForm, err
}
