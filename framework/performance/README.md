# Laravel-Go 性能监控系统

Laravel-Go 框架的性能监控系统提供了完整的性能指标收集、监控和优化功能，帮助开发者实时了解应用程序的性能状况并进行优化。

## 功能特性

### 1. 指标收集

- **计数器 (Counter)**: 统计事件发生的次数
- **仪表 (Gauge)**: 测量可增可减的数值
- **直方图 (Histogram)**: 统计数值分布情况
- **标签支持**: 为指标添加维度信息

### 2. 系统监控

- **CPU 监控**: CPU 使用率、核心数
- **内存监控**: 内存使用情况、Go 运行时内存
- **磁盘监控**: 磁盘使用情况、I/O 统计
- **网络监控**: 网络连接状态
- **Go 运行时**: 协程数量、垃圾回收统计

### 3. HTTP 监控

- **请求统计**: 请求数量、响应时间
- **错误监控**: 错误率、状态码分布
- **连接管理**: 活跃连接数
- **性能分析**: 请求大小、响应大小分布

### 4. 数据库监控

- **查询性能**: 监控数据库查询时间、类型分布
- **慢查询检测**: 自动识别和记录慢查询
- **连接池监控**: 监控连接池使用情况和性能
- **错误监控**: 跟踪数据库错误和失败查询
- **事务监控**: 监控事务执行时间和成功率

### 5. 缓存监控

- **命中率监控**: 实时监控缓存命中率
- **操作性能**: 监控 GET、SET、DELETE 操作性能
- **存储统计**: 监控缓存存储使用情况
- **错误监控**: 跟踪缓存操作错误
- **驱逐监控**: 监控缓存项驱逐情况

### 6. 告警系统

- **规则配置**: 支持多种告警规则和条件
- **多级告警**: 支持 info、warning、error、critical 级别
- **多种通知**: 支持日志、邮件、webhook 通知
- **自动恢复**: 支持告警自动恢复检测
- **历史记录**: 保存告警历史和处理状态

### 7. 性能报告

- **多种报告**: 支持摘要、详细、趋势、对比报告
- **智能建议**: 基于性能数据生成优化建议
- **多格式导出**: 支持 JSON 和文本格式导出
- **定期生成**: 支持定期自动生成报告
- **可视化数据**: 提供结构化的性能数据

### 8. 性能优化

- **连接池优化**: 数据库连接池配置建议
- **缓存优化**: 缓存命中率分析和建议
- **内存优化**: 内存使用优化和垃圾回收
- **并发优化**: 协程数量监控和优化建议
- **自动优化**: 定期自动执行优化分析

## 快速开始

### 1. 基础使用

```go
package main

import (
    "context"
    "time"
    "laravel-go/framework/performance"
)

func main() {
    // 创建性能监控器
    monitor := performance.NewPerformanceMonitor()

    // 启动监控
    ctx := context.Background()
    monitor.Start(ctx)
    defer monitor.Stop()

    // 创建指标
    counter := performance.NewCounter("requests_total", map[string]string{"service": "api"})
    gauge := performance.NewGauge("memory_usage", map[string]string{"unit": "bytes"})

    // 注册指标
    monitor.RegisterMetric(counter)
    monitor.RegisterMetric(gauge)

    // 更新指标
    counter.Increment(1)
    gauge.Set(1024.5)

    // 收集指标
    metrics := monitor.Collect()
    for _, metric := range metrics {
        fmt.Printf("%s: %v\n", metric.Name(), metric.Value())
    }
}
```

### 2. 系统监控

```go
// 创建系统监控器
systemMonitor := performance.NewSystemMonitor(monitor)
systemMonitor.Start(ctx)
defer systemMonitor.Stop()

// 获取系统指标
systemMetrics := systemMonitor.GetSystemMetrics()
for name, metric := range systemMetrics {
    fmt.Printf("%s: %v\n", name, metric)
}
```

### 3. HTTP 监控

```go
// 创建 HTTP 监控器
httpMonitor := performance.NewHTTPMonitor(monitor)

// 记录请求
httpMonitor.RecordRequest("GET", "/api/users", 150)

// 记录响应
httpMonitor.RecordResponse("GET", "/api/users", 200, 1024, 50*time.Millisecond)

// 记录错误
httpMonitor.RecordError("GET", "/api/error")
```

### 4. 性能优化

```go
// 创建性能优化器
optimizer := performance.NewPerformanceOptimizer(monitor)

// 执行所有优化
ctx := context.Background()
results, err := optimizer.Optimize(ctx)
if err != nil {
    log.Fatal(err)
}

// 查看优化结果
for _, result := range results {
    fmt.Printf("%s: %s (改进: %.1f%%)\n",
        result.Type, result.Message, result.Improvement)
}

// 执行特定类型优化
memoryResult, err := optimizer.OptimizeByType(ctx, performance.OptimizationTypeMemory)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("内存优化: %s\n", memoryResult.Message)
```

