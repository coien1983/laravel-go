package event

import (
	"context"
	"testing"
	"time"
)

func TestBaseEvent(t *testing.T) {
	// 测试创建事件
	payload := map[string]interface{}{"key": "value"}
	event := NewEvent("test.event", payload)

	// 测试基础属性
	if event.GetName() != "test.event" {
		t.Errorf("Expected event name 'test.event', got '%s'", event.GetName())
	}

	// 注意：map不能直接比较，这里只检查是否为nil
	if event.GetPayload() == nil {
		t.Error("Expected payload to not be nil")
	}

	if event.GetID() == "" {
		t.Error("Expected event ID to be set")
	}

	// 测试时间戳
	if event.GetTimestamp().IsZero() {
		t.Error("Expected timestamp to be set")
	}

	// 测试数据操作
	event.SetData("custom_key", "custom_value")
	if event.GetDataByKey("custom_key") != "custom_value" {
		t.Errorf("Expected custom_value, got %v", event.GetDataByKey("custom_key"))
	}

	// 测试传播状态
	if event.IsPropagated() {
		t.Error("Expected event to not be propagated initially")
	}

	event.SetPropagated(true)
	if !event.IsPropagated() {
		t.Error("Expected event to be propagated after setting")
	}

	// 测试序列化
	data, err := event.Serialize()
	if err != nil {
		t.Errorf("Failed to serialize event: %v", err)
	}

	// 测试反序列化
	newEvent := &BaseEvent{}
	err = newEvent.Deserialize(data)
	if err != nil {
		t.Errorf("Failed to deserialize event: %v", err)
	}

	if newEvent.GetName() != event.GetName() {
		t.Errorf("Expected name %s, got %s", event.GetName(), newEvent.GetName())
	}
}

func TestBaseListener(t *testing.T) {
	// 测试基础监听器
	handler := func(event Event) error {
		return nil
	}
	listener := NewListener("test.listener", handler)
	
	if listener.GetName() != "test.listener" {
		t.Errorf("Expected listener name 'test.listener', got '%s'", listener.GetName())
	}

	if listener.GetPriority() != 0 {
		t.Errorf("Expected priority 0, got %d", listener.GetPriority())
	}

	if listener.ShouldQueue() {
		t.Error("Expected listener to not be queued")
	}

	if listener.GetQueue() != "" {
		t.Errorf("Expected empty queue, got '%s'", listener.GetQueue())
	}

	// 测试处理事件
	event := NewEvent("test.event", nil)
	err := listener.Handle(event)
	if err != nil {
		t.Errorf("Failed to handle event: %v", err)
	}
}

func TestEventDispatcher(t *testing.T) {
	// 创建队列和分发器
	queue := NewMemoryEventQueue()
	dispatcher := NewEventDispatcher(queue)
	defer dispatcher.Close()

	// 创建测试监听器
	handled := false
	listener := NewListener("test.listener", func(event Event) error {
		handled = true
		return nil
	})

	// 监听事件
	dispatcher.Listen("test.event", listener)

	// 检查是否有监听器
	if !dispatcher.HasListeners("test.event") {
		t.Error("Expected to have listeners for test.event")
	}

	// 分发事件
	event := NewEvent("test.event", nil)
	err := dispatcher.Dispatch(event)
	if err != nil {
		t.Errorf("Failed to dispatch event: %v", err)
	}

	// 检查事件是否被处理
	if !handled {
		t.Error("Expected event to be handled")
	}

	// 测试异步分发
	handled = false
	err = dispatcher.DispatchAsync(event)
	if err != nil {
		t.Errorf("Failed to dispatch async event: %v", err)
	}

	// 等待异步处理
	time.Sleep(100 * time.Millisecond)
	if !handled {
		t.Error("Expected async event to be handled")
	}

	// 测试批量分发
	handled = false
	events := []Event{
		NewEvent("test.event", nil),
		NewEvent("test.event", nil),
	}
	err = dispatcher.DispatchBatch(events)
	if err != nil {
		t.Errorf("Failed to dispatch batch events: %v", err)
	}

	if !handled {
		t.Error("Expected batch events to be handled")
	}

	// 测试忘记监听器
	dispatcher.Forget("test.event", "test.listener")
	if dispatcher.HasListeners("test.event") {
		t.Error("Expected no listeners after forgetting")
	}
}

func TestEventQueue(t *testing.T) {
	// 创建内存队列
	queue := NewMemoryEventQueue()
	defer queue.Close()

	// 测试推送事件
	event := NewEvent("test.event", nil)
	err := queue.Push(event)
	if err != nil {
		t.Errorf("Failed to push event: %v", err)
	}

	// 检查队列大小
	size, err := queue.Size()
	if err != nil {
		t.Errorf("Failed to get queue size: %v", err)
	}
	if size != 1 {
		t.Errorf("Expected queue size 1, got %d", size)
	}

	// 测试弹出事件
	ctx := context.Background()
	poppedEvent, err := queue.Pop(ctx)
	if err != nil {
		t.Errorf("Failed to pop event: %v", err)
	}

	if poppedEvent.GetName() != event.GetName() {
		t.Errorf("Expected event name %s, got %s", event.GetName(), poppedEvent.GetName())
	}

	// 测试批量操作
	events := []Event{
		NewEvent("test.event1", nil),
		NewEvent("test.event2", nil),
	}
	err = queue.PushBatch(events)
	if err != nil {
		t.Errorf("Failed to push batch events: %v", err)
	}

	poppedEvents, err := queue.PopBatch(ctx, 2)
	if err != nil {
		t.Errorf("Failed to pop batch events: %v", err)
	}

	if len(poppedEvents) != 2 {
		t.Errorf("Expected 2 events, got %d", len(poppedEvents))
	}

	// 测试清空队列
	err = queue.Clear()
	if err != nil {
		t.Errorf("Failed to clear queue: %v", err)
	}

	size, err = queue.Size()
	if err != nil {
		t.Errorf("Failed to get queue size: %v", err)
	}
	if size != 0 {
		t.Errorf("Expected queue size 0 after clear, got %d", size)
	}
}

