package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// NacosServiceRegistry Nacos 服务注册中心
type NacosServiceRegistry struct {
	client     naming_client.INamingClient
	namespace  string
	group      string
	watchers   map[string]chan ServiceEvent
	watcherMutex sync.RWMutex
	closed     bool
}

// NacosConfig Nacos 配置
type NacosConfig struct {
	ServerAddr string        `json:"server_addr"`
	Namespace  string        `json:"namespace"`
	Group      string        `json:"group"`
	Username   string        `json:"username"`
	Password   string        `json:"password"`
	TTL        time.Duration `json:"ttl"`
}

// NewNacosServiceRegistry 创建 Nacos 服务注册中心
func NewNacosServiceRegistry(config *NacosConfig) (*NacosServiceRegistry, error) {
	if config == nil {
		config = &NacosConfig{
			ServerAddr: "localhost:8848",
			Namespace:  "public",
			Group:      "DEFAULT_GROUP",
			TTL:        30 * time.Second,
		}
	}

	// 创建 Nacos 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         config.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            config.Username,
		Password:            config.Password,
	}

	// 创建服务器配置
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: config.ServerAddr,
			Port:   8848,
		},
	}

	// 创建命名服务客户端
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create nacos client: %w", err)
	}

	registry := &NacosServiceRegistry{
		client:   client,
		namespace: config.Namespace,
		group:     config.Group,
		watchers:  make(map[string]chan ServiceEvent),
	}

	return registry, nil
}

// Register 注册服务
func (n *NacosServiceRegistry) Register(ctx context.Context, service *ServiceInfo) error {
	if n.closed {
		return fmt.Errorf("registry is closed")
	}

	// 创建 Nacos 服务实例
	instance := vo.RegisterInstanceParam{
		Ip:          service.Address,
		Port:        uint64(service.Port),
		ServiceName: service.Name,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    n.convertMetadata(service),
	}

	// 注册服务
	success, err := n.client.RegisterInstance(instance)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	if !success {
		return fmt.Errorf("failed to register service: registration failed")
	}

	// 通知监听器
	n.notifyWatchers(ServiceEvent{
		Type:    ServiceEventCreated,
		Service: service,
	})

	return nil
}

// Deregister 注销服务
func (n *NacosServiceRegistry) Deregister(ctx context.Context, serviceID string) error {
	if n.closed {
		return fmt.Errorf("registry is closed")
	}

	// 先获取服务信息
	service, err := n.GetService(ctx, serviceID)
	if err != nil {
		return fmt.Errorf("service not found: %s", serviceID)
	}

	// 注销服务
	instance := vo.DeregisterInstanceParam{
		Ip:          service.Address,
		Port:        uint64(service.Port),
		ServiceName: service.Name,
		Ephemeral:   true,
	}

	success, err := n.client.DeregisterInstance(instance)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	if !success {
		return fmt.Errorf("failed to deregister service: deregistration failed")
	}

	// 通知监听器
	n.notifyWatchers(ServiceEvent{
		Type:    ServiceEventDeleted,
		Service: service,
	})

	return nil
}

// Update 更新服务信息
func (n *NacosServiceRegistry) Update(ctx context.Context, service *ServiceInfo) error {
	if n.closed {
		return fmt.Errorf("registry is closed")
	}

	// 先注销再注册
	err := n.Deregister(ctx, service.ID)
	if err != nil {
		return fmt.Errorf("failed to deregister service for update: %w", err)
	}

	// 重新注册
	err = n.Register(ctx, service)
	if err != nil {
		return fmt.Errorf("failed to register service for update: %w", err)
	}

	// 通知监听器
	n.notifyWatchers(ServiceEvent{
		Type:    ServiceEventUpdated,
		Service: service,
	})

	return nil
}

