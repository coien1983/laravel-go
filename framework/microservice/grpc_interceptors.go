package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"crypto/sha256"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor gRPC 日志拦截器
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// 记录请求开始
		log.Printf("gRPC Request: %s", info.FullMethod)

		// 调用处理器
		resp, err := handler(ctx, req)

		// 记录响应时间和结果
		duration := time.Since(start)
		if err != nil {
			log.Printf("gRPC Error: %s - %v (duration: %v)", info.FullMethod, err, duration)
		} else {
			log.Printf("gRPC Success: %s (duration: %v)", info.FullMethod, duration)
		}

		return resp, err
	}
}

// AuthInterceptor gRPC 认证拦截器
func AuthInterceptor(authFunc func(ctx context.Context) error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 跳过健康检查的认证
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		// 执行认证
		if err := authFunc(ctx); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
		}

		return handler(ctx, req)
	}
}

// TokenAuthInterceptor 基于 Token 的认证拦截器
func TokenAuthInterceptor(validTokens map[string]bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 跳过健康检查的认证
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		// 获取元数据
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		// 获取 token
		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
		}

		token := tokens[0]
		if !validTokens[token] {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}

// RateLimitInterceptor gRPC 限流拦截器
func RateLimitInterceptor(limiter RateLimiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 检查限流
		if !limiter.Allow(info.FullMethod) {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}

// CircuitBreakerInterceptor gRPC 熔断器拦截器
func CircuitBreakerInterceptor(circuitBreaker CircuitBreaker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var result interface{}
		var resultErr error

		err := circuitBreaker.Execute(ctx, func() error {
			result, resultErr = handler(ctx, req)
			return resultErr
		})

		if err != nil {
			return nil, err
		}

		return result, resultErr
	}
}

// MetricsInterceptor gRPC 指标拦截器
func MetricsInterceptor(metrics MetricsCollector) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// 记录请求计数
		metrics.Increment("grpc_requests_total", map[string]string{
			"method": info.FullMethod,
		})

		// 调用处理器
		resp, err := handler(ctx, req)

		// 记录响应时间
		duration := time.Since(start)
		metrics.Histogram("grpc_request_duration_seconds", duration.Seconds(), map[string]string{
			"method": info.FullMethod,
		})

		// 记录错误计数
		if err != nil {
			metrics.Increment("grpc_errors_total", map[string]string{
				"method": info.FullMethod,
				"error":  err.Error(),
			})
		}

		return resp, err
	}
}

// TracingInterceptor gRPC 追踪拦截器
func TracingInterceptor(tracer Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 创建 span
		span := tracer.StartSpan("grpc_server", info.FullMethod)
		defer span.Finish()

		// 设置 span 标签
		span.SetTag("grpc.method", info.FullMethod)
		span.SetTag("grpc.service", info.Server)

		// 将 span 添加到上下文
		ctx = tracer.ContextWithSpan(ctx, span)

		// 调用处理器
		resp, err := handler(ctx, req)

		// 记录错误
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}

		return resp, err
	}
}

// ValidationInterceptor gRPC 验证拦截器
func ValidationInterceptor(validator func(interface{}) error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 验证请求
		if err := validator(req); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
		}

		return handler(ctx, req)
	}
}

// TimeoutInterceptor gRPC 超时拦截器
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 创建超时上下文
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// 创建响应通道
		respChan := make(chan interface{}, 1)
		errChan := make(chan error, 1)

		// 异步执行处理器
		go func() {
			resp, err := handler(ctx, req)
			if err != nil {
				errChan <- err
			} else {
				respChan <- resp
			}
		}()

		// 等待结果或超时
		select {
		case resp := <-respChan:
			return resp, nil
		case err := <-errChan:
			return nil, err
		case <-ctx.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "request timeout")
		}
	}
}

// RecoveryInterceptor gRPC 恢复拦截器
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("gRPC panic recovered: %v", r)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// MetadataInterceptor gRPC 元数据拦截器
func MetadataInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 获取元数据
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			// 添加自定义元数据
			md.Set("server-timestamp", time.Now().Format(time.RFC3339))
			md.Set("server-method", info.FullMethod)

			// 创建新的上下文
			ctx = metadata.NewIncomingContext(ctx, md)
		}

		return handler(ctx, req)
	}
}

// StreamLoggingInterceptor gRPC 流日志拦截器
func StreamLoggingInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		log.Printf("gRPC Stream Request: %s", info.FullMethod)

		err := handler(srv, ss)

		duration := time.Since(start)
		if err != nil {
			log.Printf("gRPC Stream Error: %s - %v (duration: %v)", info.FullMethod, err, duration)
		} else {
			log.Printf("gRPC Stream Success: %s (duration: %v)", info.FullMethod, duration)
		}

		return err
	}
}

