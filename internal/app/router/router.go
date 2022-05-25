// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package router

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model/translation"
	"net/http"
	"time"

	_ "gitlab.zcorp.cc/pangu/cne-api/docs"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model/field"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"k8s.io/klog/v2"
)

var trans ut.Translator

func initTrans(v *validator.Validate) (err error) {
	enT := en.New()
	uni := ut.New(enT, enT)

	var (
		ok    bool
		local = "en"
	)

	trans, ok = uni.GetTranslator(local)
	if !ok {
		return fmt.Errorf("uni.GetTranslator(%s) failed", local)
	}

	if err = enTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return
	}
	if err = translation.RegisterCustomENTranslations(v, trans); err != nil {
		return
	}
	return
}

func Config(r *gin.Engine) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var err error
		for k, f := range field.New() {
			if err = v.RegisterValidation(k, f); err != nil {
				klog.ErrorS(err, "register validation failed", "name", k)
			} else {
				klog.InfoS("register validation successful", "name", k)
			}
		}

		if err = initTrans(v); err != nil {
			klog.ErrorS(err, "setup translate failed")
		}
	}

	r.Use(Cors())
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health", "/metrics"},
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf(`time="%s" client=%s method=%s path=%s proto=%s status=%d cost=%s user-agent="%s" error="%s" traceId=%s`+"\n",
				param.TimeStamp.Format(time.RFC3339),
				param.ClientIP,
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency.String(),
				param.Request.UserAgent(),
				param.ErrorMessage,
				param.Request.Header.Get(HeaderTraceId),
			)
		},
	}))
	r.Use(gin.Recovery())
	r.GET("/ping", ping)
	r.GET("/health", health)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/cne", Auth(), Trace())
	{
		api.POST("/app/install", AppInstall)
		api.POST("/app/uninstall", AppUnInstall)
		api.POST("/app/start", AppStart)
		api.POST("/app/stop", AppStop)
		api.POST("/app/settings", AppPatchSettings)
		api.GET("/app/settings/simple", AppSimpleSettings)
		api.GET("/app/status", AppStatus)
		api.GET("/app/metric", AppMetric)
		api.GET("/test", AppTest)

		api.POST("/namespace/create", NamespaceCreate)
		api.POST("/namespace/recycle", NamespaceRecycle)
		api.GET("/namespace", NamespaceGet)

		api.POST("/middleware/install", MiddlewareInstall)
		api.POST("/middleware/uninstall", MiddleWareUninstall)

		api.POST("/statistics/app", AppListStatistics)
		api.GET("/statistics/cluster", ClusterStatistics)
	}

	r.NoMethod(func(c *gin.Context) {
		msg := fmt.Sprintf("not found: %v", c.Request.Method)
		renderMessage(c, http.StatusBadRequest, msg)
	})
	r.NoRoute(func(c *gin.Context) {
		msg := fmt.Sprintf("not found: %v", c.Request.URL.Path)
		renderMessage(c, http.StatusBadRequest, msg)
	})
}
