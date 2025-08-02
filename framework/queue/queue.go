package queue

import (
	"context"
	"time"
)

// Job 表示一个队列任务
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

// Queue 表示队列接口
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

// QueueStats 队列统计信息
type QueueStats struct {
	TotalJobs     int64     `json:"total_jobs"`
	PendingJobs   int64     `json:"pending_jobs"`
	ReservedJobs  int64     `json:"reserved_jobs"`
	FailedJobs    int64     `json:"failed_jobs"`
	CompletedJobs int64     `json:"completed_jobs"`
	LastJobAt     time.Time `json:"last_job_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// Worker 工作进程接口
type Worker interface {
	// 工作进程管理
	Start() error
	Stop() error
	Pause() error
	Resume() error
	
	// 任务处理
	Process(job Job) error
	HandleFailed(job Job, err error) error
	
	// 监控
	GetStatus() WorkerStatus
	GetMetrics() WorkerMetrics
}

// WorkerStatus 工作进程状态
type WorkerStatus struct {
	Status      string    `json:"status"`       // running, paused, stopped
	StartedAt   time.Time `json:"started_at"`
	Processed   int64     `json:"processed"`
	Failed      int64     `json:"failed"`
	CurrentJob  *Job      `json:"current_job"`
	Queue       string    `json:"queue"`
	WorkerID    string    `json:"worker_id"`
}

// WorkerMetrics 工作进程指标
type WorkerMetrics struct {
	TotalProcessed int64         `json:"total_processed"`
	TotalFailed    int64         `json:"total_failed"`
	AverageTime    time.Duration `json:"average_time"`
	LastJobTime    time.Time     `json:"last_job_time"`
	MemoryUsage    int64         `json:"memory_usage"`
	CPUUsage       float64       `json:"cpu_usage"`
}

// Manager 队列管理器
type Manager struct {
	queues map[string]Queue
	defaultQueue string
}

// NewManager 创建队列管理器
func NewManager() *Manager {
	return &Manager{
		queues: make(map[string]Queue),
		defaultQueue: "default",
	}
}

// Extend 扩展队列驱动
func (m *Manager) Extend(name string, queue Queue) {
	m.queues[name] = queue
}

// SetDefaultQueue 设置默认队列
func (m *Manager) SetDefaultQueue(name string) {
	m.defaultQueue = name
}

// GetQueue 获取队列实例
func (m *Manager) GetQueue(name string) (Queue, error) {
	if name == "" {
		name = m.defaultQueue
	}
	
	queue, exists := m.queues[name]
	if !exists {
		return nil, ErrQueueNotFound
	}
	
	return queue, nil
}

// Push 推送任务到默认队列
func (m *Manager) Push(job Job) error {
	queue, err := m.GetQueue("")
	if err != nil {
		return err
	}
	return queue.Push(job)
}

// PushTo 推送任务到指定队列
func (m *Manager) PushTo(queueName string, job Job) error {
	queue, err := m.GetQueue(queueName)
	if err != nil {
		return err
	}
	return queue.Push(job)
}

// Later 延迟推送任务到默认队列
func (m *Manager) Later(job Job, delay time.Duration) error {
	queue, err := m.GetQueue("")
	if err != nil {
		return err
	}
	return queue.Later(job, delay)
}

// LaterTo 延迟推送任务到指定队列
func (m *Manager) LaterTo(queueName string, job Job, delay time.Duration) error {
	queue, err := m.GetQueue(queueName)
	if err != nil {
		return err
	}
	return queue.Later(job, delay)
}

// Pop 从默认队列弹出任务
func (m *Manager) Pop(ctx context.Context) (Job, error) {
	queue, err := m.GetQueue("")
	if err != nil {
		return nil, err
	}
	return queue.Pop(ctx)
}

// PopFrom 从指定队列弹出任务
func (m *Manager) PopFrom(ctx context.Context, queueName string) (Job, error) {
	queue, err := m.GetQueue(queueName)
	if err != nil {
		return nil, err
	}
	return queue.Pop(ctx)
}

// Size 获取默认队列大小
func (m *Manager) Size() (int, error) {
	queue, err := m.GetQueue("")
	if err != nil {
		return 0, err
	}
	return queue.Size()
}

// SizeOf 获取指定队列大小
func (m *Manager) SizeOf(queueName string) (int, error) {
	queue, err := m.GetQueue(queueName)
	if err != nil {
		return 0, err
	}
	return queue.Size()
}

// Clear 清空默认队列
func (m *Manager) Clear() error {
	queue, err := m.GetQueue("")
	if err != nil {
		return err
	}
	return queue.Clear()
}

// ClearQueue 清空指定队列
func (m *Manager) ClearQueue(queueName string) error {
	queue, err := m.GetQueue(queueName)
	if err != nil {
		return err
	}
	return queue.Clear()
}

// GetStats 获取默认队列统计
func (m *Manager) GetStats() (QueueStats, error) {
	queue, err := m.GetQueue("")
	if err != nil {
		return QueueStats{}, err
	}
	return queue.GetStats()
}

// GetQueueStats 获取指定队列统计
func (m *Manager) GetQueueStats(queueName string) (QueueStats, error) {
	queue, err := m.GetQueue(queueName)
	if err != nil {
		return QueueStats{}, err
	}
	return queue.GetStats()
}

// Close 关闭所有队列
func (m *Manager) Close() error {
	for _, queue := range m.queues {
		if err := queue.Close(); err != nil {
			return err
		}
	}
	return nil
} 