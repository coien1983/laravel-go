package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"laravel-go/framework/queue"
)

func main() {
	fmt.Println("=== Laravel-Go 分布式队列示例 ===")

	// 获取节点ID
	nodeID := getNodeID()
	fmt.Printf("节点ID: %s\n", nodeID)

	// 创建Redis集群（需要先启动Redis服务）
	cluster, err := createRedisCluster(nodeID)
	if err != nil {
		log.Printf("创建Redis集群失败: %v", err)
		log.Println("使用内存模式运行...")
		runSingleNodeMode()
		return
	}
	defer cluster.Close()

	// 运行分布式模式
	runDistributedMode(nodeID, cluster)
}

// getNodeID 获取节点ID
func getNodeID() string {
	if nodeID := os.Getenv("NODE_ID"); nodeID != "" {
		return nodeID
	}
	return fmt.Sprintf("node-%d", time.Now().Unix())
}

// createRedisCluster 创建Redis集群
func createRedisCluster(nodeID string) (queue.Cluster, error) {
	// 注意：这里需要先安装Redis依赖
	// go get github.com/go-redis/redis/v8
	
	config := queue.RedisClusterConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		NodeID:   nodeID,
	}

	return queue.NewRedisCluster(config)
}

// runSingleNodeMode 运行单节点模式
func runSingleNodeMode() {
	fmt.Println("运行单节点模式...")

	// 初始化队列管理器
	queue.Init()

	// 注册内存队列
	memoryQueue := queue.NewMemoryQueue()
	queue.QueueManager.Extend("memory", memoryQueue)
	queue.QueueManager.SetDefaultQueue("memory")

	// 创建任务
	createJobs()

	// 启动工作进程
	worker := queue.NewWorker(memoryQueue, "default")
	worker.SetOnCompleted(func(job queue.Job) {
		fmt.Printf("任务完成: %s\n", string(job.GetPayload()))
	})
	worker.SetOnFailed(func(job queue.Job, err error) {
		fmt.Printf("任务失败: %s - %v\n", string(job.GetPayload()), err)
	})

	if err := worker.Start(); err != nil {
		log.Fatal("启动工作进程失败:", err)
	}

	// 等待中断信号
	waitForInterrupt()

	// 停止工作进程
	if err := worker.Stop(); err != nil {
		log.Fatal("停止工作进程失败:", err)
	}
}

// runDistributedMode 运行分布式模式
func runDistributedMode(nodeID string, cluster queue.Cluster) {
	fmt.Println("运行分布式模式...")

	// 创建分布式配置
	config := queue.DistributedConfig{
		NodeID:                nodeID,
		Cluster:               cluster,
		ElectionTimeout:       30 * time.Second,
		LockTimeout:           10 * time.Second,
		HeartbeatInterval:     5 * time.Second,
		EnableLeaderElection:  true,
		EnableJobDistribution: true,
		WorkerCount:           3,
		MaxConcurrency:        5,
	}

	// 创建分布式队列
	dq := queue.NewDistributedQueue(config)

	// 创建任务
	createJobsForDistributed(dq)

	// 设置回调
	dq.workerPool.SetOnCompleted(func(job queue.Job) {
		fmt.Printf("[%s] 分布式任务完成: %s (节点: %s)\n", 
			time.Now().Format("15:04:05"), string(job.GetPayload()), dq.GetDistributedStats().NodeID)
	})

	dq.workerPool.SetOnFailed(func(job queue.Job, err error) {
		fmt.Printf("[%s] 分布式任务失败: %s - %v (节点: %s)\n", 
			time.Now().Format("15:04:05"), string(job.GetPayload()), err, dq.GetDistributedStats().NodeID)
	})

	// 启动分布式队列
	if err := dq.Start(); err != nil {
		log.Fatal("启动分布式队列失败:", err)
	}

	// 启动监控
	go monitorDistributedQueue(dq)

	// 演示集群功能
	go demonstrateClusterFeatures(cluster)

	// 等待中断信号
	waitForInterrupt()

	// 停止分布式队列
	if err := dq.Stop(); err != nil {
		log.Fatal("停止分布式队列失败:", err)
	}
}

