# 验证系统指南

## 📖 概述

Laravel-Go Framework 提供了强大的数据验证系统，支持多种验证规则、自定义验证器、错误消息本地化等功能，确保应用程序数据的完整性和安全性。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [验证系统 API 参考](../api/validation.md)

## 🚀 快速开始

### 1. 基本使用

```go
// 创建验证器
type UserRequest struct {
    Name     string `json:"name" validate:"required,min:2,max:50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min:8"`
    Age      int    `json:"age" validate:"required,min:18,max:100"`
}

// 验证数据
func (c *UserController) Store(request http.Request) http.Response {
    var userRequest UserRequest

    // 绑定请求数据
    if err := request.Bind(&userRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    // 验证数据
    if err := validator.Validate(&userRequest); err != nil {
        return c.JsonError(err.Error(), 422)
    }

    // 创建用户
    user, err := c.userService.CreateUser(userRequest)
    if err != nil {
        return c.JsonError("Failed to create user", 500)
    }

    return c.Json(user).Status(201)
}
```

### 2. 使用验证器实例

```go
// 创建验证器实例
func (c *UserController) Store(request http.Request) http.Response {
    var userRequest UserRequest

    if err := request.Bind(&userRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    // 创建验证器
    v := validator.New()

    // 添加验证规则
    v.Required("name", userRequest.Name, "Name is required")
    v.MinLength("name", userRequest.Name, 2, "Name must be at least 2 characters")
    v.MaxLength("name", userRequest.Name, 50, "Name must not exceed 50 characters")
    v.Email("email", userRequest.Email, "Invalid email format")
    v.MinLength("password", userRequest.Password, 8, "Password must be at least 8 characters")
    v.Range("age", userRequest.Age, 18, 100, "Age must be between 18 and 100")

    // 执行验证
    if !v.Passes() {
        return c.JsonError(v.Errors(), 422)
    }

    // 处理业务逻辑
    user, err := c.userService.CreateUser(userRequest)
    if err != nil {
        return c.JsonError("Failed to create user", 500)
    }

    return c.Json(user).Status(201)
}
```

## 📋 验证规则

### 1. 基础验证规则

```go
// 必填验证
v.Required("field", value, "Field is required")

// 字符串验证
v.MinLength("field", value, 2, "Field must be at least 2 characters")
v.MaxLength("field", value, 50, "Field must not exceed 50 characters")
v.Length("field", value, 10, "Field must be exactly 10 characters")
v.Email("field", value, "Invalid email format")
v.URL("field", value, "Invalid URL format")
v.Alpha("field", value, "Field must contain only letters")
v.AlphaNum("field", value, "Field must contain only letters and numbers")
v.Numeric("field", value, "Field must be numeric")
v.Integer("field", value, "Field must be an integer")

// 数值验证
v.Min("field", value, 0, "Field must be at least 0")
v.Max("field", value, 100, "Field must not exceed 100")
v.Range("field", value, 1, 10, "Field must be between 1 and 10")
v.Positive("field", value, "Field must be positive")
v.Negative("field", value, "Field must be negative")

// 日期验证
v.Date("field", value, "Invalid date format")
v.DateAfter("field", value, "2023-01-01", "Date must be after 2023-01-01")
v.DateBefore("field", value, "2024-12-31", "Date must be before 2024-12-31")
v.DateBetween("field", value, "2023-01-01", "2024-12-31", "Date must be between 2023-01-01 and 2024-12-31")
```

### 2. 高级验证规则

```go
// 正则表达式验证
v.Regex("field", value, `^[A-Za-z0-9]+$`, "Field must contain only alphanumeric characters")

// 唯一性验证
v.Unique("email", value, "users", "email", "Email already exists")

// 存在性验证
v.Exists("user_id", value, "users", "id", "User does not exist")

// 条件验证
v.RequiredIf("field", value, "other_field", "other_value", "Field is required when other_field is other_value")
v.RequiredUnless("field", value, "other_field", "other_value", "Field is required unless other_field is other_value")

// 数组验证
v.Array("field", value, "Field must be an array")
v.MinArray("field", value, 1, "Field must have at least 1 item")
v.MaxArray("field", value, 10, "Field must not have more than 10 items")
v.ArraySize("field", value, 5, "Field must have exactly 5 items")

// 文件验证
v.File("field", value, "Field must be a file")
v.Image("field", value, "Field must be an image")
v.MaxFileSize("field", value, 1024*1024, "File size must not exceed 1MB")
v.FileType("field", value, []string{".jpg", ".png", ".gif"}, "Invalid file type")
```

