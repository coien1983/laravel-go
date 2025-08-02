# Laravel-Go 定时器模块 API 参考

## 概述

本文档提供了 Laravel-Go 定时器模块的完整 API 参考，包括所有接口、方法、结构体和常量的详细说明。

## 包导入

```go
import "laravel-go/framework/scheduler"
```

## 核心接口

### Task 接口

任务接口定义了定时任务的基本行为。

```go
type Task interface {
    // 基本信息
    GetID() string
    GetName() string
    GetDescription() string
    GetSchedule() string
    GetHandler() TaskHandler
    GetEnabled() bool
    GetCreatedAt() time.Time
    GetUpdatedAt() time.Time

    // 执行信息
    GetLastRunAt() *time.Time
    GetNextRunAt() *time.Time
    GetRunCount() int64
    GetFailedCount() int64
    GetLastError() string

    // 配置
    GetTimeout() time.Duration
    GetRetryCount() int
    GetRetryDelay() time.Duration
    GetMaxRetries() int
    GetTags() map[string]string

    // 状态管理
    Enable()
    Disable()
    UpdateNextRun()
    IncrementRunCount()
    IncrementFailedCount()

    // 序列化
    Serialize() ([]byte, error)
    Deserialize(data []byte) error
    Validate() error
    Clone() Task
    String() string
}
```

### TaskHandler 接口

任务处理器接口定义了任务执行逻辑。

```go
type TaskHandler interface {
    GetName() string
    Handle(ctx context.Context) error
}
```

### Scheduler 接口

调度器接口定义了任务调度的核心功能。

```go
type Scheduler interface {
    // 任务管理
    Add(task Task) error
    Remove(taskID string) error
    Update(task Task) error
    Get(taskID string) (Task, error)
    GetAll() []Task
    GetEnabled() []Task

    // 调度控制
    Start() error
    Stop() error
    Pause() error
    Resume() error

    // 任务执行
    RunNow(taskID string) error
    RunAll() error

    // 监控
    GetStatus() SchedulerStatus
    GetStats() SchedulerStats
    GetTaskStats(taskID string) (TaskStats, error)
}
```

### Store 接口

存储接口定义了任务持久化的抽象。

```go
type Store interface {
    // 基本操作
    Save(task Task) error
    Get(taskID string) (Task, error)
    GetAll() ([]Task, error)
    Delete(taskID string) error
    Clear() error

    // 批量操作
    SaveBatch(tasks []Task) error
    GetByTags(tags map[string]string) ([]Task, error)

    // 统计
    GetStats() (StoreStats, error)
    Close() error
}
```

### Monitor 接口

监控接口定义了任务执行监控功能。

```go
type Monitor interface {
    // 任务监控
    RecordTaskStart(taskID string)
    RecordTaskComplete(taskID string, duration time.Duration, err error)
    RecordTaskFailed(taskID string, duration time.Duration, err error)

    // 调度器监控
    RecordSchedulerStart()
    RecordSchedulerStop()
    RecordSchedulerPause()
    RecordSchedulerResume()

    // 统计查询
    GetTaskMetrics(taskID string) (TaskMetrics, error)
    GetSchedulerMetrics() SchedulerMetrics
    GetPerformanceMetrics() PerformanceMetrics

    // 清理
    Cleanup()
}
```

## 核心结构体

### DefaultTask

默认任务实现，提供了 Task 接口的完整实现。

```go
type DefaultTask struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Schedule    string            `json:"schedule"`
    Handler     TaskHandler       `json:"-"`
    Enabled     bool              `json:"enabled"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`

    // 执行信息
    LastRunAt   *time.Time        `json:"last_run_at"`
    NextRunAt   *time.Time        `json:"next_run_at"`
    RunCount    int64             `json:"run_count"`
    FailedCount int64             `json:"failed_count"`
    LastError   string            `json:"last_error"`

    // 配置
    Timeout     time.Duration     `json:"timeout"`
    RetryCount  int               `json:"retry_count"`
    RetryDelay  time.Duration     `json:"retry_delay"`
    MaxRetries  int               `json:"max_retries"`
    Tags        map[string]string `json:"tags"`
}
```

