// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package serve

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/router"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/cron"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm"
	"k8s.io/klog/v2"
)

const (
	listenPort = 8087
)

func Serve(ctx context.Context) error {
	stopCh := make(chan struct{})

	klog.Info("Initialize clusters")
	err := cluster.Init(stopCh)
	if err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}

	klog.Info("Setup cron tasks")
	defer cron.Cron.Stop()
	cron.Cron.Start()
	cron.Cron.Add("0 */2 * * *", func() {
		err = helm.RepoUpdate()
		if err != nil {
			klog.Warningf("cron helm repo update err: %v", err)
		}
	})

	klog.Info("Starting cne-api...")

	klog.Info("Setup gin engine")
	r := gin.New()
	router.Config(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: r,
	}
	go func() {
		defer close(stopCh)
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			klog.Errorf("Failed to stop server, error: %s", err)
		}
		klog.Info("server exited.")
	}()
	klog.Infof("start application server, Listen on port: %d, pid is %v", listenPort, os.Getpid())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		klog.Errorf("Failed to start http server, error: %s", err)
		return err
	}
	<-stopCh
	return nil
}
