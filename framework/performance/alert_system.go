package performance

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// AlertLevel 告警级别
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "info"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelError    AlertLevel = "error"
	AlertLevelCritical AlertLevel = "critical"
)

// AlertRule 告警规则
type AlertRule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	MetricName  string            `json:"metric_name"`
	Condition   string            `json:"condition"` // >, <, >=, <=, ==, !=
	Threshold   float64           `json:"threshold"`
	Duration    time.Duration     `json:"duration"` // 持续时间
	Level       AlertLevel        `json:"level"`
	Enabled     bool              `json:"enabled"`
	Labels      map[string]string `json:"labels"`
	Actions     []string          `json:"actions"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// AlertAction 告警动作
type AlertAction interface {
	Execute(alert *Alert) error
	GetType() string
	GetDescription() string
}

// Alert 告警
type Alert struct {
	ID         string            `json:"id"`
	RuleID     string            `json:"rule_id"`
	RuleName   string            `json:"rule_name"`
	Level      AlertLevel        `json:"level"`
	Message    string            `json:"message"`
	MetricName string            `json:"metric_name"`
	Value      interface{}       `json:"value"`
	Threshold  float64           `json:"threshold"`
	Labels     map[string]string `json:"labels"`
	Timestamp  time.Time         `json:"timestamp"`
	Resolved   bool              `json:"resolved"`
	ResolvedAt *time.Time        `json:"resolved_at,omitempty"`
}

// AlertSystem 告警系统
type AlertSystem struct {
	monitor       Monitor
	rules         map[string]*AlertRule
	alerts        map[string]*Alert
	mu            sync.RWMutex
	running       bool
	ctx           context.Context
	cancel        context.CancelFunc
	checkInterval time.Duration
	actions       map[string]AlertAction
}

// NewAlertSystem 创建告警系统
func NewAlertSystem(monitor Monitor) *AlertSystem {
	as := &AlertSystem{
		monitor:       monitor,
		rules:         make(map[string]*AlertRule),
		alerts:        make(map[string]*Alert),
		checkInterval: 30 * time.Second,
		actions:       make(map[string]AlertAction),
	}

	// 注册默认动作
	as.RegisterAction(NewLogAlertAction())
	as.RegisterAction(NewEmailAlertAction())
	as.RegisterAction(NewWebhookAlertAction())

	return as
}

// RegisterAction 注册告警动作
func (as *AlertSystem) RegisterAction(action AlertAction) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.actions[action.GetType()] = action
}

// AddRule 添加告警规则
func (as *AlertSystem) AddRule(rule *AlertRule) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if rule.ID == "" {
		return fmt.Errorf("rule ID cannot be empty")
	}

	if _, exists := as.rules[rule.ID]; exists {
		return fmt.Errorf("rule with ID %s already exists", rule.ID)
	}

	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()
	as.rules[rule.ID] = rule

	return nil
}

// UpdateRule 更新告警规则
func (as *AlertSystem) UpdateRule(rule *AlertRule) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if _, exists := as.rules[rule.ID]; !exists {
		return fmt.Errorf("rule with ID %s does not exist", rule.ID)
	}

	rule.UpdatedAt = time.Now()
	as.rules[rule.ID] = rule

	return nil
}

// RemoveRule 移除告警规则
func (as *AlertSystem) RemoveRule(ruleID string) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if _, exists := as.rules[ruleID]; !exists {
		return fmt.Errorf("rule with ID %s does not exist", ruleID)
	}

	delete(as.rules, ruleID)

	return nil
}

// GetRule 获取告警规则
func (as *AlertSystem) GetRule(ruleID string) (*AlertRule, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	rule, exists := as.rules[ruleID]
	if !exists {
		return nil, fmt.Errorf("rule with ID %s does not exist", ruleID)
	}

	return rule, nil
}

// GetRules 获取所有告警规则
func (as *AlertSystem) GetRules() []*AlertRule {
	as.mu.RLock()
	defer as.mu.RUnlock()

	rules := make([]*AlertRule, 0, len(as.rules))
	for _, rule := range as.rules {
		rules = append(rules, rule)
	}

	return rules
}

// GetAlerts 获取所有告警
func (as *AlertSystem) GetAlerts() []*Alert {
	as.mu.RLock()
	defer as.mu.RUnlock()

	alerts := make([]*Alert, 0, len(as.alerts))
	for _, alert := range as.alerts {
		alerts = append(alerts, alert)
	}

	return alerts
}

// GetActiveAlerts 获取活跃告警
func (as *AlertSystem) GetActiveAlerts() []*Alert {
	as.mu.RLock()
	defer as.mu.RUnlock()

	var activeAlerts []*Alert
	for _, alert := range as.alerts {
		if !alert.Resolved {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	return activeAlerts
}

// ResolveAlert 解决告警
func (as *AlertSystem) ResolveAlert(alertID string) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	alert, exists := as.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert with ID %s does not exist", alertID)
	}

	if alert.Resolved {
		return fmt.Errorf("alert with ID %s is already resolved", alertID)
	}

	now := time.Now()
	alert.Resolved = true
	alert.ResolvedAt = &now

	return nil
}

// Start 启动告警系统
func (as *AlertSystem) Start(ctx context.Context) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if as.running {
		return fmt.Errorf("alert system is already running")
	}

	as.ctx, as.cancel = context.WithCancel(ctx)
	as.running = true

	go as.checkLoop()

	return nil
}

// Stop 停止告警系统
func (as *AlertSystem) Stop() error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if !as.running {
		return fmt.Errorf("alert system is not running")
	}

	as.cancel()
	as.running = false

	return nil
}

// checkLoop 检查循环
func (as *AlertSystem) checkLoop() {
	ticker := time.NewTicker(as.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-as.ctx.Done():
			return
		case <-ticker.C:
			as.checkRules()
		}
	}
}

// checkRules 检查所有规则
func (as *AlertSystem) checkRules() {
	as.mu.RLock()
	rules := make([]*AlertRule, 0, len(as.rules))
	for _, rule := range as.rules {
		if rule.Enabled {
			rules = append(rules, rule)
		}
	}
	as.mu.RUnlock()

	for _, rule := range rules {
		as.checkRule(rule)
	}
}

// checkRule 检查单个规则
func (as *AlertSystem) checkRule(rule *AlertRule) {
	metric := as.monitor.GetMetric(rule.MetricName)
	if metric == nil {
		return
	}

	value := metric.Value()
	var floatValue float64

	switch v := value.(type) {
	case int:
		floatValue = float64(v)
	case int64:
		floatValue = float64(v)
	case float64:
		floatValue = v
	case float32:
		floatValue = float64(v)
	default:
		return
	}

	// 检查条件
	triggered := false
	switch rule.Condition {
	case ">":
		triggered = floatValue > rule.Threshold
	case "<":
		triggered = floatValue < rule.Threshold
	case ">=":
		triggered = floatValue >= rule.Threshold
	case "<=":
		triggered = floatValue <= rule.Threshold
	case "==":
		triggered = floatValue == rule.Threshold
	case "!=":
		triggered = floatValue != rule.Threshold
	}

	if triggered {
		as.createAlert(rule, floatValue)
	} else {
		as.resolveAlert(rule.ID)
	}
}

// createAlert 创建告警
func (as *AlertSystem) createAlert(rule *AlertRule, value float64) {
	as.mu.Lock()
	defer as.mu.Unlock()

	// 检查是否已存在未解决的告警
	if alert, exists := as.alerts[rule.ID]; exists && !alert.Resolved {
		return
	}

	alert := &Alert{
		ID:         fmt.Sprintf("%s_%d", rule.ID, time.Now().Unix()),
		RuleID:     rule.ID,
		RuleName:   rule.Name,
		Level:      rule.Level,
		Message:    fmt.Sprintf("Metric %s value %.2f %s threshold %.2f", rule.MetricName, value, rule.Condition, rule.Threshold),
		MetricName: rule.MetricName,
		Value:      value,
		Threshold:  rule.Threshold,
		Labels:     rule.Labels,
		Timestamp:  time.Now(),
		Resolved:   false,
	}

	as.alerts[alert.ID] = alert

	// 执行动作
	go as.executeActions(alert, rule.Actions)
}

// resolveAlert 解决告警
func (as *AlertSystem) resolveAlert(ruleID string) {
	as.mu.Lock()
	defer as.mu.Unlock()

	for _, alert := range as.alerts {
		if alert.RuleID == ruleID && !alert.Resolved {
			now := time.Now()
			alert.Resolved = true
			alert.ResolvedAt = &now
		}
	}
}

// executeActions 执行告警动作
func (as *AlertSystem) executeActions(alert *Alert, actionTypes []string) {
	for _, actionType := range actionTypes {
		if action, exists := as.actions[actionType]; exists {
			if err := action.Execute(alert); err != nil {
				// 记录错误但不中断其他动作
				fmt.Printf("Failed to execute alert action %s: %v\n", actionType, err)
			}
		}
	}
}

// LogAlertAction 日志告警动作
type LogAlertAction struct{}

func NewLogAlertAction() *LogAlertAction {
	return &LogAlertAction{}
}

func (la *LogAlertAction) Execute(alert *Alert) error {
	fmt.Printf("[ALERT] %s - %s: %s\n", alert.Level, alert.RuleName, alert.Message)
	return nil
}

func (la *LogAlertAction) GetType() string {
	return "log"
}

func (la *LogAlertAction) GetDescription() string {
	return "Log alert to console"
}

// EmailAlertAction 邮件告警动作
type EmailAlertAction struct {
	SMTPHost  string
	SMTPPort  int
	Username  string
	Password  string
	FromEmail string
	ToEmails  []string
}

func NewEmailAlertAction() *EmailAlertAction {
	return &EmailAlertAction{
		SMTPHost:  "localhost",
		SMTPPort:  587,
		FromEmail: "alerts@example.com",
		ToEmails:  []string{"admin@example.com"},
	}
}

func (ea *EmailAlertAction) Execute(alert *Alert) error {
	// 这里实现邮件发送逻辑
	// 为了简化，这里只是打印
	fmt.Printf("[EMAIL] Sending alert email for %s\n", alert.RuleName)
	return nil
}

func (ea *EmailAlertAction) GetType() string {
	return "email"
}

func (ea *EmailAlertAction) GetDescription() string {
	return "Send alert via email"
}

// WebhookAlertAction Webhook告警动作
type WebhookAlertAction struct {
	URL     string
	Method  string
	Headers map[string]string
}

func NewWebhookAlertAction() *WebhookAlertAction {
	return &WebhookAlertAction{
		URL:     "http://localhost:8080/webhook",
		Method:  "POST",
		Headers: make(map[string]string),
	}
}

func (wa *WebhookAlertAction) Execute(alert *Alert) error {
	// 这里实现webhook调用逻辑
	// 为了简化，这里只是打印
	fmt.Printf("[WEBHOOK] Sending alert to webhook for %s\n", alert.RuleName)
	return nil
}

func (wa *WebhookAlertAction) GetType() string {
	return "webhook"
}

func (wa *WebhookAlertAction) GetDescription() string {
	return "Send alert via webhook"
}
