package event

import (
	"context"
	"time"
)

// Event 事件接口
type Event interface {
	// 基础信息
	GetName() string
	GetPayload() interface{}
	GetTimestamp() time.Time
	GetID() string

	// 事件属性
	GetData() map[string]interface{}
	SetData(key string, value interface{})
	GetDataByKey(key string) interface{}

	// 事件状态
	IsPropagated() bool
	SetPropagated(propagated bool)

	// 序列化
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}

// Listener 监听器接口
type Listener interface {
	// 处理事件
	Handle(event Event) error

	// 监听器信息
	GetName() string
	GetPriority() int
	ShouldQueue() bool
	GetQueue() string
}

// Dispatcher 事件分发器接口
type Dispatcher interface {
	// 事件监听
	Listen(eventName string, listener Listener)
	ListenMany(eventNames []string, listener Listener)
	Forget(eventName string, listenerName string)
	ForgetMany(eventNames []string)

	// 事件分发
	Dispatch(event Event) error
	DispatchAsync(event Event) error
	DispatchBatch(events []Event) error

	// 事件订阅
	Subscribe(subscriber EventSubscriber)
	Unsubscribe(subscriber EventSubscriber)

	// 事件队列
	Queue(event Event, queue string) error
	QueueBatch(events []Event, queue string) error

	// 监听器管理
	HasListeners(eventName string) bool
	GetListeners(eventName string) []Listener
	GetAllListeners() map[string][]Listener

	// 资源管理
	Close() error
}

// EventSubscriber 事件订阅者接口
type EventSubscriber interface {
	// 订阅事件
	Subscribe(dispatcher Dispatcher)

	// 获取订阅者名称
	GetName() string
}

// EventQueue 事件队列接口
type EventQueue interface {
	// 队列操作
	Push(event Event) error
	PushBatch(events []Event) error
	Pop(ctx context.Context) (Event, error)
	PopBatch(ctx context.Context, count int) ([]Event, error)

	// 队列管理
	Size() (int, error)
	Clear() error
	Close() error
}

// EventWorker 事件工作进程接口
type EventWorker interface {
	// 工作进程管理
	Start() error
	Stop() error
	Pause() error
	Resume() error

	// 事件处理
	Process(event Event) error
	HandleFailed(event Event, err error) error

	// 监控
	GetStatus() EventWorkerStatus
	GetMetrics() EventWorkerMetrics
}

// EventWorkerStatus 事件工作进程状态
type EventWorkerStatus struct {
	Status       string    `json:"status"` // running, paused, stopped
	StartedAt    time.Time `json:"started_at"`
	Processed    int64     `json:"processed"`
	Failed       int64     `json:"failed"`
	CurrentEvent *Event    `json:"current_event"`
	Queue        string    `json:"queue"`
	WorkerID     string    `json:"worker_id"`
}

// EventWorkerMetrics 事件工作进程指标
type EventWorkerMetrics struct {
	TotalProcessed int64         `json:"total_processed"`
	TotalFailed    int64         `json:"total_failed"`
	AverageTime    time.Duration `json:"average_time"`
	LastEventTime  time.Time     `json:"last_event_time"`
	MemoryUsage    int64         `json:"memory_usage"`
	CPUUsage       float64       `json:"cpu_usage"`
}

// EventStats 事件统计信息
type EventStats struct {
	TotalEvents      int64     `json:"total_events"`
	DispatchedEvents int64     `json:"dispatched_events"`
	QueuedEvents     int64     `json:"queued_events"`
	FailedEvents     int64     `json:"failed_events"`
	LastEventAt      time.Time `json:"last_event_at"`
	CreatedAt        time.Time `json:"created_at"`
}

// EventManager 事件管理器
type EventManager struct {
	dispatcher Dispatcher
	queue      EventQueue
	workers    map[string]EventWorker
	stats      *EventStats
}

// NewEventManager 创建事件管理器
func NewEventManager(dispatcher Dispatcher, queue EventQueue) *EventManager {
	return &EventManager{
		dispatcher: dispatcher,
		queue:      queue,
		workers:    make(map[string]EventWorker),
		stats: &EventStats{
			CreatedAt: time.Now(),
		},
	}
}

// Listen 监听事件
func (em *EventManager) Listen(eventName string, listener Listener) {
	em.dispatcher.Listen(eventName, listener)
}

// ListenMany 监听多个事件
func (em *EventManager) ListenMany(eventNames []string, listener Listener) {
	em.dispatcher.ListenMany(eventNames, listener)
}

// Dispatch 分发事件
func (em *EventManager) Dispatch(event Event) error {
	em.stats.TotalEvents++
	em.stats.DispatchedEvents++
	em.stats.LastEventAt = time.Now()
	return em.dispatcher.Dispatch(event)
}

// DispatchAsync 异步分发事件
func (em *EventManager) DispatchAsync(event Event) error {
	em.stats.TotalEvents++
	em.stats.DispatchedEvents++
	em.stats.LastEventAt = time.Now()
	return em.dispatcher.DispatchAsync(event)
}

// Queue 队列事件
func (em *EventManager) Queue(event Event, queue string) error {
	em.stats.TotalEvents++
	em.stats.QueuedEvents++
	em.stats.LastEventAt = time.Now()
	return em.dispatcher.Queue(event, queue)
}

// Subscribe 订阅事件
func (em *EventManager) Subscribe(subscriber EventSubscriber) {
	em.dispatcher.Subscribe(subscriber)
}

// Unsubscribe 取消订阅
func (em *EventManager) Unsubscribe(subscriber EventSubscriber) {
	em.dispatcher.Unsubscribe(subscriber)
}

// HasListeners 检查是否有监听器
func (em *EventManager) HasListeners(eventName string) bool {
	return em.dispatcher.HasListeners(eventName)
}

// GetListeners 获取监听器
func (em *EventManager) GetListeners(eventName string) []Listener {
	return em.dispatcher.GetListeners(eventName)
}

// GetAllListeners 获取所有监听器
func (em *EventManager) GetAllListeners() map[string][]Listener {
	return em.dispatcher.GetAllListeners()
}

// GetStats 获取统计信息
func (em *EventManager) GetStats() EventStats {
	return *em.stats
}

// StartWorker 启动工作进程
func (em *EventManager) StartWorker(queueName string, worker EventWorker) error {
	em.workers[queueName] = worker
	return worker.Start()
}

// StopWorker 停止工作进程
func (em *EventManager) StopWorker(queueName string) error {
	worker, exists := em.workers[queueName]
	if !exists {
		return nil
	}
	return worker.Stop()
}

// GetWorker 获取工作进程
func (em *EventManager) GetWorker(queueName string) (EventWorker, bool) {
	worker, exists := em.workers[queueName]
	return worker, exists
}

// GetAllWorkers 获取所有工作进程
func (em *EventManager) GetAllWorkers() map[string]EventWorker {
	return em.workers
}
