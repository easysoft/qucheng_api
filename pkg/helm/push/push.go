package push

import (
	cm "github.com/chartmuseum/helm-push/pkg/chartmuseum"
	"github.com/chartmuseum/helm-push/pkg/helm"
	"strings"
)

func New() (*cm.Client, error) {
	var repo *helm.Repo
	var err error
	repo, err = helm.GetRepoByName("qucheng-test")
	if err != nil {
		return nil, err
	}

	var url = strings.Replace(repo.Config.URL, "cm://", "https://", 1)
	return cm.NewClient(
		cm.URL(url),
		cm.Username(repo.Config.Username),
		cm.Password(repo.Config.Password),
	)
}
