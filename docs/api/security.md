# 安全系统 API 参考

## 📋 概述

Laravel-Go Framework 的安全系统提供了全面的安全防护功能，包括 CSRF 保护、XSS 防护、SQL 注入防护、密码哈希、加密解密、安全头部设置等。安全系统旨在保护应用程序免受各种安全威胁，确保数据和用户安全。

## 🏗️ 核心概念

### 安全中间件 (Security Middleware)

- 自动应用安全防护措施
- 处理安全头部和策略
- 监控和记录安全事件

### 加密服务 (Encryption Service)

- 数据加密和解密
- 安全密钥管理
- 哈希算法支持

### 审计日志 (Audit Log)

- 记录安全相关事件
- 用户行为追踪
- 安全事件分析

## 🔧 基础用法

### 1. 基本安全配置

```go
// 配置安全选项
security := security.NewSecurity()

// 设置安全头部
security.SetHeaders(map[string]string{
    "X-Frame-Options":           "DENY",
    "X-Content-Type-Options":    "nosniff",
    "X-XSS-Protection":          "1; mode=block",
    "Strict-Transport-Security": "max-age=31536000; includeSubDomains",
    "Content-Security-Policy":   "default-src 'self'",
})

// 启用 CSRF 保护
security.EnableCSRF(true)

// 启用 XSS 防护
security.EnableXSSProtection(true)

// 设置密码策略
security.SetPasswordPolicy(security.PasswordPolicy{
    MinLength:      8,
    RequireUppercase: true,
    RequireLowercase: true,
    RequireNumbers:   true,
    RequireSymbols:   true,
})
```

### 2. 在中间件中使用

```go
// app/Http/Middleware/SecurityMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "laravel-go/framework/security"
)

type SecurityMiddleware struct {
    http.Middleware
    security *security.Security
}

func (m *SecurityMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 应用安全头部
    response := next(request)

    // 添加安全头部
    for key, value := range m.security.GetHeaders() {
        response.Headers[key] = value
    }

    // CSRF 保护
    if m.security.IsCSRFEnabled() && request.Method != "GET" {
        if !m.security.ValidateCSRFToken(request) {
            return http.Response{
                StatusCode: 403,
                Body:       `{"error": "CSRF token validation failed"}`,
                Headers: map[string]string{
                    "Content-Type": "application/json",
                },
            }
        }
    }

    // XSS 防护
    if m.security.IsXSSProtectionEnabled() {
        response.Body = m.security.SanitizeHTML(response.Body)
    }

    return response
}
```

### 3. 密码哈希和验证

```go
// 密码哈希
password := "mySecurePassword123!"
hashedPassword, err := security.HashPassword(password)
if err != nil {
    log.Fatal(err)
}

// 密码验证
isValid := security.CheckPassword(password, hashedPassword)
if isValid {
    fmt.Println("Password is valid")
} else {
    fmt.Println("Password is invalid")
}

// 检查密码强度
strength := security.CheckPasswordStrength(password)
fmt.Printf("Password strength: %s\n", strength)
```

## 📚 API 参考

### Security 接口

```go
type Security interface {
    SetHeaders(headers map[string]string)
    GetHeaders() map[string]string
    AddHeader(key, value string)
    RemoveHeader(key string)

    EnableCSRF(enabled bool)
    IsCSRFEnabled() bool
    GenerateCSRFToken() string
    ValidateCSRFToken(request Request) bool

    EnableXSSProtection(enabled bool)
    IsXSSProtectionEnabled() bool
    SanitizeHTML(html string) string

    SetPasswordPolicy(policy PasswordPolicy)
    GetPasswordPolicy() PasswordPolicy
    CheckPasswordStrength(password string) string

    Encrypt(data []byte) ([]byte, error)
    Decrypt(data []byte) ([]byte, error)

    HashPassword(password string) (string, error)
    CheckPassword(password, hash string) bool

    GenerateToken(length int) string
    ValidateToken(token string) bool

    LogSecurityEvent(event SecurityEvent)
    GetSecurityEvents() []SecurityEvent
}
```

#### 方法说明

