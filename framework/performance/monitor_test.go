package performance

import (
	"context"
	"testing"
	"time"
)

func TestCounter(t *testing.T) {
	counter := NewCounter("test_counter", map[string]string{"test": "value"})
	
	// 测试初始值
	if counter.Value() != int64(0) {
		t.Errorf("Expected initial value 0, got %v", counter.Value())
	}
	
	// 测试增加
	counter.Increment(5)
	if counter.Value() != int64(5) {
		t.Errorf("Expected value 5 after increment, got %v", counter.Value())
	}
	
	// 测试重置
	counter.Reset()
	if counter.Value() != int64(0) {
		t.Errorf("Expected value 0 after reset, got %v", counter.Value())
	}
	
	// 测试类型
	if counter.Type() != MetricTypeCounter {
		t.Errorf("Expected type counter, got %s", counter.Type())
	}
	
	// 测试名称
	if counter.Name() != "test_counter" {
		t.Errorf("Expected name test_counter, got %s", counter.Name())
	}
	
	// 测试标签
	labels := counter.Labels()
	if labels["test"] != "value" {
		t.Errorf("Expected label test=value, got %v", labels)
	}
}

func TestGauge(t *testing.T) {
	gauge := NewGauge("test_gauge", map[string]string{"test": "value"})
	
	// 测试初始值
	if gauge.Value() != float64(0) {
		t.Errorf("Expected initial value 0, got %v", gauge.Value())
	}
	
	// 测试设置值
	gauge.Set(10.5)
	if gauge.Value() != 10.5 {
		t.Errorf("Expected value 10.5 after set, got %v", gauge.Value())
	}
	
	// 测试增加值
	gauge.Add(5.5)
	if gauge.Value() != 16.0 {
		t.Errorf("Expected value 16.0 after add, got %v", gauge.Value())
	}
	
	// 测试类型
	if gauge.Type() != MetricTypeGauge {
		t.Errorf("Expected type gauge, got %s", gauge.Type())
	}
	
	// 测试名称
	if gauge.Name() != "test_gauge" {
		t.Errorf("Expected name test_gauge, got %s", gauge.Name())
	}
}

func TestHistogram(t *testing.T) {
	buckets := []float64{10, 50, 100}
	histogram := NewHistogram("test_histogram", buckets, map[string]string{"test": "value"})
	
	// 测试初始值
	value := histogram.Value().(map[string]interface{})
	if value["count"] != int64(0) {
		t.Errorf("Expected initial count 0, got %v", value["count"])
	}
	if value["sum"] != float64(0) {
		t.Errorf("Expected initial sum 0, got %v", value["sum"])
	}
	
	// 测试观察值
	histogram.Observe(25)
	histogram.Observe(75)
	histogram.Observe(150)
	
	value = histogram.Value().(map[string]interface{})
	if value["count"] != int64(3) {
		t.Errorf("Expected count 3, got %v", value["count"])
	}
	if value["sum"] != float64(250) {
		t.Errorf("Expected sum 250, got %v", value["sum"])
	}
	
	// 测试类型
	if histogram.Type() != MetricTypeHistogram {
		t.Errorf("Expected type histogram, got %s", histogram.Type())
	}
	
	// 测试名称
	if histogram.Name() != "test_histogram" {
		t.Errorf("Expected name test_histogram, got %s", histogram.Name())
	}
}

func TestPerformanceMonitor(t *testing.T) {
	monitor := NewPerformanceMonitor()
	
	// 测试注册指标
	counter := NewCounter("test_counter", nil)
	monitor.RegisterMetric(counter)
	
	// 测试获取指标
	retrieved := monitor.GetMetric("test_counter")
	if retrieved != counter {
		t.Errorf("Expected to retrieve the same counter")
	}
	
	// 测试获取所有指标
	allMetrics := monitor.GetAllMetrics()
	if len(allMetrics) != 1 {
		t.Errorf("Expected 1 metric, got %d", len(allMetrics))
	}
	
	// 测试收集指标
	collected := monitor.Collect()
	if len(collected) != 1 {
		t.Errorf("Expected 1 collected metric, got %d", len(collected))
	}
	
	// 测试重置
	counter.Increment(5)
	monitor.Reset()
	if counter.Value() != int64(0) {
		t.Errorf("Expected counter to be reset to 0")
	}
}

func TestPerformanceMonitorStartStop(t *testing.T) {
	monitor := NewPerformanceMonitor()
	ctx := context.Background()
	
	// 测试启动
	err := monitor.Start(ctx)
	if err != nil {
		t.Errorf("Expected no error on start, got %v", err)
	}
	
	// 测试重复启动
	err = monitor.Start(ctx)
	if err != nil {
		t.Errorf("Expected no error on duplicate start, got %v", err)
	}
	
	// 等待一段时间让收集循环运行
	time.Sleep(100 * time.Millisecond)
	
	// 测试停止
	err = monitor.Stop()
	if err != nil {
		t.Errorf("Expected no error on stop, got %v", err)
	}
	
	// 测试重复停止
	err = monitor.Stop()
	if err != nil {
		t.Errorf("Expected no error on duplicate stop, got %v", err)
	}
}

