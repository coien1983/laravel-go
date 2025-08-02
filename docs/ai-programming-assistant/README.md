# Laravel-Go Framework AI 编程助手配置指南

## 🚀 概述

本指南帮助开发者配置AI编程助手（如GitHub Copilot、Claude、GPT等），使其能够快速理解Laravel-Go框架的架构、编码规范和最佳实践，从而提供更准确、更符合框架标准的代码建议。

## 📋 目录

1. [AI助手配置文件](#ai助手配置文件)
2. [框架架构说明](#框架架构说明)
3. [编码规范指南](#编码规范指南)
4. [常用代码模板](#常用代码模板)
5. [最佳实践提示](#最佳实践提示)
6. [调试和测试指南](#调试和测试指南)

## 🤖 AI助手配置文件

### 1. GitHub Copilot 配置

在项目根目录创建 `.copilot/` 目录：

```bash
mkdir -p .copilot
```

#### `.copilot/settings.json`

```json
{
  "framework": "laravel-go",
  "language": "go",
  "architecture": "layered",
  "patterns": [
    "repository",
    "service",
    "controller",
    "middleware",
    "validation"
  ],
  "conventions": {
    "naming": "snake_case",
    "file_structure": "feature_based",
    "error_handling": "wrapped_errors",
    "logging": "structured"
  }
}
```

#### `.copilot/prompts.md`

```markdown
# Laravel-Go Framework 开发指南

## 框架概述
Laravel-Go 是一个受 Laravel PHP 启发的 Go Web 框架，提供完整的 Web 开发解决方案。

## 核心原则
1. **约定优于配置**：遵循框架约定，减少配置代码
2. **依赖注入**：使用容器管理依赖关系
3. **中间件模式**：请求处理管道化
4. **错误处理**：统一的错误处理机制
5. **性能监控**：内置性能监控和告警

## 项目结构
```
app/
├── Http/
│   ├── Controllers/    # 控制器层
│   ├── Middleware/     # 中间件
│   └── Requests/       # 请求验证
├── Models/             # 数据模型
├── Services/           # 业务服务层
└── Providers/          # 服务提供者

framework/              # 框架核心
├── http/              # HTTP 处理
├── database/          # 数据库操作
├── cache/             # 缓存系统
├── queue/             # 队列系统
├── events/            # 事件系统
├── performance/       # 性能监控
└── errors/            # 错误处理

config/                # 配置文件
database/              # 数据库迁移和种子
routes/                # 路由定义
storage/               # 文件存储
tests/                 # 测试文件
```

## 编码规范

### 命名约定
- 文件名：`snake_case.go`
- 结构体：`PascalCase`
- 方法：`PascalCase`
- 变量：`camelCase`
- 常量：`UPPER_SNAKE_CASE`
- 包名：`lowercase`

### 错误处理
```go
// 使用框架的错误处理
import "laravel-go/framework/errors"

// 创建错误
err := errors.New("error message")

// 包装错误
err = errors.Wrap(err, "additional context")

// 带状态码的错误
err = errors.NewWithCode(400, "bad request")
```

### 控制器模式
```go
type UserController struct {
    userService *UserService
    errorHandler errors.ErrorHandler
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
    // 1. 参数验证
    // 2. 业务逻辑
    // 3. 错误处理
    // 4. 响应返回
}
```

### 服务层模式
```go
type UserService struct {
    userRepo *UserRepository
    cache    *CacheService
    errorHandler errors.ErrorHandler
}

func (s *UserService) GetUser(id int) (*User, error) {
    // 使用安全执行包装器
    var user *User
    var err error
    
    errors.SafeExecuteWithContext(context.Background(), func() error {
        // 业务逻辑
        return nil
    })
    
    return user, err
}
```

## 常用代码模板

### 控制器模板
```go
package controllers

import (
    "net/http"
    "laravel-go/framework/errors"
    "laravel-go/app/Services"
)

type {{ControllerName}}Controller struct {
    {{serviceName}}Service *Services.{{ServiceName}}Service
    errorHandler errors.ErrorHandler
}

func New{{ControllerName}}Controller({{serviceName}}Service *Services.{{ServiceName}}Service, errorHandler errors.ErrorHandler) *{{ControllerName}}Controller {
    return &{{ControllerName}}Controller{
        {{serviceName}}Service: {{serviceName}}Service,
        errorHandler: errorHandler,
    }
}

func (c *{{ControllerName}}Controller) Index(w http.ResponseWriter, r *http.Request) {
    // 实现列表逻辑
}

func (c *{{ControllerName}}Controller) Show(w http.ResponseWriter, r *http.Request) {
    // 实现详情逻辑
}

func (c *{{ControllerName}}Controller) Store(w http.ResponseWriter, r *http.Request) {
    // 实现创建逻辑
}

func (c *{{ControllerName}}Controller) Update(w http.ResponseWriter, r *http.Request) {
    // 实现更新逻辑
}

func (c *{{ControllerName}}Controller) Destroy(w http.ResponseWriter, r *http.Request) {
    // 实现删除逻辑
}

func (c *{{ControllerName}}Controller) handleError(w http.ResponseWriter, err error) {
    processedErr := c.errorHandler.Handle(err)
    
    if appErr := errors.GetAppError(processedErr); appErr != nil {
        http.Error(w, appErr.Message, appErr.Code)
    } else {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
```

### 服务层模板
```go
package Services

import (
    "context"
    "laravel-go/framework/errors"
    "laravel-go/app/Models"
)

type {{ServiceName}}Service struct {
    {{serviceName}}Repo *{{ServiceName}}Repository
    cacheService *CacheService
    errorHandler errors.ErrorHandler
}

func New{{ServiceName}}Service({{serviceName}}Repo *{{ServiceName}}Repository, cacheService *CacheService, errorHandler errors.ErrorHandler) *{{ServiceName}}Service {
    return &{{ServiceName}}Service{
        {{serviceName}}Repo: {{serviceName}}Repo,
        cacheService: cacheService,
        errorHandler: errorHandler,
    }
}

func (s *{{ServiceName}}Service) Get{{ServiceName}}(id int) (*Models.{{ServiceName}}, error) {
    var {{serviceName}} *Models.{{ServiceName}}
    var err error
    
    errors.SafeExecuteWithContext(context.Background(), func() error {
        if id <= 0 {
            err = errors.Wrap(errors.New("invalid id"), "invalid {{serviceName}} id")
            return err
        }
        
        {{serviceName}}, err = s.{{serviceName}}Repo.FindByID(id)
        if err != nil {
            return errors.Wrap(err, "failed to get {{serviceName}}")
        }
        
        return nil
    })
    
    return {{serviceName}}, err
}

func (s *{{ServiceName}}Service) Create{{ServiceName}}({{serviceName}} *Models.{{ServiceName}}) error {
    return errors.SafeExecuteWithContext(context.Background(), func() error {
        if {{serviceName}} == nil {
            return errors.New("{{serviceName}} cannot be nil")
        }
        
        return s.{{serviceName}}Repo.Create({{serviceName}})
    })
}
```

### 模型模板
```go
package Models

import (
    "time"
    "laravel-go/framework/database"
)

type {{ModelName}} struct {
    ID        int       `json:"id" db:"id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
    // 添加其他字段
}

func (m *{{ModelName}}) TableName() string {
    return "{{table_name}}"
}

func (m *{{ModelName}}) BeforeCreate() error {
    m.CreatedAt = time.Now()
    m.UpdatedAt = time.Now()
    return nil
}

func (m *{{ModelName}}) BeforeUpdate() error {
    m.UpdatedAt = time.Now()
    return nil
}
```

### 中间件模板
```go
package middleware

import (
    "net/http"
    "laravel-go/framework/errors"
)

type {{MiddlewareName}}Middleware struct {
    errorHandler errors.ErrorHandler
}

func New{{MiddlewareName}}Middleware(errorHandler errors.ErrorHandler) *{{MiddlewareName}}Middleware {
    return &{{MiddlewareName}}Middleware{
        errorHandler: errorHandler,
    }
}

func (m *{{MiddlewareName}}Middleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 中间件逻辑
        next.ServeHTTP(w, r)
    })
}
```

## 最佳实践提示

### 1. 错误处理
- 始终使用框架的错误处理机制
- 使用 `errors.Wrap` 添加上下文信息
- 在服务层使用 `SafeExecuteWithContext`
- 在控制器层统一处理错误响应

### 2. 性能优化
- 使用缓存减少数据库查询
- 实现数据库连接池
- 使用异步处理处理耗时操作
- 监控关键性能指标

### 3. 安全性
- 验证所有用户输入
- 使用参数化查询防止SQL注入
- 实现适当的认证和授权
- 记录安全相关事件

### 4. 测试
- 为每个服务编写单元测试
- 使用集成测试验证API
- 模拟外部依赖
- 测试错误场景

## 调试和测试指南

### 调试技巧
```go
// 使用框架的日志系统
logger := &CustomLogger{}
logger.Info("debug message", map[string]interface{}{
    "user_id": 123,
    "action": "login",
})

// 使用性能监控
monitor := performance.NewPerformanceMonitor()
httpMonitor := performance.NewHTTPMonitor(monitor)
```

### 测试模式
```go
// 单元测试模板
func Test{{ServiceName}}Service_Get{{ServiceName}}(t *testing.T) {
    // 设置测试环境
    // 执行测试
    // 验证结果
}
```

## 常用命令

```bash
# 运行测试
go test ./...

# 运行特定测试
go test ./app/Services -v

# 构建项目
go build -o app cmd/main.go

# 运行项目
go run cmd/main.go

# 代码格式化
go fmt ./...

# 代码检查
go vet ./...
```

## 注意事项

1. **框架版本**：确保使用最新的稳定版本
2. **依赖管理**：使用 `go mod` 管理依赖
3. **配置管理**：使用环境变量管理配置
4. **日志记录**：使用结构化日志
5. **错误处理**：遵循框架的错误处理模式
6. **性能监控**：集成性能监控系统
7. **测试覆盖**：保持高测试覆盖率
8. **文档更新**：及时更新API文档

## 获取帮助

- 查看框架文档：`docs/`
- 查看示例代码：`examples/`
- 查看测试用例：`tests/`
- 查看最佳实践：`docs/best-practices/`
```

### 2. VS Code 配置

#### `.vscode/settings.json`

```json
{
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.lintFlags": [
    "--fast"
  ],
  "go.testFlags": [
    "-v"
  ],
  "go.buildTags": "",
  "go.toolsManagement.checkForUpdates": "local",
  "go.useLanguageServer": true,
  "go.languageServerExperimentalFeatures": {
    "diagnostics": true,
    "documentLink": true
  },
  "files.associations": {
    "*.go": "go"
  },
  "emmet.includeLanguages": {
    "go": "html"
  },
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  }
}
```

#### `.vscode/extensions.json`

```json
{
  "recommendations": [
    "golang.go",
    "ms-vscode.go",
    "ms-vscode.vscode-json",
    "redhat.vscode-yaml",
    "ms-vscode.vscode-markdown"
  ]
}
```

### 3. 项目级配置

#### `go.mod` 模板

```go
module your-project-name

go 1.21

require (
    laravel-go/framework v0.1.0
    github.com/gorilla/mux v1.8.0
    github.com/go-sql-driver/mysql v1.7.1
    github.com/redis/go-redis/v9 v9.0.5
    github.com/stretchr/testify v1.8.4
)

require (
    // 其他依赖
)
```

#### `.gitignore`

```gitignore
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs
*.log
logs/

# Environment files
.env
.env.local
.env.*.local

# Build output
build/
dist/

# Temporary files
tmp/
temp/
```

## 🎯 使用指南

### 1. 初始化项目

```bash
# 1. 创建项目目录
mkdir my-laravel-go-project
cd my-laravel-go-project

# 2. 复制配置文件
cp -r laravel-go/docs/ai-programming-assistant/.copilot ./
cp -r laravel-go/docs/ai-programming-assistant/.vscode ./

# 3. 初始化 Go 模块
go mod init my-project-name

# 4. 添加框架依赖
go get laravel-go/framework

# 5. 创建项目结构
mkdir -p app/{Http/{Controllers,Middleware,Requests},Models,Services,Providers}
mkdir -p config database/{migrations,seeders} routes storage/{cache,logs,uploads} tests
```

### 2. AI助手提示词

在与AI助手对话时，使用以下提示词：

```
我正在使用 Laravel-Go Framework 开发项目。这是一个受 Laravel PHP 启发的 Go Web 框架。

请按照以下规范为我生成代码：

1. 使用框架的错误处理机制（errors.Wrap, SafeExecuteWithContext）
2. 遵循分层架构（Controller -> Service -> Repository）
3. 使用框架的中间件模式
4. 实现统一的错误处理
5. 集成性能监控
6. 使用结构化日志
7. 遵循 Go 编码规范

项目结构：
- app/Http/Controllers/ - 控制器层
- app/Services/ - 业务服务层
- app/Models/ - 数据模型
- framework/ - 框架核心

请为 [具体功能] 生成符合框架标准的代码。
```

### 3. 代码生成模板

#### 快速生成控制器

```bash
# 使用 AI 助手生成控制器
echo "请为 User 模块生成符合 Laravel-Go Framework 标准的控制器，包含 CRUD 操作和错误处理"
```

#### 快速生成服务层

```bash
# 使用 AI 助手生成服务层
echo "请为 User 模块生成符合 Laravel-Go Framework 标准的服务层，包含业务逻辑和缓存处理"
```

#### 快速生成模型

```bash
# 使用 AI 助手生成模型
echo "请为 User 模型生成符合 Laravel-Go Framework 标准的数据模型，包含数据库操作和验证"
```

## 📚 学习资源

### 1. 框架文档
- [快速开始](../guides/quickstart.md)
- [核心概念](../guides/concepts.md)
- [错误处理](../guides/error-handling.md)
- [性能监控](../guides/performance.md)

### 2. 示例代码
- [API示例](../examples/api_example/)
- [博客示例](../examples/blog_example/)
- [微服务示例](../examples/microservice_example/)

### 3. 最佳实践
- [编码规范](../best-practices/coding-standards.md)
- [错误处理](../best-practices/error-handling.md)
- [性能优化](../best-practices/performance.md)

## 🔧 工具集成

### 1. 代码生成器

```bash
# 安装代码生成工具
go install laravel-go/tools/dev-tools/code-generator

# 生成控制器
code-generator controller User

# 生成服务
code-generator service User

# 生成模型
code-generator model User
```

### 2. 调试工具

```bash
# 安装调试工具
go install laravel-go/tools/dev-tools/debug-tool

# 启动调试服务器
debug-tool serve
```

### 3. 性能分析器

```bash
# 安装性能分析工具
go install laravel-go/tools/dev-tools/performance-analyzer

# 分析性能
performance-analyzer analyze
```

## 🎉 总结

通过以上配置，AI编程助手将能够：

1. **理解框架架构**：快速掌握分层架构和设计模式
2. **遵循编码规范**：生成符合框架标准的代码
3. **应用最佳实践**：使用推荐的错误处理和性能优化方案
4. **提高开发效率**：减少重复代码，专注于业务逻辑
5. **保证代码质量**：生成可测试、可维护的代码

记住，AI助手是强大的工具，但理解框架原理和最佳实践仍然是开发者的核心能力。通过合理配置和使用，AI助手将成为你开发 Laravel-Go 项目的得力助手！ 