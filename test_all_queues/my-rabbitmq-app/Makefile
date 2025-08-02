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