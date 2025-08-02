package queue

import (
	"errors"
	"fmt"
)

// 队列系统错误定义
var (
	ErrQueueNotFound     = errors.New("queue not found")
	ErrJobNotFound       = errors.New("job not found")
	ErrJobTimeout        = errors.New("job timeout")
	ErrJobMaxAttempts    = errors.New("job exceeded max attempts")
	ErrJobSerialization  = errors.New("job serialization failed")
	ErrJobDeserialization = errors.New("job deserialization failed")
	ErrQueueClosed       = errors.New("queue is closed")
	ErrWorkerStopped     = errors.New("worker is stopped")
	ErrInvalidJob        = errors.New("invalid job")
	ErrQueueFull         = errors.New("queue is full")
	ErrQueueEmpty        = errors.New("queue is empty")
)

// QueueError 队列错误
type QueueError struct {
	Queue   string
	Message string
	Err     error
}

func (e *QueueError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("queue error [%s]: %s: %v", e.Queue, e.Message, e.Err)
	}
	return fmt.Sprintf("queue error [%s]: %s", e.Queue, e.Message)
}

func (e *QueueError) Unwrap() error {
	return e.Err
}

// JobError 任务错误
type JobError struct {
	JobID   string
	Message string
	Err     error
}

func (e *JobError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("job error [%s]: %s: %v", e.JobID, e.Message, e.Err)
	}
	return fmt.Sprintf("job error [%s]: %s", e.JobID, e.Message)
}

func (e *JobError) Unwrap() error {
	return e.Err
}

// WorkerError 工作进程错误
type WorkerError struct {
	WorkerID string
	Message  string
	Err      error
}

func (e *WorkerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("worker error [%s]: %s: %v", e.WorkerID, e.Message, e.Err)
	}
	return fmt.Sprintf("worker error [%s]: %s", e.WorkerID, e.Message)
}

func (e *WorkerError) Unwrap() error {
	return e.Err
} 