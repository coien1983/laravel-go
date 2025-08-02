package cache

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// MemoryItem 内存缓存项
type MemoryItem struct {
	Value      interface{}
	Expiration time.Time
	// 添加原子计数器用于引用计数
	refCount int32
}

// IsExpired 检查是否过期
func (item *MemoryItem) IsExpired() bool {
	return !item.Expiration.IsZero() && time.Now().After(item.Expiration)
}

// IncrementRef 增加引用计数
func (item *MemoryItem) IncrementRef() {
	atomic.AddInt32(&item.refCount, 1)
}

// DecrementRef 减少引用计数
func (item *MemoryItem) DecrementRef() {
	atomic.AddInt32(&item.refCount, -1)
}

// GetRefCount 获取引用计数
func (item *MemoryItem) GetRefCount() int32 {
	return atomic.LoadInt32(&item.refCount)
}

// MemoryStore 内存缓存存储
type MemoryStore struct {
	items  map[string]*MemoryItem
	mutex  sync.RWMutex
	prefix string
	// 添加统计信息
	stats struct {
		hits    int64
		misses  int64
		sets    int64
		deletes int64
	}
	// 添加清理控制
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

// NewMemoryStore 创建新的内存缓存存储
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		items:         make(map[string]*MemoryItem),
		prefix:        "",
		cleanupTicker: time.NewTicker(time.Minute),
		stopChan:      make(chan struct{}),
	}

	// 启动清理过期项的goroutine
	go store.cleanupExpired()

	return store
}

// cleanupExpired 清理过期的缓存项
func (store *MemoryStore) cleanupExpired() {
	for {
		select {
		case <-store.cleanupTicker.C:
			store.cleanupExpiredItems()
		case <-store.stopChan:
			return
		}
	}
}

