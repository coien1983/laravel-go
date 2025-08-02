# 日志系统 API 参考

## 📋 概述

Laravel-Go Framework 的日志系统提供了强大而灵活的日志记录功能，支持多种日志级别、多种输出目标、结构化日志、日志轮转等特性。日志系统可以帮助开发者监控应用程序运行状态、调试问题、记录用户行为等。

## 🏗️ 核心概念

### 日志级别 (Log Levels)

- **Emergency**: 系统不可用
- **Alert**: 必须立即采取行动
- **Critical**: 严重错误
- **Error**: 错误信息
- **Warning**: 警告信息
- **Notice**: 一般性通知
- **Info**: 一般信息
- **Debug**: 调试信息

### 日志通道 (Channels)

- **File**: 文件日志
- **Console**: 控制台日志
- **Database**: 数据库日志
- **Redis**: Redis 日志
- **Syslog**: 系统日志
- **Custom**: 自定义通道

### 日志格式化器 (Formatters)

- **JSON**: JSON 格式
- **Text**: 文本格式
- **Custom**: 自定义格式

## 🔧 基础用法

### 1. 基本日志记录

```go
// 获取日志实例
logger := log.GetLogger()

// 记录不同级别的日志
logger.Emergency("System is down!")
logger.Alert("Database connection lost")
logger.Critical("Application crashed")
logger.Error("Failed to process request")
logger.Warning("High memory usage detected")
logger.Notice("User logged in")
logger.Info("Request processed successfully")
logger.Debug("Processing user data")

// 带上下文的日志
logger.Info("User action", map[string]interface{}{
    "user_id": 123,
    "action":  "login",
    "ip":      "192.168.1.1",
    "time":    time.Now(),
})
```

### 2. 在控制器中使用

```go
// app/Http/Controllers/UserController.go
package controllers

import (
    "laravel-go/framework/http"
    "laravel-go/framework/log"
)

type UserController struct {
    http.Controller
    logger log.Logger
}

func (c *UserController) Store(request http.Request) http.Response {
    // 记录请求开始
    c.logger.Info("Creating new user", map[string]interface{}{
        "email": request.Body["email"],
        "ip":    request.IP,
    })

    // 创建用户
    user, err := c.userService.CreateUser(request.Body)
    if err != nil {
        // 记录错误
        c.logger.Error("Failed to create user", map[string]interface{}{
            "error": err.Error(),
            "email": request.Body["email"],
        })
        return c.JsonError("Failed to create user", 500)
    }

    // 记录成功
    c.logger.Info("User created successfully", map[string]interface{}{
        "user_id": user.ID,
        "email":   user.Email,
    })

    return c.Json(user).Status(201)
}
```

### 3. 结构化日志

```go
// 使用结构化日志
logger := log.GetLogger()

// 记录用户活动
logger.Info("User activity", log.Fields{
    "user_id":    123,
    "action":     "purchase",
    "product_id": 456,
    "amount":     99.99,
    "currency":   "USD",
    "timestamp":  time.Now(),
})

// 记录系统指标
logger.Info("System metrics", log.Fields{
    "cpu_usage":    45.2,
    "memory_usage": 67.8,
    "disk_usage":   23.1,
    "active_users": 1250,
    "timestamp":    time.Now(),
})

// 记录错误详情
logger.Error("Database error", log.Fields{
    "error":       err.Error(),
    "query":       "SELECT * FROM users WHERE id = ?",
    "parameters":  []interface{}{123},
    "connection":  "mysql_primary",
    "timestamp":   time.Now(),
})
```

## 📚 API 参考

### Logger 接口

```go
type Logger interface {
    Emergency(message string, context ...map[string]interface{})
    Alert(message string, context ...map[string]interface{})
    Critical(message string, context ...map[string]interface{})
    Error(message string, context ...map[string]interface{})
    Warning(message string, context ...map[string]interface{})
    Notice(message string, context ...map[string]interface{})
    Info(message string, context ...map[string]interface{})
    Debug(message string, context ...map[string]interface{})

    Log(level Level, message string, context ...map[string]interface{})
    SetLevel(level Level)
    GetLevel() Level
    SetChannel(channel string)
    GetChannel() string
    AddContext(key string, value interface{})
    GetContext() map[string]interface{}
    ClearContext()
}
```

#### 方法说明

