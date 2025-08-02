# Laravel-Go 定时器模块指南

## 概述

Laravel-Go 定时器模块提供了强大而灵活的任务调度功能，支持多种调度表达式、任务持久化、监控统计和性能优化。该模块完全自主开发，不依赖第三方库，为应用程序提供可靠的任务调度解决方案。

## 核心特性

- ✅ **多种调度表达式**: 支持标准 Cron 表达式、特殊表达式和简单时间格式
- ✅ **任务持久化**: 支持内存存储和数据库存储
- ✅ **监控统计**: 提供详细的执行统计和性能指标
- ✅ **错误处理**: 完善的错误处理和重试机制
- ✅ **并发控制**: 支持任务并发执行和资源管理
- ✅ **生命周期管理**: 完整的调度器启动、停止、暂停、恢复功能

## 快速开始

### 1. 初始化调度器

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
    // 初始化存储
    store := scheduler.NewMemoryStore()

    // 初始化调度器
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

### 2. 基本任务创建

```go
// 创建任务处理器
handler := scheduler.NewFuncHandler("my-task", func(ctx context.Context) error {
    // 任务逻辑
    fmt.Println("执行任务:", time.Now())
    return nil
})

// 创建任务
task := scheduler.NewTask(
    "my-task",                    // 任务名称
    "我的定时任务",                // 任务描述
    "0 */5 * * * *",             // 调度表达式（每5分钟）
    handler,                      // 任务处理器
)

// 配置任务
task.SetTimeout(30 * time.Second)  // 设置超时时间
task.SetMaxRetries(3)              // 设置最大重试次数
task.SetRetryDelay(5 * time.Second) // 设置重试延迟
task.AddTag("priority", "high")     // 添加标签

// 添加到调度器
scheduler.AddTask(task)
```

## 调度表达式

### 1. 标准 Cron 表达式

支持 6-7 字段的 Cron 表达式：

```
秒 分 时 日 月 周 [年]
```

示例：

- `0 * * * * *` - 每分钟执行
- `0 0 * * * *` - 每小时执行
- `0 0 0 * * *` - 每天执行
- `0 0 0 * * 0` - 每周执行
- `0 0 0 1 * *` - 每月执行
- `0 0 2 * * *` - 每天凌晨 2 点执行

### 2. 特殊表达式

支持常用的特殊表达式：

- `@yearly` / `@annually` - 每年执行一次
- `@monthly` - 每月执行一次
- `@weekly` - 每周执行一次
- `@daily` / `@midnight` - 每天执行一次
- `@hourly` - 每小时执行一次
- `@every` - 每隔指定时间执行

### 3. 简单时间格式

支持简单的时间格式：

- `15:30` - 每天 15:30 执行
- `15:30:00` - 每天 15:30:00 执行
- `12-25 15:30` - 每年 12 月 25 日 15:30 执行

### 4. 高级表达式

- `0 */5 * * * *` - 每 5 分钟执行
- `0 0 */2 * * *` - 每 2 小时执行
- `0 0 0 */2 * *` - 每 2 天执行
- `0 0 0 1,15 * *` - 每月 1 号和 15 号执行
- `0 0 0 1-5 * *` - 每月 1-5 号执行
- `0 0 0 * * 1-5` - 周一到周五执行

## 便捷方法

### 1. 常用调度方法

```go
// 每隔指定分钟执行
task := scheduler.Every(5, handler)  // 每5分钟

// 每小时执行
task := scheduler.EveryHour(handler)

// 每天执行
task := scheduler.EveryDay(handler)

// 指定时间执行
task := scheduler.Daily(9, 30, handler)  // 每天9:30

// 使用 Cron 表达式
task := scheduler.Cron("0 0 2 * * *", handler)  // 每天凌晨2点

// 指定时间执行
task := scheduler.At("15:30", handler)  // 每天15:30
```

### 2. 任务构建器

```go
task := scheduler.NewTaskBuilder("backup", "数据库备份", "0 2 * * *", handler).
    SetTimeout(30 * time.Minute).
    SetMaxRetries(3).
    SetRetryDelay(5 * time.Minute).
    AddTag("type", "backup").
    AddTag("priority", "high").
    Build()
```

## 任务管理

### 1. 添加和删除任务

```go
// 添加任务
err := scheduler.AddTask(task)

// 删除任务
err := scheduler.RemoveTask(taskID)

// 获取任务
task, err := scheduler.GetTask(taskID)

// 获取所有任务
tasks := scheduler.GetAllTasks()

// 获取启用的任务
enabledTasks := scheduler.GetEnabledTasks()
```

