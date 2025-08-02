package main

import (
	"fmt"
	"os"

	"laravel-go/framework/config"
	"laravel-go/framework/container"
	"laravel-go/framework/errors"
)

// UserService 用户服务接口
type UserService interface {
	GetName() string
	GetEmail() string
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	Name  string
	Email string
}

func (u *UserServiceImpl) GetName() string {
	return u.Name
}

func (u *UserServiceImpl) GetEmail() string {
	return u.Email
}

// UserController 用户控制器
type UserController struct {
	UserService UserService `inject:"service"`
}

func (c *UserController) Show() string {
	return fmt.Sprintf("User: %s (%s)", c.UserService.GetName(), c.UserService.GetEmail())
}

func main() {
	fmt.Println("=== Laravel-Go Framework Core Demo ===")

	// 1. 演示容器功能
	demoContainer()

	// 2. 演示配置管理
	demoConfig()

	// 3. 演示错误处理
	demoErrorHandling()

	fmt.Println("\n=== Demo Completed ===")
}

func demoContainer() {
	fmt.Println("1. 容器功能演示:")
	fmt.Println("----------------")

	// 创建容器
	app := container.NewContainer()

	// 注册服务
	userService := &UserServiceImpl{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	app.Bind((*UserService)(nil), userService)

	// 解析服务
	resolvedService := app.Make((*UserService)(nil))
	if userSvc, ok := resolvedService.(UserService); ok {
		fmt.Printf("  解析的服务: %s (%s)\n", userSvc.GetName(), userSvc.GetEmail())
	}

	// 使用依赖注入
	controller := &UserController{}
	// 手动注入依赖
	controller.UserService = app.Make((*UserService)(nil)).(UserService)

	if controller.UserService != nil {
		fmt.Printf("  控制器注入结果: %s\n", controller.Show())
	}

	// 使用回调函数
	app.BindCallback((*UserService)(nil), func(c container.Container) interface{} {
		return &UserServiceImpl{
			Name:  "Callback User",
			Email: "callback@example.com",
		}
	})

	callbackService := app.Make((*UserService)(nil))
	if userSvc, ok := callbackService.(UserService); ok {
		fmt.Printf("  回调函数结果: %s (%s)\n", userSvc.GetName(), userSvc.GetEmail())
	}

	fmt.Println()
}

func demoConfig() {
	fmt.Println("2. 配置管理演示:")
	fmt.Println("----------------")

	// 创建配置管理器
	cfg := config.NewConfig()

	// 设置配置
	cfg.Set("app.name", "Laravel-Go Demo")
	cfg.Set("app.debug", true)
	cfg.Set("database.host", "localhost")
	cfg.Set("database.port", 3306)
	cfg.Set("app.middleware", []string{"auth", "cors", "log"})

	// 获取配置
	fmt.Printf("  应用名称: %s\n", cfg.GetString("app.name"))
	fmt.Printf("  调试模式: %t\n", cfg.GetBool("app.debug"))
	fmt.Printf("  数据库主机: %s\n", cfg.GetString("database.host"))
	fmt.Printf("  数据库端口: %d\n", cfg.GetInt("database.port"))

	// 获取带默认值的配置
	fmt.Printf("  不存在的配置(带默认值): %s\n", cfg.GetString("app.version", "1.0.0"))

	// 获取嵌套配置
	fmt.Printf("  嵌套配置: %s\n", cfg.GetString("database.host"))

	// 获取切片配置
	middleware := cfg.GetStringSlice("app.middleware")
	fmt.Printf("  中间件: %v\n", middleware)

	// 从结构体加载配置
	type AppConfig struct {
		Name  string `env:"APP_NAME" default:"Laravel-Go"`
		Debug bool   `env:"APP_DEBUG" default:"false"`
		Port  int    `env:"APP_PORT" default:"8080"`
	}

	// 设置环境变量
	os.Setenv("APP_NAME", "Demo App")
	os.Setenv("APP_DEBUG", "true")

	appConfig := &AppConfig{}
	cfg.LoadFromStruct(appConfig)

	fmt.Printf("  从结构体加载 - 名称: %s\n", cfg.GetString("Name"))
	fmt.Printf("  从结构体加载 - 调试: %t\n", cfg.GetBool("Debug"))
	fmt.Printf("  从结构体加载 - 端口: %d\n", cfg.GetInt("Port"))

	// 验证配置
	rules := map[string]string{
		"app.name":      "required|string",
		"app.debug":     "required|bool",
		"database.host": "required|string",
	}

	if err := cfg.Validate(rules); err != nil {
		fmt.Printf("  配置验证失败: %v\n", err)
	} else {
		fmt.Println("  配置验证通过")
	}

	fmt.Println()
}

func demoErrorHandling() {
	fmt.Println("3. 错误处理演示:")
	fmt.Println("----------------")

	// 创建错误
	appErr := errors.New("这是一个应用错误")
	fmt.Printf("  基础错误: %s (代码: %d)\n", appErr.Message, appErr.Code)

	// 错误链式调用
	chainedErr := errors.New("基础错误").
		WithCode(400).
		WithMessage("自定义错误消息").
		WithStack()
	fmt.Printf("  链式错误: %s (代码: %d)\n", chainedErr.Message, chainedErr.Code)
	fmt.Printf("  堆栈信息长度: %d\n", len(chainedErr.Stack))

	// 包装错误
	originalErr := fmt.Errorf("原始错误")
	wrappedErr := errors.Wrap(originalErr, "包装后的错误")
	fmt.Printf("  包装错误: %s\n", wrappedErr.Message)
	fmt.Printf("  原始错误: %v\n", wrappedErr.Err)

	// 预定义错误
	fmt.Printf("  404错误: %s (代码: %d)\n", errors.ErrNotFound.Message, errors.ErrNotFound.Code)
	fmt.Printf("  401错误: %s (代码: %d)\n", errors.ErrUnauthorized.Message, errors.ErrUnauthorized.Code)
	fmt.Printf("  500错误: %s (代码: %d)\n", errors.ErrInternalServer.Message, errors.ErrInternalServer.Code)

	// 验证错误
	validationErrors := errors.ValidationErrors{}
	validationErrors.Add("email", "邮箱格式不正确", "invalid-email")
	validationErrors.Add("password", "密码长度不足", "123")

	if validationErrors.HasErrors() {
		fmt.Printf("  验证错误: %s\n", validationErrors.Error())
		fmt.Printf("  错误数量: %d\n", len(validationErrors.GetErrors()))

		emailErrors := validationErrors.GetErrorsByField("email")
		fmt.Printf("  邮箱错误: %s\n", emailErrors[0].Message)
	}

	// 错误处理器
	logger := &MockLogger{}
	handler := errors.NewDefaultErrorHandler(logger)

	// 处理错误
	handledErr := handler.Handle(appErr)
	fmt.Printf("  处理后的错误: %s\n", handledErr.Error())
	fmt.Printf("  日志记录数量: %d\n", len(logger.ErrorLogs))

	fmt.Println()
}

// MockLogger 模拟日志记录器
type MockLogger struct {
	ErrorLogs   []map[string]interface{}
	WarningLogs []map[string]interface{}
	InfoLogs    []map[string]interface{}
	DebugLogs   []map[string]interface{}
}

func (m *MockLogger) Error(message string, context map[string]interface{}) {
	m.ErrorLogs = append(m.ErrorLogs, context)
}

func (m *MockLogger) Warning(message string, context map[string]interface{}) {
	m.WarningLogs = append(m.WarningLogs, context)
}

func (m *MockLogger) Info(message string, context map[string]interface{}) {
	m.InfoLogs = append(m.InfoLogs, context)
}

func (m *MockLogger) Debug(message string, context map[string]interface{}) {
	m.DebugLogs = append(m.DebugLogs, context)
}
