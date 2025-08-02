# 数据库指南

## 📖 概述

Laravel-Go Framework 提供了完整的数据库支持，包括数据库连接管理、查询构建器、事务处理、迁移和种子数据等功能，支持多种数据库驱动。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [数据库 API 参考](../api/database.md)

## 🚀 快速开始

### 数据库配置

```go
// config/database.go
package config

type Database struct {
    Default string `env:"DB_CONNECTION" default:"mysql"`

    Connections map[string]Connection `env:"DB_CONNECTIONS"`
}

type Connection struct {
    Driver   string `env:"DB_DRIVER" default:"mysql"`
    Host     string `env:"DB_HOST" default:"localhost"`
    Port     int    `env:"DB_PORT" default:"3306"`
    Database string `env:"DB_DATABASE"`
    Username string `env:"DB_USERNAME"`
    Password string `env:"DB_PASSWORD"`
    Charset  string `env:"DB_CHARSET" default:"utf8mb4"`
    Timezone string `env:"DB_TIMEZONE" default:"UTC"`

    // 连接池配置
    MaxOpenConns    int `env:"DB_MAX_OPEN_CONNS" default:"100"`
    MaxIdleConns    int `env:"DB_MAX_IDLE_CONNS" default:"10"`
    ConnMaxLifetime int `env:"DB_CONN_MAX_LIFETIME" default:"3600"`
}
```

### 环境变量配置

```bash
# .env
DB_CONNECTION=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=laravel_go
DB_USERNAME=root
DB_PASSWORD=password
DB_CHARSET=utf8mb4
DB_TIMEZONE=UTC

# 连接池配置
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=3600
```

### 数据库连接

```go
// 获取数据库连接
db := database.Connection("default")

// 执行简单查询
rows, err := db.Query("SELECT * FROM users")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

// 使用查询构建器
users := db.Table("users").Get()
```

## 📋 查询构建器

### 基本查询

```go
// 获取所有用户
users := db.Table("users").Get()

// 获取单个用户
user := db.Table("users").Where("id", 1).First()

// 获取指定字段
users := db.Table("users").Select("id", "name", "email").Get()

// 条件查询
users := db.Table("users").
    Where("age", ">", 18).
    Where("status", "active").
    Get()

// 排序
users := db.Table("users").
    OrderBy("created_at", "desc").
    Get()

// 分页
users := db.Table("users").
    Offset(10).
    Limit(10).
    Get()
```

### 插入数据

```go
// 插入单条记录
id, err := db.Table("users").Insert(map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
    "age":   25,
})

// 插入多条记录
ids, err := db.Table("users").Insert([]map[string]interface{}{
    {"name": "John Doe", "email": "john@example.com"},
    {"name": "Jane Smith", "email": "jane@example.com"},
})

// 使用结构体插入
user := User{
    Name:  "John Doe",
    Email: "john@example.com",
    Age:   25,
}
id, err := db.Table("users").Insert(user)
```

### 更新数据

```go
// 更新记录
affected, err := db.Table("users").
    Where("id", 1).
    Update(map[string]interface{}{
        "name": "John Updated",
        "age":  26,
    })

// 更新所有记录
affected, err := db.Table("users").
    Update(map[string]interface{}{
        "status": "inactive",
    })

// 条件更新
affected, err := db.Table("users").
    Where("age", "<", 18).
    Update(map[string]interface{}{
        "status": "minor",
    })
```

### 删除数据

```go
// 删除指定记录
affected, err := db.Table("users").Where("id", 1).Delete()

// 条件删除
affected, err := db.Table("users").
    Where("status", "inactive").
    Delete()

// 删除所有记录
affected, err := db.Table("users").Delete()
```

### 复杂查询

