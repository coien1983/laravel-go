# Laravel-Go 队列系统

## 概述

Laravel-Go 队列系统提供统一的消息队列接口，支持多种队列驱动，从轻量级开发环境到企业级生产环境。系统包含完整的任务生命周期管理、工作进程、重试机制和监控功能。

## 核心特性

### ✅ 已实现功能

- **统一队列接口**: 支持多种队列驱动的统一接口
- **内存队列驱动**: 高性能内存队列，适用于开发测试
- **分布式队列支持**: 支持多节点集群，包括领导者选举和任务分发
- **Redis 集群支持**: 基于 Redis 的分布式队列实现
- **etcd 集群支持**: 基于 etcd 的分布式队列实现
- **Consul 集群支持**: 基于 Consul 的分布式队列实现
- **ZooKeeper 集群支持**: 基于 ZooKeeper 的分布式队列实现
- **任务序列化**: 完整的任务序列化和反序列化支持
- **延迟队列**: 支持延迟执行的任务
- **批量操作**: 批量推送和弹出任务
- **工作进程**: 完整的任务处理生命周期管理
- **分布式工作进程池**: 多节点多进程并发处理，支持负载均衡
- **重试机制**: 自动重试失败的任务
- **失败处理**: 完善的失败任务处理机制
- **统计监控**: 队列和工作进程的统计信息
- **任务属性**: 支持优先级、标签、超时等属性

### 🚧 计划中功能

- **数据库队列驱动**: 持久化队列支持
- **RabbitMQ 驱动**: 企业级消息队列
- **Kafka 驱动**: 高吞吐量流处理
- **RocketMQ 驱动**: 阿里云开源消息队列
- **ActiveMQ 驱动**: 传统企业消息中间件

## 快速开始

### 1. 基本使用

```go
package main

import (
    "context"
    "laravel-go/framework/queue"
)

func main() {
    // 初始化队列系统
    queue.Init()

    // 注册内存队列
    memoryQueue := queue.NewMemoryQueue()
    queue.QueueManager.Extend("memory", memoryQueue)
    queue.QueueManager.SetDefaultQueue("memory")

    // 推送任务
    job := queue.NewJob([]byte("Hello Queue!"), "default")
    err := queue.Push(job)
    if err != nil {
        panic(err)
    }

    // 弹出任务
    ctx := context.Background()
    poppedJob, err := queue.Pop(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Printf("处理任务: %s\n", string(poppedJob.GetPayload()))
}
```

### 2. 分布式队列使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "laravel-go/framework/queue"
)

