# Laravel-Go Framework 技术设计文档

## 1. 总体架构设计

### 1.1 架构概述

Laravel-Go Framework 采用分层架构设计，借鉴 Laravel 的设计理念，同时充分利用 Go 语言的特性。整体架构分为以下几个层次：

```
┌─────────────────────────────────────────────────────────────┐
│                    应用层 (Application Layer)                │
├─────────────────────────────────────────────────────────────┤
│                    服务层 (Service Layer)                   │
├─────────────────────────────────────────────────────────────┤
│                    业务层 (Business Layer)                  │
├─────────────────────────────────────────────────────────────┤
│                    数据层 (Data Layer)                      │
├─────────────────────────────────────────────────────────────┤
│                    基础设施层 (Infrastructure Layer)         │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 核心组件架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Server   │    │   CLI Server    │    │  Queue Worker   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Application   │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │  Service Container │
                    └─────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Router        │    │   Middleware    │    │   Controller    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │      Model      │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Database      │
                    └─────────────────┘
```

## 2. 核心模块设计

### 2.1 应用容器 (Application Container)

#### 2.1.1 设计目标

- 提供依赖注入容器
- 管理服务生命周期
- 支持服务注册和解析
- 实现单例模式

#### 2.1.2 接口设计

```go
// Container 接口定义
type Container interface {
    // 注册服务
    Bind(abstract interface{}, concrete interface{})
    BindSingleton(abstract interface{}, concrete interface{})

    // 解析服务
    Make(abstract interface{}) interface{}
    Resolve(abstract interface{}) error

    // 检查服务是否存在
    Has(abstract interface{}) bool

    // 调用方法并注入依赖
    Call(callback interface{}, parameters ...interface{}) ([]interface{}, error)
}

// 服务提供者接口
type ServiceProvider interface {
    Register(container Container)
    Boot(container Container)
}
```

#### 2.1.3 实现方案

```go
// container.go
type container struct {
    bindings    map[reflect.Type]binding
    singletons  map[reflect.Type]interface{}
    instances   map[reflect.Type]interface{}
    mutex       sync.RWMutex
}

type binding struct {
    concrete interface{}
    shared   bool
    callback func(Container) interface{}
}

func NewContainer() Container {
    return &container{
        bindings:   make(map[reflect.Type]binding),
        singletons: make(map[reflect.Type]interface{}),
        instances:  make(map[reflect.Type]interface{}),
    }
}
```

### 2.2 路由系统 (Router)

#### 2.2.1 设计目标

- 支持 RESTful 路由定义
- 支持路由分组和中间件
- 支持路由参数和查询参数
- 高性能路由匹配

#### 2.2.2 接口设计

```go
// Router 接口
type Router interface {
    // HTTP 方法路由
    Get(path string, handler interface{}) Route
    Post(path string, handler interface{}) Route
    Put(path string, handler interface{}) Route
    Delete(path string, handler interface{}) Route
    Patch(path string, handler interface{}) Route

    // 路由分组
    Group(prefix string, callback func(Router)) Router

    // 中间件
    Middleware(middleware ...interface{}) Router

    // 路由参数
    Where(name, pattern string) Route
}

// Route 接口
type Route interface {
    Name(name string) Route
    Middleware(middleware ...interface{}) Route
    Where(name, pattern string) Route
}
```

#### 2.2.3 实现方案

```go
// router.go
type router struct {
    routes     []*route
    groups     []*routeGroup
    middleware []interface{}
    container  Container
}

type route struct {
    method     string
    path       string
    handler    interface{}
    name       string
    middleware []interface{}
    patterns   map[string]string
}

type routeGroup struct {
    prefix     string
    middleware []interface{}
    routes     []*route
}

// 使用 Radix Tree 进行高性能路由匹配
type radixNode struct {
    path     string
    children map[string]*radixNode
    handler  interface{}
    params   []string
}
```

### 2.3 中间件系统 (Middleware)

#### 2.3.1 设计目标

- 支持中间件链式调用
- 支持全局和路由中间件
- 支持中间件优先级
- 支持条件执行

#### 2.3.2 接口设计

```go
// Middleware 接口
type Middleware interface {
    Handle(request Request, next Next) Response
}

// Next 函数类型
type Next func(Request) Response

// Request 接口
type Request interface {
    Method() string
    Path() string
    Header(name string) string
    Body() []byte
    Param(name string) string
    Query(name string) string
    Context() context.Context
}

// Response 接口
type Response interface {
    Status() int
    Header(name, value string) Response
    Body(data interface{}) Response
    JSON(data interface{}) Response
}
```

#### 2.3.3 实现方案

