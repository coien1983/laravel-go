package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// DistributedWorkerPool 分布式工作进程池
type DistributedWorkerPool struct {
	queue        *DistributedQueue
	workers      []*DistributedWorker
	workerCount  int
	maxConcurrency int
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.RWMutex
	status       string // running, stopped, paused
}

// DistributedWorker 分布式工作进程
type DistributedWorker struct {
	id           string
	queue        *DistributedQueue
	ctx          context.Context
	cancel       context.CancelFunc
	status       string // idle, processing, stopped
	currentJob   Job
	processed    int64
	failed       int64
	startedAt    time.Time
	lastJobAt    time.Time
	onCompleted  func(Job)
	onFailed     func(Job, error)
	mu           sync.RWMutex
}

// NewDistributedWorkerPool 创建分布式工作进程池
func NewDistributedWorkerPool(queue *DistributedQueue, workerCount, maxConcurrency int) *DistributedWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &DistributedWorkerPool{
		queue:         queue,
		workerCount:   workerCount,
		maxConcurrency: maxConcurrency,
		ctx:           ctx,
		cancel:        cancel,
		status:        "stopped",
		workers:       make([]*DistributedWorker, 0, workerCount),
	}

	// 创建工作进程
	for i := 0; i < workerCount; i++ {
		worker := NewDistributedWorker(fmt.Sprintf("worker-%d", i), queue)
		pool.workers = append(pool.workers, worker)
	}

	return pool
}

// Start 启动工作进程池
func (pool *DistributedWorkerPool) Start() error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.status == "running" {
		return fmt.Errorf("worker pool is already running")
	}

	pool.status = "running"

	// 启动所有工作进程
	for _, worker := range pool.workers {
		pool.wg.Add(1)
		go func(w *DistributedWorker) {
			defer pool.wg.Done()
			w.Start()
		}(worker)
	}

	return nil
}

// Stop 停止工作进程池
func (pool *DistributedWorkerPool) Stop() error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.status == "stopped" {
		return nil
	}

	pool.status = "stopped"
	pool.cancel()

	// 等待所有工作进程停止
	pool.wg.Wait()

	return nil
}

// Pause 暂停工作进程池
func (pool *DistributedWorkerPool) Pause() error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.status != "running" {
		return fmt.Errorf("worker pool is not running")
	}

	pool.status = "paused"

	// 暂停所有工作进程
	for _, worker := range pool.workers {
		worker.Pause()
	}

	return nil
}

// Resume 恢复工作进程池
func (pool *DistributedWorkerPool) Resume() error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.status != "paused" {
		return fmt.Errorf("worker pool is not paused")
	}

	pool.status = "running"

	// 恢复所有工作进程
	for _, worker := range pool.workers {
		worker.Resume()
	}

	return nil
}

// GetStatus 获取工作进程池状态
func (pool *DistributedWorkerPool) GetStatus() string {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	return pool.status
}

// GetStats 获取工作进程池统计
func (pool *DistributedWorkerPool) GetStats() WorkerPoolStats {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	var totalProcessed, totalFailed int64
	var activeWorkers int

	for _, worker := range pool.workers {
		stats := worker.GetStats()
		totalProcessed += stats.Processed
		totalFailed += stats.Failed
		if worker.GetStatus() == "processing" {
			activeWorkers++
		}
	}

	return WorkerPoolStats{
		TotalWorkers:   len(pool.workers),
		ActiveWorkers:  activeWorkers,
		IdleWorkers:    len(pool.workers) - activeWorkers,
		TotalProcessed: totalProcessed,
		TotalFailed:    totalFailed,
		Status:         pool.status,
	}
}

// SetOnCompleted 设置任务完成回调
func (pool *DistributedWorkerPool) SetOnCompleted(callback func(Job)) {
	for _, worker := range pool.workers {
		worker.SetOnCompleted(callback)
	}
}

// SetOnFailed 设置任务失败回调
func (pool *DistributedWorkerPool) SetOnFailed(callback func(Job, error)) {
	for _, worker := range pool.workers {
		worker.SetOnFailed(callback)
	}
}

// WorkerPoolStats 工作进程池统计
type WorkerPoolStats struct {
	TotalWorkers   int   `json:"total_workers"`
	ActiveWorkers  int   `json:"active_workers"`
	IdleWorkers    int   `json:"idle_workers"`
	TotalProcessed int64 `json:"total_processed"`
	TotalFailed    int64 `json:"total_failed"`
	Status         string `json:"status"`
}