### 3. 自定义验证规则

```go
// 自定义验证器
type CustomValidator struct {
    validator.Validator
}

// 验证手机号
func (v *CustomValidator) Phone(field, value, message string) bool {
    if message == "" {
        message = "Invalid phone number format"
    }

    phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
    if !phoneRegex.MatchString(value) {
        v.AddError(field, message)
        return false
    }

    return true
}

// 验证身份证号
func (v *CustomValidator) IDCard(field, value, message string) bool {
    if message == "" {
        message = "Invalid ID card number"
    }

    idCardRegex := regexp.MustCompile(`^\d{17}[\dXx]$`)
    if !idCardRegex.MatchString(value) {
        v.AddError(field, message)
        return false
    }

    return true
}

// 验证密码强度
func (v *CustomValidator) PasswordStrength(field, value, message string) bool {
    if message == "" {
        message = "Password must contain uppercase, lowercase, number and special character"
    }

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(value)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(value)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(value)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(value)

    if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
        v.AddError(field, message)
        return false
    }

    return true
}

// 使用自定义验证器
func (c *UserController) Store(request http.Request) http.Response {
    var userRequest UserRequest

    if err := request.Bind(&userRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    v := &CustomValidator{}

    // 使用自定义验证规则
    v.Required("name", userRequest.Name, "Name is required")
    v.Email("email", userRequest.Email, "Invalid email format")
    v.Phone("phone", userRequest.Phone, "Invalid phone number")
    v.PasswordStrength("password", userRequest.Password, "Password is too weak")

    if !v.Passes() {
        return c.JsonError(v.Errors(), 422)
    }

    // 处理业务逻辑...
}
```

## 🏷️ 标签验证

### 1. 结构体标签验证

```go
// 使用结构体标签定义验证规则
type UserRequest struct {
    Name     string `json:"name" validate:"required,min:2,max:50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min:8,password_strength"`
    Phone    string `json:"phone" validate:"required,phone"`
    Age      int    `json:"age" validate:"required,min:18,max:100"`
    Avatar   string `json:"avatar" validate:"omitempty,url"`
    Bio      string `json:"bio" validate:"omitempty,max:500"`
}

// 验证结构体
func ValidateStruct(data interface{}) error {
    v := validator.New()

    // 解析结构体标签
    if err := v.Struct(data); err != nil {
        return err
    }

    return nil
}

// 在控制器中使用
func (c *UserController) Store(request http.Request) http.Response {
    var userRequest UserRequest

    if err := request.Bind(&userRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    if err := ValidateStruct(&userRequest); err != nil {
        return c.JsonError(err.Error(), 422)
    }

    // 处理业务逻辑...
}
```

### 2. 嵌套结构体验证

```go
// 嵌套结构体
type Address struct {
    Street  string `json:"street" validate:"required"`
    City    string `json:"city" validate:"required"`
    State   string `json:"state" validate:"required"`
    ZipCode string `json:"zip_code" validate:"required,regex:^\\d{5}$"`
}

type UserRequest struct {
    Name    string  `json:"name" validate:"required,min:2,max:50"`
    Email   string  `json:"email" validate:"required,email"`
    Address Address `json:"address" validate:"required"`
}

// 验证嵌套结构体
func (c *UserController) Store(request http.Request) http.Response {
    var userRequest UserRequest

    if err := request.Bind(&userRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    if err := ValidateStruct(&userRequest); err != nil {
        return c.JsonError(err.Error(), 422)
    }

    // 处理业务逻辑...
}
```

## 🔄 条件验证

### 1. 条件验证规则

```go
// 条件验证
type OrderRequest struct {
    PaymentMethod string  `json:"payment_method" validate:"required,in:credit_card,bank_transfer,cash"`
    CreditCard    *string `json:"credit_card" validate:"required_if:payment_method,credit_card"`
    BankAccount   *string `json:"bank_account" validate:"required_if:payment_method,bank_transfer"`
    CashAmount    *float64 `json:"cash_amount" validate:"required_if:payment_method,cash"`
}

// 使用条件验证
func (c *OrderController) Store(request http.Request) http.Response {
    var orderRequest OrderRequest

    if err := request.Bind(&orderRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    v := validator.New()

    // 基础验证
    v.Required("payment_method", orderRequest.PaymentMethod, "Payment method is required")
    v.In("payment_method", orderRequest.PaymentMethod, []string{"credit_card", "bank_transfer", "cash"}, "Invalid payment method")

    // 条件验证
    if orderRequest.PaymentMethod == "credit_card" {
        v.Required("credit_card", orderRequest.CreditCard, "Credit card is required for credit card payment")
    }

    if orderRequest.PaymentMethod == "bank_transfer" {
        v.Required("bank_account", orderRequest.BankAccount, "Bank account is required for bank transfer")
    }

    if orderRequest.PaymentMethod == "cash" {
        v.Required("cash_amount", orderRequest.CashAmount, "Cash amount is required for cash payment")
    }

    if !v.Passes() {
        return c.JsonError(v.Errors(), 422)
    }

    // 处理业务逻辑...
}
```

