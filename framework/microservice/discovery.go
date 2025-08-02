package microservice

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryServiceDiscovery 内存服务发现
type MemoryServiceDiscovery struct {
	registry     ServiceRegistry
	loadBalancer LoadBalancer
	cache        map[string][]*ServiceInfo
	cacheMutex   sync.RWMutex
	watchers     map[string]chan ServiceEvent
	watcherMutex sync.RWMutex
	closed       bool
}

// NewMemoryServiceDiscovery 创建内存服务发现
func NewMemoryServiceDiscovery(registry ServiceRegistry, loadBalancer LoadBalancer) *MemoryServiceDiscovery {
	if loadBalancer == nil {
		loadBalancer = NewRoundRobinLoadBalancer()
	}

	return &MemoryServiceDiscovery{
		registry:     registry,
		loadBalancer: loadBalancer,
		cache:        make(map[string][]*ServiceInfo),
		watchers:     make(map[string]chan ServiceEvent),
	}
}

// Discover 发现服务
func (d *MemoryServiceDiscovery) Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
	if d.closed {
		return nil, fmt.Errorf("discovery is closed")
	}

	// 先从缓存获取
	d.cacheMutex.RLock()
	if services, exists := d.cache[serviceName]; exists {
		d.cacheMutex.RUnlock()
		return services, nil
	}
	d.cacheMutex.RUnlock()

	// 从注册中心获取所有服务
	allServices, err := d.registry.ListServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	// 过滤指定名称的服务
	services := make([]*ServiceInfo, 0)
	for _, service := range allServices {
		if service.Name == serviceName {
			services = append(services, service)
		}
	}

	// 更新缓存
	d.cacheMutex.Lock()
	d.cache[serviceName] = services
	d.cacheMutex.Unlock()

	return services, nil
}

// DiscoverOne 发现单个服务（负载均衡）
func (d *MemoryServiceDiscovery) DiscoverOne(ctx context.Context, serviceName string) (*ServiceInfo, error) {
	services, err := d.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no service found with name: %s", serviceName)
	}

	selected := d.loadBalancer.Select(services)
	if selected == nil {
		return nil, fmt.Errorf("no healthy service available for: %s", serviceName)
	}

	return selected, nil
}

// Watch 监听服务变化
func (d *MemoryServiceDiscovery) Watch(ctx context.Context, serviceName string) (<-chan ServiceEvent, error) {
	if d.closed {
		return nil, fmt.Errorf("discovery is closed")
	}

	// 创建监听器
	watcherID := fmt.Sprintf("%s_%d", serviceName, time.Now().UnixNano())
	eventChan := make(chan ServiceEvent, 100)

	d.watcherMutex.Lock()
	d.watchers[watcherID] = eventChan
	d.watcherMutex.Unlock()

	// 启动监听协程
	go func() {
		defer func() {
			d.watcherMutex.Lock()
			delete(d.watchers, watcherID)
			close(eventChan)
			d.watcherMutex.Unlock()
		}()

		// 监听注册中心的变化
		registryEvents, err := d.registry.Watch(ctx)
		if err != nil {
			return
		}

		for {
			select {
			case event := <-registryEvents:
				if event.Service.Name == serviceName {
					// 更新缓存
					d.updateCache(event)

					// 转发事件
					select {
					case eventChan <- event:
					default:
						// 如果通道满了，跳过这个事件
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventChan, nil
}

// Close 关闭发现服务
func (d *MemoryServiceDiscovery) Close() error {
	d.watcherMutex.Lock()
	defer d.watcherMutex.Unlock()

	if d.closed {
		return nil
	}

	d.closed = true

	// 关闭所有监听器
	for _, watcher := range d.watchers {
		close(watcher)
	}
	d.watchers = make(map[string]chan ServiceEvent)

	return nil
}

// updateCache 更新缓存
func (d *MemoryServiceDiscovery) updateCache(event ServiceEvent) {
	d.cacheMutex.Lock()
	defer d.cacheMutex.Unlock()

	serviceName := event.Service.Name

	switch event.Type {
	case ServiceEventCreated:
		// 添加到缓存
		if services, exists := d.cache[serviceName]; exists {
			d.cache[serviceName] = append(services, event.Service)
		} else {
			d.cache[serviceName] = []*ServiceInfo{event.Service}
		}

	case ServiceEventUpdated:
		// 更新缓存中的服务
		if services, exists := d.cache[serviceName]; exists {
			for i, service := range services {
				if service.ID == event.Service.ID {
					services[i] = event.Service
					break
				}
			}
		}

	case ServiceEventDeleted:
		// 从缓存中删除服务
		if services, exists := d.cache[serviceName]; exists {
			newServices := make([]*ServiceInfo, 0)
			for _, service := range services {
				if service.ID != event.Service.ID {
					newServices = append(newServices, service)
				}
			}
			d.cache[serviceName] = newServices
		}
	}
}

// SetLoadBalancer 设置负载均衡器
func (d *MemoryServiceDiscovery) SetLoadBalancer(loadBalancer LoadBalancer) {
	d.loadBalancer = loadBalancer
}

// ClearCache 清除缓存
func (d *MemoryServiceDiscovery) ClearCache() {
	d.cacheMutex.Lock()
	defer d.cacheMutex.Unlock()

	d.cache = make(map[string][]*ServiceInfo)
}

// GetCacheStats 获取缓存统计信息
func (d *MemoryServiceDiscovery) GetCacheStats() map[string]int {
	d.cacheMutex.RLock()
	defer d.cacheMutex.RUnlock()

	stats := make(map[string]int)
	for serviceName, services := range d.cache {
		stats[serviceName] = len(services)
	}

	return stats
}
