# ORM API 参考

本文档提供 Laravel-Go Framework ORM（对象关系映射）组件的 API 参考。

## 📦 Model

ORM 模型基类，提供数据库操作的抽象接口。

### 模型定义

#### 基础模型

```go
type User struct {
    orm.Model
    Name     string `json:"name" gorm:"column:name;type:varchar(255);not null"`
    Email    string `json:"email" gorm:"column:email;type:varchar(255);unique;not null"`
    Password string `json:"-" gorm:"column:password;type:varchar(255);not null"`
    Age      int    `json:"age" gorm:"column:age;type:int"`
    Active   bool   `json:"active" gorm:"column:active;type:boolean;default:true"`
}
```

#### 表名配置

```go
// 自定义表名
func (User) TableName() string {
    return "users"
}

// 或使用标签
type User struct {
    orm.Model
    Name string `gorm:"table:users"`
}
```

#### 主键配置

```go
type User struct {
    orm.Model
    ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
    Name string `json:"name"`
}
```

### 查询操作

#### Find(id interface{}) \*Model

根据主键查找。

```go
user := User{}.Find(1)
if user != nil {
    fmt.Printf("User: %s\n", user.Name)
}
```

#### First(conditions ...interface{}) \*Model

获取第一条记录。

```go
// 简单条件
user := User{}.Where("email", "john@example.com").First()

// 多条件
user := User{}.Where("active", true).Where("age", ">", 18).First()

// 使用结构体
user := User{}.Where(User{Email: "john@example.com", Active: true}).First()
```

#### Get() []\*Model

获取多条记录。

```go
users := User{}.Get()
for _, user := range users {
    fmt.Printf("User: %s\n", user.Name)
}
```

#### Where(column, operator string, value interface{}) \*QueryBuilder

添加 WHERE 条件。

```go
// 简单条件
users := User{}.Where("active", true).Get()

// 操作符条件
users := User{}.Where("age", ">", 18).Get()

// 多条件
users := User{}.Where("active", true).Where("age", ">", 18).Get()
```

#### OrWhere(column, operator string, value interface{}) \*QueryBuilder

添加 OR WHERE 条件。

```go
users := User{}.Where("role", "admin").OrWhere("role", "moderator").Get()
```

#### WhereIn(column string, values interface{}) \*QueryBuilder

添加 IN 条件。

```go
users := User{}.WhereIn("id", []int{1, 2, 3}).Get()
```

#### WhereNotIn(column string, values interface{}) \*QueryBuilder

添加 NOT IN 条件。

```go
users := User{}.WhereNotIn("id", []int{1, 2, 3}).Get()
```

#### WhereNull(column string) \*QueryBuilder

添加 IS NULL 条件。

```go
users := User{}.WhereNull("deleted_at").Get()
```

#### WhereNotNull(column string) \*QueryBuilder

添加 IS NOT NULL 条件。

```go
users := User{}.WhereNotNull("email").Get()
```

#### WhereBetween(column string, min, max interface{}) \*QueryBuilder

添加 BETWEEN 条件。

```go
users := User{}.WhereBetween("age", 18, 65).Get()
```

#### WhereNotBetween(column string, min, max interface{}) \*QueryBuilder

添加 NOT BETWEEN 条件。

```go
users := User{}.WhereNotBetween("age", 18, 65).Get()
```

### 排序和分页

#### OrderBy(column, direction string) \*QueryBuilder

添加排序。

```go
users := User{}.OrderBy("created_at", "desc").Get()
```

#### OrderByDesc(column string) \*QueryBuilder

添加降序排序。

```go
users := User{}.OrderByDesc("created_at").Get()
```

#### OrderByAsc(column string) \*QueryBuilder

添加升序排序。

```go
users := User{}.OrderByAsc("name").Get()
```

#### Limit(limit int) \*QueryBuilder

限制结果数量。

```go
users := User{}.Limit(10).Get()
```

#### Offset(offset int) \*QueryBuilder

设置偏移量。

```go
users := User{}.Offset(20).Limit(10).Get()
```

#### ForPage(page, perPage int) \*QueryBuilder

分页查询。

```go
users := User{}.ForPage(1, 10).Get()
```

### 聚合函数

#### Count() int64

统计记录数。

```go
count := User{}.Count()
```

#### Sum(column string) float64

求和。

```go
total := Order{}.Sum("amount")
```

#### Avg(column string) float64

平均值。

```go
average := Product{}.Avg("price")
```

#### Max(column string) interface{}

最大值。

```go
maxPrice := Product{}.Max("price")
```

#### Min(column string) interface{}

最小值。

```go
minPrice := Product{}.Min("price")
```

