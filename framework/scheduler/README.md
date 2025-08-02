# Laravel-Go 定时器模块

## 概述

Laravel-Go 定时器模块提供完整的任务调度功能，支持多种调度表达式、任务持久化、监控统计和性能优化。模块设计遵循 Laravel-Go 框架的设计理念，提供简洁易用的 API 和强大的功能。

## 核心特性

### ✅ 已实现功能

- **多种调度表达式**: 支持标准 Cron 表达式、特殊表达式和简单时间格式
- **任务持久化**: 支持内存存储和数据库存储（可扩展）
- **任务生命周期管理**: 完整的任务创建、更新、删除、启用、禁用
- **任务执行监控**: 详细的执行统计、性能指标和错误追踪
- **调度器控制**: 启动、停止、暂停、恢复调度器
- **并发控制**: 支持任务并发执行和资源管理
- **错误处理**: 完善的错误处理和重试机制
- **监控统计**: 实时监控和性能统计
- **便捷 API**: 提供丰富的便捷方法和构建器模式

### 🚧 计划中功能

- **分布式调度**: 支持多实例环境下的任务调度
- **任务依赖**: 任务之间的依赖关系管理
- **条件任务**: 基于条件触发的任务
- **任务优先级**: 任务优先级管理
- **Web 管理界面**: 可视化的任务管理界面

## 快速开始

### 1. 基本使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "laravel-go/framework/scheduler"
)

func main() {
    // 初始化调度器
    store := scheduler.NewMemoryStore()
    scheduler.Init(store)

    // 创建任务处理器
    handler := scheduler.NewFuncHandler("hello", func(ctx context.Context) error {
        fmt.Println("Hello, Scheduler!", time.Now())
        return nil
    })

    // 创建任务
    task := scheduler.NewTask("hello-task", "Say hello every minute", "0 * * * * *", handler)

    // 添加任务到调度器
    if err := scheduler.AddTask(task); err != nil {
        log.Fatal(err)
    }

    // 启动调度器
    if err := scheduler.StartScheduler(); err != nil {
        log.Fatal(err)
    }

    // 保持程序运行
    select {}
}
```

### 2. 便捷方法

```go
// 每分钟执行
task1 := scheduler.Every(1, handler)

// 每小时执行
task2 := scheduler.EveryHour(handler)

// 每天执行
task3 := scheduler.EveryDay(handler)

// 每天指定时间执行
task4 := scheduler.Daily(9, 30, handler) // 每天 9:30

// 每周指定时间执行
task5 := scheduler.Weekly(time.Monday, 10, 0, handler) // 每周一 10:00

// 每月指定时间执行
task6 := scheduler.Monthly(1, 12, 0, handler) // 每月1号 12:00

// 自定义 Cron 表达式
task7 := scheduler.Cron("0 0 2 * * *", handler) // 每天凌晨2点

// 指定时间执行
task8 := scheduler.At("15:30", handler) // 每天 15:30
```

### 3. 任务构建器

```go
task := scheduler.NewTaskBuilder("backup", "Database backup", "0 2 * * *", handler).
    SetTimeout(5 * time.Minute).
    SetMaxRetries(3).
    SetRetryDelay(30 * time.Second).
    AddTag("type", "backup").
    AddTag("priority", "high").
    Build()
```

### 4. 调度器配置

```go
config := scheduler.NewSchedulerConfig().
    WithStore(store).
    WithMonitor(monitor).
    WithCheckInterval(2 * time.Second).
    WithMaxConcurrency(5).
    WithMetrics(true).
    WithLogging(true)

scheduler := config.Build()
```

## 调度表达式

### 1. 标准 Cron 表达式

支持 6-7 个字段的 Cron 表达式：

```
秒 分 时 日 月 周 [年]
* * * * * * *
```

示例：

- `0 * * * * *` - 每分钟执行
- `0 0 * * * *` - 每小时执行
- `0 0 0 * * *` - 每天执行
- `0 0 0 * * 0` - 每周执行
- `0 0 0 1 * *` - 每月执行
- `0 0 2 * * *` - 每天凌晨 2 点执行

### 2. 特殊表达式

- `@yearly` / `@annually` - 每年执行
- `@monthly` - 每月执行
- `@weekly` - 每周执行
- `@daily` / `@midnight` - 每天执行
- `@hourly` - 每小时执行

### 3. 简单时间格式

- `15:04` - 每天指定时间
- `15:04:05` - 每天指定时间（包含秒）
- `2006-01-02 15:04` - 指定日期和时间
- `01-02 15:04` - 每年指定日期和时间

### 4. 高级表达式

- `0 */5 * * * *` - 每 5 分钟执行
- `0 0 */2 * * *` - 每 2 小时执行
- `0 0 0 */2 * *` - 每 2 天执行
- `0 0 0 1,15 * *` - 每月 1 号和 15 号执行
- `0 0 0 1-5 * *` - 每月 1-5 号执行
- `0 0 0 * * 1-5` - 周一到周五执行

## 任务处理器

### 1. 函数处理器

```go
handler := scheduler.NewFuncHandler("my-task", func(ctx context.Context) error {
    // 任务逻辑
    return nil
})
```

### 2. 自定义处理器

```go
type MyHandler struct {
    name string
}

func (h *MyHandler) Handle(ctx context.Context) error {
    // 任务逻辑
    return nil
}

