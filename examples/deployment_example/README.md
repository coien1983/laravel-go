# Laravel-Go 部署示例

## 📝 项目概览

这是一个完整的 Laravel-Go Framework 部署示例，支持 Docker 和 Kubernetes 两种部署方式，包含完整的生产环境配置。

## 🚀 功能特性

- ✅ Docker 容器化部署
- ✅ Kubernetes 集群部署
- ✅ Nginx 反向代理
- ✅ PostgreSQL 数据库
- ✅ Redis 缓存
- ✅ Prometheus 监控
- ✅ Grafana 可视化
- ✅ 自动化部署脚本
- ✅ 健康检查
- ✅ 负载均衡
- ✅ SSL/TLS 支持

## 📁 项目结构

```
deployment_example/
├── docker/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── nginx.conf
├── kubernetes/
│   ├── deployment.yaml
│   └── monitoring.yaml
├── scripts/
│   └── deploy.sh
└── README.md
```

## 🏗️ 核心组件

### 1. Docker 部署

#### Dockerfile
- 多阶段构建
- 非 root 用户运行
- 健康检查
- 最小化镜像大小

#### Docker Compose
- 应用服务
- PostgreSQL 数据库
- Redis 缓存
- Nginx 反向代理
- Prometheus 监控
- Grafana 可视化

#### Nginx 配置
- 反向代理
- 负载均衡
- Gzip 压缩
- 安全头设置
- SSL/TLS 支持

### 2. Kubernetes 部署

#### 应用部署
- Deployment 配置
- Service 配置
- Ingress 配置
- 健康检查
- 资源限制

#### 数据库部署
- PostgreSQL StatefulSet
- 持久化存储
- 服务发现

#### 监控部署
- Prometheus 配置
- Grafana 配置
- 指标收集
- 可视化面板

## 🚀 快速开始

### 1. Docker 部署

#### 使用部署脚本

```bash
# 构建并部署到开发环境
./examples/deployment_example/scripts/deploy.sh -e dev -p docker -b -d

# 构建并部署到生产环境
./examples/deployment_example/scripts/deploy.sh -e prod -p docker -b -d

# 查看日志
./examples/deployment_example/scripts/deploy.sh -p docker -l

# 停止服务
./examples/deployment_example/scripts/deploy.sh -p docker -s

# 重启服务
./examples/deployment_example/scripts/deploy.sh -p docker -r

# 清理资源
./examples/deployment_example/scripts/deploy.sh -p docker -c
```

#### 手动部署

```bash
# 构建镜像
docker build -f examples/deployment_example/docker/Dockerfile -t laravel-go-app:latest .

# 启动服务
cd examples/deployment_example/docker
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 2. Kubernetes 部署

#### 使用部署脚本

```bash
# 部署到开发环境
./examples/deployment_example/scripts/deploy.sh -e dev -p k8s -d

# 部署到生产环境
./examples/deployment_example/scripts/deploy.sh -e prod -p k8s -d

# 查看日志
./examples/deployment_example/scripts/deploy.sh -p k8s -l

# 停止服务
./examples/deployment_example/scripts/deploy.sh -p k8s -s

# 重启服务
./examples/deployment_example/scripts/deploy.sh -p k8s -r

# 清理资源
./examples/deployment_example/scripts/deploy.sh -p k8s -c
```

#### 手动部署

```bash
# 创建命名空间
kubectl create namespace laravel-go

# 部署应用
kubectl apply -f examples/deployment_example/kubernetes/deployment.yaml

# 部署监控
kubectl apply -f examples/deployment_example/kubernetes/monitoring.yaml

# 查看部署状态
kubectl get all -n laravel-go

# 查看日志
kubectl logs -f deployment/laravel-go-app -n laravel-go

# 删除部署
kubectl delete -f examples/deployment_example/kubernetes/monitoring.yaml
kubectl delete -f examples/deployment_example/kubernetes/deployment.yaml
```

## 📊 服务架构图

### Docker 架构

```
┌─────────────────┐
│   客户端应用     │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐
│   Nginx (80/443) │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐
│ Laravel-Go App  │
│   (8080)        │
└─────────┬───────┘
          │
    ┌─────┼─────┐
    │     │     │
    ▼     ▼     ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│PostgreSQL│ │  Redis  │ │Prometheus│
