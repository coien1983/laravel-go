# 错误处理增强最佳实践

## 🚨 概述

本文档介绍了 Laravel-Go Framework 中错误处理的增强功能，包括性能监控集成、恢复机制和最佳实践。

## 🔧 核心组件

### 1. 错误处理器 (ErrorHandler)

```go
// 创建默认错误处理器
logger := &CustomLogger{}
errorHandler := errors.NewDefaultErrorHandler(logger)

// 错误处理器接口
type ErrorHandler interface {
    Handle(err error) error
    Log(err error)
    Report(err error)
}
```

### 2. 安全执行包装器

```go
// 基本安全执行
err := errors.SafeExecute(func() error {
    // 可能发生panic的代码
    return nil
})

// 带上下文的安全执行
err := errors.SafeExecuteWithContext(ctx, func() error {
    // 检查上下文是否已取消
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    // 执行代码
    return nil
})
```

### 3. 恢复中间件

```go
// HTTP恢复中间件
recoveryMiddleware := middleware.NewRecoveryMiddleware(errorHandler, logger)

// 包装处理器
http.HandleFunc("/api", recoveryMiddleware.Handle(http.HandlerFunc(handler)).ServeHTTP)

// 安全处理器包装器
http.HandleFunc("/api", middleware.SafeHandler(handler, errorHandler))
```

## 📊 性能监控集成

### 1. 增强的HTTP监控器

```go
type EnhancedHTTPMonitor struct {
    *performance.HTTPMonitor
    errorHandler errors.ErrorHandler
    errorRate    float64
}

func (ehm *EnhancedHTTPMonitor) RecordRequestWithErrorHandling(method, path string, size int64) {
    defer func() {
        if r := recover(); r != nil {
            if ehm.errorHandler != nil {
                err := errors.New(fmt.Sprintf("HTTP monitor panic: %v", r))
                ehm.errorHandler.Handle(err)
            }
        }
    }()

    ehm.RecordRequest(method, path, size)
}
```

### 2. 增强的数据库监控器

```go
type EnhancedDatabaseMonitor struct {
    *performance.DatabaseMonitor
    errorHandler errors.ErrorHandler
    timeoutRate  float64
}

func (edm *EnhancedDatabaseMonitor) RecordQueryWithErrorHandling(query string, duration time.Duration, success bool, err error) {
    defer func() {
        if r := recover(); r != nil {
            if edm.errorHandler != nil {
                panicErr := errors.New(fmt.Sprintf("Database monitor panic: %v", r))
                edm.errorHandler.Handle(panicErr)
            }
        }
    }()

    // 模拟超时
    if time.Now().UnixNano()%100 < int64(edm.timeoutRate*100) {
        success = false
        err = errors.Wrap(ErrResourceExhausted, "database query timeout")
    }

    edm.RecordQuery(query, duration, success, err)
}
```

### 3. 增强的缓存监控器

```go
type EnhancedCacheMonitor struct {
    *performance.CacheMonitor
    errorHandler errors.ErrorHandler
    unavailable  bool
}

func (ecm *EnhancedCacheMonitor) RecordGetWithErrorHandling(key string, duration time.Duration, hit bool, err error) {
    defer func() {
        if r := recover(); r != nil {
            if ecm.errorHandler != nil {
                panicErr := errors.New(fmt.Sprintf("Cache monitor panic: %v", r))
                ecm.errorHandler.Handle(panicErr)
            }
        }
    }()

    // 模拟缓存服务不可用
    if ecm.unavailable {
        hit = false
        err = errors.Wrap(ErrResourceExhausted, "cache service unavailable")
    }

    ecm.RecordGet(key, duration, hit, err)
}
```

## 🚨 告警系统增强

### 1. 增强的告警系统

```go
type EnhancedAlertSystem struct {
    *performance.AlertSystem
    errorHandler errors.ErrorHandler
}

func (eas *EnhancedAlertSystem) AddRuleWithErrorHandling(rule *performance.AlertRule) error {
    defer func() {
        if r := recover(); r != nil {
            if eas.errorHandler != nil {
                err := errors.New(fmt.Sprintf("Alert system panic: %v", r))
                eas.errorHandler.Handle(err)
            }
        }
    }()

    return eas.AddRule(rule)
}
```

### 2. 优化的告警规则

```go
// 降低错误率阈值，提高敏感度
errorRule := &performance.AlertRule{
    ID:          "error_rate_high",
    Name:        "错误率过高",
    Description: "HTTP错误率超过3%", // 从5%降低到3%
    MetricName:  "http_errors_total",
    Condition:   ">",
    Threshold:   3.0,
    Level:       performance.AlertLevelCritical,
    Enabled:     true,
    Actions:     []string{"log", "email", "webhook"},
}

// 添加响应时间告警
responseTimeRule := &performance.AlertRule{
    ID:          "response_time_high",
    Name:        "响应时间过长",
    Description: "平均响应时间超过500ms",
    MetricName:  "http_response_time",
    Condition:   ">",
    Threshold:   500.0,
    Level:       performance.AlertLevelWarning,
    Enabled:     true,
    Actions:     []string{"log"},
}
```

