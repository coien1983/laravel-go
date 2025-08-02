package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQQueue RabbitMQ 队列实现
type RabbitMQQueue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  RabbitMQConfig
	stats   QueueStats
}

// RabbitMQConfig RabbitMQ 配置
type RabbitMQConfig struct {
	URL          string
	Exchange     string
	QueueName    string
	RoutingKey   string
	PrefetchCount int
	AutoDelete   bool
	Durable      bool
}

// NewRabbitMQQueue 创建 RabbitMQ 队列
func NewRabbitMQQueue(config RabbitMQConfig) (*RabbitMQQueue, error) {
	// 连接 RabbitMQ
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// 设置预取数量
	if config.PrefetchCount > 0 {
		err = channel.Qos(config.PrefetchCount, 0, false)
		if err != nil {
			channel.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to set QoS: %w", err)
		}
	}

	// 声明交换机
	err = channel.ExchangeDeclare(
		config.Exchange, // name
		"direct",        // type
		config.Durable,  // durable
		config.AutoDelete, // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// 声明队列
	_, err = channel.QueueDeclare(
		config.QueueName,  // name
		config.Durable,    // durable
		config.AutoDelete, // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// 绑定队列到交换机
	err = channel.QueueBind(
		config.QueueName,  // queue name
		config.RoutingKey, // routing key
		config.Exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &RabbitMQQueue{
		conn:    conn,
		channel: channel,
		config:  config,
		stats:   QueueStats{CreatedAt: time.Now()},
	}, nil
}

// Push 推送任务
func (rq *RabbitMQQueue) Push(job Job) error {
	payload, err := job.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}

	// 设置消息属性
	headers := amqp.Table{
		"job_id":       job.GetID(),
		"queue":        job.GetQueue(),
		"attempts":     job.GetAttempts(),
		"max_attempts": job.GetMaxAttempts(),
		"priority":     job.GetPriority(),
		"created_at":   job.GetCreatedAt().Unix(),
	}

	// 添加标签
	if len(job.GetTags()) > 0 {
		headers["tags"] = job.GetTags()
	}

	// 发布消息
	err = rq.channel.PublishWithContext(context.Background(),
		rq.config.Exchange,   // exchange
		rq.config.RoutingKey, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
			Priority:     uint8(job.GetPriority()),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// 更新统计
	rq.stats.TotalJobs++
	rq.stats.LastJobAt = time.Now()

	return nil
}

// PushBatch 批量推送任务
func (rq *RabbitMQQueue) PushBatch(jobs []Job) error {
	for _, job := range jobs {
		if err := rq.Push(job); err != nil {
			return err
		}
	}
	return nil
}

// Pop 弹出任务
func (rq *RabbitMQQueue) Pop(ctx context.Context) (Job, error) {
	msg, ok, err := rq.channel.Get(rq.config.QueueName, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	if !ok {
		return nil, ErrQueueEmpty
	}

	// 解析消息头
	jobID, _ := msg.Headers["job_id"].(string)
	queue, _ := msg.Headers["queue"].(string)
	attempts, _ := msg.Headers["attempts"].(int)
	maxAttempts, _ := msg.Headers["max_attempts"].(int)
	priority, _ := msg.Headers["priority"].(int)
	createdAt, _ := msg.Headers["created_at"].(int64)
	tags, _ := msg.Headers["tags"].(amqp.Table)

	// 创建任务
	job := NewJob(msg.Body, queue)
	job.SetID(jobID)
	job.SetAttempts(attempts)
	job.SetMaxAttempts(maxAttempts)
	job.SetPriority(priority)
	job.SetCreatedAt(time.Unix(createdAt, 0))
	job.SetReservedAt(time.Now())

	// 设置标签
	if tags != nil {
		tagMap := make(map[string]string)
		for k, v := range tags {
			if str, ok := v.(string); ok {
				tagMap[k] = str
			}
		}
		job.SetTags(tagMap)
	}

	// 标记为已保留
	job.MarkAsReserved()

	// 更新统计
	rq.stats.ReservedJobs++

	return job, nil
}

// PopBatch 批量弹出任务
func (rq *RabbitMQQueue) PopBatch(ctx context.Context, count int) ([]Job, error) {
	var jobs []Job
	for i := 0; i < count; i++ {
		job, err := rq.Pop(ctx)
		if err != nil {
			if err == ErrQueueEmpty {
				break
			}
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// Delete 删除任务
func (rq *RabbitMQQueue) Delete(job Job) error {
	// RabbitMQ 中消息被消费后自动删除
	job.MarkAsCompleted()
	rq.stats.CompletedJobs++
	rq.stats.ReservedJobs--
	return nil
}

// Release 释放任务
func (rq *RabbitMQQueue) Release(job Job, delay time.Duration) error {
	// 重新发布消息
	payload, err := job.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}

	headers := amqp.Table{
		"job_id":       job.GetID(),
		"queue":        job.GetQueue(),
		"attempts":     job.GetAttempts(),
		"max_attempts": job.GetMaxAttempts(),
		"priority":     job.GetPriority(),
		"created_at":   job.GetCreatedAt().Unix(),
	}

	if len(job.GetTags()) > 0 {
		headers["tags"] = job.GetTags()
	}

	err = rq.channel.PublishWithContext(context.Background(),
		rq.config.Exchange,
		rq.config.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
			Priority:     uint8(job.GetPriority()),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to release job: %w", err)
	}

	rq.stats.ReservedJobs--
	return nil
}

// Later 延迟推送任务
func (rq *RabbitMQQueue) Later(job Job, delay time.Duration) error {
	// 设置延迟时间
	job.SetDelay(delay)
	job.SetAvailableAt(time.Now().Add(delay))
	return rq.Push(job)
}

// LaterBatch 批量延迟推送任务
func (rq *RabbitMQQueue) LaterBatch(jobs []Job, delay time.Duration) error {
	for _, job := range jobs {
		if err := rq.Later(job, delay); err != nil {
			return err
		}
	}
	return nil
}

// Size 获取队列大小
func (rq *RabbitMQQueue) Size() (int, error) {
	queue, err := rq.channel.QueueInspect(rq.config.QueueName)
	if err != nil {
		return 0, fmt.Errorf("failed to inspect queue: %w", err)
	}
	return queue.Messages, nil
}

// Clear 清空队列
func (rq *RabbitMQQueue) Clear() error {
	_, err := rq.channel.QueuePurge(rq.config.QueueName, false)
	if err != nil {
		return fmt.Errorf("failed to purge queue: %w", err)
	}
	
	// 重置统计
	rq.stats = QueueStats{CreatedAt: time.Now()}
	return nil
}

// Close 关闭连接
func (rq *RabbitMQQueue) Close() error {
	if rq.channel != nil {
		rq.channel.Close()
	}
	if rq.conn != nil {
		rq.conn.Close()
	}
	return nil
}

// GetStats 获取统计信息
func (rq *RabbitMQQueue) GetStats() (QueueStats, error) {
	size, err := rq.Size()
	if err != nil {
		return rq.stats, err
	}
	
	stats := rq.stats
	stats.PendingJobs = int64(size)
	return stats, nil
} 