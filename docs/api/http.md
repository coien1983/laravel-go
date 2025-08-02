# HTTP API 参考

本文档提供 Laravel-Go Framework HTTP 相关组件的 API 参考。

## 📦 Router

路由管理器，处理 URL 路由和请求分发。

### 基础路由

#### Get(path string, handler interface{}) *Route
注册 GET 路由。

```go
app.Router().Get("/users", func(c http.Context) http.Response {
    return c.Json(map[string]string{"message": "Get users"})
})
```

#### Post(path string, handler interface{}) *Route
注册 POST 路由。

```go
app.Router().Post("/users", func(c http.Context) http.Response {
    return c.Json(map[string]string{"message": "Create user"})
})
```

#### Put(path string, handler interface{}) *Route
注册 PUT 路由。

```go
app.Router().Put("/users/{id}", func(c http.Context) http.Response {
    return c.Json(map[string]string{"message": "Update user"})
})
```

#### Delete(path string, handler interface{}) *Route
注册 DELETE 路由。

```go
app.Router().Delete("/users/{id}", func(c http.Context) http.Response {
    return c.Json(map[string]string{"message": "Delete user"})
})
```

#### Patch(path string, handler interface{}) *Route
注册 PATCH 路由。

```go
app.Router().Patch("/users/{id}", func(c http.Context) http.Response {
    return c.Json(map[string]string{"message": "Patch user"})
})
```

#### Any(path string, handler interface{}) *Route
注册支持所有 HTTP 方法的路由。

```go
app.Router().Any("/webhook", func(c http.Context) http.Response {
    return c.Json(map[string]string{"message": "Webhook received"})
})
```

### 路由组

#### Group(prefix string) *RouterGroup
创建路由组。

```go
api := app.Router().Group("/api")
{
    api.Get("/users", controllers.UserController{}.Index)
    api.Post("/users", controllers.UserController{}.Store)
    api.Get("/users/{id}", controllers.UserController{}.Show)
    api.Put("/users/{id}", controllers.UserController{}.Update)
    api.Delete("/users/{id}", controllers.UserController{}.Destroy)
}
```

#### Use(middleware ...interface{}) *RouterGroup
为路由组应用中间件。

```go
api := app.Router().Group("/api")
api.Use(middleware.AuthMiddleware{}, middleware.RateLimitMiddleware{})
```

### 路由参数

#### 路径参数
```go
app.Router().Get("/users/{id}", func(c http.Context) http.Response {
    id := c.Param("id")
    return c.Json(map[string]string{"id": id})
})
```

#### 可选参数
```go
app.Router().Get("/users/{id?}", func(c http.Context) http.Response {
    id := c.Param("id") // 可能为空
    return c.Json(map[string]string{"id": id})
})
```

#### 正则表达式参数
```go
app.Router().Get("/users/{id:[0-9]+}", func(c http.Context) http.Response {
    id := c.Param("id")
    return c.Json(map[string]string{"id": id})
})
```

### 路由方法

#### Name(name string) *Route
为路由命名。

```go
app.Router().Get("/users", handler).Name("users.index")
```

#### Middleware(middleware ...interface{}) *Route
为单个路由应用中间件。

```go
app.Router().Get("/admin", handler).Middleware(middleware.AuthMiddleware{})
```

#### Where(name, pattern string) *Route
为路由参数添加约束。

```go
app.Router().Get("/users/{id}", handler).Where("id", "[0-9]+")
```

## 🌐 Context

HTTP 上下文，提供请求和响应处理功能。

### 请求信息

#### Request() *http.Request
获取原始 HTTP 请求。

```go
func handler(c http.Context) http.Response {
    req := c.Request()
    userAgent := req.Header.Get("User-Agent")
    return c.Json(map[string]string{"user_agent": userAgent})
}
```

#### Method() string
获取请求方法。

```go
func handler(c http.Context) http.Response {
    method := c.Method()
    return c.Json(map[string]string{"method": method})
}
```

