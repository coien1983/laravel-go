# Laravel-Go gRPC 扩展功能

## 📋 概述

Laravel-Go Framework 的 gRPC 扩展功能提供了完整的 gRPC 微服务支持，包括服务器、客户端、拦截器、流式通信、健康检查等功能。这些功能帮助开发者构建高性能、可扩展的 gRPC 微服务应用。

## 🏗️ 核心组件

### 1. gRPC 服务器 (GRPCServer)

gRPC 服务器提供了完整的服务端功能，支持服务注册、健康检查、拦截器等。

#### 基本用法

```go
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

// 启动服务器
if err := server.Start(); err != nil {
    log.Fatal(err)
}

// 停止服务器
defer server.Stop()
```

#### 配置选项

- `WithGRPCAddress(address)`: 设置服务器地址
- `WithGRPCPort(port)`: 设置服务器端口
- `WithGRPCTLS(certFile, keyFile)`: 启用 TLS
- `WithGRPCRegistry(registry)`: 设置服务注册中心
- `WithGRPCServiceInfo(name, version)`: 设置服务信息
- `WithGRPCMetadata(metadata)`: 设置服务元数据
- `WithGRPCHealthCheck(enabled, path)`: 配置健康检查
- `WithGRPCReflection(enabled)`: 启用反射
- `WithGRPCLogging(enabled)`: 启用日志

### 2. gRPC 客户端 (GRPCServiceClient)

gRPC 客户端提供了服务发现、连接管理、重试机制等功能。

#### 基本用法

```go
// 创建服务发现
registry, err := microservice.NewServiceRegistry(&microservice.RegistryConfig{
    Type: microservice.RegistryTypeMemory,
})
if err != nil {
    log.Fatal(err)
}

loadBalancer := microservice.NewRoundRobinLoadBalancer()
discovery := microservice.NewServiceDiscovery(registry, loadBalancer)

// 创建 gRPC 客户端
client := microservice.NewGRPCServiceClient(
    discovery,
    microservice.WithGRPCTimeout(30*time.Second),
    microservice.WithGRPCRetry(3, time.Second),
)

// 调用 gRPC 服务
ctx := context.Background()
request := map[string]interface{}{"id": "1"}
response := map[string]interface{}{}

err = client.CallGRPC(ctx, "user-service", "/user.UserService/GetUser", request, response, nil)
if err != nil {
    log.Printf("gRPC call failed: %v", err)
    return
}
```

#### 配置选项

- `WithGRPCTimeout(timeout)`: 设置超时时间
- `WithGRPCRetry(count, delay)`: 设置重试参数

### 3. gRPC 拦截器 (Interceptors)

gRPC 拦截器提供了中间件功能，包括日志、认证、限流、追踪等。

#### 内置拦截器

```go
// 日志拦截器
server.AddUnaryInterceptor(microservice.LoggingInterceptor())

// 认证拦截器
validTokens := map[string]bool{"token1": true, "token2": true}
authInterceptor := microservice.TokenAuthInterceptor(validTokens)
server.AddUnaryInterceptor(authInterceptor)

// 限流拦截器
limiter := microservice.NewSimpleRateLimiter(map[string]int{
    "/user.UserService/GetUser": 100,
})
server.AddUnaryInterceptor(microservice.RateLimitInterceptor(limiter))

// 熔断器拦截器
circuitBreaker := microservice.NewCircuitBreaker(microservice.CircuitBreakerConfig{
    Threshold: 5,
    Timeout:   30 * time.Second,
})
server.AddUnaryInterceptor(microservice.CircuitBreakerInterceptor(circuitBreaker))

// 指标拦截器
metrics := microservice.NewSimpleMetricsCollector()
server.AddUnaryInterceptor(microservice.MetricsInterceptor(metrics))

// 追踪拦截器
tracer := microservice.NewTracer(microservice.TracerConfig{
    ServiceName: "user-service",
})
server.AddUnaryInterceptor(microservice.TracingInterceptor(tracer))

// 验证拦截器
validationInterceptor := microservice.ValidationInterceptor(func(req interface{}) error {
    // 实现请求验证逻辑
    return nil
})
server.AddUnaryInterceptor(validationInterceptor)

// 超时拦截器
server.AddUnaryInterceptor(microservice.TimeoutInterceptor(10 * time.Second))

// 恢复拦截器
server.AddUnaryInterceptor(microservice.RecoveryInterceptor())

// 元数据拦截器
server.AddUnaryInterceptor(microservice.MetadataInterceptor())
```