- `Emergency(message, context)`: 记录紧急日志
- `Alert(message, context)`: 记录警报日志
- `Critical(message, context)`: 记录严重错误日志
- `Error(message, context)`: 记录错误日志
- `Warning(message, context)`: 记录警告日志
- `Notice(message, context)`: 记录通知日志
- `Info(message, context)`: 记录信息日志
- `Debug(message, context)`: 记录调试日志
- `Log(level, message, context)`: 记录指定级别的日志
- `SetLevel(level)`: 设置日志级别
- `GetLevel()`: 获取当前日志级别
- `SetChannel(channel)`: 设置日志通道
- `GetChannel()`: 获取当前日志通道
- `AddContext(key, value)`: 添加上下文信息
- `GetContext()`: 获取上下文信息
- `ClearContext()`: 清除上下文信息

### Level 枚举

```go
type Level int

const (
    EmergencyLevel Level = iota
    AlertLevel
    CriticalLevel
    ErrorLevel
    WarningLevel
    NoticeLevel
    InfoLevel
    DebugLevel
)

func (l Level) String() string {
    switch l {
    case EmergencyLevel:
        return "emergency"
    case AlertLevel:
        return "alert"
    case CriticalLevel:
        return "critical"
    case ErrorLevel:
        return "error"
    case WarningLevel:
        return "warning"
    case NoticeLevel:
        return "notice"
    case InfoLevel:
        return "info"
    case DebugLevel:
        return "debug"
    default:
        return "unknown"
    }
}
```

### Channel 接口

```go
type Channel interface {
    Write(level Level, message string, context map[string]interface{}) error
    SetFormatter(formatter Formatter)
    GetFormatter() Formatter
    SetLevel(level Level)
    GetLevel() Level
    Close() error
}
```

#### 方法说明

- `Write(level, message, context)`: 写入日志
- `SetFormatter(formatter)`: 设置格式化器
- `GetFormatter()`: 获取格式化器
- `SetLevel(level)`: 设置日志级别
- `GetLevel()`: 获取日志级别
- `Close()`: 关闭通道

### Formatter 接口

```go
type Formatter interface {
    Format(level Level, message string, context map[string]interface{}) ([]byte, error)
    SetDateFormat(format string)
    GetDateFormat() string
}
```

#### 方法说明

- `Format(level, message, context)`: 格式化日志
- `SetDateFormat(format)`: 设置日期格式
- `GetDateFormat()`: 获取日期格式

## 🎯 高级功能

### 1. 多通道日志

```go
// 配置多个日志通道
config := log.Config{
    Default: "stack",
    Channels: map[string]log.ChannelConfig{
        "stack": {
            Driver: "stack",
            Channels: []string{"single", "daily"},
        },
        "single": {
            Driver: "file",
            Path:   "storage/logs/laravel.log",
            Level:  "debug",
        },
        "daily": {
            Driver: "file",
            Path:   "storage/logs/laravel-{date}.log",
            Level:  "info",
            Days:   14,
        },
        "database": {
            Driver: "database",
            Table:  "logs",
            Level:  "error",
        },
        "redis": {
            Driver: "redis",
            Key:    "laravel:logs",
            Level:  "debug",
        },
    },
}

// 初始化日志系统
log.SetConfig(config)

// 使用特定通道
logger := log.Channel("daily")
logger.Info("This will be written to daily log file")
```

### 2. 日志轮转

```go
// 配置日志轮转
dailyConfig := log.ChannelConfig{
    Driver: "file",
    Path:   "storage/logs/laravel-{date}.log",
    Level:  "info",
    Days:   14, // 保留14天的日志
    MaxSize: "100MB", // 单个文件最大100MB
    Compress: true, // 压缩旧日志
}

// 按大小轮转
sizeConfig := log.ChannelConfig{
    Driver: "file",
    Path:   "storage/logs/laravel.log",
    Level:  "debug",
    MaxSize: "50MB", // 单个文件最大50MB
    MaxFiles: 5, // 最多保留5个文件
    Compress: true,
}
```

### 3. 自定义格式化器