### DefaultScheduler

默认调度器实现，提供了 Scheduler 接口的完整实现。

```go
type DefaultScheduler struct {
    store      Store
    tasks      map[string]Task
    mu         sync.RWMutex
    status     SchedulerStatus
    stats      SchedulerStats
    stopChan   chan struct{}
    pauseChan  chan struct{}
    resumeChan chan struct{}
    ctx        context.Context
    cancel     context.CancelFunc
}
```

### SchedulerStatus

调度器状态信息。

```go
type SchedulerStatus struct {
    Status    string    `json:"status"` // running, paused, stopped
    StartedAt time.Time `json:"started_at"`
    TaskCount int       `json:"task_count"`
}
```

### SchedulerStats

调度器统计信息。

```go
type SchedulerStats struct {
    TotalTasks    int64     `json:"total_tasks"`
    EnabledTasks  int64     `json:"enabled_tasks"`
    DisabledTasks int64     `json:"disabled_tasks"`
    TotalRuns     int64     `json:"total_runs"`
    TotalFailed   int64     `json:"total_failed"`
    SuccessRate   float64   `json:"success_rate"`
    LastRunAt     time.Time `json:"last_run_at"`
    CreatedAt     time.Time `json:"created_at"`
}
```

### TaskStats

任务统计信息。

```go
type TaskStats struct {
    TaskID      string        `json:"task_id"`
    TaskName    string        `json:"task_name"`
    RunCount    int64         `json:"run_count"`
    FailedCount int64         `json:"failed_count"`
    SuccessRate float64       `json:"success_rate"`
    LastRunAt   time.Time     `json:"last_run_at"`
    NextRunAt   time.Time     `json:"next_run_at"`
    AverageTime time.Duration `json:"average_time"`
    LastError   string        `json:"last_error"`
}
```

### StoreStats

存储统计信息。

```go
type StoreStats struct {
    TotalTasks   int64     `json:"total_tasks"`
    EnabledTasks int64     `json:"enabled_tasks"`
    LastSync     time.Time `json:"last_sync"`
}
```

### TaskMetrics

任务监控指标。

```go
type TaskMetrics struct {
    TaskID          string        `json:"task_id"`
    TaskName        string        `json:"task_name"`
    TotalExecutions int64         `json:"total_executions"`
    SuccessfulRuns  int64         `json:"successful_runs"`
    FailedRuns      int64         `json:"failed_runs"`
    SuccessRate     float64       `json:"success_rate"`
    AverageDuration time.Duration `json:"average_duration"`
    MinDuration     time.Duration `json:"min_duration"`
    MaxDuration     time.Duration `json:"max_duration"`
    LastExecution   time.Time     `json:"last_execution"`
    LastError       string        `json:"last_error"`
    TotalDuration   time.Duration `json:"total_duration"`
}
```

### SchedulerMetrics

调度器监控指标。

```go
type SchedulerMetrics struct {
    TotalTasks      int64         `json:"total_tasks"`
    EnabledTasks    int64         `json:"enabled_tasks"`
    RunningTasks    int64         `json:"running_tasks"`
    TotalExecutions int64         `json:"total_executions"`
    SuccessfulRuns  int64         `json:"successful_runs"`
    FailedRuns      int64         `json:"failed_runs"`
    SuccessRate     float64       `json:"success_rate"`
    Uptime          time.Duration `json:"uptime"`
    StartedAt       time.Time     `json:"started_at"`
    LastActivity    time.Time     `json:"last_activity"`
}
```

### PerformanceMetrics

性能监控指标。

