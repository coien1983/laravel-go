# Largo 使用文档

## 📖 概述

`largo` 是 Laravel-Go Framework 的命令行脚手架工具，提供完整的项目管理和代码生成功能。它受到 Laravel Artisan 的启发，为 Go 开发者提供类似的开发体验。

## 🚀 快速开始

### 安装

```bash
# 克隆项目
git clone git@github.com:coien1983/laravel-go.git

# 进入项目目录
cd laravel-go

# 构建并安装 largo
make install
```

### 验证安装

```bash
# 检查版本
largo version

# 查看所有可用命令
largo list
```

## 📋 命令总览

### 命令分类

| 分类           | 命令                               | 描述               |
| -------------- | ---------------------------------- | ------------------ |
| **项目初始化** | `init`                             | 交互式项目初始化   |
| **代码生成**   | `make:controller`                  | 生成控制器         |
|                | `make:model`                       | 生成模型           |
|                | `make:middleware`                  | 生成中间件         |
|                | `make:migration`                   | 生成迁移文件       |
|                | `make:test`                        | 生成测试文件       |
|                | `make:api`                         | 快速生成 API 组件  |
|                | `make:crud`                        | 快速生成 CRUD 组件 |
| **模块管理**   | `add:module`                       | 添加完整模块       |
|                | `add:service`                      | 添加服务层         |
|                | `add:repository`                   | 添加数据仓库层     |
|                | `add:validator`                    | 添加验证器         |
|                | `add:event`                        | 添加事件和监听器   |
| **部署配置**   | (已移除 Docker 和 Kubernetes 支持) | 保持框架轻量级     |
| **项目维护**   | `cache:clear`                      | 清除缓存           |
|                | `route:list`                       | 列出路由           |
| **项目信息**   | `project:info`                     | 显示项目信息       |
|                | `version`                          | 显示版本信息       |

## 🏗️ 项目初始化

### 基本初始化

```bash
# 交互式初始化项目
largo init

# 指定项目名称
largo init my-project

# 使用选项指定名称
largo init --name=my-project
```

### 预设配置类型

`largo init` 现在提供多种预设配置，简化项目初始化过程：

#### 🚀 预设配置选项

1. **API 服务** (`api`)

   - 前后端分离架构
   - JWT 认证
   - Redis 缓存和队列
   - PostgreSQL 数据库
   - Prometheus 监控
   - 完整的 Docker 和 Kubernetes 支持

2. **Web 应用** (`web`)

   - 传统 MVC 架构
   - Session 认证
   - MySQL 数据库
   - Blade 模板引擎
   - 基础 Docker 支持

3. **微服务** (`microservice`)

   - 微服务架构
   - gRPC 支持
   - Kafka 消息队列
   - 服务发现和分布式队列
   - 完整的容器化支持

4. **全栈应用** (`fullstack`)

   - Vue.js 前端集成
   - 完整功能栈
   - 企业级配置
   - 完整的部署支持

5. **最小化应用** (`minimal`)

   - 基础功能
   - SQLite 数据库
   - 内存缓存和队列
   - 快速原型开发

6. **自定义配置** (`custom`)
   - 手动选择所有选项
   - 完全自定义的项目配置

## 🔧 模块管理

### 添加完整模块

`add:module` 命令可以快速生成一个完整的模块，包含模型、服务、仓库、控制器等组件：

```bash
# 基础模块（包含模型、服务、仓库、API控制器）
largo add:module User

# 包含Web控制器的模块
largo add:module Product --web

# 完整模块（包含所有组件：验证器、事件、监听器、测试）
largo add:module Order --full

# 使用短标志
largo add:module Category -f
```

### 添加服务层

```bash
# 基础服务
largo add:service UserService

# 带接口的服务
largo add:service ProductService --interface
```

### 添加数据仓库层

```bash
# 基础仓库
largo add:repository UserRepository

# 指定模型的仓库
largo add:repository OrderRepository --model=Order

# 带接口的仓库
largo add:repository ProductRepository --model=Product --interface
```

### 添加验证器

