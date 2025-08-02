# gRPC 和 API Gateway 使用指南

## 概述

本项目包含完整的微服务架构，包括：

- **gRPC 服务**: 高性能的 RPC 服务
- **API Gateway**: HTTP 到 gRPC 的网关
- **Proto 文件**: 服务定义和消息格式

## 目录结构

```
microservice-demo/
├── proto/                    # Protocol Buffers 定义
│   └── user.proto           # 用户服务定义
├── grpc/                    # gRPC 相关代码
│   ├── server/              # gRPC 服务器
│   │   └── server.go        # 用户服务实现
│   ├── client/              # gRPC 客户端
│   │   └── client.go        # 用户客户端
│   └── interceptors/        # gRPC 拦截器
│       ├── logging.go       # 日志拦截器
│       └── auth.go          # 认证拦截器
└── gateway/                 # API Gateway
    ├── main.go              # 网关主程序
    ├── middleware/          # 网关中间件
    │   └── auth.go          # 认证中间件
    ├── routes/              # 路由定义
    │   └── routes.go        # 路由处理器
    └── plugins/             # 网关插件
        └── rate_limit.go    # 限流插件
```

## 快速开始

### 1. 生成 gRPC 代码

首先需要安装 Protocol Buffers 编译器：

```bash
# macOS
brew install protobuf

# Ubuntu/Debian
sudo apt-get install protobuf-compiler

# 安装 Go 的 protobuf 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

生成 Go 代码：

```bash
cd microservice-demo
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/user.proto
```

### 2. 启动 gRPC 服务器

```bash
cd microservice-demo
go run grpc/server/server.go
```

gRPC 服务器将在 `localhost:9090` 启动。

### 3. 启动 API Gateway

```bash
cd microservice-demo
go run gateway/main.go
```

API Gateway 将在 `localhost:8080` 启动。

### 4. 测试 API

#### 通过 API Gateway 访问（HTTP）

```bash
# 获取用户列表
curl http://localhost:8080/api/v1/users

# 获取单个用户
curl http://localhost:8080/api/v1/users/1

# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "password": "password123"
  }'

# 更新用户
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "李四",
    "email": "lisi@example.com",
    "phone": "13800138001"
  }'

# 删除用户
curl -X DELETE http://localhost:8080/api/v1/users/1

# 健康检查
curl http://localhost:8080/health
```

#### 直接访问 gRPC 服务

使用 gRPC 客户端：

```bash
cd microservice-demo
go run grpc/client/client.go
```

## 服务定义

### UserService (proto/user.proto)

```protobuf
service UserService {
  // 获取用户信息
  rpc GetUser(GetUserRequest) returns (GetUserResponse);

  // 创建用户
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);

  // 更新用户
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);

  // 删除用户
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);

  // 获取用户列表
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}
```

## API 路由

### API Gateway 路由

| 方法   | 路径                 | 描述         |
| ------ | -------------------- | ------------ |
| GET    | `/api/v1/users`      | 获取用户列表 |
| GET    | `/api/v1/users/{id}` | 获取单个用户 |
| POST   | `/api/v1/users`      | 创建用户     |
| PUT    | `/api/v1/users/{id}` | 更新用户     |
| DELETE | `/api/v1/users/{id}` | 删除用户     |
| GET    | `/health`            | 健康检查     |

## 中间件和拦截器

### gRPC 拦截器

- **LoggingInterceptor**: 记录请求日志
- **AuthInterceptor**: 认证拦截器
- **RecoveryInterceptor**: 错误恢复

### API Gateway 中间件

- **loggingMiddleware**: 请求日志
- **corsMiddleware**: CORS 支持
- **AuthMiddleware**: 认证中间件
- **RateLimitMiddleware**: 限流中间件

## 配置

### 环境变量

```bash
# gRPC 服务器端口
GRPC_PORT=9090

# API Gateway 端口
GATEWAY_PORT=8080

# 数据库配置
DB_CONNECTION=postgresql
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=microservice
DB_USERNAME=postgres
DB_PASSWORD=password

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
```

## 开发指南

### 添加新的 gRPC 服务

1. 在 `proto/` 目录下创建新的 `.proto` 文件
2. 定义服务和消息
3. 生成 Go 代码
4. 在 `grpc/server/` 下实现服务
5. 在 `grpc/client/` 下创建客户端
6. 在 `gateway/` 下添加 HTTP 路由

### 添加新的拦截器

1. 在 `grpc/interceptors/` 下创建新的拦截器
2. 在 `grpc/server/server.go` 中注册拦截器

### 添加新的中间件

1. 在 `gateway/middleware/` 下创建新的中间件
2. 在 `gateway/main.go` 中注册中间件

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t microservice-demo .

# 运行容器
docker run -p 8080:8080 -p 9090:9090 microservice-demo
```

### Kubernetes 部署

```bash
# 应用 Kubernetes 配置
kubectl apply -f k8s/

# 查看服务状态
kubectl get pods
kubectl get services
```

## 监控和日志

### 健康检查

- gRPC 服务: `grpc.health.v1.Health/Check`
- API Gateway: `GET /health`

### 日志

- gRPC 服务日志在控制台输出
- API Gateway 日志包含请求时间和状态码

### 指标

- Prometheus 指标收集
- 自定义业务指标

## 故障排除

### 常见问题

1. **gRPC 连接失败**

   - 检查端口是否正确
   - 确认服务是否启动

2. **Proto 文件编译错误**

   - 检查 protoc 版本
   - 确认 Go 插件已安装

3. **API Gateway 无法连接 gRPC**
   - 检查 gRPC 服务地址
   - 确认网络连接

### 调试工具

- **grpcurl**: gRPC 调试工具
- **grpcui**: gRPC Web UI
- **Postman**: API 测试

## 扩展功能

### 服务发现

集成 Consul 或 etcd 进行服务发现：

```go
// 示例：Consul 服务发现
consulClient, err := consul.NewClient(consul.DefaultConfig())
if err != nil {
    log.Fatal(err)
}

// 注册服务
err = consulClient.Agent().ServiceRegister(&consul.AgentServiceRegistration{
    Name: "user-service",
    Port: 9090,
})
```

### 负载均衡

使用 gRPC 内置的负载均衡：

```go
// 示例：轮询负载均衡
conn, err := grpc.Dial("consul://user-service",
    grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
)
```

### 熔断器

集成 Hystrix 或 Sentinel 进行熔断：

```go
// 示例：熔断器配置
circuitBreaker := hystrix.NewCircuitBreaker(hystrix.CommandConfig{
    Timeout:               1000,
    MaxConcurrentRequests: 100,
    ErrorPercentThreshold: 25,
})
```

## 性能优化

### gRPC 优化

1. **连接池**: 复用 gRPC 连接
2. **流式处理**: 使用流式 RPC
3. **压缩**: 启用 gRPC 压缩

### API Gateway 优化

1. **缓存**: Redis 缓存
2. **限流**: 请求限流
3. **压缩**: 响应压缩

## 安全

### 认证

- JWT Token 认证
- gRPC 元数据认证
- API Gateway 认证中间件

### 授权

- 基于角色的访问控制 (RBAC)
- 方法级别的权限控制

### 加密

- TLS/SSL 加密
- 敏感数据加密存储

## 测试

### 单元测试

```bash
go test ./grpc/server/
go test ./gateway/
```

### 集成测试

```bash
go test ./tests/integration/
```

### 性能测试

```bash
# 使用 grpcurl 进行性能测试
grpcurl -plaintext -d '{"id": 1}' localhost:9090 user.UserService/GetUser
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
