package auth

import (
	"fmt"
	"testing"
)

func TestPermission(t *testing.T) {
	permission := NewPermission("测试权限", "test.permission", "这是一个测试权限", "test", "read")

	if permission.GetName() != "测试权限" {
		t.Errorf("Expected permission name to be '测试权限', got %s", permission.GetName())
	}

	if permission.GetSlug() != "test.permission" {
		t.Errorf("Expected permission slug to be 'test.permission', got %s", permission.GetSlug())
	}

	if permission.GetDescription() != "这是一个测试权限" {
		t.Errorf("Expected permission description to be '这是一个测试权限', got %s", permission.GetDescription())
	}

	if permission.GetResource() != "test" {
		t.Errorf("Expected permission resource to be 'test', got %s", permission.GetResource())
	}

	if permission.GetAction() != "read" {
		t.Errorf("Expected permission action to be 'read', got %s", permission.GetAction())
	}
}

func TestRole(t *testing.T) {
	role := NewRole("测试角色", "test-role", "这是一个测试角色")

	if role.GetName() != "测试角色" {
		t.Errorf("Expected role name to be '测试角色', got %s", role.GetName())
	}

	if role.GetSlug() != "test-role" {
		t.Errorf("Expected role slug to be 'test-role', got %s", role.GetSlug())
	}

	if role.GetDescription() != "这是一个测试角色" {
		t.Errorf("Expected role description to be '这是一个测试角色', got %s", role.GetDescription())
	}

	// 测试权限管理
	permission1 := NewPermission("权限1", "permission1", "第一个权限", "test", "read")
	permission2 := NewPermission("权限2", "permission2", "第二个权限", "test", "write")

	role.AddPermission(permission1)
	role.AddPermission(permission2)

	permissions := role.GetPermissions()
	if len(permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(permissions))
	}

	if !role.HasPermission(permission1) {
		t.Error("Expected role to have permission1")
	}

	if !role.HasPermissionByName("权限1") {
		t.Error("Expected role to have permission by name '权限1'")
	}

	if role.HasPermissionByName("不存在的权限") {
		t.Error("Expected role to not have non-existent permission")
	}

	// 测试移除权限
	role.RemovePermission(permission1)
	permissions = role.GetPermissions()
	if len(permissions) != 1 {
		t.Errorf("Expected 1 permission after removal, got %d", len(permissions))
	}

	if role.HasPermission(permission1) {
		t.Error("Expected role to not have permission1 after removal")
	}
}

func TestAuthorizationManager(t *testing.T) {
	authManager := NewAuthorizationManager()

	// 测试角色注册和获取
	role := NewRole("测试角色", "test-role", "测试角色描述")
	authManager.RegisterRole(role)

	retrievedRole, err := authManager.GetRole("test-role")
	if err != nil {
		t.Errorf("Expected no error when getting registered role, got %v", err)
	}

	if retrievedRole.GetName() != "测试角色" {
		t.Errorf("Expected role name to be '测试角色', got %s", retrievedRole.GetName())
	}

	// 测试获取不存在的角色
	_, err = authManager.GetRole("non-existent-role")
	if err == nil {
		t.Error("Expected error when getting non-existent role")
	}

	// 测试策略注册和获取
	policy := NewPolicy("test-policy")
	authManager.RegisterPolicy("test-policy", policy)

	retrievedPolicy, err := authManager.GetPolicy("test-policy")
	if err != nil {
		t.Errorf("Expected no error when getting registered policy, got %v", err)
	}

	if retrievedPolicy == nil {
		t.Error("Expected policy to be returned")
	}

	// 测试获取不存在的策略
	_, err = authManager.GetPolicy("non-existent-policy")
	if err == nil {
		t.Error("Expected error when getting non-existent policy")
	}
}

func TestBasePolicy(t *testing.T) {
	policy := NewPolicy("test-policy")

	// 创建测试用户
	user := &BaseUser{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}

	// 基础策略默认拒绝所有操作
	if policy.Can(user, "read", nil) {
		t.Error("Expected base policy to deny all operations")
	}

	if policy.CanView(user, nil) {
		t.Error("Expected base policy to deny view operation")
	}

	if policy.CanCreate(user, nil) {
		t.Error("Expected base policy to deny create operation")
	}

	if policy.CanUpdate(user, nil) {
		t.Error("Expected base policy to deny update operation")
	}

	if policy.CanDelete(user, nil) {
		t.Error("Expected base policy to deny delete operation")
	}
}

func TestResourcePolicy(t *testing.T) {
	policy := NewResourcePolicy("test-resource-policy", "test-resource")

	// 创建测试用户
	user := &BaseUser{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}

	// 资源策略默认拒绝所有操作
	if policy.Can(user, "read", nil) {
		t.Error("Expected resource policy to deny all operations by default")
	}
}

