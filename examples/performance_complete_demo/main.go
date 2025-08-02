package main

import (
	"context"
	"fmt"
	"laravel-go/framework/performance"
	"log"
	"net/http"
	"time"
)

func main() {
	// 创建性能监控器
	monitor := performance.NewPerformanceMonitor()

	// 启动监控
	ctx := context.Background()
	monitor.Start(ctx)
	defer monitor.Stop()

	// 创建各种监控器
	httpMonitor := performance.NewHTTPMonitor(monitor)
	dbMonitor := performance.NewDatabaseMonitor(monitor, 100*time.Millisecond)
	cacheMonitor := performance.NewCacheMonitor(monitor)

	// 创建告警系统
	alertSystem := performance.NewAlertSystem(monitor)

	// 添加告警规则
	addAlertRules(alertSystem)

	// 启动告警系统
	alertSystem.Start(ctx)
	defer alertSystem.Stop()

	// 创建系统监控器
	systemMonitor := performance.NewSystemMonitor(monitor)
	systemMonitor.Start(ctx)
	defer systemMonitor.Stop()

	// 创建报告生成器
	reportGenerator := performance.NewReportGenerator(monitor, httpMonitor, dbMonitor, cacheMonitor, alertSystem)

	// 模拟应用程序运行
	go simulateApplication(httpMonitor, dbMonitor, cacheMonitor)

	// 启动HTTP服务器提供监控接口
	go startMonitoringServer(monitor, httpMonitor, dbMonitor, cacheMonitor, alertSystem, reportGenerator)

	// 定期生成报告
	go generatePeriodicReports(reportGenerator)

	// 保持运行
	select {}
}

// addAlertRules 添加告警规则
func addAlertRules(alertSystem *performance.AlertSystem) {
	// CPU使用率告警
	cpuRule := &performance.AlertRule{
		ID:          "cpu_high",
		Name:        "CPU使用率过高",
		Description: "CPU使用率超过80%",
		MetricName:  "cpu_usage",
		Condition:   ">",
		Threshold:   80.0,
		Level:       performance.AlertLevelWarning,
		Enabled:     true,
		Actions:     []string{"log", "email"},
	}
	alertSystem.AddRule(cpuRule)

	// 内存使用率告警
	memoryRule := &performance.AlertRule{
		ID:          "memory_high",
		Name:        "内存使用率过高",
		Description: "内存使用率超过85%",
		MetricName:  "memory_usage_percent",
		Condition:   ">",
		Threshold:   85.0,
		Level:       performance.AlertLevelError,
		Enabled:     true,
		Actions:     []string{"log", "webhook"},
	}
	alertSystem.AddRule(memoryRule)

	// 错误率告警
	errorRule := &performance.AlertRule{
		ID:          "error_rate_high",
		Name:        "错误率过高",
		Description: "HTTP错误率超过5%",
		MetricName:  "http_errors_total",
		Condition:   ">",
		Threshold:   5.0,
		Level:       performance.AlertLevelCritical,
		Enabled:     true,
		Actions:     []string{"log", "email", "webhook"},
	}
	alertSystem.AddRule(errorRule)
}

// simulateApplication 模拟应用程序运行
func simulateApplication(httpMonitor *performance.HTTPMonitor, dbMonitor *performance.DatabaseMonitor, cacheMonitor *performance.CacheMonitor) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 模拟HTTP请求
		simulateHTTPRequests(httpMonitor)

		// 模拟数据库查询
		simulateDatabaseQueries(dbMonitor)

		// 模拟缓存操作
		simulateCacheOperations(cacheMonitor)
	}
}

// simulateHTTPRequests 模拟HTTP请求
func simulateHTTPRequests(httpMonitor *performance.HTTPMonitor) {
	endpoints := []string{"/api/users", "/api/products", "/api/orders", "/api/health"}
	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for i := 0; i < 5; i++ {
		endpoint := endpoints[i%len(endpoints)]
		method := methods[i%len(methods)]

		// 记录请求
		httpMonitor.RecordRequest(method, endpoint, 1024)

		// 模拟处理时间
		time.Sleep(time.Duration(10+i*5) * time.Millisecond)

		// 记录响应
		statusCode := 200
		if i%10 == 0 { // 偶尔产生错误
			statusCode = 500
		}
		httpMonitor.RecordResponse(method, endpoint, statusCode, 2048, time.Duration(10+i*5)*time.Millisecond)

		if statusCode >= 400 {
			httpMonitor.RecordError(method, endpoint)
		}
	}
}

// simulateDatabaseQueries 模拟数据库查询
func simulateDatabaseQueries(dbMonitor *performance.DatabaseMonitor) {
	queries := []string{
		"SELECT * FROM users WHERE id = ?",
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"UPDATE users SET name = ? WHERE id = ?",
		"DELETE FROM users WHERE id = ?",
	}

	for i := 0; i < 3; i++ {
		query := queries[i%len(queries)]

		// 模拟查询时间
		queryTime := time.Duration(5+i*10) * time.Millisecond
		time.Sleep(queryTime)

		// 偶尔产生错误
		success := i%20 != 0
		var err error
		if !success {
			err = fmt.Errorf("database error")
		}

		dbMonitor.RecordQuery(query, queryTime, success, err)
	}

	// 更新连接池状态
	dbMonitor.UpdateConnectionPool(5, 10, 15)
}

