package console

import (
	"fmt"
	"os"
	"path/filepath"
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
	return "init [project-name] [--name=]"
}

// GetArguments 获取命令参数
func (cmd *InitCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "project-name",
			Description: "The name of the project (optional)",
			Required:    false,
		},
	}
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
	// 获取项目名称，优先使用参数，其次使用选项
	var projectName string
	if arg := input.GetArgument("project-name"); arg != nil {
		projectName = arg.(string)
	} else {
		projectName = input.GetOption("name").(string)
	}

	// 交互式配置
	config := InteractiveConfig(projectName, cmd.output)

	// 显示配置信息
	cmd.output.Info(fmt.Sprintf("正在使用配置创建项目: %s", config.Name))

	// 确定项目目录
	var projectDir string
	if projectName != "" {
		// 如果提供了项目名称，创建新目录
		projectDir = projectName
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			return fmt.Errorf("failed to create project directory %s: %w", projectDir, err)
		}
		cmd.output.Success(fmt.Sprintf("Created project directory: %s", projectDir))
	} else {
		// 如果没有提供项目名称，使用默认名称创建目录
		projectDir = "laravel-go-app"
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			return fmt.Errorf("failed to create project directory %s: %w", projectDir, err)
		}
		cmd.output.Success(fmt.Sprintf("Created project directory: %s", projectDir))
	}

	// 创建项目目录结构 (Laravel标准结构)
	dirs := []string{
		"app/Console",
		"app/Exceptions",
		"app/Http/Controllers",
		"app/Http/Middleware",
		"app/Http/Requests",
		"app/Models",
		"app/Providers",
		"app/Services",
		"bootstrap",
		"config",
		"database/factories",
		"database/migrations",
		"database/seeders",
		"public",
		"resources/css",
		"resources/js",
		"resources/lang",
		"resources/views",
		"routes",
		"storage/app",
		"storage/framework/cache",
		"storage/framework/sessions",
		"storage/framework/views",
		"storage/logs",
		"tests",
	}

	// 如果是微服务架构，添加gRPC和API网关相关目录
	if config.Architecture == "microservice" {
		// gRPC相关目录
		if config.GRPC != "none" {
			grpcDirs := []string{
				"proto",
				"grpc/server",
				"grpc/client",
				"grpc/interceptors",
				"grpc/services",
			}
			dirs = append(dirs, grpcDirs...)
		}

		// API网关相关目录
		if config.APIGateway != "none" {
			gatewayDirs := []string{
				"gateway",
				"gateway/middleware",
				"gateway/routes",
				"gateway/plugins",
			}
			dirs = append(dirs, gatewayDirs...)
		}
	}

	// 如果有项目目录，在项目目录下创建结构
	if projectDir != "" {
		for i, dir := range dirs {
			dirs[i] = filepath.Join(projectDir, dir)
		}
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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 设置服务器
	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	// 创建 HTTP 服务器
	mux := http.NewServeMux()
	
	// 注册路由
	registerRoutes(mux)
	
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// 启动服务器
	go func() {
		fmt.Printf("🚀 Server starting on http://localhost%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	fmt.Println("\n🛑 Shutting down server...")
	fmt.Println("✅ Server stopped gracefully")
}

// registerRoutes 注册路由
func registerRoutes(mux *http.ServeMux) {
	// 导入路由包
	// 这里会在运行时动态加载路由
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message": "Welcome to Laravel-Go!",
			"version": "1.0.0",
			"status":  "running",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status": "ok",
			"time":   "2024-01-01T00:00:00Z",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}`,
		"go.mod": fmt.Sprintf(`module %s

go 1.21

require (
	github.com/coien1983/laravel-go/framework v0.1.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
)

replace github.com/coien1983/laravel-go/framework => ./framework`, projectName),
		".env": `# Application Configuration
APP_NAME=Laravel-Go App
APP_ENV=local
APP_DEBUG=true
APP_URL=http://localhost:8080

# Server Configuration
PORT=8080

# Database Configuration
DB_CONNECTION=sqlite
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=app.db
DB_USERNAME=
DB_PASSWORD=

# Cache Configuration
CACHE_DRIVER=memory

# Session Configuration
SESSION_DRIVER=memory
SESSION_LIFETIME=120

# Logging Configuration
LOG_LEVEL=debug
LOG_CHANNEL=stack`,
		".env.example": getEnvExampleTemplate(),
		"config/app.go": `package config

import (
	"os"
	"strconv"
)

// AppConfig 应用配置
type AppConfig struct {
	Name   string
	Env    string
	Debug  bool
	URL    string
	Port   string
}

// LoadAppConfig 加载应用配置
func LoadAppConfig() *AppConfig {
	debug, _ := strconv.ParseBool(getEnv("APP_DEBUG", "true"))
	
	return &AppConfig{
		Name:  getEnv("APP_NAME", "Laravel-Go App"),
		Env:   getEnv("APP_ENV", "local"),
		Debug: debug,
		URL:   getEnv("APP_URL", "http://localhost:8080"),
		Port:  getEnv("PORT", "8080"),
	}
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}`,
		"config/database.go": `package config

import (
	"os"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Connection string
	Host       string
	Port       string
	Database   string
	Username   string
	Password   string
}

// LoadDatabaseConfig 加载数据库配置
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Connection: getEnv("DB_CONNECTION", "sqlite"),
		Host:       getEnv("DB_HOST", "127.0.0.1"),
		Port:       getEnv("DB_PORT", "3306"),
		Database:   getEnv("DB_DATABASE", "app.db"),
		Username:   getEnv("DB_USERNAME", ""),
		Password:   getEnv("DB_PASSWORD", ""),
	}
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}`,
		"app/Http/Controllers/HomeController.go": `package controllers

import (
	"net/http"
	"encoding/json"
)

// HomeController 首页控制器
type HomeController struct{}

// NewHomeController 创建新的首页控制器
func NewHomeController() *HomeController {
	return &HomeController{}
}

// Index 首页
func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Welcome to Laravel-Go!",
		"version": "1.0.0",
		"status":  "running",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Health 健康检查
func (c *HomeController) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "ok",
		"time":   "2024-01-01T00:00:00Z",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}`,
		"app/Http/Controllers/UserController.go": `package controllers

import (
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/gorilla/mux"
)

