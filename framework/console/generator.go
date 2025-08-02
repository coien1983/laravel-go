package console

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
	"github.com/coien1983/laravel-go/framework/http"
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
	"github.com/coien1983/laravel-go/framework/database"
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
	"github.com/coien1983/laravel-go/framework/http"
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

// toPascalCase 转换为PascalCase
func (g *Generator) toPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

// Docker 和 Kubernetes 辅助方法已移除

// GenerateModule 生成完整模块
func (g *Generator) GenerateModule(name string, api, web, full bool) error {
	moduleName := g.toPascalCase(name)

	g.output.Info(fmt.Sprintf("正在生成模块: %s", moduleName))

	// 生成模型
	if err := g.GenerateModel(name, []string{}); err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}

	// 生成服务
	if err := g.GenerateService(name, true); err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}

	// 生成仓库
	if err := g.GenerateRepository(name, name, true); err != nil {
		return fmt.Errorf("failed to generate repository: %w", err)
	}

	// 生成API控制器
	if api {
		if err := g.GenerateController(name, "api"); err != nil {
			return fmt.Errorf("failed to generate API controller: %w", err)
		}
	}

	// 生成Web控制器
	if web {
		if err := g.GenerateController(name+"Web", "web"); err != nil {
			return fmt.Errorf("failed to generate web controller: %w", err)
		}
	}

	// 生成完整模块的额外组件
	if full {
		// 生成验证器
		if err := g.GenerateValidator(name, ""); err != nil {
			return fmt.Errorf("failed to generate validator: %w", err)
		}

		// 生成事件
		if err := g.GenerateEvent(name+"Created", true, false); err != nil {
			return fmt.Errorf("failed to generate event: %w", err)
		}

		// 生成测试
		if err := g.GenerateTest(name, "integration"); err != nil {
			return fmt.Errorf("failed to generate test: %w", err)
		}
	}

	g.output.Success(fmt.Sprintf("✅ 模块 %s 生成完成！", moduleName))
	return nil
}