func (h *MyHandler) GetName() string {
    return h.name
}
```

## 存储

### 1. 内存存储

```go
store := scheduler.NewMemoryStore()
```

适用于开发和测试环境，数据不持久化。

### 2. 数据库存储

```go
// 需要实现具体的数据库连接
db := getDatabaseConnection()
store := scheduler.NewDatabaseStore(db, "scheduled_tasks")
```

适用于生产环境，数据持久化到数据库。

## 监控和统计

### 1. 任务统计

```go
// 获取任务统计
stats, err := scheduler.GetTaskStats(taskID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("任务运行次数: %d\n", stats.RunCount)
fmt.Printf("任务失败次数: %d\n", stats.FailedCount)
fmt.Printf("成功率: %.2f%%\n", stats.SuccessRate)
```

### 2. 调度器统计

```go
// 获取调度器统计
stats := scheduler.GetSchedulerStats()
fmt.Printf("总任务数: %d\n", stats.TotalTasks)
fmt.Printf("启用任务数: %d\n", stats.EnabledTasks)
fmt.Printf("总执行次数: %d\n", stats.TotalRuns)
fmt.Printf("成功率: %.2f%%\n", stats.SuccessRate)
```

### 3. 性能指标

```go
// 获取性能指标
metrics := scheduler.GetPerformanceMetrics()
fmt.Printf("内存使用: %d bytes\n", metrics.MemoryUsage)
fmt.Printf("CPU 使用率: %.2f%%\n", metrics.CPUUsage)
fmt.Printf("活跃协程数: %d\n", metrics.ActiveGoroutines)
fmt.Printf("吞吐量: %.2f 任务/秒\n", metrics.Throughput)
```

## 高级功能

### 1. 任务管理

```go
// 获取所有任务
tasks := scheduler.GetAllTasks()

// 获取启用的任务
enabledTasks := scheduler.GetEnabledTasks()

// 更新任务
task.SetTimeout(10 * time.Minute)
scheduler.UpdateTask(task)

// 启用/禁用任务
task.Enable()
task.Disable()

// 删除任务
scheduler.RemoveTask(taskID)
```

### 2. 调度器控制

```go
// 启动调度器
scheduler.StartScheduler()

// 停止调度器
scheduler.StopScheduler()

// 暂停调度器
scheduler.PauseScheduler()

// 恢复调度器
scheduler.ResumeScheduler()

// 立即运行任务
scheduler.RunTaskNow(taskID)

// 运行所有启用的任务
scheduler.RunAllTasks()
```

### 3. 任务标签

```go
// 添加标签
task.AddTag("environment", "production")
task.AddTag("priority", "high")

// 根据标签获取任务
tasks, err := store.GetByTags(map[string]string{
    "environment": "production",
    "priority": "high",
})
```

## 错误处理

### 1. 任务执行错误

```go
handler := scheduler.NewFuncHandler("error-task", func(ctx context.Context) error {
    // 模拟错误
    return fmt.Errorf("task execution failed")
})

task := scheduler.NewTask("error-task", "Task with error", "0 * * * * *", handler)
task.SetMaxRetries(3)
task.SetRetryDelay(30 * time.Second)
```

### 2. 调度器错误

```go
// 检查调度器状态
status := scheduler.GetSchedulerStatus()
switch status.Status {
case "running":
    fmt.Println("调度器正在运行")
case "paused":
    fmt.Println("调度器已暂停")
case "stopped":
    fmt.Println("调度器已停止")
}
```

## 最佳实践

### 1. 任务设计

- 任务应该是幂等的，多次执行不会产生副作用
- 任务应该处理超时和取消
- 任务应该记录详细的日志
- 任务应该返回有意义的错误信息

### 2. 调度设计

- 避免过于频繁的任务调度
- 合理设置任务超时时间
- 使用标签组织任务
- 监控任务执行情况

### 3. 性能优化

- 使用合适的存储后端
- 合理设置并发数
- 定期清理过期的监控数据
- 监控系统资源使用情况

## 示例应用

### 1. 数据备份任务

```go
func createBackupTask() *scheduler.DefaultTask {
    handler := scheduler.NewFuncHandler("backup", func(ctx context.Context) error {
        // 执行数据库备份
        return performDatabaseBackup(ctx)
    })

    return scheduler.NewTaskBuilder("database-backup", "Daily database backup", "0 2 * * *", handler).
        SetTimeout(30 * time.Minute).
        SetMaxRetries(2).
        AddTag("type", "backup").
        AddTag("priority", "high").
        Build()
}
```

### 2. 清理任务

```go
func createCleanupTask() *scheduler.DefaultTask {
    handler := scheduler.NewFuncHandler("cleanup", func(ctx context.Context) error {
        // 清理过期数据
        return cleanupExpiredData(ctx)
    })

    return scheduler.NewTaskBuilder("data-cleanup", "Clean expired data", "0 3 * * *", handler).
        SetTimeout(10 * time.Minute).
        AddTag("type", "cleanup").
        Build()
}
```

### 3. 报告生成任务

```go
func createReportTask() *scheduler.DefaultTask {
    handler := scheduler.NewFuncHandler("report", func(ctx context.Context) error {
        // 生成每日报告
        return generateDailyReport(ctx)
    })

    return scheduler.NewTaskBuilder("daily-report", "Generate daily report", "0 9 * * *", handler).
        SetTimeout(5 * time.Minute).
        AddTag("type", "report").
        Build()
}
```

## 总结

Laravel-Go 定时器模块提供了完整、灵活、易用的任务调度功能，支持多种调度表达式、任务持久化、监控统计和性能优化。通过简洁的 API 和丰富的功能，可以轻松实现各种定时任务需求。
