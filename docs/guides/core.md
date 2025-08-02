# 核心系统指南

## 📖 概述

Laravel-Go Framework 的核心系统提供了应用程序的基础架构，包括配置管理、依赖注入容器、服务提供者、应用程序生命周期管理等核心功能。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [核心系统 API 参考](../api/core.md)

## 🚀 快速开始

### 1. 应用程序启动

```go
// 主应用程序
package main

import (
    "laravel-go/framework/core"
    "laravel-go/framework/config"
    "laravel-go/framework/database"
    "laravel-go/framework/http"
)

func main() {
    // 创建应用程序实例
    app := core.NewApplication()

    // 启动应用程序
    if err := app.Start(); err != nil {
        panic(err)
    }
}

// 应用程序结构
type Application struct {
    core.Application
    config   *config.Config
    database *database.Connection
    server   *http.Server
}

func NewApplication() *Application {
    app := &Application{}

    // 初始化配置
    app.config = config.New()

    // 初始化数据库
    app.database = database.NewConnection()

    // 初始化 HTTP 服务器
    app.server = http.NewServer()

    return app
}

func (app *Application) Start() error {
    // 加载配置
    if err := app.loadConfig(); err != nil {
        return err
    }

    // 初始化数据库
    if err := app.initDatabase(); err != nil {
        return err
    }

    // 注册服务提供者
    app.registerServiceProviders()

    // 启动 HTTP 服务器
    return app.server.Start(":8080")
}
```

### 2. 配置管理

```go
// 配置加载
func (app *Application) loadConfig() error {
    // 加载环境变量
    config.LoadEnv(".env")

    // 加载配置文件
    configFiles := []string{
        "config/app.json",
        "config/database.json",
        "config/cache.json",
        "config/queue.json",
    }

    for _, file := range configFiles {
        if err := config.LoadFile(file); err != nil {
            return err
        }
    }

    return nil
}

// 配置使用
func (app *Application) initDatabase() error {
    // 获取数据库配置
    dbConfig := config.Get("database")

    // 设置数据库连接
    app.database.SetConfig(dbConfig)

    // 测试连接
    return app.database.Ping()
}

// 环境配置
type EnvironmentConfig struct {
    AppName    string `json:"app_name"`
    AppEnv     string `json:"app_env"`
    AppDebug   bool   `json:"app_debug"`
    AppURL     string `json:"app_url"`
    AppVersion string `json:"app_version"`
}

func LoadEnvironmentConfig() (*EnvironmentConfig, error) {
    var config EnvironmentConfig

    // 从环境变量加载
    config.AppName = config.Get("APP_NAME", "Laravel-Go")
    config.AppEnv = config.Get("APP_ENV", "production")
    config.AppDebug = config.GetBool("APP_DEBUG", false)
    config.AppURL = config.Get("APP_URL", "http://localhost")
    config.AppVersion = config.Get("APP_VERSION", "1.0.0")

    return &config, nil
}

// 配置验证
func ValidateConfig() error {
    requiredKeys := []string{
        "database.host",
        "database.port",
        "database.name",
        "database.username",
        "database.password",
    }

    for _, key := range requiredKeys {
        if config.Get(key) == "" {
            return fmt.Errorf("missing required config: %s", key)
        }
    }

    return nil
}
```

### 3. 依赖注入容器

```go
// 服务容器
type ServiceContainer struct {
    container *container.Container
}

func NewServiceContainer() *ServiceContainer {
    return &ServiceContainer{
        container: container.New(),
    }
}

// 注册服务
func (sc *ServiceContainer) RegisterServices() {
    // 注册数据库连接
    sc.container.Singleton("database", func() *database.Connection {
        return database.NewConnection()
    })

    // 注册缓存服务
    sc.container.Singleton("cache", func() *cache.Cache {
        return cache.New()
    })

    // 注册队列服务
    sc.container.Singleton("queue", func() *queue.Queue {
        return queue.New()
    })

    // 注册日志服务
    sc.container.Singleton("logger", func() *log.Logger {
        return log.New()
    })

    // 注册用户服务
    sc.container.Bind("user.service", func() *Services.UserService {
        db := sc.container.Make("database").(*database.Connection)
        return Services.NewUserService(db)
    })

    // 注册邮件服务
    sc.container.Bind("email.service", func() *Services.EmailService {
        return Services.NewEmailService()
    })
}

// 解析依赖
func (sc *ServiceContainer) Resolve(abstract interface{}) interface{} {
    return sc.container.Make(abstract)
}

// 在控制器中使用
type UserController struct {
    http.Controller
    userService *Services.UserService
    emailService *Services.EmailService
}

func NewUserController(container *ServiceContainer) *UserController {
    return &UserController{
        userService:  container.Resolve("user.service").(*Services.UserService),
        emailService: container.Resolve("email.service").(*Services.EmailService),
    }
}
```

