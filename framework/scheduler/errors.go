package scheduler

import "errors"

// 调度器相关错误
var (
	ErrTaskNotFound            = errors.New("task not found")
	ErrSchedulerAlreadyRunning = errors.New("scheduler is already running")
	ErrSchedulerAlreadyStopped = errors.New("scheduler is already stopped")
	ErrSchedulerNotRunning     = errors.New("scheduler is not running")
	ErrSchedulerNotPaused      = errors.New("scheduler is not paused")
	ErrInvalidSchedule         = errors.New("invalid schedule format")
	ErrTaskHandlerRequired     = errors.New("task handler is required")
	ErrTaskNameRequired        = errors.New("task name is required")
	ErrTaskIDRequired          = errors.New("task ID is required")
	ErrStoreNotInitialized     = errors.New("store is not initialized")
	ErrTaskExecutionTimeout    = errors.New("task execution timeout")
	ErrTaskExecutionFailed     = errors.New("task execution failed")
	ErrInvalidCronExpression   = errors.New("invalid cron expression")
	ErrInvalidTimeFormat       = errors.New("invalid time format")
	ErrTaskAlreadyExists       = errors.New("task already exists")
	ErrTaskDisabled            = errors.New("task is disabled")
	ErrTaskMaxRetriesExceeded  = errors.New("task max retries exceeded")
)
