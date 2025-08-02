package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"laravel-go/framework/microservice"
)

// 示例：使用高级拦截器
func useAdvancedInterceptors() {
	fmt.Println("=== 高级拦截器示例 ===")

	// 创建缓存
	cache := microservice.NewSimpleCache()

	// 创建退避策略
	backoff := microservice.NewExponentialBackoff(time.Second, 30*time.Second)

	// 创建节流器
	throttler := microservice.NewSimpleThrottler(map[string]int{
		"/user.UserService/GetUser":    100,
		"/user.UserService/CreateUser": 50,
	})

	// 创建安全检查器
	securityChecker := microservice.NewSimpleSecurityChecker(map[string]bool{
		"/user.UserService/GetUser":    true,
		"/user.UserService/CreateUser": true,
	})

	// 创建审计器
	auditor := microservice.NewSimpleAuditor(log.New(log.Writer(), "AUDIT: ", log.LstdFlags))

	// 创建性能监控器
	perfMonitor := microservice.NewSimplePerformanceMonitor()

	// 创建熔断器
	circuitBreaker := microservice.NewSimpleCircuitBreaker(5, 30*time.Second)

	// 创建拦截器
	cacheInterceptor := microservice.CacheInterceptor(cache, 5*time.Minute)
	retryInterceptor := microservice.RetryInterceptor(3, backoff)
	throttlingInterceptor := microservice.ThrottlingInterceptor(throttler)
	securityInterceptor := microservice.SecurityInterceptor(securityChecker)
	auditInterceptor := microservice.AuditInterceptor(auditor)
	perfInterceptor := microservice.PerformanceInterceptor(perfMonitor)
	circuitBreakerInterceptor := microservice.CircuitBreakerInterceptor(circuitBreaker)

	fmt.Printf("创建了 %d 个高级拦截器\n", 7)
	fmt.Printf("缓存拦截器: %v\n", cacheInterceptor != nil)
	fmt.Printf("重试拦截器: %v\n", retryInterceptor != nil)
	fmt.Printf("节流拦截器: %v\n", throttlingInterceptor != nil)
	fmt.Printf("安全拦截器: %v\n", securityInterceptor != nil)
	fmt.Printf("审计拦截器: %v\n", auditInterceptor != nil)
	fmt.Printf("性能拦截器: %v\n", perfInterceptor != nil)
	fmt.Printf("熔断器拦截器: %v\n", circuitBreakerInterceptor != nil)
}

// 示例：使用流处理功能
func useStreamProcessing() {
	fmt.Println("\n=== 流处理功能示例 ===")

	// 创建流转换器
	transformer := microservice.NewStreamTransformer()
	transformer.RegisterTransformer("json", microservice.JSONTransformer())
	transformer.RegisterTransformer("base64", microservice.Base64Transformer())

	// 创建流过滤器
	filter := microservice.NewStreamFilter()
	filter.RegisterFilter("size", microservice.SizeFilter(1024))
	filter.RegisterFilter("type", microservice.TypeFilter("string", "[]byte"))

	// 创建流聚合器
	aggregator := microservice.NewStreamAggregator()
	aggregator.RegisterAggregator("count", microservice.CountAggregator())
	aggregator.RegisterAggregator("sum", microservice.SumAggregator())
	aggregator.RegisterAggregator("average", microservice.AverageAggregator())

	// 创建流限流器
	rateLimiter := microservice.NewStreamRateLimiter()
	rateLimiter.AddLimiter("chat", 100, 10.0) // 容量100，速率10/s

	// 创建流缓冲区
	buffer := microservice.NewStreamBuffer()
	buffer.AddBuffer("messages", 1000)

	// 创建流分区器
	partitioner := microservice.NewStreamPartitioner()
	partitioner.RegisterPartition("user", func(data interface{}) (string, error) {
		// 根据用户ID分区
		return "partition1", nil
	})

	// 创建流调度器
	scheduler := microservice.NewStreamScheduler()
	scheduler.RegisterScheduler("priority", func(data interface{}) (int, error) {
		// 根据数据优先级返回优先级
		return 1, nil
	}, 100)

	// 创建流验证器
	validator := microservice.NewStreamValidator()
	validator.RegisterValidator("range", microservice.RangeValidator(0, 100))

	// 创建流丰富器
	enricher := microservice.NewStreamEnricher()
	enricher.RegisterEnricher("timestamp", func(data interface{}) (interface{}, error) {
		// 添加时间戳
		return map[string]interface{}{
			"data":      data,
			"timestamp": time.Now(),
		}, nil
	})

	// 创建流路由器
	router := microservice.NewStreamRouter()
	router.RegisterRoute("type", func(data interface{}) (string, error) {
		// 根据数据类型路由
		return "route1", nil
	})

	fmt.Printf("创建了 %d 个流处理组件\n", 9)
	fmt.Printf("流转换器: %v\n", transformer != nil)
	fmt.Printf("流过滤器: %v\n", filter != nil)
	fmt.Printf("流聚合器: %v\n", aggregator != nil)
	fmt.Printf("流限流器: %v\n", rateLimiter != nil)
	fmt.Printf("流缓冲区: %v\n", buffer != nil)
	fmt.Printf("流分区器: %v\n", partitioner != nil)
	fmt.Printf("流调度器: %v\n", scheduler != nil)
	fmt.Printf("流验证器: %v\n", validator != nil)
	fmt.Printf("流丰富器: %v\n", enricher != nil)
	fmt.Printf("流路由器: %v\n", router != nil)
}

