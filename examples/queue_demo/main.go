package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"laravel-go/framework/queue"
)

// EmailJob 邮件任务示例
type EmailJob struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Process 处理邮件任务
func (j *EmailJob) Process() error {
	fmt.Printf("发送邮件到 %s: %s\n", j.To, j.Subject)
	time.Sleep(100 * time.Millisecond) // 模拟发送邮件
	return nil
}

func main() {
	fmt.Println("=== Laravel-Go 队列系统演示 ===\n")

	// 1. 初始化队列系统
	fmt.Println("1. 初始化队列系统:")
	queue.Init()
	
	// 注册内存队列
	memoryQueue := queue.NewMemoryQueue()
	queue.QueueManager.Extend("memory", memoryQueue)
	queue.QueueManager.SetDefaultQueue("memory")
	fmt.Println("   ✅ 内存队列已注册")

	// 2. 基本队列操作
	fmt.Println("\n2. 基本队列操作:")
	
	// 推送任务
	job1 := queue.NewJob([]byte("Hello Queue!"), "default")
	err := queue.Push(job1)
	if err != nil {
		log.Fatalf("推送任务失败: %v", err)
	}
	fmt.Println("   ✅ 任务已推送")

	job2 := queue.NewJob([]byte("Another job"), "default")
	err = queue.Push(job2)
	if err != nil {
		log.Fatalf("推送任务失败: %v", err)
	}
	fmt.Println("   ✅ 第二个任务已推送")

	// 获取队列大小
	size, err := queue.Size()
	if err != nil {
		log.Fatalf("获取队列大小失败: %v", err)
	}
	fmt.Printf("   📊 队列大小: %d\n", size)

	// 弹出任务
	ctx := context.Background()
	poppedJob, err := queue.Pop(ctx)
	if err != nil {
		log.Fatalf("弹出任务失败: %v", err)
	}
	fmt.Printf("   📤 弹出任务: %s\n", string(poppedJob.GetPayload()))

	// 3. 延迟队列
	fmt.Println("\n3. 延迟队列:")
	
	delayedJob := queue.NewJob([]byte("延迟任务"), "default")
	delayedJob.SetDelay(2 * time.Second)
	
	err = queue.Push(delayedJob)
	if err != nil {
		log.Fatalf("推送延迟任务失败: %v", err)
	}
	fmt.Println("   ⏰ 延迟任务已推送 (2秒后执行)")

	// 立即尝试弹出，应该失败
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	_, err = queue.Pop(ctx)
	if err != nil {
		fmt.Println("   ✅ 延迟任务未到执行时间")
	}

	// 等待延迟时间后弹出
	fmt.Println("   ⏳ 等待延迟任务...")
	time.Sleep(2 * time.Second)
	
	ctx = context.Background()
	delayedPoppedJob, err := queue.Pop(ctx)
	if err != nil {
		log.Fatalf("弹出延迟任务失败: %v", err)
	}
	fmt.Printf("   📤 延迟任务已弹出: %s\n", string(delayedPoppedJob.GetPayload()))

	// 4. 批量操作
	fmt.Println("\n4. 批量操作:")
	
	// 批量推送
	jobs := []queue.Job{
		queue.NewJob([]byte("批量任务1"), "default"),
		queue.NewJob([]byte("批量任务2"), "default"),
		queue.NewJob([]byte("批量任务3"), "default"),
	}
	
	err = memoryQueue.PushBatch(jobs)
	if err != nil {
		log.Fatalf("批量推送失败: %v", err)
	}
	fmt.Println("   ✅ 批量任务已推送")

	// 批量弹出
	poppedJobs, err := memoryQueue.PopBatch(ctx, 2)
	if err != nil {
		log.Fatalf("批量弹出失败: %v", err)
	}
	fmt.Printf("   📤 批量弹出 %d 个任务\n", len(poppedJobs))

	// 5. 工作进程
	fmt.Println("\n5. 工作进程:")
	
	// 创建一些任务
	for i := 0; i < 5; i++ {
		job := queue.NewJob([]byte(fmt.Sprintf("工作进程任务%d", i+1)), "default")
		err := queue.Push(job)
		if err != nil {
			log.Fatalf("推送任务失败: %v", err)
		}
	}

	// 创建工作进程
	worker := queue.NewWorker(memoryQueue, "default")
	
	// 设置回调
	worker.SetOnCompleted(func(job queue.Job) {
		fmt.Printf("   ✅ 任务完成: %s\n", string(job.GetPayload()))
	})
	
	worker.SetOnFailed(func(job queue.Job, err error) {
		fmt.Printf("   ❌ 任务失败: %s - %v\n", string(job.GetPayload()), err)
	})

	// 启动工作进程
	err = worker.Start()
	if err != nil {
		log.Fatalf("启动工作进程失败: %v", err)
	}
	fmt.Println("   🚀 工作进程已启动")

	// 等待一段时间让工作进程处理任务
	time.Sleep(1 * time.Second)

	// 获取工作进程状态
	status := worker.GetStatus()
	fmt.Printf("   📊 工作进程状态: %s\n", status.Status)
	fmt.Printf("   📊 已处理任务: %d\n", status.Processed)
	fmt.Printf("   📊 失败任务: %d\n", status.Failed)

	// 停止工作进程
	err = worker.Stop()
	if err != nil {
		log.Fatalf("停止工作进程失败: %v", err)
	}
	fmt.Println("   🛑 工作进程已停止")

	// 6. 工作进程池
	fmt.Println("\n6. 工作进程池:")
	
	// 创建更多任务
	for i := 0; i < 10; i++ {
		job := queue.NewJob([]byte(fmt.Sprintf("池任务%d", i+1)), "default")
		err := queue.Push(job)
		if err != nil {
			log.Fatalf("推送任务失败: %v", err)
		}
	}

	// 创建工作进程池
	pool := queue.NewWorkerPool(memoryQueue, "default", 3)
	
	// 启动工作进程池
	err = pool.Start()
	if err != nil {
		log.Fatalf("启动工作进程池失败: %v", err)
	}
	fmt.Println("   🚀 工作进程池已启动 (3个工作进程)")

	// 等待处理任务
	time.Sleep(2 * time.Second)

	// 获取池统计信息
	poolStats, err := pool.GetStats()
	if err != nil {
		log.Fatalf("获取池统计失败: %v", err)
	}
	
	totalProcessed := int64(0)
	totalFailed := int64(0)
	for _, stat := range poolStats {
		totalProcessed += stat.Processed
		totalFailed += stat.Failed
	}
	
	fmt.Printf("   📊 池总处理任务: %d\n", totalProcessed)
	fmt.Printf("   📊 池总失败任务: %d\n", totalFailed)

	// 停止工作进程池
	err = pool.Stop()
	if err != nil {
		log.Fatalf("停止工作进程池失败: %v", err)
	}
	fmt.Println("   🛑 工作进程池已停止")

	// 7. 队列统计
	fmt.Println("\n7. 队列统计:")
	
	stats, err := queue.GetStats()
	if err != nil {
		log.Fatalf("获取队列统计失败: %v", err)
	}
	
	fmt.Printf("   📊 总任务数: %d\n", stats.TotalJobs)
	fmt.Printf("   📊 待处理任务: %d\n", stats.PendingJobs)
	fmt.Printf("   📊 保留任务: %d\n", stats.ReservedJobs)
	fmt.Printf("   📊 失败任务: %d\n", stats.FailedJobs)
	fmt.Printf("   📊 完成任务: %d\n", stats.CompletedJobs)

	// 8. 清空队列
	fmt.Println("\n8. 清空队列:")
	
	err = queue.Clear()
	if err != nil {
		log.Fatalf("清空队列失败: %v", err)
	}
	fmt.Println("   ✅ 队列已清空")

	size, err = queue.Size()
	if err != nil {
		log.Fatalf("获取队列大小失败: %v", err)
	}
	fmt.Printf("   📊 清空后队列大小: %d\n", size)

	// 9. 多队列支持
	fmt.Println("\n9. 多队列支持:")
	
	// 创建第二个队列
	queue2 := queue.NewMemoryQueue()
	queue.QueueManager.Extend("queue2", queue2)
	
	// 推送任务到不同队列
	jobA := queue.NewJob([]byte("队列1任务"), "default")
	jobB := queue.NewJob([]byte("队列2任务"), "queue2")
	
	err = queue.PushTo("memory", jobA)
	if err != nil {
		log.Fatalf("推送任务到队列1失败: %v", err)
	}
	
	err = queue.PushTo("queue2", jobB)
	if err != nil {
		log.Fatalf("推送任务到队列2失败: %v", err)
	}
	
	fmt.Println("   ✅ 任务已推送到不同队列")
	
	// 获取不同队列的大小
	size1, _ := queue.SizeOf("memory")
	size2, _ := queue.SizeOf("queue2")
	fmt.Printf("   📊 队列1大小: %d, 队列2大小: %d\n", size1, size2)

	// 10. 任务属性
	fmt.Println("\n10. 任务属性:")
	
	advancedJob := queue.NewJob([]byte("高级任务"), "default")
	advancedJob.SetPriority(10)
	advancedJob.SetMaxAttempts(5)
	advancedJob.SetTimeout(60 * time.Second)
	advancedJob.AddTag("type", "email")
	advancedJob.AddTag("priority", "high")
	
	fmt.Printf("   📋 任务ID: %s\n", advancedJob.GetID())
	fmt.Printf("   📋 优先级: %d\n", advancedJob.GetPriority())
	fmt.Printf("   📋 最大尝试次数: %d\n", advancedJob.GetMaxAttempts())
	fmt.Printf("   📋 超时时间: %v\n", advancedJob.GetTimeout())
	fmt.Printf("   📋 标签: %v\n", advancedJob.GetTags())

	fmt.Println("\n=== 演示完成 ===")
} 