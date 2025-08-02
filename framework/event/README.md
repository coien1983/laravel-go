# Laravel-Go 事件系统

Laravel-Go 事件系统是一个完整的事件驱动架构实现，提供了事件分发、监听、队列处理等功能，支持同步和异步事件处理。

## 功能特性

- ✅ **事件接口**: 统一的事件、监听器、分发器接口
- ✅ **基础实现**: 基础事件、监听器、分发器实现
- ✅ **事件管理**: 全局事件管理器，简化 API 使用
- ✅ **异步处理**: 支持异步事件分发和处理
- ✅ **批量处理**: 支持批量事件分发
- ✅ **事件队列**: 内存队列支持，可扩展其他驱动
- ✅ **事件订阅**: 支持事件订阅者模式
- ✅ **工作进程**: 事件工作进程和进程池
- ✅ **统计监控**: 事件统计和性能监控
- ✅ **错误处理**: 完善的错误处理机制

## 核心组件

### 1. 事件 (Event)

事件是系统中发生的动作或状态变化的表示。

```go
// 基础事件
event := event.NewEvent("user.registered", payload)

// 自定义事件
type UserRegisteredEvent struct {
    *event.BaseEvent
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
```

### 2. 监听器 (Listener)

监听器负责处理特定类型的事件。

```go
// 基础监听器
listener := event.NewListener("email.notification", func(e event.Event) error {
    // 处理事件
    return nil
})

// 带优先级的监听器
listener := event.NewListenerWithPriority("high.priority", 10, handler)

// 队列化监听器
listener := event.NewQueuedListener("queued.notification", "notifications", handler)
```

### 3. 事件分发器 (Dispatcher)

事件分发器负责将事件分发给相应的监听器。

```go
// 创建分发器
queue := event.NewMemoryEventQueue()
dispatcher := event.NewEventDispatcher(queue)

// 监听事件
dispatcher.Listen("user.registered", listener)

// 分发事件
err := dispatcher.Dispatch(event)
```

### 4. 事件管理器 (Manager)

事件管理器提供高级的事件管理功能。

```go
// 创建管理器
manager := event.NewEventManager(dispatcher, queue)

// 监听事件
manager.Listen("user.registered", listener)

// 分发事件
err := manager.Dispatch(event)

// 获取统计信息
stats := manager.GetStats()
```

### 5. 全局事件管理器

提供简化的全局 API。

```go
// 初始化
event.Init()

// 监听事件
event.Listen("user.registered", listener)

// 分发事件
err := event.Dispatch(event)

// 异步分发
err := event.DispatchAsync(event)

// 批量分发
err := event.DispatchBatch(events)
```

## 使用示例

### 基本用法

```go
package main

import (
    "fmt"
    "laravel-go/framework/event"
)

func main() {
    // 初始化事件系统
    event.Init()

    // 创建监听器
    listener := event.NewListener("user.registered", func(e event.Event) error {
        fmt.Printf("用户注册事件: %s\n", e.GetName())
        return nil
    })

    // 监听事件
    event.Listen("user.registered", listener)

    // 创建并分发事件
    userEvent := event.NewEvent("user.registered", map[string]interface{}{
        "user_id": 1,
        "name":    "John Doe",
    })

    err := event.Dispatch(userEvent)
    if err != nil {
        panic(err)
    }
}
```

### 自定义事件

```go
// 定义自定义事件
type OrderCreatedEvent struct {
    *event.BaseEvent
    OrderID   int64   `json:"order_id"`
    UserID    int64   `json:"user_id"`
    Amount    float64 `json:"amount"`
    Products  []string `json:"products"`
}

// 创建事件
orderEvent := &OrderCreatedEvent{
    BaseEvent: event.NewEvent("order.created", nil),
    OrderID:   1001,
    UserID:    1,
    Amount:    299.99,
    Products:  []string{"iPhone 15", "AirPods Pro"},
}

// 分发事件
err := event.Dispatch(orderEvent)
```

