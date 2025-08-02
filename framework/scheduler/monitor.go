package scheduler

import (
	"fmt"
	"sync"
	"time"
)

// Monitor 监控器接口
type Monitor interface {
	// 任务执行监控
	RecordTaskStart(taskID string)
	RecordTaskComplete(taskID string, duration time.Duration, err error)
	RecordTaskFailed(taskID string, duration time.Duration, err error)

	// 调度器监控
	RecordSchedulerStart()
	RecordSchedulerStop()
	RecordSchedulerPause()
	RecordSchedulerResume()

	// 统计查询
	GetTaskMetrics(taskID string) (TaskMetrics, error)
	GetSchedulerMetrics() SchedulerMetrics
	GetPerformanceMetrics() PerformanceMetrics

	// 清理
	Cleanup()
}

// TaskMetrics 任务指标
type TaskMetrics struct {
	TaskID          string        `json:"task_id"`
	TaskName        string        `json:"task_name"`
	TotalExecutions int64         `json:"total_executions"`
	SuccessfulRuns  int64         `json:"successful_runs"`
	FailedRuns      int64         `json:"failed_runs"`
	SuccessRate     float64       `json:"success_rate"`
	AverageDuration time.Duration `json:"average_duration"`
	MinDuration     time.Duration `json:"min_duration"`
	MaxDuration     time.Duration `json:"max_duration"`
	LastExecution   time.Time     `json:"last_execution"`
	LastError       string        `json:"last_error"`
	TotalDuration   time.Duration `json:"total_duration"`
}

// SchedulerMetrics 调度器指标
type SchedulerMetrics struct {
	TotalTasks      int64         `json:"total_tasks"`
	EnabledTasks    int64         `json:"enabled_tasks"`
	RunningTasks    int64         `json:"running_tasks"`
	TotalExecutions int64         `json:"total_executions"`
	SuccessfulRuns  int64         `json:"successful_runs"`
	FailedRuns      int64         `json:"failed_runs"`
	SuccessRate     float64       `json:"success_rate"`
	Uptime          time.Duration `json:"uptime"`
	StartedAt       time.Time     `json:"started_at"`
	LastActivity    time.Time     `json:"last_activity"`
}

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	MemoryUsage      int64         `json:"memory_usage"`
	CPUUsage         float64       `json:"cpu_usage"`
	ActiveGoroutines int           `json:"active_goroutines"`
	TaskQueueSize    int           `json:"task_queue_size"`
	AverageLatency   time.Duration `json:"average_latency"`
	Throughput       float64       `json:"throughput"` // 任务/秒
}

// DefaultMonitor 默认监控器实现
type DefaultMonitor struct {
	taskMetrics        map[string]*TaskMetrics
	schedulerMetrics   SchedulerMetrics
	performanceMetrics PerformanceMetrics
	mu                 sync.RWMutex
	startedAt          time.Time
	lastActivity       time.Time
}

// NewMonitor 创建监控器
func NewMonitor() *DefaultMonitor {
	return &DefaultMonitor{
		taskMetrics:  make(map[string]*TaskMetrics),
		startedAt:    time.Now(),
		lastActivity: time.Now(),
	}
}

// RecordTaskStart 记录任务开始
func (m *DefaultMonitor) RecordTaskStart(taskID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.taskMetrics[taskID]; !exists {
		m.taskMetrics[taskID] = &TaskMetrics{
			TaskID: taskID,
		}
	}

	m.schedulerMetrics.RunningTasks++
	m.lastActivity = time.Now()
}

// RecordTaskComplete 记录任务完成
func (m *DefaultMonitor) RecordTaskComplete(taskID string, duration time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics, exists := m.taskMetrics[taskID]
	if !exists {
		metrics = &TaskMetrics{TaskID: taskID}
		m.taskMetrics[taskID] = metrics
	}

	metrics.TotalExecutions++
	metrics.TotalDuration += duration
	metrics.LastExecution = time.Now()

	if err == nil {
		metrics.SuccessfulRuns++
		m.schedulerMetrics.SuccessfulRuns++
	} else {
		metrics.FailedRuns++
		metrics.LastError = err.Error()
		m.schedulerMetrics.FailedRuns++
	}

	// 更新平均执行时间
	metrics.AverageDuration = metrics.TotalDuration / time.Duration(metrics.TotalExecutions)

	// 更新最小/最大执行时间
	if metrics.MinDuration == 0 || duration < metrics.MinDuration {
		metrics.MinDuration = duration
	}
	if duration > metrics.MaxDuration {
		metrics.MaxDuration = duration
	}

	// 更新成功率
	if metrics.TotalExecutions > 0 {
		metrics.SuccessRate = float64(metrics.SuccessfulRuns) / float64(metrics.TotalExecutions) * 100
	}

	// 更新调度器指标
	m.schedulerMetrics.TotalExecutions++
	m.schedulerMetrics.RunningTasks--
	if m.schedulerMetrics.TotalExecutions > 0 {
		m.schedulerMetrics.SuccessRate = float64(m.schedulerMetrics.SuccessfulRuns) / float64(m.schedulerMetrics.TotalExecutions) * 100
	}

	m.lastActivity = time.Now()
}

