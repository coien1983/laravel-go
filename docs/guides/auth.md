# 认证授权指南

## 📖 概述

Laravel-Go Framework 提供了完整的认证和授权系统，支持多种认证方式和权限控制机制。

## 🚀 快速开始

### 1. 基本认证

```go
// 用户登录
func (c *AuthController) Login(request http.Request) http.Response {
    email := request.Body["email"].(string)
    password := request.Body["password"].(string)

    // 验证用户凭据
    user, token, err := auth.Attempt(email, password)
    if err != nil {
        return c.JsonError("Invalid credentials", 401)
    }

    return c.Json(map[string]interface{}{
        "user":  user,
        "token": token,
    })
}

// 用户注册
func (c *AuthController) Register(request http.Request) http.Response {
    data := request.Body

    // 创建用户
    user, err := auth.Register(data)
    if err != nil {
        return c.JsonError(err.Error(), 422)
    }

    // 生成认证 token
    token, err := auth.CreateToken(user)
    if err != nil {
        return c.JsonError("Failed to create token", 500)
    }

    return c.Json(map[string]interface{}{
        "user":  user,
        "token": token,
    }).Status(201)
}
```

### 2. 中间件保护

```go
// 认证中间件
type AuthMiddleware struct {
    http.Middleware
}

func (m *AuthMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    token := request.Headers["Authorization"]
    if token == "" {
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "Unauthorized"}`,
        }
    }

    // 验证 token
    user, err := auth.ValidateToken(token)
    if err != nil {
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "Invalid token"}`,
        }
    }

    // 将用户信息添加到请求上下文
    request.Context["user"] = user

    return next(request)
}

// 在路由中使用
router.Group("/api", func(group *routing.Router) {
    group.Use(&middleware.AuthMiddleware{})

    group.Get("/profile", &Controllers.UserController{}, "Profile")
    group.Post("/posts", &Controllers.PostController{}, "Store")
})
```

## 🔐 认证方式

### 1. JWT 认证

```go
// 配置 JWT
config.Set("auth.jwt.secret", "your-secret-key")
config.Set("auth.jwt.expire", 3600) // 1小时

// 生成 JWT token
func (s *AuthService) CreateToken(user *Models.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.Get("auth.jwt.secret")))
}

// 验证 JWT token
func (s *AuthService) ValidateToken(tokenString string) (*Models.User, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(config.Get("auth.jwt.secret")), nil
    })

    if err != nil || !token.Valid {
        return nil, errors.New("invalid token")
    }

    claims := token.Claims.(jwt.MapClaims)
    userID := uint(claims["user_id"].(float64))

    var user Models.User
    err = s.db.First(&user, userID).Error
    if err != nil {
        return nil, err
    }

    return &user, nil
}
```

### 2. Session 认证

```go
// 配置 Session
config.Set("session.driver", "redis")
config.Set("session.lifetime", 120)

// 启动 Session
func (c *AuthController) Login(request http.Request) http.Response {
    // 验证用户凭据...

    // 创建 Session
    session := request.Session()
    session.Put("user_id", user.ID)
    session.Put("user_email", user.Email)

    return c.Json(map[string]string{
        "message": "Login successful",
    })
}

// 验证 Session
func (m *AuthMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    session := request.Session()
    userID := session.Get("user_id")

    if userID == nil {
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "Unauthorized"}`,
        }
    }

    var user Models.User
    err := s.db.First(&user, userID).Error
    if err != nil {
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "User not found"}`,
        }
    }

    request.Context["user"] = &user
    return next(request)
}
```

### 3. API Token 认证

```go
// API Token 模型
type ApiToken struct {
    database.Model
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id" gorm:"not null"`
    Name      string    `json:"name" gorm:"size:255;not null"`
    Token     string    `json:"token" gorm:"size:64;unique;not null"`
    Abilities string    `json:"abilities" gorm:"type:text"`
    LastUsedAt *time.Time `json:"last_used_at"`
    ExpiresAt  *time.Time `json:"expires_at"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`

    User User `json:"user" gorm:"foreignKey:UserID"`
}

