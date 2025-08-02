package main

import (
	"context"
	"fmt"
	"laravel-go/framework/errors"
	"laravel-go/framework/http/middleware"
	"laravel-go/framework/performance"
	"log"
	"net/http"
	"time"
)

// 自定义错误类型
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrDatabaseTimeout  = errors.New("database timeout")
	ErrCacheUnavailable = errors.New("cache unavailable")
)

// User 用户模型
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserService 用户服务
type UserService struct {
	errorHandler errors.ErrorHandler
}

// NewUserService 创建用户服务
func NewUserService(errorHandler errors.ErrorHandler) *UserService {
	return &UserService{
		errorHandler: errorHandler,
	}
}

// GetUser 获取用户
func (s *UserService) GetUser(id int) (*User, error) {
	// 使用安全执行包装器
	var user *User
	var err error

	errors.SafeExecuteWithContext(context.Background(), func() error {
		if id <= 0 {
			err = errors.Wrap(ErrInvalidInput, "invalid user id")
			return err
		}

		// 模拟数据库查询
		if id == 999 {
			// 模拟数据库超时
			time.Sleep(3 * time.Second)
			err = errors.Wrap(ErrDatabaseTimeout, "database query timeout")
			return err
		}

		if id > 100 {
			err = errors.Wrap(ErrUserNotFound, fmt.Sprintf("user %d not found", id))
			return err
		}

		// 模拟成功返回
		user = &User{
			ID:    id,
			Name:  fmt.Sprintf("User %d", id),
			Email: fmt.Sprintf("user%d@example.com", id),
		}
		return nil
	})

	return user, err
}

// CreateUser 创建用户
func (s *UserService) CreateUser(name, email string) (*User, error) {
	var user *User
	var err error

	errors.SafeExecuteWithContext(context.Background(), func() error {
		if name == "" || email == "" {
			err = errors.Wrap(ErrInvalidInput, "name and email are required")
			return err
		}

		// 模拟验证失败
		if name == "error" {
			err = errors.Wrap(ErrInvalidInput, "invalid user name")
			return err
		}

		// 模拟成功创建
		user = &User{
			ID:    123,
			Name:  name,
			Email: email,
		}
		return nil
	})

	return user, err
}

// CacheService 缓存服务
type CacheService struct {
	errorHandler errors.ErrorHandler
	available    bool
}

// NewCacheService 创建缓存服务
func NewCacheService(errorHandler errors.ErrorHandler) *CacheService {
	return &CacheService{
		errorHandler: errorHandler,
		available:    true,
	}
}

// Get 获取缓存
func (s *CacheService) Get(key string) (interface{}, error) {
	var result interface{}
	var err error

	errors.SafeExecuteWithContext(context.Background(), func() error {
		if !s.available {
			err = errors.Wrap(ErrCacheUnavailable, "cache service is down")
			return err
		}

		// 模拟缓存未命中
		if key == "miss" {
			err = errors.Wrap(errors.New("cache miss"), "key not found in cache")
			return err
		}

		result = fmt.Sprintf("cached_value_for_%s", key)
		return nil
	})

	return result, err
}

// Set 设置缓存
func (s *CacheService) Set(key string, value interface{}) error {
	return errors.SafeExecuteWithContext(context.Background(), func() error {
		if !s.available {
			return errors.Wrap(ErrCacheUnavailable, "cache service is down")
		}

		// 模拟设置失败
		if key == "error" {
			return errors.Wrap(errors.New("cache set failed"), "failed to set cache")
		}

		return nil
	})
}

// UserController 用户控制器
type UserController struct {
	userService  *UserService
	cacheService *CacheService
	errorHandler errors.ErrorHandler
}

// NewUserController 创建用户控制器
func NewUserController(userService *UserService, cacheService *CacheService, errorHandler errors.ErrorHandler) *UserController {
	return &UserController{
		userService:  userService,
		cacheService: cacheService,
		errorHandler: errorHandler,
	}
}

// GetUserHandler 获取用户处理器
func (c *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// 解析用户ID
	id := 1 // 简化处理，实际应该从URL参数获取

	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("user:%d", id)
	if cached, err := c.cacheService.Get(cacheKey); err == nil {
		// 缓存命中
		fmt.Fprintf(w, "Cache hit: %v\n", cached)
		return
	}

	// 从数据库获取
	user, err := c.userService.GetUser(id)
	if err != nil {
		// 处理错误
		c.handleError(w, err)
		return
	}

	// 缓存结果
	if err := c.cacheService.Set(cacheKey, user); err != nil {
		// 记录缓存错误，但不影响主流程
		c.errorHandler.Handle(errors.Wrap(err, "failed to cache user"))
	}

	// 返回成功响应
	fmt.Fprintf(w, "User: %+v\n", user)
}

