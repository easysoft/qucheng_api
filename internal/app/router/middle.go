package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service"
	"net/http"
)

func MiddlewareInstall(c *gin.Context) {
	var (
		err  error
		body model.Middleware
		res  interface{}
	)

	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	if res, err = service.Middlewares().Mysql().CreateDB(&body); err != nil {
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	renderJson(c, http.StatusOK, res)
}

func MiddleWareUninstall(c *gin.Context) {

	var (
		err  error
		body model.Middleware
	)

	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	if err = service.Middlewares().Mysql().RecycleDB(&body); err != nil {
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	renderSuccess(c, http.StatusOK)
}
