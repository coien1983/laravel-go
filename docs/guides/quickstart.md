# 快速开始

本指南将帮助你在 5 分钟内创建第一个 Laravel-Go 应用。

## 🚀 创建第一个应用

### 1. 初始化项目

```bash
# 使用 Makefile 初始化项目
make init-custom
# 输入项目名称: my-first-app
```

或者手动创建：

```bash
# 创建项目目录
mkdir my-first-app
cd my-first-app

# 初始化 Go 模块
go mod init my-first-app

# 添加框架依赖
go get laravel-go/framework
```

### 2. 创建主程序

创建 `main.go` 文件：

```go
package main

import (
    "laravel-go/framework"
    "laravel-go/framework/http"
)

func main() {
    // 创建应用实例
    app := framework.NewApplication()

    // 注册路由
    app.Router().Get("/", func(c http.Context) http.Response {
        return c.Json(map[string]string{
            "message": "Hello Laravel-Go!",
        })
    })

    // 启动服务器
    app.Run(":8080")
}
```

### 3. 运行应用

```bash
# 运行应用
go run main.go

# 或使用 Makefile
make run
```

### 4. 测试应用

```bash
# 测试 API
curl http://localhost:8080/
# 输出: {"message":"Hello Laravel-Go!"}
```

## 📝 创建完整的 CRUD 应用

### 1. 生成用户管理组件

```bash
# 生成完整的用户 CRUD 组件
make crud
# 输入资源名称: user
```

这将生成：

- 用户控制器
- 用户模型
- 数据库迁移
- 单元测试
- 集成测试

### 2. 配置数据库

创建 `.env` 文件：

```env
APP_NAME=my-first-app
APP_ENV=development
APP_DEBUG=true

DB_CONNECTION=sqlite
DB_DATABASE=database/app.db
```

### 3. 运行迁移

```bash
# 运行数据库迁移
go run cmd/artisan/main.go migrate:run
```

### 4. 启动应用

```bash
# 启动应用
go run main.go
```

### 5. 测试 CRUD 功能

```bash
# 创建用户
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# 获取用户列表
curl http://localhost:8080/users

# 获取单个用户
curl http://localhost:8080/users/1

# 更新用户
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'

# 删除用户
curl -X DELETE http://localhost:8080/users/1
```

## 🎯 创建 API 服务

### 1. 生成 API 组件

```bash
# 生成 API 组件
make api
# 输入资源名称: product
```

### 2. 创建 API 路由

在 `main.go` 中添加 API 路由：

```go
package main

import (
    "laravel-go/framework"
    "laravel-go/framework/http"
    "my-first-app/app/controllers"
)

func main() {
    app := framework.NewApplication()

    // Web 路由
    app.Router().Get("/", func(c http.Context) http.Response {
        return c.Json(map[string]string{
            "message": "Welcome to Laravel-Go API",
        })
    })

    // API 路由组
    api := app.Router().Group("/api")
    {
        // 产品 API
        api.Get("/products", controllers.ProductController{}.Index)
        api.Post("/products", controllers.ProductController{}.Store)
        api.Get("/products/{id}", controllers.ProductController{}.Show)
        api.Put("/products/{id}", controllers.ProductController{}.Update)
        api.Delete("/products/{id}", controllers.ProductController{}.Destroy)

        // 用户 API
        api.Get("/users", controllers.UserController{}.Index)
        api.Post("/users", controllers.UserController{}.Store)
        api.Get("/users/{id}", controllers.UserController{}.Show)
        api.Put("/users/{id}", controllers.UserController{}.Update)
        api.Delete("/users/{id}", controllers.UserController{}.Destroy)
    }

    app.Run(":8080")
}
```

### 3. 测试 API

```bash
# 测试 API 根路径
curl http://localhost:8080/api

# 创建产品
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone 15","price":999.99,"description":"Latest iPhone"}'

# 获取产品列表
curl http://localhost:8080/api/products
```

## 🔧 添加中间件

### 1. 生成中间件

```bash
# 生成认证中间件
make middleware
# 输入中间件名称: auth
```

### 2. 实现中间件

编辑 `app/middleware/auth_middleware.go`：

```go
package middleware

import (
    "laravel-go/framework/http"
)

type AuthMiddleware struct{}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(c http.Context) http.Response {
        // 检查认证令牌
        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return c.Json(map[string]string{
                "error": "Unauthorized",
            }).Status(401)
        }

        // 验证令牌逻辑...

        // 继续处理请求
        return next(c)
    }
}
```

### 3. 应用中间件

```go
// 在路由组中应用中间件
api := app.Router().Group("/api")
api.Use(middleware.AuthMiddleware{})

// 或在单个路由上应用
api.Get("/protected", func(c http.Context) http.Response {
    return c.Json(map[string]string{"message": "Protected resource"})
}).Use(middleware.AuthMiddleware{})
```