```go
// 连接查询
users := db.Table("users").
    Join("posts", "users.id", "=", "posts.user_id").
    Select("users.*", "posts.title").
    Get()

// 左连接
users := db.Table("users").
    LeftJoin("profiles", "users.id", "=", "profiles.user_id").
    Get()

// 子查询
users := db.Table("users").
    WhereIn("id", func(query *database.Query) {
        query.Select("user_id").From("posts").Where("status", "published")
    }).
    Get()

// 聚合查询
result := db.Table("users").
    Select("status", db.Raw("COUNT(*) as count")).
    GroupBy("status").
    Get()

// 原始 SQL
users := db.Raw("SELECT * FROM users WHERE age > ?", 18).Get()
```

## 🏗️ 数据库迁移

### 创建迁移

```bash
# 使用 Artisan 命令创建迁移
go run cmd/artisan/main.go make:migration create_users_table
```

### 迁移文件结构

```go
// database/migrations/2024_01_01_000000_create_users_table.go
package migrations

import (
    "laravel-go/framework/database"
    "laravel-go/framework/database/migration"
)

type CreateUsersTable struct {
    migration.Migration
}

func (m *CreateUsersTable) Up() error {
    return m.Schema.CreateTable("users", func(table *database.Blueprint) {
        table.Id("id")
        table.String("name", 255).NotNull()
        table.String("email", 255).Unique().NotNull()
        table.String("password", 255).NotNull()
        table.Integer("age").Nullable()
        table.Boolean("is_active").Default(true)
        table.Timestamps()

        // 索引
        table.Index("email")
        table.Index("name", "email")
    })
}

func (m *CreateUsersTable) Down() error {
    return m.Schema.DropTable("users")
}
```

### 字段类型

```go
func (m *CreateUsersTable) Up() error {
    return m.Schema.CreateTable("users", func(table *database.Blueprint) {
        // 主键
        table.Id("id")
        table.Uuid("uuid")

        // 字符串类型
        table.String("name", 255)
        table.Text("description")
        table.LongText("content")
        table.Char("code", 10)

        // 数字类型
        table.Integer("age")
        table.BigInteger("big_id")
        table.SmallInteger("small_id")
        table.TinyInteger("tiny_id")
        table.Decimal("price", 10, 2)
        table.Float("score", 8, 2)
        table.Double("amount", 15, 2)

        // 布尔类型
        table.Boolean("is_active")

        // 日期时间类型
        table.Date("birth_date")
        table.DateTime("created_at")
        table.Time("start_time")
        table.Timestamp("updated_at")

        // 二进制类型
        table.Binary("file_data")
        table.Blob("large_data")

        // JSON 类型
        table.Json("metadata")

        // 枚举类型
        table.Enum("status", []string{"active", "inactive", "pending"})

        // 几何类型
        table.Geometry("location")
        table.Point("coordinates")

        // 时间戳
        table.Timestamps()
        table.TimestampsTz()
        table.SoftDeletes()
    })
}
```

### 字段修饰符

```go
func (m *CreateUsersTable) Up() error {
    return m.Schema.CreateTable("users", func(table *database.Blueprint) {
        table.Id("id")

        // 基本修饰符
        table.String("name", 255).NotNull()
        table.String("email", 255).Unique()
        table.String("code", 10).Default("ABC123")
        table.Integer("age").Nullable()

        // 索引
        table.String("username", 50).Index()
        table.String("email", 255).Unique()
        table.Index("name", "email")
        table.Unique("username", "email")

        // 外键
        table.Integer("role_id").Unsigned()
        table.ForeignKey("role_id").References("id").On("roles").OnDelete("cascade")

        // 其他修饰符
        table.String("comment", 1000).Comment("用户备注")
        table.String("status", 20).Collation("utf8mb4_unicode_ci")
    })
}
```

### 表操作

