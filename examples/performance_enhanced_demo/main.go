package main

import (
	"context"
	"fmt"
	"laravel-go/framework/errors"
	"laravel-go/framework/performance"
	"log"
	"net/http"
	"time"
)

// 自定义错误类型
var (
	ErrSimulatedError = errors.New("simulated error for testing")
	ErrHighLoad       = errors.New("high load detected")
	ErrResourceExhausted = errors.New("resource exhausted")
)

// EnhancedHTTPMonitor 增强的HTTP监控器
type EnhancedHTTPMonitor struct {
	*performance.HTTPMonitor
	errorHandler errors.ErrorHandler
	errorRate    float64
}

// NewEnhancedHTTPMonitor 创建增强的HTTP监控器
func NewEnhancedHTTPMonitor(monitor performance.Monitor, errorHandler errors.ErrorHandler) *EnhancedHTTPMonitor {
	return &EnhancedHTTPMonitor{
		HTTPMonitor:  performance.NewHTTPMonitor(monitor),
		errorHandler: errorHandler,
		errorRate:    0.1, // 10% 错误率
	}
}

// RecordRequestWithErrorHandling 记录请求（带错误处理）
func (ehm *EnhancedHTTPMonitor) RecordRequestWithErrorHandling(method, path string, size int64) {
	defer func() {
		if r := recover(); r != nil {
			if ehm.errorHandler != nil {
				err := errors.New(fmt.Sprintf("HTTP monitor panic: %v", r))
				ehm.errorHandler.Handle(err)
			}
		}
	}()

	ehm.RecordRequest(method, path, size)
}

// RecordResponseWithErrorHandling 记录响应（带错误处理）
func (ehm *EnhancedHTTPMonitor) RecordResponseWithErrorHandling(method, path string, statusCode int, size int64, duration time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			if ehm.errorHandler != nil {
				err := errors.New(fmt.Sprintf("HTTP monitor panic: %v", r))
				ehm.errorHandler.Handle(err)
			}
		}
	}()

	// 模拟错误率
	if time.Now().UnixNano()%100 < int64(ehm.errorRate*100) {
		statusCode = 500
		ehm.RecordError(method, path)
	}

	ehm.RecordResponse(method, path, statusCode, size, duration)
}

// EnhancedDatabaseMonitor 增强的数据库监控器
type EnhancedDatabaseMonitor struct {
	*performance.DatabaseMonitor
	errorHandler errors.ErrorHandler
	timeoutRate  float64
}

// NewEnhancedDatabaseMonitor 创建增强的数据库监控器
func NewEnhancedDatabaseMonitor(monitor performance.Monitor, checkInterval time.Duration, errorHandler errors.ErrorHandler) *EnhancedDatabaseMonitor {
	return &EnhancedDatabaseMonitor{
		DatabaseMonitor: performance.NewDatabaseMonitor(monitor, checkInterval),
		errorHandler:    errorHandler,
		timeoutRate:     0.05, // 5% 超时率
	}
}

// RecordQueryWithErrorHandling 记录查询（带错误处理）
func (edm *EnhancedDatabaseMonitor) RecordQueryWithErrorHandling(query string, duration time.Duration, success bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			if edm.errorHandler != nil {
				panicErr := errors.New(fmt.Sprintf("Database monitor panic: %v", r))
				edm.errorHandler.Handle(panicErr)
			}
		}
	}()

	// 模拟超时
	if time.Now().UnixNano()%100 < int64(edm.timeoutRate*100) {
		success = false
		err = errors.Wrap(ErrResourceExhausted, "database query timeout")
	}

	edm.RecordQuery(query, duration, success, err)
}

// EnhancedCacheMonitor 增强的缓存监控器
type EnhancedCacheMonitor struct {
	*performance.CacheMonitor
	errorHandler errors.ErrorHandler
	unavailable  bool
}

// NewEnhancedCacheMonitor 创建增强的缓存监控器
func NewEnhancedCacheMonitor(monitor performance.Monitor, errorHandler errors.ErrorHandler) *EnhancedCacheMonitor {
	return &EnhancedCacheMonitor{
		CacheMonitor: performance.NewCacheMonitor(monitor),
		errorHandler: errorHandler,
		unavailable:  false,
	}
}

// RecordGetWithErrorHandling 记录GET操作（带错误处理）
func (ecm *EnhancedCacheMonitor) RecordGetWithErrorHandling(key string, duration time.Duration, hit bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			if ecm.errorHandler != nil {
				panicErr := errors.New(fmt.Sprintf("Cache monitor panic: %v", r))
				ecm.errorHandler.Handle(panicErr)
			}
		}
	}()

	// 模拟缓存服务不可用
	if ecm.unavailable {
		hit = false
		err = errors.Wrap(ErrResourceExhausted, "cache service unavailable")
	}

	ecm.RecordGet(key, duration, hit, err)
}

// SetUnavailable 设置缓存不可用状态
func (ecm *EnhancedCacheMonitor) SetUnavailable(unavailable bool) {
	ecm.unavailable = unavailable
}

