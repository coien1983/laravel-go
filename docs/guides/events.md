# 事件系统指南

## 📖 概述

Laravel-Go Framework 提供了完整的事件系统，支持事件驱动编程、事件监听器、事件订阅者、事件广播等功能，帮助构建松耦合、可扩展的应用程序架构。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [事件系统 API 参考](../api/event.md)

## 🚀 快速开始

### 1. 基本使用

```go
// 定义事件
type UserRegistered struct {
    User *User `json:"user"`
}

// 定义监听器
type SendWelcomeEmailListener struct {
    emailService *EmailService
}

func NewSendWelcomeEmailListener(emailService *EmailService) *SendWelcomeEmailListener {
    return &SendWelcomeEmailListener{
        emailService: emailService,
    }
}

func (l *SendWelcomeEmailListener) Handle(event interface{}) error {
    userRegistered := event.(*UserRegistered)
    return l.emailService.SendWelcomeEmail(userRegistered.User)
}

// 注册事件监听器
func RegisterEventListeners() {
    dispatcher := event.GetDispatcher()

    dispatcher.Listen(&UserRegistered{}, []event.Listener{
        NewSendWelcomeEmailListener(emailService),
        NewCreateUserProfileListener(profileService),
    })
}

// 触发事件
func (s *UserService) RegisterUser(data map[string]interface{}) (*User, error) {
    user, err := s.createUser(data)
    if err != nil {
        return nil, err
    }

    // 触发用户注册事件
    event.Dispatch(&UserRegistered{User: user})

    return user, nil
}
```

### 2. 事件监听器注册

```go
// 在应用启动时注册事件
func main() {
    // 注册事件监听器
    RegisterEventListeners()

    // 启动服务器
    server := http.NewServer()
    server.Start(":8080")
}
```

## 📡 事件类型

### 1. 同步事件

```go
// 同步事件（默认）
type OrderCreated struct {
    Order *Order `json:"order"`
}

type UpdateInventoryListener struct {
    inventoryService *InventoryService
}

func (l *UpdateInventoryListener) Handle(event interface{}) error {
    orderCreated := event.(*OrderCreated)

    // 同步更新库存
    return l.inventoryService.UpdateStock(orderCreated.Order)
}

// 注册同步监听器
dispatcher.Listen(&OrderCreated{}, []event.Listener{
    &UpdateInventoryListener{inventoryService},
})
```

### 2. 异步事件

```go
// 异步事件
type OrderShipped struct {
    Order *Order `json:"order"`
    TrackingNumber string `json:"tracking_number"`
}

type SendShippingNotificationListener struct {
    emailService *EmailService
}

func (l *SendShippingNotificationListener) Handle(event interface{}) error {
    orderShipped := event.(*OrderShipped)

    // 异步发送邮件通知
    return l.emailService.SendShippingNotification(
        orderShipped.Order,
        orderShipped.TrackingNumber,
    )
}

// 注册异步监听器
dispatcher.Listen(&OrderShipped{}, []event.Listener{
    &SendShippingNotificationListener{emailService},
}).Async() // 标记为异步
```

### 3. 队列事件

```go
// 队列事件
type ProcessPaymentEvent struct {
    Payment *Payment `json:"payment"`
}

type ProcessPaymentListener struct {
    paymentService *PaymentService
}

func (l *ProcessPaymentListener) Handle(event interface{}) error {
    paymentEvent := event.(*ProcessPaymentEvent)

    // 处理支付（可能耗时较长）
    return l.paymentService.ProcessPayment(paymentEvent.Payment)
}

// 注册队列监听器
dispatcher.Listen(&ProcessPaymentEvent{}, []event.Listener{
    &ProcessPaymentListener{paymentService},
}).Queue("payments") // 推送到指定队列
```

## 🎯 事件监听器

### 1. 单个监听器

```go
// 单个监听器
type LogUserActivityListener struct {
    logger *Logger
}

func (l *LogUserActivityListener) Handle(event interface{}) error {
    userRegistered := event.(*UserRegistered)

    l.logger.Info("User registered", map[string]interface{}{
        "user_id": userRegistered.User.ID,
        "email":   userRegistered.User.Email,
        "time":    time.Now(),
    })

    return nil
}
```

