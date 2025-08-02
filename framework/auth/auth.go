package auth

import (
	"errors"
	"time"
)

// User 用户接口
type User interface {
	GetID() interface{}
	GetEmail() string
	GetPassword() string
	GetRememberToken() string
	SetRememberToken(token string)
	GetAuthIdentifierName() string
	GetAuthIdentifier() interface{}
	GetAuthPassword() string
}

// Guard 认证守卫接口
type Guard interface {
	// 用户认证
	Authenticate(credentials map[string]interface{}) (User, error)
	// 检查用户是否已认证
	Check() bool
	// 获取当前认证用户
	User() User
	// 获取用户ID
	ID() interface{}
	// 登录用户
	Login(user User) error
	// 登录用户并记住
	LoginWithRemember(user User) error
	// 登出用户
	Logout() error
	// 验证凭据
	Validate(credentials map[string]interface{}) bool
	// 设置用户
	SetUser(user User)
	// 获取认证器
	GetProvider() UserProvider
}

// UserProvider 用户提供者接口
type UserProvider interface {
	// 通过ID检索用户
	RetrieveById(identifier interface{}) (User, error)
	// 通过凭据检索用户
	RetrieveByCredentials(credentials map[string]interface{}) (User, error)
	// 通过令牌检索用户
	RetrieveByToken(identifier interface{}, token string) (User, error)
	// 更新用户的记住令牌
	UpdateRememberToken(user User, token string) error
	// 验证凭据
	ValidateCredentials(user User, credentials map[string]interface{}) bool
}

// AuthManager 认证管理器
type AuthManager struct {
	guards    map[string]Guard
	providers map[string]UserProvider
	defaultGuard string
}

// NewAuthManager 创建认证管理器
func NewAuthManager() *AuthManager {
	return &AuthManager{
		guards:    make(map[string]Guard),
		providers: make(map[string]UserProvider),
		defaultGuard: "web",
	}
}

// Guard 获取指定的守卫
func (am *AuthManager) Guard(name string) Guard {
	if guard, exists := am.guards[name]; exists {
		return guard
	}
	return am.guards[am.defaultGuard]
}

// DefaultGuard 获取默认守卫
func (am *AuthManager) DefaultGuard() Guard {
	return am.Guard(am.defaultGuard)
}

// SetDefaultGuard 设置默认守卫
func (am *AuthManager) SetDefaultGuard(name string) {
	am.defaultGuard = name
}

// ExtendGuard 扩展守卫
func (am *AuthManager) ExtendGuard(name string, guard Guard) {
	am.guards[name] = guard
}

// ExtendProvider 扩展用户提供者
func (am *AuthManager) ExtendProvider(name string, provider UserProvider) {
	am.providers[name] = provider
}

// GetProvider 获取用户提供者
func (am *AuthManager) GetProvider(name string) UserProvider {
	return am.providers[name]
}

// 认证相关错误
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotAuthenticated = errors.New("user not authenticated")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)

// SessionGuard Session认证守卫
type SessionGuard struct {
	provider UserProvider
	user     User
	session  SessionStore
}

// NewSessionGuard 创建Session认证守卫
func NewSessionGuard(provider UserProvider, session SessionStore) *SessionGuard {
	return &SessionGuard{
		provider: provider,
		session:  session,
	}
}

// Authenticate 认证用户
func (sg *SessionGuard) Authenticate(credentials map[string]interface{}) (User, error) {
	user, err := sg.provider.RetrieveByCredentials(credentials)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !sg.provider.ValidateCredentials(user, credentials) {
		return nil, ErrInvalidCredentials
	}

	sg.user = user
	return user, nil
}

// Check 检查是否已认证
func (sg *SessionGuard) Check() bool {
	if sg.user != nil {
		return true
	}

	// 从session中获取用户ID
	userID := sg.session.Get("auth_user_id")
	if userID == nil {
		return false
	}

	user, err := sg.provider.RetrieveById(userID)
	if err != nil {
		return false
	}

	sg.user = user
	return true
}

// User 获取当前用户
func (sg *SessionGuard) User() User {
	if !sg.Check() {
		return nil
	}
	return sg.user
}

// ID 获取用户ID
func (sg *SessionGuard) ID() interface{} {
	if user := sg.User(); user != nil {
		return user.GetID()
	}
	return nil
}

// Login 登录用户
func (sg *SessionGuard) Login(user User) error {
	sg.user = user
	sg.session.Put("auth_user_id", user.GetID())
	return nil
}

// LoginWithRemember 登录并记住用户
func (sg *SessionGuard) LoginWithRemember(user User) error {
	// 生成记住令牌
	token := generateRememberToken()
	user.SetRememberToken(token)
	
	// 更新用户的记住令牌
	if err := sg.provider.UpdateRememberToken(user, token); err != nil {
		return err
	}

	return sg.Login(user)
}

// Logout 登出用户
func (sg *SessionGuard) Logout() error {
	sg.user = nil
	sg.session.Forget("auth_user_id")
	return nil
}

// Validate 验证凭据
func (sg *SessionGuard) Validate(credentials map[string]interface{}) bool {
	user, err := sg.provider.RetrieveByCredentials(credentials)
	if err != nil {
		return false
	}
	return sg.provider.ValidateCredentials(user, credentials)
}

// SetUser 设置用户
func (sg *SessionGuard) SetUser(user User) {
	sg.user = user
}

// GetProvider 获取用户提供者
func (sg *SessionGuard) GetProvider() UserProvider {
	return sg.provider
}

// SessionStore Session存储接口
type SessionStore interface {
	Get(key string) interface{}
	Put(key string, value interface{})
	Forget(key string)
	Has(key string) bool
}

// 生成记住令牌
func generateRememberToken() string {
	// 这里应该使用更安全的随机字符串生成
	return "remember_token_" + time.Now().Format("20060102150405")
} 