```go
// middleware.go
type middlewareStack struct {
    middlewares []Middleware
    index       int
}

func (m *middlewareStack) Next(request Request) Response {
    if m.index >= len(m.middlewares) {
        return nil // 到达控制器
    }

    middleware := m.middlewares[m.index]
    m.index++

    return middleware.Handle(request, m.Next)
}

// 常用中间件实现
type corsMiddleware struct{}
type authMiddleware struct{}
type logMiddleware struct{}
type cacheMiddleware struct{}
```

### 2.4 数据库 ORM

#### 2.4.1 设计目标

- 类似 Eloquent 的链式查询
- 支持模型关联关系
- 支持数据库迁移
- 支持多种数据库驱动

#### 2.4.2 接口设计

```go
// Model 接口
type Model interface {
    TableName() string
    PrimaryKey() string
    Timestamps() bool
    Fillable() []string
    Hidden() []string
    Casts() map[string]string
}

// Query 接口
type Query interface {
    Where(column, operator string, value interface{}) Query
    OrWhere(column, operator string, value interface{}) Query
    WhereIn(column string, values []interface{}) Query
    OrderBy(column, direction string) Query
    Limit(limit int) Query
    Offset(offset int) Query
    Get() ([]Model, error)
    First() (Model, error)
    Find(id interface{}) (Model, error)
    Create(data map[string]interface{}) (Model, error)
    Update(data map[string]interface{}) error
    Delete() error
}

// Database 接口
type Database interface {
    Connection(name string) Connection
    DefaultConnection() Connection
    Begin() (Transaction, error)
}

// Connection 接口
type Connection interface {
    Query() Query
    Raw(sql string, args ...interface{}) RawQuery
    Exec(sql string, args ...interface{}) (Result, error)
}
```

#### 2.4.3 实现方案

```go
// model.go
type BaseModel struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (m *BaseModel) TableName() string {
    return inflection.Plural(reflect.TypeOf(m).Elem().Name())
}

func (m *BaseModel) PrimaryKey() string {
    return "id"
}

func (m *BaseModel) Timestamps() bool {
    return true
}

// query.go
type query struct {
    connection Connection
    table      string
    wheres     []whereClause
    orders     []orderClause
    limit      int
    offset     int
    model      Model
}

type whereClause struct {
    column   string
    operator string
    value    interface{}
    boolean  string // "and" or "or"
}

type orderClause struct {
    column    string
    direction string
}
```

### 2.5 模板引擎

#### 2.5.1 设计目标

- 类似 Blade 的模板语法
- 支持模板继承和组件
- 支持模板缓存
- 支持助手函数

#### 2.5.2 接口设计

```go
// Template 接口
type Template interface {
    Render(data interface{}) (string, error)
    RenderString(template string, data interface{}) (string, error)
    AddGlobal(key string, value interface{})
    AddFunction(name string, fn interface{})
}

// View 接口
type View interface {
    Make(name string, data interface{}) (string, error)
    Exists(name string) bool
    AddLocation(path string)
    AddNamespace(namespace, path string)
}
```

#### 2.5.3 实现方案

```go
// template.go
type templateEngine struct {
    templates map[string]*template.Template
    globals   map[string]interface{}
    functions template.FuncMap
    cache     Cache
    mutex     sync.RWMutex
}

// 自定义模板语法解析器
type bladeParser struct {
    tokens    []token
    position  int
    functions template.FuncMap
}

type token struct {
    type_    tokenType
    value    string
    position int
}

type tokenType int

const (
    tokenText tokenType = iota
    tokenDirective
    tokenExpression
    tokenComment
)
```

### 2.6 命令行工具

#### 2.6.1 设计目标

- 类似 Artisan 的命令行体验
- 支持代码生成
- 支持项目管理
- 支持自定义命令

#### 2.6.2 接口设计

```go
// Command 接口
type Command interface {
    Signature() string
    Description() string
    Handle(args []string) error
}

// Console 接口
type Console interface {
    Register(command Command)
    Run(args []string) error
    Call(command string, args []string) error
}
```

#### 2.6.3 实现方案

```go
// console.go
type console struct {
    commands map[string]Command
    output   Output
}

type output struct {
    writer io.Writer
}

// 内置命令
type makeControllerCommand struct{}
type makeModelCommand struct{}
type makeMiddlewareCommand struct{}
type migrateCommand struct{}
type serveCommand struct{}
```

### 2.7 缓存系统

#### 2.7.1 设计目标

- 统一的缓存接口
- 支持多种缓存驱动
- 支持缓存标签
- 支持自动过期

#### 2.7.2 接口设计

