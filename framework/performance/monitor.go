package performance

import (
	"context"
	"sync"
	"time"
)

// MetricType 指标类型
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
	MetricTypeSummary   MetricType = "summary"
)

// Metric 指标接口
type Metric interface {
	// Name 获取指标名称
	Name() string
	// Type 获取指标类型
	Type() MetricType
	// Value 获取指标值
	Value() interface{}
	// Labels 获取标签
	Labels() map[string]string
	// Timestamp 获取时间戳
	Timestamp() time.Time
}

// Counter 计数器指标
type Counter struct {
	name      string
	value     int64
	labels    map[string]string
	timestamp time.Time
	mu        sync.RWMutex
}

// NewCounter 创建计数器
func NewCounter(name string, labels map[string]string) *Counter {
	return &Counter{
		name:      name,
		labels:    labels,
		timestamp: time.Now(),
	}
}

func (c *Counter) Name() string {
	return c.name
}

func (c *Counter) Type() MetricType {
	return MetricTypeCounter
}

func (c *Counter) Value() interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

func (c *Counter) Labels() map[string]string {
	return c.labels
}

func (c *Counter) Timestamp() time.Time {
	return c.timestamp
}

// Increment 增加计数器
func (c *Counter) Increment(delta int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
	c.timestamp = time.Now()
}

// Reset 重置计数器
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
	c.timestamp = time.Now()
}

// Gauge 仪表指标
type Gauge struct {
	name      string
	value     float64
	labels    map[string]string
	timestamp time.Time
	mu        sync.RWMutex
}

// NewGauge 创建仪表
func NewGauge(name string, labels map[string]string) *Gauge {
	return &Gauge{
		name:      name,
		labels:    labels,
		timestamp: time.Now(),
	}
}

func (g *Gauge) Name() string {
	return g.name
}

func (g *Gauge) Type() MetricType {
	return MetricTypeGauge
}

func (g *Gauge) Value() interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.value
}

func (g *Gauge) Labels() map[string]string {
	return g.labels
}

func (g *Gauge) Timestamp() time.Time {
	return g.timestamp
}

// Set 设置仪表值
func (g *Gauge) Set(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = value
	g.timestamp = time.Now()
}

// Add 增加仪表值
func (g *Gauge) Add(delta float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value += delta
	g.timestamp = time.Now()
}

// Histogram 直方图指标
type Histogram struct {
	name      string
	buckets   map[float64]int64
	sum       float64
	count     int64
	labels    map[string]string
	timestamp time.Time
	mu        sync.RWMutex
}

// NewHistogram 创建直方图
func NewHistogram(name string, buckets []float64, labels map[string]string) *Histogram {
	bucketMap := make(map[float64]int64)
	for _, bucket := range buckets {
		bucketMap[bucket] = 0
	}
	
	return &Histogram{
		name:      name,
		buckets:   bucketMap,
		labels:    labels,
		timestamp: time.Now(),
	}
}

func (h *Histogram) Name() string {
	return h.name
}

func (h *Histogram) Type() MetricType {
	return MetricTypeHistogram
}

func (h *Histogram) Value() interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return map[string]interface{}{
		"buckets": h.buckets,
		"sum":     h.sum,
		"count":   h.count,
	}
}

func (h *Histogram) Labels() map[string]string {
	return h.labels
}

func (h *Histogram) Timestamp() time.Time {
	return h.timestamp
}

// Observe 观察值
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.sum += value
	h.count++
	h.timestamp = time.Now()
	
	for bucket := range h.buckets {
		if value <= bucket {
			h.buckets[bucket]++
		}
	}
}

// Monitor 性能监控器接口
type Monitor interface {
	// RegisterMetric 注册指标
	RegisterMetric(metric Metric)
	// GetMetric 获取指标
	GetMetric(name string) Metric
	// GetAllMetrics 获取所有指标
	GetAllMetrics() map[string]Metric
	// Collect 收集指标
	Collect() []Metric
	// Reset 重置所有指标
	Reset()
	// Start 启动监控
	Start(ctx context.Context) error
	// Stop 停止监控
	Stop() error
}

// PerformanceMonitor 性能监控器实现
type PerformanceMonitor struct {
	metrics map[string]Metric
	mu      sync.RWMutex
	running bool
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewPerformanceMonitor 创建性能监控器
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		metrics: make(map[string]Metric),
	}
}

// RegisterMetric 注册指标
func (pm *PerformanceMonitor) RegisterMetric(metric Metric) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.metrics[metric.Name()] = metric
}

// GetMetric 获取指标
func (pm *PerformanceMonitor) GetMetric(name string) Metric {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.metrics[name]
}

// GetAllMetrics 获取所有指标
func (pm *PerformanceMonitor) GetAllMetrics() map[string]Metric {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	result := make(map[string]Metric)
	for name, metric := range pm.metrics {
		result[name] = metric
	}
	return result
}

// Collect 收集指标
func (pm *PerformanceMonitor) Collect() []Metric {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	var metrics []Metric
	for _, metric := range pm.metrics {
		metrics = append(metrics, metric)
	}
	return metrics
}

// Reset 重置所有指标
func (pm *PerformanceMonitor) Reset() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	for _, metric := range pm.metrics {
		switch m := metric.(type) {
		case *Counter:
			m.Reset()
		case *Gauge:
			m.Set(0)
		case *Histogram:
			// 重置直方图
			for bucket := range m.buckets {
				m.buckets[bucket] = 0
			}
			m.sum = 0
			m.count = 0
		}
	}
}

// Start 启动监控
func (pm *PerformanceMonitor) Start(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if pm.running {
		return nil
	}
	
	pm.ctx, pm.cancel = context.WithCancel(ctx)
	pm.running = true
	
	// 启动后台收集任务
	go pm.collectLoop()
	
	return nil
}

// Stop 停止监控
func (pm *PerformanceMonitor) Stop() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if !pm.running {
		return nil
	}
	
	if pm.cancel != nil {
		pm.cancel()
	}
	pm.running = false
	
	return nil
}

// collectLoop 收集循环
func (pm *PerformanceMonitor) collectLoop() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒收集一次
	defer ticker.Stop()
	
	for {
		select {
		case <-pm.ctx.Done():
			return
		case <-ticker.C:
			// 这里可以添加指标持久化或发送到监控系统的逻辑
			_ = pm.Collect()
		}
	}
} 