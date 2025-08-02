package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// QueueWorker 工作进程实现
type QueueWorker struct {
	mu           sync.RWMutex
	queue        Queue
	queueName    string
	workerID     string
	status       string
	startedAt    time.Time
	processed    int64
	failed       int64
	currentJob   *Job
	stopChan     chan struct{}
	pauseChan    chan struct{}
	resumeChan   chan struct{}
	onFailed     func(Job, error)
	onCompleted  func(Job)
	timeout      time.Duration
	maxAttempts  int
	metrics      *WorkerMetrics
}

// NewWorker 创建工作进程
func NewWorker(queue Queue, queueName string) *QueueWorker {
	return &QueueWorker{
		queue:       queue,
		queueName:   queueName,
		workerID:    uuid.New().String(),
		status:      "stopped",
		stopChan:    make(chan struct{}),
		pauseChan:   make(chan struct{}),
		resumeChan:  make(chan struct{}),
		timeout:     30 * time.Second,
		maxAttempts: 3,
		metrics: &WorkerMetrics{
			LastJobTime: time.Now(),
		},
	}
}

// Start 启动工作进程
func (w *QueueWorker) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == "running" {
		return fmt.Errorf("worker is already running")
	}

	w.status = "running"
	w.startedAt = time.Now()
	w.stopChan = make(chan struct{})
	w.pauseChan = make(chan struct{})
	w.resumeChan = make(chan struct{})

	go w.run()
	return nil
}

// Stop 停止工作进程
func (w *QueueWorker) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == "stopped" {
		return fmt.Errorf("worker is already stopped")
	}

	w.status = "stopped"
	close(w.stopChan)
	return nil
}

// Pause 暂停工作进程
func (w *QueueWorker) Pause() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status != "running" {
		return fmt.Errorf("worker is not running")
	}

	w.status = "paused"
	w.pauseChan <- struct{}{}
	return nil
}

// Resume 恢复工作进程
func (w *QueueWorker) Resume() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status != "paused" {
		return fmt.Errorf("worker is not paused")
	}

	w.status = "running"
	w.resumeChan <- struct{}{}
	return nil
}

// Process 处理任务
func (w *QueueWorker) Process(job Job) error {
	startTime := time.Now()
	
	// 设置当前任务
	w.mu.Lock()
	w.currentJob = &job
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		w.currentJob = nil
		w.mu.Unlock()
	}()

	// 检查任务是否过期
	if job.(*BaseJob).IsExpired() {
		err := fmt.Errorf("job expired")
		w.handleFailed(job, err)
		return err
	}

	// 检查最大尝试次数
	if job.GetAttempts() >= w.maxAttempts {
		err := fmt.Errorf("job exceeded max attempts")
		w.handleFailed(job, err)
		return err
	}

	// 处理任务
	err := w.processJob(job)
	if err != nil {
		w.handleFailed(job, err)
		return err
	}

	// 标记为完成
	job.(*BaseJob).MarkAsCompleted()
	w.handleCompleted(job)

	// 更新指标
	w.updateMetrics(time.Since(startTime))

	return nil
}

// HandleFailed 处理失败的任务
func (w *QueueWorker) HandleFailed(job Job, err error) error {
	w.handleFailed(job, err)
	return nil
}

// GetStatus 获取工作进程状态
func (w *QueueWorker) GetStatus() WorkerStatus {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return WorkerStatus{
		Status:     w.status,
		StartedAt:  w.startedAt,
		Processed:  w.processed,
		Failed:     w.failed,
		CurrentJob: w.currentJob,
		Queue:      w.queueName,
		WorkerID:   w.workerID,
	}
}

// GetMetrics 获取工作进程指标
func (w *QueueWorker) GetMetrics() WorkerMetrics {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return *w.metrics
}

// SetOnFailed 设置失败回调
func (w *QueueWorker) SetOnFailed(callback func(Job, error)) {
	w.onFailed = callback
}

// SetOnCompleted 设置完成回调
func (w *QueueWorker) SetOnCompleted(callback func(Job)) {
	w.onCompleted = callback
}

// SetTimeout 设置超时时间
func (w *QueueWorker) SetTimeout(timeout time.Duration) {
	w.timeout = timeout
}

