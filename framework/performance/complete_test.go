package performance

import (
	"context"
	"testing"
	"time"
)

func TestCompletePerformanceMonitoring(t *testing.T) {
	// 创建性能监控器
	monitor := NewPerformanceMonitor()

	// 启动监控
	ctx := context.Background()
	monitor.Start(ctx)
	defer monitor.Stop()

	// 创建各种监控器
	httpMonitor := NewHTTPMonitor(monitor)
	dbMonitor := NewDatabaseMonitor(monitor, 100*time.Millisecond)
	cacheMonitor := NewCacheMonitor(monitor)

	// 创建告警系统
	alertSystem := NewAlertSystem(monitor)

	// 创建系统监控器
	systemMonitor := NewSystemMonitor(monitor)
	systemMonitor.Start(ctx)
	defer systemMonitor.Stop()

	// 创建报告生成器
	reportGenerator := NewReportGenerator(monitor, httpMonitor, dbMonitor, cacheMonitor, alertSystem)

	// 测试HTTP监控
	t.Run("HTTP监控", func(t *testing.T) {
		testHTTPMonitoring(t, httpMonitor)
	})

	// 测试数据库监控
	t.Run("数据库监控", func(t *testing.T) {
		testDatabaseMonitoring(t, dbMonitor)
	})

	// 测试缓存监控
	t.Run("缓存监控", func(t *testing.T) {
		testCacheMonitoring(t, cacheMonitor)
	})

	// 测试告警系统
	t.Run("告警系统", func(t *testing.T) {
		testAlertSystem(t, alertSystem)
	})

	// 测试报告生成
	t.Run("报告生成", func(t *testing.T) {
		testReportGeneration(t, reportGenerator)
	})

	// 测试系统监控
	t.Run("系统监控", func(t *testing.T) {
		testSystemMonitoring(t, systemMonitor)
	})
}

func testHTTPMonitoring(t *testing.T, httpMonitor *HTTPMonitor) {
	// 记录一些HTTP请求
	httpMonitor.RecordRequest("GET", "/api/users", 1024)
	httpMonitor.RecordResponse("GET", "/api/users", 200, 2048, 50*time.Millisecond)

	httpMonitor.RecordRequest("POST", "/api/users", 512)
	httpMonitor.RecordResponse("POST", "/api/users", 201, 1024, 100*time.Millisecond)

	httpMonitor.RecordRequest("GET", "/api/error", 256)
	httpMonitor.RecordResponse("GET", "/api/error", 500, 128, 10*time.Millisecond)

	// 验证指标
	metrics := httpMonitor.GetMetrics()
	if metrics.requestCounter.Value().(int64) != 3 {
		t.Errorf("期望请求数为3，实际为%d", metrics.requestCounter.Value().(int64))
	}

	if metrics.errorCounter.Value().(int64) != 1 {
		t.Errorf("期望错误数为1，实际为%d", metrics.errorCounter.Value().(int64))
	}
}

func testDatabaseMonitoring(t *testing.T, dbMonitor *DatabaseMonitor) {
	// 记录一些数据库查询
	dbMonitor.RecordQuery("SELECT * FROM users", 50*time.Millisecond, true, nil)
	dbMonitor.RecordQuery("INSERT INTO users (name, email) VALUES (?, ?)", 30*time.Millisecond, true, nil)
	dbMonitor.RecordQuery("UPDATE users SET name = ? WHERE id = ?", 150*time.Millisecond, true, nil) // 慢查询
	dbMonitor.RecordQuery("DELETE FROM users WHERE id = ?", 20*time.Millisecond, false, nil)         // 错误查询

	// 记录事务
	dbMonitor.RecordTransaction(200*time.Millisecond, true)

	// 更新连接池状态
	dbMonitor.UpdateConnectionPool(5, 10, 15)

	// 验证指标
	metrics := dbMonitor.GetMetrics()
	if metrics.queryCounter.Value().(int64) != 4 {
		t.Errorf("期望查询数为4，实际为%d", metrics.queryCounter.Value().(int64))
	}

	if metrics.slowQueryCounter.Value().(int64) != 1 {
		t.Errorf("期望慢查询数为1，实际为%d", metrics.slowQueryCounter.Value().(int64))
	}

	if metrics.errorCounter.Value().(int64) != 1 {
		t.Errorf("期望错误查询数为1，实际为%d", metrics.errorCounter.Value().(int64))
	}

	// 验证慢查询
	slowQueries := dbMonitor.GetSlowQueries()
	if len(slowQueries) != 1 {
		t.Errorf("期望慢查询数量为1，实际为%d", len(slowQueries))
	}

	// 验证错误查询
	errorQueries := dbMonitor.GetErrorQueries()
	if len(errorQueries) != 1 {
		t.Errorf("期望错误查询数量为1，实际为%d", len(errorQueries))
	}

	// 验证查询类型分布
	distribution := dbMonitor.GetQueryTypeDistribution()
	if distribution["SELECT"] != 1 {
		t.Errorf("期望SELECT查询数为1，实际为%d", distribution["SELECT"])
	}
	if distribution["INSERT"] != 1 {
		t.Errorf("期望INSERT查询数为1，实际为%d", distribution["INSERT"])
	}
	if distribution["UPDATE"] != 1 {
		t.Errorf("期望UPDATE查询数为1，实际为%d", distribution["UPDATE"])
	}
	if distribution["DELETE"] != 1 {
		t.Errorf("期望DELETE查询数为1，实际为%d", distribution["DELETE"])
	}
}

