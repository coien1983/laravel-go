package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"laravel-go/framework/microservice"
)

// 示例：使用 gRPC 优化功能
func main() {
	fmt.Println("=== Laravel-Go gRPC 优化功能示例 ===")

	// 示例 1: 创建优化的 gRPC 客户端
	fmt.Println("\n1. 创建优化的 gRPC 客户端")
	optimizedClient := createOptimizedGRPCClient()

	// 示例 2: 使用连接池优化
	fmt.Println("\n2. 使用连接池优化")
	demonstrateConnectionPool(optimizedClient)

	// 示例 3: 使用响应缓存
	fmt.Println("\n3. 使用响应缓存")
	demonstrateResponseCache(optimizedClient)

	// 示例 4: 使用并发控制
	fmt.Println("\n4. 使用并发控制")
	demonstrateConcurrencyControl(optimizedClient)

	// 示例 5: 性能监控
	fmt.Println("\n5. 性能监控")
	demonstratePerformanceMonitoring(optimizedClient)

	// 示例 6: 批量优化
	fmt.Println("\n6. 批量优化")
	demonstrateBatchOptimization(optimizedClient)

	fmt.Println("\n=== gRPC 优化功能示例完成 ===")
}

// createOptimizedGRPCClient 创建优化的 gRPC 客户端
func createOptimizedGRPCClient() *microservice.GRPCOptimizer {
	// 创建服务发现
	discovery := microservice.NewMemoryServiceDiscovery()

	// 注册示例服务
	discovery.Register(&microservice.ServiceInfo{
		Name:    "user-service",
		Address: "localhost",
		Port:    50051,
		Version: "1.0.0",
	})

	// 创建优化的 gRPC 客户端
	optimizer := microservice.NewGRPCOptimizer(
		microservice.WithConnectionPoolSize(20),        // 连接池大小
		microservice.WithResponseCacheTTL(10*time.Minute), // 缓存TTL
		microservice.WithConcurrencyLimit(50),          // 并发限制
	)

	fmt.Printf("创建了优化的 gRPC 客户端:\n")
	fmt.Printf("- 连接池大小: 20\n")
	fmt.Printf("- 响应缓存TTL: 10分钟\n")
	fmt.Printf("- 并发限制: 50\n")

	return optimizer
}

// demonstrateConnectionPool 演示连接池优化
func demonstrateConnectionPool(optimizer *microservice.GRPCOptimizer) {
	fmt.Println("连接池优化演示:")

	// 模拟多个并发请求
	for i := 0; i < 5; i++ {
		go func(id int) {
			ctx := context.Background()
			start := time.Now()

			// 模拟 gRPC 调用
			err := optimizer.OptimizedCallGRPC(ctx, "user-service", "/user.UserService/GetUser", 
				map[string]interface{}{"id": id}, 
				map[string]interface{}{}, 
				nil)

			duration := time.Since(start)
			if err != nil {
				fmt.Printf("  请求 %d: 失败 (%v) - 耗时: %v\n", id, err, duration)
			} else {
				fmt.Printf("  请求 %d: 成功 - 耗时: %v\n", id, duration)
			}
		}(i)
	}

	// 等待所有请求完成
	time.Sleep(2 * time.Second)
	fmt.Println("  连接池优化: 复用连接，减少连接建立开销")
}

// demonstrateResponseCache 演示响应缓存
func demonstrateResponseCache(optimizer *microservice.GRPCOptimizer) {
	fmt.Println("响应缓存演示:")

	request := map[string]interface{}{"id": 1}
	response := map[string]interface{}{}

	// 第一次调用（缓存未命中）
	start := time.Now()
	err := optimizer.OptimizedCallGRPC(context.Background(), "user-service", "/user.UserService/GetUser", 
		request, response, nil)
	firstCallDuration := time.Since(start)

	if err != nil {
		fmt.Printf("  第一次调用: 失败 (%v)\n", err)
		return
	}

	// 第二次调用（缓存命中）
	start = time.Now()
	err = optimizer.OptimizedCallGRPC(context.Background(), "user-service", "/user.UserService/GetUser", 
		request, response, nil)
	secondCallDuration := time.Since(start)

	if err != nil {
		fmt.Printf("  第二次调用: 失败 (%v)\n", err)
		return
	}

	fmt.Printf("  第一次调用: %v (缓存未命中)\n", firstCallDuration)
	fmt.Printf("  第二次调用: %v (缓存命中)\n", secondCallDuration)
	fmt.Printf("  性能提升: %.2f%%\n", float64(firstCallDuration-secondCallDuration)/float64(firstCallDuration)*100)
}

