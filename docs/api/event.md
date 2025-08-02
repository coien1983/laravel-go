# 事件系统 API 参考

## 📋 概述

Laravel-Go Framework 的事件系统提供了强大的事件驱动架构，支持事件的注册、分发、监听和处理。事件系统基于发布-订阅模式，实现了组件间的松耦合通信。

## 🏗️ 核心概念

### 事件 (Event)

- 应用程序中发生的动作或状态变化
- 包含相关的数据信息
- 可以被多个监听器处理

### 监听器 (Listener)

- 响应特定事件的处理器
- 执行具体的业务逻辑
- 可以异步或同步执行

### 事件分发器 (Dispatcher)

- 管理事件的注册和分发
- 协调事件和监听器的关系
- 提供事件队列支持

## 🔧 基础用法

### 1. 事件定义

```go
// app/Events/UserRegistered.go
package events

import (
    "laravel-go/framework/event"
    "laravel-go/app/Models"
)

type UserRegistered struct {
    event.BaseEvent
    User *Models.User
}

func NewUserRegistered(user *Models.User) *UserRegistered {
    return &UserRegistered{
        BaseEvent: event.BaseEvent{
            Name: "user.registered",
        },
        User: user,
    }
}

func (e *UserRegistered) GetData() interface{} {
    return e.User
}
```

### 2. 监听器定义

```go
// app/Listeners/SendWelcomeEmail.go
package listeners

import (
    "laravel-go/framework/event"
    "laravel-go/app/Events"
)

type SendWelcomeEmail struct {
    event.BaseListener
}

func NewSendWelcomeEmail() *SendWelcomeEmail {
    return &SendWelcomeEmail{
        BaseListener: event.BaseListener{
            Event: "user.registered",
        },
    }
}

func (l *SendWelcomeEmail) Handle(event event.Event) error {
    userRegistered := event.(*Events.UserRegistered)
    user := userRegistered.User

    // 发送欢迎邮件
    return sendWelcomeEmail(user.Email, user.Name)
}

func (l *SendWelcomeEmail) ShouldQueue() bool {
    return true // 异步执行
}

func (l *SendWelcomeEmail) GetQueue() string {
    return "emails"
}
```

### 3. 事件注册

```go
// app/Providers/EventServiceProvider.go
package providers

import (
    "laravel-go/framework/event"
    "laravel-go/app/Events"
    "laravel-go/app/Listeners"
)

type EventServiceProvider struct {
    event.ServiceProvider
}

func (p *EventServiceProvider) Register() {
    // 注册事件监听器
    p.Listen(&Events.UserRegistered{}, []event.Listener{
        &Listeners.SendWelcomeEmail{},
        &Listeners.CreateUserProfile{},
        &Listeners.SendAdminNotification{},
    })

    // 注册通配符监听器
    p.Listen("*", []event.Listener{
        &Listeners.LogAllEvents{},
    })

    // 注册订阅者
    p.Subscribe(&Listeners.UserEventSubscriber{})
}
```

### 4. 事件分发

```go
// 在控制器中分发事件
func (c *AuthController) Register(request http.Request) http.Response {
    // 创建用户
    user, err := c.userService.CreateUser(request.Body)
    if err != nil {
        return c.JsonError("Registration failed", 500)
    }

    // 分发用户注册事件
    event := events.NewUserRegistered(user)
    c.eventDispatcher.Dispatch(event)

    return c.Json(user).Status(201)
}
```

## 📚 API 参考

### Event 接口

```go
type Event interface {
    GetName() string
    GetData() interface{}
    GetTimestamp() time.Time
    IsPropagationStopped() bool
    StopPropagation()
}
```

#### 方法说明

- `GetName()`: 获取事件名称
- `GetData()`: 获取事件数据
- `GetTimestamp()`: 获取事件时间戳
- `IsPropagationStopped()`: 检查事件传播是否已停止
- `StopPropagation()`: 停止事件传播

### BaseEvent 结构体