### 2. 任务状态管理

```go
// 启用任务
task.Enable()

// 禁用任务
task.Disable()

// 检查任务状态
if task.GetEnabled() {
    fmt.Println("任务已启用")
}

// 立即执行任务
err := scheduler.RunNow(taskID)

// 执行所有任务
err := scheduler.RunAll()
```

### 3. 调度器控制

```go
// 启动调度器
err := scheduler.StartScheduler()

// 停止调度器
err := scheduler.StopScheduler()

// 暂停调度器
err := scheduler.PauseScheduler()

// 恢复调度器
err := scheduler.ResumeScheduler()

// 获取调度器状态
status := scheduler.GetSchedulerStatus()
fmt.Printf("状态: %s, 任务数: %d\n", status.Status, status.TaskCount)
```

## 存储选项

### 1. 内存存储

```go
// 创建内存存储
store := scheduler.NewMemoryStore()

// 初始化调度器
scheduler.Init(store)
```

### 2. 数据库存储

```go
// 创建数据库存储（需要实现具体的数据库驱动）
dbStore := scheduler.NewDatabaseStore(db)

// 初始化调度器
scheduler.Init(dbStore)
```

## 监控和统计

### 1. 获取统计信息

```go
// 获取调度器统计
stats := scheduler.GetSchedulerStats()
fmt.Printf("总任务数: %d\n", stats.TotalTasks)
fmt.Printf("启用任务数: %d\n", stats.EnabledTasks)
fmt.Printf("总执行次数: %d\n", stats.TotalRuns)
fmt.Printf("失败次数: %d\n", stats.TotalFailed)
fmt.Printf("成功率: %.2f%%\n", stats.SuccessRate)

// 获取任务统计
taskStats, err := scheduler.GetTaskStats(taskID)
if err == nil {
    fmt.Printf("任务运行次数: %d\n", taskStats.RunCount)
    fmt.Printf("任务失败次数: %d\n", taskStats.FailedCount)
    fmt.Printf("任务成功率: %.2f%%\n", taskStats.SuccessRate)
}
```

### 2. 监控指标

```go
// 获取监控器
monitor := scheduler.GetMonitor()

// 获取调度器指标
metrics := monitor.GetSchedulerMetrics()
fmt.Printf("运行时间: %v\n", metrics.Uptime)
fmt.Printf("总执行次数: %d\n", metrics.TotalExecutions)
fmt.Printf("成功率: %.2f%%\n", metrics.SuccessRate)

// 获取任务指标
taskMetrics, err := monitor.GetTaskMetrics(taskID)
if err == nil {
    fmt.Printf("平均执行时间: %v\n", taskMetrics.AverageDuration)
    fmt.Printf("最小执行时间: %v\n", taskMetrics.MinDuration)
    fmt.Printf("最大执行时间: %v\n", taskMetrics.MaxDuration)
}

// 获取性能指标
perfMetrics := scheduler.GetPerformanceMetrics()
fmt.Printf("吞吐量: %.2f 任务/秒\n", perfMetrics.Throughput)
```

## 错误处理

### 1. 任务错误处理

```go
handler := scheduler.NewFuncHandler("error-handler", func(ctx context.Context) error {
    // 检查上下文是否被取消
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // 执行任务逻辑
    if err := doSomething(); err != nil {
        // 记录错误
        log.Printf("任务执行失败: %v", err)
        return err
    }

    return nil
})

// 配置重试机制
task := scheduler.NewTask("retry-task", "重试任务", "0 * * * * *", handler)
task.SetMaxRetries(3)
task.SetRetryDelay(10 * time.Second)
```

### 2. 调度器错误处理

```go
// 启动调度器
if err := scheduler.StartScheduler(); err != nil {
    switch err {
    case scheduler.ErrSchedulerAlreadyRunning:
        log.Println("调度器已在运行")
    case scheduler.ErrStoreNotInitialized:
        log.Println("存储未初始化")
    default:
        log.Fatalf("启动调度器失败: %v", err)
    }
}

// 添加任务
if err := scheduler.AddTask(task); err != nil {
    switch err {
    case scheduler.ErrTaskAlreadyExists:
        log.Println("任务已存在")
    case scheduler.ErrTaskNameRequired:
        log.Println("任务名称不能为空")
    default:
        log.Printf("添加任务失败: %v", err)
    }
}
```

