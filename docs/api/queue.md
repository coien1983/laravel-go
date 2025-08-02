# 队列 API 参考

## 📋 队列系统概览

Laravel-Go Framework 提供了强大的队列系统，支持异步任务处理、任务调度、失败重试等功能。系统支持单机模式和分布式模式，分布式模式支持多种集群后端（Redis、etcd、Consul、ZooKeeper）。

## 🚀 快速开始

### 基本使用

```go
import "laravel-go/framework/queue"

// 获取队列实例
q := queue.Driver("default")

// 推送任务到队列
job := &SendEmailJob{
    To:      "user@example.com",
    Subject: "Welcome",
    Body:    "Welcome to our platform!",
}
q.Push(job)

// 处理队列任务
worker := queue.NewWorker(q)
worker.Start()

// 分布式队列使用
cluster, err := queue.NewRedisCluster(queue.RedisClusterConfig{
    Addr:   "localhost:6379",
    NodeID: "node-1",
})
if err != nil {
    log.Fatal(err)
}

config := queue.DistributedConfig{
    NodeID:   "node-1",
    Cluster:  cluster,
    WorkerCount: 3,
}

dq := queue.NewDistributedQueue(config)
dq.Start()
```

## 📋 API 参考

### 核心方法

#### Push - 推送任务

```go
// 推送任务到队列
func (q *Queue) Push(job Job) error

// 示例
job := &SendEmailJob{To: "user@example.com"}
err := q.Push(job)

// 延迟推送
func (q *Queue) Later(delay time.Duration, job Job) error

// 示例
err := q.Later(time.Hour, &SendEmailJob{To: "user@example.com"})
```

#### Job - 任务接口

```go
// 任务接口
type Job interface {
    Handle() error
    GetQueue() string
    GetDelay() time.Duration
    GetAttempts() int
    GetMaxAttempts() int
    GetTimeout() time.Duration
    GetRetryAfter() time.Duration
    Failed(error)
}

// 基本任务结构
type BaseJob struct {
    Queue       string        `json:"queue"`
    Delay       time.Duration `json:"delay"`
    Attempts    int           `json:"attempts"`
    MaxAttempts int           `json:"max_attempts"`
    Timeout     time.Duration `json:"timeout"`
    RetryAfter  time.Duration `json:"retry_after"`
}

func (j *BaseJob) GetQueue() string {
    return j.Queue
}

func (j *BaseJob) GetDelay() time.Duration {
    return j.Delay
}

func (j *BaseJob) GetAttempts() int {
    return j.Attempts
}

func (j *BaseJob) GetMaxAttempts() int {
    return j.MaxAttempts
}

func (j *BaseJob) GetTimeout() time.Duration {
    return j.Timeout
}

func (j *BaseJob) GetRetryAfter() time.Duration {
    return j.RetryAfter
}
```

### 任务类型

#### 邮件发送任务

```go
type SendEmailJob struct {
    BaseJob
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

func (j *SendEmailJob) Handle() error {
    // 发送邮件逻辑
    return sendEmail(j.To, j.Subject, j.Body)
}

func (j *SendEmailJob) Failed(err error) {
    // 任务失败处理
    log.Printf("Failed to send email to %s: %v", j.To, err)
}
```

#### 数据处理任务

```go
type ProcessDataJob struct {
    BaseJob
    DataID int    `json:"data_id"`
    Action string `json:"action"`
}

func (j *ProcessDataJob) Handle() error {
    // 数据处理逻辑
    return processData(j.DataID, j.Action)
}

func (j *ProcessDataJob) Failed(err error) {
    // 任务失败处理
    log.Printf("Failed to process data %d: %v", j.DataID, err)
}
```

#### 文件处理任务

```go
type ProcessFileJob struct {
    BaseJob
    FilePath string `json:"file_path"`
    Action   string `json:"action"`
}

func (j *ProcessFileJob) Handle() error {
    // 文件处理逻辑
    return processFile(j.FilePath, j.Action)
}

func (j *ProcessFileJob) Failed(err error) {
    // 任务失败处理
    log.Printf("Failed to process file %s: %v", j.FilePath, err)
}
```

### 队列驱动

#### 内存驱动

