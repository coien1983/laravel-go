package performance

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// OptimizationType 优化类型
type OptimizationType string

const (
	OptimizationTypeConnectionPool OptimizationType = "connection_pool"
	OptimizationTypeCache          OptimizationType = "cache"
	OptimizationTypeMemory         OptimizationType = "memory"
	OptimizationTypeConcurrency    OptimizationType = "concurrency"
)

// OptimizationResult 优化结果
type OptimizationResult struct {
	Type        OptimizationType `json:"type"`
	Success     bool             `json:"success"`
	Message     string           `json:"message"`
	Improvement float64          `json:"improvement"` // 改进百分比
	Timestamp   time.Time        `json:"timestamp"`
}

// Optimizer 性能优化器接口
type Optimizer interface {
	// Optimize 执行优化
	Optimize(ctx context.Context) (*OptimizationResult, error)
	// GetType 获取优化类型
	GetType() OptimizationType
	// GetDescription 获取优化描述
	GetDescription() string
}

// PerformanceOptimizer 性能优化器
type PerformanceOptimizer struct {
	optimizers []Optimizer
	monitor    Monitor
	mu         sync.RWMutex
}

// NewPerformanceOptimizer 创建性能优化器
func NewPerformanceOptimizer(monitor Monitor) *PerformanceOptimizer {
	po := &PerformanceOptimizer{
		monitor: monitor,
	}
	
	// 添加默认优化器
	po.optimizers = []Optimizer{
		NewConnectionPoolOptimizer(monitor),
		NewCacheOptimizer(monitor),
		NewMemoryOptimizer(monitor),
		NewConcurrencyOptimizer(monitor),
	}
	
	return po
}

// AddOptimizer 添加优化器
func (po *PerformanceOptimizer) AddOptimizer(optimizer Optimizer) {
	po.mu.Lock()
	defer po.mu.Unlock()
	po.optimizers = append(po.optimizers, optimizer)
}

// Optimize 执行所有优化
func (po *PerformanceOptimizer) Optimize(ctx context.Context) ([]*OptimizationResult, error) {
	po.mu.RLock()
	optimizers := make([]Optimizer, len(po.optimizers))
	copy(optimizers, po.optimizers)
	po.mu.RUnlock()
	
	var results []*OptimizationResult
	
	for _, optimizer := range optimizers {
		result, err := optimizer.Optimize(ctx)
		if err != nil {
			result = &OptimizationResult{
				Type:      optimizer.GetType(),
				Success:   false,
				Message:   err.Error(),
				Timestamp: time.Now(),
			}
		}
		results = append(results, result)
	}
	
	return results, nil
}

// OptimizeByType 根据类型执行优化
func (po *PerformanceOptimizer) OptimizeByType(ctx context.Context, optType OptimizationType) (*OptimizationResult, error) {
	po.mu.RLock()
	defer po.mu.RUnlock()
	
	for _, optimizer := range po.optimizers {
		if optimizer.GetType() == optType {
			return optimizer.Optimize(ctx)
		}
	}
	
	return nil, fmt.Errorf("optimizer not found for type: %s", optType)
}

// ConnectionPoolOptimizer 连接池优化器
type ConnectionPoolOptimizer struct {
	monitor Monitor
}

// NewConnectionPoolOptimizer 创建连接池优化器
func NewConnectionPoolOptimizer(monitor Monitor) *ConnectionPoolOptimizer {
	return &ConnectionPoolOptimizer{
		monitor: monitor,
	}
}

func (cpo *ConnectionPoolOptimizer) GetType() OptimizationType {
	return OptimizationTypeConnectionPool
}

func (cpo *ConnectionPoolOptimizer) GetDescription() string {
	return "优化数据库连接池配置"
}

func (cpo *ConnectionPoolOptimizer) Optimize(ctx context.Context) (*OptimizationResult, error) {
	// 获取当前连接池指标
	activeConnections := cpo.monitor.GetMetric("db_active_connections")
	maxConnections := cpo.monitor.GetMetric("db_max_connections")
	
	if activeConnections == nil || maxConnections == nil {
		return &OptimizationResult{
			Type:      cpo.GetType(),
			Success:   false,
			Message:   "连接池指标不可用",
			Timestamp: time.Now(),
		}, nil
	}
	
	activeValue := activeConnections.Value().(int64)
	maxValue := maxConnections.Value().(int64)
	
	// 计算连接池使用率
	usageRate := float64(activeValue) / float64(maxValue)
	
	var message string
	var improvement float64
	
	if usageRate > 0.8 {
		// 连接池使用率过高，建议增加最大连接数
		message = fmt.Sprintf("连接池使用率过高 (%.1f%%)，建议增加最大连接数", usageRate*100)
		improvement = 15.0 // 预计改进15%
	} else if usageRate < 0.2 {
		// 连接池使用率过低，建议减少最大连接数
		message = fmt.Sprintf("连接池使用率过低 (%.1f%%)，建议减少最大连接数以节省资源", usageRate*100)
		improvement = 10.0 // 预计改进10%
	} else {
		// 连接池使用率正常
		message = fmt.Sprintf("连接池使用率正常 (%.1f%%)", usageRate*100)
		improvement = 0.0
	}
	
	return &OptimizationResult{
		Type:        cpo.GetType(),
		Success:     true,
		Message:     message,
		Improvement: improvement,
		Timestamp:   time.Now(),
	}, nil
}

