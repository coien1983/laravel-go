package microservice

import (
	"context"
	"time"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Address   string            `json:"address"`
	Port      int               `json:"port"`
	Protocol  string            `json:"protocol"` // http, grpc, tcp
	Health    string            `json:"health"`   // healthy, unhealthy, unknown
	Metadata  map[string]string `json:"metadata"`
	Tags      []string          `json:"tags"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	LastCheck time.Time         `json:"last_check"`
	TTL       time.Duration     `json:"ttl"`
}

// ServiceRegistry 服务注册接口
type ServiceRegistry interface {
	// Register 注册服务
	Register(ctx context.Context, service *ServiceInfo) error

	// Deregister 注销服务
	Deregister(ctx context.Context, serviceID string) error

	// Update 更新服务信息
	Update(ctx context.Context, service *ServiceInfo) error

	// GetService 获取服务信息
	GetService(ctx context.Context, serviceID string) (*ServiceInfo, error)

	// ListServices 列出所有服务
	ListServices(ctx context.Context) ([]*ServiceInfo, error)

	// Watch 监听服务变化
	Watch(ctx context.Context) (<-chan ServiceEvent, error)

	// Close 关闭注册中心
	Close() error
}

// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
	// Discover 发现服务
	Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error)

	// DiscoverOne 发现单个服务（负载均衡）
	DiscoverOne(ctx context.Context, serviceName string) (*ServiceInfo, error)

	// Watch 监听服务变化
	Watch(ctx context.Context, serviceName string) (<-chan ServiceEvent, error)

	// Close 关闭发现服务
	Close() error
}

// ServiceEvent 服务事件
type ServiceEvent struct {
	Type    ServiceEventType `json:"type"`
	Service *ServiceInfo     `json:"service"`
}

// ServiceEventType 服务事件类型
type ServiceEventType string

const (
	ServiceEventCreated ServiceEventType = "created"
	ServiceEventUpdated ServiceEventType = "updated"
	ServiceEventDeleted ServiceEventType = "deleted"
)

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
	// Select 选择一个服务实例
	Select(services []*ServiceInfo) *ServiceInfo
}

// RoundRobinLoadBalancer 轮询负载均衡器
type RoundRobinLoadBalancer struct {
	current int
}

// NewRoundRobinLoadBalancer 创建轮询负载均衡器
func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{current: 0}
}

// Select 轮询选择服务
func (rr *RoundRobinLoadBalancer) Select(services []*ServiceInfo) *ServiceInfo {
	if len(services) == 0 {
		return nil
	}

	// 过滤健康服务
	healthyServices := make([]*ServiceInfo, 0)
	for _, service := range services {
		if service.Health == "healthy" {
			healthyServices = append(healthyServices, service)
		}
	}

	if len(healthyServices) == 0 {
		return nil
	}

	service := healthyServices[rr.current%len(healthyServices)]
	rr.current++
	return service
}

// RandomLoadBalancer 随机负载均衡器
type RandomLoadBalancer struct{}

// NewRandomLoadBalancer 创建随机负载均衡器
func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{}
}

// Select 随机选择服务
func (rr *RandomLoadBalancer) Select(services []*ServiceInfo) *ServiceInfo {
	if len(services) == 0 {
		return nil
	}

	// 过滤健康服务
	healthyServices := make([]*ServiceInfo, 0)
	for _, service := range services {
		if service.Health == "healthy" {
			healthyServices = append(healthyServices, service)
		}
	}

	if len(healthyServices) == 0 {
		return nil
	}

	// 简单的随机选择（实际应用中可以使用更好的随机算法）
	index := time.Now().UnixNano() % int64(len(healthyServices))
	return healthyServices[index]
}

// HealthChecker 健康检查接口
type HealthChecker interface {
	// Check 检查服务健康状态
	Check(ctx context.Context, service *ServiceInfo) error

	// IsHealthy 判断服务是否健康
	IsHealthy(service *ServiceInfo) bool
}

// HTTPHealthChecker HTTP 健康检查器
type HTTPHealthChecker struct {
	Timeout time.Duration
}

// NewHTTPHealthChecker 创建 HTTP 健康检查器
func NewHTTPHealthChecker(timeout time.Duration) *HTTPHealthChecker {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &HTTPHealthChecker{Timeout: timeout}
}

// Check 检查 HTTP 服务健康状态
func (h *HTTPHealthChecker) Check(ctx context.Context, service *ServiceInfo) error {
	// 这里应该实现实际的 HTTP 健康检查逻辑
	// 为了简化，这里只是标记为健康
	service.Health = "healthy"
	service.LastCheck = time.Now()
	return nil
}

// IsHealthy 判断服务是否健康
func (h *HTTPHealthChecker) IsHealthy(service *ServiceInfo) bool {
	return service.Health == "healthy"
}
