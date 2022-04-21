// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package router

import "github.com/gin-gonic/gin"

func ping(c *gin.Context) {
	c.String(200, "pong")
}

func health(c *gin.Context) {
	c.String(200, "OK")
}
