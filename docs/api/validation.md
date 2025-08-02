# 验证系统 API 参考

## 📋 概述

Laravel-Go Framework 的验证系统提供了强大而灵活的数据验证功能，支持多种验证规则、自定义验证器、错误消息本地化等特性。验证系统可以用于验证 HTTP 请求、表单数据、API 输入等各种场景。

## 🏗️ 核心概念

### 验证器 (Validator)

- 定义验证规则和逻辑
- 处理验证结果和错误
- 支持自定义验证规则

### 验证规则 (Rules)

- 预定义的验证规则集合
- 支持链式调用和组合
- 可扩展的自定义规则

### 验证请求 (Validation Request)

- 封装验证逻辑的请求类
- 自动处理验证和错误响应
- 支持授权和自定义逻辑

## 🔧 基础用法

### 1. 基本验证

```go
// 创建验证器
validator := validation.NewValidator()

// 定义验证规则
rules := map[string]string{
    "name":     "required|string|max:255",
    "email":    "required|email|unique:users",
    "password": "required|min:8|confirmed",
    "age":      "integer|min:18|max:100",
}

// 验证数据
data := map[string]interface{}{
    "name":     "John Doe",
    "email":    "john@example.com",
    "password": "password123",
    "age":      25,
}

// 执行验证
errors := validator.Validate(data, rules)

// 检查验证结果
if len(errors) > 0 {
    // 处理验证错误
    for field, fieldErrors := range errors {
        for _, error := range fieldErrors {
            fmt.Printf("Field %s: %s\n", field, error)
        }
    }
} else {
    // 验证通过
    fmt.Println("Validation passed")
}
```

### 2. 在控制器中使用

```go
// app/Http/Controllers/UserController.go
package controllers

import (
    "laravel-go/framework/http"
    "laravel-go/framework/validation"
)

type UserController struct {
    http.Controller
    validator *validation.Validator
}

func (c *UserController) Store(request http.Request) http.Response {
    // 定义验证规则
    rules := map[string]string{
        "name":     "required|string|max:255",
        "email":    "required|email|unique:users",
        "password": "required|min:8|confirmed",
        "role":     "in:user,admin,moderator",
    }

    // 执行验证
    errors := c.validator.Validate(request.Body, rules)

    if len(errors) > 0 {
        return c.JsonError("Validation failed", 422).WithErrors(errors)
    }

    // 验证通过，继续处理
    user, err := c.userService.CreateUser(request.Body)
    if err != nil {
        return c.JsonError("Failed to create user", 500)
    }

    return c.Json(user).Status(201)
}
```

### 3. 验证请求类

```go
// app/Http/Requests/CreateUserRequest.go
package requests

import (
    "laravel-go/framework/http"
    "laravel-go/framework/validation"
)

type CreateUserRequest struct {
    http.Request
    validator *validation.Validator
}

func (r *CreateUserRequest) Rules() map[string]string {
    return map[string]string{
        "name":     "required|string|max:255",
        "email":    "required|email|unique:users",
        "password": "required|min:8|confirmed",
        "role":     "in:user,admin,moderator",
    }
}

func (r *CreateUserRequest) Messages() map[string]map[string]string {
    return map[string]map[string]string{
        "name": {
            "required": "用户名是必填项",
            "string":   "用户名必须是字符串",
            "max":      "用户名不能超过255个字符",
        },
        "email": {
            "required": "邮箱是必填项",
            "email":    "邮箱格式不正确",
            "unique":   "该邮箱已被注册",
        },
        "password": {
            "required":  "密码是必填项",
            "min":       "密码至少需要8个字符",
            "confirmed": "密码确认不匹配",
        },
    }
}

func (r *CreateUserRequest) Authorize() bool {
    // 检查用户是否有权限创建用户
    return r.Context["user"].(*Models.User).IsAdmin()
}

func (r *CreateUserRequest) Handle() http.Response {
    // 验证失败时会自动返回错误响应
    // 验证通过后执行此方法

    user, err := r.userService.CreateUser(r.Body)
    if err != nil {
        return r.JsonError("Failed to create user", 500)
    }

    return r.Json(user).Status(201)
}
```

## 📚 API 参考

### Validator 接口

```go
type Validator interface {
    Validate(data map[string]interface{}, rules map[string]string) map[string][]string
    ValidateStruct(data interface{}, rules map[string]string) map[string][]string
    AddRule(name string, rule Rule)
    AddRules(rules map[string]Rule)
    SetLocale(locale string)
    GetLocale() string
    SetCustomMessages(messages map[string]map[string]string)
    GetCustomMessages() map[string]map[string]string
}
```

