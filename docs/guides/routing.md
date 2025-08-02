# 路由系统指南

## 🛣️ 路由系统概览

Laravel-Go Framework 提供了强大而灵活的路由系统，支持多种路由类型、中间件、参数绑定等功能。

## 🚀 快速开始

### 基本路由定义

```go
// routes/web.go
package routes

import (
    "laravel-go/framework/http"
    "laravel-go/framework/routing"
)

func WebRoutes(router *routing.Router) {
    // GET 路由
    router.Get("/", func(request http.Request) http.Response {
        return http.Response{
            Body: "Welcome to Laravel-Go!",
        }
    })
    
    // POST 路由
    router.Post("/users", func(request http.Request) http.Response {
        return http.Response{
            Body: "User created",
        }
    })
    
    // PUT 路由
    router.Put("/users/{id}", func(request http.Request) http.Response {
        id := request.Params["id"]
        return http.Response{
            Body: "User " + id + " updated",
        }
    })
    
    // DELETE 路由
    router.Delete("/users/{id}", func(request http.Request) http.Response {
        id := request.Params["id"]
        return http.Response{
            Body: "User " + id + " deleted",
        }
    })
}
```

### 控制器路由

```go
// 定义控制器
type UserController struct {
    http.Controller
}

func (c *UserController) Index() http.Response {
    return c.Json(map[string]interface{}{
        "users": []string{"user1", "user2"},
    })
}

func (c *UserController) Show(id string) http.Response {
    return c.Json(map[string]interface{}{
        "id": id,
        "name": "John Doe",
    })
}

// 路由定义
func WebRoutes(router *routing.Router) {
    // 控制器路由
    router.Get("/users", &UserController{}, "Index")
    router.Get("/users/{id}", &UserController{}, "Show")
}
```

## 📋 路由类型

### 1. GET 路由
```go
router.Get("/path", handler)
```

### 2. POST 路由
```go
router.Post("/path", handler)
```

### 3. PUT 路由
```go
router.Put("/path", handler)
```

### 4. PATCH 路由
```go
router.Patch("/path", handler)
```

### 5. DELETE 路由
```go
router.Delete("/path", handler)
```

### 6. 多方法路由
```go
router.Match([]string{"GET", "POST"}, "/path", handler)
```

### 7. 任意方法路由
```go
router.Any("/path", handler)
```

## 🔗 路由参数

### 基本参数
```go
router.Get("/users/{id}", func(request http.Request) http.Response {
    id := request.Params["id"]
    return http.Response{Body: "User ID: " + id}
})
```

### 可选参数
```go
router.Get("/users/{id?}", func(request http.Request) http.Response {
    id, exists := request.Params["id"]
    if !exists {
        id = "default"
    }
    return http.Response{Body: "User ID: " + id}
})
```

### 正则表达式约束
```go
router.Get("/users/{id:[0-9]+}", func(request http.Request) http.Response {
    id := request.Params["id"]
    return http.Response{Body: "User ID: " + id}
})
```

### 多个参数
```go
router.Get("/users/{id}/posts/{postId}", func(request http.Request) http.Response {
    id := request.Params["id"]
    postId := request.Params["postId"]
    return http.Response{Body: "User " + id + ", Post " + postId}
})
```

## 🎯 路由组

### 基本路由组
```go
func WebRoutes(router *routing.Router) {
    // API 路由组
    router.Group("/api", func(group *routing.Router) {
        group.Get("/users", &UserController{}, "Index")
        group.Post("/users", &UserController{}, "Store")
        group.Get("/users/{id}", &UserController{}, "Show")
        group.Put("/users/{id}", &UserController{}, "Update")
        group.Delete("/users/{id}", &UserController{}, "Delete")
    })
    
    // Admin 路由组
    router.Group("/admin", func(group *routing.Router) {
        group.Get("/dashboard", &AdminController{}, "Dashboard")
        group.Get("/users", &AdminController{}, "Users")
    })
}
```

### 带中间件的路由组
```go
func WebRoutes(router *routing.Router) {
    // 需要认证的路由组
    router.Group("/api", func(group *routing.Router) {
        group.Use(&AuthMiddleware{})
        
        group.Get("/profile", &UserController{}, "Profile")
        group.Put("/profile", &UserController{}, "UpdateProfile")
    })
    
    // 需要管理员权限的路由组
    router.Group("/admin", func(group *routing.Router) {
        group.Use(&AuthMiddleware{})
        group.Use(&AdminMiddleware{})
        
        group.Get("/users", &AdminController{}, "Users")
        group.Delete("/users/{id}", &AdminController{}, "DeleteUser")
    })
}
```

