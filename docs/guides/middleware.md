# 中间件指南

## 🛡️ 中间件概览

中间件是 Laravel-Go Framework 中处理 HTTP 请求的重要组件，它允许你在请求到达控制器之前或响应返回客户端之前执行代码。

## 🚀 快速开始

### 基本中间件结构

```go
// app/Http/Middleware/LoggingMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "log"
    "time"
)

type LoggingMiddleware struct {
    http.Middleware
}

func (m *LoggingMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 请求开始时间
    start := time.Now()
    
    // 记录请求信息
    log.Printf("Request: %s %s", request.Method, request.Path)
    
    // 调用下一个中间件或控制器
    response := next(request)
    
    // 计算处理时间
    duration := time.Since(start)
    
    // 记录响应信息
    log.Printf("Response: %d - %s (%v)", response.StatusCode, request.Path, duration)
    
    return response
}
```

### 注册中间件

```go
// bootstrap/app.go
func Bootstrap() {
    router := routing.NewRouter()
    
    // 注册全局中间件
    router.Use(&middleware.LoggingMiddleware{})
    router.Use(&middleware.CorsMiddleware{})
    
    // 注册路由
    routes.WebRoutes(router)
}
```

## 📋 中间件类型

### 1. 全局中间件

应用于所有请求的中间件。

```go
// 在 bootstrap/app.go 中注册
func Bootstrap() {
    router := routing.NewRouter()
    
    // 全局中间件
    router.Use(&middleware.LoggingMiddleware{})
    router.Use(&middleware.CorsMiddleware{})
    router.Use(&middleware.SecurityMiddleware{})
    
    routes.WebRoutes(router)
}
```

### 2. 路由组中间件

应用于特定路由组的中间件。

```go
// routes/web.go
func WebRoutes(router *routing.Router) {
    // API 路由组
    router.Group("/api", func(group *routing.Router) {
        group.Use(&middleware.AuthMiddleware{})
        group.Use(&middleware.RateLimitMiddleware{})
        
        group.Get("/users", &UserController{}, "Index")
        group.Post("/users", &UserController{}, "Store")
    })
    
    // Admin 路由组
    router.Group("/admin", func(group *routing.Router) {
        group.Use(&middleware.AuthMiddleware{})
        group.Use(&middleware.AdminMiddleware{})
        
        group.Get("/dashboard", &AdminController{}, "Dashboard")
    })
}
```

### 3. 单个路由中间件

应用于特定路由的中间件。

```go
func WebRoutes(router *routing.Router) {
    // 单个路由的中间件
    router.Get("/admin", &AdminController{}, "Dashboard").
        Use(&middleware.AuthMiddleware{}).
        Use(&middleware.AdminMiddleware{})
    
    // 多个中间件
    router.Post("/upload", &FileController{}, "Upload").
        Use(&middleware.AuthMiddleware{}).
        Use(&middleware.FileSizeMiddleware{}).
        Use(&middleware.FileTypeMiddleware{})
}
```

## 🔧 常用中间件

### 1. 认证中间件

```go
// app/Http/Middleware/AuthMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "laravel-go/framework/auth"
)

type AuthMiddleware struct {
    http.Middleware
}

func (m *AuthMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 获取认证 token
    token := request.Headers["Authorization"]
    if token == "" {
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "Unauthorized"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }
    
    // 验证 token
    user, err := auth.ValidateToken(token)
    if err != nil {
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "Invalid token"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }
    
    // 将用户信息添加到请求上下文
    request.Context["user"] = user
    
    return next(request)
}
```

### 2. CORS 中间件

```go
// app/Http/Middleware/CorsMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
)

type CorsMiddleware struct {
    http.Middleware
}

func (m *CorsMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    response := next(request)
    
    // 添加 CORS 头
    response.Headers["Access-Control-Allow-Origin"] = "*"
    response.Headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE, OPTIONS"
    response.Headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization"
    response.Headers["Access-Control-Max-Age"] = "86400"
    
    // 处理预检请求
    if request.Method == "OPTIONS" {
        response.StatusCode = 200
        response.Body = ""
    }
    
    return response
}
```

### 3. 日志中间件

```go
// app/Http/Middleware/LoggingMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "laravel-go/framework/log"
    "time"
)

type LoggingMiddleware struct {
    http.Middleware
}

func (m *LoggingMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    start := time.Now()
    
    // 记录请求信息
    log.Info("Request started", map[string]interface{}{
        "method": request.Method,
        "path":   request.Path,
        "ip":     request.IP,
        "user_agent": request.Headers["User-Agent"],
    })
    
    response := next(request)
    
    duration := time.Since(start)
    
    // 记录响应信息
    log.Info("Request completed", map[string]interface{}{
        "method":     request.Method,
        "path":       request.Path,
        "status":     response.StatusCode,
        "duration":   duration.String(),
        "size":       len(response.Body),
    })
    
    return response
}
```