### 数据操作

#### Save() error

保存模型（创建或更新）。

```go
user := &User{
    Name:  "John Doe",
    Email: "john@example.com",
}

err := user.Save()
if err != nil {
    log.Printf("Save error: %v", err)
}
```

#### Create(data interface{}) \*Model

创建新记录。

```go
user := User{}.Create(&User{
    Name:  "John Doe",
    Email: "john@example.com",
})
```

#### Update(data interface{}) error

更新记录。

```go
user := User{}.Find(1)
if user != nil {
    user.Name = "Jane Doe"
    err := user.Update()
    if err != nil {
        log.Printf("Update error: %v", err)
    }
}
```

#### Delete() error

删除记录。

```go
user := User{}.Find(1)
if user != nil {
    err := user.Delete()
    if err != nil {
        log.Printf("Delete error: %v", err)
    }
}
```

#### SoftDelete() error

软删除记录。

```go
user := User{}.Find(1)
if user != nil {
    err := user.SoftDelete()
    if err != nil {
        log.Printf("Soft delete error: %v", err)
    }
}
```

#### Restore() error

恢复软删除的记录。

```go
user := User{}.WithTrashed().Find(1)
if user != nil {
    err := user.Restore()
    if err != nil {
        log.Printf("Restore error: %v", err)
    }
}
```

### 关联关系

#### HasOne(related interface{}, foreignKey, localKey string) \*HasOne

一对一关联。

```go
type User struct {
    orm.Model
    Name string
}

type Profile struct {
    orm.Model
    UserID uint
    Bio    string
}

// 在 User 模型中定义关联
func (u *User) Profile() *HasOne {
    return u.HasOne(&Profile{}, "user_id", "id")
}

// 使用关联
user := User{}.Find(1)
profile := user.Profile().First()
```

#### HasMany(related interface{}, foreignKey, localKey string) \*HasMany

一对多关联。

```go
type User struct {
    orm.Model
    Name string
}

type Post struct {
    orm.Model
    UserID uint
    Title  string
}

// 在 User 模型中定义关联
func (u *User) Posts() *HasMany {
    return u.HasMany(&Post{}, "user_id", "id")
}

// 使用关联
user := User{}.Find(1)
posts := user.Posts().Get()
```

#### BelongsTo(related interface{}, foreignKey, ownerKey string) \*BelongsTo

多对一关联。

```go
type Post struct {
    orm.Model
    UserID uint
    Title  string
}

type User struct {
    orm.Model
    Name string
}

// 在 Post 模型中定义关联
func (p *Post) User() *BelongsTo {
    return p.BelongsTo(&User{}, "user_id", "id")
}

// 使用关联
post := Post{}.Find(1)
user := post.User().First()
```

#### BelongsToMany(related interface{}, pivotTable, foreignKey, relatedKey string) \*BelongsToMany

多对多关联。

```go
type User struct {
    orm.Model
    Name string
}

type Role struct {
    orm.Model
    Name string
}

// 在 User 模型中定义关联
func (u *User) Roles() *BelongsToMany {
    return u.BelongsToMany(&Role{}, "user_roles", "user_id", "role_id")
}

// 使用关联
user := User{}.Find(1)
roles := user.Roles().Get()
```

### 关联查询

#### With(relations ...string) \*QueryBuilder

预加载关联。

```go
// 预加载单个关联
users := User{}.With("profile").Get()

// 预加载多个关联
users := User{}.With("profile", "posts").Get()

// 预加载嵌套关联
users := User{}.With("posts.comments").Get()
```

#### Load(relations ...string) error

延迟加载关联。

```go
user := User{}.Find(1)
err := user.Load("profile", "posts")
if err != nil {
    log.Printf("Load error: %v", err)
}
```

### 模型钩子

#### BeforeCreate() error

创建前钩子。

```go
func (u *User) BeforeCreate() error {
    u.CreatedAt = time.Now()
    return nil
}
```

#### AfterCreate() error

创建后钩子。

```go
func (u *User) AfterCreate() error {
    // 发送欢迎邮件
    return sendWelcomeEmail(u.Email)
}
```

#### BeforeUpdate() error

更新前钩子。

```go
func (u *User) BeforeUpdate() error {
    u.UpdatedAt = time.Now()
    return nil
}
```

#### AfterUpdate() error

更新后钩子。

```go
func (u *User) AfterUpdate() error {
    // 记录更新日志
    return logUserUpdate(u.ID)
}
```

#### BeforeDelete() error

删除前钩子。

```go
func (u *User) BeforeDelete() error {
    // 检查是否可以删除
    if u.HasActivePosts() {
        return errors.New("cannot delete user with active posts")
    }
    return nil
}
```

