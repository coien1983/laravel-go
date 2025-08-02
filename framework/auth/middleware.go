package auth

import (
	"net/http"
	"strings"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	guard Guard
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(guard Guard) *AuthMiddleware {
	return &AuthMiddleware{
		guard: guard,
	}
}

// Handle 处理HTTP请求
func (am *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查用户是否已认证
		if !am.guard.Check() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// JWTMiddleware JWT认证中间件
type JWTMiddleware struct {
	guard Guard
}

// NewJWTMiddleware 创建JWT认证中间件
func NewJWTMiddleware(guard Guard) *JWTMiddleware {
	return &JWTMiddleware{
		guard: guard,
	}
}

// Handle 处理HTTP请求
func (jm *JWTMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取JWT令牌
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// 提取令牌
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 验证令牌
		if jwtGuard, ok := jm.guard.(*JWTGuard); ok {
			user, err := jwtGuard.GetUserFromToken(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// 设置用户到守卫
			jm.guard.SetUser(user)
		} else {
			http.Error(w, "JWT guard required", http.StatusInternalServerError)
			return
		}

		next(w, r)
	}
}

// GuestMiddleware 访客中间件（确保用户未认证）
type GuestMiddleware struct {
	guard Guard
}

// NewGuestMiddleware 创建访客中间件
func NewGuestMiddleware(guard Guard) *GuestMiddleware {
	return &GuestMiddleware{
		guard: guard,
	}
}

// Handle 处理HTTP请求
func (gm *GuestMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查用户是否已认证
		if gm.guard.Check() {
			http.Error(w, "Already authenticated", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// OptionalAuthMiddleware 可选认证中间件（认证用户可选）
type OptionalAuthMiddleware struct {
	guard Guard
}

// NewOptionalAuthMiddleware 创建可选认证中间件
func NewOptionalAuthMiddleware(guard Guard) *OptionalAuthMiddleware {
	return &OptionalAuthMiddleware{
		guard: guard,
	}
}

// Handle 处理HTTP请求
func (oam *OptionalAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 尝试从请求中获取认证信息
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")

			if jwtGuard, ok := oam.guard.(*JWTGuard); ok {
				if user, err := jwtGuard.GetUserFromToken(token); err == nil {
					oam.guard.SetUser(user)
				}
			}
		}

		next(w, r)
	}
}

// RoleMiddleware 角色中间件
type RoleMiddleware struct {
	guard                Guard
	authorizationManager *AuthorizationManager
	roles                []string
	requireAll           bool // true: 需要所有角色, false: 需要任一角色
}

// NewRoleMiddleware 创建角色中间件
func NewRoleMiddleware(guard Guard, authManager *AuthorizationManager, roles []string, requireAll bool) *RoleMiddleware {
	return &RoleMiddleware{
		guard:                guard,
		authorizationManager: authManager,
		roles:                roles,
		requireAll:           requireAll,
	}
}

// Handle 处理HTTP请求
func (rm *RoleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查用户是否已认证
		if !rm.guard.Check() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 检查用户角色
		user := rm.guard.User()
		if user == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// 检查角色权限
		if !rm.checkRoles(user) {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// checkRoles 检查用户角色
func (rm *RoleMiddleware) checkRoles(user User) bool {
	if len(rm.roles) == 0 {
		return true
	}

	// 获取用户角色（这里需要根据实际的用户模型来获取）
	userRoles := rm.getUserRoles(user)

	if rm.requireAll {
		// 需要所有角色
		for _, requiredRole := range rm.roles {
			found := false
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	} else {
		// 需要任一角色
		for _, requiredRole := range rm.roles {
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					return true
				}
			}
		}
		return false
	}
}

// getUserRoles 获取用户角色
func (rm *RoleMiddleware) getUserRoles(user User) []string {
	// 这里需要根据实际的用户模型来获取用户角色
	// 暂时返回空切片，需要在实际使用时实现
	return []string{}
}

// PermissionMiddleware 权限中间件
type PermissionMiddleware struct {
	guard                Guard
	authorizationManager *AuthorizationManager
	permissions          []string
	requireAll           bool // true: 需要所有权限, false: 需要任一权限
}

// NewPermissionMiddleware 创建权限中间件
func NewPermissionMiddleware(guard Guard, authManager *AuthorizationManager, permissions []string, requireAll bool) *PermissionMiddleware {
	return &PermissionMiddleware{
		guard:                guard,
		authorizationManager: authManager,
		permissions:          permissions,
		requireAll:           requireAll,
	}
}

// Handle 处理HTTP请求
func (pm *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 检查用户是否已认证
		if !pm.guard.Check() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 检查用户权限
		user := pm.guard.User()
		if user == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// 检查权限
		if !pm.checkPermissions(user) {
			http.Error(w, "Insufficient permissions", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// checkPermissions 检查用户权限
func (pm *PermissionMiddleware) checkPermissions(user User) bool {
	if len(pm.permissions) == 0 {
		return true
	}

	if pm.requireAll {
		// 需要所有权限
		for _, permission := range pm.permissions {
			if !pm.authorizationManager.Can(user, permission, nil) {
				return false
			}
		}
		return true
	} else {
		// 需要任一权限
		for _, permission := range pm.permissions {
			if pm.authorizationManager.Can(user, permission, nil) {
				return true
			}
		}
		return false
	}
}

// 便捷的中间件创建函数

// RequirePermission 创建需要指定权限的中间件
func RequirePermission(guard Guard, authManager *AuthorizationManager, permission string) *PermissionMiddleware {
	return NewPermissionMiddleware(guard, authManager, []string{permission}, false)
}

// RequireAllPermissions 创建需要所有指定权限的中间件
func RequireAllPermissions(guard Guard, authManager *AuthorizationManager, permissions []string) *PermissionMiddleware {
	return NewPermissionMiddleware(guard, authManager, permissions, true)
}

// RequireAnyPermission 创建需要任一指定权限的中间件
func RequireAnyPermission(guard Guard, authManager *AuthorizationManager, permissions []string) *PermissionMiddleware {
	return NewPermissionMiddleware(guard, authManager, permissions, false)
}

// RequireRole 创建需要指定角色的中间件
func RequireRole(guard Guard, authManager *AuthorizationManager, role string) *RoleMiddleware {
	return NewRoleMiddleware(guard, authManager, []string{role}, false)
}

// RequireAllRoles 创建需要所有指定角色的中间件
func RequireAllRoles(guard Guard, authManager *AuthorizationManager, roles []string) *RoleMiddleware {
	return NewRoleMiddleware(guard, authManager, roles, true)
}

// RequireAnyRole 创建需要任一指定角色的中间件
func RequireAnyRole(guard Guard, authManager *AuthorizationManager, roles []string) *RoleMiddleware {
	return NewRoleMiddleware(guard, authManager, roles, false)
} 