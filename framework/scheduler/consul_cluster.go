package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

// ConsulCluster Consul集群实现
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
	Token   string
	NodeID  string
}

// NewConsulCluster 创建Consul集群
func NewConsulCluster(config ConsulClusterConfig) (*ConsulCluster, error) {
	client, err := api.NewClient(&api.Config{
		Address: config.Address,
		Token:   config.Token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cc := &ConsulCluster{
		client:   client,
		nodeID:   config.NodeID,
		ctx:      ctx,
		cancel:   cancel,
		stopChan: make(chan struct{}),
	}

	// 创建会话
	session, _, err := client.Session().Create(&api.SessionEntry{
		Name:     fmt.Sprintf("scheduler-%s", config.NodeID),
		Behavior: "delete",
		TTL:      "30s",
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	cc.sessionID = session

	return cc, nil
}

// Register 注册节点
func (cc *ConsulCluster) Register(nodeID string, info NodeInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("scheduler/nodes/%s", nodeID)
	_, err = cc.client.KV().Put(&api.KVPair{
		Key:     key,
		Value:   data,
		Session: cc.sessionID,
	}, nil)
	return err
}

// Unregister 注销节点
func (cc *ConsulCluster) Unregister(nodeID string) error {
	key := fmt.Sprintf("scheduler/nodes/%s", nodeID)
	_, err := cc.client.KV().Delete(key, nil)
	return err
}

// GetNodes 获取所有节点
func (cc *ConsulCluster) GetNodes() ([]NodeInfo, error) {
	pairs, _, err := cc.client.KV().List("scheduler/nodes/", nil)
	if err != nil {
		return nil, err
	}

	var nodes []NodeInfo
	for _, pair := range pairs {
		var node NodeInfo
		if err := json.Unmarshal(pair.Value, &node); err != nil {
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// AcquireLock 获取分布式锁
func (cc *ConsulCluster) AcquireLock(key string, ttl time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("scheduler/locks/%s", key)
	value := fmt.Sprintf("%s:%d", cc.nodeID, time.Now().UnixNano())

	acquired, _, err := cc.client.KV().Acquire(&api.KVPair{
		Key:     lockKey,
		Value:   []byte(value),
		Session: cc.sessionID,
	}, nil)

	return acquired, err
}

// ReleaseLock 释放分布式锁
func (cc *ConsulCluster) ReleaseLock(key string) error {
	lockKey := fmt.Sprintf("scheduler/locks/%s", key)
	_, err := cc.client.KV().Release(&api.KVPair{
		Key:     lockKey,
		Session: cc.sessionID,
	}, nil)
	return err
}

// StartElection 启动选举
func (cc *ConsulCluster) StartElection(callback func(bool)) error {
	cc.electionChan = make(chan bool, 1)
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
	leaderKey := "scheduler/leader"
	value := fmt.Sprintf("%s:%d", cc.nodeID, time.Now().UnixNano())

	acquired, _, err := cc.client.KV().Acquire(&api.KVPair{
		Key:     leaderKey,
		Value:   []byte(value),
		Session: cc.sessionID,
	}, nil)

	if err != nil {
		return false
	}

	if acquired {
		// 成功成为领导者，定期续期会话
		go cc.renewSession()
		return true
	}

	// 检查当前领导者是否是自己
	pair, _, err := cc.client.KV().Get(leaderKey, nil)
	if err != nil || pair == nil {
		return false
	}

	currentLeader := string(pair.Value)
	var leaderID string
	if _, err := fmt.Sscanf(currentLeader, "%s:", &leaderID); err == nil {
		return leaderID == cc.nodeID
	}

	return false
}

// renewSession 续期会话
func (cc *ConsulCluster) renewSession() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 续期会话
			_, _, err := cc.client.Session().Renew(cc.sessionID, nil)
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

	key := fmt.Sprintf("scheduler/messages/%d", time.Now().UnixNano())
	_, err = cc.client.KV().Put(&api.KVPair{
		Key:   key,
		Value: data,
	}, nil)
	return err
}

// Subscribe 订阅消息
func (cc *ConsulCluster) Subscribe(callback func(ClusterMessage)) error {
	// 使用Consul的Watch机制
	params := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": "scheduler/messages/",
		"handler": func(idx uint64, raw interface{}) {
			if raw == nil {
				return
			}

			pairs, ok := raw.(api.KVPairs)
			if !ok {
				return
			}

			for _, pair := range pairs {
				var clusterMsg ClusterMessage
				if err := json.Unmarshal(pair.Value, &clusterMsg); err != nil {
					continue
				}

				// 忽略自己发送的消息
				if clusterMsg.NodeID == cc.nodeID {
					continue
				}

				callback(clusterMsg)
			}
		},
	}

	plan, err := api.Watch(params, nil)
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

	// 销毁会话
	if cc.sessionID != "" {
		cc.client.Session().Destroy(cc.sessionID, nil)
	}

	return nil
}

// GetLeader 获取当前领导者
func (cc *ConsulCluster) GetLeader() (string, error) {
	pair, _, err := cc.client.KV().Get("scheduler/leader", nil)
	if err != nil {
		return "", err
	}

	if pair == nil {
		return "", fmt.Errorf("no leader found")
	}

	leader := string(pair.Value)
	var leaderID string
	if _, err := fmt.Sscanf(leader, "%s:", &leaderID); err != nil {
		return "", err
	}

	return leaderID, nil
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
