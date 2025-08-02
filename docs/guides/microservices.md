# 微服务指南

## 📖 概述

Laravel-Go Framework 提供了完整的微服务架构支持，包括服务注册与发现、负载均衡、服务间通信、配置管理和监控等功能。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [微服务系统 API 参考](../api/microservice.md)

## 🚀 快速开始

### 1. 基本微服务

```go
// 用户微服务
type UserService struct {
    microservice.Service
    db *database.Connection
}

func NewUserService() *UserService {
    return &UserService{
        Service: microservice.NewService("user-service", "1.0.0"),
        db:      database.NewConnection(),
    }
}

// 启动用户服务
func (s *UserService) Start() error {
    // 注册服务
    if err := s.Register(); err != nil {
        return err
    }

    // 启动 HTTP 服务器
    server := http.NewServer()

    // 注册路由
    s.registerRoutes(server)

    // 启动服务器
    return server.Start(":8081")
}

// 注册路由
func (s *UserService) registerRoutes(server *http.Server) {
    router := server.Router()

    // 用户相关路由
    router.Get("/users", s.GetUsers)
    router.Get("/users/{id}", s.GetUser)
    router.Post("/users", s.CreateUser)
    router.Put("/users/{id}", s.UpdateUser)
    router.Delete("/users/{id}", s.DeleteUser)

    // 健康检查
    router.Get("/health", s.HealthCheck)
}

// 健康检查
func (s *UserService) HealthCheck(request http.Request) http.Response {
    return http.Response{
        StatusCode: 200,
        Body:       `{"status": "healthy", "service": "user-service"}`,
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }
}

// 主函数
func main() {
    service := NewUserService()

    if err := service.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### 2. 服务注册与发现

```go
// 服务注册
func (s *UserService) Register() error {
    // 配置服务注册中心
    registry := microservice.NewConsulRegistry(&microservice.ConsulConfig{
        Address: "localhost:8500",
        Token:   "",
    })

    // 注册服务
    serviceInfo := &microservice.ServiceInfo{
        Name:    "user-service",
        Version: "1.0.0",
        Address: "localhost",
        Port:    8081,
        Tags:    []string{"user", "api"},
        Metadata: map[string]string{
            "protocol": "http",
            "health":   "/health",
        },
    }

    return registry.Register(serviceInfo)
}

// 服务发现
func (s *UserService) DiscoverService(serviceName string) ([]*microservice.ServiceInfo, error) {
    registry := microservice.NewConsulRegistry(&microservice.ConsulConfig{
        Address: "localhost:8500",
    })

    return registry.Discover(serviceName)
}
```

## 🔧 服务通信

### 1. HTTP 客户端

```go
// HTTP 客户端
type HttpClient struct {
    client *http.Client
    baseURL string
}

func NewHttpClient(baseURL string) *HttpClient {
    return &HttpClient{
        client: &http.Client{
            Timeout: time.Second * 30,
        },
        baseURL: baseURL,
    }
}

// 调用其他服务
func (c *HttpClient) CallUserService(userID uint) (*User, error) {
    url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("service call failed: %d", resp.StatusCode)
    }

    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }

    return &user, nil
}

// 在订单服务中调用用户服务
type OrderService struct {
    microservice.Service
    userClient *HttpClient
}

func (s *OrderService) GetOrderWithUser(orderID uint) (*OrderWithUser, error) {
    // 获取订单
    order, err := s.GetOrder(orderID)
    if err != nil {
        return nil, err
    }

    // 调用用户服务获取用户信息
    user, err := s.userClient.CallUserService(order.UserID)
    if err != nil {
        return nil, err
    }

    return &OrderWithUser{
        Order: order,
        User:  user,
    }, nil
}
```

### 2. gRPC 通信

```go
// gRPC 服务定义
type UserServiceServer struct {
    pb.UnimplementedUserServiceServer
    db *database.Connection
}

// 实现 gRPC 方法
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    var user User
    err := s.db.First(&user, req.Id).Error
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "User not found")
    }

    return &pb.GetUserResponse{
        User: &pb.User{
            Id:    uint64(user.ID),
            Name:  user.Name,
            Email: user.Email,
        },
    }, nil
}

// 启动 gRPC 服务器
func (s *UserService) StartGRPC() error {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        return err
    }

    grpcServer := grpc.NewServer()
    pb.RegisterUserServiceServer(grpcServer, &UserServiceServer{db: s.db})

    return grpcServer.Serve(lis)
}