// demonstrateConcurrencyControl 演示并发控制
func demonstrateConcurrencyControl(optimizer *microservice.GRPCOptimizer) {
	fmt.Println("并发控制演示:")

	// 模拟大量并发请求
	successCount := 0
	errorCount := 0
	totalRequests := 100

	for i := 0; i < totalRequests; i++ {
		go func(id int) {
			ctx := context.Background()
			err := optimizer.OptimizedCallGRPC(ctx, "user-service", "/user.UserService/GetUser", 
				map[string]interface{}{"id": id}, 
				map[string]interface{}{}, 
				nil)

			if err != nil {
				errorCount++
			} else {
				successCount++
			}
		}(i)
	}

	// 等待所有请求完成
	time.Sleep(3 * time.Second)

	fmt.Printf("  总请求数: %d\n", totalRequests)
	fmt.Printf("  成功请求: %d\n", successCount)
	fmt.Printf("  失败请求: %d\n", errorCount)
	fmt.Printf("  成功率: %.2f%%\n", float64(successCount)/float64(totalRequests)*100)
	fmt.Println("  并发控制: 防止系统过载，保证服务质量")
}

// demonstratePerformanceMonitoring 演示性能监控
func demonstratePerformanceMonitoring(optimizer *microservice.GRPCOptimizer) {
	fmt.Println("性能监控演示:")

	// 执行一些测试调用
	for i := 0; i < 10; i++ {
		ctx := context.Background()
		err := optimizer.OptimizedCallGRPC(ctx, "user-service", "/user.UserService/GetUser", 
			map[string]interface{}{"id": i}, 
			map[string]interface{}{}, 
			nil)

		if err != nil {
			// 记录错误
			optimizer.performanceOpt.RecordError("user-service", "/user.UserService/GetUser")
		}

		// 模拟不同的响应时间
		time.Sleep(time.Duration(i*10) * time.Millisecond)
	}

	// 获取性能指标
	metrics := optimizer.performanceOpt.GetMetrics()
	
	for key, metric := range metrics {
		fmt.Printf("  服务方法: %s\n", key)
		fmt.Printf("    总调用次数: %d\n", metric.TotalCalls)
		fmt.Printf("    平均响应时间: %v\n", metric.AverageTime)
		fmt.Printf("    最小响应时间: %v\n", metric.MinTime)
		fmt.Printf("    最大响应时间: %v\n", metric.MaxTime)
		fmt.Printf("    错误次数: %d\n", metric.ErrorCount)
		fmt.Printf("    错误率: %.2f%%\n", float64(metric.ErrorCount)/float64(metric.TotalCalls)*100)
		fmt.Printf("    最后调用时间: %v\n", metric.LastCallTime)
	}
}

