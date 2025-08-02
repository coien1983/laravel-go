# 性能监控系统 API 参考

## 📋 概述

Laravel-Go Framework 的性能监控系统提供了全面的应用程序性能监控功能，包括 HTTP 请求监控、数据库查询监控、内存使用监控、CPU 使用监控、自定义指标收集等。性能监控系统帮助开发者识别性能瓶颈、优化应用程序性能、提供实时监控和告警。

## 🏗️ 核心概念

### 性能监控器 (Performance Monitor)

- 收集和存储性能指标
- 提供性能分析工具
- 支持自定义指标

### 指标收集器 (Metrics Collector)

- 收集各种性能指标
- 支持多种数据源
- 提供聚合和分析功能

### 性能分析器 (Performance Profiler)

- 分析性能瓶颈
- 生成性能报告
- 提供优化建议

## 🔧 基础用法

### 1. 基本性能监控

```go
// 创建性能监控器
monitor := performance.NewMonitor()

// 启动监控
monitor.Start()

// 记录 HTTP 请求性能
func (c *UserController) Index(request http.Request) http.Response {
    start := time.Now()

    // 处理请求
    users, err := c.userService.GetUsers()
    if err != nil {
        return c.JsonError("Failed to get users", 500)
    }

    // 记录请求性能
    duration := time.Since(start)
    monitor.RecordHTTPRequest("GET", "/users", duration, 200)

    return c.Json(users)
}

// 记录数据库查询性能
func (s *UserService) GetUsers() ([]*Models.User, error) {
    start := time.Now()

    var users []*Models.User
    err := s.db.Find(&users).Error

    // 记录数据库查询性能
    duration := time.Since(start)
    monitor.RecordDatabaseQuery("SELECT * FROM users", duration, err == nil)

    return users, err
}
```

### 2. 在中间件中使用

```go
// app/Http/Middleware/PerformanceMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "laravel-go/framework/performance"
    "time"
)

type PerformanceMiddleware struct {
    http.Middleware
    monitor *performance.Monitor
}

func (m *PerformanceMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    start := time.Now()

    // 处理请求
    response := next(request)

    // 计算处理时间
    duration := time.Since(start)

    // 记录性能指标
    m.monitor.RecordHTTPRequest(
        request.Method,
        request.Path,
        duration,
        response.StatusCode,
    )

    // 记录内存使用
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    m.monitor.RecordMemoryUsage(memStats.Alloc, memStats.TotalAlloc, memStats.Sys)

    return response
}
```

### 3. 自定义指标收集

```go
// 收集自定义业务指标
func (c *OrderController) ProcessOrder(request http.Request) http.Response {
    start := time.Now()

    // 处理订单
    order, err := c.orderService.CreateOrder(request.Body)
    if err != nil {
        // 记录失败指标
        c.monitor.IncrementCounter("orders.failed", map[string]string{
            "error_type": reflect.TypeOf(err).String(),
        })
        return c.JsonError("Failed to create order", 500)
    }

    // 记录成功指标
    c.monitor.IncrementCounter("orders.created", map[string]string{
        "status": "success",
    })

    // 记录订单金额
    c.monitor.RecordHistogram("orders.amount", order.TotalAmount, map[string]string{
        "currency": order.Currency,
    })

    // 记录处理时间
    duration := time.Since(start)
    c.monitor.RecordHistogram("orders.processing_time", duration.Seconds(), nil)

    return c.Json(order).Status(201)
}
```

## 📚 API 参考

### Monitor 接口