### 2. 多个监听器

```go
// 多个监听器
type UserRegistered struct {
    User *User `json:"user"`
}

// 监听器1：发送欢迎邮件
type SendWelcomeEmailListener struct {
    emailService *EmailService
}

func (l *SendWelcomeEmailListener) Handle(event interface{}) error {
    userRegistered := event.(*UserRegistered)
    return l.emailService.SendWelcomeEmail(userRegistered.User)
}

// 监听器2：创建用户档案
type CreateUserProfileListener struct {
    profileService *ProfileService
}

func (l *CreateUserProfileListener) Handle(event interface{}) error {
    userRegistered := event.(*UserRegistered)
    return l.profileService.CreateProfile(userRegistered.User.ID)
}

// 监听器3：记录用户活动
type LogUserActivityListener struct {
    activityService *ActivityService
}

func (l *LogUserActivityListener) Handle(event interface{}) error {
    userRegistered := event.(*UserRegistered)
    return l.activityService.LogActivity(userRegistered.User.ID, "registered")
}

// 注册多个监听器
dispatcher.Listen(&UserRegistered{}, []event.Listener{
    &SendWelcomeEmailListener{emailService},
    &CreateUserProfileListener{profileService},
    &LogUserActivityListener{activityService},
})
```

### 3. 条件监听器

```go
// 条件监听器
type ConditionalListener struct {
    emailService *EmailService
    condition    func(event interface{}) bool
}

func (l *ConditionalListener) Handle(event interface{}) error {
    // 检查条件
    if !l.condition(event) {
        return nil // 跳过处理
    }

    userRegistered := event.(*UserRegistered)
    return l.emailService.SendWelcomeEmail(userRegistered.User)
}

// 使用条件监听器
dispatcher.Listen(&UserRegistered{}, []event.Listener{
    &ConditionalListener{
        emailService: emailService,
        condition: func(event interface{}) bool {
            userRegistered := event.(*UserRegistered)
            return userRegistered.User.Email != "" // 只有邮箱不为空才发送
        },
    },
})
```

## 🔄 事件订阅者

### 1. 事件订阅者

```go
// 事件订阅者
type UserEventSubscriber struct {
    emailService    *EmailService
    profileService  *ProfileService
    activityService *ActivityService
}

func NewUserEventSubscriber(
    emailService *EmailService,
    profileService *ProfileService,
    activityService *ActivityService,
) *UserEventSubscriber {
    return &UserEventSubscriber{
        emailService:    emailService,
        profileService:  profileService,
        activityService: activityService,
    }
}

// 订阅用户注册事件
func (s *UserEventSubscriber) OnUserRegistered(event *UserRegistered) error {
    // 发送欢迎邮件
    if err := s.emailService.SendWelcomeEmail(event.User); err != nil {
        return err
    }

    // 创建用户档案
    if err := s.profileService.CreateProfile(event.User.ID); err != nil {
        return err
    }

    // 记录活动
    return s.activityService.LogActivity(event.User.ID, "registered")
}

// 订阅用户登录事件
func (s *UserEventSubscriber) OnUserLoggedIn(event *UserLoggedIn) error {
    return s.activityService.LogActivity(event.User.ID, "logged_in")
}

// 注册订阅者
func RegisterEventSubscribers() {
    dispatcher := event.GetDispatcher()

    subscriber := NewUserEventSubscriber(emailService, profileService, activityService)
    dispatcher.Subscribe(subscriber)
}
```

### 2. 自动事件映射

```go
// 自动事件映射
type OrderEventSubscriber struct {
    inventoryService *InventoryService
    notificationService *NotificationService
}

// 方法名格式：On + 事件名
func (s *OrderEventSubscriber) OnOrderCreated(event *OrderCreated) error {
    return s.inventoryService.UpdateStock(event.Order)
}

func (s *OrderEventSubscriber) OnOrderShipped(event *OrderShipped) error {
    return s.notificationService.SendShippingNotification(event.Order)
}

func (s *OrderEventSubscriber) OnOrderCancelled(event *OrderCancelled) error {
    return s.inventoryService.RestoreStock(event.Order)
}
```

