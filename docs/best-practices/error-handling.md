# 错误处理最佳实践

## 🚨 错误处理概览

错误处理是构建可靠应用程序的关键部分。Laravel-Go Framework 提供了完善的错误处理机制，帮助开发者构建健壮的应用程序。

## 🚀 快速开始

### 基本错误处理

```go
import "laravel-go/framework/errors"

// 创建自定义错误
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
)

// 使用错误
func GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, ErrInvalidInput
    }

    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }

    if user == nil {
        return nil, ErrUserNotFound
    }

    return user, nil
}
```

## 📋 错误类型

### 1. 基础错误类型

```go
// 应用错误
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
    return e.Message
}

// 验证错误
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   interface{} `json:"value,omitempty"`
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}

// 业务错误
type BusinessError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func (e *BusinessError) Error() string {
    return e.Message
}
```

### 2. HTTP 错误

```go
// HTTP 错误
type HTTPError struct {
    StatusCode int         `json:"-"`
    Code       string      `json:"code"`
    Message    string      `json:"message"`
    Details    interface{} `json:"details,omitempty"`
}

func (e *HTTPError) Error() string {
    return e.Message
}

// 预定义 HTTP 错误
var (
    ErrBadRequest = &HTTPError{
        StatusCode: 400,
        Code:       "BAD_REQUEST",
        Message:    "Bad request",
    }

    ErrUnauthorized = &HTTPError{
        StatusCode: 401,
        Code:       "UNAUTHORIZED",
        Message:    "Unauthorized",
    }

    ErrForbidden = &HTTPError{
        StatusCode: 403,
        Code:       "FORBIDDEN",
        Message:    "Forbidden",
    }

    ErrNotFound = &HTTPError{
        StatusCode: 404,
        Code:       "NOT_FOUND",
        Message:    "Resource not found",
    }

    ErrInternalServer = &HTTPError{
        StatusCode: 500,
        Code:       "INTERNAL_SERVER_ERROR",
        Message:    "Internal server error",
    }
)
```

### 3. 数据库错误

```go
// 数据库错误
type DatabaseError struct {
    Operation string `json:"operation"`
    Table     string `json:"table,omitempty"`
    Query     string `json:"query,omitempty"`
    Err       error  `json:"-"`
}

func (e *DatabaseError) Error() string {
    return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

func (e *DatabaseError) Unwrap() error {
    return e.Err
}

// 数据库错误类型
var (
    ErrRecordNotFound = errors.New("record not found")
    ErrDuplicateKey   = errors.New("duplicate key")
    ErrForeignKey     = errors.New("foreign key constraint failed")
    ErrConnection     = errors.New("database connection failed")
)
```

## 🛠️ 错误处理策略

### 1. 错误包装

```go
// 使用 fmt.Errorf 包装错误
func GetUser(id int) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return user, nil
}

// 使用 errors.Wrap 包装错误
func ProcessUser(userID int) error {
    user, err := GetUser(userID)
    if err != nil {
        return errors.Wrap(err, "failed to process user")
    }

    err = user.Process()
    if err != nil {
        return errors.Wrap(err, "failed to process user data")
    }

    return nil
}
```

### 2. 错误检查

```go
// 检查特定错误类型
func HandleError(err error) {
    if errors.Is(err, ErrUserNotFound) {
        // 处理用户未找到错误
        log.Printf("User not found: %v", err)
        return
    }

    if errors.Is(err, ErrInvalidInput) {
        // 处理无效输入错误
        log.Printf("Invalid input: %v", err)
        return
    }

    // 处理其他错误
    log.Printf("Unknown error: %v", err)
}

// 检查错误类型
func HandleDatabaseError(err error) {
    var dbErr *DatabaseError
    if errors.As(err, &dbErr) {
        // 处理数据库错误
        log.Printf("Database error during %s: %v", dbErr.Operation, dbErr.Err)
        return
    }

    // 处理其他错误
    log.Printf("Non-database error: %v", err)
}
```

### 3. 错误恢复

```go
// 使用 defer 和 recover 恢复 panic
func SafeFunction() (result interface{}, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
            log.Printf("Panic recovered: %v", r)
        }
    }()

    // 可能发生 panic 的代码
    result = riskyOperation()
    return result, nil
}

// 中间件中的错误恢复
func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v", err)

                // 返回 500 错误
                http.Error(w, "Internal server error", http.StatusInternalServerError)
            }
        }()

        next(w, r)
    }
}
```

