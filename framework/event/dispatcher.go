package event

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
)

// EventDispatcher 事件分发器实现
type EventDispatcher struct {
	mu          sync.RWMutex
	listeners   map[string][]Listener
	subscribers map[string]EventSubscriber
	queue       EventQueue
	closed      bool
	asyncChan   chan Event
	workerCount int
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewEventDispatcher 创建事件分发器
func NewEventDispatcher(queue EventQueue) *EventDispatcher {
	ctx, cancel := context.WithCancel(context.Background())

	dispatcher := &EventDispatcher{
		listeners:   make(map[string][]Listener),
		subscribers: make(map[string]EventSubscriber),
		queue:       queue,
		asyncChan:   make(chan Event, 1000),
		workerCount: 5,
		ctx:         ctx,
		cancel:      cancel,
	}

	// 启动异步工作进程
	go dispatcher.startAsyncWorkers()

	return dispatcher
}

// Listen 监听事件
func (d *EventDispatcher) Listen(eventName string, listener Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return
	}

	if d.listeners[eventName] == nil {
		d.listeners[eventName] = make([]Listener, 0)
	}

	d.listeners[eventName] = append(d.listeners[eventName], listener)

	// 按优先级排序
	sort.Slice(d.listeners[eventName], func(i, j int) bool {
		return d.listeners[eventName][i].GetPriority() > d.listeners[eventName][j].GetPriority()
	})
}

// ListenMany 监听多个事件
func (d *EventDispatcher) ListenMany(eventNames []string, listener Listener) {
	for _, eventName := range eventNames {
		d.Listen(eventName, listener)
	}
}

// Forget 忘记监听器
func (d *EventDispatcher) Forget(eventName string, listenerName string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return
	}

	listeners, exists := d.listeners[eventName]
	if !exists {
		return
	}

	// 移除指定名称的监听器
	newListeners := make([]Listener, 0)
	for _, listener := range listeners {
		if listener.GetName() != listenerName {
			newListeners = append(newListeners, listener)
		}
	}

	d.listeners[eventName] = newListeners
}

// ForgetMany 忘记多个事件
func (d *EventDispatcher) ForgetMany(eventNames []string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return
	}

	for _, eventName := range eventNames {
		delete(d.listeners, eventName)
	}
}

// Dispatch 分发事件
func (d *EventDispatcher) Dispatch(event Event) error {
	if d.closed {
		return ErrDispatcherClosed
	}

	// 标记事件为已传播
	event.SetPropagated(true)

	// 获取监听器
	listeners := d.getListeners(event.GetName())

	// 处理队列监听器
	queuedListeners := make([]Listener, 0)
	syncListeners := make([]Listener, 0)

	for _, listener := range listeners {
		if listener.ShouldQueue() {
			queuedListeners = append(queuedListeners, listener)
		} else {
			syncListeners = append(syncListeners, listener)
		}
	}

	// 同步处理监听器
	for _, listener := range syncListeners {
		if err := d.handleListener(listener, event); err != nil {
			log.Printf("Listener %s failed to handle event %s: %v", listener.GetName(), event.GetName(), err)
		}
	}

	// 队列化监听器
	for _, listener := range queuedListeners {
		if err := d.queueListener(listener, event); err != nil {
			log.Printf("Failed to queue listener %s for event %s: %v", listener.GetName(), event.GetName(), err)
		}
	}

	return nil
}

// DispatchAsync 异步分发事件
func (d *EventDispatcher) DispatchAsync(event Event) error {
	if d.closed {
		return ErrDispatcherClosed
	}

	select {
	case d.asyncChan <- event:
		return nil
	default:
		return ErrEventQueueFull
	}
}