### 2. 复杂条件验证

```go
// 复杂条件验证
type ComplexRequest struct {
    Type        string `json:"type" validate:"required,in:individual,company"`
    Name        string `json:"name" validate:"required"`
    Email       string `json:"email" validate:"required,email"`
    Phone       string `json:"phone" validate:"required,phone"`

    // 个人用户字段
    IDCard      *string `json:"id_card"`
    Birthday    *string `json:"birthday"`

    // 公司用户字段
    CompanyName *string `json:"company_name"`
    TaxNumber   *string `json:"tax_number"`
    License     *string `json:"license"`
}

// 复杂条件验证逻辑
func (c *UserController) Store(request http.Request) http.Response {
    var complexRequest ComplexRequest

    if err := request.Bind(&complexRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    v := validator.New()

    // 基础验证
    v.Required("type", complexRequest.Type, "Type is required")
    v.In("type", complexRequest.Type, []string{"individual", "company"}, "Invalid type")
    v.Required("name", complexRequest.Name, "Name is required")
    v.Email("email", complexRequest.Email, "Invalid email format")
    v.Phone("phone", complexRequest.Phone, "Invalid phone number")

    // 根据类型进行条件验证
    if complexRequest.Type == "individual" {
        v.Required("id_card", complexRequest.IDCard, "ID card is required for individual")
        v.Required("birthday", complexRequest.Birthday, "Birthday is required for individual")

        if complexRequest.IDCard != nil {
            v.IDCard("id_card", *complexRequest.IDCard, "Invalid ID card number")
        }

        if complexRequest.Birthday != nil {
            v.Date("birthday", *complexRequest.Birthday, "Invalid birthday format")
        }
    }

    if complexRequest.Type == "company" {
        v.Required("company_name", complexRequest.CompanyName, "Company name is required for company")
        v.Required("tax_number", complexRequest.TaxNumber, "Tax number is required for company")
        v.Required("license", complexRequest.License, "License is required for company")

        if complexRequest.TaxNumber != nil {
            v.Regex("tax_number", *complexRequest.TaxNumber, `^\d{15}$`, "Invalid tax number format")
        }
    }

    if !v.Passes() {
        return c.JsonError(v.Errors(), 422)
    }

    // 处理业务逻辑...
}
```

## 🛡️ 错误处理

### 1. 错误格式化

```go
// 错误格式化
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   interface{} `json:"value"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
    if len(v) == 0 {
        return ""
    }

    messages := make([]string, len(v))
    for i, err := range v {
        messages[i] = fmt.Sprintf("%s: %s", err.Field, err.Message)
    }

    return strings.Join(messages, "; ")
}

// 格式化验证错误
func FormatValidationErrors(errors map[string]string) ValidationErrors {
    var validationErrors ValidationErrors

    for field, message := range errors {
        validationErrors = append(validationErrors, ValidationError{
            Field:   field,
            Message: message,
        })
    }

    return validationErrors
}

