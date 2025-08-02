# 日志系统指南

## 📖 概述

Laravel-Go Framework 提供了强大的日志系统，支持多种日志级别、日志驱动、日志格式化、日志轮转等功能，帮助开发者记录和监控应用程序的运行状态。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [日志系统 API 参考](../api/log.md)

## 🚀 快速开始

### 1. 基本日志使用

```go
// 创建日志实例
logger := log.New()

// 基本日志记录
logger.Info("Application started")
logger.Warning("Database connection slow")
logger.Error("Failed to process request", map[string]interface{}{
    "error": "connection timeout",
    "url":   "/api/users",
})

// 带上下文的日志
logger.Info("User logged in", map[string]interface{}{
    "user_id": 123,
    "ip":      "192.168.1.1",
    "time":    time.Now(),
})
```

### 2. 日志级别

```go
// 不同级别的日志
logger.Debug("Debug information", map[string]interface{}{
    "request_id": "abc123",
    "params":     request.Query,
})

logger.Info("User action", map[string]interface{}{
    "action": "create_post",
    "user_id": user.ID,
})

logger.Warning("Performance issue", map[string]interface{}{
    "query_time": "2.5s",
    "sql":        "SELECT * FROM users",
})

logger.Error("System error", map[string]interface{}{
    "error":   err.Error(),
    "stack":   err.StackTrace(),
    "context": "user_service",
})

logger.Critical("Critical system failure", map[string]interface{}{
    "error": "database connection lost",
    "time":  time.Now(),
})
```

### 3. 日志驱动配置

```go
// 文件驱动
config.Set("log.default", "file")
config.Set("log.file.path", "storage/logs/app.log")
config.Set("log.file.level", "info")
config.Set("log.file.max_size", "100MB")
config.Set("log.file.max_age", "30d")
config.Set("log.file.max_backups", 10)

// Redis 驱动
config.Set("log.redis.enabled", true)
config.Set("log.redis.host", "localhost")
config.Set("log.redis.port", 6379)
config.Set("log.redis.key", "app:logs")

// 数据库驱动
config.Set("log.database.enabled", true)
config.Set("log.database.table", "logs")
config.Set("log.database.connection", "mysql")

// 多驱动配置
config.Set("log.channels", map[string]interface{}{
    "file": map[string]interface{}{
        "driver": "file",
        "path":   "storage/logs/app.log",
        "level":  "info",
    },
    "error": map[string]interface{}{
        "driver": "file",
        "path":   "storage/logs/error.log",
        "level":  "error",
    },
    "debug": map[string]interface{}{
        "driver": "file",
        "path":   "storage/logs/debug.log",
        "level":  "debug",
    },
})
```

### 4. 结构化日志

```go
// 结构化日志记录
type LogContext struct {
    RequestID string                 `json:"request_id"`
    UserID    uint                   `json:"user_id,omitempty"`
    IP        string                 `json:"ip"`
    UserAgent string                 `json:"user_agent"`
    Method    string                 `json:"method"`
    Path      string                 `json:"path"`
    Duration  time.Duration          `json:"duration"`
    Extra     map[string]interface{} `json:"extra,omitempty"`
}

// 创建日志上下文
func CreateLogContext(request http.Request) *LogContext {
    return &LogContext{
        RequestID: request.Headers["X-Request-ID"],
        UserID:    getUserID(request),
        IP:        request.IP,
        UserAgent: request.Headers["User-Agent"],
        Method:    request.Method,
        Path:      request.Path,
        Extra:     make(map[string]interface{}),
    }
}

// 使用结构化日志
func (c *UserController) Show(id string, request http.Request) http.Response {
    ctx := CreateLogContext(request)
    ctx.Extra["user_id"] = id

    logger := log.WithContext(ctx)

    logger.Info("User profile accessed")

    user, err := c.userService.GetUser(uint(id))
    if err != nil {
        logger.Error("Failed to get user", map[string]interface{}{
            "error": err.Error(),
        })
        return c.JsonError("User not found", 404)
    }

    logger.Info("User profile retrieved successfully")
    return c.Json(user)
}
```

### 5. 日志中间件

