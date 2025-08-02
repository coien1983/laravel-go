package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// DistributedQueue 分布式队列
type DistributedQueue struct {
	*MemoryQueue
	cluster      Cluster
	nodeID       string
	leader       bool
	leaderMu     sync.RWMutex
	electionMu   sync.Mutex
	stopChan     chan struct{}
	workerPool   *DistributedWorkerPool
}

// Cluster 集群接口（复用定时器的集群接口）
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

// JobExecution 任务执行记录
type JobExecution struct {
	JobID     string     `json:"job_id"`
	NodeID    string     `json:"node_id"`
	Status    string     `json:"status"` // processing, completed, failed
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
	EnableJobDistribution  bool
	WorkerCount            int
	MaxConcurrency         int
}

// NewDistributedQueue 创建分布式队列
func NewDistributedQueue(config DistributedConfig) *DistributedQueue {
	if config.ElectionTimeout == 0 {
		config.ElectionTimeout = 30 * time.Second
	}
	if config.LockTimeout == 0 {
		config.LockTimeout = 10 * time.Second
	}
	if config.HeartbeatInterval == 0 {
		config.HeartbeatInterval = 5 * time.Second
	}
	if config.WorkerCount == 0 {
		config.WorkerCount = 5
	}
	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = 10
	}

	dq := &DistributedQueue{
		MemoryQueue: NewMemoryQueue(),
		nodeID:      config.NodeID,
		cluster:     config.Cluster,
		stopChan:    make(chan struct{}),
	}

	// 创建工作进程池
	dq.workerPool = NewDistributedWorkerPool(dq, config.WorkerCount, config.MaxConcurrency)

	return dq
}

// Start 启动分布式队列
func (dq *DistributedQueue) Start() error {
	// 注册节点
	if err := dq.registerNode(); err != nil {
		return fmt.Errorf("failed to register node: %w", err)
	}

	// 启动选举
	if err := dq.startElection(); err != nil {
		return fmt.Errorf("failed to start election: %w", err)
	}

	// 启动心跳
	go dq.heartbeat()

	// 启动消息订阅
	go dq.subscribeMessages()

	// 启动工作进程池
	if err := dq.workerPool.Start(); err != nil {
		return fmt.Errorf("failed to start worker pool: %w", err)
	}

	return nil
}

// Stop 停止分布式队列
func (dq *DistributedQueue) Stop() error {
	// 停止工作进程池
	if err := dq.workerPool.Stop(); err != nil {
		return fmt.Errorf("failed to stop worker pool: %w", err)
	}

	// 停止选举
	dq.stopChan <- struct{}{}

	// 注销节点
	if err := dq.cluster.Unregister(dq.nodeID); err != nil {
		return fmt.Errorf("failed to unregister node: %w", err)
	}

	return nil
}

// IsLeader 检查是否为领导者
func (dq *DistributedQueue) IsLeader() bool {
	dq.leaderMu.RLock()
	defer dq.leaderMu.RUnlock()
	return dq.leader
}

// GetClusterNodes 获取集群节点
func (dq *DistributedQueue) GetClusterNodes() ([]NodeInfo, error) {
	return dq.cluster.GetNodes()
}

// Push 推送任务（分布式版本）
func (dq *DistributedQueue) Push(job Job) error {
	// 如果是领导者，直接推送
	if dq.IsLeader() {
		return dq.MemoryQueue.Push(job)
	}

	// 如果不是领导者，广播任务到集群
	return dq.broadcastJob(job)
}

