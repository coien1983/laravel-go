# 模板引擎指南

## 📖 概述

Laravel-Go Framework 提供了强大的模板引擎，支持模板继承、组件化、数据绑定、条件渲染等功能，帮助构建动态的 Web 页面。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [模板引擎 API 参考](../api/template.md)

## 🚀 快速开始

### 1. 基本模板使用

```go
// 创建模板引擎
engine := template.NewEngine()

// 注册模板
engine.Register("layout", `
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
    <meta charset="utf-8">
</head>
<body>
    <header>
        {{template "header" .}}
    </header>

    <main>
        {{template "content" .}}
    </main>

    <footer>
        {{template "footer" .}}
    </footer>
</body>
</html>
`)

engine.Register("header", `
<nav>
    <a href="/">Home</a>
    <a href="/about">About</a>
    <a href="/contact">Contact</a>
</nav>
`)

engine.Register("content", `
<div class="container">
    <h1>{{.title}}</h1>
    <p>{{.content}}</p>
</div>
`)

engine.Register("footer", `
<div class="footer">
    <p>&copy; 2024 My App. All rights reserved.</p>
</div>
`)

// 在控制器中使用
func (c *PageController) Show(request http.Request) http.Response {
    data := map[string]interface{}{
        "title":   "Welcome",
        "content": "This is the main content",
    }

    html, err := engine.Render("layout", data)
    if err != nil {
        return c.JsonError("Template error", 500)
    }

    return http.Response{
        StatusCode: 200,
        Body:       html,
        Headers: map[string]string{
            "Content-Type": "text/html",
        },
    }
}
```

### 2. 模板继承

```go
// 基础模板
engine.Register("base", `
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
    <meta charset="utf-8">
    {{template "styles" .}}
</head>
<body>
    <header>
        {{template "header" .}}
    </header>

    <main>
        {{template "main" .}}
    </main>

    <footer>
        {{template "footer" .}}
    </footer>

    {{template "scripts" .}}
</body>
</html>
`)

// 页面特定模板
engine.Register("user.profile", `
{{template "base" .}}

{{define "main"}}
<div class="profile">
    <h1>{{.user.name}}</h1>
    <p>Email: {{.user.email}}</p>
    <p>Joined: {{.user.created_at}}</p>

    <div class="posts">
        <h2>Recent Posts</h2>
        {{range .posts}}
        <div class="post">
            <h3>{{.title}}</h3>
            <p>{{.excerpt}}</p>
            <small>{{.created_at}}</small>
        </div>
        {{end}}
    </div>
</div>
{{end}}
`)

// 在控制器中使用
func (c *UserController) Profile(id string, request http.Request) http.Response {
    user, err := c.userService.GetUser(uint(id))
    if err != nil {
        return c.JsonError("User not found", 404)
    }

    posts, err := c.postService.GetUserPosts(user.ID)
    if err != nil {
        return c.JsonError("Failed to get posts", 500)
    }

    data := map[string]interface{}{
        "title": user.Name + "'s Profile",
        "user":  user,
        "posts": posts,
    }

    html, err := engine.Render("user.profile", data)
    if err != nil {
        return c.JsonError("Template error", 500)
    }

    return http.Response{
        StatusCode: 200,
        Body:       html,
        Headers: map[string]string{
            "Content-Type": "text/html",
        },
    }
}
```

### 3. 组件化模板

```go
// 定义可重用组件
engine.Register("components.button", `
<button class="btn btn-{{.type}}" {{if .disabled}}disabled{{end}}>
    {{.text}}
</button>
`)

engine.Register("components.card", `
<div class="card">
    {{if .header}}
    <div class="card-header">
        <h3>{{.header}}</h3>
    </div>
    {{end}}

    <div class="card-body">
        {{.content}}
    </div>

    {{if .footer}}
    <div class="card-footer">
        {{.footer}}
    </div>
    {{end}}
</div>
`)

engine.Register("components.alert", `
<div class="alert alert-{{.type}}">
    {{if .dismissible}}
    <button type="button" class="close" data-dismiss="alert">
        <span>&times;</span>
    </button>
    {{end}}

    {{.message}}
