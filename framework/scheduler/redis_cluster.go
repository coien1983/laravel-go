package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCluster Redis集群实现
type RedisCluster struct {
	client       *redis.Client
	nodeID       string
	ctx          context.Context
	cancel       context.CancelFunc
	subscribers  map[string]chan ClusterMessage
	subMu        sync.RWMutex
	electionChan chan bool
	stopChan     chan struct{}
}

// RedisClusterConfig Redis集群配置
type RedisClusterConfig struct {
	Addr     string
	Password string
	DB       int
	NodeID   string
}

// NewRedisCluster 创建Redis集群
func NewRedisCluster(config RedisClusterConfig) (*RedisCluster, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithCancel(context.Background())

	rc := &RedisCluster{
		client:      client,
		nodeID:      config.NodeID,
		ctx:         ctx,
		cancel:      cancel,
		subscribers: make(map[string]chan ClusterMessage),
		stopChan:    make(chan struct{}),
	}

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rc, nil
}

// Register 注册节点
func (rc *RedisCluster) Register(nodeID string, info NodeInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("scheduler:nodes:%s", nodeID)
	return rc.client.Set(rc.ctx, key, data, 30*time.Second).Err()
}

// Unregister 注销节点
func (rc *RedisCluster) Unregister(nodeID string) error {
	key := fmt.Sprintf("scheduler:nodes:%s", nodeID)
	return rc.client.Del(rc.ctx, key).Err()
}

// GetNodes 获取所有节点
func (rc *RedisCluster) GetNodes() ([]NodeInfo, error) {
	pattern := "scheduler:nodes:*"
	keys, err := rc.client.Keys(rc.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var nodes []NodeInfo
	for _, key := range keys {
		data, err := rc.client.Get(rc.ctx, key).Result()
		if err != nil {
			continue
		}

		var node NodeInfo
		if err := json.Unmarshal([]byte(data), &node); err != nil {
			continue
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// AcquireLock 获取分布式锁
func (rc *RedisCluster) AcquireLock(key string, ttl time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("scheduler:lock:%s", key)
	value := fmt.Sprintf("%s:%d", rc.nodeID, time.Now().UnixNano())

	// 使用SET NX EX命令实现分布式锁
	result, err := rc.client.SetNX(rc.ctx, lockKey, value, ttl).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

// ReleaseLock 释放分布式锁
func (rc *RedisCluster) ReleaseLock(key string) error {
	lockKey := fmt.Sprintf("scheduler:lock:%s", key)
	return rc.client.Del(rc.ctx, lockKey).Err()
}

// StartElection 启动选举
func (rc *RedisCluster) StartElection(callback func(bool)) error {
	rc.electionChan = make(chan bool, 1)

	go rc.runElection(callback)
	return nil
}

// StopElection 停止选举
func (rc *RedisCluster) StopElection() error {
	close(rc.stopChan)
	return nil
}

// runElection 运行选举
func (rc *RedisCluster) runElection(callback func(bool)) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			isLeader := rc.tryBecomeLeader()
			select {
			case rc.electionChan <- isLeader:
			default:
			}
			callback(isLeader)
		case <-rc.stopChan:
			return
		}
	}
}

// tryBecomeLeader 尝试成为领导者
func (rc *RedisCluster) tryBecomeLeader() bool {
	leaderKey := "scheduler:leader"
	value := fmt.Sprintf("%s:%d", rc.nodeID, time.Now().UnixNano())

	// 尝试设置领导者
	result, err := rc.client.SetNX(rc.ctx, leaderKey, value, 30*time.Second).Result()
	if err != nil {
		return false
	}

	if result {
		// 成功成为领导者，定期续期
		go rc.renewLeadership(leaderKey, value)
		return true
	}

	// 检查当前领导者是否是自己
	currentLeader, err := rc.client.Get(rc.ctx, leaderKey).Result()
	if err != nil {
		return false
	}

	// 解析领导者信息
	var leaderID string
	if _, err := fmt.Sscanf(currentLeader, "%s:", &leaderID); err != nil {
		return false
	}

	return leaderID == rc.nodeID
}

// renewLeadership 续期领导权
func (rc *RedisCluster) renewLeadership(leaderKey, value string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 检查是否仍然是领导者
			currentLeader, err := rc.client.Get(rc.ctx, leaderKey).Result()
			if err != nil || currentLeader != value {
				return
			}

			// 续期
			rc.client.Expire(rc.ctx, leaderKey, 30*time.Second)
		case <-rc.stopChan:
			return
		}
	}
}

// Broadcast 广播消息
func (rc *RedisCluster) Broadcast(msg ClusterMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	channel := "scheduler:messages"
	return rc.client.Publish(rc.ctx, channel, data).Err()
}

// Subscribe 订阅消息
func (rc *RedisCluster) Subscribe(callback func(ClusterMessage)) error {
	channel := "scheduler:messages"
	pubsub := rc.client.Subscribe(rc.ctx, channel)
	defer pubsub.Close()

	for {
		select {
		case msg := <-pubsub.Channel():
			var clusterMsg ClusterMessage
			if err := json.Unmarshal([]byte(msg.Payload), &clusterMsg); err != nil {
				continue
			}

			// 忽略自己发送的消息
			if clusterMsg.NodeID == rc.nodeID {
				continue
			}

			callback(clusterMsg)
		case <-rc.stopChan:
			return nil
		}
	}
}

// Close 关闭集群连接
func (rc *RedisCluster) Close() error {
	rc.cancel()
	close(rc.stopChan)
	return rc.client.Close()
}

// GetLeader 获取当前领导者
func (rc *RedisCluster) GetLeader() (string, error) {
	leaderKey := "scheduler:leader"
	leader, err := rc.client.Get(rc.ctx, leaderKey).Result()
	if err != nil {
		return "", err
	}

	var leaderID string
	if _, err := fmt.Sscanf(leader, "%s:", &leaderID); err != nil {
		return "", err
	}

	return leaderID, nil
}

// IsLeader 检查是否为领导者
func (rc *RedisCluster) IsLeader() bool {
	leaderID, err := rc.GetLeader()
	if err != nil {
		return false
	}
	return leaderID == rc.nodeID
}

// GetClusterInfo 获取集群信息
func (rc *RedisCluster) GetClusterInfo() (map[string]interface{}, error) {
	nodes, err := rc.GetNodes()
	if err != nil {
		return nil, err
	}

	leaderID, _ := rc.GetLeader()

	info := map[string]interface{}{
		"total_nodes": len(nodes),
		"leader_id":   leaderID,
		"node_id":     rc.nodeID,
		"is_leader":   rc.IsLeader(),
		"nodes":       nodes,
	}

	return info, nil
}
