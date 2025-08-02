#!/bin/bash

# Laravel-Go Framework AI 编程助手快速设置脚本

set -e

echo "🚀 Laravel-Go Framework AI 编程助手设置开始..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 函数：打印带颜色的消息
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 函数：检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 检查必要的工具
print_message $BLUE "📋 检查必要的工具..."

if ! command_exists go; then
    print_message $RED "❌ Go 未安装，请先安装 Go"
    exit 1
fi

if ! command_exists git; then
    print_message $RED "❌ Git 未安装，请先安装 Git"
    exit 1
fi

print_message $GREEN "✅ 必要工具检查完成"

# 获取项目名称
read -p "请输入项目名称: " PROJECT_NAME
if [ -z "$PROJECT_NAME" ]; then
    print_message $RED "❌ 项目名称不能为空"
    exit 1
fi

# 创建项目目录
print_message $BLUE "📁 创建项目目录..."
mkdir -p "$PROJECT_NAME"
cd "$PROJECT_NAME"

# 复制AI助手配置文件
print_message $BLUE "📋 复制AI助手配置文件..."
mkdir -p .copilot .vscode

# 创建 .copilot/settings.json
cat > .copilot/settings.json << 'EOF'
{
  "framework": "laravel-go",
  "language": "go",
  "architecture": "layered",
  "patterns": [
    "repository",
    "service",
    "controller",
    "middleware",
    "validation"
  ],
  "conventions": {
    "naming": "snake_case",
    "file_structure": "feature_based",
    "error_handling": "wrapped_errors",
    "logging": "structured"
  },
  "code_style": {
    "indentation": "tabs",
    "line_length": 120,
    "imports": "grouped",
    "comments": "documentation"
  },
  "best_practices": [
    "use_framework_error_handling",
    "implement_layered_architecture",
    "use_dependency_injection",
    "write_unit_tests",
    "use_structured_logging",
    "implement_performance_monitoring"
  ],
  "templates": {
    "controller": "app/Http/Controllers/",
    "service": "app/Services/",
    "model": "app/Models/",
    "middleware": "app/Http/Middleware/",
    "repository": "app/Repositories/"
  }
}
EOF

# 创建 .vscode/settings.json
cat > .vscode/settings.json << 'EOF'
{
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.lintFlags": [
    "--fast"
  ],
  "go.testFlags": [
    "-v"
  ],
  "go.buildTags": "",
  "go.toolsManagement.checkForUpdates": "local",
  "go.useLanguageServer": true,
  "go.languageServerExperimentalFeatures": {
    "diagnostics": true,
    "documentLink": true
  },
  "files.associations": {
    "*.go": "go"
  },
  "emmet.includeLanguages": {
    "go": "html"
  },
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "go.template": {
    "framework": "laravel-go",
    "architecture": "layered"
  }
}
EOF

# 创建 .vscode/extensions.json
cat > .vscode/extensions.json << 'EOF'
{
  "recommendations": [
    "golang.go",
    "ms-vscode.go",
    "ms-vscode.vscode-json",
    "redhat.vscode-yaml",
    "ms-vscode.vscode-markdown",
    "GitHub.copilot",
    "GitHub.copilot-chat"
  ]
}
EOF

print_message $GREEN "✅ AI助手配置文件创建完成"

# 初始化Go模块
print_message $BLUE "🔧 初始化Go模块..."
go mod init "$PROJECT_NAME"

# 创建项目结构
print_message $BLUE "📁 创建项目结构..."
mkdir -p app/{Http/{Controllers,Middleware,Requests},Models,Services,Providers,Repositories}
mkdir -p config database/{migrations,seeders} routes storage/{cache,logs,uploads} tests
mkdir -p cmd docs examples

print_message $GREEN "✅ 项目结构创建完成"

# 创建基本的go.mod文件
print_message $BLUE "📦 设置依赖..."
cat > go.mod << EOF
module $PROJECT_NAME

go 1.21

require (
    laravel-go/framework v0.1.0
    github.com/gorilla/mux v1.8.0
    github.com/go-sql-driver/mysql v1.7.1
    github.com/redis/go-redis/v9 v9.0.5
    github.com/stretchr/testify v1.8.4
)

require (
    // 其他依赖将在运行时自动添加
)
EOF

# 创建 .gitignore
print_message $BLUE "📝 创建 .gitignore..."
cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs
*.log
logs/

# Environment files
.env
.env.local
.env.*.local

# Build output
build/
dist/

# Temporary files
tmp/
temp/
EOF

# 创建 README.md
print_message $BLUE "📖 创建 README.md..."
cat > README.md << EOF
# $PROJECT_NAME

基于 Laravel-Go Framework 的 Web 应用项目。

## 🚀 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 安装依赖