// createJobs 创建任务（单节点模式）
func createJobs() {
	fmt.Println("创建任务...")

	// 创建普通任务
	job1 := queue.NewJob([]byte("Hello Queue!"), "default")
	job1.SetPriority(1)
	job1.AddTag("type", "greeting")
	job1.AddTag("demo", "true")

	if err := queue.Push(job1); err != nil {
		log.Printf("推送任务失败: %v", err)
	} else {
		fmt.Println("✓ 推送普通任务")
	}

	// 创建延迟任务
	job2 := queue.NewJob([]byte("Delayed Job!"), "default")
	job2.SetDelay(5 * time.Second)
	job2.SetPriority(2)
	job2.AddTag("type", "delayed")
	job2.AddTag("demo", "true")

	if err := queue.Push(job2); err != nil {
		log.Printf("推送延迟任务失败: %v", err)
	} else {
		fmt.Println("✓ 推送延迟任务")
	}

	// 创建批量任务
	var jobs []queue.Job
	for i := 1; i <= 5; i++ {
		job := queue.NewJob([]byte(fmt.Sprintf("Batch Job %d", i)), "default")
		job.SetPriority(i)
		job.AddTag("type", "batch")
		job.AddTag("demo", "true")
		jobs = append(jobs, job)
	}

	if err := queue.PushBatch(jobs); err != nil {
		log.Printf("推送批量任务失败: %v", err)
	} else {
		fmt.Println("✓ 推送批量任务")
	}
}

// createJobsForDistributed 创建任务（分布式模式）
func createJobsForDistributed(dq *queue.DistributedQueue) {
	fmt.Println("创建分布式任务...")

	// 创建普通任务
	job1 := queue.NewJob([]byte("Distributed Hello!"), "default")
	job1.SetPriority(1)
	job1.AddTag("type", "distributed")
	job1.AddTag("demo", "true")

	if err := dq.Push(job1); err != nil {
		log.Printf("推送分布式任务失败: %v", err)
	} else {
		fmt.Println("✓ 推送分布式任务")
	}

	// 创建延迟任务
	job2 := queue.NewJob([]byte("Distributed Delayed Job!"), "default")
	job2.SetDelay(3 * time.Second)
	job2.SetPriority(2)
	job2.AddTag("type", "distributed-delayed")
	job2.AddTag("demo", "true")

	if err := dq.Push(job2); err != nil {
		log.Printf("推送分布式延迟任务失败: %v", err)
	} else {
		fmt.Println("✓ 推送分布式延迟任务")
	}

	// 创建批量任务
	var jobs []queue.Job
	for i := 1; i <= 10; i++ {
		job := queue.NewJob([]byte(fmt.Sprintf("Distributed Batch Job %d", i)), "default")
		job.SetPriority(i % 3)
		job.AddTag("type", "distributed-batch")
		job.AddTag("demo", "true")
		jobs = append(jobs, job)
	}

	if err := dq.PushBatch(jobs); err != nil {
		log.Printf("推送分布式批量任务失败: %v", err)
	} else {
		fmt.Println("✓ 推送分布式批量任务")
	}

	// 定期推送新任务
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				job := queue.NewJob([]byte(fmt.Sprintf("Periodic Job at %s", time.Now().Format("15:04:05"))), "default")
				job.AddTag("type", "periodic")
				job.AddTag("demo", "true")

				if err := dq.Push(job); err != nil {
					log.Printf("推送定期任务失败: %v", err)
				} else {
					fmt.Printf("✓ 推送定期任务\n")
				}
			}
		}
	}()
}

