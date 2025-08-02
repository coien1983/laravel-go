package performance

import (
	"sync"
	"time"
)

// CacheMetrics 缓存指标
type CacheMetrics struct {
	// 命中率指标
	hitCounter  *Counter
	missCounter *Counter

	// 性能指标
	getTimeHistogram    *Histogram
	setTimeHistogram    *Histogram
	deleteTimeHistogram *Histogram

	// 存储指标
	itemCount     *Gauge
	memoryUsage   *Gauge
	evictionCount *Counter

	// 错误指标
	errorCounter *Counter

	// 操作计数器
	getCounter    *Counter
	setCounter    *Counter
	deleteCounter *Counter
	clearCounter  *Counter
}

// NewCacheMetrics 创建缓存指标
func NewCacheMetrics(monitor Monitor) *CacheMetrics {
	// 创建时间直方图，单位为微秒
	getTimeBuckets := []float64{1, 5, 10, 50, 100, 200, 500, 1000}
	getTimeHistogram := NewHistogram("cache_get_time", getTimeBuckets, map[string]string{"unit": "microseconds"})
	monitor.RegisterMetric(getTimeHistogram)

	setTimeBuckets := []float64{1, 5, 10, 50, 100, 200, 500, 1000, 2000}
	setTimeHistogram := NewHistogram("cache_set_time", setTimeBuckets, map[string]string{"unit": "microseconds"})
	monitor.RegisterMetric(setTimeHistogram)

	deleteTimeBuckets := []float64{1, 5, 10, 50, 100, 200, 500}
	deleteTimeHistogram := NewHistogram("cache_delete_time", deleteTimeBuckets, map[string]string{"unit": "microseconds"})
	monitor.RegisterMetric(deleteTimeHistogram)

	// 创建计数器
	hitCounter := NewCounter("cache_hits_total", map[string]string{"type": "hit"})
	monitor.RegisterMetric(hitCounter)

	missCounter := NewCounter("cache_misses_total", map[string]string{"type": "miss"})
	monitor.RegisterMetric(missCounter)

	errorCounter := NewCounter("cache_errors_total", map[string]string{"type": "error"})
	monitor.RegisterMetric(errorCounter)

	evictionCount := NewCounter("cache_evictions_total", map[string]string{"type": "eviction"})
	monitor.RegisterMetric(evictionCount)

	// 操作计数器
	getCounter := NewCounter("cache_operations_get", map[string]string{"operation": "get"})
	monitor.RegisterMetric(getCounter)

	setCounter := NewCounter("cache_operations_set", map[string]string{"operation": "set"})
	monitor.RegisterMetric(setCounter)

	deleteCounter := NewCounter("cache_operations_delete", map[string]string{"operation": "delete"})
	monitor.RegisterMetric(deleteCounter)

	clearCounter := NewCounter("cache_operations_clear", map[string]string{"operation": "clear"})
	monitor.RegisterMetric(clearCounter)

	// 存储指标
	itemCount := NewGauge("cache_items_count", map[string]string{"type": "count"})
	monitor.RegisterMetric(itemCount)

	memoryUsage := NewGauge("cache_memory_usage", map[string]string{"unit": "bytes"})
	monitor.RegisterMetric(memoryUsage)

	return &CacheMetrics{
		hitCounter:          hitCounter,
		missCounter:         missCounter,
		getTimeHistogram:    getTimeHistogram,
		setTimeHistogram:    setTimeHistogram,
		deleteTimeHistogram: deleteTimeHistogram,
		itemCount:           itemCount,
		memoryUsage:         memoryUsage,
		evictionCount:       evictionCount,
		errorCounter:        errorCounter,
		getCounter:          getCounter,
		setCounter:          setCounter,
		deleteCounter:       deleteCounter,
		clearCounter:        clearCounter,
	}
}

// CacheMonitor 缓存监控器
type CacheMonitor struct {
	metrics          *CacheMetrics
	mu               sync.RWMutex
	operationHistory []CacheOperation
	maxHistorySize   int
}

