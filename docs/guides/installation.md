# 安装指南

本指南将帮助你安装和配置 Laravel-Go Framework。

## 📋 系统要求

### 必需软件

- **Go 1.21+** - 编程语言环境
- **Git** - 版本控制工具

### 可选软件

- **Docker** - 容器化部署
- **Kubernetes** - 容器编排
- **PostgreSQL/MySQL** - 数据库
- **Redis** - 缓存和队列

## 🚀 安装步骤

### 1. 安装 Go

#### macOS

```bash
# 使用 Homebrew
brew install go

# 或下载官方安装包
# https://golang.org/dl/
```

#### Linux

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# CentOS/RHEL
sudo yum install golang

# 或使用官方安装脚本
curl -O https://golang.org/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
```

#### Windows

```bash
# 下载官方安装包
# https://golang.org/dl/
# 运行安装程序并按照向导完成安装
```

### 2. 验证 Go 安装

```bash
go version
# 输出: go version go1.21.x darwin/amd64
```

### 3. 设置 Go 环境

```bash
# 设置 GOPATH (可选，Go 1.11+ 默认使用 modules)
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# 添加到 ~/.bashrc 或 ~/.zshrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc
```

### 4. 克隆项目

```bash
# 克隆 Laravel-Go Framework
git clone https://github.com/your-org/laravel-go.git
cd laravel-go
```

### 5. 安装依赖

```bash
# 下载依赖
go mod tidy

# 验证安装
go mod verify
```

### 6. 运行测试

```bash
# 运行所有测试
go test ./...

# 或使用 Makefile
make test-all
```

## 🔧 配置开发环境

### 1. IDE 配置

#### VS Code

```bash
# 安装 Go 扩展
code --install-extension golang.go

# 安装推荐的扩展
code --install-extension ms-vscode.go
code --install-extension bradlc.vscode-tailwindcss
```

#### GoLand

- 下载并安装 GoLand
- 导入项目
- 配置 Go 环境

### 2. 代码格式化工具

```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装 goimports
go install golang.org/x/tools/cmd/goimports@latest
```

### 3. 创建配置文件

```bash
# 复制环境配置文件
cp .env.example .env

# 编辑配置文件
nano .env
```

## 🐳 Docker 安装 (可选)

### 1. 安装 Docker

```bash
# macOS
brew install --cask docker

# Linux
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Windows
# 下载 Docker Desktop
```

### 2. 验证 Docker 安装

```bash
docker --version
docker-compose --version
```

### 3. 使用 Docker 运行

```bash
# 生成 Docker 配置
make docker

# 构建并运行
docker-compose up -d
```

## ☸️ Kubernetes 安装 (可选)

### 1. 安装 kubectl

```bash
# macOS
brew install kubectl

# Linux
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Windows
# 下载 kubectl.exe
```

### 2. 安装 Minikube (本地开发)

```bash
# macOS
brew install minikube

# Linux
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Windows
# 下载 minikube-windows-amd64.exe
```

### 3. 启动 Kubernetes 集群

```bash
# 启动 Minikube
minikube start

# 验证集群
kubectl cluster-info
```

## 🗄️ 数据库安装 (可选)

### PostgreSQL

```bash
# macOS
brew install postgresql
brew services start postgresql

# Linux
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Docker
docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:15
```

### MySQL

```bash
# macOS
brew install mysql
brew services start mysql

# Linux
sudo apt install mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# Docker
docker run --name mysql -e MYSQL_ROOT_PASSWORD=password -d -p 3306:3306 mysql:8
```

### Redis

```bash
# macOS
brew install redis
brew services start redis

# Linux
sudo apt install redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server

# Docker
docker run --name redis -d -p 6379:6379 redis:7-alpine
```

## 🚀 快速验证

### 1. 创建测试项目

```bash
# 使用框架初始化项目
make init-custom
# 输入项目名称: my-test-app
```

### 2. 生成示例组件

```bash
# 生成 API 组件
make api
# 输入资源名称: user
```

### 3. 运行应用

```bash
# 启动开发服务器
make run

# 或直接运行
go run main.go
```

### 4. 测试 API

```bash
# 测试健康检查
curl http://localhost:8080/health

# 测试用户 API
curl http://localhost:8080/api/users
```

## 🔍 故障排除

### 常见问题

#### 1. Go 版本过低

```bash
# 检查 Go 版本
go version

# 如果版本低于 1.21，请升级
# macOS
brew upgrade go

# Linux
# 下载新版本并安装
```

#### 2. 依赖下载失败

```bash
# 设置代理 (中国用户)
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod tidy
```

#### 3. 端口被占用

```bash
# 检查端口占用
lsof -i :8080

# 杀死占用进程
kill -9 <PID>

# 或使用不同端口
go run main.go -port 8081
```

#### 4. 权限问题

```bash
# 修复权限
sudo chown -R $(whoami) /usr/local/go
sudo chown -R $(whoami) $GOPATH

# 或使用用户目录
export GOPATH=$HOME/go
```

### 获取帮助

```bash
# 查看框架帮助
make help

# 查看项目信息
make info

# 查看命令行工具帮助
make list
```

## 📚 下一步

安装完成后，建议按以下顺序学习：

1. [快速开始](quickstart.md) - 创建第一个应用
2. [基础概念](concepts.md) - 了解框架核心概念
3. [项目结构](project-structure.md) - 熟悉项目组织
4. [路由系统](routing.md) - 学习 URL 路由
5. [控制器](controllers.md) - 开发业务逻辑

## 🆘 需要帮助？

如果遇到问题，可以通过以下方式获取帮助：

- 📖 [文档](https://laravel-go.dev/docs)
- 💬 [社区讨论](https://github.com/your-org/laravel-go/discussions)
- 🐛 [问题反馈](https://github.com/your-org/laravel-go/issues)
- 📧 [邮件支持](mailto:support@laravel-go.dev)

---

恭喜！你已经成功安装了 Laravel-Go Framework。现在可以开始你的开发之旅了！ 🚀