```go
type BaseEvent struct {
    Name      string
    Data      interface{}
    Timestamp time.Time
    Stopped   bool
}
```

#### 字段说明

- `Name`: 事件名称
- `Data`: 事件数据
- `Timestamp`: 事件时间戳
- `Stopped`: 是否停止传播

### Listener 接口

```go
type Listener interface {
    Handle(event Event) error
    ShouldQueue() bool
    GetQueue() string
    GetMaxAttempts() int
    GetTimeout() time.Duration
}
```

#### 方法说明

- `Handle(event)`: 处理事件
- `ShouldQueue()`: 是否应该异步执行
- `GetQueue()`: 获取队列名称
- `GetMaxAttempts()`: 获取最大重试次数
- `GetTimeout()`: 获取超时时间

### BaseListener 结构体

```go
type BaseListener struct {
    Event       string
    Queue       string
    MaxAttempts int
    Timeout     time.Duration
}
```

#### 字段说明

- `Event`: 监听的事件名称
- `Queue`: 队列名称
- `MaxAttempts`: 最大重试次数
- `Timeout`: 超时时间

### Dispatcher 接口

```go
type Dispatcher interface {
    Listen(event Event, listeners []Listener)
    ListenPattern(pattern string, listeners []Listener)
    Subscribe(subscriber Subscriber)
    Dispatch(event Event) error
    DispatchAsync(event Event) error
    Forget(event Event)
    ForgetPattern(pattern string)
    HasListeners(event Event) bool
    GetListeners(event Event) []Listener
}
```

#### 方法说明

- `Listen(event, listeners)`: 注册事件监听器
- `ListenPattern(pattern, listeners)`: 注册通配符监听器
- `Subscribe(subscriber)`: 注册事件订阅者
- `Dispatch(event)`: 同步分发事件
- `DispatchAsync(event)`: 异步分发事件
- `Forget(event)`: 移除事件监听器
- `ForgetPattern(pattern)`: 移除通配符监听器
- `HasListeners(event)`: 检查是否有监听器
- `GetListeners(event)`: 获取事件监听器列表

## 🎯 高级功能

### 1. 事件订阅者

```go
// app/Listeners/UserEventSubscriber.go
package listeners

import (
    "laravel-go/framework/event"
    "laravel-go/app/Events"
)

type UserEventSubscriber struct {
    event.BaseSubscriber
}

func (s *UserEventSubscriber) Subscribe(dispatcher event.Dispatcher) {
    // 订阅多个事件
    dispatcher.Listen(&Events.UserRegistered{}, []event.Listener{
        &SendWelcomeEmail{},
        &CreateUserProfile{},
    })

    dispatcher.Listen(&Events.UserUpdated{}, []event.Listener{
        &UpdateUserCache{},
        &NotifyFollowers{},
    })

    dispatcher.Listen(&Events.UserDeleted{}, []event.Listener{
        &CleanupUserData{},
        &NotifyAdmins{},
    })
}
```

### 2. 通配符监听器

```go
// 监听所有用户相关事件
dispatcher.ListenPattern("user.*", []event.Listener{
    &LogUserActivity{},
    &UpdateUserMetrics{},
})

// 监听所有事件
dispatcher.ListenPattern("*", []event.Listener{
    &GlobalEventLogger{},
})
```

### 3. 事件队列

```go
// 异步事件监听器
type SendWelcomeEmail struct {
    event.BaseListener
}

func (l *SendWelcomeEmail) ShouldQueue() bool {
    return true
}

func (l *SendWelcomeEmail) GetQueue() string {
    return "emails"
}

func (l *SendWelcomeEmail) GetMaxAttempts() int {
    return 3
}

func (l *SendWelcomeEmail) GetTimeout() time.Duration {
    return time.Minute * 5
}
```

### 4. 事件传播控制