```go
// 自定义JSON格式化器
type CustomJSONFormatter struct {
    log.BaseFormatter
}

func (f *CustomJSONFormatter) Format(level log.Level, message string, context map[string]interface{}) ([]byte, error) {
    logEntry := map[string]interface{}{
        "timestamp": time.Now().Format(time.RFC3339),
        "level":     level.String(),
        "message":   message,
        "context":   context,
        "service":   "laravel-go",
        "version":   "1.0.0",
    }

    return json.Marshal(logEntry)
}

// 注册自定义格式化器
log.RegisterFormatter("custom_json", &CustomJSONFormatter{})
```

### 4. 日志中间件

```go
// 创建日志中间件
type LoggingMiddleware struct {
    http.Middleware
    logger log.Logger
}

func (m *LoggingMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
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
    m.logger.Info("Request completed", map[string]interface{}{
        "method":   request.Method,
        "path":     request.Path,
        "status":   response.StatusCode,
        "duration": duration.String(),
        "size":     len(response.Body),
    })

    return response
}
```

### 5. 错误日志记录

```go
// 记录详细的错误信息
func (c *UserController) Show(id string, request http.Request) http.Response {
    userID, err := strconv.Atoi(id)
    if err != nil {
        c.logger.Error("Invalid user ID", map[string]interface{}{
            "provided_id": id,
            "error":       err.Error(),
            "ip":          request.IP,
            "user_agent":  request.Headers["User-Agent"],
        })
        return c.JsonError("Invalid user ID", 400)
    }

    user, err := c.userService.GetUser(userID)
    if err != nil {
        c.logger.Error("Failed to get user", map[string]interface{}{
            "user_id":     userID,
            "error":       err.Error(),
            "error_type":  reflect.TypeOf(err).String(),
            "stack_trace": getStackTrace(),
        })
        return c.JsonError("User not found", 404)
    }

    return c.Json(user)
}

// 获取堆栈跟踪
func getStackTrace() string {
    var buf [4096]byte
    n := runtime.Stack(buf[:], false)
    return string(buf[:n])
}
```

## 🔧 配置选项

### 日志系统配置

```go
// config/log.go
package config

type LogConfig struct {
    // 默认通道
    Default string `json:"default"`

    // 通道配置
    Channels map[string]ChannelConfig `json:"channels"`

    // 全局日志级别
    Level string `json:"level"`

    // 是否启用调试模式
    Debug bool `json:"debug"`

    // 日志格式
    Format string `json:"format"`

    // 日期格式
    DateFormat string `json:"date_format"`

    // 时区
    Timezone string `json:"timezone"`
}

type ChannelConfig struct {
    // 驱动类型
    Driver string `json:"driver"`

    // 文件路径（文件驱动）
    Path string `json:"path"`

    // 日志级别
    Level string `json:"level"`

    // 保留天数（轮转）
    Days int `json:"days"`

    // 最大文件大小
    MaxSize string `json:"max_size"`

    // 最大文件数量
    MaxFiles int `json:"max_files"`

    // 是否压缩
    Compress bool `json:"compress"`

    // 数据库表名（数据库驱动）
    Table string `json:"table"`

    // Redis键名（Redis驱动）
    Key string `json:"key"`

    // 格式化器
    Formatter string `json:"formatter"`

    // 权限
    Permissions int `json:"permissions"`
}
```

### 配置示例

```go
// config/log.go
func GetLogConfig() *LogConfig {
    return &LogConfig{
        Default: "stack",
        Level:   "info",
        Debug:   false,
        Format:  "json",
        DateFormat: "2006-01-02 15:04:05",
        Timezone: "Asia/Shanghai",
        Channels: map[string]ChannelConfig{
            "stack": {
                Driver: "stack",
                Channels: []string{"single", "daily"},
            },
            "single": {
                Driver: "file",
                Path:   "storage/logs/laravel.log",
                Level:  "debug",
                Permissions: 0644,
            },
            "daily": {
                Driver: "file",
                Path:   "storage/logs/laravel-{date}.log",
                Level:  "info",
                Days:   14,
                MaxSize: "100MB",
                Compress: true,
                Permissions: 0644,
            },
            "database": {
                Driver: "database",
                Table:  "logs",
                Level:  "error",
            },
            "redis": {
                Driver: "redis",
                Key:    "laravel:logs",
                Level:  "debug",
            },
            "syslog": {
                Driver: "syslog",
                Level:  "warning",
            },
        },
    }
}
```

## 🚀 性能优化

### 1. 异步日志记录

