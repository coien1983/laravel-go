# 微服务系统 API 参考

## 📋 概述

Laravel-Go Framework 的微服务系统提供了完整的微服务架构支持，包括服务注册与发现、服务间通信、负载均衡、熔断器、配置管理、分布式追踪等功能。微服务系统帮助开发者构建可扩展、高可用的分布式应用程序。

## 🏗️ 核心概念

### 服务注册中心 (Service Registry)

- 服务注册和注销
- 服务发现和健康检查
- 服务元数据管理

### 服务通信 (Service Communication)

- HTTP/RPC 客户端
- 消息队列集成
- 负载均衡策略

### 服务治理 (Service Governance)

- 熔断器和重试机制
- 限流和降级
- 分布式配置管理

## 🔧 基础用法

### 1. 基本微服务配置

```go
// 创建微服务客户端
client := microservice.NewClient(microservice.ClientConfig{
    Registry: "consul",
    RegistryAddress: "localhost:8500",
    LoadBalancer: "round_robin",
    Timeout: time.Second * 30,
})

// 注册服务
service := microservice.NewService(microservice.ServiceConfig{
    Name: "user-service",
    Version: "1.0.0",
    Port: 8080,
    HealthCheck: "/health",
    Metadata: map[string]string{
        "environment": "production",
        "region": "us-west-1",
    },
})

// 启动服务
service.Start()
defer service.Stop()
```

### 2. 服务间通信

```go
// 创建服务客户端
userClient := microservice.NewClient(microservice.ClientConfig{
    ServiceName: "user-service",
    Registry: "consul",
})

// 调用远程服务
func (c *OrderController) CreateOrder(request http.Request) http.Response {
    // 获取用户信息
    userResponse, err := userClient.Call("GET", "/users/123", nil, map[string]interface{}{
        "timeout": time.Second * 5,
    })
    if err != nil {
        return c.JsonError("Failed to get user info", 500)
    }

    // 处理订单创建
    order := &Models.Order{
        UserID: userResponse.Data["id"].(uint),
        Items:  request.Body["items"].([]interface{}),
    }

    // 保存订单
    err = c.orderService.CreateOrder(order)
    if err != nil {
        return c.JsonError("Failed to create order", 500)
    }

    return c.Json(order).Status(201)
}
```

### 3. 服务发现

```go
// 服务发现
func (c *OrderController) GetUserOrders(userID string, request http.Request) http.Response {
    // 发现用户服务
    userService, err := c.serviceDiscovery.GetService("user-service")
    if err != nil {
        return c.JsonError("User service not available", 503)
    }

    // 调用用户服务
    userResponse, err := userService.Call("GET", fmt.Sprintf("/users/%s", userID), nil)
    if err != nil {
        return c.JsonError("Failed to get user", 500)
    }

    // 获取用户订单
    orders, err := c.orderService.GetOrdersByUserID(userID)
    if err != nil {
        return c.JsonError("Failed to get orders", 500)
    }

    return c.Json(map[string]interface{}{
        "user":   userResponse.Data,
        "orders": orders,
    })
}
```

## 📚 API 参考

### Client 接口

```go
type Client interface {
    Call(method, path string, body interface{}, options map[string]interface{}) (*Response, error)
    CallAsync(method, path string, body interface{}, options map[string]interface{}) (chan *Response, chan error)

    SetTimeout(timeout time.Duration)
    GetTimeout() time.Duration

    SetRetries(retries int)
    GetRetries() int

    SetCircuitBreaker(circuitBreaker CircuitBreaker)
    GetCircuitBreaker() CircuitBreaker

    SetLoadBalancer(loadBalancer LoadBalancer)
    GetLoadBalancer() LoadBalancer

    SetMiddleware(middleware Middleware)
    GetMiddleware() Middleware
}
```

#### 方法说明