#### 方法说明

- `Validate(data, rules)`: 验证 map 数据
- `ValidateStruct(data, rules)`: 验证结构体数据
- `AddRule(name, rule)`: 添加自定义验证规则
- `AddRules(rules)`: 批量添加验证规则
- `SetLocale(locale)`: 设置语言环境
- `GetLocale()`: 获取当前语言环境
- `SetCustomMessages(messages)`: 设置自定义错误消息
- `GetCustomMessages()`: 获取自定义错误消息

### Rule 接口

```go
type Rule interface {
    Validate(field string, value interface{}, parameters []string) error
    GetMessage(field string, parameters []string) string
}
```

#### 方法说明

- `Validate(field, value, parameters)`: 执行验证逻辑
- `GetMessage(field, parameters)`: 获取错误消息

### 内置验证规则

#### 基础规则

```go
// 必填项
"required"

// 字符串
"string"

// 整数
"integer"

// 浮点数
"numeric"

// 布尔值
"boolean"

// 数组
"array"

// 对象
"object"

// 文件
"file"
```

#### 字符串规则

```go
// 最小长度
"min:10"

// 最大长度
"max:255"

// 长度范围
"between:5,50"

// 正则表达式
"regex:/^[a-zA-Z0-9]+$/"

// 邮箱格式
"email"

// URL 格式
"url"

// 日期格式
"date"

// 日期时间格式
"datetime"

// 时间格式
"time"
```

#### 数值规则

```go
// 最小值
"min:18"

// 最大值
"max:100"

// 数值范围
"between:1,100"

// 正数
"positive"

// 负数
"negative"

// 非零
"nonzero"
```

#### 数组规则

```go
// 数组大小
"size:5"

// 数组最小大小
"min_size:2"

// 数组最大大小
"max_size:10"

// 数组元素类型
"array_of:string"

// 数组元素验证
"array_of:email"
```

#### 比较规则

```go
// 等于
"eq:value"

// 不等于
"ne:value"

// 大于
"gt:value"

// 大于等于
"gte:value"

// 小于
"lt:value"

// 小于等于
"lte:value"
```

#### 数据库规则

```go
// 唯一性
"unique:users,email"

// 存在性
"exists:users,id"

// 不存在
"not_exists:users,email"
```

#### 文件规则

```go
// 文件大小
"file_size:2MB"

// 文件类型
"file_type:image"

// 图片尺寸
"image_size:800,600"

// 图片比例
"image_ratio:16:9"
```

## 🎯 高级功能

### 1. 自定义验证规则

```go
// app/Validation/Rules/StrongPassword.go
package rules

import (
    "laravel-go/framework/validation"
    "regexp"
)

type StrongPassword struct {
    validation.BaseRule
}

func NewStrongPassword() *StrongPassword {
    return &StrongPassword{}
}

func (r *StrongPassword) Validate(field string, value interface{}, parameters []string) error {
    password, ok := value.(string)
    if !ok {
        return validation.NewValidationError(field, "must be a string")
    }

    // 检查密码强度
    if len(password) < 8 {
        return validation.NewValidationError(field, "must be at least 8 characters")
    }

    // 检查是否包含大写字母
    if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
        return validation.NewValidationError(field, "must contain at least one uppercase letter")
    }

    // 检查是否包含小写字母
    if !regexp.MustCompile(`[a-z]`).MatchString(password) {
        return validation.NewValidationError(field, "must contain at least one lowercase letter")
    }

    // 检查是否包含数字
    if !regexp.MustCompile(`[0-9]`).MatchString(password) {
        return validation.NewValidationError(field, "must contain at least one number")
    }

    // 检查是否包含特殊字符
    if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
        return validation.NewValidationError(field, "must contain at least one special character")
    }

    return nil
}

func (r *StrongPassword) GetMessage(field string, parameters []string) string {
    return "密码必须包含大小写字母、数字和特殊字符"
}

// 注册自定义规则
func init() {
    validation.AddRule("strong_password", NewStrongPassword)
}
```

### 2. 条件验证

```go
// 条件验证规则
rules := map[string]string{
    "email":     "required|email",
    "password":  "required_if:email,admin@example.com|min:8",
    "role":      "required|in:user,admin",
    "permissions": "required_if:role,admin|array",
}

// 使用 when 方法进行条件验证
validator := validation.NewValidator()

validator.When("role", "admin", func(v *validation.Validator) {
    v.AddRule("permissions", "required|array")
    v.AddRule("admin_level", "required|integer|between:1,10")
})
```