func main() {
    // 创建Redis集群
    cluster, err := queue.NewRedisCluster(queue.RedisClusterConfig{
        Addr:   "localhost:6379",
        NodeID: "node-1",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer cluster.Close()

    // 创建分布式配置
    config := queue.DistributedConfig{
        NodeID:                "node-1",
        Cluster:               cluster,
        ElectionTimeout:       30 * time.Second,
        LockTimeout:           10 * time.Second,
        HeartbeatInterval:     5 * time.Second,
        EnableLeaderElection:  true,
        EnableJobDistribution: true,
        WorkerCount:           3,
        MaxConcurrency:        5,
    }

    // 创建分布式队列
    dq := queue.NewDistributedQueue(config)

    // 设置回调
    dq.SetOnCompleted(func(job queue.Job) {
        fmt.Printf("任务完成: %s\n", string(job.GetPayload()))
    })

    dq.SetOnFailed(func(job queue.Job, err error) {
        fmt.Printf("任务失败: %s - %v\n", string(job.GetPayload()), err)
    })

    // 启动分布式队列
    if err := dq.Start(); err != nil {
        log.Fatal(err)
    }

    // 推送任务
    job := queue.NewJob([]byte("Distributed Job!"), "default")
    dq.Push(job)

    // 保持运行
    select {}
}
```

### 2. 延迟队列

```go
// 创建延迟任务
job := queue.NewJob([]byte("延迟任务"), "default")
job.SetDelay(5 * time.Second)

// 推送延迟任务
err := queue.Push(job)
if err != nil {
    panic(err)
}

// 延迟任务会在指定时间后可用
```

### 3. 工作进程

```go
// 创建工作进程
worker := queue.NewWorker(memoryQueue, "default")

// 设置回调
worker.SetOnCompleted(func(job queue.Job) {
    fmt.Printf("任务完成: %s\n", string(job.GetPayload()))
})

worker.SetOnFailed(func(job queue.Job, err error) {
    fmt.Printf("任务失败: %s - %v\n", string(job.GetPayload()), err)
})

// 启动工作进程
err := worker.Start()
if err != nil {
    panic(err)
}

// 停止工作进程
defer worker.Stop()
```

### 4. 工作进程池

```go
// 创建工作进程池
pool := queue.NewWorkerPool(memoryQueue, "default", 3)

// 启动工作进程池
err := pool.Start()
if err != nil {
    panic(err)
}

// 获取统计信息
stats, err := pool.GetStats()
if err != nil {
    panic(err)
}

// 停止工作进程池
defer pool.Stop()
```

### 5. 分布式工作进程池

```go
// 获取分布式队列的工作进程池
workerPool := dq.GetWorkerPool()

// 设置回调
workerPool.SetOnCompleted(func(job queue.Job) {
    fmt.Printf("分布式任务完成: %s\n", string(job.GetPayload()))
})

workerPool.SetOnFailed(func(job queue.Job, err error) {
    fmt.Printf("分布式任务失败: %s - %v\n", string(job.GetPayload()), err)
})

// 获取统计信息
stats := workerPool.GetStats()
fmt.Printf("工作进程池状态: %s, 总工作进程: %d, 活跃: %d\n",
    stats.Status, stats.TotalWorkers, stats.ActiveWorkers)
```

### 6. 批量操作

```go
// 批量推送任务
jobs := []queue.Job{
    queue.NewJob([]byte("任务1"), "default"),
    queue.NewJob([]byte("任务2"), "default"),
    queue.NewJob([]byte("任务3"), "default"),
}

err := memoryQueue.PushBatch(jobs)
if err != nil {
    panic(err)
}

// 批量弹出任务
ctx := context.Background()
poppedJobs, err := memoryQueue.PopBatch(ctx, 2)
if err != nil {
    panic(err)
}
```

### 7. 任务属性

```go
// 创建高级任务
job := queue.NewJob([]byte("高级任务"), "default")

// 设置属性
job.SetPriority(10)                    // 优先级
job.SetMaxAttempts(5)                  // 最大尝试次数
job.SetTimeout(60 * time.Second)       // 超时时间
job.AddTag("type", "email")            // 添加标签
job.AddTag("priority", "high")

// 推送任务
err := queue.Push(job)
```

### 8. 队列统计

```go
// 获取队列统计
stats, err := queue.GetStats()
if err != nil {
    panic(err)
}

fmt.Printf("总任务数: %d\n", stats.TotalJobs)
fmt.Printf("待处理任务: %d\n", stats.PendingJobs)
fmt.Printf("保留任务: %d\n", stats.ReservedJobs)
fmt.Printf("失败任务: %d\n", stats.FailedJobs)
fmt.Printf("完成任务: %d\n", stats.CompletedJobs)
```

## 分布式队列

### 概述

分布式队列支持多节点集群，提供高可用性和可扩展性。主要特性包括：

- **领导者选举**: 自动选举领导者节点，确保任务分发的唯一性
- **分布式锁**: 防止任务重复处理
- **节点管理**: 自动注册和注销节点，监控节点状态
- **消息广播**: 节点间通信，同步任务状态
- **故障转移**: 领导者故障时自动重新选举

### 架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   节点 1        │    │   节点 2        │    │   节点 3        │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │ 领导者    │  │    │  │ 跟随者    │  │    │  │ 跟随者    │  │
│  │ (Leader)  │  │    │  │ (Follower)│  │    │  │ (Follower)│  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   集群协调器    │
                    │   (Redis/etcd)  │
                    └─────────────────┘
```

### 集群配置

#### Redis 集群

```go
// 创建Redis集群
cluster, err := queue.NewRedisCluster(queue.RedisClusterConfig{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
    NodeID:   "node-1",
})
if err != nil {
    log.Fatal(err)
}
defer cluster.Close()
```

#### etcd 集群

```go
// 创建etcd集群
cluster, err := queue.NewEtcdCluster(queue.EtcdClusterConfig{
    Endpoints: []string{"localhost:2379"},
    NodeID:    "node-1",
})
if err != nil {
    log.Fatal(err)
}
defer cluster.Close()
```

#### Consul 集群

```go
// 创建Consul集群
cluster, err := queue.NewConsulCluster(queue.ConsulClusterConfig{
    Address: "localhost:8500",
    NodeID:  "node-1",
})
if err != nil {
    log.Fatal(err)
}
defer cluster.Close()
```

#### ZooKeeper 集群

```go
// 创建ZooKeeper集群
cluster, err := queue.NewZookeeperCluster(queue.ZookeeperClusterConfig{
    Servers: []string{"localhost:2181"},
    NodeID:  "node-1",
})
if err != nil {
    log.Fatal(err)
}
defer cluster.Close()
```

### 分布式配置

```go
config := queue.DistributedConfig{
    NodeID:                 "node-1",           // 节点ID
    Cluster:                cluster,            // 集群实例
    ElectionTimeout:        30 * time.Second,   // 选举超时
    LockTimeout:            10 * time.Second,   // 锁超时
    HeartbeatInterval:      5 * time.Second,    // 心跳间隔
    EnableLeaderElection:   true,               // 启用领导者选举
    EnableJobDistribution:  true,               // 启用任务分发
    WorkerCount:            3,                  // 工作进程数
    MaxConcurrency:         5,                  // 最大并发数
}
```

### 分布式统计

```go
// 获取分布式统计
stats := dq.GetDistributedStats()
fmt.Printf("节点ID: %s\n", stats.NodeID)
fmt.Printf("是否为领导者: %t\n", stats.IsLeader)
fmt.Printf("总节点数: %d\n", stats.TotalNodes)
fmt.Printf("在线节点数: %d\n", stats.OnlineNodes)
fmt.Printf("领导者ID: %s\n", stats.LeaderID)

// 获取集群节点
nodes, err := dq.GetClusterNodes()
if err == nil {
    for _, node := range nodes {
        fmt.Printf("节点: %s, 状态: %s\n", node.ID, node.Status)
    }
}
```

### 最佳实践

1. **节点 ID**: 使用唯一且有意义的节点 ID，如 `web-server-1`, `worker-node-2`
2. **超时配置**: 根据网络延迟调整选举和锁超时时间
3. **工作进程数**: 根据 CPU 核心数和任务复杂度调整
4. **监控**: 定期检查集群状态和任务处理情况
5. **故障处理**: 实现优雅的故障转移和恢复机制

### 部署建议

#### Docker 部署

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o queue-worker examples/distributed_queue_demo/main.go
CMD ["./queue-worker"]
```

#### Docker Compose (Redis)

```yaml
version: "3.8"
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  queue-worker-1:
    build: .
    environment:
      - CLUSTER_TYPE=redis
      - NODE_ID=worker-1
    depends_on:
      - redis

  queue-worker-2:
    build: .
    environment:
      - CLUSTER_TYPE=redis
      - NODE_ID=worker-2
    depends_on:
      - redis

  queue-worker-3:
    build: .
    environment:
      - CLUSTER_TYPE=redis
      - NODE_ID=worker-3
    depends_on:
      - redis
```

#### Docker Compose (etcd)

```yaml
version: "3.8"
services:
  etcd:
    image: quay.io/coreos/etcd:v3.5.0
    ports:
      - "2379:2379"
    command: etcd --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379

  queue-worker-1:
    build: .
    environment:
      - CLUSTER_TYPE=etcd
      - NODE_ID=worker-1
    depends_on:
      - etcd

  queue-worker-2:
    build: .
    environment:
      - CLUSTER_TYPE=etcd
      - NODE_ID=worker-2
    depends_on:
      - etcd
```

#### Docker Compose (Consul)

```yaml
version: "3.8"
services:
  consul:
    image: consul:1.15
    ports:
      - "8500:8500"
    command: consul agent -server -bootstrap-expect=1 -ui -client=0.0.0.0

  queue-worker-1:
    build: .
    environment:
      - CLUSTER_TYPE=consul
      - NODE_ID=worker-1
    depends_on:
      - consul

  queue-worker-2:
    build: .
    environment:
      - CLUSTER_TYPE=consul
      - NODE_ID=worker-2
    depends_on:
      - consul
```

#### Docker Compose (ZooKeeper)

```yaml
version: "3.8"
services:
  zookeeper:
    image: zookeeper:3.8
    ports:
      - "2181:2181"

  queue-worker-1:
    build: .
    environment:
      - CLUSTER_TYPE=zookeeper
      - NODE_ID=worker-1
    depends_on:
      - zookeeper

  queue-worker-2:
    build: .
    environment:
      - CLUSTER_TYPE=zookeeper
      - NODE_ID=worker-2
    depends_on:
      - zookeeper
```

### 故障排除

#### Redis 集群

1. **连接失败**: 检查 Redis 服务是否正常运行
2. **选举失败**: 检查网络连接和超时配置
3. **任务丢失**: 检查分布式锁配置和任务序列化

#### etcd 集群

1. **连接失败**: 检查 etcd 服务是否正常运行
2. **租约过期**: 检查网络延迟和租约续期配置
3. **事务失败**: 检查 etcd 版本兼容性

#### Consul 集群

1. **连接失败**: 检查 Consul 服务是否正常运行
2. **会话过期**: 检查网络延迟和会话续期配置
3. **KV 操作失败**: 检查 Consul 权限配置

#### ZooKeeper 集群

1. **连接失败**: 检查 ZooKeeper 服务是否正常运行
2. **节点创建失败**: 检查路径权限和节点类型
3. **监听器失效**: 检查网络连接和事件处理

#### 通用问题

1. **性能问题**: 调整工作进程数和并发配置
2. **内存泄漏**: 检查资源清理和连接关闭
3. **网络分区**: 实现优雅的故障转移机制

## 核心接口

### Queue 接口

```go
type Queue interface {
    // 基础操作
    Push(job Job) error
    PushBatch(jobs []Job) error
    Pop(ctx context.Context) (Job, error)
    PopBatch(ctx context.Context, count int) ([]Job, error)
    Delete(job Job) error
    Release(job Job, delay time.Duration) error

    // 延迟队列
    Later(job Job, delay time.Duration) error
    LaterBatch(jobs []Job, delay time.Duration) error

    // 队列管理
    Size() (int, error)
    Clear() error
    Close() error

    // 监控和统计
    GetStats() (QueueStats, error)
}
```

### Job 接口

```go
type Job interface {
    // 基础信息
    GetID() string
    GetPayload() []byte
    GetQueue() string
    GetAttempts() int
    GetMaxAttempts() int
    GetDelay() time.Duration
    GetTimeout() time.Duration
    GetPriority() int
    GetTags() map[string]string
    GetCreatedAt() time.Time
    GetReservedAt() *time.Time
    GetAvailableAt() time.Time

    // 状态管理
    MarkAsReserved()
    MarkAsCompleted()
    MarkAsFailed(error)
    IncrementAttempts()

    // 序列化
    Serialize() ([]byte, error)
    Deserialize(data []byte) error
}
```

## 驱动实现

### 内存队列 (MemoryQueue)

内存队列是当前唯一实现的驱动，适用于开发测试环境。

**特点**:

- 高性能，无外部依赖
- 支持所有队列功能
- 重启后数据丢失
- 不支持分布式

**使用示例**:

```go
memoryQueue := queue.NewMemoryQueue()
queue.QueueManager.Extend("memory", memoryQueue)
```

## 配置示例

### 基础配置

```go
// 初始化队列管理器
queue.Init()

// 注册队列驱动
memoryQueue := queue.NewMemoryQueue()
queue.QueueManager.Extend("memory", memoryQueue)
queue.QueueManager.SetDefaultQueue("memory")

// 配置工作进程
worker := queue.NewWorker(memoryQueue, "default")
worker.SetTimeout(30 * time.Second)
worker.SetMaxAttempts(3)
```

### 多队列配置

```go
// 创建多个队列
queue1 := queue.NewMemoryQueue()
queue2 := queue.NewMemoryQueue()

queue.QueueManager.Extend("high", queue1)
queue.QueueManager.Extend("low", queue2)

// 推送到不同队列
queue.PushTo("high", highPriorityJob)
queue.PushTo("low", lowPriorityJob)
```

## 最佳实践

### 1. 任务设计

```go
// 定义具体的任务结构
type EmailJob struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// 序列化任务
func (j *EmailJob) ToJob() queue.Job {
    data, _ := json.Marshal(j)
    job := queue.NewJob(data, "emails")
    job.SetMaxAttempts(3)
    job.SetTimeout(30 * time.Second)
    return job
}

// 反序列化任务
func (j *EmailJob) FromJob(job queue.Job) error {
    return json.Unmarshal(job.GetPayload(), j)
}
```

### 2. 错误处理

```go
worker.SetOnFailed(func(job queue.Job, err error) {
    // 记录错误日志
    log.Printf("任务失败: %s - %v", job.GetID(), err)

    // 发送告警
    if job.GetAttempts() >= job.GetMaxAttempts() {
        sendAlert(job, err)
    }
})
```

### 3. 监控和统计

```go
// 定期获取统计信息
go func() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        stats, err := queue.GetStats()
        if err != nil {
            log.Printf("获取统计失败: %v", err)
            continue
        }

        // 发送监控指标
        sendMetrics(stats)
    }
}()
```

## 测试

运行队列系统测试：

```bash
go test ./framework/queue -v
```

## 示例程序

运行队列系统演示：

```bash
cd examples/queue_demo
go run main.go
```

## 依赖

当前实现仅依赖标准库和以下第三方包：

```bash
go get github.com/google/uuid
```

## 性能特性

### 内存队列性能

- **吞吐量**: 10,000+ QPS
- **延迟**: < 1ms
- **内存使用**: 低
- **并发支持**: 完全线程安全

### 工作进程性能

- **并发处理**: 支持多工作进程并发
- **负载均衡**: 自动负载均衡
- **故障恢复**: 自动重试和故障转移
- **资源管理**: 自动资源清理

## 限制和注意事项

### 当前限制

1. **数据持久化**: 内存队列重启后数据丢失（生产环境建议使用 Redis 等持久化队列）
2. **任务大小**: 建议任务载荷不超过 1MB
3. **集群依赖**: 分布式模式需要外部集群服务（Redis、etcd、Consul、ZooKeeper）

### 注意事项

1. **内存使用**: 大量任务可能占用较多内存
2. **并发限制**: 工作进程数量应根据系统资源调整
3. **网络延迟**: 分布式模式对网络延迟敏感，需要合理配置超时时间
4. **集群稳定性**: 确保集群服务的稳定性和高可用性

## 未来计划

### 短期计划 (1-2 周)

1. 实现数据库队列驱动
2. 完善错误处理和监控
3. 添加更多集群后端支持

### 中期计划 (1 个月)

1. 实现 RabbitMQ 驱动
2. 实现 Kafka 驱动
3. 添加更多企业级特性

### 长期计划 (3 个月)

1. 实现 RocketMQ 驱动
2. 实现 ActiveMQ 驱动
3. 添加更多监控和运维功能
4. 支持更多高级特性（任务优先级、任务依赖等）

## 贡献

欢迎提交 Issue 和 Pull Request 来改进队列系统。

## 许可证

本项目采用 MIT 许可证。