```go
// 创建表
func (m *CreateUsersTable) Up() error {
    return m.Schema.CreateTable("users", func(table *database.Blueprint) {
        table.Id("id")
        table.String("name", 255)
        table.Timestamps()
    })
}

// 删除表
func (m *CreateUsersTable) Down() error {
    return m.Schema.DropTable("users")
}

// 重命名表
func (m *RenameUsersTable) Up() error {
    return m.Schema.RenameTable("users", "accounts")
}

// 检查表是否存在
func (m *CreateUsersTable) Up() error {
    if m.Schema.HasTable("users") {
        return nil
    }

    return m.Schema.CreateTable("users", func(table *database.Blueprint) {
        // 表结构
    })
}
```

### 列操作

```go
// 添加列
func (m *AddColumnToUsersTable) Up() error {
    return m.Schema.Table("users", func(table *database.Blueprint) {
        table.String("phone", 20).Nullable().After("email")
        table.Boolean("is_verified").Default(false).After("phone")
    })
}

// 修改列
func (m *ModifyColumnInUsersTable) Up() error {
    return m.Schema.Table("users", func(table *database.Blueprint) {
        table.String("name", 100).Change() // 修改长度
        table.String("email").Unique().Change() // 添加唯一约束
    })
}

// 删除列
func (m *RemoveColumnFromUsersTable) Up() error {
    return m.Schema.Table("users", func(table *database.Blueprint) {
        table.DropColumn("phone")
        table.DropColumn("is_verified")
    })
}

// 重命名列
func (m *RenameColumnInUsersTable) Up() error {
    return m.Schema.Table("users", func(table *database.Blueprint) {
        table.RenameColumn("name", "full_name")
    })
}
```

## 🌱 数据填充

### 创建填充器

```bash
# 创建填充器
go run cmd/artisan/main.go make:seeder UserSeeder
```

### 填充器实现

```go
// database/seeders/UserSeeder.go
package seeders

import (
    "laravel-go/framework/database"
    "laravel-go/framework/database/seeder"
    "laravel-go/app/Models"
)

type UserSeeder struct {
    seeder.Seeder
}

func (s *UserSeeder) Run() error {
    // 清空表
    database.Table("users").Delete()

    // 插入测试数据
    users := []map[string]interface{}{
        {
            "name":     "John Doe",
            "email":    "john@example.com",
            "password": "password123",
            "age":      25,
        },
        {
            "name":     "Jane Smith",
            "email":    "jane@example.com",
            "password": "password123",
            "age":      30,
        },
        {
            "name":     "Bob Johnson",
            "email":    "bob@example.com",
            "password": "password123",
            "age":      35,
        },
    }

    for _, user := range users {
        database.Table("users").Insert(user)
    }

    return nil
}
```

### 运行填充器

```go
// 运行所有填充器
seeder.Run()

// 运行指定填充器
seeder.Run(&UserSeeder{})

// 运行指定填充器并清空数据
seeder.RunFresh(&UserSeeder{})
```

## 🔄 事务处理

### 基本事务

```go
// 开始事务
tx := db.Begin()

// 执行操作
_, err := tx.Table("users").Insert(user)
if err != nil {
    tx.Rollback()
    return err
}

_, err = tx.Table("profiles").Insert(profile)
if err != nil {
    tx.Rollback()
    return err
}

// 提交事务
return tx.Commit()
```

### 事务闭包

```go
// 使用事务闭包
err := db.Transaction(func(tx *database.Connection) error {
    // 插入用户
    _, err := tx.Table("users").Insert(user)
    if err != nil {
        return err
    }

    // 插入用户资料
    _, err = tx.Table("profiles").Insert(profile)
    if err != nil {
        return err
    }

    return nil
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

### 嵌套事务

```go
// 嵌套事务
err := db.Transaction(func(tx *database.Connection) error {
    // 外层事务
    _, err := tx.Table("users").Insert(user)
    if err != nil {
        return err
    }

    // 内层事务
    return tx.Transaction(func(tx2 *database.Connection) error {
        _, err := tx2.Table("profiles").Insert(profile)
        if err != nil {
            return err
        }

        _, err = tx2.Table("settings").Insert(setting)
        return err
    })
})
```

## 🔍 查询优化

### 索引优化

```go
// 创建索引
func (m *AddIndexesToUsersTable) Up() error {
    return m.Schema.Table("users", func(table *database.Blueprint) {
        // 单列索引
        table.Index("email")

        // 复合索引
        table.Index("name", "email")

        // 唯一索引
        table.Unique("username")

        // 前缀索引
        table.Index("name", "email").Prefix(10)
    })
}
```

### 查询优化

```go
// 使用预加载避免 N+1 问题
users := db.Table("users").
    With("posts").
    With("comments").
    Get()