### 5. 数据库监控

```go
// 创建数据库监控器
dbMonitor := performance.NewDatabaseMonitor(monitor, 100*time.Millisecond)

// 记录查询
dbMonitor.RecordQuery("SELECT * FROM users", 50*time.Millisecond, true, nil)

// 记录事务
dbMonitor.RecordTransaction(200*time.Millisecond, true)

// 更新连接池状态
dbMonitor.UpdateConnectionPool(5, 10, 15)

// 获取慢查询
slowQueries := dbMonitor.GetSlowQueries()
for _, query := range slowQueries {
    fmt.Printf("慢查询: %s, 耗时: %v\n", query.SQL, query.Duration)
}
```

### 6. 缓存监控

```go
// 创建缓存监控器
cacheMonitor := performance.NewCacheMonitor(monitor)

// 记录GET操作
cacheMonitor.RecordGet("user:1", 1*time.Microsecond, true, nil)

// 记录SET操作
cacheMonitor.RecordSet("user:1", 2*time.Microsecond, nil)

// 记录DELETE操作
cacheMonitor.RecordDelete("user:1", 1*time.Microsecond, nil)

// 更新存储指标
cacheMonitor.UpdateStorageMetrics(1000, 1024*1024*10)

// 获取命中率
hitRate := cacheMonitor.GetHitRate()
fmt.Printf("缓存命中率: %.2f%%\n", hitRate)
```

### 7. 告警系统

```go
// 创建告警系统
alertSystem := performance.NewAlertSystem(monitor)

// 添加告警规则
cpuRule := &performance.AlertRule{
    ID:          "cpu_high",
    Name:        "CPU使用率过高",
    Description: "CPU使用率超过80%",
    MetricName:  "cpu_usage",
    Condition:   ">",
    Threshold:   80.0,
    Level:       performance.AlertLevelWarning,
    Enabled:     true,
    Actions:     []string{"log", "email"},
}
alertSystem.AddRule(cpuRule)

// 启动告警系统
alertSystem.Start(ctx)
defer alertSystem.Stop()

// 获取活跃告警
activeAlerts := alertSystem.GetActiveAlerts()
for _, alert := range activeAlerts {
    fmt.Printf("告警: %s - %s\n", alert.Level, alert.Message)
}
```

### 8. 性能报告

```go
// 创建报告生成器
reportGenerator := performance.NewReportGenerator(monitor, httpMonitor, dbMonitor, cacheMonitor, alertSystem)

// 生成摘要报告
period := performance.ReportPeriod{
    Start:    time.Now().Add(-1 * time.Hour),
    End:      time.Now(),
    Duration: time.Hour,
}
report, err := reportGenerator.GenerateReport(performance.ReportTypeSummary, period)
if err != nil {
    log.Fatal(err)
}

// 导出报告
data, err := reportGenerator.ExportReport(report, "json")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("报告: %s\n", string(data))

// 查看优化建议
for _, rec := range report.Recommendations {
    fmt.Printf("建议: %s (优先级: %s, 预期改进: %.1f%%)\n",
        rec.Title, rec.Priority, rec.Impact)
}
```

### 9. 自动优化

```go
// 创建自动优化器
autoOptimizer := performance.NewAutoOptimizer(optimizer, 5*time.Minute)
autoOptimizer.Start(ctx)
defer autoOptimizer.Stop()
```

## 指标类型详解

### Counter (计数器)

计数器用于统计事件发生的次数，只能递增。

```go
counter := performance.NewCounter("api_requests", map[string]string{"endpoint": "/users"})

// 增加计数
counter.Increment(1)
counter.Increment(5)

// 重置计数
counter.Reset()

// 获取值
value := counter.Value().(int64)
```

### Gauge (仪表)

仪表用于测量可增可减的数值，如内存使用量、连接数等。

```go
gauge := performance.NewGauge("memory_usage", map[string]string{"unit": "bytes"})

// 设置值
gauge.Set(1024.5)

// 增加值
gauge.Add(100)

// 获取值
value := gauge.Value().(float64)
```

### Histogram (直方图)

直方图用于统计数值的分布情况，如响应时间、请求大小等。

```go
buckets := []float64{10, 50, 100, 200, 500, 1000}
histogram := performance.NewHistogram("response_time", buckets, map[string]string{"unit": "ms"})

// 观察值
histogram.Observe(75.5)
histogram.Observe(150.2)

// 获取统计信息
value := histogram.Value().(map[string]interface{})
count := value["count"].(int64)
sum := value["sum"].(float64)
buckets := value["buckets"].(map[float64]int64)
```

## 系统监控指标

### CPU 指标

- `cpu_usage`: CPU 使用率 (%)
- `cpu_cores`: CPU 核心数

### 内存指标