- `Call(method, path, body, options)`: 同步调用服务
- `CallAsync(method, path, body, options)`: 异步调用服务
- `SetTimeout(timeout)`: 设置超时时间
- `GetTimeout()`: 获取超时时间
- `SetRetries(retries)`: 设置重试次数
- `GetRetries()`: 获取重试次数
- `SetCircuitBreaker(circuitBreaker)`: 设置熔断器
- `GetCircuitBreaker()`: 获取熔断器
- `SetLoadBalancer(loadBalancer)`: 设置负载均衡器
- `GetLoadBalancer()`: 获取负载均衡器
- `SetMiddleware(middleware)`: 设置中间件
- `GetMiddleware()`: 获取中间件

### Service 接口

```go
type Service interface {
    Start() error
    Stop() error
    IsRunning() bool

    Register() error
    Deregister() error
    IsRegistered() bool

    AddHandler(path string, handler http.HandlerFunc)
    RemoveHandler(path string)
    GetHandlers() map[string]http.HandlerFunc

    SetHealthCheck(healthCheck HealthCheck)
    GetHealthCheck() HealthCheck

    SetMetadata(metadata map[string]string)
    GetMetadata() map[string]string

    SetVersion(version string)
    GetVersion() string
}
```

#### 方法说明

- `Start()`: 启动服务
- `Stop()`: 停止服务
- `IsRunning()`: 检查服务是否运行
- `Register()`: 注册服务
- `Deregister()`: 注销服务
- `IsRegistered()`: 检查服务是否已注册
- `AddHandler(path, handler)`: 添加处理器
- `RemoveHandler(path)`: 移除处理器
- `GetHandlers()`: 获取所有处理器
- `SetHealthCheck(healthCheck)`: 设置健康检查
- `GetHealthCheck()`: 获取健康检查
- `SetMetadata(metadata)`: 设置元数据
- `GetMetadata()`: 获取元数据
- `SetVersion(version)`: 设置版本
- `GetVersion()`: 获取版本

### Registry 接口

```go
type Registry interface {
    Register(service *ServiceInfo) error
    Deregister(service *ServiceInfo) error
    GetService(name string) (*ServiceInfo, error)
    ListServices() ([]*ServiceInfo, error)
    Watch(name string) (chan *ServiceEvent, error)
    StopWatch(name string) error
}
```

#### 方法说明

- `Register(service)`: 注册服务
- `Deregister(service)`: 注销服务
- `GetService(name)`: 获取服务信息
- `ListServices()`: 列出所有服务
- `Watch(name)`: 监听服务变化
- `StopWatch(name)`: 停止监听

## 🎯 高级功能

### 1. 熔断器模式

```go
// 熔断器配置
circuitBreaker := microservice.NewCircuitBreaker(microservice.CircuitBreakerConfig{
    Threshold: 5,           // 失败阈值
    Timeout:   time.Second * 30, // 熔断时间
    HalfOpen:  true,        // 半开状态
})

// 创建带熔断器的客户端
client := microservice.NewClient(microservice.ClientConfig{
    ServiceName: "user-service",
    CircuitBreaker: circuitBreaker,
})

// 使用熔断器
func (c *OrderController) GetUserInfo(userID string) (*Models.User, error) {
    response, err := client.Call("GET", fmt.Sprintf("/users/%s", userID), nil, nil)
    if err != nil {
        // 熔断器会自动处理失败
        return nil, err
    }

    var user Models.User
    err = json.Unmarshal(response.Data, &user)
    return &user, err
}
```

### 2. 负载均衡

```go
// 轮询负载均衡
roundRobinLB := microservice.NewLoadBalancer("round_robin")

// 随机负载均衡
randomLB := microservice.NewLoadBalancer("random")

// 加权负载均衡
weightedLB := microservice.NewLoadBalancer("weighted", map[string]interface{}{
    "weights": map[string]int{
        "instance1": 3,
        "instance2": 2,
        "instance3": 1,
    },
})

// 最少连接负载均衡
leastConnLB := microservice.NewLoadBalancer("least_connections")

// 使用负载均衡器
client := microservice.NewClient(microservice.ClientConfig{
    ServiceName: "user-service",
    LoadBalancer: roundRobinLB,
})
```

### 3. 服务网格

