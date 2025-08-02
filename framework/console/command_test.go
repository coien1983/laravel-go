package console

import (
	"os"
	"testing"
)

func TestNewApplication(t *testing.T) {
	app := NewApplication("test-app", "1.0.0")

	if app.name != "test-app" {
		t.Errorf("Expected name 'test-app', got %s", app.name)
	}

	if app.version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", app.version)
	}

	if app.commands == nil {
		t.Error("Expected commands map to be initialized")
	}

	if app.output == nil {
		t.Error("Expected output to be initialized")
	}
}

func TestAddCommand(t *testing.T) {
	app := NewApplication("test-app", "1.0.0")
	output := NewConsoleOutput()
	generator := NewGenerator(output)

	// 添加命令
	controllerCmd := NewMakeControllerCommand(generator)
	app.AddCommand(controllerCmd)

	// 验证命令已添加
	cmd, exists := app.GetCommand("make:controller")
	if !exists {
		t.Error("Command should exist after adding")
	}

	if cmd.GetName() != "make:controller" {
		t.Errorf("Expected command name 'make:controller', got %s", cmd.GetName())
	}
}

func TestShowHelp(t *testing.T) {
	app := NewApplication("test-app", "1.0.0")
	output := NewConsoleOutput()
	generator := NewGenerator(output)

	// 添加一些命令
	app.AddCommand(NewMakeControllerCommand(generator))
	app.AddCommand(NewMakeModelCommand(generator))

	// 测试帮助命令
	err := app.showHelp()
	if err != nil {
		t.Errorf("showHelp should not return error: %v", err)
	}
}

func TestShowVersion(t *testing.T) {
	app := NewApplication("test-app", "1.0.0")

	err := app.showVersion()
	if err != nil {
		t.Errorf("showVersion should not return error: %v", err)
	}
}

func TestListCommands(t *testing.T) {
	app := NewApplication("test-app", "1.0.0")
	output := NewConsoleOutput()
	generator := NewGenerator(output)

	// 添加一些命令
	app.AddCommand(NewMakeControllerCommand(generator))
	app.AddCommand(NewMakeModelCommand(generator))

	err := app.listCommands()
	if err != nil {
		t.Errorf("listCommands should not return error: %v", err)
	}
}

func TestParseInput(t *testing.T) {
	app := NewApplication("test-app", "1.0.0")
	output := NewConsoleOutput()
	generator := NewGenerator(output)
	cmd := NewMakeControllerCommand(generator)

	// 测试参数解析
	input, err := app.parseInput([]string{"user", "--namespace=app"}, cmd)
	if err != nil {
		t.Errorf("parseInput should not return error: %v", err)
	}

	if input == nil {
		t.Error("Input should not be nil")
	}

	// 验证参数
	name := input.GetArgument("name")
	if name != "user" {
		t.Errorf("Expected argument 'name' to be 'user', got %v", name)
	}

	namespace := input.GetOption("namespace")
	if namespace != "app" {
		t.Errorf("Expected option 'namespace' to be 'app', got %v", namespace)
	}
}

func TestMakeControllerCommand(t *testing.T) {
	output := NewConsoleOutput()
	generator := NewGenerator(output)
	cmd := NewMakeControllerCommand(generator)

	// 验证命令属性
	if cmd.GetName() != "make:controller" {
		t.Errorf("Expected command name 'make:controller', got %s", cmd.GetName())
	}

	if cmd.GetDescription() != "Create a new controller class" {
		t.Errorf("Expected description 'Create a new controller class', got %s", cmd.GetDescription())
	}

	// 验证参数
	args := cmd.GetArguments()
	if len(args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(args))
	}

	if args[0].Name != "name" {
		t.Errorf("Expected argument name 'name', got %s", args[0].Name)
	}

	// 验证选项
	opts := cmd.GetOptions()
	if len(opts) != 1 {
		t.Errorf("Expected 1 option, got %d", len(opts))
	}

	if opts[0].Name != "namespace" {
		t.Errorf("Expected option name 'namespace', got %s", opts[0].Name)
	}
}

