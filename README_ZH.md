# Laravel-Go Framework

[English](README.md) | [中文](README_ZH.md)

基于 Laravel 设计思路的 Go 语言开发框架，旨在为开发者提供优雅、高效的开发体验。本项目包含强大的命令行脚手架工具（类似 Laravel 的 Artisan），用于生成 Laravel-Go 应用程序。

## 特性

- 🚀 **高性能**: 基于 Go 语言的高性能特性
- 🎯 **优雅语法**: 借鉴 Laravel 的优雅设计理念
- 🔧 **完整工具链**: 包含命令行工具、ORM、模板引擎等
- 🛡️ **安全可靠**: 内置安全特性和最佳实践
- 📦 **开箱即用**: 完整的 Web 开发、API 和微服务支持
- 🐳 **容器化**: 支持 Docker 和 Kubernetes 部署

## 快速开始

### 安装

```bash
# 克隆项目
git clone git@github.com:coien1983/laravel-go.git
cd laravel-go

# 安装依赖
go mod download

# 构建 largo 命令
make build

# 查看可用命令
./bin/largo

# 初始化新的 Laravel-Go 项目
./bin/largo init

# 或者全局安装（可选）
make install
largo init
```

### 创建新项目

```bash
# 使用框架命令行工具创建新项目
laravel-go new my-project
cd my-project

# 启动开发服务器
laravel-go serve
```

## 项目结构

```
laravel-go-project/
├── app/
│   ├── Console/
│   │   └── Commands/
│   ├── Http/
│   │   ├── Controllers/
│   │   ├── Middleware/
│   │   └── Requests/
│   ├── Models/
│   ├── Services/
│   └── Providers/
├── bootstrap/
│   ├── app.go
│   └── providers.go
├── config/
│   ├── app.go
│   ├── database.go
│   ├── cache.go
│   └── queue.go
├── database/
│   ├── migrations/
│   └── seeders/
├── public/
│   ├── index.go
│   ├── css/
│   ├── js/
│   └── images/
├── resources/
│   ├── views/
│   ├── lang/
│   └── assets/
├── routes/
│   ├── web.go
│   ├── api.go
│   └── console.go
├── storage/
│   ├── logs/
│   ├── cache/
│   └── uploads/
├── tests/
├── vendor/
├── .env
├── .env.example
├── go.mod
├── go.sum
└── main.go
```

## 核心功能

### 路由系统

```go
// routes/web.go
package routes

import "laravel-go/framework/routing"

func WebRoutes(router routing.Router) {
    router.Get("/", "HomeController@index")
    router.Get("/users", "UserController@index")
    router.Post("/users", "UserController@store")

    router.Group("/api", func(router routing.Router) {
        router.Get("/users", "Api\\UserController@index")
        router.Post("/users", "Api\\UserController@store")
    }).Middleware("auth")
}
```

### 控制器

```go
// app/Http/Controllers/UserController.go
package controllers

import (
    "laravel-go/framework/http"
    "laravel-go/app/Models/User"
)

type UserController struct {
    http.Controller
}

func (c *UserController) Index(request http.Request) http.Response {
    users := User::all()
    return c.Json(users)
}

func (c *UserController) Store(request http.Request) http.Response {
    user := User::create(request.Validate(map[string]string{
        "name":  "required|string|max:255",
        "email": "required|email|unique:users",
    }))

    return c.Json(user, 201)
}
```

### 模型

```go
// app/Models/User.go
package models

import "laravel-go/framework/database"

type User struct {
    database.Model
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"-" gorm:"not null"`
}

func (u *User) TableName() string {
    return "users"
}

func (u *User) Fillable() []string {
    return []string{"name", "email", "password"}
}

func (u *User) Hidden() []string {
    return []string{"password"}
}
```

### 中间件

```go
// app/Http/Middleware/AuthMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
)

type AuthMiddleware struct{}

func (m *AuthMiddleware) Handle(request http.Request, next http.Next) http.Response {
    if !request.Auth().Check() {
        return http.Response{}.Json(map[string]string{
            "error": "Unauthenticated",
        }, 401)
    }

    return next(request)
}
```

### 命令行工具

```go
// app/Console/Commands/MakeControllerCommand.go
package commands

import (
    "laravel-go/framework/console"
)

type MakeControllerCommand struct {
    console.Command
}

func (c *MakeControllerCommand) Signature() string {
    return "make:controller {name}"
}

func (c *MakeControllerCommand) Description() string {
    return "Create a new controller class"
}

func (c *MakeControllerCommand) Handle(args []string) error {
    name := args[0]
    // 生成控制器代码
    return nil
}
```

## 配置

### 环境变量

```bash
# .env
APP_NAME=Laravel-Go
APP_ENV=local
APP_DEBUG=true
APP_URL=http://localhost:8080
APP_KEY=your-secret-key

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel_go
DB_USERNAME=root
DB_PASSWORD=

CACHE_DRIVER=redis
QUEUE_CONNECTION=redis
SESSION_DRIVER=redis
```

### 配置文件

```go
// config/app.go
package config

type AppConfig struct {
    Name        string `env:"APP_NAME" default:"Laravel-Go"`
    Environment string `env:"APP_ENV" default:"production"`
    Debug       bool   `env:"APP_DEBUG" default:"false"`
    URL         string `env:"APP_URL" default:"http://localhost"`
    Timezone    string `env:"APP_TIMEZONE" default:"UTC"`
    Locale      string `env:"APP_LOCALE" default:"en"`
    Key         string `env:"APP_KEY"`
}
```

## 部署

### Docker

```dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: laravel-go-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: laravel-go-app
  template:
    metadata:
      labels:
        app: laravel-go-app
    spec:
      containers:
        - name: app
          image: laravel-go-app:latest
          ports:
            - containerPort: 8080
          env:
            - name: APP_ENV
              value: "production"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test ./app/Http/Controllers

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 贡献

欢迎贡献代码！请阅读 [贡献指南](CONTRIBUTING.md) 了解详情。

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 支持

- 📖 [文档](https://laravel-go.dev)
- 🐛 [问题反馈](https://github.com/coien1983/laravel-go/issues)
- 📧 [邮件支持](mailto:coien1983@126.com)

## 致谢

感谢 Laravel 框架的启发，以及所有为 Go 生态系统做出贡献的开发者。
