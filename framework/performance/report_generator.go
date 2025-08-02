package performance

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ReportType 报告类型
type ReportType string

const (
	ReportTypeSummary    ReportType = "summary"
	ReportTypeDetailed   ReportType = "detailed"
	ReportTypeTrend      ReportType = "trend"
	ReportTypeComparison ReportType = "comparison"
)

// PerformanceReport 性能报告
type PerformanceReport struct {
	ID              string                 `json:"id"`
	Type            ReportType             `json:"type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	GeneratedAt     time.Time              `json:"generated_at"`
	Period          ReportPeriod           `json:"period"`
	Summary         ReportSummary          `json:"summary"`
	Details         ReportDetails          `json:"details"`
	Recommendations []Recommendation       `json:"recommendations"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// ReportPeriod 报告周期
type ReportPeriod struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

// ReportSummary 报告摘要
type ReportSummary struct {
	TotalRequests       int64         `json:"total_requests"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	ErrorRate           float64       `json:"error_rate"`
	Throughput          float64       `json:"throughput"` // requests per second
	CPUUsage            float64       `json:"cpu_usage"`
	MemoryUsage         float64       `json:"memory_usage"`
	DatabaseQueries     int64         `json:"database_queries"`
	CacheHitRate        float64       `json:"cache_hit_rate"`
	SlowQueries         int64         `json:"slow_queries"`
	ActiveAlerts        int64         `json:"active_alerts"`
}

// ReportDetails 报告详情
type ReportDetails struct {
	HTTPMetrics     HTTPReportDetails     `json:"http_metrics"`
	DatabaseMetrics DatabaseReportDetails `json:"database_metrics"`
	CacheMetrics    CacheReportDetails    `json:"cache_metrics"`
	SystemMetrics   SystemReportDetails   `json:"system_metrics"`
	AlertMetrics    AlertReportDetails    `json:"alert_metrics"`
}

// HTTPReportDetails HTTP报告详情
type HTTPReportDetails struct {
	RequestDistribution      map[string]int64 `json:"request_distribution"`
	ResponseTimeDistribution map[string]int64 `json:"response_time_distribution"`
	StatusCodeDistribution   map[int]int64    `json:"status_code_distribution"`
	TopEndpoints             []EndpointStats  `json:"top_endpoints"`
	SlowestEndpoints         []EndpointStats  `json:"slowest_endpoints"`
	ErrorEndpoints           []EndpointStats  `json:"error_endpoints"`
}

// EndpointStats 端点统计
type EndpointStats struct {
	Path        string        `json:"path"`
	Method      string        `json:"method"`
	Count       int64         `json:"count"`
	AverageTime time.Duration `json:"average_time"`
	ErrorCount  int64         `json:"error_count"`
	ErrorRate   float64       `json:"error_rate"`
}

// DatabaseReportDetails 数据库报告详情
type DatabaseReportDetails struct {
	QueryTypeDistribution map[string]int64    `json:"query_type_distribution"`
	SlowQueries           []QueryRecord       `json:"slow_queries"`
	ErrorQueries          []QueryRecord       `json:"error_queries"`
	AverageQueryTime      time.Duration       `json:"average_query_time"`
	ConnectionPoolStats   ConnectionPoolStats `json:"connection_pool_stats"`
}

// ConnectionPoolStats 连接池统计
type ConnectionPoolStats struct {
	ActiveConnections int64   `json:"active_connections"`
	IdleConnections   int64   `json:"idle_connections"`
	TotalConnections  int64   `json:"total_connections"`
	MaxConnections    int64   `json:"max_connections"`
	Utilization       float64 `json:"utilization"`
}

// CacheReportDetails 缓存报告详情
type CacheReportDetails struct {
	HitRate               float64          `json:"hit_rate"`
	OperationDistribution map[string]int64 `json:"operation_distribution"`
	SlowOperations        []CacheOperation `json:"slow_operations"`
	ErrorOperations       []CacheOperation `json:"error_operations"`
	StorageStats          StorageStats     `json:"storage_stats"`
}

// StorageStats 存储统计
type StorageStats struct {
	ItemCount   int64   `json:"item_count"`
	MemoryUsage int64   `json:"memory_usage"`
	MaxItems    int64   `json:"max_items"`
	MaxMemory   int64   `json:"max_memory"`
	Utilization float64 `json:"utilization"`
}

// SystemReportDetails 系统报告详情
type SystemReportDetails struct {
	CPUStats       CPUStats       `json:"cpu_stats"`
	MemoryStats    MemoryStats    `json:"memory_stats"`
	DiskStats      DiskStats      `json:"disk_stats"`
	GoroutineStats GoroutineStats `json:"goroutine_stats"`
}

// CPUStats CPU统计
type CPUStats struct {
	Usage       float64 `json:"usage"`
	Cores       int     `json:"cores"`
	LoadAverage float64 `json:"load_average"`
}

// MemoryStats 内存统计
type MemoryStats struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	Usage     float64 `json:"usage"`
	HeapAlloc uint64  `json:"heap_alloc"`
	HeapSys   uint64  `json:"heap_sys"`
}

// DiskStats 磁盘统计
type DiskStats struct {
	Total uint64  `json:"total"`
	Used  uint64  `json:"used"`
	Free  uint64  `json:"free"`
	Usage float64 `json:"usage"`
	IOPS  float64 `json:"iops"`
}

// GoroutineStats 协程统计
type GoroutineStats struct {
	Count      int     `json:"count"`
	MaxCount   int     `json:"max_count"`
	GrowthRate float64 `json:"growth_rate"`
}

// AlertReportDetails 告警报告详情
type AlertReportDetails struct {
	TotalAlerts    int64            `json:"total_alerts"`
	ActiveAlerts   int64            `json:"active_alerts"`
	ResolvedAlerts int64            `json:"resolved_alerts"`
	AlertByLevel   map[string]int64 `json:"alert_by_level"`
	AlertByRule    map[string]int64 `json:"alert_by_rule"`
}

// Recommendation 优化建议
type Recommendation struct {
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"` // low, medium, high, critical
	Impact      float64  `json:"impact"`   // 预期改进百分比
	Effort      string   `json:"effort"`   // low, medium, high
	Actions     []string `json:"actions"`
}

// ReportGenerator 报告生成器
type ReportGenerator struct {
	monitor      Monitor
	httpMonitor  *HTTPMonitor
	dbMonitor    *DatabaseMonitor
	cacheMonitor *CacheMonitor
	alertSystem  *AlertSystem
}

// NewReportGenerator 创建报告生成器
func NewReportGenerator(monitor Monitor, httpMonitor *HTTPMonitor, dbMonitor *DatabaseMonitor, cacheMonitor *CacheMonitor, alertSystem *AlertSystem) *ReportGenerator {
	return &ReportGenerator{
		monitor:      monitor,
		httpMonitor:  httpMonitor,
		dbMonitor:    dbMonitor,
		cacheMonitor: cacheMonitor,
		alertSystem:  alertSystem,
	}
}

// GenerateReport 生成性能报告
func (rg *ReportGenerator) GenerateReport(reportType ReportType, period ReportPeriod) (*PerformanceReport, error) {
	report := &PerformanceReport{
		ID:          fmt.Sprintf("report_%s_%d", reportType, time.Now().Unix()),
		Type:        reportType,
		GeneratedAt: time.Now(),
		Period:      period,
		Metadata:    make(map[string]interface{}),
	}

	switch reportType {
	case ReportTypeSummary:
		report.Title = "性能监控摘要报告"
		report.Description = "应用程序性能监控摘要"
		report.Summary = rg.generateSummary(period)
		report.Recommendations = rg.generateRecommendations(report.Summary)

	case ReportTypeDetailed:
		report.Title = "性能监控详细报告"
		report.Description = "应用程序性能监控详细分析"
		report.Summary = rg.generateSummary(period)
		report.Details = rg.generateDetails(period)
		report.Recommendations = rg.generateRecommendations(report.Summary)

	case ReportTypeTrend:
		report.Title = "性能趋势分析报告"
		report.Description = "应用程序性能趋势分析"
		report.Summary = rg.generateSummary(period)
		report.Metadata["trends"] = rg.generateTrends(period)
		report.Recommendations = rg.generateRecommendations(report.Summary)

	case ReportTypeComparison:
		report.Title = "性能对比分析报告"
		report.Description = "应用程序性能对比分析"
		report.Summary = rg.generateSummary(period)
		report.Metadata["comparison"] = rg.generateComparison(period)
		report.Recommendations = rg.generateRecommendations(report.Summary)
	}

	return report, nil
}

// generateSummary 生成摘要
func (rg *ReportGenerator) generateSummary(period ReportPeriod) ReportSummary {
	summary := ReportSummary{}

	// HTTP指标
	if rg.httpMonitor != nil {
		metrics := rg.httpMonitor.GetMetrics()
		summary.TotalRequests = metrics.requestCounter.Value().(int64)
		summary.ErrorRate = rg.calculateErrorRate(metrics)
		summary.Throughput = float64(summary.TotalRequests) / period.Duration.Seconds()
		summary.AverageResponseTime = rg.calculateAverageResponseTime(metrics)
	}

	// 数据库指标
	if rg.dbMonitor != nil {
		metrics := rg.dbMonitor.GetMetrics()
		summary.DatabaseQueries = metrics.queryCounter.Value().(int64)
		summary.SlowQueries = metrics.slowQueryCounter.Value().(int64)
	}

	// 缓存指标
	if rg.cacheMonitor != nil {
		summary.CacheHitRate = rg.cacheMonitor.GetHitRate()
	}

	// 系统指标
	systemMetrics := rg.monitor.GetAllMetrics()
	if cpuMetric := systemMetrics["cpu_usage"]; cpuMetric != nil {
		summary.CPUUsage = cpuMetric.Value().(float64)
	}
	if memMetric := systemMetrics["memory_usage_percent"]; memMetric != nil {
		summary.MemoryUsage = memMetric.Value().(float64)
	}

	// 告警指标
	if rg.alertSystem != nil {
		activeAlerts := rg.alertSystem.GetActiveAlerts()
		summary.ActiveAlerts = int64(len(activeAlerts))
	}

	return summary
}

// generateDetails 生成详情
func (rg *ReportGenerator) generateDetails(period ReportPeriod) ReportDetails {
	details := ReportDetails{}

	// HTTP详情
	if rg.httpMonitor != nil {
		details.HTTPMetrics = rg.generateHTTPDetails()
	}

	// 数据库详情
	if rg.dbMonitor != nil {
		details.DatabaseMetrics = rg.generateDatabaseDetails()
	}

	// 缓存详情
	if rg.cacheMonitor != nil {
		details.CacheMetrics = rg.generateCacheDetails()
	}

	// 系统详情
	details.SystemMetrics = rg.generateSystemDetails()

	// 告警详情
	if rg.alertSystem != nil {
		details.AlertMetrics = rg.generateAlertDetails()
	}

	return details
}

// generateHTTPDetails 生成HTTP详情
func (rg *ReportGenerator) generateHTTPDetails() HTTPReportDetails {
	details := HTTPReportDetails{
		RequestDistribution:      make(map[string]int64),
		ResponseTimeDistribution: make(map[string]int64),
		StatusCodeDistribution:   make(map[int]int64),
	}

	// 获取请求历史
	collector := rg.httpMonitor.GetMetrics()
	if collector != nil {
		// 这里需要根据实际的HTTP监控器实现来获取详细信息
		// 由于当前实现限制，这里提供基础结构
	}

	return details
}

// generateDatabaseDetails 生成数据库详情
func (rg *ReportGenerator) generateDatabaseDetails() DatabaseReportDetails {
	details := DatabaseReportDetails{
		QueryTypeDistribution: make(map[string]int64),
	}

	if rg.dbMonitor != nil {
		// 获取查询类型分布
		distribution := rg.dbMonitor.GetQueryTypeDistribution()
		for k, v := range distribution {
			details.QueryTypeDistribution[k] = int64(v)
		}

		// 获取慢查询
		details.SlowQueries = rg.dbMonitor.GetSlowQueries()

		// 获取错误查询
		details.ErrorQueries = rg.dbMonitor.GetErrorQueries()

		// 获取平均查询时间
		details.AverageQueryTime = rg.dbMonitor.GetAverageQueryTime()

		// 连接池统计
		metrics := rg.dbMonitor.GetMetrics()
		if metrics != nil {
			details.ConnectionPoolStats = ConnectionPoolStats{
				ActiveConnections: int64(metrics.activeConnections.Value().(float64)),
				IdleConnections:   int64(metrics.idleConnections.Value().(float64)),
				TotalConnections:  int64(metrics.totalConnections.Value().(float64)),
			}
		}
	}

	return details
}

// generateCacheDetails 生成缓存详情
func (rg *ReportGenerator) generateCacheDetails() CacheReportDetails {
	details := CacheReportDetails{
		OperationDistribution: make(map[string]int64),
	}

	if rg.cacheMonitor != nil {
		// 获取命中率
		details.HitRate = rg.cacheMonitor.GetHitRate()

		// 获取操作分布
		distribution := rg.cacheMonitor.GetOperationDistribution()
		for k, v := range distribution {
			details.OperationDistribution[k] = int64(v)
		}

		// 获取慢操作
		details.SlowOperations = rg.cacheMonitor.GetSlowOperations(100 * time.Microsecond)

		// 获取错误操作
		details.ErrorOperations = rg.cacheMonitor.GetErrorOperations()

		// 存储统计
		metrics := rg.cacheMonitor.GetMetrics()
		if metrics != nil {
			details.StorageStats = StorageStats{
				ItemCount:   int64(metrics.itemCount.Value().(float64)),
				MemoryUsage: int64(metrics.memoryUsage.Value().(float64)),
			}
		}
	}

	return details
}

// generateSystemDetails 生成系统详情
func (rg *ReportGenerator) generateSystemDetails() SystemReportDetails {
	details := SystemReportDetails{}

	systemMetrics := rg.monitor.GetAllMetrics()

	// CPU统计
	if cpuMetric := systemMetrics["cpu_usage"]; cpuMetric != nil {
		details.CPUStats.Usage = cpuMetric.Value().(float64)
	}
	if cpuCoresMetric := systemMetrics["cpu_cores"]; cpuCoresMetric != nil {
		details.CPUStats.Cores = int(cpuCoresMetric.Value().(float64))
	}

	// 内存统计
	if memTotalMetric := systemMetrics["memory_total"]; memTotalMetric != nil {
		details.MemoryStats.Total = uint64(memTotalMetric.Value().(float64))
	}
	if memUsedMetric := systemMetrics["memory_used"]; memUsedMetric != nil {
		details.MemoryStats.Used = uint64(memUsedMetric.Value().(float64))
	}
	if memUsageMetric := systemMetrics["memory_usage_percent"]; memUsageMetric != nil {
		details.MemoryStats.Usage = memUsageMetric.Value().(float64)
	}

	// Go运行时统计
	if heapAllocMetric := systemMetrics["go_heap_alloc"]; heapAllocMetric != nil {
		details.MemoryStats.HeapAlloc = uint64(heapAllocMetric.Value().(float64))
	}
	if heapSysMetric := systemMetrics["go_heap_sys"]; heapSysMetric != nil {
		details.MemoryStats.HeapSys = uint64(heapSysMetric.Value().(float64))
	}
	if goroutinesMetric := systemMetrics["go_goroutines"]; goroutinesMetric != nil {
		details.GoroutineStats.Count = int(goroutinesMetric.Value().(float64))
	}

	return details
}

// generateAlertDetails 生成告警详情
func (rg *ReportGenerator) generateAlertDetails() AlertReportDetails {
	details := AlertReportDetails{
		AlertByLevel: make(map[string]int64),
		AlertByRule:  make(map[string]int64),
	}

	if rg.alertSystem != nil {
		alerts := rg.alertSystem.GetAlerts()
		details.TotalAlerts = int64(len(alerts))

		activeAlerts := rg.alertSystem.GetActiveAlerts()
		details.ActiveAlerts = int64(len(activeAlerts))
		details.ResolvedAlerts = details.TotalAlerts - details.ActiveAlerts

		// 按级别统计
		for _, alert := range alerts {
			details.AlertByLevel[string(alert.Level)]++
			details.AlertByRule[alert.RuleName]++
		}
	}

	return details
}

// generateTrends 生成趋势分析
func (rg *ReportGenerator) generateTrends(period ReportPeriod) map[string]interface{} {
	trends := make(map[string]interface{})

	// 这里可以实现趋势分析逻辑
	// 例如：响应时间趋势、错误率趋势、CPU使用率趋势等

	return trends
}

// generateComparison 生成对比分析
func (rg *ReportGenerator) generateComparison(period ReportPeriod) map[string]interface{} {
	comparison := make(map[string]interface{})

	// 这里可以实现对比分析逻辑
	// 例如：与历史同期对比、与基准值对比等

	return comparison
}

// generateRecommendations 生成优化建议
func (rg *ReportGenerator) generateRecommendations(summary ReportSummary) []Recommendation {
	var recommendations []Recommendation

	// 响应时间建议
	if summary.AverageResponseTime > 500*time.Millisecond {
		recommendations = append(recommendations, Recommendation{
			Type:        "performance",
			Title:       "优化响应时间",
			Description: "平均响应时间超过500ms，建议优化数据库查询和缓存策略",
			Priority:    "high",
			Impact:      30.0,
			Effort:      "medium",
			Actions:     []string{"优化慢查询", "增加缓存", "优化代码逻辑"},
		})
	}

	// 错误率建议
	if summary.ErrorRate > 5.0 {
		recommendations = append(recommendations, Recommendation{
			Type:        "reliability",
			Title:       "降低错误率",
			Description: "错误率超过5%，建议检查异常处理和系统稳定性",
			Priority:    "critical",
			Impact:      50.0,
			Effort:      "high",
			Actions:     []string{"检查错误日志", "改进异常处理", "增加监控告警"},
		})
	}

	// 缓存命中率建议
	if summary.CacheHitRate < 80.0 {
		recommendations = append(recommendations, Recommendation{
			Type:        "cache",
			Title:       "提高缓存命中率",
			Description: "缓存命中率低于80%，建议优化缓存策略",
			Priority:    "medium",
			Impact:      20.0,
			Effort:      "low",
			Actions:     []string{"增加缓存预热", "优化缓存键设计", "调整缓存过期时间"},
		})
	}

	// CPU使用率建议
	if summary.CPUUsage > 80.0 {
		recommendations = append(recommendations, Recommendation{
			Type:        "resource",
			Title:       "优化CPU使用率",
			Description: "CPU使用率超过80%，建议优化计算密集型操作",
			Priority:    "high",
			Impact:      25.0,
			Effort:      "medium",
			Actions:     []string{"优化算法", "增加并发处理", "考虑水平扩展"},
		})
	}

	// 内存使用率建议
	if summary.MemoryUsage > 85.0 {
		recommendations = append(recommendations, Recommendation{
			Type:        "resource",
			Title:       "优化内存使用",
			Description: "内存使用率超过85%，建议优化内存分配和垃圾回收",
			Priority:    "high",
			Impact:      20.0,
			Effort:      "medium",
			Actions:     []string{"优化内存分配", "调整GC参数", "检查内存泄漏"},
		})
	}

	return recommendations
}

// calculateErrorRate 计算错误率
func (rg *ReportGenerator) calculateErrorRate(metrics *HTTPMetrics) float64 {
	errors := metrics.errorCounter.Value().(int64)
	total := metrics.requestCounter.Value().(int64)

	if total == 0 {
		return 0.0
	}

	return float64(errors) / float64(total) * 100.0
}

// calculateAverageResponseTime 计算平均响应时间
func (rg *ReportGenerator) calculateAverageResponseTime(metrics *HTTPMetrics) time.Duration {
	// 这里需要根据实际的直方图数据计算平均值
	// 简化实现，返回0
	return 0
}

// ExportReport 导出报告
func (rg *ReportGenerator) ExportReport(report *PerformanceReport, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(report, "", "  ")
	case "text":
		return rg.exportAsText(report)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// exportAsText 导出为文本格式
func (rg *ReportGenerator) exportAsText(report *PerformanceReport) ([]byte, error) {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("性能监控报告\n"))
	builder.WriteString(fmt.Sprintf("================\n\n"))
	builder.WriteString(fmt.Sprintf("报告ID: %s\n", report.ID))
	builder.WriteString(fmt.Sprintf("报告类型: %s\n", report.Type))
	builder.WriteString(fmt.Sprintf("生成时间: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05")))
	builder.WriteString(fmt.Sprintf("报告周期: %s - %s\n\n",
		report.Period.Start.Format("2006-01-02 15:04:05"),
		report.Period.End.Format("2006-01-02 15:04:05")))

	// 摘要
	builder.WriteString("性能摘要\n")
	builder.WriteString("--------\n")
	builder.WriteString(fmt.Sprintf("总请求数: %d\n", report.Summary.TotalRequests))
	builder.WriteString(fmt.Sprintf("平均响应时间: %v\n", report.Summary.AverageResponseTime))
	builder.WriteString(fmt.Sprintf("错误率: %.2f%%\n", report.Summary.ErrorRate))
	builder.WriteString(fmt.Sprintf("吞吐量: %.2f req/s\n", report.Summary.Throughput))
	builder.WriteString(fmt.Sprintf("CPU使用率: %.2f%%\n", report.Summary.CPUUsage))
	builder.WriteString(fmt.Sprintf("内存使用率: %.2f%%\n", report.Summary.MemoryUsage))
	builder.WriteString(fmt.Sprintf("缓存命中率: %.2f%%\n", report.Summary.CacheHitRate))
	builder.WriteString(fmt.Sprintf("活跃告警: %d\n\n", report.Summary.ActiveAlerts))

	// 建议
	if len(report.Recommendations) > 0 {
		builder.WriteString("优化建议\n")
		builder.WriteString("--------\n")
		for i, rec := range report.Recommendations {
			builder.WriteString(fmt.Sprintf("%d. %s (%s优先级)\n", i+1, rec.Title, rec.Priority))
			builder.WriteString(fmt.Sprintf("   描述: %s\n", rec.Description))
			builder.WriteString(fmt.Sprintf("   预期改进: %.1f%%\n", rec.Impact))
			builder.WriteString(fmt.Sprintf("   实施难度: %s\n", rec.Effort))
			builder.WriteString("\n")
		}
	}

	return []byte(builder.String()), nil
}
