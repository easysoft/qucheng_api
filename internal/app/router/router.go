// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package router

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "gitlab.zcorp.cc/pangu/cne-api/docs"
)

func Config(r *gin.Engine) {
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
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
	}))
	r.GET("/ping", ping)
	r.GET("/health", health)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/cne/app/install", AppInstall)
	r.POST("/api/cne/app/uninstall", AppUnInstall)
	r.POST("/api/cne/app/start", AppStart)
	r.POST("/api/cne/app/stop", AppStop)
	r.GET("/api/cne/app/status", AppStatus)
}
