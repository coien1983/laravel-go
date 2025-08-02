package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
)

// ZookeeperServiceRegistry Zookeeper 服务注册中心
type ZookeeperServiceRegistry struct {
	conn       *zk.Conn
	prefix     string
	watchers   map[string]chan ServiceEvent
	watcherMutex sync.RWMutex
	closed     bool
}

// ZookeeperConfig Zookeeper 配置
type ZookeeperConfig struct {
	Servers   []string      `json:"servers"`
	Prefix    string        `json:"prefix"`
	TTL       time.Duration `json:"ttl"`
	SessionTimeout time.Duration `json:"session_timeout"`
}

// NewZookeeperServiceRegistry 创建 Zookeeper 服务注册中心
func NewZookeeperServiceRegistry(config *ZookeeperConfig) (*ZookeeperServiceRegistry, error) {
	if config == nil {
		config = &ZookeeperConfig{
			Servers:        []string{"localhost:2181"},
			Prefix:         "/laravel-go/services",
			TTL:            30 * time.Second,
			SessionTimeout: 10 * time.Second,
		}
	}

	// 连接 Zookeeper
	conn, _, err := zk.Connect(config.Servers, config.SessionTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to zookeeper: %w", err)
	}

	registry := &ZookeeperServiceRegistry{
		conn:     conn,
		prefix:   config.Prefix,
		watchers: make(map[string]chan ServiceEvent),
	}

	// 确保根路径存在
	err = registry.ensurePath(config.Prefix)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create root path: %w", err)
	}

	return registry, nil
}

// Register 注册服务
func (z *ZookeeperServiceRegistry) Register(ctx context.Context, service *ServiceInfo) error {
	if z.closed {
		return fmt.Errorf("registry is closed")
	}

	// 生成服务路径
	servicePath := z.getServicePath(service.Name, service.ID)
	
	// 序列化服务信息
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %w", err)
	}

	// 创建服务节点（临时节点，会话结束后自动删除）
	_, err = z.conn.Create(servicePath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	// 通知监听器
	z.notifyWatchers(ServiceEvent{
		Type:    ServiceEventCreated,
		Service: service,
	})

	return nil
}

// Deregister 注销服务
func (z *ZookeeperServiceRegistry) Deregister(ctx context.Context, serviceID string) error {
	if z.closed {
		return fmt.Errorf("registry is closed")
	}

	// 先获取服务信息
	services, err := z.ListServices(ctx)
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

	// 删除服务节点
	servicePath := z.getServicePath(targetService.Name, serviceID)
	err = z.conn.Delete(servicePath, -1)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	// 通知监听器
	z.notifyWatchers(ServiceEvent{
		Type:    ServiceEventDeleted,
		Service: targetService,
	})

	return nil
}

// Update 更新服务信息
func (z *ZookeeperServiceRegistry) Update(ctx context.Context, service *ServiceInfo) error {
	if z.closed {
		return fmt.Errorf("registry is closed")
	}

	// 生成服务路径
	servicePath := z.getServicePath(service.Name, service.ID)
	
	// 序列化服务信息
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %w", err)
	}

	// 检查节点是否存在
	exists, _, err := z.conn.Exists(servicePath)
	if err != nil {
		return fmt.Errorf("failed to check service existence: %w", err)
	}

	if !exists {
		// 节点不存在，创建新节点
		_, err = z.conn.Create(servicePath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}
	} else {
		// 节点存在，更新数据
		_, err = z.conn.Set(servicePath, data, -1)
		if err != nil {
			return fmt.Errorf("failed to update service: %w", err)
		}
	}

	// 通知监听器
	z.notifyWatchers(ServiceEvent{
		Type:    ServiceEventUpdated,
		Service: service,
	})

	return nil
}