## 🛠️ 服务层错误处理

### 1. 用户服务示例

```go
type UserService struct {
    errorHandler errors.ErrorHandler
}

func (s *UserService) GetUser(id int) (*User, error) {
    var user *User
    var err error
    
    errors.SafeExecuteWithContext(context.Background(), func() error {
        if id <= 0 {
            err = errors.Wrap(ErrInvalidInput, "invalid user id")
            return err
        }

        // 模拟数据库查询
        if id == 999 {
            // 模拟数据库超时
            time.Sleep(3 * time.Second)
            err = errors.Wrap(ErrDatabaseTimeout, "database query timeout")
            return err
        }

        if id > 100 {
            err = errors.Wrap(ErrUserNotFound, fmt.Sprintf("user %d not found", id))
            return err
        }

        // 模拟成功返回
        user = &User{
            ID:    id,
            Name:  fmt.Sprintf("User %d", id),
            Email: fmt.Sprintf("user%d@example.com", id),
        }
        return nil
    })
    
    return user, err
}
```

### 2. 缓存服务示例

```go
type CacheService struct {
    errorHandler errors.ErrorHandler
    available    bool
}

func (s *CacheService) Get(key string) (interface{}, error) {
    var result interface{}
    var err error
    
    errors.SafeExecuteWithContext(context.Background(), func() error {
        if !s.available {
            err = errors.Wrap(ErrCacheUnavailable, "cache service is down")
            return err
        }

        // 模拟缓存未命中
        if key == "miss" {
            err = errors.Wrap(errors.New("cache miss"), "key not found in cache")
            return err
        }

        result = fmt.Sprintf("cached_value_for_%s", key)
        return nil
    })
    
    return result, err
}
```

## 🎯 控制器层错误处理

### 1. 错误处理方法

```go
type UserController struct {
    userService  *UserService
    cacheService *CacheService
    errorHandler errors.ErrorHandler
}

func (c *UserController) handleError(w http.ResponseWriter, err error) {
    // 使用错误处理器处理错误
    processedErr := c.errorHandler.Handle(err)

    // 根据错误类型返回相应的HTTP状态码
    if appErr := errors.GetAppError(processedErr); appErr != nil {
        http.Error(w, appErr.Message, appErr.Code)
    } else {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
```

### 2. 处理器示例

```go
func (c *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
    // 解析用户ID
    id := 1 // 简化处理，实际应该从URL参数获取

    // 尝试从缓存获取
    cacheKey := fmt.Sprintf("user:%d", id)
    if cached, err := c.cacheService.Get(cacheKey); err == nil {
        // 缓存命中
        fmt.Fprintf(w, "Cache hit: %v\n", cached)
        return
    }

    // 从数据库获取
    user, err := c.userService.GetUser(id)
    if err != nil {
        // 处理错误
        c.handleError(w, err)
        return
    }

    // 缓存结果
    if err := c.cacheService.Set(cacheKey, user); err != nil {
        // 记录缓存错误，但不影响主流程
        c.errorHandler.Handle(errors.Wrap(err, "failed to cache user"))
    }

    // 返回成功响应
    fmt.Fprintf(w, "User: %+v\n", user)
}
```

## 📈 性能优化建议

### 1. 错误率控制

- 降低错误率告警阈值（从5%到3%）
- 实现错误重试机制
- 添加熔断器模式

### 2. 响应时间优化

- 监控平均响应时间
- 设置合理的超时时间
- 实现异步处理

### 3. 资源管理

- 监控CPU和内存使用率
- 实现资源限制
- 添加自动扩缩容

## 🔍 监控和调试

### 1. 日志记录

```go
type CustomLogger struct{}

func (l *CustomLogger) Error(message string, context map[string]interface{}) {
    log.Printf("[ERROR] %s: %+v", message, context)
}

func (l *CustomLogger) Warning(message string, context map[string]interface{}) {
    log.Printf("[WARN] %s: %+v", message, context)
}

func (l *CustomLogger) Info(message string, context map[string]interface{}) {
    log.Printf("[INFO] %s: %+v", message, context)
}

func (l *CustomLogger) Debug(message string, context map[string]interface{}) {
    log.Printf("[DEBUG] %s: %+v", message, context)
}
```

### 2. 健康检查

```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{
        "status": "healthy", 
        "timestamp": "%s", 
        "error_handling": "enhanced"
    }`, time.Now().Format(time.RFC3339))
})
```

## 🚀 部署建议

### 1. 生产环境配置

- 启用所有错误处理机制
- 配置适当的告警阈值
- 设置日志轮转

### 2. 监控指标

- 错误率
- 响应时间
- 资源使用率
- 告警数量

### 3. 故障恢复

- 自动重启机制
- 降级策略
- 备份和恢复

## 📚 示例代码

完整的示例代码请参考：
- `examples/error_handling_demo/main.go` - 基础错误处理演示
- `examples/performance_enhanced_demo/main.go` - 增强性能监控演示

## 🔗 相关文档

- [错误处理基础](../guides/error-handling.md)
- [性能监控指南](../guides/performance.md)
- [HTTP中间件](../guides/http.md)
- [最佳实践](../best-practices/) 