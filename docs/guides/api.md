# API 开发指南

## 📖 概述

Laravel-Go Framework 提供了完整的 API 开发支持，包括 RESTful API 设计、API 资源、版本控制、认证授权和文档生成等功能。

## 🚀 快速开始

### 1. 基本 API 控制器

```go
// API 控制器
type UserApiController struct {
    http.Controller
    userService *UserService
}

func NewUserApiController(userService *UserService) *UserApiController {
    return &UserApiController{
        userService: userService,
    }
}

// 获取用户列表
func (c *UserApiController) Index(request http.Request) http.Response {
    page := c.getPageParam(request)
    limit := c.getLimitParam(request)

    users, total := c.userService.GetUsers(page, limit)

    return c.Json(map[string]interface{}{
        "data":  users,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// 获取单个用户
func (c *UserApiController) Show(id string) http.Response {
    userID, err := strconv.Atoi(id)
    if err != nil {
        return c.JsonError("Invalid user ID", 400)
    }

    user, err := c.userService.GetUser(uint(userID))
    if err != nil {
        return c.JsonError("User not found", 404)
    }

    return c.Json(user)
}

// 创建用户
func (c *UserApiController) Store(request http.Request) http.Response {
    var userRequest UserRequest

    if err := request.Bind(&userRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    if err := validator.Validate(&userRequest); err != nil {
        return c.JsonError(err.Error(), 422)
    }

    user, err := c.userService.CreateUser(userRequest)
    if err != nil {
        return c.JsonError("Failed to create user", 500)
    }

    return c.Json(user).Status(201)
}

// 更新用户
func (c *UserApiController) Update(id string, request http.Request) http.Response {
    userID, err := strconv.Atoi(id)
    if err != nil {
        return c.JsonError("Invalid user ID", 400)
    }

    var userRequest UserUpdateRequest

    if err := request.Bind(&userRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    if err := validator.Validate(&userRequest); err != nil {
        return c.JsonError(err.Error(), 422)
    }

    user, err := c.userService.UpdateUser(uint(userID), userRequest)
    if err != nil {
        return c.JsonError("Failed to update user", 500)
    }

    return c.Json(user)
}

// 删除用户
func (c *UserApiController) Delete(id string) http.Response {
    userID, err := strconv.Atoi(id)
    if err != nil {
        return c.JsonError("Invalid user ID", 400)
    }

    err = c.userService.DeleteUser(uint(userID))
    if err != nil {
        return c.JsonError("Failed to delete user", 500)
    }

    return c.Json(map[string]string{
        "message": "User deleted successfully",
    })
}
```

### 2. API 路由配置

```go
// API 路由配置
func RegisterApiRoutes(router *routing.Router) {
    // API 版本控制
    api := router.Group("/api/v1")

    // 公开 API
    api.Post("/auth/login", &AuthApiController{}, "Login")
    api.Post("/auth/register", &AuthApiController{}, "Register")
    api.Post("/auth/refresh", &AuthApiController{}, "Refresh")

    // 需要认证的 API
    authenticated := api.Group("")
    authenticated.Use(&middleware.AuthMiddleware{})

    // 用户相关 API
    authenticated.Get("/users", &UserApiController{}, "Index")
    authenticated.Get("/users/{id}", &UserApiController{}, "Show")
    authenticated.Post("/users", &UserApiController{}, "Store")
    authenticated.Put("/users/{id}", &UserApiController{}, "Update")
    authenticated.Delete("/users/{id}", &UserApiController{}, "Delete")

    // 文章相关 API
    authenticated.Get("/posts", &PostApiController{}, "Index")
    authenticated.Get("/posts/{id}", &PostApiController{}, "Show")
    authenticated.Post("/posts", &PostApiController{}, "Store")
    authenticated.Put("/posts/{id}", &PostApiController{}, "Update")
    authenticated.Delete("/posts/{id}", &PostApiController{}, "Delete")

    // 评论相关 API
    authenticated.Get("/posts/{id}/comments", &CommentApiController{}, "Index")
    authenticated.Post("/posts/{id}/comments", &CommentApiController{}, "Store")
    authenticated.Put("/comments/{id}", &CommentApiController{}, "Update")
    authenticated.Delete("/comments/{id}", &CommentApiController{}, "Delete")
}
```

## 📊 API 资源

### 1. API 资源类

```go
// 用户资源
type UserResource struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Avatar    string    `json:"avatar"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 用户集合资源
type UserCollection struct {
    Data  []UserResource `json:"data"`
    Total int64          `json:"total"`
    Page  int            `json:"page"`
    Limit int            `json:"limit"`
}