#### URL() *url.URL
获取请求 URL。

```go
func handler(c http.Context) http.Response {
    url := c.URL()
    return c.Json(map[string]string{
        "path": url.Path,
        "query": url.RawQuery,
    })
}
```

#### Header(name string) string
获取请求头。

```go
func handler(c http.Context) http.Response {
    contentType := c.Header("Content-Type")
    return c.Json(map[string]string{"content_type": contentType})
}
```

#### Headers() map[string]string
获取所有请求头。

```go
func handler(c http.Context) http.Response {
    headers := c.Headers()
    return c.Json(headers)
}
```

### 参数获取

#### Param(name string) string
获取路径参数。

```go
app.Router().Get("/users/{id}", func(c http.Context) http.Response {
    id := c.Param("id")
    return c.Json(map[string]string{"id": id})
})
```

#### Query(name string) string
获取查询参数。

```go
func handler(c http.Context) http.Response {
    page := c.Query("page")
    limit := c.Query("limit")
    return c.Json(map[string]string{
        "page": page,
        "limit": limit,
    })
}
```

#### QueryInt(name string) int
获取整数查询参数。

```go
func handler(c http.Context) http.Response {
    page := c.QueryInt("page")
    limit := c.QueryInt("limit")
    return c.Json(map[string]int{
        "page": page,
        "limit": limit,
    })
}
```

#### QueryBool(name string) bool
获取布尔查询参数。

```go
func handler(c http.Context) http.Response {
    active := c.QueryBool("active")
    return c.Json(map[string]bool{"active": active})
}
```

#### Form(name string) string
获取表单参数。

```go
func handler(c http.Context) http.Response {
    name := c.Form("name")
    email := c.Form("email")
    return c.Json(map[string]string{
        "name": name,
        "email": email,
    })
}
```

#### File(name string) *multipart.FileHeader
获取上传文件。

```go
func handler(c http.Context) http.Response {
    file, err := c.File("avatar")
    if err != nil {
        return c.Json(map[string]string{"error": "File not found"}).Status(400)
    }
    
    // 处理文件上传
    return c.Json(map[string]string{"filename": file.Filename})
}
```

### 请求体解析

#### BindJSON(v interface{}) error
解析 JSON 请求体。

```go
func handler(c http.Context) http.Response {
    var user User
    if err := c.BindJSON(&user); err != nil {
        return c.Json(map[string]string{"error": "Invalid JSON"}).Status(400)
    }
    
    return c.Json(user)
}
```

#### BindXML(v interface{}) error
解析 XML 请求体。

```go
func handler(c http.Context) http.Response {
    var data Data
    if err := c.BindXML(&data); err != nil {
        return c.Json(map[string]string{"error": "Invalid XML"}).Status(400)
    }
    
    return c.Json(data)
}
```

#### BindForm(v interface{}) error
解析表单数据。

```go
func handler(c http.Context) http.Response {
    var user User
    if err := c.BindForm(&user); err != nil {
        return c.Json(map[string]string{"error": "Invalid form data"}).Status(400)
    }
    
    return c.Json(user)
}
```

#### Body() []byte
获取原始请求体。

```go
func handler(c http.Context) http.Response {
    body := c.Body()
    return c.Json(map[string]string{"body": string(body)})
}
```

### 响应方法

#### Json(data interface{}) *Response
返回 JSON 响应。

```go
func handler(c http.Context) http.Response {
    return c.Json(map[string]interface{}{
        "message": "Success",
        "data": []string{"item1", "item2"},
    })
}
```

#### XML(data interface{}) *Response
返回 XML 响应。

```go
func handler(c http.Context) http.Response {
    return c.XML(Data{
        Message: "Success",
        Items: []string{"item1", "item2"},
    })
}
```

#### Text(text string) *Response
返回文本响应。

```go
func handler(c http.Context) http.Response {
    return c.Text("Hello, World!")
}
```

#### HTML(html string) *Response
返回 HTML 响应。

