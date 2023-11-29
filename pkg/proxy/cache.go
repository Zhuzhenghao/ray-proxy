package proxy

import (
	"sync"

	rayv1 "github.com/ray-project/kuberay/ray-operator/apis/ray/v1"
)

type Cache struct {
	sync.RWMutex
	rayServices map[string]*rayv1.RayService
}

func NewCache() *Cache {
	return &Cache{
		RWMutex:     sync.RWMutex{},
		rayServices: map[string]*rayv1.RayService{},
	}
}

func (c *Cache) Add(rayService *rayv1.RayService) {
	c.Lock()
	defer c.Unlock()

	c.rayServices[rayService.Name] = rayService
}

func (c *Cache) Delete(name string) {
	c.Lock()
	defer c.Unlock()

	delete(c.rayServices, name)
}

func (c *Cache) Get(name string) (*rayv1.RayService, bool) {
	c.RLock()
	defer c.RUnlock()

	cluster, ok := c.rayServices[name]
	return cluster, ok
}

func (c *Cache) Exists(name string) bool {
	c.RLock()
	defer c.RUnlock()

	_, ok := c.rayServices[name]
	return ok
}

func (c *Cache) ContainsRayService(targetRayService *rayv1.RayService, eqFunc func(a, b *rayv1.RayService) bool) bool {
	c.RLock()
	defer c.RUnlock()

	for key := range c.rayServices {
		if eqFunc(targetRayService, c.rayServices[key]) {
			return true
		}
	}
	return false
}
