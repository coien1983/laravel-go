package tests

import (
	"context"
	"testing"
	"time"

	"laravel-go/framework/errors"
	"laravel-go/framework/event"
	"laravel-go/framework/queue"
)

// TestMemoryEventQueueLeak 测试事件队列内存泄漏
func TestMemoryEventQueueLeak(t *testing.T) {
	queue := event.NewMemoryEventQueue()
	defer queue.Close()

	// 创建上下文，5秒后超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动消费者协程
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, err := queue.Pop(ctx)
				if err != nil && err != context.DeadlineExceeded {
					t.Errorf("Pop error: %v", err)
				}
			}
		}
	}()

	// 推送一些事件
	for i := 0; i < 100; i++ {
		err := queue.Push(&event.BaseEvent{
			Name: "test",
			Data: map[string]interface{}{"id": i},
		})
		if err != nil {
			t.Errorf("Push error: %v", err)
		}
	}

	// 等待一段时间
	time.Sleep(2 * time.Second)

	// 检查队列大小
	size, err := queue.Size()
	if err != nil {
		t.Errorf("Size error: %v", err)
	}

	// 队列应该被清空
	if size > 0 {
		t.Errorf("Expected queue to be empty, got size: %d", size)
	}
}

// TestMemoryQueueLeak 测试队列系统内存泄漏
func TestMemoryQueueLeak(t *testing.T) {
	q := queue.NewMemoryQueue()
	defer q.Close()

	// 创建上下文，5秒后超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动消费者协程
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				job, err := q.Pop(ctx)
				if err != nil && err != context.DeadlineExceeded {
					t.Errorf("Pop error: %v", err)
					continue
				}

				if job != nil {
					// 模拟处理任务
					time.Sleep(10 * time.Millisecond)

					// 完成任务
					err = q.Delete(job)
					if err != nil {
						t.Errorf("Delete error: %v", err)
					}
				}
			}
		}
	}()

	// 推送一些任务
	for i := 0; i < 50; i++ {
		job := &queue.BaseJob{
			ID:          queue.GenerateJobID(),
			Name:        "test-job",
			Data:        map[string]interface{}{"id": i},
			Attempts:    0,
			MaxAttempts: 3,
		}

		err := q.Push(job)
		if err != nil {
			t.Errorf("Push error: %v", err)
		}
	}

	// 等待一段时间
	time.Sleep(3 * time.Second)

	// 检查队列统计
	stats, err := q.GetStats()
	if err != nil {
		t.Errorf("GetStats error: %v", err)
	}

	// 所有任务应该被处理完
	if stats.PendingJobs > 0 || stats.ReservedJobs > 0 {
		t.Errorf("Expected all jobs to be processed, got pending: %d, reserved: %d",
			stats.PendingJobs, stats.ReservedJobs)
	}
}

// TestDatabaseConnectionLeak 测试数据库连接泄漏
func TestDatabaseConnectionLeak(t *testing.T) {
	// 这个测试需要实际的数据库连接
	// 在实际环境中运行
	t.Skip("Skipping database connection leak test - requires actual database")
}

// TestPanicRecovery 测试panic恢复机制
func TestPanicRecovery(t *testing.T) {
	// 测试SafeExecute
	err := errors.SafeExecute(func() error {
		panic("test panic")
	})

	if err == nil {
		t.Error("Expected error from panic, got nil")
	}

	// 测试SafeExecuteWithContext
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = errors.SafeExecuteWithContext(ctx, func() error {
		panic("test panic with context")
	})

	if err == nil {
		t.Error("Expected error from panic with context, got nil")
	}

	// 测试上下文取消
	ctx, cancel = context.WithCancel(context.Background())
	cancel() // 立即取消

	err = errors.SafeExecuteWithContext(ctx, func() error {
		return nil
	})

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got: %v", err)
	}
}
