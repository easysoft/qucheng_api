package service

import (
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service/app"
)

func Apps(clusterName, namespace string) *app.AppListManager {
	return app.NewApps(clusterName, namespace)
}
