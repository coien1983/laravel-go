package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BaseEvent 基础事件实现
type BaseEvent struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Payload    interface{}            `json:"payload"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
	Propagated bool                   `json:"propagated"`
}

// NewEvent 创建新事件
func NewEvent(name string, payload interface{}) *BaseEvent {
	return &BaseEvent{
		ID:         uuid.New().String(),
		Name:       name,
		Payload:    payload,
		Timestamp:  time.Now(),
		Data:       make(map[string]interface{}),
		Propagated: false,
	}
}

// GetID 获取事件ID
func (e *BaseEvent) GetID() string {
	return e.ID
}

// GetName 获取事件名称
func (e *BaseEvent) GetName() string {
	return e.Name
}

// GetPayload 获取事件载荷
func (e *BaseEvent) GetPayload() interface{} {
	return e.Payload
}

// GetTimestamp 获取事件时间戳
func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetData 获取事件数据
func (e *BaseEvent) GetData() map[string]interface{} {
	return e.Data
}

// SetData 设置事件数据
func (e *BaseEvent) SetData(key string, value interface{}) {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	e.Data[key] = value
}

// GetDataByKey 根据键获取数据
func (e *BaseEvent) GetDataByKey(key string) interface{} {
	if e.Data == nil {
		return nil
	}
	return e.Data[key]
}

// IsPropagated 检查是否已传播
func (e *BaseEvent) IsPropagated() bool {
	return e.Propagated
}

// SetPropagated 设置传播状态
func (e *BaseEvent) SetPropagated(propagated bool) {
	e.Propagated = propagated
}

// Serialize 序列化事件
func (e *BaseEvent) Serialize() ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize event: %w", err)
	}
	return data, nil
}

// Deserialize 反序列化事件
func (e *BaseEvent) Deserialize(data []byte) error {
	err := json.Unmarshal(data, e)
	if err != nil {
		return fmt.Errorf("failed to deserialize event: %w", err)
	}
	return nil
}

// String 字符串表示
func (e *BaseEvent) String() string {
	return fmt.Sprintf("Event{ID: %s, Name: %s, Timestamp: %s}", e.ID, e.Name, e.Timestamp.Format(time.RFC3339))
}
