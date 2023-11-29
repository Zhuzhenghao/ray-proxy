package proxy

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/daocloud/ray-proxy/pkg/nginx"
	"github.com/daocloud/ray-proxy/pkg/utils"
	rayv1 "github.com/ray-project/kuberay/ray-operator/apis/ray/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	controllerruntime "sigs.k8s.io/controller-runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Controller struct {
	client.Client
	EventRecorder    record.EventRecorder
	RayServicesCache *Cache

	ProxyPort       int
	ClusterTemplate string
}

func (c *Controller) SetupWithManager(mgr ctrl.Manager) error {
	return utilerrors.NewAggregate([]error{
		ctrl.NewControllerManagedBy(mgr).For(&rayv1.RayService{}).Complete(c),
		mgr.Add(c),
	})
}

func (c *Controller) Start(ctx context.Context) error {
	klog.Infof("Starting ray proxy controller")
	defer klog.Infof("Shutting ray proxy controller")

	return nginx.Reload(ctx)
}

func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	klog.InfoS("Reconciling ray service", "ray service", klog.KRef(req.Namespace, req.Name))

	rayService := &rayv1.RayService{}
	if err := c.Get(ctx, req.NamespacedName, rayService); err != nil {
		if apierrors.IsNotFound(err) {
			klog.InfoS("RayService not found, maybe deleted", "RayService", klog.KRef(req.Namespace, req.Name))
			return ctrl.Result{}, nil
		}
		klog.ErrorS(err, "Failed to get RayService", "RayService", klog.KRef(req.Namespace, req.Name))
		return ctrl.Result{}, err
	}

	if !rayService.DeletionTimestamp.IsZero() {
		result, err := c.HandleDeleteRayService(ctx, rayService.Name)
		if err != nil {
			klog.ErrorS(err, "reconcile delete cluster", "cluster", req.NamespacedName.Name)
		}
		return result, err
	}

	if !utils.IsRayServiceReady(&rayService.Status) {
		klog.InfoS("ray service reconcile cluster not ready and retry again later", "cluster", rayService.Name)
		return controllerruntime.Result{RequeueAfter: 3 * time.Second}, nil
	}

	return c.HandleUpdateRayService(ctx, rayService)
}

func (c *Controller) HandleDeleteRayService(ctx context.Context, name string) (reconcile.Result, error) {
	rayService, ok := c.RayServicesCache.Get(name)
	if !ok {
		return reconcile.Result{}, nil
	}

	err := os.Remove(fmt.Sprintf("/etc/nginx/conf.d/cluster-%s", rayService.Name))
	if err != nil {
		if os.IsNotExist(err) {
			c.RayServicesCache.Delete(name)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	err = nginx.Reload(ctx)
	if err != nil {
		return reconcile.Result{}, err
	}

	c.RayServicesCache.Delete(name)
	return reconcile.Result{}, nil
}

func (c *Controller) HandleUpdateRayService(ctx context.Context, rayService *rayv1.RayService) (reconcile.Result, error) {
	if c.RayServicesCache.ContainsRayService(rayService, utils.EqualRayService) {
		return reconcile.Result{}, nil
	}

	c.RayServicesCache.Add(rayService)

	cacheRayService, ok := c.RayServicesCache.Get(rayService.Name)
	if !ok {
		return reconcile.Result{Requeue: true}, nil
	}

	if err := c.RefreshConf(cacheRayService); err != nil {
		klog.ErrorS(err, "update ray service refresh conf", "cluster", rayService.Name)
		return reconcile.Result{}, err
	}

	if err := nginx.Reload(ctx); err != nil {
		klog.ErrorS(err, "update cluster nginx reload", "cluster", rayService.Name)
		return reconcile.Result{}, err
	}

	klog.InfoS("add cluster to cache and refresh config succeeded", rayService.Name)
	return reconcile.Result{}, nil
}

// RefreshConf refreshes the nginx config file
func (c *Controller) RefreshConf(rayService *rayv1.RayService) error {
	if rayService == nil && rayService.Spec.ServeService == nil {
		return nil
	}

	var servePort string
	for _, port := range rayService.Spec.ServeService.Spec.Ports {
		fmt.Print(port.Name)
		fmt.Print(port.Port)
		if port.Name == "serve" {
			servePort = fmt.Sprintf("%d", port.Port)
			break
		}
	}

	if servePort == "" {
		return fmt.Errorf("no serve port found")
	}

	content, err := utils.RenderTemplate("ray-service-template", c.ClusterTemplate, struct {
		Name string
		Port string
		Host string
	}{
		Name: rayService.Name,
		Host: rayService.Spec.ServeService.Name,
		Port: servePort,
	})
	if err != nil {
		return err
	}

	return utils.WriteFile(content, fmt.Sprintf("/etc/nginx/conf.d/cluster-%s", rayService.Name))
}
