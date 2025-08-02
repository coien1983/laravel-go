package main

import (
	"context"
	"fmt"
	"laravel-go/framework/scheduler"
	"log"
	"time"
)

func main() {
	fmt.Println("=== Laravel-Go 定时器模块示例 ===")

	// 初始化调度器
	store := scheduler.NewMemoryStore()
	scheduler.Init(store)

	// 创建各种示例任务
	createExampleTasks()

	// 启动调度器
	fmt.Println("启动调度器...")
	if err := scheduler.StartScheduler(); err != nil {
		log.Fatal("启动调度器失败:", err)
	}

	// 显示调度器状态
	showSchedulerStatus()

	// 运行一段时间后停止
	fmt.Println("运行 30 秒后停止...")
	time.Sleep(30 * time.Second)

	// 停止调度器
	fmt.Println("停止调度器...")
	if err := scheduler.StopScheduler(); err != nil {
		log.Fatal("停止调度器失败:", err)
	}

	// 显示最终统计
	showFinalStats()
}

// 创建示例任务
func createExampleTasks() {
	fmt.Println("\n--- 创建示例任务 ---")

	// 1. 每分钟执行的任务
	everyMinuteHandler := scheduler.NewFuncHandler("every-minute", func(ctx context.Context) error {
		fmt.Printf("[%s] 每分钟任务执行\n", time.Now().Format("15:04:05"))
		return nil
	})

	task1 := scheduler.NewTask("every-minute", "每分钟执行的任务", "0 * * * * *", everyMinuteHandler)
	task1.SetTimeout(30 * time.Second)
	task1.AddTag("frequency", "minute")
	task1.AddTag("demo", "true")

	if err := scheduler.AddTask(task1); err != nil {
		log.Printf("添加每分钟任务失败: %v", err)
	} else {
		fmt.Println("✓ 添加每分钟任务")
	}

	// 2. 每5秒执行的任务
	every5SecondsHandler := scheduler.NewFuncHandler("every-5-seconds", func(ctx context.Context) error {
		fmt.Printf("[%s] 每5秒任务执行\n", time.Now().Format("15:04:05"))
		return nil
	})

	task2 := scheduler.NewTask("every-5-seconds", "每5秒执行的任务", "*/5 * * * * *", every5SecondsHandler)
	task2.SetTimeout(10 * time.Second)
	task2.AddTag("frequency", "5-seconds")
	task2.AddTag("demo", "true")

	if err := scheduler.AddTask(task2); err != nil {
		log.Printf("添加每5秒任务失败: %v", err)
	} else {
		fmt.Println("✓ 添加每5秒任务")
	}

	// 3. 使用便捷方法创建任务
	hourlyHandler := scheduler.NewFuncHandler("hourly", func(ctx context.Context) error {
		fmt.Printf("[%s] 每小时任务执行\n", time.Now().Format("15:04:05"))
		return nil
	})

	task3 := scheduler.EveryHour(hourlyHandler)
	task3.Name = "hourly-task"
	task3.Description = "每小时执行的任务"
	task3.AddTag("frequency", "hourly")
	task3.AddTag("demo", "true")

	if err := scheduler.AddTask(task3); err != nil {
		log.Printf("添加每小时任务失败: %v", err)
	} else {
		fmt.Println("✓ 添加每小时任务")
	}

	// 4. 使用任务构建器创建任务
	backupHandler := scheduler.NewFuncHandler("backup", func(ctx context.Context) error {
		fmt.Printf("[%s] 备份任务执行\n", time.Now().Format("15:04:05"))
		// 模拟备份过程
		time.Sleep(2 * time.Second)
		return nil
	})

	task4 := scheduler.NewTaskBuilder("backup-task", "模拟备份任务", "0 */2 * * * *", backupHandler).
		SetTimeout(5*time.Minute).
		SetMaxRetries(3).
		SetRetryDelay(30*time.Second).
		AddTag("type", "backup").
		AddTag("priority", "high").
		AddTag("demo", "true").
		Build()

	if err := scheduler.AddTask(task4); err != nil {
		log.Printf("添加备份任务失败: %v", err)
	} else {
		fmt.Println("✓ 添加备份任务")
	}

	// 5. 创建一个会失败的任务（用于演示错误处理）
	errorHandler := scheduler.NewFuncHandler("error-task", func(ctx context.Context) error {
		fmt.Printf("[%s] 错误任务执行（模拟失败）\n", time.Now().Format("15:04:05"))
		return fmt.Errorf("模拟任务执行失败")
	})

	task5 := scheduler.NewTask("error-task", "模拟失败的任务", "0 */3 * * * *", errorHandler)
	task5.SetMaxRetries(2)
	task5.SetRetryDelay(10 * time.Second)
	task5.AddTag("type", "error-demo")
	task5.AddTag("demo", "true")

	if err := scheduler.AddTask(task5); err != nil {
		log.Printf("添加错误任务失败: %v", err)
	} else {
		fmt.Println("✓ 添加错误任务")
	}

	// 6. 使用特殊表达式创建任务
	dailyHandler := scheduler.NewFuncHandler("daily", func(ctx context.Context) error {
		fmt.Printf("[%s] 每日任务执行\n", time.Now().Format("15:04:05"))
		return nil
	})

	task6 := scheduler.EveryDay(dailyHandler)
	task6.Name = "daily-task"
	task6.Description = "每日执行的任务"
	task6.AddTag("frequency", "daily")
	task6.AddTag("demo", "true")

	if err := scheduler.AddTask(task6); err != nil {
		log.Printf("添加每日任务失败: %v", err)
	} else {
		fmt.Println("✓ 添加每日任务")
	}
}

