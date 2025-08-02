# ORM 指南

## 📖 概述

Laravel-Go Framework 提供了强大的 ORM（对象关系映射）系统，支持模型定义、关系管理、查询构建器、迁移和种子数据等功能，简化数据库操作。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [ORM API 参考](../api/orm.md)

## 🚀 快速开始

### 基本模型定义

```go
// app/Models/User.go
package models

import (
    "laravel-go/framework/database"
    "time"
)

type User struct {
    database.Model
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"size:255;not null"`
    Email     string    `json:"email" gorm:"size:255;unique;not null"`
    Password  string    `json:"-" gorm:"size:255;not null"`
    Age       int       `json:"age" gorm:"default:0"`
    IsActive  bool      `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}

// 模型钩子
func (u *User) BeforeCreate() error {
    // 创建前的处理逻辑
    if u.Age == 0 {
        u.Age = 18
    }
    return nil
}

func (u *User) AfterCreate() error {
    // 创建后的处理逻辑
    return nil
}
```

### 基本操作

```go
// 创建用户
user := &User{
    Name:     "John Doe",
    Email:    "john@example.com",
    Password: "password123",
    Age:      25,
}
db.Create(user)

// 查询用户
var user User
db.First(&user, 1) // 根据主键查询
db.First(&user, "email = ?", "john@example.com") // 根据条件查询

// 更新用户
db.Model(&user).Update("Name", "John Updated")
db.Model(&user).Updates(map[string]interface{}{
    "name": "John Updated",
    "age":  26,
})

// 删除用户
db.Delete(&user, 1)
```

## 📋 模型定义

### 字段标签