```go
type Monitor interface {
    Start() error
    Stop() error
    IsRunning() bool

    RecordHTTPRequest(method, path string, duration time.Duration, statusCode int)
    RecordDatabaseQuery(query string, duration time.Duration, success bool)
    RecordMemoryUsage(alloc, totalAlloc, sys uint64)
    RecordCPUUsage(usage float64)

    IncrementCounter(name string, labels map[string]string)
    DecrementCounter(name string, labels map[string]string)
    SetGauge(name string, value float64, labels map[string]string)
    RecordHistogram(name string, value float64, labels map[string]string)

    GetMetrics() map[string]interface{}
    GetMetricsByType(metricType string) []Metric
    GetMetricsByTimeRange(start, end time.Time) []Metric

    EnableProfiling(enabled bool)
    IsProfilingEnabled() bool
    GetProfilingData() *ProfilingData

    SetAlertRule(rule AlertRule)
    GetAlertRules() []AlertRule
    TriggerAlert(alert Alert)
}
```

#### 方法说明

- `Start()`: 启动监控
- `Stop()`: 停止监控
- `IsRunning()`: 检查监控是否运行
- `RecordHTTPRequest(method, path, duration, statusCode)`: 记录 HTTP 请求
- `RecordDatabaseQuery(query, duration, success)`: 记录数据库查询
- `RecordMemoryUsage(alloc, totalAlloc, sys)`: 记录内存使用
- `RecordCPUUsage(usage)`: 记录 CPU 使用
- `IncrementCounter(name, labels)`: 增加计数器
- `DecrementCounter(name, labels)`: 减少计数器
- `SetGauge(name, value, labels)`: 设置仪表
- `RecordHistogram(name, value, labels)`: 记录直方图
- `GetMetrics()`: 获取所有指标
- `GetMetricsByType(metricType)`: 按类型获取指标
- `GetMetricsByTimeRange(start, end)`: 按时间范围获取指标
- `EnableProfiling(enabled)`: 启用/禁用性能分析
- `IsProfilingEnabled()`: 检查性能分析是否启用
- `GetProfilingData()`: 获取性能分析数据
- `SetAlertRule(rule)`: 设置告警规则
- `GetAlertRules()`: 获取告警规则
- `TriggerAlert(alert)`: 触发告警

### Metric 结构体

```go
type Metric struct {
    Name      string                 `json:"name"`
    Type      string                 `json:"type"`
    Value     float64                `json:"value"`
    Labels    map[string]string      `json:"labels"`
    Timestamp time.Time              `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata"`
}
```

#### 字段说明

- `Name`: 指标名称
- `Type`: 指标类型
- `Value`: 指标值
- `Labels`: 标签
- `Timestamp`: 时间戳
- `Metadata`: 元数据

### AlertRule 结构体

```go
type AlertRule struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Metric      string            `json:"metric"`
    Condition   string            `json:"condition"`
    Threshold   float64           `json:"threshold"`
    Duration    time.Duration     `json:"duration"`
    Labels      map[string]string `json:"labels"`
    Severity    string            `json:"severity"`
    Message     string            `json:"message"`
    Enabled     bool              `json:"enabled"`
}
```

#### 字段说明

- `ID`: 规则 ID
- `Name`: 规则名称
- `Metric`: 监控指标
- `Condition`: 条件（>、<、>=、<=、==）
- `Threshold`: 阈值
- `Duration`: 持续时间
- `Labels`: 标签
- `Severity`: 严重程度
- `Message`: 告警消息
- `Enabled`: 是否启用

## 🎯 高级功能

### 1. HTTP 请求监控

```go
// 详细的 HTTP 请求监控
type HTTPMonitor struct {
    performance.Monitor
}

func (m *HTTPMonitor) RecordHTTPRequest(method, path string, duration time.Duration, statusCode int) {
    // 记录基本指标
    m.Monitor.RecordHTTPRequest(method, path, duration, statusCode)

    // 记录详细指标
    m.RecordHistogram("http.request.duration", duration.Seconds(), map[string]string{
        "method": method,
        "path":   path,
        "status": fmt.Sprintf("%d", statusCode),
    })

    // 记录状态码分布
    m.IncrementCounter("http.response.status", map[string]string{
        "status": fmt.Sprintf("%d", statusCode),
        "method": method,
    })

    // 记录慢请求
    if duration > time.Second*2 {
        m.IncrementCounter("http.slow_requests", map[string]string{
            "method": method,
            "path":   path,
        })
    }

    // 记录错误请求
    if statusCode >= 400 {
        m.IncrementCounter("http.errors", map[string]string{
            "status": fmt.Sprintf("%d", statusCode),
            "method": method,
        })
    }
}
```

### 2. 数据库查询监控

```go
// 数据库查询监控
type DatabaseMonitor struct {
    performance.Monitor
}