```go
// 服务网格配置
mesh := microservice.NewServiceMesh(microservice.ServiceMeshConfig{
    Sidecar: true,
    Proxy: microservice.ProxyConfig{
        Port: 15001,
        AdminPort: 15000,
    },
    Tracing: microservice.TracingConfig{
        Enabled: true,
        Sampler: 0.1,
    },
})

// 启动服务网格
mesh.Start()
defer mesh.Stop()

// 在服务中使用
func (c *UserController) GetUser(id string, request http.Request) http.Response {
    // 自动添加追踪信息
    span := c.tracer.StartSpan("get_user")
    defer span.Finish()

    // 添加服务网格标签
    span.SetTag("service", "user-service")
    span.SetTag("method", "GET")
    span.SetTag("user_id", id)

    user, err := c.userService.GetUser(id)
    if err != nil {
        span.SetTag("error", true)
        span.LogKV("error", err.Error())
        return c.JsonError("User not found", 404)
    }

    return c.Json(user)
}
```

### 4. 分布式配置

```go
// 配置中心客户端
configClient := microservice.NewConfigClient(microservice.ConfigClientConfig{
    Registry: "consul",
    Prefix: "config",
})

// 获取配置
func (c *UserController) GetConfig() http.Response {
    // 获取数据库配置
    dbConfig, err := configClient.Get("database")
    if err != nil {
        return c.JsonError("Failed to get database config", 500)
    }

    // 获取缓存配置
    cacheConfig, err := configClient.Get("cache")
    if err != nil {
        return c.JsonError("Failed to get cache config", 500)
    }

    return c.Json(map[string]interface{}{
        "database": dbConfig,
        "cache":    cacheConfig,
    })
}

// 监听配置变化
func (c *UserController) WatchConfig() {
    configChan, err := configClient.Watch("database")
    if err != nil {
        log.Printf("Failed to watch config: %v", err)
        return
    }

    for config := range configChan {
        // 处理配置变化
        c.handleConfigChange(config)
    }
}
```

### 5. 分布式追踪

```go
// 追踪配置
tracer := microservice.NewTracer(microservice.TracerConfig{
    ServiceName: "user-service",
    Sampler:     0.1,
    Reporter: microservice.ReporterConfig{
        Type: "jaeger",
        URL:  "http://localhost:14268/api/traces",
    },
})

// 在服务中使用追踪
func (c *UserController) CreateUser(request http.Request) http.Response {
    // 创建根 span
    span := c.tracer.StartSpan("create_user")
    defer span.Finish()

    // 添加请求信息
    span.SetTag("http.method", request.Method)
    span.SetTag("http.url", request.Path)

    // 验证用户数据
    validationSpan := c.tracer.StartSpan("validate_user", span)
    if err := c.validateUserData(request.Body); err != nil {
        validationSpan.SetTag("error", true)
        validationSpan.LogKV("error", err.Error())
        validationSpan.Finish()
        return c.JsonError("Validation failed", 422)
    }
    validationSpan.Finish()

    // 保存用户
    saveSpan := c.tracer.StartSpan("save_user", span)
    user, err := c.userService.CreateUser(request.Body)
    if err != nil {
        saveSpan.SetTag("error", true)
        saveSpan.LogKV("error", err.Error())
        saveSpan.Finish()
        return c.JsonError("Failed to create user", 500)
    }
    saveSpan.SetTag("user_id", user.ID)
    saveSpan.Finish()

    return c.Json(user).Status(201)
}
```

## 🔧 配置选项

### 微服务系统配置