// 示例：流处理演示
func demonstrateStreamProcessing() {
	fmt.Println("\n=== 流处理演示 ===")

	// 创建流转换器
	transformer := microservice.NewStreamTransformer()
	transformer.RegisterTransformer("json", microservice.JSONTransformer())

	// 创建流聚合器
	aggregator := microservice.NewStreamAggregator()
	aggregator.RegisterAggregator("count", microservice.CountAggregator())
	aggregator.RegisterAggregator("sum", microservice.SumAggregator())

	// 创建流缓冲区
	buffer := microservice.NewStreamBuffer()
	buffer.AddBuffer("numbers", 10)

	// 模拟数据流
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 处理数据流
	for _, num := range numbers {
		// 转换数据
		transformed, err := transformer.Transform("json", num)
		if err != nil {
			fmt.Printf("转换失败: %v\n", err)
			continue
		}

		// 添加到缓冲区
		if !buffer.Push("numbers", num) {
			fmt.Printf("缓冲区已满，跳过: %d\n", num)
			continue
		}

		// 添加到聚合器
		aggregator.AddData("count", num)
		aggregator.AddData("sum", num)

		fmt.Printf("处理数据: %d -> %s\n", num, transformed)
	}

	// 获取聚合结果
	count, err := aggregator.Aggregate("count")
	if err != nil {
		fmt.Printf("聚合计数失败: %v\n", err)
	} else {
		fmt.Printf("总数: %v\n", count)
	}

	sum, err := aggregator.Aggregate("sum")
	if err != nil {
		fmt.Printf("聚合求和失败: %v\n", err)
	} else {
		fmt.Printf("总和: %v\n", sum)
	}

	// 从缓冲区读取数据
	fmt.Println("缓冲区数据:")
	for {
		data, ok := buffer.Pop("numbers")
		if !ok {
			break
		}
		fmt.Printf("  %v\n", data)
	}
}