#### 流拦截器

```go
// 流日志拦截器
server.AddStreamInterceptor(microservice.StreamLoggingInterceptor())

// 流认证拦截器
server.AddStreamInterceptor(microservice.StreamAuthInterceptor(authFunc))

// 流限流拦截器
server.AddStreamInterceptor(microservice.StreamRateLimitInterceptor(limiter))

// 流恢复拦截器
server.AddStreamInterceptor(microservice.StreamRecoveryInterceptor())
```

### 4. 流式通信 (Streaming)

gRPC 流式通信支持客户端流、服务器流和双向流。

#### 流管理器

```go
// 创建流管理器
streamManager := microservice.NewStreamManager()

// 注册流
streamManager.RegisterStream("stream1", microservice.StreamTypeBidirectional, ctx)

// 获取流信息
stream, exists := streamManager.GetStream("stream1")

// 列出所有流
streams := streamManager.ListStreams()

// 关闭所有流
streamManager.CloseAllStreams()
```

#### 流指标收集

```go
// 创建流指标收集器
metrics := microservice.NewStreamMetricsCollector()

// 记录流开始
metrics.RecordStreamStart("/user.UserService/Chat")

// 记录消息发送
metrics.RecordMessageSent("/user.UserService/Chat")

// 记录消息接收
metrics.RecordMessageReceived("/user.UserService/Chat")

// 记录错误
metrics.RecordError("/user.UserService/Chat")

// 记录流结束
metrics.RecordStreamEnd("/user.UserService/Chat", duration)

// 获取指标
streamMetrics := metrics.GetMetrics()
```

### 5. 健康检查 (Health Check)

gRPC 健康检查提供了服务健康状态监控功能。

#### 健康检查服务

```go
// 创建健康检查配置
config := microservice.NewHealthConfig()
config.Interval = 10 * time.Second
config.Timeout = 5 * time.Second

// 创建健康检查服务
healthService := microservice.NewGRPCHealthService(config)

// 注册健康检查
healthService.RegisterHealthCheck("database", func(ctx context.Context) error {
    // 实现数据库健康检查
    return nil
})

healthService.RegisterHealthCheck("redis", func(ctx context.Context) error {
    // 实现 Redis 健康检查
    return nil
})

// 启动健康检查
healthService.Start()

// 获取健康状态
statuses := healthService.GetAllStatus()
for service, status := range statuses {
    fmt.Printf("Service: %s, Status: %s\n", service, status.Status)
}
```

#### 内置健康检查

```go
// 数据库健康检查
healthService.RegisterHealthCheck("database", microservice.DatabaseHealthCheck(db))

// Redis 健康检查
healthService.RegisterHealthCheck("redis", microservice.RedisHealthCheck(redis))

// HTTP 健康检查
healthService.RegisterHealthCheck("web", microservice.HTTPHealthCheck("http://localhost:8080/health", 5*time.Second))

// gRPC 健康检查
healthService.RegisterHealthCheck("grpc", microservice.GRPCHealthCheck("localhost:50051", 5*time.Second))

// 文件系统健康检查
healthService.RegisterHealthCheck("filesystem", microservice.FileSystemHealthCheck("/data"))

// 内存健康检查
healthService.RegisterHealthCheck("memory", microservice.MemoryHealthCheck(1024*1024*1024)) // 1GB

// CPU 健康检查
healthService.RegisterHealthCheck("cpu", microservice.CPUHealthCheck(80.0)) // 80%
```

#### 健康监控

```go
// 创建健康监控器
monitor := microservice.NewHealthMonitor(healthService)

// 启动监控
monitor.Start()

// 获取监控指标
metrics := monitor.GetMetrics()
fmt.Printf("Health Metrics: Total=%d, Successful=%d, Failed=%d\n",
    metrics.TotalChecks, metrics.SuccessfulChecks, metrics.FailedChecks)

// 获取告警
alerts := monitor.GetAlerts()
for _, alert := range alerts {
    fmt.Printf("Alert: Service=%s, Status=%s, Message=%s\n",
        alert.Service, alert.Status, alert.Message)
}
```

