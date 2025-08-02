# 模板系统 API 参考

## 📋 概述

Laravel-Go Framework 的模板系统提供了强大而灵活的模板渲染功能，支持模板继承、组件化、数据绑定、缓存等特性。模板系统基于 Go 的模板引擎，提供了类似 Blade 模板的语法和功能，适用于构建动态 Web 页面和邮件模板。

## 🏗️ 核心概念

### 模板引擎 (Template Engine)

- 解析和渲染模板文件
- 支持模板继承和包含
- 提供丰富的内置函数

### 模板文件 (Template Files)

- 使用 `.html` 或 `.tpl` 扩展名
- 支持模板语法和 HTML 混合
- 可嵌套和继承

### 模板缓存 (Template Cache)

- 编译模板以提高性能
- 支持缓存失效和更新
- 减少模板解析开销

## 🔧 基础用法

### 1. 基本模板渲染

```go
// 创建模板引擎
engine := template.NewEngine()

// 渲染简单模板
data := map[string]interface{}{
    "name": "John Doe",
    "age":  30,
    "city": "New York",
}

html, err := engine.Render("welcome.html", data)
if err != nil {
    log.Fatal(err)
}

fmt.Println(html)
```

### 2. 模板文件示例

```html
<!-- resources/views/welcome.html -->
<!DOCTYPE html>
<html>
  <head>
    <title>{{ .title | default "Welcome" }}</title>
    <meta charset="utf-8" />
  </head>
  <body>
    <h1>Welcome, {{ .name }}!</h1>

    {{ if .isLoggedIn }}
    <p>You are logged in as {{ .user.name }}</p>
    <a href="/logout">Logout</a>
    {{ else }}
    <p>Please <a href="/login">login</a> to continue.</p>
    {{ end }} {{ if .posts }}
    <h2>Recent Posts</h2>
    <ul>
      {{ range .posts }}
      <li>
        <h3>{{ .title }}</h3>
        <p>{{ .excerpt }}</p>
        <small
          >By {{ .author.name }} on {{ .created_at | date "2006-01-02" }}</small
        >
      </li>
      {{ end }}
    </ul>
    {{ end }} {{ include "partials/footer.html" }}
  </body>
</html>
```

### 3. 在控制器中使用

```go
// app/Http/Controllers/HomeController.go
package controllers

import (
    "laravel-go/framework/http"
    "laravel-go/framework/template"
)

type HomeController struct {
    http.Controller
    template *template.Engine
}

func (c *HomeController) Index(request http.Request) http.Response {
    // 准备数据
    data := map[string]interface{}{
        "title": "Home Page",
        "posts": []map[string]interface{}{
            {
                "title":      "First Post",
                "excerpt":    "This is the first post excerpt",
                "author":     map[string]string{"name": "John Doe"},
                "created_at": time.Now(),
            },
            {
                "title":      "Second Post",
                "excerpt":    "This is the second post excerpt",
                "author":     map[string]string{"name": "Jane Smith"},
                "created_at": time.Now().AddDate(0, 0, -1),
            },
        },
    }

    // 渲染模板
    html, err := c.template.Render("home.html", data)
    if err != nil {
        return c.JsonError("Template rendering failed", 500)
    }

    return c.Html(html)
}
```

## 📚 API 参考

### Engine 接口

```go
type Engine interface {
    Render(template string, data interface{}) (string, error)
    RenderString(template string, data interface{}) (string, error)
    AddFunction(name string, fn interface{})
    AddFunctions(functions map[string]interface{})
    SetDelimiters(left, right string)
    GetDelimiters() (string, string)
    SetCache(cache Cache)
    GetCache() Cache
    ClearCache()
    SetDebug(debug bool)
    IsDebug() bool
    SetRootPath(path string)
    GetRootPath() string
}
```

#### 方法说明

