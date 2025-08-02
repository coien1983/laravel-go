package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTGuard JWT认证守卫
type JWTGuard struct {
	provider UserProvider
	user     User
	secret   string
	ttl      time.Duration
	refreshTTL time.Duration
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID interface{} `json:"user_id"`
	Email  string      `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTGuard 创建JWT认证守卫
func NewJWTGuard(provider UserProvider, secret string, ttl time.Duration) *JWTGuard {
	return &JWTGuard{
		provider:   provider,
		secret:     secret,
		ttl:        ttl,
		refreshTTL: ttl * 24, // 刷新令牌有效期更长
	}
}

// Authenticate 认证用户
func (jg *JWTGuard) Authenticate(credentials map[string]interface{}) (User, error) {
	user, err := jg.provider.RetrieveByCredentials(credentials)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !jg.provider.ValidateCredentials(user, credentials) {
		return nil, ErrInvalidCredentials
	}

	jg.user = user
	return user, nil
}

// Check 检查是否已认证
func (jg *JWTGuard) Check() bool {
	return jg.user != nil
}

// User 获取当前用户
func (jg *JWTGuard) User() User {
	return jg.user
}

// ID 获取用户ID
func (jg *JWTGuard) ID() interface{} {
	if user := jg.User(); user != nil {
		return user.GetID()
	}
	return nil
}

// Login 登录用户
func (jg *JWTGuard) Login(user User) error {
	jg.user = user
	return nil
}

// LoginWithRemember 登录并记住用户（JWT中不适用，但保持接口一致）
func (jg *JWTGuard) LoginWithRemember(user User) error {
	return jg.Login(user)
}

// Logout 登出用户
func (jg *JWTGuard) Logout() error {
	jg.user = nil
	return nil
}

// Validate 验证凭据
func (jg *JWTGuard) Validate(credentials map[string]interface{}) bool {
	user, err := jg.provider.RetrieveByCredentials(credentials)
	if err != nil {
		return false
	}
	return jg.provider.ValidateCredentials(user, credentials)
}

// SetUser 设置用户
func (jg *JWTGuard) SetUser(user User) {
	jg.user = user
}

// GetProvider 获取用户提供者
func (jg *JWTGuard) GetProvider() UserProvider {
	return jg.provider
}

// GenerateToken 生成JWT令牌
func (jg *JWTGuard) GenerateToken(user User) (string, error) {
	claims := JWTClaims{
		UserID: user.GetID(),
		Email:  user.GetEmail(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jg.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "laravel-go",
			Subject:   fmt.Sprintf("%v", user.GetID()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jg.secret))
}

// GenerateRefreshToken 生成刷新令牌
func (jg *JWTGuard) GenerateRefreshToken(user User) (string, error) {
	claims := JWTClaims{
		UserID: user.GetID(),
		Email:  user.GetEmail(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jg.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "laravel-go",
			Subject:   fmt.Sprintf("%v", user.GetID()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jg.secret))
}

// ValidateToken 验证JWT令牌
func (jg *JWTGuard) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jg.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GetUserFromToken 从令牌获取用户
func (jg *JWTGuard) GetUserFromToken(tokenString string) (User, error) {
	claims, err := jg.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := jg.provider.RetrieveById(claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// RefreshToken 刷新令牌
func (jg *JWTGuard) RefreshToken(refreshToken string) (string, error) {
	claims, err := jg.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	user, err := jg.provider.RetrieveById(claims.UserID)
	if err != nil {
		return "", ErrUserNotFound
	}

	return jg.GenerateToken(user)
}

// 生成随机密钥
func generateSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
} 