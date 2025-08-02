# Laravel-Go Framework 示例项目

## 📝 项目概览

这是 Laravel-Go Framework 的完整示例项目集合，展示了框架的各种使用场景和最佳实践。

## 🚀 示例项目列表

### 1. 博客系统示例 (`blog_example/`)

一个完整的博客系统，展示了框架的核心功能：

- **功能特性**: 用户认证、文章管理、数据验证、队列处理
- **技术栈**: HTTP 服务器、控制器、验证、队列
- **配置文件**: 应用、数据库、缓存、队列、日志、会话
- **端口**: 8080
- **快速开始**: `cd blog_example && go run main.go`

**相关文档**: [博客示例详细文档](blog_example/README.md)

### 2. API 开发示例 (`api_example/`)

一个完整的 RESTful API 示例，展示了 API 开发的最佳实践：

- **功能特性**: RESTful API、版本控制、中间件、CORS 支持
- **技术栈**: API 网关、控制器、中间件、文档生成
- **配置文件**: 应用、API、环境变量
- **端口**: 8081
- **快速开始**: `cd api_example && go run main.go`

**相关文档**: [API 示例详细文档](api_example/README.md)

### 3. 微服务示例 (`microservice_example/`)

一个完整的微服务架构示例，展示了分布式系统的设计：

- **功能特性**: 微服务架构、API 网关、服务发现、负载均衡
- **技术栈**: 微服务、网关、服务注册、健康检查
- **配置文件**: 应用、服务、环境变量
- **端口**: 8080 (网关), 8082-8084 (微服务)
- **快速开始**: 分别启动各个服务

**相关文档**: [微服务示例详细文档](microservice_example/README.md)

### 4. 部署示例 (`deployment_example/`)

一个完整的部署解决方案，支持多种部署方式：

- **功能特性**: Docker 部署、Kubernetes 部署、监控、自动化脚本
- **技术栈**: Docker、Kubernetes、Nginx、Prometheus、Grafana
- **配置文件**: 应用、部署、环境变量
- **快速开始**: 使用部署脚本或手动部署

**相关文档**: [部署示例详细文档](deployment_example/README.md)

### 5. 配置管理工具 (`config_manager/`)

一个用于管理配置的命令行工具：

- **功能特性**: 配置读取、设置、验证、格式转换
- **技术栈**: 命令行工具、配置管理、多格式支持
- **配置文件**: 支持 JSON、YAML、ENV 格式
- **快速开始**: `cd config_manager && go build -o config-manager main.go`

**相关文档**: [配置管理工具详细文档](config_manager/README.md)

## 🏗️ 示例项目架构

```
examples/
├── blog_example/           # 博客系统示例
│   ├── app/
│   │   └── Http/
│   │       └── Controllers/
│   ├── config/             # 配置文件
│   │   ├── app.go
│   │   ├── database.go
│   │   ├── cache.go
│   │   ├── queue.go
│   │   ├── logging.go
│   │   └── session.go
│   ├── env.example         # 环境变量示例
│   ├── main.go
│   └── README.md
├── api_example/            # API 开发示例
│   ├── controllers/
│   ├── config/             # 配置文件
│   │   ├── app.go
│   │   └── api.go
│   ├── env.example         # 环境变量示例
│   ├── main.go
│   └── README.md
├── microservice_example/   # 微服务示例
│   ├── user_service/
│   ├── product_service/
│   ├── order_service/
│   ├── gateway/
│   ├── config/             # 配置文件
│   │   ├── app.go
│   │   └── services.go
│   ├── env.example         # 环境变量示例
│   └── README.md
├── deployment_example/     # 部署示例
│   ├── docker/
│   ├── kubernetes/
│   ├── scripts/
│   ├── config/             # 配置文件
│   │   ├── app.go
│   │   └── deployment.go
│   ├── env.example         # 环境变量示例
│   └── README.md
├── config_manager/         # 配置管理工具
│   ├── main.go
│   └── README.md
└── README.md              # 本文档
```

## 🚀 快速开始

### 1. 选择示例项目

根据你的需求选择合适的示例项目：

