package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

// ConsulServiceRegistry Consul 服务注册中心
type ConsulServiceRegistry struct {
	client     *api.Client
	prefix     string
	watchers   map[string]chan ServiceEvent
	watcherMutex sync.RWMutex
	closed     bool
}

// ConsulConfig Consul 配置
type ConsulConfig struct {
	Address    string        `json:"address"`
	Token      string        `json:"token"`
	Datacenter string        `json:"datacenter"`
	Prefix     string        `json:"prefix"`
	TTL        time.Duration `json:"ttl"`
}

// NewConsulServiceRegistry 创建 Consul 服务注册中心
func NewConsulServiceRegistry(config *ConsulConfig) (*ConsulServiceRegistry, error) {
	if config == nil {
		config = &ConsulConfig{
			Address: "localhost:8500",
			Prefix:  "laravel-go/services",
			TTL:     30 * time.Second,
		}
	}

	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.Address
	consulConfig.Token = config.Token
	consulConfig.Datacenter = config.Datacenter

	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	registry := &ConsulServiceRegistry{
		client:   client,
		prefix:   config.Prefix,
		watchers: make(map[string]chan ServiceEvent),
	}

	return registry, nil
}

// Register 注册服务
func (c *ConsulServiceRegistry) Register(ctx context.Context, service *ServiceInfo) error {
	if c.closed {
		return fmt.Errorf("registry is closed")
	}

	// 创建 Consul 服务注册
	registration := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Address: service.Address,
		Port:    service.Port,
		Tags:    service.Tags,
		Meta:    service.Metadata,
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("%s://%s:%d/health", service.Protocol, service.Address, service.Port),
			Interval:                       c.formatDuration(c.getCheckInterval()),
			Timeout:                        c.formatDuration(5 * time.Second),
			DeregisterCriticalServiceAfter: c.formatDuration(c.getTTL()),
		},
	}

	// 注册服务
	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	// 存储服务元数据到 KV 存储
	metadata := map[string]interface{}{
		"version":    service.Version,
		"protocol":   service.Protocol,
		"health":     service.Health,
		"created_at": service.CreatedAt,
		"updated_at": service.UpdatedAt,
		"last_check": service.LastCheck,
		"ttl":        service.TTL.String(),
	}

	metadataBytes, _ := json.Marshal(metadata)
	_, err = c.client.KV().Put(&api.KVPair{
		Key:   c.getMetadataKey(service.Name, service.ID),
		Value: metadataBytes,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to store service metadata: %w", err)
	}

	// 通知监听器
	c.notifyWatchers(ServiceEvent{
		Type:    ServiceEventCreated,
		Service: service,
	})

	return nil
}

// Deregister 注销服务
func (c *ConsulServiceRegistry) Deregister(ctx context.Context, serviceID string) error {
	if c.closed {
		return fmt.Errorf("registry is closed")
	}

	// 先获取服务信息
	service, err := c.GetService(ctx, serviceID)
	if err != nil {
		return fmt.Errorf("service not found: %s", serviceID)
	}

	// 注销服务
	err = c.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	// 删除元数据
	_, err = c.client.KV().Delete(c.getMetadataKey(service.Name, serviceID), nil)
	if err != nil {
		// 忽略元数据删除错误
	}

	// 通知监听器
	c.notifyWatchers(ServiceEvent{
		Type:    ServiceEventDeleted,
		Service: service,
	})

	return nil
}

// Update 更新服务信息
func (c *ConsulServiceRegistry) Update(ctx context.Context, service *ServiceInfo) error {
	if c.closed {
		return fmt.Errorf("registry is closed")
	}

	// 更新服务注册
	registration := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Address: service.Address,
		Port:    service.Port,
		Tags:    service.Tags,
		Meta:    service.Metadata,
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("%s://%s:%d/health", service.Protocol, service.Address, service.Port),
			Interval:                       c.formatDuration(c.getCheckInterval()),
			Timeout:                        c.formatDuration(5 * time.Second),
			DeregisterCriticalServiceAfter: c.formatDuration(c.getTTL()),
		},
	}

	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	// 更新元数据
	metadata := map[string]interface{}{
		"version":    service.Version,
		"protocol":   service.Protocol,
		"health":     service.Health,
		"created_at": service.CreatedAt,
		"updated_at": service.UpdatedAt,
		"last_check": service.LastCheck,
		"ttl":        service.TTL.String(),
	}

	metadataBytes, _ := json.Marshal(metadata)
	_, err = c.client.KV().Put(&api.KVPair{
		Key:   c.getMetadataKey(service.Name, service.ID),
		Value: metadataBytes,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to update service metadata: %w", err)
	}

	// 通知监听器
	c.notifyWatchers(ServiceEvent{
		Type:    ServiceEventUpdated,
		Service: service,
	})

	return nil
}