// StreamAuthInterceptor gRPC 流认证拦截器
func StreamAuthInterceptor(authFunc func(ctx context.Context) error) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 跳过健康检查的认证
		if info.FullMethod == "/grpc.health.v1.Health/Watch" {
			return handler(srv, ss)
		}

		// 执行认证
		if err := authFunc(ss.Context()); err != nil {
			return status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
		}

		return handler(srv, ss)
	}
}

// StreamRateLimitInterceptor gRPC 流限流拦截器
func StreamRateLimitInterceptor(limiter RateLimiter) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 检查限流
		if !limiter.Allow(info.FullMethod) {
			return status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(srv, ss)
	}
}

// StreamRecoveryInterceptor gRPC 流恢复拦截器
func StreamRecoveryInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("gRPC stream panic recovered: %v", r)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(srv, ss)
	}
}

// CacheInterceptor gRPC 缓存拦截器
func CacheInterceptor(cache Cache, ttl time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 跳过非GET方法的缓存
		if !isCacheableMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// 生成缓存键
		cacheKey := generateCacheKey(info.FullMethod, req)

		// 尝试从缓存获取
		if cached, exists := cache.Get(cacheKey); exists {
			return cached, nil
		}

		// 调用处理器
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		// 缓存响应
		cache.Set(cacheKey, resp, ttl)

		return resp, nil
	}
}

// CompressionInterceptor gRPC 压缩拦截器
func CompressionInterceptor(compressor Compressor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 检查是否需要压缩
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if acceptEncoding := md.Get("accept-encoding"); len(acceptEncoding) > 0 {
				// 处理压缩逻辑
				ctx = context.WithValue(ctx, "compression", acceptEncoding[0])
			}
		}

		// 调用处理器
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		// 如果需要压缩响应
		if compression := ctx.Value("compression"); compression != nil {
			if compressed, err := compressor.Compress(resp); err == nil {
				return compressed, nil
			}
		}

		return resp, nil
	}
}

// RetryInterceptor gRPC 重试拦截器
func RetryInterceptor(maxRetries int, backoff BackoffStrategy) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var lastErr error

		for attempt := 0; attempt <= maxRetries; attempt++ {
			// 调用处理器
			resp, err := handler(ctx, req)
			if err == nil {
				return resp, nil
			}

			lastErr = err

			// 检查是否应该重试
			if !isRetryableError(err) || attempt == maxRetries {
				break
			}

			// 等待后重试
			delay := backoff.GetDelay(attempt)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		return nil, lastErr
	}
}

// LoadBalancingInterceptor gRPC 负载均衡拦截器
func LoadBalancingInterceptor(loadBalancer LoadBalancer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 获取可用服务
		services := getAvailableServices(ctx, info.FullMethod)
		if len(services) == 0 {
			return nil, status.Errorf(codes.Unavailable, "no available services")
		}

		// 选择服务
		selected := loadBalancer.Select(services)
		if selected == nil {
			return nil, status.Errorf(codes.Unavailable, "no healthy service available")
		}

		// 将选中的服务信息添加到上下文
		ctx = context.WithValue(ctx, "selected_service", selected)

		return handler(ctx, req)
	}
}

// ThrottlingInterceptor gRPC 节流拦截器
func ThrottlingInterceptor(throttler Throttler) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 检查节流
		if !throttler.Allow(info.FullMethod) {
			return nil, status.Errorf(codes.ResourceExhausted, "throttling limit exceeded")
		}

		return handler(ctx, req)
	}
}

// SecurityInterceptor gRPC 安全拦截器
func SecurityInterceptor(security SecurityChecker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 执行安全检查
		if err := security.Check(ctx, req, info.FullMethod); err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "security check failed: %v", err)
		}

		return handler(ctx, req)
	}
}

// AuditInterceptor gRPC 审计拦截器
func AuditInterceptor(auditor Auditor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// 记录请求审计
		auditor.LogRequest(ctx, info.FullMethod, req)

		// 调用处理器
		resp, err := handler(ctx, req)

		// 记录响应审计
		duration := time.Since(start)
		auditor.LogResponse(ctx, info.FullMethod, resp, err, duration)

		return resp, err
	}
}

// PerformanceInterceptor gRPC 性能拦截器
func PerformanceInterceptor(perfMonitor PerformanceMonitor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// 记录开始
		perfMonitor.RecordStart(info.FullMethod)

		// 调用处理器
		resp, err := handler(ctx, req)

		// 记录结束
		duration := time.Since(start)
		perfMonitor.RecordEnd(info.FullMethod, duration, err)

		return resp, err
	}
}

// 辅助函数

