package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthManager(t *testing.T) {
	// 测试创建认证管理器
	manager := NewAuthManager()
	if manager == nil {
		t.Error("Expected non-nil AuthManager")
	}

	// 测试默认守卫
	defaultGuard := manager.DefaultGuard()
	if defaultGuard != nil {
		t.Error("Expected nil default guard before setup")
	}

	// 测试设置默认守卫
	manager.SetDefaultGuard("web")
	if manager.defaultGuard != "web" {
		t.Error("Expected default guard to be 'web'")
	}
}

func TestSessionGuard(t *testing.T) {
	// 创建内存用户提供者
	provider := NewMemoryUserProvider()
	
	// 创建测试用户
	user := &BaseUser{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}
	provider.AddUser(user)

	// 创建Session存储
	session := NewMemorySessionStore()

	// 创建Session守卫
	guard := NewSessionGuard(provider, session)

	// 测试未认证状态
	if guard.Check() {
		t.Error("Expected user to be unauthenticated initially")
	}

	// 测试认证
	credentials := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password",
	}

	authenticatedUser, err := guard.Authenticate(credentials)
	if err != nil {
		t.Errorf("Expected no error during authentication, got: %v", err)
	}

	if authenticatedUser == nil {
		t.Error("Expected authenticated user")
	}

	// 测试登录
	err = guard.Login(authenticatedUser)
	if err != nil {
		t.Errorf("Expected no error during login, got: %v", err)
	}

	// 测试认证状态
	if !guard.Check() {
		t.Error("Expected user to be authenticated after login")
	}

	// 测试获取用户
	currentUser := guard.User()
	if currentUser == nil {
		t.Error("Expected current user")
	}

	if currentUser.GetID() != 1 {
		t.Errorf("Expected user ID 1, got: %v", currentUser.GetID())
	}

	// 测试登出
	err = guard.Logout()
	if err != nil {
		t.Errorf("Expected no error during logout, got: %v", err)
	}

	if guard.Check() {
		t.Error("Expected user to be unauthenticated after logout")
	}
}

func TestJWTGuard(t *testing.T) {
	// 创建内存用户提供者
	provider := NewMemoryUserProvider()
	
	// 创建测试用户
	user := &BaseUser{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}
	provider.AddUser(user)

	// 创建JWT守卫
	secret := "test-secret"
	ttl := 1 * time.Hour
	guard := NewJWTGuard(provider, secret, ttl)

	// 测试认证
	credentials := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password",
	}

	authenticatedUser, err := guard.Authenticate(credentials)
	if err != nil {
		t.Errorf("Expected no error during authentication, got: %v", err)
	}

	// 测试生成令牌
	token, err := guard.GenerateToken(authenticatedUser)
	if err != nil {
		t.Errorf("Expected no error generating token, got: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	// 测试验证令牌
	claims, err := guard.ValidateToken(token)
	if err != nil {
		t.Errorf("Expected no error validating token, got: %v", err)
	}

	// 注意：interface{}比较需要特殊处理
	if claims.UserID == nil {
		t.Error("Expected non-nil user ID in claims")
	} else if fmt.Sprintf("%v", claims.UserID) != "1" {
		t.Errorf("Expected user ID 1, got: %v", claims.UserID)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got: %s", claims.Email)
	}
}

func TestMemoryUserProvider(t *testing.T) {
	provider := NewMemoryUserProvider()

	// 创建测试用户
	user := &BaseUser{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}

	// 测试添加用户
	provider.AddUser(user)

	// 测试通过ID检索用户
	retrievedUser, err := provider.RetrieveById(1)
	if err != nil {
		t.Errorf("Expected no error retrieving user by ID, got: %v", err)
	}

	if retrievedUser.GetID() != 1 {
		t.Errorf("Expected user ID 1, got: %v", retrievedUser.GetID())
	}

	// 测试通过凭据检索用户
	credentials := map[string]interface{}{
		"email": "test@example.com",
	}

	retrievedUser, err = provider.RetrieveByCredentials(credentials)
	if err != nil {
		t.Errorf("Expected no error retrieving user by credentials, got: %v", err)
	}

	if retrievedUser.GetEmail() != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got: %s", retrievedUser.GetEmail())
	}

	// 测试验证凭据
	validCredentials := map[string]interface{}{
		"password": "password",
	}

	if !provider.ValidateCredentials(user, validCredentials) {
		t.Error("Expected credentials to be valid")
	}

	invalidCredentials := map[string]interface{}{
		"password": "wrong_password",
	}

	if provider.ValidateCredentials(user, invalidCredentials) {
		t.Error("Expected credentials to be invalid")
	}
}