#### AfterDelete() error

删除后钩子。

```go
func (u *User) AfterDelete() error {
    // 清理相关数据
    return cleanupUserData(u.ID)
}
```

### 字段操作

#### Fill(data map[string]interface{}) \*Model

填充字段。

```go
user := &User{}
user.Fill(map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
})
```

#### Set(key string, value interface{}) \*Model

设置字段值。

```go
user := &User{}
user.Set("name", "John Doe")
user.Set("email", "john@example.com")
```

#### Get(key string) interface{}

获取字段值。

```go
user := User{}.Find(1)
name := user.Get("name")
email := user.Get("email")
```

#### IsDirty(key string) bool

检查字段是否已修改。

```go
user := User{}.Find(1)
user.Name = "Jane Doe"
isDirty := user.IsDirty("name") // true
```

#### GetDirty() map[string]interface{}

获取已修改的字段。

```go
user := User{}.Find(1)
user.Name = "Jane Doe"
user.Email = "jane@example.com"
dirty := user.GetDirty() // map[string]interface{}{"name": "Jane Doe", "email": "jane@example.com"}
```

#### GetOriginal(key string) interface{}

获取原始字段值。

```go
user := User{}.Find(1)
originalName := user.GetOriginal("name")
user.Name = "Jane Doe"
currentName := user.Get("name")
originalName := user.GetOriginal("name") // 原始值
```

### 批量操作

#### CreateMany(models []interface{}) error

批量创建。

```go
users := []interface{}{
    &User{Name: "John", Email: "john@example.com"},
    &User{Name: "Jane", Email: "jane@example.com"},
}

err := User{}.CreateMany(users)
```

#### UpdateMany(conditions interface{}, data interface{}) error

批量更新。

```go
err := User{}.Where("active", false).UpdateMany(map[string]interface{}{
    "status": "inactive",
})
```

#### DeleteMany(conditions interface{}) error

批量删除。

```go
err := User{}.Where("active", false).DeleteMany()
```

### 查询作用域

#### Scope(name string, callback func(*QueryBuilder) *QueryBuilder)

定义查询作用域。

```go
// 定义作用域
func (User) ScopeActive(query *QueryBuilder) *QueryBuilder {
    return query.Where("active", true)
}

func (User) ScopeOlderThan(age int) func(*QueryBuilder) *QueryBuilder {
    return func(query *QueryBuilder) *QueryBuilder {
        return query.Where("age", ">", age)
    }
}

// 使用作用域
users := User{}.Scope("active").Get()
users := User{}.Scope("olderThan", 18).Get()
```

### 软删除

#### WithTrashed() \*QueryBuilder

包含软删除的记录。

```go
users := User{}.WithTrashed().Get()
```

#### OnlyTrashed() \*QueryBuilder

只查询软删除的记录。

```go
users := User{}.OnlyTrashed().Get()
```

### 时间戳

#### Timestamps() bool

是否使用时间戳。

```go
type User struct {
    orm.Model
    Name string
}

func (User) Timestamps() bool {
    return true // 默认使用时间戳
}
```

#### CreatedAt() time.Time

获取创建时间。

```go
user := User{}.Find(1)
createdAt := user.CreatedAt()
```

#### UpdatedAt() time.Time

获取更新时间。

```go
user := User{}.Find(1)
updatedAt := user.UpdatedAt()
```

### 序列化

#### ToJSON() ([]byte, error)

转换为 JSON。

```go
user := User{}.Find(1)
jsonData, err := user.ToJSON()
if err != nil {
    log.Printf("JSON error: %v", err)
}
```

#### ToMap() map[string]interface{}

转换为 Map。

```go
user := User{}.Find(1)
userMap := user.ToMap()
```

### 验证

#### Validate(rules map[string]string) error

验证模型数据。

```go
user := &User{
    Name:  "John",
    Email: "invalid-email",
}

rules := map[string]string{
    "name":  "required|string|max:255",
    "email": "required|email",
}

err := user.Validate(rules)
if err != nil {
    log.Printf("Validation error: %v", err)
}
```

## 📚 下一步

了解更多 ORM 相关功能：

1. [模型关联](guides/orm-relationships.md) - 关联关系详解
2. [查询优化](best-practices/performance.md) - ORM 性能优化
3. [模型钩子](guides/model-hooks.md) - 模型生命周期
4. [批量操作](guides/batch-operations.md) - 批量数据处理
5. [软删除](guides/soft-deletes.md) - 软删除功能

---

这些是 Laravel-Go Framework 的 ORM API。掌握这些 API 将帮助你优雅地进行数据库操作！ 🚀
