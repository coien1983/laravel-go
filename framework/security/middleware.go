package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// SecurityMiddleware 安全中间件接口
type SecurityMiddleware interface {
	// Process 处理请求
	Process(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	// GetName 获取中间件名称
	GetName() string
}

// CSRFMiddleware CSRF 保护中间件
type CSRFMiddleware struct {
	enabled     bool
	tokenLength int
	tokenHeader string
	tokenParam  string
	exemptPaths []string
}

// NewCSRFMiddleware 创建 CSRF 中间件
func NewCSRFMiddleware(enabled bool) *CSRFMiddleware {
	return &CSRFMiddleware{
		enabled:     enabled,
		tokenLength: 32,
		tokenHeader: "X-CSRF-Token",
		tokenParam:  "_token",
		exemptPaths: []string{},
	}
}

// SetTokenLength 设置令牌长度
func (cm *CSRFMiddleware) SetTokenLength(length int) *CSRFMiddleware {
	cm.tokenLength = length
	return cm
}

// SetTokenHeader 设置令牌头
func (cm *CSRFMiddleware) SetTokenHeader(header string) *CSRFMiddleware {
	cm.tokenHeader = header
	return cm
}

// SetTokenParam 设置令牌参数
func (cm *CSRFMiddleware) SetTokenParam(param string) *CSRFMiddleware {
	cm.tokenParam = param
	return cm
}

// AddExemptPath 添加豁免路径
func (cm *CSRFMiddleware) AddExemptPath(path string) *CSRFMiddleware {
	cm.exemptPaths = append(cm.exemptPaths, path)
	return cm
}

// GetName 获取中间件名称
func (cm *CSRFMiddleware) GetName() string {
	return "csrf"
}

// Process 处理请求
func (cm *CSRFMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !cm.enabled {
		next(w, r)
		return
	}

	// 检查是否为豁免路径
	if cm.isExemptPath(r.URL.Path) {
		next(w, r)
		return
	}

	// 只对修改数据的请求进行 CSRF 检查
	if !cm.isModifyingRequest(r.Method) {
		next(w, r)
		return
	}

	// 获取令牌
	token := cm.getToken(r)
	if token == "" {
		http.Error(w, "CSRF token missing", http.StatusForbidden)
		return
	}

	// 验证令牌
	if !cm.validateToken(r, token) {
		http.Error(w, "CSRF token invalid", http.StatusForbidden)
		return
	}

	next(w, r)
}

// isExemptPath 检查是否为豁免路径
func (cm *CSRFMiddleware) isExemptPath(path string) bool {
	for _, exemptPath := range cm.exemptPaths {
		if strings.HasPrefix(path, exemptPath) {
			return true
		}
	}
	return false
}

// isModifyingRequest 检查是否为修改数据的请求
func (cm *CSRFMiddleware) isModifyingRequest(method string) bool {
	modifyingMethods := []string{"POST", "PUT", "PATCH", "DELETE"}
	for _, m := range modifyingMethods {
		if method == m {
			return true
		}
	}
	return false
}

// getToken 获取令牌
func (cm *CSRFMiddleware) getToken(r *http.Request) string {
	// 从头部获取
	if token := r.Header.Get(cm.tokenHeader); token != "" {
		return token
	}

	// 从表单参数获取
	if token := r.FormValue(cm.tokenParam); token != "" {
		return token
	}

	// 从查询参数获取
	if token := r.URL.Query().Get(cm.tokenParam); token != "" {
		return token
	}

	return ""
}

// validateToken 验证令牌
func (cm *CSRFMiddleware) validateToken(r *http.Request, token string) bool {
	// 从会话中获取存储的令牌
	sessionToken := cm.getSessionToken(r)
	if sessionToken == "" {
		return false
	}

	// 使用恒定时间比较防止时序攻击
	return subtle.ConstantTimeCompare([]byte(token), []byte(sessionToken)) == 1
}

// getSessionToken 从会话获取令牌
func (cm *CSRFMiddleware) getSessionToken(r *http.Request) string {
	// 这里应该从实际的会话中获取令牌
	// 暂时返回一个固定值用于演示
	return "session_token"
}

// GenerateToken 生成 CSRF 令牌
func (cm *CSRFMiddleware) GenerateToken() (string, error) {
	bytes := make([]byte, cm.tokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// XSSMiddleware XSS 防护中间件
type XSSMiddleware struct {
	enabled     bool
	contentType string
}

// NewXSSMiddleware 创建 XSS 中间件
func NewXSSMiddleware(enabled bool) *XSSMiddleware {
	return &XSSMiddleware{
		enabled:     enabled,
		contentType: "text/html",
	}
}

// GetName 获取中间件名称
func (xm *XSSMiddleware) GetName() string {
	return "xss"
}

// Process 处理请求
func (xm *XSSMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !xm.enabled {
		next(w, r)
		return
	}

	// 设置 XSS 防护头
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	next(w, r)
}

// SQLInjectionMiddleware SQL 注入防护中间件
type SQLInjectionMiddleware struct {
	enabled     bool
	blockedSQL  []string
	blockedPath []string
}

// NewSQLInjectionMiddleware 创建 SQL 注入防护中间件
func NewSQLInjectionMiddleware(enabled bool) *SQLInjectionMiddleware {
	return &SQLInjectionMiddleware{
		enabled: enabled,
		blockedSQL: []string{
			"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
			"UNION", "EXEC", "EXECUTE", "SCRIPT", "EVAL", "EXPRESSION",
		},
		blockedPath: []string{},
	}
}

// AddBlockedSQL 添加被阻止的 SQL 关键字
func (sim *SQLInjectionMiddleware) AddBlockedSQL(keyword string) *SQLInjectionMiddleware {
	sim.blockedSQL = append(sim.blockedSQL, keyword)
	return sim
}

// AddBlockedPath 添加被阻止的路径
func (sim *SQLInjectionMiddleware) AddBlockedPath(path string) *SQLInjectionMiddleware {
	sim.blockedPath = append(sim.blockedPath, path)
	return sim
}

// GetName 获取中间件名称
func (sim *SQLInjectionMiddleware) GetName() string {
	return "sql_injection"
}

// Process 处理请求
func (sim *SQLInjectionMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !sim.enabled {
		next(w, r)
		return
	}

	// 检查路径
	if sim.isBlockedPath(r.URL.Path) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// 检查查询参数
	if sim.containsBlockedSQL(r.URL.RawQuery) {
		http.Error(w, "SQL injection attempt detected", http.StatusForbidden)
		return
	}

	// 检查表单数据
	if err := r.ParseForm(); err == nil {
		for _, values := range r.Form {
			for _, value := range values {
				if sim.containsBlockedSQL(value) {
					http.Error(w, "SQL injection attempt detected", http.StatusForbidden)
					return
				}
			}
		}
	}

	next(w, r)
}

// isBlockedPath 检查是否为被阻止的路径
func (sim *SQLInjectionMiddleware) isBlockedPath(path string) bool {
	for _, blockedPath := range sim.blockedPath {
		if strings.Contains(path, blockedPath) {
			return true
		}
	}
	return false
}

// containsBlockedSQL 检查是否包含被阻止的 SQL 关键字
func (sim *SQLInjectionMiddleware) containsBlockedSQL(input string) bool {
	input = strings.ToUpper(input)
	for _, keyword := range sim.blockedSQL {
		if strings.Contains(input, strings.ToUpper(keyword)) {
			return true
		}
	}
	return false
}

// SecurityHeadersMiddleware 安全头部中间件
type SecurityHeadersMiddleware struct {
	enabled bool
	headers map[string]string
}

// NewSecurityHeadersMiddleware 创建安全头部中间件
func NewSecurityHeadersMiddleware(enabled bool) *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{
		enabled: enabled,
		headers: map[string]string{
			"X-Frame-Options":           "DENY",
			"X-Content-Type-Options":    "nosniff",
			"X-XSS-Protection":          "1; mode=block",
			"Referrer-Policy":           "strict-origin-when-cross-origin",
			"Content-Security-Policy":   "default-src 'self'",
			"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
		},
	}
}

// AddHeader 添加安全头部
func (shm *SecurityHeadersMiddleware) AddHeader(key, value string) *SecurityHeadersMiddleware {
	shm.headers[key] = value
	return shm
}

// RemoveHeader 移除安全头部
func (shm *SecurityHeadersMiddleware) RemoveHeader(key string) *SecurityHeadersMiddleware {
	delete(shm.headers, key)
	return shm
}

// GetName 获取中间件名称
func (shm *SecurityHeadersMiddleware) GetName() string {
	return "security_headers"
}

// Process 处理请求
func (shm *SecurityHeadersMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !shm.enabled {
		next(w, r)
		return
	}

	// 设置安全头部
	for key, value := range shm.headers {
		w.Header().Set(key, value)
	}

	next(w, r)
}

// RateLimitMiddleware 速率限制中间件
type RateLimitMiddleware struct {
	enabled bool
	limit   int
	window  time.Duration
	store   RateLimitStore
	keyFunc RateLimitKeyFunc
}

// RateLimitStore 速率限制存储接口
type RateLimitStore interface {
	// Get 获取当前计数
	Get(key string) (int, error)
	// Increment 增加计数
	Increment(key string, window time.Duration) (int, error)
	// Reset 重置计数
	Reset(key string) error
}

// RateLimitKeyFunc 速率限制键生成函数
type RateLimitKeyFunc func(r *http.Request) string

// NewRateLimitMiddleware 创建速率限制中间件
func NewRateLimitMiddleware(enabled bool, limit int, window time.Duration) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		enabled: enabled,
		limit:   limit,
		window:  window,
		store:   NewMemoryRateLimitStore(),
		keyFunc: func(r *http.Request) string {
			return r.RemoteAddr
		},
	}
}

