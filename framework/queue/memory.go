package queue

import (
	"context"
	"sort"
	"sync"
	"time"
)

// MemoryQueue 内存队列实现
type MemoryQueue struct {
	mu           sync.RWMutex
	jobs         []*BaseJob
	reservedJobs map[string]*BaseJob
	closed       bool
	stats        *QueueStats
}

// NewMemoryQueue 创建内存队列
func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		jobs:         make([]*BaseJob, 0),
		reservedJobs: make(map[string]*BaseJob),
		stats: &QueueStats{
			CreatedAt: time.Now(),
		},
	}
}

// Push 推送任务
func (q *MemoryQueue) Push(job Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	baseJob, ok := job.(*BaseJob)
	if !ok {
		return ErrInvalidJob
	}

	// 设置可用时间
	if baseJob.GetDelay() > 0 {
		baseJob.SetDelay(baseJob.GetDelay())
	}

	q.jobs = append(q.jobs, baseJob)
	q.stats.TotalJobs++
	q.stats.PendingJobs++
	q.stats.LastJobAt = time.Now()

	return nil
}

// PushBatch 批量推送任务
func (q *MemoryQueue) PushBatch(jobs []Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	for _, job := range jobs {
		baseJob, ok := job.(*BaseJob)
		if !ok {
			return ErrInvalidJob
		}

		if baseJob.GetDelay() > 0 {
			baseJob.SetDelay(baseJob.GetDelay())
		}

		q.jobs = append(q.jobs, baseJob)
		q.stats.TotalJobs++
		q.stats.PendingJobs++
	}

	q.stats.LastJobAt = time.Now()
	return nil
}

// Pop 弹出任务
func (q *MemoryQueue) Pop(ctx context.Context) (Job, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			// 定期检查，避免无限循环
		}

		q.mu.Lock()

		if q.closed {
			q.mu.Unlock()
			return nil, ErrQueueClosed
		}

		// 清理过期的保留任务
		q.cleanupExpiredJobs()

		// 查找可用的任务
		var job *BaseJob
		var jobIndex int = -1

		for i, j := range q.jobs {
			if j.IsAvailable() && !j.IsReserved() {
				job = j
				jobIndex = i
				break
			}
		}

		if job == nil {
			q.mu.Unlock()
			// 没有可用任务，继续等待
			continue
		}

		// 标记为已保留
		job.MarkAsReserved()
		q.reservedJobs[job.GetID()] = job

		// 从队列中移除
		q.jobs = append(q.jobs[:jobIndex], q.jobs[jobIndex+1:]...)
		q.stats.PendingJobs--
		q.stats.ReservedJobs++

		q.mu.Unlock()
		return job, nil
	}
}

// PopBatch 批量弹出任务
func (q *MemoryQueue) PopBatch(ctx context.Context, count int) ([]Job, error) {
	var jobs []Job

	for i := 0; i < count; i++ {
		job, err := q.Pop(ctx)
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				break
			}
			return jobs, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// Delete 删除任务
func (q *MemoryQueue) Delete(job Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	jobID := job.GetID()

	// 从保留任务中删除
	if _, exists := q.reservedJobs[jobID]; exists {
		delete(q.reservedJobs, jobID)
		q.stats.ReservedJobs--
		return nil
	}

	// 从队列中删除
	for i, j := range q.jobs {
		if j.GetID() == jobID {
			q.jobs = append(q.jobs[:i], q.jobs[i+1:]...)
			q.stats.PendingJobs--
			return nil
		}
	}

	return ErrJobNotFound
}

// Release 释放任务
func (q *MemoryQueue) Release(job Job, delay time.Duration) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	baseJob, ok := job.(*BaseJob)
	if !ok {
		return ErrInvalidJob
	}

	jobID := job.GetID()

	// 从保留任务中移除
	if _, exists := q.reservedJobs[jobID]; exists {
		delete(q.reservedJobs, jobID)
		q.stats.ReservedJobs--
	} else {
		return ErrJobNotFound
	}

	// 设置延迟时间
	if delay > 0 {
		baseJob.SetDelay(delay)
	} else {
		baseJob.AvailableAt = time.Now()
	}

	// 重置保留状态
	baseJob.ReservedAt = nil

	// 重新加入队列
	q.jobs = append(q.jobs, baseJob)
	q.stats.PendingJobs++

	return nil
}

// Later 延迟推送任务
func (q *MemoryQueue) Later(job Job, delay time.Duration) error {
	baseJob, ok := job.(*BaseJob)
	if !ok {
		return ErrInvalidJob
	}

	baseJob.SetDelay(delay)
	return q.Push(job)
}

// LaterBatch 批量延迟推送任务
func (q *MemoryQueue) LaterBatch(jobs []Job, delay time.Duration) error {
	for _, job := range jobs {
		baseJob, ok := job.(*BaseJob)
		if !ok {
			return ErrInvalidJob
		}
		baseJob.SetDelay(delay)
	}
	return q.PushBatch(jobs)
}

// Size 获取队列大小
func (q *MemoryQueue) Size() (int, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.closed {
		return 0, ErrQueueClosed
	}

	return len(q.jobs), nil
}

// Clear 清空队列
func (q *MemoryQueue) Clear() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	q.jobs = make([]*BaseJob, 0)
	q.reservedJobs = make(map[string]*BaseJob)
	q.stats.PendingJobs = 0
	q.stats.ReservedJobs = 0

	return nil
}

// Close 关闭队列
func (q *MemoryQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	if q.closed {
		return nil
	}
	
	q.closed = true
	
	// 清空任务，释放内存
	q.jobs = nil
	q.reservedJobs = nil
	
	// 重置统计信息
	q.stats = &QueueStats{}
	
	return nil
}

// GetStats 获取统计信息
func (q *MemoryQueue) GetStats() (QueueStats, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.closed {
		return QueueStats{}, ErrQueueClosed
	}

	stats := *q.stats
	stats.PendingJobs = int64(len(q.jobs))
	stats.ReservedJobs = int64(len(q.reservedJobs))

	return stats, nil
}

// cleanupExpiredJobs 清理过期的保留任务
func (q *MemoryQueue) cleanupExpiredJobs() {
	var expiredJobs []string

	for jobID, job := range q.reservedJobs {
		if job.IsExpired() {
			expiredJobs = append(expiredJobs, jobID)
		}
	}

	for _, jobID := range expiredJobs {
		job := q.reservedJobs[jobID]
		delete(q.reservedJobs, jobID)

		// 如果任务还可以重试，重新加入队列
		if job.CanRetry() {
			job.IncrementAttempts()
			job.ReservedAt = nil
			job.AvailableAt = time.Now().Add(5 * time.Second) // 重试延迟
			q.jobs = append(q.jobs, job)
			q.stats.PendingJobs++
		} else {
			q.stats.FailedJobs++
		}

		q.stats.ReservedJobs--
	}
}

// sortJobs 按优先级排序任务
func (q *MemoryQueue) sortJobs() {
	sort.Slice(q.jobs, func(i, j int) bool {
		// 首先按可用时间排序
		if !q.jobs[i].AvailableAt.Equal(q.jobs[j].AvailableAt) {
			return q.jobs[i].AvailableAt.Before(q.jobs[j].AvailableAt)
		}
		// 然后按优先级排序（优先级高的在前）
		return q.jobs[i].Priority > q.jobs[j].Priority
	})
}