```go
// Cache 接口
type Cache interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    Has(key string) bool
    Tags(names ...string) TaggedCache
    Flush() error
}

// TaggedCache 接口
type TaggedCache interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Flush() error
}
```

#### 2.7.3 实现方案

```go
// cache.go
type cacheManager struct {
    stores map[string]Cache
    defaultStore string
}

type memoryStore struct {
    data map[string]cacheItem
    mutex sync.RWMutex
}

type redisStore struct {
    client *redis.Client
}

type cacheItem struct {
    value      interface{}
    expiration time.Time
    tags       []string
}
```

### 2.8 队列系统

#### 2.8.1 设计目标

- 异步任务处理
- 支持多种队列驱动
- 支持任务重试
- 支持任务监控

#### 2.8.2 接口设计

```go
// Queue 接口
type Queue interface {
    Push(job Job) error
    PushDelayed(job Job, delay time.Duration) error
    PushToQueue(job Job, queue string) error
}

// Job 接口
type Job interface {
    Handle() error
    Failed(err error)
    Retries() int
    MaxRetries() int
    Delay() time.Duration
}

// Worker 接口
type Worker interface {
    Start(queue string) error
    Stop() error
    Process(job Job) error
}
```

#### 2.8.3 实现方案

```go
// queue.go
type queueManager struct {
    connections map[string]Queue
    defaultConnection string
}

type databaseQueue struct {
    connection Connection
    table      string
}

type redisQueue struct {
    client *redis.Client
    prefix string
}

type job struct {
    ID        string                 `json:"id"`
    Queue     string                 `json:"queue"`
    Payload   map[string]interface{} `json:"payload"`
    Attempts  int                    `json:"attempts"`
    MaxRetries int                   `json:"max_retries"`
    Delay     time.Duration          `json:"delay"`
    CreatedAt time.Time              `json:"created_at"`
}
```

## 3. 项目结构设计

### 3.1 标准项目结构

```
laravel-go-project/
├── app/
│   ├── Console/
│   │   └── Commands/
│   ├── Http/
│   │   ├── Controllers/
│   │   ├── Middleware/
│   │   └── Requests/
│   ├── Models/
│   ├── Services/
│   └── Providers/
├── bootstrap/
│   ├── app.go
│   └── providers.go
├── config/
│   ├── app.go
│   ├── database.go
│   ├── cache.go
│   └── queue.go
├── database/
│   ├── migrations/
│   └── seeders/
├── public/
│   ├── index.go
│   ├── css/
│   ├── js/
│   └── images/
├── resources/
│   ├── views/
│   ├── lang/
│   └── assets/
├── routes/
│   ├── web.go
│   ├── api.go
│   └── console.go
├── storage/
│   ├── logs/
│   ├── cache/
│   └── uploads/
├── tests/
├── vendor/
├── .env
├── .env.example
├── go.mod
├── go.sum
└── main.go
```

### 3.2 核心文件设计

#### 3.2.1 main.go

```go
package main

import (
    "log"
    "laravel-go/bootstrap"
)

func main() {
    app := bootstrap.App()

    // 启动 HTTP 服务器
    go func() {
        if err := app.Serve(); err != nil {
            log.Fatal(err)
        }
    }()

    // 启动队列工作进程
    go func() {
        if err := app.QueueWorker().Start("default"); err != nil {
            log.Fatal(err)
        }
    }()

    // 等待信号
    app.WaitForShutdown()
}
```

#### 3.2.2 bootstrap/app.go

```go
package bootstrap

import (
    "laravel-go/framework"
)

func App() *framework.Application {
    app := framework.NewApplication()

    // 注册服务提供者
    app.RegisterProviders([]framework.ServiceProvider{
        &AppServiceProvider{},
        &RouteServiceProvider{},
        &DatabaseServiceProvider{},
        &CacheServiceProvider{},
        &QueueServiceProvider{},
    })

    // 启动应用
    app.Boot()

    return app
}
```

## 4. 配置管理设计

### 4.1 配置结构

```go
// config/app.go
type AppConfig struct {
    Name        string `env:"APP_NAME" default:"Laravel-Go"`
    Environment string `env:"APP_ENV" default:"production"`
    Debug       bool   `env:"APP_DEBUG" default:"false"`
    URL         string `env:"APP_URL" default:"http://localhost"`
    Timezone    string `env:"APP_TIMEZONE" default:"UTC"`
    Locale      string `env:"APP_LOCALE" default:"en"`
    Key         string `env:"APP_KEY"`
}

// config/database.go
type DatabaseConfig struct {
    Default     string            `env:"DB_CONNECTION" default:"mysql"`
    Connections map[string]DBConn `env:"DB_CONNECTIONS"`
}

type DBConn struct {
    Driver   string `env:"DB_DRIVER" default:"mysql"`
    Host     string `env:"DB_HOST" default:"127.0.0.1"`
    Port     int    `env:"DB_PORT" default:"3306"`
    Database string `env:"DB_DATABASE"`
    Username string `env:"DB_USERNAME"`
    Password string `env:"DB_PASSWORD"`
    Charset  string `env:"DB_CHARSET" default:"utf8mb4"`
}
```

