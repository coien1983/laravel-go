# 队列系统指南

## 📖 概述

Laravel-Go Framework 提供了强大的队列系统，支持异步任务处理、任务调度、失败重试、队列监控等功能，帮助提升应用程序的性能和可靠性。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [队列系统 API 参考](../api/queue.md)

## 🚀 快速开始

### 1. 基本使用

```go
// 创建任务
type SendEmailJob struct {
    queue.BaseJob
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

func NewSendEmailJob(to, subject, body string) *SendEmailJob {
    return &SendEmailJob{
        BaseJob: queue.BaseJob{
            Queue:       "emails",
            MaxAttempts: 3,
            Timeout:     time.Minute * 5,
        },
        To:      to,
        Subject: subject,
        Body:    body,
    }
}

// 实现任务处理逻辑
func (j *SendEmailJob) Handle() error {
    return sendEmail(j.To, j.Subject, j.Body)
}

// 推送任务到队列
func (s *EmailService) SendWelcomeEmail(user *User) error {
    job := NewSendEmailJob(
        user.Email,
        "Welcome to our platform!",
        "Thank you for joining us.",
    )

    return queue.Push(job)
}
```

### 2. 启动队列工作进程

```go
// 启动队列工作进程
func main() {
    // 启动队列工作进程
    go func() {
        worker := queue.NewWorker("emails", 5) // 5个并发工作进程
        worker.Start()
    }()

    // 启动 Web 服务器
    server := http.NewServer()
    server.Start(":8080")
}
```

## 🔧 队列驱动

### 1. Redis 驱动

```go
// 配置 Redis 队列
config.Set("queue.driver", "redis")
config.Set("queue.redis.host", "localhost")
config.Set("queue.redis.port", 6379)
config.Set("queue.redis.password", "")
config.Set("queue.redis.database", 0)

// 使用 Redis 队列
redisQueue := queue.NewRedisDriver(config.Get("queue.redis"))
redisQueue.Push(job)
```

### 2. 内存驱动

```go
// 配置内存队列
config.Set("queue.driver", "memory")

// 使用内存队列
memoryQueue := queue.NewMemoryDriver()
memoryQueue.Push(job)
```

### 3. 数据库驱动

```go
// 配置数据库队列
config.Set("queue.driver", "database")
config.Set("queue.database.table", "jobs")

// 使用数据库队列
dbQueue := queue.NewDatabaseDriver(db, "jobs")
dbQueue.Push(job)
```

### 4. RabbitMQ 驱动

```go
// 配置 RabbitMQ 队列
config.Set("queue.driver", "rabbitmq")
config.Set("queue.rabbitmq.url", "amqp://guest:guest@localhost:5672/")

// 使用 RabbitMQ 队列
rabbitQueue := queue.NewRabbitMQDriver(config.Get("queue.rabbitmq"))
rabbitQueue.Push(job)
```

## 📋 任务类型

### 1. 同步任务

```go
// 同步执行任务
type ProcessImageJob struct {
    queue.BaseJob
    ImagePath string `json:"image_path"`
}

func (j *ProcessImageJob) Handle() error {
    // 处理图片
    return processImage(j.ImagePath)
}

// 同步执行
func (s *ImageService) ProcessImage(path string) error {
    job := &ProcessImageJob{ImagePath: path}
    return job.Handle() // 直接执行，不推送到队列
}
```

### 2. 延迟任务

```go
// 延迟执行任务
type ReminderJob struct {
    queue.BaseJob
    UserID uint   `json:"user_id"`
    Message string `json:"message"`
}

func (j *ReminderJob) Handle() error {
    // 发送提醒
    return sendReminder(j.UserID, j.Message)
}

// 延迟推送任务
func (s *ReminderService) ScheduleReminder(userID uint, message string, delay time.Duration) error {
    job := &ReminderJob{
        UserID:  userID,
        Message: message,
    }

    return queue.Later(job, delay)
}
```

### 3. 链式任务

