package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"laravel-go/framework/event"
)

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	*event.BaseEvent
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// NewUserRegisteredEvent 创建用户注册事件
func NewUserRegisteredEvent(userID int64, username, email string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: event.NewEvent("user.registered", nil),
		UserID:    userID,
		Username:  username,
		Email:     email,
	}
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	*event.BaseEvent
	OrderID  int64    `json:"order_id"`
	UserID   int64    `json:"user_id"`
	Amount   float64  `json:"amount"`
	Products []string `json:"products"`
}

// NewOrderCreatedEvent 创建订单创建事件
func NewOrderCreatedEvent(orderID, userID int64, amount float64, products []string) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		BaseEvent: event.NewEvent("order.created", nil),
		OrderID:   orderID,
		UserID:    userID,
		Amount:    amount,
		Products:  products,
	}
}

// EmailNotificationListener 邮件通知监听器
type EmailNotificationListener struct {
	*event.BaseListener
}

// NewEmailNotificationListener 创建邮件通知监听器
func NewEmailNotificationListener() *EmailNotificationListener {
	return &EmailNotificationListener{
		BaseListener: event.NewListener("email.notification", func(e event.Event) error {
			fmt.Printf("📧 发送邮件通知: %s\n", e.GetName())

			// 根据事件类型发送不同的邮件
			switch e.GetName() {
			case "user.registered":
				if userEvent, ok := e.(*UserRegisteredEvent); ok {
					fmt.Printf("   欢迎邮件发送给: %s (%s)\n", userEvent.Username, userEvent.Email)
				}
			case "order.created":
				if orderEvent, ok := e.(*OrderCreatedEvent); ok {
					fmt.Printf("   订单确认邮件发送给用户: %d, 订单金额: %.2f\n", orderEvent.UserID, orderEvent.Amount)
				}
			}

			return nil
		}),
	}
}

// DatabaseLogListener 数据库日志监听器
type DatabaseLogListener struct {
	*event.BaseListener
}

// NewDatabaseLogListener 创建数据库日志监听器
func NewDatabaseLogListener() *DatabaseLogListener {
	return &DatabaseLogListener{
		BaseListener: event.NewListener("database.log", func(e event.Event) error {
			fmt.Printf("💾 记录数据库日志: %s\n", e.GetName())
			fmt.Printf("   事件ID: %s, 时间: %s\n", e.GetID(), e.GetTimestamp().Format("2006-01-02 15:04:05"))
			return nil
		}),
	}
}

// AnalyticsListener 数据分析监听器
type AnalyticsListener struct {
	*event.BaseListener
}

// NewAnalyticsListener 创建数据分析监听器
func NewAnalyticsListener() *AnalyticsListener {
	return &AnalyticsListener{
		BaseListener: event.NewListener("analytics.track", func(e event.Event) error {
			fmt.Printf("📊 数据分析追踪: %s\n", e.GetName())

			// 模拟数据分析处理
			time.Sleep(50 * time.Millisecond)

			switch e.GetName() {
			case "user.registered":
				fmt.Printf("   新用户注册统计更新\n")
			case "order.created":
				fmt.Printf("   订单数据统计更新\n")
			}

			return nil
		}),
	}
}

// QueuedNotificationListener 队列化通知监听器
type QueuedNotificationListener struct {
	*event.BaseListener
}

// NewQueuedNotificationListener 创建队列化通知监听器
func NewQueuedNotificationListener() *QueuedNotificationListener {
	return &QueuedNotificationListener{
		BaseListener: event.NewQueuedListener("queued.notification", "notifications", func(e event.Event) error {
			fmt.Printf("⏳ 队列化通知处理: %s\n", e.GetName())

			// 模拟耗时操作
			time.Sleep(200 * time.Millisecond)

			fmt.Printf("   ✅ 队列化通知处理完成\n")
			return nil
		}),
	}
}

// EventSubscriber 事件订阅者
type EventSubscriber struct {
	name string
}

// NewEventSubscriber 创建事件订阅者
func NewEventSubscriber(name string) *EventSubscriber {
	return &EventSubscriber{name: name}
}

// Subscribe 订阅事件
func (s *EventSubscriber) Subscribe(dispatcher event.Dispatcher) {
	// 订阅用户相关事件
	dispatcher.Listen("user.registered", NewEmailNotificationListener())
	dispatcher.Listen("user.registered", NewDatabaseLogListener())

	// 订阅订单相关事件
	dispatcher.Listen("order.created", NewEmailNotificationListener())
	dispatcher.Listen("order.created", NewDatabaseLogListener())
	dispatcher.Listen("order.created", NewAnalyticsListener())

	fmt.Printf("📋 事件订阅者 '%s' 已注册\n", s.name)
}

