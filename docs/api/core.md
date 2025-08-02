# 核心 API 参考

本文档提供 Laravel-Go Framework 核心组件的 API 参考。

## 📦 Application

应用实例是框架的核心，管理所有服务和组件。

### 创建应用

```go
// 创建新的应用实例
app := framework.NewApplication()

// 创建带配置的应用实例
app := framework.NewApplicationWithConfig(config)
```

### 应用方法

#### Run(addr string) error

启动 HTTP 服务器。

```go
// 启动服务器
err := app.Run(":8080")
if err != nil {
    log.Fatal(err)
}
```

#### Container() \*Container

获取应用容器实例。

```go
container := app.Container()
```

#### Config() \*Config

获取配置管理器实例。

```go
config := app.Config()
```

#### Router() \*Router

获取路由管理器实例。

```go
router := app.Router()
```

#### DB() \*Database

获取数据库连接实例。

```go
db := app.DB()
```

#### Cache() \*Cache

获取缓存管理器实例。

```go
cache := app.Cache()
```

#### Queue() \*Queue

获取队列管理器实例。

```go
queue := app.Queue()
```

#### Events() \*EventDispatcher

获取事件分发器实例。

```go
events := app.Events()
```

#### Auth() \*Auth

获取认证管理器实例。

```go
auth := app.Auth()
```

#### Metrics() \*Metrics

获取性能监控实例。

```go
metrics := app.Metrics()
```

## 🔧 Container

依赖注入容器，管理服务注册和解析。

### 注册服务

#### Singleton(abstract string, concrete interface{})

注册单例服务。

```go
app.Container().Singleton("db", func() interface{} {
    return database.NewConnection()
})
```

#### Bind(abstract string, concrete interface{})

注册绑定服务（每次解析都创建新实例）。

```go
app.Container().Bind("mailer", func() interface{} {
    return mail.NewMailer()
})
```

#### Instance(abstract string, instance interface{})

注册已存在的实例。

```go
db := database.NewConnection()
app.Container().Instance("db", db)
```

### 解析服务

#### Make(abstract string) interface{}

解析服务实例。

```go
db := app.Container().Make("db").(database.Connection)
```

#### Call(abstract string, method string, args ...interface{}) ([]interface{}, error)

调用服务方法。

```go
result, err := app.Container().Call("mailer", "Send", "user@example.com", "Welcome")
```

### 检查服务

#### Bound(abstract string) bool

检查服务是否已注册。

```go
if app.Container().Bound("db") {
    // 服务已注册
}
```

#### Resolved(abstract string) bool

检查服务是否已解析。

```go
if app.Container().Resolved("db") {
    // 服务已解析
}
```

## ⚙️ Config

配置管理器，处理应用配置。

### 获取配置

#### Get(key string) interface{}

获取配置值。

```go
// 获取简单配置
debug := app.Config().Get("app.debug").(bool)

// 获取嵌套配置
dbHost := app.Config().Get("database.host").(string)
```

#### GetString(key string) string

获取字符串配置。

```go
appName := app.Config().GetString("app.name")
```

#### GetInt(key string) int

获取整数配置。

```go
port := app.Config().GetInt("app.port")
```

#### GetBool(key string) bool

获取布尔配置。

```go
debug := app.Config().GetBool("app.debug")
```

#### GetFloat(key string) float64

获取浮点数配置。

```go
version := app.Config().GetFloat("app.version")
```

### 设置配置

#### Set(key string, value interface{})

设置配置值。

```go
app.Config().Set("app.name", "My App")
app.Config().Set("database.host", "localhost")
```

### 环境配置

#### LoadFromFile(path string) error

从文件加载配置。

```go
err := app.Config().LoadFromFile("config/app.json")
```

#### LoadFromEnv() error

从环境变量加载配置。

```go
err := app.Config().LoadFromEnv()
```

#### Has(key string) bool

检查配置是否存在。

```go
if app.Config().Has("database.host") {
    // 配置存在
}
```

## 🗄️ Database

数据库管理器，提供数据库操作接口。

### 连接管理

#### Connection(name string) \*Connection

获取指定连接。

```go
db := app.DB().Connection("mysql")
```