// NewDistributedWorker 创建分布式工作进程
func NewDistributedWorker(id string, queue *DistributedQueue) *DistributedWorker {
	ctx, cancel := context.WithCancel(context.Background())

	return &DistributedWorker{
		id:        id,
		queue:     queue,
		ctx:       ctx,
		cancel:    cancel,
		status:    "idle",
		startedAt: time.Now(),
	}
}

// Start 启动工作进程
func (w *DistributedWorker) Start() {
	w.mu.Lock()
	w.status = "idle"
	w.mu.Unlock()

	for {
		select {
		case <-w.ctx.Done():
			w.mu.Lock()
			w.status = "stopped"
			w.mu.Unlock()
			return
		default:
			w.processNextJob()
		}
	}
}

// processNextJob 处理下一个任务
func (w *DistributedWorker) processNextJob() {
	// 检查是否为领导者（只有领导者才处理任务）
	if !w.queue.IsLeader() {
		time.Sleep(1 * time.Second)
		return
	}

	// 尝试获取任务
	job, err := w.queue.Pop(w.ctx)
	if err != nil {
		// 没有任务，等待一段时间
		time.Sleep(100 * time.Millisecond)
		return
	}

	// 开始处理任务
	w.mu.Lock()
	w.status = "processing"
	w.currentJob = job
	w.mu.Unlock()

	// 广播任务开始执行
	execution := JobExecution{
		JobID:     job.GetID(),
		NodeID:    w.queue.nodeID,
		Status:    "processing",
		StartedAt: time.Now(),
	}
	w.queue.broadcastJobExecution(execution)

	// 处理任务
	err = w.processJob(job)

	// 更新统计
	w.mu.Lock()
	w.processed++
	w.lastJobAt = time.Now()
	w.currentJob = nil
	w.status = "idle"
	w.mu.Unlock()

	// 广播任务执行完成
	endedAt := time.Now()
	execution.EndedAt = &endedAt
	if err != nil {
		execution.Status = "failed"
		execution.Error = err.Error()
		w.mu.Lock()
		w.failed++
		w.mu.Unlock()
	} else {
		execution.Status = "completed"
	}
	w.queue.broadcastJobExecution(execution)

	// 调用回调
	if err != nil {
		if w.onFailed != nil {
			w.onFailed(job, err)
		}
	} else {
		if w.onCompleted != nil {
			w.onCompleted(job)
		}
	}
}

// processJob 处理单个任务
func (w *DistributedWorker) processJob(job Job) error {
	// 这里应该调用任务处理器
	// 目前只是模拟处理
	time.Sleep(100 * time.Millisecond)
	
	// 模拟随机失败
	if time.Now().UnixNano()%10 == 0 {
		return fmt.Errorf("模拟任务处理失败")
	}

	return nil
}

// Pause 暂停工作进程
func (w *DistributedWorker) Pause() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.status == "processing" {
		w.status = "paused"
	}
}

// Resume 恢复工作进程
func (w *DistributedWorker) Resume() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.status == "paused" {
		w.status = "idle"
	}
}

// GetStatus 获取工作进程状态
func (w *DistributedWorker) GetStatus() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.status
}

// GetStats 获取工作进程统计
func (w *DistributedWorker) GetStats() WorkerStats {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return WorkerStats{
		ID:           w.id,
		Status:       w.status,
		Processed:    w.processed,
		Failed:       w.failed,
		StartedAt:    w.startedAt,
		LastJobAt:    w.lastJobAt,
		CurrentJobID: w.getCurrentJobID(),
	}
}

// SetOnCompleted 设置任务完成回调
func (w *DistributedWorker) SetOnCompleted(callback func(Job)) {
	w.onCompleted = callback
}

// SetOnFailed 设置任务失败回调
func (w *DistributedWorker) SetOnFailed(callback func(Job, error)) {
	w.onFailed = callback
}

// getCurrentJobID 获取当前任务ID
func (w *DistributedWorker) getCurrentJobID() string {
	if w.currentJob != nil {
		return w.currentJob.GetID()
	}
	return ""
}

// WorkerStats 工作进程统计
type WorkerStats struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	Processed    int64     `json:"processed"`
	Failed       int64     `json:"failed"`
	StartedAt    time.Time `json:"started_at"`
	LastJobAt    time.Time `json:"last_job_at"`
	CurrentJobID string    `json:"current_job_id"`
} 