// GenerateService 生成服务类
func (g *Generator) GenerateService(name string, withInterface bool) error {
	serviceName := g.toPascalCase(name)

	// 创建服务目录
	serviceDir := filepath.Join("app", "services")
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// 生成服务文件名
	fileName := strings.ToLower(name) + "_service.go"
	filePath := filepath.Join(serviceDir, fileName)

	// 服务模板
	serviceTemplate := `package services

import (
	"{{ .ProjectName }}/app/models"
	"{{ .ProjectName }}/app/repositories"
)

{{ if .WithInterface }}
// {{ .ServiceName }}ServiceInterface 服务接口
type {{ .ServiceName }}ServiceInterface interface {
	GetAll() ([]models.{{ .ServiceName }}, error)
	GetByID(id uint) (*models.{{ .ServiceName }}, error)
	Create(data map[string]interface{}) (*models.{{ .ServiceName }}, error)
	Update(id uint, data map[string]interface{}) (*models.{{ .ServiceName }}, error)
	Delete(id uint) error
}
{{ end }}

// {{ .ServiceName }}Service 服务实现
type {{ .ServiceName }}Service struct {
	repo repositories.{{ .ServiceName }}RepositoryInterface
}

// New{{ .ServiceName }}Service 创建新的服务实例
func New{{ .ServiceName }}Service(repo repositories.{{ .ServiceName }}RepositoryInterface) {{ if .WithInterface }}{{ .ServiceName }}ServiceInterface{{ else }}*{{ .ServiceName }}Service{{ end }} {
	return &{{ .ServiceName }}Service{
		repo: repo,
	}
}

// GetAll 获取所有记录
func (s *{{ .ServiceName }}Service) GetAll() ([]models.{{ .ServiceName }}, error) {
	return s.repo.GetAll()
}

// GetByID 根据ID获取记录
func (s *{{ .ServiceName }}Service) GetByID(id uint) (*models.{{ .ServiceName }}, error) {
	return s.repo.GetByID(id)
}

// Create 创建新记录
func (s *{{ .ServiceName }}Service) Create(data map[string]interface{}) (*models.{{ .ServiceName }}, error) {
	return s.repo.Create(data)
}

// Update 更新记录
func (s *{{ .ServiceName }}Service) Update(id uint, data map[string]interface{}) (*models.{{ .ServiceName }}, error) {
	return s.repo.Update(id, data)
}

// Delete 删除记录
func (s *{{ .ServiceName }}Service) Delete(id uint) error {
	return s.repo.Delete(id)
}
`

	// 解析模板
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create service file: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := map[string]interface{}{
		"ServiceName":   serviceName,
		"ProjectName":   g.getProjectName(),
		"WithInterface": withInterface,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute service template: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ 服务 %s 生成完成: %s", serviceName, filePath))
	return nil
}

// GenerateRepository 生成仓库类
func (g *Generator) GenerateRepository(name, modelName string, withInterface bool) error {
	repoName := g.toPascalCase(name)
	modelClassName := g.toPascalCase(modelName)

	// 创建仓库目录
	repoDir := filepath.Join("app", "repositories")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return fmt.Errorf("failed to create repository directory: %w", err)
	}

	// 生成仓库文件名
	fileName := strings.ToLower(name) + "_repository.go"
	filePath := filepath.Join(repoDir, fileName)

	// 仓库模板
	repoTemplate := `package repositories

import (
	"{{ .ProjectName }}/app/models"
	"gorm.io/gorm"
)

{{ if .WithInterface }}
// {{ .RepoName }}RepositoryInterface 仓库接口
type {{ .RepoName }}RepositoryInterface interface {
	GetAll() ([]models.{{ .ModelClassName }}, error)
	GetByID(id uint) (*models.{{ .ModelClassName }}, error)
	Create(data map[string]interface{}) (*models.{{ .ModelClassName }}, error)
	Update(id uint, data map[string]interface{}) (*models.{{ .ModelClassName }}, error)
	Delete(id uint) error
}
{{ end }}

// {{ .RepoName }}Repository 仓库实现
type {{ .RepoName }}Repository struct {
	db *gorm.DB
}

// New{{ .RepoName }}Repository 创建新的仓库实例
func New{{ .RepoName }}Repository(db *gorm.DB) {{ if .WithInterface }}{{ .RepoName }}RepositoryInterface{{ else }}*{{ .RepoName }}Repository{{ end }} {
	return &{{ .RepoName }}Repository{
		db: db,
	}
}

// GetAll 获取所有记录
func (r *{{ .RepoName }}Repository) GetAll() ([]models.{{ .ModelClassName }}, error) {
	var items []models.{{ .ModelClassName }}
	err := r.db.Find(&items).Error
	return items, err
}

// GetByID 根据ID获取记录
func (r *{{ .RepoName }}Repository) GetByID(id uint) (*models.{{ .ModelClassName }}, error) {
	var item models.{{ .ModelClassName }}
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 创建新记录
func (r *{{ .RepoName }}Repository) Create(data map[string]interface{}) (*models.{{ .ModelClassName }}, error) {
	item := models.{{ .ModelClassName }}{}
	
	// 设置字段值
	for key, value := range data {
		switch key {
		// 在这里添加字段映射
		}
	}
	
	err := r.db.Create(&item).Error
	return &item, err
}

// Update 更新记录
func (r *{{ .RepoName }}Repository) Update(id uint, data map[string]interface{}) (*models.{{ .ModelClassName }}, error) {
	var item models.{{ .ModelClassName }}
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	
	// 更新字段值
	for key, value := range data {
		switch key {
		// 在这里添加字段映射
		}
	}
	
	err := r.db.Save(&item).Error
	return &item, err
}

// Delete 删除记录
func (r *{{ .RepoName }}Repository) Delete(id uint) error {
	return r.db.Delete(&models.{{ .ModelClassName }}{}, id).Error
}
`

	// 解析模板
	tmpl, err := template.New("repository").Parse(repoTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse repository template: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create repository file: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := map[string]interface{}{
		"RepoName":       repoName,
		"ModelClassName": modelClassName,
		"ProjectName":    g.getProjectName(),
		"WithInterface":  withInterface,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute repository template: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ 仓库 %s 生成完成: %s", repoName, filePath))
	return nil
}

// GenerateValidator 生成验证器
func (g *Generator) GenerateValidator(name, rules string) error {
	validatorName := g.toPascalCase(name)

	// 创建验证器目录
	validatorDir := filepath.Join("app", "validators")
	if err := os.MkdirAll(validatorDir, 0755); err != nil {
		return fmt.Errorf("failed to create validator directory: %w", err)
	}

	// 生成验证器文件名
	fileName := strings.ToLower(name) + "_validator.go"
	filePath := filepath.Join(validatorDir, fileName)

	// 验证器模板
	validatorTemplate := `package validators

import (
	"github.com/go-playground/validator/v10"
)

// {{ .ValidatorName }}Validator {{ .ValidatorName }} 验证器
type {{ .ValidatorName }}Validator struct {
	validate *validator.Validate
}

// New{{ .ValidatorName }}Validator 创建新的验证器实例
func New{{ .ValidatorName }}Validator() *{{ .ValidatorName }}Validator {
	return &{{ .ValidatorName }}Validator{
		validate: validator.New(),
	}
}

// ValidateCreate 验证创建请求
func (v *{{ .ValidatorName }}Validator) ValidateCreate(data map[string]interface{}) error {
	// 在这里添加验证规则
	{{ if .Rules }}
	// 示例验证规则
	if err := v.validate.Var(data["name"], "required,min=2,max=50"); err != nil {
		return err
	}
	{{ end }}
	
	return nil
}

// ValidateUpdate 验证更新请求
func (v *{{ .ValidatorName }}Validator) ValidateUpdate(data map[string]interface{}) error {
	// 在这里添加验证规则
	{{ if .Rules }}
	// 示例验证规则
	if name, exists := data["name"]; exists {
		if err := v.validate.Var(name, "min=2,max=50"); err != nil {
			return err
		}
	}
	{{ end }}
	
	return nil
}
`

	// 解析模板
	tmpl, err := template.New("validator").Parse(validatorTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse validator template: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create validator file: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := map[string]interface{}{
		"ValidatorName": validatorName,
		"Rules":         rules != "",
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute validator template: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ 验证器 %s 生成完成: %s", validatorName, filePath))
	return nil
}

// GenerateEvent 生成事件和监听器
func (g *Generator) GenerateEvent(name string, withListener, queued bool) error {
	eventName := g.toPascalCase(name)

	// 创建事件目录
	eventDir := filepath.Join("app", "events")
	if err := os.MkdirAll(eventDir, 0755); err != nil {
		return fmt.Errorf("failed to create event directory: %w", err)
	}

	// 生成事件文件名
	fileName := strings.ToLower(name) + "_event.go"
	filePath := filepath.Join(eventDir, fileName)

	// 事件模板
	eventTemplate := `package events

import (
	"time"
)

// {{ .EventName }}Event {{ .EventName }} 事件
type {{ .EventName }}Event struct {
	Data      interface{}
	Timestamp time.Time
}

// New{{ .EventName }}Event 创建新的事件实例
func New{{ .EventName }}Event(data interface{}) *{{ .EventName }}Event {
	return &{{ .EventName }}Event{
		Data:      data,
		Timestamp: time.Now(),
	}
}

// GetData 获取事件数据
func (e *{{ .EventName }}Event) GetData() interface{} {
	return e.Data
}

// GetTimestamp 获取事件时间戳
func (e *{{ .EventName }}Event) GetTimestamp() time.Time {
	return e.Timestamp
}
`

	// 解析模板
	tmpl, err := template.New("event").Parse(eventTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse event template: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create event file: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := map[string]interface{}{
		"EventName": eventName,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute event template: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ 事件 %s 生成完成: %s", eventName, filePath))

	// 生成监听器
	if withListener {
		if err := g.generateListener(name, queued); err != nil {
			return fmt.Errorf("failed to generate listener: %w", err)
		}
	}

	return nil
}

// generateListener 生成监听器
func (g *Generator) generateListener(eventName string, queued bool) error {
	listenerName := g.toPascalCase(eventName) + "Listener"

	// 创建监听器目录
	listenerDir := filepath.Join("app", "listeners")
	if err := os.MkdirAll(listenerDir, 0755); err != nil {
		return fmt.Errorf("failed to create listener directory: %w", err)
	}

	// 生成监听器文件名
	fileName := strings.ToLower(eventName) + "_listener.go"
	filePath := filepath.Join(listenerDir, fileName)

	// 监听器模板
	listenerTemplate := `package listeners

import (
	"{{ .ProjectName }}/app/events"
	"log"
)

// {{ .ListenerName }} {{ .EventName }} 监听器
type {{ .ListenerName }} struct {
	{{ if .Queued }}queued bool{{ end }}
}

// New{{ .ListenerName }} 创建新的监听器实例
func New{{ .ListenerName }}() *{{ .ListenerName }} {
	return &{{ .ListenerName }}{
		{{ if .Queued }}queued: true,{{ end }}
	}
}

// Handle 处理事件
func (l *{{ .ListenerName }}) Handle(event *events.{{ .EventName }}Event) error {
	log.Printf("处理事件: %s", event.GetData())
	
	// 在这里添加事件处理逻辑
	
	return nil
}

{{ if .Queued }}
// ShouldQueue 是否应该排队处理
func (l *{{ .ListenerName }}) ShouldQueue() bool {
	return l.queued
}
{{ end }}
`

	// 解析模板
	tmpl, err := template.New("listener").Parse(listenerTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse listener template: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create listener file: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := map[string]interface{}{
		"ListenerName": listenerName,
		"EventName":    g.toPascalCase(eventName),
		"ProjectName":  g.getProjectName(),
		"Queued":       queued,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute listener template: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ 监听器 %s 生成完成: %s", listenerName, filePath))
	return nil
}

// getProjectName 获取项目名称
func (g *Generator) getProjectName() string {
	// 这里可以从 go.mod 文件读取项目名称
	// 暂时返回默认值
	return "laravel-go"
}

// GenerateRpcService 从 proto 文件生成 RPC 服务
func (g *Generator) GenerateRpcService(protoFile, serviceName, outputDir, goOut, goGrpcOut string) error {
	g.output.Info("🚀 开始生成 RPC 服务...")

	// 检查 proto 文件是否存在
	if _, err := os.Stat(protoFile); os.IsNotExist(err) {
		return fmt.Errorf("proto file not found: %s", protoFile)
	}

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 创建 Go 输出目录
	if err := os.MkdirAll(goOut, 0755); err != nil {
		return fmt.Errorf("failed to create go output directory: %w", err)
	}

	// 如果服务名为空，从proto文件中解析
	if serviceName == "" {
		serviceName = g.extractServiceNameFromProto(protoFile)
	}

	// 生成 protoc 命令
	cmd := exec.Command("protoc",
		"--go_out="+goOut,
		"--go_opt=paths=source_relative",
		"--go-grpc_out="+goGrpcOut,
		"--go-grpc_opt=paths=source_relative",
		protoFile)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run protoc: %w", err)
	}

	g.output.Success("✅ Protocol Buffers 代码生成完成")

	// 生成 RPC 服务结构
	if err := g.generateRpcServiceStructure(serviceName, outputDir); err != nil {
		return fmt.Errorf("failed to generate RPC service structure: %w", err)
	}

	g.output.Success("✅ RPC 服务生成完成")
	return nil
}

// extractServiceNameFromProto 从proto文件中提取服务名
func (g *Generator) extractServiceNameFromProto(protoFile string) string {
	file, err := os.Open(protoFile)
	if err != nil {
		return "service"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	serviceRegex := regexp.MustCompile(`service\s+(\w+)\s*{`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := serviceRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return "service"
}

// generateRpcServiceStructure 生成 RPC 服务结构
func (g *Generator) generateRpcServiceStructure(serviceName, outputDir string) error {
	// 创建服务目录结构
	dirs := []string{
		"internal/logic",
		"internal/svc",
		"internal/config",
		"internal/server",
		"etc",
	}

	for _, dir := range dirs {
		fullDir := filepath.Join(outputDir, dir)
		if err := os.MkdirAll(fullDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullDir, err)
		}
	}

	// 生成主服务文件
	mainFile := filepath.Join(outputDir, serviceName+".go")
	mainTemplate := `package main

import (
	"flag"
	"fmt"
	"log"

	"{{ .ProjectName }}/internal/config"
	"{{ .ProjectName }}/internal/server"
	"{{ .ProjectName }}/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/{{ .ServiceName }}.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	srv := server.New{{ .ServiceName }}Server(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		{{ .ServiceName }}.Register{{ .ServiceName }}Server(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
`

	data := map[string]interface{}{
		"ProjectName": g.getProjectName(),
		"ServiceName": g.toPascalCase(serviceName),
	}

	if err := g.writeTemplateToFile(mainFile, mainTemplate, data); err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	// 生成配置文件
	configFile := filepath.Join(outputDir, "etc", serviceName+".yaml")
	configTemplate := `Name: {{ .ServiceName }}
Host: 0.0.0.0
Port: 8080
Mode: dev

RpcServerConf:
  Endpoints:
    - 0.0.0.0:8080
  Timeout: 30000
`

	if err := g.writeTemplateToFile(configFile, configTemplate, data); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// 生成配置结构
	configStructFile := filepath.Join(outputDir, "internal/config/config.go")
	configStructTemplate := `package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mode string
}
`

	if err := g.writeTemplateToFile(configStructFile, configStructTemplate, data); err != nil {
		return fmt.Errorf("failed to write config struct file: %w", err)
	}

	// 生成服务上下文
	svcFile := filepath.Join(outputDir, "internal/svc/servicecontext.go")
	svcTemplate := `package svc

import (
	"{{ .ProjectName }}/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
`

	if err := g.writeTemplateToFile(svcFile, svcTemplate, data); err != nil {
		return fmt.Errorf("failed to write service context file: %w", err)
	}

	// 生成服务器接口
	serverFile := filepath.Join(outputDir, "internal/server/"+serviceName+"server.go")
	serverTemplate := `package server

import (
	"context"
	"{{ .ProjectName }}/internal/logic"
	"{{ .ProjectName }}/internal/svc"
	"{{ .ProjectName }}/types"
)

type {{ .ServiceName }}Server struct {
	svcCtx *svc.ServiceContext
	types.Unimplemented{{ .ServiceName }}Server
}

func New{{ .ServiceName }}Server(svcCtx *svc.ServiceContext) *{{ .ServiceName }}Server {
	return &{{ .ServiceName }}Server{
		svcCtx: svcCtx,
	}
}

// 在这里添加你的RPC方法实现
// 例如：
// func (s *{{ .ServiceName }}Server) YourMethod(ctx context.Context, req *types.YourRequest) (*types.YourResponse, error) {
//     l := logic.NewYourMethodLogic(ctx, s.svcCtx)
//     return l.YourMethod(req)
// }
`

	if err := g.writeTemplateToFile(serverFile, serverTemplate, data); err != nil {
		return fmt.Errorf("failed to write server file: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ RPC 服务结构生成完成: %s", outputDir))
	return nil
}

// writeTemplateToFile 将模板写入文件
func (g *Generator) writeTemplateToFile(filePath, templateStr string, data interface{}) error {
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// GenerateApiFromFile 从 .api 文件生成 API 服务
func (g *Generator) GenerateApiFromFile(apiFile, outputDir, handlerDir, logicDir, svcDir string) error {
	g.output.Info("🚀 开始从 .api 文件生成 API 服务...")

	// 检查 .api 文件是否存在
	if _, err := os.Stat(apiFile); os.IsNotExist(err) {
		return fmt.Errorf(".api file not found: %s", apiFile)
	}

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 创建目录结构
	dirs := []string{handlerDir, logicDir, svcDir, "internal/config", "etc"}
	for _, dir := range dirs {
		fullDir := filepath.Join(outputDir, dir)
		if err := os.MkdirAll(fullDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullDir, err)
		}
	}

	// 解析 .api 文件
	apiContent, err := os.ReadFile(apiFile)
	if err != nil {
		return fmt.Errorf("failed to read .api file: %w", err)
	}

	// 解析API文件内容
	apiInfo := g.parseApiFile(string(apiContent))

	// 生成 API 服务文件
	if err := g.generateApiServiceFiles(apiInfo, outputDir, handlerDir, logicDir, svcDir); err != nil {
		return fmt.Errorf("failed to generate API service files: %w", err)
	}

	g.output.Success("✅ API 服务生成完成")
	return nil
}

// ApiInfo API文件信息
type ApiInfo struct {
	Info     map[string]string
	Types    []TypeInfo
	Services []ServiceInfo
}

// TypeInfo 类型信息
type TypeInfo struct {
	Name   string
	Fields []FieldInfo
}

// FieldInfo 字段信息
type FieldInfo struct {
	Name string
	Type string
	Tag  string
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name    string
	Methods []MethodInfo
}

// MethodInfo 方法信息
type MethodInfo struct {
	Name   string
	Path   string
	Method string
	Req    string
	Resp   string
}

// parseApiFile 解析.api文件
func (g *Generator) parseApiFile(content string) *ApiInfo {
	apiInfo := &ApiInfo{
		Info:     make(map[string]string),
		Types:    []TypeInfo{},
		Services: []ServiceInfo{},
	}

	lines := strings.Split(content, "\n")
	var currentType *TypeInfo
	var currentService *ServiceInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// 解析info块
		if strings.HasPrefix(line, "info") {
			// 解析info信息
			continue
		}

		// 解析type块
		if strings.HasPrefix(line, "type") {
			if currentType != nil {
				apiInfo.Types = append(apiInfo.Types, *currentType)
			}
			typeName := strings.TrimSpace(strings.TrimPrefix(line, "type"))
			typeName = strings.TrimSpace(strings.TrimSuffix(typeName, "{"))
			currentType = &TypeInfo{Name: typeName, Fields: []FieldInfo{}}
			continue
		}

		// 解析service块
		if strings.HasPrefix(line, "service") {
			if currentService != nil {
				apiInfo.Services = append(apiInfo.Services, *currentService)
			}
			serviceName := strings.TrimSpace(strings.TrimPrefix(line, "service"))
			serviceName = strings.TrimSpace(strings.TrimSuffix(serviceName, "{"))
			currentService = &ServiceInfo{Name: serviceName, Methods: []MethodInfo{}}
			continue
		}

		// 解析字段
		if currentType != nil && strings.Contains(line, " ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				fieldName := parts[0]
				fieldType := parts[1]
				tag := ""
				if len(parts) > 2 {
					tag = strings.Join(parts[2:], " ")
				}
				currentType.Fields = append(currentType.Fields, FieldInfo{
					Name: fieldName,
					Type: fieldType,
					Tag:  tag,
				})
			}
		}

		// 解析方法
		if currentService != nil && strings.Contains(line, "(") && strings.Contains(line, ")") {
			// 简单的HTTP方法解析
			if strings.Contains(line, "get") || strings.Contains(line, "post") ||
				strings.Contains(line, "put") || strings.Contains(line, "delete") {
				// 这里可以添加更复杂的HTTP方法解析逻辑
			}
		}
	}

	// 添加最后一个type和service
	if currentType != nil {
		apiInfo.Types = append(apiInfo.Types, *currentType)
	}
	if currentService != nil {
		apiInfo.Services = append(apiInfo.Services, *currentService)
	}

	return apiInfo
}

// generateApiServiceFiles 生成 API 服务文件
func (g *Generator) generateApiServiceFiles(apiInfo *ApiInfo, outputDir, handlerDir, logicDir, svcDir string) error {
	// 生成主 API 文件
	mainFile := filepath.Join(outputDir, "main.go")
	mainTemplate := `package main

import (
	"flag"
	"fmt"

	"{{ .ProjectName }}/internal/config"
	"{{ .ProjectName }}/internal/handler"
	"{{ .ProjectName }}/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
`

	data := map[string]interface{}{
		"ProjectName": g.getProjectName(),
	}

	if err := g.writeTemplateToFile(mainFile, mainTemplate, data); err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	// 生成配置文件
	configFile := filepath.Join(outputDir, "etc", "api.yaml")
	configTemplate := `Name: api
Host: 0.0.0.0
Port: 8888
Mode: dev
`

	if err := g.writeTemplateToFile(configFile, configTemplate, data); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// 生成配置结构
	configStructFile := filepath.Join(outputDir, "internal/config/config.go")
	configStructTemplate := `package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
}
`

	if err := g.writeTemplateToFile(configStructFile, configStructTemplate, data); err != nil {
		return fmt.Errorf("failed to write config struct file: %w", err)
	}

	// 生成服务上下文
	svcFile := filepath.Join(outputDir, "internal/svc/servicecontext.go")
	svcTemplate := `package svc

import (
	"{{ .ProjectName }}/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
`

	if err := g.writeTemplateToFile(svcFile, svcTemplate, data); err != nil {
		return fmt.Errorf("failed to write service context file: %w", err)
	}

	// 生成类型定义
	if len(apiInfo.Types) > 0 {
		typesFile := filepath.Join(outputDir, "internal/types/types.go")
		typesTemplate := `package types

{{ range .Types }}
type {{ .Name }} struct {
{{ range .Fields }}	{{ .Name }} {{ .Type }} {{ .Tag }}
{{ end }}}
{{ end }}
`

		if err := g.writeTemplateToFile(typesFile, typesTemplate, apiInfo); err != nil {
			return fmt.Errorf("failed to write types file: %w", err)
		}
	}

	// 生成处理器
	handlerFile := filepath.Join(outputDir, "internal/handler/handlers.go")
	handlerTemplate := `package handler

import (
	"net/http"

	"{{ .ProjectName }}/internal/logic"
	"{{ .ProjectName }}/internal/svc"
	"{{ .ProjectName }}/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 在这里注册你的处理器
	// 例如：
	// server.AddRoutes(
	//     []rest.Route{
	//         {
	//             Method:  http.MethodGet,
	//             Path:    "/api/hello",
	//             Handler: HelloHandler(serverCtx),
	//         },
	//     },
	// )
}

// 示例处理器
func HelloHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HelloReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewHelloLogic(r.Context(), serverCtx)
		resp, err := l.Hello(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
`

	if err := g.writeTemplateToFile(handlerFile, handlerTemplate, data); err != nil {
		return fmt.Errorf("failed to write handler file: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ API 服务文件生成完成: %s", outputDir))
	return nil
}

// GenerateMicroservice 生成完整的微服务
func (g *Generator) GenerateMicroservice(name, serviceType, protoFile, apiFile, outputDir string) error {
	g.output.Info("🚀 开始生成微服务...")

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 根据服务类型生成不同的文件
	switch serviceType {
	case "rpc":
		if protoFile != "" {
			if err := g.GenerateRpcService(protoFile, name, outputDir, "./types", "./types"); err != nil {
				return fmt.Errorf("failed to generate RPC service: %w", err)
			}
		}
	case "api":
		if apiFile != "" {
			if err := g.GenerateApiFromFile(apiFile, outputDir, "./internal/handler", "./internal/logic", "./internal/svc"); err != nil {
				return fmt.Errorf("failed to generate API service: %w", err)
			}
		}
	case "both":
		if protoFile != "" {
			if err := g.GenerateRpcService(protoFile, name, outputDir, "./types", "./types"); err != nil {
				return fmt.Errorf("failed to generate RPC service: %w", err)
			}
		}
		if apiFile != "" {
			if err := g.GenerateApiFromFile(apiFile, outputDir, "./internal/handler", "./internal/logic", "./internal/svc"); err != nil {
				return fmt.Errorf("failed to generate API service: %w", err)
			}
		}
	}

	// 生成微服务配置文件
	if err := g.generateMicroserviceConfig(name, outputDir); err != nil {
		return fmt.Errorf("failed to generate microservice config: %w", err)
	}

	g.output.Success("✅ 微服务生成完成")
	return nil
}

// generateMicroserviceConfig 生成微服务配置
func (g *Generator) generateMicroserviceConfig(name, outputDir string) error {
	// 生成 Dockerfile
	dockerfile := filepath.Join(outputDir, "Dockerfile")
	dockerTemplate := `FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
`

	data := map[string]interface{}{
		"ServiceName": name,
	}

	if err := g.writeTemplateToFile(dockerfile, dockerTemplate, data); err != nil {
		return fmt.Errorf("failed to write dockerfile: %w", err)
	}

	// 生成 docker-compose.yml
	composeFile := filepath.Join(outputDir, "docker-compose.yml")
	composeTemplate := `version: '3.8'

services:
  {{ .ServiceName }}:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV=production
`

	if err := g.writeTemplateToFile(composeFile, composeTemplate, data); err != nil {
		return fmt.Errorf("failed to write docker-compose file: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ 微服务配置生成完成: %s", outputDir))
	return nil
}