// cleanupExpiredItems 清理过期项的具体实现
func (store *MemoryStore) cleanupExpiredItems() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	var expiredKeys []string

	for key, item := range store.items {
		if item.IsExpired() {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// 批量删除过期项
	for _, key := range expiredKeys {
		delete(store.items, key)
		atomic.AddInt64(&store.stats.deletes, 1)
	}
}

// Get 获取缓存值
func (store *MemoryStore) Get(key string) (interface{}, error) {
	fullKey := store.prefix + key

	store.mutex.RLock()
	item, exists := store.items[fullKey]
	store.mutex.RUnlock()

	if !exists {
		atomic.AddInt64(&store.stats.misses, 1)
		return nil, fmt.Errorf("cache key not found: %s", key)
	}

	if item.IsExpired() {
		// 异步删除过期项
		go store.deleteExpiredItem(fullKey)
		atomic.AddInt64(&store.stats.misses, 1)
		return nil, fmt.Errorf("cache key expired: %s", key)
	}

	// 增加引用计数
	item.IncrementRef()
	atomic.AddInt64(&store.stats.hits, 1)

	return item.Value, nil
}

// deleteExpiredItem 异步删除过期项
func (store *MemoryStore) deleteExpiredItem(key string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if item, exists := store.items[key]; exists && item.IsExpired() {
		delete(store.items, key)
		atomic.AddInt64(&store.stats.deletes, 1)
	}
}

// GetString 获取字符串缓存值
func (store *MemoryStore) GetString(key string) (string, error) {
	value, err := store.Get(key)
	if err != nil {
		return "", err
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// GetInt 获取整数缓存值
func (store *MemoryStore) GetInt(key string) (int, error) {
	value, err := store.Get(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %v to int", value)
	}
}

// GetFloat 获取浮点数缓存值
func (store *MemoryStore) GetFloat(key string) (float64, error) {
	value, err := store.Get(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to float64", value)
	}
}

// GetBool 获取布尔值缓存值
func (store *MemoryStore) GetBool(key string) (bool, error) {
	value, err := store.Get(key)
	if err != nil {
		return false, err
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("cannot convert %v to bool", value)
	}
}

// GetBytes 获取字节数组缓存值
func (store *MemoryStore) GetBytes(key string) ([]byte, error) {
	value, err := store.Get(key)
	if err != nil {
		return nil, err
	}

	switch v := value.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		// 尝试JSON序列化
		return json.Marshal(v)
	}
}

// Set 设置缓存值
func (store *MemoryStore) Set(key string, value interface{}, ttl time.Duration) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	item := &MemoryItem{
		Value:      value,
		Expiration: expiration,
		refCount:   1,
	}

	store.items[store.prefix+key] = item
	atomic.AddInt64(&store.stats.sets, 1)

	return nil
}

// SetString 设置字符串缓存值
func (store *MemoryStore) SetString(key string, value string, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetInt 设置整数缓存值
func (store *MemoryStore) SetInt(key string, value int, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetFloat 设置浮点数缓存值
func (store *MemoryStore) SetFloat(key string, value float64, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetBool 设置布尔值缓存值
func (store *MemoryStore) SetBool(key string, value bool, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetBytes 设置字节数组缓存值
func (store *MemoryStore) SetBytes(key string, value []byte, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// Delete 删除缓存项
func (store *MemoryStore) Delete(key string) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	fullKey := store.prefix + key
	if item, exists := store.items[fullKey]; exists {
		// 减少引用计数
		item.DecrementRef()
		delete(store.items, fullKey)
		atomic.AddInt64(&store.stats.deletes, 1)
	}

	return nil
}

// DeleteMultiple 批量删除缓存项
func (store *MemoryStore) DeleteMultiple(keys []string) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for _, key := range keys {
		fullKey := store.prefix + key
		if item, exists := store.items[fullKey]; exists {
			item.DecrementRef()
			delete(store.items, fullKey)
			atomic.AddInt64(&store.stats.deletes, 1)
		}
	}

	return nil
}

// Clear 清空缓存
func (store *MemoryStore) Clear() error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	// 减少所有项的引用计数
	for _, item := range store.items {
		item.DecrementRef()
	}

	store.items = make(map[string]*MemoryItem)
	return nil
}

// Has 检查键是否存在
func (store *MemoryStore) Has(key string) bool {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	item, exists := store.items[store.prefix+key]
	if !exists {
		return false
	}

	return !item.IsExpired()
}

// Missing 检查缓存是否不存在
func (store *MemoryStore) Missing(key string) bool {
	return !store.Has(key)
}

// Increment 递增缓存值
func (store *MemoryStore) Increment(key string, value int) (int, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	// 直接获取当前值，避免锁嵌套
	item, exists := store.items[store.prefix+key]
	var current int
	if exists && !item.IsExpired() {
		switch v := item.Value.(type) {
		case int:
			current = v
		case int64:
			current = int(v)
		case float64:
			current = int(v)
		case string:
			if parsed, err := strconv.Atoi(v); err == nil {
				current = parsed
			}
		}
	}

	newValue := current + value
	store.items[store.prefix+key] = &MemoryItem{
		Value: newValue,
	}

	return newValue, nil
}

// Decrement 递减缓存值
func (store *MemoryStore) Decrement(key string, value int) (int, error) {
	return store.Increment(key, -value)
}

// Remember 记住缓存值
func (store *MemoryStore) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// 先尝试获取缓存
	if value, err := store.Get(key); err == nil {
		return value, nil
	}

	// 缓存不存在，执行回调函数
	value, err := callback()
	if err != nil {
		return nil, err
	}

	// 设置缓存
	if err := store.Set(key, value, ttl); err != nil {
		return nil, err
	}

	return value, nil
}

// RememberForever 永久记住缓存值
func (store *MemoryStore) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return store.Remember(key, 0, callback)
}

// Tags 获取标签管理器
func (store *MemoryStore) Tags(names ...string) TaggedStore {
	return NewMemoryTaggedStore(store, names...)
}

// Flush 刷新缓存
func (store *MemoryStore) Flush() error {
	return store.Clear()
}

// GetPrefix 获取缓存键前缀
func (store *MemoryStore) GetPrefix() string {
	return store.prefix
}

// SetPrefix 设置缓存键前缀
func (store *MemoryStore) SetPrefix(prefix string) {
	store.prefix = prefix
}

// GetStats 获取缓存统计信息
func (store *MemoryStore) GetStats() map[string]int64 {
	return map[string]int64{
		"hits":    atomic.LoadInt64(&store.stats.hits),
		"misses":  atomic.LoadInt64(&store.stats.misses),
		"sets":    atomic.LoadInt64(&store.stats.sets),
		"deletes": atomic.LoadInt64(&store.stats.deletes),
		"items":   int64(len(store.items)),
	}
}

// Close 关闭缓存存储
func (store *MemoryStore) Close() error {
	close(store.stopChan)
	store.cleanupTicker.Stop()
	return store.Clear()
}
