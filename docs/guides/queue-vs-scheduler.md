# Queue vs Scheduler 详细对比

## 🎯 核心概念

### Queue（队列）

队列是一个**异步任务处理系统**，用于处理需要立即或延迟执行的任务，但不阻塞主业务流程。

### Scheduler（调度器）

调度器是一个**定时任务执行系统**，用于按照预定义的时间计划执行任务。

## 📊 详细对比表

| 特性         | Queue（队列）      | Scheduler（调度器） |
| ------------ | ------------------ | ------------------- |
| **触发方式** | 事件驱动、手动触发 | 时间驱动、规则驱动  |
| **执行时机** | 立即或延迟执行     | 按计划时间执行      |
| **任务来源** | 实时产生           | 预定义计划          |
| **处理模式** | 生产者-消费者      | 定时器模式          |
| **执行频率** | 按需执行           | 定期执行            |
| **任务类型** | 异步任务           | 定时任务            |
| **阻塞性**   | 非阻塞             | 非阻塞              |
| **优先级**   | 支持优先级         | 支持优先级          |
| **重试机制** | 支持重试           | 支持重试            |
| **监控统计** | 队列统计           | 执行统计            |

## 💡 使用场景对比

### 🚀 Queue 适用场景

#### 1. **用户操作触发的任务**

```go
// 用户注册后发送欢迎邮件
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
    // 创建用户
    user := createUser(r)

    // 推送邮件任务到队列
    job := queue.NewJob([]byte(user.Email), "emails")
    queue.Push(job) // 异步处理，不阻塞用户注册

    // 立即返回成功响应
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
```

#### 2. **文件处理任务**

```go
// 用户上传文件后异步处理
func (c *FileController) Upload(w http.ResponseWriter, r *http.Request) {
    file := saveUploadedFile(r)

    // 推送文件处理任务
    job := queue.NewJob([]byte(file.Path), "file-processing")
    queue.Push(job) // 异步处理文件压缩、格式转换等

    // 立即返回上传成功
    json.NewEncoder(w).Encode(map[string]string{"file_id": file.ID})
}
```

#### 3. **通知发送**

```go
// 系统事件触发通知
func (c *NotificationController) SendNotification(event Event) {
    // 推送通知任务
    job := queue.NewJob([]byte(event.Data), "notifications")
    queue.Later(job, 5*time.Minute) // 5分钟后发送

    // 继续处理其他逻辑
}
```

### ⏰ Scheduler 适用场景

#### 1. **数据备份任务**

```go
// 每天凌晨2点备份数据库
func (s *Scheduler) ScheduleBackup() {
    s.Cron("0 2 * * *", func() {
        // 执行数据库备份
        backupDatabase()
        log.Println("Database backup completed at", time.Now())
    })
}
```

#### 2. **报表生成**

```go
// 每周一生成周报表
func (s *Scheduler) ScheduleWeeklyReport() {
    s.Weekly(time.Monday, 9, 0, func() {
        // 生成周报表
        generateWeeklyReport()
        // 发送邮件通知
        sendReportEmail()
    })
}
```

#### 3. **系统维护**

```go
// 每小时清理临时文件
func (s *Scheduler) ScheduleCleanup() {
    s.Every(1).Hour().Do(func() {
        // 清理临时文件
        cleanupTempFiles()
        // 清理过期日志
        cleanupOldLogs()
    })
}
```

## 🔧 Laravel-Go Framework 实现

### Queue 实现示例

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

### Scheduler 实现示例

```go
package main

import (
    "context"
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

## 🎯 选择指南

### 何时使用 Queue？

- ✅ **需要异步处理**：不阻塞主业务流程
- ✅ **事件驱动**：由用户操作或系统事件触发
- ✅ **实时响应**：需要立即处理或延迟处理
- ✅ **高并发**：需要处理大量并发任务
- ✅ **任务优先级**：不同任务有不同的优先级

### 何时使用 Scheduler？

- ✅ **定时执行**：按计划时间执行任务
- ✅ **周期性任务**：每天、每周、每月执行
- ✅ **系统维护**：数据备份、清理、统计
- ✅ **报表生成**：定期生成报表和数据
- ✅ **监控检查**：定期检查系统状态

## 🔄 结合使用

在实际项目中，Queue 和 Scheduler 经常结合使用：

```go
// Scheduler 定时触发任务，将任务推送到 Queue
func (s *Scheduler) ScheduleDataProcessing() {
    s.Every(1).Hour().Do(func() {
        // 定时检查需要处理的数据
        pendingData := getPendingData()

        // 将每个数据项推送到队列
        for _, data := range pendingData {
            job := queue.NewJob([]byte(data.ID), "data-processing")
            queue.Push(job)
        }
    })
}

// Queue 异步处理具体的任务
func processDataJob(job queue.Job) error {
    dataID := string(job.GetPayload())
    // 处理具体的数据
    return processData(dataID)
}
```

## 📈 性能考虑

### Queue 性能特点

- **高吞吐量**：支持大量并发任务
- **低延迟**：任务可以立即处理
- **资源隔离**：任务处理不影响主业务
- **可扩展性**：支持分布式队列

### Scheduler 性能特点

- **精确调度**：按计划时间精确执行
- **资源控制**：可以控制并发执行数量
- **持久化**：任务计划可以持久化存储
- **监控完善**：详细的执行统计和监控

## 🎉 总结

- **Queue** = 异步任务处理系统，适合事件驱动的任务
- **Scheduler** = 定时任务执行系统，适合时间驱动的任务
- **结合使用** = 完整的任务处理解决方案

选择合适的系统取决于你的具体需求：如果需要异步处理用户操作产生的任务，使用 Queue；如果需要按计划执行系统维护任务，使用 Scheduler。