// Pop 弹出任务（分布式版本）
func (dq *DistributedQueue) Pop(ctx context.Context) (Job, error) {
	// 获取分布式锁
	lockKey := fmt.Sprintf("queue_pop_%s", dq.nodeID)
	acquired, err := dq.cluster.AcquireLock(lockKey, 30*time.Second)
	if err != nil {
		return nil, err
	}
	if !acquired {
		// 其他节点正在处理，等待
		time.Sleep(100 * time.Millisecond)
		return dq.Pop(ctx)
	}
	defer dq.cluster.ReleaseLock(lockKey)

	// 从本地队列弹出任务
	job, err := dq.MemoryQueue.Pop(ctx)
	if err != nil {
		return nil, err
	}

	// 广播任务开始执行
	execution := JobExecution{
		JobID:     job.GetID(),
		NodeID:    dq.nodeID,
		Status:    "processing",
		StartedAt: time.Now(),
	}
	dq.broadcastJobExecution(execution)

	return job, nil
}

// registerNode 注册节点
func (dq *DistributedQueue) registerNode() error {
	info := NodeInfo{
		ID:        dq.nodeID,
		Status:    "online",
		StartedAt: time.Now(),
		LastSeen:  time.Now(),
		Metadata: map[string]string{
			"version": "1.0.0",
			"type":    "queue",
			"queue":   dq.nodeID,
		},
	}

	return dq.cluster.Register(dq.nodeID, info)
}

// startElection 启动选举
func (dq *DistributedQueue) startElection() error {
	return dq.cluster.StartElection(func(isLeader bool) {
		dq.leaderMu.Lock()
		dq.leader = isLeader
		dq.leaderMu.Unlock()

		if isLeader {
			dq.onBecomeLeader()
		} else {
			dq.onLoseLeadership()
		}
	})
}

// onBecomeLeader 成为领导者
func (dq *DistributedQueue) onBecomeLeader() {
	// 更新节点状态
	dq.updateNodeStatus("leader")

	// 广播领导者变更消息
	msg := ClusterMessage{
		Type:      "leader_changed",
		NodeID:    dq.nodeID,
		Timestamp: time.Now(),
	}
	dq.cluster.Broadcast(msg)
}

// onLoseLeadership 失去领导权
func (dq *DistributedQueue) onLoseLeadership() {
	// 更新节点状态
	dq.updateNodeStatus("online")
}

// updateNodeStatus 更新节点状态
func (dq *DistributedQueue) updateNodeStatus(status string) {
	// 这里应该更新节点信息到集群
	// 具体实现取决于集群接口
}

// heartbeat 心跳
func (dq *DistributedQueue) heartbeat() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dq.updateNodeStatus("online")
		case <-dq.stopChan:
			return
		}
	}
}

// subscribeMessages 订阅消息
func (dq *DistributedQueue) subscribeMessages() {
	dq.cluster.Subscribe(func(msg ClusterMessage) {
		dq.handleClusterMessage(msg)
	})
}

// handleClusterMessage 处理集群消息
func (dq *DistributedQueue) handleClusterMessage(msg ClusterMessage) {
	switch msg.Type {
	case "job_push":
		dq.handleJobPush(msg)
	case "job_execution_start":
		dq.handleJobExecutionStart(msg)
	case "job_execution_complete":
		dq.handleJobExecutionComplete(msg)
	case "leader_changed":
		dq.handleLeaderChanged(msg)
	}
}

// handleJobPush 处理任务推送
func (dq *DistributedQueue) handleJobPush(msg ClusterMessage) {
	var jobData JobData
	if err := json.Unmarshal(msg.Data, &jobData); err != nil {
		return
	}

	// 创建任务
	job := NewJob(jobData.Payload, jobData.Queue)
	// 设置任务属性
	job.SetDelay(jobData.Delay)
	job.SetTimeout(jobData.Timeout)
	job.SetPriority(jobData.Priority)
	for key, value := range jobData.Tags {
		job.AddTag(key, value)
	}

	// 添加到本地队列
	dq.MemoryQueue.Push(job)
}

// handleJobExecutionStart 处理任务开始执行
func (dq *DistributedQueue) handleJobExecutionStart(msg ClusterMessage) {
	var execution JobExecution
	if err := json.Unmarshal(msg.Data, &execution); err != nil {
		return
	}

	// 记录任务执行状态
	dq.recordJobExecution(execution)
}