// DispatchBatch 批量分发事件
func (d *EventDispatcher) DispatchBatch(events []Event) error {
	for _, event := range events {
		if err := d.Dispatch(event); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe 订阅事件
func (d *EventDispatcher) Subscribe(subscriber EventSubscriber) {
	d.mu.Lock()
	if d.closed {
		d.mu.Unlock()
		return
	}
	d.subscribers[subscriber.GetName()] = subscriber
	d.mu.Unlock()

	// 在锁外调用订阅者的Subscribe方法，避免死锁
	subscriber.Subscribe(d)
}

// Unsubscribe 取消订阅
func (d *EventDispatcher) Unsubscribe(subscriber EventSubscriber) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return
	}

	delete(d.subscribers, subscriber.GetName())
}

// Queue 队列事件
func (d *EventDispatcher) Queue(event Event, queue string) error {
	if d.closed {
		return ErrDispatcherClosed
	}

	if d.queue == nil {
		return fmt.Errorf("event queue not configured")
	}

	return d.queue.Push(event)
}

// QueueBatch 批量队列事件
func (d *EventDispatcher) QueueBatch(events []Event, queue string) error {
	if d.closed {
		return ErrDispatcherClosed
	}

	if d.queue == nil {
		return fmt.Errorf("event queue not configured")
	}

	return d.queue.PushBatch(events)
}

// HasListeners 检查是否有监听器
func (d *EventDispatcher) HasListeners(eventName string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	listeners, exists := d.listeners[eventName]
	return exists && len(listeners) > 0
}

// GetListeners 获取监听器
func (d *EventDispatcher) GetListeners(eventName string) []Listener {
	d.mu.RLock()
	defer d.mu.RUnlock()

	listeners, exists := d.listeners[eventName]
	if !exists {
		return make([]Listener, 0)
	}

	// 返回副本
	result := make([]Listener, len(listeners))
	copy(result, listeners)
	return result
}

// GetAllListeners 获取所有监听器
func (d *EventDispatcher) GetAllListeners() map[string][]Listener {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make(map[string][]Listener)
	for eventName, listeners := range d.listeners {
		result[eventName] = make([]Listener, len(listeners))
		copy(result[eventName], listeners)
	}
	return result
}

// Close 关闭分发器
func (d *EventDispatcher) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return nil
	}

	d.closed = true
	d.cancel()
	close(d.asyncChan)

	return nil
}

// SetWorkerCount 设置工作进程数量
func (d *EventDispatcher) SetWorkerCount(count int) {
	d.workerCount = count
}

// getListeners 获取监听器（内部方法）
func (d *EventDispatcher) getListeners(eventName string) []Listener {
	d.mu.RLock()
	defer d.mu.RUnlock()

	listeners, exists := d.listeners[eventName]
	if !exists {
		return make([]Listener, 0)
	}

	result := make([]Listener, len(listeners))
	copy(result, listeners)
	return result
}

// handleListener 处理监听器（内部方法）
func (d *EventDispatcher) handleListener(listener Listener, event Event) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Listener %s panicked while handling event %s: %v", listener.GetName(), event.GetName(), r)
		}
	}()

	return listener.Handle(event)
}

// queueListener 队列化监听器（内部方法）
func (d *EventDispatcher) queueListener(listener Listener, event Event) error {
	if d.queue == nil {
		return fmt.Errorf("event queue not configured")
	}

	// 创建队列事件
	queueEvent := NewEvent("queued.listener", map[string]interface{}{
		"listener_name": listener.GetName(),
		"event":         event,
	})

	return d.queue.Push(queueEvent)
}

// startAsyncWorkers 启动异步工作进程（内部方法）
func (d *EventDispatcher) startAsyncWorkers() {
	for i := 0; i < d.workerCount; i++ {
		go d.asyncWorker(i)
	}
}

// asyncWorker 异步工作进程（内部方法）
func (d *EventDispatcher) asyncWorker(id int) {
	for {
		select {
		case event, ok := <-d.asyncChan:
			if !ok {
				return
			}

			if err := d.Dispatch(event); err != nil {
				log.Printf("Async worker %d failed to dispatch event %s: %v", id, event.GetName(), err)
			}

		case <-d.ctx.Done():
			return
		}
	}
}
