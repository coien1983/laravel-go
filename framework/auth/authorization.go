package auth

import (
	"errors"
	"sync"
)

// Permission 权限接口
type Permission interface {
	GetName() string
	GetSlug() string
	GetDescription() string
	GetResource() string
	GetAction() string
}

// Role 角色接口
type Role interface {
	GetName() string
	GetSlug() string
	GetDescription() string
	GetPermissions() []Permission
	HasPermission(permission Permission) bool
	HasPermissionByName(name string) bool
	AddPermission(permission Permission)
	RemovePermission(permission Permission)
}

// Policy 策略接口
type Policy interface {
	// 检查用户是否有权限执行操作
	Can(user User, action string, resource interface{}) bool
	// 检查用户是否可以查看资源
	CanView(user User, resource interface{}) bool
	// 检查用户是否可以创建资源
	CanCreate(user User, resource interface{}) bool
	// 检查用户是否可以更新资源
	CanUpdate(user User, resource interface{}) bool
	// 检查用户是否可以删除资源
	CanDelete(user User, resource interface{}) bool
}

// AuthorizationManager 授权管理器
type AuthorizationManager struct {
	roles    map[string]Role
	policies map[string]Policy
	mu       sync.RWMutex
}

// NewAuthorizationManager 创建授权管理器
func NewAuthorizationManager() *AuthorizationManager {
	return &AuthorizationManager{
		roles:    make(map[string]Role),
		policies: make(map[string]Policy),
	}
}

// RegisterRole 注册角色
func (am *AuthorizationManager) RegisterRole(role Role) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.roles[role.GetSlug()] = role
}

// GetRole 获取角色
func (am *AuthorizationManager) GetRole(slug string) (Role, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	if role, exists := am.roles[slug]; exists {
		return role, nil
	}
	return nil, errors.New("role not found")
}

// RegisterPolicy 注册策略
func (am *AuthorizationManager) RegisterPolicy(name string, policy Policy) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.policies[name] = policy
}

// GetPolicy 获取策略
func (am *AuthorizationManager) GetPolicy(name string) (Policy, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	if policy, exists := am.policies[name]; exists {
		return policy, nil
	}
	return nil, errors.New("policy not found")
}

// Can 检查用户是否有权限执行操作
func (am *AuthorizationManager) Can(user User, action string, resource interface{}) bool {
	// 首先检查用户是否有超级管理员权限
	if am.isSuperUser(user) {
		return true
	}

	// 检查用户角色权限
	if am.checkRolePermissions(user, action, resource) {
		return true
	}

	// 检查策略
	if am.checkPolicies(user, action, resource) {
		return true
	}

	return false
}

// CanView 检查用户是否可以查看资源
func (am *AuthorizationManager) CanView(user User, resource interface{}) bool {
	return am.Can(user, "view", resource)
}

// CanCreate 检查用户是否可以创建资源
func (am *AuthorizationManager) CanCreate(user User, resource interface{}) bool {
	return am.Can(user, "create", resource)
}

// CanUpdate 检查用户是否可以更新资源
func (am *AuthorizationManager) CanUpdate(user User, resource interface{}) bool {
	return am.Can(user, "update", resource)
}

// CanDelete 检查用户是否可以删除资源
func (am *AuthorizationManager) CanDelete(user User, resource interface{}) bool {
	return am.Can(user, "delete", resource)
}

// isSuperUser 检查是否为超级用户
func (am *AuthorizationManager) isSuperUser(user User) bool {
	// 这里可以根据实际需求实现超级用户检查逻辑
	// 例如检查用户是否有特定的超级用户标识
	return false
}

// checkRolePermissions 检查角色权限
func (am *AuthorizationManager) checkRolePermissions(user User, action string, resource interface{}) bool {
	// 这里需要根据实际的用户模型来获取用户角色
	// 暂时返回 false，需要在实际使用时实现
	return false
}

// checkPolicies 检查策略
func (am *AuthorizationManager) checkPolicies(user User, action string, resource interface{}) bool {
	am.mu.RLock()
	defer am.mu.RUnlock()

	for _, policy := range am.policies {
		if policy.Can(user, action, resource) {
			return true
		}
	}
	return false
}

// BasePermission 基础权限实现
type BasePermission struct {
	name        string
	slug        string
	description string
	resource    string
	action      string
}

// NewPermission 创建权限
func NewPermission(name, slug, description, resource, action string) *BasePermission {
	return &BasePermission{
		name:        name,
		slug:        slug,
		description: description,
		resource:    resource,
		action:      action,
	}
}

func (p *BasePermission) GetName() string {
	return p.name
}

func (p *BasePermission) GetSlug() string {
	return p.slug
}

func (p *BasePermission) GetDescription() string {
	return p.description
}

func (p *BasePermission) GetResource() string {
	return p.resource
}

func (p *BasePermission) GetAction() string {
	return p.action
}

// BaseRole 基础角色实现
type BaseRole struct {
	name        string
	slug        string
	description string
	permissions []Permission
	mu          sync.RWMutex
}

// NewRole 创建角色
func NewRole(name, slug, description string) *BaseRole {
	return &BaseRole{
		name:        name,
		slug:        slug,
		description: description,
		permissions: make([]Permission, 0),
	}
}

func (r *BaseRole) GetName() string {
	return r.name
}

func (r *BaseRole) GetSlug() string {
	return r.slug
}