## 最佳实践

### 1. 任务设计

```go
// 使用有意义的任务名称和描述
task := scheduler.NewTask(
    "user-cleanup",           // 清晰的名称
    "清理过期用户数据",        // 详细的描述
    "0 2 * * *",             // 每天凌晨2点执行
    cleanupHandler,
)

// 添加标签便于管理
task.AddTag("category", "maintenance")
task.AddTag("priority", "low")
task.AddTag("team", "backend")
```

### 2. 资源管理

```go
// 设置合理的超时时间
task.SetTimeout(5 * time.Minute)

// 配置重试机制
task.SetMaxRetries(3)
task.SetRetryDelay(30 * time.Second)

// 使用上下文控制任务执行
handler := scheduler.NewFuncHandler("ctx-task", func(ctx context.Context) error {
    // 定期检查上下文
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            // 执行工作
            if err := doWork(); err != nil {
                return err
            }
        }
    }
})
```

### 3. 监控和日志

```go
// 记录任务执行情况
handler := scheduler.NewFuncHandler("logged-task", func(ctx context.Context) error {
    start := time.Now()
    log.Printf("开始执行任务: %s", taskName)

    defer func() {
        duration := time.Since(start)
        log.Printf("任务执行完成: %s, 耗时: %v", taskName, duration)
    }()

    // 执行任务逻辑
    return doWork()
})

// 定期检查调度器状态
go func() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            status := scheduler.GetSchedulerStatus()
            stats := scheduler.GetSchedulerStats()

            log.Printf("调度器状态: %s, 任务数: %d, 成功率: %.2f%%",
                status.Status, stats.TotalTasks, stats.SuccessRate)
        }
    }
}()
```

### 4. 生产环境配置

```go
// 使用数据库存储确保任务持久化
dbStore := scheduler.NewDatabaseStore(db)
scheduler.Init(dbStore)

// 配置监控
monitor := scheduler.NewMonitor()
scheduler.SetMonitor(monitor)

// 启动调度器
if err := scheduler.StartScheduler(); err != nil {
    log.Fatalf("启动调度器失败: %v", err)
}

// 优雅关闭
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)

go func() {
    <-c
    log.Println("正在关闭调度器...")

    if err := scheduler.StopScheduler(); err != nil {
        log.Printf("停止调度器失败: %v", err)
    }

    log.Println("调度器已关闭")
    os.Exit(0)
}()
```

## 常见用例

### 1. 数据清理任务

```go
func createCleanupTask() *scheduler.DefaultTask {
    handler := scheduler.NewFuncHandler("cleanup", func(ctx context.Context) error {
        return cleanupExpiredData(ctx)
    })

    return scheduler.NewTaskBuilder("data-cleanup", "清理过期数据", "0 3 * * *", handler).
        SetTimeout(10 * time.Minute).
        SetMaxRetries(2).
        AddTag("type", "maintenance").
        Build()
}
```

### 2. 备份任务

```go
func createBackupTask() *scheduler.DefaultTask {
    handler := scheduler.NewFuncHandler("backup", func(ctx context.Context) error {
        return performDatabaseBackup(ctx)
    })

    return scheduler.NewTaskBuilder("database-backup", "数据库备份", "0 2 * * *", handler).
        SetTimeout(30 * time.Minute).
        SetMaxRetries(3).
        AddTag("type", "backup").
        AddTag("priority", "high").
        Build()
}
```

### 3. 报告生成任务

```go
func createReportTask() *scheduler.DefaultTask {
    handler := scheduler.NewFuncHandler("report", func(ctx context.Context) error {
        return generateDailyReport(ctx)
    })

    return scheduler.NewTaskBuilder("daily-report", "生成日报", "0 9 * * *", handler).
        SetTimeout(5 * time.Minute).
        SetMaxRetries(1).
        AddTag("type", "report").
        Build()
}
```

## 故障排除

### 1. 常见问题

**问题**: 任务不执行

- 检查调度表达式是否正确
- 确认任务已启用
- 检查调度器是否正在运行

**问题**: 任务执行失败

- 检查任务处理器逻辑
- 查看错误日志
- 确认超时设置是否合理

**问题**: 调度器无法启动

- 检查存储是否初始化
- 确认没有重复启动
- 查看错误信息

### 2. 调试技巧