// isCacheableMethod 检查方法是否可缓存
func isCacheableMethod(method string) bool {
	// 只缓存GET方法
	return strings.Contains(method, "Get") || strings.Contains(method, "List")
}

// generateCacheKey 生成缓存键
func generateCacheKey(method string, req interface{}) string {
	// 简单的缓存键生成，实际应用中可能需要更复杂的逻辑
	data, _ := json.Marshal(req)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%s:%x", method, hash[:8])
}

// isRetryableError 检查错误是否可重试
func isRetryableError(err error) bool {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
			return true
		}
	}
	return false
}

// getAvailableServices 获取可用服务
func getAvailableServices(ctx context.Context, method string) []*ServiceInfo {
	// 这里需要从服务发现中获取可用服务
	// 实际实现中可能需要注入服务发现客户端
	return nil
}

// 接口定义

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(method string) bool
}

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
	Increment(name string, labels map[string]string)
	Histogram(name string, value float64, labels map[string]string)
}

// Tracer 追踪器接口
type Tracer interface {
	StartSpan(operationName, method string) Span
	ContextWithSpan(ctx context.Context, span Span) context.Context
}

// Span 追踪 span 接口
type Span interface {
	SetTag(key string, value interface{})
	LogKV(keyValues ...interface{})
	Finish()
}

// Cache 缓存接口
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
}

// Compressor 压缩器接口
type Compressor interface {
	Compress(data interface{}) (interface{}, error)
	Decompress(data interface{}) (interface{}, error)
}

// BackoffStrategy 退避策略接口
type BackoffStrategy interface {
	GetDelay(attempt int) time.Duration
}

// Throttler 节流器接口
type Throttler interface {
	Allow(method string) bool
}

// SecurityChecker 安全检查器接口
type SecurityChecker interface {
	Check(ctx context.Context, req interface{}, method string) error
}

// Auditor 审计器接口
type Auditor interface {
	LogRequest(ctx context.Context, method string, req interface{})
	LogResponse(ctx context.Context, method string, resp interface{}, err error, duration time.Duration)
}

// PerformanceMonitor 性能监控器接口
type PerformanceMonitor interface {
	RecordStart(method string)
	RecordEnd(method string, duration time.Duration, err error)
}

// 简单实现

// SimpleRateLimiter 简单限流器
type SimpleRateLimiter struct {
	limits map[string]int
	counts map[string]int
}

// NewSimpleRateLimiter 创建简单限流器
func NewSimpleRateLimiter(limits map[string]int) *SimpleRateLimiter {
	return &SimpleRateLimiter{
		limits: limits,
		counts: make(map[string]int),
	}
}

// Allow 检查是否允许请求
func (r *SimpleRateLimiter) Allow(method string) bool {
	limit, exists := r.limits[method]
	if !exists {
		return true // 没有限制
	}

	if r.counts[method] >= limit {
		return false
	}

	r.counts[method]++
	return true
}

// SimpleMetricsCollector 简单指标收集器
type SimpleMetricsCollector struct {
	metrics map[string]float64
}

// NewSimpleMetricsCollector 创建简单指标收集器
func NewSimpleMetricsCollector() *SimpleMetricsCollector {
	return &SimpleMetricsCollector{
		metrics: make(map[string]float64),
	}
}

// Increment 增加计数
func (m *SimpleMetricsCollector) Increment(name string, labels map[string]string) {
	key := fmt.Sprintf("%s_%v", name, labels)
	m.metrics[key]++
}

// Histogram 记录直方图
func (m *SimpleMetricsCollector) Histogram(name string, value float64, labels map[string]string) {
	key := fmt.Sprintf("%s_%v", name, labels)
	m.metrics[key] = value
}

// GetMetrics 获取指标
func (m *SimpleMetricsCollector) GetMetrics() map[string]float64 {
	return m.metrics
}

// SimpleCache 简单缓存实现
type SimpleCache struct {
	data  map[string]cacheItem
	mutex sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewSimpleCache 创建简单缓存
func NewSimpleCache() *SimpleCache {
	cache := &SimpleCache{
		data: make(map[string]cacheItem),
	}

	// 启动清理协程
	go cache.cleanup()

	return cache
}

// Get 获取缓存
func (c *SimpleCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		delete(c.data, key)
		return nil, false
	}

	return item.value, true
}