- `SetHeaders(headers)`: 设置安全头部
- `GetHeaders()`: 获取安全头部
- `AddHeader(key, value)`: 添加安全头部
- `RemoveHeader(key)`: 移除安全头部
- `EnableCSRF(enabled)`: 启用/禁用 CSRF 保护
- `IsCSRFEnabled()`: 检查 CSRF 是否启用
- `GenerateCSRFToken()`: 生成 CSRF 令牌
- `ValidateCSRFToken(request)`: 验证 CSRF 令牌
- `EnableXSSProtection(enabled)`: 启用/禁用 XSS 防护
- `IsXSSProtectionEnabled()`: 检查 XSS 防护是否启用
- `SanitizeHTML(html)`: 清理 HTML
- `SetPasswordPolicy(policy)`: 设置密码策略
- `GetPasswordPolicy()`: 获取密码策略
- `CheckPasswordStrength(password)`: 检查密码强度
- `Encrypt(data)`: 加密数据
- `Decrypt(data)`: 解密数据
- `HashPassword(password)`: 哈希密码
- `CheckPassword(password, hash)`: 验证密码
- `GenerateToken(length)`: 生成安全令牌
- `ValidateToken(token)`: 验证令牌
- `LogSecurityEvent(event)`: 记录安全事件
- `GetSecurityEvents()`: 获取安全事件

### PasswordPolicy 结构体

```go
type PasswordPolicy struct {
    MinLength        int  `json:"min_length"`
    RequireUppercase bool `json:"require_uppercase"`
    RequireLowercase bool `json:"require_lowercase"`
    RequireNumbers   bool `json:"require_numbers"`
    RequireSymbols   bool `json:"require_symbols"`
    MaxLength        int  `json:"max_length"`
    PreventCommon    bool `json:"prevent_common"`
}
```

#### 字段说明

- `MinLength`: 最小长度
- `RequireUppercase`: 要求大写字母
- `RequireLowercase`: 要求小写字母
- `RequireNumbers`: 要求数字
- `RequireSymbols`: 要求特殊字符
- `MaxLength`: 最大长度
- `PreventCommon`: 防止常见密码

### SecurityEvent 结构体

```go
type SecurityEvent struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    UserID    string                 `json:"user_id"`
    IP        string                 `json:"ip"`
    UserAgent string                 `json:"user_agent"`
    Data      map[string]interface{} `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
}
```

#### 字段说明

- `ID`: 事件 ID
- `Type`: 事件类型
- `Level`: 事件级别
- `Message`: 事件消息
- `UserID`: 用户 ID
- `IP`: IP 地址
- `UserAgent`: 用户代理
- `Data`: 额外数据
- `Timestamp`: 时间戳

## 🎯 高级功能

### 1. CSRF 保护

```go
// 生成 CSRF 令牌
func (c *AuthController) ShowLoginForm(request http.Request) http.Response {
    csrfToken := c.security.GenerateCSRFToken()

    data := map[string]interface{}{
        "csrf_token": csrfToken,
    }

    return c.View("auth/login.html", data)
}

// 验证 CSRF 令牌
func (c *AuthController) Login(request http.Request) http.Response {
    // 验证 CSRF 令牌
    if !c.security.ValidateCSRFToken(request) {
        c.security.LogSecurityEvent(security.SecurityEvent{
            Type:    "csrf_attack",
            Level:   "high",
            Message: "CSRF token validation failed",
            IP:      request.IP,
            Data: map[string]interface{}{
                "method": request.Method,
                "path":   request.Path,
            },
        })

        return c.JsonError("CSRF token validation failed", 403)
    }

    // 处理登录逻辑
    // ...
}
```

### 2. XSS 防护

```go
// HTML 清理
func (c *PostController) Store(request http.Request) http.Response {
    content := request.Body["content"].(string)

    // 清理 HTML 内容
    sanitizedContent := c.security.SanitizeHTML(content)

    // 创建文章
    post := &Models.Post{
        Title:   request.Body["title"].(string),
        Content: sanitizedContent,
    }

    // 保存到数据库
    // ...
}

