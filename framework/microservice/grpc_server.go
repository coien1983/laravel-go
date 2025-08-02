package microservice

import (
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// GRPCServer gRPC 服务器
type GRPCServer struct {
	server             *grpc.Server
	address            string
	port               int
	registry           ServiceRegistry
	healthServer       *health.Server
	interceptors       []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	services           map[string]interface{}
	servicesMux        sync.RWMutex
	options            *GRPCServerOptions
	started            bool
	startedMux         sync.RWMutex
}

// GRPCServerOptions gRPC 服务器选项
type GRPCServerOptions struct {
	// 基础配置
	Address               string
	Port                  int
	MaxConcurrentStreams  uint32
	MaxConnectionIdle     time.Duration
	MaxConnectionAge      time.Duration
	MaxConnectionAgeGrace time.Duration
	Time                  time.Duration
	Timeout               time.Duration

	// 安全配置
	TLSEnabled bool
	CertFile   string
	KeyFile    string

	// 注册中心配置
	Registry       ServiceRegistry
	ServiceName    string
	ServiceVersion string
	Metadata       map[string]string

	// 健康检查配置
	HealthCheckEnabled bool
	HealthCheckPath    string

	// 反射配置
	ReflectionEnabled bool

	// 日志配置
	LoggingEnabled bool
}

// NewGRPCServerOptions 创建 gRPC 服务器选项
func NewGRPCServerOptions() *GRPCServerOptions {
	return &GRPCServerOptions{
		Address:               "0.0.0.0",
		Port:                  50051,
		MaxConcurrentStreams:  100,
		MaxConnectionIdle:     30 * time.Second,
		MaxConnectionAge:      5 * time.Minute,
		MaxConnectionAgeGrace: 10 * time.Second,
		Time:                  2 * time.Hour,
		Timeout:               20 * time.Second,
		TLSEnabled:            false,
		HealthCheckEnabled:    true,
		HealthCheckPath:       "/health",
		ReflectionEnabled:     true,
		LoggingEnabled:        true,
		Metadata:              make(map[string]string),
	}
}

// GRPCServerOption gRPC 服务器选项函数
type GRPCServerOption func(*GRPCServerOptions)

// WithGRPCAddress 设置 gRPC 服务器地址
func WithGRPCAddress(address string) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		o.Address = address
	}
}

// WithGRPCPort 设置 gRPC 服务器端口
func WithGRPCPort(port int) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		o.Port = port
	}
}

// WithGRPCTLS 设置 gRPC TLS 配置
func WithGRPCTLS(certFile, keyFile string) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		o.TLSEnabled = true
		o.CertFile = certFile
		o.KeyFile = keyFile
	}
}

// WithGRPCRegistry 设置 gRPC 注册中心
func WithGRPCRegistry(registry ServiceRegistry) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		o.Registry = registry
	}
}

// WithGRPCServiceInfo 设置 gRPC 服务信息
func WithGRPCServiceInfo(name, version string) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		o.ServiceName = name
		o.ServiceVersion = version
	}
}

// WithGRPCMetadata 设置 gRPC 服务元数据
func WithGRPCMetadata(metadata map[string]string) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		for k, v := range metadata {
			o.Metadata[k] = v
		}
	}
}

// WithGRPCHealthCheck 设置 gRPC 健康检查
func WithGRPCHealthCheck(enabled bool, path string) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		o.HealthCheckEnabled = enabled
		if path != "" {
			o.HealthCheckPath = path
		}
	}
}

// WithGRPCReflection 设置 gRPC 反射
func WithGRPCReflection(enabled bool) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		o.ReflectionEnabled = enabled
	}
}

// WithGRPCLogging 设置 gRPC 日志
func WithGRPCLogging(enabled bool) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		o.LoggingEnabled = enabled
	}
}

// WithGRPCUnaryInterceptor 添加 gRPC 一元拦截器
func WithGRPCUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		// 这里会在创建服务器时应用
	}
}

// WithGRPCStreamInterceptor 添加 gRPC 流拦截器
func WithGRPCStreamInterceptor(interceptor grpc.StreamServerInterceptor) GRPCServerOption {
	return func(o *GRPCServerOptions) {
		// 这里会在创建服务器时应用
	}
}

