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

	// 如果提供了项目名称，创建项目目录
	var projectDir string
	if projectName != "" && projectName != "laravel-go-app" {
		projectDir = projectName
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			return fmt.Errorf("failed to create project directory %s: %w", projectDir, err)
		}
		cmd.output.Success(fmt.Sprintf("Created project directory: %s", projectDir))
	}

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
		".env.example": `# Application Configuration
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
		"app/controllers/home_controller.go": `package controllers

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
		"app/controllers/user_controller.go": `package controllers

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
	"` + projectName + `/app/controllers"
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

- app/controllers/ - 控制器
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
		"Makefile": `# Laravel-Go Project Makefile

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
	APP_ENV=production go run main.go`,
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

	cmd.output.Success(fmt.Sprintf("Project '%s' initialized successfully!", projectName))
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