func testCacheMonitoring(t *testing.T, cacheMonitor *CacheMonitor) {
	// 记录一些缓存操作
	cacheMonitor.RecordGet("user:1", 1*time.Microsecond, true, nil)  // 命中
	cacheMonitor.RecordGet("user:2", 2*time.Microsecond, false, nil) // 未命中
	cacheMonitor.RecordGet("user:3", 1*time.Microsecond, true, nil)  // 命中

	cacheMonitor.RecordSet("user:1", 5*time.Microsecond, nil)
	cacheMonitor.RecordSet("user:2", 3*time.Microsecond, nil)

	cacheMonitor.RecordDelete("user:1", 1*time.Microsecond, nil)

	// 记录驱逐
	cacheMonitor.RecordEviction("old:key")

	// 更新存储指标
	cacheMonitor.UpdateStorageMetrics(1000, 1024*1024*10)

	// 验证指标
	metrics := cacheMonitor.GetMetrics()
	if metrics.getCounter.Value().(int64) != 3 {
		t.Errorf("期望GET操作数为3，实际为%d", metrics.getCounter.Value().(int64))
	}

	if metrics.setCounter.Value().(int64) != 2 {
		t.Errorf("期望SET操作数为2，实际为%d", metrics.setCounter.Value().(int64))
	}

	if metrics.deleteCounter.Value().(int64) != 1 {
		t.Errorf("期望DELETE操作数为1，实际为%d", metrics.deleteCounter.Value().(int64))
	}

	if metrics.hitCounter.Value().(int64) != 2 {
		t.Errorf("期望命中数为2，实际为%d", metrics.hitCounter.Value().(int64))
	}

	if metrics.missCounter.Value().(int64) != 1 {
		t.Errorf("期望未命中数为1，实际为%d", metrics.missCounter.Value().(int64))
	}

	if metrics.evictionCount.Value().(int64) != 1 {
		t.Errorf("期望驱逐数为1，实际为%d", metrics.evictionCount.Value().(int64))
	}

	// 验证命中率
	hitRate := cacheMonitor.GetHitRate()
	expectedHitRate := 66.67 // 2/3 * 100
	if hitRate < expectedHitRate-1 || hitRate > expectedHitRate+1 {
		t.Errorf("期望命中率为%.2f%%，实际为%.2f%%", expectedHitRate, hitRate)
	}

	// 验证操作分布
	distribution := cacheMonitor.GetOperationDistribution()
	if distribution["GET"] != 3 {
		t.Errorf("期望GET操作数为3，实际为%d", distribution["GET"])
	}
	if distribution["SET"] != 2 {
		t.Errorf("期望SET操作数为2，实际为%d", distribution["SET"])
	}
	if distribution["DELETE"] != 1 {
		t.Errorf("期望DELETE操作数为1，实际为%d", distribution["DELETE"])
	}
}

func testAlertSystem(t *testing.T, alertSystem *AlertSystem) {
	// 添加告警规则
	cpuRule := &AlertRule{
		ID:          "cpu_high",
		Name:        "CPU使用率过高",
		Description: "CPU使用率超过80%",
		MetricName:  "cpu_usage",
		Condition:   ">",
		Threshold:   80.0,
		Level:       AlertLevelWarning,
		Enabled:     true,
		Actions:     []string{"log"},
	}

	err := alertSystem.AddRule(cpuRule)
	if err != nil {
		t.Errorf("添加告警规则失败: %v", err)
	}

	// 验证规则
	rule, err := alertSystem.GetRule("cpu_high")
	if err != nil {
		t.Errorf("获取告警规则失败: %v", err)
	}

	if rule.Name != "CPU使用率过高" {
		t.Errorf("期望规则名称为'CPU使用率过高'，实际为'%s'", rule.Name)
	}

	// 启动告警系统
	ctx := context.Background()
	err = alertSystem.Start(ctx)
	if err != nil {
		t.Errorf("启动告警系统失败: %v", err)
	}
	defer alertSystem.Stop()

	// 验证规则列表
	rules := alertSystem.GetRules()
	if len(rules) != 1 {
		t.Errorf("期望规则数量为1，实际为%d", len(rules))
	}
}