- **学习框架基础**: 选择 `blog_example`
- **开发 API**: 选择 `api_example`
- **构建微服务**: 选择 `microservice_example`
- **部署应用**: 选择 `deployment_example`
- **管理配置**: 选择 `config_manager`

### 2. 运行示例

```bash
# 博客示例
cd examples/blog_example
cp env.example .env  # 复制环境变量文件
go run main.go

# API 示例
cd examples/api_example
cp env.example .env  # 复制环境变量文件
go run main.go

# 微服务示例 (需要启动多个服务)
cd examples/microservice_example
cp env.example .env  # 复制环境变量文件
cd user_service && go run main.go &
cd product_service && go run main.go &
cd order_service && go run main.go &
cd gateway && go run main.go

# 部署示例
cd examples/deployment_example
cp env.example .env  # 复制环境变量文件
./scripts/deploy.sh -e dev -p docker -b -d

# 配置管理工具
cd examples/config_manager
go build -o config-manager main.go
./config-manager -h
```

### 3. 访问示例

- **博客系统**: http://localhost:8080
- **API 服务**: http://localhost:8081
- **微服务网关**: http://localhost:8080
- **Docker 部署**: http://localhost

## 📚 学习路径

### 初学者路径

1. **第一步**: 运行 `blog_example`，了解框架基础
2. **第二步**: 运行 `api_example`，学习 API 开发
3. **第三步**: 运行 `microservice_example`，理解微服务架构
4. **第四步**: 使用 `deployment_example`，学习部署
5. **第五步**: 使用 `config_manager`，学习配置管理

### 进阶学习路径

1. **深入框架**: 阅读框架源码，理解设计模式
2. **定制开发**: 基于示例项目进行定制开发
3. **生产部署**: 使用部署示例进行生产环境部署
4. **性能优化**: 学习监控和性能优化

## 🔧 开发环境要求

### 基础要求

- **Go**: 1.21 或更高版本
- **Git**: 版本控制
- **编辑器**: VS Code、GoLand 等

### 可选要求

- **Docker**: 容器化部署
- **Kubernetes**: 集群部署
- **PostgreSQL**: 数据库
- **Redis**: 缓存

## 📊 示例项目对比

| 示例项目             | 复杂度 | 适用场景     | 学习重点             | 配置文件                             |
| -------------------- | ------ | ------------ | -------------------- | ------------------------------------ |
| blog_example         | 简单   | 学习框架基础 | 控制器、验证、队列   | 应用、数据库、缓存、队列、日志、会话 |
| api_example          | 中等   | API 开发     | RESTful API、中间件  | 应用、API、环境变量                  |
| microservice_example | 复杂   | 微服务架构   | 服务拆分、网关、发现 | 应用、服务、环境变量                 |
| deployment_example   | 复杂   | 生产部署     | 容器化、编排、监控   | 应用、部署、环境变量                 |
| config_manager       | 简单   | 配置管理     | 配置读取、验证、转换 | 支持多种格式                         |

## 🎯 最佳实践

### 1. 代码组织

- 遵循 MVC 架构
- 使用依赖注入
- 实现接口分离
- 编写单元测试

### 2. 配置管理

- 使用环境变量
- 分离开发和生产配置
- 使用配置验证
- 实现配置热重载

### 3. 错误处理

- 统一错误格式
- 实现错误中间件
- 记录错误日志
- 提供错误监控

### 4. 性能优化

- 使用缓存
- 实现数据库优化
- 配置负载均衡
- 监控性能指标

## 🔗 相关资源

### 框架文档

- [Laravel-Go Framework 文档](../docs/)
- [API 参考](../docs/api/)
- [用户指南](../docs/guides/)

### 外部资源

- [Go 官方文档](https://golang.org/doc/)
- [Docker 官方文档](https://docs.docker.com/)
- [Kubernetes 官方文档](https://kubernetes.io/docs/)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这些示例项目。

### 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

### 代码规范

- 遵循 Go 代码规范
- 添加必要的注释
- 编写单元测试
- 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证。

## 📞 支持

如果你在使用这些示例项目时遇到问题：

1. 查看相关文档
2. 搜索现有 Issue
3. 创建新的 Issue
4. 联系维护团队

---

**Happy Coding! 🚀**
