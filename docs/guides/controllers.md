# 控制器指南

## 🎮 控制器概览

控制器是 Laravel-Go Framework 中处理 HTTP 请求的核心组件，它负责接收请求、处理业务逻辑并返回响应。

## 🚀 快速开始

### 基本控制器结构

```go
// app/Http/Controllers/UserController.go
package controllers

import (
    "laravel-go/framework/http"
    "laravel-go/app/Models"
)

type UserController struct {
    http.Controller
}

// Index 方法 - 获取用户列表
func (c *UserController) Index() http.Response {
    users := []Models.User{
        {ID: 1, Name: "John Doe", Email: "john@example.com"},
        {ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
    }
    
    return c.Json(map[string]interface{}{
        "users": users,
        "total": len(users),
    })
}

// Show 方法 - 获取单个用户
func (c *UserController) Show(id string) http.Response {
    // 这里应该从数据库获取用户
    user := Models.User{
        ID:    1,
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    return c.Json(user)
}

// Store 方法 - 创建用户
func (c *UserController) Store(request http.Request) http.Response {
    // 获取请求数据
    data := request.Body
    
    // 创建用户逻辑
    user := Models.User{
        Name:  data["name"].(string),
        Email: data["email"].(string),
    }
    
    return c.Json(user).Status(201)
}

// Update 方法 - 更新用户
func (c *UserController) Update(id string, request http.Request) http.Response {
    // 更新用户逻辑
    return c.Json(map[string]string{
        "message": "User updated successfully",
    })
}

// Delete 方法 - 删除用户
func (c *UserController) Delete(id string) http.Response {
    // 删除用户逻辑
    return c.Json(map[string]string{
        "message": "User deleted successfully",
    })
}
```

### 路由绑定

```go
// routes/web.go
func WebRoutes(router *routing.Router) {
    // 资源路由
    router.Resource("/users", &UserController{})
    
    // 或者手动定义
    router.Get("/users", &UserController{}, "Index")
    router.Get("/users/{id}", &UserController{}, "Show")
    router.Post("/users", &UserController{}, "Store")
    router.Put("/users/{id}", &UserController{}, "Update")
    router.Delete("/users/{id}", &UserController{}, "Delete")
}
```

## 📋 控制器方法

### 1. Index - 列表方法

```go
func (c *UserController) Index(request http.Request) http.Response {
    // 获取查询参数
    page := request.Query["page"]
    limit := request.Query["limit"]
    
    // 分页逻辑
    users, total := c.userService.GetUsers(page, limit)
    
    return c.Json(map[string]interface{}{
        "data":  users,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}
```

### 2. Show - 详情方法

```go
func (c *UserController) Show(id string) http.Response {
    // 参数验证
    if id == "" {
        return c.JsonError("User ID is required", 400)
    }
    
    // 获取用户
    user, err := c.userService.GetUser(id)
    if err != nil {
        return c.JsonError("User not found", 404)
    }
    
    return c.Json(user)
}
```

### 3. Store - 创建方法

```go
func (c *UserController) Store(request http.Request) http.Response {
    // 数据验证
    data := request.Body
    if err := c.validateUserData(data); err != nil {
        return c.JsonError(err.Error(), 422)
    }
    
    // 创建用户
    user, err := c.userService.CreateUser(data)
    if err != nil {
        return c.JsonError("Failed to create user", 500)
    }
    
    return c.Json(user).Status(201)
}
```

### 4. Update - 更新方法

```go
func (c *UserController) Update(id string, request http.Request) http.Response {
    // 数据验证
    data := request.Body
    if err := c.validateUserData(data); err != nil {
        return c.JsonError(err.Error(), 422)
    }
    
    // 更新用户
    user, err := c.userService.UpdateUser(id, data)
    if err != nil {
        return c.JsonError("Failed to update user", 500)
    }
    
    return c.Json(user)
}
```

### 5. Delete - 删除方法

```go
func (c *UserController) Delete(id string) http.Response {
    // 删除用户
    err := c.userService.DeleteUser(id)
    if err != nil {
        return c.JsonError("Failed to delete user", 500)
    }
    
    return c.Json(map[string]string{
        "message": "User deleted successfully",
    })
}
```

## 🔧 控制器功能

### 依赖注入

```go
type UserController struct {
    http.Controller
    userService *Services.UserService
    authService *Services.AuthService
}

func NewUserController(
    userService *Services.UserService,
    authService *Services.AuthService,
) *UserController {
    return &UserController{
        userService: userService,
        authService: authService,
    }
}

// 使用依赖注入
func (c *UserController) Index() http.Response {
    users := c.userService.GetAllUsers()
    return c.Json(users)
}
```

### 请求处理