func TestMakeModelCommand(t *testing.T) {
	output := NewConsoleOutput()
	generator := NewGenerator(output)
	cmd := NewMakeModelCommand(generator)

	// 验证命令属性
	if cmd.GetName() != "make:model" {
		t.Errorf("Expected command name 'make:model', got %s", cmd.GetName())
	}

	if cmd.GetDescription() != "Create a new model class" {
		t.Errorf("Expected description 'Create a new model class', got %s", cmd.GetDescription())
	}

	// 验证参数
	args := cmd.GetArguments()
	if len(args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(args))
	}

	// 验证选项
	opts := cmd.GetOptions()
	if len(opts) != 1 {
		t.Errorf("Expected 1 option, got %d", len(opts))
	}

	if opts[0].Name != "fields" {
		t.Errorf("Expected option name 'fields', got %s", opts[0].Name)
	}
}

func TestMakeMiddlewareCommand(t *testing.T) {
	output := NewConsoleOutput()
	generator := NewGenerator(output)
	cmd := NewMakeMiddlewareCommand(generator)

	// 验证命令属性
	if cmd.GetName() != "make:middleware" {
		t.Errorf("Expected command name 'make:middleware', got %s", cmd.GetName())
	}

	if cmd.GetDescription() != "Create a new middleware class" {
		t.Errorf("Expected description 'Create a new middleware class', got %s", cmd.GetDescription())
	}

	// 验证参数
	args := cmd.GetArguments()
	if len(args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(args))
	}

	// 验证选项
	opts := cmd.GetOptions()
	if len(opts) != 0 {
		t.Errorf("Expected 0 options, got %d", len(opts))
	}
}

func TestMakeMigrationCommand(t *testing.T) {
	output := NewConsoleOutput()
	generator := NewGenerator(output)
	cmd := NewMakeMigrationCommand(generator)

	// 验证命令属性
	if cmd.GetName() != "make:migration" {
		t.Errorf("Expected command name 'make:migration', got %s", cmd.GetName())
	}

	if cmd.GetDescription() != "Create a new migration file" {
		t.Errorf("Expected description 'Create a new migration file', got %s", cmd.GetDescription())
	}

	// 验证参数
	args := cmd.GetArguments()
	if len(args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(args))
	}

	// 验证选项
	opts := cmd.GetOptions()
	if len(opts) != 2 {
		t.Errorf("Expected 2 options, got %d", len(opts))
	}
}

func TestMakeTestCommand(t *testing.T) {
	output := NewConsoleOutput()
	generator := NewGenerator(output)
	cmd := NewMakeTestCommand(generator)

	// 验证命令属性
	if cmd.GetName() != "make:test" {
		t.Errorf("Expected command name 'make:test', got %s", cmd.GetName())
	}

	if cmd.GetDescription() != "Create a new test class" {
		t.Errorf("Expected description 'Create a new test class', got %s", cmd.GetDescription())
	}

	// 验证参数
	args := cmd.GetArguments()
	if len(args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(args))
	}

	// 验证选项
	opts := cmd.GetOptions()
	if len(opts) != 1 {
		t.Errorf("Expected 1 option, got %d", len(opts))
	}

	if opts[0].Name != "type" {
		t.Errorf("Expected option name 'type', got %s", opts[0].Name)
	}
}

func TestInitCommand(t *testing.T) {
	output := NewConsoleOutput()
	cmd := NewInitCommand(output)

	// 验证命令属性
	if cmd.GetName() != "init" {
		t.Errorf("Expected command name 'init', got %s", cmd.GetName())
	}

	if cmd.GetDescription() != "Initialize a new Laravel-Go project" {
		t.Errorf("Expected description 'Initialize a new Laravel-Go project', got %s", cmd.GetDescription())
	}

	// 验证参数
	args := cmd.GetArguments()
	if len(args) != 0 {
		t.Errorf("Expected 0 arguments, got %d", len(args))
	}

	// 验证选项
	opts := cmd.GetOptions()
	if len(opts) != 1 {
		t.Errorf("Expected 1 option, got %d", len(opts))
	}

	if opts[0].Name != "name" {
		t.Errorf("Expected option name 'name', got %s", opts[0].Name)
	}
}

func TestClearCacheCommand(t *testing.T) {
	output := NewConsoleOutput()
	cmd := NewClearCacheCommand(output)

	// 验证命令属性
	if cmd.GetName() != "cache:clear" {
		t.Errorf("Expected command name 'cache:clear', got %s", cmd.GetName())
	}

	if cmd.GetDescription() != "Clear application cache" {
		t.Errorf("Expected description 'Clear application cache', got %s", cmd.GetDescription())
	}

	// 验证参数
	args := cmd.GetArguments()
	if len(args) != 0 {
		t.Errorf("Expected 0 arguments, got %d", len(args))
	}

	// 验证选项
	opts := cmd.GetOptions()
	if len(opts) != 0 {
		t.Errorf("Expected 0 options, got %d", len(opts))
	}
}

