package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdServiceRegistry etcd 服务注册中心
type EtcdServiceRegistry struct {
	client     *clientv3.Client
	prefix     string
	leaseID    clientv3.LeaseID
	watchers   map[string]chan ServiceEvent
	watcherMutex sync.RWMutex
	closed     bool
}

// EtcdConfig etcd 配置
type EtcdConfig struct {
	Endpoints []string      `json:"endpoints"`
	Username  string        `json:"username"`
	Password  string        `json:"password"`
	Prefix    string        `json:"prefix"`
	TTL       time.Duration `json:"ttl"`
}

// NewEtcdServiceRegistry 创建 etcd 服务注册中心
func NewEtcdServiceRegistry(config *EtcdConfig) (*EtcdServiceRegistry, error) {
	if config == nil {
		config = &EtcdConfig{
			Endpoints: []string{"localhost:2379"},
			Prefix:    "/laravel-go/services",
			TTL:       30 * time.Second,
		}
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    config.Username,
		Password:    config.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	registry := &EtcdServiceRegistry{
		client:   client,
		prefix:   config.Prefix,
		watchers: make(map[string]chan ServiceEvent),
	}

	// 创建租约
	lease, err := client.Grant(context.Background(), int64(config.TTL.Seconds()))
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create lease: %w", err)
	}
	registry.leaseID = lease.ID

	return registry, nil
}

// Register 注册服务
func (e *EtcdServiceRegistry) Register(ctx context.Context, service *ServiceInfo) error {
	if e.closed {
		return fmt.Errorf("registry is closed")
	}

	// 生成服务路径
	servicePath := e.getServicePath(service.Name, service.ID)
	
	// 序列化服务信息
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %w", err)
	}

	// 注册服务到 etcd
	_, err = e.client.Put(ctx, servicePath, string(data), clientv3.WithLease(e.leaseID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	// 保持租约活跃
	go e.keepAlive()

	// 通知监听器
	e.notifyWatchers(ServiceEvent{
		Type:    ServiceEventCreated,
		Service: service,
	})

	return nil
}

// Deregister 注销服务
func (e *EtcdServiceRegistry) Deregister(ctx context.Context, serviceID string) error {
	if e.closed {
		return fmt.Errorf("registry is closed")
	}

	// 先获取服务信息
	services, err := e.ListServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	var targetService *ServiceInfo
	for _, service := range services {
		if service.ID == serviceID {
			targetService = service
			break
		}
	}

	if targetService == nil {
		return fmt.Errorf("service not found: %s", serviceID)
	}

	// 删除服务
	servicePath := e.getServicePath(targetService.Name, serviceID)
	_, err = e.client.Delete(ctx, servicePath)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	// 通知监听器
	e.notifyWatchers(ServiceEvent{
		Type:    ServiceEventDeleted,
		Service: targetService,
	})

	return nil
}

// Update 更新服务信息
func (e *EtcdServiceRegistry) Update(ctx context.Context, service *ServiceInfo) error {
	if e.closed {
		return fmt.Errorf("registry is closed")
	}

	// 生成服务路径
	servicePath := e.getServicePath(service.Name, service.ID)
	
	// 序列化服务信息
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %w", err)
	}

	// 更新服务信息
	_, err = e.client.Put(ctx, servicePath, string(data), clientv3.WithLease(e.leaseID))
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	// 通知监听器
	e.notifyWatchers(ServiceEvent{
		Type:    ServiceEventUpdated,
		Service: service,
	})

	return nil
}

// GetService 获取服务信息
func (e *EtcdServiceRegistry) GetService(ctx context.Context, serviceID string) (*ServiceInfo, error) {
	if e.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	services, err := e.ListServices(ctx)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		if service.ID == serviceID {
			return service, nil
		}
	}

	return nil, fmt.Errorf("service not found: %s", serviceID)
}

