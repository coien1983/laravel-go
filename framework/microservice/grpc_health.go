package microservice

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// GRPCHealthService gRPC 健康检查服务
type GRPCHealthService struct {
	server     *health.Server
	checks     map[string]HealthCheckFunc
	status     map[string]healthpb.HealthCheckResponse_ServingStatus
	mutex      sync.RWMutex
	interval   time.Duration
	timeout    time.Duration
	stopChan   chan struct{}
	running    bool
	runningMux sync.RWMutex
}

// HealthCheckFunc 健康检查函数
type HealthCheckFunc func(ctx context.Context) error

// HealthStatus 健康状态
type HealthStatus struct {
	Service   string                                     `json:"service"`
	Status    healthpb.HealthCheckResponse_ServingStatus `json:"status"`
	Timestamp time.Time                                  `json:"timestamp"`
	Details   map[string]interface{}                     `json:"details,omitempty"`
	Error     string                                     `json:"error,omitempty"`
}

// HealthConfig 健康检查配置
type HealthConfig struct {
	Interval      time.Duration
	Timeout       time.Duration
	AutoStart     bool
	DefaultStatus healthpb.HealthCheckResponse_ServingStatus
}

// NewHealthConfig 创建健康检查配置
func NewHealthConfig() *HealthConfig {
	return &HealthConfig{
		Interval:      30 * time.Second,
		Timeout:       5 * time.Second,
		AutoStart:     true,
		DefaultStatus: healthpb.HealthCheckResponse_SERVING,
	}
}

// NewGRPCHealthService 创建健康检查服务
func NewGRPCHealthService(config *HealthConfig) *GRPCHealthService {
	if config == nil {
		config = NewHealthConfig()
	}

	hc := &GRPCHealthService{
		server:   health.NewServer(),
		checks:   make(map[string]HealthCheckFunc),
		status:   make(map[string]healthpb.HealthCheckResponse_ServingStatus),
		interval: config.Interval,
		timeout:  config.Timeout,
		stopChan: make(chan struct{}),
	}

	// 设置默认状态
	hc.server.SetServingStatus("", config.DefaultStatus)

	if config.AutoStart {
		hc.Start()
	}

	return hc
}

// RegisterHealthCheck 注册健康检查
func (hc *GRPCHealthService) RegisterHealthCheck(service string, check HealthCheckFunc) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	hc.checks[service] = check
	hc.status[service] = healthpb.HealthCheckResponse_SERVING
	hc.server.SetServingStatus(service, healthpb.HealthCheckResponse_SERVING)
}

// UnregisterHealthCheck 注销健康检查
func (hc *GRPCHealthService) UnregisterHealthCheck(service string) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	delete(hc.checks, service)
	delete(hc.status, service)
	hc.server.SetServingStatus(service, healthpb.HealthCheckResponse_SERVICE_UNKNOWN)
}

// SetStatus 设置服务状态
func (hc *GRPCHealthService) SetStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	hc.status[service] = status
	hc.server.SetServingStatus(service, status)
}

// GetStatus 获取服务状态
func (hc *GRPCHealthService) GetStatus(service string) healthpb.HealthCheckResponse_ServingStatus {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()

	if status, exists := hc.status[service]; exists {
		return status
	}
	return healthpb.HealthCheckResponse_SERVICE_UNKNOWN
}

// GetAllStatus 获取所有服务状态
func (hc *GRPCHealthService) GetAllStatus() map[string]HealthStatus {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()

	statuses := make(map[string]HealthStatus)
	for service, status := range hc.status {
		statuses[service] = HealthStatus{
			Service:   service,
			Status:    status,
			Timestamp: time.Now(),
		}
	}
	return statuses
}

