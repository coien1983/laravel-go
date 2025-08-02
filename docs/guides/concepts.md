# 基础概念

本文档介绍 Laravel-Go Framework 的核心概念和设计理念。

## 🎯 框架理念

Laravel-Go Framework 基于 Laravel 的设计思路，结合 Go 语言的特性，提供优雅、高效的 Web 开发体验。

### 设计原则

1. **优雅简洁** - 提供简洁易用的 API
2. **约定优于配置** - 减少配置，提高开发效率
3. **依赖注入** - 松耦合的组件设计
4. **中间件模式** - 灵活的请求处理管道
5. **ORM 抽象** - 优雅的数据库操作

## 🏗️ 核心架构

### 应用生命周期

```
启动应用 → 加载配置 → 注册服务 → 启动服务器 → 处理请求 → 返回响应
```

### 请求处理流程

```
HTTP 请求 → 路由匹配 → 中间件处理 → 控制器执行 → 响应返回
```

## 📦 核心组件

### 1. 应用容器 (Application Container)

应用容器是框架的核心，负责管理所有服务和依赖关系。

```go
// 创建应用实例
app := framework.NewApplication()

// 注册服务
app.Container().Singleton("db", func() interface{} {
    return database.NewConnection()
})

// 解析服务
db := app.Container().Make("db").(database.Connection)
```

**特性：**

- 依赖注入
- 服务生命周期管理
- 单例模式支持
- 接口绑定

### 2. 配置管理 (Configuration)

统一的配置管理系统，支持多环境配置。

```go
// 获取配置
dbHost := app.Config().Get("database.host")
debug := app.Config().Get("app.debug").(bool)

// 环境变量覆盖
app.Config().Set("database.host", os.Getenv("DB_HOST"))
```

**特性：**

- 多环境支持
- 环境变量覆盖
- 嵌套配置访问
- 类型安全

### 3. 路由系统 (Routing)

高性能的路由系统，基于 Radix Tree 实现。

```go
// 基础路由
app.Router().Get("/users", controllers.UserController{}.Index)
app.Router().Post("/users", controllers.UserController{}.Store)

// 路由参数
app.Router().Get("/users/{id}", controllers.UserController{}.Show)

// 路由组
api := app.Router().Group("/api")
api.Get("/users", controllers.UserController{}.Index)
```

**特性：**

- RESTful 路由
- 路由参数
- 路由分组
- 路由中间件
- 路由缓存

### 4. 中间件系统 (Middleware)

灵活的中间件系统，支持请求处理管道。

```go
// 定义中间件
type AuthMiddleware struct{}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(c http.Context) http.Response {
        // 认证逻辑
        if !isAuthenticated(c) {
            return c.Json(map[string]string{"error": "Unauthorized"}).Status(401)
        }

        // 继续处理
        return next(c)
    }
}

// 应用中间件
app.Router().Get("/protected", handler).Use(AuthMiddleware{})
```

**特性：**

- 链式调用
- 全局中间件
- 路由中间件
- 中间件组

### 5. 控制器 (Controllers)

MVC 模式中的控制器，处理业务逻辑。

```go
type UserController struct {
    http.BaseController
}

func (c *UserController) Index() http.Response {
    users := c.DB().Table("users").Get()
    return c.Json(users)
}

func (c *UserController) Store() http.Response {
    // 验证输入
    if err := c.Validate(c.Request().Body); err != nil {
        return c.Json(err).Status(422)
    }

    // 创建用户
    user := c.DB().Table("users").Create(c.Request().Body)
    return c.Json(user).Status(201)
}
```

**特性：**

- 基础控制器
- 依赖注入
- 请求验证
- 响应格式化

### 6. 模型 (Models)

ORM 模型，提供优雅的数据库操作。

```go
type User struct {
    orm.Model
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"` // 隐藏字段
}

// 查询
users := User{}.Where("active", true).Get()

// 创建
user := User{
    Name:  "John Doe",
    Email: "john@example.com",
}
user.Save()

// 更新
user.Name = "Jane Doe"
user.Update()

// 删除
user.Delete()
```

**特性：**

- 自动时间戳
- 软删除
- 关联关系
- 模型钩子
- 字段隐藏

### 7. 数据库 (Database)

数据库抽象层，支持多种数据库。

```go
// 查询构建器
users := app.DB().Table("users").
    Select("id", "name", "email").
    Where("active", true).
    OrderBy("created_at", "desc").
    Get()

// 事务
app.DB().Transaction(func(tx *database.Transaction) {
    tx.Table("users").Create(userData)
    tx.Table("profiles").Create(profileData)
})
```

**特性：**

- 查询构建器
- 事务支持
- 连接池
- 多数据库支持

### 8. 缓存系统 (Cache)

多驱动的缓存系统。

```go
// 设置缓存
app.Cache().Set("user:1", userData, 3600)

// 获取缓存
user := app.Cache().Get("user:1")

// 删除缓存
app.Cache().Delete("user:1")