func TestAuthMiddleware(t *testing.T) {
	// 创建内存用户提供者
	provider := NewMemoryUserProvider()
	
	// 创建测试用户
	user := &BaseUser{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}
	provider.AddUser(user)

	// 创建Session存储
	session := NewMemorySessionStore()

	// 创建Session守卫
	guard := NewSessionGuard(provider, session)

	// 创建认证中间件
	middleware := NewAuthMiddleware(guard)

	// 测试未认证状态
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	middleware.Handle(handler)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got: %d", w.Code)
	}

	if handlerCalled {
		t.Error("Expected handler not to be called")
	}

	// 测试已认证状态
	guard.Login(user)

	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	handlerCalled = false
	middleware.Handle(handler)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", w.Code)
	}

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
}

func TestJWTMiddleware(t *testing.T) {
	// 创建内存用户提供者
	provider := NewMemoryUserProvider()
	
	// 创建JWT守卫
	secret := "test-secret"
	ttl := 1 * time.Hour
	guard := NewJWTGuard(provider, secret, ttl)

	// 创建JWT中间件
	middleware := NewJWTMiddleware(guard)

	// 测试无效令牌
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	middleware.Handle(handler)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got: %d", w.Code)
	}

	if handlerCalled {
		t.Error("Expected handler not to be called")
	}

	// 测试缺少Authorization头
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	handlerCalled = false
	middleware.Handle(handler)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got: %d", w.Code)
	}

	if handlerCalled {
		t.Error("Expected handler not to be called")
	}
}

func TestGuestMiddleware(t *testing.T) {
	// 创建内存用户提供者
	provider := NewMemoryUserProvider()
	
	// 创建Session存储
	session := NewMemorySessionStore()

	// 创建Session守卫
	guard := NewSessionGuard(provider, session)

	// 创建访客中间件
	middleware := NewGuestMiddleware(guard)

	// 测试未认证状态（应该通过）
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	middleware.Handle(handler)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", w.Code)
	}

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}

	// 测试已认证状态（应该被拒绝）
	user := &BaseUser{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}
	guard.Login(user)

	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	handlerCalled = false
	middleware.Handle(handler)(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got: %d", w.Code)
	}

	if handlerCalled {
		t.Error("Expected handler not to be called")
	}
}

func TestBaseUser(t *testing.T) {
	user := &BaseUser{
		ID:            1,
		Email:         "test@example.com",
		Password:      "password",
		RememberToken: "remember_token",
	}

	// 测试基本属性
	if user.GetID() != 1 {
		t.Errorf("Expected ID 1, got: %v", user.GetID())
	}

	if user.GetEmail() != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got: %s", user.GetEmail())
	}

	if user.GetPassword() != "password" {
		t.Errorf("Expected password 'password', got: %s", user.GetPassword())
	}

	if user.GetRememberToken() != "remember_token" {
		t.Errorf("Expected remember token 'remember_token', got: %s", user.GetRememberToken())
	}

	// 测试设置记住令牌
	user.SetRememberToken("new_token")
	if user.GetRememberToken() != "new_token" {
		t.Errorf("Expected remember token 'new_token', got: %s", user.GetRememberToken())
	}

	// 测试认证相关方法
	if user.GetAuthIdentifierName() != "id" {
		t.Errorf("Expected auth identifier name 'id', got: %s", user.GetAuthIdentifierName())
	}

	if user.GetAuthIdentifier() != 1 {
		t.Errorf("Expected auth identifier 1, got: %v", user.GetAuthIdentifier())
	}

	if user.GetAuthPassword() != "password" {
		t.Errorf("Expected auth password 'password', got: %s", user.GetAuthPassword())
	}
} 