// CacheOperation 缓存操作记录
type CacheOperation struct {
	Operation string        `json:"operation"`
	Key       string        `json:"key"`
	Duration  time.Duration `json:"duration"`
	Success   bool          `json:"success"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	Hit       bool          `json:"hit,omitempty"`
}

// NewCacheMonitor 创建缓存监控器
func NewCacheMonitor(monitor Monitor) *CacheMonitor {
	return &CacheMonitor{
		metrics:          NewCacheMetrics(monitor),
		operationHistory: make([]CacheOperation, 0),
		maxHistorySize:   1000,
	}
}

// RecordGet 记录获取操作
func (cm *CacheMonitor) RecordGet(key string, duration time.Duration, hit bool, err error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 增加操作计数器
	cm.metrics.getCounter.Increment(1)

	// 记录命中情况
	if hit {
		cm.metrics.hitCounter.Increment(1)
	} else {
		cm.metrics.missCounter.Increment(1)
	}

	// 记录时间
	cm.metrics.getTimeHistogram.Observe(float64(duration.Microseconds()))

	// 记录错误
	if err != nil {
		cm.metrics.errorCounter.Increment(1)
	}

	// 添加到历史记录
	operation := CacheOperation{
		Operation: "GET",
		Key:       key,
		Duration:  duration,
		Success:   err == nil,
		Hit:       hit,
		Timestamp: time.Now(),
	}
	if err != nil {
		operation.Error = err.Error()
	}

	cm.addToHistory(operation)
}

// RecordSet 记录设置操作
func (cm *CacheMonitor) RecordSet(key string, duration time.Duration, err error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 增加操作计数器
	cm.metrics.setCounter.Increment(1)

	// 记录时间
	cm.metrics.setTimeHistogram.Observe(float64(duration.Microseconds()))

	// 记录错误
	if err != nil {
		cm.metrics.errorCounter.Increment(1)
	}

	// 添加到历史记录
	operation := CacheOperation{
		Operation: "SET",
		Key:       key,
		Duration:  duration,
		Success:   err == nil,
		Timestamp: time.Now(),
	}
	if err != nil {
		operation.Error = err.Error()
	}

	cm.addToHistory(operation)
}

// RecordDelete 记录删除操作
func (cm *CacheMonitor) RecordDelete(key string, duration time.Duration, err error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 增加操作计数器
	cm.metrics.deleteCounter.Increment(1)

	// 记录时间
	cm.metrics.deleteTimeHistogram.Observe(float64(duration.Microseconds()))

	// 记录错误
	if err != nil {
		cm.metrics.errorCounter.Increment(1)
	}

	// 添加到历史记录
	operation := CacheOperation{
		Operation: "DELETE",
		Key:       key,
		Duration:  duration,
		Success:   err == nil,
		Timestamp: time.Now(),
	}
	if err != nil {
		operation.Error = err.Error()
	}

	cm.addToHistory(operation)
}

// RecordClear 记录清空操作
func (cm *CacheMonitor) RecordClear(duration time.Duration, err error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 增加操作计数器
	cm.metrics.clearCounter.Increment(1)

	// 记录错误
	if err != nil {
		cm.metrics.errorCounter.Increment(1)
	}

	// 添加到历史记录
	operation := CacheOperation{
		Operation: "CLEAR",
		Key:       "",
		Duration:  duration,
		Success:   err == nil,
		Timestamp: time.Now(),
	}
	if err != nil {
		operation.Error = err.Error()
	}

	cm.addToHistory(operation)
}

// RecordEviction 记录驱逐操作
func (cm *CacheMonitor) RecordEviction(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.metrics.evictionCount.Increment(1)
}

// UpdateStorageMetrics 更新存储指标
func (cm *CacheMonitor) UpdateStorageMetrics(itemCount int, memoryUsage int64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.metrics.itemCount.Set(float64(itemCount))
	cm.metrics.memoryUsage.Set(float64(memoryUsage))
}

// GetMetrics 获取指标
func (cm *CacheMonitor) GetMetrics() *CacheMetrics {
	return cm.metrics
}

// GetHitRate 获取命中率
func (cm *CacheMonitor) GetHitRate() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	hits := cm.metrics.hitCounter.Value().(int64)
	misses := cm.metrics.missCounter.Value().(int64)
	total := hits + misses

	if total == 0 {
		return 0.0
	}

	return float64(hits) / float64(total) * 100.0
}

// GetOperationHistory 获取操作历史
func (cm *CacheMonitor) GetOperationHistory() []CacheOperation {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	result := make([]CacheOperation, len(cm.operationHistory))
	copy(result, cm.operationHistory)
	return result
}

// GetSlowOperations 获取慢操作
func (cm *CacheMonitor) GetSlowOperations(threshold time.Duration) []CacheOperation {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var slowOperations []CacheOperation
	for _, operation := range cm.operationHistory {
		if operation.Duration > threshold {
			slowOperations = append(slowOperations, operation)
		}
	}
	return slowOperations
}

// GetErrorOperations 获取错误操作
func (cm *CacheMonitor) GetErrorOperations() []CacheOperation {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var errorOperations []CacheOperation
	for _, operation := range cm.operationHistory {
		if !operation.Success {
			errorOperations = append(errorOperations, operation)
		}
	}
	return errorOperations
}

// GetAverageOperationTime 获取平均操作时间
func (cm *CacheMonitor) GetAverageOperationTime(operationType string) time.Duration {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if len(cm.operationHistory) == 0 {
		return 0
	}

	var total time.Duration
	count := 0

	for _, operation := range cm.operationHistory {
		if operation.Operation == operationType {
			total += operation.Duration
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return total / time.Duration(count)
}

// GetOperationDistribution 获取操作分布
func (cm *CacheMonitor) GetOperationDistribution() map[string]int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	distribution := make(map[string]int)
	for _, operation := range cm.operationHistory {
		distribution[operation.Operation]++
	}
	return distribution
}

// ClearHistory 清空历史记录
func (cm *CacheMonitor) ClearHistory() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.operationHistory = make([]CacheOperation, 0)
}

// addToHistory 添加到历史记录
func (cm *CacheMonitor) addToHistory(operation CacheOperation) {
	cm.operationHistory = append(cm.operationHistory, operation)

	// 限制历史记录大小
	if len(cm.operationHistory) > cm.maxHistorySize {
		cm.operationHistory = cm.operationHistory[1:]
	}
}

// CacheMonitorMiddleware 缓存监控中间件
type CacheMonitorMiddleware struct {
	monitor *CacheMonitor
}

// NewCacheMonitorMiddleware 创建缓存监控中间件
func NewCacheMonitorMiddleware(monitor *CacheMonitor) *CacheMonitorMiddleware {
	return &CacheMonitorMiddleware{
		monitor: monitor,
	}
}

// WrapGet 包装获取操作
func (cm *CacheMonitorMiddleware) WrapGet(key string, getFunc func() (interface{}, bool, error)) (interface{}, bool, error) {
	start := time.Now()

	value, hit, err := getFunc()

	duration := time.Since(start)
	cm.monitor.RecordGet(key, duration, hit, err)

	return value, hit, err
}

// WrapSet 包装设置操作
func (cm *CacheMonitorMiddleware) WrapSet(key string, value interface{}, setFunc func() error) error {
	start := time.Now()

	err := setFunc()

	duration := time.Since(start)
	cm.monitor.RecordSet(key, duration, err)

	return err
}

// WrapDelete 包装删除操作
func (cm *CacheMonitorMiddleware) WrapDelete(key string, deleteFunc func() error) error {
	start := time.Now()

	err := deleteFunc()

	duration := time.Since(start)
	cm.monitor.RecordDelete(key, duration, err)

	return err
}