// 创建 API Token
func (s *AuthService) CreateApiToken(userID uint, name string, abilities []string) (*ApiToken, error) {
    token := &ApiToken{
        UserID:    userID,
        Name:      name,
        Token:     generateRandomToken(64),
        Abilities: strings.Join(abilities, ","),
    }

    err := s.db.Create(token).Error
    return token, err
}

// 验证 API Token
func (s *AuthService) ValidateApiToken(tokenString string) (*Models.User, error) {
    var apiToken ApiToken
    err := s.db.Where("token = ?", tokenString).First(&apiToken).Error
    if err != nil {
        return nil, err
    }

    // 检查是否过期
    if apiToken.ExpiresAt != nil && time.Now().After(*apiToken.ExpiresAt) {
        return nil, errors.New("token expired")
    }

    // 更新最后使用时间
    now := time.Now()
    s.db.Model(&apiToken).Update("last_used_at", &now)

    var user Models.User
    err = s.db.First(&user, apiToken.UserID).Error
    return &user, err
}
```

## 🔒 授权系统

### 1. 基于角色的权限控制 (RBAC)

```go
// 角色模型
type Role struct {
    database.Model
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"size:255;unique;not null"`
    DisplayName string    `json:"display_name" gorm:"size:255"`
    Description string    `json:"description" gorm:"type:text"`
    Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// 权限模型
type Permission struct {
    database.Model
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"size:255;unique;not null"`
    DisplayName string    `json:"display_name" gorm:"size:255"`
    Description string    `json:"description" gorm:"type:text"`
    Guard       string    `json:"guard" gorm:"size:255;default:'web'"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// 用户角色关联
type UserRole struct {
    UserID uint `json:"user_id" gorm:"primaryKey"`
    RoleID uint `json:"role_id" gorm:"primaryKey"`
}

// 检查用户权限
func (u *User) HasPermission(permission string) bool {
    var count int64
    s.db.Model(&Permission{}).
        Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
        Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
        Where("user_roles.user_id = ? AND permissions.name = ?", u.ID, permission).
        Count(&count)

    return count > 0
}

// 检查用户角色
func (u *User) HasRole(role string) bool {
    var count int64
    s.db.Model(&Role{}).
        Joins("JOIN user_roles ON roles.id = user_roles.role_id").
        Where("user_roles.user_id = ? AND roles.name = ?", u.ID, role).
        Count(&count)

    return count > 0
}
```

### 2. 权限中间件

```go
// 权限中间件
type PermissionMiddleware struct {
    http.Middleware
    Permission string
}

func (m *PermissionMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    user := request.Context["user"].(*Models.User)

    if !user.HasPermission(m.Permission) {
        return http.Response{
            StatusCode: 403,
            Body:       `{"error": "Insufficient permissions"}`,
        }
    }

    return next(request)
}

// 角色中间件
type RoleMiddleware struct {
    http.Middleware
    Role string
}

func (m *RoleMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    user := request.Context["user"].(*Models.User)

    if !user.HasRole(m.Role) {
        return http.Response{
            StatusCode: 403,
            Body:       `{"error": "Insufficient role"}`,
        }
    }

    return next(request)
}

// 在路由中使用
router.Group("/admin", func(group *routing.Router) {
    group.Use(&middleware.AuthMiddleware{})
    group.Use(&middleware.RoleMiddleware{Role: "admin"})

    group.Get("/users", &Controllers.AdminController{}, "Users")
    group.Post("/users", &Controllers.AdminController{}, "CreateUser").
        Use(&middleware.PermissionMiddleware{Permission: "users.create"})
})
```

### 3. 策略授权

```go
// 策略接口
type Policy interface {
    View(user *Models.User, model interface{}) bool
    Create(user *Models.User) bool
    Update(user *Models.User, model interface{}) bool
    Delete(user *Models.User, model interface{}) bool
}

// 文章策略
type PostPolicy struct{}

func (p *PostPolicy) View(user *Models.User, model interface{}) bool {
    post := model.(*Models.Post)
    return post.Status == "published" || post.UserID == user.ID || user.HasRole("admin")
}

func (p *PostPolicy) Create(user *Models.User) bool {
    return user.HasPermission("posts.create")
}

func (p *PostPolicy) Update(user *Models.User, model interface{}) bool {
    post := model.(*Models.Post)
    return post.UserID == user.ID || user.HasRole("admin")
}

func (p *PostPolicy) Delete(user *Models.User, model interface{}) bool {
    post := model.(*Models.Post)
    return post.UserID == user.ID || user.HasRole("admin")
}

// 策略门面
type Gate struct {
    policies map[string]Policy
}

func (g *Gate) Define(model string, policy Policy) {
    g.policies[model] = policy
}

func (g *Gate) Allows(user *Models.User, action string, model interface{}) bool {
    policy := g.policies[getModelName(model)]
    if policy == nil {
        return false
    }

    switch action {
    case "view":
        return policy.View(user, model)
    case "create":
        return policy.Create(user)
    case "update":
        return policy.Update(user, model)
    case "delete":
        return policy.Delete(user, model)
    default:
        return false
    }
}

// 在控制器中使用
func (c *PostController) Update(id string, request http.Request) http.Response {
    user := request.Context["user"].(*Models.User)

    post, err := c.postService.GetPost(postID)
    if err != nil {
        return c.JsonError("Post not found", 404)
    }

    if !gate.Allows(user, "update", post) {
        return c.JsonError("Unauthorized", 403)
    }

    // 更新文章...
}
```

## 🔄 密码管理

### 1. 密码哈希

```go
// 密码哈希
func HashPassword(password string) string {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    return string(hashedPassword)
}

// 验证密码
func CheckPassword(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

// 在用户模型中
func (u *User) SetPassword(password string) {
    u.Password = HashPassword(password)
}

func (u *User) CheckPassword(password string) bool {
    return CheckPassword(password, u.Password)
}
```

### 2. 密码重置

```go
// 密码重置令牌
type PasswordResetToken struct {
    database.Model
    Email     string    `json:"email" gorm:"size:255;not null"`
    Token     string    `json:"token" gorm:"size:255;unique;not null"`
    CreatedAt time.Time `json:"created_at"`
}

// 发送密码重置邮件
func (s *AuthService) SendPasswordResetLink(email string) error {
    var user Models.User
    err := s.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return err
    }

    // 生成重置令牌
    token := generateRandomToken(64)

    // 保存令牌
    resetToken := &PasswordResetToken{
        Email: email,
        Token: token,
    }
    s.db.Create(resetToken)

    // 发送邮件
    return s.mailer.SendPasswordReset(user.Email, token)
}

// 重置密码
func (s *AuthService) ResetPassword(token, password string) error {
    var resetToken PasswordResetToken
    err := s.db.Where("token = ?", token).First(&resetToken).Error
    if err != nil {
        return errors.New("invalid token")
    }

    // 检查令牌是否过期（24小时）
    if time.Since(resetToken.CreatedAt) > 24*time.Hour {
        return errors.New("token expired")
    }

    // 更新用户密码
    var user Models.User
    err = s.db.Where("email = ?", resetToken.Email).First(&user).Error
    if err != nil {
        return err
    }

    user.SetPassword(password)
    s.db.Save(&user)

    // 删除重置令牌
    s.db.Delete(&resetToken)

    return nil
}
```

## 🔐 多因素认证 (MFA)

### 1. TOTP 认证

```go
// 用户 MFA 设置
type UserMfa struct {
    database.Model
    UserID    uint      `json:"user_id" gorm:"primaryKey"`
    Secret    string    `json:"secret" gorm:"size:255;not null"`
    Enabled   bool      `json:"enabled" gorm:"default:false"`
    BackupCodes []string `json:"backup_codes" gorm:"type:text"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 生成 TOTP 密钥
func (s *AuthService) GenerateMfaSecret(userID uint) (string, string, error) {
    secret := generateRandomSecret(32)
    qrCode := generateQrCode(secret, userID)

    return secret, qrCode, nil
}

// 验证 TOTP 代码
func (s *AuthService) VerifyMfaCode(userID uint, code string) bool {
    var userMfa UserMfa
    err := s.db.Where("user_id = ? AND enabled = ?", userID, true).First(&userMfa).Error
    if err != nil {
        return false
    }

    return totp.Validate(code, userMfa.Secret)
}

// 启用 MFA
func (s *AuthService) EnableMfa(userID uint, secret, code string) error {
    if !s.VerifyMfaCode(userID, code) {
        return errors.New("invalid code")
    }

    backupCodes := generateBackupCodes(8)

    userMfa := &UserMfa{
        UserID:      userID,
        Secret:      secret,
        Enabled:     true,
        BackupCodes: backupCodes,
    }

    return s.db.Save(userMfa).Error
}
```

## 📊 认证统计

### 1. 登录历史

```go
// 登录历史模型
type LoginHistory struct {
    database.Model
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id" gorm:"not null"`
    IP        string    `json:"ip" gorm:"size:45"`
    UserAgent string    `json:"user_agent" gorm:"size:500"`
    Location  string    `json:"location" gorm:"size:255"`
    Success   bool      `json:"success" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at"`

    User User `json:"user" gorm:"foreignKey:UserID"`
}

// 记录登录历史
func (s *AuthService) RecordLogin(userID uint, ip, userAgent string, success bool) {
    location := s.getLocationByIP(ip)

    history := &LoginHistory{
        UserID:    userID,
        IP:        ip,
        UserAgent: userAgent,
        Location:  location,
        Success:   success,
    }

    s.db.Create(history)
}

// 获取用户登录历史
func (s *AuthService) GetLoginHistory(userID uint, limit int) ([]*LoginHistory, error) {
    var history []*LoginHistory
    err := s.db.Where("user_id = ?", userID).
        Order("created_at desc").
        Limit(limit).
        Find(&history).Error

    return history, err
}
```

## 🛡️ 安全最佳实践

### 1. 密码策略

```go
// 密码验证规则
func ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }

    if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
        return errors.New("password must contain at least one uppercase letter")
    }

    if !regexp.MustCompile(`[a-z]`).MatchString(password) {
        return errors.New("password must contain at least one lowercase letter")
    }

    if !regexp.MustCompile(`[0-9]`).MatchString(password) {
        return errors.New("password must contain at least one number")
    }

    if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
        return errors.New("password must contain at least one special character")
    }

    return nil
}
```

### 2. 登录限制

```go
// 登录尝试限制
type LoginAttempt struct {
    database.Model
    ID        uint      `json:"id" gorm:"primaryKey"`
    Email     string    `json:"email" gorm:"size:255;not null"`
    IP        string    `json:"ip" gorm:"size:45;not null"`
    Success   bool      `json:"success" gorm:"default:false"`
    CreatedAt time.Time `json:"created_at"`
}

