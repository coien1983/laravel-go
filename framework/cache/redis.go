package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore Redis缓存存储
type RedisStore struct {
	client *redis.Client
	prefix string
}

// NewRedisStore 创建新的Redis缓存存储
func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
		prefix: "",
	}
}

// Get 获取缓存值
func (store *RedisStore) Get(key string) (interface{}, error) {
	ctx := context.Background()

	// 获取值
	value, err := store.client.Get(ctx, store.prefix+key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("cache key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	// 尝试解析JSON
	var result interface{}
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		// 如果不是JSON，返回原始字符串
		return value, nil
	}

	return result, nil
}

// GetString 获取字符串缓存值
func (store *RedisStore) GetString(key string) (string, error) {
	ctx := context.Background()

	value, err := store.client.Get(ctx, store.prefix+key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("cache key not found: %s", key)
		}
		return "", fmt.Errorf("failed to get cache: %w", err)
	}

	return value, nil
}

// GetInt 获取整数缓存值
func (store *RedisStore) GetInt(key string) (int, error) {
	ctx := context.Background()

	value, err := store.client.Get(ctx, store.prefix+key).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("cache key not found: %s", key)
		}
		return 0, fmt.Errorf("failed to get cache: %w", err)
	}

	return value, nil
}

// GetFloat 获取浮点数缓存值
func (store *RedisStore) GetFloat(key string) (float64, error) {
	ctx := context.Background()

	value, err := store.client.Get(ctx, store.prefix+key).Float64()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("cache key not found: %s", key)
		}
		return 0, fmt.Errorf("failed to get cache: %w", err)
	}

	return value, nil
}

// GetBool 获取布尔值缓存值
func (store *RedisStore) GetBool(key string) (bool, error) {
	ctx := context.Background()

	value, err := store.client.Get(ctx, store.prefix+key).Bool()
	if err != nil {
		if err == redis.Nil {
			return false, fmt.Errorf("cache key not found: %s", key)
		}
		return false, fmt.Errorf("failed to get cache: %w", err)
	}

	return value, nil
}

// GetBytes 获取字节数组缓存值
func (store *RedisStore) GetBytes(key string) ([]byte, error) {
	ctx := context.Background()

	value, err := store.client.Get(ctx, store.prefix+key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("cache key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	return value, nil
}

// Set 设置缓存值
func (store *RedisStore) Set(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()

	// 序列化为JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	err = store.client.Set(ctx, store.prefix+key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// SetString 设置字符串缓存值
func (store *RedisStore) SetString(key string, value string, ttl time.Duration) error {
	ctx := context.Background()

	err := store.client.Set(ctx, store.prefix+key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// SetInt 设置整数缓存值
func (store *RedisStore) SetInt(key string, value int, ttl time.Duration) error {
	ctx := context.Background()

	err := store.client.Set(ctx, store.prefix+key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// SetFloat 设置浮点数缓存值
func (store *RedisStore) SetFloat(key string, value float64, ttl time.Duration) error {
	ctx := context.Background()

	err := store.client.Set(ctx, store.prefix+key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// SetBool 设置布尔值缓存值
func (store *RedisStore) SetBool(key string, value bool, ttl time.Duration) error {
	ctx := context.Background()

	err := store.client.Set(ctx, store.prefix+key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// SetBytes 设置字节数组缓存值
func (store *RedisStore) SetBytes(key string, value []byte, ttl time.Duration) error {
	ctx := context.Background()

	err := store.client.Set(ctx, store.prefix+key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete 删除缓存
func (store *RedisStore) Delete(key string) error {
	ctx := context.Background()

	err := store.client.Del(ctx, store.prefix+key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	return nil
}

// DeleteMultiple 批量删除缓存
func (store *RedisStore) DeleteMultiple(keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	ctx := context.Background()

	// 构建key列表
	keyList := make([]string, len(keys))
	for i, key := range keys {
		keyList[i] = store.prefix + key
	}

	err := store.client.Del(ctx, keyList...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete multiple cache: %w", err)
	}

	return nil
}

// Clear 清空所有缓存
func (store *RedisStore) Clear() error {
	ctx := context.Background()

	// 获取所有匹配的key
	pattern := store.prefix + "*"
	keys, err := store.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	if len(keys) > 0 {
		err = store.client.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}
	}

	return nil
}

// Has 检查缓存是否存在
func (store *RedisStore) Has(key string) bool {
	ctx := context.Background()

	exists, err := store.client.Exists(ctx, store.prefix+key).Result()
	return err == nil && exists > 0
}

// Missing 检查缓存是否不存在
func (store *RedisStore) Missing(key string) bool {
	return !store.Has(key)
}

// Increment 递增缓存值
func (store *RedisStore) Increment(key string, value int) (int, error) {
	ctx := context.Background()

	result, err := store.client.IncrBy(ctx, store.prefix+key, int64(value)).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment cache: %w", err)
	}

	return int(result), nil
}

// Decrement 递减缓存值
func (store *RedisStore) Decrement(key string, value int) (int, error) {
	ctx := context.Background()

	result, err := store.client.DecrBy(ctx, store.prefix+key, int64(value)).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement cache: %w", err)
	}

	return int(result), nil
}

// Remember 记住缓存值
func (store *RedisStore) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
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
func (store *RedisStore) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return store.Remember(key, 0, callback)
}

// Tags 获取标签管理器
func (store *RedisStore) Tags(names ...string) TaggedStore {
	return NewMemoryTaggedStore(store, names...)
}

// Flush 刷新缓存
func (store *RedisStore) Flush() error {
	ctx := context.Background()

	err := store.client.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush cache: %w", err)
	}

	return nil
}

// GetPrefix 获取缓存键前缀
func (store *RedisStore) GetPrefix() string {
	return store.prefix
}

// SetPrefix 设置缓存键前缀
func (store *RedisStore) SetPrefix(prefix string) {
	store.prefix = prefix
}

// GetStats 获取缓存统计信息
func (store *RedisStore) GetStats() (map[string]interface{}, error) {
	ctx := context.Background()

	info, err := store.client.Info(ctx, "memory").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get redis info: %w", err)
	}

	// 获取数据库大小
	dbSize, err := store.client.DBSize(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get db size: %w", err)
	}

	return map[string]interface{}{
		"info":    info,
		"db_size": dbSize,
	}, nil
}

// GetTTL 获取缓存剩余时间
func (store *RedisStore) GetTTL(key string) (time.Duration, error) {
	ctx := context.Background()

	ttl, err := store.client.TTL(ctx, store.prefix+key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get ttl: %w", err)
	}

	return ttl, nil
}

// SetTTL 设置缓存过期时间
func (store *RedisStore) SetTTL(key string, ttl time.Duration) error {
	ctx := context.Background()

	err := store.client.Expire(ctx, store.prefix+key, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set ttl: %w", err)
	}

	return nil
}