// 显示调度器状态
func showSchedulerStatus() {
	fmt.Println("\n--- 调度器状态 ---")

	status := scheduler.GetSchedulerStatus()
	fmt.Printf("状态: %s\n", status.Status)
	fmt.Printf("启动时间: %s\n", status.StartedAt.Format("15:04:05"))
	fmt.Printf("任务数量: %d\n", status.TaskCount)

	stats := scheduler.GetSchedulerStats()
	fmt.Printf("总任务数: %d\n", stats.TotalTasks)
	fmt.Printf("启用任务数: %d\n", stats.EnabledTasks)
	fmt.Printf("禁用任务数: %d\n", stats.DisabledTasks)
}

// 显示最终统计
func showFinalStats() {
	fmt.Println("\n--- 最终统计 ---")

	// 调度器统计
	stats := scheduler.GetSchedulerStats()
	fmt.Printf("总执行次数: %d\n", stats.TotalRuns)
	fmt.Printf("失败次数: %d\n", stats.TotalFailed)
	fmt.Printf("成功率: %.2f%%\n", stats.SuccessRate)

	// 监控指标
	monitor := scheduler.GetMonitor()
	schedulerMetrics := monitor.GetSchedulerMetrics()
	fmt.Printf("运行时间: %v\n", schedulerMetrics.Uptime)

	// 性能指标
	metrics := scheduler.GetPerformanceMetrics()
	fmt.Printf("吞吐量: %.2f 任务/秒\n", metrics.Throughput)

	// 任务详细统计
	fmt.Println("\n--- 任务详细统计 ---")
	tasks := scheduler.GetAllTasks()
	for _, task := range tasks {
		taskStats, err := scheduler.GetTaskStats(task.GetID())
		if err != nil {
			continue
		}

		fmt.Printf("任务: %s\n", task.GetName())
		fmt.Printf("  运行次数: %d\n", taskStats.RunCount)
		fmt.Printf("  失败次数: %d\n", taskStats.FailedCount)
		fmt.Printf("  成功率: %.2f%%\n", taskStats.SuccessRate)
		if taskStats.LastError != "" {
			fmt.Printf("  最后错误: %s\n", taskStats.LastError)
		}
		fmt.Println()
	}

	// 监控指标
	fmt.Println("--- 监控指标 ---")
	fmt.Printf("调度器运行时间: %v\n", schedulerMetrics.Uptime)
	fmt.Printf("最后活动时间: %s\n", schedulerMetrics.LastActivity.Format("15:04:05"))
}

// 演示任务管理功能
func demonstrateTaskManagement() {
	fmt.Println("\n--- 任务管理演示 ---")

	// 获取所有任务
	tasks := scheduler.GetAllTasks()
	fmt.Printf("总任务数: %d\n", len(tasks))

	// 获取启用的任务
	enabledTasks := scheduler.GetEnabledTasks()
	fmt.Printf("启用任务数: %d\n", len(enabledTasks))

	// 演示暂停和恢复调度器
	fmt.Println("暂停调度器...")
	scheduler.PauseScheduler()
	time.Sleep(5 * time.Second)

	fmt.Println("恢复调度器...")
	scheduler.ResumeScheduler()

	// 演示立即运行任务
	if len(tasks) > 0 {
		fmt.Printf("立即运行任务: %s\n", tasks[0].GetName())
		scheduler.RunTaskNow(tasks[0].GetID())
	}
}

// 演示存储功能
func demonstrateStorage() {
	fmt.Println("\n--- 存储功能演示 ---")

	// 获取存储统计
	store := scheduler.NewMemoryStore()
	storeStats, err := store.GetStats()
	if err != nil {
		fmt.Printf("获取存储统计失败: %v\n", err)
	} else {
		fmt.Printf("存储统计: 总任务数=%d, 启用任务数=%d\n",
			storeStats.TotalTasks, storeStats.EnabledTasks)
	}
}