func (m *DatabaseMonitor) RecordDatabaseQuery(query string, duration time.Duration, success bool) {
    // 记录基本指标
    m.Monitor.RecordDatabaseQuery(query, duration, success)

    // 解析查询类型
    queryType := m.parseQueryType(query)

    // 记录查询类型分布
    m.IncrementCounter("database.queries", map[string]string{
        "type": queryType,
    })

    // 记录查询时间
    m.RecordHistogram("database.query.duration", duration.Seconds(), map[string]string{
        "type": queryType,
    })

    // 记录慢查询
    if duration > time.Millisecond*100 {
        m.IncrementCounter("database.slow_queries", map[string]string{
            "type": queryType,
        })
    }

    // 记录失败查询
    if !success {
        m.IncrementCounter("database.query.errors", map[string]string{
            "type": queryType,
        })
    }
}

func (m *DatabaseMonitor) parseQueryType(query string) string {
    query = strings.TrimSpace(strings.ToUpper(query))
    if strings.HasPrefix(query, "SELECT") {
        return "SELECT"
    } else if strings.HasPrefix(query, "INSERT") {
        return "INSERT"
    } else if strings.HasPrefix(query, "UPDATE") {
        return "UPDATE"
    } else if strings.HasPrefix(query, "DELETE") {
        return "DELETE"
    }
    return "OTHER"
}
```

### 3. 内存和 CPU 监控

```go
// 系统资源监控
type SystemMonitor struct {
    performance.Monitor
    lastCPUUsage float64
    lastCPUTime  time.Time
}

func (m *SystemMonitor) StartSystemMonitoring() {
    go func() {
        ticker := time.NewTicker(time.Second * 5)
        defer ticker.Stop()

        for range ticker.C {
            m.recordSystemMetrics()
        }
    }()
}

func (m *SystemMonitor) recordSystemMetrics() {
    // 记录内存使用
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    m.RecordMemoryUsage(memStats.Alloc, memStats.TotalAlloc, memStats.Sys)

    // 记录内存指标
    m.SetGauge("memory.alloc", float64(memStats.Alloc), nil)
    m.SetGauge("memory.total_alloc", float64(memStats.TotalAlloc), nil)
    m.SetGauge("memory.sys", float64(memStats.Sys), nil)
    m.SetGauge("memory.num_gc", float64(memStats.NumGC), nil)

    // 记录 CPU 使用
    cpuUsage := m.calculateCPUUsage()
    m.RecordCPUUsage(cpuUsage)
    m.SetGauge("cpu.usage", cpuUsage, nil)

    // 记录 Goroutine 数量
    m.SetGauge("goroutines.count", float64(runtime.NumGoroutine()), nil)
}

func (m *SystemMonitor) calculateCPUUsage() float64 {
    now := time.Now()
    var usage runtime.MemStats
    runtime.ReadMemStats(&usage)

    if m.lastCPUTime.IsZero() {
        m.lastCPUUsage = 0
        m.lastCPUTime = now
        return 0
    }

    // 简化的 CPU 使用率计算
    duration := now.Sub(m.lastCPUTime).Seconds()
    cpuUsage := (float64(usage.PauseTotalNs) - m.lastCPUUsage) / duration / 1e9 * 100

    m.lastCPUUsage = float64(usage.PauseTotalNs)
    m.lastCPUTime = now

    return cpuUsage
}
```

### 4. 自定义业务指标

```go
// 业务指标监控
type BusinessMonitor struct {
    performance.Monitor
}