```go
func (c *UserController) Store(request http.Request) http.Response {
    // 获取请求数据
    data := request.Body
    
    // 获取查询参数
    query := request.Query
    
    // 获取请求头
    headers := request.Headers
    
    // 获取文件上传
    files := request.Files
    
    // 获取认证用户
    user := request.Context["user"]
    
    // 处理逻辑...
    return c.Json(data)
}
```

### 响应处理

```go
func (c *UserController) Show(id string) http.Response {
    user := c.userService.GetUser(id)
    
    // JSON 响应
    return c.Json(user)
    
    // 带状态码的响应
    return c.Json(user).Status(200)
    
    // 带头的响应
    return c.Json(user).Header("X-Custom-Header", "value")
    
    // 错误响应
    return c.JsonError("User not found", 404)
    
    // 重定向
    return c.Redirect("/users")
    
    // 视图响应
    return c.View("users.show", map[string]interface{}{
        "user": user,
    })
}
```

## 🛠️ 高级功能

### 1. 资源控制器

```go
// 完整的资源控制器
type PostController struct {
    http.Controller
    postService *Services.PostService
}

// 资源方法
func (c *PostController) Index() http.Response {
    posts := c.postService.GetAllPosts()
    return c.Json(posts)
}

func (c *PostController) Create() http.Response {
    // 显示创建表单
    return c.View("posts.create")
}

func (c *PostController) Store(request http.Request) http.Response {
    // 创建文章
    post := c.postService.CreatePost(request.Body)
    return c.Json(post).Status(201)
}

func (c *PostController) Show(id string) http.Response {
    post := c.postService.GetPost(id)
    return c.Json(post)
}

func (c *PostController) Edit(id string) http.Response {
    // 显示编辑表单
    post := c.postService.GetPost(id)
    return c.View("posts.edit", map[string]interface{}{
        "post": post,
    })
}

func (c *PostController) Update(id string, request http.Request) http.Response {
    post := c.postService.UpdatePost(id, request.Body)
    return c.Json(post)
}

func (c *PostController) Delete(id string) http.Response {
    c.postService.DeletePost(id)
    return c.Json(map[string]string{"message": "Post deleted"})
}
```

### 2. API 资源控制器

```go
// API 资源控制器
type ApiUserController struct {
    http.Controller
    userService *Services.UserService
}

func (c *ApiUserController) Index() http.Response {
    users := c.userService.GetAllUsers()
    
    // 使用 API 资源
    return c.JsonResource(users, &Resources.UserResource{})
}

func (c *ApiUserController) Show(id string) http.Response {
    user := c.userService.GetUser(id)
    
    // 使用单个 API 资源
    return c.JsonResource(user, &Resources.UserResource{})
}

func (c *ApiUserController) Store(request http.Request) http.Response {
    user := c.userService.CreateUser(request.Body)
    
    // 使用 API 资源
    return c.JsonResource(user, &Resources.UserResource{}).Status(201)
}
```

### 3. 表单请求验证

```go
// 表单请求类
type CreateUserRequest struct {
    http.FormRequest
}

func (r *CreateUserRequest) Rules() map[string]string {
    return map[string]string{
        "name":  "required|string|max:255",
        "email": "required|email|unique:users,email",
        "password": "required|string|min:8",
    }
}

func (r *CreateUserRequest) Messages() map[string]string {
    return map[string]string{
        "name.required": "Name is required",
        "email.required": "Email is required",
        "email.email": "Email must be valid",
        "password.required": "Password is required",
        "password.min": "Password must be at least 8 characters",
    }
}

// 在控制器中使用
func (c *UserController) Store(request http.Request) http.Response {
    formRequest := &CreateUserRequest{}
    
    if err := formRequest.Validate(request); err != nil {
        return c.JsonError(err.Error(), 422)
    }
    
    // 创建用户
    user := c.userService.CreateUser(formRequest.Validated())
    return c.Json(user).Status(201)
}
```

## 🎯 控制器最佳实践

### 1. 单一职责

```go
// ✅ 好的做法：控制器只处理 HTTP 请求
type UserController struct {
    http.Controller
    userService *Services.UserService
}

func (c *UserController) Store(request http.Request) http.Response {
    // 验证数据
    if err := c.validate(request.Body); err != nil {
        return c.JsonError(err.Error(), 422)
    }
    
    // 调用服务层
    user, err := c.userService.CreateUser(request.Body)
    if err != nil {
        return c.JsonError("Failed to create user", 500)
    }
    
    return c.Json(user).Status(201)
}

// ❌ 不好的做法：控制器包含业务逻辑
func (c *UserController) Store(request http.Request) http.Response {
    // 直接在控制器中处理业务逻辑
    // 数据库操作、业务规则等
}
```

### 2. 错误处理

```go
func (c *UserController) Show(id string) http.Response {
    user, err := c.userService.GetUser(id)
    if err != nil {
        // 根据错误类型返回不同的状态码
        switch err.(type) {
        case *errors.NotFoundError:
            return c.JsonError("User not found", 404)
        case *errors.ValidationError:
            return c.JsonError(err.Error(), 422)
        default:
            return c.JsonError("Internal server error", 500)
        }
    }
    
    return c.Json(user)
}
```