```go
type PerformanceMetrics struct {
    Throughput      float64 `json:"throughput"`       // 任务/秒
    AverageLatency  float64 `json:"average_latency"`  // 毫秒
    ErrorRate       float64 `json:"error_rate"`       // 百分比
    MemoryUsage     float64 `json:"memory_usage"`     // MB
    CPUUsage        float64 `json:"cpu_usage"`        // 百分比
}
```

## 构造函数

### NewTask

创建新的任务实例。

```go
func NewTask(name, description, schedule string, handler TaskHandler) *DefaultTask
```

**参数**:

- `name`: 任务名称
- `description`: 任务描述
- `schedule`: 调度表达式
- `handler`: 任务处理器

**返回值**: `*DefaultTask`

**示例**:

```go
handler := NewFuncHandler("test", func(ctx context.Context) error {
    return nil
})
task := NewTask("test-task", "Test task", "0 * * * * *", handler)
```

### NewScheduler

创建新的调度器实例。

```go
func NewScheduler(store Store) *DefaultScheduler
```

**参数**:

- `store`: 任务存储接口

**返回值**: `*DefaultScheduler`

**示例**:

```go
store := NewMemoryStore()
scheduler := NewScheduler(store)
```

### NewMemoryStore

创建内存存储实例。

```go
func NewMemoryStore() *MemoryStore
```

**返回值**: `*MemoryStore`

**示例**:

```go
store := NewMemoryStore()
```

### NewMonitor

创建监控器实例。

```go
func NewMonitor() *DefaultMonitor
```

**返回值**: `*DefaultMonitor`

**示例**:

```go
monitor := NewMonitor()
```

## 便捷方法

### Every

创建每隔指定分钟执行的任务。

```go
func Every(minutes int, handler TaskHandler) *DefaultTask
```

**参数**:

- `minutes`: 间隔分钟数
- `handler`: 任务处理器

**返回值**: `*DefaultTask`

**示例**:

```go
task := Every(5, handler) // 每5分钟执行
```

### EveryHour

创建每小时执行的任务。

```go
func EveryHour(handler TaskHandler) *DefaultTask
```

**参数**:

- `handler`: 任务处理器

**返回值**: `*DefaultTask`

**示例**:

```go
task := EveryHour(handler)
```

### EveryDay

创建每天执行的任务。

```go
func EveryDay(handler TaskHandler) *DefaultTask
```

**参数**:

- `handler`: 任务处理器

**返回值**: `*DefaultTask`

**示例**:

```go
task := EveryDay(handler)
```

### Daily

创建每天指定时间执行的任务。

```go
func Daily(hour, minute int, handler TaskHandler) *DefaultTask
```

**参数**:

- `hour`: 小时 (0-23)
- `minute`: 分钟 (0-59)
- `handler`: 任务处理器

**返回值**: `*DefaultTask`

**示例**:

```go
task := Daily(9, 30, handler) // 每天9:30执行
```

### Cron

使用 Cron 表达式创建任务。

```go
func Cron(schedule string, handler TaskHandler) *DefaultTask
```

**参数**:

- `schedule`: Cron 表达式
- `handler`: 任务处理器

**返回值**: `*DefaultTask`

**示例**:

```go
task := Cron("0 0 2 * * *", handler) // 每天凌晨2点执行
```

### At

创建每天指定时间执行的任务。

```go
func At(timeStr string, handler TaskHandler) *DefaultTask
```

**参数**:

- `timeStr`: 时间字符串 (格式: "HH:MM" 或 "HH:MM:SS")
- `handler`: 任务处理器

**返回值**: `*DefaultTask`

**示例**:

```go
task := At("15:30", handler) // 每天15:30执行
```

## 全局函数

### Init

初始化全局调度器。

```go
func Init(store Store)
```

**参数**:

- `store`: 任务存储接口

**示例**:

```go
store := NewMemoryStore()
Init(store)
```

### GetScheduler

获取全局调度器实例。

```go
func GetScheduler() Scheduler
```