// CacheOptimizer 缓存优化器
type CacheOptimizer struct {
	monitor Monitor
}

// NewCacheOptimizer 创建缓存优化器
func NewCacheOptimizer(monitor Monitor) *CacheOptimizer {
	return &CacheOptimizer{
		monitor: monitor,
	}
}

func (co *CacheOptimizer) GetType() OptimizationType {
	return OptimizationTypeCache
}

func (co *CacheOptimizer) GetDescription() string {
	return "优化缓存策略和配置"
}

func (co *CacheOptimizer) Optimize(ctx context.Context) (*OptimizationResult, error) {
	// 获取缓存命中率指标
	cacheHits := co.monitor.GetMetric("cache_hits")
	cacheMisses := co.monitor.GetMetric("cache_misses")
	
	if cacheHits == nil || cacheMisses == nil {
		return &OptimizationResult{
			Type:      co.GetType(),
			Success:   false,
			Message:   "缓存指标不可用",
			Timestamp: time.Now(),
		}, nil
	}
	
	hits := cacheHits.Value().(int64)
	misses := cacheMisses.Value().(int64)
	
	if hits+misses == 0 {
		return &OptimizationResult{
			Type:      co.GetType(),
			Success:   false,
			Message:   "没有缓存访问记录",
			Timestamp: time.Now(),
		}, nil
	}
	
	hitRate := float64(hits) / float64(hits+misses)
	
	var message string
	var improvement float64
	
	if hitRate < 0.7 {
		// 缓存命中率过低
		message = fmt.Sprintf("缓存命中率过低 (%.1f%%)，建议增加缓存大小或优化缓存策略", hitRate*100)
		improvement = 20.0 // 预计改进20%
	} else if hitRate < 0.85 {
		// 缓存命中率一般
		message = fmt.Sprintf("缓存命中率一般 (%.1f%%)，建议优化缓存策略", hitRate*100)
		improvement = 10.0 // 预计改进10%
	} else {
		// 缓存命中率良好
		message = fmt.Sprintf("缓存命中率良好 (%.1f%%)", hitRate*100)
		improvement = 0.0
	}
	
	return &OptimizationResult{
		Type:        co.GetType(),
		Success:     true,
		Message:     message,
		Improvement: improvement,
		Timestamp:   time.Now(),
	}, nil
}

// MemoryOptimizer 内存优化器
type MemoryOptimizer struct {
	monitor Monitor
}

// NewMemoryOptimizer 创建内存优化器
func NewMemoryOptimizer(monitor Monitor) *MemoryOptimizer {
	return &MemoryOptimizer{
		monitor: monitor,
	}
}

func (mo *MemoryOptimizer) GetType() OptimizationType {
	return OptimizationTypeMemory
}

func (mo *MemoryOptimizer) GetDescription() string {
	return "优化内存使用和垃圾回收"
}

