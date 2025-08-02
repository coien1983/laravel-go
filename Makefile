# Laravel-Go Framework Makefile
# 提供便捷的命令行工具操作

# 变量定义
ARTISAN := go run cmd/artisan/main.go
APP_NAME := laravel-go-app
PORT := 8080
NAMESPACE := default
REPLICAS := 3

# 默认目标
.PHONY: help
help: ## 显示帮助信息
	@echo "Laravel-Go Framework 可用命令:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# 项目初始化
# =============================================================================

.PHONY: init
init: ## 初始化新项目
	$(ARTISAN) init

.PHONY: init-custom
init-custom: ## 使用自定义名称初始化项目
	@read -p "请输入项目名称: " name; \
	$(ARTISAN) init --name=$$name

# =============================================================================
# 代码生成
# =============================================================================

.PHONY: controller
controller: ## 生成控制器 (需要指定名称)
	@read -p "请输入控制器名称: " name; \
	$(ARTISAN) make:controller $$name

.PHONY: controller-custom
controller-custom: ## 生成控制器 (指定名称和命名空间)
	@read -p "请输入控制器名称: " name; \
	read -p "请输入命名空间 (默认: app): " namespace; \
	$(ARTISAN) make:controller $$name --namespace=$${namespace:-app}

.PHONY: model
model: ## 生成模型 (需要指定名称)
	@read -p "请输入模型名称: " name; \
	$(ARTISAN) make:model $$name

.PHONY: model-fields
model-fields: ## 生成模型 (指定名称和字段)
	@read -p "请输入模型名称: " name; \
	read -p "请输入字段 (格式: name:string,email:string,age:int): " fields; \
	$(ARTISAN) make:model $$name --fields=$$fields

.PHONY: middleware
middleware: ## 生成中间件 (需要指定名称)
	@read -p "请输入中间件名称: " name; \
	$(ARTISAN) make:middleware $$name

.PHONY: migration
migration: ## 生成迁移文件 (需要指定名称)
	@read -p "请输入迁移名称: " name; \
	$(ARTISAN) make:migration $$name

.PHONY: migration-table
migration-table: ## 生成迁移文件 (指定名称和表名)
	@read -p "请输入迁移名称: " name; \
	read -p "请输入表名: " table; \
	$(ARTISAN) make:migration $$name --table=$$table

.PHONY: migration-fields
migration-fields: ## 生成迁移文件 (指定名称、表名和字段)
	@read -p "请输入迁移名称: " name; \
	read -p "请输入表名: " table; \
	read -p "请输入字段 (格式: name:string,email:string,age:int): " fields; \
	$(ARTISAN) make:migration $$name --table=$$table --fields=$$fields

.PHONY: test
test: ## 生成测试文件 (需要指定名称)
	@read -p "请输入测试名称: " name; \
	$(ARTISAN) make:test $$name

.PHONY: test-type
test-type: ## 生成测试文件 (指定名称和类型)
	@read -p "请输入测试名称: " name; \
	read -p "请输入测试类型 (unit/integration): " type; \
	$(ARTISAN) make:test $$name --type=$$type

# =============================================================================
# 部署配置生成
# =============================================================================

.PHONY: docker
docker: ## 生成 Docker 配置文件 (使用默认配置)
	$(ARTISAN) make:docker

.PHONY: docker-custom
docker-custom: ## 生成 Docker 配置文件 (自定义配置)
	@read -p "请输入应用名称 (默认: $(APP_NAME)): " name; \
	read -p "请输入端口 (默认: $(PORT)): " port; \
	read -p "请输入环境 (development/production, 默认: development): " env; \
	$(ARTISAN) make:docker --name=$${name:-$(APP_NAME)} --port=$${port:-$(PORT)} --env=$${env:-development}

.PHONY: k8s
k8s: ## 生成 Kubernetes 配置文件 (使用默认配置)
	$(ARTISAN) make:k8s

.PHONY: k8s-custom
k8s-custom: ## 生成 Kubernetes 配置文件 (自定义配置)
	@read -p "请输入应用名称 (默认: $(APP_NAME)): " name; \
	read -p "请输入副本数量 (默认: $(REPLICAS)): " replicas; \
	read -p "请输入端口 (默认: $(PORT)): " port; \
	read -p "请输入命名空间 (默认: $(NAMESPACE)): " namespace; \
	$(ARTISAN) make:k8s --name=$${name:-$(APP_NAME)} --replicas=$${replicas:-$(REPLICAS)} --port=$${port:-$(PORT)} --namespace=$${namespace:-$(NAMESPACE)}

# =============================================================================
# 项目维护
# =============================================================================

.PHONY: routes
routes: ## 列出所有路由
	$(ARTISAN) route:list

.PHONY: cache-clear
cache-clear: ## 清除应用缓存
	$(ARTISAN) cache:clear

.PHONY: list
list: ## 列出所有可用命令
	$(ARTISAN) list

.PHONY: version
version: ## 显示版本信息
	$(ARTISAN) --version

# =============================================================================
# 快速生成常用组件
# =============================================================================

.PHONY: api
api: ## 快速生成 API 控制器和模型
	@read -p "请输入资源名称 (如: user): " name; \
	echo "生成 $$name 的 API 组件..."; \
	$(ARTISAN) make:controller $$name --namespace=api; \
	$(ARTISAN) make:model $$name; \
	$(ARTISAN) make:migration create_$${name}s_table --table=$${name}s; \
	echo "✅ $$name API 组件生成完成!"

