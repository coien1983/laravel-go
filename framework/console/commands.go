package console

import (
	"fmt"
	"os"
	"strings"
)

// MakeControllerCommand 生成控制器命令
type MakeControllerCommand struct {
	generator *Generator
}

// NewMakeControllerCommand 创建新的生成控制器命令
func NewMakeControllerCommand(generator *Generator) *MakeControllerCommand {
	return &MakeControllerCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *MakeControllerCommand) GetName() string {
	return "make:controller"
}

// GetDescription 获取命令描述
func (cmd *MakeControllerCommand) GetDescription() string {
	return "Create a new controller class"
}

// GetSignature 获取命令签名
func (cmd *MakeControllerCommand) GetSignature() string {
	return "make:controller <name> [--namespace=]"
}

// GetArguments 获取命令参数
func (cmd *MakeControllerCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the controller",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *MakeControllerCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "namespace",
			ShortName:   "n",
			Description: "The namespace for the controller",
			Required:    false,
			Default:     "app",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *MakeControllerCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	namespace := input.GetOption("namespace").(string)

	return cmd.generator.GenerateController(name, namespace)
}

// MakeModelCommand 生成模型命令
type MakeModelCommand struct {
	generator *Generator
}

// NewMakeModelCommand 创建新的生成模型命令
func NewMakeModelCommand(generator *Generator) *MakeModelCommand {
	return &MakeModelCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *MakeModelCommand) GetName() string {
	return "make:model"
}

// GetDescription 获取命令描述
func (cmd *MakeModelCommand) GetDescription() string {
	return "Create a new model class"
}

// GetSignature 获取命令签名
func (cmd *MakeModelCommand) GetSignature() string {
	return "make:model <name> [--fields=]"
}

// GetArguments 获取命令参数
func (cmd *MakeModelCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the model",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *MakeModelCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "fields",
			ShortName:   "f",
			Description: "The fields for the model (format: name:type,name:type)",
			Required:    false,
			Default:     "",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *MakeModelCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	fieldsStr := input.GetOption("fields").(string)

	var fields []string
	if fieldsStr != "" {
		fields = strings.Split(fieldsStr, ",")
	}

	return cmd.generator.GenerateModel(name, fields)
}

// MakeMiddlewareCommand 生成中间件命令
type MakeMiddlewareCommand struct {
	generator *Generator
}

// NewMakeMiddlewareCommand 创建新的生成中间件命令
func NewMakeMiddlewareCommand(generator *Generator) *MakeMiddlewareCommand {
	return &MakeMiddlewareCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *MakeMiddlewareCommand) GetName() string {
	return "make:middleware"
}

// GetDescription 获取命令描述
func (cmd *MakeMiddlewareCommand) GetDescription() string {
	return "Create a new middleware class"
}

// GetSignature 获取命令签名
func (cmd *MakeMiddlewareCommand) GetSignature() string {
	return "make:middleware <name>"
}

// GetArguments 获取命令参数
func (cmd *MakeMiddlewareCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the middleware",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *MakeMiddlewareCommand) GetOptions() []Option {
	return []Option{}
}

// Execute 执行命令
func (cmd *MakeMiddlewareCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	return cmd.generator.GenerateMiddleware(name)
}

// MakeMigrationCommand 生成迁移命令
type MakeMigrationCommand struct {
	generator *Generator
}

// NewMakeMigrationCommand 创建新的生成迁移命令
func NewMakeMigrationCommand(generator *Generator) *MakeMigrationCommand {
	return &MakeMigrationCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *MakeMigrationCommand) GetName() string {
	return "make:migration"
}

// GetDescription 获取命令描述
func (cmd *MakeMigrationCommand) GetDescription() string {
	return "Create a new migration file"
}

// GetSignature 获取命令签名
func (cmd *MakeMigrationCommand) GetSignature() string {
	return "make:migration <name> [--table=] [--fields=]"
}

// GetArguments 获取命令参数
func (cmd *MakeMigrationCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the migration",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *MakeMigrationCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "table",
			ShortName:   "t",
			Description: "The table name",
			Required:    false,
			Default:     "",
			Type:        "string",
		},
		{
			Name:        "fields",
			ShortName:   "f",
			Description: "The fields for the table (format: name:type,name:type)",
			Required:    false,
			Default:     "",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *MakeMigrationCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	table := input.GetOption("table").(string)
	fieldsStr := input.GetOption("fields").(string)

	if table == "" {
		table = strings.ToLower(name) + "s"
	}

	var fields []string
	if fieldsStr != "" {
		fields = strings.Split(fieldsStr, ",")
	}

	return cmd.generator.GenerateMigration(name, table, fields)
}

// MakeTestCommand 生成测试命令
type MakeTestCommand struct {
	generator *Generator
}