// SetMaxAttempts 设置最大尝试次数
func (w *QueueWorker) SetMaxAttempts(maxAttempts int) {
	w.maxAttempts = maxAttempts
}

// run 运行工作进程
func (w *QueueWorker) run() {
	ctx := context.Background()
	
	for {
		select {
		case <-w.stopChan:
			return
		case <-w.pauseChan:
			// 等待恢复信号
			<-w.resumeChan
			continue
		default:
			// 弹出任务
			job, err := w.queue.Pop(ctx)
			if err != nil {
				// 没有任务，等待一段时间
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// 处理任务
			w.Process(job)
		}
	}
}

// processJob 处理单个任务
func (w *QueueWorker) processJob(job Job) error {
	// 这里应该根据任务类型调用相应的处理器
	// 目前只是一个示例实现
	
	// 模拟任务处理
	time.Sleep(10 * time.Millisecond)
	
	// 检查任务载荷
	payload := job.GetPayload()
	if len(payload) == 0 {
		return fmt.Errorf("empty job payload")
	}

	// 这里可以添加具体的任务处理逻辑
	// 例如：解析任务类型，调用相应的处理器等
	
	return nil
}

// handleFailed 处理失败的任务
func (w *QueueWorker) handleFailed(job Job, err error) {
	w.mu.Lock()
	w.failed++
	w.mu.Unlock()

	// 标记任务为失败
	job.(*BaseJob).MarkAsFailed(err)

	// 调用失败回调
	if w.onFailed != nil {
		w.onFailed(job, err)
	}

	// 记录日志
	log.Printf("Worker %s failed to process job %s: %v", w.workerID, job.GetID(), err)
}

// handleCompleted 处理完成的任务
func (w *QueueWorker) handleCompleted(job Job) {
	w.mu.Lock()
	w.processed++
	w.mu.Unlock()

	// 调用完成回调
	if w.onCompleted != nil {
		w.onCompleted(job)
	}

	// 记录日志
	log.Printf("Worker %s completed job %s", w.workerID, job.GetID())
}

// updateMetrics 更新指标
func (w *QueueWorker) updateMetrics(duration time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.metrics.TotalProcessed++
	w.metrics.LastJobTime = time.Now()
	
	// 计算平均处理时间
	if w.metrics.TotalProcessed > 1 {
		totalTime := w.metrics.AverageTime * time.Duration(w.metrics.TotalProcessed-1)
		w.metrics.AverageTime = (totalTime + duration) / time.Duration(w.metrics.TotalProcessed)
	} else {
		w.metrics.AverageTime = duration
	}
}

// WorkerPool 工作进程池
type WorkerPool struct {
	workers []*QueueWorker
	queue   Queue
	queueName string
	poolSize int
	mu      sync.RWMutex
}

// NewWorkerPool 创建工作进程池
func NewWorkerPool(queue Queue, queueName string, poolSize int) *WorkerPool {
	return &WorkerPool{
		queue:     queue,
		queueName: queueName,
		poolSize:  poolSize,
		workers:   make([]*QueueWorker, 0, poolSize),
	}
}

// Start 启动工作进程池
func (wp *WorkerPool) Start() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for i := 0; i < wp.poolSize; i++ {
		worker := NewWorker(wp.queue, wp.queueName)
		wp.workers = append(wp.workers, worker)
		
		if err := worker.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Stop 停止工作进程池
func (wp *WorkerPool) Stop() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for _, worker := range wp.workers {
		if err := worker.Stop(); err != nil {
			return err
		}
	}

	return nil
}

// GetWorkers 获取所有工作进程
func (wp *WorkerPool) GetWorkers() []*QueueWorker {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	workers := make([]*QueueWorker, len(wp.workers))
	copy(workers, wp.workers)
	return workers
}

// GetStats 获取工作进程池统计
func (wp *WorkerPool) GetStats() ([]WorkerStatus, error) {
	workers := wp.GetWorkers()
	stats := make([]WorkerStatus, len(workers))
	
	for i, worker := range workers {
		stats[i] = worker.GetStatus()
	}
	
	return stats, nil
} 