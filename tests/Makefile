# Laravel-Go Framework 集成测试 Makefile

.PHONY: help test test-all test-http test-database test-cache-queue test-unit clean coverage report

# 默认目标
help:
	@echo "Laravel-Go Framework 集成测试管理"
	@echo ""
	@echo "可用命令:"
	@echo "  make test-all        - 运行所有集成测试"
	@echo "  make test-http       - 运行HTTP集成测试"
	@echo "  make test-database   - 运行数据库集成测试"
	@echo "  make test-cache-queue - 运行缓存队列集成测试"
	@echo "  make test-unit       - 运行单元测试"
	@echo "  make coverage        - 生成测试覆盖率报告"
	@echo "  make report          - 生成HTML测试报告"
	@echo "  make clean           - 清理测试文件"
	@echo "  make install-deps    - 安装测试依赖"

# 运行所有集成测试
test-all:
	@echo "🚀 运行所有集成测试..."
	go test -v ./tests/... -run "Test.*Integration.*Suite"

# 运行HTTP集成测试
test-http:
	@echo "🌐 运行HTTP集成测试..."
	go test -v ./tests/ -run "TestHTTPIntegrationTestSuite"

# 运行数据库集成测试
test-database:
	@echo "🗄️  运行数据库集成测试..."
	go test -v ./tests/ -run "TestDatabaseIntegrationTestSuite"

# 运行缓存队列集成测试
test-cache-queue:
	@echo "💾 运行缓存队列集成测试..."
	go test -v ./tests/ -run "TestCacheQueueIntegrationTestSuite"

# 运行单元测试
test-unit:
	@echo "🧪 运行单元测试..."
	go test -v ../framework/... -short

# 生成测试覆盖率报告
coverage:
	@echo "📊 生成测试覆盖率报告..."
	go test -v -coverprofile=coverage.out ./... ../framework/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 生成HTML测试报告
report:
	@echo "📄 生成HTML测试报告..."
	go run tests/run_integration_tests.go

# 清理测试文件
clean:
	@echo "🧹 清理测试文件..."
	rm -rf test_reports/
	rm -f coverage.out coverage.html
	rm -f *.db
	rm -f *.sqlite
	@echo "清理完成"

# 安装测试依赖
install-deps:
	@echo "📦 安装测试依赖..."
	go get github.com/stretchr/testify/suite
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/require
	@echo "依赖安装完成"

# 运行性能测试
benchmark:
	@echo "⚡ 运行性能测试..."
	go test -v -bench=. -benchmem ./... ../framework/...

# 运行特定测试
test-specific:
	@echo "🎯 运行特定测试..."
	@read -p "请输入测试套件名称: " suite; \
	read -p "请输入测试名称 (可选): " test; \
	if [ -z "$$test" ]; then \
		go test -v ./tests/ -run "$$suite"; \
	else \
		go test -v ./tests/ -run "$$test"; \
	fi

# 检查代码质量
lint:
	@echo "🔍 检查代码质量..."
	golangci-lint run ./... ../framework/...

# 格式化代码
format:
	@echo "🎨 格式化代码..."
	go fmt ./... ../framework/...
	go vet ./... ../framework/...

# 运行测试并生成报告
test-with-report: test-all report
	@echo "✅ 测试完成，报告已生成"

# 快速测试（跳过慢速测试）
test-quick:
	@echo "⚡ 运行快速测试..."
	go test -v -short ./... ../framework/...

# 显示测试帮助
test-help:
	@echo "测试命令帮助:"
	@echo "  -v          : 详细输出"
	@echo "  -short      : 跳过慢速测试"
	@echo "  -run        : 运行特定测试"
	@echo "  -cover      : 生成覆盖率"
	@echo "  -bench      : 运行基准测试"
	@echo "  -timeout    : 设置超时时间" 