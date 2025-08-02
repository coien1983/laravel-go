# Laravel-Go Framework

[English](README.md) | [中文](README_ZH.md)

A Go language development framework based on Laravel design principles, designed to provide developers with an elegant and efficient development experience. This project includes a powerful command-line scaffolding tool (similar to Laravel's Artisan) for generating Laravel-Go applications.

## Features

- 🚀 **High Performance**: Leveraging Go language's high-performance characteristics
- 🎯 **Elegant Syntax**: Drawing inspiration from Laravel's elegant design philosophy
- 🔧 **Complete Toolchain**: Including command-line tools, ORM, template engine, and more
- 🛡️ **Secure & Reliable**: Built-in security features and best practices
- 📦 **Ready to Use**: Complete support for Web development, APIs, and microservices
- 🐳 **Containerized**: Support for Docker and Kubernetes deployment

## Quick Start

### Installation

```bash
# Clone the project
git clone git@github.com:coien1983/laravel-go.git
cd laravel-go

# Install dependencies
go mod download

# View available commands
go run cmd/artisan/main.go

# Initialize a new Laravel-Go project
go run cmd/artisan/main.go init
```

### Create New Project

```bash
# Use framework command-line tool to create new project
laravel-go new my-project
cd my-project

# Start development server
laravel-go serve
```

## Project Structure

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

## Core Features

### Routing System

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

### Controllers

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

### Models

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

### Middleware

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

### Command Line Tools

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
    // Generate controller code
    return nil
}
```

## Configuration

### Environment Variables

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

### Configuration Files

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

## Deployment

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

## Testing

```bash
# Run all tests
go test ./...

# Run specific tests
go test ./app/Http/Controllers

# Generate test coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

We welcome contributions! Please read the [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](https://laravel-go.dev)
- 🐛 [Issue Tracker](https://github.com/coien1983/laravel-go/issues)
- 📧 [Email Support](mailto:coien1983@126.com)

## Acknowledgments

Thanks to the Laravel framework for inspiration, and to all developers who have contributed to the Go ecosystem.
