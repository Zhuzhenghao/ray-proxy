package main

import (
	"os"

	"github.com/daocloud/ray-proxy/cmd/app"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	ctx := server.SetupSignalContext()

	if err := app.NewRayProxyCommand(ctx).Execute(); err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}

}