### 3. 嵌套验证

```go
// 验证嵌套结构
type User struct {
    Name     string    `json:"name" validate:"required|string|max:255"`
    Email    string    `json:"email" validate:"required|email"`
    Profile  Profile   `json:"profile" validate:"required"`
    Posts    []Post    `json:"posts" validate:"array"`
}

type Profile struct {
    Bio      string `json:"bio" validate:"string|max:500"`
    Avatar   string `json:"avatar" validate:"url"`
    Location string `json:"location" validate:"string|max:100"`
}

type Post struct {
    Title   string `json:"title" validate:"required|string|max:255"`
    Content string `json:"content" validate:"required|string"`
}

// 验证嵌套结构
user := &User{
    Name:  "John Doe",
    Email: "john@example.com",
    Profile: Profile{
        Bio:      "Software Developer",
        Avatar:   "https://example.com/avatar.jpg",
        Location: "New York",
    },
    Posts: []Post{
        {
            Title:   "My First Post",
            Content: "This is my first post content",
        },
    },
}

errors := validator.ValidateStruct(user, nil)
```

### 4. 数组验证

```go
// 验证数组元素
rules := map[string]string{
    "tags":           "required|array|min_size:1|max_size:10",
    "tags.*":         "string|max:50",
    "emails":         "array|unique",
    "emails.*":       "email",
    "scores":         "array|between:1,10",
    "scores.*":       "integer|between:0,100",
    "files":          "array|max_size:5",
    "files.*":        "file|file_size:5MB|file_type:image",
}
```

### 5. 自定义错误消息

```go
// 设置自定义错误消息
messages := map[string]map[string]string{
    "name": {
        "required": "用户名是必填项",
        "string":   "用户名必须是字符串",
        "max":      "用户名不能超过255个字符",
    },
    "email": {
        "required": "邮箱是必填项",
        "email":    "邮箱格式不正确",
        "unique":   "该邮箱已被注册",
    },
    "password": {
        "required":  "密码是必填项",
        "min":       "密码至少需要8个字符",
        "confirmed": "密码确认不匹配",
    },
}

validator.SetCustomMessages(messages)
```

## 🔧 配置选项

### 验证系统配置

```go
// config/validation.go
package config

type ValidationConfig struct {
    // 默认语言环境
    DefaultLocale string `json:"default_locale"`

    // 支持的语言环境
    SupportedLocales []string `json:"supported_locales"`

    // 错误消息文件路径
    MessagesPath string `json:"messages_path"`

    // 是否启用快速失败
    FastFail bool `json:"fast_fail"`

    // 最大错误数量
    MaxErrors int `json:"max_errors"`

    // 自定义规则路径
    CustomRulesPath string `json:"custom_rules_path"`

    // 数据库配置
    Database DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
    // 数据库连接
    Connection string `json:"connection"`

    // 表前缀
    TablePrefix string `json:"table_prefix"`

    // 缓存验证结果
    CacheResults bool `json:"cache_results"`

    // 缓存时间
    CacheTTL time.Duration `json:"cache_ttl"`
}
```

### 配置示例

```go
// config/validation.go
func GetValidationConfig() *ValidationConfig {
    return &ValidationConfig{
        DefaultLocale:     "zh-CN",
        SupportedLocales:  []string{"zh-CN", "en-US", "ja-JP"},
        MessagesPath:      "resources/lang/validation",
        FastFail:          false,
        MaxErrors:         100,
        CustomRulesPath:   "app/Validation/Rules",
        Database: DatabaseConfig{
            Connection:    "mysql",
            TablePrefix:   "",
            CacheResults:  true,
            CacheTTL:      time.Hour,
        },
    }
}
```

## 🚀 性能优化

### 1. 验证结果缓存

```go
// 缓存验证结果
type CachedValidator struct {
    validation.Validator
    cache cache.Cache
}

func (v *CachedValidator) Validate(data map[string]interface{}, rules map[string]string) map[string][]string {
    // 生成缓存键
    cacheKey := v.generateCacheKey(data, rules)

    // 尝试从缓存获取
    if cached, exists := v.cache.Get(cacheKey); exists {
        return cached.(map[string][]string)
    }

    // 执行验证
    errors := v.Validator.Validate(data, rules)

    // 缓存结果
    v.cache.Put(cacheKey, errors, time.Minute*5)

    return errors
}

func (v *CachedValidator) generateCacheKey(data map[string]interface{}, rules map[string]string) string {
    // 生成基于数据和规则的缓存键
    dataHash := hashData(data)
    rulesHash := hashRules(rules)
    return fmt.Sprintf("validation:%s:%s", dataHash, rulesHash)
}
```

