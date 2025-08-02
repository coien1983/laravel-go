package queue

import (
	"context"
	"time"
	"sync"
)

// 全局队列管理器
var (
	QueueManager *Manager
	once         sync.Once
)

// Init 初始化队列管理器
func Init() {
	once.Do(func() {
		QueueManager = NewManager()
	})
}

// Push 推送任务到默认队列
func Push(job Job) error {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.Push(job)
}

// PushTo 推送任务到指定队列
func PushTo(queueName string, job Job) error {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.PushTo(queueName, job)
}

// Later 延迟推送任务到默认队列
func Later(job Job, delay time.Duration) error {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.Later(job, delay)
}

// LaterTo 延迟推送任务到指定队列
func LaterTo(queueName string, job Job, delay time.Duration) error {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.LaterTo(queueName, job, delay)
}

// Pop 从默认队列弹出任务
func Pop(ctx context.Context) (Job, error) {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.Pop(ctx)
}

// PopFrom 从指定队列弹出任务
func PopFrom(ctx context.Context, queueName string) (Job, error) {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.PopFrom(ctx, queueName)
}

// Size 获取默认队列大小
func Size() (int, error) {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.Size()
}

// SizeOf 获取指定队列大小
func SizeOf(queueName string) (int, error) {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.SizeOf(queueName)
}

// Clear 清空默认队列
func Clear() error {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.Clear()
}

// ClearQueue 清空指定队列
func ClearQueue(queueName string) error {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.ClearQueue(queueName)
}

// GetStats 获取默认队列统计
func GetStats() (QueueStats, error) {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.GetStats()
}

// GetQueueStats 获取指定队列统计
func GetQueueStats(queueName string) (QueueStats, error) {
	if QueueManager == nil {
		Init()
	}
	return QueueManager.GetQueueStats(queueName)
}

// Close 关闭所有队列
func Close() error {
	if QueueManager == nil {
		return nil
	}
	return QueueManager.Close()
} 