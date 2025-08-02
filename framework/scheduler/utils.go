package scheduler

import (
	"context"
	"fmt"
	"time"
)

// 全局调度器实例
var (
	GlobalScheduler Scheduler
	GlobalMonitor   Monitor
	initialized     bool
)

// Init 初始化全局调度器
func Init(store Store) {
	if initialized {
		return
	}

	GlobalScheduler = NewScheduler(store)
	GlobalMonitor = NewMonitor()
	initialized = true
}

// GetScheduler 获取全局调度器
func GetScheduler() Scheduler {
	if !initialized {
		panic("scheduler not initialized, call scheduler.Init() first")
	}
	return GlobalScheduler
}

// GetMonitor 获取全局监控器
func GetMonitor() Monitor {
	if !initialized {
		panic("scheduler not initialized, call scheduler.Init() first")
	}
	return GlobalMonitor
}

// 便捷的任务创建函数

// Every 创建每分钟执行的任务
func Every(minutes int, handler TaskHandler) *DefaultTask {
	if minutes <= 0 {
		minutes = 1
	}

	schedule := fmt.Sprintf("0 */%d * * * *", minutes)
	return NewTask("", "", schedule, handler)
}

// EveryHour 创建每小时执行的任务
func EveryHour(handler TaskHandler) *DefaultTask {
	return NewTask("", "", "0 0 * * * *", handler)
}

// EveryDay 创建每天执行的任务
func EveryDay(handler TaskHandler) *DefaultTask {
	return NewTask("", "", "0 0 0 * * *", handler)
}

// EveryWeek 创建每周执行的任务
func EveryWeek(handler TaskHandler) *DefaultTask {
	return NewTask("", "", "0 0 0 * * 0", handler)
}

// EveryMonth 创建每月执行的任务
func EveryMonth(handler TaskHandler) *DefaultTask {
	return NewTask("", "", "0 0 0 1 * *", handler)
}

// Daily 创建每天指定时间执行的任务
func Daily(hour, minute int, handler TaskHandler) *DefaultTask {
	schedule := fmt.Sprintf("0 %d %d * * *", minute, hour)
	return NewTask("", "", schedule, handler)
}

// Weekly 创建每周指定时间执行的任务
func Weekly(weekday time.Weekday, hour, minute int, handler TaskHandler) *DefaultTask {
	schedule := fmt.Sprintf("0 %d %d * * %d", minute, hour, int(weekday))
	return NewTask("", "", schedule, handler)
}

// Monthly 创建每月指定日期和时间执行的任务
func Monthly(day, hour, minute int, handler TaskHandler) *DefaultTask {
	schedule := fmt.Sprintf("0 %d %d %d * *", minute, hour, day)
	return NewTask("", "", schedule, handler)
}

// Cron 创建自定义 Cron 表达式的任务
func Cron(expression string, handler TaskHandler) *DefaultTask {
	return NewTask("", "", expression, handler)
}

// At 创建指定时间执行的任务
func At(timeStr string, handler TaskHandler) *DefaultTask {
	return NewTask("", "", timeStr, handler)
}

// 便捷的任务处理器

// FuncHandler 函数处理器
type FuncHandler struct {
	name string
	fn   func(ctx context.Context) error
}

// NewFuncHandler 创建函数处理器
func NewFuncHandler(name string, fn func(ctx context.Context) error) *FuncHandler {
	return &FuncHandler{
		name: name,
		fn:   fn,
	}
}

// Handle 执行函数
func (h *FuncHandler) Handle(ctx context.Context) error {
	return h.fn(ctx)
}

// GetName 获取处理器名称
func (h *FuncHandler) GetName() string {
	return h.name
}

// 便捷的调度器操作

// AddTask 添加任务到全局调度器
func AddTask(task Task) error {
	return GetScheduler().Add(task)
}

// RemoveTask 从全局调度器移除任务
func RemoveTask(taskID string) error {
	return GetScheduler().Remove(taskID)
}

// UpdateTask 更新全局调度器中的任务
func UpdateTask(task Task) error {
	return GetScheduler().Update(task)
}

// GetTask 获取全局调度器中的任务
func GetTask(taskID string) (Task, error) {
	return GetScheduler().Get(taskID)
}

// GetAllTasks 获取全局调度器中的所有任务
func GetAllTasks() []Task {
	return GetScheduler().GetAll()
}

// GetEnabledTasks 获取全局调度器中的启用任务
func GetEnabledTasks() []Task {
	return GetScheduler().GetEnabled()
}

// StartScheduler 启动全局调度器
func StartScheduler() error {
	monitor := GetMonitor()
	monitor.RecordSchedulerStart()
	return GetScheduler().Start()
}

// StopScheduler 停止全局调度器
func StopScheduler() error {
	monitor := GetMonitor()
	monitor.RecordSchedulerStop()
	return GetScheduler().Stop()
}

