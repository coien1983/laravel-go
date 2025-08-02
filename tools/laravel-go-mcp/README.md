# Laravel-Go MCP 服务

这是一个基于 Model Context Protocol (MCP) 的 Laravel-Go 框架管理服务，提供了完整的项目生命周期管理功能。

## 功能特性

- 🚀 **项目初始化**: 快速创建新的 Laravel-Go 项目
- 📝 **代码生成**: 自动生成控制器、模型、服务等模块
- 🔨 **项目构建**: 自动化构建和编译
- 🧪 **测试运行**: 执行单元测试和集成测试
- 🚀 **项目部署**: 支持多环境部署
- 📊 **性能监控**: 实时性能数据收集
- 🔍 **代码分析**: 代码质量检查和统计
- ⚡ **性能优化**: 自动性能优化建议

## 快速开始

### 1. 启动 MCP 服务器

```bash
cd tools/laravel-go-mcp
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

### 2. 使用客户端示例

```bash
# 运行客户端演示
go run client_example.go
```

## API 接口

### 1. 初始化项目

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "name": "my-api",
    "description": "我的API项目",
    "version": "1.0.0",
    "modules": ["user", "product", "order"],
    "database": "mysql",
    "cache": "redis",
    "queue": "redis"
  }
}
```

### 2. 生成模块

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "generate",
  "params": {
    "type": "api",
    "name": "category"
  }
}
```

### 3. 构建项目

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "build",
  "params": {}
}
```

### 4. 运行测试

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "test",
  "params": {}
}
```

### 5. 部署项目

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "deploy",
  "params": {
    "environment": "production"
  }
}
```

### 6. 性能监控

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "monitor",
  "params": {}
}
```

### 7. 代码分析

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "analyze",
  "params": {}
}
```

### 8. 性能优化

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 8,
  "method": "optimize",
  "params": {}
}
```

### 9. 获取项目信息

```json
POST / HTTP/1.1
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 9,
  "method": "info",
  "params": {}
}
```

## 响应格式

所有接口都返回标准的 JSON-RPC 2.0 格式响应：

### 成功响应

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "success": true,
    "message": "操作成功",
    "data": {}
  }
}
```

### 错误响应

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32601,
    "message": "方法不存在",
    "data": null
  }
}
```

## 错误代码

| 代码   | 描述       |
| ------ | ---------- |
| -32700 | 解析错误   |
| -32600 | 无效请求   |
| -32601 | 方法不存在 |
| -32602 | 无效参数   |
| -32603 | 内部错误   |
| -32000 | 服务器错误 |

## 项目结构

```
laravel-go-mcp/
├── main.go              # MCP 服务器主文件
├── client_example.go    # 客户端示例
└── README.md           # 说明文档
```

## 环境变量

| 变量名   | 默认值 | 描述           |
| -------- | ------ | -------------- |
| MCP_PORT | 8080   | MCP 服务器端口 |

## 使用示例

### Go 客户端

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    client := NewMCPClientExample("http://localhost:8080")

    // 初始化项目
    resp, err := client.Initialize(&ClientInitializeRequest{
        Name:        "my-api",
        Description: "我的API项目",
        Version:     "1.0.0",
        Modules:     []string{"user", "product"},
        Database:    "mysql",
        Cache:       "redis",
        Queue:       "redis",
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("初始化结果: %v\n", resp["result"])
}
```

### cURL 示例

```bash
# 初始化项目
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "name": "my-api",
      "description": "我的API项目",
      "version": "1.0.0",
      "modules": ["user", "product"],
      "database": "mysql",
      "cache": "redis",
      "queue": "redis"
    }
  }'

# 生成模块
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "generate",
    "params": {
      "type": "api",
      "name": "category"
    }
  }'
```

## 开发指南

### 添加新的 MCP 方法

1. 在 `main.go` 中添加新的处理方法
2. 在 `client_example.go` 中添加对应的客户端方法
3. 更新文档说明

### 扩展功能

- 支持更多模块类型
- 添加数据库迁移功能
- 集成 CI/CD 流程
- 添加更多性能监控指标

## 许可证

MIT License