// GetName 获取订阅者名称
func (s *EventSubscriber) GetName() string {
	return s.name
}

func main() {
	fmt.Println("🚀 Laravel-Go 事件系统演示")
	fmt.Println("==================================================")

	// 初始化事件系统
	event.Init()

	// 创建事件订阅者
	subscriber := NewEventSubscriber("main.subscriber")
	event.Subscribe(subscriber)

	// 创建队列化监听器
	queuedListener := NewQueuedNotificationListener()
	event.Listen("user.registered", queuedListener)
	event.Listen("order.created", queuedListener)

	fmt.Println("\n📝 创建用户注册事件...")
	userEvent := NewUserRegisteredEvent(1, "john_doe", "john@example.com")

	// 分发用户注册事件
	err := event.Dispatch(userEvent)
	if err != nil {
		log.Fatalf("Failed to dispatch user event: %v", err)
	}

	fmt.Println("\n📝 创建订单创建事件...")
	orderEvent := NewOrderCreatedEvent(1001, 1, 299.99, []string{"iPhone 15", "AirPods Pro"})

	// 分发订单创建事件
	err = event.Dispatch(orderEvent)
	if err != nil {
		log.Fatalf("Failed to dispatch order event: %v", err)
	}

	// 等待队列化事件处理
	fmt.Println("\n⏳ 等待队列化事件处理...")
	time.Sleep(1 * time.Second)

	// 测试异步事件分发
	fmt.Println("\n🔄 测试异步事件分发...")
	asyncEvent := NewUserRegisteredEvent(2, "jane_doe", "jane@example.com")
	err = event.DispatchAsync(asyncEvent)
	if err != nil {
		log.Fatalf("Failed to dispatch async event: %v", err)
	}

	// 等待异步处理
	time.Sleep(200 * time.Millisecond)

	// 测试批量事件分发
	fmt.Println("\n📦 测试批量事件分发...")
	batchEvents := []event.Event{
		NewUserRegisteredEvent(3, "bob_smith", "bob@example.com"),
		NewOrderCreatedEvent(1002, 3, 199.99, []string{"MacBook Air"}),
		NewUserRegisteredEvent(4, "alice_jones", "alice@example.com"),
	}

	err = event.DispatchBatch(batchEvents)
	if err != nil {
		log.Fatalf("Failed to dispatch batch events: %v", err)
	}

	// 等待所有事件处理完成
	time.Sleep(500 * time.Millisecond)

	// 显示统计信息
	fmt.Println("\n📊 事件系统统计信息:")
	stats := event.GetStats()
	fmt.Printf("   总事件数: %d\n", stats.TotalEvents)
	fmt.Printf("   已分发事件: %d\n", stats.DispatchedEvents)
	fmt.Printf("   队列化事件: %d\n", stats.QueuedEvents)
	fmt.Printf("   失败事件: %d\n", stats.FailedEvents)
	fmt.Printf("   最后事件时间: %s\n", stats.LastEventAt.Format("2006-01-02 15:04:05"))

	// 显示监听器信息
	fmt.Println("\n👂 监听器信息:")
	allListeners := event.GetAllListeners()
	for eventName, listeners := range allListeners {
		fmt.Printf("   事件 '%s': %d 个监听器\n", eventName, len(listeners))
		for _, listener := range listeners {
			fmt.Printf("     - %s (优先级: %d, 队列: %s)\n",
				listener.GetName(),
				listener.GetPriority(),
				listener.GetQueue())
		}
	}

	// 测试事件队列功能
	fmt.Println("\n🔄 测试事件队列功能...")
	queue := event.NewMemoryEventQueue()
	defer queue.Close()

	// 推送事件到队列
	queueEvent := NewUserRegisteredEvent(5, "queue_user", "queue@example.com")
	err = queue.Push(queueEvent)
	if err != nil {
		log.Fatalf("Failed to push event to queue: %v", err)
	}

	// 从队列弹出事件
	ctx := context.Background()
	poppedEvent, err := queue.Pop(ctx)
	if err != nil {
		log.Fatalf("Failed to pop event from queue: %v", err)
	}

	fmt.Printf("   从队列弹出事件: %s\n", poppedEvent.GetName())

	// 清理资源
	event.Close()

	fmt.Println("\n✅ 事件系统演示完成!")
}
