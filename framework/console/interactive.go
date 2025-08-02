package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ProjectConfig 项目配置
type ProjectConfig struct {
	Name          string
	Architecture  string // "monolithic" | "microservice"
	Database      string // "sqlite" | "mysql" | "postgresql"
	Cache         string // "memory" | "redis" | "memcached"
	Queue         string // "memory" | "redis" | "rabbitmq"
	Frontend      string // "api" | "blade" | "vue" | "react"
	Auth          string // "jwt" | "session" | "none"
	API           string // "rest" | "graphql" | "both"
	Testing       string // "unit" | "integration" | "both" | "none"
	Documentation string // "swagger" | "none"
	Monitoring    string // "prometheus" | "none"
	Logging       string // "file" | "json" | "both"

	// 框架核心功能
	Console              string // "basic" | "full" | "custom"
	Events               string // "none" | "basic" | "full"
	Validation           string // "none" | "basic" | "full"
	Middleware           string // "none" | "basic" | "full"
	Routing              string // "basic" | "advanced" | "full"
	Session              string // "none" | "file" | "redis" | "database"
	Mail                 string // "none" | "smtp" | "mailgun" | "sendgrid"
	Notifications        string // "none" | "database" | "mail" | "slack"
	FileStorage          string // "local" | "s3" | "oss" | "cos"
	Encryption           string // "none" | "basic" | "full"
	Hashing              string // "none" | "bcrypt" | "argon2"
	Pagination           string // "none" | "basic" | "advanced"
	RateLimiting         string // "none" | "basic" | "advanced"
	CORS                 string // "none" | "basic" | "full"
	Compression          string // "none" | "gzip" | "brotli"
	WebSockets           string // "none" | "basic" | "full"
	TaskScheduling       string // "none" | "basic" | "full"
	Timer                string // "none" | "cron" | "interval" | "full"
	HealthChecks         string // "none" | "basic" | "full"
	Metrics              string // "none" | "basic" | "prometheus"
	Profiling            string // "none" | "pprof" | "full"
	Internationalization string // "none" | "basic" | "full"
	Localization         string // "none" | "basic" | "full"
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
		"memcached - Memcached (高性能缓存)",
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
	output.Info("🔧 框架核心功能配置")
	output.Info("==================")

	// 控制台功能
	config.Console = askChoice("请选择控制台功能:", []string{
		"basic - 基础命令 (make:controller, make:model)",
		"full - 完整命令集 (包含所有生成器)",
		"custom - 自定义命令 (可扩展)",
	}, "full", output)

	// 事件系统
	config.Events = askChoice("请选择事件系统:", []string{
		"none - 无事件系统",
		"basic - 基础事件 (应用启动/关闭)",
		"full - 完整事件系统 (自定义事件)",
	}, "basic", output)

	// 数据验证
	config.Validation = askChoice("请选择数据验证:", []string{
		"none - 无验证",
		"basic - 基础验证 (必填、长度等)",
		"full - 完整验证 (自定义规则)",
	}, "basic", output)

	// 中间件
	config.Middleware = askChoice("请选择中间件:", []string{
		"none - 无中间件",
		"basic - 基础中间件 (日志、CORS)",
		"full - 完整中间件 (认证、限流等)",
	}, "basic", output)

	// 路由系统
	config.Routing = askChoice("请选择路由系统:", []string{
		"basic - 基础路由 (GET/POST)",
		"advanced - 高级路由 (参数、中间件)",
		"full - 完整路由 (分组、命名路由)",
	}, "advanced", output)

	// 会话管理
	config.Session = askChoice("请选择会话管理:", []string{
		"none - 无会话",
		"file - 文件会话 (开发环境)",
		"redis - Redis 会话 (生产环境)",
		"database - 数据库会话 (企业级)",
	}, "file", output)

	// 邮件系统
	config.Mail = askChoice("请选择邮件系统:", []string{
		"none - 无邮件功能",
		"smtp - SMTP 邮件 (传统)",
		"mailgun - Mailgun 服务",
		"sendgrid - SendGrid 服务",
	}, "none", output)

	// 通知系统
	config.Notifications = askChoice("请选择通知系统:", []string{
		"none - 无通知功能",
		"database - 数据库通知",
		"mail - 邮件通知",
		"slack - Slack 通知",
	}, "none", output)

	// 文件存储
	config.FileStorage = askChoice("请选择文件存储:", []string{
		"local - 本地存储 (开发环境)",
		"s3 - AWS S3 存储",
		"oss - 阿里云 OSS",
		"cos - 腾讯云 COS",
	}, "local", output)

	// 加密功能
	config.Encryption = askChoice("请选择加密功能:", []string{
		"none - 无加密",
		"basic - 基础加密 (AES)",
		"full - 完整加密 (多种算法)",
	}, "basic", output)

	// 密码哈希
	config.Hashing = askChoice("请选择密码哈希:", []string{
		"none - 无哈希",
		"bcrypt - Bcrypt 哈希",
		"argon2 - Argon2 哈希 (推荐)",
	}, "bcrypt", output)

	// 分页功能
	config.Pagination = askChoice("请选择分页功能:", []string{
		"none - 无分页",
		"basic - 基础分页",
		"advanced - 高级分页 (自定义)",
	}, "basic", output)

	// 限流功能
	config.RateLimiting = askChoice("请选择限流功能:", []string{
		"none - 无限流",
		"basic - 基础限流 (IP)",
		"advanced - 高级限流 (用户、API)",
	}, "basic", output)

	// CORS 支持
	config.CORS = askChoice("请选择 CORS 支持:", []string{
		"none - 无 CORS",
		"basic - 基础 CORS",
		"full - 完整 CORS (自定义)",
	}, "basic", output)

	// 压缩功能
	config.Compression = askChoice("请选择压缩功能:", []string{
		"none - 无压缩",
		"gzip - Gzip 压缩",
		"brotli - Brotli 压缩 (现代)",
	}, "gzip", output)

	// WebSocket 支持
	config.WebSockets = askChoice("请选择 WebSocket 支持:", []string{
		"none - 无 WebSocket",
		"basic - 基础 WebSocket",
		"full - 完整 WebSocket (房间、广播)",
	}, "none", output)

	// 任务调度
	config.TaskScheduling = askChoice("请选择任务调度:", []string{
		"none - 无任务调度",
		"basic - 基础调度 (定时任务)",
		"full - 完整调度 (复杂表达式)",
	}, "none", output)

	// 定时器
	config.Timer = askChoice("请选择定时器功能:", []string{
		"none - 无定时器",
		"cron - Cron 表达式定时器 (Unix 风格)",
		"interval - 间隔定时器 (简单间隔)",
		"full - 完整定时器 (Cron + 间隔 + 自定义)",
	}, "cron", output)

	// 健康检查
	config.HealthChecks = askChoice("请选择健康检查:", []string{
		"none - 无健康检查",
		"basic - 基础检查 (应用状态)",
		"full - 完整检查 (数据库、缓存等)",
	}, "basic", output)

	// 指标监控
	config.Metrics = askChoice("请选择指标监控:", []string{
		"none - 无监控",
		"basic - 基础指标 (请求数、响应时间)",
		"prometheus - Prometheus 指标",
	}, "basic", output)

	// 性能分析
	config.Profiling = askChoice("请选择性能分析:", []string{
		"none - 无分析",
		"pprof - Go pprof 分析",
		"full - 完整分析 (CPU、内存、GC)",
	}, "none", output)

	// 国际化
	config.Internationalization = askChoice("请选择国际化支持:", []string{
		"none - 无国际化",
		"basic - 基础国际化 (多语言)",
		"full - 完整国际化 (时区、货币)",
	}, "none", output)

	// 本地化
	config.Localization = askChoice("请选择本地化支持:", []string{
		"none - 无本地化",
		"basic - 基础本地化 (日期格式)",
		"full - 完整本地化 (数字、货币)",
	}, "none", output)

	output.Info("")
	output.Success("✅ 配置完成！")
	output.Info("")
	output.Info("📋 项目配置摘要:")
	output.Info("")
	output.Info("🏗️  基础架构:")
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
	output.Info("🔧 框架功能:")
	output.Info(fmt.Sprintf("   控制台: %s", config.Console))
	output.Info(fmt.Sprintf("   事件系统: %s", config.Events))
	output.Info(fmt.Sprintf("   数据验证: %s", config.Validation))
	output.Info(fmt.Sprintf("   中间件: %s", config.Middleware))
	output.Info(fmt.Sprintf("   路由系统: %s", config.Routing))
	output.Info(fmt.Sprintf("   会话管理: %s", config.Session))
	output.Info(fmt.Sprintf("   邮件系统: %s", config.Mail))
	output.Info(fmt.Sprintf("   通知系统: %s", config.Notifications))
	output.Info(fmt.Sprintf("   文件存储: %s", config.FileStorage))
	output.Info(fmt.Sprintf("   加密功能: %s", config.Encryption))
	output.Info(fmt.Sprintf("   密码哈希: %s", config.Hashing))
	output.Info(fmt.Sprintf("   分页功能: %s", config.Pagination))
	output.Info(fmt.Sprintf("   限流功能: %s", config.RateLimiting))
	output.Info(fmt.Sprintf("   CORS 支持: %s", config.CORS))
	output.Info(fmt.Sprintf("   压缩功能: %s", config.Compression))
	output.Info(fmt.Sprintf("   WebSocket: %s", config.WebSockets))
	output.Info(fmt.Sprintf("   任务调度: %s", config.TaskScheduling))
	output.Info(fmt.Sprintf("   定时器: %s", config.Timer))
	output.Info(fmt.Sprintf("   健康检查: %s", config.HealthChecks))
	output.Info(fmt.Sprintf("   指标监控: %s", config.Metrics))
	output.Info(fmt.Sprintf("   性能分析: %s", config.Profiling))
	output.Info(fmt.Sprintf("   国际化: %s", config.Internationalization))
	output.Info(fmt.Sprintf("   本地化: %s", config.Localization))
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
