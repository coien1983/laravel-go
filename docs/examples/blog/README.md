# 博客系统示例

## 📝 项目概览

这是一个使用 Laravel-Go Framework 构建的完整博客系统示例，展示了框架的核心功能和最佳实践。

## 🚀 功能特性

- ✅ 用户认证和授权
- ✅ 文章管理（CRUD）
- ✅ 分类和标签管理
- ✅ 评论系统
- ✅ 文件上传
- ✅ 搜索功能
- ✅ 缓存优化
- ✅ API 接口

## 📁 项目结构

```
blog/
├── app/
│   ├── Http/
│   │   ├── Controllers/
│   │   │   ├── AuthController.go
│   │   │   ├── PostController.go
│   │   │   ├── CategoryController.go
│   │   │   ├── CommentController.go
│   │   │   └── UserController.go
│   │   ├── Middleware/
│   │   │   ├── AuthMiddleware.go
│   │   │   └── AdminMiddleware.go
│   │   └── Requests/
│   │       ├── LoginRequest.go
│   │       ├── RegisterRequest.go
│   │       └── PostRequest.go
│   ├── Models/
│   │   ├── User.go
│   │   ├── Post.go
│   │   ├── Category.go
│   │   ├── Tag.go
│   │   └── Comment.go
│   ├── Services/
│   │   ├── AuthService.go
│   │   ├── PostService.go
│   │   ├── FileService.go
│   │   └── SearchService.go
│   └── Jobs/
│       ├── SendWelcomeEmailJob.go
│       └── ProcessImageJob.go
├── config/
│   ├── app.go
│   ├── database.go
│   ├── cache.go
│   └── queue.go
├── database/
│   ├── migrations/
│   │   ├── create_users_table.go
│   │   ├── create_posts_table.go
│   │   ├── create_categories_table.go
│   │   ├── create_tags_table.go
│   │   ├── create_comments_table.go
│   │   └── create_post_tag_table.go
│   └── seeders/
│       ├── UserSeeder.go
│       ├── CategorySeeder.go
│       └── PostSeeder.go
├── routes/
│   ├── web.go
│   └── api.go
├── storage/
│   ├── uploads/
│   └── cache/
├── tests/
│   ├── auth_test.go
│   ├── post_test.go
│   └── integration_test.go
└── main.go
```

## 🏗️ 核心组件

### 1. 数据模型

#### User 模型

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
    Avatar    string    `json:"avatar" gorm:"size:500"`
    Bio       string    `json:"bio" gorm:"type:text"`
    Role      string    `json:"role" gorm:"size:20;default:'user'"`
    IsActive  bool      `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // 关联关系
    Posts    []Post    `json:"posts,omitempty" gorm:"foreignKey:UserID"`
    Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:UserID"`
}

func (u *User) SetPassword(password string) {
    u.Password = hashPassword(password)
}

func (u *User) CheckPassword(password string) bool {
    return checkPassword(password, u.Password)
}

func (u *User) IsAdmin() bool {
    return u.Role == "admin"
}
```

#### Post 模型

```go
// app/Models/Post.go
package models

import (
    "laravel-go/framework/database"
    "time"
)