```go
// 启用详细日志
log.SetLevel(log.DebugLevel)

// 检查任务配置
task, err := scheduler.GetTask(taskID)
if err == nil {
    fmt.Printf("任务配置: %+v\n", task)
    fmt.Printf("下次运行时间: %v\n", task.GetNextRunAt())
}

// 监控任务执行
monitor := scheduler.GetMonitor()
metrics, err := monitor.GetTaskMetrics(taskID)
if err == nil {
    fmt.Printf("任务指标: %+v\n", metrics)
}
```

## 分布式支持

### 概述

Laravel-Go 定时器模块支持分布式运行，可以在多个节点间协调任务执行，避免重复执行和单点故障。

### 分布式特性

- ✅ **领导者选举**: 自动选举领导者节点，只有领导者负责调度任务
- ✅ **分布式锁**: 确保同一任务不会被多个节点同时执行
- ✅ **节点管理**: 自动注册和注销节点，监控节点状态
- ✅ **消息广播**: 节点间通信，同步任务执行状态
- ✅ **故障转移**: 领导者故障时自动重新选举
- ✅ **负载均衡**: 支持任务分发到不同节点

### 分布式架构

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   节点1     │    │   节点2     │    │   节点3     │
│ (Leader)    │    │ (Follower)  │    │ (Follower)  │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌─────────────┐
                    │   Redis     │
                    │ (协调中心)   │
                    └─────────────┘
```

### 基本使用

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
    // 创建Redis集群
    cluster, err := scheduler.NewRedisCluster(scheduler.RedisClusterConfig{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
        NodeID:   "node-1",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer cluster.Close()

    // 创建分布式配置
    config := scheduler.DistributedConfig{
        NodeID:          "node-1",
        Cluster:         cluster,
        ElectionTimeout:  30 * time.Second,
        LockTimeout:     10 * time.Second,
        HeartbeatInterval: 5 * time.Second,
        EnableLeaderElection: true,
        EnableTaskDistribution: true,
    }

    // 创建存储
    store := scheduler.NewMemoryStore()

    // 创建分布式调度器
    ds := scheduler.NewDistributedScheduler(store, config)

    // 创建任务
    handler := scheduler.NewFuncHandler("distributed-task", func(ctx context.Context) error {
        fmt.Printf("分布式任务执行 (节点: %s)\n", ds.GetDistributedStats().NodeID)
        return nil
    })

    task := scheduler.NewTask("distributed-task", "分布式任务", "0 * * * * *", handler)
    ds.Add(task)

    // 启动分布式调度器
    if err := ds.Start(); err != nil {
        log.Fatal(err)
    }

    // 监控状态
    go func() {
        ticker := time.NewTicker(10 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                stats := ds.GetDistributedStats()
                fmt.Printf("节点: %s, 领导者: %t, 总节点数: %d\n",
                    stats.NodeID, stats.IsLeader, stats.TotalNodes)
            }
        }
    }()

    // 保持运行
    select {}
}
```

### 集群接口

```go
// Cluster 集群接口
type Cluster interface {
    // 节点管理
    Register(nodeID string, info NodeInfo) error
    Unregister(nodeID string) error
    GetNodes() ([]NodeInfo, error)

    // 分布式锁
    AcquireLock(key string, ttl time.Duration) (bool, error)
    ReleaseLock(key string) error

    // 选举
    StartElection(callback func(bool)) error
    StopElection() error

    // 消息广播
    Broadcast(msg ClusterMessage) error
    Subscribe(callback func(ClusterMessage)) error
}
```

### 节点信息

```go
// NodeInfo 节点信息
type NodeInfo struct {
    ID        string            `json:"id"`
    Address   string            `json:"address"`
    Port      int               `json:"port"`
    Status    string            `json:"status"` // online, offline, leader
    StartedAt time.Time         `json:"started_at"`
    LastSeen  time.Time         `json:"last_seen"`
    Metadata  map[string]string `json:"metadata"`
}
```

### 分布式统计

```go
// DistributedStats 分布式统计
type DistributedStats struct {
    NodeID      string `json:"node_id"`
    IsLeader    bool   `json:"is_leader"`
    TotalNodes  int    `json:"total_nodes"`
    OnlineNodes int    `json:"online_nodes"`
    LeaderID    string `json:"leader_id"`
}
```

### 分布式最佳实践

#### 1. 节点配置