// 使用分页
users := db.Table("users").
    Offset(0).
    Limit(20).
    Get()

// 使用缓存
cacheKey := "users:list"
if cached, found := cache.Get(cacheKey); found {
    return cached.([]User)
}

users := db.Table("users").Get()
cache.Set(cacheKey, users, time.Hour)

// 使用原始 SQL 优化复杂查询
users := db.Raw(`
    SELECT u.*, COUNT(p.id) as post_count
    FROM users u
    LEFT JOIN posts p ON u.id = p.user_id
    GROUP BY u.id
    HAVING post_count > 0
`).Get()
```

### 连接池优化

```go
// 配置连接池
db := database.Connection("default")
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)
```

## 🛡️ 安全性

### SQL 注入防护

```go
// 使用参数化查询
users := db.Table("users").
    Where("name", "LIKE", "%" + searchTerm + "%").
    Get()

// 使用原始 SQL 时使用参数
users := db.Raw("SELECT * FROM users WHERE name LIKE ?", "%"+searchTerm+"%").Get()

// 避免直接拼接 SQL
// ❌ 危险的做法
users := db.Raw("SELECT * FROM users WHERE name LIKE '%" + searchTerm + "%'").Get()
```

### 数据验证

```go
// 验证输入数据
func validateUserData(data map[string]interface{}) error {
    if name, ok := data["name"].(string); !ok || len(name) == 0 {
        return errors.New("name is required")
    }

    if email, ok := data["email"].(string); !ok || !isValidEmail(email) {
        return errors.New("valid email is required")
    }

    return nil
}

// 使用验证器
func (c *UserController) Store(request http.Request) http.Response {
    validator := validation.NewValidator()

    rules := map[string]string{
        "name":  "required|string|max:255",
        "email": "required|email|unique:users,email",
        "age":   "integer|min:0|max:150",
    }

    if err := validator.Validate(request.Body, rules); err != nil {
        return c.JsonError(err.Error(), 422)
    }

    // 创建用户
    user := c.userService.CreateUser(request.Body)
    return c.Json(user)
}
```

## 📊 监控和调试

### 查询日志

```go
// 启用查询日志
db.EnableQueryLog()

// 执行查询
users := db.Table("users").Get()

// 获取查询日志
queries := db.GetQueryLog()
for _, query := range queries {
    log.Printf("SQL: %s, Time: %v", query.SQL, query.Time)
}
```

### 性能监控

```go
// 监控查询性能
start := time.Now()
users := db.Table("users").Get()
duration := time.Since(start)

if duration > time.Second {
    log.Printf("Slow query detected: %v", duration)
}
```

### 数据库健康检查

```go
// 检查数据库连接
func checkDatabaseHealth() error {
    db := database.Connection("default")

    // 执行简单查询测试连接
    _, err := db.Raw("SELECT 1").First()
    if err != nil {
        return fmt.Errorf("database connection failed: %v", err)
    }

    return nil
}
```

## 📝 总结

Laravel-Go Framework 的数据库系统提供了：

1. **灵活性**: 支持多种数据库驱动和查询方式
2. **安全性**: 内置 SQL 注入防护和数据验证
3. **可维护性**: 支持迁移和种子管理
4. **性能优化**: 提供查询优化和连接池管理
5. **可观测性**: 支持查询日志和性能监控

通过合理使用数据库系统的各种功能，可以构建出高效、安全、可维护的数据驱动应用程序。