type Post struct {
    database.Model
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title" gorm:"size:255;not null"`
    Slug        string    `json:"slug" gorm:"size:255;unique;not null"`
    Content     string    `json:"content" gorm:"type:text;not null"`
    Excerpt     string    `json:"excerpt" gorm:"type:text"`
    FeaturedImage string  `json:"featured_image" gorm:"size:500"`
    Status      string    `json:"status" gorm:"size:20;default:'draft'"`
    ViewCount   int       `json:"view_count" gorm:"default:0"`
    UserID      uint      `json:"user_id" gorm:"not null"`
    CategoryID  uint      `json:"category_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // 关联关系
    User      User       `json:"user" gorm:"foreignKey:UserID"`
    Category  Category   `json:"category" gorm:"foreignKey:CategoryID"`
    Tags      []Tag      `json:"tags" gorm:"many2many:post_tags;"`
    Comments  []Comment  `json:"comments" gorm:"foreignKey:PostID"`
}

func (p *Post) BeforeCreate() error {
    if p.Slug == "" {
        p.Slug = generateSlug(p.Title)
    }
    return nil
}

func (p *Post) IncrementViewCount() {
    p.ViewCount++
}
```

### 2. 控制器

#### PostController

```go
// app/Http/Controllers/PostController.go
package controllers

import (
    "laravel-go/framework/http"
    "laravel-go/app/Models"
    "laravel-go/app/Services"
)

type PostController struct {
    http.Controller
    postService *Services.PostService
}

func NewPostController(postService *Services.PostService) *PostController {
    return &PostController{
        postService: postService,
    }
}

// Index - 获取文章列表
func (c *PostController) Index(request http.Request) http.Response {
    page := c.getPageParam(request)
    limit := c.getLimitParam(request)
    category := request.Query["category"]
    tag := request.Query["tag"]
    search := request.Query["search"]
    
    posts, total := c.postService.GetPosts(page, limit, category, tag, search)
    
    return c.Json(map[string]interface{}{
        "data":  posts,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// Show - 获取单个文章
func (c *PostController) Show(slug string) http.Response {
    post, err := c.postService.GetPostBySlug(slug)
    if err != nil {
        return c.JsonError("Post not found", 404)
    }
    
    // 增加浏览次数
    c.postService.IncrementViewCount(post.ID)
    
    return c.Json(post)
}

// Store - 创建文章
func (c *PostController) Store(request http.Request) http.Response {
    // 验证输入
    if err := c.validatePostData(request.Body); err != nil {
        return c.JsonError(err.Error(), 422)
    }
    
    // 获取当前用户
    user := request.Context["user"].(*Models.User)
    
    // 创建文章
    post, err := c.postService.CreatePost(request.Body, user.ID)
    if err != nil {
        return c.JsonError("Failed to create post", 500)
    }
    
    return c.Json(post).Status(201)
}

// Update - 更新文章
func (c *PostController) Update(id string, request http.Request) http.Response {
    postID, err := strconv.Atoi(id)
    if err != nil {
        return c.JsonError("Invalid post ID", 400)
    }
    
    // 验证输入
    if err := c.validatePostData(request.Body); err != nil {
        return c.JsonError(err.Error(), 422)
    }
    
    // 获取当前用户
    user := request.Context["user"].(*Models.User)
    
    // 检查权限
    post, err := c.postService.GetPost(postID)
    if err != nil {
        return c.JsonError("Post not found", 404)
    }
    
    if post.UserID != user.ID && !user.IsAdmin() {
        return c.JsonError("Unauthorized", 403)
    }
    
    // 更新文章
    updatedPost, err := c.postService.UpdatePost(postID, request.Body)
    if err != nil {
        return c.JsonError("Failed to update post", 500)
    }
    
    return c.Json(updatedPost)
}

// Delete - 删除文章
func (c *PostController) Delete(id string, request http.Request) http.Response {
    postID, err := strconv.Atoi(id)
    if err != nil {
        return c.JsonError("Invalid post ID", 400)
    }
    
    // 获取当前用户
    user := request.Context["user"].(*Models.User)
    
    // 检查权限
    post, err := c.postService.GetPost(postID)
    if err != nil {
        return c.JsonError("Post not found", 404)
    }
    
    if post.UserID != user.ID && !user.IsAdmin() {
        return c.JsonError("Unauthorized", 403)
    }
    
    // 删除文章
    err = c.postService.DeletePost(postID)
    if err != nil {
        return c.JsonError("Failed to delete post", 500)
    }
    
    return c.Json(map[string]string{
        "message": "Post deleted successfully",
    })
}
```

### 3. 服务层

#### PostService

```go
// app/Services/PostService.go
package services

import (
    "laravel-go/framework/database"
    "laravel-go/framework/cache"
    "laravel-go/app/Models"
    "time"
)

type PostService struct {
    db    *database.Connection
    cache cache.Cache
}

func NewPostService(db *database.Connection, cache cache.Cache) *PostService {
    return &PostService{
        db:    db,
        cache: cache,
    }
}

// GetPosts - 获取文章列表
func (s *PostService) GetPosts(page, limit int, category, tag, search string) ([]*Models.Post, int64) {
    cacheKey := fmt.Sprintf("posts:list:%d:%d:%s:%s:%s", page, limit, category, tag, search)
    
    // 尝试从缓存获取
    if cached, exists := s.cache.Get(cacheKey); exists {
        return cached.([]*Models.Post), 0 // 这里应该返回总数，简化处理
    }
    
    query := s.db.Model(&Models.Post{}).Preload("User").Preload("Category").Preload("Tags")
    
    // 应用过滤条件
    if category != "" {
        query = query.Where("category_id = ?", category)
    }
    
    if tag != "" {
        query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
            Joins("JOIN tags ON post_tags.tag_id = tags.id").
            Where("tags.name = ?", tag)
    }
    
    if search != "" {
        query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
    }
    
    // 只显示已发布的文章
    query = query.Where("status = ?", "published")
    
    // 获取总数
    var total int64
    query.Count(&total)
    
    // 分页
    offset := (page - 1) * limit
    var posts []*Models.Post
    query.Offset(offset).Limit(limit).Order("created_at desc").Find(&posts)
    
    // 缓存结果
    s.cache.Put(cacheKey, posts, time.Hour)
    
    return posts, total
}

// GetPostBySlug - 根据 slug 获取文章
func (s *PostService) GetPostBySlug(slug string) (*Models.Post, error) {
    cacheKey := fmt.Sprintf("post:slug:%s", slug)
    
    // 尝试从缓存获取
    if cached, exists := s.cache.Get(cacheKey); exists {
        return cached.(*Models.Post), nil
    }
    
    var post Models.Post
    err := s.db.Preload("User").Preload("Category").Preload("Tags").
        Preload("Comments.User").Where("slug = ? AND status = ?", slug, "published").
        First(&post).Error
    
    if err != nil {
        return nil, err
    }
    
    // 缓存结果
    s.cache.Put(cacheKey, &post, time.Hour)
    
    return &post, nil
}

// CreatePost - 创建文章
func (s *PostService) CreatePost(data map[string]interface{}, userID uint) (*Models.Post, error) {
    post := &Models.Post{
        Title:      data["title"].(string),
        Content:    data["content"].(string),
        Excerpt:    data["excerpt"].(string),
        Status:     data["status"].(string),
        UserID:     userID,
        CategoryID: uint(data["category_id"].(float64)),
    }
    
    // 处理特色图片
    if featuredImage, ok := data["featured_image"].(string); ok {
        post.FeaturedImage = featuredImage
    }
    
    err := s.db.Create(post).Error
    if err != nil {
        return nil, err
    }
    
    // 处理标签
    if tags, ok := data["tags"].([]interface{}); ok {
        s.attachTags(post, tags)
    }
    
    // 清除相关缓存
    s.clearPostCache()
    
    return post, nil
}

// UpdatePost - 更新文章
func (s *PostService) UpdatePost(id int, data map[string]interface{}) (*Models.Post, error) {
    var post Models.Post
    err := s.db.First(&post, id).Error
    if err != nil {
        return nil, err
    }
    
    // 更新字段
    if title, ok := data["title"].(string); ok {
        post.Title = title
    }
    
    if content, ok := data["content"].(string); ok {
        post.Content = content
    }
    
    if excerpt, ok := data["excerpt"].(string); ok {
        post.Excerpt = excerpt
    }
    
    if status, ok := data["status"].(string); ok {
        post.Status = status
    }
    
    if categoryID, ok := data["category_id"].(float64); ok {
        post.CategoryID = uint(categoryID)
    }
    
    if featuredImage, ok := data["featured_image"].(string); ok {
        post.FeaturedImage = featuredImage
    }
    
    err = s.db.Save(&post).Error
    if err != nil {
        return nil, err
    }
    
    // 处理标签
    if tags, ok := data["tags"].([]interface{}); ok {
        s.detachTags(&post)
        s.attachTags(&post, tags)
    }
    
    // 清除相关缓存
    s.clearPostCache()
    
    return &post, nil
}

// DeletePost - 删除文章
func (s *PostService) DeletePost(id int) error {
    var post Models.Post
    err := s.db.First(&post, id).Error
    if err != nil {
        return err
    }
    
    // 删除关联的标签
    s.detachTags(&post)
    
    // 删除文章
    err = s.db.Delete(&post).Error
    if err != nil {
        return err
    }
    
    // 清除相关缓存
    s.clearPostCache()
    
    return nil
}

// IncrementViewCount - 增加浏览次数
func (s *PostService) IncrementViewCount(id uint) {
    s.db.Model(&Models.Post{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
}

// 辅助方法
func (s *PostService) attachTags(post *Models.Post, tags []interface{}) {
    for _, tagName := range tags {
        var tag Models.Tag
        s.db.Where("name = ?", tagName).FirstOrCreate(&tag, Models.Tag{Name: tagName.(string)})
        s.db.Model(post).Association("Tags").Append(&tag)
    }
}

func (s *PostService) detachTags(post *Models.Post) {
    s.db.Model(post).Association("Tags").Clear()
}

func (s *PostService) clearPostCache() {
    // 清除文章相关的缓存
    s.cache.Forget("posts:list:*")
    s.cache.Forget("post:slug:*")
}
```

### 4. 中间件

#### AuthMiddleware

```go
// app/Http/Middleware/AuthMiddleware.go
package middleware

import (
    "laravel-go/framework/http"
    "laravel-go/framework/auth"
)

type AuthMiddleware struct {
    http.Middleware
}

func (m *AuthMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    // 获取认证 token
    token := request.Headers["Authorization"]
    if token == "" {
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "Unauthorized"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }
    
    // 验证 token
    user, err := auth.ValidateToken(token)
    if err != nil {
        return http.Response{
            StatusCode: 401,
            Body:       `{"error": "Invalid token"}`,
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        }
    }
    
    // 将用户信息添加到请求上下文
    request.Context["user"] = user
    
    return next(request)
}
```

### 5. 路由配置

```go
// routes/web.go
package routes

import (
    "laravel-go/framework/routing"
    "laravel-go/app/Http/Controllers"
    "laravel-go/app/Http/Middleware"
)

func WebRoutes(router *routing.Router) {
    // 公开路由
    router.Get("/", &Controllers.HomeController{}, "Index")
    router.Get("/posts", &Controllers.PostController{}, "Index")
    router.Get("/posts/{slug}", &Controllers.PostController{}, "Show")
    router.Get("/categories", &Controllers.CategoryController{}, "Index")
    router.Get("/categories/{id}", &Controllers.CategoryController{}, "Show")
    router.Get("/tags", &Controllers.TagController{}, "Index")
    router.Get("/tags/{id}", &Controllers.TagController{}, "Show")
    
    // 认证路由
    router.Post("/auth/login", &Controllers.AuthController{}, "Login")
    router.Post("/auth/register", &Controllers.AuthController{}, "Register")
    router.Post("/auth/logout", &Controllers.AuthController{}, "Logout")
    
    // 需要认证的路由
    router.Group("/api", func(group *routing.Router) {
        group.Use(&middleware.AuthMiddleware{})
        
        // 用户相关
        group.Get("/profile", &Controllers.UserController{}, "Profile")
        group.Put("/profile", &Controllers.UserController{}, "UpdateProfile")
        
        // 文章相关
        group.Post("/posts", &Controllers.PostController{}, "Store")
        group.Put("/posts/{id}", &Controllers.PostController{}, "Update")
        group.Delete("/posts/{id}", &Controllers.PostController{}, "Delete")
        
        // 评论相关
        group.Post("/posts/{id}/comments", &Controllers.CommentController{}, "Store")
        group.Put("/comments/{id}", &Controllers.CommentController{}, "Update")
        group.Delete("/comments/{id}", &Controllers.CommentController{}, "Delete")
    })
    
    // 管理员路由
    router.Group("/admin", func(group *routing.Router) {
        group.Use(&middleware.AuthMiddleware{})
        group.Use(&middleware.AdminMiddleware{})
        
        group.Get("/dashboard", &Controllers.AdminController{}, "Dashboard")
        group.Get("/users", &Controllers.AdminController{}, "Users")
        group.Get("/posts", &Controllers.AdminController{}, "Posts")
        group.Get("/categories", &Controllers.AdminController{}, "Categories")
        group.Get("/comments", &Controllers.AdminController{}, "Comments")
    })
}
```

### 6. 数据库迁移

```go
// database/migrations/create_posts_table.go
package migrations

import (
    "laravel-go/framework/database"
    "laravel-go/framework/database/migration"
)

type CreatePostsTable struct {
    migration.Migration
}

func (m *CreatePostsTable) Up() error {
    return m.Schema.CreateTable("posts", func(table *database.Blueprint) {
        table.Id("id")
        table.String("title", 255).NotNull()
        table.String("slug", 255).Unique().NotNull()
        table.Text("content").NotNull()
        table.Text("excerpt").Nullable()
        table.String("featured_image", 500).Nullable()
        table.Enum("status", []string{"draft", "published", "archived"}).Default("draft")
        table.Integer("view_count").Default(0)
        table.Integer("user_id").Unsigned().NotNull()
        table.Integer("category_id").Unsigned().Nullable()
        table.Timestamps()
        
        // 索引
        table.Index("user_id")
        table.Index("category_id")
        table.Index("status")
        table.Index("created_at")
        
        // 外键
        table.ForeignKey("user_id").References("id").On("users").OnDelete("cascade")
        table.ForeignKey("category_id").References("id").On("categories").OnDelete("set null")
    })
}

func (m *CreatePostsTable) Down() error {
    return m.Schema.DropTable("posts")
}
```

### 7. 队列任务

```go
// app/Jobs/SendWelcomeEmailJob.go
package jobs

import (
    "laravel-go/framework/queue"
    "time"
)

type SendWelcomeEmailJob struct {
    queue.BaseJob
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
}

func NewSendWelcomeEmailJob(userID int, email, name string) *SendWelcomeEmailJob {
    return &SendWelcomeEmailJob{
        BaseJob: queue.BaseJob{
            Queue:       "emails",
            MaxAttempts: 3,
            Timeout:     time.Minute * 5,
        },
        UserID: userID,
        Email:  email,
        Name:   name,
    }
}

func (j *SendWelcomeEmailJob) Handle() error {
    // 发送欢迎邮件
    return sendWelcomeEmail(j.Email, j.Name)
}

func (j *SendWelcomeEmailJob) Failed(err error) {
    log.Printf("Failed to send welcome email to %s: %v", j.Email, err)
}
```

## 🚀 运行项目

### 1. 环境配置

```bash
# 复制环境配置文件
cp .env.example .env

# 配置数据库
DB_CONNECTION=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=blog
DB_USERNAME=root
DB_PASSWORD=password

# 配置缓存
CACHE_DRIVER=redis
CACHE_REDIS_HOST=localhost
CACHE_REDIS_PORT=6379

# 配置队列
QUEUE_DRIVER=redis
QUEUE_REDIS_HOST=localhost
QUEUE_REDIS_PORT=6379
```

### 2. 安装依赖

```bash
# 安装 Go 依赖
go mod tidy

# 安装前端依赖（如果有）
npm install
```

### 3. 数据库设置

```bash
# 运行迁移
go run cmd/artisan/main.go migrate

# 运行数据填充
go run cmd/artisan/main.go db:seed
```

### 4. 启动服务

```bash
# 启动 Web 服务
go run main.go

# 启动队列工作进程
go run cmd/artisan/main.go queue:work

# 启动调度器
go run cmd/artisan/main.go schedule:run
```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test ./tests/auth_test.go

# 运行集成测试
go test ./tests/integration_test.go
```

### API 测试

```bash
# 测试用户注册
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'

# 测试用户登录
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'

# 测试获取文章列表
curl http://localhost:8080/posts

# 测试创建文章（需要认证）
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"title":"My First Post","content":"This is my first post content","status":"published"}'
```

## 📊 性能优化

### 1. 缓存策略

```go
// 文章列表缓存
func (s *PostService) GetPosts(page, limit int, category, tag, search string) ([]*Models.Post, int64) {
    cacheKey := fmt.Sprintf("posts:list:%d:%d:%s:%s:%s", page, limit, category, tag, search)
    
    if cached, exists := s.cache.Get(cacheKey); exists {
        return cached.([]*Models.Post), 0
    }
    
    // 查询数据库...
    
    // 缓存结果
    s.cache.Put(cacheKey, posts, time.Hour)
    return posts, total
}
```

### 2. 数据库优化

```go
// 使用预加载避免 N+1 问题
query := s.db.Model(&Models.Post{}).
    Preload("User").
    Preload("Category").
    Preload("Tags").
    Preload("Comments.User")
