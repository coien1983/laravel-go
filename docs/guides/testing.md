# 测试指南

## 📖 概述

Laravel-Go Framework 提供了完整的测试支持，包括单元测试、集成测试、功能测试和测试工具，帮助确保代码质量和应用程序的可靠性。

## 🚀 快速开始

### 1. 基本测试

```go
// 用户服务测试
package tests

import (
    "testing"
    "laravel-go/app/Services"
    "laravel-go/app/Models"
)

// 用户服务测试
func TestUserService_CreateUser(t *testing.T) {
    // 设置测试环境
    db := setupTestDatabase()
    userService := Services.NewUserService(db)

    // 测试数据
    userData := map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "password123",
    }

    // 执行测试
    user, err := userService.CreateUser(userData)

    // 断言结果
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    if user == nil {
        t.Error("Expected user to be created")
    }

    if user.Name != userData["name"] {
        t.Errorf("Expected name %s, got %s", userData["name"], user.Name)
    }

    if user.Email != userData["email"] {
        t.Errorf("Expected email %s, got %s", userData["email"], user.Email)
    }

    // 清理测试数据
    cleanupTestData(db)
}

// 设置测试数据库
func setupTestDatabase() *database.Connection {
    // 使用测试数据库配置
    config.Set("database.default", "sqlite")
    config.Set("database.sqlite.database", ":memory:")

    db := database.NewConnection()

    // 运行迁移
    db.AutoMigrate(&Models.User{})

    return db
}

// 清理测试数据
func cleanupTestData(db *database.Connection) {
    db.Exec("DELETE FROM users")
}
```

### 2. 测试套件

```go
// 测试套件
type UserServiceTestSuite struct {
    testing.TestSuite
    db          *database.Connection
    userService *Services.UserService
}

// 设置测试套件
func (suite *UserServiceTestSuite) SetupSuite() {
    suite.db = setupTestDatabase()
    suite.userService = Services.NewUserService(suite.db)
}

// 每个测试前的设置
func (suite *UserServiceTestSuite) SetupTest() {
    // 清理数据
    suite.db.Exec("DELETE FROM users")
}

// 每个测试后的清理
func (suite *UserServiceTestSuite) TearDownTest() {
    // 清理数据
    suite.db.Exec("DELETE FROM users")
}

// 测试套件清理
func (suite *UserServiceTestSuite) TearDownSuite() {
    suite.db.Close()
}

// 测试创建用户
func (suite *UserServiceTestSuite) TestCreateUser() {
    userData := map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "password123",
    }

    user, err := suite.userService.CreateUser(userData)

    suite.NoError(err)
    suite.NotNil(user)
    suite.Equal(userData["name"], user.Name)
    suite.Equal(userData["email"], user.Email)
}

// 测试创建重复邮箱用户
func (suite *UserServiceTestSuite) TestCreateUserWithDuplicateEmail() {
    userData := map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "password123",
    }

    // 创建第一个用户
    _, err := suite.userService.CreateUser(userData)
    suite.NoError(err)

    // 尝试创建第二个相同邮箱的用户
    _, err = suite.userService.CreateUser(userData)
    suite.Error(err)
    suite.Contains(err.Error(), "email already exists")
}

// 运行测试套件
func TestUserServiceSuite(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}
```

## 🔧 测试类型

### 1. 单元测试

