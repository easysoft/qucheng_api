package router

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
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

	r.POST("/api/cne/app/install", AppInstall)
}