```go
type User struct {
    database.Model
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Name      string    `json:"name" gorm:"size:255;not null;index"`
    Email     string    `json:"email" gorm:"size:255;unique;not null;index"`
    Password  string    `json:"-" gorm:"size:255;not null"` // json:"-" 隐藏字段
    Age       int       `json:"age" gorm:"default:18;check:age >= 0"`
    Score     float64   `json:"score" gorm:"type:decimal(5,2);default:0.00"`
    Status    string    `json:"status" gorm:"type:enum('active','inactive','pending');default:'active'"`
    Metadata  string    `json:"metadata" gorm:"type:json"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt *time.Time `json:"deleted_at" gorm:"index"` // 软删除
}
```

### 字段类型

```go
type Product struct {
    database.Model
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"size:255;not null"`
    Description string    `gorm:"type:text"`
    Price       float64   `gorm:"type:decimal(10,2)"`
    Stock       int       `gorm:"default:0"`
    IsActive    bool      `gorm:"default:true"`
    Category    string    `gorm:"size:100;index"`
    Tags        string    `gorm:"type:json"`
    ImageURL    string    `gorm:"size:500"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
```

### 模型方法

```go
type User struct {
    database.Model
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"`
}

// 自定义方法
func (u *User) FullName() string {
    return u.Name
}

func (u *User) IsAdult() bool {
    return u.Age >= 18
}

// 密码加密
func (u *User) SetPassword(password string) {
    u.Password = hashPassword(password)
}

// 验证密码
func (u *User) CheckPassword(password string) bool {
    return checkPassword(password, u.Password)
}
```

## 🔗 关联关系

### 一对一关系

```go
type User struct {
    database.Model
    ID      uint   `json:"id"`
    Name    string `json:"name"`
    Email   string `json:"email"`
    Profile Profile `json:"profile" gorm:"foreignKey:UserID"`
}

type Profile struct {
    database.Model
    ID       uint   `json:"id"`
    UserID   uint   `json:"user_id"`
    Avatar   string `json:"avatar"`
    Bio      string `json:"bio"`
    Location string `json:"location"`
    User     User   `json:"user" gorm:"foreignKey:UserID"`
}

// 使用关联
var user User
db.Preload("Profile").First(&user, 1)

// 创建关联
user := &User{Name: "John", Email: "john@example.com"}
profile := &Profile{Avatar: "avatar.jpg", Bio: "Hello World"}
user.Profile = *profile
db.Create(&user)
```

### 一对多关系

```go
type User struct {
    database.Model
    ID    uint    `json:"id"`
    Name  string  `json:"name"`
    Email string  `json:"email"`
    Posts []Post  `json:"posts" gorm:"foreignKey:UserID"`
}

type Post struct {
    database.Model
    ID      uint   `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
    UserID  uint   `json:"user_id"`
    User    User   `json:"user" gorm:"foreignKey:UserID"`
}

// 使用关联
var user User
db.Preload("Posts").First(&user, 1)

// 创建关联
user := &User{Name: "John", Email: "john@example.com"}
posts := []Post{
    {Title: "First Post", Content: "Hello World"},
    {Title: "Second Post", Content: "Another post"},
}
user.Posts = posts
db.Create(&user)
```

### 多对多关系

```go
type User struct {
    database.Model
    ID       uint     `json:"id"`
    Name     string   `json:"name"`
    Email    string   `json:"email"`
    Roles    []Role   `json:"roles" gorm:"many2many:user_roles;"`
}

type Role struct {
    database.Model
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Users       []User `json:"users" gorm:"many2many:user_roles;"`
}

// 使用关联
var user User
db.Preload("Roles").First(&user, 1)

// 创建关联
user := &User{Name: "John", Email: "john@example.com"}
roles := []Role{
    {Name: "admin", Description: "Administrator"},
    {Name: "user", Description: "Regular user"},
}
user.Roles = roles
db.Create(&user)

// 添加角色
var user User
var role Role
db.First(&user, 1)
db.First(&role, 1)
db.Model(&user).Association("Roles").Append(&role)

// 移除角色
db.Model(&user).Association("Roles").Delete(&role)
```

### 关联查询

```go
// 预加载关联
var users []User
db.Preload("Profile").Preload("Posts").Find(&users)

// 条件预加载
db.Preload("Posts", "status = ?", "published").Find(&users)

// 嵌套预加载
db.Preload("Posts.Comments").Find(&users)

// 关联计数
var user User
db.First(&user, 1)
postCount := db.Model(&user).Association("Posts").Count()

// 关联查询
var posts []Post
db.Model(&user).Association("Posts").Find(&posts)
```

## 🔍 查询操作

### 基本查询

```go
// 查找所有用户
var users []User
db.Find(&users)

// 查找单个用户
var user User
db.First(&user, 1) // 根据主键查找
db.First(&user, "email = ?", "john@example.com") // 根据条件查找

// 查找最后一个用户
db.Last(&user)

// 查找指定数量的用户
db.Limit(10).Find(&users)

// 跳过指定数量的用户
db.Offset(10).Find(&users)
```

### 条件查询

```go
// Where 条件
db.Where("age > ?", 18).Find(&users)
db.Where("name IN ?", []string{"John", "Jane"}).Find(&users)
db.Where("email LIKE ?", "%@example.com").Find(&users)

// 链式条件
db.Where("age > ?", 18).Where("status = ?", "active").Find(&users)

// Or 条件
db.Where("age > ?", 18).Or("status = ?", "admin").Find(&users)

// Not 条件
db.Not("status = ?", "inactive").Find(&users)

// 结构体条件
db.Where(&User{Age: 25, Status: "active"}).Find(&users)

// Map 条件
db.Where(map[string]interface{}{"age": 25, "status": "active"}).Find(&users)
```

### 高级查询

```go
// 选择字段
db.Select("id", "name", "email").Find(&users)

// 排除字段
db.Omit("password", "created_at").Find(&users)

// 排序
db.Order("created_at desc").Find(&users)
db.Order("age asc, name desc").Find(&users)

// 分组
type Result struct {
    Status string
    Count  int
}
var results []Result
db.Model(&User{}).Select("status, count(*) as count").Group("status").Find(&results)

// 聚合
var count int64
db.Model(&User{}).Count(&count)

var avgAge float64
db.Model(&User{}).Select("avg(age)").Scan(&avgAge)

// 子查询
var users []User
db.Where("age > (?)", db.Model(&User{}).Select("avg(age)")).Find(&users)
```

### 分页查询

```go
// 基本分页
var users []User
page := 1
pageSize := 10
offset := (page - 1) * pageSize

db.Offset(offset).Limit(pageSize).Find(&users)

// 分页结构
type Pagination struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalPages int         `json:"total_pages"`
}

func Paginate(page, pageSize int, model interface{}) Pagination {
    var total int64
    db.Model(model).Count(&total)

    offset := (page - 1) * pageSize
    db.Offset(offset).Limit(pageSize).Find(model)

    totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

    return Pagination{
        Data:       model,
        Total:      total,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
    }
}
```

## ✏️ 创建和更新

### 创建记录

```go
// 创建单个记录
user := &User{
    Name:     "John Doe",
    Email:    "john@example.com",
    Password: "password123",
    Age:      25,
}
db.Create(user)

// 批量创建
users := []User{
    {Name: "John", Email: "john@example.com"},
    {Name: "Jane", Email: "jane@example.com"},
    {Name: "Bob", Email: "bob@example.com"},
}
db.Create(&users)

// 创建时忽略某些字段
db.Omit("CreatedAt", "UpdatedAt").Create(&user)

// 创建时选择某些字段
db.Select("Name", "Email").Create(&user)
```

### 更新记录

```go
// 更新单个字段
db.Model(&user).Update("Name", "John Updated")

// 更新多个字段
db.Model(&user).Updates(map[string]interface{}{
    "name": "John Updated",
    "age":  26,
})

// 使用结构体更新
updates := User{Name: "John Updated", Age: 26}
db.Model(&user).Updates(updates)

// 条件更新
db.Model(&User{}).Where("age < ?", 18).Update("status", "minor")

// 更新所有字段
db.Save(&user)

// 更新时忽略某些字段
db.Model(&user).Omit("CreatedAt").Updates(updates)
```

### 删除记录

```go
// 删除单个记录
db.Delete(&user, 1)

// 条件删除
db.Where("age < ?", 18).Delete(&User{})

// 批量删除
db.Delete(&User{}, []int{1, 2, 3})

// 软删除（如果模型支持）
db.Delete(&user) // 设置 DeletedAt 字段

// 强制删除（忽略软删除）
db.Unscoped().Delete(&user, 1)
```

## 🔄 事务处理

### 基本事务

```go
// 开始事务
tx := db.Begin()

// 执行操作
user := &User{Name: "John", Email: "john@example.com"}
if err := tx.Create(user).Error; err != nil {
    tx.Rollback()
    return err
}

profile := &Profile{UserID: user.ID, Bio: "Hello World"}
if err := tx.Create(profile).Error; err != nil {
    tx.Rollback()
    return err
}

// 提交事务
return tx.Commit().Error
```

### 事务闭包

```go
// 使用事务闭包
err := db.Transaction(func(tx *gorm.DB) error {
    // 创建用户
    user := &User{Name: "John", Email: "john@example.com"}
    if err := tx.Create(user).Error; err != nil {
        return err
    }

    // 创建用户资料
    profile := &Profile{UserID: user.ID, Bio: "Hello World"}
    if err := tx.Create(profile).Error; err != nil {
        return err
    }

    return nil
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

## 🎯 模型钩子

### 生命周期钩子

```go
type User struct {
    database.Model
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"`
}

// 创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 加密密码
    u.Password = hashPassword(u.Password)
    return nil
}

// 创建后
func (u *User) AfterCreate(tx *gorm.DB) error {
    // 发送欢迎邮件
    go sendWelcomeEmail(u.Email)
    return nil
}

// 更新前
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    // 更新前的处理逻辑
    return nil
}