### 4. 限流中间件

```go
// app/Http/Middleware/RateLimitMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "laravel-go/framework/cache"
    "time"
    "strconv"
)

type RateLimitMiddleware struct {
    http.Middleware
    Limit     int           // 请求限制
    Window    time.Duration // 时间窗口
}

func NewRateLimitMiddleware(limit int, window time.Duration) *RateLimitMiddleware {
    return &RateLimitMiddleware{
        Limit:  limit,
        Window: window,
    }
}

func (m *RateLimitMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 获取客户端 IP
    clientIP := request.IP
    key := "rate_limit:" + clientIP
    
    // 检查当前请求次数
    current, _ := cache.Get(key)
    count := 0
    if current != nil {
        count = current.(int)
    }
    
    // 检查是否超过限制
    if count >= m.Limit {
        return http.Response{
            StatusCode: 429,
            Body:       `{"error": "Too many requests"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
                "Retry-After":  strconv.Itoa(int(m.Window.Seconds())),
            },
        }
    }
    
    // 增加计数
    cache.Set(key, count+1, m.Window)
    
    response := next(request)
    
    // 添加限流头
    response.Headers["X-RateLimit-Limit"] = strconv.Itoa(m.Limit)
    response.Headers["X-RateLimit-Remaining"] = strconv.Itoa(m.Limit - count - 1)
    response.Headers["X-RateLimit-Reset"] = strconv.FormatInt(time.Now().Add(m.Window).Unix(), 10)
    
    return response
}
```

### 5. 缓存中间件

```go
// app/Http/Middleware/CacheMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "laravel-go/framework/cache"
    "crypto/md5"
    "encoding/hex"
    "time"
)

type CacheMiddleware struct {
    http.Middleware
    TTL time.Duration // 缓存时间
}

func NewCacheMiddleware(ttl time.Duration) *CacheMiddleware {
    return &CacheMiddleware{TTL: ttl}
}

func (m *CacheMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 只缓存 GET 请求
    if request.Method != "GET" {
        return next(request)
    }
    
    // 生成缓存键
    cacheKey := m.generateCacheKey(request)
    
    // 尝试从缓存获取
    if cached, found := cache.Get(cacheKey); found {
        return cached.(http.Response)
    }
    
    // 执行请求
    response := next(request)
    
    // 缓存响应
    if response.StatusCode == 200 {
        cache.Set(cacheKey, response, m.TTL)
    }
    
    return response
}

func (m *CacheMiddleware) generateCacheKey(request http.Request) string {
    data := request.Method + ":" + request.Path + ":" + request.QueryString
    hash := md5.Sum([]byte(data))
    return "cache:" + hex.EncodeToString(hash[:])
}
```

## 🛠️ 自定义中间件

### 创建中间件

```bash
# 创建中间件文件
mkdir -p app/Http/Middleware
touch app/Http/Middleware/CustomMiddleware.go
```

### 中间件模板

```go
// app/Http/Middleware/CustomMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
)

type CustomMiddleware struct {
    http.Middleware
    // 添加中间件配置
    Config map[string]interface{}
}

func NewCustomMiddleware(config map[string]interface{}) *CustomMiddleware {
    return &CustomMiddleware{
        Config: config,
    }
}

func (m *CustomMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 请求前处理
    m.beforeRequest(&request)
    
    // 调用下一个中间件或控制器
    response := next(request)
    
    // 响应后处理
    m.afterResponse(&response)
    
    return response
}

func (m *CustomMiddleware) beforeRequest(request *http.Request) {
    // 请求前的处理逻辑
}

func (m *CustomMiddleware) afterResponse(response *http.Response) {
    // 响应后的处理逻辑
}
```

### 带配置的中间件

```go
// app/Http/Middleware/ConfigurableMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "laravel-go/framework/config"
)

type ConfigurableMiddleware struct {
    http.Middleware
    enabled bool
    options map[string]interface{}
}

func NewConfigurableMiddleware() *ConfigurableMiddleware {
    return &ConfigurableMiddleware{
        enabled: config.Get("middleware.configurable.enabled", true).(bool),
        options: config.Get("middleware.configurable.options", map[string]interface{}{}).(map[string]interface{}),
    }
}