```bash
# 基础验证器
largo add:validator UserValidator

# 带验证规则的验证器
largo add:validator ProductValidator --rules="required,min=2,max=100"
```

### 添加事件和监听器

```bash
# 基础事件和监听器
largo add:event UserRegistered

# 队列监听器
largo add:event OrderCreated --queue

# 只生成事件，不生成监听器
largo add:event ProductUpdated --listener=false
```

### 交互式配置选项（自定义模式）

#### 🏗️ 基础架构

- **项目架构**: 单体应用 vs 微服务架构
- **数据库**: SQLite、MySQL、PostgreSQL
- **缓存系统**: Memory、Redis、Memcached
- **队列系统**: Memory、Redis、RabbitMQ、Kafka、SQS、Beanstalkd、Database、etcd、Consul、ZooKeeper
- **前端方案**: API、Blade、Vue.js、React
- **认证方式**: None、JWT、Session
- **API 类型**: REST、GraphQL、Both
- **测试策略**: None、Unit、Integration、Both
- **API 文档**: None、Swagger
- **监控方案**: None、Prometheus
- **日志方案**: File、JSON、Both

#### 🔧 框架功能

- **控制台功能**: Basic、Full、Custom
- **事件系统**: None、Basic、Full
- **数据验证**: None、Basic、Full
- **中间件**: None、Basic、Full
- **路由系统**: Basic、Advanced、Full
- **会话管理**: None、File、Redis、Database
- **邮件系统**: None、SMTP、Mailgun、SendGrid
- **通知系统**: None、Database、Mail、Slack
- **文件存储**: Local、S3、OSS、COS
- **加密功能**: None、Basic、Full
- **密码哈希**: None、Bcrypt、Argon2
- **分页功能**: None、Basic、Advanced
- **限流功能**: None、Basic、Advanced
- **CORS 支持**: None、Basic、Full
- **压缩功能**: None、Gzip、Brotli
- **WebSocket**: None、Basic、Full
- **任务调度**: None、Basic、Full
- **定时器**: None、Cron、Interval、Full
- **健康检查**: None、Basic、Full
- **指标监控**: None、Basic、Prometheus
- **性能分析**: None、Pprof、Full
- **国际化**: None、Basic、Full
- **本地化**: None、Basic、Full

## 🔧 代码生成

### 基础生成命令

#### 生成控制器

```bash
# 基本控制器
largo make:controller User

# 指定命名空间
largo make:controller User --namespace=api

# 使用短选项
largo make:controller User -n api
```

#### 生成模型

```bash
# 基本模型
largo make:model User

# 指定字段
largo make:model User --fields=name:string,email:string,age:int

# 使用短选项
largo make:model User -f name:string,email:string,age:int
```

#### 生成中间件

```bash
# 生成中间件
largo make:middleware Auth

# 生成认证中间件
largo make:middleware Cors
```

#### 生成迁移文件

```bash
# 基本迁移
largo make:migration create_users_table

# 指定表名
largo make:migration create_users_table --table=users

# 指定字段
largo make:migration create_users_table --table=users --fields=name:string,email:string,age:int
```

#### 生成测试文件

```bash
# 基本测试
largo make:test User

# 指定测试类型
largo make:test User --type=unit

# 生成集成测试
largo make:test User --type=integration
```

### 快速生成命令

#### make:api - 快速生成 API 组件

```bash
# 生成用户 API 组件
largo make:api user --fields=name:string,email:string,age:int

# 生成产品 API 组件
largo make:api product --fields=name:string,price:decimal,description:text

# 简单模式（不指定字段）
largo make:api user
```

**生成的文件**:

- `app/controllers/user_controller.go` - API 控制器
- `app/models/user.go` - 用户模型
- `database/migrations/xxx_create_users_table.sql` - 数据库迁移

#### make:crud - 快速生成 CRUD 组件

```bash
# 生成用户 CRUD 组件
largo make:crud user --fields=name:string,email:string,age:int

# 生成产品 CRUD 组件
largo make:crud product --fields=name:string,price:decimal,description:text

# 简单模式（不指定字段）
largo make:crud user
```

**生成的文件**:

- `app/controllers/user_controller.go` - CRUD 控制器
- `app/models/user.go` - 用户模型
- `database/migrations/xxx_create_users_table.sql` - 数据库迁移
- `tests/user_test.go` - 单元测试
- `tests/user_test.go` - 集成测试

## 🐳 部署配置

### 部署配置

Docker 和 Kubernetes 支持已被移除，以保持框架的轻量级和专注于核心功能。

## 🛠️ 项目维护

### 缓存管理

```bash
# 清除应用缓存
largo cache:clear
```

### 路由管理

```bash
# 列出所有注册的路由
largo route:list
```

## 📊 项目信息

### 版本信息

```bash
# 显示版本信息
largo version
```

**输出示例**:

```
Laravel-Go Framework v1.0.0
A modern Go web framework inspired by Laravel
GitHub: https://github.com/coien1983/laravel-go
```

### 项目信息

```bash
# 显示项目信息
largo project:info
```

**输出示例**:

```
Laravel-Go Framework 项目信息:
  应用名称: laravel-go-app
  默认端口: 8080
  默认命名空间: default
  默认副本数: 3

可用命令:
  largo list          - 显示所有命令
  largo init          - 初始化项目
  largo make:controller - 生成控制器
  largo make:model    - 生成模型
  # Docker 和 Kubernetes 支持已移除
  largo make:api      - 快速生成 API 组件
  largo make:crud     - 快速生成 CRUD 组件
```

## 📝 Makefile 支持

### 基础操作

```bash
# 显示帮助信息
make help

# 构建 largo 可执行文件
make build

# 安装到 Go bin 目录
make install

# 显示脚手架工具帮助
make run
```

### 项目初始化

```bash
# 交互式初始化项目
make init

# 使用自定义名称初始化
make init-custom
```

### 代码生成

```bash
# 生成控制器
make controller
make controller-custom

# 生成模型
make model
make model-fields

# 生成中间件
make middleware

# 生成迁移文件
make migration
make migration-table
make migration-fields

# 生成测试文件
make test
make test-type
```

### 快速生成

```bash
# 快速生成 API 组件
make api
make api-simple

# 快速生成 CRUD 组件
make crud
make crud-simple
```

### 部署配置

```bash
# 生成 Docker 配置
make docker
make docker-custom

# 生成 Kubernetes 配置
make k8s
make k8s-custom
```

### 项目维护

```bash
# 列出所有路由
make routes

# 清除应用缓存
make cache-clear

# 列出所有可用命令
make list

# 显示版本信息
make version

# 显示项目信息
make info
```

### 开发工具

```bash
# 设置开发环境
make dev-setup

# 运行所有测试
make test-all

# 运行测试并生成覆盖率报告
make test-coverage

# 代码检查
make lint

# 格式化代码
make fmt

# 代码静态分析
make vet
```

### Docker 操作

```bash
# 构建 Docker 镜像
make docker-build

# 运行 Docker 容器
make docker-run

# 启动 Docker Compose 服务
make docker-compose-up

# 停止 Docker Compose 服务
make docker-compose-down

# 查看 Docker Compose 日志
make docker-compose-logs
```

### Kubernetes 操作

```bash
# 部署到 Kubernetes
make k8s-apply

# 从 Kubernetes 删除
make k8s-delete

# 查看 Kubernetes 部署状态
make k8s-status

# 查看 Kubernetes 日志
make k8s-logs
```

### 清理操作

```bash
# 清理构建文件
make clean

# 清理 Docker 文件
make clean-docker

# 清理 Kubernetes 文件
make clean-k8s

# 清理所有生成的文件
make clean-all
```

## 🎯 使用示例

### 示例 1: 创建完整的用户管理系统

```bash
# 1. 初始化项目
largo init user-management

# 2. 进入项目目录
cd user-management

# 3. 生成用户 CRUD 组件
largo make:crud user --fields=name:string,email:string,password:string,age:int

# 4. 生成认证中间件
largo make:middleware auth

# 5. Docker 和 Kubernetes 支持已移除
```

### 示例 2: 创建 API 服务