// 自定义清理规则
func (c *PostController) StoreWithCustomSanitization(request http.Request) http.Response {
    content := request.Body["content"].(string)

    // 自定义清理规则
    sanitizer := security.NewHTMLSanitizer()
    sanitizer.AllowTags("p", "br", "strong", "em", "a")
    sanitizer.AllowAttributes("href", "target")

    sanitizedContent := sanitizer.Sanitize(content)

    // 保存文章
    // ...
}
```

### 3. 密码安全

```go
// 密码策略验证
func (c *UserController) ChangePassword(request http.Request) http.Response {
    newPassword := request.Body["new_password"].(string)

    // 检查密码强度
    strength := c.security.CheckPasswordStrength(newPassword)
    if strength == "weak" {
        return c.JsonError("Password is too weak", 422)
    }

    // 验证密码策略
    policy := c.security.GetPasswordPolicy()
    if !c.security.ValidatePasswordPolicy(newPassword, policy) {
        return c.JsonError("Password does not meet requirements", 422)
    }

    // 哈希密码
    hashedPassword, err := c.security.HashPassword(newPassword)
    if err != nil {
        return c.JsonError("Failed to hash password", 500)
    }

    // 更新用户密码
    user := request.Context["user"].(*Models.User)
    user.Password = hashedPassword

    // 保存到数据库
    // ...
}
```

### 4. 数据加密

```go
// 加密敏感数据
func (c *UserController) Store(request http.Request) http.Response {
    // 加密敏感信息
    encryptedSSN, err := c.security.Encrypt([]byte(request.Body["ssn"].(string)))
    if err != nil {
        return c.JsonError("Failed to encrypt data", 500)
    }

    // 创建用户
    user := &Models.User{
        Name:     request.Body["name"].(string),
        Email:    request.Body["email"].(string),
        SSN:      string(encryptedSSN), // 存储加密后的数据
    }

    // 保存到数据库
    // ...
}

// 解密数据
func (c *UserController) Show(id string, request http.Request) http.Response {
    user, err := c.userService.GetUser(id)
    if err != nil {
        return c.JsonError("User not found", 404)
    }

    // 解密敏感信息
    decryptedSSN, err := c.security.Decrypt([]byte(user.SSN))
    if err != nil {
        return c.JsonError("Failed to decrypt data", 500)
    }

    // 返回用户信息（不包含敏感数据）
    userData := map[string]interface{}{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
        // 不返回 SSN
    }

    return c.Json(userData)
}
```

### 5. 安全审计

```go
// 记录安全事件
func (c *AuthController) Login(request http.Request) http.Response {
    email := request.Body["email"].(string)

    // 尝试登录
    user, err := c.authService.Authenticate(email, request.Body["password"].(string))
    if err != nil {
        // 记录登录失败事件
        c.security.LogSecurityEvent(security.SecurityEvent{
            Type:    "login_failed",
            Level:   "medium",
            Message: "Failed login attempt",
            IP:      request.IP,
            Data: map[string]interface{}{
                "email": email,
                "error": err.Error(),
            },
        })

        return c.JsonError("Invalid credentials", 401)
    }

    // 记录登录成功事件
    c.security.LogSecurityEvent(security.SecurityEvent{
        Type:    "login_success",
        Level:   "low",
        Message: "User logged in successfully",
        UserID:  fmt.Sprintf("%d", user.ID),
        IP:      request.IP,
        Data: map[string]interface{}{
            "email": user.Email,
        },
    })

    return c.Json(map[string]interface{}{
        "token": generateToken(user),
        "user":  user,
    })
}
```

## 🔧 配置选项

### 安全系统配置

```go
// config/security.go
package config

type SecurityConfig struct {
    // CSRF 配置
    CSRF CSRFConfig `json:"csrf"`

    // XSS 配置
    XSS XSSConfig `json:"xss"`

    // 密码策略
    Password PasswordConfig `json:"password"`

    // 加密配置
    Encryption EncryptionConfig `json:"encryption"`

    // 安全头部
    Headers HeadersConfig `json:"headers"`

    // 审计配置
    Audit AuditConfig `json:"audit"`

    // 速率限制
    RateLimit RateLimitConfig `json:"rate_limit"`
}

type CSRFConfig struct {
    Enabled bool `json:"enabled"`
    TokenLength int `json:"token_length"`
    ExpireTime time.Duration `json:"expire_time"`
}

type XSSConfig struct {
    Enabled bool `json:"enabled"`
    Strict bool `json:"strict"`
    AllowedTags []string `json:"allowed_tags"`
    AllowedAttributes []string `json:"allowed_attributes"`
}

type PasswordConfig struct {
    MinLength int `json:"min_length"`
    RequireUppercase bool `json:"require_uppercase"`
    RequireLowercase bool `json:"require_lowercase"`
    RequireNumbers bool `json:"require_numbers"`
    RequireSymbols bool `json:"require_symbols"`
    MaxLength int `json:"max_length"`
    PreventCommon bool `json:"prevent_common"`
}

type EncryptionConfig struct {
    Key string `json:"key"`
    Algorithm string `json:"algorithm"`
    KeySize int `json:"key_size"`
}

