package console

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Generator 代码生成器
type Generator struct {
	output Output
}

// NewGenerator 创建新的代码生成器
func NewGenerator(output Output) *Generator {
	return &Generator{
		output: output,
	}
}

// GenerateController 生成控制器
func (g *Generator) GenerateController(name, namespace string) error {
	// 创建控制器目录
	controllerDir := filepath.Join("app", "controllers")
	if err := os.MkdirAll(controllerDir, 0755); err != nil {
		return fmt.Errorf("failed to create controller directory: %w", err)
	}

	// 生成控制器文件名
	controllerName := g.toPascalCase(name)
	fileName := strings.ToLower(name) + "_controller.go"
	filePath := filepath.Join(controllerDir, fileName)

	// 控制器模板
	controllerTemplate := `package controllers

import (
	"laravel-go/framework/http"
)

// {{ .ControllerName }} 控制器
type {{ .ControllerName }} struct {
	http.BaseController
}

// New{{ .ControllerName }} 创建新的控制器实例
func New{{ .ControllerName }}() *{{ .ControllerName }} {
	return &{{ .ControllerName }}{}
}

// Index 显示资源列表
func (c *{{ .ControllerName }}) Index() http.Response {
	return c.Json(map[string]interface{}{
		"message": "{{ .ControllerName }} Index",
	})
}

// Show 显示指定资源
func (c *{{ .ControllerName }}) Show(id string) http.Response {
	return c.Json(map[string]interface{}{
		"message": "{{ .ControllerName }} Show",
		"id":      id,
	})
}

// Store 存储新创建的资源
func (c *{{ .ControllerName }}) Store() http.Response {
	return c.Json(map[string]interface{}{
		"message": "{{ .ControllerName }} Store",
	})
}

// Update 更新指定资源
func (c *{{ .ControllerName }}) Update(id string) http.Response {
	return c.Json(map[string]interface{}{
		"message": "{{ .ControllerName }} Update",
		"id":      id,
	})
}

// Delete 删除指定资源
func (c *{{ .ControllerName }}) Delete(id string) http.Response {
	return c.Json(map[string]interface{}{
		"message": "{{ .ControllerName }} Delete",
		"id":      id,
	})
}
`

	// 解析模板
	tmpl, err := template.New("controller").Parse(controllerTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse controller template: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create controller file: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := map[string]interface{}{
		"ControllerName": controllerName,
		"Namespace":      namespace,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute controller template: %w", err)
	}

	g.output.Success(fmt.Sprintf("Controller created successfully: %s", filePath))
	return nil
}

// GenerateModel 生成模型
func (g *Generator) GenerateModel(name string, fields []string) error {
	// 创建模型目录
	modelDir := filepath.Join("app", "models")
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return fmt.Errorf("failed to create model directory: %w", err)
	}

	// 生成模型文件名
	modelName := g.toPascalCase(name)
	fileName := strings.ToLower(name) + ".go"
	filePath := filepath.Join(modelDir, fileName)

	// 模型模板
	modelTemplate := `package models

import (
	"laravel-go/framework/database"
)

// {{ .ModelName }} 模型
type {{ .ModelName }} struct {
	database.Model
	{{ range .Fields }}
	{{ .Name }} {{ .Type }} ` + "`" + `json:"{{ .JsonName }}"` + "`" + `
	{{ end }}
}

// TableName 获取表名
func (m *{{ .ModelName }}) TableName() string {
	return "{{ .TableName }}"
}

// New{{ .ModelName }} 创建新的模型实例
func New{{ .ModelName }}() *{{ .ModelName }} {
	return &{{ .ModelName }}{}
}
`

	// 解析字段
	var modelFields []map[string]string
	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) >= 2 {
			fieldName := g.toPascalCase(parts[0])
			fieldType := parts[1]
			jsonName := parts[0]

			modelFields = append(modelFields, map[string]string{
				"Name":     fieldName,
				"Type":     fieldType,
				"JsonName": jsonName,
			})
		}
	}

	// 解析模板
	tmpl, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse model template: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create model file: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := map[string]interface{}{
		"ModelName": modelName,
		"TableName": strings.ToLower(name) + "s",
		"Fields":    modelFields,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute model template: %w", err)
	}

	g.output.Success(fmt.Sprintf("Model created successfully: %s", filePath))
	return nil
}