// 创建用户资源
func NewUserResource(user *User) UserResource {
    return UserResource{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Avatar:    user.Avatar,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
}

// 创建用户集合资源
func NewUserCollection(users []*User, total int64, page, limit int) UserCollection {
    resources := make([]UserResource, len(users))
    for i, user := range users {
        resources[i] = NewUserResource(user)
    }

    return UserCollection{
        Data:  resources,
        Total: total,
        Page:  page,
        Limit: limit,
    }
}

// 在控制器中使用资源
func (c *UserApiController) Index(request http.Request) http.Response {
    page := c.getPageParam(request)
    limit := c.getLimitParam(request)

    users, total := c.userService.GetUsers(page, limit)

    collection := NewUserCollection(users, total, page, limit)
    return c.Json(collection)
}

func (c *UserApiController) Show(id string) http.Response {
    userID, err := strconv.Atoi(id)
    if err != nil {
        return c.JsonError("Invalid user ID", 400)
    }

    user, err := c.userService.GetUser(uint(userID))
    if err != nil {
        return c.JsonError("User not found", 404)
    }

    resource := NewUserResource(user)
    return c.Json(resource)
}
```

### 2. 条件资源

```go
// 条件用户资源
type UserResource struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email,omitempty"` // 条件显示
    Avatar    string    `json:"avatar,omitempty"`
    Posts     []PostResource `json:"posts,omitempty"` // 关联数据
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 创建条件资源
func NewUserResource(user *User, includeEmail bool, includePosts bool) UserResource {
    resource := UserResource{
        ID:        user.ID,
        Name:      user.Name,
        Avatar:    user.Avatar,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }

    if includeEmail {
        resource.Email = user.Email
    }

    if includePosts && len(user.Posts) > 0 {
        posts := make([]PostResource, len(user.Posts))
        for i, post := range user.Posts {
            posts[i] = NewPostResource(post, false, false)
        }
        resource.Posts = posts
    }

    return resource
}

// 在控制器中使用条件资源
func (c *UserApiController) Show(id string, request http.Request) http.Response {
    userID, err := strconv.Atoi(id)
    if err != nil {
        return c.JsonError("Invalid user ID", 400)
    }

    user, err := c.userService.GetUserWithPosts(uint(userID))
    if err != nil {
        return c.JsonError("User not found", 404)
    }

    // 根据查询参数决定包含哪些字段
    includeEmail := request.Query["include_email"] == "true"
    includePosts := request.Query["include_posts"] == "true"

    resource := NewUserResource(user, includeEmail, includePosts)
    return c.Json(resource)
}
```

## 🔐 API 认证

### 1. JWT 认证

```go
// JWT 认证控制器
type AuthApiController struct {
    http.Controller
    authService *AuthService
}

func NewAuthApiController(authService *AuthService) *AuthApiController {
    return &AuthApiController{
        authService: authService,
    }
}

// 用户登录
func (c *AuthApiController) Login(request http.Request) http.Response {
    var loginRequest LoginRequest

    if err := request.Bind(&loginRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    if err := validator.Validate(&loginRequest); err != nil {
        return c.JsonError(err.Error(), 422)
    }

    user, token, err := c.authService.Login(loginRequest.Email, loginRequest.Password)
    if err != nil {
        return c.JsonError("Invalid credentials", 401)
    }

    return c.Json(map[string]interface{}{
        "user":  NewUserResource(user, true, false),
        "token": token,
    })
}

// 用户注册
func (c *AuthApiController) Register(request http.Request) http.Response {
    var registerRequest RegisterRequest

    if err := request.Bind(&registerRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    if err := validator.Validate(&registerRequest); err != nil {
        return c.JsonError(err.Error(), 422)
    }

    user, token, err := c.authService.Register(registerRequest)
    if err != nil {
        return c.JsonError("Failed to register user", 500)
    }

    return c.Json(map[string]interface{}{
        "user":  NewUserResource(user, true, false),
        "token": token,
    }).Status(201)
}

// 刷新令牌
func (c *AuthApiController) Refresh(request http.Request) http.Response {
    token := request.Headers["Authorization"]
    if token == "" {
        return c.JsonError("Token is required", 401)
    }

    newToken, err := c.authService.RefreshToken(token)
    if err != nil {
        return c.JsonError("Invalid token", 401)
    }

    return c.Json(map[string]string{
        "token": newToken,
    })
}

// 用户登出
func (c *AuthApiController) Logout(request http.Request) http.Response {
    token := request.Headers["Authorization"]
    if token == "" {
        return c.JsonError("Token is required", 401)
    }

    err := c.authService.Logout(token)
    if err != nil {
        return c.JsonError("Failed to logout", 500)
    }

    return c.Json(map[string]string{
        "message": "Logged out successfully",
    })
}
```

### 2. API Token 认证

```go
// API Token 认证
type ApiTokenController struct {
    http.Controller
    tokenService *ApiTokenService
}

// 创建 API Token
func (c *ApiTokenController) Store(request http.Request) http.Response {
    user := request.Context["user"].(*User)

    var tokenRequest CreateTokenRequest

    if err := request.Bind(&tokenRequest); err != nil {
        return c.JsonError("Invalid request data", 400)
    }

    if err := validator.Validate(&tokenRequest); err != nil {
        return c.JsonError(err.Error(), 422)
    }

    token, err := c.tokenService.CreateToken(user.ID, tokenRequest.Name, tokenRequest.Abilities)
    if err != nil {
        return c.JsonError("Failed to create token", 500)
    }

    return c.Json(map[string]interface{}{
        "token": token.Token,
        "name":  token.Name,
        "abilities": strings.Split(token.Abilities, ","),
    }).Status(201)
}

// 获取用户的所有 Token
func (c *ApiTokenController) Index(request http.Request) http.Response {
    user := request.Context["user"].(*User)

    tokens, err := c.tokenService.GetUserTokens(user.ID)
    if err != nil {
        return c.JsonError("Failed to get tokens", 500)
    }

    return c.Json(tokens)
}

// 删除 API Token
func (c *ApiTokenController) Delete(id string, request http.Request) http.Response {
    user := request.Context["user"].(*User)

    tokenID, err := strconv.Atoi(id)
    if err != nil {
        return c.JsonError("Invalid token ID", 400)
    }

    err = c.tokenService.DeleteToken(uint(tokenID), user.ID)
    if err != nil {
        return c.JsonError("Failed to delete token", 500)
    }

    return c.Json(map[string]string{
        "message": "Token deleted successfully",
    })
}
```

## 📄 API 版本控制

### 1. 版本控制策略

```go
// API 版本控制
func RegisterApiRoutes(router *routing.Router) {
    // v1 API
    v1 := router.Group("/api/v1")
    registerV1Routes(v1)

    // v2 API
    v2 := router.Group("/api/v2")
    registerV2Routes(v2)
}

func registerV1Routes(router *routing.Router) {
    // v1 路由
    router.Get("/users", &UserApiControllerV1{}, "Index")
    router.Get("/users/{id}", &UserApiControllerV1{}, "Show")
    router.Post("/users", &UserApiControllerV1{}, "Store")
    router.Put("/users/{id}", &UserApiControllerV1{}, "Update")
    router.Delete("/users/{id}", &UserApiControllerV1{}, "Delete")
}

func registerV2Routes(router *routing.Router) {
    // v2 路由（新版本）
    router.Get("/users", &UserApiControllerV2{}, "Index")
    router.Get("/users/{id}", &UserApiControllerV2{}, "Show")
    router.Post("/users", &UserApiControllerV2{}, "Store")
    router.Put("/users/{id}", &UserApiControllerV2{}, "Update")
    router.Delete("/users/{id}", &UserApiControllerV2{}, "Delete")

    // v2 新增功能
    router.Get("/users/{id}/analytics", &UserApiControllerV2{}, "Analytics")
    router.Post("/users/{id}/avatar", &UserApiControllerV2{}, "UploadAvatar")
}
```

### 2. 版本兼容性

```go
// 版本兼容性中间件
type ApiVersionMiddleware struct {
    http.Middleware
    supportedVersions []string
}

func (m *ApiVersionMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    version := request.Headers["Accept-Version"]
    if version == "" {
        version = "v1" // 默认版本
    }

    // 检查版本是否支持
    supported := false
    for _, v := range m.supportedVersions {
        if v == version {
            supported = true
            break
        }
    }

    if !supported {
        return http.Response{
            StatusCode: 400,
            Body:       fmt.Sprintf(`{"error": "Unsupported API version: %s"}`, version),
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }

    // 将版本信息添加到请求上下文
    request.Context["api_version"] = version

    return next(request)
}

// 使用版本中间件
func RegisterApiRoutes(router *routing.Router) {
    api := router.Group("/api")
    api.Use(&ApiVersionMiddleware{
        supportedVersions: []string{"v1", "v2"},
    })

    v1 := api.Group("/v1")
    registerV1Routes(v1)

    v2 := api.Group("/v2")
    registerV2Routes(v2)
}
```

## 📊 API 响应格式

### 1. 统一响应格式

```go
// 统一响应格式
type ApiResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Errors  interface{} `json:"errors,omitempty"`
    Meta    interface{} `json:"meta,omitempty"`
}

// 成功响应
func (c *ApiController) SuccessResponse(data interface{}, message string) http.Response {
    response := ApiResponse{
        Success: true,
        Message: message,
        Data:    data,
    }

    return c.Json(response)
}

// 错误响应
func (c *ApiController) ErrorResponse(message string, errors interface{}, statusCode int) http.Response {
    response := ApiResponse{
        Success: false,
        Message: message,
        Errors:  errors,
    }

    return c.Json(response).Status(statusCode)
}

// 分页响应
func (c *ApiController) PaginatedResponse(data interface{}, total int64, page, limit int) http.Response {
    response := ApiResponse{
        Success: true,
        Data:    data,
        Meta: map[string]interface{}{
            "total": total,
            "page":  page,
            "limit": limit,
            "pages": int(math.Ceil(float64(total) / float64(limit))),
        },
    }

    return c.Json(response)
}

// 在控制器中使用
func (c *UserApiController) Index(request http.Request) http.Response {
    page := c.getPageParam(request)
    limit := c.getLimitParam(request)

    users, total := c.userService.GetUsers(page, limit)
    collection := NewUserCollection(users, total, page, limit)

    return c.PaginatedResponse(collection.Data, collection.Total, collection.Page, collection.Limit)
}

func (c *UserApiController) Show(id string) http.Response {
    userID, err := strconv.Atoi(id)
    if err != nil {
        return c.ErrorResponse("Invalid user ID", nil, 400)
    }

    user, err := c.userService.GetUser(uint(userID))
    if err != nil {
        return c.ErrorResponse("User not found", nil, 404)
    }

    resource := NewUserResource(user, true, false)
    return c.SuccessResponse(resource, "User retrieved successfully")
}
```

### 2. 错误处理

```go
// API 错误处理
type ApiError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Field   string `json:"field,omitempty"`
}

type ApiErrors []ApiError

// 验证错误响应
func (c *ApiController) ValidationErrorResponse(errors map[string]string) http.Response {
    apiErrors := make(ApiErrors, 0)

    for field, message := range errors {
        apiErrors = append(apiErrors, ApiError{
            Code:    "VALIDATION_ERROR",
            Message: message,
            Field:   field,
        })
    }

    return c.ErrorResponse("Validation failed", apiErrors, 422)
}

// 业务错误响应
func (c *ApiController) BusinessErrorResponse(code, message string) http.Response {
    apiError := ApiError{
        Code:    code,
        Message: message,
    }

    return c.ErrorResponse(message, apiError, 400)
}

// 在控制器中使用
func (c *UserApiController) Store(request http.Request) http.Response {
    var userRequest UserRequest

    if err := request.Bind(&userRequest); err != nil {
        return c.ErrorResponse("Invalid request data", nil, 400)
    }

    v := validator.New()
    // 添加验证规则...

    if !v.Passes() {
        return c.ValidationErrorResponse(v.Errors())
    }

    user, err := c.userService.CreateUser(userRequest)
    if err != nil {
        if err.Error() == "email_exists" {
            return c.BusinessErrorResponse("EMAIL_EXISTS", "Email already exists")
        }
        return c.ErrorResponse("Failed to create user", nil, 500)
    }

    resource := NewUserResource(user, true, false)
    return c.SuccessResponse(resource, "User created successfully").Status(201)
}
```

## 📈 API 监控

### 1. API 监控中间件

```go
// API 监控中间件
type ApiMonitorMiddleware struct {
    http.Middleware
    logger *Logger
}

func (m *ApiMonitorMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    start := time.Now()

    // 记录请求信息
    m.logger.Info("API Request", map[string]interface{}{
        "method":     request.Method,
        "path":       request.Path,
        "ip":         request.IP,
        "user_agent": request.Headers["User-Agent"],
        "timestamp":  start,
    })

    response := next(request)

    duration := time.Since(start)

    // 记录响应信息
    m.logger.Info("API Response", map[string]interface{}{
        "method":     request.Method,
        "path":       request.Path,
        "status":     response.StatusCode,
        "duration":   duration,
        "timestamp":  time.Now(),
    })

    // 添加响应头
    response.Headers["X-Response-Time"] = duration.String()
    response.Headers["X-API-Version"] = request.Context["api_version"].(string)

    return response
}
```

### 2. API 限流

```go
// API 限流中间件
type RateLimitMiddleware struct {
    http.Middleware
    limiter *RateLimiter
}

func (m *RateLimitMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    key := m.getRateLimitKey(request)

    if !m.limiter.Allow(key) {
        return http.Response{
            StatusCode: 429,
            Body:       `{"error": "Rate limit exceeded"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
                "Retry-After":  "60",
            },
        }
    }

    return next(request)
}

func (m *RateLimitMiddleware) getRateLimitKey(request http.Request) string {
    // 根据用户 ID 或 IP 生成限流键
    if user := request.Context["user"]; user != nil {
        return fmt.Sprintf("user:%d", user.(*User).ID)
    }

    return fmt.Sprintf("ip:%s", request.IP)
}
```

## 📚 API 文档

### 1. 自动生成 API 文档

```go
// API 文档生成器
type ApiDocumentation struct {
    Title       string                 `json:"title"`
    Version     string                 `json:"version"`
    Description string                 `json:"description"`
    BaseURL     string                 `json:"base_url"`
    Endpoints   []ApiEndpoint          `json:"endpoints"`
    Schemas     map[string]interface{} `json:"schemas"`
}

type ApiEndpoint struct {
    Method      string                 `json:"method"`
    Path        string                 `json:"path"`
    Summary     string                 `json:"summary"`
    Description string                 `json:"description"`
    Parameters  []ApiParameter         `json:"parameters"`
    Responses   map[int]ApiResponse    `json:"responses"`
    Tags        []string               `json:"tags"`
}

type ApiParameter struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    Required    bool   `json:"required"`
    Description string `json:"description"`
    Example     string `json:"example"`
}

// 生成 API 文档
func GenerateApiDocumentation() *ApiDocumentation {
    doc := &ApiDocumentation{
        Title:       "User Management API",
        Version:     "v1",
        Description: "API for managing users",
        BaseURL:     "https://api.example.com",
        Endpoints:   []ApiEndpoint{},
        Schemas:     make(map[string]interface{}),
    }

    // 添加端点
    doc.Endpoints = append(doc.Endpoints, ApiEndpoint{
        Method:      "GET",
        Path:        "/api/v1/users",
        Summary:     "Get users list",
        Description: "Retrieve a paginated list of users",
        Parameters: []ApiParameter{
            {
                Name:        "page",
                Type:        "integer",
                Required:    false,
                Description: "Page number",
                Example:     "1",
            },
            {
                Name:        "limit",
                Type:        "integer",
                Required:    false,
                Description: "Number of items per page",
                Example:     "10",
            },
        },
        Responses: map[int]ApiResponse{
            200: {
                Description: "Success",
                Schema:      "UserCollection",
            },
        },
        Tags: []string{"users"},
    })

    return doc
}

// API 文档控制器
type ApiDocController struct {
    http.Controller
}

func (c *ApiDocController) Index() http.Response {
    doc := GenerateApiDocumentation()
    return c.Json(doc)
}

func (c *ApiDocController) Swagger() http.Response {
    // 返回 Swagger UI HTML
    swaggerHTML := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>API Documentation</title>
        <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4/swagger-ui.css" />
    </head>
    <body>
        <div id="swagger-ui"></div>
        <script src="https://unpkg.com/swagger-ui-dist@4/swagger-ui-bundle.js"></script>
        <script>
            window.onload = function() {
                SwaggerUIBundle({
                    url: '/api/docs',
                    dom_id: '#swagger-ui',
                });
            };
        </script>
    </body>
    </html>
    `

    return http.Response{
        StatusCode: 200,
        Body:       swaggerHTML,
        Headers: map[string]string{
            "Content-Type": "text/html",
        },
    }
}
```

## 📚 总结

Laravel-Go Framework 的 API 开发系统提供了：

1. **RESTful API**: 完整的 RESTful API 设计
2. **API 资源**: 数据转换和格式化
3. **认证授权**: JWT、API Token 认证
4. **版本控制**: API 版本管理和兼容性
5. **响应格式**: 统一的响应格式和错误处理
6. **监控功能**: API 监控和限流
7. **文档生成**: 自动生成 API 文档

通过合理使用 API 开发系统，可以构建高质量、可维护的 API 服务。
