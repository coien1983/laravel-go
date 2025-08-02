# Laravel-Go Framework 配置模块

## 📝 模块概览

配置模块是 Laravel-Go Framework 的核心组件之一，提供了完整的配置管理功能，包括环境变量、配置文件、配置验证等。

## 🚀 功能特性

- ✅ 环境变量管理
- ✅ 配置文件加载 (JSON, YAML)
- ✅ 配置验证
- ✅ 默认配置
- ✅ 配置热重载
- ✅ 类型安全
- ✅ 嵌套配置支持

## 📁 文件结构

```
framework/config/
├── config.go      # 核心配置管理器
├── app.go         # 应用配置结构
├── init.go        # 配置初始化工具
├── env.example    # 环境变量示例
└── README.md      # 本文档
```

## 🚀 快速开始

### 1. 基本使用

```go
package main

import (
    "laravel-go/framework/config"
)

func main() {
    // 创建配置管理器
    cfg := config.NewConfig()

    // 加载环境变量
    cfg.LoadEnv()

    // 加载配置文件
    cfg.LoadFromFile("config/app.json")

    // 获取配置值
    appName := cfg.GetString("app.name", "Laravel-Go")
    debug := cfg.GetBool("app.debug", false)

    fmt.Printf("应用名称: %s, 调试模式: %t\n", appName, debug)
}
```

### 2. 项目初始化

```go
package main

import (
    "laravel-go/framework/config"
)

func main() {
    // 初始化项目配置
    if err := config.InitConfig(); err != nil {
        log.Fatalf("初始化配置失败: %v", err)
    }

    fmt.Println("✅ 项目配置初始化完成")
}
```

### 3. 使用默认配置

```go
package main

import (
    "laravel-go/framework/config"
)

func main() {
    // 加载应用配置
    appConfig := config.LoadAppConfig()
    fmt.Printf("应用名称: %s\n", appConfig.Name)

    // 加载数据库配置
    dbConfig := config.LoadDatabaseConfig()
    fmt.Printf("默认数据库: %s\n", dbConfig.Default)

    // 加载缓存配置
    cacheConfig := config.LoadCacheConfig()
    fmt.Printf("默认缓存: %s\n", cacheConfig.Default)
}
```

## 🔧 配置结构

### 1. 应用配置 (AppConfig)

```go
type AppConfig struct {
    Name      string   `json:"name"`
    Version   string   `json:"version"`
    Env       string   `json:"env"`
    Debug     bool     `json:"debug"`
    URL       string   `json:"url"`
    Port      string   `json:"port"`
    Timezone  string   `json:"timezone"`
    Locale    string   `json:"locale"`
    Key       string   `json:"key"`
    Providers []string `json:"providers"`
}
```

### 2. 数据库配置 (DatabaseConfig)

```go
type DatabaseConfig struct {
    Default     string                     `json:"default"`
    Connections map[string]ConnectionConfig `json:"connections"`
    Migrations  string                     `json:"migrations"`
    Redis       RedisConfig                `json:"redis"`
}
```

### 3. 缓存配置 (CacheConfig)

```go
type CacheConfig struct {
    Default string                  `json:"default"`
    Stores  map[string]StoreConfig `json:"stores"`
    Prefix  string                 `json:"prefix"`
}
```

### 4. 队列配置 (QueueConfig)

```go
type QueueConfig struct {
    Default     string                    `json:"default"`
    Connections map[string]QueueConnection `json:"connections"`
    Failed      FailedJobConfig           `json:"failed"`
}
```

### 5. 会话配置 (SessionConfig)

```go
type SessionConfig struct {
    Driver        string `json:"driver"`
    Lifetime      int    `json:"lifetime"`
    ExpireOnClose bool   `json:"expire_on_close"`
    Encrypt       bool   `json:"encrypt"`
    Files         string `json:"files"`
    Cookie        string `json:"cookie"`
    Path          string `json:"path"`
    Domain        string `json:"domain"`
    Secure        bool   `json:"secure"`
    HTTPOnly      bool   `json:"http_only"`
    SameSite      string `json:"same_site"`
}
```

### 6. 日志配置 (LoggingConfig)

```go
type LoggingConfig struct {
    Default      string                    `json:"default"`
    Deprecations DeprecationConfig         `json:"deprecations"`
    Channels     map[string]ChannelConfig  `json:"channels"`
}
```

## 📋 环境变量

### 应用配置

```env
APP_NAME="Laravel-Go"
APP_VERSION=1.0.0
APP_ENV=production
APP_DEBUG=false
APP_URL=http://localhost:8080
APP_PORT=8080
APP_TIMEZONE=UTC
APP_LOCALE=en
APP_KEY=
```