```go
// 单元测试 - 用户模型
func TestUser_SetPassword(t *testing.T) {
    user := &Models.User{
        Name:  "John Doe",
        Email: "john@example.com",
    }

    password := "password123"
    user.SetPassword(password)

    // 验证密码已哈希
    if user.Password == password {
        t.Error("Password should be hashed")
    }

    // 验证密码检查
    if !user.CheckPassword(password) {
        t.Error("Password check should pass")
    }

    if user.CheckPassword("wrongpassword") {
        t.Error("Wrong password should fail")
    }
}

// 单元测试 - 验证器
func TestValidator_Email(t *testing.T) {
    validator := validation.New()

    // 测试有效邮箱
    validEmails := []string{
        "test@example.com",
        "user.name@domain.co.uk",
        "user+tag@example.org",
    }

    for _, email := range validEmails {
        if !validator.Email("email", email, "") {
            t.Errorf("Email %s should be valid", email)
        }
    }

    // 测试无效邮箱
    invalidEmails := []string{
        "invalid-email",
        "@example.com",
        "user@",
        "user@.com",
    }

    for _, email := range invalidEmails {
        if validator.Email("email", email, "") {
            t.Errorf("Email %s should be invalid", email)
        }
    }
}

// 单元测试 - 工具函数
func TestGenerateSlug(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"Hello World", "hello-world"},
        {"Test Title 123", "test-title-123"},
        {"Special@#$%^&*()", "special"},
        {"Multiple   Spaces", "multiple-spaces"},
        {"", ""},
    }

    for _, test := range tests {
        result := generateSlug(test.input)
        if result != test.expected {
            t.Errorf("generateSlug(%s) = %s, expected %s", test.input, result, test.expected)
        }
    }
}
```

### 2. 集成测试

```go
// 集成测试 - 用户控制器
func TestUserController_Store(t *testing.T) {
    // 设置测试环境
    db := setupTestDatabase()
    userService := Services.NewUserService(db)
    controller := Controllers.NewUserController(userService)

    // 创建测试请求
    requestBody := map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "password123",
    }

    request := http.Request{
        Method: "POST",
        Path:   "/users",
        Body:   requestBody,
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }

    // 执行请求
    response := controller.Store(request)

    // 验证响应
    if response.StatusCode != 201 {
        t.Errorf("Expected status 201, got %d", response.StatusCode)
    }

    // 验证响应体
    var responseData map[string]interface{}
    if err := json.Unmarshal([]byte(response.Body), &responseData); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if responseData["name"] != requestBody["name"] {
        t.Errorf("Expected name %s, got %s", requestBody["name"], responseData["name"])
    }

    // 验证数据库中的数据
    var user Models.User
    if err := db.Where("email = ?", requestBody["email"]).First(&user).Error; err != nil {
        t.Errorf("User not found in database: %v", err)
    }

    cleanupTestData(db)
}

// 集成测试 - 认证流程
func TestAuthFlow(t *testing.T) {
    db := setupTestDatabase()
    authService := Services.NewAuthService(db)
    authController := Controllers.NewAuthController(authService)

    // 1. 注册用户
    registerRequest := http.Request{
        Method: "POST",
        Path:   "/auth/register",
        Body: map[string]interface{}{
            "name":     "John Doe",
            "email":    "john@example.com",
            "password": "password123",
        },
    }

    registerResponse := authController.Register(registerRequest)

    if registerResponse.StatusCode != 201 {
        t.Errorf("Registration failed: %d", registerResponse.StatusCode)
    }

    // 2. 登录用户
    loginRequest := http.Request{
        Method: "POST",
        Path:   "/auth/login",
        Body: map[string]interface{}{
            "email":    "john@example.com",
            "password": "password123",
        },
    }

    loginResponse := authController.Login(loginRequest)

    if loginResponse.StatusCode != 200 {
        t.Errorf("Login failed: %d", loginResponse.StatusCode)
    }

    // 3. 验证令牌
    var loginData map[string]interface{}
    json.Unmarshal([]byte(loginResponse.Body), &loginData)

    token := loginData["token"].(string)

    // 使用令牌访问受保护的资源
    protectedRequest := http.Request{
        Method: "GET",
        Path:   "/api/profile",
        Headers: map[string]string{
            "Authorization": "Bearer " + token,
        },
    }

    // 这里需要设置认证中间件
    // profileResponse := profileController.Show(protectedRequest)

    cleanupTestData(db)
}
```

### 3. 功能测试

