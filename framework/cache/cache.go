package cache

import (
	"time"
)

// Store 缓存存储接口
type Store interface {
	// Get 获取缓存值
	Get(key string) (interface{}, error)
	// GetString 获取字符串缓存值
	GetString(key string) (string, error)
	// GetInt 获取整数缓存值
	GetInt(key string) (int, error)
	// GetFloat 获取浮点数缓存值
	GetFloat(key string) (float64, error)
	// GetBool 获取布尔值缓存值
	GetBool(key string) (bool, error)
	// GetBytes 获取字节数组缓存值
	GetBytes(key string) ([]byte, error)

	// Set 设置缓存值
	Set(key string, value interface{}, ttl time.Duration) error
	// SetString 设置字符串缓存值
	SetString(key string, value string, ttl time.Duration) error
	// SetInt 设置整数缓存值
	SetInt(key string, value int, ttl time.Duration) error
	// SetFloat 设置浮点数缓存值
	SetFloat(key string, value float64, ttl time.Duration) error
	// SetBool 设置布尔值缓存值
	SetBool(key string, value bool, ttl time.Duration) error
	// SetBytes 设置字节数组缓存值
	SetBytes(key string, value []byte, ttl time.Duration) error

	// Delete 删除缓存
	Delete(key string) error
	// DeleteMultiple 批量删除缓存
	DeleteMultiple(keys []string) error
	// Clear 清空所有缓存
	Clear() error

	// Has 检查缓存是否存在
	Has(key string) bool
	// Missing 检查缓存是否不存在
	Missing(key string) bool

	// Increment 递增缓存值
	Increment(key string, value int) (int, error)
	// Decrement 递减缓存值
	Decrement(key string, value int) (int, error)

	// Remember 记住缓存值（如果不存在则设置）
	Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error)
	// RememberForever 永久记住缓存值
	RememberForever(key string, callback func() (interface{}, error)) (interface{}, error)

	// Tags 获取标签管理器
	Tags(names ...string) TaggedStore
	// Flush 刷新缓存
	Flush() error

	// GetPrefix 获取缓存键前缀
	GetPrefix() string
	// SetPrefix 设置缓存键前缀
	SetPrefix(prefix string)
}

// TaggedStore 带标签的缓存存储接口
type TaggedStore interface {
	Store
	// Flush 刷新标签下的所有缓存
	Flush() error
	// GetTags 获取标签列表
	GetTags() []string
	// AddTags 添加标签
	AddTags(names ...string) TaggedStore
}

// Manager 缓存管理器
type Manager struct {
	stores       map[string]Store
	defaultStore string
	config       map[string]interface{}
}

// NewManager 创建新的缓存管理器
func NewManager() *Manager {
	return &Manager{
		stores:       make(map[string]Store),
		defaultStore: "memory",
		config:       make(map[string]interface{}),
	}
}

// Store 获取缓存存储
func (m *Manager) Store(name string) Store {
	if store, exists := m.stores[name]; exists {
		return store
	}
	return m.stores[m.defaultStore]
}

// DefaultStore 获取默认缓存存储
func (m *Manager) DefaultStore() Store {
	return m.Store(m.defaultStore)
}

// Extend 扩展缓存存储
func (m *Manager) Extend(name string, store Store) {
	m.stores[name] = store
}

// SetDefaultStore 设置默认缓存存储
func (m *Manager) SetDefaultStore(name string) {
	m.defaultStore = name
}

// Get 获取缓存值
func (m *Manager) Get(key string) (interface{}, error) {
	return m.DefaultStore().Get(key)
}

// GetString 获取字符串缓存值
func (m *Manager) GetString(key string) (string, error) {
	return m.DefaultStore().GetString(key)
}

// GetInt 获取整数缓存值
func (m *Manager) GetInt(key string) (int, error) {
	return m.DefaultStore().GetInt(key)
}

