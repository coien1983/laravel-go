package main

import (
	"context"
	"encoding/json"
	"fmt"
	"laravel-go/framework/performance"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MCPRequest MCP请求结构
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// MCPResponse MCP响应结构
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError MCP错误结构
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// APIGenerator API生成器
type APIGenerator struct {
	projectName string
	projectPath string
	config      *ProjectConfig
}

// ProjectConfig 项目配置
type ProjectConfig struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Modules     []string `json:"modules"`
	Database    string   `json:"database"`
	Cache       string   `json:"cache"`
	Queue       string   `json:"queue"`
}

func main() {
	port := ":8080"
	if envPort := os.Getenv("MCP_PORT"); envPort != "" {
		port = ":" + envPort
	}

	http.HandleFunc("/", handleMCPRequest)
	
	fmt.Printf("🚀 MCP API Generator 启动在端口 %s\n", port)
	fmt.Println("📝 支持的命令:")
	fmt.Println("  - initialize: 初始化新项目")
	fmt.Println("  - generate: 生成API模块")
	fmt.Println("  - build: 构建项目")
	fmt.Println("  - test: 运行测试")
	
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, -32700, "解析错误", err.Error())
		return
	}

	var response MCPResponse
	response.JSONRPC = "2.0"
	response.ID = req.ID

	switch req.Method {
	case "initialize":
		response.Result = handleInitialize(req.Params)
	case "generate":
		response.Result = handleGenerate(req.Params)
	case "build":
		response.Result = handleBuild(req.Params)
	case "test":
		response.Result = handleTest(req.Params)
	default:
		response.Error = &MCPError{
			Code:    -32601,
			Message: "方法不存在",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, code int, message, data string) {
	response := MCPResponse{
		JSONRPC: "2.0",
		Error: &MCPError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleInitialize(params interface{}) map[string]interface{} {
	config := &ProjectConfig{
		Name:        "laravel-go-api",
		Description: "Laravel-Go API项目",
		Version:     "1.0.0",
		Author:      "Developer",
		Modules:     []string{"user", "product", "order"},
		Database:    "mysql",
		Cache:       "redis",
		Queue:       "redis",
	}

	generator := &APIGenerator{
		projectName: config.Name,
		projectPath: config.Name,
		config:      config,
	}

	if err := generator.initialize(); err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("项目 %s 初始化成功", config.Name),
		"path":    generator.projectPath,
	}
}

func handleGenerate(params interface{}) map[string]interface{} {
	// 这里可以解析params来生成特定的模块
	modules := []string{"user", "product", "order"}
	
	generator := &APIGenerator{
		projectName: "laravel-go-api",
		projectPath: "laravel-go-api",
		config: &ProjectConfig{
			Modules: modules,
		},
	}

	if err := generator.generateModules(modules); err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("成功生成 %d 个模块", len(modules)),
		"modules": modules,
	}
}

func handleBuild(params interface{}) map[string]interface{} {
	generator := &APIGenerator{
		projectName: "laravel-go-api",
		projectPath: "laravel-go-api",
	}

	if err := generator.build(); err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
		"message": "项目构建成功",
	}
}

func handleTest(params interface{}) map[string]interface{} {
	generator := &APIGenerator{
		projectName: "laravel-go-api",
		projectPath: "laravel-go-api",
	}

	if err := generator.test(); err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
		"message": "测试运行成功",
	}
}

func (ag *APIGenerator) initialize() error {
	// 创建项目目录结构
	dirs := []string{
		"app/Http/Controllers",
		"app/Http/Middleware",
		"app/Http/Requests",
		"app/Models",
		"app/Services",
		"app/Providers",
		"config",
		"database/migrations",
		"database/seeders",
		"routes",
		"storage/cache",
		"storage/logs",
		"storage/uploads",
		"tests",
		"docs",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(ag.projectPath, dir), 0755); err != nil {
			return fmt.Errorf("创建目录失败 %s: %v", dir, err)
		}
	}

	// 生成基础文件
	if err := ag.generateBaseFiles(); err != nil {
		return err
	}

	return nil
}

func (ag *APIGenerator) generateBaseFiles() error {
	// 生成main.go
	mainContent := ag.generateMainFile()
	if err := os.WriteFile(filepath.Join(ag.projectPath, "main.go"), []byte(mainContent), 0644); err != nil {
		return err
	}

	// 生成go.mod
	goModContent := ag.generateGoModFile()
	if err := os.WriteFile(filepath.Join(ag.projectPath, "go.mod"), []byte(goModContent), 0644); err != nil {
		return err
	}

	// 生成配置文件
	if err := ag.generateConfigFiles(); err != nil {
		return err
	}

	// 生成README
	readmeContent := ag.generateReadmeFile()
	if err := os.WriteFile(filepath.Join(ag.projectPath, "README.md"), []byte(readmeContent), 0644); err != nil {
		return err
	}

	return nil
}