```go
// 请求日志中间件
type RequestLogMiddleware struct {
    logger *log.Logger
}

func (m *RequestLogMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    start := time.Now()

    // 记录请求开始
    m.logger.Info("Request started", map[string]interface{}{
        "method": request.Method,
        "path":   request.Path,
        "ip":     request.IP,
        "user_agent": request.Headers["User-Agent"],
    })

    // 处理请求
    response := next(request)

    // 计算处理时间
    duration := time.Since(start)

    // 记录请求完成
    logLevel := "info"
    if response.StatusCode >= 400 {
        logLevel = "warning"
    }
    if response.StatusCode >= 500 {
        logLevel = "error"
    }

    m.logger.Log(logLevel, "Request completed", map[string]interface{}{
        "method":     request.Method,
        "path":       request.Path,
        "status":     response.StatusCode,
        "duration":   duration.String(),
        "size":       len(response.Body),
    })

    return response
}
```

### 6. 错误日志记录

```go
// 错误日志记录器
type ErrorLogger struct {
    logger *log.Logger
}

func (l *ErrorLogger) LogError(err error, context map[string]interface{}) {
    // 获取错误堆栈
    stack := make([]string, 0)
    for _, frame := range err.StackTrace() {
        stack = append(stack, frame.String())
    }

    // 记录错误详情
    l.logger.Error("Application error", map[string]interface{}{
        "error":   err.Error(),
        "type":    reflect.TypeOf(err).String(),
        "stack":   stack,
        "context": context,
        "time":    time.Now(),
    })
}

// 在服务中使用
func (s *UserService) CreateUser(data map[string]interface{}) (*Models.User, error) {
    user := &Models.User{
        Name:  data["name"].(string),
        Email: data["email"].(string),
    }

    if err := user.SetPassword(data["password"].(string)); err != nil {
        s.errorLogger.LogError(err, map[string]interface{}{
            "action": "create_user",
            "email":  data["email"],
        })
        return nil, err
    }

    if err := s.db.Create(user).Error; err != nil {
        s.errorLogger.LogError(err, map[string]interface{}{
            "action": "save_user",
            "user_id": user.ID,
        })
        return nil, err
    }

    return user, nil
}
```

### 7. 性能日志

```go
// 性能监控日志
type PerformanceLogger struct {
    logger *log.Logger
}

func (l *PerformanceLogger) LogDatabaseQuery(sql string, duration time.Duration, rows int) {
    if duration > time.Second {
        l.logger.Warning("Slow database query", map[string]interface{}{
            "sql":      sql,
            "duration": duration.String(),
            "rows":     rows,
        })
    } else {
        l.logger.Debug("Database query", map[string]interface{}{
            "sql":      sql,
            "duration": duration.String(),
            "rows":     rows,
        })
    }
}

func (l *PerformanceLogger) LogCacheOperation(operation string, key string, duration time.Duration, hit bool) {
    l.logger.Debug("Cache operation", map[string]interface{}{
        "operation": operation,
        "key":       key,
        "duration":  duration.String(),
        "hit":       hit,
    })
}

// 在数据库中间件中使用
type DatabaseLogMiddleware struct {
    performanceLogger *PerformanceLogger
}

func (m *DatabaseLogMiddleware) BeforeQuery(sql string, args []interface{}) {
    // 记录查询开始
}

func (m *DatabaseLogMiddleware) AfterQuery(sql string, args []interface{}, duration time.Duration, rows int, err error) {
    if err != nil {
        // 记录错误
        return
    }

    m.performanceLogger.LogDatabaseQuery(sql, duration, rows)
}
```

### 8. 日志轮转

```go
// 日志轮转配置
type LogRotationConfig struct {
    MaxSize    int           `json:"max_size"`    // 最大文件大小 (MB)
    MaxAge     time.Duration `json:"max_age"`     // 最大保留时间
    MaxBackups int           `json:"max_backups"` // 最大备份数量
    Compress   bool          `json:"compress"`    // 是否压缩
}

// 配置日志轮转
func ConfigureLogRotation() {
    config.Set("log.rotation", LogRotationConfig{
        MaxSize:    100,           // 100MB
        MaxAge:     30 * 24 * time.Hour, // 30天
        MaxBackups: 10,            // 保留10个备份
        Compress:   true,          // 压缩旧日志
    })
}

// 手动触发日志轮转
func (c *AdminController) RotateLogs(request http.Request) http.Response {
    err := log.Rotate()
    if err != nil {
        return c.JsonError("Failed to rotate logs", 500)
    }

    return c.Json(map[string]string{
        "message": "Logs rotated successfully",
    })
}
```