// EnhancedAlertSystem 增强的告警系统
type EnhancedAlertSystem struct {
	*performance.AlertSystem
	errorHandler errors.ErrorHandler
}

// NewEnhancedAlertSystem 创建增强的告警系统
func NewEnhancedAlertSystem(monitor performance.Monitor, errorHandler errors.ErrorHandler) *EnhancedAlertSystem {
	return &EnhancedAlertSystem{
		AlertSystem:  performance.NewAlertSystem(monitor),
		errorHandler: errorHandler,
	}
}

// AddRuleWithErrorHandling 添加告警规则（带错误处理）
func (eas *EnhancedAlertSystem) AddRuleWithErrorHandling(rule *performance.AlertRule) error {
	defer func() {
		if r := recover(); r != nil {
			if eas.errorHandler != nil {
				err := errors.New(fmt.Sprintf("Alert system panic: %v", r))
				eas.errorHandler.Handle(err)
			}
		}
	}()

	return eas.AddRule(rule)
}

// CustomLogger 自定义日志器
type CustomLogger struct{}

func (l *CustomLogger) Error(message string, context map[string]interface{}) {
	log.Printf("[ERROR] %s: %+v", message, context)
}

func (l *CustomLogger) Warning(message string, context map[string]interface{}) {
	log.Printf("[WARN] %s: %+v", message, context)
}

func (l *CustomLogger) Info(message string, context map[string]interface{}) {
	log.Printf("[INFO] %s: %+v", message, context)
}

func (l *CustomLogger) Debug(message string, context map[string]interface{}) {
	log.Printf("[DEBUG] %s: %+v", message, context)
}

// addEnhancedAlertRules 添加增强的告警规则
func addEnhancedAlertRules(alertSystem *EnhancedAlertSystem) {
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
	alertSystem.AddRuleWithErrorHandling(cpuRule)

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
	alertSystem.AddRuleWithErrorHandling(memoryRule)

	// 错误率告警（降低阈值）
	errorRule := &performance.AlertRule{
		ID:          "error_rate_high",
		Name:        "错误率过高",
		Description: "HTTP错误率超过3%",
		MetricName:  "http_errors_total",
		Condition:   ">",
		Threshold:   3.0,
		Level:       performance.AlertLevelCritical,
		Enabled:     true,
		Actions:     []string{"log", "email", "webhook"},
	}
	alertSystem.AddRuleWithErrorHandling(errorRule)

	// 响应时间告警
	responseTimeRule := &performance.AlertRule{
		ID:          "response_time_high",
		Name:        "响应时间过长",
		Description: "平均响应时间超过500ms",
		MetricName:  "http_response_time",
		Condition:   ">",
		Threshold:   500.0,
		Level:       performance.AlertLevelWarning,
		Enabled:     true,
		Actions:     []string{"log"},
	}
	alertSystem.AddRuleWithErrorHandling(responseTimeRule)
}

// simulateEnhancedApplication 模拟增强的应用程序运行
func simulateEnhancedApplication(httpMonitor *EnhancedHTTPMonitor, dbMonitor *EnhancedDatabaseMonitor, cacheMonitor *EnhancedCacheMonitor) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	cycle := 0
	for range ticker.C {
		cycle++

		// 模拟HTTP请求
		simulateEnhancedHTTPRequests(httpMonitor)

		// 模拟数据库查询
		simulateEnhancedDatabaseQueries(dbMonitor)

		// 模拟缓存操作
		simulateEnhancedCacheOperations(cacheMonitor)

		// 每10个周期切换缓存可用性
		if cycle%10 == 0 {
			cacheMonitor.SetUnavailable(!cacheMonitor.unavailable)
			log.Printf("Cache availability changed to: %v", !cacheMonitor.unavailable)
		}
	}
}

// simulateEnhancedHTTPRequests 模拟增强的HTTP请求
func simulateEnhancedHTTPRequests(httpMonitor *EnhancedHTTPMonitor) {
	endpoints := []string{"/api/users", "/api/products", "/api/orders", "/api/health"}
	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for i := 0; i < 3; i++ {
		endpoint := endpoints[i%len(endpoints)]
		method := methods[i%len(methods)]

		// 记录请求
		httpMonitor.RecordRequestWithErrorHandling(method, endpoint, 1024)

		// 模拟处理时间
		time.Sleep(time.Duration(10+i*5) * time.Millisecond)

		// 记录响应
		statusCode := 200
		responseSize := int64(2048)
		duration := time.Duration(10+i*5) * time.Millisecond

		httpMonitor.RecordResponseWithErrorHandling(method, endpoint, statusCode, responseSize, duration)
	}
}