// Check 执行健康检查
func (hc *GRPCHealthService) Check(ctx context.Context, service string) error {
	hc.mutex.RLock()
	check, exists := hc.checks[service]
	hc.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("health check not registered for service: %s", service)
	}

	// 创建超时上下文
	checkCtx, cancel := context.WithTimeout(ctx, hc.timeout)
	defer cancel()

	// 执行检查
	if err := check(checkCtx); err != nil {
		hc.SetStatus(service, healthpb.HealthCheckResponse_NOT_SERVING)
		return err
	}

	hc.SetStatus(service, healthpb.HealthCheckResponse_SERVING)
	return nil
}

// CheckAll 检查所有服务
func (hc *GRPCHealthService) CheckAll(ctx context.Context) map[string]error {
	hc.mutex.RLock()
	services := make([]string, 0, len(hc.checks))
	for service := range hc.checks {
		services = append(services, service)
	}
	hc.mutex.RUnlock()

	errors := make(map[string]error)
	for _, service := range services {
		if err := hc.Check(ctx, service); err != nil {
			errors[service] = err
		}
	}

	return errors
}

// Start 启动健康检查
func (hc *GRPCHealthService) Start() {
	hc.runningMux.Lock()
	if hc.running {
		hc.runningMux.Unlock()
		return
	}
	hc.running = true
	hc.runningMux.Unlock()

	go hc.run()
}

// Stop 停止健康检查
func (hc *GRPCHealthService) Stop() {
	hc.runningMux.Lock()
	if !hc.running {
		hc.runningMux.Unlock()
		return
	}
	hc.running = false
	hc.runningMux.Unlock()

	close(hc.stopChan)
}

// IsRunning 检查是否运行
func (hc *GRPCHealthService) IsRunning() bool {
	hc.runningMux.RLock()
	defer hc.runningMux.RUnlock()
	return hc.running
}

// run 运行健康检查循环
func (hc *GRPCHealthService) run() {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hc.CheckAll(context.Background())
		case <-hc.stopChan:
			return
		}
	}
}

// GetServer 获取健康检查服务器
func (hc *GRPCHealthService) GetServer() *health.Server {
	return hc.server
}

// RegisterWithGRPC 注册到 gRPC 服务器
func (hc *GRPCHealthService) RegisterWithGRPC(server *grpc.Server) {
	healthpb.RegisterHealthServer(server, hc.server)
}

// HealthCheckService gRPC 健康检查服务实现
type HealthCheckService struct {
	healthpb.UnimplementedHealthServer
	checker *GRPCHealthService
}

// NewHealthCheckService 创建健康检查服务
func NewHealthCheckService(checker *GRPCHealthService) *HealthCheckService {
	return &HealthCheckService{
		checker: checker,
	}
}

// Check 实现健康检查
func (s *HealthCheckService) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	service := req.Service
	if service == "" {
		service = ""
	}

	status := s.checker.GetStatus(service)
	return &healthpb.HealthCheckResponse{
		Status: status,
	}, nil
}