.PHONY: crud
crud: ## 快速生成完整的 CRUD 组件
	@read -p "请输入资源名称 (如: user): " name; \
	echo "生成 $$name 的完整 CRUD 组件..."; \
	$(ARTISAN) make:controller $$name --namespace=app; \
	$(ARTISAN) make:model $$name; \
	$(ARTISAN) make:migration create_$${name}s_table --table=$${name}s; \
	$(ARTISAN) make:test $$name --type=unit; \
	$(ARTISAN) make:test $$name --type=integration; \
	echo "✅ $$name CRUD 组件生成完成!"

# =============================================================================
# 开发工具
# =============================================================================

.PHONY: dev-setup
dev-setup: ## 设置开发环境
	@echo "设置 Laravel-Go 开发环境..."
	go mod tidy
	go mod download
	@echo "✅ 开发环境设置完成!"

.PHONY: build
build: ## 构建应用
	@echo "构建 Laravel-Go 应用..."
	go build -o bin/laravel-go .
	@echo "✅ 应用构建完成: bin/laravel-go"

.PHONY: run
run: ## 运行应用
	@echo "启动 Laravel-Go 应用..."
	$(ARTISAN)

.PHONY: test-all
test-all: ## 运行所有测试
	@echo "运行所有测试..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率报告..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告生成完成: coverage.html"

.PHONY: lint
lint: ## 代码检查
	@echo "运行代码检查..."
	golangci-lint run

.PHONY: fmt
fmt: ## 格式化代码
	@echo "格式化代码..."
	go fmt ./...

.PHONY: vet
vet: ## 代码静态分析
	@echo "运行代码静态分析..."
	go vet ./...

# =============================================================================
# Docker 操作
# =============================================================================

.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	@read -p "请输入镜像名称 (默认: $(APP_NAME)): " name; \
	docker build -t $${name:-$(APP_NAME)} .

.PHONY: docker-run
docker-run: ## 运行 Docker 容器
	@read -p "请输入镜像名称 (默认: $(APP_NAME)): " name; \
	read -p "请输入端口映射 (默认: $(PORT):$(PORT)): " port; \
	docker run -p $${port:-$(PORT):$(PORT)} $${name:-$(APP_NAME)}

.PHONY: docker-compose-up
docker-compose-up: ## 启动 Docker Compose 服务
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## 停止 Docker Compose 服务
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## 查看 Docker Compose 日志
	docker-compose logs -f

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
	@read -p "请输入 Pod 名称: " pod; \
	kubectl logs -f $$pod

# =============================================================================
# 清理操作
# =============================================================================

.PHONY: clean
clean: ## 清理构建文件
	@echo "清理构建文件..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "✅ 清理完成!"

.PHONY: clean-docker
clean-docker: ## 清理 Docker 文件
	@echo "清理 Docker 文件..."
	rm -f Dockerfile docker-compose.yml .dockerignore
	@echo "✅ Docker 文件清理完成!"

.PHONY: clean-k8s
clean-k8s: ## 清理 Kubernetes 文件
	@echo "清理 Kubernetes 文件..."
	rm -rf k8s/
	@echo "✅ Kubernetes 文件清理完成!"

.PHONY: clean-all
clean-all: clean clean-docker clean-k8s ## 清理所有生成的文件

# =============================================================================
# 项目信息
# =============================================================================

.PHONY: info
info: ## 显示项目信息
	@echo "Laravel-Go Framework 项目信息:"
	@echo "  应用名称: $(APP_NAME)"
	@echo "  默认端口: $(PORT)"
	@echo "  默认命名空间: $(NAMESPACE)"
	@echo "  默认副本数: $(REPLICAS)"
	@echo ""
	@echo "可用命令:"
	@echo "  make help          - 显示所有命令"
	@echo "  make init          - 初始化项目"
	@echo "  make controller    - 生成控制器"
	@echo "  make model         - 生成模型"
	@echo "  make docker        - 生成 Docker 配置"
	@echo "  make k8s           - 生成 Kubernetes 配置"
	@echo "  make api           - 快速生成 API 组件"
	@echo "  make crud          - 快速生成 CRUD 组件"

# =============================================================================
# 示例用法
# =============================================================================

.PHONY: example-api
example-api: ## 生成示例 API 项目
	@echo "生成示例 API 项目..."
	$(ARTISAN) make:controller user --namespace=api
	$(ARTISAN) make:model user --fields=name:string,email:string,age:int
	$(ARTISAN) make:migration create_users_table --table=users --fields=name:string,email:string,age:int
	$(ARTISAN) make:test user --type=unit
	$(ARTISAN) make:docker --name=user-api --port=3000
	$(ARTISAN) make:k8s --name=user-api --replicas=3 --port=3000
	@echo "✅ 示例 API 项目生成完成!"

.PHONY: example-crud
example-crud: ## 生成示例 CRUD 项目
	@echo "生成示例 CRUD 项目..."
	$(ARTISAN) make:controller product --namespace=app
	$(ARTISAN) make:model product --fields=name:string,price:decimal,description:text
	$(ARTISAN) make:migration create_products_table --table=products --fields=name:string,price:decimal,description:text
	$(ARTISAN) make:middleware auth
	$(ARTISAN) make:test product --type=integration
	$(ARTISAN) make:docker --name=product-crud --port=8080 --env=development
	$(ARTISAN) make:k8s --name=product-crud --replicas=2 --port=8080 --namespace=development
	@echo "✅ 示例 CRUD 项目生成完成!" 