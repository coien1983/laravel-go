package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// DistributedScheduler 分布式调度器
type DistributedScheduler struct {
	*DefaultScheduler
	nodeID       string
	cluster      Cluster
	leader       bool
	leaderMu     sync.RWMutex
	electionMu   sync.Mutex
	stopElection chan struct{}
}

// Cluster 集群接口
type Cluster interface {
	// 节点管理
	Register(nodeID string, info NodeInfo) error
	Unregister(nodeID string) error
	GetNodes() ([]NodeInfo, error)

	// 分布式锁
	AcquireLock(key string, ttl time.Duration) (bool, error)
	ReleaseLock(key string) error

	// 选举
	StartElection(callback func(bool)) error
	StopElection() error

	// 消息广播
	Broadcast(msg ClusterMessage) error
	Subscribe(callback func(ClusterMessage)) error
}

// NodeInfo 节点信息
type NodeInfo struct {
	ID        string            `json:"id"`
	Address   string            `json:"address"`
	Port      int               `json:"port"`
	Status    string            `json:"status"` // online, offline, leader
	StartedAt time.Time         `json:"started_at"`
	LastSeen  time.Time         `json:"last_seen"`
	Metadata  map[string]string `json:"metadata"`
}

// ClusterMessage 集群消息
type ClusterMessage struct {
	Type      string    `json:"type"`
	NodeID    string    `json:"node_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"`
}

// TaskExecution 任务执行记录
type TaskExecution struct {
	TaskID    string     `json:"task_id"`
	NodeID    string     `json:"node_id"`
	Status    string     `json:"status"` // running, completed, failed
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Error     string     `json:"error,omitempty"`
}

// DistributedConfig 分布式配置
type DistributedConfig struct {
	NodeID                 string
	Cluster                Cluster
	ElectionTimeout        time.Duration
	LockTimeout            time.Duration
	HeartbeatInterval      time.Duration
	EnableLeaderElection   bool
	EnableTaskDistribution bool
}

// NewDistributedScheduler 创建分布式调度器
func NewDistributedScheduler(store Store, config DistributedConfig) *DistributedScheduler {
	if config.ElectionTimeout == 0 {
		config.ElectionTimeout = 30 * time.Second
	}
	if config.LockTimeout == 0 {
		config.LockTimeout = 10 * time.Second
	}
	if config.HeartbeatInterval == 0 {
		config.HeartbeatInterval = 5 * time.Second
	}

	ds := &DistributedScheduler{
		DefaultScheduler: NewScheduler(store),
		nodeID:           config.NodeID,
		cluster:          config.Cluster,
		stopElection:     make(chan struct{}),
	}

	return ds
}

// Start 启动分布式调度器
func (ds *DistributedScheduler) Start() error {
	// 注册节点
	if err := ds.registerNode(); err != nil {
		return fmt.Errorf("failed to register node: %w", err)
	}

	// 启动选举
	if err := ds.startElection(); err != nil {
		return fmt.Errorf("failed to start election: %w", err)
	}

	// 启动心跳
	go ds.heartbeat()

	// 启动消息订阅
	go ds.subscribeMessages()

	// 启动基础调度器
	return ds.DefaultScheduler.Start()
}

// Stop 停止分布式调度器
func (ds *DistributedScheduler) Stop() error {
	// 停止选举
	ds.stopElection <- struct{}{}

	// 注销节点
	if err := ds.cluster.Unregister(ds.nodeID); err != nil {
		return fmt.Errorf("failed to unregister node: %w", err)
	}

	// 停止基础调度器
	return ds.DefaultScheduler.Stop()
}

// IsLeader 检查是否为领导者
func (ds *DistributedScheduler) IsLeader() bool {
	ds.leaderMu.RLock()
	defer ds.leaderMu.RUnlock()
	return ds.leader
}

// GetClusterNodes 获取集群节点
func (ds *DistributedScheduler) GetClusterNodes() ([]NodeInfo, error) {
	return ds.cluster.GetNodes()
}

// registerNode 注册节点
func (ds *DistributedScheduler) registerNode() error {
	info := NodeInfo{
		ID:        ds.nodeID,
		Status:    "online",
		StartedAt: time.Now(),
		LastSeen:  time.Now(),
		Metadata: map[string]string{
			"version": "1.0.0",
			"type":    "scheduler",
		},
	}

	return ds.cluster.Register(ds.nodeID, info)
}

// startElection 启动选举
func (ds *DistributedScheduler) startElection() error {
	return ds.cluster.StartElection(func(isLeader bool) {
		ds.leaderMu.Lock()
		ds.leader = isLeader
		ds.leaderMu.Unlock()

		if isLeader {
			ds.onBecomeLeader()
		} else {
			ds.onLoseLeadership()
		}
	})
}

// onBecomeLeader 成为领导者
func (ds *DistributedScheduler) onBecomeLeader() {
	// 更新节点状态
	ds.updateNodeStatus("leader")

	// 广播领导者变更消息
	msg := ClusterMessage{
		Type:      "leader_changed",
		NodeID:    ds.nodeID,
		Timestamp: time.Now(),
	}
	ds.cluster.Broadcast(msg)
}

// onLoseLeadership 失去领导权
func (ds *DistributedScheduler) onLoseLeadership() {
	// 更新节点状态
	ds.updateNodeStatus("online")
}

