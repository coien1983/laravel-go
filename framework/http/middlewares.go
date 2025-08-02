package http

import (
	"fmt"
	"time"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	authService interface{} // 认证服务接口
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(authService interface{}) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Handle 实现 Middleware 接口
func (am *AuthMiddleware) Handle(request Request, next Next) Response {
	// 检查认证头
	authHeader := request.Header("Authorization")
	if authHeader == "" {
		return NewResponse(401, "Unauthorized")
	}

	// 这里可以添加具体的认证逻辑
	// 例如：验证 JWT token、检查 session 等

	// 如果认证通过，继续处理
	return next(request)
}

// CORSMiddleware CORS 中间件
type CORSMiddleware struct {
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
}

// NewCORSMiddleware 创建 CORS 中间件
func NewCORSMiddleware(allowedOrigins, allowedMethods, allowedHeaders []string) *CORSMiddleware {
	return &CORSMiddleware{
		allowedOrigins: allowedOrigins,
		allowedMethods: allowedMethods,
		allowedHeaders: allowedHeaders,
	}
}

// Handle 实现 Middleware 接口
func (cm *CORSMiddleware) Handle(request Request, next Next) Response {
	response := next(request)

	// 设置 CORS 头
	origin := request.Header("Origin")
	if cm.isOriginAllowed(origin) {
		// 注意：这里需要扩展 Response 接口以支持设置头部
		// 暂时跳过头部设置，在实际使用时需要修改 Response 接口
	}

	return response
}

// isOriginAllowed 检查源是否允许
func (cm *CORSMiddleware) isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range cm.allowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}
	return false
}

// getAllowedMethodsString 获取允许的方法字符串
func (cm *CORSMiddleware) getAllowedMethodsString() string {
	if len(cm.allowedMethods) == 0 {
		return "GET, POST, PUT, DELETE, OPTIONS"
	}

	result := ""
	for i, method := range cm.allowedMethods {
		if i > 0 {
			result += ", "
		}
		result += method
	}
	return result
}

// getAllowedHeadersString 获取允许的头字符串
func (cm *CORSMiddleware) getAllowedHeadersString() string {
	if len(cm.allowedHeaders) == 0 {
		return "Content-Type, Authorization"
	}

	result := ""
	for i, header := range cm.allowedHeaders {
		if i > 0 {
			result += ", "
		}
		result += header
	}
	return result
}

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct {
	logger interface{} // 日志服务接口
}

// NewLoggingMiddleware 创建日志中间件
func NewLoggingMiddleware(logger interface{}) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Handle 实现 Middleware 接口
func (lm *LoggingMiddleware) Handle(request Request, next Next) Response {
	startTime := time.Now()

	// 记录请求开始
	lm.logRequest(request, "Request started")

	// 处理请求
	response := next(request)

	// 计算处理时间
	duration := time.Since(startTime)

	// 记录请求结束
	lm.logRequest(request, fmt.Sprintf("Request completed in %v", duration))

	return response
}

// logRequest 记录请求日志
func (lm *LoggingMiddleware) logRequest(request Request, message string) {
	// 这里可以添加具体的日志记录逻辑
	// 例如：记录到文件、数据库或日志服务
	fmt.Printf("[%s] %s %s - %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		request.Method(),
		request.Path(),
		message)
}

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
	limit  int           // 限制次数
	window time.Duration // 时间窗口
	store  interface{}   // 存储接口（Redis、内存等）
}

// NewRateLimitMiddleware 创建限流中间件
func NewRateLimitMiddleware(limit int, window time.Duration, store interface{}) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limit:  limit,
		window: window,
		store:  store,
	}
}

// Handle 实现 Middleware 接口
func (rlm *RateLimitMiddleware) Handle(request Request, next Next) Response {
	// 获取客户端标识（IP、用户ID等）
	clientID := rlm.getClientID(request)

	// 检查是否超过限制
	if rlm.isRateLimited(clientID) {
		return NewResponse(429, "Too Many Requests")
	}

	// 增加计数
	rlm.incrementCounter(clientID)

	// 继续处理
	return next(request)
}

// getClientID 获取客户端标识
func (rlm *RateLimitMiddleware) getClientID(request Request) string {
	// 优先使用用户ID，否则使用IP地址
	userID := request.Param("user_id")
	if userID != "" {
		return "user:" + userID
	}
	return "ip:" + request.IP()
}

// isRateLimited 检查是否被限流
func (rlm *RateLimitMiddleware) isRateLimited(clientID string) bool {
	// 这里可以添加具体的限流逻辑
	// 例如：检查 Redis 中的计数器
	return false
}

// incrementCounter 增加计数器
func (rlm *RateLimitMiddleware) incrementCounter(clientID string) {
	// 这里可以添加具体的计数逻辑
	// 例如：在 Redis 中增加计数器
}

// CacheMiddleware 缓存中间件
type CacheMiddleware struct {
	cache interface{}   // 缓存服务接口
	ttl   time.Duration // 缓存时间
}

// NewCacheMiddleware 创建缓存中间件
func NewCacheMiddleware(cache interface{}, ttl time.Duration) *CacheMiddleware {
	return &CacheMiddleware{
		cache: cache,
		ttl:   ttl,
	}
}

// Handle 实现 Middleware 接口
func (cm *CacheMiddleware) Handle(request Request, next Next) Response {
	// 只缓存 GET 请求
	if request.Method() != "GET" {
		return next(request)
	}

	// 生成缓存键
	cacheKey := cm.generateCacheKey(request)

	// 尝试从缓存获取
	if cachedResponse := cm.getFromCache(cacheKey); cachedResponse != nil {
		return cachedResponse
	}

	// 处理请求
	response := next(request)

	// 缓存响应
	cm.setCache(cacheKey, response)

	return response
}

// generateCacheKey 生成缓存键
func (cm *CacheMiddleware) generateCacheKey(request Request) string {
	return fmt.Sprintf("cache:%s:%s", request.Method(), request.Path())
}

// getFromCache 从缓存获取
func (cm *CacheMiddleware) getFromCache(key string) Response {
	// 这里可以添加具体的缓存获取逻辑
	return nil
}

// setCache 设置缓存
func (cm *CacheMiddleware) setCache(key string, response Response) {
	// 这里可以添加具体的缓存设置逻辑
}

// SecurityMiddleware 安全中间件
type SecurityMiddleware struct {
	enableCSRF bool
	enableXSS  bool
	enableHSTS bool
}

// NewSecurityMiddleware 创建安全中间件
func NewSecurityMiddleware(enableCSRF, enableXSS, enableHSTS bool) *SecurityMiddleware {
	return &SecurityMiddleware{
		enableCSRF: enableCSRF,
		enableXSS:  enableXSS,
		enableHSTS: enableHSTS,
	}
}

// Handle 实现 Middleware 接口
func (sm *SecurityMiddleware) Handle(request Request, next Next) Response {
	response := next(request)

	// 设置安全头
	sm.setSecurityHeaders(response)

	return response
}

// setSecurityHeaders 设置安全头
func (sm *SecurityMiddleware) setSecurityHeaders(response Response) {
	// 设置 XSS 保护
	if sm.enableXSS {
		response.SetHeader("X-XSS-Protection", "1; mode=block")
	}

	// 设置内容类型选项
	response.SetHeader("X-Content-Type-Options", "nosniff")

	// 设置框架选项
	response.SetHeader("X-Frame-Options", "DENY")

	// 设置 HSTS
	if sm.enableHSTS {
		response.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}

	// 设置引用策略
	response.SetHeader("Referrer-Policy", "strict-origin-when-cross-origin")
}
