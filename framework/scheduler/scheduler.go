package scheduler

import (
	"context"
	"sync"
	"time"
)

// Task 表示一个定时任务
type Task interface {
	// 基础信息
	GetID() string
	GetName() string
	GetDescription() string
	GetSchedule() string
	GetHandler() TaskHandler
	GetEnabled() bool
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time

	// 执行信息
	GetLastRunAt() *time.Time
	GetNextRunAt() *time.Time
	GetRunCount() int64
	GetFailedCount() int64
	GetLastError() string

	// 配置
	GetTimeout() time.Duration
	GetRetryCount() int
	GetRetryDelay() time.Duration
	GetMaxRetries() int
	GetTags() map[string]string

	// 状态管理
	Enable()
	Disable()
	UpdateNextRun()
	MarkAsRun()
	MarkAsFailed(error)
	IncrementRunCount()
	IncrementFailedCount()

	// 序列化
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}

// TaskHandler 任务处理器接口
type TaskHandler interface {
	Handle(ctx context.Context) error
	GetName() string
}

// Scheduler 调度器接口
type Scheduler interface {
	// 任务管理
	Add(task Task) error
	Remove(taskID string) error
	Update(task Task) error
	Get(taskID string) (Task, error)
	GetAll() []Task
	GetEnabled() []Task

	// 调度控制
	Start() error
	Stop() error
	Pause() error
	Resume() error

	// 任务执行
	RunNow(taskID string) error
	RunAll() error

	// 监控
	GetStatus() SchedulerStatus
	GetStats() SchedulerStats
	GetTaskStats(taskID string) (TaskStats, error)
}

// SchedulerStatus 调度器状态
type SchedulerStatus struct {
	Status    string    `json:"status"` // running, paused, stopped
	StartedAt time.Time `json:"started_at"`
	TaskCount int       `json:"task_count"`
	Running   int       `json:"running"`
	Failed    int       `json:"failed"`
}

