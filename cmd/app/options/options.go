package options

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/daocloud/ray-proxy/pkg/utils"
	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	cliflag "k8s.io/component-base/cli/flag"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/klog/v2"
)

type Options struct {
	// LeaderElection defines the configuration of leader election client.
	LeaderElection componentbaseconfig.LeaderElectionConfiguration
	// SecurePort is the port that the server serves at.
	// Note: We hope support https in the future once controller-runtime provides the functionality.
	SecurePort int

	ProxyPort       int
	ClusterTemplate string
}

func NewOptions() *Options {
	opts := &Options{
		LeaderElection: componentbaseconfig.LeaderElectionConfiguration{
			ResourceLock:      resourcelock.LeasesResourceLock,
			ResourceNamespace: utils.GetCurrentNSOrDefault(),
			ResourceName:      "ray-proxy-manager",
		},
		SecurePort: 8443,
		ProxyPort:  8000,
	}

	proxyPort, err := strconv.Atoi(utils.NginxPortForProxyKpandaAPIServer)
	if err == nil {
		opts.ProxyPort = proxyPort
	}

	clusterTpl, err := InitTemplate()
	if err == nil {
		opts.ClusterTemplate = clusterTpl
	}

	return opts
}

func InitTemplate() (string, error) {
	clusterTpl, err := os.ReadFile("/etc/template/ray-service.tmpl")
	if err != nil {
		return "", err
	}
	return string(clusterTpl), err
}

func (o *Options) AddFlag(flags *pflag.FlagSet) {
	flags.IntVar(&o.ProxyPort, "proxy-port", o.ProxyPort, "kpanda ingress nginx proxy port")
}

func (o *Options) Flag() *cliflag.NamedFlagSets {
	fss := &cliflag.NamedFlagSets{}
	fs := fss.FlagSet("generic")
	o.AddFlag(fs)

	fs = fss.FlagSet("klog")
	local := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(local)
	local.VisitAll(func(fl *flag.Flag) {
		fl.Name = strings.Replace(fl.Name, "_", "-", -1)
		fs.AddGoFlag(fl)
	})
	return fss
}
