# Laravel-Go Framework 文档

欢迎使用 Laravel-Go Framework！这是一个基于 Laravel 设计思路的 Go 语言开发框架，提供完整的 Web 开发、API 和微服务功能。

## 📚 文档导航

### 🚀 快速开始

- [安装指南](guides/installation.md) - 如何安装和配置框架
- [快速开始](guides/quickstart.md) - 5 分钟创建你的第一个应用
- [项目结构](guides/project-structure.md) - 了解项目目录结构

### 📖 用户指南

- [基础概念](guides/concepts.md) - 框架核心概念
- [应用容器](guides/container.md) - 依赖注入和服务管理
- [配置管理](guides/configuration.md) - 环境配置和变量管理
- [路由系统](guides/routing.md) - URL 路由和参数处理
- [中间件](guides/middleware.md) - 请求处理和中间件开发
- [控制器](guides/controllers.md) - MVC 控制器开发
- [数据库](guides/database.md) - 数据库连接和操作
- [ORM](guides/orm.md) - 对象关系映射
- [模板引擎](guides/templates.md) - 视图模板系统
- [缓存系统](guides/cache.md) - 数据缓存管理
- [队列系统](guides/queue.md) - 异步任务处理
- [事件系统](guides/events.md) - 事件驱动编程
- [认证授权](guides/auth.md) - 用户认证和权限控制
- [验证系统](guides/validation.md) - 数据验证
- [API 开发](guides/api.md) - RESTful API 开发
- [微服务](guides/microservices.md) - 微服务架构支持
- [gRPC 扩展](guides/grpc_extension.md) - gRPC 微服务扩展
- [命令行工具](guides/console.md) - Artisan 命令行工具
- [测试指南](guides/testing.md) - 单元测试和集成测试
- [性能优化](guides/performance.md) - 性能优化指南
- [定时器](guides/scheduler.md) - 定时任务调度
- [部署指南](guides/deployment.md) - 生产环境部署

### 🔧 API 参考

- [核心 API](api/core.md) - 框架核心 API
- [HTTP API](api/http.md) - HTTP 相关 API
- [数据库 API](api/database.md) - 数据库操作 API
- [ORM API](api/orm.md) - ORM 模型 API
- [缓存 API](api/cache.md) - 缓存操作 API
- [队列 API](api/queue.md) - 队列操作 API
- [事件 API](api/events.md) - 事件系统 API
- [认证 API](api/auth.md) - 认证授权 API
- [验证 API](api/validation.md) - 数据验证 API
- [命令行 API](api/console.md) - 命令行工具 API

### 💡 最佳实践

- [代码规范](best-practices/coding-standards.md) - 代码编写规范
- [项目结构](best-practices/project-structure.md) - 项目组织最佳实践
- [性能优化](best-practices/performance.md) - 性能优化技巧
- [安全实践](best-practices/security.md) - 安全开发指南
- [错误处理](best-practices/error-handling.md) - 错误处理最佳实践
- [日志记录](best-practices/logging.md) - 日志记录策略
- [测试策略](best-practices/testing.md) - 测试最佳实践
- [部署策略](best-practices/deployment.md) - 部署最佳实践

### 🎯 示例项目

- [博客系统](examples/blog/README.md) - 完整的博客应用示例
- [API 服务](examples/api/README.md) - RESTful API 示例
- [微服务](examples/microservices/README.md) - 微服务架构示例
- [认证系统](examples/auth/README.md) - 用户认证示例
- [文件上传](examples/file-upload/README.md) - 文件处理示例
- [WebSocket](examples/websocket/README.md) - 实时通信示例

### 🛠️ 工具和插件

- [Makefile 使用](makefile-usage.md) - Makefile 命令大全
- [部署命令](deployment-commands.md) - Docker 和 K8s 部署命令
- [IDE 配置](guides/ide-setup.md) - IDE 开发环境配置

### 📋 最新更新

- [文档索引](INDEX.md) - 完整文档索引和快速导航
- [最新更新](LATEST_UPDATES.md) - 最新功能更新和文档完善
- [文档总结](DOCUMENTATION_SUMMARY.md) - 完整文档体系总结
- [新文档总结](NEW_DOCUMENTS_SUMMARY.md) - 新增文档详细说明
- [更新日志](CHANGELOG.md) - 文档更新历史记录

## 🎯 框架特性

### ✨ 核心特性

- **依赖注入容器** - 强大的服务管理和依赖注入
- **路由系统** - 高性能的 Radix Tree 路由
- **中间件系统** - 灵活的请求处理管道
- **ORM 系统** - 优雅的数据库操作
- **模板引擎** - 强大的视图渲染
- **缓存系统** - 多驱动缓存支持
- **队列系统** - 异步任务处理
- **事件系统** - 事件驱动架构
- **认证授权** - 完整的用户认证系统
- **验证系统** - 强大的数据验证
- **命令行工具** - Artisan 命令行工具
- **测试支持** - 完整的测试框架

### 🚀 高级特性

- **微服务支持** - 服务发现和通信
- **gRPC 扩展** - 完整的 gRPC 微服务支持
- **API 开发** - RESTful API 开发工具
- **性能监控** - 内置性能监控和 Prometheus 集成
- **定时器系统** - 完整的任务调度系统
- **安全特性** - CSRF、XSS、SQL 注入防护
- **部署支持** - Docker 和 Kubernetes 支持
- **监控集成** - Prometheus 和 Grafana 监控

## 📦 安装

```bash
# 克隆项目
git clone https://github.com/your-org/laravel-go.git
cd laravel-go

# 安装依赖
go mod tidy

# 运行示例
go run main.go
```

## 🚀 快速开始

```go
package main

import (
    "laravel-go/framework"
    "laravel-go/framework/http"
)

func main() {
    app := framework.NewApplication()

    // 注册路由
    app.Router().Get("/", func(c http.Context) http.Response {
        return c.Json(map[string]string{
            "message": "Hello Laravel-Go!",
        })
    })

    // 启动服务器
    app.Run(":8080")
}
```

## 📖 学习路径

### 初学者

1. [安装指南](guides/installation.md)
2. [快速开始](guides/quickstart.md)
3. [基础概念](guides/concepts.md)
4. [路由系统](guides/routing.md)
5. [控制器](guides/controllers.md)

### 进阶用户

1. [应用容器](guides/container.md)
2. [ORM](guides/orm.md)
3. [中间件](guides/middleware.md)
4. [认证授权](guides/auth.md)
5. [API 开发](guides/api.md)

### 高级用户

1. [微服务](guides/microservices.md)
2. [性能优化](best-practices/performance.md)
3. [安全实践](best-practices/security.md)
4. [部署指南](guides/deployment.md)

## 🤝 贡献

我们欢迎所有形式的贡献！请查看 [贡献指南](CONTRIBUTING.md) 了解如何参与项目开发。

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 支持

- 📖 [文档](https://laravel-go.dev/docs)
- 💬 [社区讨论](https://github.com/your-org/laravel-go/discussions)
- 🐛 [问题反馈](https://github.com/your-org/laravel-go/issues)
- 📧 [邮件支持](mailto:support@laravel-go.dev)

## 🔄 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本更新历史。

---

**Laravel-Go Framework** - 让 Go 开发更优雅！ 🚀
