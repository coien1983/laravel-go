package microservice

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryServiceRegistry 内存服务注册中心
type MemoryServiceRegistry struct {
	services map[string]*ServiceInfo
	watchers map[string]chan ServiceEvent
	mutex    sync.RWMutex
	closed   bool
}

// NewMemoryServiceRegistry 创建内存服务注册中心
func NewMemoryServiceRegistry() *MemoryServiceRegistry {
	return &MemoryServiceRegistry{
		services: make(map[string]*ServiceInfo),
		watchers: make(map[string]chan ServiceEvent),
	}
}

// Register 注册服务
func (r *MemoryServiceRegistry) Register(ctx context.Context, service *ServiceInfo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.closed {
		return fmt.Errorf("registry is closed")
	}

	// 设置创建时间和更新时间
	now := time.Now()
	service.CreatedAt = now
	service.UpdatedAt = now
	service.LastCheck = now

	// 如果没有设置 TTL，设置默认值
	if service.TTL == 0 {
		service.TTL = 30 * time.Second
	}

	r.services[service.ID] = service

	// 通知监听器
	r.notifyWatchers(ServiceEvent{
		Type:    ServiceEventCreated,
		Service: service,
	})

	return nil
}

// Deregister 注销服务
func (r *MemoryServiceRegistry) Deregister(ctx context.Context, serviceID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.closed {
		return fmt.Errorf("registry is closed")
	}

	service, exists := r.services[serviceID]
	if !exists {
		return fmt.Errorf("service %s not found", serviceID)
	}

	delete(r.services, serviceID)

	// 通知监听器
	r.notifyWatchers(ServiceEvent{
		Type:    ServiceEventDeleted,
		Service: service,
	})

	return nil
}

// Update 更新服务信息
func (r *MemoryServiceRegistry) Update(ctx context.Context, service *ServiceInfo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.closed {
		return fmt.Errorf("registry is closed")
	}

	existing, exists := r.services[service.ID]
	if !exists {
		return fmt.Errorf("service %s not found", service.ID)
	}

	// 保留创建时间
	service.CreatedAt = existing.CreatedAt
	service.UpdatedAt = time.Now()

	r.services[service.ID] = service

	// 通知监听器
	r.notifyWatchers(ServiceEvent{
		Type:    ServiceEventUpdated,
		Service: service,
	})

	return nil
}

// GetService 获取服务信息
func (r *MemoryServiceRegistry) GetService(ctx context.Context, serviceID string) (*ServiceInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if r.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	service, exists := r.services[serviceID]
	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceID)
	}

	return service, nil
}

// ListServices 列出所有服务
func (r *MemoryServiceRegistry) ListServices(ctx context.Context) ([]*ServiceInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if r.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	services := make([]*ServiceInfo, 0, len(r.services))
	for _, service := range r.services {
		services = append(services, service)
	}

	return services, nil
}

// Watch 监听服务变化
func (r *MemoryServiceRegistry) Watch(ctx context.Context) (<-chan ServiceEvent, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	watcherID := fmt.Sprintf("watcher_%d", time.Now().UnixNano())
	eventChan := make(chan ServiceEvent, 100)
	r.watchers[watcherID] = eventChan

	// 启动清理协程
	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.watchers, watcherID)
		close(eventChan)
		r.mutex.Unlock()
	}()

	return eventChan, nil
}

// Close 关闭注册中心
func (r *MemoryServiceRegistry) Close() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.closed {
		return nil
	}

	r.closed = true

	// 关闭所有监听器
	for _, watcher := range r.watchers {
		close(watcher)
	}
	r.watchers = make(map[string]chan ServiceEvent)

	return nil
}

// notifyWatchers 通知所有监听器
func (r *MemoryServiceRegistry) notifyWatchers(event ServiceEvent) {
	for _, watcher := range r.watchers {
		select {
		case watcher <- event:
		default:
			// 如果通道满了，跳过这个事件
		}
	}
}

// CleanupExpiredServices 清理过期服务
func (r *MemoryServiceRegistry) CleanupExpiredServices() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	expiredServices := make([]string, 0)

	for id, service := range r.services {
		if now.Sub(service.LastCheck) > service.TTL {
			expiredServices = append(expiredServices, id)
		}
	}

	for _, id := range expiredServices {
		service := r.services[id]
		delete(r.services, id)

		// 通知监听器
		r.notifyWatchers(ServiceEvent{
			Type:    ServiceEventDeleted,
			Service: service,
		})
	}
}

// StartCleanupWorker 启动清理工作协程
func (r *MemoryServiceRegistry) StartCleanupWorker(interval time.Duration) {
	if interval == 0 {
		interval = 10 * time.Second
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				r.CleanupExpiredServices()
			}
		}
	}()
}