// gRPC 客户端
type GrpcClient struct {
    conn   *grpc.ClientConn
    client pb.UserServiceClient
}

func NewGrpcClient(address string) (*GrpcClient, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }

    return &GrpcClient{
        conn:   conn,
        client: pb.NewUserServiceClient(conn),
    }, nil
}

func (c *GrpcClient) GetUser(userID uint) (*pb.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{Id: uint64(userID)})
    if err != nil {
        return nil, err
    }

    return resp.User, nil
}
```

## 🔄 负载均衡

### 1. 客户端负载均衡

```go
// 负载均衡器
type LoadBalancer struct {
    services []*microservice.ServiceInfo
    strategy LoadBalanceStrategy
    mutex    sync.RWMutex
}

type LoadBalanceStrategy interface {
    Select(services []*microservice.ServiceInfo) *microservice.ServiceInfo
}

// 轮询策略
type RoundRobinStrategy struct {
    current int
    mutex   sync.Mutex
}

func (r *RoundRobinStrategy) Select(services []*microservice.ServiceInfo) *microservice.ServiceInfo {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    if len(services) == 0 {
        return nil
    }

    service := services[r.current]
    r.current = (r.current + 1) % len(services)

    return service
}

// 随机策略
type RandomStrategy struct{}

func (r *RandomStrategy) Select(services []*microservice.ServiceInfo) *microservice.ServiceInfo {
    if len(services) == 0 {
        return nil
    }

    return services[rand.Intn(len(services))]
}

// 加权轮询策略
type WeightedRoundRobinStrategy struct {
    current int
    mutex   sync.Mutex
}

func (w *WeightedRoundRobinStrategy) Select(services []*microservice.ServiceInfo) *microservice.ServiceInfo {
    w.mutex.Lock()
    defer w.mutex.Unlock()

    if len(services) == 0 {
        return nil
    }

    // 根据权重选择服务
    totalWeight := 0
    for _, service := range services {
        weight := 1
        if w, ok := service.Metadata["weight"]; ok {
            if weight, err := strconv.Atoi(w); err == nil {
                totalWeight += weight
            }
        }
    }

    // 实现加权轮询逻辑
    // ...

    return services[w.current]
}

// 使用负载均衡器
type ServiceClient struct {
    registry    microservice.Registry
    balancer    *LoadBalancer
    httpClient  *http.Client
}

func (c *ServiceClient) CallService(serviceName, path string) (*http.Response, error) {
    // 发现服务
    services, err := c.registry.Discover(serviceName)
    if err != nil {
        return nil, err
    }

    // 更新负载均衡器的服务列表
    c.balancer.UpdateServices(services)

    // 选择服务
    service := c.balancer.Select()
    if service == nil {
        return nil, errors.New("no available service")
    }

    // 构建请求 URL
    url := fmt.Sprintf("http://%s:%d%s", service.Address, service.Port, path)

    // 发送请求
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    return c.httpClient.Do(req)
}
```

### 2. 服务端负载均衡

```go
// 反向代理负载均衡
type ReverseProxy struct {
    services []*microservice.ServiceInfo
    proxy    *httputil.ReverseProxy
}

func NewReverseProxy(services []*microservice.ServiceInfo) *ReverseProxy {
    director := func(req *http.Request) {
        // 选择后端服务
        service := selectService(services)

        req.URL.Scheme = "http"
        req.URL.Host = fmt.Sprintf("%s:%d", service.Address, service.Port)
    }

    return &ReverseProxy{
        services: services,
        proxy:    &httputil.ReverseProxy{Director: director},
    }
}

func (r *ReverseProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    r.proxy.ServeHTTP(w, req)
}

// API 网关
type ApiGateway struct {
    router *routing.Router
    registry microservice.Registry
}

func (g *ApiGateway) Start() error {
    // 注册路由
    g.registerRoutes()

    // 启动服务器
    server := http.NewServer()
    server.Router = g.router

    return server.Start(":8080")
}

func (g *ApiGateway) registerRoutes() {
    // 用户服务路由
    g.router.Get("/api/users", g.proxyToService("user-service", "/users"))
    g.router.Get("/api/users/{id}", g.proxyToService("user-service", "/users/{id}"))

    // 订单服务路由
    g.router.Get("/api/orders", g.proxyToService("order-service", "/orders"))
    g.router.Get("/api/orders/{id}", g.proxyToService("order-service", "/orders/{id}"))
}