```go
// config/microservice.go
package config

type MicroserviceConfig struct {
    // 服务配置
    Service ServiceConfig `json:"service"`

    // 客户端配置
    Client ClientConfig `json:"client"`

    // 注册中心配置
    Registry RegistryConfig `json:"registry"`

    // 负载均衡配置
    LoadBalancer LoadBalancerConfig `json:"load_balancer"`

    // 熔断器配置
    CircuitBreaker CircuitBreakerConfig `json:"circuit_breaker"`

    // 追踪配置
    Tracing TracingConfig `json:"tracing"`

    // 配置中心配置
    Config ConfigCenterConfig `json:"config"`

    // 服务网格配置
    ServiceMesh ServiceMeshConfig `json:"service_mesh"`
}

type ServiceConfig struct {
    Name     string            `json:"name"`
    Version  string            `json:"version"`
    Port     int               `json:"port"`
    Metadata map[string]string `json:"metadata"`
    HealthCheck HealthCheckConfig `json:"health_check"`
}

type ClientConfig struct {
    Timeout       time.Duration `json:"timeout"`
    Retries       int           `json:"retries"`
    RetryDelay    time.Duration `json:"retry_delay"`
    MaxRetries    int           `json:"max_retries"`
    LoadBalancer  string        `json:"load_balancer"`
    CircuitBreaker bool         `json:"circuit_breaker"`
}

type RegistryConfig struct {
    Driver  string `json:"driver"`
    Address string `json:"address"`
    Timeout time.Duration `json:"timeout"`
    TTL     time.Duration `json:"ttl"`
}

type LoadBalancerConfig struct {
    Strategy string                 `json:"strategy"`
    Options  map[string]interface{} `json:"options"`
}

type CircuitBreakerConfig struct {
    Threshold int           `json:"threshold"`
    Timeout   time.Duration `json:"timeout"`
    HalfOpen  bool          `json:"half_open"`
}

type TracingConfig struct {
    Enabled    bool   `json:"enabled"`
    ServiceName string `json:"service_name"`
    Sampler    float64 `json:"sampler"`
    Reporter   ReporterConfig `json:"reporter"`
}

type ConfigCenterConfig struct {
    Driver string `json:"driver"`
    Address string `json:"address"`
    Prefix  string `json:"prefix"`
    Timeout time.Duration `json:"timeout"`
}

type ServiceMeshConfig struct {
    Enabled bool `json:"enabled"`
    Sidecar bool `json:"sidecar"`
    Proxy   ProxyConfig `json:"proxy"`
}
```

### 配置示例

```go
// config/microservice.go
func GetMicroserviceConfig() *MicroserviceConfig {
    return &MicroserviceConfig{
        Service: ServiceConfig{
            Name:    "user-service",
            Version: "1.0.0",
            Port:    8080,
            Metadata: map[string]string{
                "environment": "production",
                "region":      "us-west-1",
            },
            HealthCheck: HealthCheckConfig{
                Path:     "/health",
                Interval: time.Second * 30,
                Timeout:  time.Second * 5,
            },
        },
        Client: ClientConfig{
            Timeout:        time.Second * 30,
            Retries:        3,
            RetryDelay:     time.Second * 1,
            MaxRetries:     5,
            LoadBalancer:   "round_robin",
            CircuitBreaker: true,
        },
        Registry: RegistryConfig{
            Driver:  "consul",
            Address: "localhost:8500",
            Timeout: time.Second * 10,
            TTL:     time.Second * 30,
        },
        LoadBalancer: LoadBalancerConfig{
            Strategy: "round_robin",
            Options: map[string]interface{}{
                "max_retries": 3,
            },
        },
        CircuitBreaker: CircuitBreakerConfig{
            Threshold: 5,
            Timeout:   time.Second * 30,
            HalfOpen:  true,
        },
        Tracing: TracingConfig{
            Enabled:     true,
            ServiceName: "user-service",
            Sampler:     0.1,
            Reporter: ReporterConfig{
                Type: "jaeger",
                URL:  "http://localhost:14268/api/traces",
            },
        },
        Config: ConfigCenterConfig{
            Driver:  "consul",
            Address: "localhost:8500",
            Prefix:  "config",
            Timeout: time.Second * 10,
        },
        ServiceMesh: ServiceMeshConfig{
            Enabled: true,
            Sidecar: true,
            Proxy: ProxyConfig{
                Port:      15001,
                AdminPort: 15000,
            },
        },
    }
}
```

## 🚀 性能优化

### 1. 连接池管理