// GetService 获取服务信息
func (c *ConsulServiceRegistry) GetService(ctx context.Context, serviceID string) (*ServiceInfo, error) {
	if c.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	services, err := c.ListServices(ctx)
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
func (c *ConsulServiceRegistry) ListServices(ctx context.Context) ([]*ServiceInfo, error) {
	if c.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	// 获取所有服务
	services, err := c.client.Agent().Services()
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	result := make([]*ServiceInfo, 0)
	for _, consulService := range services {
		// 跳过 Consul 内部服务
		if consulService.Service == "consul" {
			continue
		}

		service := &ServiceInfo{
			ID:       consulService.ID,
			Name:     consulService.Service,
			Address:  consulService.Address,
			Port:     consulService.Port,
			Tags:     consulService.Tags,
			Metadata: consulService.Meta,
		}

		// 获取元数据
		metadata, err := c.getServiceMetadata(consulService.Service, consulService.ID)
		if err == nil {
			service.Version = metadata["version"].(string)
			service.Protocol = metadata["protocol"].(string)
			service.Health = metadata["health"].(string)
			if createdAt, ok := metadata["created_at"].(string); ok {
				service.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
			}
			if updatedAt, ok := metadata["updated_at"].(string); ok {
				service.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
			}
			if lastCheck, ok := metadata["last_check"].(string); ok {
				service.LastCheck, _ = time.Parse(time.RFC3339, lastCheck)
			}
			if ttlStr, ok := metadata["ttl"].(string); ok {
				service.TTL, _ = time.ParseDuration(ttlStr)
			}
		}

		result = append(result, service)
	}

	return result, nil
}

// Watch 监听服务变化
func (c *ConsulServiceRegistry) Watch(ctx context.Context) (<-chan ServiceEvent, error) {
	if c.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	eventChan := make(chan ServiceEvent, 100)
	watcherID := fmt.Sprintf("watcher_%d", time.Now().UnixNano())

	c.watcherMutex.Lock()
	c.watchers[watcherID] = eventChan
	c.watcherMutex.Unlock()

	// 启动 Consul 监听
	go func() {
		defer func() {
			c.watcherMutex.Lock()
			delete(c.watchers, watcherID)
			c.watcherMutex.Unlock()
			close(eventChan)
		}()

		// 简化的监听机制（实际应用中可以使用 Consul 的 watch 功能）
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 定期检查服务变化
				services, err := c.ListServices(ctx)
				if err == nil {
					for _, service := range services {
						eventChan <- ServiceEvent{
							Type:    ServiceEventUpdated,
							Service: service,
						}
					}
				}
			}
		}
	}()

	return eventChan, nil
}

// Close 关闭注册中心
func (c *ConsulServiceRegistry) Close() error {
	if c.closed {
		return nil
	}

	c.closed = true

	// 关闭所有监听器
	c.watcherMutex.Lock()
	for _, ch := range c.watchers {
		close(ch)
	}
	c.watchers = make(map[string]chan ServiceEvent)
	c.watcherMutex.Unlock()

	return nil
}

// getServiceMetadata 获取服务元数据
func (c *ConsulServiceRegistry) getServiceMetadata(serviceName, serviceID string) (map[string]interface{}, error) {
	pair, _, err := c.client.KV().Get(c.getMetadataKey(serviceName, serviceID), nil)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return make(map[string]interface{}), nil
	}

	var metadata map[string]interface{}
	err = json.Unmarshal(pair.Value, &metadata)
	return metadata, err
}

// getMetadataKey 获取元数据键
func (c *ConsulServiceRegistry) getMetadataKey(serviceName, serviceID string) string {
	return fmt.Sprintf("%s/%s/%s/metadata", c.prefix, serviceName, serviceID)
}

// getCheckInterval 获取检查间隔
func (c *ConsulServiceRegistry) getCheckInterval() time.Duration {
	return 10 * time.Second
}

// getTTL 获取 TTL
func (c *ConsulServiceRegistry) getTTL() time.Duration {
	return 30 * time.Second
}

// formatDuration 格式化持续时间
func (c *ConsulServiceRegistry) formatDuration(d time.Duration) string {
	return d.String()
}

// notifyWatchers 通知监听器
func (c *ConsulServiceRegistry) notifyWatchers(event ServiceEvent) {
	c.watcherMutex.RLock()
	defer c.watcherMutex.RUnlock()

	for _, ch := range c.watchers {
		select {
		case ch <- event:
		default:
			// 通道已满，跳过
		}
	}
} 