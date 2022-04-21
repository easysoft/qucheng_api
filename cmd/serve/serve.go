package serve

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/robfig/cron"
	"gitlab.zcorp.cc/pangu/cne-api/internal/app/router"
	"gitlab.zcorp.cc/pangu/cne-api/internal/pkg/kube/cluster"
	"gitlab.zcorp.cc/pangu/cne-api/pkg/helm"
	"k8s.io/klog/v2"
)

const (
	listenPort = 8087
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

	stopCh := make(chan struct{})

	klog.Info("Initialize clusters")
	err := cluster.Init(stopCh)
	if err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}

	klog.Info("Setup cron tasks")
	cron := cron.New()
	err = cron.AddFunc("0 */5 * * *", func() {
		err = helm.RepoUpdate()
		if err != nil {
			fmt.Println(err)
		}
	})
	cron.Start()

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

func initKubeLogs() {
	gofs := flag.NewFlagSet("klog", flag.ExitOnError)
	_ = gofs.Set("add_dir_header", "true")
	klog.InitFlags(gofs)
}
