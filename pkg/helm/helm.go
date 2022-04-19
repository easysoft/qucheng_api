package helm

import (
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"

	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm/push"
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

	client.Namespace = h.namespace
	client.ReleaseName = name

	rel, err := client.Run(chartRequested, vals)
	if err != nil {
		return nil, err
	}
	return rel, nil
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
