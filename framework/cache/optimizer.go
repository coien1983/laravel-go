package cache

import (
	"fmt"
	"sync"
	"time"
)

// CacheStats 缓存统计信息
type CacheStats struct {
	Hits        int64     `json:"hits"`
	Misses      int64     `json:"misses"`
	HitRate     float64   `json:"hit_rate"`
	TotalKeys   int       `json:"total_keys"`
	MemoryUsage int64     `json:"memory_usage"`
	LastReset   time.Time `json:"last_reset"`
	mutex       sync.RWMutex
}

// IncrementHits 增加命中次数
func (stats *CacheStats) IncrementHits() {
	stats.mutex.Lock()
	defer stats.mutex.Unlock()
	stats.Hits++
	stats.updateHitRate()
}

// IncrementMisses 增加未命中次数
func (stats *CacheStats) IncrementMisses() {
	stats.mutex.Lock()
	defer stats.mutex.Unlock()
	stats.Misses++
	stats.updateHitRate()
}

// updateHitRate 更新命中率
func (stats *CacheStats) updateHitRate() {
	total := stats.Hits + stats.Misses
	if total > 0 {
		stats.HitRate = float64(stats.Hits) / float64(total) * 100
	}
}

// Reset 重置统计信息
func (stats *CacheStats) Reset() {
	stats.mutex.Lock()
	defer stats.mutex.Unlock()
	stats.Hits = 0
	stats.Misses = 0
	stats.HitRate = 0
	stats.LastReset = time.Now()
}

// GetStats 获取统计信息
func (stats *CacheStats) GetStats() CacheStats {
	stats.mutex.RLock()
	defer stats.mutex.RUnlock()
	return *stats
}

// Optimizer 缓存优化器
type Optimizer struct {
	store Store
	stats *CacheStats
}

// NewOptimizer 创建新的缓存优化器
func NewOptimizer(store Store) *Optimizer {
	return &Optimizer{
		store: store,
		stats: &CacheStats{
			LastReset: time.Now(),
		},
	}
}

// WarmUp 缓存预热
func (opt *Optimizer) WarmUp(items map[string]interface{}, ttl time.Duration) error {
	var wg sync.WaitGroup
	errors := make(chan error, len(items))

	for key, value := range items {
		wg.Add(1)
		go func(k string, v interface{}) {
			defer wg.Done()
			if err := opt.store.Set(k, v, ttl); err != nil {
				errors <- fmt.Errorf("failed to warm up cache for key %s: %w", k, err)
			}
		}(key, value)
	}

	wg.Wait()
	close(errors)

	// 收集错误
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("cache warm up failed: %v", errs)
	}

	return nil
}

// WarmUpWithCallback 使用回调函数进行缓存预热
func (opt *Optimizer) WarmUpWithCallback(keys []string, ttl time.Duration, callback func(string) (interface{}, error)) error {
	var wg sync.WaitGroup
	errors := make(chan error, len(keys))

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			// 执行回调函数获取值
			value, err := callback(k)
			if err != nil {
				errors <- fmt.Errorf("callback failed for key %s: %w", k, err)
				return
			}

			// 设置缓存
			if err := opt.store.Set(k, value, ttl); err != nil {
				errors <- fmt.Errorf("failed to warm up cache for key %s: %w", k, err)
			}
		}(key)
	}

	wg.Wait()
	close(errors)

	// 收集错误
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("cache warm up failed: %v", errs)
	}

	return nil
}

// BatchGet 批量获取缓存
func (opt *Optimizer) BatchGet(keys []string) (map[string]interface{}, error) {
	var wg sync.WaitGroup
	results := make(map[string]interface{})
	errors := make(chan error, len(keys))
	resultChan := make(chan struct {
		key   string
		value interface{}
	}, len(keys))

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			value, err := opt.store.Get(k)
			if err != nil {
				opt.stats.IncrementMisses()
				errors <- fmt.Errorf("failed to get cache for key %s: %w", k, err)
				return
			}

			opt.stats.IncrementHits()
			resultChan <- struct {
				key   string
				value interface{}
			}{k, value}
		}(key)
	}

	wg.Wait()
	close(errors)
	close(resultChan)

	// 收集结果
	for result := range resultChan {
		results[result.key] = result.value
	}

	// 检查错误
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return results, fmt.Errorf("batch get failed: %v", errs)
	}

	return results, nil
}