// 更新后
func (u *User) AfterUpdate(tx *gorm.DB) error {
    // 更新后的处理逻辑
    return nil
}

// 删除前
func (u *User) BeforeDelete(tx *gorm.DB) error {
    // 删除前的处理逻辑
    return nil
}

// 删除后
func (u *User) AfterDelete(tx *gorm.DB) error {
    // 删除后的处理逻辑
    return nil
}

// 查询后
func (u *User) AfterFind(tx *gorm.DB) error {
    // 查询后的处理逻辑
    return nil
}
```

## 🛠️ 高级功能

### 作用域

```go
type User struct {
    database.Model
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Status   string `json:"status"`
    Age      int    `json:"age"`
}

// 定义作用域
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "active")
}

func AdultUsers(db *gorm.DB) *gorm.DB {
    return db.Where("age >= ?", 18)
}

func RecentUsers(db *gorm.DB) *gorm.DB {
    return db.Where("created_at >= ?", time.Now().AddDate(0, 0, -7))
}

// 使用作用域
var users []User
db.Scopes(ActiveUsers, AdultUsers).Find(&users)
```

### 模型方法

```go
type User struct {
    database.Model
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
}

// 自定义查询方法
func (u *User) FindByEmail(email string) error {
    return db.Where("email = ?", email).First(u).Error
}

