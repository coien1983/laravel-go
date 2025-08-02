package middleware

import (
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"laravel-go/framework/performance"
)

// PerformanceMiddleware HTTP性能监控中间件
type PerformanceMiddleware struct {
	monitor performance.Monitor
	metrics *HTTPMetrics
	enabled bool
	mu      sync.RWMutex
}

// HTTPMetrics HTTP指标
type HTTPMetrics struct {
	requestCount      int64
	responseCount     int64
	errorCount        int64
	activeRequests    int64
	totalResponseTime int64
	avgResponseTime   float64
	mu                sync.RWMutex
}

// NewPerformanceMiddleware 创建性能监控中间件
func NewPerformanceMiddleware(monitor performance.Monitor) *PerformanceMiddleware {
	pm := &PerformanceMiddleware{
		monitor: monitor,
		metrics: &HTTPMetrics{},
		enabled: true,
	}

	// 注册指标
	if monitor != nil {
		pm.registerMetrics()
	}

	return pm
}

// registerMetrics 注册性能指标
func (pm *PerformanceMiddleware) registerMetrics() {
	// 请求计数器
	requestCounter := performance.NewCounter("http_requests_total", map[string]string{"type": "total"})
	pm.monitor.RegisterMetric(requestCounter)

	// 响应计数器
	responseCounter := performance.NewCounter("http_responses_total", map[string]string{"type": "total"})
	pm.monitor.RegisterMetric(responseCounter)

	// 错误计数器
	errorCounter := performance.NewCounter("http_errors_total", map[string]string{"type": "total"})
	pm.monitor.RegisterMetric(errorCounter)

	// 活跃请求数
	activeRequests := performance.NewGauge("http_active_requests", map[string]string{"type": "count"})
	pm.monitor.RegisterMetric(activeRequests)

	// 响应时间直方图
	responseTimeBuckets := []float64{10, 50, 100, 200, 500, 1000, 2000, 5000}
	responseTimeHistogram := performance.NewHistogram("http_response_time", responseTimeBuckets, map[string]string{"unit": "milliseconds"})
	pm.monitor.RegisterMetric(responseTimeHistogram)

	// 内存使用仪表
	memoryUsage := performance.NewGauge("http_memory_usage", map[string]string{"unit": "bytes"})
	pm.monitor.RegisterMetric(memoryUsage)

	// Goroutine数量仪表
	goroutineCount := performance.NewGauge("http_goroutines", map[string]string{"type": "count"})
	pm.monitor.RegisterMetric(goroutineCount)
}

// Handle 处理HTTP请求
func (pm *PerformanceMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !pm.enabled {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// 增加活跃请求数
		atomic.AddInt64(&pm.metrics.activeRequests, 1)
		atomic.AddInt64(&pm.metrics.requestCount, 1)

		// 更新指标
		pm.updateMetrics("request", 1)

		// 包装响应写入器以捕获状态码
		responseWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// 执行请求
		next.ServeHTTP(responseWriter, r)

		// 计算响应时间
		duration := time.Since(start)
		durationMs := float64(duration.Nanoseconds()) / 1e6

		// 减少活跃请求数
		atomic.AddInt64(&pm.metrics.activeRequests, -1)
		atomic.AddInt64(&pm.metrics.responseCount, 1)

		// 更新总响应时间
		atomic.AddInt64(&pm.metrics.totalResponseTime, int64(duration))

		// 计算平均响应时间
		pm.mu.Lock()
		responseCount := atomic.LoadInt64(&pm.metrics.responseCount)
		if responseCount > 0 {
			pm.metrics.avgResponseTime = float64(atomic.LoadInt64(&pm.metrics.totalResponseTime)) / float64(responseCount) / 1e6
		}
		pm.mu.Unlock()

		// 检查是否为错误响应
		if responseWriter.statusCode >= 400 {
			atomic.AddInt64(&pm.metrics.errorCount, 1)
			pm.updateMetrics("error", 1)
		}

		// 更新响应时间指标
		pm.updateResponseTimeMetrics(durationMs)

		// 更新系统指标
		pm.updateSystemMetrics()
	})
}

// updateMetrics 更新指标
func (pm *PerformanceMiddleware) updateMetrics(metricType string, value int64) {
	if pm.monitor == nil {
		return
	}

	switch metricType {
	case "request":
		if counter := pm.monitor.GetMetric("http_requests_total"); counter != nil {
			if c, ok := counter.(*performance.Counter); ok {
				c.Increment(value)
			}
		}
	case "response":
		if counter := pm.monitor.GetMetric("http_responses_total"); counter != nil {
			if c, ok := counter.(*performance.Counter); ok {
				c.Increment(value)
			}
		}
	case "error":
		if counter := pm.monitor.GetMetric("http_errors_total"); counter != nil {
			if c, ok := counter.(*performance.Counter); ok {
				c.Increment(value)
			}
		}
	}
}

// updateResponseTimeMetrics 更新响应时间指标
func (pm *PerformanceMiddleware) updateResponseTimeMetrics(durationMs float64) {
	if pm.monitor == nil {
		return
	}

	// 更新活跃请求数
	if gauge := pm.monitor.GetMetric("http_active_requests"); gauge != nil {
		if g, ok := gauge.(*performance.Gauge); ok {
			g.Set(float64(atomic.LoadInt64(&pm.metrics.activeRequests)))
		}
	}

	// 更新响应时间直方图
	if histogram := pm.monitor.GetMetric("http_response_time"); histogram != nil {
		if h, ok := histogram.(*performance.Histogram); ok {
			h.Observe(durationMs)
		}
	}
}

// updateSystemMetrics 更新系统指标
func (pm *PerformanceMiddleware) updateSystemMetrics() {
	if pm.monitor == nil {
		return
	}

	// 更新内存使用
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	if gauge := pm.monitor.GetMetric("http_memory_usage"); gauge != nil {
		if g, ok := gauge.(*performance.Gauge); ok {
			g.Set(float64(m.Alloc))
		}
	}

	// 更新Goroutine数量
	if gauge := pm.monitor.GetMetric("http_goroutines"); gauge != nil {
		if g, ok := gauge.(*performance.Gauge); ok {
			g.Set(float64(runtime.NumGoroutine()))
		}
	}
}

// GetMetrics 获取性能指标
func (pm *PerformanceMiddleware) GetMetrics() map[string]interface{} {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return map[string]interface{}{
		"request_count":       atomic.LoadInt64(&pm.metrics.requestCount),
		"response_count":      atomic.LoadInt64(&pm.metrics.responseCount),
		"error_count":         atomic.LoadInt64(&pm.metrics.errorCount),
		"active_requests":     atomic.LoadInt64(&pm.metrics.activeRequests),
		"avg_response_time":   pm.metrics.avgResponseTime,
		"total_response_time": atomic.LoadInt64(&pm.metrics.totalResponseTime),
	}
}

// Enable 启用性能监控
func (pm *PerformanceMiddleware) Enable() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.enabled = true
}

// Disable 禁用性能监控
func (pm *PerformanceMiddleware) Disable() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.enabled = false
}

// responseWriter 包装响应写入器
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 写入状态码
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write 写入响应
func (rw *responseWriter) Write(data []byte) (int, error) {
	return rw.ResponseWriter.Write(data)
}