```go
// 使用内存队列
q := queue.Driver("memory")

// 配置
type MemoryConfig struct {
    MaxJobs int `env:"QUEUE_MEMORY_MAX_JOBS" default:"1000"`
}
```

#### 数据库驱动

```go
// 使用数据库队列
q := queue.Driver("database")

// 配置
type DatabaseConfig struct {
    Connection string `env:"QUEUE_DATABASE_CONNECTION" default:"default"`
    Table      string `env:"QUEUE_DATABASE_TABLE" default:"jobs"`
    Queue      string `env:"QUEUE_DATABASE_QUEUE" default:"default"`
}
```

#### Redis 驱动

```go
// 使用 Redis 队列
q := queue.Driver("redis")

// 配置
type RedisConfig struct {
    Host     string `env:"QUEUE_REDIS_HOST" default:"localhost"`
    Port     int    `env:"QUEUE_REDIS_PORT" default:"6379"`
    Password string `env:"QUEUE_REDIS_PASSWORD"`
    Database int    `env:"QUEUE_REDIS_DB" default:"0"`
    Queue    string `env:"QUEUE_REDIS_QUEUE" default:"default"`
}
```

### 工作进程

#### 基本工作进程

```go
// 创建工作进程
worker := queue.NewWorker(q)

// 启动工作进程
worker.Start()

// 停止工作进程
worker.Stop()

// 设置并发数
worker.SetConcurrency(5)

// 设置超时时间
worker.SetTimeout(time.Minute * 5)
```

#### 高级工作进程

```go
// 创建工作进程配置
config := &queue.WorkerConfig{
    Concurrency: 5,
    Timeout:     time.Minute * 5,
    Sleep:       time.Second,
    MaxAttempts: 3,
    RetryAfter:  time.Minute,
}

worker := queue.NewWorkerWithConfig(q, config)
worker.Start()
```

### 任务调度

#### 基本调度

```go
// 创建调度器
scheduler := queue.NewScheduler()

// 添加定时任务
scheduler.Add(&SendEmailJob{To: "user@example.com"}, "0 9 * * *") // 每天上午9点

// 启动调度器
scheduler.Start()
```

#### 高级调度

```go
// 创建调度器配置
config := &queue.SchedulerConfig{
    Timezone: "Asia/Shanghai",
    LogLevel: "info",
}

scheduler := queue.NewSchedulerWithConfig(config)

// 添加多种定时任务
scheduler.Add(&SendEmailJob{To: "user@example.com"}, "0 9 * * *")     // 每天上午9点
scheduler.Add(&ProcessDataJob{DataID: 1}, "*/5 * * * *")              // 每5分钟
scheduler.Add(&ProcessFileJob{FilePath: "/tmp/file"}, "0 2 * * *")    // 每天凌晨2点

scheduler.Start()
```

## 🎯 使用示例

### 邮件队列

```go
type EmailService struct {
    queue queue.Queue
}

func (s *EmailService) SendWelcomeEmail(user *User) error {
    job := &SendEmailJob{
        BaseJob: BaseJob{
            Queue:       "emails",
            MaxAttempts: 3,
            Timeout:     time.Minute * 5,
        },
        To:      user.Email,
        Subject: "Welcome to our platform",
        Body:    fmt.Sprintf("Hello %s, welcome to our platform!", user.Name),
    }

    return s.queue.Push(job)
}

func (s *EmailService) SendPasswordResetEmail(email, token string) error {
    job := &SendEmailJob{
        BaseJob: BaseJob{
            Queue:       "emails",
            MaxAttempts: 3,
            Timeout:     time.Minute * 5,
        },
        To:      email,
        Subject: "Password Reset",
        Body:    fmt.Sprintf("Your password reset token is: %s", token),
    }

    return s.queue.Push(job)
}
```

### 数据处理队列

```go
type DataService struct {
    queue queue.Queue
}

func (s *DataService) ProcessUserData(userID int) error {
    job := &ProcessDataJob{
        BaseJob: BaseJob{
            Queue:       "data",
            MaxAttempts: 5,
            Timeout:     time.Minute * 10,
        },
        DataID: userID,
        Action: "process",
    }

    return s.queue.Push(job)
}

func (s *DataService) ExportUserData(userID int) error {
    job := &ProcessDataJob{
        BaseJob: BaseJob{
            Queue:       "exports",
            MaxAttempts: 3,
            Timeout:     time.Minute * 30,
        },
        DataID: userID,
        Action: "export",
    }

    return s.queue.Push(job)
}
```