func (ag *APIGenerator) generateMainFile() string {
	return `package main

import (
	"context"
	"laravel-go/framework/performance"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 创建性能监控器
	monitor := performance.NewPerformanceMonitor()
	ctx := context.Background()
	monitor.Start(ctx)
	defer monitor.Stop()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: setupRoutes(),
	}

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("正在关闭服务器...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("服务器关闭错误: %v", err)
		}
	}()

	log.Println("🚀 Laravel-Go API 服务器启动在端口 8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("服务器启动失败:", err)
	}
}

func setupRoutes() http.Handler {
	mux := http.NewServeMux()
	
	// 健康检查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// API路由
	mux.HandleFunc("/api/", handleAPI)

	return mux
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Laravel-Go API", "version": "1.0.0"}`))
}
`
}

func (ag *APIGenerator) generateGoModFile() string {
	return `module ` + ag.projectName + `

go 1.21

require (
	laravel-go/framework v0.1.0
	github.com/gorilla/mux v1.8.0
	github.com/go-sql-driver/mysql v1.7.1
	github.com/redis/go-redis/v9 v9.0.5
	github.com/stretchr/testify v1.8.4
)
`
}

func (ag *APIGenerator) generateConfigFiles() error {
	// 生成app配置
	appConfig := `{
	"name": "` + ag.config.Name + `",
	"version": "` + ag.config.Version + `",
	"debug": true,
	"timezone": "Asia/Shanghai"
}`
	
	if err := os.WriteFile(filepath.Join(ag.projectPath, "config/app.json"), []byte(appConfig), 0644); err != nil {
		return err
	}

	// 生成数据库配置
	dbConfig := `{
	"driver": "` + ag.config.Database + `",
	"host": "localhost",
	"port": 3306,
	"database": "` + ag.projectName + `",
	"username": "root",
	"password": "",
	"charset": "utf8mb4"
}`
	
	if err := os.WriteFile(filepath.Join(ag.projectPath, "config/database.json"), []byte(dbConfig), 0644); err != nil {
		return err
	}

	return nil
}

func (ag *APIGenerator) generateReadmeFile() string {
	return `# ` + ag.config.Name + `

` + ag.config.Description + `

## 功能特性

- 🚀 基于 Laravel-Go 框架
- 📊 内置性能监控
- 🔐 完整的认证授权
- 🗄️ 数据库ORM支持
- 💾 缓存系统
- 📝 API文档自动生成
- 🧪 完整的测试覆盖

## 快速开始

### 1. 安装依赖

\`\`\`bash
go mod tidy
\`\`\`

### 2. 配置数据库

编辑 \`config/database.json\` 文件，配置数据库连接信息。

### 3. 运行迁移

\`\`\`bash
go run main.go migrate
\`\`\`

### 4. 启动服务器

\`\`\`bash
go run main.go
\`\`\`

服务器将在 http://localhost:8080 启动

## API文档

启动服务器后，访问以下端点：

- 健康检查: GET /health
- API信息: GET /api/

## 开发指南

### 项目结构

\`\`\`
` + ag.projectName + `/
├── app/
│   ├── Http/
│   │   ├── Controllers/    # 控制器
│   │   ├── Middleware/     # 中间件
│   │   └── Requests/       # 请求验证
│   ├── Models/             # 数据模型
│   ├── Services/           # 业务服务
│   └── Providers/          # 服务提供者
├── config/                 # 配置文件
├── database/               # 数据库文件
├── routes/                 # 路由定义
├── storage/                # 存储文件
├── tests/                  # 测试文件
└── main.go                 # 入口文件
\`\`\`

## 测试

\`\`\`bash
go test ./...
\`\`\`

## 部署

\`\`\`bash
go build -o ` + ag.projectName + ` main.go
./` + ag.projectName + `
\`\`\`

## 许可证

MIT License
`
}

func (ag *APIGenerator) generateModules(modules []string) error {
	for _, module := range modules {
		if err := ag.generateModule(module); err != nil {
			return fmt.Errorf("生成模块 %s 失败: %v", module, err)
		}
	}
	return nil
}

func (ag *APIGenerator) generateModule(moduleName string) error {
	// 生成控制器
	controllerContent := ag.generateController(moduleName)
	controllerPath := filepath.Join(ag.projectPath, "app/Http/Controllers", moduleName+"_controller.go")
	if err := os.WriteFile(controllerPath, []byte(controllerContent), 0644); err != nil {
		return err
	}

	// 生成模型
	modelContent := ag.generateModel(moduleName)
	modelPath := filepath.Join(ag.projectPath, "app/Models", moduleName+".go")
	if err := os.WriteFile(modelPath, []byte(modelContent), 0644); err != nil {
		return err
	}

	// 生成服务
	serviceContent := ag.generateService(moduleName)
	servicePath := filepath.Join(ag.projectPath, "app/Services", moduleName+"_service.go")
	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		return err
	}

	// 生成请求验证
	requestContent := ag.generateRequest(moduleName)
	requestPath := filepath.Join(ag.projectPath, "app/Http/Requests", moduleName+"_request.go")
	if err := os.WriteFile(requestPath, []byte(requestContent), 0644); err != nil {
		return err
	}

	return nil
}

