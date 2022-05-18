// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/model"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/service"
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