```go
// 功能测试 - HTTP 服务器
func TestHTTPServer_UserEndpoints(t *testing.T) {
    // 启动测试服务器
    server := setupTestServer()
    defer server.Close()

    // 测试创建用户
    createUserData := map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "password123",
    }

    createResp, err := http.PostJSON(server.URL+"/api/users", createUserData)
    if err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }

    if createResp.StatusCode != 201 {
        t.Errorf("Expected status 201, got %d", createResp.StatusCode)
    }

    // 测试获取用户列表
    listResp, err := http.Get(server.URL + "/api/users")
    if err != nil {
        t.Fatalf("Failed to get users: %v", err)
    }

    if listResp.StatusCode != 200 {
        t.Errorf("Expected status 200, got %d", listResp.StatusCode)
    }

    var users []map[string]interface{}
    if err := json.Unmarshal(listResp.Body, &users); err != nil {
        t.Fatalf("Failed to parse users: %v", err)
    }

    if len(users) == 0 {
        t.Error("Expected at least one user")
    }
}

// 设置测试服务器
func setupTestServer() *httptest.Server {
    // 设置测试数据库
    db := setupTestDatabase()

    // 创建服务
    userService := Services.NewUserService(db)
    userController := Controllers.NewUserController(userService)

    // 设置路由
    router := routing.NewRouter()
    router.Post("/api/users", userController.Store)
    router.Get("/api/users", userController.Index)

    // 创建服务器
    server := httptest.NewServer(router)

    return server
}
```

## 🛠️ 测试工具

### 1. 测试辅助函数

```go
// 测试辅助函数
package testhelpers

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "laravel-go/framework/database"
    "laravel-go/framework/config"
)

// 创建测试数据库连接
func CreateTestDatabase() *database.Connection {
    config.Set("database.default", "sqlite")
    config.Set("database.sqlite.database", ":memory:")

    db := database.NewConnection()
    return db
}

// 创建测试 HTTP 请求
func CreateTestRequest(method, path string, body interface{}) *http.Request {
    var reqBody []byte
    if body != nil {
        reqBody, _ = json.Marshal(body)
    }

    req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    return req
}

// 创建测试 HTTP 响应记录器
func CreateTestResponseRecorder() *httptest.ResponseRecorder {
    return httptest.NewRecorder()
}

// 断言 HTTP 状态码
func AssertStatusCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected status code %d, got %d", expected, actual)
    }
}

// 断言响应体包含字段
func AssertResponseContains(t *testing.T, response *httptest.ResponseRecorder, field string, expected interface{}) {
    var responseData map[string]interface{}
    if err := json.Unmarshal(response.Body.Bytes(), &responseData); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if responseData[field] != expected {
        t.Errorf("Expected %s to be %v, got %v", field, expected, responseData[field])
    }
}

// 创建测试用户
func CreateTestUser(db *database.Connection, name, email string) *Models.User {
    user := &Models.User{
        Name:  name,
        Email: email,
    }
    user.SetPassword("password123")

    db.Create(user)
    return user
}

// 清理测试数据
func CleanupTestData(db *database.Connection) {
    db.Exec("DELETE FROM users")
    db.Exec("DELETE FROM posts")
    db.Exec("DELETE FROM comments")
}
```

### 2. 模拟和存根

```go
// 模拟用户服务
type MockUserService struct {
    users map[uint]*Models.User
    nextID uint
}

func NewMockUserService() *MockUserService {
    return &MockUserService{
        users: make(map[uint]*Models.User),
        nextID: 1,
    }
}

func (m *MockUserService) CreateUser(data map[string]interface{}) (*Models.User, error) {
    user := &Models.User{
        ID:    m.nextID,
        Name:  data["name"].(string),
        Email: data["email"].(string),
    }

    m.users[user.ID] = user
    m.nextID++

    return user, nil
}

func (m *MockUserService) GetUser(id uint) (*Models.User, error) {
    user, exists := m.users[id]
    if !exists {
        return nil, errors.New("user not found")
    }

    return user, nil
}

func (m *MockUserService) GetUsers(page, limit int) ([]*Models.User, int64) {
    users := make([]*Models.User, 0)
    for _, user := range m.users {
        users = append(users, user)
    }

    return users, int64(len(users))
}

// 使用模拟服务进行测试
func TestUserControllerWithMock(t *testing.T) {
    mockService := NewMockUserService()
    controller := Controllers.NewUserController(mockService)

    request := http.Request{
        Method: "POST",
        Path:   "/users",
        Body: map[string]interface{}{
            "name":     "John Doe",
            "email":    "john@example.com",
            "password": "password123",
        },
    }

    response := controller.Store(request)

    if response.StatusCode != 201 {
        t.Errorf("Expected status 201, got %d", response.StatusCode)
    }

    // 验证模拟服务中的数据
    users, _ := mockService.GetUsers(1, 10)
    if len(users) != 1 {
        t.Error("Expected one user to be created")
    }
}
```