// demonstrateBatchOptimization 演示批量优化
func demonstrateBatchOptimization(optimizer *microservice.GRPCOptimizer) {
	fmt.Println("批量优化演示:")

	// 模拟批量请求
	batchSize := 10
	start := time.Now()

	// 使用 goroutine 并发处理批量请求
	results := make(chan error, batchSize)
	for i := 0; i < batchSize; i++ {
		go func(id int) {
			ctx := context.Background()
			err := optimizer.OptimizedCallGRPC(ctx, "user-service", "/user.UserService/GetUser", 
				map[string]interface{}{"id": id}, 
				map[string]interface{}{}, 
				nil)
			results <- err
		}(i)
	}

	// 收集结果
	successCount := 0
	errorCount := 0
	for i := 0; i < batchSize; i++ {
		if err := <-results; err != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	totalDuration := time.Since(start)
	averageDuration := totalDuration / time.Duration(batchSize)

	fmt.Printf("  批量大小: %d\n", batchSize)
	fmt.Printf("  总耗时: %v\n", totalDuration)
	fmt.Printf("  平均耗时: %v\n", averageDuration)
	fmt.Printf("  成功请求: %d\n", successCount)
	fmt.Printf("  失败请求: %d\n", errorCount)
	fmt.Printf("  吞吐量: %.2f 请求/秒\n", float64(batchSize)/totalDuration.Seconds())
	fmt.Println("  批量优化: 提高吞吐量，减少网络开销")
}

// 示例：高级优化配置
func advancedOptimizationExample() {
	fmt.Println("\n=== 高级优化配置示例 ===")

	// 创建自定义优化配置
	optimizer := microservice.NewGRPCOptimizer(
		// 连接池配置
		microservice.WithConnectionPoolSize(50),
		
		// 缓存配置
		microservice.WithResponseCacheTTL(30*time.Minute),
		
		// 并发配置
		microservice.WithConcurrencyLimit(200),
	)

	fmt.Println("高级优化配置:")
	fmt.Println("- 大连接池: 支持高并发")
	fmt.Println("- 长缓存TTL: 减少重复请求")
	fmt.Println("- 高并发限制: 提高吞吐量")

	// 使用优化器
	ctx := context.Background()
	err := optimizer.OptimizedCallGRPC(ctx, "user-service", "/user.UserService/GetUser", 
		map[string]interface{}{"id": 1}, 
		map[string]interface{}{}, 
		nil)

	if err != nil {
		log.Printf("优化调用失败: %v", err)
	} else {
		fmt.Println("优化调用成功")
	}
}

// 示例：性能基准测试
func performanceBenchmarkExample() {
	fmt.Println("\n=== 性能基准测试示例 ===")

	// 创建优化器
	optimizer := microservice.NewGRPCOptimizer(
		microservice.WithConnectionPoolSize(10),
		microservice.WithResponseCacheTTL(5*time.Minute),
		microservice.WithConcurrencyLimit(100),
	)

	// 基准测试参数
	iterations := 1000
	concurrent := 10

	fmt.Printf("基准测试配置:\n")
	fmt.Printf("- 总迭代次数: %d\n", iterations)
	fmt.Printf("- 并发数: %d\n", concurrent)
	fmt.Printf("- 连接池大小: 10\n")
	fmt.Printf("- 缓存TTL: 5分钟\n")
	fmt.Printf("- 并发限制: 100\n")

	// 执行基准测试
	start := time.Now()
	
	// 使用工作池执行测试
	semaphore := make(chan struct{}, concurrent)
	results := make(chan error, iterations)

	for i := 0; i < iterations; i++ {
		go func(id int) {
			semaphore <- struct{}{} // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			ctx := context.Background()
			err := optimizer.OptimizedCallGRPC(ctx, "user-service", "/user.UserService/GetUser", 
				map[string]interface{}{"id": id}, 
				map[string]interface{}{}, 
				nil)
			results <- err
		}(i)
	}

	// 收集结果
	successCount := 0
	errorCount := 0
	for i := 0; i < iterations; i++ {
		if err := <-results; err != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	totalDuration := time.Since(start)
	throughput := float64(iterations) / totalDuration.Seconds()

	fmt.Printf("\n基准测试结果:\n")
	fmt.Printf("- 总耗时: %v\n", totalDuration)
	fmt.Printf("- 成功请求: %d\n", successCount)
	fmt.Printf("- 失败请求: %d\n", errorCount)
	fmt.Printf("- 成功率: %.2f%%\n", float64(successCount)/float64(iterations)*100)
	fmt.Printf("- 吞吐量: %.2f 请求/秒\n", throughput)
	fmt.Printf("- 平均响应时间: %v\n", totalDuration/time.Duration(iterations))
} 