### 9. 日志分析

```go
// 日志分析器
type LogAnalyzer struct {
    logger *log.Logger
}

func (a *LogAnalyzer) AnalyzeErrors(timeRange time.Duration) (*ErrorAnalysis, error) {
    // 分析错误日志
    errors, err := a.getErrors(timeRange)
    if err != nil {
        return nil, err
    }

    analysis := &ErrorAnalysis{
        TotalErrors: len(errors),
        ErrorTypes:  make(map[string]int),
        TopErrors:   make([]ErrorSummary, 0),
    }

    // 统计错误类型
    for _, err := range errors {
        analysis.ErrorTypes[err.Type]++
    }

    // 获取最常见的错误
    analysis.TopErrors = a.getTopErrors(errors, 10)

    return analysis, nil
}

func (a *LogAnalyzer) AnalyzePerformance(timeRange time.Duration) (*PerformanceAnalysis, error) {
    // 分析性能日志
    queries, err := a.getSlowQueries(timeRange)
    if err != nil {
        return nil, err
    }

    analysis := &PerformanceAnalysis{
        TotalQueries: len(queries),
        AvgDuration:  a.calculateAverageDuration(queries),
        SlowQueries:  queries,
    }

    return analysis, nil
}

// 在管理控制器中使用
func (c *AdminController) LogAnalytics(request http.Request) http.Response {
    analyzer := &LogAnalyzer{logger: log.GetLogger()}

    // 获取时间范围
    timeRange := 24 * time.Hour // 默认24小时

    // 分析错误
    errorAnalysis, err := analyzer.AnalyzeErrors(timeRange)
    if err != nil {
        return c.JsonError("Failed to analyze errors", 500)
    }

    // 分析性能
    performanceAnalysis, err := analyzer.AnalyzePerformance(timeRange)
    if err != nil {
        return c.JsonError("Failed to analyze performance", 500)
    }

    return c.Json(map[string]interface{}{
        "errors":      errorAnalysis,
        "performance": performanceAnalysis,
        "time_range":  timeRange.String(),
    })
}
```

### 10. 日志导出

```go
// 日志导出器
type LogExporter struct {
    logger *log.Logger
}

func (e *LogExporter) ExportToFile(startTime, endTime time.Time, format string) (string, error) {
    // 获取指定时间范围的日志
    logs, err := e.getLogs(startTime, endTime)
    if err != nil {
        return "", err
    }

    // 根据格式导出
    switch format {
    case "json":
        return e.exportToJSON(logs)
    case "csv":
        return e.exportToCSV(logs)
    case "xml":
        return e.exportToXML(logs)
    default:
        return "", errors.New("unsupported format")
    }
}

func (e *LogExporter) ExportToJSON(logs []LogEntry) (string, error) {
    filename := fmt.Sprintf("logs_%s.json", time.Now().Format("20060102_150405"))
    filepath := "storage/exports/" + filename

    data, err := json.MarshalIndent(logs, "", "  ")
    if err != nil {
        return "", err
    }

    err = ioutil.WriteFile(filepath, data, 0644)
    if err != nil {
        return "", err
    }

    return filepath, nil
}

// 在控制器中使用
func (c *AdminController) ExportLogs(request http.Request) http.Response {
    startTime := request.Query["start_time"]
    endTime := request.Query["end_time"]
    format := request.Query["format"]

    if format == "" {
        format = "json"
    }

    exporter := &LogExporter{logger: log.GetLogger()}

    filepath, err := exporter.ExportToFile(startTime, endTime, format)
    if err != nil {
        return c.JsonError("Failed to export logs", 500)
    }

    return c.Json(map[string]string{
        "filepath": filepath,
        "message":  "Logs exported successfully",
    })
}
```

## 📚 总结

Laravel-Go Framework 的日志系统提供了：

1. **多级别日志**: Debug、Info、Warning、Error、Critical
2. **多驱动支持**: 文件、Redis、数据库等
3. **结构化日志**: 支持上下文和元数据
4. **日志中间件**: 自动记录请求和错误
5. **性能监控**: 数据库查询和缓存操作日志
6. **日志轮转**: 自动管理日志文件大小和保留
7. **日志分析**: 错误统计和性能分析
8. **日志导出**: 支持多种格式的日志导出

通过合理使用日志系统，可以有效地监控和调试应用程序。
