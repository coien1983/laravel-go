# Laravel-Go 认证系统

Laravel-Go 认证系统提供了完整的用户认证和授权功能，支持多种认证方式，包括 Session 认证、JWT 认证等。

## 功能特性

### 🔐 认证方式
- **Session 认证**: 基于 Session 的传统 Web 应用认证
- **JWT 认证**: 基于 JSON Web Token 的 API 认证
- **多守卫支持**: 支持多个认证守卫同时使用
- **用户提供者**: 灵活的用户数据源支持

### 🛡️ 中间件
- **认证中间件**: 保护需要登录的路由
- **JWT 中间件**: 验证 JWT 令牌
- **访客中间件**: 确保只有未认证用户可以访问
- **可选认证中间件**: 支持可选的认证状态
- **角色中间件**: 基于角色的访问控制
- **权限中间件**: 基于权限的访问控制

### 👥 用户管理
- **用户接口**: 统一的用户模型接口
- **用户提供者**: 支持多种用户数据源
- **凭据验证**: 安全的密码验证机制
- **记住我**: 长期登录支持

## 核心组件

### AuthManager (认证管理器)
认证系统的核心管理器，负责管理多个守卫和用户提供者。

```go
// 创建认证管理器
manager := auth.NewAuthManager()

// 设置默认守卫
manager.SetDefaultGuard("web")

// 获取守卫
guard := manager.Guard("web")
```

### Guard (认证守卫)
认证守卫负责处理具体的认证逻辑。

#### SessionGuard (Session 守卫)
```go
// 创建 Session 守卫
provider := auth.NewMemoryUserProvider()
session := auth.NewMemorySessionStore()
guard := auth.NewSessionGuard(provider, session)

// 认证用户
credentials := map[string]interface{}{
    "email":    "user@example.com",
    "password": "password",
}
user, err := guard.Authenticate(credentials)

// 登录用户
err = guard.Login(user)

// 检查认证状态
if guard.Check() {
    // 用户已认证
}

// 获取当前用户
currentUser := guard.User()

// 登出用户
err = guard.Logout()
```

#### JWTGuard (JWT 守卫)
```go
// 创建 JWT 守卫
provider := auth.NewMemoryUserProvider()
secret := "your-secret-key"
ttl := 1 * time.Hour
guard := auth.NewJWTGuard(provider, secret, ttl)

// 生成 JWT 令牌
token, err := guard.GenerateToken(user)

// 验证 JWT 令牌
claims, err := guard.ValidateToken(token)

// 从令牌获取用户
user, err := guard.GetUserFromToken(token)

// 生成刷新令牌
refreshToken, err := guard.GenerateRefreshToken(user)

// 刷新令牌
newToken, err := guard.RefreshToken(refreshToken)
```

### UserProvider (用户提供者)
用户提供者负责从数据源检索和验证用户。

#### MemoryUserProvider (内存用户提供者)
```go
// 创建内存用户提供者
provider := auth.NewMemoryUserProvider()

// 添加用户
user := &auth.BaseUser{
    ID:       1,
    Email:    "user@example.com",
    Password: "password",
}
provider.AddUser(user)

// 通过 ID 检索用户
user, err := provider.RetrieveById(1)

// 通过凭据检索用户
credentials := map[string]interface{}{
    "email": "user@example.com",
}
user, err := provider.RetrieveByCredentials(credentials)

// 验证凭据
valid := provider.ValidateCredentials(user, credentials)
```

#### DatabaseUserProvider (数据库用户提供者)
```go
// 创建数据库用户提供者
connection := // 数据库连接
table := "users"
hashKey := "password"
provider := auth.NewDatabaseUserProvider(connection, table, hashKey)
```

### User (用户接口)
用户模型需要实现 User 接口。

```go
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
```

#### BaseUser (基础用户实现)
```go
user := &auth.BaseUser{
    ID:            1,
    Email:         "user@example.com",
    Password:      "password",
    RememberToken: "remember_token",
}
```

