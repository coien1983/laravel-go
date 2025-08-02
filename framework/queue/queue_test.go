package queue

import (
	"context"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Fatal("NewManager should not return nil")
	}
	if manager.defaultQueue != "default" {
		t.Errorf("Expected default queue to be 'default', got %s", manager.defaultQueue)
	}
}

func TestManagerExtend(t *testing.T) {
	manager := NewManager()
	memoryQueue := NewMemoryQueue()
	
	manager.Extend("memory", memoryQueue)
	
	if len(manager.queues) != 1 {
		t.Errorf("Expected 1 queue, got %d", len(manager.queues))
	}
	
	if manager.queues["memory"] != memoryQueue {
		t.Error("Queue not properly extended")
	}
}

func TestManagerSetDefaultQueue(t *testing.T) {
	manager := NewManager()
	manager.SetDefaultQueue("test")
	
	if manager.defaultQueue != "test" {
		t.Errorf("Expected default queue to be 'test', got %s", manager.defaultQueue)
	}
}

func TestNewJob(t *testing.T) {
	payload := []byte("test payload")
	job := NewJob(payload, "test-queue")
	
	if job.GetID() == "" {
		t.Error("Job ID should not be empty")
	}
	
	if string(job.GetPayload()) != "test payload" {
		t.Errorf("Expected payload 'test payload', got %s", string(job.GetPayload()))
	}
	
	if job.GetQueue() != "test-queue" {
		t.Errorf("Expected queue 'test-queue', got %s", job.GetQueue())
	}
	
	if job.GetMaxAttempts() != 3 {
		t.Errorf("Expected max attempts 3, got %d", job.GetMaxAttempts())
	}
	
	if job.GetTimeout() != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", job.GetTimeout())
	}
}

func TestJobSerialization(t *testing.T) {
	job := NewJob([]byte("test"), "test-queue")
	job.SetPriority(10)
	job.AddTag("test-tag", "test-value")
	
	data, err := job.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize job: %v", err)
	}
	
	newJob := &BaseJob{}
	err = newJob.Deserialize(data)
	if err != nil {
		t.Fatalf("Failed to deserialize job: %v", err)
	}
	
	if newJob.GetID() != job.GetID() {
		t.Errorf("Expected ID %s, got %s", job.GetID(), newJob.GetID())
	}
	
	if string(newJob.GetPayload()) != string(job.GetPayload()) {
		t.Errorf("Expected payload %s, got %s", string(job.GetPayload()), string(newJob.GetPayload()))
	}
	
	if newJob.GetPriority() != job.GetPriority() {
		t.Errorf("Expected priority %d, got %d", job.GetPriority(), newJob.GetPriority())
	}
	
	if newJob.GetTags()["test-tag"] != job.GetTags()["test-tag"] {
		t.Errorf("Expected tag value %s, got %s", job.GetTags()["test-tag"], newJob.GetTags()["test-tag"])
	}
}

func TestMemoryQueue(t *testing.T) {
	queue := NewMemoryQueue()
	
	// 测试推送任务
	job1 := NewJob([]byte("job1"), "test-queue")
	err := queue.Push(job1)
	if err != nil {
		t.Fatalf("Failed to push job: %v", err)
	}
	
	job2 := NewJob([]byte("job2"), "test-queue")
	err = queue.Push(job2)
	if err != nil {
		t.Fatalf("Failed to push job: %v", err)
	}
	
	// 测试队列大小
	size, err := queue.Size()
	if err != nil {
		t.Fatalf("Failed to get queue size: %v", err)
	}
	if size != 2 {
		t.Errorf("Expected queue size 2, got %d", size)
	}
	
	// 测试弹出任务
	ctx := context.Background()
	poppedJob, err := queue.Pop(ctx)
	if err != nil {
		t.Fatalf("Failed to pop job: %v", err)
	}
	
	if poppedJob == nil {
		t.Fatal("Popped job should not be nil")
	}
	
	// 测试删除任务
	err = queue.Delete(poppedJob)
	if err != nil {
		t.Fatalf("Failed to delete job: %v", err)
	}
	
	// 测试清空队列
	err = queue.Clear()
	if err != nil {
		t.Fatalf("Failed to clear queue: %v", err)
	}
	
	size, err = queue.Size()
	if err != nil {
		t.Fatalf("Failed to get queue size: %v", err)
	}
	if size != 0 {
		t.Errorf("Expected queue size 0 after clear, got %d", size)
	}
}

func TestMemoryQueueDelay(t *testing.T) {
	queue := NewMemoryQueue()
	
	// 测试延迟任务
	job := NewJob([]byte("delayed job"), "test-queue")
	job.SetDelay(100 * time.Millisecond)
	
	err := queue.Push(job)
	if err != nil {
		t.Fatalf("Failed to push delayed job: %v", err)
	}
	
	// 立即尝试弹出，应该失败
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	
	_, err = queue.Pop(ctx)
	if err == nil {
		t.Error("Expected timeout error for delayed job")
	}
	
	// 等待延迟时间后，应该能弹出
	ctx = context.Background()
	poppedJob, err := queue.Pop(ctx)
	if err != nil {
		t.Fatalf("Failed to pop delayed job: %v", err)
	}
	
	if poppedJob == nil {
		t.Fatal("Popped delayed job should not be nil")
	}
}