### 事件订阅者

```go
// 定义事件订阅者
type OrderEventSubscriber struct {
    name string
}

func (s *OrderEventSubscriber) Subscribe(dispatcher event.Dispatcher) {
    // 订阅订单相关事件
    dispatcher.Listen("order.created", NewEmailNotificationListener())
    dispatcher.Listen("order.created", NewDatabaseLogListener())
    dispatcher.Listen("order.created", NewAnalyticsListener())
}

func (s *OrderEventSubscriber) GetName() string {
    return s.name
}

// 使用订阅者
subscriber := &OrderEventSubscriber{name: "order.subscriber"}
event.Subscribe(subscriber)
```

### 异步事件处理

```go
// 异步分发事件
err := event.DispatchAsync(event)
if err != nil {
    panic(err)
}

// 等待异步处理完成
time.Sleep(100 * time.Millisecond)
```

### 批量事件处理

```go
// 创建批量事件
events := []event.Event{
    event.NewEvent("user.registered", user1),
    event.NewEvent("user.registered", user2),
    event.NewEvent("order.created", order1),
}

// 批量分发
err := event.DispatchBatch(events)
if err != nil {
    panic(err)
}
```

### 事件队列

```go
// 创建队列
queue := event.NewMemoryEventQueue()
defer queue.Close()

// 推送事件到队列
err := queue.Push(event)
if err != nil {
    panic(err)
}

// 从队列弹出事件
ctx := context.Background()
poppedEvent, err := queue.Pop(ctx)
if err != nil {
    panic(err)
}
```

### 事件工作进程

```go
// 创建工作进程
queue := event.NewMemoryEventQueue()
worker := event.NewEventWorker(queue, "notifications")

// 启动工作进程
err := worker.Start()
if err != nil {
    panic(err)
}

// 检查状态
status := worker.GetStatus()
fmt.Printf("工作进程状态: %s\n", status.Status)

// 停止工作进程
err = worker.Stop()
if err != nil {
    panic(err)
}
```

## API 参考

### Event 接口

```go
type Event interface {
    GetName() string
    GetPayload() interface{}
    GetTimestamp() time.Time
    GetID() string
    GetData() map[string]interface{}
    SetData(key string, value interface{})
    GetDataByKey(key string) interface{}
    IsPropagated() bool
    SetPropagated(propagated bool)
    Serialize() ([]byte, error)
    Deserialize(data []byte) error
}
```

### Listener 接口

```go
type Listener interface {
    Handle(event Event) error
    GetName() string
    GetPriority() int
    ShouldQueue() bool
    GetQueue() string
}
```

### Dispatcher 接口

```go
type Dispatcher interface {
    Listen(eventName string, listener Listener)
    ListenMany(eventNames []string, listener Listener)
    Forget(eventName string, listenerName string)
    ForgetMany(eventNames []string)
    Dispatch(event Event) error
    DispatchAsync(event Event) error
    DispatchBatch(events []Event) error
    Subscribe(subscriber EventSubscriber)
    Unsubscribe(subscriber EventSubscriber)
    Queue(event Event, queue string) error
    QueueBatch(events []Event, queue string) error
    HasListeners(eventName string) bool
    GetListeners(eventName string) []Listener
    GetAllListeners() map[string][]Listener
    Close() error
}
```

### 全局函数

```go
// 初始化
func Init()

// 事件监听
func Listen(eventName string, listener Listener)
func ListenMany(eventNames []string, listener Listener)
func Forget(eventName string, listenerName string)
func ForgetMany(eventNames []string)

// 事件分发
func Dispatch(event Event) error
func DispatchAsync(event Event) error
func DispatchBatch(events []Event) error

// 事件队列
func Queue(event Event, queue string) error
func QueueBatch(events []Event, queue string) error

// 事件订阅
func Subscribe(subscriber EventSubscriber)
func Unsubscribe(subscriber EventSubscriber)

// 监听器管理
func HasListeners(eventName string) bool
func GetListeners(eventName string) []Listener
func GetAllListeners() map[string][]Listener

// 统计信息
func GetStats() EventStats

// 工作进程管理
func StartWorker(queueName string, worker EventWorker) error
func StopWorker(queueName string) error
func GetWorker(queueName string) (EventWorker, bool)
func GetAllWorkers() map[string]EventWorker

// 资源清理
func Close() error
```

