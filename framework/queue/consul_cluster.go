package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

// ConsulCluster Consul集群实现（复用定时器的实现）
type ConsulCluster struct {
	client       *api.Client
	nodeID       string
	ctx          context.Context
	cancel       context.CancelFunc
	stopChan     chan struct{}
	electionChan chan bool
	sessionID    string
}

// ConsulClusterConfig Consul集群配置
type ConsulClusterConfig struct {
	Address string
	NodeID  string
}

// NewConsulCluster 创建Consul集群
func NewConsulCluster(config ConsulClusterConfig) (*ConsulCluster, error) {
	client, err := api.NewClient(&api.Config{
		Address: config.Address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cc := &ConsulCluster{
		client:       client,
		nodeID:       config.NodeID,
		ctx:          ctx,
		cancel:       cancel,
		stopChan:     make(chan struct{}),
		electionChan: make(chan bool, 1),
	}

	// 测试连接
	_, err = client.Agent().Self()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Consul: %w", err)
	}

	return cc, nil
}

// Register 注册节点
func (cc *ConsulCluster) Register(nodeID string, info NodeInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("queue/nodes/%s", nodeID)

	// 创建会话
	session, _, err := cc.client.Session().Create(&api.SessionEntry{
		Name:     fmt.Sprintf("queue-node-%s", nodeID),
		Behavior: "delete",
		TTL:      "30s",
	}, nil)
	if err != nil {
		return err
	}

	// 设置节点信息，使用会话
	_, err = cc.client.KV().Put(&api.KVPair{
		Key:     key,
		Value:   data,
		Session: session,
	}, nil)
	if err != nil {
		return err
	}

	// 保持会话活跃
	go cc.keepSessionAlive(session)

	return nil
}

// Unregister 注销节点
func (cc *ConsulCluster) Unregister(nodeID string) error {
	key := fmt.Sprintf("queue/nodes/%s", nodeID)
	_, err := cc.client.KV().Delete(key, nil)
	return err
}

// GetNodes 获取所有节点
func (cc *ConsulCluster) GetNodes() ([]NodeInfo, error) {
	key := "queue/nodes/"
	resp, _, err := cc.client.KV().List(key, nil)
	if err != nil {
		return nil, err
	}

	var nodes []NodeInfo
	for _, kv := range resp {
		var node NodeInfo
		if err := json.Unmarshal(kv.Value, &node); err != nil {
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// AcquireLock 获取分布式锁
func (cc *ConsulCluster) AcquireLock(key string, ttl time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("queue/locks/%s", key)

	// 创建会话
	session, _, err := cc.client.Session().Create(&api.SessionEntry{
		Name:     fmt.Sprintf("queue-lock-%s", key),
		Behavior: "delete",
		TTL:      fmt.Sprintf("%ds", int(ttl.Seconds())),
	}, nil)
	if err != nil {
		return false, err
	}

	// 尝试获取锁
	acquired, _, err := cc.client.KV().Acquire(&api.KVPair{
		Key:     lockKey,
		Value:   []byte(cc.nodeID),
		Session: session,
	}, nil)
	if err != nil {
		return false, err
	}

	if acquired {
		// 成功获取锁，保持会话活跃
		go cc.keepSessionAlive(session)
		return true, nil
	}

	return false, nil
}

// ReleaseLock 释放分布式锁
func (cc *ConsulCluster) ReleaseLock(key string) error {
	lockKey := fmt.Sprintf("queue/locks/%s", key)

	// 检查锁是否属于当前节点
	resp, _, err := cc.client.KV().Get(lockKey, nil)
	if err != nil {
		return err
	}

	if resp != nil && string(resp.Value) == cc.nodeID {
		_, err = cc.client.KV().Release(&api.KVPair{
			Key:     lockKey,
			Value:   []byte(cc.nodeID),
			Session: resp.Session,
		}, nil)
		return err
	}

	return nil
}

// StartElection 启动选举
func (cc *ConsulCluster) StartElection(callback func(bool)) error {
	go cc.runElection(callback)
	return nil
}

// StopElection 停止选举
func (cc *ConsulCluster) StopElection() error {
	close(cc.stopChan)
	return nil
}

// runElection 运行选举
func (cc *ConsulCluster) runElection(callback func(bool)) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			isLeader := cc.tryBecomeLeader()
			select {
			case cc.electionChan <- isLeader:
			default:
			}
			callback(isLeader)
		case <-cc.stopChan:
			return
		}
	}
}

// tryBecomeLeader 尝试成为领导者
func (cc *ConsulCluster) tryBecomeLeader() bool {
	leaderKey := "queue/leader"

	// 创建会话
	session, _, err := cc.client.Session().Create(&api.SessionEntry{
		Name:     "queue-leader",
		Behavior: "delete",
		TTL:      "30s",
	}, nil)
	if err != nil {
		return false
	}

	// 尝试成为领导者
	acquired, _, err := cc.client.KV().Acquire(&api.KVPair{
		Key:     leaderKey,
		Value:   []byte(cc.nodeID),
		Session: session,
	}, nil)
	if err != nil {
		return false
	}

	if acquired {
		// 成功成为领导者，保持会话活跃
		cc.sessionID = session
		go cc.renewLeadership(session)
		return true
	}

	// 检查当前领导者是否是自己
	resp, _, err := cc.client.KV().Get(leaderKey, nil)
	if err != nil {
		return false
	}

	if resp != nil {
		currentLeader := string(resp.Value)
		return currentLeader == cc.nodeID
	}

	return false
}

// renewLeadership 续期领导权
func (cc *ConsulCluster) renewLeadership(sessionID string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 续期会话
			_, _, err := cc.client.Session().Renew(sessionID, nil)
			if err != nil {
				return
			}
		case <-cc.stopChan:
			return
		}
	}
}