func TestMemoryQueueBatch(t *testing.T) {
	queue := NewMemoryQueue()
	
	// 批量推送任务
	jobs := []Job{
		NewJob([]byte("job1"), "test-queue"),
		NewJob([]byte("job2"), "test-queue"),
		NewJob([]byte("job3"), "test-queue"),
	}
	
	err := queue.PushBatch(jobs)
	if err != nil {
		t.Fatalf("Failed to push batch jobs: %v", err)
	}
	
	// 批量弹出任务
	ctx := context.Background()
	poppedJobs, err := queue.PopBatch(ctx, 2)
	if err != nil {
		t.Fatalf("Failed to pop batch jobs: %v", err)
	}
	
	if len(poppedJobs) != 2 {
		t.Errorf("Expected 2 popped jobs, got %d", len(poppedJobs))
	}
}

func TestMemoryQueueStats(t *testing.T) {
	queue := NewMemoryQueue()
	
	// 推送一些任务
	for i := 0; i < 5; i++ {
		job := NewJob([]byte("test"), "test-queue")
		err := queue.Push(job)
		if err != nil {
			t.Fatalf("Failed to push job: %v", err)
		}
	}
	
	stats, err := queue.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}
	
	if stats.TotalJobs != 5 {
		t.Errorf("Expected total jobs 5, got %d", stats.TotalJobs)
	}
	
	if stats.PendingJobs != 5 {
		t.Errorf("Expected pending jobs 5, got %d", stats.PendingJobs)
	}
}

func TestNewWorker(t *testing.T) {
	queue := NewMemoryQueue()
	worker := NewWorker(queue, "test-queue")
	
	if worker == nil {
		t.Fatal("NewWorker should not return nil")
	}
	
	if worker.queueName != "test-queue" {
		t.Errorf("Expected queue name 'test-queue', got %s", worker.queueName)
	}
	
	if worker.workerID == "" {
		t.Error("Worker ID should not be empty")
	}
	
	if worker.status != "stopped" {
		t.Errorf("Expected status 'stopped', got %s", worker.status)
	}
}

func TestWorkerLifecycle(t *testing.T) {
	queue := NewMemoryQueue()
	worker := NewWorker(queue, "test-queue")
	
	// 测试启动
	err := worker.Start()
	if err != nil {
		t.Fatalf("Failed to start worker: %v", err)
	}
	
	status := worker.GetStatus()
	if status.Status != "running" {
		t.Errorf("Expected status 'running', got %s", status.Status)
	}
	
	// 测试暂停
	err = worker.Pause()
	if err != nil {
		t.Fatalf("Failed to pause worker: %v", err)
	}
	
	status = worker.GetStatus()
	if status.Status != "paused" {
		t.Errorf("Expected status 'paused', got %s", status.Status)
	}
	
	// 测试恢复
	err = worker.Resume()
	if err != nil {
		t.Fatalf("Failed to resume worker: %v", err)
	}
	
	status = worker.GetStatus()
	if status.Status != "running" {
		t.Errorf("Expected status 'running', got %s", status.Status)
	}
	
	// 测试停止
	err = worker.Stop()
	if err != nil {
		t.Fatalf("Failed to stop worker: %v", err)
	}
	
	status = worker.GetStatus()
	if status.Status != "stopped" {
		t.Errorf("Expected status 'stopped', got %s", status.Status)
	}
}

func TestWorkerPool(t *testing.T) {
	queue := NewMemoryQueue()
	pool := NewWorkerPool(queue, "test-queue", 3)
	
	// 测试启动工作进程池
	err := pool.Start()
	if err != nil {
		t.Fatalf("Failed to start worker pool: %v", err)
	}
	
	workers := pool.GetWorkers()
	if len(workers) != 3 {
		t.Errorf("Expected 3 workers, got %d", len(workers))
	}
	
	// 测试获取统计信息
	stats, err := pool.GetStats()
	if err != nil {
		t.Fatalf("Failed to get pool stats: %v", err)
	}
	
	if len(stats) != 3 {
		t.Errorf("Expected 3 worker stats, got %d", len(stats))
	}
	
	// 测试停止工作进程池
	err = pool.Stop()
	if err != nil {
		t.Fatalf("Failed to stop worker pool: %v", err)
	}
}

func TestGlobalFunctions(t *testing.T) {
	// 初始化全局管理器
	Init()
	
	// 注册内存队列
	memoryQueue := NewMemoryQueue()
	QueueManager.Extend("memory", memoryQueue)
	QueueManager.SetDefaultQueue("memory")
	
	// 测试全局推送
	job := NewJob([]byte("global test"), "test-queue")
	err := Push(job)
	if err != nil {
		t.Fatalf("Failed to push job globally: %v", err)
	}
	
	// 测试全局大小
	size, err := Size()
	if err != nil {
		t.Fatalf("Failed to get global size: %v", err)
	}
	if size != 1 {
		t.Errorf("Expected global size 1, got %d", size)
	}
	
	// 测试全局弹出
	ctx := context.Background()
	poppedJob, err := Pop(ctx)
	if err != nil {
		t.Fatalf("Failed to pop job globally: %v", err)
	}
	
	if poppedJob == nil {
		t.Fatal("Popped job should not be nil")
	}
	
	// 测试全局清空
	err = Clear()
	if err != nil {
		t.Fatalf("Failed to clear globally: %v", err)
	}
} 