// handleJobExecutionComplete 处理任务执行完成
func (dq *DistributedQueue) handleJobExecutionComplete(msg ClusterMessage) {
	var execution JobExecution
	if err := json.Unmarshal(msg.Data, &execution); err != nil {
		return
	}

	// 更新任务执行状态
	dq.updateJobExecution(execution)
}

// handleLeaderChanged 处理领导者变更
func (dq *DistributedQueue) handleLeaderChanged(msg ClusterMessage) {
	// 可以在这里处理领导者变更逻辑
}

// broadcastJob 广播任务
func (dq *DistributedQueue) broadcastJob(job Job) error {
	jobData := JobData{
		ID:      job.GetID(),
		Payload: job.GetPayload(),
		Queue:   job.GetQueue(),
		Delay:   job.GetDelay(),
		Timeout: job.GetTimeout(),
		Priority: job.GetPriority(),
		Tags:    job.GetTags(),
	}

	data, err := json.Marshal(jobData)
	if err != nil {
		return err
	}

	msg := ClusterMessage{
		Type:      "job_push",
		NodeID:    dq.nodeID,
		Timestamp: time.Now(),
		Data:      data,
	}

	return dq.cluster.Broadcast(msg)
}

// broadcastJobExecution 广播任务执行状态
func (dq *DistributedQueue) broadcastJobExecution(execution JobExecution) {
	data, err := json.Marshal(execution)
	if err != nil {
		return
	}

	msgType := "job_execution_start"
	if execution.EndedAt != nil {
		msgType = "job_execution_complete"
	}

	msg := ClusterMessage{
		Type:      msgType,
		NodeID:    dq.nodeID,
		Timestamp: time.Now(),
		Data:      data,
	}

	dq.cluster.Broadcast(msg)
}

// recordJobExecution 记录任务执行
func (dq *DistributedQueue) recordJobExecution(execution JobExecution) {
	// 这里可以记录到数据库或内存中
	// 用于跟踪任务执行状态
}

// updateJobExecution 更新任务执行状态
func (dq *DistributedQueue) updateJobExecution(execution JobExecution) {
	// 更新任务执行状态
	// 可以用于统计和监控
}

// GetDistributedStats 获取分布式统计
func (dq *DistributedQueue) GetDistributedStats() DistributedStats {
	nodes, _ := dq.GetClusterNodes()
	stats, _ := dq.GetStats()
	
	return DistributedStats{
		NodeID:      dq.nodeID,
		IsLeader:    dq.IsLeader(),
		TotalNodes:  len(nodes),
		OnlineNodes: dq.countOnlineNodes(nodes),
		LeaderID:    dq.getLeaderID(nodes),
		QueueStats:  stats,
	}
}

// DistributedStats 分布式统计
type DistributedStats struct {
	NodeID      string     `json:"node_id"`
	IsLeader    bool       `json:"is_leader"`
	TotalNodes  int        `json:"total_nodes"`
	OnlineNodes int        `json:"online_nodes"`
	LeaderID    string     `json:"leader_id"`
	QueueStats  QueueStats `json:"queue_stats"`
}

// JobData 任务数据
type JobData struct {
	ID       string            `json:"id"`
	Payload  []byte            `json:"payload"`
	Queue    string            `json:"queue"`
	Delay    time.Duration     `json:"delay"`
	Timeout  time.Duration     `json:"timeout"`
	Priority int               `json:"priority"`
	Tags     map[string]string `json:"tags"`
}

// countOnlineNodes 统计在线节点
func (dq *DistributedQueue) countOnlineNodes(nodes []NodeInfo) int {
	count := 0
	for _, node := range nodes {
		if node.Status == "online" || node.Status == "leader" {
			count++
		}
	}
	return count
}

// getLeaderID 获取领导者ID
func (dq *DistributedQueue) getLeaderID(nodes []NodeInfo) string {
	for _, node := range nodes {
		if node.Status == "leader" {
			return node.ID
		}
	}
	return ""
} 