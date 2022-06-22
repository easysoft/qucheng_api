// Copyright (c) 2022 北京渠成软件有限公司 All rights reserved.
// Use of this source code is governed by Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// license that can be found in the LICENSE file.

package serve

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	gins "gitlab.zcorp.cc/pangu/cne-api/internal/app/serve"
)

func NewCmdServe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "serve apiserver",
		Run:   serve,
	}
	return cmd
}

func serve(cmd *cobra.Command, args []string) {
	initKubeLogs()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ctx.Done()
		stop()
	}()

	if err := gins.Serve(ctx); err != nil {
		klog.Fatal("run serve: %v", err)
	}
}

func initKubeLogs() {
	gofs := flag.NewFlagSet("klog", flag.ExitOnError)
	_ = gofs.Set("add_dir_header", "true")
	klog.InitFlags(gofs)
}