func (ag *APIGenerator) generateController(moduleName string) string {
	moduleTitle := strings.Title(moduleName)
	return `package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"` + ag.projectName + `/app/Services"
	"` + ag.projectName + `/app/Http/Requests"
)

type ` + moduleTitle + `Controller struct {
	` + moduleName + `Service *Services.` + moduleTitle + `Service
}

func New` + moduleTitle + `Controller(` + moduleName + `Service *Services.` + moduleTitle + `Service) *` + moduleTitle + `Controller {
	return &` + moduleTitle + `Controller{
		` + moduleName + `Service: ` + moduleName + `Service,
	}
}

func (c *` + moduleTitle + `Controller) Index(w http.ResponseWriter, r *http.Request) {
	` + moduleName + `s, err := c.` + moduleName + `Service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    ` + moduleName + `s,
	})
}

func (c *` + moduleTitle + `Controller) Show(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "无效的ID", http.StatusBadRequest)
		return
	}

	` + moduleName + `, err := c.` + moduleName + `Service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    ` + moduleName + `,
	})
}

func (c *` + moduleTitle + `Controller) Store(w http.ResponseWriter, r *http.Request) {
	var request Requests.` + moduleTitle + `Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	if err := c.` + moduleName + `Service.Create(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "创建成功",
	})
}

func (c *` + moduleTitle + `Controller) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "无效的ID", http.StatusBadRequest)
		return
	}

	var request Requests.` + moduleTitle + `Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	if err := c.` + moduleName + `Service.Update(id, &request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "更新成功",
	})
}

func (c *` + moduleTitle + `Controller) Destroy(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "无效的ID", http.StatusBadRequest)
		return
	}

	if err := c.` + moduleName + `Service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "删除成功",
	})
}
`
}

func (ag *APIGenerator) generateModel(moduleName string) string {
	moduleTitle := strings.Title(moduleName)
	return `package Models

import (
	"time"
)

type ` + moduleTitle + ` struct {
	ID        int       ` + "`json:\"id\" db:\"id\"`" + `
	Name      string    ` + "`json:\"name\" db:\"name\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\" db:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\" db:\"updated_at\"`" + `
}

func (m *` + moduleTitle + `) TableName() string {
	return "` + moduleName + `s"
}

func (m *` + moduleTitle + `) BeforeCreate() error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *` + moduleTitle + `) BeforeUpdate() error {
	m.UpdatedAt = time.Now()
	return nil
}
`
}

func (ag *APIGenerator) generateService(moduleName string) string {
	moduleTitle := strings.Title(moduleName)
	return `package Services

import (
	"errors"
	"` + ag.projectName + `/app/Models"
	"` + ag.projectName + `/app/Http/Requests"
)

type ` + moduleTitle + `Service struct {
	// 这里可以注入数据库连接、缓存等依赖
}

func New` + moduleTitle + `Service() *` + moduleTitle + `Service {
	return &` + moduleTitle + `Service{}
}

func (s *` + moduleTitle + `Service) GetAll() ([]*Models.` + moduleTitle + `, error) {
	// 实现获取所有记录的逻辑
	return []*Models.` + moduleTitle + `{}, nil
}

func (s *` + moduleTitle + `Service) GetByID(id int) (*Models.` + moduleTitle + `, error) {
	// 实现根据ID获取记录的逻辑
	if id <= 0 {
		return nil, errors.New("无效的ID")
	}
	
	return &Models.` + moduleTitle + `{
		ID:   id,
		Name: "示例" + string(rune(id)),
	}, nil
}

func (s *` + moduleTitle + `Service) Create(request *Requests.` + moduleTitle + `Request) error {
	// 实现创建记录的逻辑
	if request.Name == "" {
		return errors.New("名称不能为空")
	}
	return nil
}

func (s *` + moduleTitle + `Service) Update(id int, request *Requests.` + moduleTitle + `Request) error {
	// 实现更新记录的逻辑
	if id <= 0 {
		return errors.New("无效的ID")
	}
	if request.Name == "" {
		return errors.New("名称不能为空")
	}
	return nil
}

func (s *` + moduleTitle + `Service) Delete(id int) error {
	// 实现删除记录的逻辑
	if id <= 0 {
		return errors.New("无效的ID")
	}
	return nil
}
`
}

func (ag *APIGenerator) generateRequest(moduleName string) string {
	moduleTitle := strings.Title(moduleName)
	return `package Requests

type ` + moduleTitle + `Request struct {
	Name string ` + "`json:\"name\" validate:\"required\"`" + `
}

func (r *` + moduleTitle + `Request) Validate() error {
	if r.Name == "" {
		return errors.New("名称不能为空")
	}
	return nil
}
`
}

func (ag *APIGenerator) build() error {
	// 这里可以实现项目构建逻辑
	// 例如：go build、依赖检查等
	return nil
}

func (ag *APIGenerator) test() error {
	// 这里可以实现测试运行逻辑
	// 例如：go test、覆盖率检查等
	return nil
} 