func TestRouteListCommand(t *testing.T) {
	output := NewConsoleOutput()
	cmd := NewRouteListCommand(output)

	// 验证命令属性
	if cmd.GetName() != "route:list" {
		t.Errorf("Expected command name 'route:list', got %s", cmd.GetName())
	}

	if cmd.GetDescription() != "List all registered routes" {
		t.Errorf("Expected description 'List all registered routes', got %s", cmd.GetDescription())
	}

	// 验证参数
	args := cmd.GetArguments()
	if len(args) != 0 {
		t.Errorf("Expected 0 arguments, got %d", len(args))
	}

	// 验证选项
	opts := cmd.GetOptions()
	if len(opts) != 0 {
		t.Errorf("Expected 0 options, got %d", len(opts))
	}
}

func TestGeneratorToPascalCase(t *testing.T) {
	output := NewConsoleOutput()
	generator := NewGenerator(output)

	testCases := []struct {
		input    string
		expected string
	}{
		{"user", "User"},
		{"user_profile", "UserProfile"},
		{"api_user", "ApiUser"},
		{"", ""},
	}

	for _, tc := range testCases {
		result := generator.toPascalCase(tc.input)
		if result != tc.expected {
			t.Errorf("toPascalCase(%s) = %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestConsoleOutput(t *testing.T) {
	output := NewConsoleOutput()

	// 测试基本输出方法
	output.Write("test")
	output.WriteLine("test line")

	// 测试彩色输出方法
	output.Error("error message")
	output.Success("success message")
	output.Warning("warning message")
	output.Info("info message")

	// 测试表格输出
	headers := []string{"Name", "Age"}
	rows := [][]string{
		{"John", "25"},
		{"Jane", "30"},
	}
	output.Table(headers, rows)
}

func TestGeneratorGenerateController(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	output := NewConsoleOutput()
	generator := NewGenerator(output)

	// 生成控制器
	err := generator.GenerateController("user", "app")
	if err != nil {
		t.Errorf("GenerateController should not return error: %v", err)
	}

	// 验证文件是否创建
	filePath := "app/controllers/user_controller.go"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Controller file should be created: %s", filePath)
	}
}

func TestGeneratorGenerateModel(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	output := NewConsoleOutput()
	generator := NewGenerator(output)

	// 生成模型
	fields := []string{"name:string", "email:string", "age:int"}
	err := generator.GenerateModel("user", fields)
	if err != nil {
		t.Errorf("GenerateModel should not return error: %v", err)
	}

	// 验证文件是否创建
	filePath := "app/models/user.go"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Model file should be created: %s", filePath)
	}
}

func TestGeneratorGenerateMiddleware(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	output := NewConsoleOutput()
	generator := NewGenerator(output)

	// 生成中间件
	err := generator.GenerateMiddleware("auth")
	if err != nil {
		t.Errorf("GenerateMiddleware should not return error: %v", err)
	}

	// 验证文件是否创建
	filePath := "app/middleware/auth_middleware.go"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Middleware file should be created: %s", filePath)
	}
}

func TestGeneratorGenerateMigration(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	output := NewConsoleOutput()
	generator := NewGenerator(output)

	// 生成迁移
	fields := []string{"name:string", "email:string:not_null", "age:int:not_null:DEFAULT 0"}
	err := generator.GenerateMigration("create_users_table", "users", fields)
	if err != nil {
		t.Errorf("GenerateMigration should not return error: %v", err)
	}

	// 验证文件是否创建
	files, err := os.ReadDir("database/migrations")
	if err != nil {
		t.Errorf("Failed to read migrations directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("Migration file should be created")
	}
}

func TestGeneratorGenerateTest(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	output := NewConsoleOutput()
	generator := NewGenerator(output)

	// 生成测试
	err := generator.GenerateTest("user", "unit")
	if err != nil {
		t.Errorf("GenerateTest should not return error: %v", err)
	}

	// 验证文件是否创建
	filePath := "tests/user_test.go"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Test file should be created: %s", filePath)
	}
}
