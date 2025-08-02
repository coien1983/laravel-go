# Laravel-Go Framework 集成测试系统

## 概述

本目录包含了 Laravel-Go Framework 的完整集成测试系统，提供了全面的测试覆盖，包括 HTTP、数据库、缓存、队列等核心功能的集成测试。

## 目录结构

```
tests/
├── README.md                    # 测试系统说明文档
├── Makefile                     # 测试运行脚本
├── test_config.go              # 测试配置管理
├── test_helpers.go             # 测试辅助工具
├── run_integration_tests.go    # 测试运行器
├── integration_test.go         # 综合集成测试
├── http_integration_test.go    # HTTP集成测试
├── database_integration_test.go # 数据库集成测试
├── cache_queue_integration_test.go # 缓存队列集成测试
└── testdata/                   # 测试数据目录
```

## 快速开始

### 1. 安装依赖

```bash
# 安装测试依赖
make install-deps

# 或者手动安装
go get github.com/stretchr/testify/suite
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
```

### 2. 运行所有测试

```bash
# 运行所有集成测试
make test-all

# 或者使用 go test
go test -v ./tests/... -run "Test.*Integration.*Suite"
```

### 3. 运行特定测试

```bash
# 运行HTTP集成测试
make test-http

# 运行数据库集成测试
make test-database

# 运行缓存队列集成测试
make test-cache-queue

# 运行单元测试
make test-unit
```

### 4. 生成测试报告

```bash
# 生成HTML测试报告
make report

# 生成覆盖率报告
make coverage
```

## 测试套件说明

### 1. 综合集成测试 (integration_test.go)

测试框架各个组件之间的集成，包括：

- HTTP 请求处理
- 数据库操作
- 缓存系统
- 队列处理
- 事件系统
- 认证授权
- 中间件功能
- 数据验证

### 2. HTTP 集成测试 (http_integration_test.go)

专门测试 HTTP 相关功能：

- 基础路由功能
- 路由参数处理
- 路由分组
- 中间件链式调用
- 控制器 CRUD 操作
- 验证中间件
- 错误处理
- 内容类型处理
- 请求头处理

### 3. 数据库集成测试 (database_integration_test.go)

测试数据库和 ORM 功能：

- 基础 CRUD 操作
- 查询构建器
- 模型关联关系
- 多对多关联
- 软删除功能
- 模型钩子
- 事务处理
- 原始查询
- 模型验证

### 4. 缓存队列集成测试 (cache_queue_integration_test.go)

测试缓存和队列系统：

- 缓存基础操作
- 缓存过期机制
- 缓存标签
- 批量操作
- 队列基础操作
- 延迟任务
- 批量任务处理
- 统计功能

## 配置管理

### 环境变量配置

可以通过环境变量来配置测试环境：

```bash
# 数据库配置
export TEST_DB_DRIVER=sqlite
export TEST_DB_DATABASE=:memory:
export TEST_DB_HOST=localhost
export TEST_DB_PORT=3306
export TEST_DB_USERNAME=root
export TEST_DB_PASSWORD=password

# 缓存配置
export TEST_CACHE_DRIVER=memory
export TEST_CACHE_HOST=localhost
export TEST_CACHE_PORT=6379
export TEST_CACHE_DB=0

# 队列配置
export TEST_QUEUE_DRIVER=memory
export TEST_QUEUE_HOST=localhost
export TEST_QUEUE_PORT=6379
export TEST_QUEUE_DB=1

# HTTP配置
export TEST_HTTP_PORT=8080
export TEST_HTTP_TIMEOUT=30s

# 测试配置
export TEST_PARALLEL=false
export TEST_TIMEOUT=60s
export TEST_RETRY_COUNT=3
export TEST_CLEANUP_DATA=true

# 日志配置
export TEST_LOG_LEVEL=info
export TEST_LOG_OUTPUT=stdout
```

### 配置文件

可以创建 `tests/test_config.json` 文件来配置测试：

```json
{
  "database": {
    "driver": "sqlite",
    "database": ":memory:",
    "charset": "utf8mb4"
  },
  "cache": {
    "driver": "memory",
    "host": "localhost",
    "port": 6379,
    "db": 0
  },
  "queue": {
    "driver": "memory",
    "host": "localhost",
    "port": 6379,
    "db": 1
  },
  "http": {
    "port": 8080,
    "timeout": "30s"
  },
  "test": {
    "parallel": false,
    "timeout": "60s",
    "retry_count": 3,
    "cleanup_data": true
  },
  "log": {
    "level": "info",
    "output": "stdout"
  }
}
```

## 测试辅助工具

### TestHelper 类