// GetFloat 获取浮点数缓存值
func (m *Manager) GetFloat(key string) (float64, error) {
	return m.DefaultStore().GetFloat(key)
}

// GetBool 获取布尔值缓存值
func (m *Manager) GetBool(key string) (bool, error) {
	return m.DefaultStore().GetBool(key)
}

// GetBytes 获取字节数组缓存值
func (m *Manager) GetBytes(key string) ([]byte, error) {
	return m.DefaultStore().GetBytes(key)
}

// Set 设置缓存值
func (m *Manager) Set(key string, value interface{}, ttl time.Duration) error {
	return m.DefaultStore().Set(key, value, ttl)
}

// SetString 设置字符串缓存值
func (m *Manager) SetString(key string, value string, ttl time.Duration) error {
	return m.DefaultStore().SetString(key, value, ttl)
}

// SetInt 设置整数缓存值
func (m *Manager) SetInt(key string, value int, ttl time.Duration) error {
	return m.DefaultStore().SetInt(key, value, ttl)
}

// SetFloat 设置浮点数缓存值
func (m *Manager) SetFloat(key string, value float64, ttl time.Duration) error {
	return m.DefaultStore().SetFloat(key, value, ttl)
}

// SetBool 设置布尔值缓存值
func (m *Manager) SetBool(key string, value bool, ttl time.Duration) error {
	return m.DefaultStore().SetBool(key, value, ttl)
}

// SetBytes 设置字节数组缓存值
func (m *Manager) SetBytes(key string, value []byte, ttl time.Duration) error {
	return m.DefaultStore().SetBytes(key, value, ttl)
}

// Delete 删除缓存
func (m *Manager) Delete(key string) error {
	return m.DefaultStore().Delete(key)
}

// DeleteMultiple 批量删除缓存
func (m *Manager) DeleteMultiple(keys []string) error {
	return m.DefaultStore().DeleteMultiple(keys)
}

// Clear 清空所有缓存
func (m *Manager) Clear() error {
	return m.DefaultStore().Clear()
}

// Has 检查缓存是否存在
func (m *Manager) Has(key string) bool {
	return m.DefaultStore().Has(key)
}

// Missing 检查缓存是否不存在
func (m *Manager) Missing(key string) bool {
	return m.DefaultStore().Missing(key)
}

// Increment 递增缓存值
func (m *Manager) Increment(key string, value int) (int, error) {
	return m.DefaultStore().Increment(key, value)
}

// Decrement 递减缓存值
func (m *Manager) Decrement(key string, value int) (int, error) {
	return m.DefaultStore().Decrement(key, value)
}

// Remember 记住缓存值
func (m *Manager) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	return m.DefaultStore().Remember(key, ttl, callback)
}

// RememberForever 永久记住缓存值
func (m *Manager) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return m.DefaultStore().RememberForever(key, callback)
}

// Tags 获取标签管理器
func (m *Manager) Tags(names ...string) TaggedStore {
	return m.DefaultStore().Tags(names...)
}

// Flush 刷新缓存
func (m *Manager) Flush() error {
	return m.DefaultStore().Flush()
}

// Cache 全局缓存实例
var Cache *Manager

// Init 初始化全局缓存
func Init() {
	Cache = NewManager()

	// 注册默认的内存缓存驱动
	memoryStore := NewMemoryStore()
	Cache.Extend("memory", memoryStore)
	Cache.SetDefaultStore("memory")
}

// Get 全局获取缓存值
func Get(key string) (interface{}, error) {
	return Cache.Get(key)
}

// Set 全局设置缓存值
func Set(key string, value interface{}, ttl time.Duration) error {
	return Cache.Set(key, value, ttl)
}

// Delete 全局删除缓存
func Delete(key string) error {
	return Cache.Delete(key)
}

// Has 全局检查缓存是否存在
func Has(key string) bool {
	return Cache.Has(key)
}

// Remember 全局记住缓存值
func Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	return Cache.Remember(key, ttl, callback)
}