### 4. 服务提供者

```go
// 服务提供者接口
type ServiceProvider interface {
    Register()
    Boot()
}

// 应用服务提供者
type AppServiceProvider struct {
    app *Application
}

func NewAppServiceProvider(app *Application) *AppServiceProvider {
    return &AppServiceProvider{app: app}
}

func (p *AppServiceProvider) Register() {
    // 注册核心服务
    p.app.Container.Singleton("app", func() *Application {
        return p.app
    })

    p.app.Container.Singleton("config", func() *config.Config {
        return p.app.config
    })

    p.app.Container.Singleton("database", func() *database.Connection {
        return p.app.database
    })
}

func (p *AppServiceProvider) Boot() {
    // 启动时的初始化工作
    p.loadRoutes()
    p.registerMiddleware()
    p.setupErrorHandling()
}

// 数据库服务提供者
type DatabaseServiceProvider struct {
    app *Application
}

func NewDatabaseServiceProvider(app *Application) *DatabaseServiceProvider {
    return &DatabaseServiceProvider{app: app}
}

func (p *DatabaseServiceProvider) Register() {
    // 注册数据库相关服务
    p.app.Container.Singleton("db.connection", func() *database.Connection {
        return p.app.database
    })

    p.app.Container.Bind("db.migrator", func() *database.Migrator {
        return database.NewMigrator(p.app.database)
    })

    p.app.Container.Bind("db.seeder", func() *database.Seeder {
        return database.NewSeeder(p.app.database)
    })
}

func (p *DatabaseServiceProvider) Boot() {
    // 运行数据库迁移
    if p.app.config.GetBool("database.auto_migrate", false) {
        migrator := p.app.Container.Make("db.migrator").(*database.Migrator)
        migrator.Run()
    }
}

// 注册服务提供者
func (app *Application) registerServiceProviders() {
    providers := []ServiceProvider{
        NewAppServiceProvider(app),
        NewDatabaseServiceProvider(app),
        NewCacheServiceProvider(app),
        NewQueueServiceProvider(app),
        NewLogServiceProvider(app),
    }

    for _, provider := range providers {
        provider.Register()
    }

    for _, provider := range providers {
        provider.Boot()
    }
}
```

### 5. 应用程序生命周期

```go
// 应用程序生命周期管理
type ApplicationLifecycle struct {
    app *Application
}

func NewApplicationLifecycle(app *Application) *ApplicationLifecycle {
    return &ApplicationLifecycle{app: app}
}

// 启动前
func (lc *ApplicationLifecycle) BeforeStart() error {
    // 验证配置
    if err := ValidateConfig(); err != nil {
        return err
    }

    // 检查数据库连接
    if err := lc.app.database.Ping(); err != nil {
        return err
    }

    // 初始化缓存
    if err := lc.initCache(); err != nil {
        return err
    }

    return nil
}

// 启动后
func (lc *ApplicationLifecycle) AfterStart() error {
    // 启动队列工作进程
    go lc.startQueueWorkers()

    // 启动定时任务
    go lc.startScheduledTasks()

    // 预热缓存
    go lc.warmupCache()

    return nil
}

// 关闭前
func (lc *ApplicationLifecycle) BeforeShutdown() error {
    // 停止队列工作进程
    lc.stopQueueWorkers()

    // 停止定时任务
    lc.stopScheduledTasks()

    // 保存缓存
    lc.saveCache()

    return nil
}

// 关闭后
func (lc *ApplicationLifecycle) AfterShutdown() error {
    // 关闭数据库连接
    if err := lc.app.database.Close(); err != nil {
        return err
    }

    // 清理临时文件
    lc.cleanupTempFiles()

    return nil
}

// 启动应用程序
func (app *Application) Start() error {
    lifecycle := NewApplicationLifecycle(app)

    // 启动前
    if err := lifecycle.BeforeStart(); err != nil {
        return err
    }

    // 启动 HTTP 服务器
    go func() {
        if err := app.server.Start(":8080"); err != nil {
            panic(err)
        }
    }()

    // 启动后
    if err := lifecycle.AfterStart(); err != nil {
        return err
    }

    // 等待信号
    app.waitForShutdownSignal()

    return nil
}

// 等待关闭信号
func (app *Application) waitForShutdownSignal() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan

    lifecycle := NewApplicationLifecycle(app)

    // 关闭前
    lifecycle.BeforeShutdown()

    // 关闭 HTTP 服务器
    app.server.Shutdown()

    // 关闭后
    lifecycle.AfterShutdown()
}
```

### 6. 错误处理

