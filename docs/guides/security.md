# 安全实践指南

## 📖 概述

Laravel-Go Framework 提供了全面的安全功能，包括认证授权、数据加密、CSRF 保护、XSS 防护、SQL 注入防护等，帮助构建安全可靠的应用程序。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [安全系统 API 参考](../api/security.md)

## 🚀 快速开始

### 1. 基本安全配置

```go
// 安全配置
config.Set("security.csrf.enabled", true)
config.Set("security.csrf.token_name", "_token")
config.Set("security.csrf.expire", 3600)

config.Set("security.xss.enabled", true)
config.Set("security.xss.whitelist", []string{"p", "br", "strong"})

config.Set("security.rate_limit.enabled", true)
config.Set("security.rate_limit.max_attempts", 60)
config.Set("security.rate_limit.decay_minutes", 1)
```

### 2. CSRF 保护

```go
// CSRF 中间件
type CSRFMiddleware struct{}

func (m *CSRFMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 跳过 GET 请求
    if request.Method == "GET" {
        return next(request)
    }

    // 验证 CSRF 令牌
    token := request.Headers["X-CSRF-TOKEN"]
    if !security.ValidateCSRFToken(token) {
        return http.Response{
            StatusCode: 419,
            Body:       `{"error": "CSRF token mismatch"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }

    return next(request)
}

// 在控制器中生成 CSRF 令牌
func (c *UserController) CreateForm(request http.Request) http.Response {
    token := security.GenerateCSRFToken()

    return c.Json(map[string]interface{}{
        "csrf_token": token,
        "form_data":  map[string]interface{}{},
    })
}
```

### 3. XSS 防护

```go
// XSS 防护中间件
type XSSMiddleware struct{}

func (m *XSSMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 清理请求数据
    cleanedBody := security.SanitizeHTML(request.Body)
    request.Body = cleanedBody

    response := next(request)

    // 清理响应数据
    if response.Headers["Content-Type"] == "text/html" {
        response.Body = security.SanitizeHTML(response.Body)
    }

    return response
}

// 手动清理 HTML
func (c *PostController) Store(request http.Request) http.Response {
    content := request.Body["content"].(string)

    // 清理 HTML 内容
    cleanContent := security.SanitizeHTML(content)

    post := &Models.Post{
        Title:   request.Body["title"].(string),
        Content: cleanContent,
    }

    // 保存文章
    err := c.postService.CreatePost(post)
    if err != nil {
        return c.JsonError("Failed to create post", 500)
    }

    return c.Json(post).Status(201)
}
```

### 4. SQL 注入防护

```go
// 使用参数化查询
func (s *UserService) GetUserByEmail(email string) (*Models.User, error) {
    var user Models.User

    // 使用参数化查询防止 SQL 注入
    err := s.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }

    return &user, nil
}

// 使用查询构建器
func (s *UserService) SearchUsers(query string, filters map[string]interface{}) ([]*Models.User, error) {
    db := s.db.Model(&Models.User{})

    // 安全的搜索查询
    if query != "" {
        db = db.Where("name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%")
    }

    // 安全的过滤条件
    for field, value := range filters {
        if security.IsAllowedField(field) {
            db = db.Where(field+" = ?", value)
        }
    }

    var users []*Models.User
    err := db.Find(&users).Error

    return users, err
}
```

### 5. 输入验证

```go
// 严格的输入验证
func (c *UserController) Store(request http.Request) http.Response {
    v := validator.New()

    // 验证必填字段
    v.Required("name", request.Body["name"], "Name is required")
    v.Required("email", request.Body["email"], "Email is required")
    v.Required("password", request.Body["password"], "Password is required")

    // 验证格式
    v.Email("email", request.Body["email"], "Invalid email format")
    v.MinLength("password", request.Body["password"], 8, "Password too short")
    v.MaxLength("name", request.Body["name"], 255, "Name too long")

    // 验证内容
    v.Regex("name", request.Body["name"], `^[a-zA-Z\s]+$`, "Name contains invalid characters")

    if !v.Passes() {
        return c.JsonError("Validation failed", 422).WithErrors(v.Errors())
    }

    // 处理业务逻辑
    user, err := c.userService.CreateUser(request.Body)
    if err != nil {
        return c.JsonError("Failed to create user", 500)
    }

    return c.Json(user).Status(201)
}
```

### 6. 文件上传安全

```go
// 安全的文件上传
func (c *FileController) Upload(request http.Request) http.Response {
    file := request.Files["file"]

    // 验证文件类型
    if !security.IsAllowedFileType(file.Filename) {
        return c.JsonError("File type not allowed", 400)
    }

    // 验证文件大小
    if file.Size > 10*1024*1024 { // 10MB
        return c.JsonError("File too large", 400)
    }

    // 生成安全的文件名
    safeName := security.GenerateSafeFileName(file.Filename)

    // 保存文件到安全位置
    uploadPath := "storage/uploads/" + safeName
    if err := file.Save(uploadPath); err != nil {
        return c.JsonError("Failed to save file", 500)
    }

    return c.Json(map[string]string{
        "filename": safeName,
        "path":     uploadPath,
    })
}

