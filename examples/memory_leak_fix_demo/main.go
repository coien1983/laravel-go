package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"laravel-go/framework/database"
	"laravel-go/framework/errors"
	"laravel-go/framework/event"
	"laravel-go/framework/microservice"
	"laravel-go/framework/queue"
)

func main() {
	fmt.Println("=== Laravel-Go Framework 内存泄漏修复演示 ===")

	// 1. 演示事件队列修复
	demoEventQueue()

	// 2. 演示队列系统修复
	demoQueueSystem()

	// 3. 演示数据库连接管理修复
	demoDatabaseConnection()

	// 4. 演示gRPC流管理修复
	demoGRPCStreamManager()

	// 5. 演示panic恢复机制
	demoPanicRecovery()

	fmt.Println("=== 演示完成 ===")
}

// demoEventQueue 演示事件队列修复
func demoEventQueue() {
	fmt.Println("\n1. 事件队列内存泄漏修复演示")

	queue := event.NewMemoryEventQueue()
	defer queue.Close()

	// 启动消费者
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				event, err := queue.Pop(ctx)
				if err != nil {
					if err != context.DeadlineExceeded {
						log.Printf("Pop error: %v", err)
					}
					continue
				}
				fmt.Printf("处理事件: %s\n", event.GetName())
			}
		}
	}()

	// 推送事件
	for i := 0; i < 10; i++ {
		err := queue.Push(&event.BaseEvent{
			Name: fmt.Sprintf("event-%d", i),
			Data: map[string]interface{}{"id": i},
		})
		if err != nil {
			log.Printf("Push error: %v", err)
		}
	}

	time.Sleep(2 * time.Second)
	fmt.Println("事件队列演示完成")
}

// demoQueueSystem 演示队列系统修复
func demoQueueSystem() {
	fmt.Println("\n2. 队列系统内存泄漏修复演示")

	q := queue.NewMemoryQueue()
	defer q.Close()

	// 启动工作进程
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				job, err := q.Pop(ctx)
				if err != nil {
					if err != context.DeadlineExceeded {
						log.Printf("Pop error: %v", err)
					}
					continue
				}

				fmt.Printf("处理任务: %s\n", job.GetID())

				// 模拟处理时间
				time.Sleep(100 * time.Millisecond)

				// 完成任务
				err = q.Delete(job)
				if err != nil {
					log.Printf("Delete error: %v", err)
				}
			}
		}
	}()

	// 推送任务
	for i := 0; i < 5; i++ {
		// 创建任务载荷
		payload := fmt.Sprintf(`{"id": %d, "message": "test job %d"}`, i, i)

		job := queue.NewJob([]byte(payload), "default")
		job.SetMaxAttempts(3)

		err := q.Push(job)
		if err != nil {
			log.Printf("Push error: %v", err)
		}
	}

	time.Sleep(2 * time.Second)
	fmt.Println("队列系统演示完成")
}

// demoDatabaseConnection 演示数据库连接管理修复
func demoDatabaseConnection() {
	fmt.Println("\n3. 数据库连接管理修复演示")

	// 创建连接管理器
	cm := database.NewConnectionManager()
	defer cm.CloseAll()

	// 添加连接配置
	config := &database.ConnectionConfig{
		Driver:          database.SQLite,
		Database:        ":memory:",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		log.Printf("配置验证失败: %v", err)
		return
	}

	cm.AddConnection("test", config)

	// 获取连接（带超时）
	conn, err := cm.GetConnectionWithTimeout("test", 5*time.Second)
	if err != nil {
		log.Printf("获取连接失败: %v", err)
		return
	}

	// 测试连接
	err = conn.Ping()
	if err != nil {
		log.Printf("连接ping失败: %v", err)
		return
	}

	fmt.Println("数据库连接管理演示完成")
}

// demoGRPCStreamManager 演示gRPC流管理修复
func demoGRPCStreamManager() {
	fmt.Println("\n4. gRPC流管理修复演示")

	// 创建流管理器
	sm := microservice.NewStreamManager()
	defer sm.Close()

	// 注册流（带超时）
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sm.RegisterStreamWithTimeout("test-stream", "bidirectional", ctx, 2*time.Minute)

	// 模拟流活动
	go func() {
		time.Sleep(1 * time.Second)
		sm.UnregisterStream("test-stream")
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("gRPC流管理演示完成")
}

// demoPanicRecovery 演示panic恢复机制
func demoPanicRecovery() {
	fmt.Println("\n5. Panic恢复机制演示")

	// 测试SafeExecute
	fmt.Println("测试SafeExecute...")
	err := errors.SafeExecute(func() error {
		panic("测试panic")
	})

	if err != nil {
		fmt.Printf("成功捕获panic: %v\n", err)
	}

	// 测试SafeExecuteWithContext
	fmt.Println("测试SafeExecuteWithContext...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = errors.SafeExecuteWithContext(ctx, func() error {
		panic("测试带上下文的panic")
	})

	if err != nil {
		fmt.Printf("成功捕获带上下文的panic: %v\n", err)
	}

	// 测试上下文取消
	fmt.Println("测试上下文取消...")
	ctx, cancel = context.WithCancel(context.Background())
	cancel()

	err = errors.SafeExecuteWithContext(ctx, func() error {
		return nil
	})

	if err == context.Canceled {
		fmt.Println("成功处理上下文取消")
	}

	fmt.Println("Panic恢复机制演示完成")
}
