# 数据库 API 参考

本文档提供 Laravel-Go Framework 数据库操作组件的 API 参考。

## 📦 Database

数据库管理器，提供数据库连接和操作接口。

### 连接管理

#### Connection(name string) *Connection
获取指定连接。

```go
db := app.DB().Connection("mysql")
```

#### DefaultConnection() *Connection
获取默认连接。

```go
db := app.DB().DefaultConnection()
```

#### Disconnect(name string) error
断开指定连接。

```go
err := app.DB().Disconnect("mysql")
```

#### DisconnectAll() error
断开所有连接。

```go
err := app.DB().DisconnectAll()
```

### 查询构建器

#### Table(name string) *QueryBuilder
开始查询构建。

```go
users := app.DB().Table("users").Get()
```

#### Raw(sql string, args ...interface{}) *QueryBuilder
执行原始 SQL。

```go
result := app.DB().Raw("SELECT * FROM users WHERE active = ?", true).Get()
```

#### Select(columns ...string) *QueryBuilder
选择指定列。

```go
users := app.DB().Table("users").Select("id", "name", "email").Get()
```

#### Where(column, operator string, value interface{}) *QueryBuilder
添加 WHERE 条件。

```go
// 简单条件
users := app.DB().Table("users").Where("active", true).Get()

// 操作符条件
users := app.DB().Table("users").Where("age", ">", 18).Get()

// 多条件
users := app.DB().Table("users").
    Where("active", true).
    Where("age", ">", 18).
    Get()
```

#### OrWhere(column, operator string, value interface{}) *QueryBuilder
添加 OR WHERE 条件。

```go
users := app.DB().Table("users").
    Where("role", "admin").
    OrWhere("role", "moderator").
    Get()
```

#### WhereIn(column string, values interface{}) *QueryBuilder
添加 IN 条件。

```go
users := app.DB().Table("users").WhereIn("id", []int{1, 2, 3}).Get()
```

#### WhereNotIn(column string, values interface{}) *QueryBuilder
添加 NOT IN 条件。

```go
users := app.DB().Table("users").WhereNotIn("id", []int{1, 2, 3}).Get()
```

#### WhereNull(column string) *QueryBuilder
添加 IS NULL 条件。

```go
users := app.DB().Table("users").WhereNull("deleted_at").Get()
```

#### WhereNotNull(column string) *QueryBuilder
添加 IS NOT NULL 条件。

```go
users := app.DB().Table("users").WhereNotNull("email").Get()
```

#### WhereBetween(column string, min, max interface{}) *QueryBuilder
添加 BETWEEN 条件。

```go
users := app.DB().Table("users").WhereBetween("age", 18, 65).Get()
```

#### WhereNotBetween(column string, min, max interface{}) *QueryBuilder
添加 NOT BETWEEN 条件。

```go
users := app.DB().Table("users").WhereNotBetween("age", 18, 65).Get()
```

#### WhereExists(callback func(*QueryBuilder)) *QueryBuilder
添加 EXISTS 条件。

```go
users := app.DB().Table("users").
    WhereExists(func(qb *QueryBuilder) {
        qb.Table("posts").WhereRaw("posts.user_id = users.id")
    }).
    Get()
```

#### WhereNotExists(callback func(*QueryBuilder)) *QueryBuilder
添加 NOT EXISTS 条件。

```go
users := app.DB().Table("users").
    WhereNotExists(func(qb *QueryBuilder) {
        qb.Table("posts").WhereRaw("posts.user_id = users.id")
    }).
    Get()
```

### 排序和分页

#### OrderBy(column, direction string) *QueryBuilder
添加排序。

```go
users := app.DB().Table("users").OrderBy("created_at", "desc").Get()
```

#### OrderByDesc(column string) *QueryBuilder
添加降序排序。

```go
users := app.DB().Table("users").OrderByDesc("created_at").Get()
```

#### OrderByAsc(column string) *QueryBuilder
添加升序排序。

