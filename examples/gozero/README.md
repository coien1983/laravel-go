# Laravel-Go go-zero 集成示例

## 快速开始

### 1. 从 proto 文件生成 RPC 服务

```bash
# 生成用户 RPC 服务
largo gozero:proto user.proto --output=./user-service
```

### 2. 从 .api 文件生成 API 服务

```bash
# 生成用户 API 服务
largo gozero:api user.api --output=./user-api
```

### 3. 生成完整微服务

```bash
# 生成包含 RPC 和 API 的完整微服务
largo gozero:microservice user-service --proto=user.proto --api=user.api --output=./user-microservice
```

## 生成的文件结构

### RPC 服务

```
user-service/
├── main.go
├── etc/
│   └── userservice.yaml
└── internal/
    ├── config/
    ├── svc/
    └── server/
```

### API 服务

```
user-api/
├── main.go
├── etc/
│   └── api.yaml
└── internal/
    ├── config/
    ├── svc/
    ├── types/
    └── handler/
```

### 完整微服务

```
user-microservice/
├── gateway.yaml
├── rpc/
└── api/
```

## 运行服务

```bash
# 运行 RPC 服务
cd user-service
go run main.go -f etc/userservice.yaml

# 运行 API 服务
cd user-api
go run main.go -f etc/api.yaml
```
