package main

import (
	"fmt"
	"laravel-go/framework/queue"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("=== Laravel-Go 多集群队列示例 ===")

	// 获取集群类型和节点ID
	clusterType := getEnv("CLUSTER_TYPE", "redis")
	nodeID := getEnv("NODE_ID", fmt.Sprintf("node-%d", time.Now().Unix()))

	fmt.Printf("集群类型: %s\n", clusterType)
	fmt.Printf("节点ID: %s\n", nodeID)

	// 创建集群
	cluster, err := createCluster(clusterType, nodeID)
	if err != nil {
		log.Printf("创建集群失败: %v", err)
		log.Println("使用内存模式运行...")
		runSingleNodeMode()
		return
	}
	defer cluster.Close()

	// 运行分布式模式
	runDistributedMode(nodeID, cluster, clusterType)
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// createCluster 根据类型创建集群
func createCluster(clusterType, nodeID string) (queue.Cluster, error) {
	switch clusterType {
	case "redis":
		return createRedisCluster(nodeID)
	case "etcd":
		return createEtcdCluster(nodeID)
	case "consul":
		return createConsulCluster(nodeID)
	case "zookeeper":
		return createZookeeperCluster(nodeID)
	default:
		return nil, fmt.Errorf("unsupported cluster type: %s", clusterType)
	}
}

// createRedisCluster 创建Redis集群
func createRedisCluster(nodeID string) (queue.Cluster, error) {
	config := queue.RedisClusterConfig{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
		NodeID:   nodeID,
	}

	return queue.NewRedisCluster(config)
}

// createEtcdCluster 创建etcd集群
func createEtcdCluster(nodeID string) (queue.Cluster, error) {
	endpoints := []string{getEnv("ETCD_ENDPOINTS", "localhost:2379")}

	config := queue.EtcdClusterConfig{
		Endpoints: endpoints,
		NodeID:    nodeID,
	}

	return queue.NewEtcdCluster(config)
}

// createConsulCluster 创建Consul集群
func createConsulCluster(nodeID string) (queue.Cluster, error) {
	config := queue.ConsulClusterConfig{
		Address: getEnv("CONSUL_ADDRESS", "localhost:8500"),
		NodeID:  nodeID,
	}

	return queue.NewConsulCluster(config)
}

// createZookeeperCluster 创建ZooKeeper集群
func createZookeeperCluster(nodeID string) (queue.Cluster, error) {
	servers := []string{getEnv("ZOOKEEPER_SERVERS", "localhost:2181")}

	config := queue.ZookeeperClusterConfig{
		Servers: servers,
		NodeID:  nodeID,
	}

	return queue.NewZookeeperCluster(config)
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
func runDistributedMode(nodeID string, cluster queue.Cluster, clusterType string) {
	fmt.Printf("运行分布式模式 (集群: %s)...\n", clusterType)

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
	createJobsForDistributed(dq, clusterType)

	// 设置回调
	dq.SetOnCompleted(func(job queue.Job) {
		fmt.Printf("[%s] 分布式任务完成: %s (节点: %s, 集群: %s)\n",
			time.Now().Format("15:04:05"), string(job.GetPayload()), dq.GetDistributedStats().NodeID, clusterType)
	})

	dq.SetOnFailed(func(job queue.Job, err error) {
		fmt.Printf("[%s] 分布式任务失败: %s - %v (节点: %s, 集群: %s)\n",
			time.Now().Format("15:04:05"), string(job.GetPayload()), err, dq.GetDistributedStats().NodeID, clusterType)
	})

	// 启动分布式队列
	if err := dq.Start(); err != nil {
		log.Fatal("启动分布式队列失败:", err)
	}

	// 启动监控
	go monitorDistributedQueue(dq, clusterType)

	// 演示集群功能
	go demonstrateClusterFeatures(cluster, clusterType)

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
}

// createJobsForDistributed 创建任务（分布式模式）
func createJobsForDistributed(dq *queue.DistributedQueue, clusterType string) {
	fmt.Printf("创建分布式任务 (集群: %s)...\n", clusterType)

	// 创建普通任务
	job1 := queue.NewJob([]byte(fmt.Sprintf("Distributed Hello from %s!", clusterType)), "default")
	job1.SetPriority(1)
	job1.AddTag("type", "distributed")
	job1.AddTag("cluster", clusterType)
	job1.AddTag("demo", "true")

	if err := dq.Push(job1); err != nil {
		log.Printf("推送分布式任务失败: %v", err)
	} else {
		fmt.Println("✓ 推送分布式任务")
	}

	// 创建延迟任务
	job2 := queue.NewJob([]byte(fmt.Sprintf("Distributed Delayed Job from %s!", clusterType)), "default")
	job2.SetDelay(3 * time.Second)
	job2.SetPriority(2)
	job2.AddTag("type", "distributed-delayed")
	job2.AddTag("cluster", clusterType)
	job2.AddTag("demo", "true")

	if err := dq.Push(job2); err != nil {
		log.Printf("推送分布式延迟任务失败: %v", err)
	} else {
		fmt.Println("✓ 推送分布式延迟任务")
	}

	// 创建批量任务
	var jobs []queue.Job
	for i := 1; i <= 5; i++ {
		job := queue.NewJob([]byte(fmt.Sprintf("Distributed Batch Job %d from %s", i, clusterType)), "default")
		job.SetPriority(i % 3)
		job.AddTag("type", "distributed-batch")
		job.AddTag("cluster", clusterType)
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
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				job := queue.NewJob([]byte(fmt.Sprintf("Periodic Job from %s at %s", clusterType, time.Now().Format("15:04:05"))), "default")
				job.AddTag("type", "periodic")
				job.AddTag("cluster", clusterType)
				job.AddTag("demo", "true")

				if err := dq.Push(job); err != nil {
					log.Printf("推送定期任务失败: %v", err)
				} else {
					fmt.Printf("✓ 推送定期任务 (集群: %s)\n", clusterType)
				}
			}
		}
	}()
}