```go
// 链式任务
type ProcessOrderJob struct {
    queue.BaseJob
    OrderID uint `json:"order_id"`
}

func (j *ProcessOrderJob) Handle() error {
    // 处理订单
    return processOrder(j.OrderID)
}

type SendOrderConfirmationJob struct {
    queue.BaseJob
    OrderID uint `json:"order_id"`
}

func (j *SendOrderConfirmationJob) Handle() error {
    // 发送订单确认邮件
    return sendOrderConfirmation(j.OrderID)
}

// 创建任务链
func (s *OrderService) ProcessOrder(orderID uint) error {
    chain := queue.NewChain(
        &ProcessOrderJob{OrderID: orderID},
        &SendOrderConfirmationJob{OrderID: orderID},
    )

    return chain.Dispatch()
}
```

### 4. 批量任务

```go
// 批量任务
type BatchEmailJob struct {
    queue.BaseJob
    Emails []string `json:"emails"`
    Subject string  `json:"subject"`
    Body    string  `json:"body"`
}

func (j *BatchEmailJob) Handle() error {
    for _, email := range j.Emails {
        if err := sendEmail(email, j.Subject, j.Body); err != nil {
            log.Printf("Failed to send email to %s: %v", email, err)
        }
    }
    return nil
}

// 批量发送邮件
func (s *EmailService) SendBatchEmails(emails []string, subject, body string) error {
    job := &BatchEmailJob{
        Emails:  emails,
        Subject: subject,
        Body:    body,
    }

    return queue.Push(job)
}
```

## 🔄 任务调度

### 1. 定时任务

```go
// 定时任务
type DailyReportJob struct {
    queue.BaseJob
}

func (j *DailyReportJob) Handle() error {
    // 生成每日报告
    return generateDailyReport()
}

// 注册定时任务
func RegisterScheduledJobs() {
    scheduler := queue.NewScheduler()

    // 每天凌晨2点执行
    scheduler.DailyAt("02:00", &DailyReportJob{})

    // 每小时执行
    scheduler.Hourly(&CleanupJob{})

    // 每周一执行
    scheduler.Weekly(&WeeklyReportJob{})

    // 自定义 Cron 表达式
    scheduler.Cron("0 */6 * * *", &CleanupJob{}) // 每6小时执行

    scheduler.Start()
}
```

### 2. 任务频率控制

```go
// 频率限制任务
type RateLimitedJob struct {
    queue.BaseJob
    UserID uint `json:"user_id"`
}

func (j *RateLimitedJob) Handle() error {
    // 检查频率限制
    if !j.checkRateLimit(j.UserID) {
        return errors.New("rate limit exceeded")
    }

    // 执行任务
    return j.processTask()
}

// 使用频率限制
func (s *Service) ProcessWithRateLimit(userID uint) error {
    job := &RateLimitedJob{UserID: userID}
    return queue.Push(job)
}
```

## 🛡️ 错误处理

### 1. 失败重试

```go
// 配置重试策略
type RetryableJob struct {
    queue.BaseJob
    Data string `json:"data"`
}

func NewRetryableJob(data string) *RetryableJob {
    return &RetryableJob{
        BaseJob: queue.BaseJob{
            MaxAttempts: 5,
            Backoff:     time.Second * 30, // 重试间隔
        },
        Data: data,
    }
}

func (j *RetryableJob) Handle() error {
    // 尝试处理任务
    return j.processData(j.Data)
}

// 失败回调
func (j *RetryableJob) Failed(err error) {
    log.Printf("Job failed after %d attempts: %v", j.Attempts, err)

    // 发送告警
    sendAlert(fmt.Sprintf("Job failed: %v", err))
}
```

### 2. 异常处理

```go
// 异常处理任务
type SafeJob struct {
    queue.BaseJob
    Critical bool `json:"critical"`
}

func (j *SafeJob) Handle() error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Job panicked: %v", r)

            if j.Critical {
                // 发送紧急告警
                sendEmergencyAlert(fmt.Sprintf("Critical job panicked: %v", r))
            }
        }
    }()

    return j.processTask()
}
```

## 📊 队列监控

### 1. 队列统计

```go
// 队列统计信息
type QueueStats struct {
    QueueName    string `json:"queue_name"`
    PendingJobs  int64  `json:"pending_jobs"`
    ProcessingJobs int64 `json:"processing_jobs"`
    FailedJobs   int64  `json:"failed_jobs"`
    CompletedJobs int64 `json:"completed_jobs"`
}

// 获取队列统计
func (s *QueueService) GetStats(queueName string) (*QueueStats, error) {
    stats := &QueueStats{QueueName: queueName}

    // 获取待处理任务数
    stats.PendingJobs = s.queue.Size(queueName)

    // 获取处理中任务数
    stats.ProcessingJobs = s.queue.Processing(queueName)

    // 获取失败任务数
    stats.FailedJobs = s.queue.Failed(queueName)

    // 获取已完成任务数
    stats.CompletedJobs = s.queue.Completed(queueName)

    return stats, nil
}
```