func (m *BusinessMonitor) RecordUserRegistration(user *Models.User) {
    // 记录用户注册
    m.IncrementCounter("users.registered", map[string]string{
        "source": user.RegistrationSource,
    })

    // 记录用户来源分布
    m.IncrementCounter("users.registration_source", map[string]string{
        "source": user.RegistrationSource,
    })
}

func (m *BusinessMonitor) RecordOrder(order *Models.Order) {
    // 记录订单创建
    m.IncrementCounter("orders.created", map[string]string{
        "status": order.Status,
    })

    // 记录订单金额
    m.RecordHistogram("orders.amount", order.TotalAmount, map[string]string{
        "currency": order.Currency,
        "status":   order.Status,
    })

    // 记录订单来源
    m.IncrementCounter("orders.source", map[string]string{
        "source": order.Source,
    })
}

func (m *BusinessMonitor) RecordPayment(payment *Models.Payment) {
    // 记录支付
    m.IncrementCounter("payments.processed", map[string]string{
        "status": payment.Status,
        "method": payment.Method,
    })

    // 记录支付金额
    m.RecordHistogram("payments.amount", payment.Amount, map[string]string{
        "currency": payment.Currency,
        "method":   payment.Method,
    })

    // 记录支付成功率
    if payment.Status == "success" {
        m.IncrementCounter("payments.success", map[string]string{
            "method": payment.Method,
        })
    } else {
        m.IncrementCounter("payments.failed", map[string]string{
            "method": payment.Method,
            "reason": payment.FailureReason,
        })
    }
}
```

### 5. 性能分析

```go
// 性能分析器
type PerformanceProfiler struct {
    performance.Monitor
    profiles map[string]*Profile
    mutex    sync.RWMutex
}

type Profile struct {
    Name      string                 `json:"name"`
    StartTime time.Time              `json:"start_time"`
    EndTime   time.Time              `json:"end_time"`
    Duration  time.Duration          `json:"duration"`
    Metrics   map[string]interface{} `json:"metrics"`
    Traces    []Trace                `json:"traces"`
}

type Trace struct {
    Function  string        `json:"function"`
    File      string        `json:"file"`
    Line      int           `json:"line"`
    Duration  time.Duration `json:"duration"`
    StartTime time.Time     `json:"start_time"`
}

func (p *PerformanceProfiler) StartProfile(name string) {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    p.profiles[name] = &Profile{
        Name:      name,
        StartTime: time.Now(),
        Metrics:   make(map[string]interface{}),
        Traces:    make([]Trace, 0),
    }
}

func (p *PerformanceProfiler) EndProfile(name string) *Profile {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    profile, exists := p.profiles[name]
    if !exists {
        return nil
    }

    profile.EndTime = time.Now()
    profile.Duration = profile.EndTime.Sub(profile.StartTime)

    // 记录性能指标
    p.RecordHistogram("profile.duration", profile.Duration.Seconds(), map[string]string{
        "profile": name,
    })

    return profile
}

func (p *PerformanceProfiler) AddTrace(profileName, function, file string, line int, duration time.Duration) {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    profile, exists := p.profiles[profileName]
    if !exists {
        return
    }

    trace := Trace{
        Function:  function,
        File:      file,
        Line:      line,
        Duration:  duration,
        StartTime: time.Now(),
    }

    profile.Traces = append(profile.Traces, trace)
}
```

## 🔧 配置选项

### 性能监控配置

```go
// config/performance.go
package config

type PerformanceConfig struct {
    // 基本配置
    Enabled bool `json:"enabled"`
    Interval time.Duration `json:"interval"`

    // 存储配置
    Storage StorageConfig `json:"storage"`

    // 指标配置
    Metrics MetricsConfig `json:"metrics"`

    // 告警配置
    Alerts AlertsConfig `json:"alerts"`

    // 性能分析配置
    Profiling ProfilingConfig `json:"profiling"`

    // 导出配置
    Export ExportConfig `json:"export"`
}

type StorageConfig struct {
    Driver string `json:"driver"`
    Path   string `json:"path"`
    TTL    time.Duration `json:"ttl"`
}