```go
// 异步日志记录器
type AsyncLogger struct {
    log.Logger
    queue chan LogEntry
    done  chan bool
}

type LogEntry struct {
    Level   log.Level
    Message string
    Context map[string]interface{}
    Time    time.Time
}

func NewAsyncLogger(logger log.Logger, bufferSize int) *AsyncLogger {
    al := &AsyncLogger{
        Logger: logger,
        queue:  make(chan LogEntry, bufferSize),
        done:   make(chan bool),
    }

    go al.process()
    return al
}

func (al *AsyncLogger) process() {
    for entry := range al.queue {
        al.Logger.Log(entry.Level, entry.Message, entry.Context)
    }
    al.done <- true
}

func (al *AsyncLogger) Log(level log.Level, message string, context map[string]interface{}) {
    select {
    case al.queue <- LogEntry{
        Level:   level,
        Message: message,
        Context: context,
        Time:    time.Now(),
    }:
    default:
        // 队列满了，直接记录
        al.Logger.Log(level, message, context)
    }
}

func (al *AsyncLogger) Close() {
    close(al.queue)
    <-al.done
}
```

### 2. 日志缓冲

```go
// 缓冲日志记录器
type BufferedLogger struct {
    log.Logger
    buffer    []LogEntry
    bufferSize int
    mutex     sync.Mutex
    flushInterval time.Duration
}

func NewBufferedLogger(logger log.Logger, bufferSize int, flushInterval time.Duration) *BufferedLogger {
    bl := &BufferedLogger{
        Logger:        logger,
        buffer:        make([]LogEntry, 0, bufferSize),
        bufferSize:    bufferSize,
        flushInterval: flushInterval,
    }

    go bl.periodicFlush()
    return bl
}

func (bl *BufferedLogger) Log(level log.Level, message string, context map[string]interface{}) {
    bl.mutex.Lock()
    defer bl.mutex.Unlock()

    bl.buffer = append(bl.buffer, LogEntry{
        Level:   level,
        Message: message,
        Context: context,
        Time:    time.Now(),
    })

    if len(bl.buffer) >= bl.bufferSize {
        bl.flush()
    }
}

func (bl *BufferedLogger) flush() {
    for _, entry := range bl.buffer {
        bl.Logger.Log(entry.Level, entry.Message, entry.Context)
    }
    bl.buffer = bl.buffer[:0]
}

func (bl *BufferedLogger) periodicFlush() {
    ticker := time.NewTicker(bl.flushInterval)
    defer ticker.Stop()

    for range ticker.C {
        bl.mutex.Lock()
        if len(bl.buffer) > 0 {
            bl.flush()
        }
        bl.mutex.Unlock()
    }
}
```

### 3. 日志级别过滤

```go
// 级别过滤日志记录器
type LevelFilterLogger struct {
    log.Logger
    minLevel log.Level
}

func NewLevelFilterLogger(logger log.Logger, minLevel log.Level) *LevelFilterLogger {
    return &LevelFilterLogger{
        Logger:   logger,
        minLevel: minLevel,
    }
}

func (l *LevelFilterLogger) Log(level log.Level, message string, context map[string]interface{}) {
    if level >= l.minLevel {
        l.Logger.Log(level, message, context)
    }
}
```

## 🧪 测试

### 1. 日志测试

```go
// tests/log_test.go
package tests

import (
    "testing"
    "laravel-go/framework/log"
    "os"
    "strings"
)

func TestFileLogger(t *testing.T) {
    // 创建临时日志文件
    tempFile, err := os.CreateTemp("", "test_log_*.log")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tempFile.Name())
    defer tempFile.Close()

    // 创建文件日志记录器
    config := log.ChannelConfig{
        Driver: "file",
        Path:   tempFile.Name(),
        Level:  "debug",
    }

    channel, err := log.NewFileChannel(config)
    if err != nil {
        t.Fatal(err)
    }

    // 记录日志
    testMessage := "Test log message"
    err = channel.Write(log.InfoLevel, testMessage, map[string]interface{}{
        "test": true,
    })
    if err != nil {
        t.Fatal(err)
    }

    // 读取日志文件
    content, err := os.ReadFile(tempFile.Name())
    if err != nil {
        t.Fatal(err)
    }

    // 验证日志内容
    if !strings.Contains(string(content), testMessage) {
        t.Error("Log message not found in file")
    }
}

func TestLogLevels(t *testing.T) {
    logger := log.GetLogger()

    // 测试所有日志级别
    levels := []log.Level{
        log.EmergencyLevel,
        log.AlertLevel,
        log.CriticalLevel,
        log.ErrorLevel,
        log.WarningLevel,
        log.NoticeLevel,
        log.InfoLevel,
        log.DebugLevel,
    }

    for _, level := range levels {
        logger.Log(level, "Test message", map[string]interface{}{
            "level": level.String(),
        })
    }
}
```