```go
func handler(c http.Context) http.Response {
    return c.HTML("<h1>Hello, World!</h1>")
}
```

#### View(name string, data interface{}) *Response
渲染视图模板。

```go
func handler(c http.Context) http.Response {
    return c.View("users.index", map[string]interface{}{
        "users": users,
        "title": "User List",
    })
}
```

#### Redirect(url string) *Response
重定向响应。

```go
func handler(c http.Context) http.Response {
    return c.Redirect("/users")
}
```

#### Status(code int) *Response
设置响应状态码。

```go
func handler(c http.Context) http.Response {
    return c.Json(map[string]string{"error": "Not found"}).Status(404)
}
```

#### Header(name, value string) *Response
设置响应头。

```go
func handler(c http.Context) http.Response {
    return c.Json(data).Header("Cache-Control", "no-cache")
}
```

#### Cookie(name, value string, options ...CookieOption) *Response
设置 Cookie。

```go
func handler(c http.Context) http.Response {
    return c.Json(data).Cookie("session", "token123", 
        http.CookieOption{
            MaxAge: 3600,
            HttpOnly: true,
            Secure: true,
        })
}
```

## 🔧 Middleware

中间件接口和实现。

### 中间件接口

```go
type Middleware interface {
    Handle(next HandlerFunc) HandlerFunc
}
```

### 基础中间件

#### LoggingMiddleware
日志记录中间件。

```go
type LoggingMiddleware struct{}

func (m *LoggingMiddleware) Handle(next HandlerFunc) HandlerFunc {
    return func(c Context) Response {
        start := time.Now()
        
        response := next(c)
        
        duration := time.Since(start)
        log.Printf("%s %s - %d - %v", c.Method(), c.URL().Path, response.Status(), duration)
        
        return response
    }
}
```

#### CORSMiddleware
CORS 中间件。

```go
type CORSMiddleware struct{}

func (m *CORSMiddleware) Handle(next HandlerFunc) HandlerFunc {
    return func(c Context) Response {
        response := next(c)
        
        return response.Header("Access-Control-Allow-Origin", "*").
            Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS").
            Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
    }
}
```

#### RateLimitMiddleware
速率限制中间件。

```go
type RateLimitMiddleware struct {
    limit int
    window time.Duration
}

func (m *RateLimitMiddleware) Handle(next HandlerFunc) HandlerFunc {
    return func(c Context) Response {
        // 实现速率限制逻辑
        return next(c)
    }
}
```

### 中间件链

#### 全局中间件
```go
app.Router().Use(middleware.LoggingMiddleware{})
app.Router().Use(middleware.CORSMiddleware{})
```

#### 路由组中间件
```go
api := app.Router().Group("/api")
api.Use(middleware.AuthMiddleware{})
api.Use(middleware.RateLimitMiddleware{Limit: 100, Window: time.Minute})
```

#### 单个路由中间件
```go
app.Router().Get("/admin", handler).Middleware(middleware.AuthMiddleware{})
```

## 🎯 Controller

控制器基类和辅助方法。

### BaseController

```go
type BaseController struct {
    app *framework.Application
}

func (c *BaseController) App() *framework.Application {
    return c.app
}

func (c *BaseController) DB() *database.Database {
    return c.app.DB()
}

func (c *BaseController) Cache() *cache.Cache {
    return c.app.Cache()
}

func (c *BaseController) Auth() *auth.Auth {
    return c.app.Auth()
}

func (c *BaseController) Validate(data interface{}, rules map[string]string) error {
    return validation.Validate(data, rules)
}
```

### 控制器示例

