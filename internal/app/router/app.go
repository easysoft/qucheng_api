// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package router

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm/form"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/tlog"
	"gopkg.in/yaml.v3"
	"net/http"

	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service/app"

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
		ctx  = c.Request.Context()
		err  error
		body model.AppCreateModel
	)
	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	i, err := service.Apps(ctx, body.Cluster, body.Namespace).GetApp(body.Name)
	if i != nil {
		tlog.WithCtx(ctx).ErrorS(nil, "app already exists, install can't continue",
			"cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name, "chart", body.Chart)
		renderError(c, http.StatusInternalServerError, errors.New("app already installed"))
		return
	}

	if err = service.Apps(ctx, body.Cluster, body.Namespace).Install(body.Name, body); err != nil {
		tlog.WithCtx(ctx).ErrorS(err, "install app failed",
			"cluster", body.Cluster, "namespace", body.Namespace,
			"name", body.Name, "chart", body.Chart)
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	tlog.WithCtx(ctx).InfoS("install app successful",
		"cluster", body.Cluster, "namespace", body.Namespace,
		"name", body.Name, "chart", body.Chart)
	renderSuccess(c, http.StatusCreated)
}

// AppUnInstall 卸载接口
// @Summary 卸载接口
// @Tags 应用管理
// @Description 卸载接口
// @Accept json
// @Produce json
// @Param Authorization header string false "jwtToken"
// @Param X-Auth-Token header string false "staticToken"
// @Security ApiKeyAuth
// @Param body body model.AppModel true "meta"
// @Success 201 {object} response2xx
// @Failure 500 {object} response5xx
// @Router /api/cne/app/uninstall [post]
func AppUnInstall(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		err  error
		body model.AppModel
	)
	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	_, err = service.Apps(ctx, body.Cluster, body.Namespace).GetApp(body.Name)
	if err != nil {
		tlog.WithCtx(ctx).ErrorS(err, errGetAppFailed, "cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name)
		if errors.As(err, &app.ErrAppNotFound{}) {
			renderError(c, http.StatusNotFound, err)
			return
		}
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	if err = service.Apps(ctx, body.Cluster, body.Namespace).UnInstall(body.Name); err != nil {
		tlog.WithCtx(ctx).ErrorS(err, "uninstall app failed",
			"cluster", body.Cluster, "namespace", body.Namespace,
			"app", body.Name)
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	tlog.WithCtx(ctx).InfoS("uninstall app successful",
		"cluster", body.Cluster, "namespace", body.Namespace,
		"app", body.Name)
	renderSuccess(c, http.StatusOK)
}

// AppStart 启动应用
// @Summary 启动应用
// @Tags 应用管理
// @Description 启动应用
// @Accept json
// @Produce json
// @Param Authorization header string false "jwtToken"
// @Param X-Auth-Token header string false "staticToken"
// @Security ApiKeyAuth
// @Param body body model.AppManageModel true "meta"
// @Success 201 {object} response2xx
// @Failure 500 {object} response5xx
// @Router /api/cne/app/start [post]
func AppStart(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		err  error
		body model.AppManageModel
	)

	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}
	a, err := service.Apps(ctx, body.Cluster, body.Namespace).GetApp(body.Name)
	if err != nil {
		tlog.WithCtx(ctx).ErrorS(err, errGetAppFailed, "cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name)
		if errors.As(err, &app.ErrAppNotFound{}) {
			renderError(c, http.StatusNotFound, err)
			return
		}
		renderError(c, http.StatusInternalServerError, errors.New(errStartAppFailed))
		return
	}

	err = a.Start(body.Chart, body.Channel)
	if err != nil {
		tlog.WithCtx(ctx).ErrorS(err, errStartAppFailed, "cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name)
		renderError(c, http.StatusInternalServerError, errors.New(errStartAppFailed))
		return
	}
	tlog.WithCtx(ctx).InfoS("start app successful", "cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name)
	renderSuccess(c, http.StatusOK)
}

// AppStop 关闭应用
// @Summary 关闭应用
// @Tags 应用管理
// @Description 关闭应用
// @Accept json
// @Produce json
// @Param Authorization header string false "jwtToken"
// @Param X-Auth-Token header string false "staticToken"
// @Security ApiKeyAuth
// @Param body body model.AppManageModel true "meta"
// @Success 201 {object} response2xx
// @Failure 500 {object} response5xx
// @Router /api/cne/app/stop [post]
func AppStop(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		err  error
		body model.AppManageModel
	)
	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	a, err := service.Apps(ctx, body.Cluster, body.Namespace).GetApp(body.Name)
	if err != nil {
		tlog.WithCtx(ctx).ErrorS(err, errGetAppFailed, "cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name)
		if errors.As(err, &app.ErrAppNotFound{}) {
			renderError(c, http.StatusNotFound, err)
			return
		}
		renderError(c, http.StatusInternalServerError, errors.New(errStopAppFailed))
		return
	}

	err = a.Stop(body.Chart, body.Channel)
	if err != nil {
		tlog.WithCtx(ctx).ErrorS(err, errStopAppFailed, "cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name)
		renderError(c, http.StatusInternalServerError, errors.New(errStopAppFailed))
		return
	}
	tlog.WithCtx(ctx).InfoS("stop app successful", "cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name)
	renderSuccess(c, http.StatusOK)
}

// AppStop 设置应用
// @Summary 设置应用
// @Tags 应用管理
// @Description 设置应用
// @Accept json
// @Produce json
// @Param Authorization header string false "jwtToken"
// @Param X-Auth-Token header string false "staticToken"
// @Security ApiKeyAuth
// @Param body body model.AppCreateModel true "meta"
// @Success 201 {object} response2xx
// @Failure 500 {object} response5xx
// @Router /api/cne/app/settings [post]
func AppPatchSettings(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		err  error
		body model.AppCreateModel
	)
	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	a, err := service.Apps(ctx, body.Cluster, body.Namespace).GetApp(body.Name)
	if err != nil {
		tlog.WithCtx(ctx).ErrorS(err, errGetAppFailed, "cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name)
		if errors.As(err, &app.ErrAppNotFound{}) {
			renderError(c, http.StatusNotFound, err)
			return
		}
		renderError(c, http.StatusInternalServerError, errors.New(errPatchAppFailed))
		return
	}

	err = a.PatchSettings(body.Chart, body)
	if err != nil {
		tlog.WithCtx(ctx).ErrorS(err, errPatchAppFailed, "cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name)
		renderError(c, http.StatusInternalServerError, errors.New(errPatchAppFailed))
		return
	}
	tlog.WithCtx(ctx).InfoS("patch app settings failed", "cluster", body.Cluster, "namespace", body.Namespace, "name", body.Name)
	renderSuccess(c, http.StatusOK)
}

// AppStatus 应用状态
// @Summary 应用状态
// @Tags 应用管理
// @Description 应用状态
// @Accept json
// @Produce json
// @Param Authorization header string false "jwtToken"
// @Param X-Auth-Token header string false "staticToken"
// @Security ApiKeyAuth
// @Param body query model.AppModel true "meta"
// @Success 201 {object} response2xx
// @Failure 500 {object} response5xx
// @Router /api/cne/app/status [get]
func AppStatus(c *gin.Context) {
	var (
		ctx = c.Request.Context()

		err   error
		query model.AppModel
		i     *app.Instance
		data  *model.AppRespStatus
	)
	if err = c.ShouldBindQuery(&query); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	i, err = service.Apps(ctx, query.Cluster, query.Namespace).GetApp(query.Name)
	if err != nil {
		tlog.WithCtx(ctx).ErrorS(err, errGetAppFailed, "cluster", query.Cluster, "namespace", query.Namespace, "name", query.Name)
		if errors.As(err, &app.ErrAppNotFound{}) {
			renderError(c, http.StatusNotFound, err)
			return
		}
		renderError(c, http.StatusInternalServerError, errors.New(errGetAppStatusFailed))
		return
	}

	data = i.ParseStatus()

	/*
		parse App Uri
	*/
	data.AccessHost = ""
	nodePort := i.ParseNodePort()
	if nodePort > 0 {
		nodePortIPS := service.Nodes(ctx, query.Cluster).ListNodePortIPS()
		if len(nodePortIPS) != 0 {
			accessHost := fmt.Sprintf("%s:%d", nodePortIPS[0], nodePort)
			data.AccessHost = accessHost
		}
	}
	renderJson(c, http.StatusOK, data)
}

func AppSimpleSettings(c *gin.Context) {
	var (
		ctx = c.Request.Context()

		err   error
		query model.AppSettingMode
		i     *app.Instance
	)

	if err = c.ShouldBindQuery(&query); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	i, err = service.Apps(ctx, query.Cluster, query.Namespace).GetApp(query.Name)
	if err != nil {
		tlog.WithCtx(ctx).ErrorS(err, errGetAppFailed, "cluster", query.Cluster, "namespace", query.Namespace, "name", query.Name)
		if errors.As(err, &app.ErrAppNotFound{}) {
			renderError(c, http.StatusNotFound, err)
			return
		}
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	settings, err := i.Settings().Simple().Mode(query.Mode).Parse()
	if err != nil {
		tlog.WithCtx(ctx).ErrorS(err, "get simple settings failed", "cluster", query.Cluster, "namespace", query.Namespace, "name", query.Name)
		renderError(c, http.StatusInternalServerError, err)
		return
	}
	renderJson(c, http.StatusOK, settings)
}

func AppTest(c *gin.Context) {
	ch, err := helm.GetChart("qucheng-test/cne-market-api")
	if err != nil {
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	fmt.Println(helm.ParseValues(ch.Values))
	var dynForm form.DynamicForm
	for _, f := range ch.Files {
		if f.Name == "form.yaml" {
			err = yaml.Unmarshal(f.Data, &dynForm)
			if err != nil {
				klog.ErrorS(err, "parse dynform failed")
				renderError(c, http.StatusInternalServerError, errors.New("parse dynform failed"))
				return
			}
		}
	}
	fmt.Println(dynForm)
	renderSuccess(c, 200)
}