```go
// 连接池配置
type ConnectionPool struct {
    maxConnections int
    maxIdle        int
    idleTimeout    time.Duration
    connections    chan *Connection
    mutex          sync.Mutex
}

func NewConnectionPool(maxConn, maxIdle int, idleTimeout time.Duration) *ConnectionPool {
    pool := &ConnectionPool{
        maxConnections: maxConn,
        maxIdle:        maxIdle,
        idleTimeout:    idleTimeout,
        connections:    make(chan *Connection, maxConn),
    }

    // 初始化连接池
    for i := 0; i < maxIdle; i++ {
        conn := pool.createConnection()
        pool.connections <- conn
    }

    return pool
}

func (p *ConnectionPool) GetConnection() (*Connection, error) {
    select {
    case conn := <-p.connections:
        if conn.IsValid() {
            return conn, nil
        }
        // 连接无效，创建新连接
        return p.createConnection(), nil
    case <-time.After(time.Second * 5):
        return nil, errors.New("connection pool timeout")
    }
}

func (p *ConnectionPool) ReturnConnection(conn *Connection) {
    if conn.IsValid() {
        select {
        case p.connections <- conn:
        default:
            // 连接池满了，关闭连接
            conn.Close()
        }
    } else {
        conn.Close()
    }
}
```

### 2. 请求缓存

```go
// 请求缓存
type RequestCache struct {
    cache cache.Cache
    ttl   time.Duration
}

func (rc *RequestCache) GetCachedResponse(key string) (*Response, bool) {
    if cached, exists := rc.cache.Get(key); exists {
        return cached.(*Response), true
    }
    return nil, false
}

func (rc *RequestCache) CacheResponse(key string, response *Response) {
    rc.cache.Set(key, response, rc.ttl)
}

// 在客户端中使用缓存
func (c *CachedClient) Call(method, path string, body interface{}, options map[string]interface{}) (*Response, error) {
    // 生成缓存键
    cacheKey := c.generateCacheKey(method, path, body)

    // 检查缓存
    if cached, exists := c.cache.GetCachedResponse(cacheKey); exists {
        return cached, nil
    }

    // 调用远程服务
    response, err := c.Client.Call(method, path, body, options)
    if err != nil {
        return nil, err
    }

    // 缓存响应
    c.cache.CacheResponse(cacheKey, response)

    return response, nil
}
```

### 3. 批量请求

```go
// 批量请求处理
type BatchClient struct {
    microservice.Client
    batchSize int
    batchTimeout time.Duration
}

func (bc *BatchClient) BatchCall(requests []Request) ([]*Response, error) {
    if len(requests) == 0 {
        return nil, nil
    }

    // 分批处理
    batches := bc.splitIntoBatches(requests)
    results := make([]*Response, 0, len(requests))

    for _, batch := range batches {
        batchResults, err := bc.processBatch(batch)
        if err != nil {
            return nil, err
        }
        results = append(results, batchResults...)
    }

    return results, nil
}

func (bc *BatchClient) processBatch(requests []Request) ([]*Response, error) {
    responses := make([]*Response, len(requests))
    errors := make([]error, len(requests))

    var wg sync.WaitGroup
    for i, req := range requests {
        wg.Add(1)
        go func(index int, request Request) {
            defer wg.Done()

            response, err := bc.Client.Call(request.Method, request.Path, request.Body, request.Options)
            responses[index] = response
            errors[index] = err
        }(i, req)
    }

    wg.Wait()

    // 检查错误
    for _, err := range errors {
        if err != nil {
            return nil, err
        }
    }

    return responses, nil
}
```

## 🧪 测试

### 1. 微服务测试