// BatchSet 批量设置缓存
func (opt *Optimizer) BatchSet(items map[string]interface{}, ttl time.Duration) error {
	var wg sync.WaitGroup
	errors := make(chan error, len(items))

	for key, value := range items {
		wg.Add(1)
		go func(k string, v interface{}) {
			defer wg.Done()
			if err := opt.store.Set(k, v, ttl); err != nil {
				errors <- fmt.Errorf("failed to set cache for key %s: %w", k, err)
			}
		}(key, value)
	}

	wg.Wait()
	close(errors)

	// 收集错误
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("batch set failed: %v", errs)
	}

	return nil
}

// BatchDelete 批量删除缓存
func (opt *Optimizer) BatchDelete(keys []string) error {
	return opt.store.DeleteMultiple(keys)
}

// GetStats 获取缓存统计信息
func (opt *Optimizer) GetStats() CacheStats {
	return opt.stats.GetStats()
}

// ResetStats 重置统计信息
func (opt *Optimizer) ResetStats() {
	opt.stats.Reset()
}

// CacheWithStats 带统计的缓存包装器
type CacheWithStats struct {
	store Store
	stats *CacheStats
}

// NewCacheWithStats 创建带统计的缓存包装器
func NewCacheWithStats(store Store) *CacheWithStats {
	return &CacheWithStats{
		store: store,
		stats: &CacheStats{
			LastReset: time.Now(),
		},
	}
}

// Get 获取缓存值（带统计）
func (c *CacheWithStats) Get(key string) (interface{}, error) {
	value, err := c.store.Get(key)
	if err != nil {
		c.stats.IncrementMisses()
		return nil, err
	}

	c.stats.IncrementHits()
	return value, nil
}

// Set 设置缓存值（带统计）
func (c *CacheWithStats) Set(key string, value interface{}, ttl time.Duration) error {
	return c.store.Set(key, value, ttl)
}

// Delete 删除缓存（带统计）
func (c *CacheWithStats) Delete(key string) error {
	return c.store.Delete(key)
}

// Has 检查缓存是否存在（带统计）
func (c *CacheWithStats) Has(key string) bool {
	exists := c.store.Has(key)
	if exists {
		c.stats.IncrementHits()
	} else {
		c.stats.IncrementMisses()
	}
	return exists
}

// GetStats 获取统计信息
func (c *CacheWithStats) GetStats() CacheStats {
	return c.stats.GetStats()
}

// ResetStats 重置统计信息
func (c *CacheWithStats) ResetStats() {
	c.stats.Reset()
}

// 实现Store接口的其他方法
func (c *CacheWithStats) GetString(key string) (string, error) {
	return c.store.GetString(key)
}

func (c *CacheWithStats) GetInt(key string) (int, error) {
	return c.store.GetInt(key)
}

func (c *CacheWithStats) GetFloat(key string) (float64, error) {
	return c.store.GetFloat(key)
}

func (c *CacheWithStats) GetBool(key string) (bool, error) {
	return c.store.GetBool(key)
}

func (c *CacheWithStats) GetBytes(key string) ([]byte, error) {
	return c.store.GetBytes(key)
}

func (c *CacheWithStats) SetString(key string, value string, ttl time.Duration) error {
	return c.store.SetString(key, value, ttl)
}

func (c *CacheWithStats) SetInt(key string, value int, ttl time.Duration) error {
	return c.store.SetInt(key, value, ttl)
}

func (c *CacheWithStats) SetFloat(key string, value float64, ttl time.Duration) error {
	return c.store.SetFloat(key, value, ttl)
}

func (c *CacheWithStats) SetBool(key string, value bool, ttl time.Duration) error {
	return c.store.SetBool(key, value, ttl)
}

func (c *CacheWithStats) SetBytes(key string, value []byte, ttl time.Duration) error {
	return c.store.SetBytes(key, value, ttl)
}

func (c *CacheWithStats) DeleteMultiple(keys []string) error {
	return c.store.DeleteMultiple(keys)
}

func (c *CacheWithStats) Clear() error {
	return c.store.Clear()
}

func (c *CacheWithStats) Missing(key string) bool {
	return c.store.Missing(key)
}

func (c *CacheWithStats) Increment(key string, value int) (int, error) {
	return c.store.Increment(key, value)
}

func (c *CacheWithStats) Decrement(key string, value int) (int, error) {
	return c.store.Decrement(key, value)
}

func (c *CacheWithStats) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	return c.store.Remember(key, ttl, callback)
}

func (c *CacheWithStats) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return c.store.RememberForever(key, callback)
}

func (c *CacheWithStats) Tags(names ...string) TaggedStore {
	return c.store.Tags(names...)
}

func (c *CacheWithStats) Flush() error {
	return c.store.Flush()
}

func (c *CacheWithStats) GetPrefix() string {
	return c.store.GetPrefix()
}

func (c *CacheWithStats) SetPrefix(prefix string) {
	c.store.SetPrefix(prefix)
}
