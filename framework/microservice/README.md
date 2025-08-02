# Laravel-Go 微服务系统

基于 Laravel-Go 框架的微服务支持系统，提供完整的服务发现、注册、负载均衡和服务间通信功能。

## 功能特性

### 1. 服务注册与发现

- **多种注册中心支持**: 
  - 内存注册中心 (开发测试用)
  - Etcd 注册中心 (生产环境推荐)
  - Consul 注册中心 (企业级服务发现)
  - Nacos 注册中心 (阿里云开源)
  - Zookeeper 注册中心 (传统分布式协调)
- **服务发现**: 自动发现服务实例，支持缓存和实时更新
- **健康检查**: 内置健康检查机制，自动过滤不健康服务
- **TTL 机制**: 支持服务过期自动清理

### 2. 负载均衡

- **轮询负载均衡器**: 按顺序分配请求到健康服务实例
- **随机负载均衡器**: 随机选择健康服务实例
- **健康过滤**: 自动过滤不健康的服务实例

### 3. 服务通信

- **HTTP 客户端**: 支持 HTTP/HTTPS 服务间通信
- **重试机制**: 内置重试和超时机制
- **熔断器**: 防止级联故障的熔断器模式
- **JSON 支持**: 自动序列化和反序列化 JSON 数据

### 4. 监控和事件

- **事件监听**: 支持服务变化事件监听
- **缓存统计**: 提供缓存使用情况统计
- **并发安全**: 所有操作都是线程安全的

## 快速开始

### 1. 创建服务注册中心

#### 内存注册中心 (开发测试)

```go
package main

import (
    "context"
    "time"
    "laravel-go/framework/microservice"
)

func main() {
    // 创建内存服务注册中心
    registry := microservice.NewMemoryServiceRegistry()

    // 启动清理工作协程（可选）
    registry.StartCleanupWorker(10 * time.Second)

    ctx := context.Background()

    // 注册服务
    service := &microservice.ServiceInfo{
        ID:       "user-service-1",
        Name:     "user-service",
        Version:  "1.0.0",
        Address:  "localhost",
        Port:     8080,
        Protocol: "http",
        Health:   "healthy",
        Metadata: map[string]string{
            "environment": "production",
            "region":      "us-west-1",
        },
        Tags: []string{"api", "user"},
        TTL:   30 * time.Second,
    }

    err := registry.Register(ctx, service)
    if err != nil {
        panic(err)
    }
}
```

#### Etcd 注册中心 (生产环境)

```go
// 使用配置模式
config := &microservice.RegistryConfig{
    Type: microservice.RegistryTypeEtcd,
    Etcd: &microservice.EtcdConfig{
        Endpoints: []string{"localhost:2379", "localhost:2380"},
        Username:  "admin",
        Password:  "password",
        Prefix:    "/laravel-go/services",
        TTL:       30 * time.Second,
    },
}

registry, err := microservice.NewServiceRegistry(config)
if err != nil {
    panic(err)
}

// 使用构建器模式
registry, err := microservice.NewRegistryBuilder().
    WithType(microservice.RegistryTypeEtcd).
    WithEtcd([]string{"localhost:2379"}, "/laravel-go/services").
    Build()
```

#### Consul 注册中心

```go
// 使用配置模式
config := &microservice.RegistryConfig{
    Type: microservice.RegistryTypeConsul,
    Consul: &microservice.ConsulConfig{
        Address:    "localhost:8500",
        Token:      "consul-token",
        Datacenter: "dc1",
        Prefix:     "laravel-go/services",
        TTL:        30 * time.Second,
    },
}

registry, err := microservice.NewServiceRegistry(config)

// 使用构建器模式
registry, err := microservice.NewRegistryBuilder().
    WithType(microservice.RegistryTypeConsul).
    WithConsul("localhost:8500", "laravel-go/services").
    Build()
```

#### Nacos 注册中心

```go
// 使用配置模式
config := &microservice.RegistryConfig{
    Type: microservice.RegistryTypeNacos,
    Nacos: &microservice.NacosConfig{
        ServerAddr: "localhost:8848",
        Namespace:  "public",
        Group:      "DEFAULT_GROUP",
        Username:   "nacos",
        Password:   "nacos",
        TTL:        30 * time.Second,
    },
}

registry, err := microservice.NewServiceRegistry(config)

// 使用构建器模式
registry, err := microservice.NewRegistryBuilder().
    WithType(microservice.RegistryTypeNacos).
    WithNacos("localhost:8848", "public", "DEFAULT_GROUP").
    Build()
```

#### Zookeeper 注册中心