// GetService 获取服务信息
func (n *NacosServiceRegistry) GetService(ctx context.Context, serviceID string) (*ServiceInfo, error) {
	if n.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	services, err := n.ListServices(ctx)
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
func (n *NacosServiceRegistry) ListServices(ctx context.Context) ([]*ServiceInfo, error) {
	if n.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	// 获取服务列表
	services, err := n.client.GetService(vo.GetServiceParam{
		ServiceName: "",
		GroupName:   n.group,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	result := make([]*ServiceInfo, 0)
	for _, nacosService := range services.Hosts {
		service := &ServiceInfo{
			ID:       fmt.Sprintf("%s-%s-%d", nacosService.ServiceName, nacosService.Ip, nacosService.Port),
			Name:     nacosService.ServiceName,
			Address:  nacosService.Ip,
			Port:     int(nacosService.Port),
			Health:   n.convertHealth(nacosService.Healthy),
			Metadata: nacosService.Metadata,
		}

		// 从元数据中恢复其他字段
		if version, ok := nacosService.Metadata["version"]; ok {
			service.Version = version
		}
		if protocol, ok := nacosService.Metadata["protocol"]; ok {
			service.Protocol = protocol
		}
		if tags, ok := nacosService.Metadata["tags"]; ok {
			// 解析标签
			var tagList []string
			if err := json.Unmarshal([]byte(tags), &tagList); err == nil {
				service.Tags = tagList
			}
		}

		result = append(result, service)
	}

	return result, nil
}

// Watch 监听服务变化
func (n *NacosServiceRegistry) Watch(ctx context.Context) (<-chan ServiceEvent, error) {
	if n.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	eventChan := make(chan ServiceEvent, 100)
	watcherID := fmt.Sprintf("watcher_%d", time.Now().UnixNano())

	n.watcherMutex.Lock()
	n.watchers[watcherID] = eventChan
	n.watcherMutex.Unlock()

	// 启动 Nacos 监听
	go func() {
		defer func() {
			n.watcherMutex.Lock()
			delete(n.watchers, watcherID)
			n.watcherMutex.Unlock()
			close(eventChan)
		}()

		// 简化的监听机制（实际应用中可以使用 Nacos 的订阅功能）
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 定期检查服务变化
				services, err := n.ListServices(ctx)
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
func (n *NacosServiceRegistry) Close() error {
	if n.closed {
		return nil
	}

	n.closed = true

	// 关闭所有监听器
	n.watcherMutex.Lock()
	for _, ch := range n.watchers {
		close(ch)
	}
	n.watchers = make(map[string]chan ServiceEvent)
	n.watcherMutex.Unlock()

	return nil
}

// convertMetadata 转换元数据
func (n *NacosServiceRegistry) convertMetadata(service *ServiceInfo) map[string]string {
	metadata := make(map[string]string)
	
	// 复制现有元数据
	for k, v := range service.Metadata {
		metadata[k] = v
	}

	// 添加标准字段
	metadata["version"] = service.Version
	metadata["protocol"] = service.Protocol
	metadata["health"] = service.Health
	metadata["created_at"] = service.CreatedAt.Format(time.RFC3339)
	metadata["updated_at"] = service.UpdatedAt.Format(time.RFC3339)
	metadata["last_check"] = service.LastCheck.Format(time.RFC3339)
	metadata["ttl"] = service.TTL.String()

	// 序列化标签
	if len(service.Tags) > 0 {
		tagsBytes, _ := json.Marshal(service.Tags)
		metadata["tags"] = string(tagsBytes)
	}

	return metadata
}

// convertHealth 转换健康状态
func (n *NacosServiceRegistry) convertHealth(healthy bool) string {
	if healthy {
		return "healthy"
	}
	return "unhealthy"
}

// notifyWatchers 通知监听器
func (n *NacosServiceRegistry) notifyWatchers(event ServiceEvent) {
	n.watcherMutex.RLock()
	defer n.watcherMutex.RUnlock()

	for _, ch := range n.watchers {
		select {
		case ch <- event:
		default:
			// 通道已满，跳过
		}
	}
} 