// simulateEnhancedDatabaseQueries 模拟增强的数据库查询
func simulateEnhancedDatabaseQueries(dbMonitor *EnhancedDatabaseMonitor) {
	queries := []string{
		"SELECT * FROM users WHERE id = ?",
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"UPDATE users SET name = ? WHERE id = ?",
		"DELETE FROM users WHERE id = ?",
	}

	for i := 0; i < 2; i++ {
		query := queries[i%len(queries)]

		// 模拟查询时间
		queryTime := time.Duration(5+i*10) * time.Millisecond
		time.Sleep(queryTime)

		// 偶尔产生错误
		success := i%15 != 0
		var err error
		if !success {
			err = errors.Wrap(ErrSimulatedError, "database query failed")
		}

		dbMonitor.RecordQueryWithErrorHandling(query, queryTime, success, err)
	}

	// 更新连接池状态
	dbMonitor.UpdateConnectionPool(5, 10, 15)
}

// simulateEnhancedCacheOperations 模拟增强的缓存操作
func simulateEnhancedCacheOperations(cacheMonitor *EnhancedCacheMonitor) {
	keys := []string{"user:1", "product:123", "order:456", "config:app"}

	for i := 0; i < 3; i++ {
		key := keys[i%len(keys)]

		// 模拟GET操作
		start := time.Now()
		time.Sleep(time.Duration(1+i) * time.Microsecond)
		hit := i%3 != 0 // 2/3的命中率
		var err error
		if !hit {
			err = errors.New("cache miss")
		}
		cacheMonitor.RecordGetWithErrorHandling(key, time.Since(start), hit, err)

		// 模拟SET操作
		start = time.Now()
		time.Sleep(time.Duration(2+i) * time.Microsecond)
		cacheMonitor.RecordSet(key, time.Since(start), nil)
	}

	// 更新存储指标
	cacheMonitor.UpdateStorageMetrics(1000, 1024*1024*10) // 1000个条目，10MB
}

// startEnhancedMonitoringServer 启动增强的监控服务器
func startEnhancedMonitoringServer(monitor performance.Monitor, httpMonitor *EnhancedHTTPMonitor, dbMonitor *EnhancedDatabaseMonitor, cacheMonitor *EnhancedCacheMonitor, alertSystem *EnhancedAlertSystem, reportGenerator *performance.ReportGenerator) {
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
			"uptime": "%v",
			"error_handling": "enhanced"
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

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "healthy", "timestamp": "%s", "error_handling": "enhanced"}`, time.Now().Format(time.RFC3339))
	})

	port := ":8090"
	log.Printf("增强性能监控服务器启动在端口 %s", port)
	log.Printf("可用端点:")
	log.Printf("  GET /metrics - 获取所有指标")
	log.Printf("  GET /status - 获取系统状态")
	log.Printf("  GET /alerts - 获取活跃告警")
	log.Printf("  GET /health - 健康检查")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

// generateEnhancedPeriodicReports 生成增强的定期报告
func generateEnhancedPeriodicReports(reportGenerator *performance.ReportGenerator) {
	ticker := time.NewTicker(3 * time.Minute) // 缩短报告间隔
	defer ticker.Stop()

	for range ticker.C {
		period := performance.ReportPeriod{
			Start:    time.Now().Add(-3 * time.Minute),
			End:      time.Now(),
			Duration: 3 * time.Minute,
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

		log.Printf("=== 增强性能监控报告 ===\n%s", string(data))

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

func main() {
	// 创建错误处理器
	logger := &CustomLogger{}
	errorHandler := errors.NewDefaultErrorHandler(logger)

	// 创建性能监控器
	monitor := performance.NewPerformanceMonitor()

	// 启动监控
	ctx := context.Background()
	monitor.Start(ctx)
	defer monitor.Stop()

	// 创建增强的监控器
	httpMonitor := NewEnhancedHTTPMonitor(monitor, errorHandler)
	dbMonitor := NewEnhancedDatabaseMonitor(monitor, 100*time.Millisecond, errorHandler)
	cacheMonitor := NewEnhancedCacheMonitor(monitor, errorHandler)

	// 创建增强的告警系统
	alertSystem := NewEnhancedAlertSystem(monitor, errorHandler)

	// 添加告警规则
	addEnhancedAlertRules(alertSystem)

	// 启动告警系统
	alertSystem.Start(ctx)
	defer alertSystem.Stop()

	// 创建系统监控器
	systemMonitor := performance.NewSystemMonitor(monitor)
	systemMonitor.Start(ctx)
	defer systemMonitor.Stop()

	// 创建报告生成器
	reportGenerator := performance.NewReportGenerator(monitor, httpMonitor.HTTPMonitor, dbMonitor.DatabaseMonitor, cacheMonitor.CacheMonitor, alertSystem.AlertSystem)

	// 模拟应用程序运行
	go simulateEnhancedApplication(httpMonitor, dbMonitor, cacheMonitor)

	// 启动HTTP服务器提供监控接口
	go startEnhancedMonitoringServer(monitor, httpMonitor, dbMonitor, cacheMonitor, alertSystem, reportGenerator)

	// 定期生成报告
	go generateEnhancedPeriodicReports(reportGenerator)

	// 保持运行
	select {}
} 