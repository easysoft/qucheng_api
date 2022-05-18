package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service"
	"net/http"
)

func NamespaceCreate(c *gin.Context) {
	var (
		err  error
		body model.NamespaceBase
	)

	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	if err = service.Namespaces(body.Cluster).Create(body.Name); err != nil {
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	renderSuccess(c, http.StatusOK)
}

func NamespaceRecycle(c *gin.Context) {
	var (
		err  error
		body model.NamespaceBase
	)

	if err = c.ShouldBindJSON(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	if err = service.Namespaces(body.Cluster).Recycle(body.Name); err != nil {
		renderError(c, http.StatusInternalServerError, err)
		return
	}

	renderSuccess(c, http.StatusOK)
}

func NamespaceGet(c *gin.Context) {
	var (
		err  error
		body model.NamespaceBase
	)

	if err = c.ShouldBindQuery(&body); err != nil {
		renderError(c, http.StatusBadRequest, err)
		return
	}

	if ok := service.Namespaces(body.Cluster).Has(body.Name); !ok {
		renderError(c, http.StatusNotFound, errors.New("namespace not found"))
		return
	}

	renderSuccess(c, http.StatusOK)
}
