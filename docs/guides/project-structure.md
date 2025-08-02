# 项目结构指南

## 📁 目录结构概览

Laravel-Go Framework 采用清晰的项目结构，遵循 Go 语言最佳实践和 Laravel 的设计理念。

```
laravel-go/
├── app/                    # 应用程序核心代码
│   ├── Console/           # 命令行命令
│   ├── Http/              # HTTP 层
│   │   ├── Controllers/   # 控制器
│   │   ├── Middleware/    # 中间件
│   │   └── Requests/      # 请求验证
│   ├── Models/            # 数据模型
│   ├── Providers/         # 服务提供者
│   └── Services/          # 业务服务
├── bootstrap/             # 框架启动文件
├── cmd/                   # 可执行文件
│   └── artisan/           # Artisan 命令行工具
├── config/                # 配置文件
├── database/              # 数据库相关
│   ├── migrations/        # 数据库迁移
│   └── seeders/           # 数据填充
├── docs/                  # 项目文档
├── examples/              # 示例代码
├── framework/             # 框架核心代码
├── public/                # 公共资源
├── resources/             # 资源文件
├── routes/                # 路由定义
├── storage/               # 存储目录
├── tests/                 # 测试文件
├── go.mod                 # Go 模块文件
├── go.sum                 # 依赖校验文件
├── Makefile               # 构建脚本
└── README.md              # 项目说明
```

## 🏗️ 核心目录详解

### `/app` - 应用程序代码

应用程序的核心代码目录，包含所有业务逻辑。

#### `/app/Console/Commands/`
```go
// 自定义 Artisan 命令
type CreateUserCommand struct {
    console.Command
}

func (c *CreateUserCommand) Handle() error {
    // 命令逻辑
    return nil
}
```

#### `/app/Http/Controllers/`
```go
// HTTP 控制器
type UserController struct {
    http.Controller
}

func (c *UserController) Index() http.Response {
    // 控制器逻辑
    return c.Json(map[string]interface{}{
        "users": []string{},
    })
}
```

#### `/app/Http/Middleware/`
```go
// 自定义中间件
type AuthMiddleware struct {
    http.Middleware
}

func (m *AuthMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 中间件逻辑
    return next(request)
}
```

#### `/app/Models/`
```go
// 数据模型
type User struct {
    database.Model
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"`
}
```

#### `/app/Services/`
```go
// 业务服务
type UserService struct {
    db *database.Connection
}

func (s *UserService) CreateUser(data map[string]interface{}) (*User, error) {
    // 业务逻辑
    return user, nil
}
```

### `/config` - 配置文件

应用程序的配置文件目录。

```go
// config/app.go
package config

type App struct {
    Name    string `env:"APP_NAME" default:"Laravel-Go"`
    Env     string `env:"APP_ENV" default:"local"`
    Debug   bool   `env:"APP_DEBUG" default:"true"`
    URL     string `env:"APP_URL" default:"http://localhost"`
    Timezone string `env:"APP_TIMEZONE" default:"UTC"`
}
```

### `/database` - 数据库相关

数据库迁移和种子文件。

#### `/database/migrations/`
```go
// 用户表迁移
type CreateUsersTable struct {
    migration.Migration
}

func (m *CreateUsersTable) Up() error {
    return m.Schema.CreateTable("users", func(table *database.Blueprint) {
        table.Id("id")
        table.String("name")
        table.String("email").Unique()
        table.String("password")
        table.Timestamps()
    })
}
```

#### `/database/seeders/`
```go
// 数据填充
type UserSeeder struct {
    seeder.Seeder
}

func (s *UserSeeder) Run() error {
    // 填充数据
    return nil
}
```

### `/framework` - 框架核心

框架的核心代码，包含所有功能模块。

#### 主要模块
- `api/` - API 资源处理
- `auth/` - 认证授权
- `cache/` - 缓存系统
- `config/` - 配置管理
- `console/` - 命令行工具
- `container/` - 依赖注入容器
- `database/` - 数据库操作
- `event/` - 事件系统
- `http/` - HTTP 处理
- `queue/` - 队列系统
- `routing/` - 路由系统
- `validation/` - 数据验证

### `/routes` - 路由定义

应用程序的路由配置。

```go
// routes/web.go
package routes