type HeadersConfig struct {
    XFrameOptions string `json:"x_frame_options"`
    XContentTypeOptions string `json:"x_content_type_options"`
    XSSProtection string `json:"xss_protection"`
    HSTS string `json:"hsts"`
    CSP string `json:"csp"`
}

type AuditConfig struct {
    Enabled bool `json:"enabled"`
    LogLevel string `json:"log_level"`
    Storage string `json:"storage"`
    Retention time.Duration `json:"retention"`
}

type RateLimitConfig struct {
    Enabled bool `json:"enabled"`
    Requests int `json:"requests"`
    Window time.Duration `json:"window"`
    BlockDuration time.Duration `json:"block_duration"`
}
```

### 配置示例

```go
// config/security.go
func GetSecurityConfig() *SecurityConfig {
    return &SecurityConfig{
        CSRF: CSRFConfig{
            Enabled:     true,
            TokenLength: 32,
            ExpireTime:  time.Hour,
        },
        XSS: XSSConfig{
            Enabled: true,
            Strict:  false,
            AllowedTags: []string{"p", "br", "strong", "em", "a", "ul", "ol", "li"},
            AllowedAttributes: []string{"href", "target"},
        },
        Password: PasswordConfig{
            MinLength:        8,
            RequireUppercase: true,
            RequireLowercase: true,
            RequireNumbers:   true,
            RequireSymbols:   true,
            MaxLength:        128,
            PreventCommon:    true,
        },
        Encryption: EncryptionConfig{
            Key:        "your-secret-key-here",
            Algorithm:  "AES-256-GCM",
            KeySize:    256,
        },
        Headers: HeadersConfig{
            XFrameOptions:        "DENY",
            XContentTypeOptions:  "nosniff",
            XSSProtection:        "1; mode=block",
            HSTS:                 "max-age=31536000; includeSubDomains",
            CSP:                  "default-src 'self'; script-src 'self' 'unsafe-inline'",
        },
        Audit: AuditConfig{
            Enabled:   true,
            LogLevel:  "medium",
            Storage:   "database",
            Retention: time.Hour * 24 * 30, // 30 days
        },
        RateLimit: RateLimitConfig{
            Enabled:       true,
            Requests:      100,
            Window:        time.Minute,
            BlockDuration: time.Minute * 5,
        },
    }
}
```

## 🚀 性能优化

### 1. 安全缓存

```go
// 缓存安全检查结果
type CachedSecurity struct {
    security.Security
    cache cache.Cache
}

func (s *CachedSecurity) ValidateCSRFToken(request Request) bool {
    token := request.Headers["X-CSRF-Token"]
    cacheKey := fmt.Sprintf("csrf:%s", token)

    if cached, exists := s.cache.Get(cacheKey); exists {
        return cached.(bool)
    }

    result := s.Security.ValidateCSRFToken(request)
    s.cache.Set(cacheKey, result, time.Minute*5)

    return result
}
```

### 2. 批量安全检查

```go
// 批量处理安全检查
func (s *Security) BatchSecurityCheck(requests []Request) []SecurityResult {
    results := make([]SecurityResult, len(requests))

    var wg sync.WaitGroup
    for i, request := range requests {
        wg.Add(1)
        go func(index int, req Request) {
            defer wg.Done()

            results[index] = SecurityResult{
                RequestID: req.ID,
                CSRFValid: s.ValidateCSRFToken(req),
                XSSRisk:   s.DetectXSSRisk(req),
                IPRisk:    s.DetectIPRisk(req.IP),
            }
        }(i, request)
    }

    wg.Wait()
    return results
}
```

### 3. 安全事件聚合

```go
// 聚合安全事件
type SecurityEventAggregator struct {
    events []SecurityEvent
    mutex  sync.RWMutex
}

func (a *SecurityEventAggregator) AddEvent(event SecurityEvent) {
    a.mutex.Lock()
    defer a.mutex.Unlock()

    a.events = append(a.events, event)

    // 如果事件数量超过阈值，进行聚合
    if len(a.events) > 1000 {
        a.aggregateEvents()
    }
}

func (a *SecurityEventAggregator) aggregateEvents() {
    // 按类型和时间窗口聚合事件
    aggregated := make(map[string]int)

    for _, event := range a.events {
        key := fmt.Sprintf("%s:%s", event.Type, event.IP)
        aggregated[key]++
    }

    // 保存聚合结果
    // ...

    // 清空事件列表
    a.events = a.events[:0]
}
```

## 🧪 测试

### 1. 安全功能测试

```go
// tests/security_test.go
package tests

