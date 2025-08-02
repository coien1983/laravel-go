package microservice

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcmetadata "google.golang.org/grpc/metadata"
)

// GRPCServiceClient gRPC 服务通信客户端
type GRPCServiceClient struct {
	discovery     ServiceDiscovery
	connections   map[string]*grpc.ClientConn
	connectionMux sync.RWMutex
	timeout       time.Duration
	retryCount    int
	retryDelay    time.Duration
}

// NewGRPCServiceClient 创建 gRPC 服务通信客户端
func NewGRPCServiceClient(discovery ServiceDiscovery, options ...GRPCServiceClientOption) *GRPCServiceClient {
	client := &GRPCServiceClient{
		discovery:   discovery,
		connections: make(map[string]*grpc.ClientConn),
		timeout:     30 * time.Second,
		retryCount:  3,
		retryDelay:  1 * time.Second,
	}

	// 应用选项
	for _, option := range options {
		option(client)
	}

	return client
}

// GRPCServiceClientOption gRPC 服务客户端选项
type GRPCServiceClientOption func(*GRPCServiceClient)

// WithGRPCTimeout 设置 gRPC 超时时间
func WithGRPCTimeout(timeout time.Duration) GRPCServiceClientOption {
	return func(c *GRPCServiceClient) {
		c.timeout = timeout
	}
}

// WithGRPCRetry 设置 gRPC 重试参数
func WithGRPCRetry(count int, delay time.Duration) GRPCServiceClientOption {
	return func(c *GRPCServiceClient) {
		c.retryCount = count
		c.retryDelay = delay
	}
}

// getConnection 获取或创建 gRPC 连接
func (c *GRPCServiceClient) getConnection(ctx context.Context, serviceName string) (*grpc.ClientConn, error) {
	c.connectionMux.RLock()
	if conn, exists := c.connections[serviceName]; exists {
		c.connectionMux.RUnlock()
		return conn, nil
	}
	c.connectionMux.RUnlock()

	// 发现服务
	service, err := c.discovery.DiscoverOne(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	// 创建连接
	address := fmt.Sprintf("%s:%d", service.Address, service.Port)
	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(c.timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	// 缓存连接
	c.connectionMux.Lock()
	c.connections[serviceName] = conn
	c.connectionMux.Unlock()

	return conn, nil
}

// CallGRPC 调用 gRPC 服务
func (c *GRPCServiceClient) CallGRPC(ctx context.Context, serviceName, method string, request, response interface{}, metadata map[string]string) error {
	// 添加超时上下文
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// 添加元数据
	if metadata != nil {
		md := grpcmetadata.New(metadata)
		ctx = grpcmetadata.NewOutgoingContext(ctx, md)
	}

	// 获取连接
	conn, err := c.getConnection(ctx, serviceName)
	if err != nil {
		return err
	}

	// 执行 gRPC 调用
	var lastErr error
	for i := 0; i <= c.retryCount; i++ {
		err := conn.Invoke(ctx, method, request, response)
		if err == nil {
			return nil
		}
		lastErr = err

		// 如果不是最后一次重试，等待后继续
		if i < c.retryCount {
			time.Sleep(c.retryDelay)
		}
	}

	return fmt.Errorf("gRPC call failed after %d retries: %w", c.retryCount, lastErr)
}

// StreamGRPC 流式 gRPC 调用
func (c *GRPCServiceClient) StreamGRPC(ctx context.Context, serviceName, method string, request interface{}, metadata map[string]string) (grpc.ClientStream, error) {
	// 添加超时上下文
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// 添加元数据
	if metadata != nil {
		md := grpcmetadata.New(metadata)
		ctx = grpcmetadata.NewOutgoingContext(ctx, md)
	}

	// 获取连接
	conn, err := c.getConnection(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	// 创建流
	stream, err := conn.NewStream(ctx, &grpc.StreamDesc{}, method)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	return stream, nil
}

// Close 关闭所有连接
func (c *GRPCServiceClient) Close() error {
	c.connectionMux.Lock()
	defer c.connectionMux.Unlock()

	var lastErr error
	for serviceName, conn := range c.connections {
		if err := conn.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close connection for %s: %w", serviceName, err)
		}
		delete(c.connections, serviceName)
	}

	return lastErr
}

// GRPCHealthChecker gRPC 健康检查器
type GRPCHealthChecker struct {
	Timeout time.Duration
}

// NewGRPCHealthChecker 创建 gRPC 健康检查器
func NewGRPCHealthChecker(timeout time.Duration) *GRPCHealthChecker {
	return &GRPCHealthChecker{
		Timeout: timeout,
	}
}

// Check 检查 gRPC 服务健康状态
func (h *GRPCHealthChecker) Check(ctx context.Context, service *ServiceInfo) error {
	if service.Protocol != "grpc" {
		return fmt.Errorf("service %s is not a gRPC service", service.Name)
	}

	// 创建连接
	address := fmt.Sprintf("%s:%d", service.Address, service.Port)
	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(h.Timeout),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC service %s: %w", service.Name, err)
	}
	defer conn.Close()

	// 这里可以调用 gRPC 健康检查服务
	// 例如: healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{})

	return nil
}

// IsHealthy 判断 gRPC 服务是否健康
func (h *GRPCHealthChecker) IsHealthy(service *ServiceInfo) bool {
	ctx, cancel := context.WithTimeout(context.Background(), h.Timeout)
	defer cancel()

	return h.Check(ctx, service) == nil
}