func WebRoutes(router *routing.Router) {
    router.Get("/", func(request http.Request) http.Response {
        return http.Response{
            Body: "Welcome to Laravel-Go!",
        }
    })
    
    router.Get("/users", &UserController{}, "Index")
}
```

### `/storage` - 存储目录

应用程序的存储文件。

```
storage/
├── cache/     # 缓存文件
├── logs/      # 日志文件
└── uploads/   # 上传文件
```

### `/tests` - 测试文件

应用程序的测试代码。

```go
// tests/basic_test.go
package tests

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestBasicFunctionality(t *testing.T) {
    // 测试逻辑
    assert.True(t, true)
}
```

## 🔧 配置文件说明

### `go.mod`
```go
module laravel-go

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/go-sql-driver/mysql v1.7.1
    github.com/redis/go-redis/v9 v9.3.0
    // 其他依赖
)
```

### `Makefile`
```makefile
# 构建命令
build:
    go build -o bin/laravel-go cmd/artisan/main.go

# 运行测试
test:
    go test ./...

# 代码格式化
fmt:
    go fmt ./...
```

## 📋 目录命名规范

### 1. 使用小写字母
- ✅ `controllers/`
- ✅ `middleware/`
- ❌ `Controllers/`
- ❌ `Middleware/`

### 2. 使用复数形式
- ✅ `controllers/`
- ✅ `models/`
- ✅ `services/`
- ❌ `controller/`
- ❌ `model/`

### 3. 使用连字符分隔
- ✅ `user-controller.go`
- ✅ `auth-middleware.go`
- ❌ `userController.go`
- ❌ `authMiddleware.go`

## 🎯 最佳实践

### 1. 模块化组织
```go
// 按功能模块组织代码
app/
├── User/              # 用户模块
│   ├── Controllers/
│   ├── Models/
│   └── Services/
├── Product/           # 产品模块
│   ├── Controllers/
│   ├── Models/
│   └── Services/
```

### 2. 依赖注入
```go
// 使用容器管理依赖
type UserController struct {
    http.Controller
    userService *UserService
}

func NewUserController(userService *UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}
```

### 3. 接口分离
```go
// 定义接口
type UserRepository interface {
    Find(id int) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id int) error
}

// 实现接口
type UserRepositoryImpl struct {
    db *database.Connection
}
```

### 4. 错误处理
```go
// 统一错误处理
func (c *UserController) Show(id int) http.Response {
    user, err := c.userService.Find(id)
    if err != nil {
        return c.JsonError(err, 404)
    }
    
    return c.Json(user)
}
```

## 🔄 扩展项目结构

### 添加新模块
```bash
# 创建新模块目录
mkdir -p app/NewModule/{Controllers,Models,Services}

# 创建控制器
touch app/NewModule/Controllers/new-module-controller.go

# 创建模型
touch app/NewModule/Models/new-module.go

# 创建服务
touch app/NewModule/Services/new-module-service.go
```

### 添加新配置
```bash
# 创建配置文件
touch config/new-module.go

# 添加环境变量
echo "NEW_MODULE_ENABLED=true" >> .env
```

## 📊 目录大小建议

### 控制器
- 单个控制器文件不超过 500 行
- 每个方法不超过 50 行
- 使用服务层处理复杂逻辑

### 模型
- 单个模型文件不超过 300 行
- 使用关联方法组织代码
- 避免在模型中放置业务逻辑

### 服务
- 单个服务文件不超过 1000 行
- 按功能方法组织代码
- 使用接口定义契约

## 🚀 性能优化

### 1. 文件组织
- 避免过深的目录嵌套
- 使用扁平化结构
- 按功能而非类型组织

### 2. 导入优化
```go
// 避免循环导入
package controllers

import (
    "app/Services"  // 导入服务层
    "app/Models"    // 导入模型层
)
```

### 3. 缓存策略
```go
// 使用缓存减少数据库查询
func (s *UserService) GetUser(id int) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    
    if cached, found := s.cache.Get(cacheKey); found {
        return cached.(*User), nil
    }
    
    user, err := s.repository.Find(id)
    if err != nil {
        return nil, err
    }
    
    s.cache.Set(cacheKey, user, time.Hour)
    return user, nil
}
```

## 📝 总结

Laravel-Go Framework 的项目结构设计遵循以下原则：

1. **清晰性**: 目录结构清晰易懂
2. **模块化**: 按功能模块组织代码
3. **可扩展性**: 易于添加新功能
4. **可维护性**: 便于维护和重构
5. **性能优化**: 考虑编译和运行性能

通过遵循这些规范和最佳实践，可以构建出高质量、可维护的 Laravel-Go 应用程序。 