type MetricsConfig struct {
    HTTPRequests    bool `json:"http_requests"`
    DatabaseQueries bool `json:"database_queries"`
    MemoryUsage     bool `json:"memory_usage"`
    CPUUsage        bool `json:"cpu_usage"`
    CustomMetrics   bool `json:"custom_metrics"`
}

type AlertsConfig struct {
    Enabled bool `json:"enabled"`
    Rules   []AlertRule `json:"rules"`
    Notifications []NotificationConfig `json:"notifications"`
}

type ProfilingConfig struct {
    Enabled bool `json:"enabled"`
    Sampling float64 `json:"sampling"`
    MaxProfiles int `json:"max_profiles"`
}

type ExportConfig struct {
    Prometheus PrometheusConfig `json:"prometheus"`
    StatsD     StatsDConfig     `json:"statsd"`
    InfluxDB   InfluxDBConfig   `json:"influxdb"`
}

type PrometheusConfig struct {
    Enabled bool `json:"enabled"`
    Port    int  `json:"port"`
    Path    string `json:"path"`
}

type StatsDConfig struct {
    Enabled bool `json:"enabled"`
    Host    string `json:"host"`
    Port    int    `json:"port"`
    Prefix  string `json:"prefix"`
}

type InfluxDBConfig struct {
    Enabled  bool   `json:"enabled"`
    URL      string `json:"url"`
    Database string `json:"database"`
    Username string `json:"username"`
    Password string `json:"password"`
}
```

### 配置示例

```go
// config/performance.go
func GetPerformanceConfig() *PerformanceConfig {
    return &PerformanceConfig{
        Enabled:  true,
        Interval: time.Second * 5,
        Storage: StorageConfig{
            Driver: "memory",
            Path:   "storage/metrics",
            TTL:    time.Hour * 24, // 24 hours
        },
        Metrics: MetricsConfig{
            HTTPRequests:    true,
            DatabaseQueries: true,
            MemoryUsage:     true,
            CPUUsage:        true,
            CustomMetrics:   true,
        },
        Alerts: AlertsConfig{
            Enabled: true,
            Rules: []AlertRule{
                {
                    ID:        "slow_requests",
                    Name:      "Slow HTTP Requests",
                    Metric:    "http.request.duration",
                    Condition: ">",
                    Threshold: 2.0,
                    Duration:  time.Minute,
                    Severity:  "warning",
                    Message:   "HTTP requests are taking longer than 2 seconds",
                    Enabled:   true,
                },
                {
                    ID:        "high_memory",
                    Name:      "High Memory Usage",
                    Metric:    "memory.alloc",
                    Condition: ">",
                    Threshold: 100 * 1024 * 1024, // 100MB
                    Duration:  time.Minute * 5,
                    Severity:  "critical",
                    Message:   "Memory usage is above 100MB",
                    Enabled:   true,
                },
            },
        },
        Profiling: ProfilingConfig{
            Enabled:     true,
            Sampling:    0.1, // 10% sampling
            MaxProfiles: 100,
        },
        Export: ExportConfig{
            Prometheus: PrometheusConfig{
                Enabled: true,
                Port:    9090,
                Path:    "/metrics",
            },
            StatsD: StatsDConfig{
                Enabled: false,
                Host:    "localhost",
                Port:    8125,
                Prefix:  "laravel_go",
            },
            InfluxDB: InfluxDBConfig{
                Enabled:  false,
                URL:      "http://localhost:8086",
                Database: "laravel_go_metrics",
                Username: "",
                Password: "",
            },
        },
    }
}
```

## 🚀 性能优化

### 1. 指标缓存

```go
// 缓存性能指标
type CachedMonitor struct {
    performance.Monitor
    cache cache.Cache
}