```go
// tests/microservice_test.go
package tests

import (
    "testing"
    "time"
    "laravel-go/framework/microservice"
)

func TestServiceRegistration(t *testing.T) {
    // 创建注册中心
    registry := microservice.NewRegistry("consul", "localhost:8500")

    // 创建服务
    service := microservice.NewService(microservice.ServiceConfig{
        Name: "test-service",
        Port: 8081,
    })

    // 注册服务
    err := registry.Register(service)
    if err != nil {
        t.Fatal(err)
    }

    // 验证服务注册
    registeredService, err := registry.GetService("test-service")
    if err != nil {
        t.Fatal(err)
    }

    if registeredService.Name != "test-service" {
        t.Error("Service name mismatch")
    }

    // 注销服务
    err = registry.Deregister(service)
    if err != nil {
        t.Fatal(err)
    }
}

func TestServiceDiscovery(t *testing.T) {
    // 创建服务发现
    discovery := microservice.NewServiceDiscovery("consul", "localhost:8500")

    // 发现服务
    services, err := discovery.GetServices("user-service")
    if err != nil {
        t.Fatal(err)
    }

    if len(services) == 0 {
        t.Error("No services found")
    }

    // 验证服务信息
    service := services[0]
    if service.Name != "user-service" {
        t.Error("Service name mismatch")
    }
}

func TestServiceCommunication(t *testing.T) {
    // 创建客户端
    client := microservice.NewClient(microservice.ClientConfig{
        ServiceName: "user-service",
        Timeout:     time.Second * 5,
    })

    // 调用服务
    response, err := client.Call("GET", "/users/123", nil, nil)
    if err != nil {
        t.Fatal(err)
    }

    if response.StatusCode != 200 {
        t.Errorf("Expected status 200, got %d", response.StatusCode)
    }
}
```

### 2. 熔断器测试

```go
func TestCircuitBreaker(t *testing.T) {
    // 创建熔断器
    circuitBreaker := microservice.NewCircuitBreaker(microservice.CircuitBreakerConfig{
        Threshold: 3,
        Timeout:   time.Second * 10,
    })

    // 模拟失败
    for i := 0; i < 3; i++ {
        err := circuitBreaker.Execute(func() error {
            return errors.New("service error")
        })
        if err == nil {
            t.Error("Expected error")
        }
    }

    // 验证熔断器状态
    if !circuitBreaker.IsOpen() {
        t.Error("Circuit breaker should be open")
    }

    // 等待超时
    time.Sleep(time.Second * 11)

    // 验证半开状态
    if !circuitBreaker.IsHalfOpen() {
        t.Error("Circuit breaker should be half-open")
    }
}
```

## 🔍 调试和监控

### 1. 服务监控

```go
// 服务监控
type ServiceMonitor struct {
    microservice.Service
    metrics metrics.Collector
}

func (sm *ServiceMonitor) RecordRequest(method, path string, duration time.Duration, statusCode int) {
    // 记录请求指标
    sm.metrics.Increment("service.requests", map[string]string{
        "method": method,
        "path":   path,
        "status": fmt.Sprintf("%d", statusCode),
    })

    // 记录响应时间
    sm.metrics.Histogram("service.response_time", duration.Seconds(), map[string]string{
        "method": method,
        "path":   path,
    })

    // 记录错误率
    if statusCode >= 400 {
        sm.metrics.Increment("service.errors", map[string]string{
            "method": method,
            "path":   path,
        })
    }
}

func (sm *ServiceMonitor) RecordServiceCall(targetService, method, path string, duration time.Duration, success bool) {
    // 记录服务调用指标
    sm.metrics.Increment("service.calls", map[string]string{
        "target_service": targetService,
        "method":         method,
        "path":           path,
        "success":        fmt.Sprintf("%t", success),
    })

    // 记录调用时间
    sm.metrics.Histogram("service.call_duration", duration.Seconds(), map[string]string{
        "target_service": targetService,
    })
}
```

### 2. 分布式追踪

```go
// 分布式追踪监控
type TracingMonitor struct {
    tracer microservice.Tracer
}

func (tm *TracingMonitor) MonitorServiceCall(span microservice.Span, targetService, method, path string) {
    // 添加服务调用标签
    span.SetTag("service.call.target", targetService)
    span.SetTag("service.call.method", method)
    span.SetTag("service.call.path", path)

    // 记录调用开始时间
    span.LogKV("event", "service_call_start")
}

func (tm *TracingMonitor) MonitorServiceResponse(span microservice.Span, response *microservice.Response, err error) {
    if err != nil {
        span.SetTag("error", true)
        span.LogKV("error", err.Error())
    } else {
        span.SetTag("response.status", response.StatusCode)
        span.LogKV("response.size", len(response.Data))
    }

    span.LogKV("event", "service_call_end")
}
```

