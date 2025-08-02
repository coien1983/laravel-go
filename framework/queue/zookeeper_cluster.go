package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
)

// ZookeeperCluster ZooKeeper集群实现（复用定时器的实现）
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
	conn, _, err := zk.Connect(config.Servers, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ZooKeeper: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	zc := &ZookeeperCluster{
		conn:         conn,
		nodeID:       config.NodeID,
		ctx:          ctx,
		cancel:       cancel,
		stopChan:     make(chan struct{}),
		electionChan: make(chan bool, 1),
		leaderPath:   "/queue/leader",
	}

	// 创建必要的路径
	paths := []string{"/queue", "/queue/nodes", "/queue/locks", "/queue/messages"}
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

	path := fmt.Sprintf("/queue/nodes/%s", nodeID)

	// 创建临时节点
	_, err = zc.conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {
		return err
	}

	return nil
}

// Unregister 注销节点
func (zc *ZookeeperCluster) Unregister(nodeID string) error {
	path := fmt.Sprintf("/queue/nodes/%s", nodeID)
	return zc.conn.Delete(path, -1)
}

// GetNodes 获取所有节点
func (zc *ZookeeperCluster) GetNodes() ([]NodeInfo, error) {
	path := "/queue/nodes/"
	children, _, err := zc.conn.Children(path)
	if err != nil {
		return nil, err
	}

	var nodes []NodeInfo
	for _, child := range children {
		nodePath := path + child
		data, _, err := zc.conn.Get(nodePath)
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
	lockPath := fmt.Sprintf("/queue/locks/%s", key)

	// 创建临时顺序节点
	path, err := zc.conn.Create(lockPath, []byte(zc.nodeID), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		return false, err
	}

	// 获取所有锁节点
	children, _, err := zc.conn.Children("/queue/locks")
	if err != nil {
		return false, err
	}

	// 检查是否是最小的节点（获得锁）
	lockName := path[len("/queue/locks/"):]
	isFirst := true
	for _, child := range children {
		if child < lockName {
			isFirst = false
			break
		}
	}

	if isFirst {
		return true, nil
	}

	// 不是第一个，删除节点
	zc.conn.Delete(path, -1)
	return false, nil
}

// ReleaseLock 释放分布式锁
func (zc *ZookeeperCluster) ReleaseLock(key string) error {
	lockPath := fmt.Sprintf("/queue/locks/%s", key)

	// 获取所有锁节点
	children, _, err := zc.conn.Children("/queue/locks")
	if err != nil {
		return err
	}

	// 找到属于当前节点的锁
	for _, child := range children {
		if child == key {
			path := "/queue/locks/" + child
			data, _, err := zc.conn.Get(path)
			if err != nil {
				continue
			}

			if string(data) == zc.nodeID {
				return zc.conn.Delete(path, -1)
			}
		}
	}

	return nil
}

// StartElection 启动选举
func (zc *ZookeeperCluster) StartElection(callback func(bool)) error {
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
	// 创建临时顺序节点
	path, err := zc.conn.Create(zc.leaderPath, []byte(zc.nodeID), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		// 如果节点已存在，检查是否是自己
		children, _, err := zc.conn.Children("/queue")
		if err != nil {
			return false
		}

		for _, child := range children {
			if child == "leader" {
				data, _, err := zc.conn.Get("/queue/leader")
				if err != nil {
					return false
				}
				return string(data) == zc.nodeID
			}
		}
		return false
	}

	// 检查是否是最小的节点
	children, _, err := zc.conn.Children("/queue")
	if err != nil {
		return false
	}

	leaderName := path[len("/queue/"):]
	isFirst := true
	for _, child := range children {
		if child == "leader" && child < leaderName {
			isFirst = false
			break
		}
	}

	if isFirst {
		return true
	}

	// 不是第一个，删除节点
	zc.conn.Delete(path, -1)
	return false
}

// Broadcast 广播消息
func (zc *ZookeeperCluster) Broadcast(msg ClusterMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/queue/messages/%d", time.Now().UnixNano())
	_, err = zc.conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
	return err
}

// Subscribe 订阅消息
func (zc *ZookeeperCluster) Subscribe(callback func(ClusterMessage)) error {
	path := "/queue/messages/"

	// 监听消息
	go func() {
		for {
			select {
			case <-zc.stopChan:
				return
			default:
				children, _, events, err := zc.conn.ChildrenW(path)
				if err != nil {
					time.Sleep(1 * time.Second)
					continue
				}

				// 处理现有消息
				for _, child := range children {
					msgPath := path + child
					data, _, err := zc.conn.Get(msgPath)
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
				case event := <-events:
					if event.Type == zk.EventNodeChildrenChanged {
						// 有新消息，继续处理
						continue
					}
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
		return "", nil
	}

	data, _, err := zc.conn.Get(zc.leaderPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
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