// Watch 实现健康检查监听
func (s *HealthCheckService) Watch(req *healthpb.HealthCheckRequest, stream healthpb.Health_WatchServer) error {
	service := req.Service
	if service == "" {
		service = ""
	}

	// 发送初始状态
	status := s.checker.GetStatus(service)
	if err := stream.Send(&healthpb.HealthCheckResponse{Status: status}); err != nil {
		return err
	}

	// 监听状态变化
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentStatus := s.checker.GetStatus(service)
			if currentStatus != status {
				status = currentStatus
				if err := stream.Send(&healthpb.HealthCheckResponse{Status: status}); err != nil {
					return err
				}
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

// Built-in Health Checks

// DatabaseHealthCheck 数据库健康检查
func DatabaseHealthCheck(db interface{}) HealthCheckFunc {
	return func(ctx context.Context) error {
		// 这里需要根据具体的数据库类型来实现
		// 例如：db.PingContext(ctx)
		return nil
	}
}

// RedisHealthCheck Redis 健康检查
func RedisHealthCheck(redis interface{}) HealthCheckFunc {
	return func(ctx context.Context) error {
		// 这里需要根据具体的 Redis 客户端来实现
		// 例如：redis.Ping(ctx)
		return nil
	}
}

// HTTPHealthCheck HTTP 健康检查
func HTTPHealthCheck(url string, timeout time.Duration) HealthCheckFunc {
	return func(ctx context.Context) error {
		// 这里需要实现 HTTP 健康检查
		// 例如：http.Get(url)
		return nil
	}
}

// GRPCHealthCheck gRPC 健康检查
func GRPCHealthCheck(address string, timeout time.Duration) HealthCheckFunc {
	return func(ctx context.Context) error {
		// 创建 gRPC 连接
		conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()

		// 创建健康检查客户端
		client := healthpb.NewHealthClient(conn)

		// 执行健康检查
		_, err = client.Check(ctx, &healthpb.HealthCheckRequest{})
		return err
	}
}

// FileSystemHealthCheck 文件系统健康检查
func FileSystemHealthCheck(path string) HealthCheckFunc {
	return func(ctx context.Context) error {
		// 这里需要实现文件系统健康检查
		// 例如：检查磁盘空间、文件权限等
		return nil
	}
}

// MemoryHealthCheck 内存健康检查
func MemoryHealthCheck(threshold uint64) HealthCheckFunc {
	return func(ctx context.Context) error {
		// 这里需要实现内存健康检查
		// 例如：检查内存使用量是否超过阈值
		return nil
	}
}

// CPUHealthCheck CPU 健康检查
func CPUHealthCheck(threshold float64) HealthCheckFunc {
	return func(ctx context.Context) error {
		// 这里需要实现 CPU 健康检查
		// 例如：检查 CPU 使用率是否超过阈值
		return nil
	}
}

// HealthMonitor 健康监控器
type HealthMonitor struct {
	checker   *GRPCHealthService
	metrics   *HealthMetrics
	alerts    []HealthAlert
	alertChan chan HealthAlert
	mutex     sync.RWMutex
}

// HealthMetrics 健康指标
type HealthMetrics struct {
	TotalChecks         int64
	SuccessfulChecks    int64
	FailedChecks        int64
	LastCheckTime       time.Time
	AverageResponseTime time.Duration
}

// HealthAlert 健康告警
type HealthAlert struct {
	Service   string
	Status    healthpb.HealthCheckResponse_ServingStatus
	Message   string
	Timestamp time.Time
	Severity  string
}

// NewHealthMonitor 创建健康监控器
func NewHealthMonitor(checker *GRPCHealthService) *HealthMonitor {
	return &HealthMonitor{
		checker:   checker,
		metrics:   &HealthMetrics{},
		alerts:    make([]HealthAlert, 0),
		alertChan: make(chan HealthAlert, 100),
	}
}

// Start 启动监控
func (hm *HealthMonitor) Start() {
	go hm.monitor()
}

// Stop 停止监控
func (hm *HealthMonitor) Stop() {
	close(hm.alertChan)
}

// monitor 监控循环
func (hm *HealthMonitor) monitor() {
	for alert := range hm.alertChan {
		hm.mutex.Lock()
		hm.alerts = append(hm.alerts, alert)
		hm.mutex.Unlock()

		// 处理告警
		hm.handleAlert(alert)
	}
}

// handleAlert 处理告警
func (hm *HealthMonitor) handleAlert(alert HealthAlert) {
	// 这里可以实现告警处理逻辑
	// 例如：发送邮件、短信、Slack 通知等
}

// GetMetrics 获取指标
func (hm *HealthMonitor) GetMetrics() *HealthMetrics {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	return hm.metrics
}

// GetAlerts 获取告警
func (hm *HealthMonitor) GetAlerts() []HealthAlert {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	alerts := make([]HealthAlert, len(hm.alerts))
	copy(alerts, hm.alerts)
	return alerts
}

// ClearAlerts 清除告警
func (hm *HealthMonitor) ClearAlerts() {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	hm.alerts = make([]HealthAlert, 0)
}