func (g *ApiGateway) proxyToService(serviceName, path string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 发现服务
        services, err := g.registry.Discover(serviceName)
        if err != nil {
            http.Error(w, "Service not available", 503)
            return
        }

        if len(services) == 0 {
            http.Error(w, "No available service", 503)
            return
        }

        // 选择服务（简单轮询）
        service := services[0]

        // 构建目标 URL
        targetURL := fmt.Sprintf("http://%s:%d%s", service.Address, service.Port, path)

        // 创建反向代理
        proxy := httputil.NewSingleHostReverseProxy(&url.URL{
            Scheme: "http",
            Host:   fmt.Sprintf("%s:%d", service.Address, service.Port),
        })

        // 转发请求
        proxy.ServeHTTP(w, r)
    }
}
```

## 🔧 配置管理

### 1. 配置中心

```go
// 配置管理器
type ConfigManager struct {
    registry microservice.Registry
    configs  map[string]interface{}
    mutex    sync.RWMutex
}

func NewConfigManager(registry microservice.Registry) *ConfigManager {
    return &ConfigManager{
        registry: registry,
        configs:  make(map[string]interface{}),
    }
}

// 获取配置
func (c *ConfigManager) GetConfig(key string) (interface{}, error) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()

    if value, exists := c.configs[key]; exists {
        return value, nil
    }

    return nil, fmt.Errorf("config not found: %s", key)
}

// 设置配置
func (c *ConfigManager) SetConfig(key string, value interface{}) error {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    c.configs[key] = value

    // 通知其他服务配置变更
    return c.notifyConfigChange(key, value)
}

// 监听配置变更
func (c *ConfigManager) WatchConfig(key string, callback func(interface{})) {
    // 实现配置监听逻辑
    go func() {
        for {
            // 监听配置变更事件
            // 当配置变更时调用 callback
            time.Sleep(time.Second * 5)
        }
    }()
}

// 在服务中使用配置
type UserService struct {
    microservice.Service
    config *ConfigManager
}

func (s *UserService) GetDatabaseConfig() (*DatabaseConfig, error) {
    config, err := s.config.GetConfig("database")
    if err != nil {
        return nil, err
    }

    var dbConfig DatabaseConfig
    if err := mapstructure.Decode(config, &dbConfig); err != nil {
        return nil, err
    }

    return &dbConfig, nil
}
```

### 2. 环境配置

```go
// 环境配置
type EnvironmentConfig struct {
    Environment string            `json:"environment"`
    Services    map[string]string `json:"services"`
    Database    DatabaseConfig    `json:"database"`
    Redis       RedisConfig       `json:"redis"`
    Logging     LoggingConfig     `json:"logging"`
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Database string `json:"database"`
    Username string `json:"username"`
    Password string `json:"password"`
}

// 加载环境配置
func LoadEnvironmentConfig() (*EnvironmentConfig, error) {
    env := os.Getenv("ENVIRONMENT")
    if env == "" {
        env = "development"
    }

    configPath := fmt.Sprintf("config/%s.json", env)

    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var config EnvironmentConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

## 📊 监控和追踪

### 1. 服务监控

```go
// 服务监控
type ServiceMonitor struct {
    metrics map[string]interface{}
    mutex   sync.RWMutex
}

func NewServiceMonitor() *ServiceMonitor {
    return &ServiceMonitor{
        metrics: make(map[string]interface{}),
    }
}

// 记录指标
func (m *ServiceMonitor) RecordMetric(name string, value interface{}) {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    m.metrics[name] = value
}

// 获取指标
func (m *ServiceMonitor) GetMetrics() map[string]interface{} {
    m.mutex.RLock()
    defer m.mutex.RUnlock()

    metrics := make(map[string]interface{})
    for k, v := range m.metrics {
        metrics[k] = v
    }

    return metrics
}

// 监控中间件
type MonitorMiddleware struct {
    http.Middleware
    monitor *ServiceMonitor
}

func (m *MonitorMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    start := time.Now()

    response := next(request)

    duration := time.Since(start)

    // 记录请求指标
    m.monitor.RecordMetric("request_count", m.monitor.GetMetrics()["request_count"].(int)+1)
    m.monitor.RecordMetric("response_time", duration)

    if response.StatusCode >= 400 {
        m.monitor.RecordMetric("error_count", m.monitor.GetMetrics()["error_count"].(int)+1)
    }

    return response
}

// 监控端点
func (s *UserService) Metrics(request http.Request) http.Response {
    metrics := s.monitor.GetMetrics()

    return http.Response{
        StatusCode: 200,
        Body:       toJSON(metrics),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }
}
```

### 2. 分布式追踪

```go
// 分布式追踪
type TraceContext struct {
    TraceID   string `json:"trace_id"`
    SpanID    string `json:"span_id"`
    ParentID  string `json:"parent_id"`
    Service   string `json:"service"`
    Operation string `json:"operation"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
}

