package main

import (
	"fmt"
	"net/http"

	"laravel-go/framework/auth"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("Laravel-Go 授权系统演示")
	fmt.Println("==================================================")

	// 创建授权管理器
	authManager := auth.NewAuthorizationManager()

	// 注册预定义角色
	authManager.RegisterRole(auth.SuperAdminRole)
	authManager.RegisterRole(auth.AdminRole)
	authManager.RegisterRole(auth.UserRole)

	// 创建自定义策略
	userPolicy := &UserPolicy{
		BasePolicy: auth.NewPolicy("user-policy"),
	}
	authManager.RegisterPolicy("user-policy", userPolicy)

	// 创建用户提供者
	provider := auth.NewMemoryUserProvider()

	// 创建不同角色的用户
	superAdmin := &auth.BaseUser{
		ID:       1,
		Email:    "superadmin@example.com",
		Password: "password",
	}
	provider.AddUser(superAdmin)

	admin := &auth.BaseUser{
		ID:       2,
		Email:    "admin@example.com",
		Password: "password",
	}
	provider.AddUser(admin)

	user := &auth.BaseUser{
		ID:       3,
		Email:    "user@example.com",
		Password: "password",
	}
	provider.AddUser(user)

	// 创建认证守卫
	sessionStore := auth.NewMemorySessionStore()
	guard := auth.NewSessionGuard(provider, sessionStore)

	// 演示权限检查
	fmt.Println("\n1. 权限检查演示")
	fmt.Println("--------------------------------------------------")

	// 超级管理员权限检查
	guard.SetUser(superAdmin)
	fmt.Printf("超级管理员 (%s) 权限检查:\n", superAdmin.GetEmail())
	fmt.Printf("  查看用户: %v\n", authManager.CanView(superAdmin, nil))
	fmt.Printf("  创建用户: %v\n", authManager.CanCreate(superAdmin, nil))
	fmt.Printf("  删除用户: %v\n", authManager.CanDelete(superAdmin, nil))
	fmt.Printf("  查看角色: %v\n", authManager.CanView(superAdmin, "role"))

	// 管理员权限检查
	guard.SetUser(admin)
	fmt.Printf("\n管理员 (%s) 权限检查:\n", admin.GetEmail())
	fmt.Printf("  查看用户: %v\n", authManager.CanView(admin, nil))
	fmt.Printf("  创建用户: %v\n", authManager.CanCreate(admin, nil))
	fmt.Printf("  删除用户: %v\n", authManager.CanDelete(admin, nil))
	fmt.Printf("  查看角色: %v\n", authManager.CanView(admin, "role"))

	// 普通用户权限检查
	guard.SetUser(user)
	fmt.Printf("\n普通用户 (%s) 权限检查:\n", user.GetEmail())
	fmt.Printf("  查看用户: %v\n", authManager.CanView(user, nil))
	fmt.Printf("  创建用户: %v\n", authManager.CanCreate(user, nil))
	fmt.Printf("  删除用户: %v\n", authManager.CanDelete(user, nil))
	fmt.Printf("  查看角色: %v\n", authManager.CanView(user, "role"))

	// 演示角色管理
	fmt.Println("\n2. 角色管理演示")
	fmt.Println("--------------------------------------------------")

	// 获取角色信息
	superAdminRole, _ := authManager.GetRole("super-admin")
	fmt.Printf("超级管理员角色: %s (%s)\n", superAdminRole.GetName(), superAdminRole.GetDescription())
	fmt.Printf("  权限数量: %d\n", len(superAdminRole.GetPermissions()))

	adminRole, _ := authManager.GetRole("admin")
	fmt.Printf("管理员角色: %s (%s)\n", adminRole.GetName(), adminRole.GetDescription())
	fmt.Printf("  权限数量: %d\n", len(adminRole.GetPermissions()))

	userRole, _ := authManager.GetRole("user")
	fmt.Printf("普通用户角色: %s (%s)\n", userRole.GetName(), userRole.GetDescription())
	fmt.Printf("  权限数量: %d\n", len(userRole.GetPermissions()))

	// 演示权限管理
	fmt.Println("\n3. 权限管理演示")
	fmt.Println("--------------------------------------------------")

	// 创建自定义权限
	customPermission := auth.NewPermission("自定义权限", "custom.permission", "这是一个自定义权限", "custom", "execute")
	fmt.Printf("创建自定义权限: %s (%s)\n", customPermission.GetName(), customPermission.GetSlug())

	// 创建自定义角色
	customRole := auth.NewRole("自定义角色", "custom-role", "这是一个自定义角色")
	customRole.AddPermission(customPermission)
	customRole.AddPermission(auth.UserViewPermission)
	authManager.RegisterRole(customRole)

	fmt.Printf("创建自定义角色: %s (%s)\n", customRole.GetName(), customRole.GetSlug())
	fmt.Printf("  权限数量: %d\n", len(customRole.GetPermissions()))

	// 演示策略检查
	fmt.Println("\n4. 策略检查演示")
	fmt.Println("--------------------------------------------------")

	// 创建测试资源
	testResource := &TestResource{
		ID:      1,
		OwnerID: 3, // 属于用户3
		Name:    "测试资源",
	}

	// 不同用户对同一资源的权限
	guard.SetUser(superAdmin)
	fmt.Printf("超级管理员对资源 '%s' 的权限:\n", testResource.Name)
	fmt.Printf("  查看: %v\n", userPolicy.CanView(superAdmin, testResource))
	fmt.Printf("  更新: %v\n", userPolicy.CanUpdate(superAdmin, testResource))
	fmt.Printf("  删除: %v\n", userPolicy.CanDelete(superAdmin, testResource))

	guard.SetUser(admin)
	fmt.Printf("\n管理员对资源 '%s' 的权限:\n", testResource.Name)
	fmt.Printf("  查看: %v\n", userPolicy.CanView(admin, testResource))
	fmt.Printf("  更新: %v\n", userPolicy.CanUpdate(admin, testResource))
	fmt.Printf("  删除: %v\n", userPolicy.CanDelete(admin, testResource))

	guard.SetUser(user)
	fmt.Printf("\n资源所有者对资源 '%s' 的权限:\n", testResource.Name)
	fmt.Printf("  查看: %v\n", userPolicy.CanView(user, testResource))
	fmt.Printf("  更新: %v\n", userPolicy.CanUpdate(user, testResource))
	fmt.Printf("  删除: %v\n", userPolicy.CanDelete(user, testResource))

	// 演示中间件创建
	fmt.Println("\n5. 中间件创建演示")
	fmt.Println("--------------------------------------------------")

	// 创建各种中间件
	permissionMiddleware := auth.RequirePermission(guard, authManager, "user.view")
	fmt.Printf("创建权限中间件: 需要 user.view 权限\n")

	_ = auth.RequireRole(guard, authManager, "admin")
	fmt.Printf("创建角色中间件: 需要 admin 角色\n")

	_ = auth.RequireAllPermissions(guard, authManager, []string{"user.view", "user.create"})
	fmt.Printf("创建多权限中间件: 需要 user.view 和 user.create 权限\n")

	_ = auth.RequireAnyRole(guard, authManager, []string{"admin", "super-admin"})
	fmt.Printf("创建多角色中间件: 需要 admin 或 super-admin 角色\n")

	// 演示HTTP处理
	fmt.Println("\n6. HTTP处理演示")
	fmt.Println("--------------------------------------------------")

	// 创建简单的HTTP处理器
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "访问成功！用户: %s\n", guard.User().GetEmail())
	}

	// 使用权限中间件包装处理器
	_ = permissionMiddleware.Handle(handler)

	// 模拟HTTP请求（这里只是演示，不实际启动服务器）
	fmt.Println("创建受保护的HTTP处理器:")
	fmt.Println("  - 需要 user.view 权限")
	fmt.Println("  - 当前用户:", guard.User().GetEmail())
	fmt.Printf("  - 用户是否有权限: %v\n", authManager.CanView(guard.User(), nil))

	// 演示错误处理
	fmt.Println("\n7. 错误处理演示")
	fmt.Println("--------------------------------------------------")

	// 尝试获取不存在的角色
	_, err := authManager.GetRole("non-existent-role")
	if err != nil {
		fmt.Printf("获取不存在的角色错误: %v\n", err)
	}

	// 尝试获取不存在的策略
	_, err = authManager.GetPolicy("non-existent-policy")
	if err != nil {
		fmt.Printf("获取不存在的策略错误: %v\n", err)
	}

	fmt.Println("\n==================================================")
	fmt.Println("授权系统演示完成")
	fmt.Println("==================================================")
}