// keepSessionAlive 保持会话活跃
func (cc *ConsulCluster) keepSessionAlive(sessionID string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, _, err := cc.client.Session().Renew(sessionID, nil)
			if err != nil {
				return
			}
		case <-cc.stopChan:
			return
		}
	}
}

// Broadcast 广播消息
func (cc *ConsulCluster) Broadcast(msg ClusterMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("queue/messages/%d", time.Now().UnixNano())
	_, err = cc.client.KV().Put(&api.KVPair{
		Key:   key,
		Value: data,
	}, nil)
	return err
}

// Subscribe 订阅消息
func (cc *ConsulCluster) Subscribe(callback func(ClusterMessage)) error {
	key := "queue/messages/"

	// 监听消息
	plan, err := api.Watch(&api.WatchParams{
		Type:   "keyprefix",
		Prefix: key,
		Handler: func(idx uint64, raw interface{}) {
			if raw == nil {
				return
			}

			v, ok := raw.(api.KVPairs)
			if !ok {
				return
			}

			for _, kv := range v {
				var clusterMsg ClusterMessage
				if err := json.Unmarshal(kv.Value, &clusterMsg); err != nil {
					continue
				}

				// 忽略自己发送的消息
				if clusterMsg.NodeID == cc.nodeID {
					continue
				}

				callback(clusterMsg)
			}
		},
	})
	if err != nil {
		return err
	}

	go func() {
		plan.Run("localhost:8500")
	}()

	return nil
}

// Close 关闭集群连接
func (cc *ConsulCluster) Close() error {
	cc.cancel()
	close(cc.stopChan)
	return nil
}

// GetLeader 获取当前领导者
func (cc *ConsulCluster) GetLeader() (string, error) {
	leaderKey := "queue/leader"
	resp, _, err := cc.client.KV().Get(leaderKey, nil)
	if err != nil {
		return "", err
	}

	if resp != nil {
		return string(resp.Value), nil
	}

	return "", nil
}

// IsLeader 检查是否为领导者
func (cc *ConsulCluster) IsLeader() bool {
	leaderID, err := cc.GetLeader()
	if err != nil {
		return false
	}
	return leaderID == cc.nodeID
}

// GetClusterInfo 获取集群信息
func (cc *ConsulCluster) GetClusterInfo() (map[string]interface{}, error) {
	nodes, err := cc.GetNodes()
	if err != nil {
		return nil, err
	}

	leaderID, _ := cc.GetLeader()

	info := map[string]interface{}{
		"total_nodes": len(nodes),
		"leader_id":   leaderID,
		"node_id":     cc.nodeID,
		"is_leader":   cc.IsLeader(),
		"nodes":       nodes,
	}

	return info, nil
}
