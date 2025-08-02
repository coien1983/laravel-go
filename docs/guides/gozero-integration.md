# Laravel-Go 与 go-zero 集成指南

## 📖 概述

Laravel-Go 框架集成了类似 go-zero goctl 的功能，可以根据 proto 文件和 .api 文件自动生成完整的微服务代码，包括 RPC 服务、API 服务、网关配置等。

## 🚀 快速开始

### 1. 从 proto 文件生成 RPC 服务

```bash
# 生成完整的 go-zero RPC 服务
largo gozero:proto user.proto --output=./user-service

# 指定输出目录
largo gozero:proto user.proto -o ./services/user
```

**生成的文件结构：**

```
user-service/
├── main.go                    # 主服务文件
├── etc/
│   └── userservice.yaml       # 配置文件
└── internal/
    ├── config/
    │   └── config.go          # 配置结构
    ├── svc/
    │   └── servicecontext.go  # 服务上下文
    └── server/
        └── userserviceserver.go # 服务实现
```

### 2. 从 .api 文件生成 API 服务

```bash
# 生成完整的 go-zero API 服务
largo gozero:api user.api --output=./user-api

# 指定输出目录
largo gozero:api user.api -o ./apis/user
```

**生成的文件结构：**

```
user-api/
├── main.go                    # 主服务文件
├── etc/
│   └── api.yaml              # 配置文件
└── internal/
    ├── config/
    │   └── config.go         # 配置结构
    ├── svc/
    │   └── servicecontext.go # 服务上下文
    ├── types/
    │   └── types.go          # 类型定义
    └── handler/
        └── handlers.go       # 处理器
```

### 3. 生成完整的微服务

```bash
# 生成包含 RPC 和 API 的完整微服务
largo gozero:microservice user-service --proto=user.proto --api=user.api --output=./user-microservice

# 只生成 RPC 服务
largo gozero:microservice user-service --proto=user.proto --output=./user-rpc

# 只生成 API 服务
largo gozero:microservice user-service --api=user.api --output=./user-api
```

**生成的文件结构：**

```
user-microservice/
├── gateway.yaml              # 网关配置
├── rpc/                      # RPC 服务
│   ├── main.go
│   ├── etc/
│   └── internal/
└── api/                      # API 服务
    ├── main.go
    ├── etc/
    └── internal/
```

## 📋 命令详解

### gozero:proto - 从 proto 文件生成 RPC 服务

```bash
largo gozero:proto <proto_file> [--output=]
```

**参数：**

- `proto_file`: proto 文件路径（必需）
- `--output, -o`: 输出目录（可选，默认为当前目录）

**功能：**

- 解析 proto 文件中的服务定义
- 生成完整的 go-zero RPC 服务结构
- 自动生成服务实现框架
- 生成配置文件

### gozero:api - 从 .api 文件生成 API 服务

```bash
largo gozero:api <api_file> [--output=]
```

**参数：**

- `api_file`: .api 文件路径（必需）
- `--output, -o`: 输出目录（可选，默认为当前目录）

**功能：**

- 解析 .api 文件中的 API 定义
- 生成完整的 go-zero API 服务结构
- 自动生成处理器和路由
- 生成类型定义

### gozero:microservice - 生成完整微服务

```bash
largo gozero:microservice <name> [--proto=] [--api=] [--output=]
```

**参数：**

- `name`: 微服务名称（必需）
- `--proto, -p`: proto 文件路径（可选）
- `--api, -a`: .api 文件路径（可选）
- `--output, -o`: 输出目录（可选，默认为当前目录）

**功能：**

- 根据提供的文件生成 RPC 和/或 API 服务
- 生成网关配置文件
- 创建完整的微服务项目结构

### gozero:logic - 生成 logic 层

```bash
largo gozero:logic <method_name> [--service=] [--output=]
```

**参数：**

- `method_name`: 方法名称（必需）
- `--service, -s`: 服务名称（可选）
- `--output, -o`: 输出目录（可选，默认为 ./internal/logic）

**功能：**

- 生成指定方法的 logic 层代码
- 创建业务逻辑框架

### gozero:handler - 生成 handler 层

```bash
largo gozero:handler <endpoint_name> [--method=] [--path=] [--output=]
```

**参数：**

