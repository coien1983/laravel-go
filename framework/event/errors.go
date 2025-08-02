package event

import (
	"errors"
	"fmt"
)

// 事件系统错误定义
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

// EventError 事件错误
type EventError struct {
	EventName string
	Message   string
	Err       error
}

func (e *EventError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("event error [%s]: %s: %v", e.EventName, e.Message, e.Err)
	}
	return fmt.Sprintf("event error [%s]: %s", e.EventName, e.Message)
}

func (e *EventError) Unwrap() error {
	return e.Err
}

// ListenerError 监听器错误
type ListenerError struct {
	ListenerName string
	EventName    string
	Message      string
	Err          error
}

func (e *ListenerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("listener error [%s] for event [%s]: %s: %v", e.ListenerName, e.EventName, e.Message, e.Err)
	}
	return fmt.Sprintf("listener error [%s] for event [%s]: %s", e.ListenerName, e.EventName, e.Message)
}

func (e *ListenerError) Unwrap() error {
	return e.Err
}

// DispatcherError 分发器错误
type DispatcherError struct {
	DispatcherName string
	Message        string
	Err            error
}

func (e *DispatcherError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("dispatcher error [%s]: %s: %v", e.DispatcherName, e.Message, e.Err)
	}
	return fmt.Sprintf("dispatcher error [%s]: %s", e.DispatcherName, e.Message)
}

func (e *DispatcherError) Unwrap() error {
	return e.Err
}
