// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"

	_ "go.uber.org/automaxprocs"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	"gitlab.zcorp.cc/pangu/cne-api/internal/app/router"
)

const (
	listenPort = 8087
)

func initKubeLogs() {
	gofs := flag.NewFlagSet("klog", flag.ExitOnError)
	_ = gofs.Set("add_dir_header", "true")
	klog.InitFlags(gofs)
}

// @title CNE API
// @version 1.0.0
// @description CNE API.
// @contact.name QuCheng Pangu Team
// @license.name Z PUBLIC LICENSE 1.2
func main() {
	initKubeLogs()

	klog.Info("Starting cne-api...")

	klog.Info("Setup gin engine")
	r := gin.New()
	router.Config(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: r,
	}
	klog.Infof("start application server, Listen on port: %d", listenPort)
	_ = srv.ListenAndServe()
}