// monitorDistributedQueue 监控分布式队列
func monitorDistributedQueue(dq *queue.DistributedQueue, clusterType string) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := dq.GetDistributedStats()
			fmt.Printf("\n=== 分布式队列状态 (集群: %s) ===\n", clusterType)
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
			poolStats := dq.GetWorkerPool().GetStats()
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
func demonstrateClusterFeatures(cluster queue.Cluster, clusterType string) {
	time.Sleep(5 * time.Second) // 等待集群启动

	// 演示分布式锁
	demonstrateDistributedLock(cluster, clusterType)

	// 演示消息广播
	demonstrateMessageBroadcast(cluster, clusterType)
}

// demonstrateDistributedLock 演示分布式锁功能
func demonstrateDistributedLock(cluster queue.Cluster, clusterType string) {
	fmt.Printf("\n=== 演示分布式锁 (集群: %s) ===\n", clusterType)

	// 尝试获取锁
	acquired, err := cluster.AcquireLock("demo-lock", 10*time.Second)
	if err != nil {
		log.Printf("获取锁失败: %v", err)
		return
	}

	if acquired {
		fmt.Printf("✓ 成功获取分布式锁 (集群: %s)\n", clusterType)
		time.Sleep(5 * time.Second)

		// 释放锁
		if err := cluster.ReleaseLock("demo-lock"); err != nil {
			log.Printf("释放锁失败: %v", err)
		} else {
			fmt.Printf("✓ 成功释放分布式锁 (集群: %s)\n", clusterType)
		}
	} else {
		fmt.Printf("✗ 无法获取分布式锁 (集群: %s, 可能被其他节点持有)\n", clusterType)
	}
}

// demonstrateMessageBroadcast 演示集群消息广播
func demonstrateMessageBroadcast(cluster queue.Cluster, clusterType string) {
	fmt.Printf("\n=== 演示消息广播 (集群: %s) ===\n", clusterType)

	// 订阅消息
	go func() {
		cluster.Subscribe(func(msg queue.ClusterMessage) {
			fmt.Printf("收到消息: %s 来自 %s (集群: %s)\n", msg.Type, msg.NodeID, clusterType)
		})
	}()

	// 广播消息
	msg := queue.ClusterMessage{
		Type:      "demo",
		NodeID:    "demo-node",
		Timestamp: time.Now(),
		Data:      []byte(fmt.Sprintf("Hello, Distributed Queue World from %s!", clusterType)),
	}

	if err := cluster.Broadcast(msg); err != nil {
		log.Printf("广播消息失败: %v", err)
	} else {
		fmt.Printf("✓ 成功广播消息 (集群: %s)\n", clusterType)
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
