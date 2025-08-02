#!/bin/bash

# Laravel-Go Framework 安装脚本
# 用于安装 largo 命令行工具

set -e

echo "🚀 Laravel-Go Framework 安装脚本"
echo "=================================="

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到 Go 环境，请先安装 Go"
    exit 1
fi

echo "✅ Go 环境检查通过"

# 获取 Go bin 目录
GOBIN=$(go env GOPATH)/bin
echo "📁 Go bin 目录: $GOBIN"

# 构建 largo 命令
echo "🔨 构建 largo 命令..."
go build -o bin/largo cmd/artisan/main.go

if [ ! -f "bin/largo" ]; then
    echo "❌ 构建失败"
    exit 1
fi

echo "✅ 构建完成"

# 安装到 Go bin 目录
echo "📦 安装 largo 到 $GOBIN..."
cp bin/largo "$GOBIN/"

if [ -f "$GOBIN/largo" ]; then
    echo "✅ 安装成功！"
    echo ""
    echo "🎉 largo 命令已安装到: $GOBIN/largo"
    echo "现在可以在任何地方使用 'largo' 命令"
    echo ""
    echo "📋 使用示例:"
    echo "  largo --version          # 查看版本"
    echo "  largo list               # 查看所有命令"
    echo "  largo init               # 初始化新项目"
    echo "  largo make:controller    # 创建控制器"
    echo "  largo make:model         # 创建模型"
    echo ""
    echo "📖 更多信息请查看 README.md"
else
    echo "❌ 安装失败"
    exit 1
fi 