#### DefaultConnection() \*Connection

获取默认连接。

```go
db := app.DB().DefaultConnection()
```

### 查询构建器

#### Table(name string) \*QueryBuilder

开始查询构建。

```go
users := app.DB().Table("users").Get()
```

#### Raw(sql string, args ...interface{}) \*QueryBuilder

执行原始 SQL。

```go
result := app.DB().Raw("SELECT * FROM users WHERE active = ?", true).Get()
```

### 事务

#### Transaction(callback func(\*Transaction) error) error

执行事务。

```go
err := app.DB().Transaction(func(tx *database.Transaction) error {
    // 在事务中执行操作
    tx.Table("users").Create(userData)
    tx.Table("profiles").Create(profileData)
    return nil
})
```

#### Begin() \*Transaction

开始事务。

```go
tx := app.DB().Begin()
defer tx.Rollback()

// 执行操作
tx.Table("users").Create(userData)
tx.Table("profiles").Create(profileData)

// 提交事务
tx.Commit()
```

## 💾 Cache

缓存管理器，提供缓存操作接口。

### 基本操作

#### Set(key string, value interface{}, ttl ...int) error

设置缓存。

```go
// 设置缓存，默认过期时间
app.Cache().Set("user:1", userData)

// 设置缓存，指定过期时间（秒）
app.Cache().Set("user:1", userData, 3600)
```

#### Get(key string) interface{}

获取缓存。

```go
user := app.Cache().Get("user:1")
```

#### Delete(key string) error

删除缓存。

```go
app.Cache().Delete("user:1")
```

#### Has(key string) bool

检查缓存是否存在。

```go
if app.Cache().Has("user:1") {
    // 缓存存在
}
```

#### Flush() error

清空所有缓存。

```go
app.Cache().Flush()
```

### 缓存标签

#### Tags(names ...string) \*TaggedCache

使用缓存标签。

```go
// 设置带标签的缓存
app.Cache().Tags("users", "profiles").Set("user:1", userData)

// 清空标签下的所有缓存
app.Cache().Tags("users").Flush()
```

### 缓存驱动

#### Driver(name string) \*Cache

使用指定驱动。

```go
redisCache := app.Cache().Driver("redis")
memoryCache := app.Cache().Driver("memory")
```

## 📨 Queue

队列管理器，处理异步任务。

### 任务操作

#### Push(job interface{}) error

推送任务到队列。

```go
err := app.Queue().Push(&SendEmailJob{
    To:      "user@example.com",
    Subject: "Welcome",
    Body:    "Welcome to our platform!",
})
```

#### Later(delay time.Duration, job interface{}) error

延迟推送任务。

```go
err := app.Queue().Later(5*time.Minute, &SendEmailJob{
    To:      "user@example.com",
    Subject: "Reminder",
    Body:    "Don't forget to complete your profile!",
})
```

#### PushOn(queue string, job interface{}) error

推送到指定队列。

```go
err := app.Queue().PushOn("emails", &SendEmailJob{
    To:      "user@example.com",
    Subject: "Welcome",
    Body:    "Welcome to our platform!",
})
```

### 队列处理

#### Work(queue string, options ...WorkOption) error

开始处理队列任务。

```go
// 处理默认队列
app.Queue().Work()

// 处理指定队列
app.Queue().Work("emails")

// 处理多个队列
app.Queue().Work("emails", "notifications")
```

### 队列驱动

#### Driver(name string) \*Queue

使用指定驱动。

```go
redisQueue := app.Queue().Driver("redis")
databaseQueue := app.Queue().Driver("database")
```

## 📡 Events

事件分发器，处理事件驱动编程。

### 事件监听

#### Listen(event string, listener interface{})

注册事件监听器。

```go
app.Events().Listen("user.registered", func(event *UserRegisteredEvent) {
    // 处理用户注册事件
    sendWelcomeEmail(event.User)
})
```

#### ListenOnce(event string, listener interface{})

注册一次性事件监听器。

```go
app.Events().ListenOnce("user.registered", func(event *UserRegisteredEvent) {
    // 只执行一次
    sendWelcomeEmail(event.User)
})
```

### 事件分发

#### Dispatch(event interface{}) error