func testReportGeneration(t *testing.T, reportGenerator *ReportGenerator) {
	// 生成报告
	period := ReportPeriod{
		Start:    time.Now().Add(-1 * time.Hour),
		End:      time.Now(),
		Duration: time.Hour,
	}

	report, err := reportGenerator.GenerateReport(ReportTypeSummary, period)
	if err != nil {
		t.Errorf("生成报告失败: %v", err)
	}

	if report.ID == "" {
		t.Error("报告ID不能为空")
	}

	if report.Type != ReportTypeSummary {
		t.Errorf("期望报告类型为summary，实际为%s", report.Type)
	}

	// 验证摘要数据
	summary := report.Summary
	if summary.TotalRequests < 0 {
		t.Error("总请求数不能为负数")
	}

	if summary.ErrorRate < 0 || summary.ErrorRate > 100 {
		t.Error("错误率必须在0-100之间")
	}

	// 导出报告
	data, err := reportGenerator.ExportReport(report, "json")
	if err != nil {
		t.Errorf("导出报告失败: %v", err)
	}

	if len(data) == 0 {
		t.Error("导出的报告数据不能为空")
	}

	// 导出文本格式
	textData, err := reportGenerator.ExportReport(report, "text")
	if err != nil {
		t.Errorf("导出文本报告失败: %v", err)
	}

	if len(textData) == 0 {
		t.Error("导出的文本报告数据不能为空")
	}
}

func testSystemMonitoring(t *testing.T, systemMonitor *SystemMonitor) {
	// 等待系统监控收集一些数据
	time.Sleep(2 * time.Second)

	// 获取系统指标
	systemMetrics := systemMonitor.GetSystemMetrics()

	// 验证基本指标存在
	if len(systemMetrics) == 0 {
		t.Error("系统指标不能为空")
	}

	// 验证CPU指标
	if cpuMetric, exists := systemMetrics["cpu_usage"]; exists {
		cpuValue := cpuMetric.(float64)
		if cpuValue < 0 || cpuValue > 100 {
			t.Errorf("CPU使用率必须在0-100之间，实际为%.2f", cpuValue)
		}
	}

	// 验证内存指标
	if memMetric, exists := systemMetrics["memory_usage_percent"]; exists {
		memValue := memMetric.(float64)
		if memValue < 0 || memValue > 100 {
			t.Errorf("内存使用率必须在0-100之间，实际为%.2f", memValue)
		}
	}
}

func TestCompletePerformanceOptimizer(t *testing.T) {
	// 创建性能监控器
	monitor := NewPerformanceMonitor()

	// 创建性能优化器
	optimizer := NewPerformanceOptimizer(monitor)

	// 执行优化
	ctx := context.Background()
	results, err := optimizer.Optimize(ctx)
	if err != nil {
		t.Errorf("执行优化失败: %v", err)
	}

	// 验证优化结果
	if len(results) == 0 {
		t.Error("优化结果不能为空")
	}

	for _, result := range results {
		if result.Type == "" {
			t.Error("优化类型不能为空")
		}

		if result.Message == "" {
			t.Error("优化消息不能为空")
		}

		if result.Improvement < 0 {
			t.Error("改进百分比不能为负数")
		}
	}
}

func TestCompleteAutoOptimizer(t *testing.T) {
	// 创建性能监控器
	monitor := NewPerformanceMonitor()

	// 创建性能优化器
	optimizer := NewPerformanceOptimizer(monitor)

	// 创建自动优化器
	autoOptimizer := NewAutoOptimizer(optimizer, 1*time.Second)

	// 启动自动优化器
	ctx := context.Background()
	err := autoOptimizer.Start(ctx)
	if err != nil {
		t.Errorf("启动自动优化器失败: %v", err)
	}

	// 等待一段时间
	time.Sleep(2 * time.Second)

	// 停止自动优化器
	err = autoOptimizer.Stop()
	if err != nil {
		t.Errorf("停止自动优化器失败: %v", err)
	}
}

func BenchmarkPerformanceMonitoring(b *testing.B) {
	// 创建性能监控器
	monitor := NewPerformanceMonitor()

	// 创建各种监控器
	httpMonitor := NewHTTPMonitor(monitor)
	dbMonitor := NewDatabaseMonitor(monitor, 100*time.Millisecond)
	cacheMonitor := NewCacheMonitor(monitor)

	b.ResetTimer()

	b.Run("HTTP监控", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			httpMonitor.RecordRequest("GET", "/api/test", 1024)
			httpMonitor.RecordResponse("GET", "/api/test", 200, 2048, 50*time.Millisecond)
		}
	})

	b.Run("数据库监控", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dbMonitor.RecordQuery("SELECT * FROM test", 50*time.Millisecond, true, nil)
		}
	})

	b.Run("缓存监控", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cacheMonitor.RecordGet("test:key", 1*time.Microsecond, true, nil)
		}
	})
}