### 文件处理队列

```go
type FileService struct {
    queue queue.Queue
}

func (s *FileService) ProcessUploadedFile(filePath string) error {
    job := &ProcessFileJob{
        BaseJob: BaseJob{
            Queue:       "files",
            MaxAttempts: 3,
            Timeout:     time.Minute * 15,
        },
        FilePath: filePath,
        Action:   "process",
    }

    return s.queue.Push(job)
}

func (s *FileService) GenerateThumbnail(filePath string) error {
    job := &ProcessFileJob{
        BaseJob: BaseJob{
            Queue:       "thumbnails",
            MaxAttempts: 3,
            Timeout:     time.Minute * 5,
        },
        FilePath: filePath,
        Action:   "thumbnail",
    }

    return s.queue.Push(job)
}
```

## 🔄 任务重试

### 重试机制

```go
type RetryableJob struct {
    BaseJob
    Data interface{} `json:"data"`
}

func (j *RetryableJob) Handle() error {
    // 任务处理逻辑
    err := processData(j.Data)
    if err != nil {
        // 如果失败，增加重试次数
        j.Attempts++

        // 如果还有重试机会，抛出错误让队列重试
        if j.Attempts < j.MaxAttempts {
            return err
        }
    }

    return nil
}

func (j *RetryableJob) Failed(err error) {
    // 任务最终失败处理
    log.Printf("Job failed after %d attempts: %v", j.Attempts, err)

    // 可以发送通知、记录日志等
    sendFailureNotification(j, err)
}
```

### 指数退避

```go
type ExponentialBackoffJob struct {
    BaseJob
    Data interface{} `json:"data"`
}

func (j *ExponentialBackoffJob) GetRetryAfter() time.Duration {
    // 指数退避：1秒、2秒、4秒、8秒...
    return time.Duration(math.Pow(2, float64(j.Attempts))) * time.Second
}

func (j *ExponentialBackoffJob) Handle() error {
    // 任务处理逻辑
    return processData(j.Data)
}
```

## 📊 队列监控

### 队列统计

```go
// 获取队列统计信息
func (q *Queue) Stats() *QueueStats

type QueueStats struct {
    TotalJobs     int64 `json:"total_jobs"`
    PendingJobs   int64 `json:"pending_jobs"`
    ProcessingJobs int64 `json:"processing_jobs"`
    FailedJobs    int64 `json:"failed_jobs"`
    CompletedJobs int64 `json:"completed_jobs"`
}

// 示例
stats := q.Stats()
fmt.Printf("Queue stats: total=%d, pending=%d, processing=%d, failed=%d, completed=%d\n",
    stats.TotalJobs, stats.PendingJobs, stats.ProcessingJobs, stats.FailedJobs, stats.CompletedJobs)
```

### 队列监控

```go
type QueueMonitor struct {
    queue queue.Queue
}

func (m *QueueMonitor) GetQueueHealth() *QueueHealth {
    stats := m.queue.Stats()

    health := &QueueHealth{
        Status: "healthy",
        Stats:  stats,
    }

    // 检查队列健康状态
    if stats.FailedJobs > 100 {
        health.Status = "warning"
    }

    if stats.FailedJobs > 1000 {
        health.Status = "critical"
    }

    return health
}

type QueueHealth struct {
    Status string      `json:"status"`
    Stats  *QueueStats `json:"stats"`
}
```

## 🛡️ 错误处理

### 错误类型

```go
// 队列错误类型
type QueueError struct {
    Message string
    Job     Job
    Err     error
}

func (e *QueueError) Error() string {
    return fmt.Sprintf("queue error: %s", e.Message)
}

// 处理队列错误
func handleQueueError(err error, job Job) {
    if queueErr, ok := err.(*QueueError); ok {
        log.Printf("Queue error for job: %v", queueErr.Err)

        // 可以发送告警、记录日志等
        sendQueueErrorAlert(queueErr)
    } else {
        log.Printf("Unknown queue error: %v", err)
    }
}
```

### 错误处理示例