func TestHTTPMonitor(t *testing.T) {
	monitor := NewPerformanceMonitor()
	httpMonitor := NewHTTPMonitor(monitor)
	
	// 测试记录请求
	httpMonitor.RecordRequest("GET", "/test", 100)
	
	// 验证请求计数器
	requestCounter := monitor.GetMetric("http_requests_total")
	if requestCounter == nil {
		t.Fatal("Expected http_requests_total metric to exist")
	}
	if requestCounter.Value() != int64(1) {
		t.Errorf("Expected 1 request, got %v", requestCounter.Value())
	}
	
	// 测试记录响应
	httpMonitor.RecordResponse("GET", "/test", 200, 500, 50*time.Millisecond)
	
	// 验证响应计数器
	responseCounter := monitor.GetMetric("http_responses_total")
	if responseCounter == nil {
		t.Fatal("Expected http_responses_total metric to exist")
	}
	if responseCounter.Value() != int64(1) {
		t.Errorf("Expected 1 response, got %v", responseCounter.Value())
	}
	
	// 测试记录错误
	httpMonitor.RecordError("GET", "/test")
	
	// 验证错误计数器
	errorCounter := monitor.GetMetric("http_errors_total")
	if errorCounter == nil {
		t.Fatal("Expected http_errors_total metric to exist")
	}
	if errorCounter.Value() != int64(1) {
		t.Errorf("Expected 1 error, got %v", errorCounter.Value())
	}
}

func TestHTTPMetricsCollector(t *testing.T) {
	monitor := NewPerformanceMonitor()
	collector := NewHTTPMetricsCollector(monitor, 10)
	
	// 测试收集指标
	metrics := RequestMetrics{
		Method:     "GET",
		Path:       "/test",
		StatusCode: 200,
		Duration:   50 * time.Millisecond,
		Size:       500,
		Timestamp:  time.Now(),
	}
	
	collector.Collect(metrics)
	
	// 验证收集的指标
	collected := collector.GetMetrics()
	if len(collected) != 1 {
		t.Errorf("Expected 1 collected metric, got %d", len(collected))
	}
	
	if collected[0].Method != "GET" {
		t.Errorf("Expected method GET, got %s", collected[0].Method)
	}
	
	if collected[0].Path != "/test" {
		t.Errorf("Expected path /test, got %s", collected[0].Path)
	}
	
	// 测试按路径过滤
	pathMetrics := collector.GetMetricsByPath("/test")
	if len(pathMetrics) != 1 {
		t.Errorf("Expected 1 metric for path /test, got %d", len(pathMetrics))
	}
	
	// 测试按方法过滤
	methodMetrics := collector.GetMetricsByMethod("GET")
	if len(methodMetrics) != 1 {
		t.Errorf("Expected 1 metric for method GET, got %d", len(methodMetrics))
	}
	
	// 测试按状态码过滤
	statusMetrics := collector.GetMetricsByStatusCode(200)
	if len(statusMetrics) != 1 {
		t.Errorf("Expected 1 metric for status 200, got %d", len(statusMetrics))
	}
}

func TestPerformanceOptimizer(t *testing.T) {
	monitor := NewPerformanceMonitor()
	optimizer := NewPerformanceOptimizer(monitor)
	
	// 测试执行所有优化
	ctx := context.Background()
	results, err := optimizer.Optimize(ctx)
	if err != nil {
		t.Errorf("Expected no error on optimize, got %v", err)
	}
	
	// 应该有4个默认优化器
	if len(results) != 4 {
		t.Errorf("Expected 4 optimization results, got %d", len(results))
	}
	
	// 测试按类型优化
	result, err := optimizer.OptimizeByType(ctx, OptimizationTypeMemory)
	if err != nil {
		t.Errorf("Expected no error on memory optimization, got %v", err)
	}
	
	if result.Type != OptimizationTypeMemory {
		t.Errorf("Expected memory optimization type, got %s", result.Type)
	}
	
	// 测试不存在的优化类型
	_, err = optimizer.OptimizeByType(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent optimization type")
	}
}

func TestAutoOptimizer(t *testing.T) {
	monitor := NewPerformanceMonitor()
	optimizer := NewPerformanceOptimizer(monitor)
	autoOptimizer := NewAutoOptimizer(optimizer, 100*time.Millisecond)
	
	ctx := context.Background()
	
	// 测试启动
	err := autoOptimizer.Start(ctx)
	if err != nil {
		t.Errorf("Expected no error on start, got %v", err)
	}
	
	// 等待一段时间让自动优化运行
	time.Sleep(200 * time.Millisecond)
	
	// 测试停止
	err = autoOptimizer.Stop()
	if err != nil {
		t.Errorf("Expected no error on stop, got %v", err)
	}
}

func TestOptimizationResults(t *testing.T) {
	result := &OptimizationResult{
		Type:        OptimizationTypeCache,
		Success:     true,
		Message:     "Cache optimization completed",
		Improvement: 15.5,
		Timestamp:   time.Now(),
	}
	
	if result.Type != OptimizationTypeCache {
		t.Errorf("Expected cache optimization type, got %s", result.Type)
	}
	
	if !result.Success {
		t.Error("Expected optimization to be successful")
	}
	
	if result.Improvement != 15.5 {
		t.Errorf("Expected improvement 15.5, got %f", result.Improvement)
	}
}

func BenchmarkCounterIncrement(b *testing.B) {
	counter := NewCounter("benchmark_counter", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Increment(1)
	}
}

func BenchmarkGaugeSet(b *testing.B) {
	gauge := NewGauge("benchmark_gauge", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gauge.Set(float64(i))
	}
}

func BenchmarkHistogramObserve(b *testing.B) {
	buckets := []float64{10, 50, 100, 200, 500}
	histogram := NewHistogram("benchmark_histogram", buckets, nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		histogram.Observe(float64(i % 1000))
	}
}

func BenchmarkPerformanceMonitorCollect(b *testing.B) {
	monitor := NewPerformanceMonitor()
	
	// 添加一些指标
	for i := 0; i < 100; i++ {
		counter := NewCounter("counter_"+string(rune(i)), nil)
		monitor.RegisterMetric(counter)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitor.Collect()
	}
} 