```go
// 使用配置模式
config := &microservice.RegistryConfig{
    Type: microservice.RegistryTypeZookeeper,
    Zookeeper: &microservice.ZookeeperConfig{
        Servers:        []string{"localhost:2181", "localhost:2182"},
        Prefix:         "/laravel-go/services",
        TTL:            30 * time.Second,
        SessionTimeout: 10 * time.Second,
    },
}

registry, err := microservice.NewServiceRegistry(config)

// 使用构建器模式
registry, err := microservice.NewRegistryBuilder().
    WithType(microservice.RegistryTypeZookeeper).
    WithZookeeper([]string{"localhost:2181"}, "/laravel-go/services").
    Build()
```

### 2. 服务发现

```go
// 创建服务发现
discovery := microservice.NewMemoryServiceDiscovery(registry, microservice.NewRoundRobinLoadBalancer())

// 发现所有服务实例
services, err := discovery.Discover(ctx, "user-service")
if err != nil {
    panic(err)
}

// 发现单个服务实例（负载均衡）
service, err := discovery.DiscoverOne(ctx, "user-service")
if err != nil {
    panic(err)
}

fmt.Printf("选择的服务: %s:%d\n", service.Address, service.Port)
```

### 3. 服务间通信

```go
// 创建服务客户端
client := microservice.NewServiceClient(
    discovery,
    microservice.WithTimeout(5*time.Second),
    microservice.WithRetry(3, 1*time.Second),
)

// GET 请求
response, err := client.Get(ctx, "user-service", "/users")
if err != nil {
    panic(err)
}

// POST 请求
user := User{Name: "张三", Email: "zhangsan@example.com"}
response, err = client.Post(ctx, "user-service", "/users", user)
if err != nil {
    panic(err)
}

// JSON 请求
var responseUser User
err = client.PostJSON(ctx, "user-service", "/users", user, &responseUser)
if err != nil {
    panic(err)
}
```

### 4. 熔断器

```go
// 创建熔断器
cb := microservice.NewSimpleCircuitBreaker(3, 5*time.Second)

// 使用熔断器执行操作
err := cb.Execute(ctx, func() error {
    // 执行可能失败的操作
    return someOperation()
})

if err != nil {
    fmt.Printf("操作失败: %v\n", err)
}

// 检查熔断器状态
if cb.IsOpen() {
    fmt.Println("熔断器已开启")
}
```

## API 参考

### ServiceInfo

服务信息结构体，包含服务的所有元数据。

```go
type ServiceInfo struct {
    ID        string            `json:"id"`         // 服务唯一标识
    Name      string            `json:"name"`       // 服务名称
    Version   string            `json:"version"`    // 服务版本
    Address   string            `json:"address"`    // 服务地址
    Port      int               `json:"port"`       // 服务端口
    Protocol  string            `json:"protocol"`   // 协议 (http, grpc, tcp)
    Health    string            `json:"health"`     // 健康状态 (healthy, unhealthy, unknown)
    Metadata  map[string]string `json:"metadata"`   // 元数据
    Tags      []string          `json:"tags"`       // 标签
    CreatedAt time.Time         `json:"created_at"` // 创建时间
    UpdatedAt time.Time         `json:"updated_at"` // 更新时间
    LastCheck time.Time         `json:"last_check"` // 最后检查时间
    TTL       time.Duration     `json:"ttl"`        // 生存时间
}
```

### ServiceRegistry

服务注册中心接口。

```go
type ServiceRegistry interface {
    Register(ctx context.Context, service *ServiceInfo) error
    Deregister(ctx context.Context, serviceID string) error
    Update(ctx context.Context, service *ServiceInfo) error
    GetService(ctx context.Context, serviceID string) (*ServiceInfo, error)
    ListServices(ctx context.Context) ([]*ServiceInfo, error)
    Watch(ctx context.Context) (<-chan ServiceEvent, error)
    Close() error
}
```

### ServiceDiscovery

服务发现接口。

```go
type ServiceDiscovery interface {
    Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error)
    DiscoverOne(ctx context.Context, serviceName string) (*ServiceInfo, error)
    Watch(ctx context.Context, serviceName string) (<-chan ServiceEvent, error)
    Close() error
}
```

### ServiceClient

服务通信客户端。

```go
type ServiceClient struct {
    discovery  ServiceDiscovery
    httpClient *http.Client
    timeout    time.Duration
    retryCount int
    retryDelay time.Duration
}

// 主要方法
func (c *ServiceClient) Call(ctx context.Context, serviceName, method, path string, data interface{}) ([]byte, error)
func (c *ServiceClient) CallJSON(ctx context.Context, serviceName, method, path string, requestData, responseData interface{}) error
func (c *ServiceClient) Get(ctx context.Context, serviceName, path string) ([]byte, error)
func (c *ServiceClient) Post(ctx context.Context, serviceName, path string, data interface{}) ([]byte, error)
func (c *ServiceClient) Put(ctx context.Context, serviceName, path string, data interface{}) ([]byte, error)
func (c *ServiceClient) Delete(ctx context.Context, serviceName, path string) ([]byte, error)
```

