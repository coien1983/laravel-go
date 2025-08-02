package performance

import (
	"strings"
	"sync"
	"time"
)

// DatabaseMetrics 数据库指标
type DatabaseMetrics struct {
	// 查询计数器
	queryCounter     *Counter
	slowQueryCounter *Counter
	errorCounter     *Counter

	// 查询时间直方图
	queryTimeHistogram *Histogram

	// 连接池指标
	activeConnections *Gauge
	idleConnections   *Gauge
	totalConnections  *Gauge

	// 查询类型分布
	selectCounter *Counter
	insertCounter *Counter
	updateCounter *Counter
	deleteCounter *Counter

	// 事务指标
	transactionCounter       *Counter
	transactionTimeHistogram *Histogram
}

// NewDatabaseMetrics 创建数据库指标
func NewDatabaseMetrics(monitor Monitor) *DatabaseMetrics {
	// 创建查询时间直方图，单位为毫秒
	queryTimeBuckets := []float64{1, 5, 10, 50, 100, 200, 500, 1000, 2000, 5000}
	queryTimeHistogram := NewHistogram("database_query_time", queryTimeBuckets, map[string]string{"unit": "milliseconds"})
	monitor.RegisterMetric(queryTimeHistogram)

	// 创建事务时间直方图
	transactionTimeBuckets := []float64{10, 50, 100, 200, 500, 1000, 2000, 5000, 10000}
	transactionTimeHistogram := NewHistogram("database_transaction_time", transactionTimeBuckets, map[string]string{"unit": "milliseconds"})
	monitor.RegisterMetric(transactionTimeHistogram)

	// 创建计数器
	queryCounter := NewCounter("database_queries_total", map[string]string{"type": "total"})
	monitor.RegisterMetric(queryCounter)

	slowQueryCounter := NewCounter("database_slow_queries_total", map[string]string{"type": "slow"})
	monitor.RegisterMetric(slowQueryCounter)

	errorCounter := NewCounter("database_errors_total", map[string]string{"type": "error"})
	monitor.RegisterMetric(errorCounter)

	// 查询类型计数器
	selectCounter := NewCounter("database_queries_select", map[string]string{"type": "select"})
	monitor.RegisterMetric(selectCounter)

	insertCounter := NewCounter("database_queries_insert", map[string]string{"type": "insert"})
	monitor.RegisterMetric(insertCounter)

	updateCounter := NewCounter("database_queries_update", map[string]string{"type": "update"})
	monitor.RegisterMetric(updateCounter)

	deleteCounter := NewCounter("database_queries_delete", map[string]string{"type": "delete"})
	monitor.RegisterMetric(deleteCounter)

	// 事务计数器
	transactionCounter := NewCounter("database_transactions_total", map[string]string{"type": "total"})
	monitor.RegisterMetric(transactionCounter)

	// 连接池指标
	activeConnections := NewGauge("database_active_connections", map[string]string{"type": "active"})
	monitor.RegisterMetric(activeConnections)

	idleConnections := NewGauge("database_idle_connections", map[string]string{"type": "idle"})
	monitor.RegisterMetric(idleConnections)

	totalConnections := NewGauge("database_total_connections", map[string]string{"type": "total"})
	monitor.RegisterMetric(totalConnections)

	return &DatabaseMetrics{
		queryCounter:             queryCounter,
		slowQueryCounter:         slowQueryCounter,
		errorCounter:             errorCounter,
		queryTimeHistogram:       queryTimeHistogram,
		activeConnections:        activeConnections,
		idleConnections:          idleConnections,
		totalConnections:         totalConnections,
		selectCounter:            selectCounter,
		insertCounter:            insertCounter,
		updateCounter:            updateCounter,
		deleteCounter:            deleteCounter,
		transactionCounter:       transactionCounter,
		transactionTimeHistogram: transactionTimeHistogram,
	}
}

// DatabaseMonitor 数据库监控器
type DatabaseMonitor struct {
	metrics            *DatabaseMetrics
	slowQueryThreshold time.Duration
	mu                 sync.RWMutex
	queryHistory       []QueryRecord
	maxHistorySize     int
}

