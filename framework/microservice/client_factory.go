package microservice

import (
	"context"
	"fmt"
	"time"
)

// ProtocolType 协议类型
type ProtocolType string

const (
	ProtocolHTTP ProtocolType = "http"
	ProtocolGRPC ProtocolType = "grpc"
	ProtocolTCP  ProtocolType = "tcp"
)

// ServiceClientFactory 服务客户端工厂
type ServiceClientFactory struct {
	discovery ServiceDiscovery
	clients   map[ProtocolType]interface{}
}

// NewServiceClientFactory 创建服务客户端工厂
func NewServiceClientFactory(discovery ServiceDiscovery) *ServiceClientFactory {
	return &ServiceClientFactory{
		discovery: discovery,
		clients:   make(map[ProtocolType]interface{}),
	}
}

// GetHTTPClient 获取 HTTP 客户端
func (f *ServiceClientFactory) GetHTTPClient(options ...ServiceClientOption) *ServiceClient {
	if client, exists := f.clients[ProtocolHTTP]; exists {
		return client.(*ServiceClient)
	}

	client := NewServiceClient(f.discovery, options...)
	f.clients[ProtocolHTTP] = client
	return client
}

// GetGRPCClient 获取 gRPC 客户端
func (f *ServiceClientFactory) GetGRPCClient(options ...GRPCServiceClientOption) *GRPCServiceClient {
	if client, exists := f.clients[ProtocolGRPC]; exists {
		return client.(*GRPCServiceClient)
	}

	client := NewGRPCServiceClient(f.discovery, options...)
	f.clients[ProtocolGRPC] = client
	return client
}

// GetClient 根据协议类型获取客户端
func (f *ServiceClientFactory) GetClient(protocol ProtocolType) (interface{}, error) {
	switch protocol {
	case ProtocolHTTP:
		return f.GetHTTPClient(), nil
	case ProtocolGRPC:
		return f.GetGRPCClient(), nil
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// Close 关闭所有客户端
func (f *ServiceClientFactory) Close() error {
	var lastErr error

	// 关闭 HTTP 客户端
	if httpClient, exists := f.clients[ProtocolHTTP]; exists {
		// HTTP 客户端没有 Close 方法，跳过
		delete(f.clients, ProtocolHTTP)
	}

	// 关闭 gRPC 客户端
	if grpcClient, exists := f.clients[ProtocolGRPC]; exists {
		if err := grpcClient.(*GRPCServiceClient).Close(); err != nil {
			lastErr = err
		}
		delete(f.clients, ProtocolGRPC)
	}

	return lastErr
}

// UnifiedServiceClient 统一服务客户端
type UnifiedServiceClient struct {
	factory *ServiceClientFactory
}

// NewUnifiedServiceClient 创建统一服务客户端
func NewUnifiedServiceClient(discovery ServiceDiscovery) *UnifiedServiceClient {
	return &UnifiedServiceClient{
		factory: NewServiceClientFactory(discovery),
	}
}

// Call 统一调用服务
func (u *UnifiedServiceClient) Call(ctx context.Context, serviceName, protocol, method, path string, request, response interface{}) error {
	switch ProtocolType(protocol) {
	case ProtocolHTTP:
		client := u.factory.GetHTTPClient()
		return client.CallJSON(ctx, serviceName, method, path, request, response)
	case ProtocolGRPC:
		client := u.factory.GetGRPCClient()
		return client.CallGRPC(ctx, serviceName, method, request, response, nil)
	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// Get 统一 GET 请求
func (u *UnifiedServiceClient) Get(ctx context.Context, serviceName, protocol, path string, response interface{}) error {
	switch ProtocolType(protocol) {
	case ProtocolHTTP:
		client := u.factory.GetHTTPClient()
		return client.GetJSON(ctx, serviceName, path, response)
	case ProtocolGRPC:
		client := u.factory.GetGRPCClient()
		return client.CallGRPC(ctx, serviceName, "GET", nil, response, nil)
	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// Post 统一 POST 请求
func (u *UnifiedServiceClient) Post(ctx context.Context, serviceName, protocol, path string, request, response interface{}) error {
	switch ProtocolType(protocol) {
	case ProtocolHTTP:
		client := u.factory.GetHTTPClient()
		return client.PostJSON(ctx, serviceName, path, request, response)
	case ProtocolGRPC:
		client := u.factory.GetGRPCClient()
		return client.CallGRPC(ctx, serviceName, "POST", request, response, nil)
	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// Put 统一 PUT 请求
func (u *UnifiedServiceClient) Put(ctx context.Context, serviceName, protocol, path string, request, response interface{}) error {
	switch ProtocolType(protocol) {
	case ProtocolHTTP:
		client := u.factory.GetHTTPClient()
		return client.PutJSON(ctx, serviceName, path, request, response)
	case ProtocolGRPC:
		client := u.factory.GetGRPCClient()
		return client.CallGRPC(ctx, serviceName, "PUT", request, response, nil)
	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// Delete 统一 DELETE 请求
func (u *UnifiedServiceClient) Delete(ctx context.Context, serviceName, protocol, path string, response interface{}) error {
	switch ProtocolType(protocol) {
	case ProtocolHTTP:
		client := u.factory.GetHTTPClient()
		return client.DeleteJSON(ctx, serviceName, path, response)
	case ProtocolGRPC:
		client := u.factory.GetGRPCClient()
		return client.CallGRPC(ctx, serviceName, "DELETE", nil, response, nil)
	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// Close 关闭客户端
func (u *UnifiedServiceClient) Close() error {
	return u.factory.Close()
}

// ServiceHealthChecker 服务健康检查器
type ServiceHealthChecker struct {
	httpChecker *HTTPHealthChecker
	grpcChecker *GRPCHealthChecker
}

// NewServiceHealthChecker 创建服务健康检查器
func NewServiceHealthChecker(timeout time.Duration) *ServiceHealthChecker {
	return &ServiceHealthChecker{
		httpChecker: NewHTTPHealthChecker(timeout),
		grpcChecker: NewGRPCHealthChecker(timeout),
	}
}

// Check 检查服务健康状态
func (h *ServiceHealthChecker) Check(ctx context.Context, service *ServiceInfo) error {
	switch ProtocolType(service.Protocol) {
	case ProtocolHTTP:
		return h.httpChecker.Check(ctx, service)
	case ProtocolGRPC:
		return h.grpcChecker.Check(ctx, service)
	default:
		return fmt.Errorf("unsupported protocol: %s", service.Protocol)
	}
}

// IsHealthy 判断服务是否健康
func (h *ServiceHealthChecker) IsHealthy(service *ServiceInfo) bool {
	switch ProtocolType(service.Protocol) {
	case ProtocolHTTP:
		return h.httpChecker.IsHealthy(service)
	case ProtocolGRPC:
		return h.grpcChecker.IsHealthy(service)
	default:
		return false
	}
}