## 🔧 高级功能

### 1. TLS 支持

```go
// 启用 TLS
server := microservice.NewGRPCServer(
    microservice.WithGRPCTLS("cert.pem", "key.pem"),
)
```

### 2. 服务注册与发现

```go
// 创建注册中心
registry, err := microservice.NewServiceRegistry(&microservice.RegistryConfig{
    Type: microservice.RegistryTypeConsul,
    Consul: &microservice.ConsulConfig{
        Address: "localhost:8500",
        Prefix:  "laravel-go/services",
    },
})

// 创建 gRPC 服务器并注册服务
server := microservice.NewGRPCServer(
    microservice.WithGRPCRegistry(registry),
    microservice.WithGRPCServiceInfo("user-service", "1.0.0"),
)
```

### 3. 负载均衡

```go
// 轮询负载均衡
loadBalancer := microservice.NewRoundRobinLoadBalancer()

// 随机负载均衡
loadBalancer := microservice.NewRandomLoadBalancer()

// 创建服务发现
discovery := microservice.NewServiceDiscovery(registry, loadBalancer)
```

### 4. 熔断器

```go
// 创建熔断器
circuitBreaker := microservice.NewCircuitBreaker(microservice.CircuitBreakerConfig{
    Threshold: 5,
    Timeout:   30 * time.Second,
    HalfOpen:  true,
})

// 在拦截器中使用
server.AddUnaryInterceptor(microservice.CircuitBreakerInterceptor(circuitBreaker))
```

### 5. 分布式追踪

```go
// 创建追踪器
tracer := microservice.NewTracer(microservice.TracerConfig{
    ServiceName: "user-service",
    Sampler:     0.1,
    Reporter: microservice.ReporterConfig{
        Type: "jaeger",
        URL:  "http://localhost:14268/api/traces",
    },
})

// 在拦截器中使用
server.AddUnaryInterceptor(microservice.TracingInterceptor(tracer))
```

## 📊 监控和指标

### 1. 指标收集

```go
// 创建指标收集器
metrics := microservice.NewSimpleMetricsCollector()

// 在拦截器中使用
server.AddUnaryInterceptor(microservice.MetricsInterceptor(metrics))

// 获取指标
allMetrics := metrics.GetMetrics()
for key, value := range allMetrics {
    fmt.Printf("Metric: %s = %f\n", key, value)
}
```

### 2. 流指标

```go
// 创建流指标收集器
streamMetrics := microservice.NewStreamMetricsCollector()

// 记录流指标
streamMetrics.RecordStreamStart("/user.UserService/Chat")
streamMetrics.RecordMessageSent("/user.UserService/Chat")
streamMetrics.RecordMessageReceived("/user.UserService/Chat")
streamMetrics.RecordError("/user.UserService/Chat")
streamMetrics.RecordStreamEnd("/user.UserService/Chat", duration)

// 获取流指标
metrics := streamMetrics.GetMetrics()
for method, metric := range metrics {
    fmt.Printf("Method: %s, Active Streams: %d, Total Streams: %d\n",
        method, metric.ActiveStreams, metric.TotalStreams)
}
```

## 🚀 性能优化

### 1. 连接池管理

gRPC 客户端自动管理连接池，复用连接以提高性能。

### 2. 缓存策略

服务发现支持本地缓存，减少网络请求。

### 3. 批量处理

支持批量 gRPC 调用，提高吞吐量。

### 4. 异步处理

支持异步 gRPC 调用，提高并发性能。

## 🧪 测试

### 1. 单元测试

```go
func TestGRPCServer(t *testing.T) {
    // 创建测试服务器
    server := microservice.NewGRPCServer(
        microservice.WithGRPCPort(0), // 随机端口
    )

    // 启动服务器
    err := server.Start()
    if err != nil {
        t.Fatal(err)
    }
    defer server.Stop()

    // 测试服务器功能
    if !server.IsRunning() {
        t.Error("Server should be running")
    }
}
```

### 2. 集成测试