func (m *ConfigurableMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    if !m.enabled {
        return next(request)
    }
    
    // 使用配置选项
    if m.options["log_requests"] == true {
        // 记录请求
    }
    
    response := next(request)
    
    // 使用配置选项
    if m.options["add_headers"] == true {
        response.Headers["X-Custom-Header"] = "Custom Value"
    }
    
    return response
}
```

## 🔄 中间件链

### 中间件执行顺序

```go
func WebRoutes(router *routing.Router) {
    // 中间件按注册顺序执行
    router.Group("/api", func(group *routing.Router) {
        group.Use(&middleware.LoggingMiddleware{})      // 1. 日志
        group.Use(&middleware.CorsMiddleware{})         // 2. CORS
        group.Use(&middleware.AuthMiddleware{})         // 3. 认证
        group.Use(&middleware.RateLimitMiddleware{})    // 4. 限流
        
        group.Get("/users", &UserController{}, "Index")
    })
}
```

### 中间件执行流程

```
请求 → 中间件1 → 中间件2 → 中间件3 → 控制器 → 中间件3 → 中间件2 → 中间件1 → 响应
```

### 提前返回

```go
func (m *AuthMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 检查认证
    if !m.isAuthenticated(request) {
        // 提前返回，不调用 next()
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "Unauthorized"}`,
        }
    }
    
    // 继续执行
    return next(request)
}
```

## 🎯 中间件最佳实践

### 1. 单一职责

```go
// ✅ 好的做法：每个中间件只做一件事
type LoggingMiddleware struct {
    http.Middleware
}

type AuthMiddleware struct {
    http.Middleware
}

type CorsMiddleware struct {
    http.Middleware
}

// ❌ 不好的做法：一个中间件做多件事
type AllInOneMiddleware struct {
    http.Middleware
}
```

### 2. 错误处理

```go
func (m *CustomMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Middleware panic: %v", r)
        }
    }()
    
    return next(request)
}
```

### 3. 性能优化

```go
func (m *CacheMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 只对特定路径缓存
    if !m.shouldCache(request.Path) {
        return next(request)
    }
    
    // 缓存逻辑...
}
```

### 4. 配置管理

```go
// config/middleware.go
type MiddlewareConfig struct {
    Auth struct {
        Enabled bool `env:"AUTH_MIDDLEWARE_ENABLED" default:"true"`
        Secret  string `env:"AUTH_SECRET"`
    }
    
    RateLimit struct {
        Enabled bool `env:"RATE_LIMIT_ENABLED" default:"true"`
        Limit    int `env:"RATE_LIMIT_LIMIT" default:"100"`
        Window   int `env:"RATE_LIMIT_WINDOW" default:"3600"`
    }
}
```

## 🔍 调试中间件

### 中间件调试

```go
// app/Http/Middleware/DebugMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "fmt"
)

type DebugMiddleware struct {
    http.Middleware
}

func (m *DebugMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    fmt.Printf("Debug: Request to %s %s\n", request.Method, request.Path)
    
    response := next(request)
    
    fmt.Printf("Debug: Response %d for %s %s\n", response.StatusCode, request.Method, request.Path)
    
    return response
}
```

### 中间件测试

```go
// tests/middleware_test.go
package tests

import (
    "testing"
    "laravel-go/framework/http"
    "laravel-go/app/Http/Middleware"
)

func TestAuthMiddleware(t *testing.T) {
    middleware := &middleware.AuthMiddleware{}
    
    request := http.Request{
        Method: "GET",
        Path:   "/api/users",
        Headers: map[string]string{
            "Authorization": "Bearer valid-token",
        },
    }
    
    response := middleware.Handle(request, func(req http.Request) http.Response {
        return http.Response{StatusCode: 200}
    })
    
    if response.StatusCode != 200 {
        t.Errorf("Expected status 200, got %d", response.StatusCode)
    }
}
```

## 📊 中间件监控

### 中间件性能监控

```go
// app/Http/Middleware/MonitoringMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "laravel-go/framework/metrics"
    "time"
)

type MonitoringMiddleware struct {
    http.Middleware
}

func (m *MonitoringMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    start := time.Now()
    
    response := next(request)
    
    duration := time.Since(start)
    
    // 记录指标
    metrics.Histogram("middleware_duration", duration.Seconds(), map[string]string{
        "middleware": "monitoring",
        "path":       request.Path,
        "method":     request.Method,
    })
    
    metrics.Counter("requests_total", 1, map[string]string{
        "path":   request.Path,
        "method": request.Method,
        "status": string(response.StatusCode),
    })
    
    return response
}
```

## 📝 总结

Laravel-Go Framework 的中间件系统提供了：

1. **灵活性**: 支持全局、路由组和单个路由的中间件
2. **可扩展性**: 易于创建自定义中间件
3. **性能优化**: 支持缓存和限流等性能优化
4. **安全性**: 提供认证、CORS 等安全功能
5. **可观测性**: 支持日志记录和性能监控

通过合理使用中间件，可以构建出安全、高效、可维护的 Web 应用程序。 