func (r *BaseRole) GetDescription() string {
	return r.description
}

func (r *BaseRole) GetPermissions() []Permission {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.permissions
}

func (r *BaseRole) HasPermission(permission Permission) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.permissions {
		if p.GetSlug() == permission.GetSlug() {
			return true
		}
	}
	return false
}

func (r *BaseRole) HasPermissionByName(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.permissions {
		if p.GetName() == name {
			return true
		}
	}
	return false
}

func (r *BaseRole) AddPermission(permission Permission) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// 检查是否已存在该权限
	exists := false
	for _, p := range r.permissions {
		if p.GetSlug() == permission.GetSlug() {
			exists = true
			break
		}
	}
	
	if !exists {
		r.permissions = append(r.permissions, permission)
	}
}

func (r *BaseRole) RemovePermission(permission Permission) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, p := range r.permissions {
		if p.GetSlug() == permission.GetSlug() {
			r.permissions = append(r.permissions[:i], r.permissions[i+1:]...)
			break
		}
	}
}

// BasePolicy 基础策略实现
type BasePolicy struct {
	name string
}

// NewPolicy 创建策略
func NewPolicy(name string) *BasePolicy {
	return &BasePolicy{
		name: name,
	}
}

func (p *BasePolicy) Can(user User, action string, resource interface{}) bool {
	// 基础策略默认拒绝所有操作
	return false
}

func (p *BasePolicy) CanView(user User, resource interface{}) bool {
	return p.Can(user, "view", resource)
}

func (p *BasePolicy) CanCreate(user User, resource interface{}) bool {
	return p.Can(user, "create", resource)
}

func (p *BasePolicy) CanUpdate(user User, resource interface{}) bool {
	return p.Can(user, "update", resource)
}

func (p *BasePolicy) CanDelete(user User, resource interface{}) bool {
	return p.Can(user, "delete", resource)
}

// ResourcePolicy 资源策略
type ResourcePolicy struct {
	*BasePolicy
	resourceType string
}

// NewResourcePolicy 创建资源策略
func NewResourcePolicy(name, resourceType string) *ResourcePolicy {
	return &ResourcePolicy{
		BasePolicy:   NewPolicy(name),
		resourceType: resourceType,
	}
}

func (p *ResourcePolicy) Can(user User, action string, resource interface{}) bool {
	// 这里可以实现具体的资源权限检查逻辑
	// 例如检查资源的所有者、用户角色等
	return false
}

// 预定义权限
var (
	// 用户管理权限
	UserViewPermission   = NewPermission("查看用户", "user.view", "查看用户信息", "user", "view")
	UserCreatePermission = NewPermission("创建用户", "user.create", "创建新用户", "user", "create")
	UserUpdatePermission = NewPermission("更新用户", "user.update", "更新用户信息", "user", "update")
	UserDeletePermission = NewPermission("删除用户", "user.delete", "删除用户", "user", "delete")

	// 角色管理权限
	RoleViewPermission   = NewPermission("查看角色", "role.view", "查看角色信息", "role", "view")
	RoleCreatePermission = NewPermission("创建角色", "role.create", "创建新角色", "role", "create")
	RoleUpdatePermission = NewPermission("更新角色", "role.update", "更新角色信息", "role", "update")
	RoleDeletePermission = NewPermission("删除角色", "role.delete", "删除角色", "role", "delete")

	// 权限管理权限
	PermissionViewPermission   = NewPermission("查看权限", "permission.view", "查看权限信息", "permission", "view")
	PermissionCreatePermission = NewPermission("创建权限", "permission.create", "创建新权限", "permission", "create")
	PermissionUpdatePermission = NewPermission("更新权限", "permission.update", "更新权限信息", "permission", "update")
	PermissionDeletePermission = NewPermission("删除权限", "permission.delete", "删除权限", "permission", "delete")
)

// 预定义角色
var (
	// 超级管理员角色
	SuperAdminRole = func() Role {
		role := NewRole("超级管理员", "super-admin", "拥有所有权限的超级管理员")
		role.AddPermission(UserViewPermission)
		role.AddPermission(UserCreatePermission)
		role.AddPermission(UserUpdatePermission)
		role.AddPermission(UserDeletePermission)
		role.AddPermission(RoleViewPermission)
		role.AddPermission(RoleCreatePermission)
		role.AddPermission(RoleUpdatePermission)
		role.AddPermission(RoleDeletePermission)
		role.AddPermission(PermissionViewPermission)
		role.AddPermission(PermissionCreatePermission)
		role.AddPermission(PermissionUpdatePermission)
		role.AddPermission(PermissionDeletePermission)
		return role
	}()

	// 管理员角色
	AdminRole = func() Role {
		role := NewRole("管理员", "admin", "系统管理员")
		role.AddPermission(UserViewPermission)
		role.AddPermission(UserCreatePermission)
		role.AddPermission(UserUpdatePermission)
		role.AddPermission(RoleViewPermission)
		role.AddPermission(PermissionViewPermission)
		return role
	}()

	// 普通用户角色
	UserRole = func() Role {
		role := NewRole("普通用户", "user", "普通用户")
		role.AddPermission(UserViewPermission)
		return role
	}()
)

// 错误定义
var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrRoleNotFound     = errors.New("role not found")
	ErrPolicyNotFound   = errors.New("policy not found")
) 