```go
// 全局错误处理
type ErrorHandler struct {
    logger *log.Logger
}

func NewErrorHandler() *ErrorHandler {
    return &ErrorHandler{
        logger: log.New(),
    }
}

// 处理错误
func (h *ErrorHandler) Handle(err error) {
    // 记录错误
    h.logger.Error("Application error", map[string]interface{}{
        "error": err.Error(),
        "stack": string(debug.Stack()),
        "time":  time.Now(),
    })

    // 根据错误类型处理
    switch e := err.(type) {
    case *DatabaseError:
        h.handleDatabaseError(e)
    case *ValidationError:
        h.handleValidationError(e)
    case *AuthenticationError:
        h.handleAuthenticationError(e)
    default:
        h.handleGenericError(err)
    }
}

// 处理数据库错误
func (h *ErrorHandler) handleDatabaseError(err *DatabaseError) {
    // 记录数据库错误
    h.logger.Error("Database error", map[string]interface{}{
        "sql":   err.SQL,
        "error": err.Error(),
    })

    // 发送告警
    h.sendAlert("Database error occurred", err)
}

// 处理验证错误
func (h *ErrorHandler) handleValidationError(err *ValidationError) {
    // 记录验证错误
    h.logger.Warning("Validation error", map[string]interface{}{
        "field":   err.Field,
        "message": err.Message,
    })
}

// 处理认证错误
func (h *ErrorHandler) handleAuthenticationError(err *AuthenticationError) {
    // 记录认证错误
    h.logger.Warning("Authentication error", map[string]interface{}{
        "user_id": err.UserID,
        "ip":      err.IP,
        "error":   err.Error(),
    })
}

// 发送告警
func (h *ErrorHandler) sendAlert(message string, err error) {
    // 实现告警逻辑
    // 例如发送邮件、短信或推送到监控系统
}
```

### 7. 中间件管理

```go
// 中间件管理器
type MiddlewareManager struct {
    middlewares map[string]http.Middleware
}

func NewMiddlewareManager() *MiddlewareManager {
    return &MiddlewareManager{
        middlewares: make(map[string]http.Middleware),
    }
}

// 注册中间件
func (mm *MiddlewareManager) Register(name string, middleware http.Middleware) {
    mm.middlewares[name] = middleware
}

// 获取中间件
func (mm *MiddlewareManager) Get(name string) (http.Middleware, bool) {
    middleware, exists := mm.middlewares[name]
    return middleware, exists
}

// 注册全局中间件
func (mm *MiddlewareManager) RegisterGlobalMiddlewares() {
    mm.Register("auth", &middleware.AuthMiddleware{})
    mm.Register("cors", &middleware.CORSMiddleware{})
    mm.Register("log", &middleware.LogMiddleware{})
    mm.Register("cache", &middleware.CacheMiddleware{})
    mm.Register("rate_limit", &middleware.RateLimitMiddleware{})
}

// 应用中间件到路由
func (mm *MiddlewareManager) ApplyToRoute(route *routing.Route, middlewareNames []string) {
    for _, name := range middlewareNames {
        if middleware, exists := mm.Get(name); exists {
            route.Use(middleware)
        }
    }
}
```

### 8. 事件系统

```go
// 事件调度器
type EventDispatcher struct {
    listeners map[string][]event.Listener
}

func NewEventDispatcher() *EventDispatcher {
    return &EventDispatcher{
        listeners: make(map[string][]event.Listener),
    }
}

// 注册事件监听器
func (ed *EventDispatcher) Listen(eventName string, listener event.Listener) {
    ed.listeners[eventName] = append(ed.listeners[eventName], listener)
}

// 触发事件
func (ed *EventDispatcher) Dispatch(eventName string, event interface{}) error {
    listeners, exists := ed.listeners[eventName]
    if !exists {
        return nil
    }

    for _, listener := range listeners {
        if err := listener.Handle(event); err != nil {
            return err
        }
    }

    return nil
}

// 注册应用程序事件
func (ed *EventDispatcher) RegisterAppEvents() {
    // 应用程序启动事件
    ed.Listen("app.started", &AppStartedListener{})

    // 应用程序关闭事件
    ed.Listen("app.shutdown", &AppShutdownListener{})

    // 用户注册事件
    ed.Listen("user.registered", &UserRegisteredListener{})

    // 用户登录事件
    ed.Listen("user.logged_in", &UserLoggedInListener{})
}

// 事件监听器
type AppStartedListener struct{}

func (l *AppStartedListener) Handle(event interface{}) error {
    log.Info("Application started")
    return nil
}

type AppShutdownListener struct{}

func (l *AppShutdownListener) Handle(event interface{}) error {
    log.Info("Application shutting down")
    return nil
}
```

### 9. 缓存管理