// SetStore 设置存储
func (rlm *RateLimitMiddleware) SetStore(store RateLimitStore) *RateLimitMiddleware {
	rlm.store = store
	return rlm
}

// SetKeyFunc 设置键生成函数
func (rlm *RateLimitMiddleware) SetKeyFunc(keyFunc RateLimitKeyFunc) *RateLimitMiddleware {
	rlm.keyFunc = keyFunc
	return rlm
}

// GetName 获取中间件名称
func (rlm *RateLimitMiddleware) GetName() string {
	return "rate_limit"
}

// Process 处理请求
func (rlm *RateLimitMiddleware) Process(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !rlm.enabled {
		next(w, r)
		return
	}

	key := rlm.keyFunc(r)
	count, err := rlm.store.Increment(key, rlm.window)
	if err != nil {
		http.Error(w, "Rate limit error", http.StatusInternalServerError)
		return
	}

	if count > rlm.limit {
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rlm.limit))
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(rlm.window).Unix()))
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rlm.limit))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", rlm.limit-count))
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(rlm.window).Unix()))

	next(w, r)
}

// MemoryRateLimitStore 内存速率限制存储
type MemoryRateLimitStore struct {
	store map[string]*rateLimitEntry
}

type rateLimitEntry struct {
	count     int
	lastReset time.Time
}

// NewMemoryRateLimitStore 创建内存速率限制存储
func NewMemoryRateLimitStore() *MemoryRateLimitStore {
	return &MemoryRateLimitStore{
		store: make(map[string]*rateLimitEntry),
	}
}

// Get 获取当前计数
func (mrs *MemoryRateLimitStore) Get(key string) (int, error) {
	if entry, exists := mrs.store[key]; exists {
		return entry.count, nil
	}
	return 0, nil
}

// Increment 增加计数
func (mrs *MemoryRateLimitStore) Increment(key string, window time.Duration) (int, error) {
	now := time.Now()
	entry, exists := mrs.store[key]

	if !exists || now.Sub(entry.lastReset) >= window {
		entry = &rateLimitEntry{
			count:     1,
			lastReset: now,
		}
		mrs.store[key] = entry
		return 1, nil
	}

	entry.count++
	return entry.count, nil
}

// Reset 重置计数
func (mrs *MemoryRateLimitStore) Reset(key string) error {
	delete(mrs.store, key)
	return nil
}