// PauseScheduler 暂停全局调度器
func PauseScheduler() error {
	monitor := GetMonitor()
	monitor.RecordSchedulerPause()
	return GetScheduler().Pause()
}

// ResumeScheduler 恢复全局调度器
func ResumeScheduler() error {
	monitor := GetMonitor()
	monitor.RecordSchedulerResume()
	return GetScheduler().Resume()
}

// RunTaskNow 立即运行任务
func RunTaskNow(taskID string) error {
	return GetScheduler().RunNow(taskID)
}

// RunAllTasks 运行所有启用的任务
func RunAllTasks() error {
	return GetScheduler().RunAll()
}

// GetSchedulerStatus 获取全局调度器状态
func GetSchedulerStatus() SchedulerStatus {
	return GetScheduler().GetStatus()
}

// GetSchedulerStats 获取全局调度器统计
func GetSchedulerStats() SchedulerStats {
	return GetScheduler().GetStats()
}

// GetTaskStats 获取任务统计
func GetTaskStats(taskID string) (TaskStats, error) {
	return GetScheduler().GetTaskStats(taskID)
}

// GetTaskMetrics 获取任务指标
func GetTaskMetrics(taskID string) (TaskMetrics, error) {
	return GetMonitor().GetTaskMetrics(taskID)
}

// GetSchedulerMetrics 获取调度器指标
func GetSchedulerMetrics() SchedulerMetrics {
	return GetMonitor().GetSchedulerMetrics()
}

// GetPerformanceMetrics 获取性能指标
func GetPerformanceMetrics() PerformanceMetrics {
	return GetMonitor().GetPerformanceMetrics()
}

// 任务构建器

// TaskBuilder 任务构建器
type TaskBuilder struct {
	task *DefaultTask
}

// NewTaskBuilder 创建任务构建器
func NewTaskBuilder(name, description, schedule string, handler TaskHandler) *TaskBuilder {
	return &TaskBuilder{
		task: NewTask(name, description, schedule, handler),
	}
}

// SetTimeout 设置超时时间
func (b *TaskBuilder) SetTimeout(timeout time.Duration) *TaskBuilder {
	b.task.SetTimeout(timeout)
	return b
}

// SetMaxRetries 设置最大重试次数
func (b *TaskBuilder) SetMaxRetries(maxRetries int) *TaskBuilder {
	b.task.SetMaxRetries(maxRetries)
	return b
}

// SetRetryDelay 设置重试延迟
func (b *TaskBuilder) SetRetryDelay(retryDelay time.Duration) *TaskBuilder {
	b.task.SetRetryDelay(retryDelay)
	return b
}

// AddTag 添加标签
func (b *TaskBuilder) AddTag(key, value string) *TaskBuilder {
	b.task.AddTag(key, value)
	return b
}

// Disable 禁用任务
func (b *TaskBuilder) Disable() *TaskBuilder {
	b.task.Disable()
	return b
}

// Build 构建任务
func (b *TaskBuilder) Build() *DefaultTask {
	return b.task
}

// 调度器配置

// SchedulerConfig 调度器配置
type SchedulerConfig struct {
	Store          Store
	Monitor        Monitor
	CheckInterval  time.Duration
	MaxConcurrency int
	EnableMetrics  bool
	EnableLogging  bool
}

// NewSchedulerConfig 创建调度器配置
func NewSchedulerConfig() *SchedulerConfig {
	return &SchedulerConfig{
		CheckInterval:  time.Second,
		MaxConcurrency: 10,
		EnableMetrics:  true,
		EnableLogging:  true,
	}
}

// WithStore 设置存储
func (c *SchedulerConfig) WithStore(store Store) *SchedulerConfig {
	c.Store = store
	return c
}

// WithMonitor 设置监控器
func (c *SchedulerConfig) WithMonitor(monitor Monitor) *SchedulerConfig {
	c.Monitor = monitor
	return c
}

// WithCheckInterval 设置检查间隔
func (c *SchedulerConfig) WithCheckInterval(interval time.Duration) *SchedulerConfig {
	c.CheckInterval = interval
	return c
}

// WithMaxConcurrency 设置最大并发数
func (c *SchedulerConfig) WithMaxConcurrency(max int) *SchedulerConfig {
	c.MaxConcurrency = max
	return c
}

// WithMetrics 设置是否启用指标
func (c *SchedulerConfig) WithMetrics(enable bool) *SchedulerConfig {
	c.EnableMetrics = enable
	return c
}

// WithLogging 设置是否启用日志
func (c *SchedulerConfig) WithLogging(enable bool) *SchedulerConfig {
	c.EnableLogging = enable
	return c
}

// Build 构建调度器
func (c *SchedulerConfig) Build() Scheduler {
	if c.Store == nil {
		c.Store = NewMemoryStore()
	}

	if c.Monitor == nil && c.EnableMetrics {
		c.Monitor = NewMonitor()
	}

	return NewScheduler(c.Store)
}