### SessionStore (Session 存储)
Session 存储接口定义了 Session 数据的存储方式。

```go
type SessionStore interface {
    Get(key string) interface{}
    Put(key string, value interface{})
    Forget(key string)
    Has(key string) bool
}
```

#### MemorySessionStore (内存 Session 存储)
```go
// 创建内存 Session 存储
session := auth.NewMemorySessionStore()

// 存储数据
session.Put("user_id", 1)

// 获取数据
userID := session.Get("user_id")

// 检查数据是否存在
if session.Has("user_id") {
    // 数据存在
}

// 删除数据
session.Forget("user_id")
```

## 中间件使用

### 认证中间件
```go
// 创建认证中间件
guard := auth.NewSessionGuard(provider, session)
middleware := auth.NewAuthMiddleware(guard)

// 在路由中使用
http.HandleFunc("/protected", middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
    // 只有认证用户才能访问
    w.Write([]byte("Protected content"))
}))
```

### JWT 中间件
```go
// 创建 JWT 中间件
guard := auth.NewJWTGuard(provider, secret, ttl)
middleware := auth.NewJWTMiddleware(guard)

// 在路由中使用
http.HandleFunc("/api/protected", middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
    // 只有有效 JWT 令牌才能访问
    w.Write([]byte("API protected content"))
}))
```

### 访客中间件
```go
// 创建访客中间件
guard := auth.NewSessionGuard(provider, session)
middleware := auth.NewGuestMiddleware(guard)

// 在路由中使用
http.HandleFunc("/login", middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
    // 只有未认证用户才能访问
    w.Write([]byte("Login page"))
}))
```

### 可选认证中间件
```go
// 创建可选认证中间件
guard := auth.NewJWTGuard(provider, secret, ttl)
middleware := auth.NewOptionalAuthMiddleware(guard)

// 在路由中使用
http.HandleFunc("/public", middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
    // 认证用户和访客都可以访问
    w.Write([]byte("Public content"))
}))
```

## 完整示例

### Web 应用认证
```go
package main

import (
    "net/http"
    "laravel-go/framework/auth"
)

func main() {
    // 创建用户提供者
    provider := auth.NewMemoryUserProvider()
    
    // 添加测试用户
    user := &auth.BaseUser{
        ID:       1,
        Email:    "admin@example.com",
        Password: "password",
    }
    provider.AddUser(user)

    // 创建 Session 存储
    session := auth.NewMemorySessionStore()

    // 创建 Session 守卫
    guard := auth.NewSessionGuard(provider, session)

    // 创建认证中间件
    authMiddleware := auth.NewAuthMiddleware(guard)

    // 登录路由
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            // 处理登录
            credentials := map[string]interface{}{
                "email":    r.FormValue("email"),
                "password": r.FormValue("password"),
            }
            
            user, err := guard.Authenticate(credentials)
            if err != nil {
                http.Error(w, "Invalid credentials", http.StatusUnauthorized)
                return
            }
            
            guard.Login(user)
            http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
        } else {
            // 显示登录页面
            w.Write([]byte(`
                <form method="POST">
                    <input name="email" placeholder="Email" required><br>
                    <input name="password" type="password" placeholder="Password" required><br>
                    <button type="submit">Login</button>
                </form>
            `))
        }
    })

    // 受保护的路由
    http.HandleFunc("/dashboard", authMiddleware.Handle(func(w http.ResponseWriter, r *http.Request) {
        currentUser := guard.User()
        w.Write([]byte("Welcome, " + currentUser.GetEmail()))
    }))

    // 登出路由
    http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
        guard.Logout()
        http.Redirect(w, r, "/login", http.StatusSeeOther)
    })

    http.ListenAndServe(":8080", nil)
}
```