```go
func (l *SendWelcomeEmail) Handle(event event.Event) error {
    // 检查事件是否已被其他监听器停止
    if event.IsPropagationStopped() {
        return nil
    }

    // 处理事件
    userRegistered := event.(*Events.UserRegistered)

    // 在某些条件下停止事件传播
    if userRegistered.User.Email == "" {
        event.StopPropagation()
        return errors.New("user email is empty")
    }

    return sendWelcomeEmail(userRegistered.User.Email, userRegistered.User.Name)
}
```

## 🔧 配置选项

### 事件系统配置

```go
// config/event.go
package config

type EventConfig struct {
    // 默认队列名称
    DefaultQueue string `json:"default_queue"`

    // 默认最大重试次数
    DefaultMaxAttempts int `json:"default_max_attempts"`

    // 默认超时时间
    DefaultTimeout time.Duration `json:"default_timeout"`

    // 是否启用事件日志
    EnableLogging bool `json:"enable_logging"`

    // 事件日志级别
    LogLevel string `json:"log_level"`

    // 队列配置
    Queue QueueConfig `json:"queue"`
}

type QueueConfig struct {
    // 队列驱动
    Driver string `json:"driver"`

    // Redis 配置
    Redis RedisConfig `json:"redis"`

    // 数据库配置
    Database DatabaseConfig `json:"database"`
}
```

### 配置示例

```go
// config/event.go
func GetEventConfig() *EventConfig {
    return &EventConfig{
        DefaultQueue:       "default",
        DefaultMaxAttempts: 3,
        DefaultTimeout:     time.Minute * 5,
        EnableLogging:      true,
        LogLevel:           "info",
        Queue: QueueConfig{
            Driver: "redis",
            Redis: RedisConfig{
                Host:     "localhost",
                Port:     6379,
                Database: 0,
            },
        },
    }
}
```

## 🚀 性能优化

### 1. 事件缓存

```go
// 缓存事件监听器
type CachedDispatcher struct {
    event.Dispatcher
    cache cache.Cache
}

func (d *CachedDispatcher) GetListeners(event Event) []Listener {
    cacheKey := fmt.Sprintf("event:listeners:%s", event.GetName())

    if cached, exists := d.cache.Get(cacheKey); exists {
        return cached.([]Listener)
    }

    listeners := d.Dispatcher.GetListeners(event)
    d.cache.Put(cacheKey, listeners, time.Hour)

    return listeners
}
```

### 2. 批量事件处理

```go
// 批量处理事件
func (d *Dispatcher) DispatchBatch(events []Event) error {
    for _, event := range events {
        if err := d.Dispatch(event); err != nil {
            return err
        }
    }
    return nil
}
```

### 3. 事件优先级

```go
type PriorityListener struct {
    event.BaseListener
    Priority int
}

func (l *PriorityListener) GetPriority() int {
    return l.Priority
}

// 按优先级排序监听器
func sortListenersByPriority(listeners []Listener) []Listener {
    sort.Slice(listeners, func(i, j int) bool {
        if p1, ok := listeners[i].(*PriorityListener); ok {
            if p2, ok := listeners[j].(*PriorityListener); ok {
                return p1.GetPriority() > p2.GetPriority()
            }
        }
        return false
    })
    return listeners
}
```

## 🧪 测试

### 1. 事件测试

```go
// tests/event_test.go
package tests

import (
    "testing"
    "laravel-go/framework/event"
    "laravel-go/app/Events"
    "laravel-go/app/Listeners"
)

func TestUserRegisteredEvent(t *testing.T) {
    // 创建事件分发器
    dispatcher := event.NewDispatcher()

    // 注册监听器
    listener := &Listeners.SendWelcomeEmail{}
    dispatcher.Listen(&Events.UserRegistered{}, []event.Listener{listener})

    // 创建事件
    user := &Models.User{Name: "John", Email: "john@example.com"}
    event := events.NewUserRegistered(user)

    // 分发事件
    err := dispatcher.Dispatch(event)
    if err != nil {
        t.Errorf("Failed to dispatch event: %v", err)
    }

    // 验证监听器被调用
    // 这里需要添加监听器调用验证逻辑
}
```

