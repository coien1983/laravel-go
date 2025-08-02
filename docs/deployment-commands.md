# 部署配置生成命令

Laravel-Go Framework 提供了两个便捷的命令来生成 Docker 和 Kubernetes 部署配置文件。

## Docker 部署配置

### 命令语法

```bash
go run cmd/artisan/main.go make:docker [选项]
```

### 可用选项

- `--name, -n`: 应用名称 (默认: laravel-go-app)
- `--port, -p`: 应用端口 (默认: 8080)
- `--env, -e`: 环境 (development/production, 默认: development)

### 使用示例

```bash
# 使用默认配置生成 Docker 文件
go run cmd/artisan/main.go make:docker

# 自定义应用名称和端口
go run cmd/artisan/main.go make:docker --name=my-api --port=3000

# 生成生产环境配置
go run cmd/artisan/main.go make:docker --name=my-api --port=3000 --env=production
```

### 生成的文件

- `Dockerfile`: 多阶段构建的 Docker 镜像配置
- `docker-compose.yml`: 包含应用、Redis、PostgreSQL 的完整服务配置
- `.dockerignore`: 排除不需要的文件和目录

## Kubernetes 部署配置

### 命令语法

```bash
go run cmd/artisan/main.go make:k8s [选项]
```

### 可用选项

- `--name, -n`: 应用名称 (默认: laravel-go-app)
- `--replicas, -r`: 副本数量 (默认: 3)
- `--port, -p`: 应用端口 (默认: 8080)
- `--namespace, -ns`: Kubernetes 命名空间 (默认: default)

### 使用示例

```bash
# 使用默认配置生成 K8s 文件
go run cmd/artisan/main.go make:k8s

# 自定义应用配置
go run cmd/artisan/main.go make:k8s --name=my-api --replicas=5 --port=3000

# 指定命名空间
go run cmd/artisan/main.go make:k8s --name=my-api --namespace=production
```

### 生成的文件

- `k8s/deployment.yaml`: 应用部署配置，包含健康检查和资源限制
- `k8s/service.yaml`: 服务配置，用于内部通信
- `k8s/ingress.yaml`: 入口配置，用于外部访问
- `k8s/configmap.yaml`: 配置映射，包含环境变量

## 部署步骤

### Docker 部署

1. 生成配置文件：

```bash
go run cmd/artisan/main.go make:docker --name=my-app --port=3000
```

2. 构建镜像：

```bash
docker build -t my-app .
```

3. 启动服务：

```bash
docker-compose up -d
```

### Kubernetes 部署

1. 生成配置文件：

```bash
go run cmd/artisan/main.go make:k8s --name=my-app --replicas=3
```

2. 构建镜像：

```bash
docker build -t my-app .
```

3. 推送镜像到仓库：

```bash
docker tag my-app your-registry/my-app:latest
docker push your-registry/my-app:latest
```

4. 部署到 Kubernetes：

```bash
kubectl apply -f k8s/
```

## 配置说明

### Docker 配置特性

- **多阶段构建**: 优化镜像大小
- **Alpine Linux**: 轻量级基础镜像
- **健康检查**: 内置健康检查端点
- **环境变量**: 支持开发和生产环境
- **网络配置**: 自动配置服务网络

### Kubernetes 配置特性

- **资源限制**: CPU 和内存限制
- **健康检查**: 存活性和就绪性探针
- **自动扩缩容**: 支持 HPA 配置
- **服务发现**: 内部服务通信
- **负载均衡**: Ingress 配置
- **配置管理**: ConfigMap 支持

## 自定义配置

生成的配置文件可以根据需要进行自定义：

1. **修改资源限制**: 在 `deployment.yaml` 中调整 CPU 和内存
2. **添加环境变量**: 在 `configmap.yaml` 中添加更多配置
3. **调整健康检查**: 修改探针路径和间隔
4. **配置持久化**: 添加 PVC 配置
5. **设置安全策略**: 添加 SecurityContext

## 最佳实践

1. **镜像标签**: 使用语义化版本标签
2. **资源限制**: 根据实际需求设置合理的资源限制
3. **健康检查**: 确保健康检查端点正确实现
4. **环境分离**: 为不同环境使用不同的配置
5. **监控配置**: 添加监控和日志收集配置