## 📡 事件广播

### 1. 本地广播

```go
// 本地广播事件
type UserStatusChanged struct {
    UserID uint   `json:"user_id"`
    Status string `json:"status"`
    Time   time.Time `json:"time"`
}

type BroadcastUserStatusListener struct {
    broadcaster *Broadcaster
}

func (l *BroadcastUserStatusListener) Handle(event interface{}) error {
    statusChanged := event.(*UserStatusChanged)

    // 广播到本地频道
    return l.broadcaster.Broadcast("user.status", statusChanged)
}

// 注册广播监听器
dispatcher.Listen(&UserStatusChanged{}, []event.Listener{
    &BroadcastUserStatusListener{broadcaster},
})
```

### 2. WebSocket 广播

```go
// WebSocket 广播
type MessageSent struct {
    RoomID  string `json:"room_id"`
    UserID  uint   `json:"user_id"`
    Message string `json:"message"`
}

type BroadcastMessageListener struct {
    websocketService *WebSocketService
}

func (l *BroadcastMessageListener) Handle(event interface{}) error {
    messageSent := event.(*MessageSent)

    // 广播到 WebSocket 频道
    return l.websocketService.BroadcastToRoom(
        messageSent.RoomID,
        "message.sent",
        messageSent,
    )
}
```

### 3. Redis 广播

```go
// Redis 广播
type CacheUpdated struct {
    Key   string      `json:"key"`
    Value interface{} `json:"value"`
}

type RedisBroadcastListener struct {
    redisClient *RedisClient
}

func (l *RedisBroadcastListener) Handle(event interface{}) error {
    cacheUpdated := event.(*CacheUpdated)

    // 通过 Redis 发布事件
    return l.redisClient.Publish("cache.updated", cacheUpdated)
}
```

## 🛡️ 错误处理

### 1. 监听器错误处理

```go
// 错误处理监听器
type SafeListener struct {
    listener event.Listener
    logger   *Logger
}

func (l *SafeListener) Handle(event interface{}) error {
    defer func() {
        if r := recover(); r != nil {
            l.logger.Error("Listener panicked", map[string]interface{}{
                "error": r,
                "event": event,
            })
        }
    }()

    return l.listener.Handle(event)
}

// 包装监听器
dispatcher.Listen(&UserRegistered{}, []event.Listener{
    &SafeListener{
        listener: &SendWelcomeEmailListener{emailService},
        logger:   logger,
    },
})
```

### 2. 事件失败处理

```go
// 事件失败处理
type EventFailureHandler struct {
    logger *Logger
    queue  *Queue
}

func (h *EventFailureHandler) HandleFailure(event interface{}, err error) {
    h.logger.Error("Event handling failed", map[string]interface{}{
        "event": event,
        "error": err.Error(),
    })

    // 将失败的事件推送到队列重试
    h.queue.Push(&RetryEventJob{
        Event: event,
        Error: err.Error(),
    })
}

// 注册失败处理器
dispatcher.SetFailureHandler(&EventFailureHandler{
    logger: logger,
    queue:  queue,
})
```

## 📊 事件监控

### 1. 事件统计

```go
// 事件统计
type EventStats struct {
    EventName    string `json:"event_name"`
    TotalFired   int64  `json:"total_fired"`
    TotalHandled int64  `json:"total_handled"`
    FailedCount  int64  `json:"failed_count"`
    AvgDuration  time.Duration `json:"avg_duration"`
}

// 事件监控监听器
type EventMonitorListener struct {
    stats map[string]*EventStats
    mutex sync.RWMutex
}

func (l *EventMonitorListener) Before(event interface{}) {
    eventName := reflect.TypeOf(event).String()

    l.mutex.Lock()
    defer l.mutex.Unlock()

    if l.stats[eventName] == nil {
        l.stats[eventName] = &EventStats{EventName: eventName}
    }

    l.stats[eventName].TotalFired++
}

func (l *EventMonitorListener) After(event interface{}, err error) {
    eventName := reflect.TypeOf(event).String()

    l.mutex.Lock()
    defer l.mutex.Unlock()

    if l.stats[eventName] != nil {
        l.stats[eventName].TotalHandled++
        if err != nil {
            l.stats[eventName].FailedCount++
        }
    }
}

// 获取事件统计
func (l *EventMonitorListener) GetStats() map[string]*EventStats {
    l.mutex.RLock()
    defer l.mutex.RUnlock()

    stats := make(map[string]*EventStats)
    for k, v := range l.stats {
        stats[k] = v
    }

    return stats
}
```

