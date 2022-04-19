package router

import (
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
)

func AppInstall(c *gin.Context) {
	var (
		err  error
		body model.AppCreateModel
	)
	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	if err = service.Apps(body.Cluster, body.Namespace).Install(body.Name, body); err != nil {
		klog.ErrorS(err, "install app failed",
			"cluster", body.Cluster, "namespace", body.Namespace,
			"app", body.Name)
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	klog.InfoS("install app successful",
		"cluster", body.Cluster, "namespace", body.Namespace,
		"app", body.Name)
	renderSuccess(c, http.StatusCreated)
}
