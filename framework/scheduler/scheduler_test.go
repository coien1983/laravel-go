package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	handler := NewFuncHandler("test", func(ctx context.Context) error {
		return nil
	})

	task := NewTask("test-task", "Test task", "0 * * * * *", handler)

	if task.GetID() == "" {
		t.Error("Task ID should not be empty")
	}

	if task.GetName() != "test-task" {
		t.Errorf("Expected name 'test-task', got '%s'", task.GetName())
	}

	if task.GetDescription() != "Test task" {
		t.Errorf("Expected description 'Test task', got '%s'", task.GetDescription())
	}

	if task.GetSchedule() != "0 * * * * *" {
		t.Errorf("Expected schedule '0 * * * * *', got '%s'", task.GetSchedule())
	}

	if !task.GetEnabled() {
		t.Error("Task should be enabled by default")
	}
}

func TestTaskValidation(t *testing.T) {
	handler := NewFuncHandler("test", func(ctx context.Context) error {
		return nil
	})

	// 测试空名称
	task := NewTask("", "Test task", "0 * * * * *", handler)
	if err := task.Validate(); err != ErrTaskNameRequired {
		t.Errorf("Expected ErrTaskNameRequired, got %v", err)
	}

	// 测试空处理器
	task = NewTask("test", "Test task", "0 * * * * *", nil)
	if err := task.Validate(); err != ErrTaskHandlerRequired {
		t.Errorf("Expected ErrTaskHandlerRequired, got %v", err)
	}

	// 测试空调度表达式
	task = NewTask("test", "Test task", "", handler)
	if err := task.Validate(); err == nil {
		t.Error("Expected error for empty schedule")
	}

	// 测试无效调度表达式
	task = NewTask("test", "Test task", "invalid", handler)
	if err := task.Validate(); err == nil {
		t.Error("Expected error for invalid schedule")
	}
}