```go
// 缓存管理器
type CacheManager struct {
    cache *cache.Cache
}

func NewCacheManager() *CacheManager {
    return &CacheManager{
        cache: cache.New(),
    }
}

// 初始化缓存
func (cm *CacheManager) Init() error {
    // 配置缓存驱动
    driver := config.Get("cache.default", "file")

    switch driver {
    case "redis":
        return cm.initRedisCache()
    case "memory":
        return cm.initMemoryCache()
    case "file":
        return cm.initFileCache()
    default:
        return fmt.Errorf("unsupported cache driver: %s", driver)
    }
}

// 初始化 Redis 缓存
func (cm *CacheManager) initRedisCache() error {
    redisConfig := config.Get("cache.redis")

    cm.cache.SetDriver("redis", redisConfig)
    return cm.cache.Connect()
}

// 初始化内存缓存
func (cm *CacheManager) initMemoryCache() error {
    cm.cache.SetDriver("memory", nil)
    return nil
}

// 初始化文件缓存
func (cm *CacheManager) initFileCache() error {
    fileConfig := config.Get("cache.file")

    cm.cache.SetDriver("file", fileConfig)
    return nil
}

// 预热缓存
func (cm *CacheManager) Warmup() error {
    // 预热用户缓存
    if err := cm.warmupUserCache(); err != nil {
        return err
    }

    // 预热配置缓存
    if err := cm.warmupConfigCache(); err != nil {
        return err
    }

    return nil
}

// 预热用户缓存
func (cm *CacheManager) warmupUserCache() error {
    // 获取热门用户
    users, err := userService.GetPopularUsers(100)
    if err != nil {
        return err
    }

    // 缓存用户数据
    for _, user := range users {
        key := fmt.Sprintf("user:%d", user.ID)
        cm.cache.Set(key, user, time.Hour)
    }

    return nil
}
```

### 10. 监控和健康检查

```go
// 健康检查器
type HealthChecker struct {
    checks map[string]HealthCheck
}

type HealthCheck func() error

func NewHealthChecker() *HealthChecker {
    return &HealthChecker{
        checks: make(map[string]HealthCheck),
    }
}

// 注册健康检查
func (hc *HealthChecker) Register(name string, check HealthCheck) {
    hc.checks[name] = check
}

// 执行健康检查
func (hc *HealthChecker) Check() map[string]interface{} {
    results := make(map[string]interface{})

    for name, check := range hc.checks {
        start := time.Now()
        err := check()
        duration := time.Since(start)

        results[name] = map[string]interface{}{
            "status":   err == nil,
            "duration": duration.String(),
            "error":    err,
        }
    }

    return results
}

// 注册应用程序健康检查
func (hc *HealthChecker) RegisterAppChecks() {
    // 数据库健康检查
    hc.Register("database", func() error {
        return database.Ping()
    })

    // Redis 健康检查
    hc.Register("redis", func() error {
        return cache.Ping()
    })

    // 磁盘空间检查
    hc.Register("disk", func() error {
        return hc.checkDiskSpace()
    })

    // 内存使用检查
    hc.Register("memory", func() error {
        return hc.checkMemoryUsage()
    })
}

// 检查磁盘空间
func (hc *HealthChecker) checkDiskSpace() error {
    // 实现磁盘空间检查逻辑
    return nil
}

// 检查内存使用
func (hc *HealthChecker) checkMemoryUsage() error {
    // 实现内存使用检查逻辑
    return nil
}

// 健康检查端点
func (hc *HealthChecker) HealthCheckHandler(request http.Request) http.Response {
    results := hc.Check()

    // 检查是否有失败的健康检查
    hasFailure := false
    for _, result := range results {
        if !result.(map[string]interface{})["status"].(bool) {
            hasFailure = true
            break
        }
    }

    statusCode := 200
    if hasFailure {
        statusCode = 503
    }

    return http.Response{
        StatusCode: statusCode,
        Body:       toJSON(results),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }
}
```

## 📚 总结

Laravel-Go Framework 的核心系统提供了：

1. **应用程序管理**: 应用程序生命周期、启动和关闭
2. **配置管理**: 环境配置、配置文件加载、配置验证
3. **依赖注入**: 服务容器、服务注册、依赖解析
4. **服务提供者**: 服务注册、启动初始化
5. **错误处理**: 全局错误处理、错误分类处理
6. **中间件管理**: 中间件注册、应用和管理
7. **事件系统**: 事件调度、事件监听器
8. **缓存管理**: 缓存初始化、预热、驱动管理
9. **监控检查**: 健康检查、系统监控
10. **生命周期**: 应用程序启动、运行、关闭流程

通过合理使用核心系统，可以构建稳定、可维护的应用程序架构。
