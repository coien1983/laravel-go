# HTTP 系统指南

## 📖 概述

Laravel-Go Framework 提供了完整的 HTTP 系统，包括请求处理、响应生成、中间件、路由管理、文件上传等功能，帮助构建强大的 Web 应用程序。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [HTTP 系统 API 参考](../api/http.md)

## 🚀 快速开始

### 1. 基本 HTTP 服务器

```go
// 创建 HTTP 服务器
server := http.NewServer()

// 注册路由
server.Router.Get("/", func(request http.Request) http.Response {
    return http.Response{
        StatusCode: 200,
        Body:       "Hello, Laravel-Go!",
        Headers: map[string]string{
            "Content-Type": "text/plain",
        },
    }
})

// 启动服务器
server.Start(":8080")
```

### 2. 控制器使用

```go
// 用户控制器
type UserController struct {
    http.Controller
    userService *Services.UserService
}

func NewUserController() *UserController {
    return &UserController{
        userService: Services.NewUserService(),
    }
}

// 获取用户列表
func (c *UserController) Index(request http.Request) http.Response {
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
func (c *UserController) Show(id string, request http.Request) http.Response {
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
func (c *UserController) Store(request http.Request) http.Response {
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
func (c *UserController) Update(id string, request http.Request) http.Response {
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
func (c *UserController) Delete(id string, request http.Request) http.Response {
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

// 辅助方法
func (c *UserController) getPageParam(request http.Request) int {
    page := request.Query["page"]
    if page == "" {
        return 1
    }

    if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
        return pageNum
    }

    return 1
}

func (c *UserController) getLimitParam(request http.Request) int {
    limit := request.Query["limit"]
    if limit == "" {
        return 10
    }

    if limitNum, err := strconv.Atoi(limit); err == nil && limitNum > 0 {
        return limitNum
    }

    return 10
}
```

### 3. 请求处理

```go
// 请求数据绑定
type UserRequest struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min:8"`
}

type UserUpdateRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

// 请求验证
func (c *UserController) Store(request http.Request) http.Response {
    var userRequest UserRequest

    // 绑定 JSON 数据
    if err := request.BindJSON(&userRequest); err != nil {
        return c.JsonError("Invalid JSON data", 400)
    }

    // 验证数据
    if err := validator.Validate(&userRequest); err != nil {
        return c.JsonError(err.Error(), 422)
    }

    // 处理业务逻辑
    user, err := c.userService.CreateUser(userRequest)
    if err != nil {
        return c.JsonError("Failed to create user", 500)
    }

    return c.Json(user).Status(201)
}

// 处理表单数据
func (c *UserController) UpdateProfile(request http.Request) http.Response {
    // 获取表单数据
    name := request.Form["name"]
    email := request.Form["email"]

    // 验证数据
    if name == "" {
        return c.JsonError("Name is required", 422)
    }

    if email == "" {
        return c.JsonError("Email is required", 422)
    }

    // 更新用户资料
    user := request.Context["user"].(*Models.User)
    user.Name = name
    user.Email = email

    if err := c.userService.UpdateUser(user.ID, user); err != nil {
        return c.JsonError("Failed to update profile", 500)
    }

    return c.Json(user)
}

// 处理查询参数
func (c *UserController) Search(request http.Request) http.Response {
    query := request.Query["q"]
    category := request.Query["category"]
    sort := request.Query["sort"]

    // 构建搜索条件
    filters := make(map[string]interface{})
    if query != "" {
        filters["query"] = query
    }
    if category != "" {
        filters["category"] = category
    }
    if sort != "" {
        filters["sort"] = sort
    }

    // 执行搜索
    users, err := c.userService.SearchUsers(filters)
    if err != nil {
        return c.JsonError("Search failed", 500)
    }

    return c.Json(users)
}
```

### 4. 响应生成

```go
// JSON 响应
func (c *UserController) GetUser(id string, request http.Request) http.Response {
    user, err := c.userService.GetUser(uint(id))
    if err != nil {
        return c.JsonError("User not found", 404)
    }

    return c.Json(user)
}

// 自定义状态码
func (c *UserController) CreateUser(request http.Request) http.Response {
    user, err := c.userService.CreateUser(request.Body)
    if err != nil {
        return c.JsonError("Failed to create user", 500)
    }

    return c.Json(user).Status(201)
}

// 带头的响应
func (c *UserController) DownloadFile(request http.Request) http.Response {
    filename := request.Query["file"]
    if filename == "" {
        return c.JsonError("File name is required", 400)
    }

    filepath := "storage/files/" + filename
    data, err := ioutil.ReadFile(filepath)
    if err != nil {
        return c.JsonError("File not found", 404)
    }

    return http.Response{
        StatusCode: 200,
        Body:       string(data),
        Headers: map[string]string{
            "Content-Type":        "application/octet-stream",
            "Content-Disposition": "attachment; filename=" + filename,
        },
    }
}