### 2. 规则预编译

```go
// 预编译验证规则
type CompiledValidator struct {
    validation.Validator
    compiledRules map[string]*CompiledRule
}

type CompiledRule struct {
    Rule       validation.Rule
    Parameters []string
    Compiled   interface{}
}

func (v *CompiledValidator) CompileRules(rules map[string]string) {
    for field, ruleString := range rules {
        rule, parameters := v.parseRule(ruleString)
        v.compiledRules[field] = &CompiledRule{
            Rule:       rule,
            Parameters: parameters,
            Compiled:   v.compileRule(rule, parameters),
        }
    }
}
```

### 3. 并行验证

```go
// 并行执行验证
func (v *Validator) ValidateParallel(data map[string]interface{}, rules map[string]string) map[string][]string {
    var wg sync.WaitGroup
    errors := make(map[string][]string)
    errorMutex := sync.Mutex{}

    for field, rule := range rules {
        wg.Add(1)
        go func(field, rule string) {
            defer wg.Done()

            fieldErrors := v.validateField(field, data[field], rule)
            if len(fieldErrors) > 0 {
                errorMutex.Lock()
                errors[field] = fieldErrors
                errorMutex.Unlock()
            }
        }(field, rule)
    }

    wg.Wait()
    return errors
}
```

## 🧪 测试

### 1. 验证规则测试

```go
// tests/validation_test.go
package tests

import (
    "testing"
    "laravel-go/framework/validation"
)

func TestRequiredRule(t *testing.T) {
    rule := validation.NewRequiredRule()

    // 测试空值
    err := rule.Validate("name", "", nil)
    if err == nil {
        t.Error("Required rule should fail for empty string")
    }

    // 测试非空值
    err = rule.Validate("name", "John", nil)
    if err != nil {
        t.Errorf("Required rule should pass for non-empty string: %v", err)
    }
}

func TestEmailRule(t *testing.T) {
    rule := validation.NewEmailRule()

    // 测试有效邮箱
    err := rule.Validate("email", "john@example.com", nil)
    if err != nil {
        t.Errorf("Email rule should pass for valid email: %v", err)
    }

    // 测试无效邮箱
    err = rule.Validate("email", "invalid-email", nil)
    if err == nil {
        t.Error("Email rule should fail for invalid email")
    }
}
```

### 2. 验证器测试

```go
func TestValidator(t *testing.T) {
    validator := validation.NewValidator()

    // 测试基本验证
    data := map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
        "age":   25,
    }

    rules := map[string]string{
        "name":  "required|string|max:255",
        "email": "required|email",
        "age":   "integer|min:18",
    }

    errors := validator.Validate(data, rules)
    if len(errors) > 0 {
        t.Errorf("Validation should pass: %v", errors)
    }

    // 测试验证失败
    invalidData := map[string]interface{}{
        "name":  "",
        "email": "invalid-email",
        "age":   15,
    }

    errors = validator.Validate(invalidData, rules)
    if len(errors) == 0 {
        t.Error("Validation should fail for invalid data")
    }
}
```

### 3. 自定义规则测试

```go
func TestStrongPasswordRule(t *testing.T) {
    rule := rules.NewStrongPassword()

    // 测试强密码
    err := rule.Validate("password", "StrongPass123!", nil)
    if err != nil {
        t.Errorf("Strong password should pass: %v", err)
    }

    // 测试弱密码
    err = rule.Validate("password", "weak", nil)
    if err == nil {
        t.Error("Weak password should fail")
    }
}
```

## 🔍 调试和监控

### 1. 验证日志

```go
type ValidationLogger struct {
    validation.Validator
    logger log.Logger
}

func (v *ValidationLogger) Validate(data map[string]interface{}, rules map[string]string) map[string][]string {
    start := time.Now()

    errors := v.Validator.Validate(data, rules)

    duration := time.Since(start)

    v.logger.Info("Validation completed", map[string]interface{}{
        "duration": duration,
        "fields":   len(rules),
        "errors":   len(errors),
        "data":     data,
    })

    return errors
}
```

### 2. 验证监控