// QueryRecord 查询记录
type QueryRecord struct {
	SQL       string        `json:"sql"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Success   bool          `json:"success"`
	Error     string        `json:"error,omitempty"`
	Type      string        `json:"type"`
}

// NewDatabaseMonitor 创建数据库监控器
func NewDatabaseMonitor(monitor Monitor, slowQueryThreshold time.Duration) *DatabaseMonitor {
	return &DatabaseMonitor{
		metrics:            NewDatabaseMetrics(monitor),
		slowQueryThreshold: slowQueryThreshold,
		queryHistory:       make([]QueryRecord, 0),
		maxHistorySize:     1000,
	}
}

// RecordQuery 记录查询
func (dm *DatabaseMonitor) RecordQuery(sql string, duration time.Duration, success bool, err error) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// 解析查询类型
	queryType := dm.parseQueryType(sql)

	// 增加查询计数器
	dm.metrics.queryCounter.Increment(1)

	// 记录查询时间
	dm.metrics.queryTimeHistogram.Observe(float64(duration.Milliseconds()))

	// 记录查询类型
	switch queryType {
	case "SELECT":
		dm.metrics.selectCounter.Increment(1)
	case "INSERT":
		dm.metrics.insertCounter.Increment(1)
	case "UPDATE":
		dm.metrics.updateCounter.Increment(1)
	case "DELETE":
		dm.metrics.deleteCounter.Increment(1)
	}

	// 检查慢查询
	if duration > dm.slowQueryThreshold {
		dm.metrics.slowQueryCounter.Increment(1)
	}

	// 记录错误
	if !success {
		dm.metrics.errorCounter.Increment(1)
	}

	// 添加到历史记录
	record := QueryRecord{
		SQL:       sql,
		Duration:  duration,
		Timestamp: time.Now(),
		Success:   success,
		Type:      queryType,
	}
	if err != nil {
		record.Error = err.Error()
	}

	dm.addToHistory(record)
}

// RecordTransaction 记录事务
func (dm *DatabaseMonitor) RecordTransaction(duration time.Duration, success bool) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// 增加事务计数器
	dm.metrics.transactionCounter.Increment(1)

	// 记录事务时间
	dm.metrics.transactionTimeHistogram.Observe(float64(duration.Milliseconds()))

	// 记录错误
	if !success {
		dm.metrics.errorCounter.Increment(1)
	}
}

// UpdateConnectionPool 更新连接池状态
func (dm *DatabaseMonitor) UpdateConnectionPool(active, idle, total int) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.metrics.activeConnections.Set(float64(active))
	dm.metrics.idleConnections.Set(float64(idle))
	dm.metrics.totalConnections.Set(float64(total))
}

// GetMetrics 获取指标
func (dm *DatabaseMonitor) GetMetrics() *DatabaseMetrics {
	return dm.metrics
}

// GetQueryHistory 获取查询历史
func (dm *DatabaseMonitor) GetQueryHistory() []QueryRecord {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	result := make([]QueryRecord, len(dm.queryHistory))
	copy(result, dm.queryHistory)
	return result
}

// GetSlowQueries 获取慢查询
func (dm *DatabaseMonitor) GetSlowQueries() []QueryRecord {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	var slowQueries []QueryRecord
	for _, record := range dm.queryHistory {
		if record.Duration > dm.slowQueryThreshold {
			slowQueries = append(slowQueries, record)
		}
	}
	return slowQueries
}

// GetErrorQueries 获取错误查询
func (dm *DatabaseMonitor) GetErrorQueries() []QueryRecord {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	var errorQueries []QueryRecord
	for _, record := range dm.queryHistory {
		if !record.Success {
			errorQueries = append(errorQueries, record)
		}
	}
	return errorQueries
}

// GetAverageQueryTime 获取平均查询时间
func (dm *DatabaseMonitor) GetAverageQueryTime() time.Duration {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if len(dm.queryHistory) == 0 {
		return 0
	}

	var total time.Duration
	for _, record := range dm.queryHistory {
		total += record.Duration
	}

	return total / time.Duration(len(dm.queryHistory))
}

// GetQueryTypeDistribution 获取查询类型分布
func (dm *DatabaseMonitor) GetQueryTypeDistribution() map[string]int {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	distribution := make(map[string]int)
	for _, record := range dm.queryHistory {
		distribution[record.Type]++
	}
	return distribution
}

// ClearHistory 清空历史记录
func (dm *DatabaseMonitor) ClearHistory() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.queryHistory = make([]QueryRecord, 0)
}

// parseQueryType 解析查询类型
func (dm *DatabaseMonitor) parseQueryType(sql string) string {
	sql = strings.TrimSpace(strings.ToUpper(sql))

	if strings.HasPrefix(sql, "SELECT") {
		return "SELECT"
	} else if strings.HasPrefix(sql, "INSERT") {
		return "INSERT"
	} else if strings.HasPrefix(sql, "UPDATE") {
		return "UPDATE"
	} else if strings.HasPrefix(sql, "DELETE") {
		return "DELETE"
	} else if strings.HasPrefix(sql, "CREATE") {
		return "CREATE"
	} else if strings.HasPrefix(sql, "ALTER") {
		return "ALTER"
	} else if strings.HasPrefix(sql, "DROP") {
		return "DROP"
	} else if strings.HasPrefix(sql, "BEGIN") || strings.HasPrefix(sql, "COMMIT") || strings.HasPrefix(sql, "ROLLBACK") {
		return "TRANSACTION"
	}

	return "OTHER"
}

// addToHistory 添加到历史记录
func (dm *DatabaseMonitor) addToHistory(record QueryRecord) {
	dm.queryHistory = append(dm.queryHistory, record)

	// 限制历史记录大小
	if len(dm.queryHistory) > dm.maxHistorySize {
		dm.queryHistory = dm.queryHistory[1:]
	}
}

// DatabaseMonitorMiddleware 数据库监控中间件
type DatabaseMonitorMiddleware struct {
	monitor *DatabaseMonitor
}

// NewDatabaseMonitorMiddleware 创建数据库监控中间件
func NewDatabaseMonitorMiddleware(monitor *DatabaseMonitor) *DatabaseMonitorMiddleware {
	return &DatabaseMonitorMiddleware{
		monitor: monitor,
	}
}

// WrapQuery 包装查询函数
func (dm *DatabaseMonitorMiddleware) WrapQuery(queryFunc func() error) error {
	start := time.Now()

	err := queryFunc()

	duration := time.Since(start)
	dm.monitor.RecordQuery("", duration, err == nil, err)

	return err
}

// WrapTransaction 包装事务函数
func (dm *DatabaseMonitorMiddleware) WrapTransaction(txFunc func() error) error {
	start := time.Now()

	err := txFunc()

	duration := time.Since(start)
	dm.monitor.RecordTransaction(duration, err == nil)

	return err
}