```go
func (j *SendEmailJob) Handle() error {
    // 发送邮件逻辑
    err := sendEmail(j.To, j.Subject, j.Body)
    if err != nil {
        // 记录错误日志
        log.Printf("Failed to send email to %s: %v", j.To, err)

        // 如果是临时错误，可以重试
        if isTemporaryError(err) {
            return err
        }

        // 如果是永久错误，标记为失败
        j.Failed(err)
        return nil
    }

    return nil
}

func isTemporaryError(err error) bool {
    // 判断是否为临时错误（网络问题、服务暂时不可用等）
    return strings.Contains(err.Error(), "connection refused") ||
           strings.Contains(err.Error(), "timeout")
}
```

## 📝 最佳实践

### 1. 任务设计

```go
// 任务应该小而专注
type SendWelcomeEmailJob struct {
    BaseJob
    UserID int `json:"user_id"`
}

func (j *SendWelcomeEmailJob) Handle() error {
    // 获取用户信息
    user, err := getUser(j.UserID)
    if err != nil {
        return err
    }

    // 发送邮件
    return sendEmail(user.Email, "Welcome", "Welcome to our platform!")
}

// 避免在任务中做太多事情
type ProcessUserJob struct {
    BaseJob
    UserID int `json:"user_id"`
}

func (j *ProcessUserJob) Handle() error {
    // 只做一件事：处理用户数据
    return processUserData(j.UserID)
}
```

### 2. 队列配置

```go
// 根据任务类型配置不同的队列
const (
    EmailQueue     = "emails"
    DataQueue      = "data"
    FileQueue      = "files"
    NotificationQueue = "notifications"
)

// 配置不同队列的工作进程
func startWorkers() {
    // 邮件队列：高并发，短超时
    emailWorker := queue.NewWorkerWithConfig(queue.Driver(EmailQueue), &queue.WorkerConfig{
        Concurrency: 10,
        Timeout:     time.Minute * 2,
    })
    go emailWorker.Start()

    // 数据处理队列：低并发，长超时
    dataWorker := queue.NewWorkerWithConfig(queue.Driver(DataQueue), &queue.WorkerConfig{
        Concurrency: 3,
        Timeout:     time.Minute * 30,
    })
    go dataWorker.Start()

    // 文件处理队列：中等并发，中等超时
    fileWorker := queue.NewWorkerWithConfig(queue.Driver(FileQueue), &queue.WorkerConfig{
        Concurrency: 5,
        Timeout:     time.Minute * 10,
    })
    go fileWorker.Start()
}
```

### 3. 任务重试策略

```go
// 根据任务类型设置不同的重试策略
type EmailJob struct {
    BaseJob
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

func NewEmailJob(to, subject, body string) *EmailJob {
    return &EmailJob{
        BaseJob: BaseJob{
            Queue:       "emails",
            MaxAttempts: 3,        // 邮件任务重试3次
            Timeout:     time.Minute * 2,
            RetryAfter:  time.Minute * 5, // 5分钟后重试
        },
        To:      to,
        Subject: subject,
        Body:    body,
    }
}

type DataJob struct {
    BaseJob
    DataID int    `json:"data_id"`
    Action string `json:"action"`
}

func NewDataJob(dataID int, action string) *DataJob {
    return &DataJob{
        BaseJob: BaseJob{
            Queue:       "data",
            MaxAttempts: 5,        // 数据处理任务重试5次
            Timeout:     time.Minute * 30,
            RetryAfter:  time.Minute * 10, // 10分钟后重试
        },
        DataID: dataID,
        Action: action,
    }
}
```

### 4. 监控和告警

```go
type QueueHealthChecker struct {
    queues map[string]queue.Queue
}

func (c *QueueHealthChecker) CheckHealth() {
    for name, q := range c.queues {
        stats := q.Stats()

        // 检查失败任务数量
        if stats.FailedJobs > 100 {
            sendAlert(fmt.Sprintf("Queue %s has too many failed jobs: %d", name, stats.FailedJobs))
        }

        // 检查队列积压
        if stats.PendingJobs > 1000 {
            sendAlert(fmt.Sprintf("Queue %s has too many pending jobs: %d", name, stats.PendingJobs))
        }

        // 检查处理中的任务
        if stats.ProcessingJobs > 100 {
            sendAlert(fmt.Sprintf("Queue %s has too many processing jobs: %d", name, stats.ProcessingJobs))
        }
    }
}

func (c *QueueHealthChecker) StartMonitoring() {
    ticker := time.NewTicker(time.Minute * 5)
    go func() {
        for range ticker.C {
            c.CheckHealth()
        }
    }()
}
```