**返回值**: `Scheduler`

**示例**:

```go
scheduler := GetScheduler()
```

### GetMonitor

获取全局监控器实例。

```go
func GetMonitor() Monitor
```

**返回值**: `Monitor`

**示例**:

```go
monitor := GetMonitor()
```

### AddTask

添加任务到全局调度器。

```go
func AddTask(task Task) error
```

**参数**:

- `task`: 任务实例

**返回值**: `error`

**示例**:

```go
err := AddTask(task)
```

### RemoveTask

从全局调度器移除任务。

```go
func RemoveTask(taskID string) error
```

**参数**:

- `taskID`: 任务 ID

**返回值**: `error`

**示例**:

```go
err := RemoveTask("task-id")
```

### GetTask

从全局调度器获取任务。

```go
func GetTask(taskID string) (Task, error)
```

**参数**:

- `taskID`: 任务 ID

**返回值**: `(Task, error)`

**示例**:

```go
task, err := GetTask("task-id")
```

### GetAllTasks

获取全局调度器的所有任务。

```go
func GetAllTasks() []Task
```

**返回值**: `[]Task`

**示例**:

```go
tasks := GetAllTasks()
```

### GetEnabledTasks

获取全局调度器的启用任务。

```go
func GetEnabledTasks() []Task
```

**返回值**: `[]Task`

**示例**:

```go
tasks := GetEnabledTasks()
```

### StartScheduler

启动全局调度器。

```go
func StartScheduler() error
```

**返回值**: `error`

**示例**:

```go
err := StartScheduler()
```

### StopScheduler

停止全局调度器。

```go
func StopScheduler() error
```

**返回值**: `error`

**示例**:

```go
err := StopScheduler()
```

### PauseScheduler

暂停全局调度器。

```go
func PauseScheduler() error
```

**返回值**: `error`

**示例**:

```go
err := PauseScheduler()
```

### ResumeScheduler

恢复全局调度器。

```go
func ResumeScheduler() error
```

**返回值**: `error`

**示例**:

```go
err := ResumeScheduler()
```

### GetSchedulerStatus

获取全局调度器状态。

```go
func GetSchedulerStatus() SchedulerStatus
```

**返回值**: `SchedulerStatus`

**示例**:

```go
status := GetSchedulerStatus()
```

### GetSchedulerStats

获取全局调度器统计。

```go
func GetSchedulerStats() SchedulerStats
```

**返回值**: `SchedulerStats`

**示例**:

```go
stats := GetSchedulerStats()
```

### GetTaskStats

获取任务统计。

```go
func GetTaskStats(taskID string) (TaskStats, error)
```

**参数**:

- `taskID`: 任务 ID

**返回值**: `(TaskStats, error)`

**示例**:

```go
stats, err := GetTaskStats("task-id")
```

### GetPerformanceMetrics

获取性能指标。

```go
func GetPerformanceMetrics() PerformanceMetrics
```

**返回值**: `PerformanceMetrics`

**示例**:

```go
metrics := GetPerformanceMetrics()
```

### RunNow

立即执行任务。

```go
func RunNow(taskID string) error
```

**参数**:

- `taskID`: 任务 ID

**返回值**: `error`

**示例**:

```go
err := RunNow("task-id")
```

### RunAll

执行所有启用任务。

```go
func RunAll() error
```

**返回值**: `error`

**示例**:

```go
err := RunAll()
```

## 任务构建器

### TaskBuilder

任务构建器提供了流式 API 来配置任务。

```go
type TaskBuilder struct {
    task *DefaultTask
}
```

### NewTaskBuilder

创建任务构建器。

```go
func NewTaskBuilder(name, description, schedule string, handler TaskHandler) *TaskBuilder
```

**参数**:

- `name`: 任务名称
- `description`: 任务描述
- `schedule`: 调度表达式
- `handler`: 任务处理器

**返回值**: `*TaskBuilder`

