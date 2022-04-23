// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package router

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	successMessage = "请求成功"
)

func renderError(c *gin.Context, code int, err error) {
	_ = c.Error(err)
	c.JSON(code, gin.H{
		"code":      code,
		"message":   err.Error(),
		"timestamp": time.Now().Unix(),
	})
}

func renderSuccess(c *gin.Context, code int) {
	c.JSON(code, gin.H{
		"code":      200,
		"message":   successMessage,
		"timestamp": time.Now().Unix(),
	})
}

func renderJson(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"code":      200,
		"message":   successMessage,
		"timestamp": time.Now().Unix(),
		"data":      data,
	})
}

func renderJsonWithPagination(c *gin.Context, code int, data interface{}, p interface{}) {
	c.JSON(code, gin.H{
		"code":       200,
		"message":    successMessage,
		"timestamp":  time.Now().Unix(),
		"data":       data,
		"pagination": p,
	})
}

type response2xx struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
}

type response5xx struct {
	Success    bool        `json:"success" default:"false"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
}
