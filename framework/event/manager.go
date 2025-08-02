package event

import (
	"sync"
)

// 全局事件管理器
var (
	GlobalEventManager *EventManager
	once               sync.Once
)

// Init 初始化事件管理器
func Init() {
	once.Do(func() {
		queue := NewMemoryEventQueue()
		dispatcher := NewEventDispatcher(queue)
		GlobalEventManager = NewEventManager(dispatcher, queue)
	})
}

// Listen 监听事件
func Listen(eventName string, listener Listener) {
	if GlobalEventManager == nil {
		Init()
	}
	GlobalEventManager.Listen(eventName, listener)
}

// ListenMany 监听多个事件
func ListenMany(eventNames []string, listener Listener) {
	if GlobalEventManager == nil {
		Init()
	}
	GlobalEventManager.ListenMany(eventNames, listener)
}

// Forget 忘记监听器
func Forget(eventName string, listenerName string) {
	if GlobalEventManager == nil {
		return
	}
	GlobalEventManager.dispatcher.Forget(eventName, listenerName)
}

// ForgetMany 忘记多个事件
func ForgetMany(eventNames []string) {
	if GlobalEventManager == nil {
		return
	}
	GlobalEventManager.dispatcher.ForgetMany(eventNames)
}

// Dispatch 分发事件
func Dispatch(event Event) error {
	if GlobalEventManager == nil {
		Init()
	}
	return GlobalEventManager.Dispatch(event)
}

// DispatchAsync 异步分发事件
func DispatchAsync(event Event) error {
	if GlobalEventManager == nil {
		Init()
	}
	return GlobalEventManager.DispatchAsync(event)
}

// DispatchBatch 批量分发事件
func DispatchBatch(events []Event) error {
	if GlobalEventManager == nil {
		Init()
	}
	return GlobalEventManager.dispatcher.DispatchBatch(events)
}

// Queue 队列事件
func Queue(event Event, queue string) error {
	if GlobalEventManager == nil {
		Init()
	}
	return GlobalEventManager.Queue(event, queue)
}

// QueueBatch 批量队列事件
func QueueBatch(events []Event, queue string) error {
	if GlobalEventManager == nil {
		Init()
	}
	return GlobalEventManager.dispatcher.QueueBatch(events, queue)
}

// Subscribe 订阅事件
func Subscribe(subscriber EventSubscriber) {
	if GlobalEventManager == nil {
		Init()
	}
	GlobalEventManager.Subscribe(subscriber)
}

// Unsubscribe 取消订阅
func Unsubscribe(subscriber EventSubscriber) {
	if GlobalEventManager == nil {
		return
	}
	GlobalEventManager.Unsubscribe(subscriber)
}

// HasListeners 检查是否有监听器
func HasListeners(eventName string) bool {
	if GlobalEventManager == nil {
		return false
	}
	return GlobalEventManager.HasListeners(eventName)
}

// GetListeners 获取监听器
func GetListeners(eventName string) []Listener {
	if GlobalEventManager == nil {
		return make([]Listener, 0)
	}
	return GlobalEventManager.GetListeners(eventName)
}

// GetAllListeners 获取所有监听器
func GetAllListeners() map[string][]Listener {
	if GlobalEventManager == nil {
		return make(map[string][]Listener)
	}
	return GlobalEventManager.GetAllListeners()
}

// GetStats 获取统计信息
func GetStats() EventStats {
	if GlobalEventManager == nil {
		return EventStats{}
	}
	return GlobalEventManager.GetStats()
}

// StartWorker 启动工作进程
func StartWorker(queueName string, worker EventWorker) error {
	if GlobalEventManager == nil {
		Init()
	}
	return GlobalEventManager.StartWorker(queueName, worker)
}

// StopWorker 停止工作进程
func StopWorker(queueName string) error {
	if GlobalEventManager == nil {
		return nil
	}
	return GlobalEventManager.StopWorker(queueName)
}

// GetWorker 获取工作进程
func GetWorker(queueName string) (EventWorker, bool) {
	if GlobalEventManager == nil {
		return nil, false
	}
	return GlobalEventManager.GetWorker(queueName)
}

// GetAllWorkers 获取所有工作进程
func GetAllWorkers() map[string]EventWorker {
	if GlobalEventManager == nil {
		return make(map[string]EventWorker)
	}
	return GlobalEventManager.GetAllWorkers()
}

// Close 关闭事件管理器
func Close() error {
	if GlobalEventManager == nil {
		return nil
	}
	return GlobalEventManager.dispatcher.Close()
}
