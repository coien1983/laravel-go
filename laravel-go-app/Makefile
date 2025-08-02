# Laravel-Go Project Makefile

.PHONY: help
help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## 运行应用
	go run main.go

.PHONY: build
build: ## 构建应用
	go build -o bin/app main.go

.PHONY: test
test: ## 运行测试
	go test ./...

.PHONY: clean
clean: ## 清理构建文件
	rm -rf bin/
	rm -f *.db

.PHONY: deps
deps: ## 安装依赖
	go mod tidy
	go mod download

.PHONY: dev
dev: ## 开发模式运行
	APP_ENV=local go run main.go

.PHONY: prod
prod: ## 生产模式运行
	APP_ENV=production go run main.go

# =============================================================================
# Docker 操作
# =============================================================================

.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	docker build -t laravel-go-app .

.PHONY: docker-run
docker-run: ## 运行 Docker 容器
	docker run -p 8080:8080 laravel-go-app

.PHONY: docker-compose-up
docker-compose-up: ## 启动 Docker Compose 服务
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## 停止 Docker Compose 服务
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## 查看 Docker Compose 日志
	docker-compose logs -f

.PHONY: docker-clean
docker-clean: ## 清理 Docker 资源
	docker-compose down -v --remove-orphans
	docker system prune -f

# =============================================================================
# Kubernetes 操作
# =============================================================================

.PHONY: k8s-apply
k8s-apply: ## 部署到 Kubernetes
	kubectl apply -f k8s/

.PHONY: k8s-delete
k8s-delete: ## 从 Kubernetes 删除
	kubectl delete -f k8s/

.PHONY: k8s-status
k8s-status: ## 查看 Kubernetes 部署状态
	kubectl get pods,services,ingress

.PHONY: k8s-logs
k8s-logs: ## 查看 Kubernetes 日志
	kubectl logs -f deployment/laravel-go-app

.PHONY: k8s-clean
k8s-clean: ## 清理 Kubernetes 资源
	kubectl delete -f k8s/ --ignore-not-found=true