func (mo *MemoryOptimizer) Optimize(ctx context.Context) (*OptimizationResult, error) {
	// 获取内存使用指标
	heapAlloc := mo.monitor.GetMetric("go_heap_alloc")
	heapSys := mo.monitor.GetMetric("go_heap_sys")
	
	if heapAlloc == nil || heapSys == nil {
		return &OptimizationResult{
			Type:      mo.GetType(),
			Success:   false,
			Message:   "内存指标不可用",
			Timestamp: time.Now(),
		}, nil
	}
	
	alloc := heapAlloc.Value().(float64)
	sys := heapSys.Value().(float64)
	
	if sys == 0 {
		return &OptimizationResult{
			Type:      mo.GetType(),
			Success:   false,
			Message:   "内存系统指标无效",
			Timestamp: time.Now(),
		}, nil
	}
	
	// 计算内存使用率
	memoryUsage := alloc / sys
	
	var message string
	var improvement float64
	
	if memoryUsage > 0.8 {
		// 内存使用率过高
		message = fmt.Sprintf("内存使用率过高 (%.1f%%)，建议增加内存或优化内存使用", memoryUsage*100)
		improvement = 25.0 // 预计改进25%
	} else if memoryUsage < 0.3 {
		// 内存使用率过低
		message = fmt.Sprintf("内存使用率过低 (%.1f%%)，可以考虑减少内存分配", memoryUsage*100)
		improvement = 5.0 // 预计改进5%
	} else {
		// 内存使用率正常
		message = fmt.Sprintf("内存使用率正常 (%.1f%%)", memoryUsage*100)
		improvement = 0.0
	}
	
	// 强制垃圾回收
	runtime.GC()
	
	return &OptimizationResult{
		Type:        mo.GetType(),
		Success:     true,
		Message:     message,
		Improvement: improvement,
		Timestamp:   time.Now(),
	}, nil
}

// ConcurrencyOptimizer 并发优化器
type ConcurrencyOptimizer struct {
	monitor Monitor
}

// NewConcurrencyOptimizer 创建并发优化器
func NewConcurrencyOptimizer(monitor Monitor) *ConcurrencyOptimizer {
	return &ConcurrencyOptimizer{
		monitor: monitor,
	}
}

func (co *ConcurrencyOptimizer) GetType() OptimizationType {
	return OptimizationTypeConcurrency
}

func (co *ConcurrencyOptimizer) GetDescription() string {
	return "优化并发处理和协程使用"
}

func (co *ConcurrencyOptimizer) Optimize(ctx context.Context) (*OptimizationResult, error) {
	// 获取协程数量指标
	goroutines := co.monitor.GetMetric("go_goroutines")
	
	if goroutines == nil {
		return &OptimizationResult{
			Type:      co.GetType(),
			Success:   false,
			Message:   "协程指标不可用",
			Timestamp: time.Now(),
		}, nil
	}
	
	goroutineCount := goroutines.Value().(float64)
	
	var message string
	var improvement float64
	
	if goroutineCount > 10000 {
		// 协程数量过多
		message = fmt.Sprintf("协程数量过多 (%.0f)，建议优化并发处理逻辑", goroutineCount)
		improvement = 30.0 // 预计改进30%
	} else if goroutineCount > 5000 {
		// 协程数量较多
		message = fmt.Sprintf("协程数量较多 (%.0f)，建议检查是否有协程泄漏", goroutineCount)
		improvement = 15.0 // 预计改进15%
	} else {
		// 协程数量正常
		message = fmt.Sprintf("协程数量正常 (%.0f)", goroutineCount)
		improvement = 0.0
	}
	
	return &OptimizationResult{
		Type:        co.GetType(),
		Success:     true,
		Message:     message,
		Improvement: improvement,
		Timestamp:   time.Now(),
	}, nil
}

// AutoOptimizer 自动优化器
type AutoOptimizer struct {
	optimizer *PerformanceOptimizer
	interval  time.Duration
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.RWMutex
}

// NewAutoOptimizer 创建自动优化器
func NewAutoOptimizer(optimizer *PerformanceOptimizer, interval time.Duration) *AutoOptimizer {
	if interval <= 0 {
		interval = 5 * time.Minute // 默认5分钟
	}
	
	return &AutoOptimizer{
		optimizer: optimizer,
		interval:  interval,
	}
}

// Start 启动自动优化
func (ao *AutoOptimizer) Start(ctx context.Context) error {
	ao.mu.Lock()
	defer ao.mu.Unlock()
	
	if ao.running {
		return nil
	}
	
	ao.ctx, ao.cancel = context.WithCancel(ctx)
	ao.running = true
	
	// 启动自动优化循环
	go ao.optimizationLoop()
	
	return nil
}

// Stop 停止自动优化
func (ao *AutoOptimizer) Stop() error {
	ao.mu.Lock()
	defer ao.mu.Unlock()
	
	if !ao.running {
		return nil
	}
	
	if ao.cancel != nil {
		ao.cancel()
	}
	ao.running = false
	
	return nil
}

// optimizationLoop 优化循环
func (ao *AutoOptimizer) optimizationLoop() {
	ticker := time.NewTicker(ao.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ao.ctx.Done():
			return
		case <-ticker.C:
			// 执行自动优化
			results, err := ao.optimizer.Optimize(ao.ctx)
			if err != nil {
				// 记录错误但不中断
				continue
			}
			
			// 这里可以添加优化结果的日志记录或通知
			_ = results
		}
	}
} 