### 2. 格式化器测试

```go
func TestJSONFormatter(t *testing.T) {
    formatter := log.NewJSONFormatter()

    message := "Test message"
    context := map[string]interface{}{
        "user_id": 123,
        "action":  "test",
    }

    formatted, err := formatter.Format(log.InfoLevel, message, context)
    if err != nil {
        t.Fatal(err)
    }

    // 验证JSON格式
    var result map[string]interface{}
    err = json.Unmarshal(formatted, &result)
    if err != nil {
        t.Fatal(err)
    }

    if result["message"] != message {
        t.Error("Message not found in formatted log")
    }

    if result["level"] != "info" {
        t.Error("Level not found in formatted log")
    }
}
```

## 🔍 调试和监控

### 1. 日志监控

```go
type LogMonitor struct {
    log.Logger
    metrics metrics.Collector
    counters map[log.Level]int64
    mutex    sync.RWMutex
}

func NewLogMonitor(logger log.Logger, metrics metrics.Collector) *LogMonitor {
    return &LogMonitor{
        Logger:  logger,
        metrics: metrics,
        counters: make(map[log.Level]int64),
    }
}

func (m *LogMonitor) Log(level log.Level, message string, context map[string]interface{}) {
    // 记录指标
    m.mutex.Lock()
    m.counters[level]++
    m.mutex.Unlock()

    m.metrics.Increment("logs.total", map[string]string{
        "level": level.String(),
    })

    // 记录错误率
    if level >= log.ErrorLevel {
        m.metrics.Increment("logs.errors")
    }

    // 调用原始记录器
    m.Logger.Log(level, message, context)
}

func (m *LogMonitor) GetStats() map[string]int64 {
    m.mutex.RLock()
    defer m.mutex.RUnlock()

    stats := make(map[string]int64)
    for level, count := range m.counters {
        stats[level.String()] = count
    }
    return stats
}
```

### 2. 日志分析

```go
type LogAnalyzer struct {
    patterns map[string]*regexp.Regexp
    alerts   []LogAlert
}

type LogAlert struct {
    Pattern string
    Count   int
    Level   log.Level
    Message string
}

func NewLogAnalyzer() *LogAnalyzer {
    return &LogAnalyzer{
        patterns: make(map[string]*regexp.Regexp),
        alerts:   make([]LogAlert, 0),
    }
}

func (la *LogAnalyzer) AddPattern(name, pattern string) error {
    regex, err := regexp.Compile(pattern)
    if err != nil {
        return err
    }

    la.patterns[name] = regex
    return nil
}

func (la *LogAnalyzer) Analyze(logs []string) []LogAlert {
    alerts := make([]LogAlert, 0)

    for _, logLine := range logs {
        for name, pattern := range la.patterns {
            if pattern.MatchString(logLine) {
                alerts = append(alerts, LogAlert{
                    Pattern: name,
                    Count:   1,
                    Level:   log.WarningLevel,
                    Message: logLine,
                })
            }
        }
    }

    return alerts
}
```

## 📝 最佳实践

### 1. 日志级别使用

```go
// 正确使用日志级别
logger := log.GetLogger()

// Emergency: 系统不可用
logger.Emergency("Database connection failed, system shutting down")

// Alert: 需要立即关注
logger.Alert("High CPU usage detected: 95%")

// Critical: 严重错误
logger.Critical("Application crashed due to memory overflow")

// Error: 错误但系统仍可运行
logger.Error("Failed to process user request", map[string]interface{}{
    "user_id": 123,
    "error":   err.Error(),
})

// Warning: 潜在问题
logger.Warning("Database query took longer than expected", map[string]interface{}{
    "query_time": "2.5s",
    "threshold":  "1s",
})

// Notice: 重要事件
logger.Notice("User registered", map[string]interface{}{
    "user_id": 123,
    "email":   "user@example.com",
})

// Info: 一般信息
logger.Info("Request processed", map[string]interface{}{
    "method": "POST",
    "path":   "/api/users",
    "status": 201,
})

// Debug: 调试信息
logger.Debug("Processing user data", map[string]interface{}{
    "user_id": 123,
    "step":    "validation",
})
```