// NewGRPCServer 创建 gRPC 服务器
func NewGRPCServer(options ...GRPCServerOption) *GRPCServer {
	opts := NewGRPCServerOptions()

	// 应用选项
	for _, option := range options {
		option(opts)
	}

	// 创建 gRPC 服务器选项
	grpcOpts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(opts.MaxConcurrentStreams),
		grpc.MaxConnectionIdle(opts.MaxConnectionIdle),
		grpc.MaxConnectionAge(opts.MaxConnectionAge),
		grpc.MaxConnectionAgeGrace(opts.MaxConnectionAgeGrace),
		grpc.ConnectionTimeout(opts.Timeout),
	}

	// 添加 TLS 配置
	if opts.TLSEnabled {
		creds, err := credentials.NewServerTLSFromFile(opts.CertFile, opts.KeyFile)
		if err != nil {
			panic(fmt.Sprintf("failed to load TLS credentials: %v", err))
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	// 创建 gRPC 服务器
	server := grpc.NewServer(grpcOpts...)

	// 创建健康检查服务器
	var healthServer *health.Server
	if opts.HealthCheckEnabled {
		healthServer = health.NewServer()
		healthpb.RegisterHealthServer(server, healthServer)
	}

	// 启用反射
	if opts.ReflectionEnabled {
		reflection.Register(server)
	}

	return &GRPCServer{
		server:       server,
		address:      opts.Address,
		port:         opts.Port,
		registry:     opts.Registry,
		healthServer: healthServer,
		services:     make(map[string]interface{}),
		options:      opts,
	}
}

// RegisterService 注册 gRPC 服务
func (s *GRPCServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.servicesMux.Lock()
	defer s.servicesMux.Unlock()

	s.server.RegisterService(desc, impl)
	s.services[desc.ServiceName] = impl
}

// AddUnaryInterceptor 添加一元拦截器
func (s *GRPCServer) AddUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) {
	s.interceptors = append(s.interceptors, interceptor)
}

// AddStreamInterceptor 添加流拦截器
func (s *GRPCServer) AddStreamInterceptor(interceptor grpc.StreamServerInterceptor) {
	s.streamInterceptors = append(s.streamInterceptors, interceptor)
}

// SetHealthStatus 设置服务健康状态
func (s *GRPCServer) SetHealthStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	if s.healthServer != nil {
		s.healthServer.SetServingStatus(service, status)
	}
}

// Start 启动 gRPC 服务器
func (s *GRPCServer) Start() error {
	s.startedMux.Lock()
	if s.started {
		s.startedMux.Unlock()
		return fmt.Errorf("gRPC server is already started")
	}
	s.started = true
	s.startedMux.Unlock()

	// 监听端口
	address := fmt.Sprintf("%s:%d", s.address, s.port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	// 注册到服务注册中心
	if s.registry != nil && s.options.ServiceName != "" {
		serviceInfo := &ServiceInfo{
			Name:     s.options.ServiceName,
			Version:  s.options.ServiceVersion,
			Address:  s.address,
			Port:     s.port,
			Protocol: "grpc",
			Metadata: s.options.Metadata,
		}

		if err := s.registry.Register(serviceInfo); err != nil {
			return fmt.Errorf("failed to register service: %w", err)
		}

		// 设置健康状态
		if s.healthServer != nil {
			s.healthServer.SetServingStatus(s.options.ServiceName, healthpb.HealthCheckResponse_SERVING)
		}
	}

	// 启动服务器
	go func() {
		if err := s.server.Serve(lis); err != nil {
			fmt.Printf("gRPC server failed to serve: %v\n", err)
		}
	}()

	fmt.Printf("gRPC server started on %s\n", address)
	return nil
}

// Stop 停止 gRPC 服务器
func (s *GRPCServer) Stop() error {
	s.startedMux.Lock()
	defer s.startedMux.Unlock()

	if !s.started {
		return nil
	}

	// 注销服务
	if s.registry != nil && s.options.ServiceName != "" {
		serviceInfo := &ServiceInfo{
			Name:     s.options.ServiceName,
			Version:  s.options.ServiceVersion,
			Address:  s.address,
			Port:     s.port,
			Protocol: "grpc",
		}

		if err := s.registry.Deregister(serviceInfo); err != nil {
			fmt.Printf("failed to deregister service: %v\n", err)
		}
	}

	// 优雅关闭
	s.server.GracefulStop()
	s.started = false

	fmt.Println("gRPC server stopped")
	return nil
}

// IsRunning 检查服务器是否运行
func (s *GRPCServer) IsRunning() bool {
	s.startedMux.RLock()
	defer s.startedMux.RUnlock()
	return s.started
}

// GetServer 获取底层 gRPC 服务器
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

// GetAddress 获取服务器地址
func (s *GRPCServer) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.address, s.port)
}

// GetServices 获取注册的服务
func (s *GRPCServer) GetServices() map[string]interface{} {
	s.servicesMux.RLock()
	defer s.servicesMux.RUnlock()

	services := make(map[string]interface{})
	for k, v := range s.services {
		services[k] = v
	}
	return services
}