### 4.2 环境变量管理

```go
// framework/config/config.go
type Config struct {
    data map[string]interface{}
    env  *Env
}

func (c *Config) Get(key string, defaultValue ...interface{}) interface{} {
    // 支持点号分隔的嵌套键
    keys := strings.Split(key, ".")
    value := c.data

    for _, k := range keys {
        if v, ok := value[k]; ok {
            value = v
        } else {
            if len(defaultValue) > 0 {
                return defaultValue[0]
            }
            return nil
        }
    }

    return value
}
```

## 5. 错误处理设计

### 5.1 错误类型定义

```go
// framework/errors/errors.go
type AppError struct {
    Code    int
    Message string
    Err     error
    Stack   []string
}

func (e *AppError) Error() string {
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// 预定义错误类型
var (
    ErrNotFound          = &AppError{Code: 404, Message: "Resource not found"}
    ErrUnauthorized      = &AppError{Code: 401, Message: "Unauthorized"}
    ErrForbidden         = &AppError{Code: 403, Message: "Forbidden"}
    ErrValidation        = &AppError{Code: 422, Message: "Validation failed"}
    ErrInternalServer    = &AppError{Code: 500, Message: "Internal server error"}
)
```

### 5.2 错误处理中间件

```go
// framework/http/middleware/error_handler.go
type ErrorHandlerMiddleware struct{}

func (m *ErrorHandlerMiddleware) Handle(request Request, next Next) Response {
    defer func() {
        if r := recover(); r != nil {
            // 记录 panic 信息
            log.Printf("Panic: %v", r)

            // 返回 500 错误
            return &response{
                status: 500,
                body:   map[string]string{"error": "Internal server error"},
            }
        }
    }()

    response := next(request)

    // 处理应用错误
    if err, ok := response.(*AppError); ok {
        return &response{
            status: err.Code,
            body:   map[string]string{"error": err.Message},
        }
    }

    return response
}
```

## 6. 日志系统设计

### 6.1 日志接口

```go
// framework/log/logger.go
type Logger interface {
    Emergency(message string, context map[string]interface{})
    Alert(message string, context map[string]interface{})
    Critical(message string, context map[string]interface{})
    Error(message string, context map[string]interface{})
    Warning(message string, context map[string]interface{})
    Notice(message string, context map[string]interface{})
    Info(message string, context map[string]interface{})
    Debug(message string, context map[string]interface{})
}

type LogManager struct {
    channels map[string]Logger
    defaultChannel string
}
```

### 6.2 日志驱动

```go
// framework/log/drivers/file.go
type FileLogger struct {
    path   string
    level  string
    writer *os.File
}

// framework/log/drivers/syslog.go
type SyslogLogger struct {
    tag   string
    level string
}

// framework/log/drivers/redis.go
type RedisLogger struct {
    client *redis.Client
    key    string
}
```

## 7. 安全设计

### 7.1 CSRF 保护

```go
// framework/http/middleware/csrf.go
type CSRFMiddleware struct {
    tokenManager TokenManager
}

func (m *CSRFMiddleware) Handle(request Request, next Next) Response {
    if request.Method() == "GET" {
        // 生成 CSRF token
        token := m.tokenManager.Generate()
        response := next(request)
        response.Header("X-CSRF-Token", token)
        return response
    }

    // 验证 CSRF token
    token := request.Header("X-CSRF-Token")
    if !m.tokenManager.Validate(token) {
        return &response{
            status: 419,
            body:   map[string]string{"error": "CSRF token mismatch"},
        }
    }

    return next(request)
}
```

### 7.2 认证中间件

```go
// framework/http/middleware/auth.go
type AuthMiddleware struct {
    auth Auth
}

func (m *AuthMiddleware) Handle(request Request, next Next) Response {
    if !m.auth.Check(request) {
        return &response{
            status: 401,
            body:   map[string]string{"error": "Unauthenticated"},
        }
    }

    return next(request)
}
```

## 8. 性能优化设计

### 8.1 连接池管理