### 嵌套路由组
```go
func WebRoutes(router *routing.Router) {
    router.Group("/api/v1", func(v1 *routing.Router) {
        v1.Group("/users", func(users *routing.Router) {
            users.Get("/", &UserController{}, "Index")
            users.Post("/", &UserController{}, "Store")
            
            users.Group("/{id}", func(user *routing.Router) {
                user.Get("/", &UserController{}, "Show")
                user.Put("/", &UserController{}, "Update")
                user.Delete("/", &UserController{}, "Delete")
                
                user.Group("/posts", func(posts *routing.Router) {
                    posts.Get("/", &PostController{}, "Index")
                    posts.Post("/", &PostController{}, "Store")
                })
            })
        })
    })
}
```

## 🛡️ 中间件

### 全局中间件
```go
// bootstrap/app.go
func Bootstrap() {
    router := routing.NewRouter()
    
    // 添加全局中间件
    router.Use(&LoggingMiddleware{})
    router.Use(&CorsMiddleware{})
    
    // 注册路由
    routes.WebRoutes(router)
    routes.ApiRoutes(router)
}
```

### 路由中间件
```go
func WebRoutes(router *routing.Router) {
    // 单个路由的中间件
    router.Get("/admin", &AdminController{}, "Dashboard").
        Use(&AuthMiddleware{}).
        Use(&AdminMiddleware{})
    
    // 路由组的中间件
    router.Group("/api", func(group *routing.Router) {
        group.Use(&AuthMiddleware{})
        
        group.Get("/users", &UserController{}, "Index")
    })
}
```

### 自定义中间件
```go
type AuthMiddleware struct {
    http.Middleware
}

func (m *AuthMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 检查认证
    token := request.Headers["Authorization"]
    if token == "" {
        return http.Response{
            StatusCode: 401,
            Body:       "Unauthorized",
        }
    }
    
    // 验证 token
    if !m.validateToken(token) {
        return http.Response{
            StatusCode: 401,
            Body:       "Invalid token",
        }
    }
    
    // 继续处理请求
    return next(request)
}

func (m *AuthMiddleware) validateToken(token string) bool {
    // 实现 token 验证逻辑
    return true
}
```

## 🔄 路由缓存

### 启用路由缓存
```go
// config/app.go
type App struct {
    // ... 其他配置
    RouteCache bool `env:"ROUTE_CACHE" default:"false"`
}

// 缓存路由
func CacheRoutes() error {
    router := routing.NewRouter()
    routes.WebRoutes(router)
    routes.ApiRoutes(router)
    
    return router.Cache("routes.cache")
}

// 加载缓存的路由
func LoadCachedRoutes() error {
    router := routing.NewRouter()
    return router.LoadCache("routes.cache")
}
```

### 清除路由缓存
```go
func ClearRouteCache() error {
    return os.Remove("routes.cache")
}
```

## 📝 路由命名

### 命名路由
```go
func WebRoutes(router *routing.Router) {
    // 命名路由
    router.Get("/users/{id}", &UserController{}, "Show").Name("users.show")
    router.Get("/posts/{id}", &PostController{}, "Show").Name("posts.show")
}
```

### 生成 URL
```go
// 在控制器中生成 URL
func (c *UserController) Index() http.Response {
    // 生成命名路由的 URL
    userUrl := c.Route("users.show", map[string]string{"id": "123"})
    postUrl := c.Route("posts.show", map[string]string{"id": "456"})
    
    return c.Json(map[string]interface{}{
        "user_url": userUrl,
        "post_url": postUrl,
    })
}
```

## 🎨 路由装饰器

### 路由前缀
```go
func WebRoutes(router *routing.Router) {
    // 添加前缀
    router.Prefix("/api/v1", func(group *routing.Router) {
        group.Get("/users", &UserController{}, "Index")
        group.Get("/posts", &PostController{}, "Index")
    })
}
```

### 域名约束
```go
func WebRoutes(router *routing.Router) {
    // 域名约束
    router.Domain("api.example.com", func(group *routing.Router) {
        group.Get("/users", &UserController{}, "Index")
    })
    
    router.Domain("admin.example.com", func(group *routing.Router) {
        group.Get("/dashboard", &AdminController{}, "Dashboard")
    })
}
```

### 子域名约束
```go
func WebRoutes(router *routing.Router) {
    // 子域名约束
    router.Subdomain("api", func(group *routing.Router) {
        group.Get("/users", &UserController{}, "Index")
    })
}
```

## 🔍 路由调试