## 🎯 控制器错误处理

### 1. 统一错误响应

```go
type ErrorResponse struct {
    Success bool        `json:"success"`
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

func (c *Controller) HandleError(err error) http.Response {
    var httpErr *HTTPError
    if errors.As(err, &httpErr) {
        return c.Json(ErrorResponse{
            Success: false,
            Code:    httpErr.Code,
            Message: httpErr.Message,
            Details: httpErr.Details,
        }).Status(httpErr.StatusCode)
    }

    // 处理其他错误类型
    if errors.Is(err, ErrUserNotFound) {
        return c.Json(ErrorResponse{
            Success: false,
            Code:    "USER_NOT_FOUND",
            Message: "User not found",
        }).Status(404)
    }

    if errors.Is(err, ErrInvalidInput) {
        return c.Json(ErrorResponse{
            Success: false,
            Code:    "INVALID_INPUT",
            Message: "Invalid input",
        }).Status(400)
    }

    // 默认错误响应
    return c.Json(ErrorResponse{
        Success: false,
        Code:    "INTERNAL_ERROR",
        Message: "Internal server error",
    }).Status(500)
}
```

### 2. 控制器方法中的错误处理

```go
func (c *UserController) Show(id string) http.Response {
    userID, err := strconv.Atoi(id)
    if err != nil {
        return c.HandleError(ErrInvalidInput)
    }

    user, err := c.userService.GetUser(userID)
    if err != nil {
        return c.HandleError(err)
    }

    return c.Json(user)
}

func (c *UserController) Store(request http.Request) http.Response {
    // 验证输入
    if err := c.validateUserData(request.Body); err != nil {
        return c.HandleError(err)
    }

    // 创建用户
    user, err := c.userService.CreateUser(request.Body)
    if err != nil {
        return c.HandleError(err)
    }

    return c.Json(user).Status(201)
}
```

## 🔧 服务层错误处理

### 1. 服务方法错误处理

```go
type UserService struct {
    db    *database.Connection
    cache cache.Cache
}

func (s *UserService) GetUser(id int) (*User, error) {
    // 参数验证
    if id <= 0 {
        return nil, ErrInvalidInput
    }

    // 尝试从缓存获取
    cacheKey := fmt.Sprintf("user:%d", id)
    if cached, exists := s.cache.Get(cacheKey); exists {
        return cached.(*User), nil
    }

    // 从数据库获取
    var user User
    err := s.db.First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, &DatabaseError{
            Operation: "find user",
            Table:     "users",
            Err:       err,
        }
    }

    // 缓存用户数据
    s.cache.Put(cacheKey, &user, time.Hour)

    return &user, nil
}

func (s *UserService) CreateUser(data map[string]interface{}) (*User, error) {
    // 验证数据
    if err := s.validateUserData(data); err != nil {
        return nil, err
    }

    // 检查邮箱是否已存在
    var existingUser User
    err := s.db.Where("email = ?", data["email"]).First(&existingUser).Error
    if err == nil {
        return nil, &BusinessError{
            Code:    "EMAIL_EXISTS",
            Message: "Email already exists",
        }
    }

    // 创建用户
    user := &User{
        Name:     data["name"].(string),
        Email:    data["email"].(string),
        Password: hashPassword(data["password"].(string)),
    }

    err = s.db.Create(user).Error
    if err != nil {
        return nil, &DatabaseError{
            Operation: "create user",
            Table:     "users",
            Err:       err,
        }
    }

    return user, nil
}
```

### 2. 事务错误处理

```go
func (s *UserService) CreateUserWithProfile(userData, profileData map[string]interface{}) (*User, error) {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 创建用户
        user := &User{
            Name:     userData["name"].(string),
            Email:    userData["email"].(string),
            Password: hashPassword(userData["password"].(string)),
        }

        err := tx.Create(user).Error
        if err != nil {
            return &DatabaseError{
                Operation: "create user",
                Table:     "users",
                Err:       err,
            }
        }

        // 创建用户资料
        profile := &Profile{
            UserID: user.ID,
            Bio:    profileData["bio"].(string),
        }

        err = tx.Create(profile).Error
        if err != nil {
            return &DatabaseError{
                Operation: "create profile",
                Table:     "profiles",
                Err:       err,
            }
        }

        return nil
    })
}
```

## 🔍 中间件错误处理

### 1. 错误处理中间件