// ListServices 列出所有服务
func (e *EtcdServiceRegistry) ListServices(ctx context.Context) ([]*ServiceInfo, error) {
	if e.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	// 获取所有服务
	resp, err := e.client.Get(ctx, e.prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	services := make([]*ServiceInfo, 0)
	for _, kv := range resp.Kvs {
		var service ServiceInfo
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			continue // 跳过无效数据
		}
		services = append(services, &service)
	}

	return services, nil
}

// Watch 监听服务变化
func (e *EtcdServiceRegistry) Watch(ctx context.Context) (<-chan ServiceEvent, error) {
	if e.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	eventChan := make(chan ServiceEvent, 100)
	watcherID := fmt.Sprintf("watcher_%d", time.Now().UnixNano())

	e.watcherMutex.Lock()
	e.watchers[watcherID] = eventChan
	e.watcherMutex.Unlock()

	// 启动 etcd 监听
	go func() {
		defer func() {
			e.watcherMutex.Lock()
			delete(e.watchers, watcherID)
			e.watcherMutex.Unlock()
			close(eventChan)
		}()

		watchChan := e.client.Watch(ctx, e.prefix, clientv3.WithPrefix())
		for {
			select {
			case <-ctx.Done():
				return
			case watchResp := <-watchChan:
				for _, ev := range watchResp.Events {
					if ev.Type == clientv3.EventTypeDelete {
						// 服务删除事件
						serviceName, serviceID := e.parseServicePath(string(ev.Kv.Key))
						eventChan <- ServiceEvent{
							Type: ServiceEventDeleted,
							Service: &ServiceInfo{
								ID:   serviceID,
								Name: serviceName,
							},
						}
					} else if ev.Type == clientv3.EventTypePut {
						// 服务创建或更新事件
						var service ServiceInfo
						if err := json.Unmarshal(ev.Kv.Value, &service); err == nil {
							eventType := ServiceEventCreated
							if ev.IsModify() {
								eventType = ServiceEventUpdated
							}
							eventChan <- ServiceEvent{
								Type:    eventType,
								Service: &service,
							}
						}
					}
				}
			}
		}
	}()

	return eventChan, nil
}

// Close 关闭注册中心
func (e *EtcdServiceRegistry) Close() error {
	if e.closed {
		return nil
	}

	e.closed = true

	// 关闭所有监听器
	e.watcherMutex.Lock()
	for _, ch := range e.watchers {
		close(ch)
	}
	e.watchers = make(map[string]chan ServiceEvent)
	e.watcherMutex.Unlock()

	// 关闭 etcd 客户端
	if e.client != nil {
		return e.client.Close()
	}

	return nil
}

// keepAlive 保持租约活跃
func (e *EtcdServiceRegistry) keepAlive() {
	keepAliveCh, err := e.client.KeepAlive(context.Background(), e.leaseID)
	if err != nil {
		return
	}

	for {
		select {
		case <-keepAliveCh:
			// 租约保持活跃
		case <-time.After(10 * time.Second):
			// 超时，重新创建租约
			lease, err := e.client.Grant(context.Background(), 30)
			if err != nil {
				return
			}
			e.leaseID = lease.ID
			return
		}
	}
}

// getServicePath 获取服务路径
func (e *EtcdServiceRegistry) getServicePath(serviceName, serviceID string) string {
	return path.Join(e.prefix, serviceName, serviceID)
}

// parseServicePath 解析服务路径
func (e *EtcdServiceRegistry) parseServicePath(servicePath string) (serviceName, serviceID string) {
	parts := strings.Split(servicePath, "/")
	if len(parts) >= 3 {
		return parts[len(parts)-2], parts[len(parts)-1]
	}
	return "", ""
}

// notifyWatchers 通知监听器
func (e *EtcdServiceRegistry) notifyWatchers(event ServiceEvent) {
	e.watcherMutex.RLock()
	defer e.watcherMutex.RUnlock()

	for _, ch := range e.watchers {
		select {
		case ch <- event:
		default:
			// 通道已满，跳过
		}
	}
} 