```go
type UserController struct {
    http.BaseController
}

func (c *UserController) Index(ctx http.Context) http.Response {
    users := c.DB().Table("users").Get()
    return ctx.Json(users)
}

func (c *UserController) Store(ctx http.Context) http.Response {
    var user User
    if err := ctx.BindJSON(&user); err != nil {
        return ctx.Json(map[string]string{"error": "Invalid JSON"}).Status(400)
    }
    
    // 验证数据
    rules := map[string]string{
        "name":  "required|string|max:255",
        "email": "required|email|unique:users",
    }
    
    if err := c.Validate(&user, rules); err != nil {
        return ctx.Json(map[string]string{"error": err.Error()}).Status(422)
    }
    
    // 创建用户
    result := c.DB().Table("users").Create(user)
    return ctx.Json(result).Status(201)
}

func (c *UserController) Show(ctx http.Context) http.Response {
    id := ctx.Param("id")
    user := c.DB().Table("users").Where("id", id).First()
    
    if user == nil {
        return ctx.Json(map[string]string{"error": "User not found"}).Status(404)
    }
    
    return ctx.Json(user)
}

func (c *UserController) Update(ctx http.Context) http.Response {
    id := ctx.Param("id")
    var user User
    
    if err := ctx.BindJSON(&user); err != nil {
        return ctx.Json(map[string]string{"error": "Invalid JSON"}).Status(400)
    }
    
    result := c.DB().Table("users").Where("id", id).Update(user)
    return ctx.Json(result)
}

func (c *UserController) Destroy(ctx http.Context) http.Response {
    id := ctx.Param("id")
    c.DB().Table("users").Where("id", id).Delete()
    return ctx.Json(map[string]string{"message": "User deleted"})
}
```

## 📊 Response

响应对象和构建器。

### 响应方法

#### Status() int
获取响应状态码。

```go
response := ctx.Json(data).Status(201)
status := response.Status() // 201
```

#### Body() []byte
获取响应体。

```go
response := ctx.Json(data)
body := response.Body()
```

#### Headers() map[string]string
获取响应头。

```go
response := ctx.Json(data).Header("Cache-Control", "no-cache")
headers := response.Headers()
```

#### SetHeader(name, value string) *Response
设置响应头。

```go
response := ctx.Json(data).SetHeader("Content-Type", "application/json")
```

#### SetCookie(name, value string, options ...CookieOption) *Response
设置 Cookie。

```go
response := ctx.Json(data).SetCookie("session", "token123", 
    http.CookieOption{
        MaxAge: 3600,
        HttpOnly: true,
    })
```

### Cookie 选项

```go
type CookieOption struct {
    Domain   string
    Path     string
    MaxAge   int
    HttpOnly bool
    Secure   bool
    SameSite string
}
```

## 🔍 错误处理

### HTTP 错误

```go
// 400 Bad Request
return ctx.Json(map[string]string{"error": "Bad request"}).Status(400)

// 401 Unauthorized
return ctx.Json(map[string]string{"error": "Unauthorized"}).Status(401)

// 403 Forbidden
return ctx.Json(map[string]string{"error": "Forbidden"}).Status(403)

// 404 Not Found
return ctx.Json(map[string]string{"error": "Not found"}).Status(404)

// 422 Unprocessable Entity
return ctx.Json(map[string]string{"error": "Validation failed"}).Status(422)

// 500 Internal Server Error
return ctx.Json(map[string]string{"error": "Internal server error"}).Status(500)
```

### 错误中间件

```go
type ErrorMiddleware struct{}

func (m *ErrorMiddleware) Handle(next HandlerFunc) HandlerFunc {
    return func(c Context) Response {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Panic: %v", r)
                c.Json(map[string]string{"error": "Internal server error"}).Status(500)
            }
        }()
        
        return next(c)
    }
}
```

## 📚 下一步

了解更多 HTTP 相关功能：

1. [中间件开发](guides/middleware.md) - 自定义中间件开发
2. [控制器开发](guides/controllers.md) - MVC 控制器开发
3. [路由系统](guides/routing.md) - 路由配置和管理
4. [模板引擎](guides/templates.md) - 视图模板系统
5. [认证授权](guides/auth.md) - 用户认证和权限控制

---

这些是 Laravel-Go Framework 的 HTTP API。掌握这些 API 将帮助你构建强大的 Web 应用！ 🚀 