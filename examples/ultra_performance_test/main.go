package main

import (
	"context"
	"encoding/json"
	"fmt"
	"laravel-go/framework/performance"
	"log"
	"time"
)

func main() {
	fmt.Println("🧪 Laravel-Go 超高性能优化功能测试")
	fmt.Println("====================================")

	// 创建性能监控器
	monitor := performance.NewPerformanceMonitor()

	// 启动监控
	ctx := context.Background()
	monitor.Start(ctx)
	defer monitor.Stop()

	// 创建超高性能优化器
	ultraOptimizer := performance.NewUltraOptimizer(monitor)

	// 执行超高性能优化
	fmt.Println("\n🔧 执行超高性能优化...")
	start := time.Now()
	results, err := ultraOptimizer.Optimize(ctx)
	duration := time.Since(start)

	if err != nil {
		log.Printf("超高性能优化失败: %v", err)
		return
	}

	fmt.Printf("✅ 超高性能优化完成，耗时: %v\n", duration)
	fmt.Printf("📊 优化结果数量: %d\n", len(results))

	// 打印优化结果
	for i, result := range results {
		fmt.Printf("\n--- 优化结果 %d ---\n", i+1)
		fmt.Printf("类型: %s\n", result.Type)
		fmt.Printf("成功: %t\n", result.Success)
		fmt.Printf("消息: %s\n", result.Message)
		fmt.Printf("改进: %.1f%%\n", result.Improvement)
		fmt.Printf("耗时: %v\n", result.Duration)

		if len(result.Metrics) > 0 {
			fmt.Printf("指标:\n")
			for key, value := range result.Metrics {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
	}

	// 创建智能缓存优化器
	fmt.Println("\n💾 创建智能缓存优化器...")
	cacheOptimizer := performance.NewSmartCacheOptimizer(monitor)

	// 执行智能缓存优化
	fmt.Println("🔧 执行智能缓存优化...")
	start = time.Now()
	cacheResults, err := cacheOptimizer.Optimize(ctx)
	duration = time.Since(start)

	if err != nil {
		log.Printf("智能缓存优化失败: %v", err)
	} else {
		fmt.Printf("✅ 智能缓存优化完成，耗时: %v\n", duration)
		fmt.Printf("📊 优化结果数量: %d\n", len(cacheResults))
	}

	// 创建数据库优化器
	fmt.Println("\n🗄️ 创建数据库优化器...")
	dbOptimizer := performance.NewDatabaseOptimizer(monitor)

	// 执行数据库优化
	fmt.Println("🔧 执行数据库优化...")
	start = time.Now()
	dbResults, err := dbOptimizer.Optimize(ctx)
	duration = time.Since(start)

	if err != nil {
		log.Printf("数据库优化失败: %v", err)
	} else {
		fmt.Printf("✅ 数据库优化完成，耗时: %v\n", duration)
		fmt.Printf("📊 优化结果数量: %d\n", len(dbResults))
	}

	// 生成综合报告
	fmt.Println("\n📈 生成综合性能报告...")

	allResults := map[string]interface{}{
		"ultra_optimization":    results,
		"cache_optimization":    cacheResults,
		"database_optimization": dbResults,
		"timestamp":             time.Now(),
		"total_optimizations":   len(results) + len(cacheResults) + len(dbResults),
	}

	// 计算平均性能提升
	var totalImprovement float64
	var totalCount int

	for _, result := range results {
		totalImprovement += result.Improvement
		totalCount++
	}
	for _, result := range cacheResults {
		totalImprovement += result.Improvement
		totalCount++
	}
	for _, result := range dbResults {
		totalImprovement += result.Improvement
		totalCount++
	}

	averageImprovement := 0.0
	if totalCount > 0 {
		averageImprovement = totalImprovement / float64(totalCount)
	}

	allResults["average_improvement"] = averageImprovement

	// 输出JSON格式的报告
	reportJSON, _ := json.MarshalIndent(allResults, "", "  ")
	fmt.Printf("\n📊 综合性能报告:\n%s\n", string(reportJSON))

	fmt.Println("\n🎉 超高性能优化功能测试完成！")
	fmt.Printf("📈 平均性能提升: %.1f%%\n", averageImprovement)
	fmt.Printf("🔧 总优化项目: %d\n", totalCount)
}
