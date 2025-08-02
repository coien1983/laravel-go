package scheduler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DefaultTask 默认任务实现
type DefaultTask struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Schedule    string      `json:"schedule"`
	Handler     TaskHandler `json:"-"`
	Enabled     bool        `json:"enabled"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`

	// 执行信息
	LastRunAt   *time.Time `json:"last_run_at"`
	NextRunAt   *time.Time `json:"next_run_at"`
	RunCount    int64      `json:"run_count"`
	FailedCount int64      `json:"failed_count"`
	LastError   string     `json:"last_error"`

	// 配置
	Timeout    time.Duration     `json:"timeout"`
	RetryCount int               `json:"retry_count"`
	RetryDelay time.Duration     `json:"retry_delay"`
	MaxRetries int               `json:"max_retries"`
	Tags       map[string]string `json:"tags"`
}

// NewTask 创建新任务
func NewTask(name, description, schedule string, handler TaskHandler) *DefaultTask {
	now := time.Now()

	task := &DefaultTask{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Schedule:    schedule,
		Handler:     handler,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
		Timeout:     30 * time.Second,
		MaxRetries:  3,
		RetryDelay:  5 * time.Second,
		Tags:        make(map[string]string),
	}

	// 计算下次运行时间
	task.UpdateNextRun()

	return task
}

// GetID 获取任务ID
func (t *DefaultTask) GetID() string {
	return t.ID
}

// GetName 获取任务名称
func (t *DefaultTask) GetName() string {
	return t.Name
}

// GetDescription 获取任务描述
func (t *DefaultTask) GetDescription() string {
	return t.Description
}

// GetSchedule 获取调度表达式
func (t *DefaultTask) GetSchedule() string {
	return t.Schedule
}

// GetHandler 获取任务处理器
func (t *DefaultTask) GetHandler() TaskHandler {
	return t.Handler
}

// GetEnabled 获取是否启用
func (t *DefaultTask) GetEnabled() bool {
	return t.Enabled
}

// GetCreatedAt 获取创建时间
func (t *DefaultTask) GetCreatedAt() time.Time {
	return t.CreatedAt
}

// GetUpdatedAt 获取更新时间
func (t *DefaultTask) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}

// GetLastRunAt 获取上次运行时间
func (t *DefaultTask) GetLastRunAt() *time.Time {
	return t.LastRunAt
}

// GetNextRunAt 获取下次运行时间
func (t *DefaultTask) GetNextRunAt() *time.Time {
	return t.NextRunAt
}

// GetRunCount 获取运行次数
func (t *DefaultTask) GetRunCount() int64 {
	return t.RunCount
}

// GetFailedCount 获取失败次数
func (t *DefaultTask) GetFailedCount() int64 {
	return t.FailedCount
}

// GetLastError 获取最后错误
func (t *DefaultTask) GetLastError() string {
	return t.LastError
}

// GetTimeout 获取超时时间
func (t *DefaultTask) GetTimeout() time.Duration {
	return t.Timeout
}

// GetRetryCount 获取重试次数
func (t *DefaultTask) GetRetryCount() int {
	return t.RetryCount
}

// GetRetryDelay 获取重试延迟
func (t *DefaultTask) GetRetryDelay() time.Duration {
	return t.RetryDelay
}

// GetMaxRetries 获取最大重试次数
func (t *DefaultTask) GetMaxRetries() int {
	return t.MaxRetries
}

// GetTags 获取标签
func (t *DefaultTask) GetTags() map[string]string {
	return t.Tags
}

// Enable 启用任务
func (t *DefaultTask) Enable() {
	t.Enabled = true
	t.UpdatedAt = time.Now()
}

// Disable 禁用任务
func (t *DefaultTask) Disable() {
	t.Enabled = false
	t.UpdatedAt = time.Now()
}

// UpdateNextRun 更新下次运行时间
func (t *DefaultTask) UpdateNextRun() {
	nextRun, err := ParseSchedule(t.Schedule)
	if err == nil {
		t.NextRunAt = &nextRun
	}
}

// MarkAsRun 标记为已运行
func (t *DefaultTask) MarkAsRun() {
	now := time.Now()
	t.LastRunAt = &now
	t.UpdatedAt = now
	t.UpdateNextRun()
}

// MarkAsFailed 标记为失败
func (t *DefaultTask) MarkAsFailed(err error) {
	t.LastError = err.Error()
	t.UpdatedAt = time.Now()
}

// IncrementRunCount 增加运行次数
func (t *DefaultTask) IncrementRunCount() {
	t.RunCount++
}

// IncrementFailedCount 增加失败次数
func (t *DefaultTask) IncrementFailedCount() {
	t.FailedCount++
}

// SetTimeout 设置超时时间
func (t *DefaultTask) SetTimeout(timeout time.Duration) {
	t.Timeout = timeout
	t.UpdatedAt = time.Now()
}

// SetMaxRetries 设置最大重试次数
func (t *DefaultTask) SetMaxRetries(maxRetries int) {
	t.MaxRetries = maxRetries
	t.UpdatedAt = time.Now()
}

// SetRetryDelay 设置重试延迟
func (t *DefaultTask) SetRetryDelay(retryDelay time.Duration) {
	t.RetryDelay = retryDelay
	t.UpdatedAt = time.Now()
}

// AddTag 添加标签
func (t *DefaultTask) AddTag(key, value string) {
	if t.Tags == nil {
		t.Tags = make(map[string]string)
	}
	t.Tags[key] = value
	t.UpdatedAt = time.Now()
}

// RemoveTag 移除标签
func (t *DefaultTask) RemoveTag(key string) {
	if t.Tags != nil {
		delete(t.Tags, key)
		t.UpdatedAt = time.Now()
	}
}

// Serialize 序列化任务
func (t *DefaultTask) Serialize() ([]byte, error) {
	return json.Marshal(t)
}

// Deserialize 反序列化任务
func (t *DefaultTask) Deserialize(data []byte) error {
	return json.Unmarshal(data, t)
}

// Validate 验证任务
func (t *DefaultTask) Validate() error {
	if t.Name == "" {
		return ErrTaskNameRequired
	}

	if t.Handler == nil {
		return ErrTaskHandlerRequired
	}

	if t.Schedule == "" {
		return ErrInvalidSchedule
	}

	// 验证调度表达式
	_, err := ParseSchedule(t.Schedule)
	if err != nil {
		return fmt.Errorf("invalid schedule: %w", err)
	}

	return nil
}

// Clone 克隆任务
func (t *DefaultTask) Clone() *DefaultTask {
	clone := *t
	clone.ID = uuid.New().String()
	clone.CreatedAt = time.Now()
	clone.UpdatedAt = time.Now()
	clone.RunCount = 0
	clone.FailedCount = 0
	clone.LastError = ""
	clone.LastRunAt = nil
	clone.UpdateNextRun()

	// 克隆标签
	if t.Tags != nil {
		clone.Tags = make(map[string]string)
		for k, v := range t.Tags {
			clone.Tags[k] = v
		}
	}

	return &clone
}

// String 字符串表示
func (t *DefaultTask) String() string {
	return fmt.Sprintf("Task{ID: %s, Name: %s, Schedule: %s, Enabled: %t}",
		t.ID, t.Name, t.Schedule, t.Enabled)
}