// updateNodeStatus 更新节点状态
func (ds *DistributedScheduler) updateNodeStatus(status string) {
	// 这里应该更新节点信息到集群
	// 具体实现取决于集群接口
}

// heartbeat 心跳
func (ds *DistributedScheduler) heartbeat() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ds.updateNodeStatus("online")
		case <-ds.stopElection:
			return
		}
	}
}

// subscribeMessages 订阅消息
func (ds *DistributedScheduler) subscribeMessages() {
	ds.cluster.Subscribe(func(msg ClusterMessage) {
		ds.handleClusterMessage(msg)
	})
}

// handleClusterMessage 处理集群消息
func (ds *DistributedScheduler) handleClusterMessage(msg ClusterMessage) {
	switch msg.Type {
	case "task_execution_start":
		ds.handleTaskExecutionStart(msg)
	case "task_execution_complete":
		ds.handleTaskExecutionComplete(msg)
	case "leader_changed":
		ds.handleLeaderChanged(msg)
	}
}

// handleTaskExecutionStart 处理任务开始执行
func (ds *DistributedScheduler) handleTaskExecutionStart(msg ClusterMessage) {
	var execution TaskExecution
	if err := json.Unmarshal(msg.Data, &execution); err != nil {
		return
	}

	// 记录任务执行状态
	ds.recordTaskExecution(execution)
}

// handleTaskExecutionComplete 处理任务执行完成
func (ds *DistributedScheduler) handleTaskExecutionComplete(msg ClusterMessage) {
	var execution TaskExecution
	if err := json.Unmarshal(msg.Data, &execution); err != nil {
		return
	}

	// 更新任务执行状态
	ds.updateTaskExecution(execution)
}

// handleLeaderChanged 处理领导者变更
func (ds *DistributedScheduler) handleLeaderChanged(msg ClusterMessage) {
	// 可以在这里处理领导者变更逻辑
}

// executeTask 执行任务（分布式版本）
func (ds *DistributedScheduler) executeTask(task Task) {
	// 获取分布式锁
	lockKey := fmt.Sprintf("task_execution_%s", task.GetID())
	acquired, err := ds.cluster.AcquireLock(lockKey, 30*time.Second)
	if err != nil {
		return
	}
	if !acquired {
		// 其他节点正在执行此任务
		return
	}
	defer ds.cluster.ReleaseLock(lockKey)

	// 广播任务开始执行
	execution := TaskExecution{
		TaskID:    task.GetID(),
		NodeID:    ds.nodeID,
		Status:    "running",
		StartedAt: time.Now(),
	}
	ds.broadcastTaskExecution(execution)

	// 执行任务
	ctx, cancel := context.WithTimeout(ds.ctx, task.GetTimeout())
	defer cancel()

	err = task.GetHandler().Handle(ctx)

	// 更新执行状态
	endedAt := time.Now()
	execution.EndedAt = &endedAt
	if err != nil {
		execution.Status = "failed"
		execution.Error = err.Error()
	} else {
		execution.Status = "completed"
	}

	// 广播任务执行完成
	ds.broadcastTaskExecution(execution)

	// 更新任务状态
	ds.DefaultScheduler.executeTask(task)
}

// broadcastTaskExecution 广播任务执行状态
func (ds *DistributedScheduler) broadcastTaskExecution(execution TaskExecution) {
	data, err := json.Marshal(execution)
	if err != nil {
		return
	}

	msgType := "task_execution_start"
	if execution.EndedAt != nil {
		msgType = "task_execution_complete"
	}

	msg := ClusterMessage{
		Type:      msgType,
		NodeID:    ds.nodeID,
		Timestamp: time.Now(),
		Data:      data,
	}

	ds.cluster.Broadcast(msg)
}

// recordTaskExecution 记录任务执行
func (ds *DistributedScheduler) recordTaskExecution(execution TaskExecution) {
	// 这里可以记录到数据库或内存中
	// 用于跟踪任务执行状态
}

// updateTaskExecution 更新任务执行状态
func (ds *DistributedScheduler) updateTaskExecution(execution TaskExecution) {
	// 更新任务执行状态
	// 可以用于统计和监控
}

// GetDistributedStats 获取分布式统计
func (ds *DistributedScheduler) GetDistributedStats() DistributedStats {
	nodes, _ := ds.GetClusterNodes()

	return DistributedStats{
		NodeID:      ds.nodeID,
		IsLeader:    ds.IsLeader(),
		TotalNodes:  len(nodes),
		OnlineNodes: ds.countOnlineNodes(nodes),
		LeaderID:    ds.getLeaderID(nodes),
	}
}

// DistributedStats 分布式统计
type DistributedStats struct {
	NodeID      string `json:"node_id"`
	IsLeader    bool   `json:"is_leader"`
	TotalNodes  int    `json:"total_nodes"`
	OnlineNodes int    `json:"online_nodes"`
	LeaderID    string `json:"leader_id"`
}

// countOnlineNodes 统计在线节点
func (ds *DistributedScheduler) countOnlineNodes(nodes []NodeInfo) int {
	count := 0
	for _, node := range nodes {
		if node.Status == "online" || node.Status == "leader" {
			count++
		}
	}
	return count
}

// getLeaderID 获取领导者ID
func (ds *DistributedScheduler) getLeaderID(nodes []NodeInfo) string {
	for _, node := range nodes {
		if node.Status == "leader" {
			return node.ID
		}
	}
	return ""
}