// 追踪中间件
type TraceMiddleware struct {
    http.Middleware
    tracer *Tracer
}

func (m *TraceMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 从请求头获取追踪信息
    traceID := request.Headers["X-Trace-ID"]
    if traceID == "" {
        traceID = generateTraceID()
    }

    spanID := request.Headers["X-Span-ID"]
    if spanID == "" {
        spanID = generateSpanID()
    }

    // 创建追踪上下文
    traceCtx := &TraceContext{
        TraceID:   traceID,
        SpanID:    spanID,
        ParentID:  request.Headers["X-Parent-ID"],
        Service:   "user-service",
        Operation: request.Method + " " + request.Path,
        StartTime: time.Now(),
    }

    // 将追踪信息添加到请求上下文
    request.Context["trace"] = traceCtx

    response := next(request)

    traceCtx.EndTime = time.Now()

    // 记录追踪信息
    m.tracer.RecordSpan(traceCtx)

    // 添加追踪头到响应
    response.Headers["X-Trace-ID"] = traceID
    response.Headers["X-Span-ID"] = spanID

    return response
}

// 追踪器
type Tracer struct {
    spans []*TraceContext
    mutex sync.RWMutex
}

func (t *Tracer) RecordSpan(span *TraceContext) {
    t.mutex.Lock()
    defer t.mutex.Unlock()

    t.spans = append(t.spans, span)
}

func (t *Tracer) GetSpans() []*TraceContext {
    t.mutex.RLock()
    defer t.mutex.RUnlock()

    spans := make([]*TraceContext, len(t.spans))
    copy(spans, t.spans)

    return spans
}
```

## 🛡️ 熔断器

### 1. 熔断器实现

```go
// 熔断器状态
type CircuitBreakerState int

const (
    StateClosed CircuitBreakerState = iota
    StateOpen
    StateHalfOpen
)

// 熔断器
type CircuitBreaker struct {
    state       CircuitBreakerState
    failureCount int
    lastFailure  time.Time
    threshold    int
    timeout      time.Duration
    mutex        sync.RWMutex
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        state:    StateClosed,
        threshold: threshold,
        timeout:   timeout,
    }
}

// 执行操作
func (c *CircuitBreaker) Execute(operation func() error) error {
    if !c.canExecute() {
        return errors.New("circuit breaker is open")
    }

    err := operation()
    c.recordResult(err)

    return err
}

func (c *CircuitBreaker) canExecute() bool {
    c.mutex.RLock()
    defer c.mutex.RUnlock()

    switch c.state {
    case StateClosed:
        return true
    case StateOpen:
        if time.Since(c.lastFailure) > c.timeout {
            c.state = StateHalfOpen
            return true
        }
        return false
    case StateHalfOpen:
        return true
    default:
        return false
    }
}

func (c *CircuitBreaker) recordResult(err error) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if err != nil {
        c.failureCount++
        c.lastFailure = time.Now()

        if c.failureCount >= c.threshold {
            c.state = StateOpen
        }
    } else {
        if c.state == StateHalfOpen {
            c.state = StateClosed
            c.failureCount = 0
        }
    }
}

// 在服务中使用熔断器
type ServiceClient struct {
    httpClient     *http.Client
    circuitBreaker *CircuitBreaker
}

func (c *ServiceClient) CallService(url string) (*http.Response, error) {
    var response *http.Response
    var err error

    err = c.circuitBreaker.Execute(func() error {
        resp, e := c.httpClient.Get(url)
        if e != nil {
            return e
        }

        if resp.StatusCode >= 500 {
            return fmt.Errorf("server error: %d", resp.StatusCode)
        }

        response = resp
        return nil
    })

    return response, err
}
```

## 📚 总结

Laravel-Go Framework 的微服务系统提供了：

1. **服务注册与发现**: Consul、Etcd、Nacos 支持
2. **服务通信**: HTTP、gRPC 通信
3. **负载均衡**: 多种负载均衡策略
4. **配置管理**: 配置中心和环境配置
5. **监控追踪**: 服务监控和分布式追踪
6. **熔断器**: 服务保护机制
7. **API 网关**: 统一入口和路由

通过合理使用微服务系统，可以构建高可用、可扩展的分布式应用程序。