```go
users := app.DB().Table("users").OrderByAsc("name").Get()
```

#### Limit(limit int) *QueryBuilder
限制结果数量。

```go
users := app.DB().Table("users").Limit(10).Get()
```

#### Offset(offset int) *QueryBuilder
设置偏移量。

```go
users := app.DB().Table("users").Offset(20).Limit(10).Get()
```

#### ForPage(page, perPage int) *QueryBuilder
分页查询。

```go
users := app.DB().Table("users").ForPage(1, 10).Get()
```

### 聚合函数

#### Count() int64
统计记录数。

```go
count := app.DB().Table("users").Count()
```

#### Sum(column string) float64
求和。

```go
total := app.DB().Table("orders").Sum("amount")
```

#### Avg(column string) float64
平均值。

```go
average := app.DB().Table("products").Avg("price")
```

#### Max(column string) interface{}
最大值。

```go
maxPrice := app.DB().Table("products").Max("price")
```

#### Min(column string) interface{}
最小值。

```go
minPrice := app.DB().Table("products").Min("price")
```

### 数据操作

#### Get() []map[string]interface{}
获取多条记录。

```go
users := app.DB().Table("users").Get()
```

#### First() map[string]interface{}
获取第一条记录。

```go
user := app.DB().Table("users").Where("id", 1).First()
```

#### Find(id interface{}) map[string]interface{}
根据主键查找。

```go
user := app.DB().Table("users").Find(1)
```

#### Create(data interface{}) map[string]interface{}
创建记录。

```go
user := app.DB().Table("users").Create(map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
})
```

#### Insert(data []map[string]interface{}) bool
批量插入。

```go
users := []map[string]interface{}{
    {"name": "John", "email": "john@example.com"},
    {"name": "Jane", "email": "jane@example.com"},
}
success := app.DB().Table("users").Insert(users)
```

#### Update(data interface{}) int64
更新记录。

```go
affected := app.DB().Table("users").
    Where("id", 1).
    Update(map[string]interface{}{
        "name": "Jane Doe",
    })
```

#### Delete() int64
删除记录。

```go
affected := app.DB().Table("users").Where("id", 1).Delete()
```

#### Truncate() error
清空表。

```go
err := app.DB().Table("users").Truncate()
```

### 连接和联合

#### Join(table, first, operator, second string) *QueryBuilder
内连接。

```go
users := app.DB().Table("users").
    Join("posts", "users.id", "=", "posts.user_id").
    Select("users.*", "posts.title").
    Get()
```

#### LeftJoin(table, first, operator, second string) *QueryBuilder
左连接。

```go
users := app.DB().Table("users").
    LeftJoin("posts", "users.id", "=", "posts.user_id").
    Select("users.*", "posts.title").
    Get()
```

#### RightJoin(table, first, operator, second string) *QueryBuilder
右连接。

```go
users := app.DB().Table("users").
    RightJoin("posts", "users.id", "=", "posts.user_id").
    Select("users.*", "posts.title").
    Get()
```

#### Union(query *QueryBuilder) *QueryBuilder
联合查询。

```go
query1 := app.DB().Table("users").Select("name")
query2 := app.DB().Table("admins").Select("name")
result := query1.Union(query2).Get()
```

### 事务

#### Transaction(callback func(*Transaction) error) error
执行事务。

```go
err := app.DB().Transaction(func(tx *database.Transaction) error {
    // 创建用户
    user := tx.Table("users").Create(userData)
    
    // 创建用户资料
    profile := tx.Table("profiles").Create(map[string]interface{}{
        "user_id": user["id"],
        "bio":     "New user",
    })
    
    return nil
})
```

#### Begin() *Transaction
开始事务。

```go
tx := app.DB().Begin()
defer tx.Rollback()

// 执行操作
user := tx.Table("users").Create(userData)
profile := tx.Table("profiles").Create(profileData)

// 提交事务
tx.Commit()
```

