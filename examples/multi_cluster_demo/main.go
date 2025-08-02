package main

import (
	"context"
	"fmt"
	"laravel-go/framework/scheduler"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("=== Laravel-Go 多集群定时器示例 ===")

	// 获取集群类型
	clusterType := getClusterType()
	fmt.Printf("集群类型: %s\n", clusterType)

	// 获取节点ID
	nodeID := getNodeID()
	fmt.Printf("节点ID: %s\n", nodeID)

	// 根据集群类型创建集群
	cluster, err := createCluster(clusterType, nodeID)
	if err != nil {
		log.Printf("创建集群失败: %v", err)
		log.Println("使用内存模式运行...")
		runSingleNodeMode()
		return
	}
	defer cluster.Close()

	// 运行分布式模式
	runDistributedMode(nodeID, cluster)
}

// getClusterType 获取集群类型
func getClusterType() string {
	if clusterType := os.Getenv("CLUSTER_TYPE"); clusterType != "" {
		return clusterType
	}
	return "redis" // 默认使用Redis
}

// getNodeID 获取节点ID
func getNodeID() string {
	if nodeID := os.Getenv("NODE_ID"); nodeID != "" {
		return nodeID
	}
	return fmt.Sprintf("node-%d", time.Now().Unix())
}

// createCluster 根据类型创建集群
func createCluster(clusterType, nodeID string) (scheduler.Cluster, error) {
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
func createRedisCluster(nodeID string) (scheduler.Cluster, error) {
	config := scheduler.RedisClusterConfig{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
		NodeID:   nodeID,
	}

	return scheduler.NewRedisCluster(config)
}

// createEtcdCluster 创建etcd集群
func createEtcdCluster(nodeID string) (scheduler.Cluster, error) {
	endpoints := []string{getEnv("ETCD_ENDPOINTS", "localhost:2379")}

	config := scheduler.EtcdClusterConfig{
		Endpoints: endpoints,
		Username:  getEnv("ETCD_USERNAME", ""),
		Password:  getEnv("ETCD_PASSWORD", ""),
		NodeID:    nodeID,
	}

	return scheduler.NewEtcdCluster(config)
}

// createConsulCluster 创建Consul集群
func createConsulCluster(nodeID string) (scheduler.Cluster, error) {
	config := scheduler.ConsulClusterConfig{
		Address: getEnv("CONSUL_ADDRESS", "localhost:8500"),
		Token:   getEnv("CONSUL_TOKEN", ""),
		NodeID:  nodeID,
	}

	return scheduler.NewConsulCluster(config)
}

// createZookeeperCluster 创建ZooKeeper集群
func createZookeeperCluster(nodeID string) (scheduler.Cluster, error) {
	servers := []string{getEnv("ZOOKEEPER_SERVERS", "localhost:2181")}

	config := scheduler.ZookeeperClusterConfig{
		Servers: servers,
		NodeID:  nodeID,
	}

	return scheduler.NewZookeeperCluster(config)
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// runSingleNodeMode 运行单节点模式
func runSingleNodeMode() {
	fmt.Println("运行单节点模式...")

	// 初始化存储
	store := scheduler.NewMemoryStore()
	scheduler.Init(store)

	// 创建任务
	createTasks()

	// 启动调度器
	if err := scheduler.StartScheduler(); err != nil {
		log.Fatal("启动调度器失败:", err)
	}

	// 等待中断信号
	waitForInterrupt()

	// 停止调度器
	if err := scheduler.StopScheduler(); err != nil {
		log.Fatal("停止调度器失败:", err)
	}
}

// runDistributedMode 运行分布式模式
func runDistributedMode(nodeID string, cluster scheduler.Cluster) {
	fmt.Println("运行分布式模式...")

	// 创建分布式配置
	config := scheduler.DistributedConfig{
		NodeID:                 nodeID,
		Cluster:                cluster,
		ElectionTimeout:        30 * time.Second,
		LockTimeout:            10 * time.Second,
		HeartbeatInterval:      5 * time.Second,
		EnableLeaderElection:   true,
		EnableTaskDistribution: true,
	}

	// 创建存储
	store := scheduler.NewMemoryStore()

	// 创建分布式调度器
	ds := scheduler.NewDistributedScheduler(store, config)

	// 创建任务
	createTasksForDistributed(ds)

	// 启动分布式调度器
	if err := ds.Start(); err != nil {
		log.Fatal("启动分布式调度器失败:", err)
	}

	// 启动监控
	go monitorDistributedScheduler(ds)

	// 演示集群功能
	go demonstrateClusterFeatures(cluster)

	// 等待中断信号
	waitForInterrupt()

	// 停止分布式调度器
	if err := ds.Stop(); err != nil {
		log.Fatal("停止分布式调度器失败:", err)
	}
}

// createTasks 创建任务（单节点模式）
func createTasks() {
	fmt.Println("创建任务...")

	// 每分钟执行的任务
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

	// 每5秒执行的任务
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
}

// createTasksForDistributed 创建任务（分布式模式）
func createTasksForDistributed(ds *scheduler.DistributedScheduler) {
	fmt.Println("创建分布式任务...")

	// 每分钟执行的任务
	everyMinuteHandler := scheduler.NewFuncHandler("distributed-minute", func(ctx context.Context) error {
		fmt.Printf("[%s] 分布式每分钟任务执行 (节点: %s)\n",
			time.Now().Format("15:04:05"), ds.GetDistributedStats().NodeID)
		return nil
	})

	task1 := scheduler.NewTask("distributed-minute", "分布式每分钟任务", "0 * * * * *", everyMinuteHandler)
	task1.SetTimeout(30 * time.Second)
	task1.AddTag("frequency", "minute")
	task1.AddTag("distributed", "true")

	if err := ds.Add(task1); err != nil {
		log.Printf("添加分布式每分钟任务失败: %v", err)
	} else {
		fmt.Println("✓ 添加分布式每分钟任务")
	}

	// 每10秒执行的任务
	every10SecondsHandler := scheduler.NewFuncHandler("distributed-10-seconds", func(ctx context.Context) error {
		fmt.Printf("[%s] 分布式每10秒任务执行 (节点: %s)\n",
			time.Now().Format("15:04:05"), ds.GetDistributedStats().NodeID)
		return nil
	})

	task2 := scheduler.NewTask("distributed-10-seconds", "分布式每10秒任务", "*/10 * * * * *", every10SecondsHandler)
	task2.SetTimeout(10 * time.Second)
	task2.AddTag("frequency", "10-seconds")
	task2.AddTag("distributed", "true")

	if err := ds.Add(task2); err != nil {
		log.Printf("添加分布式每10秒任务失败: %v", err)
	} else {
		fmt.Println("✓ 添加分布式每10秒任务")
	}
}

// monitorDistributedScheduler 监控分布式调度器
func monitorDistributedScheduler(ds *scheduler.DistributedScheduler) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := ds.GetDistributedStats()
			fmt.Printf("\n=== 分布式调度器状态 ===\n")
			fmt.Printf("节点ID: %s\n", stats.NodeID)
			fmt.Printf("是否为领导者: %t\n", stats.IsLeader)
			fmt.Printf("总节点数: %d\n", stats.TotalNodes)
			fmt.Printf("在线节点数: %d\n", stats.OnlineNodes)
			fmt.Printf("领导者ID: %s\n", stats.LeaderID)

			// 获取集群节点信息
			nodes, err := ds.GetClusterNodes()
			if err == nil {
				fmt.Printf("集群节点:\n")
				for _, node := range nodes {
					fmt.Printf("  - %s (%s) - %s\n", node.ID, node.Status, node.LastSeen.Format("15:04:05"))
				}
			}

			// 获取调度器统计
			schedulerStats := ds.GetStats()
			fmt.Printf("调度器统计:\n")
			fmt.Printf("  总任务数: %d\n", schedulerStats.TotalTasks)
			fmt.Printf("  启用任务数: %d\n", schedulerStats.EnabledTasks)
			fmt.Printf("  总执行次数: %d\n", schedulerStats.TotalRuns)
			fmt.Printf("  失败次数: %d\n", schedulerStats.TotalFailed)
			fmt.Printf("  成功率: %.2f%%\n", schedulerStats.SuccessRate)
			fmt.Println()
		}
	}
}

// demonstrateClusterFeatures 演示集群功能
func demonstrateClusterFeatures(cluster scheduler.Cluster) {
	time.Sleep(5 * time.Second) // 等待集群启动

	// 演示分布式锁
	demonstrateDistributedLock(cluster)

	// 演示消息广播
	demonstrateMessageBroadcast(cluster)
}

// demonstrateDistributedLock 演示分布式锁功能
func demonstrateDistributedLock(cluster scheduler.Cluster) {
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
func demonstrateMessageBroadcast(cluster scheduler.Cluster) {
	fmt.Println("\n=== 演示消息广播 ===")

	// 订阅消息
	go func() {
		cluster.Subscribe(func(msg scheduler.ClusterMessage) {
			fmt.Printf("收到消息: %s 来自 %s\n", msg.Type, msg.NodeID)
		})
	}()

	// 广播消息
	msg := scheduler.ClusterMessage{
		Type:      "demo",
		NodeID:    "demo-node",
		Timestamp: time.Now(),
		Data:      []byte("Hello, Multi-Cluster World!"),
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

	fmt.Println("按 Ctrl+C 停止调度器...")
	<-c
	fmt.Println("收到中断信号，正在停止...")
}