```go
type ValidationMonitor struct {
    validation.Validator
    metrics metrics.Collector
}

func (v *ValidationMonitor) Validate(data map[string]interface{}, rules map[string]string) map[string][]string {
    // 记录验证指标
    v.metrics.Increment("validation.attempts", map[string]string{
        "rules_count": fmt.Sprintf("%d", len(rules)),
    })

    start := time.Now()
    errors := v.Validator.Validate(data, rules)
    duration := time.Since(start)

    // 记录验证结果
    if len(errors) > 0 {
        v.metrics.Increment("validation.failures", map[string]string{
            "error_count": fmt.Sprintf("%d", len(errors)),
        })
    } else {
        v.metrics.Increment("validation.successes")
    }

    // 记录验证时间
    v.metrics.Histogram("validation.duration", duration.Seconds())

    return errors
}
```

## 📝 最佳实践

### 1. 验证规则组织

```go
// 将验证规则组织到单独的文件中
// app/Validation/Rules/UserRules.go
package rules

var UserRules = map[string]string{
    "name":     "required|string|max:255",
    "email":    "required|email|unique:users",
    "password": "required|min:8|confirmed",
    "role":     "in:user,admin,moderator",
}

var UserUpdateRules = map[string]string{
    "name":  "string|max:255",
    "email": "email|unique:users,email," + "{id}",
    "role":  "in:user,admin,moderator",
}

// 在控制器中使用
func (c *UserController) Store(request http.Request) http.Response {
    errors := c.validator.Validate(request.Body, rules.UserRules)
    // ...
}
```

### 2. 错误消息本地化

```go
// resources/lang/zh-CN/validation.php
{
    "required": "字段 :field 是必填项",
    "email": "字段 :field 必须是有效的邮箱地址",
    "min": "字段 :field 至少需要 :min 个字符",
    "max": "字段 :field 不能超过 :max 个字符",
    "unique": "字段 :field 的值已经存在",
    "confirmed": "字段 :field 确认不匹配",
    "in": "字段 :field 的值无效",
    "integer": "字段 :field 必须是整数",
    "string": "字段 :field 必须是字符串",
    "array": "字段 :field 必须是数组",
    "file": "字段 :field 必须是文件",
    "url": "字段 :field 必须是有效的URL",
    "date": "字段 :field 必须是有效的日期",
    "between": "字段 :field 必须在 :min 和 :max 之间",
    "exists": "字段 :field 的值不存在",
    "not_exists": "字段 :field 的值已存在",
}
```

### 3. 验证请求类

```go
// 使用验证请求类封装验证逻辑
type CreateUserRequest struct {
    http.Request
}

func (r *CreateUserRequest) Rules() map[string]string {
    return map[string]string{
        "name":     "required|string|max:255",
        "email":    "required|email|unique:users",
        "password": "required|min:8|confirmed",
        "role":     "in:user,admin,moderator",
    }
}

func (r *CreateUserRequest) Messages() map[string]map[string]string {
    return map[string]map[string]string{
        "name": {
            "required": "用户名是必填项",
            "string":   "用户名必须是字符串",
            "max":      "用户名不能超过255个字符",
        },
        "email": {
            "required": "邮箱是必填项",
            "email":    "邮箱格式不正确",
            "unique":   "该邮箱已被注册",
        },
    }
}

func (r *CreateUserRequest) Authorize() bool {
    return r.Context["user"].(*Models.User).IsAdmin()
}
```

### 4. 验证中间件

```go
// 创建验证中间件
type ValidationMiddleware struct {
    http.Middleware
    validator *validation.Validator
}

func (m *ValidationMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 检查是否需要验证
    if rules, ok := request.Context["validation_rules"]; ok {
        errors := m.validator.Validate(request.Body, rules.(map[string]string))
        if len(errors) > 0 {
            return http.Response{
                StatusCode: 422,
                Body:       m.formatErrors(errors),
                Headers: map[string]string{
                    "Content-Type": "application/json",
                },
            }
        }
    }

    return next(request)
}

func (m *ValidationMiddleware) formatErrors(errors map[string][]string) string {
    // 格式化错误响应
    response := map[string]interface{}{
        "message": "Validation failed",
        "errors":  errors,
    }

    jsonData, _ := json.Marshal(response)
    return string(jsonData)
}
```

## 🚀 总结

验证系统是 Laravel-Go Framework 中重要的功能之一，它提供了：

1. **完整的验证功能**: 支持多种验证规则和自定义规则
2. **灵活的配置**: 支持条件验证、嵌套验证等高级功能
3. **性能优化**: 提供缓存、预编译等性能优化方案
4. **错误处理**: 完善的错误消息和本地化支持
5. **测试支持**: 完整的测试框架和工具
6. **最佳实践**: 遵循验证系统的最佳实践

通过合理使用验证系统，可以确保应用程序的数据完整性和安全性，提供更好的用户体验。