### API 认证
```go
package main

import (
    "encoding/json"
    "net/http"
    "time"
    "laravel-go/framework/auth"
)

func main() {
    // 创建用户提供者
    provider := auth.NewMemoryUserProvider()
    
    // 添加测试用户
    user := &auth.BaseUser{
        ID:       1,
        Email:    "api@example.com",
        Password: "password",
    }
    provider.AddUser(user)

    // 创建 JWT 守卫
    secret := "your-secret-key"
    ttl := 1 * time.Hour
    guard := auth.NewJWTGuard(provider, secret, ttl)

    // 创建 JWT 中间件
    jwtMiddleware := auth.NewJWTMiddleware(guard)

    // 登录 API
    http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            var credentials map[string]interface{}
            json.NewDecoder(r.Body).Decode(&credentials)
            
            user, err := guard.Authenticate(credentials)
            if err != nil {
                http.Error(w, "Invalid credentials", http.StatusUnauthorized)
                return
            }
            
            token, err := guard.GenerateToken(user)
            if err != nil {
                http.Error(w, "Token generation failed", http.StatusInternalServerError)
                return
            }
            
            response := map[string]string{
                "token": token,
                "type":  "Bearer",
            }
            
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(response)
        }
    })

    // 受保护的 API 路由
    http.HandleFunc("/api/protected", jwtMiddleware.Handle(func(w http.ResponseWriter, r *http.Request) {
        currentUser := guard.User()
        response := map[string]interface{}{
            "message": "Protected API endpoint",
            "user":    currentUser.GetEmail(),
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }))

    // 刷新令牌 API
    http.HandleFunc("/api/refresh", func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        
        token := authHeader[7:] // 移除 "Bearer " 前缀
        newToken, err := guard.RefreshToken(token)
        if err != nil {
            http.Error(w, "Token refresh failed", http.StatusUnauthorized)
            return
        }
        
        response := map[string]string{
            "token": newToken,
            "type":  "Bearer",
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    http.ListenAndServe(":8080", nil)
}
```

## 错误处理

认证系统定义了以下错误类型：

```go
var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserNotFound       = errors.New("user not found")
    ErrUserNotAuthenticated = errors.New("user not authenticated")
    ErrInvalidToken       = errors.New("invalid token")
    ErrTokenExpired       = errors.New("token expired")
)
```

## 最佳实践

### 1. 密码安全
- 使用强密码哈希算法（如 bcrypt）
- 实现密码复杂度验证
- 定期要求用户更改密码

### 2. JWT 安全
- 使用强密钥
- 设置合理的令牌过期时间
- 实现令牌黑名单机制
- 使用 HTTPS 传输令牌

### 3. Session 安全
- 使用安全的 Session 存储
- 实现 Session 固定攻击防护
- 定期清理过期 Session

### 4. 错误处理
- 不要泄露敏感信息
- 记录认证失败事件
- 实现账户锁定机制

### 5. 多因素认证
- 支持 TOTP/HOTP
- 实现短信验证码
- 提供备用认证方式

## 扩展功能

### 自定义用户提供者
```go
type CustomUserProvider struct {
    // 自定义字段
}

func (p *CustomUserProvider) RetrieveById(identifier interface{}) (auth.User, error) {
    // 实现从自定义数据源检索用户
}

func (p *CustomUserProvider) RetrieveByCredentials(credentials map[string]interface{}) (auth.User, error) {
    // 实现通过凭据检索用户
}

// 实现其他接口方法...
```

### 自定义认证守卫
```go
type CustomGuard struct {
    // 自定义字段
}

func (g *CustomGuard) Authenticate(credentials map[string]interface{}) (auth.User, error) {
    // 实现自定义认证逻辑
}

// 实现其他接口方法...
```

## 测试

认证系统包含完整的单元测试：

```bash
# 运行认证系统测试
go test ./framework/auth/... -v
```

## 示例程序

运行认证系统演示程序：

```bash
# 进入示例目录
cd examples/auth_demo

# 运行演示
go run main.go
```

## 总结

Laravel-Go 认证系统提供了完整、灵活、安全的用户认证解决方案，支持多种认证方式，可以满足不同应用场景的需求。通过统一的接口设计，可以轻松扩展和定制认证功能。 