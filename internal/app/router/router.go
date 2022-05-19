// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "gitlab.zcorp.cc/pangu/cne-api/docs"
)

func Config(r *gin.Engine) {
	r.Use(Cors())
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health", "/metrics"},
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf(`time="%s" client=%s method=%s path=%s proto=%s status=%d cost=%s user-agent="%s" error="%s"`+"\n",
				param.TimeStamp.Format(time.RFC3339),
				param.ClientIP,
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency.String(),
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
	}))
	r.Use(gin.Recovery())
	r.GET("/ping", ping)
	r.GET("/health", health)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/cne", Auth())
	{
		api.POST("/app/install", AppInstall)
		api.POST("/app/uninstall", AppUnInstall)
		api.POST("/app/start", AppStart)
		api.POST("/app/stop", AppStop)
		api.POST("/app/settings", AppPatchSettings)
		api.GET("/app/status", AppStatus)

		api.POST("/namespace/create", NamespaceCreate)
		api.POST("/namespace/recycle", NamespaceRecycle)
		api.GET("/namespace", NamespaceGet)

		api.POST("/middleware/install", MiddlewareInstall)
		api.POST("/middleware/uninstall", MiddleWareUninstall)
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