// simulateCacheOperations 模拟缓存操作
func simulateCacheOperations(cacheMonitor *performance.CacheMonitor) {
	keys := []string{"user:1", "product:123", "order:456", "config:app"}

	for i := 0; i < 4; i++ {
		key := keys[i%len(keys)]

		// 模拟GET操作
		start := time.Now()
		time.Sleep(time.Duration(1+i) * time.Microsecond)
		hit := i%3 != 0 // 2/3的命中率
		cacheMonitor.RecordGet(key, time.Since(start), hit, nil)

		// 模拟SET操作
		start = time.Now()
		time.Sleep(time.Duration(2+i) * time.Microsecond)
		cacheMonitor.RecordSet(key, time.Since(start), nil)
	}

	// 更新存储指标
	cacheMonitor.UpdateStorageMetrics(1000, 1024*1024*10) // 1000个条目，10MB
}

// startMonitoringServer 启动监控服务器
func startMonitoringServer(monitor performance.Monitor, httpMonitor *performance.HTTPMonitor, dbMonitor *performance.DatabaseMonitor, cacheMonitor *performance.CacheMonitor, alertSystem *performance.AlertSystem, reportGenerator *performance.ReportGenerator) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := monitor.GetAllMetrics()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\n")
		first := true
		for name, metric := range metrics {
			if !first {
				fmt.Fprintf(w, ",\n")
			}
			fmt.Fprintf(w, "  \"%s\": %v", name, metric.Value())
			first = false
		}
		fmt.Fprintf(w, "\n}\n")
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "running",
			"timestamp": "%s",
			"uptime": "%v"
		}`, time.Now().Format(time.RFC3339), time.Since(time.Now()))
	})

	http.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		alerts := alertSystem.GetActiveAlerts()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\n  \"active_alerts\": %d,\n  \"alerts\": [\n", len(alerts))
		for i, alert := range alerts {
			if i > 0 {
				fmt.Fprintf(w, ",\n")
			}
			fmt.Fprintf(w, "    {\n")
			fmt.Fprintf(w, "      \"id\": \"%s\",\n", alert.ID)
			fmt.Fprintf(w, "      \"rule_name\": \"%s\",\n", alert.RuleName)
			fmt.Fprintf(w, "      \"level\": \"%s\",\n", alert.Level)
			fmt.Fprintf(w, "      \"message\": \"%s\",\n", alert.Message)
			fmt.Fprintf(w, "      \"timestamp\": \"%s\"\n", alert.Timestamp.Format(time.RFC3339))
			fmt.Fprintf(w, "    }")
		}
		fmt.Fprintf(w, "\n  ]\n}\n")
	})

	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		reportType := r.URL.Query().Get("type")
		if reportType == "" {
			reportType = "summary"
		}

		period := performance.ReportPeriod{
			Start:    time.Now().Add(-1 * time.Hour),
			End:      time.Now(),
			Duration: time.Hour,
		}

		report, err := reportGenerator.GenerateReport(performance.ReportType(reportType), period)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		format := r.URL.Query().Get("format")
		if format == "" {
			format = "json"
		}

		data, err := reportGenerator.ExportReport(report, format)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		contentType := "application/json"
		if format == "text" {
			contentType = "text/plain"
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "healthy", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
	})

	port := ":8088"
	log.Printf("监控服务器启动在端口 %s", port)
	log.Printf("可用端点:")
	log.Printf("  GET /metrics - 获取所有指标")
	log.Printf("  GET /status - 获取系统状态")
	log.Printf("  GET /alerts - 获取活跃告警")
	log.Printf("  GET /report?type=summary&format=json - 生成性能报告")
	log.Printf("  GET /health - 健康检查")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

// generatePeriodicReports 定期生成报告
func generatePeriodicReports(reportGenerator *performance.ReportGenerator) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		period := performance.ReportPeriod{
			Start:    time.Now().Add(-5 * time.Minute),
			End:      time.Now(),
			Duration: 5 * time.Minute,
		}

		// 生成摘要报告
		report, err := reportGenerator.GenerateReport(performance.ReportTypeSummary, period)
		if err != nil {
			log.Printf("生成报告失败: %v", err)
			continue
		}

		// 导出为文本格式
		data, err := reportGenerator.ExportReport(report, "text")
		if err != nil {
			log.Printf("导出报告失败: %v", err)
			continue
		}

		log.Printf("=== 性能监控报告 ===\n%s", string(data))

		// 如果有建议，打印建议
		if len(report.Recommendations) > 0 {
			log.Printf("=== 优化建议 ===")
			for i, rec := range report.Recommendations {
				log.Printf("%d. %s (%s优先级) - 预期改进: %.1f%%",
					i+1, rec.Title, rec.Priority, rec.Impact)
			}
		}
	}
}