- `memory_total`: 总内存 (bytes)
- `memory_used`: 已用内存 (bytes)
- `memory_available`: 可用内存 (bytes)
- `memory_usage_percent`: 内存使用率 (%)

### Go 运行时指标

- `go_heap_alloc`: 堆内存分配 (bytes)
- `go_heap_sys`: 堆内存系统 (bytes)
- `go_goroutines`: 协程数量

### 磁盘指标

- `disk_total_<path>`: 总空间 (bytes)
- `disk_used_<path>`: 已用空间 (bytes)
- `disk_free_<path>`: 可用空间 (bytes)
- `disk_usage_percent_<path>`: 使用率 (%)

## HTTP 监控指标

### 基础指标

- `http_requests_total`: 总请求数
- `http_responses_total`: 总响应数
- `http_errors_total`: 错误数
- `http_active_connections`: 活跃连接数

### 性能指标

- `http_response_time`: 响应时间分布 (直方图)
- `http_request_size`: 请求大小分布 (直方图)
- `http_response_size`: 响应大小分布 (直方图)

## 性能优化类型

### 1. 连接池优化 (OptimizationTypeConnectionPool)

- 分析数据库连接池使用率
- 提供连接池大小调整建议
- 监控连接池性能指标

### 2. 缓存优化 (OptimizationTypeCache)

- 分析缓存命中率
- 提供缓存策略优化建议
- 监控缓存性能指标

### 3. 内存优化 (OptimizationTypeMemory)

- 分析内存使用情况
- 提供内存优化建议
- 自动执行垃圾回收

### 4. 并发优化 (OptimizationTypeConcurrency)

- 监控协程数量
- 检测协程泄漏
- 提供并发处理优化建议

## 监控服务器

性能监控系统提供了内置的 HTTP 服务器，用于查看监控指标和系统状态。

```go
// 启动监控服务器
port := ":8088"
http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
    metrics := monitor.GetAllMetrics()
    // 返回指标数据
})

http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
    // 返回系统状态
})

http.HandleFunc("/optimize", func(w http.ResponseWriter, r *http.Request) {
    results, _ := optimizer.Optimize(context.Background())
    // 返回优化结果
})

http.ListenAndServe(port, nil)
```

### 可用端点

- `/metrics`: 获取所有监控指标
- `/status`: 获取系统状态
- `/optimize`: 执行性能优化
- `/health`: 健康检查

## 最佳实践

### 1. 指标命名

- 使用有意义的指标名称
- 遵循命名约定：`<namespace>_<metric>_<unit>`
- 例如：`http_requests_total`, `db_query_duration_ms`

### 2. 标签使用

- 为指标添加维度信息
- 避免标签值过多，影响性能
- 使用一致的标签命名

### 3. 监控配置

- 根据应用负载调整监控频率
- 设置合理的指标保留时间
- 配置适当的告警阈值

### 4. 性能考虑

- 避免在高频路径中创建过多指标
- 使用批量操作减少锁竞争
- 定期清理过期指标

## 扩展开发

### 自定义指标收集器

```go
type CustomCollector struct {
    monitor performance.Monitor
}

func (cc *CustomCollector) Collect(monitor performance.Monitor) error {
    // 实现自定义指标收集逻辑
    return nil
}

func (cc *CustomCollector) Name() string {
    return "custom"
}

// 添加到系统监控器
systemMonitor.AddCollector(&CustomCollector{})
```

### 自定义优化器

```go
type CustomOptimizer struct {
    monitor performance.Monitor
}

func (co *CustomOptimizer) Optimize(ctx context.Context) (*performance.OptimizationResult, error) {
    // 实现自定义优化逻辑
    return &performance.OptimizationResult{
        Type:        "custom",
        Success:     true,
        Message:     "Custom optimization completed",
        Improvement: 10.0,
        Timestamp:   time.Now(),
    }, nil
}

func (co *CustomOptimizer) GetType() performance.OptimizationType {
    return "custom"
}

func (co *CustomOptimizer) GetDescription() string {
    return "Custom optimization"
}

// 添加到性能优化器
optimizer.AddOptimizer(&CustomOptimizer{})
```

## 故障排除

### 常见问题

1. **指标不更新**

   - 检查监控器是否已启动
   - 确认指标已正确注册
   - 验证指标更新逻辑

2. **系统指标获取失败**

   - 检查 gopsutil 依赖是否正确安装
   - 确认系统权限是否足够
   - 验证系统监控器配置

3. **性能影响**

   - 调整监控频率
   - 减少指标数量
   - 使用异步收集

4. **内存泄漏**
   - 定期清理过期指标
   - 限制指标历史记录
   - 监控指标存储大小

## 总结

Laravel-Go 性能监控系统提供了完整的性能监控和优化解决方案，帮助开发者：

- 实时监控应用程序性能
- 快速定位性能瓶颈
- 自动优化系统配置
- 提供详细的性能报告

通过合理使用监控系统，可以显著提升应用程序的性能和稳定性。