**示例**:

```go
builder := NewTaskBuilder("test", "Test task", "0 * * * * *", handler)
```

### SetTimeout

设置任务超时时间。

```go
func (b *TaskBuilder) SetTimeout(timeout time.Duration) *TaskBuilder
```

**参数**:

- `timeout`: 超时时间

**返回值**: `*TaskBuilder`

**示例**:

```go
builder.SetTimeout(5 * time.Minute)
```

### SetMaxRetries

设置最大重试次数。

```go
func (b *TaskBuilder) SetMaxRetries(maxRetries int) *TaskBuilder
```

**参数**:

- `maxRetries`: 最大重试次数

**返回值**: `*TaskBuilder`

**示例**:

```go
builder.SetMaxRetries(3)
```

### SetRetryDelay

设置重试延迟。

```go
func (b *TaskBuilder) SetRetryDelay(delay time.Duration) *TaskBuilder
```

**参数**:

- `delay`: 重试延迟

**返回值**: `*TaskBuilder`

**示例**:

```go
builder.SetRetryDelay(30 * time.Second)
```

### AddTag

添加标签。

```go
func (b *TaskBuilder) AddTag(key, value string) *TaskBuilder
```

**参数**:

- `key`: 标签键
- `value`: 标签值

**返回值**: `*TaskBuilder`

**示例**:

```go
builder.AddTag("priority", "high")
```

### Build

构建任务。

```go
func (b *TaskBuilder) Build() *DefaultTask
```

**返回值**: `*DefaultTask`

**示例**:

```go
task := builder.Build()
```

## 错误定义

### 预定义错误

```go
var (
    ErrTaskNotFound            = errors.New("task not found")
    ErrSchedulerAlreadyRunning = errors.New("scheduler is already running")
    ErrSchedulerAlreadyStopped = errors.New("scheduler is already stopped")
    ErrSchedulerNotRunning     = errors.New("scheduler is not running")
    ErrSchedulerNotPaused      = errors.New("scheduler is not paused")
    ErrInvalidSchedule         = errors.New("invalid schedule format")
    ErrTaskHandlerRequired     = errors.New("task handler is required")
    ErrTaskNameRequired        = errors.New("task name is required")
    ErrTaskIDRequired          = errors.New("task ID is required")
    ErrStoreNotInitialized     = errors.New("store is not initialized")
    ErrTaskExecutionTimeout    = errors.New("task execution timeout")
    ErrTaskExecutionFailed     = errors.New("task execution failed")
    ErrInvalidCronExpression   = errors.New("invalid cron expression")
    ErrInvalidTimeFormat       = errors.New("invalid time format")
    ErrTaskAlreadyExists       = errors.New("task already exists")
    ErrTaskDisabled            = errors.New("task is disabled")
    ErrTaskMaxRetriesExceeded  = errors.New("task max retries exceeded")
)
```

## 调度表达式解析

### ParseSchedule

解析调度表达式。

```go
func ParseSchedule(schedule string) (time.Time, error)
```

**参数**:

- `schedule`: 调度表达式

**返回值**: `(time.Time, error)`

**示例**:

```go
nextRun, err := ParseSchedule("0 * * * * *")
```

### CronExpression

Cron 表达式结构体。

```go
type CronExpression struct {
    Second     []int
    Minute     []int
    Hour       []int
    DayOfMonth []int
    Month      []int
    DayOfWeek  []int
    Year       []int
}
```

### NextRun

计算下次运行时间。

```go
func (c *CronExpression) NextRun(from time.Time) (time.Time, error)
```

**参数**:

- `from`: 起始时间

**返回值**: `(time.Time, error)`

**示例**:

```go
nextRun, err := cron.NextRun(time.Now())
```

## 配置结构体

### SchedulerConfig

调度器配置。