// 重定向响应
func (c *UserController) Redirect(request http.Request) http.Response {
    return http.Response{
        StatusCode: 302,
        Headers: map[string]string{
            "Location": "/dashboard",
        },
    }
}

// 视图响应
func (c *UserController) ShowProfile(request http.Request) http.Response {
    user := request.Context["user"].(*Models.User)

    data := map[string]interface{}{
        "user": user,
        "title": "User Profile",
    }

    html, err := template.Render("user.profile", data)
    if err != nil {
        return c.JsonError("Template error", 500)
    }

    return http.Response{
        StatusCode: 200,
        Body:       html,
        Headers: map[string]string{
            "Content-Type": "text/html",
        },
    }
}
```

### 5. 文件上传

```go
// 文件上传控制器
type FileController struct {
    http.Controller
}

// 上传单个文件
func (c *FileController) Upload(request http.Request) http.Response {
    file := request.Files["file"]
    if file == nil {
        return c.JsonError("No file uploaded", 400)
    }

    // 验证文件类型
    if !c.isAllowedFileType(file.Filename) {
        return c.JsonError("File type not allowed", 400)
    }

    // 验证文件大小
    if file.Size > 10*1024*1024 { // 10MB
        return c.JsonError("File too large", 400)
    }

    // 生成安全的文件名
    safeName := c.generateSafeFileName(file.Filename)

    // 保存文件
    uploadPath := "storage/uploads/" + safeName
    if err := file.Save(uploadPath); err != nil {
        return c.JsonError("Failed to save file", 500)
    }

    return c.Json(map[string]string{
        "filename": safeName,
        "path":     uploadPath,
        "size":     fmt.Sprintf("%d", file.Size),
    })
}

// 上传多个文件
func (c *FileController) UploadMultiple(request http.Request) http.Response {
    files := request.Files["files"]
    if len(files) == 0 {
        return c.JsonError("No files uploaded", 400)
    }

    var uploadedFiles []map[string]string

    for _, file := range files {
        // 验证文件
        if !c.isAllowedFileType(file.Filename) {
            continue
        }

        if file.Size > 10*1024*1024 {
            continue
        }

        // 保存文件
        safeName := c.generateSafeFileName(file.Filename)
        uploadPath := "storage/uploads/" + safeName

        if err := file.Save(uploadPath); err != nil {
            continue
        }

        uploadedFiles = append(uploadedFiles, map[string]string{
            "filename": safeName,
            "path":     uploadPath,
            "size":     fmt.Sprintf("%d", file.Size),
        })
    }

    return c.Json(map[string]interface{}{
        "files": uploadedFiles,
        "count": len(uploadedFiles),
    })
}

// 辅助方法
func (c *FileController) isAllowedFileType(filename string) bool {
    allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx"}
    ext := strings.ToLower(filepath.Ext(filename))

    for _, allowedType := range allowedTypes {
        if ext == allowedType {
            return true
        }
    }

    return false
}

func (c *FileController) generateSafeFileName(filename string) string {
    ext := filepath.Ext(filename)
    name := strings.TrimSuffix(filename, ext)

    // 生成随机字符串
    randomStr := fmt.Sprintf("%d", time.Now().UnixNano())

    return name + "_" + randomStr + ext
}
```

### 6. 中间件

```go
// 认证中间件
type AuthMiddleware struct{}

func (m *AuthMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    token := request.Headers["Authorization"]
    if token == "" {
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "Unauthorized"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }

    // 验证令牌
    user, err := auth.ValidateToken(token)
    if err != nil {
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "Invalid token"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }

    // 将用户信息添加到请求上下文
    request.Context["user"] = user

    return next(request)
}

// 日志中间件
type LogMiddleware struct {
    logger *log.Logger
}

func (m *LogMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    start := time.Now()

    // 记录请求开始
    m.logger.Info("Request started", map[string]interface{}{
        "method": request.Method,
        "path":   request.Path,
        "ip":     request.IP,
    })

    // 处理请求
    response := next(request)

    // 记录请求完成
    duration := time.Since(start)
    m.logger.Info("Request completed", map[string]interface{}{
        "method":   request.Method,
        "path":     request.Path,
        "status":   response.StatusCode,
        "duration": duration.String(),
    })

    return response
}

// CORS 中间件
type CORSMiddleware struct{}

func (m *CORSMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    response := next(request)

    // 添加 CORS 头
    response.Headers["Access-Control-Allow-Origin"] = "*"
    response.Headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE, OPTIONS"
    response.Headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization"

    return response
}
```

### 7. 错误处理

```go
// 错误处理中间件
type ErrorHandlerMiddleware struct {
    logger *log.Logger
}

func (m *ErrorHandlerMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    defer func() {
        if r := recover(); r != nil {
            // 记录错误
            m.logger.Error("Panic recovered", map[string]interface{}{
                "error": r,
                "stack": string(debug.Stack()),
            })
        }
    }()

    return next(request)
}

