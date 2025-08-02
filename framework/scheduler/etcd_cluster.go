package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdCluster etcd集群实现
type EtcdCluster struct {
	client       *clientv3.Client
	nodeID       string
	ctx          context.Context
	cancel       context.CancelFunc
	leaseID      clientv3.LeaseID
	stopChan     chan struct{}
	electionChan chan bool
}

// EtcdClusterConfig etcd集群配置
type EtcdClusterConfig struct {
	Endpoints []string
	Username  string
	Password  string
	NodeID    string
}

// NewEtcdCluster 创建etcd集群
func NewEtcdCluster(config EtcdClusterConfig) (*EtcdCluster, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    config.Username,
		Password:    config.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	ec := &EtcdCluster{
		client:   client,
		nodeID:   config.NodeID,
		ctx:      ctx,
		cancel:   cancel,
		stopChan: make(chan struct{}),
	}

	// 创建租约
	lease, err := client.Grant(ctx, 30)
	if err != nil {
		return nil, fmt.Errorf("failed to create lease: %w", err)
	}
	ec.leaseID = lease.ID

	return ec, nil
}

// Register 注册节点
func (ec *EtcdCluster) Register(nodeID string, info NodeInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("/scheduler/nodes/%s", nodeID)
	_, err = ec.client.Put(ec.ctx, key, string(data), clientv3.WithLease(ec.leaseID))
	return err
}

// Unregister 注销节点
func (ec *EtcdCluster) Unregister(nodeID string) error {
	key := fmt.Sprintf("/scheduler/nodes/%s", nodeID)
	_, err := ec.client.Delete(ec.ctx, key)
	return err
}

// GetNodes 获取所有节点
func (ec *EtcdCluster) GetNodes() ([]NodeInfo, error) {
	resp, err := ec.client.Get(ec.ctx, "/scheduler/nodes/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var nodes []NodeInfo
	for _, kv := range resp.Kvs {
		var node NodeInfo
		if err := json.Unmarshal(kv.Value, &node); err != nil {
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// AcquireLock 获取分布式锁
func (ec *EtcdCluster) AcquireLock(key string, ttl time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("/scheduler/locks/%s", key)
	value := fmt.Sprintf("%s:%d", ec.nodeID, time.Now().UnixNano())

	// 创建租约
	lease, err := ec.client.Grant(ec.ctx, int64(ttl.Seconds()))
	if err != nil {
		return false, err
	}

	// 尝试获取锁
	txn := ec.client.Txn(ec.ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, value, clientv3.WithLease(lease.ID))).
		Else(clientv3.OpGet(lockKey))

	resp, err := txn.Commit()
	if err != nil {
		return false, err
	}

	return resp.Succeeded, nil
}

// ReleaseLock 释放分布式锁
func (ec *EtcdCluster) ReleaseLock(key string) error {
	lockKey := fmt.Sprintf("/scheduler/locks/%s", key)
	_, err := ec.client.Delete(ec.ctx, lockKey)
	return err
}

// StartElection 启动选举
func (ec *EtcdCluster) StartElection(callback func(bool)) error {
	ec.electionChan = make(chan bool, 1)
	go ec.runElection(callback)
	return nil
}

// StopElection 停止选举
func (ec *EtcdCluster) StopElection() error {
	close(ec.stopChan)
	return nil
}

// runElection 运行选举
func (ec *EtcdCluster) runElection(callback func(bool)) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			isLeader := ec.tryBecomeLeader()
			select {
			case ec.electionChan <- isLeader:
			default:
			}
			callback(isLeader)
		case <-ec.stopChan:
			return
		}
	}
}

// tryBecomeLeader 尝试成为领导者
func (ec *EtcdCluster) tryBecomeLeader() bool {
	leaderKey := "/scheduler/leader"
	value := fmt.Sprintf("%s:%d", ec.nodeID, time.Now().UnixNano())

	// 创建租约
	lease, err := ec.client.Grant(ec.ctx, 30)
	if err != nil {
		return false
	}

	// 尝试成为领导者
	txn := ec.client.Txn(ec.ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(leaderKey), "=", 0)).
		Then(clientv3.OpPut(leaderKey, value, clientv3.WithLease(lease.ID))).
		Else(clientv3.OpGet(leaderKey))

	resp, err := txn.Commit()
	if err != nil {
		return false
	}

	if resp.Succeeded {
		// 成功成为领导者，定期续期
		go ec.renewLeadership(leaderKey, value, lease.ID)
		return true
	}

	// 检查当前领导者是否是自己
	if len(resp.Responses) > 0 && len(resp.Responses[0].GetResponseRange().Kvs) > 0 {
		currentLeader := string(resp.Responses[0].GetResponseRange().Kvs[0].Value)
		var leaderID string
		if _, err := fmt.Sscanf(currentLeader, "%s:", &leaderID); err == nil {
			return leaderID == ec.nodeID
		}
	}

	return false
}

// renewLeadership 续期领导权
func (ec *EtcdCluster) renewLeadership(leaderKey, value string, leaseID clientv3.LeaseID) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 续期租约
			_, err := ec.client.KeepAliveOnce(ec.ctx, leaseID)
			if err != nil {
				return
			}
		case <-ec.stopChan:
			return
		}
	}
}

// Broadcast 广播消息
func (ec *EtcdCluster) Broadcast(msg ClusterMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("/scheduler/messages/%d", time.Now().UnixNano())
	_, err = ec.client.Put(ec.ctx, key, string(data), clientv3.WithLease(ec.leaseID))
	return err
}

// Subscribe 订阅消息
func (ec *EtcdCluster) Subscribe(callback func(ClusterMessage)) error {
	watchChan := ec.client.Watch(ec.ctx, "/scheduler/messages/", clientv3.WithPrefix())

	go func() {
		for {
			select {
			case resp := <-watchChan:
				for _, ev := range resp.Events {
					if ev.Type == clientv3.EventTypePut {
						var clusterMsg ClusterMessage
						if err := json.Unmarshal(ev.Kv.Value, &clusterMsg); err != nil {
							continue
						}

						// 忽略自己发送的消息
						if clusterMsg.NodeID == ec.nodeID {
							continue
						}

						callback(clusterMsg)
					}
				}
			case <-ec.stopChan:
				return
			}
		}
	}()

	return nil
}

// Close 关闭集群连接
func (ec *EtcdCluster) Close() error {
	ec.cancel()
	close(ec.stopChan)
	return ec.client.Close()
}

// GetLeader 获取当前领导者
func (ec *EtcdCluster) GetLeader() (string, error) {
	resp, err := ec.client.Get(ec.ctx, "/scheduler/leader")
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("no leader found")
	}

	leader := string(resp.Kvs[0].Value)
	var leaderID string
	if _, err := fmt.Sscanf(leader, "%s:", &leaderID); err != nil {
		return "", err
	}

	return leaderID, nil
}

// IsLeader 检查是否为领导者
func (ec *EtcdCluster) IsLeader() bool {
	leaderID, err := ec.GetLeader()
	if err != nil {
		return false
	}
	return leaderID == ec.nodeID
}

// GetClusterInfo 获取集群信息
func (ec *EtcdCluster) GetClusterInfo() (map[string]interface{}, error) {
	nodes, err := ec.GetNodes()
	if err != nil {
		return nil, err
	}

	leaderID, _ := ec.GetLeader()

	info := map[string]interface{}{
		"total_nodes": len(nodes),
		"leader_id":   leaderID,
		"node_id":     ec.nodeID,
		"is_leader":   ec.IsLeader(),
		"nodes":       nodes,
	}

	return info, nil
}
