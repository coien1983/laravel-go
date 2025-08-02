#!/bin/bash

# Laravel-Go 部署脚本
# 支持 Docker 和 Kubernetes 部署

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
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 未安装，请先安装 $1"
        exit 1
    fi
}

# 显示帮助信息
show_help() {
    echo "Laravel-Go 部署脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help              显示帮助信息"
    echo "  -e, --env ENV           部署环境 (dev|prod)"
    echo "  -p, --platform PLATFORM 部署平台 (docker|k8s)"
    echo "  -b, --build             构建镜像"
    echo "  -d, --deploy            部署应用"
    echo "  -s, --stop              停止服务"
    echo "  -r, --restart           重启服务"
    echo "  -l, --logs              查看日志"
    echo "  -c, --clean             清理资源"
    echo ""
    echo "示例:"
    echo "  $0 -e prod -p docker -b -d    # 构建并部署到生产环境 Docker"
    echo "  $0 -e dev -p k8s -d           # 部署到开发环境 Kubernetes"
    echo "  $0 -p docker -s               # 停止 Docker 服务"
    echo "  $0 -p k8s -l                  # 查看 Kubernetes 日志"
}

# 构建 Docker 镜像
build_docker_image() {
    print_header "构建 Docker 镜像"
    
    check_command docker
    
    # 构建应用镜像
    print_message "构建 Laravel-Go 应用镜像..."
    docker build -f examples/deployment_example/docker/Dockerfile -t laravel-go-app:latest .
    
    print_message "Docker 镜像构建完成"
}

# 部署 Docker 服务
deploy_docker() {
    print_header "部署 Docker 服务"
    
    check_command docker
    check_command docker-compose
    
    local env=${ENV:-dev}
    
    print_message "部署环境: $env"
    
    # 切换到 Docker 目录
    cd examples/deployment_example/docker
    
    # 根据环境设置环境变量
    if [ "$env" = "prod" ]; then
        export APP_ENV=production
        export APP_DEBUG=false
    else
        export APP_ENV=development
        export APP_DEBUG=true
    fi
    
    # 启动服务
    print_message "启动 Docker Compose 服务..."
    docker-compose up -d
    
    print_message "Docker 服务部署完成"
    print_message "应用地址: http://localhost"
    print_message "API 地址: http://localhost/api"
    print_message "健康检查: http://localhost/health"
    print_message "监控面板: http://localhost:3000 (Grafana)"
    print_message "指标服务: http://localhost:9090 (Prometheus)"
}

# 停止 Docker 服务
stop_docker() {
    print_header "停止 Docker 服务"
    
    cd examples/deployment_example/docker
    docker-compose down
    
    print_message "Docker 服务已停止"
}

# 重启 Docker 服务
restart_docker() {
    print_header "重启 Docker 服务"
    
    stop_docker
    deploy_docker
}

# 查看 Docker 日志
logs_docker() {
    print_header "查看 Docker 日志"
    
    cd examples/deployment_example/docker
    docker-compose logs -f
}

# 清理 Docker 资源
clean_docker() {
    print_header "清理 Docker 资源"
    
    cd examples/deployment_example/docker
    docker-compose down -v --remove-orphans
    docker system prune -f
    
    print_message "Docker 资源清理完成"
}

# 部署 Kubernetes 服务
deploy_k8s() {
    print_header "部署 Kubernetes 服务"
    
    check_command kubectl
    
    local env=${ENV:-dev}
    
    print_message "部署环境: $env"
    
    # 切换到 Kubernetes 目录
    cd examples/deployment_example/kubernetes
    
    # 创建命名空间
    print_message "创建命名空间..."
    kubectl create namespace laravel-go --dry-run=client -o yaml | kubectl apply -f -
    
    # 应用配置
    print_message "应用 Kubernetes 配置..."
    kubectl apply -f deployment.yaml
    kubectl apply -f monitoring.yaml
    
    # 等待部署完成
    print_message "等待部署完成..."
    kubectl wait --for=condition=available --timeout=300s deployment/laravel-go-app -n laravel-go
    
    print_message "Kubernetes 服务部署完成"
    print_message "查看服务状态: kubectl get all -n laravel-go"
    print_message "查看日志: kubectl logs -f deployment/laravel-go-app -n laravel-go"
}