// 文件类型验证
func IsAllowedFileType(filename string) bool {
    allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx"}
    ext := strings.ToLower(filepath.Ext(filename))

    for _, allowedType := range allowedTypes {
        if ext == allowedType {
            return true
        }
    }

    return false
}
```

### 7. 密码安全

```go
// 密码哈希
func (u *User) SetPassword(password string) error {
    // 验证密码强度
    if err := security.ValidatePasswordStrength(password); err != nil {
        return err
    }

    // 生成密码哈希
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    u.Password = string(hashedPassword)
    return nil
}

// 密码验证
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}

// 密码强度验证
func ValidatePasswordStrength(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }

    if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
        return errors.New("password must contain uppercase letter")
    }

    if !regexp.MustCompile(`[a-z]`).MatchString(password) {
        return errors.New("password must contain lowercase letter")
    }

    if !regexp.MustCompile(`[0-9]`).MatchString(password) {
        return errors.New("password must contain number")
    }

    if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
        return errors.New("password must contain special character")
    }

    return nil
}
```

### 8. 会话安全

```go
// 安全的会话配置
func ConfigureSecureSession() {
    config.Set("session.driver", "redis")
    config.Set("session.lifetime", 120) // 2小时
    config.Set("session.expire_on_close", true)
    config.Set("session.secure", true) // HTTPS only
    config.Set("session.http_only", true)
    config.Set("session.same_site", "strict")
}

// 会话劫持防护
type SessionSecurityMiddleware struct{}

func (m *SessionSecurityMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    session := request.Session()

    // 检查会话指纹
    currentFingerprint := security.GenerateSessionFingerprint(request)
    storedFingerprint := session.Get("fingerprint")

    if storedFingerprint != nil && storedFingerprint != currentFingerprint {
        // 会话可能被劫持，重新生成会话
        session.Regenerate()
        session.Put("fingerprint", currentFingerprint)
    }

    // 更新最后活动时间
    session.Put("last_activity", time.Now())

    return next(request)
}

// 生成会话指纹
func GenerateSessionFingerprint(request http.Request) string {
    userAgent := request.Headers["User-Agent"]
    ip := request.IP

    data := userAgent + "|" + ip
    hash := sha256.Sum256([]byte(data))

    return hex.EncodeToString(hash[:])
}
```

### 9. 日志安全

```go
// 安全日志记录
type SecurityLogger struct {
    logger *log.Logger
}

func (l *SecurityLogger) LogSecurityEvent(event string, details map[string]interface{}) {
    // 移除敏感信息
    sanitizedDetails := l.sanitizeDetails(details)

    l.logger.Info("Security Event", map[string]interface{}{
        "event":   event,
        "details": sanitizedDetails,
        "time":    time.Now(),
        "ip":      details["ip"],
    })
}

func (l *SecurityLogger) sanitizeDetails(details map[string]interface{}) map[string]interface{} {
    sanitized := make(map[string]interface{})

    for key, value := range details {
        if l.isSensitiveField(key) {
            sanitized[key] = "[REDACTED]"
        } else {
            sanitized[key] = value
        }
    }

    return sanitized
}

func (l *SecurityLogger) isSensitiveField(field string) bool {
    sensitiveFields := []string{"password", "token", "secret", "key", "credit_card"}

    for _, sensitive := range sensitiveFields {
        if strings.Contains(strings.ToLower(field), sensitive) {
            return true
        }
    }

    return false
}
```

## 📚 总结

Laravel-Go Framework 的安全系统提供了：

1. **CSRF 保护**: 防止跨站请求伪造攻击
2. **XSS 防护**: 防止跨站脚本攻击
3. **SQL 注入防护**: 使用参数化查询
4. **输入验证**: 严格的数据验证
5. **文件上传安全**: 安全的文件处理
6. **密码安全**: 强密码策略和哈希
7. **会话安全**: 会话劫持防护
8. **安全日志**: 敏感信息保护

通过合理使用这些安全功能，可以构建安全可靠的应用程序。