- `Render(template, data)`: 渲染模板文件
- `RenderString(template, data)`: 渲染模板字符串
- `AddFunction(name, fn)`: 添加自定义函数
- `AddFunctions(functions)`: 批量添加自定义函数
- `SetDelimiters(left, right)`: 设置分隔符
- `GetDelimiters()`: 获取分隔符
- `SetCache(cache)`: 设置缓存
- `GetCache()`: 获取缓存
- `ClearCache()`: 清除缓存
- `SetDebug(debug)`: 设置调试模式
- `IsDebug()`: 检查调试模式
- `SetRootPath(path)`: 设置模板根路径
- `GetRootPath()`: 获取模板根路径

### Cache 接口

```go
type Cache interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
    Delete(key string)
    Clear()
    Has(key string) bool
}
```

#### 方法说明

- `Get(key)`: 获取缓存值
- `Set(key, value, ttl)`: 设置缓存值
- `Delete(key)`: 删除缓存
- `Clear()`: 清除所有缓存
- `Has(key)`: 检查缓存是否存在

## 🎯 高级功能

### 1. 模板继承

```html
<!-- resources/views/layouts/app.html -->
<!DOCTYPE html>
<html>
  <head>
    <title>{{ block "title" }}Default Title{{ end }}</title>
    <meta charset="utf-8" />
    {{ block "meta" }}{{ end }}
    <link rel="stylesheet" href="/css/app.css" />
  </head>
  <body>
    <header>{{ include "partials/header.html" }}</header>

    <main>{{ block "content" }}{{ end }}</main>

    <footer>{{ include "partials/footer.html" }}</footer>

    <script src="/js/app.js"></script>
    {{ block "scripts" }}{{ end }}
  </body>
</html>
```

```html
<!-- resources/views/posts/show.html -->
{{ extends "layouts/app.html" }} {{ block "title" }}Post: {{ .post.title }}{{
end }} {{ block "meta" }}
<meta name="description" content="{{ .post.excerpt }}" />
<meta property="og:title" content="{{ .post.title }}" />
<meta property="og:description" content="{{ .post.excerpt }}" />
{{ end }} {{ block "content" }}
<article class="post">
  <header>
    <h1>{{ .post.title }}</h1>
    <div class="meta">
      <span>By {{ .post.author.name }}</span>
      <span>{{ .post.created_at | date "January 2, 2006" }}</span>
    </div>
  </header>

  <div class="content">{{ .post.content | markdown }}</div>

  {{ if .post.tags }}
  <div class="tags">
    {{ range .post.tags }}
    <span class="tag">{{ .name }}</span>
    {{ end }}
  </div>
  {{ end }}
</article>
{{ end }} {{ block "scripts" }}
<script>
  // 页面特定的 JavaScript
  console.log("Post loaded: {{ .post.title }}");
</script>
{{ end }}
```

### 2. 组件系统

```html
<!-- resources/views/components/alert.html -->
<div class="alert alert-{{ .type | default "info" }} {{ if .dismissible }}alert-dismissible{{ end }}">
    {{ if .title }}
        <h4 class="alert-heading">{{ .title }}</h4>
    {{ end }}

    <p>{{ .message }}</p>

    {{ if .dismissible }}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    {{ end }}
</div>
```

```html
<!-- 使用组件 -->
{{ component "alert" .alertData }}

<!-- 或者使用简写语法 -->
{{ alert type="success" title="Success!" message="Operation completed
successfully." dismissible=true }}
```

### 3. 自定义函数

```go
// 注册自定义函数
engine := template.NewEngine()

// 添加单个函数
engine.AddFunction("uppercase", strings.ToUpper)
engine.AddFunction("lowercase", strings.ToLower)
engine.AddFunction("truncate", func(s string, length int) string {
    if len(s) <= length {
        return s
    }
    return s[:length] + "..."
})

// 添加多个函数
functions := map[string]interface{}{
    "formatDate": func(t time.Time, layout string) string {
        return t.Format(layout)
    },
    "pluralize": func(count int, singular, plural string) string {
        if count == 1 {
            return singular
        }
        return plural
    },
    "sum": func(a, b int) int {
        return a + b
    },
}

engine.AddFunctions(functions)
```