## 📝 最佳实践

### 1. 服务设计

```go
// 服务接口设计
type UserService interface {
    GetUser(id string) (*Models.User, error)
    CreateUser(user *Models.User) error
    UpdateUser(id string, user *Models.User) error
    DeleteUser(id string) error
    ListUsers(page, limit int) ([]*Models.User, int64, error)
}

// 服务实现
type UserServiceImpl struct {
    db     *database.Connection
    cache  cache.Cache
    logger log.Logger
}

func (s *UserServiceImpl) GetUser(id string) (*Models.User, error) {
    // 检查缓存
    cacheKey := fmt.Sprintf("user:%s", id)
    if cached, exists := s.cache.Get(cacheKey); exists {
        return cached.(*Models.User), nil
    }

    // 从数据库获取
    var user Models.User
    err := s.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }

    // 缓存结果
    s.cache.Set(cacheKey, &user, time.Hour)

    return &user, nil
}
```

### 2. 错误处理

```go
// 统一的错误处理
type ServiceError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (se *ServiceError) Error() string {
    return se.Message
}

// 错误处理中间件
func ErrorHandlingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // 记录错误
                log.Printf("Panic: %v", err)

                // 返回错误响应
                response := map[string]interface{}{
                    "error": "Internal server error",
                    "code":  500,
                }

                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(500)
                json.NewEncoder(w).Encode(response)
            }
        }()

        next(w, r)
    }
}
```

### 3. 服务健康检查

```go
// 健康检查实现
type HealthCheck struct {
    db     *database.Connection
    cache  cache.Cache
    logger log.Logger
}

func (hc *HealthCheck) Check() map[string]interface{} {
    status := map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now(),
        "checks": map[string]interface{}{},
    }

    // 检查数据库连接
    if err := hc.db.Raw("SELECT 1").Error; err != nil {
        status["status"] = "unhealthy"
        status["checks"].(map[string]interface{})["database"] = map[string]interface{}{
            "status": "failed",
            "error":  err.Error(),
        }
    } else {
        status["checks"].(map[string]interface{})["database"] = map[string]interface{}{
            "status": "healthy",
        }
    }

    // 检查缓存连接
    if err := hc.cache.Ping(); err != nil {
        status["status"] = "unhealthy"
        status["checks"].(map[string]interface{})["cache"] = map[string]interface{}{
            "status": "failed",
            "error":  err.Error(),
        }
    } else {
        status["checks"].(map[string]interface{})["cache"] = map[string]interface{}{
            "status": "healthy",
        }
    }

    return status
}
```

### 4. 服务版本管理

```go
// 服务版本管理
type VersionedService struct {
    microservice.Service
    versions map[string]interface{}
}

func (vs *VersionedService) AddVersion(version string, handler interface{}) {
    vs.versions[version] = handler
}

func (vs *VersionedService) GetHandler(version string) (interface{}, error) {
    if handler, exists := vs.versions[version]; exists {
        return handler, nil
    }

    // 返回默认版本
    if handler, exists := vs.versions["latest"]; exists {
        return handler, nil
    }

    return nil, errors.New("no handler found for version")
}

// 版本路由中间件
func VersionMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 从请求头获取版本
        version := r.Header.Get("X-API-Version")
        if version == "" {
            version = "latest"
        }

        // 设置版本到上下文
        ctx := context.WithValue(r.Context(), "version", version)
        r = r.WithContext(ctx)

        next(w, r)
    }
}
```

## 🚀 总结

微服务系统是 Laravel-Go Framework 中重要的功能之一，它提供了：

1. **完整的微服务支持**: 服务注册发现、通信、治理等
2. **高可用性**: 熔断器、负载均衡、重试机制
3. **可观测性**: 分布式追踪、监控、日志
4. **配置管理**: 分布式配置中心
5. **服务网格**: 现代化的服务治理
6. **最佳实践**: 遵循微服务架构的最佳实践

通过合理使用微服务系统，可以构建出可扩展、高可用、易维护的分布式应用程序。
