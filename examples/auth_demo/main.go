package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"laravel-go/framework/auth"
)

func main() {
	fmt.Println("🚀 Laravel-Go 认证系统演示")
	fmt.Println("==================================================")

	// 演示Session认证
	demoSessionAuth()

	fmt.Println()

	// 演示JWT认证
	demoJWTAuth()

	fmt.Println()

	// 演示认证中间件
	demoAuthMiddleware()

	fmt.Println()

	// 演示用户提供者
	demoUserProvider()

	fmt.Println("✅ 认证系统演示完成!")
}

// 演示Session认证
func demoSessionAuth() {
	fmt.Println("📝 Session认证演示:")

	// 创建内存用户提供者
	provider := auth.NewMemoryUserProvider()

	// 创建测试用户
	user := &auth.BaseUser{
		ID:       1,
		Email:    "john@example.com",
		Password: "password",
	}
	provider.AddUser(user)

	// 创建Session存储
	session := auth.NewMemorySessionStore()

	// 创建Session守卫
	guard := auth.NewSessionGuard(provider, session)

	// 测试未认证状态
	fmt.Println("  检查未认证状态...")
	if guard.Check() {
		fmt.Println("   ❌ 用户已认证 (意外)")
	} else {
		fmt.Println("   ✅ 用户未认证 (正确)")
	}

	// 测试认证
	fmt.Println("  尝试认证用户...")
	credentials := map[string]interface{}{
		"email":    "john@example.com",
		"password": "password",
	}

	authenticatedUser, err := guard.Authenticate(credentials)
	if err != nil {
		fmt.Printf("   ❌ 认证失败: %v\n", err)
		return
	}
	fmt.Printf("   ✅ 认证成功: %s (ID: %v)\n", authenticatedUser.GetEmail(), authenticatedUser.GetID())

	// 测试登录
	fmt.Println("  登录用户...")
	err = guard.Login(authenticatedUser)
	if err != nil {
		fmt.Printf("   ❌ 登录失败: %v\n", err)
		return
	}
	fmt.Println("   ✅ 登录成功")

	// 测试认证状态
	fmt.Println("  检查认证状态...")
	if guard.Check() {
		fmt.Println("   ✅ 用户已认证")
	} else {
		fmt.Println("   ❌ 用户未认证 (意外)")
	}

	// 测试获取当前用户
	fmt.Println("  获取当前用户...")
	currentUser := guard.User()
	if currentUser != nil {
		fmt.Printf("   ✅ 当前用户: %s (ID: %v)\n", currentUser.GetEmail(), currentUser.GetID())
	} else {
		fmt.Println("   ❌ 无法获取当前用户")
	}

	// 测试登出
	fmt.Println("  登出用户...")
	err = guard.Logout()
	if err != nil {
		fmt.Printf("   ❌ 登出失败: %v\n", err)
		return
	}
	fmt.Println("   ✅ 登出成功")

	// 验证登出状态
	if guard.Check() {
		fmt.Println("   ❌ 用户仍已认证 (意外)")
	} else {
		fmt.Println("   ✅ 用户已登出")
	}
}

// 演示JWT认证
func demoJWTAuth() {
	fmt.Println("🔐 JWT认证演示:")

	// 创建内存用户提供者
	provider := auth.NewMemoryUserProvider()

	// 创建测试用户
	user := &auth.BaseUser{
		ID:       2,
		Email:    "jane@example.com",
		Password: "password",
	}
	provider.AddUser(user)

	// 创建JWT守卫
	secret := "your-secret-key"
	ttl := 1 * time.Hour
	guard := auth.NewJWTGuard(provider, secret, ttl)

	// 测试认证
	fmt.Println("  认证用户...")
	credentials := map[string]interface{}{
		"email":    "jane@example.com",
		"password": "password",
	}

	authenticatedUser, err := guard.Authenticate(credentials)
	if err != nil {
		fmt.Printf("   ❌ 认证失败: %v\n", err)
		return
	}
	fmt.Printf("   ✅ 认证成功: %s\n", authenticatedUser.GetEmail())

	// 生成JWT令牌
	fmt.Println("  生成JWT令牌...")
	token, err := guard.GenerateToken(authenticatedUser)
	if err != nil {
		fmt.Printf("   ❌ 生成令牌失败: %v\n", err)
		return
	}
	fmt.Printf("   ✅ 令牌生成成功: %s...\n", token[:20])

	// 验证JWT令牌
	fmt.Println("  验证JWT令牌...")
	claims, err := guard.ValidateToken(token)
	if err != nil {
		fmt.Printf("   ❌ 令牌验证失败: %v\n", err)
		return
	}
	fmt.Printf("   ✅ 令牌验证成功: 用户ID=%v, 邮箱=%s\n", claims.UserID, claims.Email)

	// 生成刷新令牌
	fmt.Println("  生成刷新令牌...")
	refreshToken, err := guard.GenerateRefreshToken(authenticatedUser)
	if err != nil {
		fmt.Printf("   ❌ 生成刷新令牌失败: %v\n", err)
		return
	}
	fmt.Printf("   ✅ 刷新令牌生成成功: %s...\n", refreshToken[:20])

	// 测试无效令牌
	fmt.Println("  测试无效令牌...")
	_, err = guard.ValidateToken("invalid-token")
	if err != nil {
		fmt.Printf("   ✅ 无效令牌被正确拒绝: %v\n", err)
	} else {
		fmt.Println("   ❌ 无效令牌被错误接受")
	}
}