## 🎨 创建 Web 页面

### 1. 创建控制器

```bash
# 生成页面控制器
make controller
# 输入控制器名称: page
```

### 2. 实现控制器

编辑 `app/controllers/page_controller.go`：

```go
package controllers

import (
    "laravel-go/framework/http"
)

type PageController struct {
    http.BaseController
}

func (c *PageController) Home() http.Response {
    return c.View("home", map[string]interface{}{
        "title": "Welcome to Laravel-Go",
        "message": "This is your first Laravel-Go application!",
    })
}

func (c *PageController) About() http.Response {
    return c.View("about", map[string]interface{}{
        "title": "About Us",
        "content": "Learn more about our application.",
    })
}
```

### 3. 创建视图模板

创建 `resources/views/home.html`：

```html
<!DOCTYPE html>
<html>
  <head>
    <title>{{.title}}</title>
    <link
      href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css"
      rel="stylesheet"
    />
  </head>
  <body class="bg-gray-100">
    <div class="container mx-auto px-4 py-8">
      <h1 class="text-4xl font-bold text-gray-800 mb-4">{{.title}}</h1>
      <p class="text-lg text-gray-600">{{.message}}</p>

      <div class="mt-8">
        <a
          href="/users"
          class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Manage Users
        </a>
        <a
          href="/api/products"
          class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 ml-4"
        >
          API Products
        </a>
      </div>
    </div>
  </body>
</html>
```

### 4. 添加路由

```go
// 添加页面路由
app.Router().Get("/", controllers.PageController{}.Home)
app.Router().Get("/about", controllers.PageController{}.About)
```

## 🧪 编写测试

### 1. 运行现有测试

```bash
# 运行所有测试
make test-all

# 运行特定测试
go test ./tests/ -v
```

### 2. 创建新测试

```bash
# 生成测试文件
make test
# 输入测试名称: page
```

### 3. 编写测试代码

编辑 `tests/page_test.go`：

```go
package tests

import (
    "net/http/httptest"
    "testing"
    "laravel-go/framework"
    "laravel-go/framework/http"
)

func TestPageController(t *testing.T) {
    app := framework.NewApplication()

    // 注册路由
    app.Router().Get("/", func(c http.Context) http.Response {
        return c.Json(map[string]string{
            "message": "Hello from test",
        })
    })

    // 创建测试请求
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()

    // 执行请求
    app.ServeHTTP(w, req)

    // 验证响应
    if w.Code != 200 {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

## 🚀 部署应用

### 1. 生成部署配置

```bash
# 生成 Docker 配置
make docker-custom
# 输入应用名称: my-first-app
# 输入端口: 8080
# 输入环境: production

# 生成 Kubernetes 配置
make k8s-custom
# 输入应用名称: my-first-app
# 输入副本数: 3
# 输入端口: 8080
# 输入命名空间: production
```

### 2. 构建和部署

```bash
# 构建 Docker 镜像
make docker-build

# 启动 Docker 服务
make docker-compose-up

# 或部署到 Kubernetes
make k8s-apply
```

## 📚 下一步

恭喜！你已经成功创建了第一个 Laravel-Go 应用。接下来建议学习：

1. [基础概念](concepts.md) - 了解框架核心概念
2. [应用容器](container.md) - 学习依赖注入
3. [路由系统](routing.md) - 深入路由功能
4. [ORM](orm.md) - 数据库操作
5. [中间件](middleware.md) - 请求处理
6. [认证授权](auth.md) - 用户认证
7. [API 开发](api.md) - RESTful API
8. [测试指南](testing.md) - 测试策略
9. [部署指南](deployment.md) - 生产部署

## 🎯 项目结构

你的项目现在应该包含：

```
my-first-app/
├── main.go                 # 应用入口
├── go.mod                  # Go 模块文件
├── .env                    # 环境配置
├── app/
│   ├── controllers/        # 控制器
│   ├── models/            # 模型
│   └── middleware/        # 中间件
├── resources/
│   └── views/             # 视图模板
├── tests/                 # 测试文件
├── database/
│   └── migrations/        # 数据库迁移
├── Dockerfile             # Docker 配置
├── docker-compose.yml     # Docker Compose
└── k8s/                   # Kubernetes 配置
```

## 🆘 需要帮助？

如果遇到问题：

- 📖 查看 [完整文档](../README.md)
- 💬 加入 [社区讨论](https://github.com/your-org/laravel-go/discussions)
- 🐛 提交 [问题反馈](https://github.com/your-org/laravel-go/issues)

---

你已经成功入门 Laravel-Go Framework！继续探索更多功能吧！ 🚀