## 负载均衡器

### RoundRobinLoadBalancer

轮询负载均衡器，按顺序分配请求。

```go
rr := microservice.NewRoundRobinLoadBalancer()
selected := rr.Select(services)
```

### RandomLoadBalancer

随机负载均衡器，随机选择服务实例。

```go
random := microservice.NewRandomLoadBalancer()
selected := random.Select(services)
```

## 健康检查

### HTTPHealthChecker

HTTP 健康检查器。

```go
checker := microservice.NewHTTPHealthChecker(5 * time.Second)
err := checker.Check(ctx, service)
isHealthy := checker.IsHealthy(service)
```

## 事件系统

### ServiceEvent

服务事件结构体。

```go
type ServiceEvent struct {
    Type    ServiceEventType `json:"type"`    // 事件类型
    Service *ServiceInfo     `json:"service"` // 服务信息
}

type ServiceEventType string

const (
    ServiceEventCreated ServiceEventType = "created"
    ServiceEventUpdated ServiceEventType = "updated"
    ServiceEventDeleted ServiceEventType = "deleted"
)
```

### 监听服务变化

```go
// 监听所有服务变化
events, err := registry.Watch(ctx)
if err != nil {
    panic(err)
}

for event := range events {
    switch event.Type {
    case microservice.ServiceEventCreated:
        fmt.Printf("服务创建: %s\n", event.Service.ID)
    case microservice.ServiceEventUpdated:
        fmt.Printf("服务更新: %s\n", event.Service.ID)
    case microservice.ServiceEventDeleted:
        fmt.Printf("服务删除: %s\n", event.Service.ID)
    }
}

// 监听特定服务变化
events, err := discovery.Watch(ctx, "user-service")
if err != nil {
    panic(err)
}

for event := range events {
    fmt.Printf("用户服务变化: %s\n", event.Type)
}
```

## 配置选项

### ServiceClient 选项

```go
// 设置超时时间
client := microservice.NewServiceClient(
    discovery,
    microservice.WithTimeout(10*time.Second),
)

// 设置重试参数
client := microservice.NewServiceClient(
    discovery,
    microservice.WithRetry(5, 2*time.Second),
)

// 设置自定义 HTTP 客户端
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        IdleConnTimeout:     90 * time.Second,
        DisableCompression:  true,
    },
}

client := microservice.NewServiceClient(
    discovery,
    microservice.WithHTTPClient(httpClient),
)
```

## 最佳实践

### 1. 服务注册

- 为每个服务设置唯一的 ID
- 合理设置 TTL 值，避免服务过期
- 定期更新服务健康状态
- 在服务关闭时主动注销

### 2. 负载均衡

- 根据业务需求选择合适的负载均衡策略
- 监控服务健康状态
- 实现自定义负载均衡器

### 3. 错误处理

- 使用熔断器防止级联故障
- 合理设置重试次数和间隔
- 监控服务调用成功率

### 4. 监控和日志

- 监听服务变化事件
- 记录服务调用日志
- 监控缓存命中率

## 示例程序

完整的示例程序请参考 `examples/microservice_demo/main.go`，包含：

1. 服务注册和发现演示
2. 负载均衡演示
3. 服务通信演示
4. 熔断器演示
5. 微服务演示服务器

运行示例：

```bash
cd examples/microservice_demo
go run main.go
```

## 扩展性

微服务系统设计为可扩展的，支持：

1. **自定义注册中心**: 实现 `ServiceRegistry` 接口
2. **自定义发现服务**: 实现 `ServiceDiscovery` 接口
3. **自定义负载均衡器**: 实现 `LoadBalancer` 接口
4. **自定义健康检查器**: 实现 `HealthChecker` 接口
5. **自定义熔断器**: 实现 `CircuitBreaker` 接口

## 性能考虑

1. **缓存**: 服务发现结果会被缓存，减少注册中心访问
2. **并发**: 所有操作都是线程安全的
3. **内存**: 内存注册中心适合中小规模部署
4. **网络**: 支持连接池和超时配置

## 故障排除

### 常见问题

1. **服务发现失败**: 检查服务是否已注册，健康状态是否正确
2. **通信超时**: 调整超时时间和重试参数
3. **熔断器开启**: 检查下游服务是否正常，考虑重置熔断器
4. **缓存不一致**: 清除缓存或等待自动更新

### 调试技巧

1. 启用详细日志
2. 监控服务变化事件
3. 检查缓存统计信息
4. 验证服务健康状态
