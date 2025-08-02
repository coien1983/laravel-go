package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
)

// MemorySessionStore 内存Session存储
type MemorySessionStore struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewMemorySessionStore 创建内存Session存储
func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		data: make(map[string]interface{}),
	}
}

// Get 获取Session值
func (mss *MemorySessionStore) Get(key string) interface{} {
	mss.mu.RLock()
	defer mss.mu.RUnlock()
	return mss.data[key]
}

// Put 设置Session值
func (mss *MemorySessionStore) Put(key string, value interface{}) {
	mss.mu.Lock()
	defer mss.mu.Unlock()
	mss.data[key] = value
}

// Forget 删除Session值
func (mss *MemorySessionStore) Forget(key string) {
	mss.mu.Lock()
	defer mss.mu.Unlock()
	delete(mss.data, key)
}

// Has 检查Session值是否存在
func (mss *MemorySessionStore) Has(key string) bool {
	mss.mu.RLock()
	defer mss.mu.RUnlock()
	_, exists := mss.data[key]
	return exists
}

// CookieSessionStore Cookie Session存储
type CookieSessionStore struct {
	request  *http.Request
	response http.ResponseWriter
	secret   string
}

// NewCookieSessionStore 创建Cookie Session存储
func NewCookieSessionStore(request *http.Request, response http.ResponseWriter, secret string) *CookieSessionStore {
	return &CookieSessionStore{
		request:  request,
		response: response,
		secret:   secret,
	}
}

// Get 获取Cookie值
func (css *CookieSessionStore) Get(key string) interface{} {
	cookie, err := css.request.Cookie(key)
	if err != nil {
		return nil
	}
	return cookie.Value
}

// Put 设置Cookie值
func (css *CookieSessionStore) Put(key string, value interface{}) {
	valueStr, ok := value.(string)
	if !ok {
		valueStr = "value"
	}

	cookie := &http.Cookie{
		Name:     key,
		Value:    valueStr,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // 在生产环境中应该设置为true
		MaxAge:   3600,  // 1小时
	}

	http.SetCookie(css.response, cookie)
}

// Forget 删除Cookie值
func (css *CookieSessionStore) Forget(key string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1, // 立即过期
	}

	http.SetCookie(css.response, cookie)
}

// Has 检查Cookie值是否存在
func (css *CookieSessionStore) Has(key string) bool {
	_, err := css.request.Cookie(key)
	return err == nil
}

// SessionManager Session管理器
type SessionManager struct {
	driver  string
	config  map[string]interface{}
	sessions map[string]SessionStore
	mu      sync.RWMutex
}

// NewSessionManager 创建Session管理器
func NewSessionManager(driver string, config map[string]interface{}) *SessionManager {
	return &SessionManager{
		driver:   driver,
		config:   config,
		sessions: make(map[string]SessionStore),
	}
}

// Store 获取Session存储
func (sm *SessionManager) Store(name string) SessionStore {
	sm.mu.RLock()
	if store, exists := sm.sessions[name]; exists {
		sm.mu.RUnlock()
		return store
	}
	sm.mu.RUnlock()

	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 双重检查
	if store, exists := sm.sessions[name]; exists {
		return store
	}

	// 创建新的Session存储
	var store SessionStore
	switch sm.driver {
	case "memory":
		store = NewMemorySessionStore()
	case "cookie":
		// Cookie存储需要request和response，这里暂时返回内存存储
		store = NewMemorySessionStore()
	default:
		store = NewMemorySessionStore()
	}

	sm.sessions[name] = store
	return store
}

// 生成Session ID
func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// 生成CSRF令牌
func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
} 