// 在控制器中使用
func (c *UserController) Store(request http.Request) http.Response {
    var userRequest UserRequest

    if err := request.Bind(&userRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    v := validator.New()

    // 添加验证规则...

    if !v.Passes() {
        validationErrors := FormatValidationErrors(v.Errors())
        return c.Json(map[string]interface{}{
            "message": "Validation failed",
            "errors":  validationErrors,
        }).Status(422)
    }

    // 处理业务逻辑...
}
```

### 2. 国际化错误消息

```go
// 国际化错误消息
type I18nValidator struct {
    validator.Validator
    locale string
}

func NewI18nValidator(locale string) *I18nValidator {
    return &I18nValidator{
        locale: locale,
    }
}

func (v *I18nValidator) getMessage(key string) string {
    messages := map[string]map[string]string{
        "en": {
            "required": "The {field} field is required",
            "email":    "The {field} field must be a valid email address",
            "min":      "The {field} field must be at least {min} characters",
            "max":      "The {field} field must not exceed {max} characters",
        },
        "zh": {
            "required": "{field} 字段是必填的",
            "email":    "{field} 字段必须是有效的邮箱地址",
            "min":      "{field} 字段至少需要 {min} 个字符",
            "max":      "{field} 字段不能超过 {max} 个字符",
        },
    }

    if localeMessages, exists := messages[v.locale]; exists {
        if message, exists := localeMessages[key]; exists {
            return message
        }
    }

    return key
}

func (v *I18nValidator) Required(field, value, message string) bool {
    if message == "" {
        message = v.getMessage("required")
        message = strings.Replace(message, "{field}", field, -1)
    }

    if value == "" {
        v.AddError(field, message)
        return false
    }

    return true
}

// 使用国际化验证器
func (c *UserController) Store(request http.Request) http.Response {
    var userRequest UserRequest

    if err := request.Bind(&userRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    // 从请求头获取语言
    locale := request.Headers["Accept-Language"]
    if locale == "" {
        locale = "en"
    }

    v := NewI18nValidator(locale)

    // 添加验证规则（消息会自动国际化）
    v.Required("name", userRequest.Name, "")
    v.Email("email", userRequest.Email, "")
    v.MinLength("password", userRequest.Password, 8, "")

    if !v.Passes() {
        return c.JsonError(v.Errors(), 422)
    }

    // 处理业务逻辑...
}
```

## 📊 验证中间件

### 1. 验证中间件

```go
// 验证中间件
type ValidationMiddleware struct {
    http.Middleware
    validator interface{}
}

func (m *ValidationMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 获取验证器类型
    validatorType := reflect.TypeOf(m.validator)
    validatorValue := reflect.New(validatorType.Elem()).Interface()

    // 绑定请求数据
    if err := request.Bind(validatorValue); err != nil {
        return http.Response{
            StatusCode: 400,
            Body:       `{"error": "Invalid request data"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }

    // 验证数据
    if err := ValidateStruct(validatorValue); err != nil {
        return http.Response{
            StatusCode: 422,
            Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }

    // 将验证后的数据添加到请求上下文
    request.Context["validated_data"] = validatorValue

    return next(request)
}

// 使用验证中间件
func RegisterRoutes() {
    router := routing.NewRouter()

    // 用户注册路由
    router.Post("/users", &UserController{}, "Store").
        Use(&ValidationMiddleware{validator: &UserRequest{}})

    // 订单创建路由
    router.Post("/orders", &OrderController{}, "Store").
        Use(&ValidationMiddleware{validator: &OrderRequest{}})
}
```

### 2. 自动验证中间件

```go
// 自动验证中间件
type AutoValidationMiddleware struct {
    http.Middleware
}

func (m *AutoValidationMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 从路由获取控制器和方法
    controller := request.Context["controller"].(http.Controller)
    method := request.Context["method"].(string)

    // 根据控制器和方法确定验证器
    validator := m.getValidator(controller, method)
    if validator == nil {
        return next(request)
    }

    // 绑定和验证数据
    if err := request.Bind(validator); err != nil {
        return http.Response{
            StatusCode: 400,
            Body:       `{"error": "Invalid request data"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }

    if err := ValidateStruct(validator); err != nil {
        return http.Response{
            StatusCode: 422,
            Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }

    request.Context["validated_data"] = validator
    return next(request)
}

func (m *AutoValidationMiddleware) getValidator(controller http.Controller, method string) interface{} {
    // 根据控制器和方法返回对应的验证器
    switch controller.(type) {
    case *UserController:
        switch method {
        case "Store":
            return &UserRequest{}
        case "Update":
            return &UserUpdateRequest{}
        }
    case *OrderController:
        switch method {
        case "Store":
            return &OrderRequest{}
        }
    }

    return nil
}
```

## 📚 总结

Laravel-Go Framework 的验证系统提供了：

1. **基础验证规则**: 必填、字符串、数值、日期验证
2. **高级验证规则**: 正则、唯一性、存在性、条件验证
3. **自定义验证器**: 扩展验证规则
4. **标签验证**: 结构体标签定义验证规则
5. **条件验证**: 根据条件进行验证
6. **错误处理**: 错误格式化和国际化
7. **验证中间件**: 自动验证请求数据

通过合理使用验证系统，可以确保应用程序数据的完整性和安全性。