```go
type SchedulerConfig struct {
    Store          Store
    Monitor        Monitor
    CheckInterval  time.Duration
    MaxConcurrency int
    EnableMetrics  bool
    EnableLogging  bool
}
```

### NewSchedulerConfig

创建调度器配置。

```go
func NewSchedulerConfig() *SchedulerConfig
```

**返回值**: `*SchedulerConfig`

**示例**:

```go
config := NewSchedulerConfig()
```

### WithStore

设置存储。

```go
func (c *SchedulerConfig) WithStore(store Store) *SchedulerConfig
```

**参数**:

- `store`: 存储接口

**返回值**: `*SchedulerConfig`

**示例**:

```go
config.WithStore(store)
```

### WithMonitor

设置监控器。

```go
func (c *SchedulerConfig) WithMonitor(monitor Monitor) *SchedulerConfig
```

**参数**:

- `monitor`: 监控器接口

**返回值**: `*SchedulerConfig`

**示例**:

```go
config.WithMonitor(monitor)
```

### WithCheckInterval

设置检查间隔。

```go
func (c *SchedulerConfig) WithCheckInterval(interval time.Duration) *SchedulerConfig
```

**参数**:

- `interval`: 检查间隔

**返回值**: `*SchedulerConfig`

**示例**:

```go
config.WithCheckInterval(time.Second)
```

### WithMaxConcurrency

设置最大并发数。

```go
func (c *SchedulerConfig) WithMaxConcurrency(max int) *SchedulerConfig
```

**参数**:

- `max`: 最大并发数

**返回值**: `*SchedulerConfig`

**示例**:

```go
config.WithMaxConcurrency(10)
```

### WithMetrics

启用指标。

```go
func (c *SchedulerConfig) WithMetrics(enable bool) *SchedulerConfig
```

**参数**:

- `enable`: 是否启用

**返回值**: `*SchedulerConfig`

**示例**:

```go
config.WithMetrics(true)
```

### WithLogging

启用日志。

```go
func (c *SchedulerConfig) WithLogging(enable bool) *SchedulerConfig
```

**参数**:

- `enable`: 是否启用

**返回值**: `*SchedulerConfig`

**示例**:

```go
config.WithLogging(true)
```

### Build

构建调度器。

```go
func (c *SchedulerConfig) Build() *DefaultScheduler
```

**返回值**: `*DefaultScheduler`

**示例**:

```go
scheduler := config.Build()
```

## 完整示例

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
    // 初始化
    store := scheduler.NewMemoryStore()
    scheduler.Init(store)

    // 创建任务
    handler := scheduler.NewFuncHandler("hello", func(ctx context.Context) error {
        fmt.Printf("Hello at %s\n", time.Now().Format("15:04:05"))
        return nil
    })

    task := scheduler.NewTask("hello", "Say hello", "0 * * * * *", handler)
    scheduler.AddTask(task)

    // 启动调度器
    if err := scheduler.StartScheduler(); err != nil {
        log.Fatal(err)
    }

    // 保持运行
    select {}
}
```

### 高级配置

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
    // 使用配置构建器
    config := scheduler.NewSchedulerConfig().
        WithStore(scheduler.NewMemoryStore()).
        WithMonitor(scheduler.NewMonitor()).
        WithCheckInterval(time.Second).
        WithMaxConcurrency(5).
        WithMetrics(true).
        WithLogging(true)

    scheduler := config.Build()

    // 创建任务
    handler := scheduler.NewFuncHandler("backup", func(ctx context.Context) error {
        fmt.Println("执行备份任务")
        return nil
    })

    task := scheduler.NewTaskBuilder("backup", "数据库备份", "0 2 * * *", handler).
        SetTimeout(30 * time.Minute).
        SetMaxRetries(3).
        SetRetryDelay(5 * time.Minute).
        AddTag("type", "backup").
        AddTag("priority", "high").
        Build()

    scheduler.Add(task)
    scheduler.Start()

    // 监控
    go func() {
        ticker := time.NewTicker(time.Minute)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                status := scheduler.GetStatus()
                stats := scheduler.GetStats()
                fmt.Printf("状态: %s, 任务数: %d, 成功率: %.2f%%\n",
                    status.Status, stats.TotalTasks, stats.SuccessRate)
            }
        }
    }()

    // 保持运行
    select {}
}
```