\`\`\`bash
go mod tidy
\`\`\`

### 运行项目

\`\`\`bash
go run cmd/main.go
\`\`\`

### 运行测试

\`\`\`bash
go test ./...
\`\`\`

## 📁 项目结构

\`\`\`
$PROJECT_NAME/
├── app/                    # 应用代码
│   ├── Http/              # HTTP 层
│   │   ├── Controllers/   # 控制器
│   │   ├── Middleware/    # 中间件
│   │   └── Requests/      # 请求验证
│   ├── Models/            # 数据模型
│   ├── Services/          # 业务服务
│   └── Providers/         # 服务提供者
├── config/                # 配置文件
├── database/              # 数据库相关
├── routes/                # 路由定义
├── storage/               # 文件存储
├── tests/                 # 测试文件
└── cmd/                   # 命令行入口
\`\`\`

## 🤖 AI 编程助手

本项目已配置 AI 编程助手支持，包括：

- GitHub Copilot 配置
- VS Code 设置
- 代码模板
- 最佳实践指南

详细配置请查看 \`.copilot/\` 目录。

## 📚 文档

- [框架文档](docs/)
- [API 文档](docs/api/)
- [最佳实践](docs/best-practices/)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
EOF

# 创建基本的 main.go
print_message $BLUE "🔧 创建基本的 main.go..."
mkdir -p cmd
cat > cmd/main.go << 'EOF'
package main

import (
	"context"
	"log"
	"net/http"
	"laravel-go/framework/errors"
	"laravel-go/framework/performance"
)

func main() {
	// 创建错误处理器
	logger := &CustomLogger{}
	errorHandler := errors.NewDefaultErrorHandler(logger)

	// 创建性能监控器
	monitor := performance.NewPerformanceMonitor()
	ctx := context.Background()
	monitor.Start(ctx)
	defer monitor.Stop()

	// 设置路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Welcome to Laravel-Go Framework!", "status": "running"}`))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy"}`))
	})

	// 启动服务器
	port := ":8080"
	log.Printf("🚀 服务器启动在端口 %s", port)
	log.Printf("📖 查看文档: http://localhost%s", port)
	log.Printf("🏥 健康检查: http://localhost%s/health", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

// CustomLogger 自定义日志器
type CustomLogger struct{}

func (l *CustomLogger) Error(message string, context map[string]interface{}) {
	log.Printf("[ERROR] %s: %+v", message, context)
}

func (l *CustomLogger) Warning(message string, context map[string]interface{}) {
	log.Printf("[WARN] %s: %+v", message, context)
}

func (l *CustomLogger) Info(message string, context map[string]interface{}) {
	log.Printf("[INFO] %s: %+v", message, context)
}

func (l *CustomLogger) Debug(message string, context map[string]interface{}) {
	log.Printf("[DEBUG] %s: %+v", message, context)
}
EOF

# 创建 AI 提示词文件
print_message $BLUE "📝 创建 AI 提示词文件..."
cat > AI_PROMPTS.md << 'EOF'
# AI 编程助手提示词

## 🚀 快速开始

在与 AI 助手对话时，使用以下提示词：

```
我正在使用 Laravel-Go Framework 开发项目。这是一个受 Laravel PHP 启发的 Go Web 框架。

请按照以下规范为我生成代码：

1. 使用框架的错误处理机制（errors.Wrap, SafeExecuteWithContext）
2. 遵循分层架构（Controller -> Service -> Repository）
3. 使用框架的中间件模式
4. 实现统一的错误处理
5. 集成性能监控
6. 使用结构化日志
7. 遵循 Go 编码规范

项目结构：
- app/Http/Controllers/ - 控制器层
- app/Services/ - 业务服务层
- app/Models/ - 数据模型
- framework/ - 框架核心

请为 [具体功能] 生成符合框架标准的代码。
```

## 📝 常用提示词模板

### 生成控制器
```
请为 [模块名称] 生成符合 Laravel-Go Framework 标准的控制器，包含 CRUD 操作和错误处理
```

### 生成服务层
```
请为 [模块名称] 生成符合 Laravel-Go Framework 标准的服务层，包含业务逻辑和缓存处理
```

### 生成模型
```
请为 [模块名称] 生成符合 Laravel-Go Framework 标准的数据模型，包含数据库操作和验证
```

## 📚 更多提示词

详细提示词模板请查看：https://github.com/your-repo/laravel-go/docs/ai-programming-assistant/ai-prompts.md
EOF

print_message $GREEN "✅ 项目初始化完成！"

# 显示后续步骤
print_message $YELLOW "
🎉 项目创建成功！

📋 后续步骤：

1. 进入项目目录：
   cd $PROJECT_NAME

2. 安装依赖：
   go mod tidy

3. 运行项目：
   go run cmd/main.go

4. 查看 AI 提示词：
   cat AI_PROMPTS.md

5. 配置 IDE：
   - 安装推荐的 VS Code 扩展
   - 启用 GitHub Copilot

📚 学习资源：
- 框架文档：docs/
- AI 提示词：AI_PROMPTS.md
- 最佳实践：docs/best-practices/

🤖 AI 助手使用：
- 复制 AI_PROMPTS.md 中的提示词
- 在 AI 助手中使用这些提示词
- 根据项目需求调整提示词

🚀 开始开发吧！
"

print_message $GREEN "✅ Laravel-Go Framework AI 编程助手设置完成！" 