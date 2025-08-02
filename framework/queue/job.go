package queue

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BaseJob 基础任务实现
type BaseJob struct {
	ID          string            `json:"id"`
	Payload     []byte            `json:"payload"`
	Queue       string            `json:"queue"`
	Attempts    int               `json:"attempts"`
	MaxAttempts int               `json:"max_attempts"`
	Delay       time.Duration     `json:"delay"`
	Timeout     time.Duration     `json:"timeout"`
	Priority    int               `json:"priority"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	ReservedAt  *time.Time        `json:"reserved_at"`
	AvailableAt time.Time         `json:"available_at"`
	CompletedAt *time.Time        `json:"completed_at"`
	FailedAt    *time.Time        `json:"failed_at"`
	Error       string            `json:"error"`
}

// NewJob 创建新任务
func NewJob(payload []byte, queue string) *BaseJob {
	now := time.Now()
	return &BaseJob{
		ID:          uuid.New().String(),
		Payload:     payload,
		Queue:       queue,
		MaxAttempts: 3,
		Delay:       0,
		Timeout:     30 * time.Second,
		Priority:    0,
		Tags:        make(map[string]string),
		CreatedAt:   now,
		AvailableAt: now,
	}
}

// GetID 获取任务ID
func (j *BaseJob) GetID() string {
	return j.ID
}

// GetPayload 获取任务载荷
func (j *BaseJob) GetPayload() []byte {
	return j.Payload
}

// GetQueue 获取队列名称
func (j *BaseJob) GetQueue() string {
	return j.Queue
}

// GetAttempts 获取尝试次数
func (j *BaseJob) GetAttempts() int {
	return j.Attempts
}

// GetMaxAttempts 获取最大尝试次数
func (j *BaseJob) GetMaxAttempts() int {
	return j.MaxAttempts
}

// GetDelay 获取延迟时间
func (j *BaseJob) GetDelay() time.Duration {
	return j.Delay
}

// GetTimeout 获取超时时间
func (j *BaseJob) GetTimeout() time.Duration {
	return j.Timeout
}

// GetPriority 获取优先级
func (j *BaseJob) GetPriority() int {
	return j.Priority
}

// GetTags 获取标签
func (j *BaseJob) GetTags() map[string]string {
	return j.Tags
}

// GetCreatedAt 获取创建时间
func (j *BaseJob) GetCreatedAt() time.Time {
	return j.CreatedAt
}

// GetReservedAt 获取保留时间
func (j *BaseJob) GetReservedAt() *time.Time {
	return j.ReservedAt
}

// GetAvailableAt 获取可用时间
func (j *BaseJob) GetAvailableAt() time.Time {
	return j.AvailableAt
}

// MarkAsReserved 标记为已保留
func (j *BaseJob) MarkAsReserved() {
	now := time.Now()
	j.ReservedAt = &now
}

// MarkAsCompleted 标记为已完成
func (j *BaseJob) MarkAsCompleted() {
	now := time.Now()
	j.CompletedAt = &now
}

// MarkAsFailed 标记为失败
func (j *BaseJob) MarkAsFailed(err error) {
	now := time.Now()
	j.FailedAt = &now
	if err != nil {
		j.Error = err.Error()
	}
}

// IncrementAttempts 增加尝试次数
func (j *BaseJob) IncrementAttempts() {
	j.Attempts++
}

// SetMaxAttempts 设置最大尝试次数
func (j *BaseJob) SetMaxAttempts(maxAttempts int) {
	j.MaxAttempts = maxAttempts
}

// SetDelay 设置延迟时间
func (j *BaseJob) SetDelay(delay time.Duration) {
	j.Delay = delay
	j.AvailableAt = time.Now().Add(delay)
}

// SetTimeout 设置超时时间
func (j *BaseJob) SetTimeout(timeout time.Duration) {
	j.Timeout = timeout
}

// SetPriority 设置优先级
func (j *BaseJob) SetPriority(priority int) {
	j.Priority = priority
}

// AddTag 添加标签
func (j *BaseJob) AddTag(key, value string) {
	if j.Tags == nil {
		j.Tags = make(map[string]string)
	}
	j.Tags[key] = value
}

// RemoveTag 移除标签
func (j *BaseJob) RemoveTag(key string) {
	if j.Tags != nil {
		delete(j.Tags, key)
	}
}

// Serialize 序列化任务
func (j *BaseJob) Serialize() ([]byte, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize job: %w", err)
	}
	return data, nil
}

// Deserialize 反序列化任务
func (j *BaseJob) Deserialize(data []byte) error {
	err := json.Unmarshal(data, j)
	if err != nil {
		return fmt.Errorf("failed to deserialize job: %w", err)
	}
	return nil
}

// IsExpired 检查是否过期
func (j *BaseJob) IsExpired() bool {
	if j.ReservedAt == nil {
		return false
	}
	return time.Now().After(j.ReservedAt.Add(j.Timeout))
}

// IsAvailable 检查是否可用
func (j *BaseJob) IsAvailable() bool {
	return time.Now().After(j.AvailableAt)
}

// CanRetry 检查是否可以重试
func (j *BaseJob) CanRetry() bool {
	return j.Attempts < j.MaxAttempts
}

// IsCompleted 检查是否已完成
func (j *BaseJob) IsCompleted() bool {
	return j.CompletedAt != nil
}

// IsFailed 检查是否失败
func (j *BaseJob) IsFailed() bool {
	return j.FailedAt != nil
}

// IsReserved 检查是否已保留
func (j *BaseJob) IsReserved() bool {
	return j.ReservedAt != nil
} 