// 缓存标签
app.Cache().Tags("users").Set("user:1", userData)
app.Cache().Tags("users").Flush()
```

**特性：**

- 多驱动支持 (Redis, Memory, File)
- 缓存标签
- 自动过期
- 缓存键前缀

### 9. 队列系统 (Queue)

异步任务处理系统。

```go
// 推送任务
app.Queue().Push(&SendEmailJob{
    To:      "user@example.com",
    Subject: "Welcome",
    Body:    "Welcome to our platform!",
})

// 处理任务
type SendEmailJob struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

func (j *SendEmailJob) Handle() error {
    return sendEmail(j.To, j.Subject, j.Body)
}
```

**特性：**

- 多驱动支持
- 任务重试
- 延迟任务
- 任务优先级

### 10. 事件系统 (Events)

事件驱动架构支持。

```go
// 注册事件监听器
app.Events().Listen("user.registered", func(event *UserRegisteredEvent) {
    // 发送欢迎邮件
    sendWelcomeEmail(event.User)
})

// 触发事件
app.Events().Dispatch(&UserRegisteredEvent{
    User: user,
})
```

**特性：**

- 事件监听
- 事件广播
- 异步事件
- 事件队列

## 🔄 服务提供者

服务提供者用于注册和启动框架服务。

```go
type DatabaseServiceProvider struct{}

func (p *DatabaseServiceProvider) Register(app *framework.Application) {
    app.Container().Singleton("db", func() interface{} {
        return database.NewConnection(app.Config())
    })
}

func (p *DatabaseServiceProvider) Boot(app *framework.Application) {
    // 启动时的初始化逻辑
}
```

## 🎨 模板引擎

内置模板引擎，支持视图渲染。

```go
// 渲染视图
return c.View("users.index", map[string]interface{}{
    "users": users,
    "title": "User List",
})

// 视图模板 (resources/views/users/index.html)
{{define "users.index"}}
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
</head>
<body>
    <h1>{{.title}}</h1>
    {{range .users}}
        <div>{{.name}} - {{.email}}</div>
    {{end}}
</body>
</html>
{{end}}
```

## 🔐 认证授权

完整的用户认证和权限控制系统。

```go
// 用户认证
if app.Auth().Attempt(credentials) {
    // 登录成功
    user := app.Auth().User()
    return c.Json(user)
}

// 权限检查
if app.Auth().User().Can("edit-posts") {
    // 允许编辑文章
}

// 中间件保护
app.Router().Get("/admin", handler).Use(middleware.Auth{})
```

## ✅ 验证系统

强大的数据验证系统。

```go
// 验证规则
rules := map[string]string{
    "name":  "required|string|max:255",
    "email": "required|email|unique:users",
    "age":   "integer|min:18|max:100",
}

// 验证请求
if err := c.Validate(c.Request().Body, rules); err != nil {
    return c.Json(err).Status(422)
}
```

## 🧪 测试支持

完整的测试框架支持。

```go
func TestUserController(t *testing.T) {
    app := framework.NewApplication()

    // 设置测试数据库
    app.Config().Set("database.connection", "sqlite")
    app.Config().Set("database.database", ":memory:")

    // 运行迁移
    app.DB().Migrate()

    // 测试请求
    req := httptest.NewRequest("GET", "/users", nil)
    w := httptest.NewRecorder()

    app.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

## 🚀 命令行工具

Artisan 命令行工具，提供代码生成和项目管理功能。

```bash
# 生成控制器
go run cmd/artisan/main.go make:controller User

# 生成模型
go run cmd/artisan/main.go make:model User

# 生成迁移
go run cmd/artisan/main.go make:migration create_users_table

# 运行迁移
go run cmd/artisan/main.go migrate:run

# 清除缓存
go run cmd/artisan/main.go cache:clear
```

## 📊 性能监控

内置性能监控和指标收集。

```go
// 记录指标
app.Metrics().Counter("http_requests_total").Inc()
app.Metrics().Histogram("http_request_duration").Observe(duration)

// 健康检查
app.Router().Get("/health", func(c http.Context) http.Response {
    return c.Json(map[string]string{
        "status": "healthy",
        "uptime": app.Uptime().String(),
    })
})
```

## 🔒 安全特性

内置安全防护功能。

```go
// CSRF 保护
app.Router().Use(middleware.CSRF{})

// XSS 防护
app.Router().Use(middleware.XSS{})

// SQL 注入防护
// 通过参数化查询自动防护

// 安全头
app.Router().Use(middleware.SecurityHeaders{})
```

## 📚 下一步

了解基础概念后，建议深入学习：

1. [应用容器](container.md) - 依赖注入和服务管理
2. [路由系统](routing.md) - URL 路由和参数处理
3. [中间件](middleware.md) - 请求处理管道
4. [ORM](orm.md) - 数据库操作和模型
5. [认证授权](auth.md) - 用户认证和权限控制
6. [API 开发](api.md) - RESTful API 开发
7. [测试指南](testing.md) - 单元测试和集成测试

---

这些核心概念构成了 Laravel-Go Framework 的基础。掌握这些概念将帮助你更好地使用框架进行开发！ 🚀