## 分布式集群支持

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

### Redis 集群

```go
// RedisClusterConfig Redis集群配置
type RedisClusterConfig struct {
    Addr     string
    Password string
    DB       int
    NodeID   string
}

// NewRedisCluster 创建Redis集群
func NewRedisCluster(config RedisClusterConfig) (*RedisCluster, error)
```

**示例**:

```go
config := scheduler.RedisClusterConfig{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
    NodeID:   "node-1",
}

cluster, err := scheduler.NewRedisCluster(config)
```

### etcd 集群

```go
// EtcdClusterConfig etcd集群配置
type EtcdClusterConfig struct {
    Endpoints []string
    Username  string
    Password  string
    NodeID    string
}

// NewEtcdCluster 创建etcd集群
func NewEtcdCluster(config EtcdClusterConfig) (*EtcdCluster, error)
```

**示例**:

```go
config := scheduler.EtcdClusterConfig{
    Endpoints: []string{"localhost:2379"},
    Username:  "",
    Password:  "",
    NodeID:    "node-1",
}

cluster, err := scheduler.NewEtcdCluster(config)
```

### Consul 集群

```go
// ConsulClusterConfig Consul集群配置
type ConsulClusterConfig struct {
    Address string
    Token   string
    NodeID  string
}

// NewConsulCluster 创建Consul集群
func NewConsulCluster(config ConsulClusterConfig) (*ConsulCluster, error)
```

**示例**:

```go
config := scheduler.ConsulClusterConfig{
    Address: "localhost:8500",
    Token:   "",
    NodeID:  "node-1",
}

cluster, err := scheduler.NewConsulCluster(config)
```

### ZooKeeper 集群

```go
// ZookeeperClusterConfig ZooKeeper集群配置
type ZookeeperClusterConfig struct {
    Servers []string
    NodeID  string
}

// NewZookeeperCluster 创建ZooKeeper集群
func NewZookeeperCluster(config ZookeeperClusterConfig) (*ZookeeperCluster, error)
```

**示例**:

```go
config := scheduler.ZookeeperClusterConfig{
    Servers: []string{"localhost:2181"},
    NodeID:  "node-1",
}

cluster, err := scheduler.NewZookeeperCluster(config)
```

### 分布式调度器

```go
// DistributedScheduler 分布式调度器
type DistributedScheduler struct {
    *DefaultScheduler
    nodeID       string
    cluster      Cluster
    leader       bool
    leaderMu     sync.RWMutex
    electionMu   sync.Mutex
    stopElection chan struct{}
}

// DistributedConfig 分布式配置
type DistributedConfig struct {
    NodeID                 string
    Cluster                Cluster
    ElectionTimeout        time.Duration
    LockTimeout            time.Duration
    HeartbeatInterval      time.Duration
    EnableLeaderElection   bool
    EnableTaskDistribution bool
}

// NewDistributedScheduler 创建分布式调度器
func NewDistributedScheduler(store Store, config DistributedConfig) *DistributedScheduler
```

**示例**:

```go
config := scheduler.DistributedConfig{
    NodeID:                 "node-1",
    Cluster:                cluster,
    ElectionTimeout:        30 * time.Second,
    LockTimeout:            10 * time.Second,
    HeartbeatInterval:      5 * time.Second,
    EnableLeaderElection:   true,
    EnableTaskDistribution: true,
}

ds := scheduler.NewDistributedScheduler(store, config)
```

### 分布式调度器方法