## 最佳实践

### 1. 事件命名

使用点分隔的命名约定：

```go
// 好的命名
"user.registered"
"order.created"
"payment.completed"
"email.sent"

// 避免的命名
"userRegistered"
"order_created"
"PaymentCompleted"
```

### 2. 事件结构

事件应该包含足够的信息供监听器处理：

```go
type UserRegisteredEvent struct {
    *event.BaseEvent
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    // 包含监听器需要的所有信息
}
```

### 3. 监听器设计

监听器应该专注于单一职责：

```go
// 邮件通知监听器
emailListener := event.NewListener("email.notification", func(e event.Event) error {
    // 只处理邮件发送逻辑
    return sendEmail(e)
})

// 数据库日志监听器
dbListener := event.NewListener("database.log", func(e event.Event) error {
    // 只处理数据库日志记录
    return logToDatabase(e)
})
```

### 4. 错误处理

监听器应该妥善处理错误：

```go
listener := event.NewListener("user.registered", func(e event.Event) error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("监听器发生panic: %v", r)
        }
    }()

    // 处理事件
    if err := processEvent(e); err != nil {
        log.Printf("处理事件失败: %v", err)
        return err
    }

    return nil
})
```

### 5. 性能优化

- 使用异步事件处理耗时操作
- 合理使用事件队列
- 避免在监听器中执行阻塞操作
- 使用批量事件处理提高效率

## 错误处理

事件系统提供了完善的错误处理机制：

```go
// 预定义错误
var (
    ErrEventNotFound        = errors.New("event not found")
    ErrListenerNotFound     = errors.New("listener not found")
    ErrDispatcherClosed     = errors.New("dispatcher is closed")
    ErrEventSerialization   = errors.New("event serialization failed")
    ErrEventDeserialization = errors.New("event deserialization failed")
    ErrInvalidEvent         = errors.New("invalid event")
    ErrInvalidListener      = errors.New("invalid listener")
    ErrEventQueueFull       = errors.New("event queue is full")
    ErrWorkerStopped        = errors.New("worker is stopped")
    ErrEventTimeout         = errors.New("event timeout")
    ErrEventPropagation     = errors.New("event propagation failed")
)

// 自定义错误
type EventError struct {
    EventName string
    Message   string
    Err       error
}

type ListenerError struct {
    ListenerName string
    EventName    string
    Message      string
    Err          error
}
```

## 测试

事件系统包含完整的单元测试：

```bash
# 运行所有测试
go test ./framework/event/... -v

# 运行特定测试
go test ./framework/event/... -run TestEventDispatcher
```

## 示例程序

查看 `examples/event_demo/` 目录获取完整的使用示例。

## 扩展

事件系统设计为可扩展的，可以轻松添加：

- 新的队列驱动 (Redis, RabbitMQ, Kafka 等)
- 新的事件类型
- 新的监听器类型
- 新的事件处理器

## 总结

Laravel-Go 事件系统提供了完整的事件驱动架构支持，具有以下优势：

- **简单易用**: 提供简洁的 API 和全局函数
- **功能完整**: 支持同步、异步、批量事件处理
- **高性能**: 内存队列和异步处理
- **可扩展**: 模块化设计，易于扩展
- **生产就绪**: 完善的错误处理和监控

事件系统是构建松耦合、可维护应用程序的重要工具，特别适用于微服务架构和事件驱动系统。