### 4. 条件渲染

```html
<!-- 基本条件 -->
{{ if .user }}
<p>Welcome back, {{ .user.name }}!</p>
{{ else }}
<p>Please <a href="/login">login</a> to continue.</p>
{{ end }}

<!-- 多条件 -->
{{ if .user.role == "admin" }}
<div class="admin-panel">
  <h3>Admin Panel</h3>
  <!-- 管理员内容 -->
</div>
{{ else if .user.role == "moderator" }}
<div class="moderator-panel">
  <h3>Moderator Panel</h3>
  <!-- 版主内容 -->
</div>
{{ else }}
<div class="user-panel">
  <h3>User Panel</h3>
  <!-- 普通用户内容 -->
</div>
{{ end }}

<!-- 循环中的条件 -->
{{ range .posts }} {{ if .isPublished }}
<article class="post">
  <h2>{{ .title }}</h2>
  <p>{{ .excerpt }}</p>
</article>
{{ end }} {{ end }}
```

### 5. 循环和迭代

```html
<!-- 基本循环 -->
<ul>
  {{ range .users }}
  <li>{{ .name }} ({{ .email }})</li>
  {{ end }}
</ul>

<!-- 带索引的循环 -->
<table>
  <thead>
    <tr>
      <th>#</th>
      <th>Name</th>
      <th>Email</th>
    </tr>
  </thead>
  <tbody>
    {{ range $index, $user := .users }}
    <tr>
      <td>{{ add $index 1 }}</td>
      <td>{{ $user.name }}</td>
      <td>{{ $user.email }}</td>
    </tr>
    {{ end }}
  </tbody>
</table>

<!-- 嵌套循环 -->
{{ range .categories }}
<div class="category">
  <h2>{{ .name }}</h2>
  <div class="products">
    {{ range .products }}
    <div class="product">
      <h3>{{ .name }}</h3>
      <p>{{ .description }}</p>
      <span class="price">${{ .price }}</span>
    </div>
    {{ end }}
  </div>
</div>
{{ end }}
```

## 🔧 配置选项

### 模板系统配置

```go
// config/template.go
package config

type TemplateConfig struct {
    // 模板根路径
    RootPath string `json:"root_path"`

    // 模板文件扩展名
    Extensions []string `json:"extensions"`

    // 默认分隔符
    Delimiters DelimitersConfig `json:"delimiters"`

    // 缓存配置
    Cache CacheConfig `json:"cache"`

    // 调试模式
    Debug bool `json:"debug"`

    // 自动重载
    AutoReload bool `json:"auto_reload"`

    // 默认布局
    DefaultLayout string `json:"default_layout"`

    // 组件路径
    ComponentsPath string `json:"components_path"`

    // 部分模板路径
    PartialsPath string `json:"partials_path"`
}

type DelimitersConfig struct {
    Left  string `json:"left"`
    Right string `json:"right"`
}

type CacheConfig struct {
    // 是否启用缓存
    Enabled bool `json:"enabled"`

    // 缓存时间
    TTL time.Duration `json:"ttl"`

    // 缓存驱动
    Driver string `json:"driver"`

    // 缓存路径（文件驱动）
    Path string `json:"path"`
}
```

### 配置示例

```go
// config/template.go
func GetTemplateConfig() *TemplateConfig {
    return &TemplateConfig{
        RootPath: "resources/views",
        Extensions: []string{".html", ".tpl"},
        Delimiters: DelimitersConfig{
            Left:  "{{",
            Right: "}}",
        },
        Cache: CacheConfig{
            Enabled: true,
            TTL:     time.Hour,
            Driver:  "file",
            Path:    "storage/cache/templates",
        },
        Debug:         false,
        AutoReload:    true,
        DefaultLayout: "layouts/app.html",
        ComponentsPath: "resources/views/components",
        PartialsPath:  "resources/views/partials",
    }
}
```