// NewMakeTestCommand 创建新的生成测试命令
func NewMakeTestCommand(generator *Generator) *MakeTestCommand {
	return &MakeTestCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *MakeTestCommand) GetName() string {
	return "make:test"
}

// GetDescription 获取命令描述
func (cmd *MakeTestCommand) GetDescription() string {
	return "Create a new test class"
}

// GetSignature 获取命令签名
func (cmd *MakeTestCommand) GetSignature() string {
	return "make:test <name> [--type=]"
}

// GetArguments 获取命令参数
func (cmd *MakeTestCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the test",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *MakeTestCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "type",
			ShortName:   "t",
			Description: "The type of test (unit, integration, feature)",
			Required:    false,
			Default:     "unit",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *MakeTestCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	type_ := input.GetOption("type").(string)
	return cmd.generator.GenerateTest(name, type_)
}

// InitCommand 项目初始化命令
type InitCommand struct {
	output Output
}

// NewInitCommand 创建新的项目初始化命令
func NewInitCommand(output Output) *InitCommand {
	return &InitCommand{
		output: output,
	}
}

// GetName 获取命令名称
func (cmd *InitCommand) GetName() string {
	return "init"
}

// GetDescription 获取命令描述
func (cmd *InitCommand) GetDescription() string {
	return "Initialize a new Laravel-Go project"
}

// GetSignature 获取命令签名
func (cmd *InitCommand) GetSignature() string {
	return "init [--name=]"
}

// GetArguments 获取命令参数
func (cmd *InitCommand) GetArguments() []Argument {
	return []Argument{}
}

// GetOptions 获取命令选项
func (cmd *InitCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "name",
			ShortName:   "n",
			Description: "The name of the project",
			Required:    false,
			Default:     "laravel-go-app",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *InitCommand) Execute(input Input) error {
	name := input.GetOption("name").(string)

	// 创建项目目录结构
	dirs := []string{
		"app/controllers",
		"app/models",
		"app/middleware",
		"config",
		"database/migrations",
		"resources/views",
		"routes",
		"storage/cache",
		"storage/logs",
		"tests",
		"public",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 创建基础文件
	files := map[string]string{
		"main.go": `package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Laravel-Go application started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}`,
		"go.mod": fmt.Sprintf(`module %s

go 1.21

require laravel-go/framework v0.1.0

replace laravel-go/framework => ./framework`, name),
		"README.md": fmt.Sprintf(`# %s

A Laravel-Go framework application.

## Getting Started

1. Run the application:
   `+"`"+`bash
   go run main.go
   `+"`"+`

2. Visit http://localhost:8080`, name),
	}

	for fileName, content := range files {
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", fileName, err)
		}
	}

	cmd.output.Success(fmt.Sprintf("Project '%s' initialized successfully!", name))
	return nil
}

// ClearCacheCommand 清理缓存命令
type ClearCacheCommand struct {
	output Output
}

// NewClearCacheCommand 创建新的清理缓存命令
func NewClearCacheCommand(output Output) *ClearCacheCommand {
	return &ClearCacheCommand{
		output: output,
	}
}

// GetName 获取命令名称
func (cmd *ClearCacheCommand) GetName() string {
	return "cache:clear"
}

// GetDescription 获取命令描述
func (cmd *ClearCacheCommand) GetDescription() string {
	return "Clear application cache"
}

// GetSignature 获取命令签名
func (cmd *ClearCacheCommand) GetSignature() string {
	return "cache:clear"
}

// GetArguments 获取命令参数
func (cmd *ClearCacheCommand) GetArguments() []Argument {
	return []Argument{}
}

// GetOptions 获取命令选项
func (cmd *ClearCacheCommand) GetOptions() []Option {
	return []Option{}
}

// Execute 执行命令
func (cmd *ClearCacheCommand) Execute(input Input) error {
	cacheDirs := []string{
		"storage/cache",
		"storage/logs",
	}

	for _, dir := range cacheDirs {
		if err := os.RemoveAll(dir); err != nil {
			cmd.output.Warning(fmt.Sprintf("Failed to remove %s: %v", dir, err))
			continue
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			cmd.output.Warning(fmt.Sprintf("Failed to recreate %s: %v", dir, err))
			continue
		}
		cmd.output.Success(fmt.Sprintf("Cleared cache directory: %s", dir))
	}

	return nil
}

// RouteListCommand 路由列表命令
type RouteListCommand struct {
	output Output
}

// NewRouteListCommand 创建新的路由列表命令
func NewRouteListCommand(output Output) *RouteListCommand {
	return &RouteListCommand{
		output: output,
	}
}

// GetName 获取命令名称
func (cmd *RouteListCommand) GetName() string {
	return "route:list"
}

