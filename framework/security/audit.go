package security

import (
	"fmt"
	"net/http"
	"time"
)

// AuditLevel 审计级别
type AuditLevel string

const (
	AuditLevelInfo     AuditLevel = "info"
	AuditLevelWarning  AuditLevel = "warning"
	AuditLevelError    AuditLevel = "error"
	AuditLevelCritical AuditLevel = "critical"
)

// AuditEvent 审计事件
type AuditEvent struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Level     AuditLevel             `json:"level"`
	Category  string                 `json:"category"`
	Action    string                 `json:"action"`
	UserID    string                 `json:"user_id,omitempty"`
	IPAddress string                 `json:"ip_address,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Resource  string                 `json:"resource,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RiskScore int                    `json:"risk_score,omitempty"`
	Blocked   bool                   `json:"blocked,omitempty"`
	Reason    string                 `json:"reason,omitempty"`
}

// AuditLogger 审计日志记录器接口
type AuditLogger interface {
	// Log 记录审计事件
	Log(event *AuditEvent) error
	// GetEvents 获取审计事件
	GetEvents(filter *AuditFilter) ([]*AuditEvent, error)
	// GetEventByID 根据ID获取事件
	GetEventByID(id string) (*AuditEvent, error)
	// Clear 清空审计日志
	Clear() error
}