## 🔄 分布式队列 API

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

### 分布式配置

```go
// DistributedConfig 分布式配置
type DistributedConfig struct {
    NodeID                 string
    Cluster                Cluster
    ElectionTimeout        time.Duration
    LockTimeout            time.Duration
    HeartbeatInterval      time.Duration
    EnableLeaderElection   bool
    EnableJobDistribution  bool
    WorkerCount            int
    MaxConcurrency         int
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

// 创建Redis集群
cluster, err := queue.NewRedisCluster(queue.RedisClusterConfig{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
    NodeID:   "node-1",
})
```

### etcd 集群

```go
// EtcdClusterConfig etcd集群配置
type EtcdClusterConfig struct {
    Endpoints []string
    NodeID    string
}

// 创建etcd集群
cluster, err := queue.NewEtcdCluster(queue.EtcdClusterConfig{
    Endpoints: []string{"localhost:2379"},
    NodeID:    "node-1",
})
```

### Consul 集群

```go
// ConsulClusterConfig Consul集群配置
type ConsulClusterConfig struct {
    Address string
    NodeID  string
}

// 创建Consul集群
cluster, err := queue.NewConsulCluster(queue.ConsulClusterConfig{
    Address: "localhost:8500",
    NodeID:  "node-1",
})
```

### ZooKeeper 集群

```go
// ZookeeperClusterConfig ZooKeeper集群配置
type ZookeeperClusterConfig struct {
    Servers []string
    NodeID  string
}

// 创建ZooKeeper集群
cluster, err := queue.NewZookeeperCluster(queue.ZookeeperClusterConfig{
    Servers: []string{"localhost:2181"},
    NodeID:  "node-1",
})
```

### 分布式队列

```go
// 创建分布式队列
dq := queue.NewDistributedQueue(queue.DistributedConfig{
    NodeID:                "node-1",
    Cluster:               cluster,
    ElectionTimeout:       30 * time.Second,
    LockTimeout:           10 * time.Second,
    HeartbeatInterval:     5 * time.Second,
    EnableLeaderElection:  true,
    EnableJobDistribution: true,
    WorkerCount:           3,
    MaxConcurrency:        5,
})

// 启动分布式队列
err := dq.Start()

// 获取分布式统计
stats := dq.GetDistributedStats()
fmt.Printf("节点ID: %s, 是否为领导者: %t\n", stats.NodeID, stats.IsLeader)

// 获取集群节点
nodes, err := dq.GetClusterNodes()
```

### 分布式工作进程池

```go
// 获取工作进程池
workerPool := dq.GetWorkerPool()

// 设置回调
workerPool.SetOnCompleted(func(job queue.Job) {
    fmt.Printf("分布式任务完成: %s\n", string(job.GetPayload()))
})

workerPool.SetOnFailed(func(job queue.Job, err error) {
    fmt.Printf("分布式任务失败: %s - %v\n", string(job.GetPayload()), err)
})

// 获取统计信息
poolStats := workerPool.GetStats()
fmt.Printf("工作进程池状态: %s, 总工作进程: %d, 活跃: %d\n",
    poolStats.Status, poolStats.TotalWorkers, poolStats.ActiveWorkers)
```

## 📚 总结

Laravel-Go Framework 的队列 API 提供了：

1. **异步处理**: 支持异步任务处理，提高系统响应速度
2. **任务调度**: 支持定时任务和延迟任务
3. **重试机制**: 内置任务重试和失败处理
4. **多种驱动**: 支持内存、数据库、Redis 等队列驱动
5. **分布式支持**: 支持多种集群后端（Redis、etcd、Consul、ZooKeeper）
6. **领导者选举**: 自动选举领导者节点，确保任务分发的唯一性
7. **分布式锁**: 防止任务重复处理
8. **监控功能**: 提供队列统计和健康检查

通过合理使用队列 API，可以构建出高效、可靠的异步任务处理系统，支持从单机到分布式集群的各种部署场景。
