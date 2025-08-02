#!/bin/bash

# Laravel-Go MCP 服务演示脚本

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 打印消息
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}  Laravel-Go MCP 服务演示${NC}"
    echo -e "${BLUE}================================${NC}"
}

# 检查服务是否运行
check_service() {
    if curl -s http://localhost:8080 > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# 发送 MCP 请求
send_mcp_request() {
    local method=$1
    local params=$2
    
    curl -s -X POST http://localhost:8080 \
        -H "Content-Type: application/json" \
        -d "{
            \"jsonrpc\": \"2.0\",
            \"id\": $(date +%s),
            \"method\": \"$method\",
            \"params\": $params
        }" | jq '.'
}

# 演示初始化项目
demo_initialize() {
    print_message "演示：初始化项目"
    
    local params='{
        "name": "demo-api",
        "description": "演示API项目",
        "version": "1.0.0",
        "modules": ["user", "product", "order"],
        "database": "mysql",
        "cache": "redis",
        "queue": "redis"
    }'
    
    send_mcp_request "initialize" "$params"
    echo ""
}

# 演示生成模块
demo_generate() {
    print_message "演示：生成模块"
    
    local params='{
        "type": "api",
        "name": "category"
    }'
    
    send_mcp_request "generate" "$params"
    echo ""
}

# 演示构建项目
demo_build() {
    print_message "演示：构建项目"
    
    send_mcp_request "build" "{}"
    echo ""
}

# 演示运行测试
demo_test() {
    print_message "演示：运行测试"
    
    send_mcp_request "test" "{}"
    echo ""
}

# 演示性能监控
demo_monitor() {
    print_message "演示：性能监控"
    
    send_mcp_request "monitor" "{}"
    echo ""
}

# 演示代码分析
demo_analyze() {
    print_message "演示：代码分析"
    
    send_mcp_request "analyze" "{}"
    echo ""
}

# 演示性能优化
demo_optimize() {
    print_message "演示：性能优化"
    
    send_mcp_request "optimize" "{}"
    echo ""
}

# 演示获取项目信息
demo_info() {
    print_message "演示：获取项目信息"
    
    send_mcp_request "info" "{}"
    echo ""
}

# 演示部署
demo_deploy() {
    print_message "演示：部署项目"
    
    local params='{
        "environment": "production"
    }'
    
    send_mcp_request "deploy" "$params"
    echo ""
}

# 主演示函数
main_demo() {
    print_header
    
    # 检查服务是否运行
    if ! check_service; then
        print_warning "MCP 服务未运行，请先启动服务："
        echo "  cd tools/laravel-go-mcp"
        echo "  go run main.go"
        echo ""
        print_warning "或者使用 Docker："
        echo "  ./deploy.sh start"
        echo ""
        exit 1
    fi
    
    print_message "MCP 服务运行正常，开始演示..."
    echo ""
    
    # 演示各个功能
    demo_initialize
    demo_generate
    demo_build
    demo_test
    demo_monitor
    demo_analyze
    demo_optimize
    demo_info
    demo_deploy
    
    echo ""
    print_message "演示完成！"
    echo ""
    echo "📝 更多信息请查看 README.md"
    echo "🔧 使用 'make help' 查看可用命令"
}

# 交互式演示
interactive_demo() {
    print_header
    print_message "交互式演示模式"
    echo ""
    
    while true; do
        echo "请选择要演示的功能："
        echo "1. 初始化项目"
        echo "2. 生成模块"
        echo "3. 构建项目"
        echo "4. 运行测试"
        echo "5. 性能监控"
        echo "6. 代码分析"
        echo "7. 性能优化"
        echo "8. 获取项目信息"
        echo "9. 部署项目"
        echo "0. 退出"
        echo ""
        read -p "请输入选项 (0-9): " choice
        
        case $choice in
            1) demo_initialize ;;
            2) demo_generate ;;
            3) demo_build ;;
            4) demo_test ;;
            5) demo_monitor ;;
            6) demo_analyze ;;
            7) demo_optimize ;;
            8) demo_info ;;
            9) demo_deploy ;;
            0) 
                print_message "退出演示"
                exit 0
                ;;
            *)
                print_warning "无效选项，请重新选择"
                ;;
        esac
        
        echo ""
        read -p "按回车键继续..."
        echo ""
    done
}

# 显示帮助
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -a, --auto      自动演示所有功能"
    echo "  -i, --interactive 交互式演示"
    echo "  -h, --help      显示帮助"
    echo ""
    echo "示例:"
    echo "  $0 -a           # 自动演示"
    echo "  $0 -i           # 交互式演示"
    echo ""
}

# 主函数
case "${1:-}" in
    "-a"|"--auto")
        main_demo
        ;;
    "-i"|"--interactive")
        interactive_demo
        ;;
    "-h"|"--help"|"")
        show_help
        ;;
    *)
        echo "未知选项: $1"
        show_help
        exit 1
        ;;
esac 