│ (5432)  │ │ (6379)  │ │ (9090)  │
└─────────┘ └─────────┘ └─────────┘
    │           │           │
    └───────────┼───────────┘
                ▼
        ┌─────────────┐
        │   Grafana   │
        │   (3000)    │
        └─────────────┘
```

### Kubernetes 架构

```
┌─────────────────┐
│   Ingress       │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐
│ Laravel-Go App  │
│   Service       │
└─────────┬───────┘
          │
    ┌─────┼─────┐
    │     │     │
    ▼     ▼     ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│PostgreSQL│ │  Redis  │ │Prometheus│
│ Service │ │ Service │ │ Service │
└─────────┘ └─────────┘ └─────────┘
    │           │           │
    └───────────┼───────────┘
                ▼
        ┌─────────────┐
        │   Grafana   │
        │   Service   │
        └─────────────┘
```

## 🔧 配置说明

### 环境变量

#### 应用配置
- `APP_ENV`: 应用环境 (development/production)
- `APP_DEBUG`: 调试模式 (true/false)
- `APP_PORT`: 应用端口 (默认: 8080)

#### 数据库配置
- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口 (默认: 5432)
- `DB_DATABASE`: 数据库名称
- `DB_USERNAME`: 数据库用户名
- `DB_PASSWORD`: 数据库密码

#### Redis 配置
- `REDIS_HOST`: Redis 主机
- `REDIS_PORT`: Redis 端口 (默认: 6379)
- `REDIS_PASSWORD`: Redis 密码 (可选)

### 端口配置

#### Docker 部署
- **应用**: 8080
- **Nginx**: 80, 443
- **PostgreSQL**: 5432
- **Redis**: 6379
- **Prometheus**: 9090
- **Grafana**: 3000

#### Kubernetes 部署
- **应用**: 80 (ClusterIP)
- **PostgreSQL**: 5432 (ClusterIP)
- **Redis**: 6379 (ClusterIP)
- **Prometheus**: 9090 (ClusterIP)
- **Grafana**: 3000 (ClusterIP)

## 🚀 生产环境部署

### 1. 安全配置

#### SSL/TLS 证书
```bash
# 生成自签名证书 (仅用于测试)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout examples/deployment_example/docker/ssl/key.pem \
  -out examples/deployment_example/docker/ssl/cert.pem
```

#### 环境变量
```bash
# 生产环境变量
export APP_ENV=production
export APP_DEBUG=false
export DB_PASSWORD=your_secure_password
export REDIS_PASSWORD=your_redis_password
```

### 2. 监控配置

#### Prometheus 告警规则
```yaml
groups:
  - name: laravel-go
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
```

#### Grafana 仪表板
- 应用性能监控
- 数据库性能监控
- 系统资源监控
- 业务指标监控

### 3. 备份策略

#### 数据库备份
```bash
# PostgreSQL 备份脚本
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec laravel-go-postgres pg_dump -U laravel_go laravel_go > backup_$DATE.sql
```

#### 应用备份
```bash
# 应用数据备份
tar -czf app_backup_$(date +%Y%m%d_%H%M%S).tar.gz \
  examples/deployment_example/docker/storage/
```

## 📚 学习要点

### 1. 容器化部署

- Docker 多阶段构建
- 镜像优化
- 容器安全
- 资源限制

### 2. 编排部署

- Kubernetes 资源管理
- 服务发现
- 负载均衡
- 自动扩缩容

### 3. 监控运维

- 指标收集
- 日志聚合
- 告警机制
- 故障排查

### 4. 安全实践

- 最小权限原则
- 网络安全
- 数据加密
- 访问控制

## 🔗 相关文档

- [Docker 官方文档](https://docs.docker.com/)
- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Grafana 官方文档](https://grafana.com/docs/)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个部署示例。

## 📄 许可证

本项目采用 MIT 许可证。 