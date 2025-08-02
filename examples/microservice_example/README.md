# Laravel-Go 微服务示例

## 📝 项目概览

这是一个使用 Laravel-Go Framework 构建的完整微服务架构示例，展示了框架的微服务开发功能和最佳实践。

## 🚀 功能特性

- ✅ 微服务架构
- ✅ API 网关
- ✅ 服务发现
- ✅ 负载均衡
- ✅ 健康检查
- ✅ CORS 支持
- ✅ 统一响应格式

## 📁 项目结构

```
microservice_example/
├── user_service/
│   └── main.go
├── product_service/
│   └── main.go
├── order_service/
│   └── main.go
├── gateway/
│   └── main.go
└── README.md
```

## 🏗️ 核心组件

### 1. 用户微服务 (User Service)

- **端口**: 8082
- **功能**: 用户管理
- **API**: 
  - `GET /health` - 健康检查
  - `GET /users` - 获取用户列表
  - `GET /users/:id` - 获取单个用户

### 2. 产品微服务 (Product Service)

- **端口**: 8083
- **功能**: 产品管理
- **API**:
  - `GET /health` - 健康检查
  - `GET /products` - 获取产品列表
  - `GET /products/:id` - 获取单个产品

### 3. 订单微服务 (Order Service)

- **端口**: 8084
- **功能**: 订单管理
- **API**:
  - `GET /health` - 健康检查
  - `GET /orders` - 获取订单列表
  - `GET /orders/:id` - 获取单个订单

### 4. API 网关 (Gateway)

- **端口**: 8080
- **功能**: 统一入口、路由转发
- **API**:
  - `GET /` - 网关信息
  - `GET /health` - 网关健康检查
  - `GET /services` - 服务列表
  - `GET /users/*` - 转发到用户服务
  - `GET /products/*` - 转发到产品服务
  - `GET /orders/*` - 转发到订单服务

## 🚀 快速开始

### 1. 启动所有服务

#### 方法一：分别启动

```bash
# 终端1：启动用户服务
cd examples/microservice_example/user_service
go run main.go

# 终端2：启动产品服务
cd examples/microservice_example/product_service
go run main.go

# 终端3：启动订单服务
cd examples/microservice_example/order_service
go run main.go

# 终端4：启动API网关
cd examples/microservice_example/gateway
go run main.go
```

#### 方法二：使用脚本启动

```bash
# 创建启动脚本
cat > start_services.sh << 'EOF'
#!/bin/bash
echo "启动 Laravel-Go 微服务..."

# 启动用户服务
cd user_service && go run main.go &
USER_PID=$!

# 启动产品服务
cd ../product_service && go run main.go &
PRODUCT_PID=$!

# 启动订单服务
cd ../order_service && go run main.go &
ORDER_PID=$!

# 启动API网关
cd ../gateway && go run main.go &
GATEWAY_PID=$!

echo "所有服务已启动"
echo "用户服务 PID: $USER_PID"
echo "产品服务 PID: $PRODUCT_PID"
echo "订单服务 PID: $ORDER_PID"
echo "API网关 PID: $GATEWAY_PID"

# 等待中断信号
trap 'echo "正在关闭所有服务..."; kill $USER_PID $PRODUCT_PID $ORDER_PID $GATEWAY_PID; exit' INT
wait
EOF

chmod +x start_services.sh
./start_services.sh
```

### 2. 测试微服务

#### 通过API网关访问

```bash
# 网关信息
curl http://localhost:8080/

# 用户服务
curl http://localhost:8080/users
curl http://localhost:8080/users/1

# 产品服务
curl http://localhost:8080/products
curl http://localhost:8080/products/1

# 订单服务
curl http://localhost:8080/orders
curl http://localhost:8080/orders/1

# 服务列表
curl http://localhost:8080/services
```

#### 直接访问微服务

```bash
# 用户服务
curl http://localhost:8082/users
curl http://localhost:8082/users/1

# 产品服务
curl http://localhost:8083/products
curl http://localhost:8083/products/1

# 订单服务
curl http://localhost:8084/orders
curl http://localhost:8084/orders/1
```

#### 健康检查

```bash
# 网关健康检查
curl http://localhost:8080/health

# 各服务健康检查
curl http://localhost:8082/health
curl http://localhost:8083/health
curl http://localhost:8084/health
```

## 📊 服务架构图

```
┌─────────────────┐
│   客户端应用     │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐
│   API 网关      │
│   (8080)        │
└─────────┬───────┘
          │
    ┌─────┼─────┐
    │     │     │
    ▼     ▼     ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│ 用户服务 │ │ 产品服务 │ │ 订单服务 │
│ (8082)  │ │ (8083)  │ │ (8084)  │
└─────────┘ └─────────┘ └─────────┘
```

## 🔧 配置说明

### 服务端口配置

- **API 网关**: 8080
- **用户服务**: 8082
- **产品服务**: 8083
- **订单服务**: 8084

### 服务发现

当前示例使用静态配置的服务发现，生产环境可以集成：

- Consul
- etcd
- ZooKeeper
- Nacos

## 🚀 部署

### 开发环境

```bash
# 分别启动各个服务
go run main.go
```

### 生产环境

```bash
# 编译各个服务
go build -o user-service main.go
go build -o product-service main.go
go build -o order-service main.go
go build -o api-gateway main.go

# 启动服务
./user-service &
./product-service &
./order-service &
./api-gateway
```

### Docker 部署

```dockerfile
# 用户服务 Dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o user-service main.go
EXPOSE 8082
CMD ["./user-service"]
```

```yaml
# docker-compose.yml
version: '3.8'
services:
  user-service:
    build: ./user_service
    ports:
      - "8082:8082"
    environment:
      - SERVICE_PORT=8082

  product-service:
    build: ./product_service
    ports:
      - "8083:8083"
    environment:
      - SERVICE_PORT=8083

  order-service:
    build: ./order_service
    ports:
      - "8084:8084"
    environment:
      - SERVICE_PORT=8084

  api-gateway:
    build: ./gateway
    ports:
      - "8080:8080"
    depends_on:
      - user-service
      - product-service
      - order-service
```

## 📚 学习要点

### 1. 微服务架构设计

- 服务拆分原则
- 服务间通信
- 数据一致性
- 故障隔离

### 2. API 网关

- 统一入口
- 路由转发
- 负载均衡
- 安全控制

### 3. 服务发现

- 服务注册
- 服务发现
- 健康检查
- 故障转移

### 4. 监控和日志

- 服务监控
- 链路追踪
- 日志聚合
- 告警机制

## 🔗 相关文档

- [Laravel-Go Framework 文档](../docs/)
- [微服务指南](../docs/guides/microservices.md)
- [API 参考](../docs/api/)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个示例项目。

## 📄 许可证

本项目采用 MIT 许可证。 