func (m *CachedMonitor) GetMetrics() map[string]interface{} {
    cacheKey := "performance:metrics"

    if cached, exists := m.cache.Get(cacheKey); exists {
        return cached.(map[string]interface{})
    }

    metrics := m.Monitor.GetMetrics()
    m.cache.Set(cacheKey, metrics, time.Minute)

    return metrics
}
```

### 2. 批量指标处理

```go
// 批量处理指标
type BatchMonitor struct {
    performance.Monitor
    metrics []Metric
    mutex   sync.Mutex
    batchSize int
}

func (m *BatchMonitor) RecordMetric(metric Metric) {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    m.metrics = append(m.metrics, metric)

    if len(m.metrics) >= m.batchSize {
        m.flushMetrics()
    }
}

func (m *BatchMonitor) flushMetrics() {
    // 批量处理指标
    for _, metric := range m.metrics {
        m.Monitor.RecordHistogram(metric.Name, metric.Value, metric.Labels)
    }

    // 清空指标列表
    m.metrics = m.metrics[:0]
}
```

### 3. 异步指标收集

```go
// 异步指标收集
type AsyncMonitor struct {
    performance.Monitor
    queue chan Metric
    done  chan bool
}

func NewAsyncMonitor(monitor performance.Monitor, bufferSize int) *AsyncMonitor {
    am := &AsyncMonitor{
        Monitor: monitor,
        queue:   make(chan Metric, bufferSize),
        done:    make(chan bool),
    }

    go am.process()
    return am
}

func (am *AsyncMonitor) process() {
    for metric := range am.queue {
        am.Monitor.RecordHistogram(metric.Name, metric.Value, metric.Labels)
    }
    am.done <- true
}

func (am *AsyncMonitor) RecordMetric(metric Metric) {
    select {
    case am.queue <- metric:
    default:
        // 队列满了，直接处理
        am.Monitor.RecordHistogram(metric.Name, metric.Value, metric.Labels)
    }
}

func (am *AsyncMonitor) Close() {
    close(am.queue)
    <-am.done
}
```

## 🧪 测试

### 1. 性能监控测试

```go
// tests/performance_test.go
package tests

import (
    "testing"
    "time"
    "laravel-go/framework/performance"
)

func TestPerformanceMonitor(t *testing.T) {
    monitor := performance.NewMonitor()
    monitor.Start()
    defer monitor.Stop()

    // 测试 HTTP 请求监控
    monitor.RecordHTTPRequest("GET", "/users", time.Millisecond*100, 200)
    monitor.RecordHTTPRequest("POST", "/users", time.Millisecond*500, 201)

    // 测试数据库查询监控
    monitor.RecordDatabaseQuery("SELECT * FROM users", time.Millisecond*50, true)
    monitor.RecordDatabaseQuery("INSERT INTO users", time.Millisecond*200, true)

    // 测试内存监控
    monitor.RecordMemoryUsage(1024*1024, 2048*1024, 4096*1024)

    // 测试 CPU 监控
    monitor.RecordCPUUsage(25.5)

    // 获取指标
    metrics := monitor.GetMetrics()
    if len(metrics) == 0 {
        t.Error("Metrics should not be empty")
    }
}

func TestCustomMetrics(t *testing.T) {
    monitor := performance.NewMonitor()
    monitor.Start()
    defer monitor.Stop()

    // 测试计数器
    monitor.IncrementCounter("test.counter", map[string]string{"label": "value"})
    monitor.IncrementCounter("test.counter", map[string]string{"label": "value"})

    // 测试仪表
    monitor.SetGauge("test.gauge", 42.5, map[string]string{"label": "value"})

    // 测试直方图
    monitor.RecordHistogram("test.histogram", 1.5, map[string]string{"label": "value"})
    monitor.RecordHistogram("test.histogram", 2.5, map[string]string{"label": "value"})

    // 验证指标
    metrics := monitor.GetMetrics()
    if len(metrics) == 0 {
        t.Error("Custom metrics should be recorded")
    }
}
```

### 2. 告警规则测试

```go
func TestAlertRules(t *testing.T) {
    monitor := performance.NewMonitor()
    monitor.Start()
    defer monitor.Stop()

    // 设置告警规则
    rule := performance.AlertRule{
        ID:        "test_alert",
        Name:      "Test Alert",
        Metric:    "test.counter",
        Condition: ">",
        Threshold: 5,
        Duration:  time.Second,
        Severity:  "warning",
        Message:   "Test alert triggered",
        Enabled:   true,
    }

    monitor.SetAlertRule(rule)

    // 触发告警条件
    for i := 0; i < 6; i++ {
        monitor.IncrementCounter("test.counter", nil)
    }

    // 等待告警触发
    time.Sleep(time.Second * 2)

    // 验证告警
    rules := monitor.GetAlertRules()
    if len(rules) == 0 {
        t.Error("Alert rules should be set")
    }
}
```

## 🔍 调试和监控

### 1. 性能监控面板

```go
// 性能监控面板
type PerformanceDashboard struct {
    monitor *performance.Monitor
}