### 数据库配置

```env
DB_CONNECTION=sqlite
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel_go
DB_USERNAME=root
DB_PASSWORD=
DB_MIGRATIONS=database/migrations
```

### Redis 配置

```env
REDIS_HOST=127.0.0.1
REDIS_PASSWORD=null
REDIS_PORT=6379
REDIS_DB=0
REDIS_CACHE_DB=1
REDIS_QUEUE_DB=2
```

### 缓存配置

```env
CACHE_DRIVER=file
CACHE_PATH=storage/framework/cache/data
CACHE_PREFIX=laravel_go_cache
```

### 队列配置

```env
QUEUE_CONNECTION=sync
REDIS_QUEUE=default
```

### 会话配置

```env
SESSION_DRIVER=file
SESSION_LIFETIME=120
SESSION_COOKIE=laravel_go_session
SESSION_DOMAIN=
SESSION_SECURE_COOKIE=false
SESSION_FILES=storage/framework/sessions
```

### 日志配置

```env
LOG_CHANNEL=single
LOG_LEVEL=debug
```

## 📁 配置文件格式

### 1. JSON 格式

```json
{
  "app": {
    "name": "Laravel-Go",
    "version": "1.0.0",
    "env": "production",
    "debug": false,
    "url": "http://localhost:8080",
    "port": "8080",
    "timezone": "UTC",
    "locale": "en"
  },
  "database": {
    "default": "sqlite",
    "connections": {
      "sqlite": {
        "driver": "sqlite",
        "database": "database/laravel-go.sqlite"
      }
    }
  }
}
```

### 2. YAML 格式

```yaml
app:
  name: Laravel-Go
  version: 1.0.0
  env: production
  debug: false
  url: http://localhost:8080
  port: 8080
  timezone: UTC
  locale: en

database:
  default: sqlite
  connections:
    sqlite:
      driver: sqlite
      database: database/laravel-go.sqlite
```

## 🔍 配置验证

### 1. 基本验证

```go
// 验证规则
rules := map[string]string{
    "app.name":     "required",
    "app.version":  "required",
    "app.port":     "required|numeric",
    "app.debug":    "required|boolean",
    "database.host": "required",
    "database.port": "required|numeric",
}

// 验证配置
if err := cfg.Validate(rules); err != nil {
    log.Fatalf("配置验证失败: %v", err)
}
```

### 2. 支持的验证规则

- `required`: 必填字段
- `numeric`: 数字类型
- `boolean`: 布尔类型
- `string`: 字符串类型
- `email`: 邮箱格式
- `url`: URL 格式

## 🛠️ 高级用法

### 1. 配置热重载

```go
// 监听配置文件变化
go func() {
    for {
        time.Sleep(5 * time.Second)
        if err := cfg.LoadFromFile("config/app.json"); err != nil {
            log.Printf("重新加载配置失败: %v", err)
        }
    }
}()
```

### 2. 配置合并

```go
// 合并多个配置文件
cfg.LoadFromFile("config/app.json")
cfg.LoadFromFile("config/database.json")
cfg.LoadFromFile("config/cache.json")
```

### 3. 配置转换

```go
// 将配置转换为结构体
type AppConfig struct {
    Name string `json:"name"`
    Port int    `json:"port"`
}

var appConfig AppConfig
if err := cfg.LoadFromStruct(&appConfig); err != nil {
    log.Fatalf("加载配置结构体失败: %v", err)
}
```

### 4. 环境特定配置

```go
// 根据环境加载不同配置
env := os.Getenv("APP_ENV")
if env == "" {
    env = "production"
}

cfg.LoadFromFile(fmt.Sprintf("config/%s.json", env))
```

## 📚 最佳实践

### 1. 配置组织

- 按功能模块组织配置文件
- 使用环境变量覆盖敏感配置
- 提供默认配置值
- 实现配置验证

### 2. 安全性

- 不要在配置文件中存储敏感信息
- 使用环境变量存储密钥和密码
- 实现配置加密（如需要）
- 限制配置文件权限

### 3. 性能优化

- 缓存配置值
- 避免频繁读取配置文件
- 使用配置热重载
- 实现配置预加载

### 4. 调试和监控

- 记录配置加载日志
- 监控配置变化
- 实现配置健康检查
- 提供配置诊断工具

## 🔗 相关模块

- [核心模块](../core/) - 应用生命周期管理
- [数据库模块](../database/) - 数据库连接和操作
- [缓存模块](../cache/) - 缓存管理
- [队列模块](../queue/) - 队列处理
- [日志模块](../logging/) - 日志记录

## 📄 许可证

本项目采用 MIT 许可证。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进配置模块。