// monitorDistributedQueue 监控分布式队列
func monitorDistributedQueue(dq *queue.DistributedQueue) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := dq.GetDistributedStats()
			fmt.Printf("\n=== 分布式队列状态 ===\n")
			fmt.Printf("节点ID: %s\n", stats.NodeID)
			fmt.Printf("是否为领导者: %t\n", stats.IsLeader)
			fmt.Printf("总节点数: %d\n", stats.TotalNodes)
			fmt.Printf("在线节点数: %d\n", stats.OnlineNodes)
			fmt.Printf("领导者ID: %s\n", stats.LeaderID)

			// 获取集群节点信息
			nodes, err := dq.GetClusterNodes()
			if err == nil {
				fmt.Printf("集群节点:\n")
				for _, node := range nodes {
					fmt.Printf("  - %s (%s) - %s\n", node.ID, node.Status, node.LastSeen.Format("15:04:05"))
				}
			}

			// 获取队列统计
			fmt.Printf("队列统计:\n")
			fmt.Printf("  总任务数: %d\n", stats.QueueStats.TotalJobs)
			fmt.Printf("  待处理任务数: %d\n", stats.QueueStats.PendingJobs)
			fmt.Printf("  保留任务数: %d\n", stats.QueueStats.ReservedJobs)
			fmt.Printf("  失败任务数: %d\n", stats.QueueStats.FailedJobs)
			fmt.Printf("  完成任务数: %d\n", stats.QueueStats.CompletedJobs)

			// 获取工作进程池统计
			poolStats := dq.workerPool.GetStats()
			fmt.Printf("工作进程池统计:\n")
			fmt.Printf("  总工作进程数: %d\n", poolStats.TotalWorkers)
			fmt.Printf("  活跃工作进程数: %d\n", poolStats.ActiveWorkers)
			fmt.Printf("  空闲工作进程数: %d\n", poolStats.IdleWorkers)
			fmt.Printf("  总处理任务数: %d\n", poolStats.TotalProcessed)
			fmt.Printf("  总失败任务数: %d\n", poolStats.TotalFailed)
			fmt.Printf("  状态: %s\n", poolStats.Status)
			fmt.Println()
		}
	}
}

// demonstrateClusterFeatures 演示集群功能
func demonstrateClusterFeatures(cluster queue.Cluster) {
	time.Sleep(5 * time.Second) // 等待集群启动

	// 演示分布式锁
	demonstrateDistributedLock(cluster)

	// 演示消息广播
	demonstrateMessageBroadcast(cluster)
}

// demonstrateDistributedLock 演示分布式锁功能
func demonstrateDistributedLock(cluster queue.Cluster) {
	fmt.Println("\n=== 演示分布式锁 ===")

	// 尝试获取锁
	acquired, err := cluster.AcquireLock("demo-lock", 10*time.Second)
	if err != nil {
		log.Printf("获取锁失败: %v", err)
		return
	}

	if acquired {
		fmt.Println("✓ 成功获取分布式锁")
		time.Sleep(5 * time.Second)

		// 释放锁
		if err := cluster.ReleaseLock("demo-lock"); err != nil {
			log.Printf("释放锁失败: %v", err)
		} else {
			fmt.Println("✓ 成功释放分布式锁")
		}
	} else {
		fmt.Println("✗ 无法获取分布式锁（可能被其他节点持有）")
	}
}

// demonstrateMessageBroadcast 演示集群消息广播
func demonstrateMessageBroadcast(cluster queue.Cluster) {
	fmt.Println("\n=== 演示消息广播 ===")

	// 订阅消息
	go func() {
		cluster.Subscribe(func(msg queue.ClusterMessage) {
			fmt.Printf("收到消息: %s 来自 %s\n", msg.Type, msg.NodeID)
		})
	}()

	// 广播消息
	msg := queue.ClusterMessage{
		Type:      "demo",
		NodeID:    "demo-node",
		Timestamp: time.Now(),
		Data:      []byte("Hello, Distributed Queue World!"),
	}

	if err := cluster.Broadcast(msg); err != nil {
		log.Printf("广播消息失败: %v", err)
	} else {
		fmt.Println("✓ 成功广播消息")
	}

	time.Sleep(2 * time.Second)
}

// waitForInterrupt 等待中断信号
func waitForInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fmt.Println("按 Ctrl+C 停止队列...")
	<-c
	fmt.Println("收到中断信号，正在停止...")
} 