// 演示认证中间件
func demoAuthMiddleware() {
	fmt.Println("🛡️ 认证中间件演示:")

	// 创建内存用户提供者
	provider := auth.NewMemoryUserProvider()

	// 创建测试用户
	user := &auth.BaseUser{
		ID:       3,
		Email:    "admin@example.com",
		Password: "password",
	}
	provider.AddUser(user)

	// 创建Session存储
	session := auth.NewMemorySessionStore()

	// 创建Session守卫
	guard := auth.NewSessionGuard(provider, session)

	// 创建认证中间件
	authMiddleware := auth.NewAuthMiddleware(guard)

	// 创建访客中间件
	guestMiddleware := auth.NewGuestMiddleware(guard)

	// 测试受保护的路由（未认证）
	fmt.Println("  测试受保护路由 (未认证)...")
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Protected content"))
	})

	authMiddleware.Handle(protectedHandler)(w, req)

	if w.Code == http.StatusUnauthorized {
		fmt.Println("   ✅ 未认证用户被正确拒绝")
	} else {
		fmt.Printf("   ❌ 未认证用户被错误接受 (状态码: %d)\n", w.Code)
	}

	// 测试受保护的路由（已认证）
	fmt.Println("  测试受保护路由 (已认证)...")
	guard.Login(user)

	req = httptest.NewRequest("GET", "/protected", nil)
	w = httptest.NewRecorder()

	authMiddleware.Handle(protectedHandler)(w, req)

	if w.Code == http.StatusOK {
		fmt.Println("   ✅ 已认证用户可以访问")
	} else {
		fmt.Printf("   ❌ 已认证用户被错误拒绝 (状态码: %d)\n", w.Code)
	}

	// 测试访客路由（未认证）
	fmt.Println("  测试访客路由 (未认证)...")
	guard.Logout()

	req = httptest.NewRequest("GET", "/guest", nil)
	w = httptest.NewRecorder()

	guestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Guest content"))
	})

	guestMiddleware.Handle(guestHandler)(w, req)

	if w.Code == http.StatusOK {
		fmt.Println("   ✅ 访客可以访问")
	} else {
		fmt.Printf("   ❌ 访客被错误拒绝 (状态码: %d)\n", w.Code)
	}

	// 测试访客路由（已认证）
	fmt.Println("  测试访客路由 (已认证)...")
	guard.Login(user)

	req = httptest.NewRequest("GET", "/guest", nil)
	w = httptest.NewRecorder()

	guestMiddleware.Handle(guestHandler)(w, req)

	if w.Code == http.StatusForbidden {
		fmt.Println("   ✅ 已认证用户被正确拒绝访问访客页面")
	} else {
		fmt.Printf("   ❌ 已认证用户被错误允许访问访客页面 (状态码: %d)\n", w.Code)
	}
}

// 演示用户提供者
func demoUserProvider() {
	fmt.Println("👥 用户提供者演示:")

	// 创建内存用户提供者
	provider := auth.NewMemoryUserProvider()

	// 添加多个用户
	users := []*auth.BaseUser{
		{ID: 1, Email: "user1@example.com", Password: "password1"},
		{ID: 2, Email: "user2@example.com", Password: "password2"},
		{ID: 3, Email: "user3@example.com", Password: "password3"},
	}

	for _, user := range users {
		provider.AddUser(user)
	}
	fmt.Printf("   ✅ 添加了 %d 个用户\n", len(users))

	// 测试通过ID检索用户
	fmt.Println("  通过ID检索用户...")
	for _, expectedUser := range users {
		user, err := provider.RetrieveById(expectedUser.GetID())
		if err != nil {
			fmt.Printf("   ❌ 无法检索用户 ID %v: %v\n", expectedUser.GetID(), err)
			continue
		}
		fmt.Printf("   ✅ 检索到用户: %s (ID: %v)\n", user.GetEmail(), user.GetID())
	}

	// 测试通过凭据检索用户
	fmt.Println("  通过凭据检索用户...")
	for _, expectedUser := range users {
		credentials := map[string]interface{}{
			"email": expectedUser.GetEmail(),
		}
		user, err := provider.RetrieveByCredentials(credentials)
		if err != nil {
			fmt.Printf("   ❌ 无法通过凭据检索用户 %s: %v\n", expectedUser.GetEmail(), err)
			continue
		}
		fmt.Printf("   ✅ 通过凭据检索到用户: %s\n", user.GetEmail())
	}

	// 测试凭据验证
	fmt.Println("  验证用户凭据...")
	testUser := users[0]
	validCredentials := map[string]interface{}{
		"password": "password1",
	}
	invalidCredentials := map[string]interface{}{
		"password": "wrong_password",
	}

	if provider.ValidateCredentials(testUser, validCredentials) {
		fmt.Println("   ✅ 有效凭据验证通过")
	} else {
		fmt.Println("   ❌ 有效凭据验证失败")
	}

	if !provider.ValidateCredentials(testUser, invalidCredentials) {
		fmt.Println("   ✅ 无效凭据被正确拒绝")
	} else {
		fmt.Println("   ❌ 无效凭据被错误接受")
	}
} 