## 🚀 性能优化

### 1. 模板缓存

```go
// 启用模板缓存
engine := template.NewEngine()
engine.SetCache(template.NewFileCache("storage/cache/templates"))

// 缓存配置
cacheConfig := template.CacheConfig{
    Enabled: true,
    TTL:     time.Hour,
    Driver:  "file",
    Path:    "storage/cache/templates",
}

engine.SetCache(template.NewCache(cacheConfig))
```

### 2. 模板预编译

```go
// 预编译模板
type PrecompiledEngine struct {
    template.Engine
    compiled map[string]*template.Template
    mutex    sync.RWMutex
}

func (e *PrecompiledEngine) Precompile() error {
    e.mutex.Lock()
    defer e.mutex.Unlock()

    // 扫描模板目录
    templates, err := e.scanTemplates()
    if err != nil {
        return err
    }

    // 编译所有模板
    for _, tmpl := range templates {
        compiled, err := e.compileTemplate(tmpl)
        if err != nil {
            return err
        }
        e.compiled[tmpl] = compiled
    }

    return nil
}

func (e *PrecompiledEngine) Render(template string, data interface{}) (string, error) {
    e.mutex.RLock()
    compiled, exists := e.compiled[template]
    e.mutex.RUnlock()

    if !exists {
        return e.Engine.Render(template, data)
    }

    // 使用预编译的模板
    var buf bytes.Buffer
    err := compiled.Execute(&buf, data)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}
```

### 3. 模板片段缓存

```html
<!-- 缓存模板片段 -->
{{ cache "user-profile" .user.id 3600 }}
<div class="user-profile">
  <img src="{{ .user.avatar }}" alt="{{ .user.name }}" />
  <h3>{{ .user.name }}</h3>
  <p>{{ .user.bio }}</p>
  <div class="stats">
    <span>{{ .user.posts_count }} posts</span>
    <span>{{ .user.followers_count }} followers</span>
  </div>
</div>
{{ end }}
```

## 🧪 测试

### 1. 模板测试

```go
// tests/template_test.go
package tests

import (
    "testing"
    "laravel-go/framework/template"
    "strings"
)

func TestTemplateRendering(t *testing.T) {
    engine := template.NewEngine()
    engine.SetRootPath("testdata/templates")

    // 测试数据
    data := map[string]interface{}{
        "name": "John Doe",
        "age":  30,
    }

    // 渲染模板
    result, err := engine.Render("test.html", data)
    if err != nil {
        t.Fatal(err)
    }

    // 验证结果
    if !strings.Contains(result, "John Doe") {
        t.Error("Template should contain user name")
    }

    if !strings.Contains(result, "30") {
        t.Error("Template should contain user age")
    }
}

func TestTemplateInheritance(t *testing.T) {
    engine := template.NewEngine()
    engine.SetRootPath("testdata/templates")

    data := map[string]interface{}{
        "title": "Test Page",
        "content": "Test content",
    }

    result, err := engine.Render("child.html", data)
    if err != nil {
        t.Fatal(err)
    }

    // 验证继承
    if !strings.Contains(result, "Test Page") {
        t.Error("Child template should inherit title")
    }

    if !strings.Contains(result, "Test content") {
        t.Error("Child template should contain content")
    }
}
```

### 2. 自定义函数测试

```go
func TestCustomFunctions(t *testing.T) {
    engine := template.NewEngine()

    // 添加自定义函数
    engine.AddFunction("uppercase", strings.ToUpper)
    engine.AddFunction("add", func(a, b int) int {
        return a + b
    })

    // 测试模板字符串
    template := "Hello {{ .name | uppercase }}! {{ add .a .b }}"
    data := map[string]interface{}{
        "name": "john",
        "a":    5,
        "b":    3,
    }

    result, err := engine.RenderString(template, data)
    if err != nil {
        t.Fatal(err)
    }

    expected := "Hello JOHN! 8"
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

## 🔍 调试和监控

### 1. 模板调试

```go
// 启用调试模式
engine := template.NewEngine()
engine.SetDebug(true)

