// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package router

import (
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/constant"
	"net/http"

	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
)

// AppInstall 安装接口
// @Summary 安装接口
// @Tags 应用管理
// @Description 安装接口
// @Accept json
// @Produce json
// @Param Authorization header string false "jwtToken"
// @Param X-Auth-Token header string false "staticToken"
// @Security ApiKeyAuth
// @Param body body model.AppCreateModel true "meta"
// @Success 201 {object} response2xx
// @Failure 500 {object} response5xx
// @Router /api/cne/app/install [post]
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

func AppUnInstall(c *gin.Context) {
	var (
		err  error
		body model.AppModel
	)
	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	if err = service.Apps(body.Cluster, body.Namespace).UnInstall(body.Name); err != nil {
		klog.ErrorS(err, "uninstall app failed",
			"cluster", body.Cluster, "namespace", body.Namespace,
			"app", body.Name)
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	klog.InfoS("uninstall app successful",
		"cluster", body.Cluster, "namespace", body.Namespace,
		"app", body.Name)
	renderSuccess(c, http.StatusOK)
}

func AppStart(c *gin.Context) {
	var (
		err  error
		body model.AppManageModel
	)

	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}
	app, err := service.Apps(body.Cluster, body.Namespace).GetApp(body.Name)
	if err != nil {
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	err = app.Start(body.Chart)
	if err != nil {
		renderError(c, http.StatusInternalServerError, err)
		return
	}
	renderSuccess(c, http.StatusOK)
}

func AppStop(c *gin.Context) {
	var (
		err  error
		body model.AppManageModel
	)
	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	app, err := service.Apps(body.Cluster, body.Namespace).GetApp(body.Name)
	if err != nil {
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	err = app.Stop(body.Chart)
	if err != nil {
		renderError(c, http.StatusInternalServerError, err)
		return
	}
	renderSuccess(c, http.StatusOK)
}

func AppStatus(c *gin.Context) {
	var (
		err error
		query model.AppModel
		//app *app.AppInstance
		data model.AppRespStatus
	)
	if err = c.ShouldBindQuery(&query); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	//if app, err = service.Apps(query.Cluster, query.Namespace).GetAppWithOutHelm(query.Name);err != nil {
	//	renderError(c, http.StatusInternalServerError, err)
	//	return
	//}
	_, err = service.Apps(query.Cluster, query.Namespace).GetAppWithOutHelm(query.Name)
	if err != nil {
		data.Status = constant.AppStatusAbnormal.String()
	} else {
		data.Status = constant.AppStatusRunning.String()
	}

	renderJson(c, http.StatusOK, data)
}
