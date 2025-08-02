# Laravel-Go 博客系统示例

## 📝 项目概览

这是一个使用 Laravel-Go Framework 构建的完整博客系统示例，展示了框架的核心功能和最佳实践。

## 🚀 功能特性

- ✅ 用户认证和授权
- ✅ 文章管理（CRUD）
- ✅ 数据验证
- ✅ 队列任务处理
- ✅ RESTful API
- ✅ 错误处理

## 📁 项目结构

```
blog_example/
├── app/
│   └── Http/
│       └── Controllers/
│           ├── AuthController.go
│           ├── PostController.go
│           └── UserController.go
├── config/
│   └── app.go
├── main.go
└── README.md
```

## 🏗️ 核心组件

### 1. 认证控制器 (AuthController)

提供用户注册、登录和退出功能：

- `Register`: 用户注册
- `Login`: 用户登录
- `Logout`: 用户退出

### 2. 文章控制器 (PostController)

提供文章的完整CRUD操作：

- `Index`: 获取文章列表
- `Show`: 获取单篇文章
- `Store`: 创建新文章
- `Update`: 更新文章
- `Destroy`: 删除文章

### 3. 用户控制器 (UserController)

提供用户信息管理：

- `Profile`: 获取用户信息
- `UpdateProfile`: 更新用户信息

## 🚀 快速开始

### 1. 运行示例

```bash
cd examples/blog_example
go run main.go
```

### 2. 测试API

#### 用户注册
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试用户",
    "email": "test@example.com",
    "password": "123456"
  }'
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }'
```

#### 获取文章列表
```bash
curl http://localhost:8080/posts
```

#### 创建文章
```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "新文章标题",
    "content": "这是文章内容...",
    "status": "published"
  }'
```

#### 获取单篇文章
```bash
curl http://localhost:8080/posts/1
```

#### 更新文章
```bash
curl -X PUT http://localhost:8080/posts/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "更新后的标题",
    "content": "更新后的内容...",
    "status": "published"
  }'
```

#### 删除文章
```bash
curl -X DELETE http://localhost:8080/posts/1
```

## 🔧 配置说明

### 应用配置

- `app.name`: 应用名称
- `app.version`: 应用版本
- `app.port`: 服务端口
- `app.debug`: 调试模式
- `app.timezone`: 时区设置

### 数据库配置

- `database.default`: 默认数据库连接
- `database.connections`: 数据库连接配置

### 缓存配置

- `cache.default`: 默认缓存驱动
- `cache.stores`: 缓存存储配置

### 队列配置

- `queue.default`: 默认队列连接
- `queue.connections`: 队列连接配置

## 📊 API 响应格式

### 成功响应

```json
{
  "message": "操作成功",
  "data": {
    // 具体数据
  }
}
```

### 错误响应

```json
{
  "error": "错误信息"
}
```

## 🧪 测试

### 单元测试

```bash
go test ./tests -v
```

### 集成测试

```bash
go test ./tests/integration_test.go -v
```

## 🔍 监控和日志

### 队列监控

系统会自动启动队列工作进程，处理异步任务：

- 任务完成时会输出日志
- 任务失败时会输出错误信息

### 性能监控

- 请求响应时间
- 内存使用情况
- 队列任务处理统计

## 🚀 部署

### 开发环境

```bash
go run main.go
```

### 生产环境

```bash
go build -o blog-server main.go
./blog-server
```

### Docker 部署

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o blog-server main.go
EXPOSE 8080
CMD ["./blog-server"]
```

## 📚 学习要点

### 1. 控制器设计

- 遵循RESTful设计原则
- 统一的响应格式
- 完善的错误处理

### 2. 数据验证

- 使用结构体标签进行验证
- 自定义验证规则
- 友好的错误信息

### 3. 队列处理

- 异步任务处理
- 任务重试机制
- 失败任务处理

### 4. 配置管理

- 环境变量支持
- 多环境配置
- 配置验证

## 🔗 相关文档

- [Laravel-Go Framework 文档](../docs/)
- [API 参考](../docs/api/)
- [用户指南](../docs/guides/)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个示例项目。

## 📄 许可证

本项目采用 MIT 许可证。 