// 检查登录限制
func (s *AuthService) IsLoginBlocked(email, ip string) bool {
    var count int64
    s.db.Model(&LoginAttempt{}).
        Where("email = ? AND ip = ? AND success = ? AND created_at > ?",
            email, ip, false, time.Now().Add(-15*time.Minute)).
        Count(&count)

    return count >= 5
}

// 记录登录尝试
func (s *AuthService) RecordLoginAttempt(email, ip string, success bool) {
    attempt := &LoginAttempt{
        Email:   email,
        IP:      ip,
        Success: success,
    }

    s.db.Create(attempt)
}
```

### 3. 会话管理

```go
// 会话配置
config.Set("session.secure", true)
config.Set("session.http_only", true)
config.Set("session.same_site", "strict")

// 强制登出
func (s *AuthService) ForceLogout(userID uint) error {
    // 删除所有会话
    return s.db.Where("user_id = ?", userID).Delete(&Session{}).Error
}

// 会话超时
func (m *AuthMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    session := request.Session()
    lastActivity := session.Get("last_activity")

    if lastActivity != nil {
        lastTime := lastActivity.(time.Time)
        if time.Since(lastTime) > 30*time.Minute {
            session.Destroy()
            return http.Response{
                StatusCode: 401,
                Body:       `{"error": "Session expired"}`,
            }
        }
    }

    session.Put("last_activity", time.Now())
    return next(request)
}
```

## 📚 总结

Laravel-Go Framework 的认证授权系统提供了：

1. **多种认证方式**: JWT、Session、API Token
2. **灵活的权限控制**: RBAC、策略授权
3. **安全特性**: 密码哈希、MFA、登录限制
4. **会话管理**: 安全配置、超时控制
5. **审计功能**: 登录历史、操作日志

通过合理配置和使用这些功能，可以构建安全可靠的用户认证和授权系统。