提供了丰富的测试辅助方法：

```go
func TestExample(t *testing.T) {
    helper := NewTestHelper(t)

    // 创建临时文件
    tempFile := helper.CreateTempFile("test content")
    defer helper.CleanupTempFile(tempFile)

    // 创建测试数据库
    db := helper.CreateTestDatabase()

    // 创建测试表
    helper.CreateTestTable(db, "users", `
        CREATE TABLE users (
            id INTEGER PRIMARY KEY,
            name VARCHAR(255),
            email VARCHAR(255)
        )
    `)

    // 插入测试数据
    helper.InsertTestData(db, "users", map[string]interface{}{
        "name": "John Doe",
        "email": "john@example.com",
    })

    // 创建HTTP请求
    req := helper.CreateHTTPRequest("POST", "/users", map[string]string{
        "name": "Jane Doe",
        "email": "jane@example.com",
    })

    // 执行HTTP请求
    recorder := helper.ExecuteHTTPRequest(handler, req)

    // 断言响应
    helper.AssertJSONResponse(recorder, http.StatusCreated, map[string]interface{}{
        "message": "User created",
    })
}
```

## 常用命令

### Makefile 命令

```bash
# 查看所有可用命令
make help

# 运行所有集成测试
make test-all

# 运行特定测试套件
make test-http
make test-database
make test-cache-queue

# 运行单元测试
make test-unit

# 生成测试报告
make report

# 生成覆盖率报告
make coverage

# 清理测试文件
make clean

# 安装依赖
make install-deps

# 运行性能测试
make benchmark

# 检查代码质量
make lint

# 格式化代码
make format
```

### Go Test 命令

```bash
# 运行所有测试
go test -v ./tests/...

# 运行特定测试套件
go test -v ./tests/ -run "TestHTTPIntegrationTestSuite"
go test -v ./tests/ -run "TestDatabaseIntegrationTestSuite"
go test -v ./tests/ -run "TestCacheQueueIntegrationTestSuite"

# 运行特定测试方法
go test -v ./tests/ -run "TestBasicRouting"

# 生成覆盖率报告
go test -v -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out -o coverage.html

# 运行基准测试
go test -v -bench=. -benchmem ./tests/...

# 并行运行测试
go test -v -parallel 4 ./tests/...

# 设置测试超时
go test -v -timeout 5m ./tests/...
```

## 测试最佳实践

### 1. 测试结构

- 使用 `suite.Suite` 作为测试套件基类
- 在 `SetupSuite()` 中初始化测试环境
- 在 `SetupTest()` 中准备每个测试的数据
- 在 `TearDownSuite()` 中清理资源

### 2. 测试数据管理

- 使用内存数据库进行测试
- 每个测试前清理数据
- 使用工厂函数创建测试数据
- 避免测试间的数据依赖

### 3. 断言使用

- 使用 `assert` 进行非致命断言
- 使用 `require` 进行致命断言
- 提供清晰的错误消息
- 测试边界条件和异常情况

### 4. 性能考虑

- 使用 `-short` 标志跳过慢速测试
- 合理设置测试超时时间
- 避免不必要的网络调用
- 使用模拟对象替代外部依赖

### 5. 测试覆盖率

- 定期检查测试覆盖率
- 重点关注核心业务逻辑
- 测试错误处理路径
- 测试边界条件

## 故障排除

### 常见问题

1. **测试失败**

   - 检查测试环境配置
   - 确保依赖服务可用
   - 查看详细的错误日志

2. **性能问题**

   - 使用 `-short` 标志跳过慢速测试
   - 检查测试超时设置
   - 优化测试数据准备

3. **环境问题**
   - 确保 Go 版本兼容
   - 检查依赖包版本
   - 清理测试缓存

### 调试技巧

1. **启用详细日志**

   ```bash
   go test -v -log.level=debug ./tests/...
   ```

2. **运行单个测试**

   ```bash
   go test -v -run "TestSpecificFunction" ./tests/...
   ```

3. **使用测试标签**
   ```bash
   go test -v -tags=integration ./tests/...
   ```

## 贡献指南

### 添加新测试

1. 创建测试文件
2. 继承适当的测试套件
3. 实现测试方法
4. 添加必要的断言
5. 更新文档

### 测试命名规范

- 测试文件：`*_test.go`
- 测试套件：`*TestSuite`
- 测试方法：`Test*`
- 基准测试：`Benchmark*`

### 代码审查

- 确保测试覆盖核心功能
- 验证测试的可重复性
- 检查测试的性能影响
- 确保测试文档完整

## 许可证

本测试系统遵循与主项目相同的许可证。