func (u *User) FindActiveUsers() ([]User, error) {
    var users []User
    err := db.Where("status = ?", "active").Find(&users).Error
    return users, err
}

// 实例方法
func (u *User) IsAdult() bool {
    return u.Age >= 18
}

func (u *User) FullName() string {
    return u.Name
}
```

### 软删除

```go
type User struct {
    database.Model
    ID        uint       `json:"id"`
    Name      string     `json:"name"`
    Email     string     `json:"email"`
    DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// 软删除
db.Delete(&user, 1) // 设置 DeletedAt 字段

// 查询时包含软删除的记录
var users []User
db.Unscoped().Find(&users)

// 强制删除
db.Unscoped().Delete(&user, 1)

// 恢复软删除的记录
db.Unscoped().Model(&user).Update("DeletedAt", nil)
```

## 📊 性能优化

### 预加载优化

```go
// 避免 N+1 问题
var users []User
db.Preload("Posts").Preload("Profile").Find(&users)

// 条件预加载
db.Preload("Posts", "status = ?", "published").Find(&users)

// 嵌套预加载
db.Preload("Posts.Comments.User").Find(&users)

// 预加载计数
db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
    return db.Select("user_id, count(*) as count").Group("user_id")
}).Find(&users)
```

### 批量操作

```go
// 批量创建
users := []User{
    {Name: "John", Email: "john@example.com"},
    {Name: "Jane", Email: "jane@example.com"},
    {Name: "Bob", Email: "bob@example.com"},
}
db.CreateInBatches(users, 100)

// 批量更新
db.Model(&User{}).Where("age < ?", 18).Update("status", "minor")

// 批量删除
db.Where("status = ?", "inactive").Delete(&User{})
```

### 查询优化

```go
// 使用索引
db.Where("email = ?", "john@example.com").Find(&user)

// 限制查询字段
db.Select("id", "name", "email").Find(&users)

// 使用原始 SQL 优化复杂查询
var users []User
db.Raw("SELECT * FROM users WHERE age > (SELECT AVG(age) FROM users)").Scan(&users)
```

## 📝 总结

Laravel-Go Framework 的 ORM 系统提供了：

1. **简洁性**: 直观的模型定义和操作方法
2. **灵活性**: 支持复杂的关联关系和查询
3. **性能优化**: 内置预加载和批量操作
4. **安全性**: 自动处理 SQL 注入防护
5. **可维护性**: 清晰的代码结构和生命周期钩子

通过合理使用 ORM 系统的各种功能，可以构建出高效、可维护的数据访问层。