### 列出所有路由
```go
// 在开发环境中列出所有路由
func ListRoutes() {
    router := routing.NewRouter()
    routes.WebRoutes(router)
    routes.ApiRoutes(router)
    
    routes := router.ListRoutes()
    for _, route := range routes {
        fmt.Printf("%s %s -> %s\n", route.Method, route.Path, route.Handler)
    }
}
```

### 路由测试
```go
func TestRoutes() {
    router := routing.NewRouter()
    routes.WebRoutes(router)
    
    // 测试路由匹配
    request := http.Request{
        Method: "GET",
        Path:   "/users/123",
    }
    
    response := router.Handle(request)
    fmt.Printf("Response: %+v\n", response)
}
```

## 🚀 性能优化

### 1. 路由缓存
```go
// 生产环境启用路由缓存
if config.Get("app.env") == "production" {
    router.LoadCache("routes.cache")
}
```

### 2. 路由压缩
```go
// 压缩相似路由
router.Compress()
```

### 3. 路由预编译
```go
// 预编译路由正则表达式
router.Precompile()
```

## 📊 路由统计

### 路由信息
```go
func GetRouteStats() map[string]interface{} {
    router := routing.NewRouter()
    routes.WebRoutes(router)
    
    stats := router.GetStats()
    return map[string]interface{}{
        "total_routes":    stats.TotalRoutes,
        "cached_routes":   stats.CachedRoutes,
        "compiled_routes": stats.CompiledRoutes,
        "memory_usage":    stats.MemoryUsage,
    }
}
```

## 🛠️ 高级功能

### 自定义路由处理器
```go
type CustomHandler struct {
    http.Handler
}

func (h *CustomHandler) Handle(request http.Request) http.Response {
    // 自定义处理逻辑
    return http.Response{
        Body: "Custom handler response",
    }
}

func WebRoutes(router *routing.Router) {
    router.Get("/custom", &CustomHandler{})
}
```

### 路由事件
```go
// 路由匹配事件
router.OnMatch(func(route *routing.Route, request http.Request) {
    log.Printf("Route matched: %s %s", request.Method, request.Path)
})

// 路由未找到事件
router.OnNotFound(func(request http.Request) http.Response {
    return http.Response{
        StatusCode: 404,
        Body:       "Route not found",
    }
})
```

### 路由重定向
```go
func WebRoutes(router *routing.Router) {
    // 永久重定向
    router.Redirect("/old-path", "/new-path", 301)
    
    // 临时重定向
    router.Redirect("/temp", "/permanent", 302)
}
```

## 📝 最佳实践

### 1. 路由组织
```go
// 按功能模块组织路由
func WebRoutes(router *routing.Router) {
    // 用户相关路由
    userRoutes(router)
    
    // 文章相关路由
    postRoutes(router)
    
    // 管理相关路由
    adminRoutes(router)
}

func userRoutes(router *routing.Router) {
    router.Group("/users", func(group *routing.Router) {
        group.Get("/", &UserController{}, "Index")
        group.Post("/", &UserController{}, "Store")
        group.Get("/{id}", &UserController{}, "Show")
        group.Put("/{id}", &UserController{}, "Update")
        group.Delete("/{id}", &UserController{}, "Delete")
    })
}
```

### 2. 中间件使用
```go
// 合理使用中间件
func WebRoutes(router *routing.Router) {
    // 公开路由
    router.Get("/", &HomeController{}, "Index")
    router.Get("/about", &HomeController{}, "About")
    
    // 需要认证的路由
    router.Group("/dashboard", func(group *routing.Router) {
        group.Use(&AuthMiddleware{})
        
        group.Get("/", &DashboardController{}, "Index")
        group.Get("/profile", &DashboardController{}, "Profile")
    })
}
```

### 3. 参数验证
```go
// 在路由中使用参数约束
func WebRoutes(router *routing.Router) {
    // 数字 ID 约束
    router.Get("/users/{id:[0-9]+}", &UserController{}, "Show")
    
    // 可选参数
    router.Get("/posts/{category?}", &PostController{}, "Index")
    
    // 多个参数约束
    router.Get("/files/{year:[0-9]{4}}/{month:[0-9]{2}}/{filename}", &FileController{}, "Show")
}
```

## 📚 总结

Laravel-Go Framework 的路由系统提供了：

1. **灵活性**: 支持多种路由类型和参数
2. **可扩展性**: 易于添加中间件和自定义处理器
3. **性能优化**: 支持路由缓存和压缩
4. **开发体验**: 提供调试工具和统计信息
5. **最佳实践**: 遵循 RESTful 设计原则

通过合理使用路由系统的各种功能，可以构建出高效、可维护的 Web 应用程序。 