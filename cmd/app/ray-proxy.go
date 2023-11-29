package app

import (
	"context"
	"fmt"

	"github.com/daocloud/ray-proxy/cmd/app/options"
	"github.com/daocloud/ray-proxy/pkg/proxy"
	"github.com/daocloud/ray-proxy/pkg/version"
	rayv1 "github.com/ray-project/kuberay/ray-operator/apis/ray/v1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

var schema = runtime.NewScheme()

func init() {
	utilruntime.Must(rayv1.AddToScheme(schema))
}

func NewRayProxyCommand(ctx context.Context) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:   "ray-proxy",
		Short: "A proxy for Ray",
		Long:  `A proxy for Ray`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(ctx, opts)
		},
	}

	fs := cmd.Flags()
	namedFlagSets := opts.Flag()
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	fmt.Print(opts.ProxyPort)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of ray-proxy",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(version.Get())
		},
	}

	cmd.AddCommand(versionCmd)

	return cmd
}

func Run(ctx context.Context, opts *options.Options) error {
	klog.Info("Starting Ray Proxy...")

	config := controllerruntime.GetConfigOrDie()
	controllerOptions := controllerruntime.Options{
		Scheme: schema,
	}

	controllerManager, err := controllerruntime.NewManager(config, controllerOptions)

	if err != nil {
		klog.Fatalf("Error creating controller manager: %v", err)
		return err
	}

	proxyController := proxy.Controller{
		Client:           controllerManager.GetClient(),
		EventRecorder:    controllerManager.GetEventRecorderFor(opts.LeaderElection.ResourceName),
		RayServicesCache: proxy.NewCache(),
		ProxyPort:        opts.ProxyPort,
		ClusterTemplate:  opts.ClusterTemplate,
	}

	if err := proxyController.SetupWithManager(controllerManager); err != nil {
		klog.ErrorS(err, "Error setting up controller")
		return err
	}

	// blocks until the context is done.
	if err := controllerManager.Start(ctx); err != nil {
		klog.ErrorS(err, "ray service controller manager exits unexpectedly")
		return err
	}

	return nil
}
