package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"laravel-go/framework/microservice"
)

// 示例 gRPC 服务实现
type UserService struct {
	users map[string]*User
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserService() *UserService {
	return &UserService{
		users: map[string]*User{
			"1": {ID: "1", Name: "Alice", Email: "alice@example.com"},
			"2": {ID: "2", Name: "Bob", Email: "bob@example.com"},
		},
	}
}

func (s *UserService) GetUser(id string) (*User, error) {
	if user, exists := s.users[id]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found: %s", id)
}

func (s *UserService) CreateUser(user *User) error {
	s.users[user.ID] = user
	return nil
}

// 示例：启动 gRPC 服务器
func startGRPCServer() {
	// 创建注册中心
	registry, err := microservice.NewServiceRegistry(&microservice.RegistryConfig{
		Type: microservice.RegistryTypeMemory,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建 gRPC 服务器
	server := microservice.NewGRPCServer(
		microservice.WithGRPCAddress("0.0.0.0"),
		microservice.WithGRPCPort(50051),
		microservice.WithGRPCRegistry(registry),
		microservice.WithGRPCServiceInfo("user-service", "1.0.0"),
		microservice.WithGRPCHealthCheck(true, "/health"),
		microservice.WithGRPCReflection(true),
		microservice.WithGRPCLogging(true),
	)

	// 添加拦截器
	server.AddUnaryInterceptor(microservice.LoggingInterceptor())
	server.AddUnaryInterceptor(microservice.RecoveryInterceptor())

	// 创建限流器
	limiter := microservice.NewSimpleRateLimiter(map[string]int{
		"/user.UserService/GetUser":    100,
		"/user.UserService/CreateUser": 50,
	})
	server.AddUnaryInterceptor(microservice.RateLimitInterceptor(limiter))

	// 创建指标收集器
	metrics := microservice.NewSimpleMetricsCollector()
	server.AddUnaryInterceptor(microservice.MetricsInterceptor(metrics))

	// 启动服务器
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("gRPC server started on :50051")

	// 保持服务器运行
	select {}
}

// 示例：使用 gRPC 客户端
func useGRPCClient() {
	// 创建注册中心
	registry, err := microservice.NewServiceRegistry(&microservice.RegistryConfig{
		Type: microservice.RegistryTypeMemory,
	})
	if err != nil {
		log.Printf("Failed to create registry: %v", err)
		return
	}

	// 创建负载均衡器
	loadBalancer := microservice.NewRoundRobinLoadBalancer()

	// 创建服务发现
	discovery := microservice.NewServiceDiscovery(registry, loadBalancer)

	// 创建 gRPC 客户端
	client := microservice.NewGRPCServiceClient(
		discovery,
		microservice.WithGRPCTimeout(30*time.Second),
		microservice.WithGRPCRetry(3, time.Second),
	)

	// 调用 gRPC 服务
	ctx := context.Background()

	// 示例请求和响应结构（实际使用时需要根据 protobuf 定义）
	request := map[string]interface{}{
		"id": "1",
	}
	response := map[string]interface{}{}

	err = client.CallGRPC(ctx, "user-service", "/user.UserService/GetUser", request, response, nil)
	if err != nil {
		log.Printf("gRPC call failed: %v", err)
		return
	}

	fmt.Printf("Response: %+v\n", response)
}

// 示例：使用健康检查
func useHealthCheck() {
	// 创建健康检查配置
	config := microservice.NewHealthConfig()
	config.Interval = 10 * time.Second
	config.Timeout = 5 * time.Second

	// 创建健康检查服务
	healthService := microservice.NewGRPCHealthService(config)

	// 注册健康检查
	healthService.RegisterHealthCheck("database", func(ctx context.Context) error {
		// 模拟数据库健康检查
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	healthService.RegisterHealthCheck("redis", func(ctx context.Context) error {
		// 模拟 Redis 健康检查
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	// 启动健康检查
	healthService.Start()

	// 获取健康状态
	statuses := healthService.GetAllStatus()
	for service, status := range statuses {
		fmt.Printf("Service: %s, Status: %s\n", service, status.Status)
	}
}

// 示例：使用流式通信
func useStreaming() {
	// 创建流指标收集器
	metrics := microservice.NewStreamMetricsCollector()

	// 记录流开始
	metrics.RecordStreamStart("/user.UserService/Chat")

	// 模拟流处理
	go func() {
		time.Sleep(5 * time.Second)
		metrics.RecordStreamEnd("/user.UserService/Chat", 5*time.Second)
	}()

	// 获取流指标
	streamMetrics := metrics.GetMetrics()
	for method, metric := range streamMetrics {
		fmt.Printf("Method: %s, Active Streams: %d, Total Streams: %d\n",
			method, metric.ActiveStreams, metric.TotalStreams)
	}
}

// 示例：使用拦截器
func useInterceptors() {
	// 创建认证拦截器
	validTokens := map[string]bool{
		"token1": true,
		"token2": true,
	}
	authInterceptor := microservice.TokenAuthInterceptor(validTokens)

	// 创建超时拦截器
	timeoutInterceptor := microservice.TimeoutInterceptor(10 * time.Second)

	// 创建验证拦截器
	validationInterceptor := microservice.ValidationInterceptor(func(req interface{}) error {
		// 实现请求验证逻辑
		return nil
	})

	// 在实际的 gRPC 服务器中使用这些拦截器
	fmt.Printf("Created interceptors: auth=%v, timeout=%v, validation=%v\n",
		authInterceptor != nil, timeoutInterceptor != nil, validationInterceptor != nil)
}

// 示例：使用健康监控
func useHealthMonitoring() {
	// 创建健康检查服务
	healthService := microservice.NewGRPCHealthService(nil)

	// 创建健康监控器
	monitor := microservice.NewHealthMonitor(healthService)

	// 启动监控
	monitor.Start()

	// 模拟健康检查
	healthService.RegisterHealthCheck("web-service", func(ctx context.Context) error {
		// 模拟健康检查
		return nil
	})

	// 获取监控指标
	metrics := monitor.GetMetrics()
	fmt.Printf("Health Metrics: Total=%d, Successful=%d, Failed=%d\n",
		metrics.TotalChecks, metrics.SuccessfulChecks, metrics.FailedChecks)
}

// 主函数
func main() {
	fmt.Println("=== gRPC 示例 ===")

	// 示例 1: 使用拦截器
	fmt.Println("\n1. 使用拦截器")
	useInterceptors()

	// 示例 2: 使用健康检查
	fmt.Println("\n2. 使用健康检查")
	useHealthCheck()

	// 示例 3: 使用流式通信
	fmt.Println("\n3. 使用流式通信")
	useStreaming()

	// 示例 4: 使用健康监控
	fmt.Println("\n4. 使用健康监控")
	useHealthMonitoring()

	// 示例 5: 使用 gRPC 客户端
	fmt.Println("\n5. 使用 gRPC 客户端")
	useGRPCClient()

	fmt.Println("\n=== 示例完成 ===")

	// 注意：实际使用时，startGRPCServer() 应该在单独的 goroutine 中运行
	// 这里只是展示代码结构
	fmt.Println("\n要启动 gRPC 服务器，请调用 startGRPCServer()")
}
