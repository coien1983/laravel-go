# my-rabbitmq-app

一个基于 Laravel-Go Framework 构建的完整 Web 应用。

## 快速开始

1. 安装依赖: go mod tidy
2. 运行应用: go run main.go
3. 访问: http://localhost:8080

## 项目结构

- app/controllers/ - 控制器
- app/models/ - 数据模型
- config/ - 配置文件
- database/ - 数据库相关
- routes/ - 路由定义
- storage/ - 存储目录

## API 接口

- GET / - 首页
- GET /health - 健康检查
- GET /api/users - 获取用户列表
- POST /api/users - 创建用户

## 开发

使用 largo 命令生成代码:
- largo make:controller ProductController
- largo make:model Product
- largo make:middleware AuthMiddleware

更多信息请参考 Laravel-Go Framework 文档