### 2. 事件日志

```go
// 事件日志监听器
type EventLogListener struct {
    logger *Logger
}

func (l *EventLogListener) Before(event interface{}) {
    l.logger.Info("Event fired", map[string]interface{}{
        "event": reflect.TypeOf(event).String(),
        "data":  event,
        "time":  time.Now(),
    })
}

func (l *EventLogListener) After(event interface{}, err error) {
    if err != nil {
        l.logger.Error("Event handling failed", map[string]interface{}{
            "event": reflect.TypeOf(event).String(),
            "error": err.Error(),
        })
    } else {
        l.logger.Info("Event handled successfully", map[string]interface{}{
            "event": reflect.TypeOf(event).String(),
        })
    }
}
```

## 🔧 高级功能

### 1. 事件中间件

```go
// 事件中间件
type EventMiddleware interface {
    Before(event interface{})
    After(event interface{}, err error)
}

// 性能监控中间件
type PerformanceMiddleware struct {
    logger *Logger
}

func (m *PerformanceMiddleware) Before(event interface{}) {
    // 记录开始时间
    event.SetMetadata("started_at", time.Now())
}

func (m *PerformanceMiddleware) After(event interface{}, err error) {
    startTime := event.GetMetadata("started_at").(time.Time)
    duration := time.Since(startTime)

    if duration > time.Millisecond*100 {
        m.logger.Warning("Slow event handling", map[string]interface{}{
            "event":    reflect.TypeOf(event).String(),
            "duration": duration,
        })
    }
}

// 注册中间件
dispatcher.Use(&PerformanceMiddleware{logger})
```

### 2. 事件优先级

```go
// 优先级监听器
type PriorityListener struct {
    listener event.Listener
    priority int
}

func (l *PriorityListener) Handle(event interface{}) error {
    return l.listener.Handle(event)
}

func (l *PriorityListener) GetPriority() int {
    return l.priority
}

// 注册优先级监听器
dispatcher.Listen(&UserRegistered{}, []event.Listener{
    &PriorityListener{
        listener: &SendWelcomeEmailListener{emailService},
        priority: 1, // 高优先级
    },
    &PriorityListener{
        listener: &LogUserActivityListener{activityService},
        priority: 10, // 低优先级
    },
})
```

### 3. 事件过滤

```go
// 事件过滤器
type EventFilter interface {
    ShouldHandle(event interface{}) bool
}

// 用户事件过滤器
type UserEventFilter struct {
    userID uint
}

func (f *UserEventFilter) ShouldHandle(event interface{}) bool {
    switch e := event.(type) {
    case *UserRegistered:
        return e.User.ID == f.userID
    case *UserLoggedIn:
        return e.User.ID == f.userID
    default:
        return false
    }
}

// 使用过滤器
dispatcher.Listen(&UserRegistered{}, []event.Listener{
    &SendWelcomeEmailListener{emailService},
}).Filter(&UserEventFilter{userID: 123})
```

## 📚 总结

Laravel-Go Framework 的事件系统提供了：

1. **事件类型**: 同步、异步、队列事件
2. **监听器**: 单个、多个、条件监听器
3. **订阅者**: 自动事件映射
4. **事件广播**: 本地、WebSocket、Redis 广播
5. **错误处理**: 监听器错误处理、事件失败处理
6. **监控功能**: 事件统计、事件日志
7. **高级功能**: 事件中间件、优先级、过滤器

通过合理使用事件系统，可以构建松耦合、可扩展的应用程序架构。