- `endpoint_name`: 端点名称（必需）
- `--method, -m`: HTTP 方法（可选，默认为 GET）
- `--path, -p`: API 路径（可选，默认为 /api/endpoint）
- `--output, -o`: 输出目录（可选，默认为 ./internal/handler）

**功能：**

- 生成指定端点的 handler 层代码
- 创建 HTTP 处理器框架

## 📝 文件格式说明

### Proto 文件格式

```protobuf
syntax = "proto3";

package user;

option go_package = "user/types";

// 服务定义
service UserService {
  // 方法定义
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}

// 消息定义
message User {
  int64 id = 1;
  string name = 2;
  string email = 3;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
}

message CreateUserResponse {
  User user = 1;
  string message = 2;
}
```

### .api 文件格式

```api
syntax = "v1"

info(
    title: "用户管理API"
    desc: "用户管理相关的API接口"
    author: "Laravel-Go"
    version: "1.0"
)

type (
    // 类型定义
    User {
        Id    int64  `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    }

    CreateUserReq {
        Name  string `json:"name"`
        Email string `json:"email"`
    }

    CreateUserResp {
        User    User   `json:"user"`
        Message string `json:"message"`
    }
)

service user-api {
    @doc "创建用户"
    @handler createUser
    post /api/users (CreateUserReq) returns (CreateUserResp)
}
```

## 🔧 高级用法

### 1. 自定义模板

你可以通过修改生成器中的模板来自定义生成的代码风格：

```go
// 在 framework/console/goctl_enhanced.go 中修改模板
logicTemplate := `package logic

import (
    "context"
    "{{ .ProjectName }}/internal/svc"
    "{{ .ProjectName }}/internal/types"
)

type {{ .MethodName }}Logic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

// 自定义你的逻辑
func (l *{{ .MethodName }}Logic) {{ .MethodName }}(req *types.{{ .MethodName }}Req) (resp *types.{{ .MethodName }}Resp, err error) {
    // 在这里实现你的业务逻辑
    return &types.{{ .MethodName }}Resp{}, nil
}
`
```

### 2. 集成到现有项目

```bash
# 在现有项目中生成 RPC 服务
cd your-project
largo gozero:proto proto/user.proto --output=./internal/rpc

# 在现有项目中生成 API 服务
largo gozero:api api/user.api --output=./internal/api
```

### 3. 批量生成

```bash
# 生成多个服务
for proto in proto/*.proto; do
    service_name=$(basename "$proto" .proto)
    largo gozero:proto "$proto" --output="./services/$service_name"
done
```

## 🚀 运行生成的代码

### 1. 运行 RPC 服务

```bash
cd user-service
go mod tidy
go run main.go -f etc/userservice.yaml
```

### 2. 运行 API 服务

```bash
cd user-api
go mod tidy
go run main.go -f etc/api.yaml
```

### 3. 使用网关

```bash
# 配置网关路由
# 编辑 gateway.yaml 文件
# 启动网关服务
```

## 📚 最佳实践

### 1. 项目结构

```
project/
├── proto/                    # proto 文件
│   ├── user.proto
│   └── order.proto
├── api/                      # .api 文件
│   ├── user.api
│   └── order.api
├── services/                 # 生成的服务
│   ├── user-service/
│   └── order-service/
└── gateway/                  # 网关配置
    └── gateway.yaml
```

### 2. 命名规范

- 服务名使用小写字母和下划线
- 方法名使用 PascalCase
- 文件路径使用小写字母

### 3. 配置管理

- 使用环境变量管理配置
- 分离开发、测试、生产环境配置
- 使用配置中心管理配置

## 🔍 故障排除

### 1. 常见错误

**错误：proto 文件解析失败**

```bash
# 检查 proto 文件语法
protoc --proto_path=. --go_out=. user.proto
```

**错误：.api 文件解析失败**

```bash
# 检查 .api 文件语法
# 确保语法版本正确
syntax = "v1"
```

### 2. 调试技巧

```bash
# 启用详细输出
largo gozero:proto user.proto --output=./debug --verbose

# 检查生成的文件
ls -la ./debug/
```

## 📖 更多资源

- [go-zero 官方文档](https://go-zero.dev/)
- [Protocol Buffers 指南](https://developers.google.com/protocol-buffers)
- [Laravel-Go 框架文档](../README.md)