func TestEventWorker(t *testing.T) {
	// 创建队列和工作进程
	queue := NewMemoryEventQueue()
	worker := NewEventWorker(queue, "test.queue")
	defer worker.Stop()
	defer queue.Close()

	// 启动工作进程
	err := worker.Start()
	if err != nil {
		t.Errorf("Failed to start worker: %v", err)
	}

	// 检查工作进程状态
	status := worker.GetStatus()
	if status.Status != "running" {
		t.Errorf("Expected worker status 'running', got '%s'", status.Status)
	}

	// 测试暂停和恢复
	err = worker.Pause()
	if err != nil {
		t.Errorf("Failed to pause worker: %v", err)
	}

	status = worker.GetStatus()
	if status.Status != "paused" {
		t.Errorf("Expected worker status 'paused', got '%s'", status.Status)
	}

	err = worker.Resume()
	if err != nil {
		t.Errorf("Failed to resume worker: %v", err)
	}

	status = worker.GetStatus()
	if status.Status != "running" {
		t.Errorf("Expected worker status 'running', got '%s'", status.Status)
	}
}

func TestEventManager(t *testing.T) {
	// 创建事件管理器
	queue := NewMemoryEventQueue()
	dispatcher := NewEventDispatcher(queue)
	manager := NewEventManager(dispatcher, queue)
	defer manager.dispatcher.Close()
	defer queue.Close()

	// 创建测试监听器
	handled := false
	listener := NewListener("test.listener", func(event Event) error {
		handled = true
		return nil
	})

	// 监听事件
	manager.Listen("test.event", listener)

	// 分发事件
	event := NewEvent("test.event", nil)
	err := manager.Dispatch(event)
	if err != nil {
		t.Errorf("Failed to dispatch event: %v", err)
	}

	if !handled {
		t.Error("Expected event to be handled")
	}

	// 检查统计信息
	stats := manager.GetStats()
	if stats.TotalEvents != 1 {
		t.Errorf("Expected total events 1, got %d", stats.TotalEvents)
	}

	if stats.DispatchedEvents != 1 {
		t.Errorf("Expected dispatched events 1, got %d", stats.DispatchedEvents)
	}

	// 测试工作进程管理
	worker := NewEventWorker(queue, "test.queue")

	err = manager.StartWorker("test.queue", worker)
	if err != nil {
		t.Errorf("Failed to start worker: %v", err)
	}

	// 检查工作进程
	retrievedWorker, exists := manager.GetWorker("test.queue")
	if !exists {
		t.Error("Expected worker to exist")
	}

	if retrievedWorker != worker {
		t.Error("Expected retrieved worker to match")
	}

	// 停止工作进程
	err = manager.StopWorker("test.queue")
	if err != nil {
		t.Errorf("Failed to stop worker: %v", err)
	}
}

func TestGlobalEventManager(t *testing.T) {
	// 测试全局事件管理器
	handled := false
	listener := NewListener("test.listener", func(event Event) error {
		handled = true
		return nil
	})

	// 监听事件
	Listen("test.event", listener)

	// 分发事件
	event := NewEvent("test.event", nil)
	err := Dispatch(event)
	if err != nil {
		t.Errorf("Failed to dispatch event: %v", err)
	}

	if !handled {
		t.Error("Expected event to be handled")
	}

	// 检查监听器
	if !HasListeners("test.event") {
		t.Error("Expected to have listeners")
	}

	listeners := GetListeners("test.event")
	if len(listeners) != 1 {
		t.Errorf("Expected 1 listener, got %d", len(listeners))
	}

	// 测试统计信息
	stats := GetStats()
	if stats.TotalEvents != 1 {
		t.Errorf("Expected total events 1, got %d", stats.TotalEvents)
	}

	// 清理
	Close()
}

func TestEventErrors(t *testing.T) {
	// 测试预定义错误
	if ErrEventNotFound.Error() != "event not found" {
		t.Errorf("Expected error message 'event not found', got '%s'", ErrEventNotFound.Error())
	}

	if ErrListenerNotFound.Error() != "listener not found" {
		t.Errorf("Expected error message 'listener not found', got '%s'", ErrListenerNotFound.Error())
	}

	if ErrDispatcherClosed.Error() != "dispatcher is closed" {
		t.Errorf("Expected error message 'dispatcher is closed', got '%s'", ErrDispatcherClosed.Error())
	}

	// 测试自定义错误
	eventErr := &EventError{
		EventName: "test.event",
		Message:   "test error",
	}
	if eventErr.Error() != "event error [test.event]: test error" {
		t.Errorf("Expected error message 'event error [test.event]: test error', got '%s'", eventErr.Error())
	}

	listenerErr := &ListenerError{
		ListenerName: "test.listener",
		EventName:    "test.event",
		Message:      "test error",
	}
	if listenerErr.Error() != "listener error [test.listener] for event [test.event]: test error" {
		t.Errorf("Expected error message 'listener error [test.listener] for event [test.event]: test error', got '%s'", listenerErr.Error())
	}
} 