// GenerateMiddleware 生成中间件
func (g *Generator) GenerateMiddleware(name string) error {
	// 创建中间件目录
	middlewareDir := filepath.Join("app", "middleware")
	if err := os.MkdirAll(middlewareDir, 0755); err != nil {
		return fmt.Errorf("failed to create middleware directory: %w", err)
	}

	// 生成中间件文件名
	middlewareName := g.toPascalCase(name)
	fileName := strings.ToLower(name) + "_middleware.go"
	filePath := filepath.Join(middlewareDir, fileName)

	// 中间件模板
	middlewareTemplate := `package middleware

import (
	"laravel-go/framework/http"
)

// {{ .MiddlewareName }} 中间件
type {{ .MiddlewareName }} struct {
	// 添加中间件配置字段
}

// New{{ .MiddlewareName }} 创建新的中间件实例
func New{{ .MiddlewareName }}() *{{ .MiddlewareName }} {
	return &{{ .MiddlewareName }}{}
}

// Handle 处理请求
func (m *{{ .MiddlewareName }}) Handle(request http.Request, next func(http.Request) http.Response) http.Response {
	// 前置处理逻辑
	// ...

	// 调用下一个中间件或处理器
	response := next(request)

	// 后置处理逻辑
	// ...

	return response
}
`

	// 解析模板
	tmpl, err := template.New("middleware").Parse(middlewareTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse middleware template: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create middleware file: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := map[string]interface{}{
		"MiddlewareName": middlewareName,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute middleware template: %w", err)
	}

	g.output.Success(fmt.Sprintf("Middleware created successfully: %s", filePath))
	return nil
}

// GenerateMigration 生成迁移文件
func (g *Generator) GenerateMigration(name string, table string, fields []string) error {
	// 创建迁移目录
	migrationDir := filepath.Join("database", "migrations")
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return fmt.Errorf("failed to create migration directory: %w", err)
	}

	// 生成迁移文件名
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.sql", timestamp, strings.ToLower(strings.ReplaceAll(name, " ", "_")))
	filePath := filepath.Join(migrationDir, fileName)

	// 迁移模板
	migrationTemplate := `-- Migration: {{ .Name }}
-- Description: {{ .Description }}
-- Version: {{ .Version }}

-- UP Migration
CREATE TABLE IF NOT EXISTS {{ .Table }} (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	{{ range .Fields }}
	{{ .Name }} {{ .Type }}{{ if .Nullable }}{{ else }} NOT NULL{{ end }}{{ if .Default }}{{ .Default }}{{ end }},
	{{ end }}
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- DOWN Migration
DROP TABLE IF EXISTS {{ .Table }};
`

	// 解析字段
	var migrationFields []map[string]string
	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) >= 2 {
			fieldName := parts[0]
			fieldType := parts[1]
			nullable := true
			defaultValue := ""

			if len(parts) > 2 {
				if parts[2] == "not_null" {
					nullable = false
				}
			}
			if len(parts) > 3 {
				defaultValue = " DEFAULT " + parts[3]
			}

			migrationFields = append(migrationFields, map[string]string{
				"Name":     fieldName,
				"Type":     fieldType,
				"Nullable": fmt.Sprintf("%t", nullable),
				"Default":  defaultValue,
			})
		}
	}

	// 解析模板
	tmpl, err := template.New("migration").Parse(migrationTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse migration template: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := map[string]interface{}{
		"Name":        name,
		"Description": fmt.Sprintf("Create %s table", table),
		"Version":     timestamp,
		"Table":       table,
		"Fields":      migrationFields,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute migration template: %w", err)
	}

	g.output.Success(fmt.Sprintf("Migration created successfully: %s", filePath))
	return nil
}

