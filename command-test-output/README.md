# 指令测试输出目录

这个目录用于存放 `largo` 指令的测试输出文件。

## 用途

- 存放 `make:rpc` 命令生成的RPC服务文件
- 存放 `make:api` 命令生成的API服务文件
- 存放 `make:microservice` 命令生成的微服务文件
- 存放其他指令的测试输出

## 使用方法

```bash
# 测试 make:rpc 命令
largo make:rpc examples/gozero/user.proto --output=command-test-output/rpc-test

# 测试 make:api 命令
largo make:api examples/gozero/user.api --output=command-test-output/api-test

# 测试 make:microservice 命令
largo make:microservice user examples/gozero/user.proto examples/gozero/user.api --output=command-test-output/microservice-test
```

## 注意事项

- 此目录已被 `.gitignore` 忽略，不会提交到版本控制
- 测试完成后可以安全删除此目录中的文件
- 建议定期清理此目录以节省磁盘空间 