```go
func ErrorHandlingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v", err)

                response := ErrorResponse{
                    Success: false,
                    Code:    "INTERNAL_ERROR",
                    Message: "Internal server error",
                }

                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusInternalServerError)
                json.NewEncoder(w).Encode(response)
            }
        }()

        next(w, r)
    }
}
```

### 2. 日志中间件

```go
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // 包装 ResponseWriter 以捕获状态码
        wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: 200}

        next(wrappedWriter, r)

        duration := time.Since(start)

        // 记录请求日志
        log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrappedWriter.statusCode, duration)

        // 记录错误
        if wrappedWriter.statusCode >= 400 {
            log.Printf("Error: %s %s returned %d", r.Method, r.URL.Path, wrappedWriter.statusCode)
        }
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

## 📊 错误监控

### 1. 错误统计

```go
type ErrorMonitor struct {
    errorCounts map[string]int64
    mu          sync.RWMutex
}

func (m *ErrorMonitor) RecordError(err error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    errorType := "unknown"

    var httpErr *HTTPError
    if errors.As(err, &httpErr) {
        errorType = httpErr.Code
    } else if errors.Is(err, ErrUserNotFound) {
        errorType = "user_not_found"
    } else if errors.Is(err, ErrInvalidInput) {
        errorType = "invalid_input"
    }

    m.errorCounts[errorType]++
}

func (m *ErrorMonitor) GetErrorStats() map[string]int64 {
    m.mu.RLock()
    defer m.mu.RUnlock()

    stats := make(map[string]int64)
    for errorType, count := range m.errorCounts {
        stats[errorType] = count
    }

    return stats
}
```

### 2. 错误告警

```go
type ErrorAlerter struct {
    threshold int64
    monitor   *ErrorMonitor
}

func (a *ErrorAlerter) CheckAlerts() {
    stats := a.monitor.GetErrorStats()

    for errorType, count := range stats {
        if count > a.threshold {
            a.sendAlert(errorType, count)
        }
    }
}

func (a *ErrorAlerter) sendAlert(errorType string, count int64) {
    message := fmt.Sprintf("High error rate detected for %s: %d errors", errorType, count)

    // 发送告警（邮件、Slack、钉钉等）
    log.Printf("ALERT: %s", message)
}
```

## 📝 最佳实践

### 1. 错误设计原则

```go
// ✅ 好的做法：错误信息清晰明确
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidEmail = errors.New("invalid email format")
var ErrPasswordTooShort = errors.New("password must be at least 8 characters")

// ❌ 不好的做法：错误信息模糊
var ErrBadInput = errors.New("bad input")
var ErrSomethingWentWrong = errors.New("something went wrong")
```

### 2. 错误处理原则

```go
// ✅ 好的做法：在错误发生的地方处理
func GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, ErrInvalidInput
    }

    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }

    return user, nil
}

// ❌ 不好的做法：忽略错误
func GetUser(id int) (*User, error) {
    user, _ := db.FindUser(id) // 忽略错误
    return user, nil
}
```

### 3. 错误日志记录

```go
// ✅ 好的做法：记录有意义的错误信息
func HandleError(err error) {
    log.Printf("Failed to process user request: %v", err)

    // 记录堆栈信息
    if debug {
        log.Printf("Stack trace: %s", debug.Stack())
    }
}

// ❌ 不好的做法：记录过多或过少信息
func HandleError(err error) {
    log.Printf("Error: %v", err) // 信息太少
    // 或者
    log.Printf("Everything failed: %v, user: %v, request: %v, context: %v", err, user, request, context) // 信息太多
}
```

### 4. 错误恢复策略

```go
// ✅ 好的做法：优雅的错误恢复
func SafeOperation() error {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
            // 执行清理操作
            cleanup()
        }
    }()

    return riskyOperation()
}

// ❌ 不好的做法：忽略 panic
func UnsafeOperation() error {
    return riskyOperation() // 可能 panic
}
```

## 📚 总结

Laravel-Go Framework 的错误处理最佳实践包括：

1. **错误类型设计**: 定义清晰的错误类型和错误码
2. **错误包装**: 使用错误包装保持错误上下文
3. **统一处理**: 在控制器层统一处理错误响应
4. **错误监控**: 实现错误统计和告警机制
5. **优雅恢复**: 使用 defer 和 recover 处理 panic

通过遵循这些最佳实践，可以构建出健壮、可维护的应用程序。
