package cache

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcachedStore Memcached缓存存储
type MemcachedStore struct {
	client *memcache.Client
	prefix string
}

// NewMemcachedStore 创建新的Memcached缓存存储
func NewMemcachedStore(servers ...string) *MemcachedStore {
	client := memcache.New(servers...)
	return &MemcachedStore{
		client: client,
		prefix: "",
	}
}

// NewMemcachedStoreWithConfig 使用配置创建Memcached缓存存储
func NewMemcachedStoreWithConfig(config map[string]interface{}) *MemcachedStore {
	host := "127.0.0.1"
	port := "11211"

	if h, ok := config["host"].(string); ok && h != "" {
		host = h
	}
	if p, ok := config["port"].(string); ok && p != "" {
		port = p
	}

	server := fmt.Sprintf("%s:%s", host, port)
	client := memcache.New(server)

	return &MemcachedStore{
		client: client,
		prefix: "",
	}
}

// Get 获取缓存值
func (store *MemcachedStore) Get(key string) (interface{}, error) {
	item, err := store.client.Get(store.prefix + key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, fmt.Errorf("cache key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	var result interface{}
	if err := json.Unmarshal(item.Value, &result); err != nil {
		return string(item.Value), nil
	}

	return result, nil
}

// GetString 获取字符串缓存值
func (store *MemcachedStore) GetString(key string) (string, error) {
	item, err := store.client.Get(store.prefix + key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return "", fmt.Errorf("cache key not found: %s", key)
		}
		return "", fmt.Errorf("failed to get cache: %w", err)
	}

	return string(item.Value), nil
}

// GetInt 获取整数缓存值
func (store *MemcachedStore) GetInt(key string) (int, error) {
	item, err := store.client.Get(store.prefix + key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return 0, fmt.Errorf("cache key not found: %s", key)
		}
		return 0, fmt.Errorf("failed to get cache: %w", err)
	}

	value, err := strconv.Atoi(string(item.Value))
	if err != nil {
		return 0, fmt.Errorf("failed to convert to int: %w", err)
	}

	return value, nil
}

// GetFloat 获取浮点数缓存值
func (store *MemcachedStore) GetFloat(key string) (float64, error) {
	item, err := store.client.Get(store.prefix + key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return 0, fmt.Errorf("cache key not found: %s", key)
		}
		return 0, fmt.Errorf("failed to get cache: %w", err)
	}

	value, err := strconv.ParseFloat(string(item.Value), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert to float: %w", err)
	}

	return value, nil
}

// GetBool 获取布尔值缓存值
func (store *MemcachedStore) GetBool(key string) (bool, error) {
	item, err := store.client.Get(store.prefix + key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return false, fmt.Errorf("cache key not found: %s", key)
		}
		return false, fmt.Errorf("failed to get cache: %w", err)
	}

	value, err := strconv.ParseBool(string(item.Value))
	if err != nil {
		return false, fmt.Errorf("failed to convert to bool: %w", err)
	}

	return value, nil
}

