package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
)

// ZookeeperCluster ZooKeeper集群实现
type ZookeeperCluster struct {
	conn         *zk.Conn
	nodeID       string
	ctx          context.Context
	cancel       context.CancelFunc
	stopChan     chan struct{}
	electionChan chan bool
	leaderPath   string
}

// ZookeeperClusterConfig ZooKeeper集群配置
type ZookeeperClusterConfig struct {
	Servers []string
	NodeID  string
}

// NewZookeeperCluster 创建ZooKeeper集群
func NewZookeeperCluster(config ZookeeperClusterConfig) (*ZookeeperCluster, error) {
	conn, _, err := zk.Connect(config.Servers, time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to zookeeper: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	zc := &ZookeeperCluster{
		conn:       conn,
		nodeID:     config.NodeID,
		ctx:        ctx,
		cancel:     cancel,
		stopChan:   make(chan struct{}),
		leaderPath: "/scheduler/leader",
	}

	// 创建必要的路径
	paths := []string{
		"/scheduler",
		"/scheduler/nodes",
		"/scheduler/locks",
		"/scheduler/messages",
	}

	for _, path := range paths {
		exists, _, err := conn.Exists(path)
		if err != nil {
			return nil, fmt.Errorf("failed to check path %s: %w", path, err)
		}
		if !exists {
			_, err = conn.Create(path, []byte{}, 0, zk.WorldACL(zk.PermAll))
			if err != nil && err != zk.ErrNodeExists {
				return nil, fmt.Errorf("failed to create path %s: %w", path, err)
			}
		}
	}

	return zc, nil
}

// Register 注册节点
func (zc *ZookeeperCluster) Register(nodeID string, info NodeInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/scheduler/nodes/%s", nodeID)
	_, err = zc.conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

// Unregister 注销节点
func (zc *ZookeeperCluster) Unregister(nodeID string) error {
	path := fmt.Sprintf("/scheduler/nodes/%s", nodeID)
	return zc.conn.Delete(path, -1)
}

// GetNodes 获取所有节点
func (zc *ZookeeperCluster) GetNodes() ([]NodeInfo, error) {
	children, _, err := zc.conn.Children("/scheduler/nodes")
	if err != nil {
		return nil, err
	}

	var nodes []NodeInfo
	for _, child := range children {
		path := fmt.Sprintf("/scheduler/nodes/%s", child)
		data, _, err := zc.conn.Get(path)
		if err != nil {
			continue
		}

		var node NodeInfo
		if err := json.Unmarshal(data, &node); err != nil {
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// AcquireLock 获取分布式锁
func (zc *ZookeeperCluster) AcquireLock(key string, ttl time.Duration) (bool, error) {
	lockPath := fmt.Sprintf("/scheduler/locks/%s", key)
	value := fmt.Sprintf("%s:%d", zc.nodeID, time.Now().UnixNano())

	// 尝试创建锁节点
	_, err := zc.conn.Create(lockPath, []byte(value), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err == nil {
		return true, nil
	}

	if err == zk.ErrNodeExists {
		return false, nil
	}

	return false, err
}

// ReleaseLock 释放分布式锁
func (zc *ZookeeperCluster) ReleaseLock(key string) error {
	lockPath := fmt.Sprintf("/scheduler/locks/%s", key)
	return zc.conn.Delete(lockPath, -1)
}

// StartElection 启动选举
func (zc *ZookeeperCluster) StartElection(callback func(bool)) error {
	zc.electionChan = make(chan bool, 1)
	go zc.runElection(callback)
	return nil
}

// StopElection 停止选举
func (zc *ZookeeperCluster) StopElection() error {
	close(zc.stopChan)
	return nil
}

// runElection 运行选举
func (zc *ZookeeperCluster) runElection(callback func(bool)) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			isLeader := zc.tryBecomeLeader()
			select {
			case zc.electionChan <- isLeader:
			default:
			}
			callback(isLeader)
		case <-zc.stopChan:
			return
		}
	}
}

// tryBecomeLeader 尝试成为领导者
func (zc *ZookeeperCluster) tryBecomeLeader() bool {
	value := fmt.Sprintf("%s:%d", zc.nodeID, time.Now().UnixNano())

	// 尝试创建领导者节点
	_, err := zc.conn.Create(zc.leaderPath, []byte(value), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err == nil {
		// 成功成为领导者
		return true
	}

	if err == zk.ErrNodeExists {
		// 检查当前领导者是否是自己
		data, _, err := zc.conn.Get(zc.leaderPath)
		if err != nil {
			return false
		}

		currentLeader := string(data)
		var leaderID string
		if _, err := fmt.Sscanf(currentLeader, "%s:", &leaderID); err == nil {
			return leaderID == zc.nodeID
		}
	}

	return false
}

// Broadcast 广播消息
func (zc *ZookeeperCluster) Broadcast(msg ClusterMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/scheduler/messages/%d", time.Now().UnixNano())
	_, err = zc.conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

// Subscribe 订阅消息
func (zc *ZookeeperCluster) Subscribe(callback func(ClusterMessage)) error {
	// 监听消息目录
	go func() {
		for {
			select {
			case <-zc.stopChan:
				return
			default:
				children, _, events, err := zc.conn.ChildrenW("/scheduler/messages")
				if err != nil {
					time.Sleep(time.Second)
					continue
				}

				// 处理现有消息
				for _, child := range children {
					path := fmt.Sprintf("/scheduler/messages/%s", child)
					data, _, err := zc.conn.Get(path)
					if err != nil {
						continue
					}

					var clusterMsg ClusterMessage
					if err := json.Unmarshal(data, &clusterMsg); err != nil {
						continue
					}

					// 忽略自己发送的消息
					if clusterMsg.NodeID == zc.nodeID {
						continue
					}

					callback(clusterMsg)
				}

				// 等待新消息
				select {
				case <-events:
					// 有新消息，继续循环处理
				case <-zc.stopChan:
					return
				}
			}
		}
	}()

	return nil
}

// Close 关闭集群连接
func (zc *ZookeeperCluster) Close() error {
	zc.cancel()
	close(zc.stopChan)
	zc.conn.Close()
	return nil
}

// GetLeader 获取当前领导者
func (zc *ZookeeperCluster) GetLeader() (string, error) {
	exists, _, err := zc.conn.Exists(zc.leaderPath)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("no leader found")
	}

	data, _, err := zc.conn.Get(zc.leaderPath)
	if err != nil {
		return "", err
	}

	leader := string(data)
	var leaderID string
	if _, err := fmt.Sscanf(leader, "%s:", &leaderID); err != nil {
		return "", err
	}

	return leaderID, nil
}

// IsLeader 检查是否为领导者
func (zc *ZookeeperCluster) IsLeader() bool {
	leaderID, err := zc.GetLeader()
	if err != nil {
		return false
	}
	return leaderID == zc.nodeID
}

// GetClusterInfo 获取集群信息
func (zc *ZookeeperCluster) GetClusterInfo() (map[string]interface{}, error) {
	nodes, err := zc.GetNodes()
	if err != nil {
		return nil, err
	}

	leaderID, _ := zc.GetLeader()

	info := map[string]interface{}{
		"total_nodes": len(nodes),
		"leader_id":   leaderID,
		"node_id":     zc.nodeID,
		"is_leader":   zc.IsLeader(),
		"nodes":       nodes,
	}

	return info, nil
}
