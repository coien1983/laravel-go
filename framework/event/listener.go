package event

import (
	"fmt"
	"reflect"
)

// BaseListener 基础监听器实现
type BaseListener struct {
	name     string
	priority int
	queue    string
	handler  func(Event) error
}

// NewListener 创建新监听器
func NewListener(name string, handler func(Event) error) *BaseListener {
	return &BaseListener{
		name:     name,
		priority: 0,
		queue:    "",
		handler:  handler,
	}
}

// NewListenerWithPriority 创建带优先级的监听器
func NewListenerWithPriority(name string, priority int, handler func(Event) error) *BaseListener {
	return &BaseListener{
		name:     name,
		priority: priority,
		queue:    "",
		handler:  handler,
	}
}

// NewQueuedListener 创建队列监听器
func NewQueuedListener(name string, queue string, handler func(Event) error) *BaseListener {
	return &BaseListener{
		name:     name,
		priority: 0,
		queue:    queue,
		handler:  handler,
	}
}

// Handle 处理事件
func (l *BaseListener) Handle(event Event) error {
	if l.handler == nil {
		return fmt.Errorf("listener handler is nil")
	}
	return l.handler(event)
}

// GetName 获取监听器名称
func (l *BaseListener) GetName() string {
	return l.name
}

// GetPriority 获取优先级
func (l *BaseListener) GetPriority() int {
	return l.priority
}

// ShouldQueue 是否应该队列化
func (l *BaseListener) ShouldQueue() bool {
	return l.queue != ""
}

// GetQueue 获取队列名称
func (l *BaseListener) GetQueue() string {
	return l.queue
}

// SetPriority 设置优先级
func (l *BaseListener) SetPriority(priority int) {
	l.priority = priority
}

// SetQueue 设置队列
func (l *BaseListener) SetQueue(queue string) {
	l.queue = queue
}

// String 字符串表示
func (l *BaseListener) String() string {
	return fmt.Sprintf("Listener{Name: %s, Priority: %d, Queue: %s}", l.name, l.priority, l.queue)
}

// FunctionListener 函数监听器
type FunctionListener struct {
	*BaseListener
}

// NewFunctionListener 创建函数监听器
func NewFunctionListener(name string, fn interface{}) *FunctionListener {
	handler := func(event Event) error {
		fnValue := reflect.ValueOf(fn)
		fnType := fnValue.Type()

		// 检查函数签名
		if fnType.Kind() != reflect.Func {
			return fmt.Errorf("handler must be a function")
		}

		// 根据函数参数数量调用
		switch fnType.NumIn() {
		case 0:
			// 无参数函数
			results := fnValue.Call(nil)
			if len(results) > 0 && !results[0].IsNil() {
				return results[0].Interface().(error)
			}
			return nil

		case 1:
			// 单参数函数 (Event)
			results := fnValue.Call([]reflect.Value{reflect.ValueOf(event)})
			if len(results) > 0 && !results[0].IsNil() {
				return results[0].Interface().(error)
			}
			return nil

		default:
			return fmt.Errorf("handler function must have 0 or 1 parameter")
		}
	}

	return &FunctionListener{
		BaseListener: NewListener(name, handler),
	}
}

// MethodListener 方法监听器
type MethodListener struct {
	*BaseListener
	receiver interface{}
	method   string
}

// NewMethodListener 创建方法监听器
func NewMethodListener(name string, receiver interface{}, method string) *MethodListener {
	handler := func(event Event) error {
		receiverValue := reflect.ValueOf(receiver)
		methodValue := receiverValue.MethodByName(method)

		if !methodValue.IsValid() {
			return fmt.Errorf("method %s not found on receiver", method)
		}

		methodType := methodValue.Type()

		// 检查方法签名
		switch methodType.NumIn() {
		case 0:
			// 无参数方法
			results := methodValue.Call(nil)
			if len(results) > 0 && !results[0].IsNil() {
				return results[0].Interface().(error)
			}
			return nil

		case 1:
			// 单参数方法 (Event)
			results := methodValue.Call([]reflect.Value{reflect.ValueOf(event)})
			if len(results) > 0 && !results[0].IsNil() {
				return results[0].Interface().(error)
			}
			return nil

		default:
			return fmt.Errorf("method %s must have 0 or 1 parameter", method)
		}
	}

	return &MethodListener{
		BaseListener: NewListener(name, handler),
		receiver:     receiver,
		method:       method,
	}
}
