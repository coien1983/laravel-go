# Laravel-Go 授权系统

## 概述

Laravel-Go 授权系统提供了完整的权限和角色管理功能，支持基于角色的访问控制（RBAC）和基于策略的授权。该系统与认证系统紧密集成，为应用程序提供强大的安全保护。

## 核心组件

### 1. 权限 (Permission)

权限是系统中最基本的授权单位，定义了用户可以执行的具体操作。

#### 权限接口

```go
type Permission interface {
    GetName() string        // 权限名称
    GetSlug() string        // 权限标识符
    GetDescription() string // 权限描述
    GetResource() string    // 资源类型
    GetAction() string      // 操作类型
}
```

#### 创建权限

```go
// 创建基础权限
permission := auth.NewPermission("查看用户", "user.view", "查看用户信息", "user", "view")

// 使用预定义权限
userViewPermission := auth.UserViewPermission
userCreatePermission := auth.UserCreatePermission
userUpdatePermission := auth.UserUpdatePermission
userDeletePermission := auth.UserDeletePermission
```

### 2. 角色 (Role)

角色是权限的集合，用于简化权限管理。

#### 角色接口

```go
type Role interface {
    GetName() string                    // 角色名称
    GetSlug() string                    // 角色标识符
    GetDescription() string             // 角色描述
    GetPermissions() []Permission       // 获取所有权限
    HasPermission(permission Permission) bool // 检查是否有指定权限
    HasPermissionByName(name string) bool     // 根据名称检查权限
    AddPermission(permission Permission)      // 添加权限
    RemovePermission(permission Permission)   // 移除权限
}
```

#### 创建角色

```go
// 创建基础角色
role := auth.NewRole("管理员", "admin", "系统管理员")

// 添加权限
role.AddPermission(auth.UserViewPermission)
role.AddPermission(auth.UserCreatePermission)
role.AddPermission(auth.UserUpdatePermission)

// 使用预定义角色
superAdminRole := auth.SuperAdminRole
adminRole := auth.AdminRole
userRole := auth.UserRole
```

### 3. 策略 (Policy)

策略定义了复杂的授权逻辑，可以基于用户、资源和操作进行细粒度的权限控制。

#### 策略接口

```go
type Policy interface {
    Can(user User, action string, resource interface{}) bool // 检查权限
    CanView(user User, resource interface{}) bool            // 检查查看权限
    CanCreate(user User, resource interface{}) bool          // 检查创建权限
    CanUpdate(user User, resource interface{}) bool          // 检查更新权限
    CanDelete(user User, resource interface{}) bool          // 检查删除权限
}
```

#### 创建策略

```go
// 创建基础策略
policy := auth.NewPolicy("user-policy")

// 创建资源策略
resourcePolicy := auth.NewResourcePolicy("post-policy", "post")

// 创建自定义策略
type UserPolicy struct {
    *auth.BasePolicy
}

func (p *UserPolicy) Can(user auth.User, action string, resource interface{}) bool {
    // 实现自定义权限逻辑
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
```

### 4. 授权管理器 (AuthorizationManager)

授权管理器是授权系统的核心，负责管理角色、策略和权限检查。

#### 创建授权管理器

```go
authManager := auth.NewAuthorizationManager()
```

#### 注册角色和策略

```go
// 注册角色
authManager.RegisterRole(superAdminRole)
authManager.RegisterRole(adminRole)
authManager.RegisterRole(userRole)

// 注册策略
authManager.RegisterPolicy("user-policy", userPolicy)
```

#### 权限检查

```go
// 检查基本权限
canView := authManager.CanView(user, resource)
canCreate := authManager.CanCreate(user, resource)
canUpdate := authManager.CanUpdate(user, resource)
canDelete := authManager.CanDelete(user, resource)

// 检查自定义权限
canExecute := authManager.Can(user, "execute", resource)
```

## 中间件

授权系统提供了多种中间件来保护 HTTP 路由。

### 权限中间件

```go
// 需要单个权限
permissionMiddleware := auth.RequirePermission(guard, authManager, "user.view")

// 需要所有权限
allPermissionsMiddleware := auth.RequireAllPermissions(guard, authManager, []string{"user.view", "user.create"})

// 需要任一权限
anyPermissionMiddleware := auth.RequireAnyPermission(guard, authManager, []string{"user.view", "user.create"})
```

### 角色中间件

```go
// 需要单个角色
roleMiddleware := auth.RequireRole(guard, authManager, "admin")

// 需要所有角色
allRolesMiddleware := auth.RequireAllRoles(guard, authManager, []string{"admin", "moderator"})

// 需要任一角色
anyRoleMiddleware := auth.RequireAnyRole(guard, authManager, []string{"admin", "super-admin"})
```

### 使用中间件

```go
// 创建处理器
handler := func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "访问成功！")
}

// 使用中间件包装处理器
protectedHandler := permissionMiddleware.Handle(handler)

// 在路由中使用
http.HandleFunc("/users", protectedHandler)
```

## 预定义权限和角色

### 预定义权限

