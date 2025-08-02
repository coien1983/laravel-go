package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"laravel-go/framework/performance"
)

func main() {
	fmt.Println("=== Laravel-Go 性能监控系统演示 ===")

	// 创建性能监控器
	monitor := performance.NewPerformanceMonitor()
	
	// 启动监控
	ctx := context.Background()
	if err := monitor.Start(ctx); err != nil {
		log.Fatalf("启动监控失败: %v", err)
	}
	defer monitor.Stop()

	// 创建系统监控器
	systemMonitor := performance.NewSystemMonitor(monitor)
	if err := systemMonitor.Start(ctx); err != nil {
		log.Fatalf("启动系统监控失败: %v", err)
	}
	defer systemMonitor.Stop()

	// 创建HTTP监控器
	httpMonitor := performance.NewHTTPMonitor(monitor)

	// 创建性能优化器
	optimizer := performance.NewPerformanceOptimizer(monitor)

	// 创建自动优化器
	autoOptimizer := performance.NewAutoOptimizer(optimizer, 30*time.Second)
	if err := autoOptimizer.Start(ctx); err != nil {
		log.Fatalf("启动自动优化失败: %v", err)
	}
	defer autoOptimizer.Stop()

	// 演示1: 基础指标监控
	demonstrateBasicMetrics(monitor)

	// 演示2: HTTP请求监控
	demonstrateHTTPMonitoring(httpMonitor, monitor)

	// 演示3: 系统资源监控
	demonstrateSystemMonitoring(systemMonitor)

	// 演示4: 性能优化
	demonstratePerformanceOptimization(optimizer)

	// 演示5: 启动监控服务器
	demonstrateMonitoringServer(monitor, systemMonitor, optimizer)

	fmt.Println("\n演示完成！")
}

// demonstrateBasicMetrics 演示基础指标监控
func demonstrateBasicMetrics(monitor performance.Monitor) {
	fmt.Println("\n1. 基础指标监控演示")
	fmt.Println("==================")

	// 创建一些测试指标
	counter := performance.NewCounter("demo_requests", map[string]string{"service": "demo"})
	monitor.RegisterMetric(counter)

	gauge := performance.NewGauge("demo_memory_usage", map[string]string{"unit": "bytes"})
	monitor.RegisterMetric(gauge)

	buckets := []float64{10, 50, 100, 200, 500, 1000}
	histogram := performance.NewHistogram("demo_response_time", buckets, map[string]string{"unit": "ms"})
	monitor.RegisterMetric(histogram)

	// 模拟一些操作
	fmt.Println("模拟请求处理...")
	for i := 0; i < 10; i++ {
		counter.Increment(1)
		gauge.Set(float64(i * 1000))
		histogram.Observe(float64(i * 50))
		time.Sleep(100 * time.Millisecond)
	}

	// 显示指标
	fmt.Println("当前指标:")
	metrics := monitor.GetAllMetrics()
	for name, metric := range metrics {
		fmt.Printf("- %s: %v (%s)\n", name, metric.Value(), metric.Type())
	}
}

// demonstrateHTTPMonitoring 演示HTTP请求监控
func demonstrateHTTPMonitoring(httpMonitor *performance.HTTPMonitor, monitor performance.Monitor) {
	fmt.Println("\n2. HTTP请求监控演示")
	fmt.Println("==================")

	// 模拟HTTP请求
	fmt.Println("模拟HTTP请求...")
	
	// 成功请求
	httpMonitor.RecordRequest("GET", "/api/users", 150)
	time.Sleep(50 * time.Millisecond)
	httpMonitor.RecordResponse("GET", "/api/users", 200, 1024, 50*time.Millisecond)

	httpMonitor.RecordRequest("POST", "/api/users", 300)
	time.Sleep(100 * time.Millisecond)
	httpMonitor.RecordResponse("POST", "/api/users", 201, 512, 100*time.Millisecond)

	// 错误请求
	httpMonitor.RecordRequest("GET", "/api/invalid", 100)
	time.Sleep(20 * time.Millisecond)
	httpMonitor.RecordResponse("GET", "/api/invalid", 404, 200, 20*time.Millisecond)

	// 记录错误
	httpMonitor.RecordError("GET", "/api/error")

	fmt.Println("HTTP监控指标:")
	// 从监控器获取指标
	requestCounter := monitor.GetMetric("http_requests_total")
	responseCounter := monitor.GetMetric("http_responses_total")
	errorCounter := monitor.GetMetric("http_errors_total")
	activeConnections := monitor.GetMetric("http_active_connections")
	
	if requestCounter != nil {
		fmt.Printf("- 总请求数: %v\n", requestCounter.Value())
	}
	
	if responseCounter != nil {
		fmt.Printf("- 总响应数: %v\n", responseCounter.Value())
	}
	
	if errorCounter != nil {
		fmt.Printf("- 错误数: %v\n", errorCounter.Value())
	}
	
	if activeConnections != nil {
		fmt.Printf("- 活跃连接数: %v\n", activeConnections.Value())
	}
}

