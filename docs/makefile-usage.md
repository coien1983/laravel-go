# Laravel-Go Framework Makefile 使用指南

这个 Makefile 提供了 Laravel-Go Framework 所有命令行工具的便捷访问方式，让开发更加高效。

## 🚀 快速开始

### 查看所有可用命令

```bash
make help
```

### 显示项目信息

```bash
make info
```

## 📁 项目初始化

### 基础初始化

```bash
make init
```

### 自定义名称初始化

```bash
make init-custom
```

## 🔧 代码生成

### 控制器生成

```bash
# 基础控制器
make controller

# 自定义命名空间控制器
make controller-custom
```

### 模型生成

```bash
# 基础模型
make model

# 带字段的模型
make model-fields
```

### 中间件生成

```bash
make middleware
```

### 迁移文件生成

```bash
# 基础迁移
make migration

# 指定表名
make migration-table

# 指定表名和字段
make migration-fields
```

### 测试文件生成

```bash
# 基础测试
make test

# 指定类型测试
make test-type
```

## 🐳 部署配置生成

### Docker 配置

```bash
# 使用默认配置
make docker

# 自定义配置
make docker-custom
```

### Kubernetes 配置

```bash
# 使用默认配置
make k8s

# 自定义配置
make k8s-custom
```

## ⚡ 快速生成常用组件

### API 组件

```bash
make api
```

生成完整的 API 组件，包括：

- 控制器 (api 命名空间)
- 模型
- 迁移文件

### CRUD 组件

```bash
make crud
```

生成完整的 CRUD 组件，包括：

- 控制器 (app 命名空间)
- 模型
- 迁移文件
- 单元测试
- 集成测试

## 🛠️ 开发工具

### 环境设置

```bash
make dev-setup
```

### 构建应用

```bash
make build
```

### 运行应用

```bash
make run
```

### 测试

```bash
# 运行所有测试
make test-all

# 生成覆盖率报告
make test-coverage
```

### 代码质量

```bash
# 代码检查
make lint

# 格式化代码
make fmt

# 静态分析
make vet
```

## 🐳 Docker 操作

### 构建镜像

```bash
make docker-build
```

### 运行容器

```bash
make docker-run
```

### Docker Compose

```bash
# 启动服务
make docker-compose-up

# 停止服务
make docker-compose-down

# 查看日志
make docker-compose-logs
```

## ☸️ Kubernetes 操作

### 部署

```bash
# 部署到集群
make k8s-apply

# 从集群删除
make k8s-delete
```

### 监控

```bash
# 查看状态
make k8s-status

# 查看日志
make k8s-logs
```

## 🧹 清理操作

### 清理构建文件

```bash
make clean
```

### 清理 Docker 文件

```bash
make clean-docker
```

### 清理 Kubernetes 文件

```bash
make clean-k8s
```

### 清理所有文件

```bash
make clean-all
```

## 📋 项目维护

### 查看路由

```bash
make routes
```

### 清除缓存

```bash
make cache-clear
```

### 列出命令

```bash
make list
```

### 查看版本

```bash
make version
```

## 🎯 示例项目

### 生成示例 API 项目

```bash
make example-api
```

生成包含用户管理的完整 API 项目。

### 生成示例 CRUD 项目

```bash
make example-crud
```

生成包含产品管理的完整 CRUD 项目。

## 🔧 自定义配置

### 修改默认变量

在 Makefile 顶部可以修改默认配置：

```makefile
# 变量定义
ARTISAN := go run cmd/artisan/main.go
APP_NAME := laravel-go-app
PORT := 8080
NAMESPACE := default
REPLICAS := 3
```

### 添加自定义命令

可以在 Makefile 中添加自己的命令：

```makefile
.PHONY: my-command
my-command: ## 我的自定义命令
	@echo "执行自定义命令..."
	# 你的命令逻辑
```

## 💡 使用技巧

### 1. 交互式输入

大部分命令都支持交互式输入，会提示你输入必要的参数。

### 2. 默认值

所有命令都有合理的默认值，可以直接使用。

### 3. 组合使用

可以组合多个命令来快速搭建项目：

```bash
# 快速搭建 API 项目
make api
make docker
make k8s

# 快速搭建 CRUD 项目
make crud
make docker-custom
make k8s-custom
```

### 4. 开发流程

推荐的开发流程：

```bash
# 1. 设置开发环境
make dev-setup

# 2. 生成组件
make api  # 或 make crud

# 3. 运行测试
make test-all

# 4. 代码检查
make lint

# 5. 生成部署配置
make docker
make k8s

# 6. 部署
make docker-compose-up
# 或
make k8s-apply
```

## 🚨 注意事项

1. **依赖检查**: 确保已安装 Go、Docker、kubectl 等必要工具
2. **权限问题**: 某些命令可能需要管理员权限
3. **网络连接**: Docker 和 Kubernetes 命令需要网络连接
4. **配置文件**: 确保在项目根目录下运行命令

## 🔍 故障排除

### 常见问题

1. **命令未找到**

   ```bash
   # 确保在项目根目录
   pwd
   ls Makefile
   ```

2. **权限错误**

   ```bash
   # 使用 sudo (如果需要)
   sudo make docker-compose-up
   ```

3. **Docker 未运行**

   ```bash
   # 启动 Docker
   docker --version
   ```

4. **Kubernetes 集群未连接**
   ```bash
   # 检查集群连接
   kubectl cluster-info
   ```

### 获取帮助

```bash
# 查看所有命令
make help

# 查看项目信息
make info

# 查看命令行工具帮助
make list
```

这个 Makefile 大大简化了 Laravel-Go Framework 的开发流程，让开发者可以专注于业务逻辑，而不用担心复杂的命令操作！