分发事件。

```go
err := app.Events().Dispatch(&UserRegisteredEvent{
    User: user,
})
```

#### Fire(event string, payload interface{}) error

触发事件（字符串形式）。

```go
err := app.Events().Fire("user.registered", map[string]interface{}{
    "user": user,
})
```

### 事件管理

#### Forget(event string)

移除事件监听器。

```go
app.Events().Forget("user.registered")
```

#### ForgetAll()

移除所有事件监听器。

```go
app.Events().ForgetAll()
```

## 🔐 Auth

认证管理器，处理用户认证和授权。

### 用户认证

#### Attempt(credentials map[string]string) bool

尝试验证用户凭据。

```go
success := app.Auth().Attempt(map[string]string{
    "email":    "user@example.com",
    "password": "password",
})
```

#### Login(user interface{}) error

登录用户。

```go
err := app.Auth().Login(user)
```

#### Logout() error

登出用户。

```go
err := app.Auth().Logout()
```

#### User() interface{}

获取当前用户。

```go
user := app.Auth().User()
```

#### Check() bool

检查用户是否已认证。

```go
if app.Auth().Check() {
    // 用户已认证
}
```

### 权限检查

#### Can(ability string) bool

检查用户是否有指定权限。

```go
if app.Auth().User().Can("edit-posts") {
    // 用户可以编辑文章
}
```

#### Cannot(ability string) bool

检查用户是否没有指定权限。

```go
if app.Auth().User().Cannot("delete-posts") {
    // 用户不能删除文章
}
```

## 📊 Metrics

性能监控，收集应用指标。

### 计数器

#### Counter(name string) \*Counter

创建或获取计数器。

```go
counter := app.Metrics().Counter("http_requests_total")
counter.Inc()
counter.Add(5)
```

### 仪表

#### Gauge(name string) \*Gauge

创建或获取仪表。

```go
gauge := app.Metrics().Gauge("active_users")
gauge.Set(100)
gauge.Inc()
gauge.Dec()
```

### 直方图

#### Histogram(name string) \*Histogram

创建或获取直方图。

```go
histogram := app.Metrics().Histogram("http_request_duration")
histogram.Observe(0.5) // 0.5 秒
```

### 指标管理

#### Reset()

重置所有指标。

```go
app.Metrics().Reset()
```

#### Export() map[string]interface{}

导出所有指标。

```go
metrics := app.Metrics().Export()
```

## 🧪 Testing

测试支持，提供测试工具和辅助函数。

### 测试应用

#### NewTestApplication() \*Application

创建测试应用实例。

```go
app := framework.NewTestApplication()
```

#### WithDatabase(config map[string]interface{}) \*Application

配置测试数据库。

```go
app := framework.NewTestApplication().WithDatabase(map[string]interface{}{
    "connection": "sqlite",
    "database":   ":memory:",
})
```

### 测试请求

#### TestRequest(method, path string, body ...interface{}) \*TestResponse

创建测试请求。

```go
response := app.TestRequest("GET", "/users")
response := app.TestRequest("POST", "/users", map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
})
```

### 测试响应

#### Status() int

获取响应状态码。

```go
status := response.Status()
assert.Equal(t, 200, status)
```

#### Body() string

获取响应体。

```go
body := response.Body()
assert.Contains(t, body, "John Doe")
```

#### Json() map[string]interface{}

解析 JSON 响应。

```go
data := response.Json()
assert.Equal(t, "John Doe", data["name"])
```

## 📚 下一步

了解更多 API 细节：

1. [HTTP API](http.md) - HTTP 相关 API
2. [数据库 API](database.md) - 数据库操作 API
3. [ORM API](orm.md) - ORM 模型 API
4. [缓存 API](cache.md) - 缓存操作 API
5. [队列 API](queue.md) - 队列操作 API
6. [事件 API](events.md) - 事件系统 API
7. [认证 API](auth.md) - 认证授权 API
8. [验证 API](validation.md) - 数据验证 API
9. [命令行 API](console.md) - 命令行工具 API

---

这些是 Laravel-Go Framework 的核心 API。掌握这些 API 将帮助你更好地使用框架进行开发！ 🚀