// User 用户模型
type User struct {
	ID    int    ` + "`json:\"id\"`" + `
	Name  string ` + "`json:\"name\"`" + `
	Email string ` + "`json:\"email\"`" + `
}

// UserController 用户控制器
type UserController struct {
	users []User
}

// NewUserController 创建新的用户控制器
func NewUserController() *UserController {
	// 初始化一些示例数据
	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	
	return &UserController{
		users: users,
	}
}

// Index 获取用户列表
func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c.users)
}

// Show 获取单个用户
func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	for _, user := range c.users {
		if user.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	
	http.Error(w, "User not found", http.StatusNotFound)
}

// Store 创建用户
func (c *UserController) Store(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// 简单的 ID 生成
	user.ID = len(c.users) + 1
	c.users = append(c.users, user)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}`,
		"app/models/user.go": `package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      ` + "`json:\"id\"`" + `
	Name      string    ` + "`json:\"name\"`" + `
	Email     string    ` + "`json:\"email\"`" + `
	Password  string    ` + "`json:\"-\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// TableName 获取表名
func (u *User) TableName() string {
	return "users"
}

// NewUser 创建新用户
func NewUser() *User {
	return &User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Fillable 可填充字段
func (u *User) Fillable() []string {
	return []string{"name", "email", "password"}
}

// Hidden 隐藏字段
func (u *User) Hidden() []string {
	return []string{"password"}
}`,
		"routes/web.go": `package routes

import (
	"net/http"
	"github.com/gorilla/mux"
			"` + projectName + `/app/Http/Controllers"
)

// RegisterWebRoutes 注册 Web 路由
func RegisterWebRoutes(router *mux.Router) {
	// 首页路由
	router.HandleFunc("/", controllers.NewHomeController().Index).Methods("GET")
	router.HandleFunc("/health", controllers.NewHomeController().Health).Methods("GET")
	
	// API 路由
	api := router.PathPrefix("/api").Subrouter()
	
	// 用户路由
	userController := controllers.NewUserController()
	api.HandleFunc("/users", userController.Index).Methods("GET")
	api.HandleFunc("/users", userController.Store).Methods("POST")
	api.HandleFunc("/users/{id}", userController.Show).Methods("GET")
	
	// 静态文件
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
}`,
		"database/migrations/001_create_users_table.sql": `-- Migration: Create Users Table
-- Description: 创建用户表
-- Version: 1.0

-- UP Migration
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入示例数据
INSERT INTO users (name, email, password) VALUES 
('John Doe', 'john@example.com', 'hashed_password_1'),
('Jane Smith', 'jane@example.com', 'hashed_password_2');

-- DOWN Migration (如果需要回滚)
-- DROP TABLE IF EXISTS users;`,
		"README.md": fmt.Sprintf(`# %s

一个基于 Laravel-Go Framework 构建的完整 Web 应用。

## 快速开始

1. 安装依赖: go mod tidy
2. 运行应用: go run main.go
3. 访问: http://localhost:8080

## 项目结构

		- app/Http/Controllers/ - 控制器
- app/models/ - 数据模型
- config/ - 配置文件
- database/ - 数据库相关
- routes/ - 路由定义
- storage/ - 存储目录

## API 接口

- GET / - 首页
- GET /health - 健康检查
- GET /api/users - 获取用户列表
- POST /api/users - 创建用户

## 开发

使用 largo 命令生成代码:
- largo make:controller ProductController
- largo make:model Product
- largo make:middleware AuthMiddleware

更多信息请参考 Laravel-Go Framework 文档`, projectName),
		".gitignore": `# 编译输出
*.exe
*.exe~
*.dll
*.so
*.dylib

# 测试二进制文件
*.test

# 覆盖率文件
*.out

# 依赖目录
vendor/

# IDE 文件
.vscode/
.idea/
*.swp
*.swo

# 环境变量文件
.env

# 日志文件
storage/logs/*.log

# 缓存文件
storage/framework/cache/*

# 会话文件
storage/framework/sessions/*

# 数据库文件
*.db
*.sqlite

# 临时文件
*.tmp
*.temp

# 系统文件
.DS_Store
Thumbs.db

# 上传文件
storage/uploads/*

# 备份文件
*.backup
*.bak`,
		"Makefile": cmd.generateMakefile(config),
	}

	for fileName, content := range files {
		// 如果有项目目录，在项目目录下创建文件
		if projectDir != "" {
			fileName = filepath.Join(projectDir, fileName)
		}
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", fileName, err)
		}
	}

	// 如果是微服务架构，生成gRPC和API网关相关文件
	if config.Architecture == "microservice" {
		if err := cmd.generateMicroserviceFiles(config, projectDir); err != nil {
			return fmt.Errorf("failed to generate microservice files: %w", err)
		}
	}

	// 根据配置生成 Docker 和 Kubernetes 文件
	if err := cmd.GenerateDeploymentFiles(config, projectDir); err != nil {
		return fmt.Errorf("failed to generate deployment files: %w", err)
	}

	cmd.output.Success(fmt.Sprintf("Project '%s' initialized successfully!", projectName))
	return nil
}

// GenerateDeploymentFiles 根据配置生成部署文件
func (cmd *InitCommand) GenerateDeploymentFiles(config *ProjectConfig, projectDir string) error {
	// 生成 Docker 文件
	if config.Docker != "none" {
		if err := cmd.generateDockerFiles(config, projectDir); err != nil {
			return fmt.Errorf("failed to generate docker files: %w", err)
		}
		cmd.output.Success("✅ Docker 配置文件已生成")
	}

	// 生成 Kubernetes 文件
	if config.Kubernetes != "none" {
		if err := cmd.generateK8sFiles(config, projectDir); err != nil {
			return fmt.Errorf("failed to generate kubernetes files: %w", err)
		}
		cmd.output.Success("✅ Kubernetes 配置文件已生成")
	}

	return nil
}

// generateDockerFiles 生成 Docker 相关文件
func (cmd *InitCommand) generateDockerFiles(config *ProjectConfig, projectDir string) error {
	dockerFiles := map[string]string{
		"Dockerfile": `# {{ .Name }} Dockerfile
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE {{ .Port }}

# 设置环境变量
ENV APP_ENV={{ .Env }}
ENV APP_DEBUG={{ if eq .Env "development" }}true{{ else }}false{{ end }}

# 运行应用
CMD ["./main"]`,
		".dockerignore": `# Git
.git
.gitignore

# IDE
.vscode
.idea
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log
logs/

# Cache
cache/
tmp/

# Test
test.db
*.test

# Documentation
README.md
docs/

# Build artifacts
main
*.exe`,
	}

	// 如果是完整配置，添加 docker-compose.yml
	if config.Docker == "full" {
		dockerFiles["docker-compose.yml"] = `version: '3.8'

services:
  {{ .Name }}:
    build: .
    ports:
      - "{{ .Port }}:{{ .Port }}"
    environment:
      - APP_ENV={{ .Env }}
      - APP_DEBUG={{ if eq .Env "development" }}true{{ else }}false{{ end }}
    depends_on:
      - redis
      - postgres
    networks:
      - {{ .Name }}-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - {{ .Name }}-network

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: {{ .Name }}
      POSTGRES_USER: {{ .Name }}
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - {{ .Name }}-network

volumes:
  redis_data:
  postgres_data:

networks:
  {{ .Name }}-network:
    driver: bridge`
	}

	// 创建文件
	for fileName, content := range dockerFiles {
		if projectDir != "" {
			fileName = filepath.Join(projectDir, fileName)
		}
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create docker file %s: %w", fileName, err)
		}
	}

	return nil
}

// generateK8sFiles 生成 Kubernetes 相关文件
func (cmd *InitCommand) generateK8sFiles(config *ProjectConfig, projectDir string) error {
	// 创建 k8s 目录
	k8sDir := "k8s"
	if projectDir != "" {
		k8sDir = filepath.Join(projectDir, k8sDir)
	}
	if err := os.MkdirAll(k8sDir, 0755); err != nil {
		return fmt.Errorf("failed to create k8s directory: %w", err)
	}

	k8sFiles := map[string]string{
		"deployment.yaml": `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    app: {{ .Name }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      app: {{ .Name }}
  template:
    metadata:
      labels:
        app: {{ .Name }}
    spec:
      containers:
      - name: {{ .Name }}
        image: {{ .Name }}:latest
        ports:
        - containerPort: {{ .Port }}
        env:
        - name: APP_ENV
          value: "production"
        - name: APP_DEBUG
          value: "false"
        livenessProbe:
          httpGet:
            path: /health
            port: {{ .Port }}
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: {{ .Port }}
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"`,
		"service.yaml": `apiVersion: v1
kind: Service
metadata:
  name: {{ .Name }}-service
  namespace: {{ .Namespace }}
spec:
  selector:
    app: {{ .Name }}
  ports:
    - protocol: TCP
      port: 80
      targetPort: {{ .Port }}
  type: ClusterIP`,
		"ingress.yaml": `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Name }}-ingress
  namespace: {{ .Namespace }}
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: {{ .Name }}.local
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: {{ .Name }}-service
            port:
              number: 80`,
		"configmap.yaml": `apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Name }}-config
  namespace: {{ .Namespace }}
data:
  APP_ENV: "production"
  APP_DEBUG: "false"
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_NAME: "{{ .Name }}"
  DB_USER: "{{ .Name }}"
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"`,
	}

	// 如果是完整配置，添加监控配置
	if config.Kubernetes == "full" {
		k8sFiles["monitoring.yaml"] = `apiVersion: v1
kind: Service
metadata:
  name: {{ .Name }}-monitoring
  namespace: {{ .Namespace }}
spec:
  selector:
    app: {{ .Name }}
  ports:
    - name: metrics
      port: 9090
      targetPort: 9090
  type: ClusterIP

---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: {{ .Name }}-monitor
  namespace: {{ .Namespace }}
spec:
  selector:
    matchLabels:
      app: {{ .Name }}
  endpoints:
    - port: metrics
      interval: 30s`
	}

	// 创建文件
	for fileName, content := range k8sFiles {
		filePath := filepath.Join(k8sDir, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create k8s file %s: %w", filePath, err)
		}
	}

	return nil
}

// generateMakefile 生成项目的 Makefile
func (cmd *InitCommand) generateMakefile(config *ProjectConfig) string {
	makefile := `# Laravel-Go Project Makefile

.PHONY: help
help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## 运行应用
	go run main.go

.PHONY: build
build: ## 构建应用
	go build -o bin/app main.go

.PHONY: test
test: ## 运行测试
	go test ./...

.PHONY: clean
clean: ## 清理构建文件
	rm -rf bin/
	rm -f *.db

.PHONY: deps
deps: ## 安装依赖
	go mod tidy
	go mod download

.PHONY: dev
dev: ## 开发模式运行
	APP_ENV=local go run main.go

.PHONY: prod
prod: ## 生产模式运行
	APP_ENV=production go run main.go`

	// 添加 Docker 相关命令
	if config.Docker != "none" {
		makefile += `

# =============================================================================
# Docker 操作
# =============================================================================

.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	docker build -t ` + config.Name + ` .

.PHONY: docker-run
docker-run: ## 运行 Docker 容器
	docker run -p 8080:8080 ` + config.Name + `

.PHONY: docker-compose-up
docker-compose-up: ## 启动 Docker Compose 服务
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## 停止 Docker Compose 服务
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## 查看 Docker Compose 日志
	docker-compose logs -f

.PHONY: docker-clean
docker-clean: ## 清理 Docker 资源
	docker-compose down -v --remove-orphans
	docker system prune -f`
	}

	// 添加 Kubernetes 相关命令
	if config.Kubernetes != "none" {
		makefile += `

# =============================================================================
# Kubernetes 操作
# =============================================================================

.PHONY: k8s-apply
k8s-apply: ## 部署到 Kubernetes
	kubectl apply -f k8s/

.PHONY: k8s-delete
k8s-delete: ## 从 Kubernetes 删除
	kubectl delete -f k8s/

.PHONY: k8s-status
k8s-status: ## 查看 Kubernetes 部署状态
	kubectl get pods,services,ingress

.PHONY: k8s-logs
k8s-logs: ## 查看 Kubernetes 日志
	kubectl logs -f deployment/` + config.Name + `

.PHONY: k8s-clean
k8s-clean: ## 清理 Kubernetes 资源
	kubectl delete -f k8s/ --ignore-not-found=true`
	}

	return makefile
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

// Docker 和 Kubernetes 命令已移除

// MakeApiCommand 快速生成 API 组件命令
type MakeApiCommand struct {
	generator *Generator
}

// NewMakeApiCommand 创建新的快速生成 API 组件命令
func NewMakeApiCommand(generator *Generator) *MakeApiCommand {
	return &MakeApiCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *MakeApiCommand) GetName() string {
	return "make:api"
}

// GetDescription 获取命令描述
func (cmd *MakeApiCommand) GetDescription() string {
	return "Quickly generate API controller and model"
}

// GetSignature 获取命令签名
func (cmd *MakeApiCommand) GetSignature() string {
	return "make:api <name> [--fields=]"
}

// GetArguments 获取命令参数
func (cmd *MakeApiCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the resource",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *MakeApiCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "fields",
			ShortName:   "f",
			Description: "The fields for the model (format: name:string,email:string,age:int)",
			Required:    false,
			Default:     "",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *MakeApiCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	fields := input.GetOption("fields").(string)

	cmd.generator.output.Info(fmt.Sprintf("生成 %s 的 API 组件...", name))

	// 生成控制器
	if err := cmd.generator.GenerateController(name, "api"); err != nil {
		return err
	}

	// 生成模型
	fieldList := []string{}
	if fields != "" {
		fieldList = strings.Split(fields, ",")
	}
	if err := cmd.generator.GenerateModel(name, fieldList); err != nil {
		return err
	}

	// 生成迁移
	migrationName := fmt.Sprintf("create_%ss_table", name)
	tableName := fmt.Sprintf("%ss", name)
	if err := cmd.generator.GenerateMigration(migrationName, tableName, fieldList); err != nil {
		return err
	}

	cmd.generator.output.Success(fmt.Sprintf("✅ %s API 组件生成完成!", name))
	return nil
}

// MakeCrudCommand 快速生成 CRUD 组件命令
type MakeCrudCommand struct {
	generator *Generator
}

// NewMakeCrudCommand 创建新的快速生成 CRUD 组件命令
func NewMakeCrudCommand(generator *Generator) *MakeCrudCommand {
	return &MakeCrudCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *MakeCrudCommand) GetName() string {
	return "make:crud"
}

// GetDescription 获取命令描述
func (cmd *MakeCrudCommand) GetDescription() string {
	return "Quickly generate complete CRUD components"
}

// GetSignature 获取命令签名
func (cmd *MakeCrudCommand) GetSignature() string {
	return "make:crud <name> [--fields=]"
}

// GetArguments 获取命令参数
func (cmd *MakeCrudCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the resource",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *MakeCrudCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "fields",
			ShortName:   "f",
			Description: "The fields for the model (format: name:string,email:string,age:int)",
			Required:    false,
			Default:     "",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *MakeCrudCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	fields := input.GetOption("fields").(string)

	cmd.generator.output.Info(fmt.Sprintf("生成 %s 的完整 CRUD 组件...", name))

	// 生成控制器
	if err := cmd.generator.GenerateController(name, "app"); err != nil {
		return err
	}

	// 生成模型
	fieldList := []string{}
	if fields != "" {
		fieldList = strings.Split(fields, ",")
	}
	if err := cmd.generator.GenerateModel(name, fieldList); err != nil {
		return err
	}

	// 生成迁移
	migrationName := fmt.Sprintf("create_%ss_table", name)
	tableName := fmt.Sprintf("%ss", name)
	if err := cmd.generator.GenerateMigration(migrationName, tableName, fieldList); err != nil {
		return err
	}

	// 生成单元测试
	if err := cmd.generator.GenerateTest(name, "unit"); err != nil {
		return err
	}

	// 生成集成测试
	if err := cmd.generator.GenerateTest(name, "integration"); err != nil {
		return err
	}

	cmd.generator.output.Success(fmt.Sprintf("✅ %s CRUD 组件生成完成!", name))
	return nil
}

// ProjectInfoCommand 项目信息命令
type ProjectInfoCommand struct {
	output Output
}

// NewProjectInfoCommand 创建新的项目信息命令
func NewProjectInfoCommand(output Output) *ProjectInfoCommand {
	return &ProjectInfoCommand{
		output: output,
	}
}

// GetName 获取命令名称
func (cmd *ProjectInfoCommand) GetName() string {
	return "project:info"
}

// GetDescription 获取命令描述
func (cmd *ProjectInfoCommand) GetDescription() string {
	return "Show project information"
}

// GetSignature 获取命令签名
func (cmd *ProjectInfoCommand) GetSignature() string {
	return "project:info"
}

// GetArguments 获取命令参数
func (cmd *ProjectInfoCommand) GetArguments() []Argument {
	return []Argument{}
}

// GetOptions 获取命令选项
func (cmd *ProjectInfoCommand) GetOptions() []Option {
	return []Option{}
}

// Execute 执行命令
func (cmd *ProjectInfoCommand) Execute(input Input) error {
	cmd.output.Info("Laravel-Go Framework 项目信息:")
	cmd.output.Info("  应用名称: laravel-go-app")
	cmd.output.Info("  默认端口: 8080")
	cmd.output.Info("  默认命名空间: default")
	cmd.output.Info("  默认副本数: 3")
	cmd.output.Info("")
	cmd.output.Info("可用命令:")
	cmd.output.Info("  largo list          - 显示所有命令")
	cmd.output.Info("  largo init          - 初始化项目")
	cmd.output.Info("  largo make:controller - 生成控制器")
	cmd.output.Info("  largo make:model    - 生成模型")
	cmd.output.Info("  # Docker 和 Kubernetes 支持已移除")
	cmd.output.Info("  largo make:api      - 快速生成 API 组件")
	cmd.output.Info("  largo make:crud     - 快速生成 CRUD 组件")
	return nil
}

// VersionCommand 版本信息命令
type VersionCommand struct {
	output Output
}

// NewVersionCommand 创建新的版本信息命令
func NewVersionCommand(output Output) *VersionCommand {
	return &VersionCommand{
		output: output,
	}
}

// GetName 获取命令名称
func (cmd *VersionCommand) GetName() string {
	return "version"
}

// GetDescription 获取命令描述
func (cmd *VersionCommand) GetDescription() string {
	return "Show version information"
}

// GetSignature 获取命令签名
func (cmd *VersionCommand) GetSignature() string {
	return "version"
}

// GetArguments 获取命令参数
func (cmd *VersionCommand) GetArguments() []Argument {
	return []Argument{}
}

// GetOptions 获取命令选项
func (cmd *VersionCommand) GetOptions() []Option {
	return []Option{}
}

// Execute 执行命令
func (cmd *VersionCommand) Execute(input Input) error {
	cmd.output.Info("Laravel-Go Framework v1.0.0")
	cmd.output.Info("A modern Go web framework inspired by Laravel")
	cmd.output.Info("GitHub: https://github.com/coien1983/laravel-go")
	return nil
}

// AddModuleCommand 添加模块命令
type AddModuleCommand struct {
	generator *Generator
}

// NewAddModuleCommand 创建新的添加模块命令
func NewAddModuleCommand(generator *Generator) *AddModuleCommand {
	return &AddModuleCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *AddModuleCommand) GetName() string {
	return "add:module"
}

// GetDescription 获取命令描述
func (cmd *AddModuleCommand) GetDescription() string {
	return "Add a new module with controller, model, service, and repository"
}

// GetSignature 获取命令签名
func (cmd *AddModuleCommand) GetSignature() string {
	return "add:module <name> [--api] [--web] [--full]"
}

// GetArguments 获取命令参数
func (cmd *AddModuleCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the module",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *AddModuleCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "api",
			ShortName:   "a",
			Description: "Generate API controller and routes",
			Required:    false,
			Default:     true,
			Type:        "bool",
		},
		{
			Name:        "web",
			ShortName:   "w",
			Description: "Generate web controller and views",
			Required:    false,
			Default:     false,
			Type:        "bool",
		},
		{
			Name:        "full",
			ShortName:   "f",
			Description: "Generate complete module with all components",
			Required:    false,
			Default:     false,
			Type:        "bool",
		},
	}
}

// Execute 执行命令
func (cmd *AddModuleCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	api := input.GetOption("api").(bool)
	web := input.GetOption("web").(bool)
	full := input.GetOption("full").(bool)

	return cmd.generator.GenerateModule(name, api, web, full)
}

// AddServiceCommand 添加服务命令
type AddServiceCommand struct {
	generator *Generator
}

// NewAddServiceCommand 创建新的添加服务命令
func NewAddServiceCommand(generator *Generator) *AddServiceCommand {
	return &AddServiceCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *AddServiceCommand) GetName() string {
	return "add:service"
}

// GetDescription 获取命令描述
func (cmd *AddServiceCommand) GetDescription() string {
	return "Add a new service class"
}

// GetSignature 获取命令签名
func (cmd *AddServiceCommand) GetSignature() string {
	return "add:service <name> [--interface]"
}

// GetArguments 获取命令参数
func (cmd *AddServiceCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the service",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *AddServiceCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "interface",
			ShortName:   "i",
			Description: "Generate interface for the service",
			Required:    false,
			Default:     false,
			Type:        "bool",
		},
	}
}

// Execute 执行命令
func (cmd *AddServiceCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	withInterface := input.GetOption("interface").(bool)

	return cmd.generator.GenerateService(name, withInterface)
}

// AddRepositoryCommand 添加仓库命令
type AddRepositoryCommand struct {
	generator *Generator
}

// NewAddRepositoryCommand 创建新的添加仓库命令
func NewAddRepositoryCommand(generator *Generator) *AddRepositoryCommand {
	return &AddRepositoryCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *AddRepositoryCommand) GetName() string {
	return "add:repository"
}

// GetDescription 获取命令描述
func (cmd *AddRepositoryCommand) GetDescription() string {
	return "Add a new repository class"
}

// GetSignature 获取命令签名
func (cmd *AddRepositoryCommand) GetSignature() string {
	return "add:repository <name> [--model=] [--interface]"
}

// GetArguments 获取命令参数
func (cmd *AddRepositoryCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the repository",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *AddRepositoryCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "model",
			ShortName:   "m",
			Description: "The model name for the repository",
			Required:    false,
			Default:     "",
			Type:        "string",
		},
		{
			Name:        "interface",
			ShortName:   "i",
			Description: "Generate interface for the repository",
			Required:    false,
			Default:     false,
			Type:        "bool",
		},
	}
}

// Execute 执行命令
func (cmd *AddRepositoryCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	model := input.GetOption("model").(string)
	withInterface := input.GetOption("interface").(bool)

	return cmd.generator.GenerateRepository(name, model, withInterface)
}

// AddValidatorCommand 添加验证器命令
type AddValidatorCommand struct {
	generator *Generator
}

// NewAddValidatorCommand 创建新的添加验证器命令
func NewAddValidatorCommand(generator *Generator) *AddValidatorCommand {
	return &AddValidatorCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *AddValidatorCommand) GetName() string {
	return "add:validator"
}

// GetDescription 获取命令描述
func (cmd *AddValidatorCommand) GetDescription() string {
	return "Add a new validator class"
}

// GetSignature 获取命令签名
func (cmd *AddValidatorCommand) GetSignature() string {
	return "add:validator <name> [--rules=]"
}

// GetArguments 获取命令参数
func (cmd *AddValidatorCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the validator",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *AddValidatorCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "rules",
			ShortName:   "r",
			Description: "Validation rules (comma-separated)",
			Required:    false,
			Default:     "",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *AddValidatorCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	rules := input.GetOption("rules").(string)

	return cmd.generator.GenerateValidator(name, rules)
}

// AddEventCommand 添加事件命令
type AddEventCommand struct {
	generator *Generator
}

// NewAddEventCommand 创建新的添加事件命令
func NewAddEventCommand(generator *Generator) *AddEventCommand {
	return &AddEventCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *AddEventCommand) GetName() string {
	return "add:event"
}

// GetDescription 获取命令描述
func (cmd *AddEventCommand) GetDescription() string {
	return "Add a new event and listener"
}

// GetSignature 获取命令签名
func (cmd *AddEventCommand) GetSignature() string {
	return "add:event <name> [--listener] [--queue]"
}

// GetArguments 获取命令参数
func (cmd *AddEventCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "name",
			Description: "The name of the event",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *AddEventCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "listener",
			ShortName:   "l",
			Description: "Generate listener for the event",
			Required:    false,
			Default:     true,
			Type:        "bool",
		},
		{
			Name:        "queue",
			ShortName:   "q",
			Description: "Make the listener queued",
			Required:    false,
			Default:     false,
			Type:        "bool",
		},
	}
}

// Execute 执行命令
func (cmd *AddEventCommand) Execute(input Input) error {
	name := input.GetArgument("name").(string)
	withListener := input.GetOption("listener").(bool)
	queued := input.GetOption("queue").(bool)

	return cmd.generator.GenerateEvent(name, withListener, queued)
}

// getEnvExampleTemplate 获取 .env.example 模板内容
func getEnvExampleTemplate() string {
	return `# =============================================================================
# Laravel-Go Application Environment Configuration
# =============================================================================

# =============================================================================
# Application Configuration
# =============================================================================
APP_NAME=Laravel-Go App
APP_ENV=local
APP_DEBUG=true
APP_URL=http://localhost:8080
APP_KEY=base64:your-32-character-app-key-here
APP_TIMEZONE=UTC
APP_LOCALE=en

# =============================================================================
# Server Configuration
# =============================================================================
PORT=8080
HOST=0.0.0.0
READ_TIMEOUT=30
WRITE_TIMEOUT=30
IDLE_TIMEOUT=120

# =============================================================================
# Database Configuration
# =============================================================================
DB_CONNECTION=sqlite
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=app.db
DB_USERNAME=
DB_PASSWORD=
DB_CHARSET=utf8mb4
DB_COLLATION=utf8mb4_unicode_ci
DB_PREFIX=

# PostgreSQL Configuration
DB_PG_HOST=127.0.0.1
DB_PG_PORT=5432
DB_PG_DATABASE=laravel_go
DB_PG_USERNAME=postgres
DB_PG_PASSWORD=
DB_PG_SSLMODE=disable

# MySQL Configuration
DB_MYSQL_HOST=127.0.0.1
DB_MYSQL_PORT=3306
DB_MYSQL_DATABASE=laravel_go
DB_MYSQL_USERNAME=root
DB_MYSQL_PASSWORD=
DB_MYSQL_CHARSET=utf8mb4
DB_MYSQL_COLLATION=utf8mb4_unicode_ci

# =============================================================================
# Cache Configuration
# =============================================================================
CACHE_DRIVER=memory
CACHE_PREFIX=laravel_go_
CACHE_TTL=3600

# Redis Cache Configuration
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_CACHE_DB=1
REDIS_SESSION_DB=2
REDIS_QUEUE_DB=3

# Memcached Configuration
MEMCACHED_HOST=127.0.0.1
MEMCACHED_PORT=11211
MEMCACHED_WEIGHT=100

# =============================================================================
# Session Configuration
# =============================================================================
SESSION_DRIVER=memory
SESSION_LIFETIME=120
SESSION_ENCRYPT=false
SESSION_FILES=/tmp/sessions
SESSION_CONNECTION=default
SESSION_TABLE=sessions
SESSION_STORE=redis

# =============================================================================
# Queue Configuration
# =============================================================================
QUEUE_CONNECTION=sync
QUEUE_DRIVER=sync
QUEUE_FAILED_DRIVER=database-uuids
QUEUE_FAILED_TABLE=failed_jobs

# Redis Queue Configuration
QUEUE_REDIS_CONNECTION=default
QUEUE_REDIS_QUEUE=default

# Database Queue Configuration
QUEUE_DB_TABLE=jobs
QUEUE_DB_RETRY_AFTER=90

# =============================================================================
# Mail Configuration
# =============================================================================
MAIL_MAILER=smtp
MAIL_HOST=smtp.mailtrap.io
MAIL_PORT=2525
MAIL_USERNAME=
MAIL_PASSWORD=
MAIL_ENCRYPTION=tls
MAIL_FROM_ADDRESS=hello@example.com
MAIL_FROM_NAME="${APP_NAME}"

# Mailgun Configuration
MAILGUN_DOMAIN=
MAILGUN_SECRET=
MAILGUN_ENDPOINT=api.mailgun.net

# SendGrid Configuration
SENDGRID_API_KEY=

# =============================================================================
# Logging Configuration
# =============================================================================
LOG_CHANNEL=stack
LOG_LEVEL=debug
LOG_DAYS=14
LOG_SLACK_WEBHOOK_URL=

# =============================================================================
# Authentication Configuration
# =============================================================================
AUTH_DRIVER=jwt
AUTH_GUARD=web
AUTH_PROVIDERS=users

# JWT Configuration
JWT_SECRET=your-jwt-secret-key-here
JWT_TTL=60
JWT_REFRESH_TTL=20160
JWT_ALGO=HS256

# =============================================================================
# File Storage Configuration
# =============================================================================
FILESYSTEM_DISK=local
FILESYSTEM_DRIVER=local

# Local Storage
STORAGE_PATH=storage/app/public
STORAGE_URL=storage

# S3 Configuration
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_DEFAULT_REGION=us-east-1
AWS_BUCKET=
AWS_USE_PATH_STYLE_ENDPOINT=false

# =============================================================================
# Security Configuration
# =============================================================================
BCRYPT_ROUNDS=12
HASH_DRIVER=bcrypt
ENCRYPTION_KEY=your-encryption-key-here

# CORS Configuration
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=*
CORS_EXPOSED_HEADERS=
CORS_MAX_AGE=86400
CORS_SUPPORTS_CREDENTIALS=false

# =============================================================================
# Rate Limiting Configuration
# =============================================================================
RATE_LIMIT_ENABLED=true
RATE_LIMIT_ATTEMPTS=60
RATE_LIMIT_DECAY_MINUTES=1
RATE_LIMIT_HEADERS=true

# =============================================================================
# Monitoring Configuration
# =============================================================================
MONITORING_ENABLED=false
PROMETHEUS_ENABLED=false
PROMETHEUS_PORT=9090

# =============================================================================
# Task Scheduling Configuration
# =============================================================================
SCHEDULER_ENABLED=true
SCHEDULER_DRIVER=cron

# =============================================================================
# WebSocket Configuration
# =============================================================================
WEBSOCKET_ENABLED=false
WEBSOCKET_PORT=8081
WEBSOCKET_HOST=0.0.0.0

# =============================================================================
# Internationalization Configuration
# =============================================================================
I18N_ENABLED=false
I18N_DEFAULT_LOCALE=en
I18N_FALLBACK_LOCALE=en
I18N_AVAILABLE_LOCALES=en,zh,ja

# =============================================================================
# API Configuration
# =============================================================================
API_VERSION=v1
API_PREFIX=api
API_RATE_LIMIT=60,1
API_THROTTLE=60,1

# =============================================================================
# Development Configuration
# =============================================================================
DEVELOPMENT_MODE=true
PROFILING_ENABLED=false
DEBUG_BAR_ENABLED=false

# =============================================================================
# Testing Configuration
# =============================================================================
TESTING_DATABASE=testing.db
TESTING_CACHE_DRIVER=array
TESTING_SESSION_DRIVER=array
TESTING_QUEUE_DRIVER=sync

# =============================================================================
# External Services Configuration
# =============================================================================

# Slack Configuration
SLACK_WEBHOOK_URL=
SLACK_CHANNEL=#general

# GitHub Configuration
GITHUB_TOKEN=
GITHUB_WEBHOOK_SECRET=

# Stripe Configuration
STRIPE_KEY=
STRIPE_SECRET=
STRIPE_WEBHOOK_SECRET=

# =============================================================================
# Custom Application Configuration
# =============================================================================
# Add your custom configuration variables below
# CUSTOM_VARIABLE=value`
}

// generateMicroserviceFiles 生成微服务相关文件
func (cmd *InitCommand) generateMicroserviceFiles(config *ProjectConfig, projectDir string) error {
	// 生成gRPC相关文件
	if config.GRPC != "none" {
		if err := cmd.generateGRPCFiles(config, projectDir); err != nil {
			return fmt.Errorf("failed to generate gRPC files: %w", err)
		}
		cmd.output.Success("✅ gRPC 文件已生成")
	}

	// 生成API网关相关文件
	if config.APIGateway != "none" {
		if err := cmd.generateAPIGatewayFiles(config, projectDir); err != nil {
			return fmt.Errorf("failed to generate API gateway files: %w", err)
		}
		cmd.output.Success("✅ API Gateway 文件已生成")
	}

	return nil
}

// generateGRPCFiles 生成gRPC相关文件
func (cmd *InitCommand) generateGRPCFiles(config *ProjectConfig, projectDir string) error {
	grpcFiles := map[string]string{
		"proto/user.proto": `syntax = "proto3";

package user;

option go_package = "` + projectDir + `/proto/user";

// 用户服务定义
service UserService {
  // 获取用户信息
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  
  // 创建用户
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  
  // 更新用户
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  
  // 删除用户
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  
  // 获取用户列表
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

// 用户信息
message User {
  int64 id = 1;
  string name = 2;
  string email = 3;
  string phone = 4;
  string avatar = 5;
  string status = 6;
  string created_at = 7;
  string updated_at = 8;
}

// 获取用户请求
message GetUserRequest {
  int64 id = 1;
}

// 获取用户响应
message GetUserResponse {
  User user = 1;
  string message = 2;
  int32 code = 3;
}

// 创建用户请求
message CreateUserRequest {
  string name = 1;
  string email = 2;
  string phone = 3;
  string password = 4;
}

// 创建用户响应
message CreateUserResponse {
  User user = 1;
  string message = 2;
  int32 code = 3;
}

// 更新用户请求
message UpdateUserRequest {
  int64 id = 1;
  string name = 2;
  string email = 3;
  string phone = 4;
  string avatar = 5;
}

// 更新用户响应
message UpdateUserResponse {
  User user = 1;
  string message = 2;
  int32 code = 3;
}

// 删除用户请求
message DeleteUserRequest {
  int64 id = 1;
}

// 删除用户响应
message DeleteUserResponse {
  string message = 1;
  int32 code = 2;
}

// 获取用户列表请求
message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
  string search = 3;
}

// 获取用户列表响应
message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
  int32 page = 3;
  int32 page_size = 4;
  string message = 5;
  int32 code = 6;
}`,
		"grpc/server/server.go": `package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "` + projectDir + `/proto/user"
)

// UserServer 用户服务实现
type UserServer struct {
	pb.UnimplementedUserServiceServer
}

// NewUserServer 创建用户服务实例
func NewUserServer() *UserServer {
	return &UserServer{}
}

// GetUser 获取用户信息
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// TODO: 实现获取用户逻辑
	user := &pb.User{
		Id:        req.Id,
		Name:      "示例用户",
		Email:     "user@example.com",
		Phone:     "13800138000",
		Status:    "active",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}

	return &pb.GetUserResponse{
		User:    user,
		Message: "获取用户成功",
		Code:    200,
	}, nil
}

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// TODO: 实现创建用户逻辑
	user := &pb.User{
		Id:        1,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    "active",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}

	return &pb.CreateUserResponse{
		User:    user,
		Message: "创建用户成功",
		Code:    201,
	}, nil
}

// UpdateUser 更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO: 实现更新用户逻辑
	user := &pb.User{
		Id:        req.Id,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		Status:    "active",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}

	return &pb.UpdateUserResponse{
		User:    user,
		Message: "更新用户成功",
		Code:    200,
	}, nil
}

// DeleteUser 删除用户
func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// TODO: 实现删除用户逻辑
	return &pb.DeleteUserResponse{
		Message: "删除用户成功",
		Code:    200,
	}, nil
}

// ListUsers 获取用户列表
func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// TODO: 实现获取用户列表逻辑
	users := []*pb.User{
		{
			Id:        1,
			Name:      "用户1",
			Email:     "user1@example.com",
			Phone:     "13800138001",
			Status:    "active",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		},
		{
			Id:        2,
			Name:      "用户2",
			Email:     "user2@example.com",
			Phone:     "13800138002",
			Status:    "active",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		},
	}

	return &pb.ListUsersResponse{
		Users:    users,
		Total:    2,
		Page:     req.Page,
		PageSize: req.PageSize,
		Message:  "获取用户列表成功",
		Code:     200,
	}, nil
}

// StartGRPCServer 启动gRPC服务器
func StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServer{})
	
	// 启用反射服务（用于调试）
	reflection.Register(s)

	log.Printf("🚀 gRPC Server starting on %s", port)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}`,
		"grpc/client/client.go": `package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "` + projectDir + `/proto/user"
)

// UserClient gRPC用户客户端
type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserClient 创建用户客户端
func NewUserClient(serverAddr string) (*UserClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	client := pb.NewUserServiceClient(conn)
	return &UserClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭连接
func (c *UserClient) Close() error {
	return c.conn.Close()
}

// GetUser 获取用户
func (c *UserClient) GetUser(id int64) (*pb.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return c.client.GetUser(ctx, &pb.GetUserRequest{Id: id})
}

// CreateUser 创建用户
func (c *UserClient) CreateUser(name, email, phone, password string) (*pb.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return c.client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: password,
	})
}

// UpdateUser 更新用户
func (c *UserClient) UpdateUser(id int64, name, email, phone, avatar string) (*pb.UpdateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return c.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:     id,
		Name:   name,
		Email:  email,
		Phone:  phone,
		Avatar: avatar,
	})
}

// DeleteUser 删除用户
func (c *UserClient) DeleteUser(id int64) (*pb.DeleteUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return c.client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
}

// ListUsers 获取用户列表
func (c *UserClient) ListUsers(page, pageSize int32, search string) (*pb.ListUsersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return c.client.ListUsers(ctx, &pb.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})
}

// ExampleUsage 使用示例
func ExampleUsage() {
	client, err := NewUserClient("localhost:9090")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// 获取用户
	user, err := client.GetUser(1)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
	} else {
		log.Printf("User: %+v", user.User)
	}

	// 创建用户
	createResp, err := client.CreateUser("张三", "zhangsan@example.com", "13800138000", "password123")
	if err != nil {
		log.Printf("Failed to create user: %v", err)
	} else {
		log.Printf("Created user: %+v", createResp.User)
	}
}`,
		"grpc/interceptors/logging.go": `package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor 日志拦截器
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	
	// 调用实际的RPC方法
	resp, err := handler(ctx, req)
	
	// 记录日志
	duration := time.Since(start)
	statusCode := codes.OK
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}
		log.Printf("gRPC: %s | %s | %v | %s", info.FullMethod, statusCode, duration, err)
	} else {
		log.Printf("gRPC: %s | %s | %v", info.FullMethod, statusCode, duration)
	}
	
	return resp, err
}

// RecoveryInterceptor 恢复拦截器
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("gRPC panic: %v", r)
			err = status.Errorf(codes.Internal, "Internal server error")
		}
	}()
	
	return handler(ctx, req)
}`,
		"grpc/interceptors/auth.go": `package interceptors

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor 认证拦截器
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 跳过认证的方法
	skipAuthMethods := map[string]bool{
		"/user.UserService/GetUser": true,
		"/user.UserService/ListUsers": true,
	}
	
	if skipAuthMethods[info.FullMethod] {
		return handler(ctx, req)
	}
	
	// 从元数据中获取token
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	
	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	
	token := authHeader[0]
	if !strings.HasPrefix(token, "Bearer ") {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
	}
	
	// TODO: 验证token
	tokenValue := strings.TrimPrefix(token, "Bearer ")
	if tokenValue == "" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	
	// 将用户信息添加到上下文
	userID := "123" // TODO: 从token中解析用户ID
	newCtx := context.WithValue(ctx, "user_id", userID)
	
	return handler(newCtx, req)
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in context")
	}
	return userID, nil
}`,
	}

	// 创建gRPC文件
	for fileName, content := range grpcFiles {
		fullPath := filepath.Join(projectDir, fileName)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", fullPath, err)
		}
	}

	return nil
}

// generateAPIGatewayFiles 生成API网关相关文件
func (cmd *InitCommand) generateAPIGatewayFiles(config *ProjectConfig, projectDir string) error {
	gatewayFiles := map[string]string{
		"gateway/main.go": `package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "` + projectDir + `/proto/user"
)

// Gateway API网关
type Gateway struct {
	userClient pb.UserServiceClient
	router     *mux.Router
}

// NewGateway 创建网关实例
func NewGateway() (*Gateway, error) {
	// 连接gRPC服务
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	userClient := pb.NewUserServiceClient(conn)

	router := mux.NewRouter()
	gateway := &Gateway{
		userClient: userClient,
		router:     router,
	}

	// 注册路由
	gateway.registerRoutes()

	return gateway, nil
}

// registerRoutes 注册路由
func (gateway *Gateway) registerRoutes() {
	// 中间件
	gateway.router.Use(gateway.loggingMiddleware)
	gateway.router.Use(gateway.corsMiddleware)

	// API路由
	api := gateway.router.PathPrefix("/api/v1").Subrouter()
	
	// 用户相关路由
	api.HandleFunc("/users", gateway.getUsers).Methods("GET")
	api.HandleFunc("/users/{id}", gateway.getUser).Methods("GET")
	api.HandleFunc("/users", gateway.createUser).Methods("POST")
	api.HandleFunc("/users/{id}", gateway.updateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", gateway.deleteUser).Methods("DELETE")

	// 健康检查
	gateway.router.HandleFunc("/health", gateway.healthCheck).Methods("GET")
}

// loggingMiddleware 日志中间件
func (gateway *Gateway) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("API Gateway: %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// corsMiddleware CORS中间件
func (gateway *Gateway) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// healthCheck 健康检查
func (gateway *Gateway) healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
		"service": "api-gateway",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getUsers 获取用户列表
func (gateway *Gateway) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 从查询参数获取分页信息
	page := int32(1)
	pageSize := int32(10)
	search := r.URL.Query().Get("search")

	resp, err := gateway.userClient.ListUsers(ctx, &pb.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// getUser 获取单个用户
func (gateway *Gateway) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	
	// TODO: 解析用户ID
	id := int64(1) // 示例

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := gateway.userClient.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// createUser 创建用户
func (gateway *Gateway) createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string ` + "`json:\"name\"`" + `
		Email    string ` + "`json:\"email\"`" + `
		Phone    string ` + "`json:\"phone\"`" + `
		Password string ` + "`json:\"password\"`" + `
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := gateway.userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// updateUser 更新用户
func (gateway *Gateway) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	
	// TODO: 解析用户ID
	id := int64(1) // 示例

	var req struct {
		Name   string ` + "`json:\"name\"`" + `
		Email  string ` + "`json:\"email\"`" + `
		Phone  string ` + "`json:\"phone\"`" + `
		Avatar string ` + "`json:\"avatar\"`" + `
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := gateway.userClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:     id,
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Avatar: req.Avatar,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// deleteUser 删除用户
func (gateway *Gateway) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	
	// TODO: 解析用户ID
	id := int64(1) // 示例

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := gateway.userClient.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	gateway, err := NewGateway()
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	port := ":8080"
	if envPort := os.Getenv("GATEWAY_PORT"); envPort != "" {
		port = ":" + envPort
	}

	server := &http.Server{
		Addr:    port,
		Handler: gateway.router,
	}

	// 启动服务器
	go func() {
		fmt.Printf("🚀 API Gateway starting on http://localhost%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gateway error: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	fmt.Println("\n🛑 Shutting down API Gateway...")
	fmt.Println("✅ API Gateway stopped gracefully")
}`,
		"gateway/middleware/auth.go": `package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 跳过认证的路径
		skipAuthPaths := map[string]bool{
			"/health": true,
			"/api/v1/users": true, // GET请求
		}
		
		if skipAuthPaths[r.URL.Path] && r.Method == "GET" {
			next.ServeHTTP(w, r)
			return
		}
		
		// 获取Authorization头
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		// 验证Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		
		// TODO: 验证token
		// 这里应该调用认证服务验证token
		
		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: 实现限流逻辑
		// 这里可以使用Redis或其他存储来实现限流
		
		next.ServeHTTP(w, r)
	})
}`,
		"gateway/routes/routes.go": `package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *mux.Router) {
	// API v1 路由组
	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	
	// 用户路由
	registerUserRoutes(apiV1)
	
	// 其他服务路由
	registerOtherRoutes(apiV1)
}

// registerUserRoutes 注册用户相关路由
func registerUserRoutes(router *mux.Router) {
	router.HandleFunc("/users", handleGetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", handleGetUser).Methods("GET")
	router.HandleFunc("/users", handleCreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", handleUpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", handleDeleteUser).Methods("DELETE")
}

// registerOtherRoutes 注册其他服务路由
func registerOtherRoutes(router *mux.Router) {
	// TODO: 添加其他微服务的路由
	router.HandleFunc("/products", handleGetProducts).Methods("GET")
	router.HandleFunc("/orders", handleGetOrders).Methods("GET")
}

// 用户路由处理器
func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现获取用户列表逻辑
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(` + "`{\"message\": \"Get users endpoint\"}`" + `))
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现获取单个用户逻辑
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(` + "`{\"message\": \"Get user endpoint\"}`" + `))
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现创建用户逻辑
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(` + "`{\"message\": \"Create user endpoint\"}`" + `))
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现更新用户逻辑
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(` + "`{\"message\": \"Update user endpoint\"}`" + `))
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现删除用户逻辑
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(` + "`{\"message\": \"Delete user endpoint\"}`" + `))
}

// 其他服务路由处理器
func handleGetProducts(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现获取产品列表逻辑
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(` + "`{\"message\": \"Get products endpoint\"}`" + `))
}

func handleGetOrders(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现获取订单列表逻辑
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(` + "`{\"message\": \"Get orders endpoint\"}`" + `))
}`,
		"gateway/plugins/rate_limit.go": `package plugins

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RateLimiter 限流器
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter 创建限流器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// 清理过期的请求记录
	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[key] = validRequests
	}

	// 检查是否超过限制
	if len(rl.requests[key]) >= rl.limit {
		return false
	}

	// 记录当前请求
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

// GetRemaining 获取剩余请求次数
func (rl *RateLimiter) GetRemaining(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if requests, exists := rl.requests[key]; exists {
		return rl.limit - len(requests)
	}
	return rl.limit
}

// Reset 重置限流器
func (rl *RateLimiter) Reset(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.requests, key)
}

// RateLimitPlugin 限流插件
type RateLimitPlugin struct {
	limiter *RateLimiter
}

// NewRateLimitPlugin 创建限流插件
func NewRateLimitPlugin(limit int, window time.Duration) *RateLimitPlugin {
	return &RateLimitPlugin{
		limiter: NewRateLimiter(limit, window),
	}
}

// Process 处理请求
func (p *RateLimitPlugin) Process(ctx context.Context, req interface{}) (interface{}, error) {
	// TODO: 从请求中提取客户端标识
	clientID := "default"
	
	if !p.limiter.Allow(clientID) {
		return nil, fmt.Errorf("rate limit exceeded")
	}
	
	return req, nil
}`,
	}

	// 创建API网关文件
	for fileName, content := range gatewayFiles {
		fullPath := filepath.Join(projectDir, fileName)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", fullPath, err)
		}
	}

	return nil
}