### 2. 结构化日志

```go
// 使用结构化日志记录用户活动
func (c *UserController) Login(request http.Request) http.Response {
    email := request.Body["email"].(string)

    // 记录登录尝试
    c.logger.Info("Login attempt", log.Fields{
        "email":     email,
        "ip":        request.IP,
        "user_agent": request.Headers["User-Agent"],
        "timestamp": time.Now(),
    })

    user, err := c.authService.Authenticate(email, request.Body["password"].(string))
    if err != nil {
        // 记录登录失败
        c.logger.Warning("Login failed", log.Fields{
            "email":     email,
            "ip":        request.IP,
            "error":     err.Error(),
            "timestamp": time.Now(),
        })
        return c.JsonError("Invalid credentials", 401)
    }

    // 记录登录成功
    c.logger.Info("Login successful", log.Fields{
        "user_id":   user.ID,
        "email":     user.Email,
        "ip":        request.IP,
        "timestamp": time.Now(),
    })

    return c.Json(map[string]interface{}{
        "token": generateToken(user),
        "user":  user,
    })
}
```

### 3. 错误日志记录

```go
// 记录详细的错误信息
func (c *UserController) Update(id string, request http.Request) http.Response {
    userID, err := strconv.Atoi(id)
    if err != nil {
        c.logger.Error("Invalid user ID format", log.Fields{
            "provided_id": id,
            "error":       err.Error(),
            "error_type":  reflect.TypeOf(err).String(),
            "ip":          request.IP,
            "user_agent":  request.Headers["User-Agent"],
        })
        return c.JsonError("Invalid user ID", 400)
    }

    user, err := c.userService.GetUser(userID)
    if err != nil {
        c.logger.Error("User not found", log.Fields{
            "user_id":    userID,
            "error":      err.Error(),
            "error_type": reflect.TypeOf(err).String(),
            "stack_trace": getStackTrace(),
        })
        return c.JsonError("User not found", 404)
    }

    // 更新用户
    updatedUser, err := c.userService.UpdateUser(userID, request.Body)
    if err != nil {
        c.logger.Error("Failed to update user", log.Fields{
            "user_id":    userID,
            "error":      err.Error(),
            "error_type": reflect.TypeOf(err).String(),
            "changes":    request.Body,
        })
        return c.JsonError("Failed to update user", 500)
    }

    c.logger.Info("User updated successfully", log.Fields{
        "user_id": userID,
        "changes": request.Body,
    })

    return c.Json(updatedUser)
}
```

### 4. 性能日志

```go
// 记录性能指标
type PerformanceMiddleware struct {
    http.Middleware
    logger log.Logger
}

func (m *PerformanceMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    start := time.Now()

    // 记录请求开始
    m.logger.Debug("Request started", log.Fields{
        "method": request.Method,
        "path":   request.Path,
        "start":  start,
    })

    // 处理请求
    response := next(request)

    // 计算处理时间
    duration := time.Since(start)

    // 记录性能指标
    m.logger.Info("Request completed", log.Fields{
        "method":     request.Method,
        "path":       request.Path,
        "status":     response.StatusCode,
        "duration":   duration.String(),
        "duration_ms": duration.Milliseconds(),
        "size":       len(response.Body),
    })

    // 记录慢请求
    if duration > time.Second*2 {
        m.logger.Warning("Slow request detected", log.Fields{
            "method":     request.Method,
            "path":       request.Path,
            "duration":   duration.String(),
            "threshold":  "2s",
        })
    }

    return response
}
```

## 🚀 总结

日志系统是 Laravel-Go Framework 中重要的功能之一，它提供了：

1. **完整的日志功能**: 支持多种日志级别和输出目标
2. **结构化日志**: 便于日志分析和监控
3. **性能优化**: 提供异步和缓冲日志记录
4. **灵活配置**: 支持多种日志通道和格式化器
5. **监控调试**: 完整的日志监控和分析功能
6. **最佳实践**: 遵循日志记录的最佳实践

通过合理使用日志系统，可以有效地监控应用程序运行状态、调试问题、记录用户行为，提高应用程序的可维护性和可靠性。