```bash
# 1. 初始化 API 项目
largo init api-service

# 2. 进入项目目录
cd api-service

# 3. 生成产品 API
largo make:api product --fields=name:string,price:decimal,description:text

# 4. 生成订单 API
largo make:api order --fields=user_id:int,product_id:int,quantity:int,total:decimal

# 5. Docker 和 Kubernetes 支持已移除
```

### 示例 3: 使用 Makefile 快速开发

```bash
# 1. 初始化项目
make init

# 2. 快速生成 API 组件
make api

# 3. 快速生成 CRUD 组件
make crud

# 4. 部署配置支持已移除

# 5. 运行测试
make test-all

# 6. 代码检查
make lint
make fmt
```

## 🔧 高级用法

### 自定义字段类型

支持以下字段类型：

| 类型       | 描述      | 示例                  |
| ---------- | --------- | --------------------- |
| `string`   | 字符串    | `name:string`         |
| `int`      | 整数      | `age:int`             |
| `decimal`  | 小数      | `price:decimal`       |
| `text`     | 长文本    | `description:text`    |
| `boolean`  | 布尔值    | `is_active:boolean`   |
| `datetime` | 日期时间  | `created_at:datetime` |
| `json`     | JSON 数据 | `metadata:json`       |

### 字段修饰符

```bash
# 可空字段
largo make:model User --fields=name:string,email:string?,age:int?

# 带默认值的字段
largo make:model User --fields=name:string,status:string:active,created_at:datetime:now
```

### 批量生成

```bash
# 生成多个控制器
largo make:controller User
largo make:controller Product
largo make:controller Order

# 生成多个模型
largo make:model User --fields=name:string,email:string
largo make:model Product --fields=name:string,price:decimal
largo make:model Order --fields=user_id:int,total:decimal
```

## 🐛 故障排除

### 常见问题

#### 1. 命令未找到

```bash
# 确保 largo 已正确安装
which largo

# 重新安装
make install
```

#### 2. 权限问题

```bash
# 确保有执行权限
chmod +x bin/largo

# 重新构建
make build
```

#### 3. 依赖问题

```bash
# 更新依赖
go mod tidy
go mod download
```

#### 4. 生成文件失败

```bash
# 检查目录权限
ls -la

# 确保目录存在
mkdir -p app/controllers app/models database/migrations tests
```

### 调试模式

```bash
# 启用详细输出
DEBUG=1 largo init

# 查看详细错误信息
largo --verbose init
```

## 📚 最佳实践

### 1. 项目结构

```
my-project/
├── app/
│   ├── controllers/
│   ├── models/
│   └── middleware/
├── database/
│   └── migrations/
├── tests/
├── config/
├── routes/
├── Dockerfile
├── docker-compose.yml
├── k8s/
└── README.md
```

### 2. 命名约定

- **控制器**: 使用 PascalCase，以 `Controller` 结尾
- **模型**: 使用 PascalCase，单数形式
- **中间件**: 使用 PascalCase，描述性名称
- **迁移**: 使用 snake_case，描述性名称
- **测试**: 使用与源文件相同的名称，以 `_test.go` 结尾

### 3. 字段设计

```bash
# 好的字段设计
largo make:model User --fields=name:string,email:string:unique,password:string,age:int,is_active:boolean,created_at:datetime

# 避免过度设计
largo make:model User --fields=id:int,name:string,email:string
```

### 4. 测试策略

```bash
# 为每个模型生成测试
largo make:test User --type=unit
largo make:test User --type=integration

# 为 API 端点生成测试
largo make:test UserController --type=integration
```

## 🔗 相关资源

- [Laravel-Go Framework 主页](https://github.com/coien1983/laravel-go)
- [Queue vs Scheduler 对比](docs/guides/queue-vs-scheduler.md)
- [API 文档](docs/api/)
- [部署指南](docs/deployment/)

## 📞 支持

- **GitHub Issues**: [报告问题](https://github.com/coien1983/laravel-go/issues)
- **邮箱支持**: coien1983@126.com
- **文档**: [完整文档](https://github.com/coien1983/laravel-go/tree/main/docs)

---

**Laravel-Go Framework** - 受 Laravel 启发的现代 Go Web 框架 🚀