import (
    "testing"
    "laravel-go/framework/security"
)

func TestCSRFProtection(t *testing.T) {
    security := security.NewSecurity()
    security.EnableCSRF(true)

    // 生成令牌
    token := security.GenerateCSRFToken()
    if token == "" {
        t.Error("CSRF token should not be empty")
    }

    // 创建请求
    request := http.Request{
        Method: "POST",
        Headers: map[string]string{
            "X-CSRF-Token": token,
        },
    }

    // 验证令牌
    if !security.ValidateCSRFToken(request) {
        t.Error("CSRF token should be valid")
    }

    // 测试无效令牌
    request.Headers["X-CSRF-Token"] = "invalid-token"
    if security.ValidateCSRFToken(request) {
        t.Error("Invalid CSRF token should be rejected")
    }
}

func TestPasswordHashing(t *testing.T) {
    security := security.NewSecurity()

    password := "testPassword123!"

    // 哈希密码
    hash, err := security.HashPassword(password)
    if err != nil {
        t.Fatal(err)
    }

    // 验证密码
    if !security.CheckPassword(password, hash) {
        t.Error("Password verification should succeed")
    }

    // 测试错误密码
    if security.CheckPassword("wrongPassword", hash) {
        t.Error("Wrong password should be rejected")
    }
}

func TestXSSProtection(t *testing.T) {
    security := security.NewSecurity()
    security.EnableXSSProtection(true)

    maliciousHTML := `<script>alert('xss')</script><p>Hello World</p>`
    sanitized := security.SanitizeHTML(maliciousHTML)

    if strings.Contains(sanitized, "<script>") {
        t.Error("Script tags should be removed")
    }

    if !strings.Contains(sanitized, "<p>Hello World</p>") {
        t.Error("Safe HTML should be preserved")
    }
}
```

### 2. 密码策略测试

```go
func TestPasswordPolicy(t *testing.T) {
    security := security.NewSecurity()

    policy := security.PasswordPolicy{
        MinLength:        8,
        RequireUppercase: true,
        RequireLowercase: true,
        RequireNumbers:   true,
        RequireSymbols:   true,
    }

    security.SetPasswordPolicy(policy)

    // 测试强密码
    strongPassword := "StrongPass123!"
    strength := security.CheckPasswordStrength(strongPassword)
    if strength != "strong" {
        t.Errorf("Expected strong password, got %s", strength)
    }

    // 测试弱密码
    weakPassword := "weak"
    strength = security.CheckPasswordStrength(weakPassword)
    if strength != "weak" {
        t.Errorf("Expected weak password, got %s", strength)
    }
}
```

## 🔍 调试和监控

### 1. 安全监控

```go
type SecurityMonitor struct {
    security.Security
    metrics metrics.Collector
}

func (m *SecurityMonitor) ValidateCSRFToken(request Request) bool {
    result := m.Security.ValidateCSRFToken(request)

    // 记录指标
    if result {
        m.metrics.Increment("security.csrf.valid")
    } else {
        m.metrics.Increment("security.csrf.invalid")
    }

    return result
}

func (m *SecurityMonitor) LogSecurityEvent(event SecurityEvent) {
    // 记录安全事件指标
    m.metrics.Increment("security.events", map[string]string{
        "type":  event.Type,
        "level": event.Level,
    })

    // 调用原始方法
    m.Security.LogSecurityEvent(event)
}
```

### 2. 安全报告

```go
type SecurityReporter struct {
    security.Security
}

