package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ProjectConfig 项目配置
type ProjectConfig struct {
	Name           string
	Architecture   string // "monolithic" | "microservice"
	Database       string // "sqlite" | "mysql" | "postgresql"
	Cache          string // "memory" | "redis"
	Queue          string // "memory" | "redis" | "rabbitmq"
	Frontend       string // "api" | "blade" | "vue" | "react"
	Auth           string // "jwt" | "session" | "none"
	API            string // "rest" | "graphql" | "both"
	Testing        string // "unit" | "integration" | "both" | "none"
	Documentation  string // "swagger" | "none"
	Monitoring     string // "prometheus" | "none"
	Logging        string // "file" | "json" | "both"
}

// InteractiveConfig 交互式配置
func InteractiveConfig(projectName string, output Output) *ProjectConfig {
	config := &ProjectConfig{
		Name: projectName,
	}

	output.Info("🚀 Laravel-Go 项目初始化向导")
	output.Info("================================")
	output.Info("")

	// 架构选择
	config.Architecture = askChoice("请选择项目架构:", []string{
		"monolithic - 单体应用 (推荐新手)",
		"microservice - 微服务架构 (适合大型项目)",
	}, "monolithic", output)

	// 数据库选择
	config.Database = askChoice("请选择数据库:", []string{
		"sqlite - SQLite (开发环境推荐)",
		"mysql - MySQL (生产环境常用)",
		"postgresql - PostgreSQL (企业级应用)",
	}, "sqlite", output)

	// 缓存选择
	config.Cache = askChoice("请选择缓存系统:", []string{
		"memory - 内存缓存 (开发环境)",
		"redis - Redis (生产环境推荐)",
	}, "memory", output)

	// 队列选择
	config.Queue = askChoice("请选择队列系统:", []string{
		"memory - 内存队列 (开发环境)",
		"redis - Redis 队列 (生产环境)",
		"rabbitmq - RabbitMQ (企业级)",
	}, "memory", output)

	// 前端选择
	config.Frontend = askChoice("请选择前端方案:", []string{
		"api - 纯 API 服务 (前后端分离)",
		"blade - Blade 模板 (传统 MVC)",
		"vue - Vue.js 集成 (现代前端)",
		"react - React 集成 (现代前端)",
	}, "api", output)

	// 认证选择
	config.Auth = askChoice("请选择认证方式:", []string{
		"none - 无认证 (简单应用)",
		"jwt - JWT 认证 (API 服务推荐)",
		"session - Session 认证 (传统 Web)",
	}, "jwt", output)

	// API 类型选择
	config.API = askChoice("请选择 API 类型:", []string{
		"rest - REST API (传统)",
		"graphql - GraphQL (现代)",
		"both - 同时支持 REST 和 GraphQL",
	}, "rest", output)

	// 测试选择
	config.Testing = askChoice("请选择测试策略:", []string{
		"none - 无测试 (快速原型)",
		"unit - 单元测试 (基础)",
		"integration - 集成测试 (推荐)",
		"both - 单元 + 集成测试 (完整)",
	}, "integration", output)

	// 文档选择
	config.Documentation = askChoice("请选择 API 文档:", []string{
		"none - 无文档",
		"swagger - Swagger/OpenAPI 文档",
	}, "swagger", output)

	// 监控选择
	config.Monitoring = askChoice("请选择监控方案:", []string{
		"none - 无监控",
		"prometheus - Prometheus 监控",
	}, "none", output)

	// 日志选择
	config.Logging = askChoice("请选择日志方案:", []string{
		"file - 文件日志 (简单)",
		"json - JSON 格式日志 (结构化)",
		"both - 文件 + JSON 日志 (完整)",
	}, "file", output)

	output.Info("")
	output.Success("✅ 配置完成！")
	output.Info("")
	output.Info("📋 项目配置摘要:")
	output.Info(fmt.Sprintf("   项目名称: %s", config.Name))
	output.Info(fmt.Sprintf("   架构模式: %s", config.Architecture))
	output.Info(fmt.Sprintf("   数据库: %s", config.Database))
	output.Info(fmt.Sprintf("   缓存: %s", config.Cache))
	output.Info(fmt.Sprintf("   队列: %s", config.Queue))
	output.Info(fmt.Sprintf("   前端: %s", config.Frontend))
	output.Info(fmt.Sprintf("   认证: %s", config.Auth))
	output.Info(fmt.Sprintf("   API: %s", config.API))
	output.Info(fmt.Sprintf("   测试: %s", config.Testing))
	output.Info(fmt.Sprintf("   文档: %s", config.Documentation))
	output.Info(fmt.Sprintf("   监控: %s", config.Monitoring))
	output.Info(fmt.Sprintf("   日志: %s", config.Logging))
	output.Info("")

	return config
}

// askChoice 询问用户选择
func askChoice(question string, options []string, defaultChoice string, output Output) string {
	output.Info(question)
	for i, option := range options {
		output.Info(fmt.Sprintf("  %d. %s", i+1, option))
	}
	
	defaultIndex := 0
	for i, option := range options {
		if strings.Contains(option, defaultChoice) {
			defaultIndex = i + 1
			break
		}
	}
	
	output.Info(fmt.Sprintf("请选择 (默认: %d): ", defaultIndex))
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "" {
		// 使用默认值
		for _, option := range options {
			if strings.Contains(option, defaultChoice) {
				return defaultChoice
			}
		}
		return defaultChoice
	}
	
	// 解析用户输入
	var choice int
	fmt.Sscanf(input, "%d", &choice)
	
	if choice > 0 && choice <= len(options) {
		selected := options[choice-1]
		// 提取选择的值
		parts := strings.Split(selected, " - ")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
		return selected
	}
	
	// 无效输入，使用默认值
	return defaultChoice
}

// askYesNo 询问是/否问题
func askYesNo(question string, defaultYes bool, output Output) bool {
	defaultText := "Y/n"
	if !defaultYes {
		defaultText = "y/N"
	}
	
	output.Info(fmt.Sprintf("%s (%s): ", question, defaultText))
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	
	if input == "" {
		return defaultYes
	}
	
	return input == "y" || input == "yes"
} 