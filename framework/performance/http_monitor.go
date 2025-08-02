package performance

import (
	"net/http"
	"sync"
	"time"
)

// HTTPMetrics HTTP指标
type HTTPMetrics struct {
	// 请求计数器
	requestCounter   *Counter
	responseCounter  *Counter
	errorCounter     *Counter
	
	// 响应时间直方图
	responseTimeHistogram *Histogram
	
	// 活跃连接数
	activeConnections *Gauge
	
	// 请求大小和响应大小
	requestSizeHistogram  *Histogram
	responseSizeHistogram *Histogram
}

// NewHTTPMetrics 创建HTTP指标
func NewHTTPMetrics(monitor Monitor) *HTTPMetrics {
	// 创建响应时间直方图，单位为毫秒
	responseTimeBuckets := []float64{10, 50, 100, 200, 500, 1000, 2000, 5000}
	responseTimeHistogram := NewHistogram("http_response_time", responseTimeBuckets, map[string]string{"unit": "milliseconds"})
	monitor.RegisterMetric(responseTimeHistogram)
	
	// 创建请求大小直方图，单位为字节
	requestSizeBuckets := []float64{100, 500, 1000, 5000, 10000, 50000, 100000}
	requestSizeHistogram := NewHistogram("http_request_size", requestSizeBuckets, map[string]string{"unit": "bytes"})
	monitor.RegisterMetric(requestSizeHistogram)
	
	// 创建响应大小直方图，单位为字节
	responseSizeBuckets := []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 1000000}
	responseSizeHistogram := NewHistogram("http_response_size", responseSizeBuckets, map[string]string{"unit": "bytes"})
	monitor.RegisterMetric(responseSizeHistogram)
	
	// 创建计数器
	requestCounter := NewCounter("http_requests_total", map[string]string{"type": "total"})
	monitor.RegisterMetric(requestCounter)
	
	responseCounter := NewCounter("http_responses_total", map[string]string{"type": "total"})
	monitor.RegisterMetric(responseCounter)
	
	errorCounter := NewCounter("http_errors_total", map[string]string{"type": "total"})
	monitor.RegisterMetric(errorCounter)
	
	// 创建活跃连接数仪表
	activeConnections := NewGauge("http_active_connections", map[string]string{"type": "count"})
	monitor.RegisterMetric(activeConnections)
	
	return &HTTPMetrics{
		requestCounter:        requestCounter,
		responseCounter:       responseCounter,
		errorCounter:          errorCounter,
		responseTimeHistogram: responseTimeHistogram,
		activeConnections:     activeConnections,
		requestSizeHistogram:  requestSizeHistogram,
		responseSizeHistogram: responseSizeHistogram,
	}
}

// HTTPMonitor HTTP监控器
type HTTPMonitor struct {
	metrics *HTTPMetrics
	mu      sync.RWMutex
}

// NewHTTPMonitor 创建HTTP监控器
func NewHTTPMonitor(monitor Monitor) *HTTPMonitor {
	return &HTTPMonitor{
		metrics: NewHTTPMetrics(monitor),
	}
}

// RecordRequest 记录请求
func (hm *HTTPMonitor) RecordRequest(method, path string, size int64) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	
	// 增加请求计数器
	hm.metrics.requestCounter.Increment(1)
	
	// 记录请求大小
	hm.metrics.requestSizeHistogram.Observe(float64(size))
	
	// 增加活跃连接数
	hm.metrics.activeConnections.Add(1)
}

// RecordResponse 记录响应
func (hm *HTTPMonitor) RecordResponse(method, path string, statusCode int, size int64, duration time.Duration) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	
	// 增加响应计数器
	hm.metrics.responseCounter.Increment(1)
	
	// 记录响应时间（毫秒）
	hm.metrics.responseTimeHistogram.Observe(float64(duration.Milliseconds()))
	
	// 记录响应大小
	hm.metrics.responseSizeHistogram.Observe(float64(size))
	
	// 减少活跃连接数
	hm.metrics.activeConnections.Add(-1)
	
	// 如果是错误响应，增加错误计数器
	if statusCode >= 400 {
		hm.metrics.errorCounter.Increment(1)
	}
}

// RecordError 记录错误
func (hm *HTTPMonitor) RecordError(method, path string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	
	hm.metrics.errorCounter.Increment(1)
	
	// 减少活跃连接数
	hm.metrics.activeConnections.Add(-1)
}

// GetMetrics 获取指标
func (hm *HTTPMonitor) GetMetrics() *HTTPMetrics {
	return hm.metrics
}

// HTTPMonitorMiddleware HTTP监控中间件
type HTTPMonitorMiddleware struct {
	monitor *HTTPMonitor
}

// NewHTTPMonitorMiddleware 创建HTTP监控中间件
func NewHTTPMonitorMiddleware(monitor Monitor) *HTTPMonitorMiddleware {
	return &HTTPMonitorMiddleware{
		monitor: NewHTTPMonitor(monitor),
	}
}

// ServeHTTP 实现http.Handler接口
func (hm *HTTPMonitorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 记录请求
	hm.monitor.RecordRequest(r.Method, r.URL.Path, r.ContentLength)
	
	// 调用下一个处理器
	// 这里需要实际的中间件链来调用下一个处理器
	// 暂时只是记录请求和响应
}