// demonstrateSystemMonitoring 演示系统资源监控
func demonstrateSystemMonitoring(systemMonitor *performance.SystemMonitor) {
	fmt.Println("\n3. 系统资源监控演示")
	fmt.Println("==================")

	// 等待系统监控收集一些数据
	fmt.Println("收集系统指标中...")
	time.Sleep(6 * time.Second)

	// 显示系统指标
	systemMetrics := systemMonitor.GetSystemMetrics()
	
	fmt.Println("系统指标:")
	for name, metric := range systemMetrics {
		metricMap := metric.(map[string]interface{})
		value := metricMap["value"]
		
		// 格式化显示
		switch name {
		case "cpu_usage":
			fmt.Printf("- CPU使用率: %.1f%%\n", value)
		case "memory_usage_percent":
			fmt.Printf("- 内存使用率: %.1f%%\n", value)
		case "go_goroutines":
			fmt.Printf("- Go协程数: %.0f\n", value)
		case "go_heap_alloc":
			fmt.Printf("- Go堆内存: %.0f bytes\n", value)
		}
	}
}

// demonstratePerformanceOptimization 演示性能优化
func demonstratePerformanceOptimization(optimizer *performance.PerformanceOptimizer) {
	fmt.Println("\n4. 性能优化演示")
	fmt.Println("================")

	ctx := context.Background()

	// 执行所有优化
	fmt.Println("执行性能优化分析...")
	results, err := optimizer.Optimize(ctx)
	if err != nil {
		fmt.Printf("优化执行失败: %v\n", err)
		return
	}

	fmt.Println("优化结果:")
	for _, result := range results {
		status := "成功"
		if !result.Success {
			status = "失败"
		}
		
		fmt.Printf("- %s (%s): %s", result.Type, status, result.Message)
		if result.Improvement > 0 {
			fmt.Printf(" (预计改进: %.1f%%)", result.Improvement)
		}
		fmt.Println()
	}

	// 执行特定类型的优化
	fmt.Println("\n执行内存优化...")
	memoryResult, err := optimizer.OptimizeByType(ctx, performance.OptimizationTypeMemory)
	if err != nil {
		fmt.Printf("内存优化失败: %v\n", err)
	} else {
		fmt.Printf("内存优化结果: %s\n", memoryResult.Message)
	}
}

// demonstrateMonitoringServer 演示监控服务器
func demonstrateMonitoringServer(monitor performance.Monitor, systemMonitor *performance.SystemMonitor, optimizer *performance.PerformanceOptimizer) {
	fmt.Println("\n5. 监控服务器演示")
	fmt.Println("==================")

	// 创建HTTP服务器
	mux := http.NewServeMux()

	// 指标端点
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// 收集所有指标
		metrics := monitor.GetAllMetrics()
		
		// 这里应该使用JSON编码，简化显示
		fmt.Fprintf(w, "Metrics endpoint - %d metrics available\n", len(metrics))
	})

	// 系统状态端点
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		status := map[string]interface{}{
			"timestamp":    time.Now(),
			"goroutines":   runtime.NumGoroutine(),
			"heap_alloc":   m.HeapAlloc,
			"heap_sys":     m.HeapSys,
			"heap_idle":    m.HeapIdle,
			"heap_inuse":   m.HeapInuse,
			"heap_released": m.HeapReleased,
		}
		
		fmt.Fprintf(w, "System Status - Goroutines: %d, Heap: %d bytes\n", 
			status["goroutines"], status["heap_alloc"])
	})

	// 优化端点
	mux.HandleFunc("/optimize", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		ctx := context.Background()
		results, err := optimizer.Optimize(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		fmt.Fprintf(w, "Optimization completed - %d results\n", len(results))
	})

	// 健康检查端点
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "OK - Performance monitoring system is running\n")
	})

	// 启动服务器
	port := ":8088"
	fmt.Printf("启动监控服务器在 http://localhost%s\n", port)
	fmt.Println("可用端点:")
	fmt.Printf("- 指标: http://localhost%s/metrics\n", port)
	fmt.Printf("- 状态: http://localhost%s/status\n", port)
	fmt.Printf("- 优化: http://localhost%s/optimize\n", port)
	fmt.Printf("- 健康检查: http://localhost%s/health\n", port)
	
	// 在后台启动服务器
	go func() {
		if err := http.ListenAndServe(port, mux); err != nil {
			log.Printf("服务器启动失败: %v", err)
		}
	}()

	// 等待一段时间让用户测试
	fmt.Println("\n服务器已启动，按 Ctrl+C 停止...")
	
	// 模拟一些负载
	go func() {
		for i := 0; i < 100; i++ {
			// 模拟一些操作来产生指标
			runtime.GC()
			time.Sleep(1 * time.Second)
		}
	}()
	
	// 等待用户中断
	select {
	case <-time.After(60 * time.Second):
		fmt.Println("演示时间结束")
	}
}

// 辅助函数：格式化字节数
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
} 