</div>
`)

// 使用组件
engine.Register("dashboard", `
{{template "base" .}}

{{define "main"}}
<div class="dashboard">
    {{template "components.alert" map[string]interface{}{
        "type": "info",
        "message": "Welcome to your dashboard!",
        "dismissible": true,
    }}}

    <div class="row">
        <div class="col-md-6">
            {{template "components.card" map[string]interface{}{
                "header": "Recent Activity",
                "content": "Your recent activity will appear here.",
            }}}
        </div>

        <div class="col-md-6">
            {{template "components.card" map[string]interface{}{
                "header": "Quick Actions",
                "content": `
                    {{template "components.button" map[string]interface{}{
                        "type": "primary",
                        "text": "Create Post",
                    }}}
                    {{template "components.button" map[string]interface{}{
                        "type": "secondary",
                        "text": "View Profile",
                    }}}
                `,
            }}}
        </div>
    </div>
</div>
{{end}}
`)
```

### 4. 条件渲染

```go
// 条件渲染示例
engine.Register("user.list", `
<div class="users">
    {{if .users}}
        <h2>Users ({{len .users}})</h2>
        <div class="user-grid">
            {{range .users}}
            <div class="user-card">
                <img src="{{.avatar}}" alt="{{.name}}" class="avatar">
                <h3>{{.name}}</h3>
                <p>{{.email}}</p>

                {{if .is_online}}
                <span class="status online">Online</span>
                {{else}}
                <span class="status offline">Offline</span>
                {{end}}

                {{if eq .role "admin"}}
                <span class="badge admin">Admin</span>
                {{else if eq .role "moderator"}}
                <span class="badge moderator">Moderator</span>
                {{end}}
            </div>
            {{end}}
        </div>
    {{else}}
        <div class="empty-state">
            <p>No users found.</p>
            {{if .can_create}}
            <a href="/users/create" class="btn btn-primary">Create User</a>
            {{end}}
        </div>
    {{end}}
</div>
`)

// 在控制器中使用
func (c *UserController) Index(request http.Request) http.Response {
    users, err := c.userService.GetUsers()
    if err != nil {
        return c.JsonError("Failed to get users", 500)
    }

    // 检查用户权限
    currentUser := request.Context["user"].(*Models.User)
    canCreate := currentUser.HasPermission("create_user")

    data := map[string]interface{}{
        "users":      users,
        "can_create": canCreate,
    }

    html, err := engine.Render("user.list", data)
    if err != nil {
        return c.JsonError("Template error", 500)
    }

    return http.Response{
        StatusCode: 200,
        Body:       html,
        Headers: map[string]string{
            "Content-Type": "text/html",
        },
    }
}
```

### 5. 循环和迭代

```go
// 循环渲染示例
engine.Register("post.list", `
<div class="posts">
    {{range $index, $post := .posts}}
    <article class="post" id="post-{{$post.id}}">
        <header class="post-header">
            <h2><a href="/posts/{{$post.id}}">{{$post.title}}</a></h2>
            <div class="post-meta">
                <span class="author">By {{$post.author.name}}</span>
                <span class="date">{{$post.created_at}}</span>
                <span class="category">{{$post.category.name}}</span>
            </div>
        </header>

        <div class="post-content">
            {{$post.excerpt}}
        </div>

        <footer class="post-footer">
            <div class="tags">
                {{range $post.tags}}
                <span class="tag">{{.}}</span>
                {{end}}
            </div>

            <div class="stats">
                <span class="comments">{{$post.comment_count}} comments</span>
                <span class="likes">{{$post.like_count}} likes</span>
            </div>
        </footer>
    </article>
    {{end}}

    {{if .pagination}}
    <div class="pagination">
        {{if .pagination.prev_page}}
        <a href="?page={{.pagination.prev_page}}" class="prev">Previous</a>
        {{end}}

        {{range .pagination.pages}}
        <a href="?page={{.}}" class="page {{if eq . $.pagination.current_page}}active{{end}}">{{.}}</a>
        {{end}}

        {{if .pagination.next_page}}
        <a href="?page={{.pagination.next_page}}" class="next">Next</a>
        {{end}}
    </div>
    {{end}}