// CreateUserHandler 创建用户处理器
func (c *UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	name := "Test User"
	email := "test@example.com"

	user, err := c.userService.CreateUser(name, email)
	if err != nil {
		c.handleError(w, err)
		return
	}

	fmt.Fprintf(w, "Created user: %+v\n", user)
}

// handleError 处理错误
func (c *UserController) handleError(w http.ResponseWriter, err error) {
	// 使用错误处理器处理错误
	processedErr := c.errorHandler.Handle(err)

	// 根据错误类型返回相应的HTTP状态码
	if appErr := errors.GetAppError(processedErr); appErr != nil {
		http.Error(w, appErr.Message, appErr.Code)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CustomLogger 自定义日志器
type CustomLogger struct{}

func (l *CustomLogger) Error(message string, context map[string]interface{}) {
	log.Printf("[ERROR] %s: %+v", message, context)
}

func (l *CustomLogger) Warning(message string, context map[string]interface{}) {
	log.Printf("[WARN] %s: %+v", message, context)
}

func (l *CustomLogger) Info(message string, context map[string]interface{}) {
	log.Printf("[INFO] %s: %+v", message, context)
}

func (l *CustomLogger) Debug(message string, context map[string]interface{}) {
	log.Printf("[DEBUG] %s: %+v", message, context)
}

// ErrorDemoHandler 错误演示处理器
func ErrorDemoHandler(w http.ResponseWriter, r *http.Request) {
	// 演示不同类型的错误
	errorType := r.URL.Query().Get("type")

	switch errorType {
	case "panic":
		// 演示panic恢复
		panic("this is a panic for testing recovery")
	case "timeout":
		// 演示超时错误
		time.Sleep(3 * time.Second)
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
	case "validation":
		// 演示验证错误
		http.Error(w, "Validation failed", http.StatusBadRequest)
	case "notfound":
		// 演示未找到错误
		http.Error(w, "Resource not found", http.StatusNotFound)
	default:
		// 演示一般错误
		http.Error(w, "General error", http.StatusInternalServerError)
	}
}

func main() {
	// 创建错误处理器
	logger := &CustomLogger{}
	errorHandler := errors.NewDefaultErrorHandler(logger)

	// 创建服务
	userService := NewUserService(errorHandler)
	cacheService := NewCacheService(errorHandler)
	userController := NewUserController(userService, cacheService, errorHandler)

	// 创建性能监控器
	monitor := performance.NewPerformanceMonitor()
	monitor.Start(context.Background())
	defer monitor.Stop()

	// 创建HTTP监控器（用于演示，实际未使用）
	_ = performance.NewHTTPMonitor(monitor)

	// 创建恢复中间件
	recoveryMiddleware := middleware.NewRecoveryMiddleware(errorHandler, logger)

	// 设置路由
	http.HandleFunc("/user", middleware.SafeHandler(userController.GetUserHandler, errorHandler))
	http.HandleFunc("/user/create", middleware.SafeHandler(userController.CreateUserHandler, errorHandler))
	http.HandleFunc("/error", ErrorDemoHandler)

	// 使用恢复中间件包装所有处理器
	http.HandleFunc("/user/safe", recoveryMiddleware.Handle(http.HandlerFunc(userController.GetUserHandler)).ServeHTTP)
	http.HandleFunc("/user/create/safe", recoveryMiddleware.Handle(http.HandlerFunc(userController.CreateUserHandler)).ServeHTTP)

	// 启动服务器
	port := ":8089"
	log.Printf("错误处理演示服务器启动在端口 %s", port)
	log.Printf("可用端点:")
	log.Printf("  GET /user - 获取用户（正常情况）")
	log.Printf("  GET /user/safe - 获取用户（带恢复中间件）")
	log.Printf("  GET /user/create - 创建用户")
	log.Printf("  GET /user/create/safe - 创建用户（带恢复中间件）")
	log.Printf("  GET /error?type=panic - 演示panic恢复")
	log.Printf("  GET /error?type=timeout - 演示超时错误")
	log.Printf("  GET /error?type=validation - 演示验证错误")
	log.Printf("  GET /error?type=notfound - 演示未找到错误")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