// 自定义错误响应
func (c *UserController) HandleError(err error, statusCode int) http.Response {
    errorResponse := map[string]interface{}{
        "error":   err.Error(),
        "code":    statusCode,
        "message": http.StatusText(statusCode),
        "time":    time.Now(),
    }

    return http.Response{
        StatusCode: statusCode,
        Body:       c.toJSON(errorResponse),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }
}

// 404 错误处理
func (c *UserController) NotFound(request http.Request) http.Response {
    return c.HandleError(errors.New("Page not found"), 404)
}

// 500 错误处理
func (c *UserController) InternalError(request http.Request) http.Response {
    return c.HandleError(errors.New("Internal server error"), 500)
}
```

### 8. 路由管理

```go
// 路由注册
func RegisterRoutes() {
    router := routing.NewRouter()

    // 公开路由
    router.Get("/", &HomeController{}, "Index")
    router.Get("/about", &PageController{}, "About")
    router.Get("/contact", &PageController{}, "Contact")

    // 认证路由
    router.Post("/auth/login", &AuthController{}, "Login")
    router.Post("/auth/register", &AuthController{}, "Register")
    router.Post("/auth/logout", &AuthController{}, "Logout")

    // 需要认证的路由
    authenticated := router.Group("/api")
    authenticated.Use(&middleware.AuthMiddleware{})

    // 用户相关路由
    authenticated.Get("/users", &UserController{}, "Index")
    authenticated.Get("/users/{id}", &UserController{}, "Show")
    authenticated.Post("/users", &UserController{}, "Store")
    authenticated.Put("/users/{id}", &UserController{}, "Update")
    authenticated.Delete("/users/{id}", &UserController{}, "Delete")

    // 文件上传路由
    authenticated.Post("/upload", &FileController{}, "Upload")
    authenticated.Post("/upload/multiple", &FileController{}, "UploadMultiple")

    // 管理路由
    admin := authenticated.Group("/admin")
    admin.Use(&middleware.AdminMiddleware{})

    admin.Get("/dashboard", &AdminController{}, "Dashboard")
    admin.Get("/users", &AdminController{}, "Users")
    admin.Get("/logs", &AdminController{}, "Logs")
}
```

### 9. 会话管理

```go
// 会话中间件
type SessionMiddleware struct{}

func (m *SessionMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 启动会话
    session := request.Session()

    // 检查会话是否过期
    if session.IsExpired() {
        session.Regenerate()
    }

    // 更新最后活动时间
    session.Put("last_activity", time.Now())

    response := next(request)

    return response
}

// 在控制器中使用会话
func (c *UserController) Login(request http.Request) http.Response {
    email := request.Body["email"].(string)
    password := request.Body["password"].(string)

    user, err := c.userService.Authenticate(email, password)
    if err != nil {
        return c.JsonError("Invalid credentials", 401)
    }

    // 创建会话
    session := request.Session()
    session.Put("user_id", user.ID)
    session.Put("user_email", user.Email)
    session.Put("logged_in_at", time.Now())

    return c.Json(map[string]interface{}{
        "user":  user,
        "token": session.Get("token"),
    })
}

func (c *UserController) Logout(request http.Request) http.Response {
    // 销毁会话
    session := request.Session()
    session.Destroy()

    return c.Json(map[string]string{
        "message": "Logged out successfully",
    })
}
```

### 10. 缓存控制

```go
// 缓存控制中间件
type CacheControlMiddleware struct {
    maxAge time.Duration
}

func NewCacheControlMiddleware(maxAge time.Duration) *CacheControlMiddleware {
    return &CacheControlMiddleware{maxAge: maxAge}
}

func (m *CacheControlMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    response := next(request)

    // 添加缓存控制头
    response.Headers["Cache-Control"] = fmt.Sprintf("public, max-age=%d", int(m.maxAge.Seconds()))
    response.Headers["Expires"] = time.Now().Add(m.maxAge).Format(time.RFC1123)

    return response
}

// 在控制器中使用缓存
func (c *UserController) GetUserProfile(id string, request http.Request) http.Response {
    cacheKey := fmt.Sprintf("user_profile:%s", id)

    // 尝试从缓存获取
    if cached, exists := cache.Get(cacheKey); exists {
        return c.Json(cached)
    }

    // 从数据库获取
    user, err := c.userService.GetUser(uint(id))
    if err != nil {
        return c.JsonError("User not found", 404)
    }

    // 缓存结果
    cache.Set(cacheKey, user, time.Hour)

    return c.Json(user)
}
```

## 📚 总结

Laravel-Go Framework 的 HTTP 系统提供了：

1. **请求处理**: 数据绑定、验证、查询参数处理
2. **响应生成**: JSON、HTML、文件下载、重定向
3. **文件上传**: 单文件、多文件、安全验证
4. **中间件**: 认证、日志、CORS、错误处理
5. **路由管理**: 路由注册、分组、参数绑定
6. **会话管理**: 会话创建、销毁、数据存储
7. **缓存控制**: 响应缓存、缓存头设置
8. **错误处理**: 统一错误处理、自定义错误响应

通过合理使用 HTTP 系统，可以构建功能完整、安全可靠的 Web 应用程序。