```go
// 使用环境变量配置节点ID
nodeID := os.Getenv("NODE_ID")
if nodeID == "" {
    nodeID = fmt.Sprintf("node-%d", time.Now().Unix())
}

// 配置Redis连接
config := scheduler.RedisClusterConfig{
    Addr:     os.Getenv("REDIS_ADDR"),
    Password: os.Getenv("REDIS_PASSWORD"),
    DB:       0,
    NodeID:   nodeID,
}
```

#### 2. 故障处理

```go
// 创建分布式调度器
ds := scheduler.NewDistributedScheduler(store, config)

// 启动调度器
if err := ds.Start(); err != nil {
    log.Printf("启动分布式调度器失败: %v", err)
    // 降级到单节点模式
    runSingleNodeMode()
    return
}

// 监控领导者状态
go func() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if ds.IsLeader() {
                log.Println("当前节点为领导者")
            } else {
                log.Println("当前节点为跟随者")
            }
        }
    }
}()
```

#### 3. 任务分发

```go
// 创建适合分布式的任务
handler := scheduler.NewFuncHandler("distributed-backup", func(ctx context.Context) error {
    // 检查是否为领导者
    if !ds.IsLeader() {
        return fmt.Errorf("只有领导者可以执行此任务")
    }

    // 执行备份逻辑
    return performBackup(ctx)
})

task := scheduler.NewTask("backup", "数据库备份", "0 2 * * *", handler)
ds.Add(task)
```

#### 4. 监控和日志

```go
// 监控集群状态
func monitorCluster(ds *scheduler.DistributedScheduler) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            stats := ds.GetDistributedStats()
            nodes, _ := ds.GetClusterNodes()

            log.Printf("集群状态: 节点数=%d, 在线=%d, 领导者=%s",
                stats.TotalNodes, stats.OnlineNodes, stats.LeaderID)

            for _, node := range nodes {
                log.Printf("节点: %s, 状态: %s, 最后活跃: %s",
                    node.ID, node.Status, node.LastSeen.Format("15:04:05"))
            }
        }
    }
}
```

### 部署建议

#### 1. 生产环境配置

```bash
# 环境变量配置
export NODE_ID="prod-node-1"
export REDIS_ADDR="redis-cluster:6379"
export REDIS_PASSWORD="your-password"

# 启动多个节点
./scheduler-app &
./scheduler-app &
./scheduler-app &
```

#### 2. 容器化部署

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o scheduler-app examples/distributed_scheduler_demo/

CMD ["./scheduler-app"]
```

```yaml
# docker-compose.yml
version: "3.8"
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  scheduler-1:
    build: .
    environment:
      - NODE_ID=scheduler-1
      - REDIS_ADDR=redis:6379
    depends_on:
      - redis

  scheduler-2:
    build: .
    environment:
      - NODE_ID=scheduler-2
      - REDIS_ADDR=redis:6379
    depends_on:
      - redis

  scheduler-3:
    build: .
    environment:
      - NODE_ID=scheduler-3
      - REDIS_ADDR=redis:6379
    depends_on:
      - redis
```

### 故障排除

#### 1. 常见问题

**问题**: 节点无法成为领导者

- 检查 Redis 连接是否正常
- 确认节点 ID 是否唯一
- 查看选举超时配置

**问题**: 任务重复执行

- 检查分布式锁是否正常工作
- 确认任务 ID 是否唯一
- 查看节点间时间同步

**问题**: 节点无法通信

- 检查网络连接
- 确认 Redis pub/sub 功能
- 查看防火墙设置

#### 2. 调试技巧

```go
// 启用详细日志
log.SetLevel(log.DebugLevel)

// 检查集群状态
info, err := cluster.GetClusterInfo()
if err == nil {
    fmt.Printf("集群信息: %+v\n", info)
}

// 检查分布式锁
acquired, err := cluster.AcquireLock("test-lock", 5*time.Second)
if err == nil {
    fmt.Printf("锁获取结果: %t\n", acquired)
    if acquired {
        cluster.ReleaseLock("test-lock")
    }
}
```

## 总结

Laravel-Go 定时器模块提供了完整、灵活、易用的任务调度功能，支持多种调度表达式、任务持久化、监控统计和性能优化。通过简洁的 API 和丰富的功能，可以轻松实现各种定时任务需求。

**分布式支持**使得定时器模块可以在多节点环境中可靠运行，通过领导者选举、分布式锁和消息广播等机制，确保任务的一致性和高可用性。

通过合理使用调度表达式、配置任务参数、实施监控和遵循最佳实践，可以构建可靠、高效的定时任务系统。