### 2. 监听器测试

```go
func TestSendWelcomeEmailListener(t *testing.T) {
    listener := &Listeners.SendWelcomeEmail{}

    // 创建测试事件
    user := &Models.User{Name: "John", Email: "john@example.com"}
    event := events.NewUserRegistered(user)

    // 测试监听器处理
    err := listener.Handle(event)
    if err != nil {
        t.Errorf("Listener failed to handle event: %v", err)
    }

    // 验证邮件发送逻辑
    // 这里需要添加邮件发送验证逻辑
}
```

## 🔍 调试和监控

### 1. 事件日志

```go
type EventLogger struct {
    event.BaseListener
    logger log.Logger
}

func (l *EventLogger) Handle(event Event) error {
    l.logger.Info("Event dispatched", map[string]interface{}{
        "event":     event.GetName(),
        "timestamp": event.GetTimestamp(),
        "data":      event.GetData(),
    })
    return nil
}
```

### 2. 事件监控

```go
type EventMonitor struct {
    event.BaseListener
    metrics metrics.Collector
}

func (m *EventMonitor) Handle(event Event) error {
    // 记录事件指标
    m.metrics.Increment("events.dispatched", map[string]string{
        "event": event.GetName(),
    })

    // 记录事件处理时间
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        m.metrics.Histogram("events.duration", duration.Seconds(), map[string]string{
            "event": event.GetName(),
        })
    }()

    return nil
}
```

## 📝 最佳实践

### 1. 事件命名规范

```go
// 使用点分隔的命名方式
type UserRegistered struct {
    event.BaseEvent
    User *Models.User
}

func (e *UserRegistered) GetName() string {
    return "user.registered"
}

type UserUpdated struct {
    event.BaseEvent
    User *Models.User
}

func (e *UserUpdated) GetName() string {
    return "user.updated"
}
```

### 2. 事件数据设计

```go
// 事件数据应该包含足够的信息
type OrderCreated struct {
    event.BaseEvent
    Order     *Models.Order
    User      *Models.User
    Timestamp time.Time
}

func (e *OrderCreated) GetData() interface{} {
    return map[string]interface{}{
        "order_id":     e.Order.ID,
        "user_id":      e.User.ID,
        "total_amount": e.Order.TotalAmount,
        "created_at":   e.Timestamp,
    }
}
```

### 3. 监听器职责分离

```go
// 每个监听器只负责一个职责
type SendWelcomeEmail struct {
    event.BaseListener
}

func (l *SendWelcomeEmail) Handle(event Event) error {
    // 只负责发送欢迎邮件
    return sendWelcomeEmail(event.GetData())
}

type CreateUserProfile struct {
    event.BaseListener
}

func (l *CreateUserProfile) Handle(event Event) error {
    // 只负责创建用户档案
    return createUserProfile(event.GetData())
}
```

### 4. 错误处理

```go
func (l *SendWelcomeEmail) Handle(event Event) error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic in SendWelcomeEmail: %v", r)
        }
    }()

    userRegistered := event.(*Events.UserRegistered)

    if err := sendWelcomeEmail(userRegistered.User.Email, userRegistered.User.Name); err != nil {
        // 记录错误但不阻止其他监听器执行
        log.Printf("Failed to send welcome email: %v", err)
        return err
    }

    return nil
}
```

## 🚀 总结

事件系统是 Laravel-Go Framework 中强大的功能之一，它提供了：

1. **松耦合架构**: 组件间通过事件通信，降低耦合度
2. **异步处理**: 支持事件队列，提高系统性能
3. **灵活扩展**: 易于添加新的事件和监听器
4. **监控调试**: 提供完整的监控和调试功能
5. **最佳实践**: 遵循事件驱动架构的最佳实践

通过合理使用事件系统，可以构建出更加灵活、可维护和可扩展的应用程序。