```go
// framework/database/connection_pool.go
type ConnectionPool struct {
    connections chan Connection
    factory     func() Connection
    maxConn     int
    timeout     time.Duration
}

func (p *ConnectionPool) Get() (Connection, error) {
    select {
    case conn := <-p.connections:
        return conn, nil
    case <-time.After(p.timeout):
        return nil, errors.New("connection pool timeout")
    }
}

func (p *ConnectionPool) Put(conn Connection) {
    select {
    case p.connections <- conn:
    default:
        // 连接池已满，关闭连接
        conn.Close()
    }
}
```

### 8.2 缓存优化

```go
// framework/cache/optimization.go
type CacheOptimizer struct {
    cache Cache
    ttl   time.Duration
}

func (o *CacheOptimizer) Remember(key string, callback func() interface{}) interface{} {
    if value, err := o.cache.Get(key); err == nil {
        return value
    }

    value := callback()
    o.cache.Set(key, value, o.ttl)
    return value
}
```

## 9. 测试设计

### 9.1 测试框架集成

```go
// framework/testing/test_case.go
type TestCase struct {
    app *Application
}

func (t *TestCase) Setup() {
    t.app = bootstrap.TestApp()
}

func (t *TestCase) TearDown() {
    // 清理测试数据
}

func (t *TestCase) AssertResponse(response Response, status int) {
    if response.Status() != status {
        t.Fatalf("Expected status %d, got %d", status, response.Status())
    }
}

func (t *TestCase) AssertDatabaseHas(table string, data map[string]interface{}) {
    // 验证数据库中的数据
}
```

### 9.2 模拟器

```go
// framework/testing/mocks.go
type MockCache struct {
    data map[string]interface{}
}

func (m *MockCache) Get(key string) (interface{}, error) {
    if value, exists := m.data[key]; exists {
        return value, nil
    }
    return nil, errors.New("key not found")
}

func (m *MockCache) Set(key string, value interface{}, ttl time.Duration) error {
    m.data[key] = value
    return nil
}
```

## 10. 部署设计

### 10.1 Docker 支持

```dockerfile
# Dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]
```

### 10.2 Kubernetes 配置

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: laravel-go-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: laravel-go-app
  template:
    metadata:
      labels:
        app: laravel-go-app
    spec:
      containers:
        - name: app
          image: laravel-go-app:latest
          ports:
            - containerPort: 8080
          env:
            - name: APP_ENV
              value: "production"
            - name: DB_HOST
              valueFrom:
                secretKeyRef:
                  name: db-secret
                  key: host
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
```

## 11. 监控和指标

### 11.1 健康检查

```go
// framework/http/controllers/health.go
type HealthController struct{}

func (c *HealthController) Check(request Request) Response {
    checks := map[string]interface{}{
        "database": c.checkDatabase(),
        "cache":    c.checkCache(),
        "queue":    c.checkQueue(),
    }

    status := 200
    for _, check := range checks {
        if err, ok := check.(error); ok {
            status = 503
            break
        }
    }

    return &response{
        status: status,
        body:   checks,
    }
}
```

### 11.2 性能指标

```go
// framework/metrics/metrics.go
type Metrics struct {
    requests    *prometheus.CounterVec
    duration    *prometheus.HistogramVec
    errors      *prometheus.CounterVec
    connections *prometheus.GaugeVec
}

func (m *Metrics) RecordRequest(method, path string, status int, duration time.Duration) {
    m.requests.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
    m.duration.WithLabelValues(method, path).Observe(duration.Seconds())
}
```

## 12. 总结

本设计文档详细描述了 Laravel-Go Framework 的技术架构和实现方案。该框架将 Laravel 的优秀设计理念与 Go 语言的高性能特性相结合，为开发者提供了一个优雅、高效的开发平台。

### 主要特点：

1. **分层架构**：清晰的分层设计，便于维护和扩展
2. **依赖注入**：松耦合的组件设计，便于测试和替换
3. **中间件系统**：灵活的请求处理管道
4. **ORM 系统**：类似 Eloquent 的数据库操作体验
5. **模板引擎**：类似 Blade 的模板语法
6. **命令行工具**：类似 Artisan 的开发体验
7. **缓存系统**：统一的缓存接口，支持多种驱动
8. **队列系统**：异步任务处理能力
9. **安全特性**：完善的认证和授权机制
10. **性能优化**：连接池、缓存等性能优化措施
11. **测试支持**：完善的测试框架和工具
12. **部署支持**：Docker 和 Kubernetes 支持

该框架将为 Go 开发者提供类似 Laravel 的开发体验，同时充分利用 Go 语言的高性能特性，是一个理想的现代化 Web 开发框架。