// Set 设置缓存
func (c *SimpleCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete 删除缓存
func (c *SimpleCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}

// Clear 清空缓存
func (c *SimpleCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = make(map[string]cacheItem)
}

// cleanup 清理过期缓存
func (c *SimpleCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.data {
			if now.After(item.expiration) {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}
}

// ExponentialBackoff 指数退避策略
type ExponentialBackoff struct {
	baseDelay time.Duration
	maxDelay  time.Duration
}

// NewExponentialBackoff 创建指数退避策略
func NewExponentialBackoff(baseDelay, maxDelay time.Duration) *ExponentialBackoff {
	return &ExponentialBackoff{
		baseDelay: baseDelay,
		maxDelay:  maxDelay,
	}
}

// GetDelay 获取延迟时间
func (eb *ExponentialBackoff) GetDelay(attempt int) time.Duration {
	delay := eb.baseDelay * time.Duration(1<<attempt)
	if delay > eb.maxDelay {
		delay = eb.maxDelay
	}
	return delay
}

// SimpleThrottler 简单节流器
type SimpleThrottler struct {
	limits map[string]int
	counts map[string]int
	mutex  sync.RWMutex
}

// NewSimpleThrottler 创建简单节流器
func NewSimpleThrottler(limits map[string]int) *SimpleThrottler {
	return &SimpleThrottler{
		limits: limits,
		counts: make(map[string]int),
	}
}

// Allow 检查是否允许
func (st *SimpleThrottler) Allow(method string) bool {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	limit, exists := st.limits[method]
	if !exists {
		return true
	}

	if st.counts[method] >= limit {
		return false
	}

	st.counts[method]++
	return true
}

// Reset 重置计数
func (st *SimpleThrottler) Reset() {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	st.counts = make(map[string]int)
}

// SimpleSecurityChecker 简单安全检查器
type SimpleSecurityChecker struct {
	allowedMethods map[string]bool
}

// NewSimpleSecurityChecker 创建简单安全检查器
func NewSimpleSecurityChecker(allowedMethods map[string]bool) *SimpleSecurityChecker {
	return &SimpleSecurityChecker{
		allowedMethods: allowedMethods,
	}
}

// Check 执行安全检查
func (ssc *SimpleSecurityChecker) Check(ctx context.Context, req interface{}, method string) error {
	if !ssc.allowedMethods[method] {
		return fmt.Errorf("method %s is not allowed", method)
	}
	return nil
}

// SimpleAuditor 简单审计器
type SimpleAuditor struct {
	logger *log.Logger
}

// NewSimpleAuditor 创建简单审计器
func NewSimpleAuditor(logger *log.Logger) *SimpleAuditor {
	return &SimpleAuditor{
		logger: logger,
	}
}

// LogRequest 记录请求
func (sa *SimpleAuditor) LogRequest(ctx context.Context, method string, req interface{}) {
	sa.logger.Printf("AUDIT REQUEST: %s, Method: %s", time.Now().Format(time.RFC3339), method)
}

// LogResponse 记录响应
func (sa *SimpleAuditor) LogResponse(ctx context.Context, method string, resp interface{}, err error, duration time.Duration) {
	status := "SUCCESS"
	if err != nil {
		status = "ERROR"
	}
	sa.logger.Printf("AUDIT RESPONSE: %s, Method: %s, Status: %s, Duration: %v",
		time.Now().Format(time.RFC3339), method, status, duration)
}

// SimplePerformanceMonitor 简单性能监控器
type SimplePerformanceMonitor struct {
	metrics map[string]*PerformanceMetrics
	mutex   sync.RWMutex
}

type PerformanceMetrics struct {
	Count     int64
	TotalTime time.Duration
	Errors    int64
	MinTime   time.Duration
	MaxTime   time.Duration
}

// NewSimplePerformanceMonitor 创建简单性能监控器
func NewSimplePerformanceMonitor() *SimplePerformanceMonitor {
	return &SimplePerformanceMonitor{
		metrics: make(map[string]*PerformanceMetrics),
	}
}

// RecordStart 记录开始
func (spm *SimplePerformanceMonitor) RecordStart(method string) {
	// 可以在这里记录开始时间
}

// RecordEnd 记录结束
func (spm *SimplePerformanceMonitor) RecordEnd(method string, duration time.Duration, err error) {
	spm.mutex.Lock()
	defer spm.mutex.Unlock()

	if _, exists := spm.metrics[method]; !exists {
		spm.metrics[method] = &PerformanceMetrics{}
	}

	metrics := spm.metrics[method]
	metrics.Count++
	metrics.TotalTime += duration

	if err != nil {
		metrics.Errors++
	}

	if metrics.MinTime == 0 || duration < metrics.MinTime {
		metrics.MinTime = duration
	}

	if duration > metrics.MaxTime {
		metrics.MaxTime = duration
	}
}

// GetMetrics 获取指标
func (spm *SimplePerformanceMonitor) GetMetrics() map[string]*PerformanceMetrics {
	spm.mutex.RLock()
	defer spm.mutex.RUnlock()

	metrics := make(map[string]*PerformanceMetrics)
	for k, v := range spm.metrics {
		metrics[k] = v
	}
	return metrics
}