// RecordTaskFailed 记录任务失败
func (m *DefaultMonitor) RecordTaskFailed(taskID string, duration time.Duration, err error) {
	m.RecordTaskComplete(taskID, duration, err)
}

// RecordSchedulerStart 记录调度器启动
func (m *DefaultMonitor) RecordSchedulerStart() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.schedulerMetrics.StartedAt = time.Now()
	m.lastActivity = time.Now()
}

// RecordSchedulerStop 记录调度器停止
func (m *DefaultMonitor) RecordSchedulerStop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastActivity = time.Now()
}

// RecordSchedulerPause 记录调度器暂停
func (m *DefaultMonitor) RecordSchedulerPause() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastActivity = time.Now()
}

// RecordSchedulerResume 记录调度器恢复
func (m *DefaultMonitor) RecordSchedulerResume() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastActivity = time.Now()
}

// GetTaskMetrics 获取任务指标
func (m *DefaultMonitor) GetTaskMetrics(taskID string) (TaskMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics, exists := m.taskMetrics[taskID]
	if !exists {
		return TaskMetrics{}, fmt.Errorf("task metrics not found: %s", taskID)
	}

	return *metrics, nil
}

// GetSchedulerMetrics 获取调度器指标
func (m *DefaultMonitor) GetSchedulerMetrics() SchedulerMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := m.schedulerMetrics
	metrics.Uptime = time.Since(m.startedAt)
	metrics.LastActivity = m.lastActivity

	return metrics
}

// GetPerformanceMetrics 获取性能指标
func (m *DefaultMonitor) GetPerformanceMetrics() PerformanceMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 这里应该实现实际的性能监控
	// 由于没有具体的性能监控库，这里只是示例
	metrics := m.performanceMetrics

	// 计算吞吐量
	if m.schedulerMetrics.Uptime > 0 {
		metrics.Throughput = float64(m.schedulerMetrics.TotalExecutions) / m.schedulerMetrics.Uptime.Seconds()
	}

	return metrics
}

// Cleanup 清理监控数据
func (m *DefaultMonitor) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 清理过期的任务指标（保留最近1000个任务）
	if len(m.taskMetrics) > 1000 {
		// 这里可以实现更复杂的清理策略
		// 比如按时间清理或按执行次数清理
		m.taskMetrics = make(map[string]*TaskMetrics)
	}
}

// SetTaskName 设置任务名称
func (m *DefaultMonitor) SetTaskName(taskID, taskName string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if metrics, exists := m.taskMetrics[taskID]; exists {
		metrics.TaskName = taskName
	}
}

// UpdateSchedulerStats 更新调度器统计
func (m *DefaultMonitor) UpdateSchedulerStats(totalTasks, enabledTasks int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.schedulerMetrics.TotalTasks = totalTasks
	m.schedulerMetrics.EnabledTasks = enabledTasks
}

// GetTopTasks 获取执行最多的任务
func (m *DefaultMonitor) GetTopTasks(limit int) []TaskMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var tasks []TaskMetrics
	for _, metrics := range m.taskMetrics {
		tasks = append(tasks, *metrics)
	}

	// 按执行次数排序
	// 这里应该实现排序逻辑
	// 为了简化，直接返回前limit个

	if len(tasks) > limit {
		return tasks[:limit]
	}

	return tasks
}

// GetFailedTasks 获取失败的任务
func (m *DefaultMonitor) GetFailedTasks() []TaskMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var failedTasks []TaskMetrics
	for _, metrics := range m.taskMetrics {
		if metrics.FailedRuns > 0 {
			failedTasks = append(failedTasks, *metrics)
		}
	}

	return failedTasks
}

// ResetMetrics 重置指标
func (m *DefaultMonitor) ResetMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.taskMetrics = make(map[string]*TaskMetrics)
	m.schedulerMetrics = SchedulerMetrics{}
	m.performanceMetrics = PerformanceMetrics{}
	m.startedAt = time.Now()
	m.lastActivity = time.Now()
}