### 2. 任务监控

```go
// 任务监控中间件
type JobMonitorMiddleware struct {
    queue.Middleware
}

func (m *JobMonitorMiddleware) Before(job queue.Job) {
    log.Printf("Starting job: %s", job.GetName())

    // 记录开始时间
    job.SetMetadata("started_at", time.Now())
}

func (m *JobMonitorMiddleware) After(job queue.Job, err error) {
    duration := time.Since(job.GetMetadata("started_at").(time.Time))

    if err != nil {
        log.Printf("Job failed: %s, duration: %v, error: %v",
            job.GetName(), duration, err)
    } else {
        log.Printf("Job completed: %s, duration: %v",
            job.GetName(), duration)
    }
}
```

## 🔧 高级功能

### 1. 任务优先级

```go
// 优先级任务
type PriorityJob struct {
    queue.BaseJob
    Priority int `json:"priority"`
}

func NewPriorityJob(priority int) *PriorityJob {
    return &PriorityJob{
        BaseJob: queue.BaseJob{
            Queue:    "priority",
            Priority: priority,
        },
    }
}

// 使用优先级队列
func (s *Service) ProcessWithPriority(priority int) error {
    job := NewPriorityJob(priority)
    return queue.Push(job)
}
```

### 2. 任务超时控制

```go
// 超时控制任务
type TimeoutJob struct {
    queue.BaseJob
    Timeout time.Duration `json:"timeout"`
}

func (j *TimeoutJob) Handle() error {
    // 创建带超时的上下文
    ctx, cancel := context.WithTimeout(context.Background(), j.Timeout)
    defer cancel()

    // 在上下文中执行任务
    done := make(chan error, 1)
    go func() {
        done <- j.processTask()
    }()

    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return errors.New("job timeout")
    }
}
```

### 3. 任务依赖

```go
// 依赖任务
type DependentJob struct {
    queue.BaseJob
    Dependencies []string `json:"dependencies"`
}

func (j *DependentJob) Handle() error {
    // 检查依赖是否完成
    for _, dep := range j.Dependencies {
        if !j.isDependencyCompleted(dep) {
            return errors.New("dependency not completed")
        }
    }

    return j.processTask()
}
```

## 🚀 性能优化

### 1. 批量处理

```go
// 批量处理任务
type BatchProcessJob struct {
    queue.BaseJob
    Items []interface{} `json:"items"`
    BatchSize int       `json:"batch_size"`
}

func (j *BatchProcessJob) Handle() error {
    for i := 0; i < len(j.Items); i += j.BatchSize {
        end := i + j.BatchSize
        if end > len(j.Items) {
            end = len(j.Items)
        }

        batch := j.Items[i:end]
        if err := j.processBatch(batch); err != nil {
            return err
        }
    }

    return nil
}
```

### 2. 并发控制

```go
// 并发控制任务
type ConcurrentJob struct {
    queue.BaseJob
    MaxConcurrency int `json:"max_concurrency"`
}

func (j *ConcurrentJob) Handle() error {
    semaphore := make(chan struct{}, j.MaxConcurrency)
    var wg sync.WaitGroup

    for _, item := range j.Items {
        wg.Add(1)
        go func(item interface{}) {
            defer wg.Done()

            semaphore <- struct{}{}
            defer func() { <-semaphore }()

            j.processItem(item)
        }(item)
    }

    wg.Wait()
    return nil
}
```

## 📚 总结

Laravel-Go Framework 的队列系统提供了：

1. **多种驱动**: Redis、内存、数据库、RabbitMQ
2. **任务类型**: 同步、延迟、链式、批量
3. **任务调度**: 定时任务、频率控制
4. **错误处理**: 失败重试、异常处理
5. **监控功能**: 队列统计、任务监控
6. **高级功能**: 优先级、超时控制、依赖管理
7. **性能优化**: 批量处理、并发控制

通过合理使用队列系统，可以提升应用性能、可靠性和可维护性。