// GetService 获取服务信息
func (z *ZookeeperServiceRegistry) GetService(ctx context.Context, serviceID string) (*ServiceInfo, error) {
	if z.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	services, err := z.ListServices(ctx)
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
func (z *ZookeeperServiceRegistry) ListServices(ctx context.Context) ([]*ServiceInfo, error) {
	if z.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	// 获取所有服务
	services, err := z.getAllServices(z.prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	return services, nil
}

// Watch 监听服务变化
func (z *ZookeeperServiceRegistry) Watch(ctx context.Context) (<-chan ServiceEvent, error) {
	if z.closed {
		return nil, fmt.Errorf("registry is closed")
	}

	eventChan := make(chan ServiceEvent, 100)
	watcherID := fmt.Sprintf("watcher_%d", time.Now().UnixNano())

	z.watcherMutex.Lock()
	z.watchers[watcherID] = eventChan
	z.watcherMutex.Unlock()

	// 启动 Zookeeper 监听
	go func() {
		defer func() {
			z.watcherMutex.Lock()
			delete(z.watchers, watcherID)
			z.watcherMutex.Unlock()
			close(eventChan)
		}()

		// 监听根路径变化
		z.watchPath(z.prefix, eventChan, ctx)
	}()

	return eventChan, nil
}

// Close 关闭注册中心
func (z *ZookeeperServiceRegistry) Close() error {
	if z.closed {
		return nil
	}

	z.closed = true

	// 关闭所有监听器
	z.watcherMutex.Lock()
	for _, ch := range z.watchers {
		close(ch)
	}
	z.watchers = make(map[string]chan ServiceEvent)
	z.watcherMutex.Unlock()

	// 关闭 Zookeeper 连接
	if z.conn != nil {
		z.conn.Close()
	}

	return nil
}

// getAllServices 递归获取所有服务
func (z *ZookeeperServiceRegistry) getAllServices(basePath string) ([]*ServiceInfo, error) {
	services := make([]*ServiceInfo, 0)

	// 获取子节点
	children, _, err := z.conn.Children(basePath)
	if err != nil {
		return services, err
	}

	for _, child := range children {
		childPath := path.Join(basePath, child)
		
		// 检查是否是服务节点（有数据的节点）
		data, _, err := z.conn.Get(childPath)
		if err != nil {
			continue
		}

		if len(data) > 0 {
			// 这是一个服务节点
			var service ServiceInfo
			if err := json.Unmarshal(data, &service); err == nil {
				services = append(services, &service)
			}
		} else {
			// 这是一个目录节点，递归获取
			subServices, err := z.getAllServices(childPath)
			if err == nil {
				services = append(services, subServices...)
			}
		}
	}

	return services, nil
}

// watchPath 监听路径变化
func (z *ZookeeperServiceRegistry) watchPath(path string, eventChan chan<- ServiceEvent, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 监听子节点变化
			children, _, events, err := z.conn.ChildrenW(path)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			// 处理当前子节点
			for _, child := range children {
				childPath := path + "/" + child
				z.watchServiceNode(childPath, eventChan)
			}

			// 等待事件
			select {
			case <-ctx.Done():
				return
			case event := <-events:
				if event.Type == zk.EventNodeChildrenChanged {
					// 子节点变化，重新获取子节点列表
					continue
				}
			}
		}
	}
}

// watchServiceNode 监听服务节点变化
func (z *ZookeeperServiceRegistry) watchServiceNode(nodePath string, eventChan chan<- ServiceEvent) {
	go func() {
		for {
			data, _, events, err := z.conn.GetW(nodePath)
			if err != nil {
				// 节点可能被删除
				return
			}

			if len(data) > 0 {
				var service ServiceInfo
				if err := json.Unmarshal(data, &service); err == nil {
					eventChan <- ServiceEvent{
						Type:    ServiceEventUpdated,
						Service: &service,
					}
				}
			}

			select {
			case event := <-events:
				if event.Type == zk.EventNodeDeleted {
					// 节点被删除
					serviceName, serviceID := z.parseServicePath(nodePath)
					eventChan <- ServiceEvent{
						Type: ServiceEventDeleted,
						Service: &ServiceInfo{
							ID:   serviceID,
							Name: serviceName,
						},
					}
					return
				}
			}
		}
	}()
}

// ensurePath 确保路径存在
func (z *ZookeeperServiceRegistry) ensurePath(path string) error {
	parts := strings.Split(path, "/")
	currentPath := ""
	
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		currentPath += "/" + part
		exists, _, err := z.conn.Exists(currentPath)
		if err != nil {
			return err
		}
		
		if !exists {
			_, err = z.conn.Create(currentPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
			if err != nil && err != zk.ErrNodeExists {
				return err
			}
		}
	}
	
	return nil
}

// getServicePath 获取服务路径
func (z *ZookeeperServiceRegistry) getServicePath(serviceName, serviceID string) string {
	return path.Join(z.prefix, serviceName, serviceID)
}

// parseServicePath 解析服务路径
func (z *ZookeeperServiceRegistry) parseServicePath(servicePath string) (serviceName, serviceID string) {
	parts := strings.Split(servicePath, "/")
	if len(parts) >= 3 {
		return parts[len(parts)-2], parts[len(parts)-1]
	}
	return "", ""
}

// notifyWatchers 通知监听器
func (z *ZookeeperServiceRegistry) notifyWatchers(event ServiceEvent) {
	z.watcherMutex.RLock()
	defer z.watcherMutex.RUnlock()

	for _, ch := range z.watchers {
		select {
		case ch <- event:
		default:
			// 通道已满，跳过
		}
	}
} 