func TestAuthorizationManagerCanMethods(t *testing.T) {
	authManager := NewAuthorizationManager()

	// 创建测试用户
	user := &BaseUser{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}

	// 测试各种权限检查方法
	if authManager.Can(user, "read", nil) {
		t.Error("Expected authorization manager to deny operation by default")
	}

	if authManager.CanView(user, nil) {
		t.Error("Expected authorization manager to deny view operation by default")
	}

	if authManager.CanCreate(user, nil) {
		t.Error("Expected authorization manager to deny create operation by default")
	}

	if authManager.CanUpdate(user, nil) {
		t.Error("Expected authorization manager to deny update operation by default")
	}

	if authManager.CanDelete(user, nil) {
		t.Error("Expected authorization manager to deny delete operation by default")
	}
}

func TestPredefinedPermissions(t *testing.T) {
	// 测试预定义权限
	if UserViewPermission.GetName() != "查看用户" {
		t.Errorf("Expected UserViewPermission name to be '查看用户', got %s", UserViewPermission.GetName())
	}

	if UserCreatePermission.GetSlug() != "user.create" {
		t.Errorf("Expected UserCreatePermission slug to be 'user.create', got %s", UserCreatePermission.GetSlug())
	}

	if RoleViewPermission.GetResource() != "role" {
		t.Errorf("Expected RoleViewPermission resource to be 'role', got %s", RoleViewPermission.GetResource())
	}

	if PermissionViewPermission.GetAction() != "view" {
		t.Errorf("Expected PermissionViewPermission action to be 'view', got %s", PermissionViewPermission.GetAction())
	}
}

func TestPredefinedRoles(t *testing.T) {
	// 测试预定义角色
	if SuperAdminRole.GetName() != "超级管理员" {
		t.Errorf("Expected SuperAdminRole name to be '超级管理员', got %s", SuperAdminRole.GetName())
	}

	if AdminRole.GetSlug() != "admin" {
		t.Errorf("Expected AdminRole slug to be 'admin', got %s", AdminRole.GetSlug())
	}

	if UserRole.GetDescription() != "普通用户" {
		t.Errorf("Expected UserRole description to be '普通用户', got %s", UserRole.GetDescription())
	}

	// 测试超级管理员角色的权限
	permissions := SuperAdminRole.GetPermissions()
	if len(permissions) == 0 {
		t.Error("Expected SuperAdminRole to have permissions")
	}

	// 检查超级管理员是否有用户查看权限
	if !SuperAdminRole.HasPermission(UserViewPermission) {
		t.Error("Expected SuperAdminRole to have user view permission")
	}

	// 检查管理员角色的权限
	if !AdminRole.HasPermission(UserViewPermission) {
		t.Error("Expected AdminRole to have user view permission")
	}

	if AdminRole.HasPermission(UserDeletePermission) {
		t.Error("Expected AdminRole to not have user delete permission")
	}

	// 检查普通用户角色的权限
	if !UserRole.HasPermission(UserViewPermission) {
		t.Error("Expected UserRole to have user view permission")
	}

	if UserRole.HasPermission(UserCreatePermission) {
		t.Error("Expected UserRole to not have user create permission")
	}
}

func TestErrorDefinitions(t *testing.T) {
	// 测试错误定义
	if ErrPermissionDenied.Error() != "permission denied" {
		t.Errorf("Expected ErrPermissionDenied message to be 'permission denied', got %s", ErrPermissionDenied.Error())
	}

	if ErrRoleNotFound.Error() != "role not found" {
		t.Errorf("Expected ErrRoleNotFound message to be 'role not found', got %s", ErrRoleNotFound.Error())
	}

	if ErrPolicyNotFound.Error() != "policy not found" {
		t.Errorf("Expected ErrPolicyNotFound message to be 'policy not found', got %s", ErrPolicyNotFound.Error())
	}
}

func TestRolePermissionManagement(t *testing.T) {
	role := NewRole("测试角色", "test-role", "测试角色描述")

	// 测试添加重复权限
	permission := NewPermission("测试权限", "test.permission", "测试权限描述", "test", "read")
	role.AddPermission(permission)
	role.AddPermission(permission) // 添加重复权限

	permissions := role.GetPermissions()
	if len(permissions) != 1 {
		t.Errorf("Expected 1 permission after adding duplicate, got %d", len(permissions))
	}

	// 测试移除不存在的权限
	nonExistentPermission := NewPermission("不存在", "non.existent", "不存在的权限", "test", "write")
	role.RemovePermission(nonExistentPermission)

	permissions = role.GetPermissions()
	if len(permissions) != 1 {
		t.Errorf("Expected 1 permission after removing non-existent permission, got %d", len(permissions))
	}
}

func TestAuthorizationManagerConcurrency(t *testing.T) {
	authManager := NewAuthorizationManager()

	// 测试并发注册角色
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			role := NewRole("并发角色", fmt.Sprintf("concurrent-role-%d", index), "并发测试角色")
			authManager.RegisterRole(role)
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有角色都被注册
	for i := 0; i < 10; i++ {
		role, err := authManager.GetRole(fmt.Sprintf("concurrent-role-%d", i))
		if err != nil {
			t.Errorf("Expected role concurrent-role-%d to be registered, got error: %v", i, err)
		}
		if role.GetName() != "并发角色" {
			t.Errorf("Expected role name to be '并发角色', got %s", role.GetName())
		}
	}
} 