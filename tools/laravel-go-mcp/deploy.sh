#!/bin/bash

# Laravel-Go MCP 服务部署脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}  Laravel-Go MCP 服务部署${NC}"
    echo -e "${BLUE}================================${NC}"
}

# 检查依赖
check_dependencies() {
    print_message "检查依赖..."
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    # 检查 Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        print_warning "Go 未安装，将使用 Docker 构建"
    fi
    
    print_message "依赖检查完成"
}

# 创建必要的目录
create_directories() {
    print_message "创建必要的目录..."
    
    mkdir -p projects
    mkdir -p logs
    mkdir -p nginx/ssl
    mkdir -p mysql/init
    mkdir -p monitoring/grafana/dashboards
    mkdir -p monitoring/grafana/datasources
    
    print_message "目录创建完成"
}

# 构建应用
build_app() {
    print_message "构建应用..."
    
    if command -v go &> /dev/null; then
        print_message "使用 Go 构建..."
        go mod tidy
        go build -o laravel-go-mcp main.go
        print_message "Go 构建完成"
    else
        print_message "使用 Docker 构建..."
        docker build -t laravel-go-mcp .
        print_message "Docker 构建完成"
    fi
}

# 启动服务
start_services() {
    print_message "启动服务..."
    
    # 检查是否已经运行
    if docker-compose ps | grep -q "laravel-go-mcp"; then
        print_warning "服务已经在运行，正在重启..."
        docker-compose down
    fi
    
    # 启动服务
    docker-compose up -d
    
    print_message "服务启动完成"
}

# 检查服务状态
check_status() {
    print_message "检查服务状态..."
    
    sleep 5
    
    # 检查 MCP 服务
    if curl -s http://localhost:8080 > /dev/null; then
        print_message "✅ MCP 服务运行正常"
    else
        print_error "❌ MCP 服务启动失败"
        return 1
    fi
    
    # 检查 Redis
    if docker-compose exec redis redis-cli ping | grep -q "PONG"; then
        print_message "✅ Redis 服务运行正常"
    else
        print_warning "⚠️  Redis 服务可能有问题"
    fi
    
    # 检查 MySQL
    if docker-compose exec mysql mysqladmin ping -h localhost -u root -proot > /dev/null 2>&1; then
        print_message "✅ MySQL 服务运行正常"
    else
        print_warning "⚠️  MySQL 服务可能有问题"
    fi
    
    print_message "状态检查完成"
}

# 运行测试
run_tests() {
    print_message "运行测试..."
    
    # 等待服务完全启动
    sleep 10
    
    # 运行客户端测试
    if command -v go &> /dev/null; then
        print_message "运行客户端测试..."
        go run client_example.go || print_warning "客户端测试失败"
    else
        print_warning "跳过客户端测试 (Go 未安装)"
    fi
    
    print_message "测试完成"
}

# 显示服务信息
show_info() {
    print_message "服务信息:"
    echo ""
    echo "🌐 MCP 服务: http://localhost:8080"
    echo "🗄️  MySQL: localhost:3306"
    echo "💾 Redis: localhost:6379"
    echo "📊 Prometheus: http://localhost:9090"
    echo "📈 Grafana: http://localhost:3000 (admin/admin)"
    echo ""
    echo "📝 查看日志: docker-compose logs -f"
    echo "🛑 停止服务: docker-compose down"
    echo "🔄 重启服务: docker-compose restart"
}

# 停止服务
stop_services() {
    print_message "停止服务..."
    docker-compose down
    print_message "服务已停止"
}

# 清理资源
cleanup() {
    print_message "清理资源..."
    docker-compose down -v
    docker system prune -f
    print_message "清理完成"
}

# 显示帮助
show_help() {
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  deploy     - 部署服务 (默认)"
    echo "  start      - 启动服务"
    echo "  stop       - 停止服务"
    echo "  restart    - 重启服务"
    echo "  status     - 检查服务状态"
    echo "  test       - 运行测试"
    echo "  logs       - 查看日志"
    echo "  cleanup    - 清理资源"
    echo "  help       - 显示帮助"
    echo ""
}

# 查看日志
show_logs() {
    print_message "查看服务日志..."
    docker-compose logs -f
}

# 主函数
main() {
    case "${1:-deploy}" in
        "deploy")
            print_header
            check_dependencies
            create_directories
            build_app
            start_services
            check_status
            run_tests
            show_info
            ;;
        "start")
            start_services
            check_status
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            stop_services
            start_services
            check_status
            ;;
        "status")
            check_status
            ;;
        "test")
            run_tests
            ;;
        "logs")
            show_logs
            ;;
        "cleanup")
            cleanup
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@" 