// GetDescription 获取命令描述
func (cmd *RouteListCommand) GetDescription() string {
	return "List all registered routes"
}

// GetSignature 获取命令签名
func (cmd *RouteListCommand) GetSignature() string {
	return "route:list"
}

// GetArguments 获取命令参数
func (cmd *RouteListCommand) GetArguments() []Argument {
	return []Argument{}
}

// GetOptions 获取命令选项
func (cmd *RouteListCommand) GetOptions() []Option {
	return []Option{}
}

// Execute 执行命令
func (cmd *RouteListCommand) Execute(input Input) error {
	// 这里应该从路由系统中获取路由列表
	// 暂时显示示例数据
	headers := []string{"Method", "URI", "Name", "Action"}
	rows := [][]string{
		{"GET", "/", "home", "HomeController@index"},
		{"GET", "/users", "users.index", "UserController@index"},
		{"POST", "/users", "users.store", "UserController@store"},
		{"GET", "/users/{id}", "users.show", "UserController@show"},
		{"PUT", "/users/{id}", "users.update", "UserController@update"},
		{"DELETE", "/users/{id}", "users.destroy", "UserController@destroy"},
	}

	cmd.output.Table(headers, rows)
	return nil
}

// MakeDockerCommand 生成Docker配置命令
type MakeDockerCommand struct {
	generator *Generator
}

// NewMakeDockerCommand 创建新的生成Docker配置命令
func NewMakeDockerCommand(generator *Generator) *MakeDockerCommand {
	return &MakeDockerCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *MakeDockerCommand) GetName() string {
	return "make:docker"
}

// GetDescription 获取命令描述
func (cmd *MakeDockerCommand) GetDescription() string {
	return "Generate Docker deployment configuration files"
}

// GetSignature 获取命令签名
func (cmd *MakeDockerCommand) GetSignature() string {
	return "make:docker [--name=] [--port=] [--env=]"
}

// GetArguments 获取命令参数
func (cmd *MakeDockerCommand) GetArguments() []Argument {
	return []Argument{}
}

// GetOptions 获取命令选项
func (cmd *MakeDockerCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "name",
			ShortName:   "n",
			Description: "Application name",
			Required:    false,
			Default:     "laravel-go-app",
			Type:        "string",
		},
		{
			Name:        "port",
			ShortName:   "p",
			Description: "Application port",
			Required:    false,
			Default:     "8080",
			Type:        "string",
		},
		{
			Name:        "env",
			ShortName:   "e",
			Description: "Environment (development/production)",
			Required:    false,
			Default:     "development",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *MakeDockerCommand) Execute(input Input) error {
	name := input.GetOption("name").(string)
	port := input.GetOption("port").(string)
	env := input.GetOption("env").(string)

	return cmd.generator.GenerateDockerConfig(name, port, env)
}

// MakeK8sCommand 生成Kubernetes配置命令
type MakeK8sCommand struct {
	generator *Generator
}

// NewMakeK8sCommand 创建新的生成Kubernetes配置命令
func NewMakeK8sCommand(generator *Generator) *MakeK8sCommand {
	return &MakeK8sCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *MakeK8sCommand) GetName() string {
	return "make:k8s"
}

// GetDescription 获取命令描述
func (cmd *MakeK8sCommand) GetDescription() string {
	return "Generate Kubernetes deployment configuration files"
}

// GetSignature 获取命令签名
func (cmd *MakeK8sCommand) GetSignature() string {
	return "make:k8s [--name=] [--replicas=] [--port=] [--namespace=]"
}

// GetArguments 获取命令参数
func (cmd *MakeK8sCommand) GetArguments() []Argument {
	return []Argument{}
}

// GetOptions 获取命令选项
func (cmd *MakeK8sCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "name",
			ShortName:   "n",
			Description: "Application name",
			Required:    false,
			Default:     "laravel-go-app",
			Type:        "string",
		},
		{
			Name:        "replicas",
			ShortName:   "r",
			Description: "Number of replicas",
			Required:    false,
			Default:     "3",
			Type:        "string",
		},
		{
			Name:        "port",
			ShortName:   "p",
			Description: "Application port",
			Required:    false,
			Default:     "8080",
			Type:        "string",
		},
		{
			Name:        "namespace",
			ShortName:   "ns",
			Description: "Kubernetes namespace",
			Required:    false,
			Default:     "default",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *MakeK8sCommand) Execute(input Input) error {
	name := input.GetOption("name").(string)
	replicas := input.GetOption("replicas").(string)
	port := input.GetOption("port").(string)
	namespace := input.GetOption("namespace").(string)

	return cmd.generator.GenerateK8sConfig(name, replicas, port, namespace)
}