```

### 3. 队列处理

```go
// 异步处理耗时操作
func (s *UserService) RegisterUser(data map[string]interface{}) (*Models.User, error) {
    user, err := s.createUser(data)
    if err != nil {
        return nil, err
    }
    
    // 异步发送欢迎邮件
    job := jobs.NewSendWelcomeEmailJob(int(user.ID), user.Email, user.Name)
    s.queue.Push(job)
    
    return user, nil
}
```

## 🔒 安全考虑

### 1. 输入验证

```go
// 验证用户输入
func (c *PostController) validatePostData(data map[string]interface{}) error {
    if title, ok := data["title"].(string); !ok || len(title) == 0 {
        return errors.New("title is required")
    }
    
    if content, ok := data["content"].(string); !ok || len(content) == 0 {
        return errors.New("content is required")
    }
    
    return nil
}
```

### 2. 权限控制

```go
// 检查用户权限
func (c *PostController) Update(id string, request http.Request) http.Response {
    user := request.Context["user"].(*Models.User)
    
    post, err := c.postService.GetPost(postID)
    if err != nil {
        return c.JsonError("Post not found", 404)
    }
    
    // 只有文章作者或管理员可以编辑
    if post.UserID != user.ID && !user.IsAdmin() {
        return c.JsonError("Unauthorized", 403)
    }
    
    // 更新文章...
}
```

### 3. SQL 注入防护

```go
// 使用参数化查询
query := s.db.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
```

## 📚 总结

这个博客系统示例展示了 Laravel-Go Framework 的核心功能：

1. **MVC 架构**: 清晰的模型、视图、控制器分离
2. **数据库操作**: ORM 和查询构建器
3. **缓存系统**: Redis 缓存优化
4. **队列处理**: 异步任务处理
5. **认证授权**: 用户认证和权限控制
6. **API 设计**: RESTful API 接口
7. **错误处理**: 统一的错误处理机制
8. **测试**: 单元测试和集成测试

通过这个示例，开发者可以学习如何使用 Laravel-Go Framework 构建完整的 Web 应用程序。 