// AuditFilter 审计过滤器
type AuditFilter struct {
	Level     AuditLevel `json:"level,omitempty"`
	Category  string     `json:"category,omitempty"`
	UserID    string     `json:"user_id,omitempty"`
	IPAddress string     `json:"ip_address,omitempty"`
	StartTime time.Time  `json:"start_time,omitempty"`
	EndTime   time.Time  `json:"end_time,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Offset    int        `json:"offset,omitempty"`
}

// SecurityAuditor 安全审计器
type SecurityAuditor struct {
	logger AuditLogger
	config *AuditorConfig
}

// AuditorConfig 审计器配置
type AuditorConfig struct {
	Enabled       bool          `json:"enabled"`
	LogLevel      AuditLevel    `json:"log_level"`
	MaxEvents     int           `json:"max_events"`
	Retention     time.Duration `json:"retention"`
	RiskThreshold int           `json:"risk_threshold"`
}

// NewSecurityAuditor 创建安全审计器
func NewSecurityAuditor(logger AuditLogger, config *AuditorConfig) *SecurityAuditor {
	if config == nil {
		config = &AuditorConfig{
			Enabled:       true,
			LogLevel:      AuditLevelInfo,
			MaxEvents:     10000,
			Retention:     30 * 24 * time.Hour, // 30天
			RiskThreshold: 80,
		}
	}

	return &SecurityAuditor{
		logger: logger,
		config: config,
	}
}

// LogSecurityEvent 记录安全事件
func (sa *SecurityAuditor) LogSecurityEvent(level AuditLevel, category, action string, r *http.Request, details map[string]interface{}) error {
	if !sa.config.Enabled {
		return nil
	}

	event := &AuditEvent{
		ID:        generateEventID(),
		Timestamp: time.Now(),
		Level:     level,
		Category:  category,
		Action:    action,
		Details:   details,
	}

	// 从请求中提取信息
	if r != nil {
		event.IPAddress = getClientIP(r)
		event.UserAgent = r.UserAgent()
		event.Resource = r.URL.Path
	}

	// 计算风险分数
	event.RiskScore = sa.calculateRiskScore(event)

	// 检查是否需要阻止
	if event.RiskScore >= sa.config.RiskThreshold {
		event.Blocked = true
		event.Reason = "Risk score exceeded threshold"
	}

	return sa.logger.Log(event)
}

// LogLoginAttempt 记录登录尝试
func (sa *SecurityAuditor) LogLoginAttempt(userID, ipAddress, userAgent string, success bool, details map[string]interface{}) error {
	level := AuditLevelInfo
	if !success {
		level = AuditLevelWarning
	}

	event := &AuditEvent{
		ID:        generateEventID(),
		Timestamp: time.Now(),
		Level:     level,
		Category:  "authentication",
		Action:    "login_attempt",
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
	}

	event.RiskScore = sa.calculateRiskScore(event)
	if event.RiskScore >= sa.config.RiskThreshold {
		event.Blocked = true
		event.Reason = "Suspicious login attempt"
	}

	return sa.logger.Log(event)
}

// LogAccessAttempt 记录访问尝试
func (sa *SecurityAuditor) LogAccessAttempt(userID, resource, ipAddress string, success bool, details map[string]interface{}) error {
	level := AuditLevelInfo
	if !success {
		level = AuditLevelWarning
	}

	event := &AuditEvent{
		ID:        generateEventID(),
		Timestamp: time.Now(),
		Level:     level,
		Category:  "authorization",
		Action:    "access_attempt",
		UserID:    userID,
		IPAddress: ipAddress,
		Resource:  resource,
		Details:   details,
	}

	event.RiskScore = sa.calculateRiskScore(event)
	if event.RiskScore >= sa.config.RiskThreshold {
		event.Blocked = true
		event.Reason = "Unauthorized access attempt"
	}

	return sa.logger.Log(event)
}

// LogSecurityViolation 记录安全违规
func (sa *SecurityAuditor) LogSecurityViolation(violationType, ipAddress, userAgent string, details map[string]interface{}) error {
	event := &AuditEvent{
		ID:        generateEventID(),
		Timestamp: time.Now(),
		Level:     AuditLevelError,
		Category:  "security_violation",
		Action:    violationType,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		Blocked:   true,
		Reason:    "Security violation detected",
	}

	event.RiskScore = sa.calculateRiskScore(event)
	return sa.logger.Log(event)
}

// GetSecurityReport 获取安全报告
func (sa *SecurityAuditor) GetSecurityReport(startTime, endTime time.Time) (*SecurityReport, error) {
	filter := &AuditFilter{
		StartTime: startTime,
		EndTime:   endTime,
	}

	events, err := sa.logger.GetEvents(filter)
	if err != nil {
		return nil, err
	}

	report := &SecurityReport{
		Period:      fmt.Sprintf("%s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")),
		TotalEvents: len(events),
		GeneratedAt: time.Now(),
	}

	// 统计各级别事件
	levelStats := make(map[AuditLevel]int)
	categoryStats := make(map[string]int)
	actionStats := make(map[string]int)
	blockedCount := 0
	highRiskCount := 0

	for _, event := range events {
		levelStats[event.Level]++
		categoryStats[event.Category]++
		actionStats[event.Action]++

		if event.Blocked {
			blockedCount++
		}

		if event.RiskScore >= sa.config.RiskThreshold {
			highRiskCount++
		}
	}

	report.LevelStats = levelStats
	report.CategoryStats = categoryStats
	report.ActionStats = actionStats
	report.BlockedEvents = blockedCount
	report.HighRiskEvents = highRiskCount

	return report, nil
}

// calculateRiskScore 计算风险分数
func (sa *SecurityAuditor) calculateRiskScore(event *AuditEvent) int {
	score := 0

	// 基于级别
	switch event.Level {
	case AuditLevelInfo:
		score += 10
	case AuditLevelWarning:
		score += 30
	case AuditLevelError:
		score += 50
	case AuditLevelCritical:
		score += 80
	}

	// 基于类别
	switch event.Category {
	case "authentication":
		score += 20
	case "authorization":
		score += 25
	case "security_violation":
		score += 60
	}

	// 基于动作
	switch event.Action {
	case "login_attempt":
		score += 15
	case "access_attempt":
		score += 20
	case "sql_injection":
		score += 70
	case "xss_attempt":
		score += 65
	case "csrf_violation":
		score += 55
	}

	// 基于详情
	if details, ok := event.Details["failed_attempts"]; ok {
		if attempts, ok := details.(int); ok && attempts > 5 {
			score += attempts * 5
		}
	}

	return score
}

// SecurityReport 安全报告
type SecurityReport struct {
	Period         string             `json:"period"`
	TotalEvents    int                `json:"total_events"`
	LevelStats     map[AuditLevel]int `json:"level_stats"`
	CategoryStats  map[string]int     `json:"category_stats"`
	ActionStats    map[string]int     `json:"action_stats"`
	BlockedEvents  int                `json:"blocked_events"`
	HighRiskEvents int                `json:"high_risk_events"`
	GeneratedAt    time.Time          `json:"generated_at"`
}

// MemoryAuditLogger 内存审计日志记录器
type MemoryAuditLogger struct {
	events    []*AuditEvent
	maxEvents int
}

// NewMemoryAuditLogger 创建内存审计日志记录器
func NewMemoryAuditLogger(maxEvents int) *MemoryAuditLogger {
	if maxEvents <= 0 {
		maxEvents = 10000
	}

	return &MemoryAuditLogger{
		events:    make([]*AuditEvent, 0),
		maxEvents: maxEvents,
	}
}

// Log 记录审计事件
func (mal *MemoryAuditLogger) Log(event *AuditEvent) error {
	mal.events = append(mal.events, event)

	// 如果超过最大事件数，移除最旧的事件
	if len(mal.events) > mal.maxEvents {
		mal.events = mal.events[1:]
	}

	return nil
}

// GetEvents 获取审计事件
func (mal *MemoryAuditLogger) GetEvents(filter *AuditFilter) ([]*AuditEvent, error) {
	var filteredEvents []*AuditEvent

	for _, event := range mal.events {
		if mal.matchesFilter(event, filter) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	// 应用分页
	if filter != nil {
		if filter.Offset > 0 && filter.Offset < len(filteredEvents) {
			filteredEvents = filteredEvents[filter.Offset:]
		}

		if filter.Limit > 0 && filter.Limit < len(filteredEvents) {
			filteredEvents = filteredEvents[:filter.Limit]
		}
	}

	return filteredEvents, nil
}

// GetEventByID 根据ID获取事件
func (mal *MemoryAuditLogger) GetEventByID(id string) (*AuditEvent, error) {
	for _, event := range mal.events {
		if event.ID == id {
			return event, nil
		}
	}
	return nil, fmt.Errorf("event not found: %s", id)
}

// Clear 清空审计日志
func (mal *MemoryAuditLogger) Clear() error {
	mal.events = make([]*AuditEvent, 0)
	return nil
}

// matchesFilter 检查事件是否匹配过滤器
func (mal *MemoryAuditLogger) matchesFilter(event *AuditEvent, filter *AuditFilter) bool {
	if filter == nil {
		return true
	}

	if filter.Level != "" && event.Level != filter.Level {
		return false
	}

	if filter.Category != "" && event.Category != filter.Category {
		return false
	}

	if filter.UserID != "" && event.UserID != filter.UserID {
		return false
	}

	if filter.IPAddress != "" && event.IPAddress != filter.IPAddress {
		return false
	}

	if !filter.StartTime.IsZero() && event.Timestamp.Before(filter.StartTime) {
		return false
	}

	if !filter.EndTime.IsZero() && event.Timestamp.After(filter.EndTime) {
		return false
	}

	return true
}

// 辅助函数

// generateEventID 生成事件ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

// getClientIP 获取客户端IP
func getClientIP(r *http.Request) string {
	// 检查代理头
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Client-IP"); ip != "" {
		return ip
	}

	return r.RemoteAddr
}

// JSONAuditLogger JSON文件审计日志记录器
type JSONAuditLogger struct {
	filePath string
}

// NewJSONAuditLogger 创建JSON文件审计日志记录器
func NewJSONAuditLogger(filePath string) *JSONAuditLogger {
	return &JSONAuditLogger{
		filePath: filePath,
	}
}

// Log 记录审计事件
func (jal *JSONAuditLogger) Log(event *AuditEvent) error {
	// 这里应该实现将事件写入JSON文件的逻辑
	// 为了简化，暂时只返回nil
	return nil
}

// GetEvents 获取审计事件
func (jal *JSONAuditLogger) GetEvents(filter *AuditFilter) ([]*AuditEvent, error) {
	// 这里应该实现从JSON文件读取事件的逻辑
	// 为了简化，暂时返回空切片
	return []*AuditEvent{}, nil
}

// GetEventByID 根据ID获取事件
func (jal *JSONAuditLogger) GetEventByID(id string) (*AuditEvent, error) {
	// 这里应该实现从JSON文件查找特定事件的逻辑
	return nil, fmt.Errorf("event not found: %s", id)
}

// Clear 清空审计日志
func (jal *JSONAuditLogger) Clear() error {
	// 这里应该实现清空JSON文件的逻辑
	return nil
}