### 事务方法

#### Commit() error
提交事务。

```go
err := tx.Commit()
```

#### Rollback() error
回滚事务。

```go
err := tx.Rollback()
```

#### Table(name string) *QueryBuilder
在事务中查询表。

```go
users := tx.Table("users").Get()
```

### 数据库管理

#### Migrate() error
运行迁移。

```go
err := app.DB().Migrate()
```

#### Rollback() error
回滚迁移。

```go
err := app.DB().Rollback()
```

#### Status() []MigrationStatus
查看迁移状态。

```go
status := app.DB().Status()
```

#### Seed() error
运行数据填充。

```go
err := app.DB().Seed()
```

### 连接配置

#### 配置示例

```go
// 数据库配置
config := map[string]interface{}{
    "default": "mysql",
    "connections": map[string]interface{}{
        "mysql": map[string]interface{}{
            "driver":   "mysql",
            "host":     "localhost",
            "port":     3306,
            "database": "laravel_go",
            "username": "root",
            "password": "password",
            "charset":  "utf8mb4",
            "collation": "utf8mb4_unicode_ci",
            "prefix":   "",
            "strict":   true,
            "engine":   "",
        },
        "postgres": map[string]interface{}{
            "driver":   "postgres",
            "host":     "localhost",
            "port":     5432,
            "database": "laravel_go",
            "username": "postgres",
            "password": "password",
            "charset":  "utf8",
            "prefix":   "",
            "schema":   "public",
            "sslmode":  "disable",
        },
        "sqlite": map[string]interface{}{
            "driver":   "sqlite",
            "database": "database/app.db",
            "prefix":   "",
        },
    },
}
```

### 查询日志

#### EnableQueryLog() *Database
启用查询日志。

```go
app.DB().EnableQueryLog()
```

#### GetQueryLog() []QueryLog
获取查询日志。

```go
logs := app.DB().GetQueryLog()
```

#### FlushQueryLog() *Database
清空查询日志。

```go
app.DB().FlushQueryLog()
```

### 性能优化

#### 索引优化

```go
// 创建索引
app.DB().Raw("CREATE INDEX idx_users_email ON users(email)").Exec()

// 复合索引
app.DB().Raw("CREATE INDEX idx_users_name_email ON users(name, email)").Exec()
```

#### 查询优化

```go
// 使用索引的查询
users := app.DB().Table("users").
    Where("email", "john@example.com").
    Select("id", "name", "email"). // 只选择需要的字段
    Get()

// 避免 N+1 问题
users := app.DB().Table("users").Get()
userIDs := make([]int, 0, len(users))
for _, user := range users {
    userIDs = append(userIDs, user["id"].(int))
}

posts := app.DB().Table("posts").
    WhereIn("user_id", userIDs).
    Get()
```

### 错误处理

#### 连接错误

```go
db := app.DB().Connection("mysql")
if db == nil {
    log.Fatal("Failed to connect to database")
}
```

#### 查询错误

```go
users, err := app.DB().Table("users").Get()
if err != nil {
    log.Printf("Query error: %v", err)
    return
}
```

#### 事务错误

```go
err := app.DB().Transaction(func(tx *database.Transaction) error {
    // 执行操作
    if err := someOperation(); err != nil {
        return err // 自动回滚
    }
    return nil
})

if err != nil {
    log.Printf("Transaction error: %v", err)
}
```

## 📚 下一步

了解更多数据库相关功能：

1. [ORM 使用](guides/orm.md) - 对象关系映射
2. [数据库迁移](guides/migrations.md) - 数据库结构管理
3. [查询优化](best-practices/performance.md) - 性能优化技巧
4. [事务管理](guides/transactions.md) - 事务处理
5. [数据填充](guides/seeding.md) - 测试数据填充

---

这些是 Laravel-Go Framework 的数据库 API。掌握这些 API 将帮助你高效地进行数据库操作！ 🚀 