// 示例：高级健康检查
func useAdvancedHealthCheck() {
	fmt.Println("\n=== 高级健康检查示例 ===")

	// 创建健康检查配置
	config := microservice.NewHealthConfig()
	config.Interval = 10 * time.Second
	config.Timeout = 5 * time.Second

	// 创建健康检查服务
	healthService := microservice.NewGRPCHealthService(config)

	// 注册自定义健康检查
	healthService.RegisterHealthCheck("custom", func(ctx context.Context) error {
		// 模拟自定义健康检查
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	// 注册依赖服务健康检查
	healthService.RegisterHealthCheck("database", func(ctx context.Context) error {
		// 模拟数据库健康检查
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	healthService.RegisterHealthCheck("redis", func(ctx context.Context) error {
		// 模拟Redis健康检查
		time.Sleep(30 * time.Millisecond)
		return nil
	})

	healthService.RegisterHealthCheck("external-api", func(ctx context.Context) error {
		// 模拟外部API健康检查
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	// 启动健康检查
	healthService.Start()

	// 等待一段时间让健康检查运行
	time.Sleep(2 * time.Second)

	// 获取健康状态
	statuses := healthService.GetAllStatus()
	fmt.Printf("健康检查状态 (%d 个服务):\n", len(statuses))
	for service, status := range statuses {
		fmt.Printf("  %s: %s\n", service, status.Status)
	}

	// 创建健康监控器
	monitor := microservice.NewHealthMonitor(healthService)
	monitor.Start()

	// 获取监控指标
	metrics := monitor.GetMetrics()
	fmt.Printf("监控指标: 总检查=%d, 成功=%d, 失败=%d\n",
		metrics.TotalChecks, metrics.SuccessfulChecks, metrics.FailedChecks)
}

// 示例：性能监控
func demonstratePerformanceMonitoring() {
	fmt.Println("\n=== 性能监控演示 ===")

	// 创建性能监控器
	perfMonitor := microservice.NewSimplePerformanceMonitor()

	// 模拟一些操作
	operations := []string{
		"/user.UserService/GetUser",
		"/user.UserService/CreateUser",
		"/user.UserService/UpdateUser",
		"/user.UserService/DeleteUser",
	}

	for _, operation := range operations {
		// 记录开始
		perfMonitor.RecordStart(operation)

		// 模拟操作执行时间
		duration := time.Duration(100+time.Now().UnixNano()%900) * time.Millisecond
		time.Sleep(duration)

		// 模拟一些操作失败
		var err error
		if time.Now().UnixNano()%10 == 0 {
			err = fmt.Errorf("simulated error")
		}

		// 记录结束
		perfMonitor.RecordEnd(operation, duration, err)
	}

	// 获取性能指标
	metrics := perfMonitor.GetMetrics()
	fmt.Printf("性能指标 (%d 个操作):\n", len(metrics))
	for operation, metric := range metrics {
		avgTime := time.Duration(0)
		if metric.Count > 0 {
			avgTime = metric.TotalTime / time.Duration(metric.Count)
		}
		fmt.Printf("  %s:\n", operation)
		fmt.Printf("    调用次数: %d\n", metric.Count)
		fmt.Printf("    总时间: %v\n", metric.TotalTime)
		fmt.Printf("    平均时间: %v\n", avgTime)
		fmt.Printf("    最小时间: %v\n", metric.MinTime)
		fmt.Printf("    最大时间: %v\n", metric.MaxTime)
		fmt.Printf("    错误次数: %d\n", metric.Errors)
		fmt.Printf("    错误率: %.2f%%\n", float64(metric.Errors)/float64(metric.Count)*100)
	}
}

// 主函数
func main() {
	fmt.Println("=== Laravel-Go gRPC 高级功能示例 ===")

	// 示例 1: 高级拦截器
	useAdvancedInterceptors()

	// 示例 2: 流处理功能
	useStreamProcessing()

	// 示例 3: 流处理演示
	demonstrateStreamProcessing()

	// 示例 4: 高级健康检查
	useAdvancedHealthCheck()

	// 示例 5: 性能监控
	demonstratePerformanceMonitoring()

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("这些示例展示了 Laravel-Go gRPC 扩展的高级功能，包括：")
	fmt.Println("1. 丰富的拦截器：缓存、重试、节流、安全、审计、性能监控、熔断器")
	fmt.Println("2. 强大的流处理：转换、过滤、聚合、限流、缓冲、分区、调度、验证、丰富、路由")
	fmt.Println("3. 完整的健康检查：自定义检查、依赖检查、监控告警")
	fmt.Println("4. 详细的性能监控：调用统计、时间分析、错误率监控")
}