func (d *PerformanceDashboard) GetDashboardData() map[string]interface{} {
    metrics := d.monitor.GetMetrics()

    // 计算关键指标
    httpRequests := d.calculateHTTPMetrics(metrics)
    databaseQueries := d.calculateDatabaseMetrics(metrics)
    systemMetrics := d.calculateSystemMetrics(metrics)

    return map[string]interface{}{
        "http_requests":     httpRequests,
        "database_queries":  databaseQueries,
        "system_metrics":    systemMetrics,
        "alerts":           d.getActiveAlerts(),
        "profiles":         d.getRecentProfiles(),
    }
}

func (d *PerformanceDashboard) calculateHTTPMetrics(metrics map[string]interface{}) map[string]interface{} {
    // 计算 HTTP 请求指标
    return map[string]interface{}{
        "total_requests":    0,
        "avg_response_time": 0.0,
        "error_rate":        0.0,
        "slow_requests":     0,
    }
}

func (d *PerformanceDashboard) calculateDatabaseMetrics(metrics map[string]interface{}) map[string]interface{} {
    // 计算数据库查询指标
    return map[string]interface{}{
        "total_queries":     0,
        "avg_query_time":    0.0,
        "slow_queries":      0,
        "error_rate":        0.0,
    }
}

func (d *PerformanceDashboard) calculateSystemMetrics(metrics map[string]interface{}) map[string]interface{} {
    // 计算系统指标
    return map[string]interface{}{
        "memory_usage":      0.0,
        "cpu_usage":         0.0,
        "goroutines":        0,
        "gc_count":          0,
    }
}
```

### 2. 性能报告生成

```go
// 性能报告生成器
type PerformanceReporter struct {
    monitor *performance.Monitor
}

func (r *PerformanceReporter) GenerateReport(startTime, endTime time.Time) *PerformanceReport {
    metrics := r.monitor.GetMetricsByTimeRange(startTime, endTime)

    report := &PerformanceReport{
        Period:     fmt.Sprintf("%s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")),
        StartTime:  startTime,
        EndTime:    endTime,
        Metrics:    metrics,
        Summary:    r.generateSummary(metrics),
        Alerts:     r.getAlertsInRange(startTime, endTime),
        Recommendations: r.generateRecommendations(metrics),
    }

    return report
}

func (r *PerformanceReporter) generateSummary(metrics []performance.Metric) map[string]interface{} {
    // 生成性能摘要
    return map[string]interface{}{
        "total_requests":    0,
        "avg_response_time": 0.0,
        "error_rate":        0.0,
        "peak_memory":       0.0,
        "peak_cpu":          0.0,
    }
}