```go
// Start 启动分布式调度器
func (ds *DistributedScheduler) Start() error

// Stop 停止分布式调度器
func (ds *DistributedScheduler) Stop() error

// IsLeader 检查是否为领导者
func (ds *DistributedScheduler) IsLeader() bool

// GetClusterNodes 获取集群节点
func (ds *DistributedScheduler) GetClusterNodes() ([]NodeInfo, error)

// GetDistributedStats 获取分布式统计
func (ds *DistributedScheduler) GetDistributedStats() DistributedStats
```

### 集群信息结构

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

// ClusterMessage 集群消息
type ClusterMessage struct {
    Type      string    `json:"type"`
    NodeID    string    `json:"node_id"`
    Timestamp time.Time `json:"timestamp"`
    Data      []byte    `json:"data"`
}

// TaskExecution 任务执行记录
type TaskExecution struct {
    TaskID    string     `json:"task_id"`
    NodeID    string     `json:"node_id"`
    Status    string     `json:"status"` // running, completed, failed
    StartedAt time.Time  `json:"started_at"`
    EndedAt   *time.Time `json:"ended_at,omitempty"`
    Error     string     `json:"error,omitempty"`
}

// DistributedStats 分布式统计
type DistributedStats struct {
    NodeID      string `json:"node_id"`
    IsLeader    bool   `json:"is_leader"`
    TotalNodes  int    `json:"total_nodes"`
    OnlineNodes int    `json:"online_nodes"`
    LeaderID    string `json:"leader_id"`
}
```

### 多集群使用示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"
    "laravel-go/framework/scheduler"
)

func main() {
    // 根据环境变量选择集群类型
    clusterType := os.Getenv("CLUSTER_TYPE")
    nodeID := os.Getenv("NODE_ID")

    var cluster scheduler.Cluster
    var err error

    switch clusterType {
    case "redis":
        cluster, err = scheduler.NewRedisCluster(scheduler.RedisClusterConfig{
            Addr:   os.Getenv("REDIS_ADDR"),
            NodeID: nodeID,
        })
    case "etcd":
        cluster, err = scheduler.NewEtcdCluster(scheduler.EtcdClusterConfig{
            Endpoints: []string{os.Getenv("ETCD_ENDPOINTS")},
            NodeID:    nodeID,
        })
    case "consul":
        cluster, err = scheduler.NewConsulCluster(scheduler.ConsulClusterConfig{
            Address: os.Getenv("CONSUL_ADDRESS"),
            NodeID:  nodeID,
        })
    case "zookeeper":
        cluster, err = scheduler.NewZookeeperCluster(scheduler.ZookeeperClusterConfig{
            Servers: []string{os.Getenv("ZOOKEEPER_SERVERS")},
            NodeID:  nodeID,
        })
    default:
        log.Fatal("Unsupported cluster type")
    }

    if err != nil {
        log.Fatal(err)
    }
    defer cluster.Close()

    // 创建分布式调度器
    config := scheduler.DistributedConfig{
        NodeID:   nodeID,
        Cluster:  cluster,
        ElectionTimeout: 30 * time.Second,
    }

    store := scheduler.NewMemoryStore()
    ds := scheduler.NewDistributedScheduler(store, config)

    // 创建任务
    handler := scheduler.NewFuncHandler("distributed-task", func(ctx context.Context) error {
        fmt.Printf("分布式任务执行 (节点: %s)\n", ds.GetDistributedStats().NodeID)
        return nil
    })

    task := scheduler.NewTask("distributed-task", "分布式任务", "0 * * * * *", handler)
    ds.Add(task)

    // 启动调度器
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

## 总结

本文档提供了 Laravel-Go 定时器模块的完整 API 参考，包括所有接口、方法、结构体和常量的详细说明。通过合理使用这些 API，可以构建强大而灵活的任务调度系统。

**分布式支持**提供了多种集群实现选项，包括 Redis、etcd、Consul 和 ZooKeeper，可以根据不同的部署环境和需求选择合适的集群方案。