func TestParseSchedule(t *testing.T) {
	// 测试标准 Cron 表达式
	nextRun, err := ParseSchedule("0 0 2 * * *")
	if err != nil {
		t.Errorf("Failed to parse schedule: %v", err)
	}

	if nextRun.Before(time.Now()) {
		t.Error("Next run time should be in the future")
	}

	// 测试特殊表达式
	nextRun, err = ParseSchedule("@hourly")
	if err != nil {
		t.Errorf("Failed to parse @hourly: %v", err)
	}

	if nextRun.Before(time.Now()) {
		t.Error("Next run time should be in the future")
	}

	// 测试简单时间格式
	nextRun, err = ParseSchedule("15:30")
	if err != nil {
		t.Errorf("Failed to parse time format: %v", err)
	}

	if nextRun.Before(time.Now()) {
		t.Error("Next run time should be in the future")
	}
}

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore()

	handler := NewFuncHandler("test", func(ctx context.Context) error {
		return nil
	})

	task := NewTask("test-task", "Test task", "0 * * * * *", handler)

	// 测试保存
	err := store.Save(task)
	if err != nil {
		t.Errorf("Failed to save task: %v", err)
	}

	// 测试获取
	retrievedTask, err := store.Get(task.GetID())
	if err != nil {
		t.Errorf("Failed to get task: %v", err)
	}

	if retrievedTask.GetID() != task.GetID() {
		t.Errorf("Task ID mismatch: expected %s, got %s", task.GetID(), retrievedTask.GetID())
	}

	// 测试获取所有任务
	tasks, err := store.GetAll()
	if err != nil {
		t.Errorf("Failed to get all tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	// 测试删除
	err = store.Delete(task.GetID())
	if err != nil {
		t.Errorf("Failed to delete task: %v", err)
	}

	// 测试删除后获取
	_, err = store.Get(task.GetID())
	if err != ErrTaskNotFound {
		t.Errorf("Expected ErrTaskNotFound, got %v", err)
	}
}

func TestScheduler(t *testing.T) {
	store := NewMemoryStore()
	scheduler := NewScheduler(store)

	handler := NewFuncHandler("test", func(ctx context.Context) error {
		return nil
	})

	task := NewTask("test-task", "Test task", "0 * * * * *", handler)

	// 测试添加任务
	err := scheduler.Add(task)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	// 测试获取任务
	retrievedTask, err := scheduler.Get(task.GetID())
	if err != nil {
		t.Errorf("Failed to get task: %v", err)
	}

	if retrievedTask.GetID() != task.GetID() {
		t.Errorf("Task ID mismatch: expected %s, got %s", task.GetID(), retrievedTask.GetID())
	}

	// 测试获取所有任务
	tasks := scheduler.GetAll()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	// 测试获取启用的任务
	enabledTasks := scheduler.GetEnabled()
	if len(enabledTasks) != 1 {
		t.Errorf("Expected 1 enabled task, got %d", len(enabledTasks))
	}

	// 测试删除任务
	err = scheduler.Remove(task.GetID())
	if err != nil {
		t.Errorf("Failed to remove task: %v", err)
	}

	// 测试删除后获取
	_, err = scheduler.Get(task.GetID())
	if err != ErrTaskNotFound {
		t.Errorf("Expected ErrTaskNotFound, got %v", err)
	}
}

func TestMonitor(t *testing.T) {
	monitor := NewMonitor()

	// 测试记录任务开始
	monitor.RecordTaskStart("test-task")

	// 测试记录任务完成
	monitor.RecordTaskComplete("test-task", time.Second, nil)

	// 测试获取任务指标
	metrics, err := monitor.GetTaskMetrics("test-task")
	if err != nil {
		t.Errorf("Failed to get task metrics: %v", err)
	}

	if metrics.TotalExecutions != 1 {
		t.Errorf("Expected 1 execution, got %d", metrics.TotalExecutions)
	}

	if metrics.SuccessfulRuns != 1 {
		t.Errorf("Expected 1 successful run, got %d", metrics.SuccessfulRuns)
	}

	// 测试记录任务失败
	monitor.RecordTaskFailed("test-task", time.Second, fmt.Errorf("test error"))

	metrics, err = monitor.GetTaskMetrics("test-task")
	if err != nil {
		t.Errorf("Failed to get task metrics: %v", err)
	}

	if metrics.TotalExecutions != 2 {
		t.Errorf("Expected 2 executions, got %d", metrics.TotalExecutions)
	}

	if metrics.FailedRuns != 1 {
		t.Errorf("Expected 1 failed run, got %d", metrics.FailedRuns)
	}
}

func TestConvenienceFunctions(t *testing.T) {
	handler := NewFuncHandler("test", func(ctx context.Context) error {
		return nil
	})

	// 测试便捷方法
	task1 := Every(5, handler)
	if task1.GetSchedule() != "0 */5 * * * *" {
		t.Errorf("Expected schedule '0 */5 * * * *', got '%s'", task1.GetSchedule())
	}

	task2 := EveryHour(handler)
	if task2.GetSchedule() != "0 0 * * * *" {
		t.Errorf("Expected schedule '0 0 * * * *', got '%s'", task2.GetSchedule())
	}

	task3 := EveryDay(handler)
	if task3.GetSchedule() != "0 0 0 * * *" {
		t.Errorf("Expected schedule '0 0 0 * * *', got '%s'", task3.GetSchedule())
	}

	task4 := Daily(9, 30, handler)
	if task4.GetSchedule() != "0 30 9 * * *" {
		t.Errorf("Expected schedule '0 30 9 * * *', got '%s'", task4.GetSchedule())
	}

	task5 := Cron("0 0 2 * * *", handler)
	if task5.GetSchedule() != "0 0 2 * * *" {
		t.Errorf("Expected schedule '0 0 2 * * *', got '%s'", task5.GetSchedule())
	}
}

func TestTaskBuilder(t *testing.T) {
	handler := NewFuncHandler("test", func(ctx context.Context) error {
		return nil
	})

	builder := NewTaskBuilder("test-task", "Test task", "0 * * * * *", handler)

	task := builder.
		SetTimeout(5*time.Minute).
		SetMaxRetries(3).
		SetRetryDelay(30*time.Second).
		AddTag("test", "true").
		Build()

	if task.GetTimeout() != 5*time.Minute {
		t.Errorf("Expected timeout 5m, got %v", task.GetTimeout())
	}

	if task.GetMaxRetries() != 3 {
		t.Errorf("Expected max retries 3, got %d", task.GetMaxRetries())
	}

	if task.GetRetryDelay() != 30*time.Second {
		t.Errorf("Expected retry delay 30s, got %v", task.GetRetryDelay())
	}

	tags := task.GetTags()
	if tags["test"] != "true" {
		t.Errorf("Expected tag 'test'='true', got '%s'", tags["test"])
	}
}