### 3. 参数验证

```go
func (c *UserController) Show(id string) http.Response {
    // 验证 ID 参数
    if id == "" || !c.isValidID(id) {
        return c.JsonError("Invalid user ID", 400)
    }
    
    user := c.userService.GetUser(id)
    return c.Json(user)
}

func (c *UserController) isValidID(id string) bool {
    // ID 验证逻辑
    return len(id) > 0 && id != "0"
}
```

### 4. 响应格式化

```go
func (c *UserController) Index() http.Response {
    users := c.userService.GetAllUsers()
    
    // 统一的响应格式
    return c.Json(map[string]interface{}{
        "success": true,
        "data":    users,
        "message": "Users retrieved successfully",
    })
}

func (c *UserController) Store(request http.Request) http.Response {
    user, err := c.userService.CreateUser(request.Body)
    if err != nil {
        return c.Json(map[string]interface{}{
            "success": false,
            "message": err.Error(),
        }).Status(422)
    }
    
    return c.Json(map[string]interface{}{
        "success": true,
        "data":    user,
        "message": "User created successfully",
    }).Status(201)
}
```

## 🔍 控制器测试

### 单元测试

```go
// tests/controllers/user_controller_test.go
package tests

import (
    "testing"
    "laravel-go/framework/http"
    "laravel-go/app/Http/Controllers"
    "laravel-go/app/Services"
)

func TestUserController_Index(t *testing.T) {
    // 创建 mock 服务
    mockUserService := &Services.MockUserService{}
    
    // 创建控制器
    controller := &Controllers.UserController{
        userService: mockUserService,
    }
    
    // 设置 mock 期望
    mockUserService.On("GetAllUsers").Return([]Models.User{
        {ID: 1, Name: "John Doe"},
        {ID: 2, Name: "Jane Smith"},
    })
    
    // 执行测试
    response := controller.Index()
    
    // 验证结果
    if response.StatusCode != 200 {
        t.Errorf("Expected status 200, got %d", response.StatusCode)
    }
    
    // 验证 mock 调用
    mockUserService.AssertExpectations(t)
}
```

### 集成测试

```go
func TestUserController_Store_Integration(t *testing.T) {
    // 设置测试数据库
    db := setupTestDatabase()
    
    // 创建服务
    userService := &Services.UserService{DB: db}
    
    // 创建控制器
    controller := &Controllers.UserController{
        userService: userService,
    }
    
    // 创建请求
    request := http.Request{
        Body: map[string]interface{}{
            "name":  "John Doe",
            "email": "john@example.com",
        },
    }
    
    // 执行测试
    response := controller.Store(request)
    
    // 验证结果
    if response.StatusCode != 201 {
        t.Errorf("Expected status 201, got %d", response.StatusCode)
    }
    
    // 验证数据库中的数据
    var user Models.User
    db.First(&user, "email = ?", "john@example.com")
    if user.Name != "John Doe" {
        t.Errorf("Expected name John Doe, got %s", user.Name)
    }
}
```

## 📊 控制器性能优化

### 1. 缓存响应

```go
func (c *UserController) Index() http.Response {
    // 检查缓存
    cacheKey := "users:list"
    if cached, found := cache.Get(cacheKey); found {
        return c.Json(cached)
    }
    
    // 获取数据
    users := c.userService.GetAllUsers()
    
    // 缓存结果
    cache.Set(cacheKey, users, time.Hour)
    
    return c.Json(users)
}
```

### 2. 数据库查询优化

```go
func (c *UserController) Index() http.Response {
    // 使用分页
    page := c.getPageParam()
    limit := c.getLimitParam()
    
    // 使用预加载
    users := c.userService.GetUsersWithRelations(page, limit, []string{"posts", "comments"})
    
    return c.Json(users)
}
```

### 3. 异步处理

```go
func (c *UserController) Store(request http.Request) http.Response {
    // 创建用户
    user := c.userService.CreateUser(request.Body)
    
    // 异步发送欢迎邮件
    go c.sendWelcomeEmail(user)
    
    return c.Json(user).Status(201)
}

func (c *UserController) sendWelcomeEmail(user *Models.User) {
    // 发送邮件的逻辑
}
```

## 📝 总结

Laravel-Go Framework 的控制器系统提供了：

1. **简洁性**: 清晰的方法结构和命名约定
2. **灵活性**: 支持多种响应类型和状态码
3. **可测试性**: 易于编写单元测试和集成测试
4. **可维护性**: 遵循单一职责原则
5. **性能优化**: 支持缓存和异步处理

通过合理使用控制器的各种功能，可以构建出高效、可维护的 Web 应用程序。 