```go
// 用户管理权限
UserViewPermission   = "user.view"
UserCreatePermission = "user.create"
UserUpdatePermission = "user.update"
UserDeletePermission = "user.delete"

// 角色管理权限
RoleViewPermission   = "role.view"
RoleCreatePermission = "role.create"
RoleUpdatePermission = "role.update"
RoleDeletePermission = "role.delete"

// 权限管理权限
PermissionViewPermission   = "permission.view"
PermissionCreatePermission = "permission.create"
PermissionUpdatePermission = "permission.update"
PermissionDeletePermission = "permission.delete"
```

### 预定义角色

```go
// 超级管理员角色
SuperAdminRole = {
    Name: "超级管理员",
    Slug: "super-admin",
    Description: "拥有所有权限的超级管理员",
    Permissions: [所有权限]
}

// 管理员角色
AdminRole = {
    Name: "管理员",
    Slug: "admin",
    Description: "系统管理员",
    Permissions: [用户查看、创建、更新, 角色查看, 权限查看]
}

// 普通用户角色
UserRole = {
    Name: "普通用户",
    Slug: "user",
    Description: "普通用户",
    Permissions: [用户查看]
}
```

## 完整示例

### 基本使用

```go
package main

import (
    "net/http"
    "laravel-go/framework/auth"
)

func main() {
    // 创建授权管理器
    authManager := auth.NewAuthorizationManager()
    
    // 注册预定义角色
    authManager.RegisterRole(auth.SuperAdminRole)
    authManager.RegisterRole(auth.AdminRole)
    authManager.RegisterRole(auth.UserRole)
    
    // 创建用户提供者和守卫
    provider := auth.NewMemoryUserProvider()
    sessionStore := auth.NewMemorySessionStore()
    guard := auth.NewSessionGuard(provider, sessionStore)
    
    // 创建用户
    user := &auth.BaseUser{
        ID:       1,
        Email:    "user@example.com",
        Password: "password",
    }
    provider.AddUser(user)
    guard.SetUser(user)
    
    // 检查权限
    canView := authManager.CanView(user, nil)
    canCreate := authManager.CanCreate(user, nil)
    
    // 创建中间件
    permissionMiddleware := auth.RequirePermission(guard, authManager, "user.view")
    
    // 创建处理器
    handler := func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "访问成功！")
    }
    
    // 使用中间件
    protectedHandler := permissionMiddleware.Handle(handler)
    
    // 启动服务器
    http.HandleFunc("/protected", protectedHandler)
    http.ListenAndServe(":8080", nil)
}
```

### 自定义策略

```go
type PostPolicy struct {
    *auth.BasePolicy
}

func (p *PostPolicy) CanView(user auth.User, resource interface{}) bool {
    // 所有用户都可以查看公开的帖子
    return true
}

func (p *PostPolicy) CanCreate(user auth.User, resource interface{}) bool {
    // 只有认证用户才能创建帖子
    return user != nil
}

func (p *PostPolicy) CanUpdate(user auth.User, resource interface{}) bool {
    // 只有帖子作者或管理员才能更新帖子
    if post, ok := resource.(*Post); ok {
        return user.GetID() == post.AuthorID || p.isAdmin(user)
    }
    return false
}

func (p *PostPolicy) CanDelete(user auth.User, resource interface{}) bool {
    // 只有管理员才能删除帖子
    return p.isAdmin(user)
}

func (p *PostPolicy) isAdmin(user auth.User) bool {
    // 实现管理员检查逻辑
    return user.GetID() == 1 // 假设ID为1的用户是管理员
}
```

## 最佳实践

### 1. 权限设计

- 使用清晰的权限命名约定（如 `resource.action`）
- 为每个资源定义完整的 CRUD 权限
- 避免过于细粒度的权限，保持权限管理的简洁性

### 2. 角色设计

- 基于业务需求设计角色层次结构
- 使用继承关系简化权限管理
- 定期审查和更新角色权限

### 3. 策略实现

- 在策略中实现复杂的业务逻辑
- 使用策略来处理资源所有者检查
- 保持策略的单一职责原则

### 4. 性能优化

- 缓存用户权限和角色信息
- 避免在策略中进行复杂的数据库查询
- 使用索引优化权限检查查询

### 5. 安全考虑

- 默认拒绝所有权限，明确授权
- 定期审计权限分配
- 实现权限变更日志记录

## 错误处理

授权系统定义了以下错误类型：

```go
var (
    ErrPermissionDenied = errors.New("permission denied")
    ErrRoleNotFound     = errors.New("role not found")
    ErrPolicyNotFound   = errors.New("policy not found")
)
```

## 扩展性

授权系统设计为高度可扩展的：

1. **自定义权限类型**：可以实现自定义的权限接口
2. **自定义角色逻辑**：可以扩展角色接口添加特定功能
3. **自定义策略**：可以实现复杂的授权逻辑
4. **中间件扩展**：可以创建自定义的授权中间件
5. **存储后端**：可以集成不同的存储后端（数据库、Redis等）

## 测试

授权系统包含完整的单元测试：

```bash
go test ./framework/auth -v
```

测试覆盖了以下方面：
- 权限创建和管理
- 角色创建和权限分配
- 策略实现和权限检查
- 授权管理器的功能
- 中间件的正确性
- 并发安全性

## 总结

Laravel-Go 授权系统提供了强大而灵活的权限管理功能，支持基于角色和策略的授权模式。通过合理的权限设计和使用，可以为应用程序提供可靠的安全保护。 