```go
func TestGRPCClientServer(t *testing.T) {
    // 创建注册中心
    registry, _ := microservice.NewServiceRegistry(&microservice.RegistryConfig{
        Type: microservice.RegistryTypeMemory,
    })

    // 创建服务器
    server := microservice.NewGRPCServer(
        microservice.WithGRPCRegistry(registry),
        microservice.WithGRPCServiceInfo("test-service", "1.0.0"),
    )

    // 启动服务器
    server.Start()
    defer server.Stop()

    // 创建客户端
    loadBalancer := microservice.NewRoundRobinLoadBalancer()
    discovery := microservice.NewServiceDiscovery(registry, loadBalancer)
    client := microservice.NewGRPCServiceClient(discovery)

    // 测试调用
    ctx := context.Background()
    request := map[string]interface{}{"test": "data"}
    response := map[string]interface{}{}

    err := client.CallGRPC(ctx, "test-service", "/test.TestService/Test", request, response, nil)
    if err != nil {
        t.Errorf("gRPC call failed: %v", err)
    }
}
```

## 📝 最佳实践

### 1. 错误处理

```go
// 使用统一的错误处理
err := client.CallGRPC(ctx, "user-service", "/user.UserService/GetUser", request, response, nil)
if err != nil {
    // 检查错误类型
    if grpcErr, ok := err.(*microservice.StreamError); ok {
        switch grpcErr.Code {
        case codes.NotFound:
            // 处理未找到错误
        case codes.Unauthenticated:
            // 处理认证错误
        case codes.ResourceExhausted:
            // 处理限流错误
        default:
            // 处理其他错误
        }
    }
}
```

### 2. 超时设置

```go
// 设置合理的超时时间
client := microservice.NewGRPCServiceClient(
    discovery,
    microservice.WithGRPCTimeout(30*time.Second),
    microservice.WithGRPCRetry(3, time.Second),
)
```

### 3. 健康检查

```go
// 定期检查服务健康状态
healthService.RegisterHealthCheck("critical-service", func(ctx context.Context) error {
    // 实现关键服务的健康检查
    return nil
})

// 监控健康状态变化
monitor := microservice.NewHealthMonitor(healthService)
monitor.Start()
```

### 4. 日志记录

```go
// 启用详细的日志记录
server.AddUnaryInterceptor(microservice.LoggingInterceptor())
server.AddStreamInterceptor(microservice.StreamLoggingInterceptor())
```

### 5. 指标监控

```go
// 收集关键指标
metrics := microservice.NewSimpleMetricsCollector()
server.AddUnaryInterceptor(microservice.MetricsInterceptor(metrics))

// 定期导出指标
go func() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        allMetrics := metrics.GetMetrics()
        // 导出指标到监控系统
    }
}()
```

## 🔍 调试和故障排除

### 1. 启用反射

```go
// 启用 gRPC 反射，便于调试
server := microservice.NewGRPCServer(
    microservice.WithGRPCReflection(true),
)
```

### 2. 详细日志

```go
// 启用详细日志
server.AddUnaryInterceptor(microservice.LoggingInterceptor())
```

### 3. 健康检查

```go
// 使用健康检查监控服务状态
healthService := microservice.NewGRPCHealthService(nil)
healthService.RegisterHealthCheck("service", healthCheckFunc)
```

### 4. 指标监控

```go
// 使用指标监控性能
metrics := microservice.NewSimpleMetricsCollector()
server.AddUnaryInterceptor(microservice.MetricsInterceptor(metrics))
```

## 🚀 总结

Laravel-Go Framework 的 gRPC 扩展功能提供了：

1. **完整的 gRPC 支持**: 服务器、客户端、拦截器、流式通信
2. **丰富的拦截器**: 日志、认证、限流、熔断器、追踪等
3. **健康检查**: 内置多种健康检查，支持监控和告警
4. **服务治理**: 服务注册发现、负载均衡、熔断器等
5. **性能优化**: 连接池、缓存、批量处理等
6. **监控指标**: 详细的性能指标和监控支持
7. **最佳实践**: 完整的错误处理、超时设置、日志记录等

通过这些功能，开发者可以构建高性能、高可用、易维护的 gRPC 微服务应用。