# 停止 Kubernetes 服务
stop_k8s() {
    print_header "停止 Kubernetes 服务"
    
    cd examples/deployment_example/kubernetes
    kubectl delete -f monitoring.yaml --ignore-not-found=true
    kubectl delete -f deployment.yaml --ignore-not-found=true
    
    print_message "Kubernetes 服务已停止"
}

# 重启 Kubernetes 服务
restart_k8s() {
    print_header "重启 Kubernetes 服务"
    
    stop_k8s
    deploy_k8s
}

# 查看 Kubernetes 日志
logs_k8s() {
    print_header "查看 Kubernetes 日志"
    
    kubectl logs -f deployment/laravel-go-app -n laravel-go
}

# 清理 Kubernetes 资源
clean_k8s() {
    print_header "清理 Kubernetes 资源"
    
    kubectl delete namespace laravel-go --ignore-not-found=true
    
    print_message "Kubernetes 资源清理完成"
}

# 主函数
main() {
    # 默认值
    ENV="dev"
    PLATFORM="docker"
    BUILD=false
    DEPLOY=false
    STOP=false
    RESTART=false
    LOGS=false
    CLEAN=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -e|--env)
                ENV="$2"
                shift 2
                ;;
            -p|--platform)
                PLATFORM="$2"
                shift 2
                ;;
            -b|--build)
                BUILD=true
                shift
                ;;
            -d|--deploy)
                DEPLOY=true
                shift
                ;;
            -s|--stop)
                STOP=true
                shift
                ;;
            -r|--restart)
                RESTART=true
                shift
                ;;
            -l|--logs)
                LOGS=true
                shift
                ;;
            -c|--clean)
                CLEAN=true
                shift
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 验证环境
    if [[ "$ENV" != "dev" && "$ENV" != "prod" ]]; then
        print_error "无效的环境: $ENV (支持: dev, prod)"
        exit 1
    fi
    
    # 验证平台
    if [[ "$PLATFORM" != "docker" && "$PLATFORM" != "k8s" ]]; then
        print_error "无效的平台: $PLATFORM (支持: docker, k8s)"
        exit 1
    fi
    
    # 执行操作
    if [ "$BUILD" = true ]; then
        if [ "$PLATFORM" = "docker" ]; then
            build_docker_image
        fi
    fi
    
    if [ "$DEPLOY" = true ]; then
        if [ "$PLATFORM" = "docker" ]; then
            deploy_docker
        elif [ "$PLATFORM" = "k8s" ]; then
            deploy_k8s
        fi
    fi
    
    if [ "$STOP" = true ]; then
        if [ "$PLATFORM" = "docker" ]; then
            stop_docker
        elif [ "$PLATFORM" = "k8s" ]; then
            stop_k8s
        fi
    fi
    
    if [ "$RESTART" = true ]; then
        if [ "$PLATFORM" = "docker" ]; then
            restart_docker
        elif [ "$PLATFORM" = "k8s" ]; then
            restart_k8s
        fi
    fi
    
    if [ "$LOGS" = true ]; then
        if [ "$PLATFORM" = "docker" ]; then
            logs_docker
        elif [ "$PLATFORM" = "k8s" ]; then
            logs_k8s
        fi
    fi
    
    if [ "$CLEAN" = true ]; then
        if [ "$PLATFORM" = "docker" ]; then
            clean_docker
        elif [ "$PLATFORM" = "k8s" ]; then
            clean_k8s
        fi
    fi
    
    # 如果没有指定任何操作，显示帮助
    if [ "$BUILD" = false ] && [ "$DEPLOY" = false ] && [ "$STOP" = false ] && [ "$RESTART" = false ] && [ "$LOGS" = false ] && [ "$CLEAN" = false ]; then
        show_help
    fi
}

# 运行主函数
main "$@" 