// responseWriter 响应写入器包装器
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	startTime  time.Time
	monitor    *HTTPMonitor
	method     string
	path       string
	written    bool
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.written {
		rw.statusCode = statusCode
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.written = true
	}
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(200)
	}
	
	// 记录响应
	duration := time.Since(rw.startTime)
	rw.monitor.RecordResponse(rw.method, rw.path, rw.statusCode, int64(len(data)), duration)
	
	return rw.ResponseWriter.Write(data)
}

// RequestMetrics 请求指标结构
type RequestMetrics struct {
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	StatusCode  int               `json:"status_code"`
	Duration    time.Duration     `json:"duration"`
	Size        int64             `json:"size"`
	Headers     map[string]string `json:"headers"`
	Timestamp   time.Time         `json:"timestamp"`
}

// HTTPMetricsCollector HTTP指标收集器
type HTTPMetricsCollector struct {
	monitor Monitor
	metrics []RequestMetrics
	mu      sync.RWMutex
	maxSize int
}

// NewHTTPMetricsCollector 创建HTTP指标收集器
func NewHTTPMetricsCollector(monitor Monitor, maxSize int) *HTTPMetricsCollector {
	if maxSize <= 0 {
		maxSize = 1000
	}
	
	return &HTTPMetricsCollector{
		monitor: monitor,
		maxSize: maxSize,
	}
}

// Collect 收集请求指标
func (hmc *HTTPMetricsCollector) Collect(metrics RequestMetrics) {
	hmc.mu.Lock()
	defer hmc.mu.Unlock()
	
	// 添加到指标列表
	hmc.metrics = append(hmc.metrics, metrics)
	
	// 如果超过最大大小，移除最旧的指标
	if len(hmc.metrics) > hmc.maxSize {
		hmc.metrics = hmc.metrics[1:]
	}
}

// GetMetrics 获取所有指标
func (hmc *HTTPMetricsCollector) GetMetrics() []RequestMetrics {
	hmc.mu.RLock()
	defer hmc.mu.RUnlock()
	
	result := make([]RequestMetrics, len(hmc.metrics))
	copy(result, hmc.metrics)
	return result
}

// GetMetricsByPath 根据路径获取指标
func (hmc *HTTPMetricsCollector) GetMetricsByPath(path string) []RequestMetrics {
	hmc.mu.RLock()
	defer hmc.mu.RUnlock()
	
	var result []RequestMetrics
	for _, metric := range hmc.metrics {
		if metric.Path == path {
			result = append(result, metric)
		}
	}
	return result
}

// GetMetricsByMethod 根据方法获取指标
func (hmc *HTTPMetricsCollector) GetMetricsByMethod(method string) []RequestMetrics {
	hmc.mu.RLock()
	defer hmc.mu.RUnlock()
	
	var result []RequestMetrics
	for _, metric := range hmc.metrics {
		if metric.Method == method {
			result = append(result, metric)
		}
	}
	return result
}

// GetMetricsByStatusCode 根据状态码获取指标
func (hmc *HTTPMetricsCollector) GetMetricsByStatusCode(statusCode int) []RequestMetrics {
	hmc.mu.RLock()
	defer hmc.mu.RUnlock()
	
	var result []RequestMetrics
	for _, metric := range hmc.metrics {
		if metric.StatusCode == statusCode {
			result = append(result, metric)
		}
	}
	return result
}

// GetAverageResponseTime 获取平均响应时间
func (hmc *HTTPMetricsCollector) GetAverageResponseTime() time.Duration {
	hmc.mu.RLock()
	defer hmc.mu.RUnlock()
	
	if len(hmc.metrics) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, metric := range hmc.metrics {
		total += metric.Duration
	}
	
	return total / time.Duration(len(hmc.metrics))
}

// GetErrorRate 获取错误率
func (hmc *HTTPMetricsCollector) GetErrorRate() float64 {
	hmc.mu.RLock()
	defer hmc.mu.RUnlock()
	
	if len(hmc.metrics) == 0 {
		return 0
	}
	
	errorCount := 0
	for _, metric := range hmc.metrics {
		if metric.StatusCode >= 400 {
			errorCount++
		}
	}
	
	return float64(errorCount) / float64(len(hmc.metrics))
}

// GetRequestRate 获取请求率（每秒请求数）
func (hmc *HTTPMetricsCollector) GetRequestRate() float64 {
	hmc.mu.RLock()
	defer hmc.mu.RUnlock()
	
	if len(hmc.metrics) == 0 {
		return 0
	}
	
	// 计算时间范围
	oldest := hmc.metrics[0].Timestamp
	newest := hmc.metrics[len(hmc.metrics)-1].Timestamp
	duration := newest.Sub(oldest)
	
	if duration <= 0 {
		return 0
	}
	
	return float64(len(hmc.metrics)) / duration.Seconds()
}

// Clear 清空指标
func (hmc *HTTPMetricsCollector) Clear() {
	hmc.mu.Lock()
	defer hmc.mu.Unlock()
	hmc.metrics = nil
} 