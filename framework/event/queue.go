package event

import (
	"context"
	"sync"
	"time"
)

// MemoryEventQueue 内存事件队列实现
type MemoryEventQueue struct {
	mu       sync.RWMutex
	events   []Event
	closed   bool
	waitChan chan struct{} // 用于通知等待的消费者
}

// NewMemoryEventQueue 创建内存事件队列
func NewMemoryEventQueue() *MemoryEventQueue {
	return &MemoryEventQueue{
		events:   make([]Event, 0),
		waitChan: make(chan struct{}, 1),
	}
}

// Push 推送事件
func (q *MemoryEventQueue) Push(event Event) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrDispatcherClosed
	}

	q.events = append(q.events, event)
	
	// 通知等待的消费者
	select {
	case q.waitChan <- struct{}{}:
	default:
		// 通道已满，说明已经有消费者被通知
	}
	
	return nil
}

// PushBatch 批量推送事件
func (q *MemoryEventQueue) PushBatch(events []Event) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrDispatcherClosed
	}

	q.events = append(q.events, events...)
	
	// 通知等待的消费者
	select {
	case q.waitChan <- struct{}{}:
	default:
		// 通道已满，说明已经有消费者被通知
	}
	
	return nil
}

// Pop 弹出事件
func (q *MemoryEventQueue) Pop(ctx context.Context) (Event, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-q.waitChan:
			// 有事件到达，立即检查
		case <-ticker.C:
			// 定期检查，避免无限等待
		}
		
		q.mu.Lock()
		
		if q.closed {
			q.mu.Unlock()
			return nil, ErrDispatcherClosed
		}
		
		if len(q.events) > 0 {
			// 弹出第一个事件
			event := q.events[0]
			q.events = q.events[1:]
			
			// 如果还有事件，继续通知其他消费者
			if len(q.events) > 0 {
				select {
				case q.waitChan <- struct{}{}:
				default:
				}
			}
			
			q.mu.Unlock()
			return event, nil
		}
		
		q.mu.Unlock()
	}
}

// PopBatch 批量弹出事件
func (q *MemoryEventQueue) PopBatch(ctx context.Context, count int) ([]Event, error) {
	var events []Event
	
	for i := 0; i < count; i++ {
		event, err := q.Pop(ctx)
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				break
			}
			return events, err
		}
		events = append(events, event)
	}
	
	return events, nil
}

// Size 获取队列大小
func (q *MemoryEventQueue) Size() (int, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	
	if q.closed {
		return 0, ErrDispatcherClosed
	}
	
	return len(q.events), nil
}

// Clear 清空队列
func (q *MemoryEventQueue) Clear() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	if q.closed {
		return ErrDispatcherClosed
	}
	
	// 清空事件数组，释放内存
	q.events = nil
	q.events = make([]Event, 0)
	
	return nil
}

// Close 关闭队列
func (q *MemoryEventQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	if q.closed {
		return nil
	}
	
	q.closed = true
	
	// 清空事件，释放内存
	q.events = nil
	
	// 关闭通知通道
	close(q.waitChan)
	
	return nil
}
