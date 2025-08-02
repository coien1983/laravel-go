package event

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// EventQueueWorker 事件工作进程实现
type EventQueueWorker struct {
	mu           sync.RWMutex
	queue        EventQueue
	queueName    string
	workerID     string
	status       string
	startedAt    time.Time
	processed    int64
	failed       int64
	currentEvent *Event
	stopChan     chan struct{}
	pauseChan    chan struct{}
	resumeChan   chan struct{}
	onFailed     func(Event, error)
	onProcessed  func(Event)
	timeout      time.Duration
	metrics      *EventWorkerMetrics
}

// NewEventWorker 创建事件工作进程
func NewEventWorker(queue EventQueue, queueName string) *EventQueueWorker {
	return &EventQueueWorker{
		queue:      queue,
		queueName:  queueName,
		workerID:   uuid.New().String(),
		status:     "stopped",
		stopChan:   make(chan struct{}),
		pauseChan:  make(chan struct{}),
		resumeChan: make(chan struct{}),
		timeout:    30 * time.Second,
		metrics: &EventWorkerMetrics{
			LastEventTime: time.Now(),
		},
	}
}

// Start 启动工作进程
func (w *EventQueueWorker) Start() error {
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
func (w *EventQueueWorker) Stop() error {
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
func (w *EventQueueWorker) Pause() error {
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
func (w *EventQueueWorker) Resume() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status != "paused" {
		return fmt.Errorf("worker is not paused")
	}

	w.status = "running"
	w.resumeChan <- struct{}{}
	return nil
}

// Process 处理事件
func (w *EventQueueWorker) Process(event Event) error {
	startTime := time.Now()

	// 设置当前事件
	w.mu.Lock()
	w.currentEvent = &event
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		w.currentEvent = nil
		w.mu.Unlock()
	}()

	// 处理事件
	err := w.processEvent(event)
	if err != nil {
		w.handleFailed(event, err)
		return err
	}

	// 标记为已处理
	w.handleProcessed(event)

	// 更新指标
	w.updateMetrics(time.Since(startTime))

	return nil
}

// HandleFailed 处理失败的事件
func (w *EventQueueWorker) HandleFailed(event Event, err error) error {
	w.handleFailed(event, err)
	return nil
}

// GetStatus 获取工作进程状态
func (w *EventQueueWorker) GetStatus() EventWorkerStatus {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return EventWorkerStatus{
		Status:       w.status,
		StartedAt:    w.startedAt,
		Processed:    w.processed,
		Failed:       w.failed,
		CurrentEvent: w.currentEvent,
		Queue:        w.queueName,
		WorkerID:     w.workerID,
	}
}

// GetMetrics 获取工作进程指标
func (w *EventQueueWorker) GetMetrics() EventWorkerMetrics {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return *w.metrics
}

// SetOnFailed 设置失败回调
func (w *EventQueueWorker) SetOnFailed(callback func(Event, error)) {
	w.onFailed = callback
}

// SetOnProcessed 设置处理完成回调
func (w *EventQueueWorker) SetOnProcessed(callback func(Event)) {
	w.onProcessed = callback
}

// SetTimeout 设置超时时间
func (w *EventQueueWorker) SetTimeout(timeout time.Duration) {
	w.timeout = timeout
}

// run 运行工作进程
func (w *EventQueueWorker) run() {
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
			// 弹出事件
			event, err := w.queue.Pop(ctx)
			if err != nil {
				// 没有事件，等待一段时间
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// 处理事件
			w.Process(event)
		}
	}
}

// processEvent 处理单个事件
func (w *EventQueueWorker) processEvent(event Event) error {
	// 这里应该根据事件类型调用相应的处理器
	// 目前只是一个示例实现

	// 模拟事件处理
	time.Sleep(10 * time.Millisecond)

	// 检查事件载荷
	payload := event.GetPayload()
	if payload == nil {
		return fmt.Errorf("empty event payload")
	}

	// 这里可以添加具体的事件处理逻辑
	// 例如：解析事件类型，调用相应的处理器等

	return nil
}

// handleFailed 处理失败的事件
func (w *EventQueueWorker) handleFailed(event Event, err error) {
	w.mu.Lock()
	w.failed++
	w.mu.Unlock()

	// 调用失败回调
	if w.onFailed != nil {
		w.onFailed(event, err)
	}

	// 记录日志
	log.Printf("Worker %s failed to process event %s: %v", w.workerID, event.GetName(), err)
}

// handleProcessed 处理完成的事件
func (w *EventQueueWorker) handleProcessed(event Event) {
	w.mu.Lock()
	w.processed++
	w.mu.Unlock()

	// 调用处理完成回调
	if w.onProcessed != nil {
		w.onProcessed(event)
	}

	// 记录日志
	log.Printf("Worker %s processed event %s", w.workerID, event.GetName())
}

// updateMetrics 更新指标
func (w *EventQueueWorker) updateMetrics(duration time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.metrics.TotalProcessed++
	w.metrics.LastEventTime = time.Now()

	// 计算平均处理时间
	if w.metrics.TotalProcessed > 1 {
		totalTime := w.metrics.AverageTime * time.Duration(w.metrics.TotalProcessed-1)
		w.metrics.AverageTime = (totalTime + duration) / time.Duration(w.metrics.TotalProcessed)
	} else {
		w.metrics.AverageTime = duration
	}
}

// EventWorkerPool 事件工作进程池
type EventWorkerPool struct {
	workers   []*EventQueueWorker
	queue     EventQueue
	queueName string
	poolSize  int
	mu        sync.RWMutex
}

// NewEventWorkerPool 创建事件工作进程池
func NewEventWorkerPool(queue EventQueue, queueName string, poolSize int) *EventWorkerPool {
	return &EventWorkerPool{
		queue:     queue,
		queueName: queueName,
		poolSize:  poolSize,
		workers:   make([]*EventQueueWorker, 0, poolSize),
	}
}

// Start 启动工作进程池
func (wp *EventWorkerPool) Start() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for i := 0; i < wp.poolSize; i++ {
		worker := NewEventWorker(wp.queue, wp.queueName)
		wp.workers = append(wp.workers, worker)

		if err := worker.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Stop 停止工作进程池
func (wp *EventWorkerPool) Stop() error {
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
func (wp *EventWorkerPool) GetWorkers() []*EventQueueWorker {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	workers := make([]*EventQueueWorker, len(wp.workers))
	copy(workers, wp.workers)
	return workers
}

// GetStats 获取工作进程池统计
func (wp *EventWorkerPool) GetStats() ([]EventWorkerStatus, error) {
	workers := wp.GetWorkers()
	stats := make([]EventWorkerStatus, len(workers))

	for i, worker := range workers {
		stats[i] = worker.GetStatus()
	}

	return stats, nil
}