// UserPolicy 用户策略实现
type UserPolicy struct {
	*auth.BasePolicy
}

// Can 检查用户是否有权限执行操作
func (p *UserPolicy) Can(user auth.User, action string, resource interface{}) bool {
	// 超级管理员拥有所有权限
	if p.isSuperUser(user) {
		return true
	}

	// 根据操作类型和资源进行权限检查
	switch action {
	case "view":
		return p.CanView(user, resource)
	case "create":
		return p.CanCreate(user, resource)
	case "update":
		return p.CanUpdate(user, resource)
	case "delete":
		return p.CanDelete(user, resource)
	default:
		return false
	}
}

// CanView 检查用户是否可以查看资源
func (p *UserPolicy) CanView(user auth.User, resource interface{}) bool {
	// 所有认证用户都可以查看
	return user != nil
}

// CanCreate 检查用户是否可以创建资源
func (p *UserPolicy) CanCreate(user auth.User, resource interface{}) bool {
	// 管理员和超级管理员可以创建
	return p.isAdmin(user) || p.isSuperUser(user)
}

// CanUpdate 检查用户是否可以更新资源
func (p *UserPolicy) CanUpdate(user auth.User, resource interface{}) bool {
	if p.isSuperUser(user) || p.isAdmin(user) {
		return true
	}

	// 资源所有者可以更新自己的资源
	if testResource, ok := resource.(*TestResource); ok {
		return p.isResourceOwner(user, testResource)
	}

	return false
}

// CanDelete 检查用户是否可以删除资源
func (p *UserPolicy) CanDelete(user auth.User, resource interface{}) bool {
	// 只有超级管理员可以删除
	return p.isSuperUser(user)
}

// isSuperUser 检查是否为超级用户
func (p *UserPolicy) isSuperUser(user auth.User) bool {
	// 这里可以根据实际需求实现超级用户检查逻辑
	// 例如检查用户ID是否为1
	return user.GetID() == 1
}

// isAdmin 检查是否为管理员
func (p *UserPolicy) isAdmin(user auth.User) bool {
	// 这里可以根据实际需求实现管理员检查逻辑
	// 例如检查用户ID是否为2
	return user.GetID() == 2
}

// isResourceOwner 检查是否为资源所有者
func (p *UserPolicy) isResourceOwner(user auth.User, resource *TestResource) bool {
	return user.GetID() == resource.OwnerID
}

// TestResource 测试资源
type TestResource struct {
	ID      int
	OwnerID int
	Name    string
} 