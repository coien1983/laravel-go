# 代码规范

本文档定义了 Laravel-Go Framework 项目的代码编写规范，确保代码质量和一致性。

## 📋 目录

- [命名规范](#命名规范)
- [代码格式](#代码格式)
- [注释规范](#注释规范)
- [错误处理](#错误处理)
- [性能优化](#性能优化)
- [安全规范](#安全规范)
- [测试规范](#测试规范)

## 🏷️ 命名规范

### 包命名

```go
// ✅ 正确：使用小写字母，简短且有意义
package http
package database
package cache

// ❌ 错误：使用下划线或大写
package HTTP
package database_utils
```

### 文件命名

```go
// ✅ 正确：使用小写字母和下划线
user_controller.go
database_connection.go
auth_middleware.go

// ❌ 错误：使用大写或连字符
UserController.go
database-connection.go
```

### 变量命名

```go
// ✅ 正确：使用驼峰命名法
var userName string
var isAuthenticated bool
var maxRetryCount int

// ❌ 错误：使用下划线或大写开头
var user_name string
var IsAuthenticated bool
```

### 常量命名

```go
// ✅ 正确：使用大写字母和下划线
const (
    MAX_RETRY_COUNT = 3
    DEFAULT_TIMEOUT = 30
    API_VERSION    = "v1"
)

// ❌ 错误：使用驼峰命名法
const (
    maxRetryCount = 3
    defaultTimeout = 30
)
```

### 函数命名

```go
// ✅ 正确：使用驼峰命名法，动词开头
func GetUser(id int) (*User, error)
func CreateUser(user *User) error
func IsValidEmail(email string) bool

// ❌ 错误：使用下划线或名词开头
func get_user(id int) (*User, error)
func UserCreate(user *User) error
```

### 结构体命名

```go
// ✅ 正确：使用驼峰命名法，名词
type UserController struct{}
type DatabaseConnection struct{}
type AuthMiddleware struct{}

// ❌ 错误：使用下划线或动词
type user_controller struct{}
type CreateUser struct{}
```

### 接口命名

```go
// ✅ 正确：使用驼峰命名法，通常以 -er 结尾
type UserRepository interface{}
type CacheDriver interface{}
type EventListener interface{}

// ❌ 错误：使用下划线或不以 -er 结尾
type user_repository interface{}
type Cache interface{}
```

## 📝 代码格式

### 导入顺序

```go
// ✅ 正确：标准库、第三方库、本地包
import (
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"

    "laravel-go/framework/http"
    "laravel-go/framework/database"
)
```

### 结构体定义

```go
// ✅ 正确：字段按重要性排序，标签对齐
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" validate:"required"`
    Email     string    `json:"email" validate:"required,email"`
    Password  string    `json:"-" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ❌ 错误：字段顺序混乱，标签不对齐
type User struct {
    CreatedAt time.Time `json:"created_at"`
    ID uint `json:"id" gorm:"primaryKey"`
    Password string `json:"-" gorm:"not null"`
    Name string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 函数定义

```go
// ✅ 正确：参数和返回值清晰
func (c *UserController) CreateUser(ctx http.Context) http.Response {
    // 函数体
}

// ❌ 错误：参数过多，返回值不清晰
func (c *UserController) CreateUser(ctx http.Context, name, email, password string, age int, isActive bool) (map[string]interface{}, error) {
    // 函数体
}
```

### 错误处理

```go
// ✅ 正确：立即处理错误
func GetUser(id int) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    return user, nil
}

// ❌ 错误：忽略错误或延迟处理
func GetUser(id int) (*User, error) {
    user, _ := db.FindUser(id) // 忽略错误
    return user, nil
}
```

## 💬 注释规范

### 包注释

```go
// Package http provides HTTP server functionality for Laravel-Go Framework.
// It includes routing, middleware, and request/response handling.
package http
```

### 函数注释

```go
// CreateUser creates a new user in the system.
// It validates the input data and returns the created user or an error.
func (c *UserController) CreateUser(ctx http.Context) http.Response {
    // 函数体
}
```

### 复杂逻辑注释

```go
// 验证用户权限
// 1. 检查用户是否已登录
// 2. 验证用户是否有编辑权限
// 3. 检查资源是否属于用户
if !c.Auth().Check() {
    return c.Json(map[string]string{"error": "Unauthorized"}).Status(401)
}
```

### TODO 注释

```go
// TODO: 实现用户权限缓存机制
// 当前每次请求都查询数据库，需要优化性能
func (c *UserController) checkPermission(userID int, action string) bool {
    // 实现逻辑
}
```

## ⚠️ 错误处理

### 错误包装

```go
// ✅ 正确：使用 fmt.Errorf 包装错误
func GetUser(id int) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to find user with id %d: %w", id, err)
    }
    return user, nil
}
```

### 自定义错误

```go
// ✅ 正确：定义自定义错误类型
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// 使用自定义错误
func ValidateUser(user *User) error {
    if user.Name == "" {
        return &ValidationError{
            Field:   "name",
            Message: "name is required",
        }
    }
    return nil
}
```

### 错误检查

```go
// ✅ 正确：使用 if err != nil 检查错误
func ProcessUser(user *User) error {
    if err := ValidateUser(user); err != nil {
        return err
    }

    if err := SaveUser(user); err != nil {
        return err
    }

    return nil
}
```

## ⚡ 性能优化

### 字符串拼接

```go
// ✅ 正确：使用 strings.Builder 进行大量字符串拼接
func BuildQuery(conditions []string) string {
    var builder strings.Builder
    builder.WriteString("SELECT * FROM users WHERE ")

    for i, condition := range conditions {
        if i > 0 {
            builder.WriteString(" AND ")
        }
        builder.WriteString(condition)
    }

    return builder.String()
}

// ❌ 错误：使用 + 进行大量字符串拼接
func BuildQuery(conditions []string) string {
    query := "SELECT * FROM users WHERE "
    for i, condition := range conditions {
        if i > 0 {
            query += " AND " // 性能较差
        }
        query += condition
    }
    return query
}
```

### 切片预分配

```go
// ✅ 正确：预分配切片容量
func GetUsers(ids []int) []*User {
    users := make([]*User, 0, len(ids)) // 预分配容量
    for _, id := range ids {
        user, err := GetUser(id)
        if err == nil {
            users = append(users, user)
        }
    }
    return users
}

// ❌ 错误：不预分配容量
func GetUsers(ids []int) []*User {
    var users []*User // 没有预分配容量
    for _, id := range ids {
        user, err := GetUser(id)
        if err == nil {
            users = append(users, user)
        }
    }
    return users
}
```

### 数据库查询优化

```go
// ✅ 正确：使用批量查询
func GetUsersByIDs(ids []int) ([]*User, error) {
    var users []*User
    err := db.Where("id IN ?", ids).Find(&users).Error
    return users, err
}

// ❌ 错误：循环查询数据库
func GetUsersByIDs(ids []int) ([]*User, error) {
    var users []*User
    for _, id := range ids {
        user, err := GetUser(id) // 每次查询数据库
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}
```

## 🔒 安全规范

### 输入验证

```go
// ✅ 正确：验证所有输入
func (c *UserController) CreateUser(ctx http.Context) http.Response {
    var user User
    if err := ctx.BindJSON(&user); err != nil {
        return c.Json(map[string]string{"error": "Invalid JSON"}).Status(400)
    }

    // 验证输入
    if err := ValidateUser(&user); err != nil {
        return c.Json(map[string]string{"error": err.Error()}).Status(422)
    }

    // 处理业务逻辑
    return c.CreateUser(&user)
}

// ❌ 错误：不验证输入
func (c *UserController) CreateUser(ctx http.Context) http.Response {
    var user User
    ctx.BindJSON(&user) // 没有检查错误

    // 直接处理，可能存在安全风险
    return c.CreateUser(&user)
}
```

### SQL 注入防护

```go
// ✅ 正确：使用参数化查询
func GetUserByEmail(email string) (*User, error) {
    var user User
    err := db.Where("email = ?", email).First(&user).Error
    return &user, err
}

// ❌ 错误：直接拼接 SQL
func GetUserByEmail(email string) (*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
    // 存在 SQL 注入风险
    return db.Raw(query).Scan(&user)
}
```

### 密码处理

```go
// ✅ 正确：使用安全的密码哈希
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

// ❌ 错误：明文存储密码
func CreateUser(user *User) error {
    // 直接存储明文密码
    user.Password = user.Password // 不安全
    return db.Create(user).Error
}
```

## 🧪 测试规范

### 测试命名

```go
// ✅ 正确：使用描述性的测试名称
func TestUserController_CreateUser_Success(t *testing.T) {
    // 测试逻辑
}

func TestUserController_CreateUser_ValidationError(t *testing.T) {
    // 测试逻辑
}

// ❌ 错误：使用不清晰的测试名称
func TestCreateUser(t *testing.T) {
    // 测试逻辑
}
```

### 测试结构

```go
// ✅ 正确：使用 AAA 模式 (Arrange, Act, Assert)
func TestUserController_CreateUser(t *testing.T) {
    // Arrange - 准备测试数据
    app := framework.NewTestApplication()
    controller := &UserController{}
    userData := map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
    }

    // Act - 执行被测试的方法
    response := controller.CreateUser(mockContext(userData))

    // Assert - 验证结果
    assert.Equal(t, 201, response.Status())

    var result map[string]interface{}
    json.Unmarshal([]byte(response.Body()), &result)
    assert.Equal(t, "John Doe", result["name"])
}
```

### 测试覆盖率

```go
// ✅ 正确：确保测试覆盖所有分支
func TestValidateUser(t *testing.T) {
    tests := []struct {
        name    string
        user    *User
        wantErr bool
    }{
        {
            name: "valid user",
            user: &User{Name: "John", Email: "john@example.com"},
            wantErr: false,
        },
        {
            name: "missing name",
            user: &User{Email: "john@example.com"},
            wantErr: true,
        },
        {
            name: "invalid email",
            user: &User{Name: "John", Email: "invalid-email"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateUser(tt.user)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Mock 使用

```go
// ✅ 正确：使用 Mock 进行单元测试
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
    // 创建 Mock
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    // 设置期望
    user := &User{Name: "John", Email: "john@example.com"}
    mockRepo.On("Create", user).Return(nil)

    // 执行测试
    err := service.CreateUser(user)

    // 验证结果
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

## 📚 工具和检查

### 代码格式化

```bash
# 格式化代码
go fmt ./...

# 使用 goimports 自动管理导入
goimports -w .

# 使用 golangci-lint 进行代码检查
golangci-lint run
```

### 预提交钩子

```bash
#!/bin/sh
# .git/hooks/pre-commit

# 运行测试
go test ./...

# 检查代码格式
go fmt ./...

# 运行 linter
golangci-lint run

# 检查测试覆盖率
go test -cover ./...
```

### IDE 配置

#### VS Code 设置

```json
{
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.testFlags": ["-v", "-cover"],
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  }
}
```

## 📖 参考资源

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Best Practices](https://github.com/golang/go/wiki/BestPractices)

---

遵循这些代码规范将帮助你编写高质量、可维护的代码！ 🚀