// GetBytes 获取字节数组缓存值
func (store *MemcachedStore) GetBytes(key string) ([]byte, error) {
	item, err := store.client.Get(store.prefix + key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, fmt.Errorf("cache key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	return item.Value, nil
}

// Set 设置缓存值
func (store *MemcachedStore) Set(key string, value interface{}, ttl time.Duration) error {
	var data []byte
	var err error

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case int, int8, int16, int32, int64:
		data = []byte(fmt.Sprintf("%d", v))
	case uint, uint8, uint16, uint32, uint64:
		data = []byte(fmt.Sprintf("%d", v))
	case float32, float64:
		data = []byte(fmt.Sprintf("%f", v))
	case bool:
		data = []byte(strconv.FormatBool(v))
	case []byte:
		data = v
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	item := &memcache.Item{
		Key:        store.prefix + key,
		Value:      data,
		Expiration: int32(ttl.Seconds()),
	}

	return store.client.Set(item)
}

// SetString 设置字符串缓存值
func (store *MemcachedStore) SetString(key string, value string, ttl time.Duration) error {
	item := &memcache.Item{
		Key:        store.prefix + key,
		Value:      []byte(value),
		Expiration: int32(ttl.Seconds()),
	}

	return store.client.Set(item)
}

// SetInt 设置整数缓存值
func (store *MemcachedStore) SetInt(key string, value int, ttl time.Duration) error {
	item := &memcache.Item{
		Key:        store.prefix + key,
		Value:      []byte(strconv.Itoa(value)),
		Expiration: int32(ttl.Seconds()),
	}

	return store.client.Set(item)
}

// SetFloat 设置浮点数缓存值
func (store *MemcachedStore) SetFloat(key string, value float64, ttl time.Duration) error {
	item := &memcache.Item{
		Key:        store.prefix + key,
		Value:      []byte(strconv.FormatFloat(value, 'f', -1, 64)),
		Expiration: int32(ttl.Seconds()),
	}

	return store.client.Set(item)
}

// SetBool 设置布尔值缓存值
func (store *MemcachedStore) SetBool(key string, value bool, ttl time.Duration) error {
	item := &memcache.Item{
		Key:        store.prefix + key,
		Value:      []byte(strconv.FormatBool(value)),
		Expiration: int32(ttl.Seconds()),
	}

	return store.client.Set(item)
}

// SetBytes 设置字节数组缓存值
func (store *MemcachedStore) SetBytes(key string, value []byte, ttl time.Duration) error {
	item := &memcache.Item{
		Key:        store.prefix + key,
		Value:      value,
		Expiration: int32(ttl.Seconds()),
	}

	return store.client.Set(item)
}

// Delete 删除缓存
func (store *MemcachedStore) Delete(key string) error {
	err := store.client.Delete(store.prefix + key)
	if err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

// DeleteMultiple 批量删除缓存
func (store *MemcachedStore) DeleteMultiple(keys []string) error {
	for _, key := range keys {
		if err := store.Delete(key); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", key, err)
		}
	}
	return nil
}

// Clear 清空所有缓存
func (store *MemcachedStore) Clear() error {
	return fmt.Errorf("memcached does not support clearing all cache, use Delete() for specific keys")
}

// Has 检查缓存是否存在
func (store *MemcachedStore) Has(key string) bool {
	_, err := store.client.Get(store.prefix + key)
	return err == nil
}

// Missing 检查缓存是否不存在
func (store *MemcachedStore) Missing(key string) bool {
	return !store.Has(key)
}

// Increment 递增缓存值
func (store *MemcachedStore) Increment(key string, value int) (int, error) {
	newValue, err := store.client.Increment(store.prefix+key, uint64(value))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			if err := store.SetInt(key, 0, 0); err != nil {
				return 0, fmt.Errorf("failed to initialize key for increment: %w", err)
			}
			newValue, err = store.client.Increment(store.prefix+key, uint64(value))
			if err != nil {
				return 0, fmt.Errorf("failed to increment after initialization: %w", err)
			}
		} else {
			return 0, fmt.Errorf("failed to increment cache: %w", err)
		}
	}

	return int(newValue), nil
}

// Decrement 递减缓存值
func (store *MemcachedStore) Decrement(key string, value int) (int, error) {
	newValue, err := store.client.Decrement(store.prefix+key, uint64(value))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			if err := store.SetInt(key, 0, 0); err != nil {
				return 0, fmt.Errorf("failed to initialize key for decrement: %w", err)
			}
			newValue, err = store.client.Decrement(store.prefix+key, uint64(value))
			if err != nil {
				return 0, fmt.Errorf("failed to decrement after initialization: %w", err)
			}
		} else {
			return 0, fmt.Errorf("failed to decrement cache: %w", err)
		}
	}

	return int(newValue), nil
}

// Remember 记住缓存值（如果不存在则设置）
func (store *MemcachedStore) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	if value, err := store.Get(key); err == nil {
		return value, nil
	}

	value, err := callback()
	if err != nil {
		return nil, fmt.Errorf("callback failed: %w", err)
	}

	if err := store.Set(key, value, ttl); err != nil {
		return nil, fmt.Errorf("failed to set cache: %w", err)
	}

	return value, nil
}

// RememberForever 永久记住缓存值
func (store *MemcachedStore) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return store.Remember(key, 0, callback)
}

// Tags 获取标签管理器
func (store *MemcachedStore) Tags(names ...string) TaggedStore {
	return &MemcachedTaggedStore{
		store: store,
		tags:  names,
	}
}

// Flush 刷新缓存
func (store *MemcachedStore) Flush() error {
	return fmt.Errorf("memcached does not support flushing all cache")
}

// GetPrefix 获取缓存键前缀
func (store *MemcachedStore) GetPrefix() string {
	return store.prefix
}

// SetPrefix 设置缓存键前缀
func (store *MemcachedStore) SetPrefix(prefix string) {
	store.prefix = prefix
}

// GetStats 获取缓存统计信息
func (store *MemcachedStore) GetStats() (map[string]interface{}, error) {
	stats, err := store.client.Stats()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	result := make(map[string]interface{})
	for server, statMap := range stats {
		serverStats := make(map[string]interface{})
		for key, value := range statMap {
			serverStats[key] = value
		}
		result[server] = serverStats
	}

	return result, nil
}