</div>
`)
```

### 6. 模板函数

```go
// 注册自定义函数
engine.RegisterFunction("formatDate", func(date time.Time) string {
    return date.Format("2006-01-02 15:04:05")
})

engine.RegisterFunction("truncate", func(text string, length int) string {
    if len(text) <= length {
        return text
    }
    return text[:length] + "..."
})

engine.RegisterFunction("pluralize", func(count int, singular, plural string) string {
    if count == 1 {
        return singular
    }
    return plural
})

engine.RegisterFunction("asset", func(path string) string {
    return "/assets/" + path
})

// 使用自定义函数
engine.Register("post.detail", `
<article class="post">
    <header>
        <h1>{{.post.title}}</h1>
        <div class="meta">
            <span class="author">{{.post.author.name}}</span>
            <span class="date">{{formatDate .post.created_at}}</span>
            <span class="category">{{.post.category.name}}</span>
        </div>
    </header>

    <div class="content">
        {{.post.content}}
    </div>

    <footer>
        <div class="tags">
            {{range .post.tags}}
            <span class="tag">{{.}}</span>
            {{end}}
        </div>

        <div class="stats">
            <span class="comments">{{pluralize .post.comment_count "comment" "comments"}}</span>
            <span class="likes">{{pluralize .post.like_count "like" "likes"}}</span>
        </div>
    </footer>
</article>

<div class="comments">
    <h3>{{pluralize (len .comments) "Comment" "Comments"}}</h3>

    {{range .comments}}
    <div class="comment">
        <div class="comment-meta">
            <span class="author">{{.author.name}}</span>
            <span class="date">{{formatDate .created_at}}</span>
        </div>
        <div class="comment-content">
            {{truncate .content 200}}
        </div>
    </div>
    {{end}}
</div>
`)
```

### 7. 模板缓存

```go
// 启用模板缓存
engine.EnableCache(true)
engine.SetCacheDir("storage/cache/templates")

// 缓存配置
engine.SetCacheOptions(map[string]interface{}{
    "enabled":     true,
    "ttl":         3600, // 1小时
    "auto_reload": true, // 开发环境自动重载
})

// 手动清除缓存
func (c *AdminController) ClearTemplateCache(request http.Request) http.Response {
    err := engine.ClearCache()
    if err != nil {
        return c.JsonError("Failed to clear cache", 500)
    }

    return c.Json(map[string]string{
        "message": "Template cache cleared successfully",
    })
}
```

### 8. 模板布局

```go
// 定义布局
engine.Register("layouts.app", `
<!DOCTYPE html>
<html lang="{{.locale}}">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{.title}} - {{.app_name}}</title>

    <link rel="stylesheet" href="{{asset "css/app.css"}}">
    {{template "styles" .}}
</head>
<body class="{{.body_class}}">
    <div id="app">
        {{template "header" .}}

        <main class="main-content">
            {{template "content" .}}
        </main>

        {{template "footer" .}}
    </div>

    <script src="{{asset "js/app.js"}}"></script>
    {{template "scripts" .}}
</body>
</html>
`)

// 使用布局
engine.Register("pages.home", `
{{template "layouts.app" .}}

{{define "content"}}
<div class="home-page">
    <section class="hero">
        <h1>{{.hero.title}}</h1>
        <p>{{.hero.subtitle}}</p>
        <a href="{{.hero.cta_url}}" class="btn btn-primary">{{.hero.cta_text}}</a>
    </section>

    <section class="features">
        {{range .features}}
        <div class="feature">
            <h3>{{.title}}</h3>
            <p>{{.description}}</p>
        </div>
        {{end}}
    </section>
</div>
{{end}}
`)
```

## 📚 总结

Laravel-Go Framework 的模板引擎提供了：

1. **模板继承**: 支持基础模板和页面特定模板
2. **组件化**: 可重用的模板组件
3. **条件渲染**: 根据数据条件显示内容
4. **循环迭代**: 处理数组和切片数据
5. **自定义函数**: 扩展模板功能
6. **模板缓存**: 提升性能
7. **布局系统**: 统一的页面结构

通过合理使用模板引擎，可以构建动态、可维护的 Web 页面。
