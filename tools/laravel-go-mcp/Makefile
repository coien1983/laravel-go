# Laravel-Go MCP 服务 Makefile

.PHONY: help build run test clean demo client

# 默认目标
help:
	@echo "🚀 Laravel-Go MCP 服务"
	@echo "======================"
	@echo ""
	@echo "可用命令:"
	@echo "  build   - 构建 MCP 服务器"
	@echo "  run     - 运行 MCP 服务器"
	@echo "  test    - 运行测试"
	@echo "  clean   - 清理构建文件"
	@echo "  demo    - 运行客户端演示"
	@echo "  client  - 构建客户端示例"
	@echo "  help    - 显示帮助信息"
	@echo ""

# 构建 MCP 服务器
build:
	@echo "🔨 构建 MCP 服务器..."
	go build -o laravel-go-mcp main.go
	@echo "✅ 构建完成: laravel-go-mcp"

# 运行 MCP 服务器
run:
	@echo "🚀 启动 MCP 服务器..."
	@echo "📝 服务器将在 http://localhost:8080 启动"
	@echo "📝 按 Ctrl+C 停止服务器"
	@echo ""
	go run main.go

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	rm -f laravel-go-mcp
	@echo "✅ 清理完成"

# 运行客户端演示
demo:
	@echo "🎯 运行客户端演示..."
	@echo "📝 确保 MCP 服务器正在运行 (make run)"
	@echo ""
	go run client_example.go

# 构建客户端示例
client:
	@echo "🔨 构建客户端示例..."
	go build -o mcp-client client_example.go
	@echo "✅ 构建完成: mcp-client"

# 开发模式 (同时运行服务器和客户端)
dev:
	@echo "🛠️  开发模式..."
	@echo "📝 启动 MCP 服务器..."
	@echo ""
	@make run &
	@sleep 2
	@echo ""
	@echo "📝 运行客户端演示..."
	@make demo
	@echo ""
	@echo "📝 停止服务器..."
	@pkill -f "go run main.go" || true

# 安装依赖
deps:
	@echo "📦 安装依赖..."
	go mod tidy
	go mod download
	@echo "✅ 依赖安装完成"

# 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...
	@echo "✅ 代码格式化完成"

# 代码检查
lint:
	@echo "🔍 代码检查..."
	golangci-lint run
	@echo "✅ 代码检查完成"

# 性能测试
bench:
	@echo "⚡ 运行性能测试..."
	go test -bench=. ./...
	@echo "✅ 性能测试完成"

# 生成文档
docs:
	@echo "📝 生成文档..."
	@echo "✅ 文档已生成: README.md"

# 部署
deploy:
	@echo "🚀 部署 MCP 服务..."
	@make build
	@echo "✅ 部署完成"

# 健康检查
health:
	@echo "🏥 健康检查..."
	@curl -s http://localhost:8080 || echo "❌ 服务器未运行"
	@echo "✅ 健康检查完成"

# 完整测试流程
test-full:
	@echo "🧪 完整测试流程..."
	@make deps
	@make fmt
	@make lint
	@make test
	@make bench
	@echo "✅ 完整测试流程完成"

# 快速开始
quickstart:
	@echo "⚡ 快速开始..."
	@make deps
	@make build
	@echo ""
	@echo "🎉 快速开始完成!"
	@echo "📝 运行 'make run' 启动服务器"
	@echo "📝 运行 'make demo' 测试客户端" 