func (r *PerformanceReporter) generateRecommendations(metrics []performance.Metric) []string {
    // 生成优化建议
    recommendations := make([]string, 0)

    // 分析指标并生成建议
    // ...

    return recommendations
}
```

## 📝 最佳实践

### 1. 指标命名规范

```go
// 使用一致的指标命名规范
func (m *BusinessMonitor) recordUserMetrics(user *Models.User) {
    // 使用点分隔的命名方式
    m.IncrementCounter("users.registered.total", map[string]string{
        "source": user.RegistrationSource,
    })

    // 使用有意义的标签
    m.RecordHistogram("users.registration.duration", user.RegistrationDuration.Seconds(), map[string]string{
        "source": user.RegistrationSource,
        "method": user.RegistrationMethod,
    })

    // 避免过于细粒度的指标
    m.IncrementCounter("users.registration.success", map[string]string{
        "source": user.RegistrationSource,
    })
}
```

### 2. 性能监控策略

```go
// 合理的监控策略
func configurePerformanceMonitoring(monitor *performance.Monitor) {
    // 设置合理的采样率
    monitor.SetSamplingRate(0.1) // 10% 采样

    // 设置指标保留时间
    monitor.SetRetentionPeriod(time.Hour * 24 * 7) // 7 天

    // 设置告警阈值
    monitor.SetAlertRule(performance.AlertRule{
        Metric:    "http.request.duration",
        Condition: ">",
        Threshold: 2.0,
        Duration:  time.Minute,
        Severity:  "warning",
    })

    // 设置关键业务指标
    monitor.SetCriticalMetrics([]string{
        "users.registered.total",
        "orders.created.total",
        "payments.processed.total",
    })
}
```

### 3. 性能优化监控

```go
// 监控性能优化效果
func (m *PerformanceMonitor) monitorOptimization(optimizationName string, before, after func()) {
    // 记录优化前的性能
    before()
    beforeMetrics := m.GetMetrics()

    // 执行优化
    after()
    afterMetrics := m.GetMetrics()

    // 计算改进
    improvement := m.calculateImprovement(beforeMetrics, afterMetrics)

    // 记录优化结果
    m.RecordHistogram("optimization.improvement", improvement, map[string]string{
        "optimization": optimizationName,
    })

    // 如果改进不明显，记录警告
    if improvement < 0.1 {
        m.IncrementCounter("optimization.ineffective", map[string]string{
            "optimization": optimizationName,
        })
    }
}
```

### 4. 实时监控告警

```go
// 实时监控和告警
func (m *PerformanceMonitor) setupRealTimeMonitoring() {
    // 设置实时告警
    m.SetAlertRule(performance.AlertRule{
        ID:        "high_error_rate",
        Name:      "High Error Rate",
        Metric:    "http.errors.rate",
        Condition: ">",
        Threshold: 0.05, // 5% 错误率
        Duration:  time.Minute * 5,
        Severity:  "critical",
        Message:   "Error rate is above 5%",
    })

    // 设置性能告警
    m.SetAlertRule(performance.AlertRule{
        ID:        "slow_response_time",
        Name:      "Slow Response Time",
        Metric:    "http.request.duration",
        Condition: ">",
        Threshold: 1.0, // 1 秒
        Duration:  time.Minute * 2,
        Severity:  "warning",
        Message:   "Response time is above 1 second",
    })

    // 设置资源告警
    m.SetAlertRule(performance.AlertRule{
        ID:        "high_memory_usage",
        Name:      "High Memory Usage",
        Metric:    "memory.usage.percentage",
        Condition: ">",
        Threshold: 80.0, // 80% 内存使用
        Duration:  time.Minute * 5,
        Severity:  "warning",
        Message:   "Memory usage is above 80%",
    })
}
```

## 🚀 总结

性能监控系统是 Laravel-Go Framework 中重要的功能之一，它提供了：

1. **全面的性能监控**: HTTP 请求、数据库查询、系统资源等
2. **自定义指标**: 支持业务指标和自定义监控
3. **实时告警**: 基于阈值的告警机制
4. **性能分析**: 详细的性能分析和优化建议
5. **数据导出**: 支持多种监控系统集成
6. **最佳实践**: 遵循性能监控的最佳实践

通过合理使用性能监控系统，可以有效地识别性能瓶颈、优化应用程序性能、提供实时监控和告警，确保应用程序的高性能和稳定性。