func (r *SecurityReporter) GenerateReport(startTime, endTime time.Time) SecurityReport {
    events := r.GetSecurityEvents()

    // 过滤时间范围内的事件
    filteredEvents := make([]SecurityEvent, 0)
    for _, event := range events {
        if event.Timestamp.After(startTime) && event.Timestamp.Before(endTime) {
            filteredEvents = append(filteredEvents, event)
        }
    }

    // 生成报告
    report := SecurityReport{
        Period:     fmt.Sprintf("%s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")),
        TotalEvents: len(filteredEvents),
        EventsByType: make(map[string]int),
        EventsByLevel: make(map[string]int),
        TopIPs:      make(map[string]int),
    }

    for _, event := range filteredEvents {
        report.EventsByType[event.Type]++
        report.EventsByLevel[event.Level]++
        report.TopIPs[event.IP]++
    }

    return report
}
```

## 📝 最佳实践

### 1. 安全头部配置

```go
// 配置完整的安全头部
func configureSecurityHeaders(security *security.Security) {
    headers := map[string]string{
        // 防止点击劫持
        "X-Frame-Options": "DENY",

        // 防止 MIME 类型嗅探
        "X-Content-Type-Options": "nosniff",

        // XSS 防护
        "X-XSS-Protection": "1; mode=block",

        // 强制 HTTPS
        "Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",

        // 内容安全策略
        "Content-Security-Policy": "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self' https:; frame-ancestors 'none';",

        // 引用策略
        "Referrer-Policy": "strict-origin-when-cross-origin",

        // 权限策略
        "Permissions-Policy": "geolocation=(), microphone=(), camera=()",
    }

    security.SetHeaders(headers)
}
```

### 2. 密码安全

```go
// 实施强密码策略
func enforcePasswordPolicy(security *security.Security) {
    policy := security.PasswordPolicy{
        MinLength:        12,
        RequireUppercase: true,
        RequireLowercase: true,
        RequireNumbers:   true,
        RequireSymbols:   true,
        MaxLength:        128,
        PreventCommon:    true,
    }

    security.SetPasswordPolicy(policy)
}

// 密码验证函数
func validatePassword(password string, security *security.Security) error {
    policy := security.GetPasswordPolicy()

    if len(password) < policy.MinLength {
        return errors.New("password too short")
    }

    if policy.RequireUppercase && !hasUppercase(password) {
        return errors.New("password must contain uppercase letter")
    }

    if policy.RequireLowercase && !hasLowercase(password) {
        return errors.New("password must contain lowercase letter")
    }

    if policy.RequireNumbers && !hasNumbers(password) {
        return errors.New("password must contain number")
    }

    if policy.RequireSymbols && !hasSymbols(password) {
        return errors.New("password must contain symbol")
    }

    if policy.PreventCommon && isCommonPassword(password) {
        return errors.New("password is too common")
    }

    return nil
}
```

### 3. 输入验证

```go
// 安全的输入验证
func validateInput(input map[string]interface{}, security *security.Security) (map[string]interface{}, error) {
    validated := make(map[string]interface{})

    for key, value := range input {
        switch v := value.(type) {
        case string:
            // 清理字符串输入
            cleaned := security.SanitizeHTML(v)
            validated[key] = cleaned

        case []string:
            // 清理字符串数组
            cleaned := make([]string, len(v))
            for i, s := range v {
                cleaned[i] = security.SanitizeHTML(s)
            }
            validated[key] = cleaned

        default:
            validated[key] = v
        }
    }

    return validated, nil
}
```

### 4. 安全日志记录

```go
// 记录安全事件
func logSecurityEvent(security *security.Security, eventType, message, userID, ip string, data map[string]interface{}) {
    event := security.SecurityEvent{
        Type:      eventType,
        Level:     "medium",
        Message:   message,
        UserID:    userID,
        IP:        ip,
        Data:      data,
        Timestamp: time.Now(),
    }

    security.LogSecurityEvent(event)
}

// 在控制器中使用
func (c *UserController) Login(request http.Request) http.Response {
    email := request.Body["email"].(string)
    ip := request.IP

    // 尝试登录
    user, err := c.authService.Authenticate(email, request.Body["password"].(string))
    if err != nil {
        // 记录登录失败
        logSecurityEvent(c.security, "login_failed", "Failed login attempt", "", ip, map[string]interface{}{
            "email": email,
            "error": err.Error(),
        })

        return c.JsonError("Invalid credentials", 401)
    }

    // 记录登录成功
    logSecurityEvent(c.security, "login_success", "User logged in", fmt.Sprintf("%d", user.ID), ip, map[string]interface{}{
        "email": user.Email,
    })

    return c.Json(map[string]interface{}{
        "token": generateToken(user),
        "user":  user,
    })
}
```

## 🚀 总结

安全系统是 Laravel-Go Framework 中重要的功能之一，它提供了：

1. **全面的安全防护**: CSRF、XSS、SQL 注入等防护
2. **密码安全**: 强密码策略和哈希算法
3. **数据加密**: 敏感数据加密和解密
4. **安全审计**: 完整的安全事件记录和分析
5. **性能优化**: 缓存和批量处理优化
6. **最佳实践**: 遵循安全开发的最佳实践

通过合理使用安全系统，可以有效地保护应用程序免受各种安全威胁，确保数据和用户安全。
