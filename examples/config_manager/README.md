# Laravel-Go 配置管理工具

## 📝 项目概览

这是一个用于管理 Laravel-Go Framework 配置的命令行工具，支持配置的读取、设置、验证和格式转换。

## 🚀 功能特性

- ✅ 配置读取和设置
- ✅ 配置文件加载
- ✅ 环境变量支持
- ✅ 配置验证
- ✅ 多格式输出 (JSON, YAML, ENV)
- ✅ 嵌套配置支持
- ✅ 类型安全

## 🚀 快速开始

### 1. 编译工具

```bash
cd examples/config_manager
go build -o config-manager main.go
```

### 2. 基本使用

```bash
# 查看帮助
./config-manager -h

# 获取配置值
./config-manager -key="app.name"

# 设置配置值
./config-manager -action=set -key="app.name" -value="My App"

# 列出所有配置
./config-manager -action=list

# 验证配置
./config-manager -action=validate
```

### 3. 加载配置文件

```bash
# 从JSON文件加载配置
./config-manager -config=config.json -action=list

# 从YAML文件加载配置
./config-manager -config=config.yaml -action=list

# 从环境变量文件加载配置
./config-manager -config=.env -action=list
```

### 4. 格式输出

```bash
# JSON格式输出
./config-manager -action=list -format=json

# YAML格式输出
./config-manager -action=list -format=yaml

# 环境变量格式输出
./config-manager -action=list -format=env
```

## 📋 命令行参数

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `-config` | string | 否 | 配置文件路径 |
| `-key` | string | 否 | 配置键 |
| `-value` | string | 否 | 配置值 |
| `-action` | string | 否 | 操作类型 (get, set, list, validate) |
| `-format` | string | 否 | 输出格式 (json, yaml, env) |

## 🔧 操作类型

### 1. get - 获取配置值

```bash
# 获取简单配置
./config-manager -action=get -key="app.name"

# 获取嵌套配置
./config-manager -action=get -key="database.connections.postgres.host"

# 指定输出格式
./config-manager -action=get -key="app" -format=json
```

### 2. set - 设置配置值

```bash
# 设置简单配置
./config-manager -action=set -key="app.name" -value="My Application"

# 设置嵌套配置
./config-manager -action=set -key="database.host" -value="localhost"
```

### 3. list - 列出所有配置

```bash
# 列出所有配置
./config-manager -action=list

# 指定输出格式
./config-manager -action=list -format=yaml
```

### 4. validate - 验证配置

```bash
# 验证配置
./config-manager -action=validate
```

## 📁 配置文件格式

### 1. JSON 格式

```json
{
  "app": {
    "name": "Laravel-Go App",
    "version": "1.0.0",
    "debug": true,
    "port": 8080
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "database": "laravel_go",
    "username": "user",
    "password": "password"
  }
}
```

### 2. YAML 格式

```yaml
app:
  name: Laravel-Go App
  version: 1.0.0
  debug: true
  port: 8080

database:
  host: localhost
  port: 5432
  database: laravel_go
  username: user
  password: password
```

### 3. 环境变量格式

```env
APP_NAME="Laravel-Go App"
APP_VERSION=1.0.0
APP_DEBUG=true
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=laravel_go
DB_USERNAME=user
DB_PASSWORD=password
```

## 🔍 配置验证规则

工具支持以下验证规则：

- `required`: 必填字段
- `numeric`: 数字类型
- `boolean`: 布尔类型
- `string`: 字符串类型
- `email`: 邮箱格式
- `url`: URL格式

### 验证示例

```bash
# 验证特定配置
./config-manager -action=validate -key="app.name" -rule="required"

# 验证多个配置
./config-manager -action=validate
```

## 🛠️ 高级用法

### 1. 批量操作

```bash
# 批量设置配置
./config-manager -action=set -key="app.name" -value="App1" && \
./config-manager -action=set -key="app.version" -value="2.0.0" && \
./config-manager -action=set -key="app.debug" -value="false"
```

### 2. 配置转换

```bash
# JSON转环境变量
./config-manager -config=config.json -action=list -format=env > .env

# 环境变量转JSON
./config-manager -config=.env -action=list -format=json > config.json
```

### 3. 配置比较

```bash
# 比较两个配置文件
diff <(./config-manager -config=config1.json -action=list -format=json) \
     <(./config-manager -config=config2.json -action=list -format=json)
```

## 📚 使用场景

### 1. 开发环境

```bash
# 加载开发环境配置
./config-manager -config=config/dev.json -action=list

# 设置开发环境变量
./config-manager -action=set -key="app.debug" -value="true"
```

### 2. 生产环境

```bash
# 加载生产环境配置
./config-manager -config=config/prod.json -action=validate

# 设置生产环境变量
./config-manager -action=set -key="app.debug" -value="false"
```

### 3. 测试环境

```bash
# 加载测试环境配置
./config-manager -config=config/test.json -action=list

# 验证测试配置
./config-manager -config=config/test.json -action=validate
```

## 🔧 集成到项目

### 1. 在脚本中使用

```bash
#!/bin/bash

# 获取应用名称
APP_NAME=$(./config-manager -action=get -key="app.name")

# 获取数据库配置
DB_HOST=$(./config-manager -action=get -key="database.host")
DB_PORT=$(./config-manager -action=get -key="database.port")

echo "应用名称: $APP_NAME"
echo "数据库主机: $DB_HOST:$DB_PORT"
```

### 2. 在CI/CD中使用

```yaml
# GitHub Actions 示例
- name: 验证配置
  run: |
    ./config-manager -config=config/prod.json -action=validate

- name: 生成环境变量
  run: |
    ./config-manager -config=config/prod.json -action=list -format=env > .env
```

## 🚨 错误处理

### 1. 常见错误

```bash
# 配置文件不存在
./config-manager -config=not-exist.json -action=list
# 错误: 加载配置文件失败: open not-exist.json: no such file or directory

# 配置键不存在
./config-manager -action=get -key="not.exist"
# 输出: <nil>

# 验证失败
./config-manager -action=validate
# 错误: 配置验证失败: app.name is required
```

### 2. 调试模式

```bash
# 启用详细输出
DEBUG=true ./config-manager -action=list
```

## 📄 许可证

本项目采用 MIT 许可证。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个配置管理工具。 