// GetTTL 获取缓存项的TTL
func (store *MemcachedStore) GetTTL(key string) (time.Duration, error) {
	_, err := store.client.Get(store.prefix + key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return 0, fmt.Errorf("cache key not found: %s", key)
		}
		return 0, fmt.Errorf("failed to get cache: %w", err)
	}

	return 0, nil
}

// SetTTL 设置缓存项的TTL
func (store *MemcachedStore) SetTTL(key string, ttl time.Duration) error {
	item, err := store.client.Get(store.prefix + key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return fmt.Errorf("cache key not found: %s", key)
		}
		return fmt.Errorf("failed to get cache: %w", err)
	}

	newItem := &memcache.Item{
		Key:        item.Key,
		Value:      item.Value,
		Expiration: int32(ttl.Seconds()),
	}

	return store.client.Set(newItem)
}

// MemcachedTaggedStore Memcached标签存储包装器
type MemcachedTaggedStore struct {
	store *MemcachedStore
	tags  []string
}

// 实现TaggedStore接口的所有方法
func (ts *MemcachedTaggedStore) Get(key string) (interface{}, error) {
	return ts.store.Get(key)
}

func (ts *MemcachedTaggedStore) GetString(key string) (string, error) {
	return ts.store.GetString(key)
}

func (ts *MemcachedTaggedStore) GetInt(key string) (int, error) {
	return ts.store.GetInt(key)
}

func (ts *MemcachedTaggedStore) GetFloat(key string) (float64, error) {
	return ts.store.GetFloat(key)
}

func (ts *MemcachedTaggedStore) GetBool(key string) (bool, error) {
	return ts.store.GetBool(key)
}

func (ts *MemcachedTaggedStore) GetBytes(key string) ([]byte, error) {
	return ts.store.GetBytes(key)
}

func (ts *MemcachedTaggedStore) Set(key string, value interface{}, ttl time.Duration) error {
	return ts.store.Set(key, value, ttl)
}

func (ts *MemcachedTaggedStore) SetString(key string, value string, ttl time.Duration) error {
	return ts.store.SetString(key, value, ttl)
}

func (ts *MemcachedTaggedStore) SetInt(key string, value int, ttl time.Duration) error {
	return ts.store.SetInt(key, value, ttl)
}

func (ts *MemcachedTaggedStore) SetFloat(key string, value float64, ttl time.Duration) error {
	return ts.store.SetFloat(key, value, ttl)
}

func (ts *MemcachedTaggedStore) SetBool(key string, value bool, ttl time.Duration) error {
	return ts.store.SetBool(key, value, ttl)
}

func (ts *MemcachedTaggedStore) SetBytes(key string, value []byte, ttl time.Duration) error {
	return ts.store.SetBytes(key, value, ttl)
}

func (ts *MemcachedTaggedStore) Delete(key string) error {
	return ts.store.Delete(key)
}

func (ts *MemcachedTaggedStore) DeleteMultiple(keys []string) error {
	return ts.store.DeleteMultiple(keys)
}

func (ts *MemcachedTaggedStore) Clear() error {
	return ts.store.Clear()
}

func (ts *MemcachedTaggedStore) Has(key string) bool {
	return ts.store.Has(key)
}

func (ts *MemcachedTaggedStore) Missing(key string) bool {
	return ts.store.Missing(key)
}

func (ts *MemcachedTaggedStore) Increment(key string, value int) (int, error) {
	return ts.store.Increment(key, value)
}

func (ts *MemcachedTaggedStore) Decrement(key string, value int) (int, error) {
	return ts.store.Decrement(key, value)
}

func (ts *MemcachedTaggedStore) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	return ts.store.Remember(key, ttl, callback)
}

func (ts *MemcachedTaggedStore) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return ts.store.RememberForever(key, callback)
}

func (ts *MemcachedTaggedStore) Tags(names ...string) TaggedStore {
	allTags := append(ts.tags, names...)
	return &MemcachedTaggedStore{
		store: ts.store,
		tags:  allTags,
	}
}

func (ts *MemcachedTaggedStore) Flush() error {
	return fmt.Errorf("memcached does not support tag-based flushing")
}

func (ts *MemcachedTaggedStore) GetPrefix() string {
	return ts.store.GetPrefix()
}

func (ts *MemcachedTaggedStore) SetPrefix(prefix string) {
	ts.store.SetPrefix(prefix)
}

func (ts *MemcachedTaggedStore) GetTags() []string {
	return ts.tags
}

func (ts *MemcachedTaggedStore) AddTags(names ...string) TaggedStore {
	allTags := append(ts.tags, names...)
	return &MemcachedTaggedStore{
		store: ts.store,
		tags:  allTags,
	}
}