// SchedulerStats 调度器统计
type SchedulerStats struct {
	TotalTasks    int64     `json:"total_tasks"`
	EnabledTasks  int64     `json:"enabled_tasks"`
	DisabledTasks int64     `json:"disabled_tasks"`
	TotalRuns     int64     `json:"total_runs"`
	TotalFailed   int64     `json:"total_failed"`
	SuccessRate   float64   `json:"success_rate"`
	LastRunAt     time.Time `json:"last_run_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// TaskStats 任务统计
type TaskStats struct {
	TaskID      string        `json:"task_id"`
	TaskName    string        `json:"task_name"`
	RunCount    int64         `json:"run_count"`
	FailedCount int64         `json:"failed_count"`
	SuccessRate float64       `json:"success_rate"`
	LastRunAt   time.Time     `json:"last_run_at"`
	NextRunAt   time.Time     `json:"next_run_at"`
	AverageTime time.Duration `json:"average_time"`
	LastError   string        `json:"last_error"`
}

// Store 任务存储接口
type Store interface {
	// 基础操作
	Save(task Task) error
	Get(taskID string) (Task, error)
	GetAll() ([]Task, error)
	Delete(taskID string) error
	Clear() error

	// 批量操作
	SaveBatch(tasks []Task) error
	GetByTags(tags map[string]string) ([]Task, error)

	// 统计
	GetStats() (StoreStats, error)
	Close() error
}

// StoreStats 存储统计
type StoreStats struct {
	TotalTasks   int64     `json:"total_tasks"`
	EnabledTasks int64     `json:"enabled_tasks"`
	LastSync     time.Time `json:"last_sync"`
}

// 主调度器实现
type DefaultScheduler struct {
	store      Store
	tasks      map[string]Task
	mu         sync.RWMutex
	status     SchedulerStatus
	stats      SchedulerStats
	stopChan   chan struct{}
	pauseChan  chan struct{}
	resumeChan chan struct{}
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewScheduler 创建新的调度器
func NewScheduler(store Store) *DefaultScheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &DefaultScheduler{
		store:      store,
		tasks:      make(map[string]Task),
		status:     SchedulerStatus{Status: "stopped"},
		stats:      SchedulerStats{CreatedAt: time.Now()},
		stopChan:   make(chan struct{}),
		pauseChan:  make(chan struct{}),
		resumeChan: make(chan struct{}),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Add 添加任务
func (s *DefaultScheduler) Add(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 保存到存储
	if err := s.store.Save(task); err != nil {
		return err
	}

	// 添加到内存
	s.tasks[task.GetID()] = task
	s.stats.TotalTasks++

	if task.GetEnabled() {
		s.stats.EnabledTasks++
	} else {
		s.stats.DisabledTasks++
	}

	return nil
}

// Remove 移除任务
func (s *DefaultScheduler) Remove(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 从存储删除
	if err := s.store.Delete(taskID); err != nil {
		return err
	}

	// 从内存删除
	if task, exists := s.tasks[taskID]; exists {
		if task.GetEnabled() {
			s.stats.EnabledTasks--
		} else {
			s.stats.DisabledTasks--
		}
		s.stats.TotalTasks--
		delete(s.tasks, taskID)
	}

	return nil
}

// Update 更新任务
func (s *DefaultScheduler) Update(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 更新存储
	if err := s.store.Save(task); err != nil {
		return err
	}

	// 更新内存
	oldTask, exists := s.tasks[task.GetID()]
	if exists {
		if oldTask.GetEnabled() && !task.GetEnabled() {
			s.stats.EnabledTasks--
			s.stats.DisabledTasks++
		} else if !oldTask.GetEnabled() && task.GetEnabled() {
			s.stats.EnabledTasks++
			s.stats.DisabledTasks--
		}
	}

	s.tasks[task.GetID()] = task
	return nil
}

// Get 获取任务
func (s *DefaultScheduler) Get(taskID string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

// GetAll 获取所有任务
func (s *DefaultScheduler) GetAll() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// GetEnabled 获取启用的任务
func (s *DefaultScheduler) GetEnabled() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0)
	for _, task := range s.tasks {
		if task.GetEnabled() {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

// Start 启动调度器
func (s *DefaultScheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.Status == "running" {
		return ErrSchedulerAlreadyRunning
	}

	// 从存储加载任务
	tasks, err := s.store.GetAll()
	if err != nil {
		return err
	}

	// 加载到内存
	for _, task := range tasks {
		s.tasks[task.GetID()] = task
		if task.GetEnabled() {
			s.stats.EnabledTasks++
		} else {
			s.stats.DisabledTasks++
		}
	}
	s.stats.TotalTasks = int64(len(tasks))

	// 启动调度循环
	go s.scheduleLoop()

	s.status.Status = "running"
	s.status.StartedAt = time.Now()
	s.status.TaskCount = len(s.tasks)

	return nil
}

// Stop 停止调度器
func (s *DefaultScheduler) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.Status == "stopped" {
		return ErrSchedulerAlreadyStopped
	}

	close(s.stopChan)
	s.cancel()

	s.status.Status = "stopped"
	return nil
}

// Pause 暂停调度器
func (s *DefaultScheduler) Pause() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.Status != "running" {
		return ErrSchedulerNotRunning
	}

	s.pauseChan <- struct{}{}
	s.status.Status = "paused"

	return nil
}

// Resume 恢复调度器
func (s *DefaultScheduler) Resume() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.Status != "paused" {
		return ErrSchedulerNotPaused
	}

	s.resumeChan <- struct{}{}
	s.status.Status = "running"

	return nil
}

// RunNow 立即运行任务
func (s *DefaultScheduler) RunNow(taskID string) error {
	task, err := s.Get(taskID)
	if err != nil {
		return err
	}

	go s.executeTask(task)
	return nil
}

// RunAll 运行所有启用的任务
func (s *DefaultScheduler) RunAll() error {
	tasks := s.GetEnabled()

	for _, task := range tasks {
		go s.executeTask(task)
	}

	return nil
}

// GetStatus 获取调度器状态
func (s *DefaultScheduler) GetStatus() SchedulerStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.status
}

// GetStats 获取调度器统计
func (s *DefaultScheduler) GetStats() SchedulerStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.stats.TotalRuns > 0 {
		s.stats.SuccessRate = float64(s.stats.TotalRuns-s.stats.TotalFailed) / float64(s.stats.TotalRuns) * 100
	}

	return s.stats
}

// GetTaskStats 获取任务统计
func (s *DefaultScheduler) GetTaskStats(taskID string) (TaskStats, error) {
	task, err := s.Get(taskID)
	if err != nil {
		return TaskStats{}, err
	}

	stats := TaskStats{
		TaskID:      task.GetID(),
		TaskName:    task.GetName(),
		RunCount:    task.GetRunCount(),
		FailedCount: task.GetFailedCount(),
		LastError:   task.GetLastError(),
	}

	if task.GetLastRunAt() != nil {
		stats.LastRunAt = *task.GetLastRunAt()
	}

	if task.GetNextRunAt() != nil {
		stats.NextRunAt = *task.GetNextRunAt()
	}

	if stats.RunCount > 0 {
		stats.SuccessRate = float64(stats.RunCount-stats.FailedCount) / float64(stats.RunCount) * 100
	}

	return stats, nil
}

// scheduleLoop 调度循环
func (s *DefaultScheduler) scheduleLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-s.pauseChan:
			<-s.resumeChan
			continue
		case <-ticker.C:
			s.checkAndRunTasks()
		}
	}
}

// checkAndRunTasks 检查并运行任务
func (s *DefaultScheduler) checkAndRunTasks() {
	s.mu.RLock()
	tasks := make([]Task, 0)
	for _, task := range s.tasks {
		if task.GetEnabled() {
			tasks = append(tasks, task)
		}
	}
	s.mu.RUnlock()

	now := time.Now()
	for _, task := range tasks {
		if task.GetNextRunAt() != nil && now.After(*task.GetNextRunAt()) {
			go s.executeTask(task)
		}
	}
}

// executeTask 执行任务
func (s *DefaultScheduler) executeTask(task Task) {
	ctx, cancel := context.WithTimeout(s.ctx, task.GetTimeout())
	defer cancel()

	// 执行任务
	err := task.GetHandler().Handle(ctx)

	// 更新任务状态
	s.mu.Lock()
	task.MarkAsRun()
	if err != nil {
		task.MarkAsFailed(err)
		task.IncrementFailedCount()
		s.stats.TotalFailed++
	} else {
		task.IncrementRunCount()
	}
	s.stats.TotalRuns++
	s.stats.LastRunAt = time.Now()

	// 更新下次运行时间
	task.UpdateNextRun()

	// 保存到存储
	s.store.Save(task)
	s.mu.Unlock()
}