// 调试信息会包含模板路径和行号
result, err := engine.Render("debug.html", data)
if err != nil {
    // 错误信息会包含详细的调试信息
    log.Printf("Template error: %v", err)
}
```

### 2. 模板监控

```go
type TemplateMonitor struct {
    template.Engine
    metrics metrics.Collector
}

func (m *TemplateMonitor) Render(template string, data interface{}) (string, error) {
    start := time.Now()

    // 记录渲染指标
    m.metrics.Increment("templates.rendered", map[string]string{
        "template": template,
    })

    result, err := m.Engine.Render(template, data)
    duration := time.Since(start)

    // 记录渲染时间
    m.metrics.Histogram("templates.duration", duration.Seconds(), map[string]string{
        "template": template,
    })

    if err != nil {
        m.metrics.Increment("templates.errors", map[string]string{
            "template": template,
        })
    }

    return result, err
}
```

## 📝 最佳实践

### 1. 模板组织

```go
// 推荐的模板目录结构
resources/
└── views/
    ├── layouts/
    │   ├── app.html
    │   ├── admin.html
    │   └── email.html
    ├── components/
    │   ├── alert.html
    │   ├── button.html
    │   └── modal.html
    ├── partials/
    │   ├── header.html
    │   ├── footer.html
    │   └── sidebar.html
    ├── pages/
    │   ├── home.html
    │   ├── about.html
    │   └── contact.html
    └── emails/
        ├── welcome.html
        ├── reset-password.html
        └── notification.html
```

### 2. 数据传递

```go
// 在控制器中准备数据
func (c *HomeController) Index(request http.Request) http.Response {
    // 获取用户信息
    user := request.Context["user"].(*Models.User)

    // 获取页面数据
    posts, err := c.postService.GetRecentPosts(10)
    if err != nil {
        return c.JsonError("Failed to load posts", 500)
    }

    // 准备模板数据
    data := map[string]interface{}{
        "user": user,
        "posts": posts,
        "page": map[string]interface{}{
            "title": "Home",
            "description": "Welcome to our blog",
        },
        "meta": map[string]interface{}{
            "canonical": "https://example.com",
            "og_image": "https://example.com/og-image.jpg",
        },
    }

    return c.View("pages/home.html", data)
}
```

### 3. 错误处理

```html
<!-- 在模板中处理错误 -->
{{ if .error }} {{ component "alert" type="error" message=.error }} {{ end }} {{
if .success }} {{ component "alert" type="success" message=.success }} {{ end }}

<!-- 处理空数据 -->
{{ if .posts }} {{ range .posts }}
<!-- 显示文章 -->
{{ end }} {{ else }}
<p>No posts found.</p>
{{ end }}
```

### 4. 安全性

```html
<!-- 自动转义 HTML -->
<p>{{ .user_input }}</p>

<!-- 安全输出 HTML -->
<div>{{ .html_content | safe }}</div>

<!-- 条件输出 -->
{{ if .is_admin }}
<div class="admin-content">{{ .admin_content | safe }}</div>
{{ end }}
```

## 🚀 总结

模板系统是 Laravel-Go Framework 中重要的功能之一，它提供了：

1. **完整的模板功能**: 支持模板继承、组件化、条件渲染等
2. **性能优化**: 提供模板缓存和预编译功能
3. **灵活配置**: 支持自定义函数、分隔符等配置
4. **安全性**: 自动 HTML 转义和安全输出
5. **调试支持**: 完整的调试和监控功能
6. **最佳实践**: 遵循模板开发的最佳实践

通过合理使用模板系统，可以构建出结构清晰、易于维护的动态 Web 页面和邮件模板。