### 3. 测试数据工厂

```go
// 测试数据工厂
type UserFactory struct{}

func (f *UserFactory) Create(attributes map[string]interface{}) *Models.User {
    user := &Models.User{
        Name:  f.getAttribute(attributes, "name", "Test User"),
        Email: f.getAttribute(attributes, "email", "test@example.com"),
    }
    user.SetPassword(f.getAttribute(attributes, "password", "password123"))

    return user
}

func (f *UserFactory) CreateMany(count int, attributes map[string]interface{}) []*Models.User {
    users := make([]*Models.User, count)
    for i := 0; i < count; i++ {
        users[i] = f.Create(attributes)
    }
    return users
}

func (f *UserFactory) getAttribute(attributes map[string]interface{}, key, defaultValue string) string {
    if value, exists := attributes[key]; exists {
        return value.(string)
    }
    return defaultValue
}

// 使用工厂创建测试数据
func TestWithFactory(t *testing.T) {
    factory := &UserFactory{}

    // 创建单个用户
    user := factory.Create(map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
    })

    if user.Name != "John Doe" {
        t.Error("Factory should create user with specified name")
    }

    // 创建多个用户
    users := factory.CreateMany(3, map[string]interface{}{
        "name": "Test User",
    })

    if len(users) != 3 {
        t.Error("Factory should create specified number of users")
    }
}
```

## 📊 测试覆盖率

### 1. 覆盖率测试

```go
// 运行覆盖率测试
func TestMain(m *testing.M) {
    // 设置测试环境
    setupTestEnvironment()

    // 运行测试
    code := m.Run()

    // 清理测试环境
    cleanupTestEnvironment()

    os.Exit(code)
}

// 生成覆盖率报告
func TestCoverage(t *testing.T) {
    // 这个测试用于生成覆盖率报告
    // 实际覆盖率测试应该在命令行运行：
    // go test -coverprofile=coverage.out ./...
    // go tool cover -html=coverage.out -o coverage.html
}

// 覆盖率检查
func TestCoverageThreshold(t *testing.T) {
    // 检查覆盖率是否达到阈值
    coverage := getCoverage()
    threshold := 80.0 // 80% 覆盖率阈值

    if coverage < threshold {
        t.Errorf("Coverage %.2f%% is below threshold %.2f%%", coverage, threshold)
    }
}

func getCoverage() float64 {
    // 这里应该实现获取实际覆盖率的逻辑
    // 可以通过解析覆盖率文件来实现
    return 85.0 // 示例值
}
```

### 2. 性能测试

```go
// 性能测试
func BenchmarkUserService_CreateUser(b *testing.B) {
    db := setupTestDatabase()
    userService := Services.NewUserService(db)

    userData := map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "password123",
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        userData["email"] = fmt.Sprintf("user%d@example.com", i)
        _, err := userService.CreateUser(userData)
        if err != nil {
            b.Fatalf("Failed to create user: %v", err)
        }
    }
}

// 内存使用测试
func BenchmarkUserService_GetUsers(b *testing.B) {
    db := setupTestDatabase()
    userService := Services.NewUserService(db)

    // 创建测试数据
    for i := 0; i < 1000; i++ {
        userData := map[string]interface{}{
            "name":     fmt.Sprintf("User %d", i),
            "email":    fmt.Sprintf("user%d@example.com", i),
            "password": "password123",
        }
        userService.CreateUser(userData)
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _, _ = userService.GetUsers(1, 10)
    }
}
```

## 📚 总结

Laravel-Go Framework 的测试系统提供了：

1. **单元测试**: 测试单个函数和方法
2. **集成测试**: 测试组件间的交互
3. **功能测试**: 测试完整的业务流程
4. **测试工具**: 辅助函数和模拟对象
5. **测试数据工厂**: 创建测试数据
6. **覆盖率测试**: 确保代码覆盖率
7. **性能测试**: 测试性能和内存使用

通过合理使用测试系统，可以确保代码质量和应用程序的可靠性。