// GenerateTest 生成测试文件
func (g *Generator) GenerateTest(name, type_ string) error {
	// 创建测试目录
	testDir := filepath.Join("tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return fmt.Errorf("failed to create test directory: %w", err)
	}

	// 生成测试文件名
	testName := g.toPascalCase(name)
	fileName := strings.ToLower(name) + "_test.go"
	filePath := filepath.Join(testDir, fileName)

	// 测试模板
	testTemplate := `package tests

import (
	"testing"
)

// Test{{ .TestName }} 测试{{ .Type }}
func Test{{ .TestName }}(t *testing.T) {
	// 设置测试环境
	// ...

	// 执行测试
	// ...

	// 验证结果
	// ...
}

// Benchmark{{ .TestName }} 基准测试{{ .Type }}
func Benchmark{{ .TestName }}(b *testing.B) {
	// 设置基准测试环境
	// ...

	for i := 0; i < b.N; i++ {
		// 执行基准测试
		// ...
	}
}
`

	// 解析模板
	tmpl, err := template.New("test").Parse(testTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse test template: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := map[string]interface{}{
		"TestName": testName,
		"Type":     type_,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute test template: %w", err)
	}

	g.output.Success(fmt.Sprintf("Test created successfully: %s", filePath))
	return nil
}

// toPascalCase 转换为帕斯卡命名法
func (g *Generator) toPascalCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

// GenerateDockerConfig 生成Docker配置文件
func (g *Generator) GenerateDockerConfig(name, port, env string) error {
	// 创建Docker目录
	dockerDir := "."
	if err := os.MkdirAll(dockerDir, 0755); err != nil {
		return fmt.Errorf("failed to create docker directory: %w", err)
	}

	// 生成Dockerfile
	dockerfilePath := filepath.Join(dockerDir, "Dockerfile")
	if err := g.generateDockerfile(dockerfilePath, name, port, env); err != nil {
		return fmt.Errorf("failed to generate Dockerfile: %w", err)
	}

	// 生成docker-compose.yml
	composePath := filepath.Join(dockerDir, "docker-compose.yml")
	if err := g.generateDockerCompose(composePath, name, port, env); err != nil {
		return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
	}

	// 生成.dockerignore
	dockerignorePath := filepath.Join(dockerDir, ".dockerignore")
	if err := g.generateDockerignore(dockerignorePath); err != nil {
		return fmt.Errorf("failed to generate .dockerignore: %w", err)
	}

	g.output.Success("Docker configuration files generated successfully!")
	g.output.WriteLine(fmt.Sprintf("  - Dockerfile: %s", dockerfilePath))
	g.output.WriteLine(fmt.Sprintf("  - docker-compose.yml: %s", composePath))
	g.output.WriteLine(fmt.Sprintf("  - .dockerignore: %s", dockerignorePath))

	return nil
}

// GenerateK8sConfig 生成Kubernetes配置文件
func (g *Generator) GenerateK8sConfig(name, replicas, port, namespace string) error {
	// 创建k8s目录
	k8sDir := "k8s"
	if err := os.MkdirAll(k8sDir, 0755); err != nil {
		return fmt.Errorf("failed to create k8s directory: %w", err)
	}

	// 生成deployment.yaml
	deploymentPath := filepath.Join(k8sDir, "deployment.yaml")
	if err := g.generateK8sDeployment(deploymentPath, name, replicas, port, namespace); err != nil {
		return fmt.Errorf("failed to generate deployment.yaml: %w", err)
	}

	// 生成service.yaml
	servicePath := filepath.Join(k8sDir, "service.yaml")
	if err := g.generateK8sService(servicePath, name, port, namespace); err != nil {
		return fmt.Errorf("failed to generate service.yaml: %w", err)
	}

	// 生成ingress.yaml
	ingressPath := filepath.Join(k8sDir, "ingress.yaml")
	if err := g.generateK8sIngress(ingressPath, name, namespace); err != nil {
		return fmt.Errorf("failed to generate ingress.yaml: %w", err)
	}

	// 生成configmap.yaml
	configmapPath := filepath.Join(k8sDir, "configmap.yaml")
	if err := g.generateK8sConfigmap(configmapPath, name, namespace); err != nil {
		return fmt.Errorf("failed to generate configmap.yaml: %w", err)
	}

	g.output.Success("Kubernetes configuration files generated successfully!")
	g.output.WriteLine(fmt.Sprintf("  - deployment.yaml: %s", deploymentPath))
	g.output.WriteLine(fmt.Sprintf("  - service.yaml: %s", servicePath))
	g.output.WriteLine(fmt.Sprintf("  - ingress.yaml: %s", ingressPath))
	g.output.WriteLine(fmt.Sprintf("  - configmap.yaml: %s", configmapPath))

	return nil
}

// generateDockerfile 生成Dockerfile
func (g *Generator) generateDockerfile(filePath, name, port, env string) error {
	dockerfileTemplate := `# {{ .Name }} Dockerfile
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
CMD ["./main"]
`

	tmpl, err := template.New("dockerfile").Parse(dockerfileTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse dockerfile template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create dockerfile: %w", err)
	}
	defer file.Close()

	data := struct {
		Name string
		Port string
		Env  string
	}{
		Name: name,
		Port: port,
		Env:  env,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute dockerfile template: %w", err)
	}

	return nil
}

// generateDockerCompose 生成docker-compose.yml
func (g *Generator) generateDockerCompose(filePath, name, port, env string) error {
	composeTemplate := `version: '3.8'

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
    driver: bridge
`

	tmpl, err := template.New("compose").Parse(composeTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse compose template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create compose file: %w", err)
	}
	defer file.Close()

	data := struct {
		Name string
		Port string
		Env  string
	}{
		Name: name,
		Port: port,
		Env:  env,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute compose template: %w", err)
	}

	return nil
}

// generateDockerignore 生成.dockerignore
func (g *Generator) generateDockerignore(filePath string) error {
	dockerignoreContent := `# Git
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

# Docker
Dockerfile
docker-compose.yml
.dockerignore

# Kubernetes
k8s/

# Build artifacts
main
*.exe
`

	return os.WriteFile(filePath, []byte(dockerignoreContent), 0644)
}

// generateK8sDeployment 生成Kubernetes Deployment
func (g *Generator) generateK8sDeployment(filePath, name, replicas, port, namespace string) error {
	deploymentTemplate := `apiVersion: apps/v1
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
            path: /ready
            port: {{ .Port }}
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
`

	tmpl, err := template.New("deployment").Parse(deploymentTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse deployment template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create deployment file: %w", err)
	}
	defer file.Close()

	data := struct {
		Name      string
		Replicas  string
		Port      string
		Namespace string
	}{
		Name:      name,
		Replicas:  replicas,
		Port:      port,
		Namespace: namespace,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute deployment template: %w", err)
	}

	return nil
}

// generateK8sService 生成Kubernetes Service
func (g *Generator) generateK8sService(filePath, name, port, namespace string) error {
	serviceTemplate := `apiVersion: v1
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
  type: ClusterIP
`

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create service file: %w", err)
	}
	defer file.Close()

	data := struct {
		Name      string
		Port      string
		Namespace string
	}{
		Name:      name,
		Port:      port,
		Namespace: namespace,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute service template: %w", err)
	}

	return nil
}

// generateK8sIngress 生成Kubernetes Ingress
func (g *Generator) generateK8sIngress(filePath, name, namespace string) error {
	ingressTemplate := `apiVersion: networking.k8s.io/v1
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
              number: 80
`

	tmpl, err := template.New("ingress").Parse(ingressTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse ingress template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create ingress file: %w", err)
	}
	defer file.Close()

	data := struct {
		Name      string
		Namespace string
	}{
		Name:      name,
		Namespace: namespace,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute ingress template: %w", err)
	}

	return nil
}

// generateK8sConfigmap 生成Kubernetes ConfigMap
func (g *Generator) generateK8sConfigmap(filePath, name, namespace string) error {
	configmapTemplate := `apiVersion: v1
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
  REDIS_PORT: "6379"
`

	tmpl, err := template.New("configmap").Parse(configmapTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse configmap template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create configmap file: %w", err)
	}
	defer file.Close()

	data := struct {
		Name      string
		Namespace string
	}{
		Name:      name,
		Namespace: namespace,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute configmap template: %w", err)
	}

	return nil
}
