# Laravel-Go API 示例

## 📝 项目概览

这是一个使用 Laravel-Go Framework 构建的完整 RESTful API 示例，展示了框架的 API 开发功能和最佳实践。

## 🚀 功能特性

- ✅ RESTful API 设计
- ✅ API 版本控制
- ✅ 中间件支持
- ✅ CORS 支持
- ✅ 统一响应格式
- ✅ 错误处理
- ✅ API 文档

## 📁 项目结构

```
api_example/
├── controllers/
│   ├── user_controller.go
│   ├── product_controller.go
│   └── order_controller.go
├── main.go
└── README.md
```

## 🏗️ 核心组件

### 1. 用户管理 API

提供用户的完整CRUD操作：

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取单个用户
- `POST /api/v1/users` - 创建新用户
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户

### 2. 产品管理 API

提供产品的完整CRUD操作：

- `GET /api/v1/products` - 获取产品列表
- `GET /api/v1/products/:id` - 获取单个产品
- `POST /api/v1/products` - 创建新产品
- `PUT /api/v1/products/:id` - 更新产品信息
- `DELETE /api/v1/products/:id` - 删除产品

### 3. 订单管理 API

提供订单的完整CRUD操作：

- `GET /api/v1/orders` - 获取订单列表
- `GET /api/v1/orders/:id` - 获取单个订单
- `POST /api/v1/orders` - 创建新订单
- `PUT /api/v1/orders/:id` - 更新订单信息
- `DELETE /api/v1/orders/:id` - 删除订单

## 🚀 快速开始

### 1. 运行示例

```bash
cd examples/api_example
go run main.go
```

### 2. 访问API文档

打开浏览器访问：http://localhost:8081/docs

### 3. 测试API

#### 健康检查
```bash
curl http://localhost:8081/api/v1/health
```

#### 用户管理
```bash
# 获取用户列表
curl http://localhost:8081/api/v1/users

# 获取单个用户
curl http://localhost:8081/api/v1/users/1

# 创建用户
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "新用户",
    "email": "newuser@example.com",
    "age": 25
  }'

# 更新用户
curl -X PUT http://localhost:8081/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新后的用户",
    "email": "updated@example.com",
    "age": 30
  }'

# 删除用户
curl -X DELETE http://localhost:8081/api/v1/users/1
```

#### 产品管理
```bash
# 获取产品列表
curl http://localhost:8081/api/v1/products

# 获取单个产品
curl http://localhost:8081/api/v1/products/1

# 创建产品
curl -X POST http://localhost:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "新产品",
    "description": "这是新产品的描述",
    "price": 999.99,
    "stock": 100,
    "category": "电子产品"
  }'

# 更新产品
curl -X PUT http://localhost:8081/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新后的产品",
    "description": "更新后的描述",
    "price": 899.99,
    "stock": 80,
    "category": "电子产品"
  }'

# 删除产品
curl -X DELETE http://localhost:8081/api/v1/products/1
```

#### 订单管理
```bash
# 获取订单列表
curl http://localhost:8081/api/v1/orders

# 获取单个订单
curl http://localhost:8081/api/v1/orders/1

# 创建订单
curl -X POST http://localhost:8081/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "product_id": 1,
    "quantity": 2
  }'

# 更新订单状态
curl -X PUT http://localhost:8081/api/v1/orders/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'

# 删除订单
curl -X DELETE http://localhost:8081/api/v1/orders/1
```

## 📊 API 响应格式

### 成功响应

```json
{
  "success": true,
  "data": {
    // 具体数据
  },
  "message": "操作成功",
  "total": 10
}
```

### 错误响应

```json
{
  "success": false,
  "error": "错误信息"
}
```

## 🔧 中间件

### 1. 日志中间件

记录所有API请求的日志信息：
- 请求方法
- 请求路径
- 客户端IP

### 2. CORS中间件

支持跨域请求：
- 允许所有来源
- 支持常用HTTP方法
- 支持自定义请求头

## 🚀 部署

### 开发环境

```bash
go run main.go
```

### 生产环境

```bash
go build -o api-server main.go
./api-server
```

### Docker 部署

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o api-server main.go
EXPOSE 8081
CMD ["./api-server"]
```

## 📚 学习要点

### 1. RESTful API 设计

- 使用标准HTTP方法
- 统一的URL结构
- 合适的HTTP状态码
- 一致的响应格式

### 2. API 版本控制

- 使用URL路径版本控制
- 向后兼容性考虑
- 版本迁移策略

### 3. 中间件使用

- 日志记录
- 跨域处理
- 认证授权
- 请求限流

### 4. 错误处理

- 统一的错误响应格式
- 合适的HTTP状态码
- 详细的错误信息

## 🔗 相关文档

- [Laravel-Go Framework 文档](../docs/)
- [API 参考](../docs/api/)
- [用户指南](../docs/guides/)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个示例项目。

## 📄 许可证

本项目采用 MIT 许可证。 