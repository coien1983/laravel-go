package plugins

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RateLimiter 限流器
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter 创建限流器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// 清理过期的请求记录
	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[key] = validRequests
	}

	// 检查是否超过限制
	if len(rl.requests[key]) >= rl.limit {
		return false
	}

	// 记录当前请求
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

// GetRemaining 获取剩余请求次数
func (rl *RateLimiter) GetRemaining(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if requests, exists := rl.requests[key]; exists {
		return rl.limit - len(requests)
	}
	return rl.limit
}

// Reset 重置限流器
func (rl *RateLimiter) Reset(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.requests, key)
}

// RateLimitPlugin 限流插件
type RateLimitPlugin struct {
	limiter *RateLimiter
}

// NewRateLimitPlugin 创建限流插件
func NewRateLimitPlugin(limit int, window time.Duration) *RateLimitPlugin {
	return &RateLimitPlugin{
		limiter: NewRateLimiter(limit, window),
	}
}

// Process 处理请求
func (p *RateLimitPlugin) Process(ctx context.Context, req interface{}) (interface{}, error) {
	// TODO: 从请求中提取客户端标识
	clientID := "default"
	
	if !p.limiter.Allow(clientID) {
		return nil, fmt.Errorf("rate limit exceeded")
	}
	
	return req, nil
}