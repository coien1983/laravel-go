package microservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// ServiceClient 服务通信客户端
type ServiceClient struct {
	discovery  ServiceDiscovery
	httpClient *http.Client
	timeout    time.Duration
	retryCount int
	retryDelay time.Duration
}

// NewServiceClient 创建服务通信客户端
func NewServiceClient(discovery ServiceDiscovery, options ...ServiceClientOption) *ServiceClient {
	client := &ServiceClient{
		discovery:  discovery,
		timeout:    30 * time.Second,
		retryCount: 3,
		retryDelay: 1 * time.Second,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// 应用选项
	for _, option := range options {
		option(client)
	}

	return client
}

// ServiceClientOption 服务客户端选项
type ServiceClientOption func(*ServiceClient)

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ServiceClientOption {
	return func(c *ServiceClient) {
		c.timeout = timeout
		c.httpClient.Timeout = timeout
	}
}

// WithRetry 设置重试参数
func WithRetry(count int, delay time.Duration) ServiceClientOption {
	return func(c *ServiceClient) {
		c.retryCount = count
		c.retryDelay = delay
	}
}

// WithHTTPClient 设置自定义 HTTP 客户端
func WithHTTPClient(httpClient *http.Client) ServiceClientOption {
	return func(c *ServiceClient) {
		c.httpClient = httpClient
	}
}

// Call 调用服务
func (c *ServiceClient) Call(ctx context.Context, serviceName, method, path string, data interface{}) ([]byte, error) {
	// 发现服务
	service, err := c.discovery.DiscoverOne(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	// 构建请求 URL
	url := fmt.Sprintf("%s://%s:%d%s", service.Protocol, service.Address, service.Port, path)

	// 序列化请求数据
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "laravel-go-microservice-client")

	// 添加服务元数据到请求头
	for key, value := range service.Metadata {
		req.Header.Set(fmt.Sprintf("X-Service-%s", key), value)
	}

	// 执行请求（带重试）
	var resp *http.Response
	var lastErr error

	for i := 0; i <= c.retryCount; i++ {
		resp, lastErr = c.httpClient.Do(req)
		if lastErr == nil && resp.StatusCode < 500 {
			break
		}

		if i < c.retryCount {
			time.Sleep(c.retryDelay)
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to call service after %d retries: %w", c.retryCount, lastErr)
	}

	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("service returned error status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// CallJSON 调用服务并解析 JSON 响应
func (c *ServiceClient) CallJSON(ctx context.Context, serviceName, method, path string, requestData, responseData interface{}) error {
	responseBody, err := c.Call(ctx, serviceName, method, path, requestData)
	if err != nil {
		return err
	}

	if responseData != nil {
		err = json.Unmarshal(responseBody, responseData)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Get 发送 GET 请求
func (c *ServiceClient) Get(ctx context.Context, serviceName, path string) ([]byte, error) {
	return c.Call(ctx, serviceName, "GET", path, nil)
}

// GetJSON 发送 GET 请求并解析 JSON 响应
func (c *ServiceClient) GetJSON(ctx context.Context, serviceName, path string, responseData interface{}) error {
	return c.CallJSON(ctx, serviceName, "GET", path, nil, responseData)
}

// Post 发送 POST 请求
func (c *ServiceClient) Post(ctx context.Context, serviceName, path string, data interface{}) ([]byte, error) {
	return c.Call(ctx, serviceName, "POST", path, data)
}

// PostJSON 发送 POST 请求并解析 JSON 响应
func (c *ServiceClient) PostJSON(ctx context.Context, serviceName, path string, requestData, responseData interface{}) error {
	return c.CallJSON(ctx, serviceName, "POST", path, requestData, responseData)
}

// Put 发送 PUT 请求
func (c *ServiceClient) Put(ctx context.Context, serviceName, path string, data interface{}) ([]byte, error) {
	return c.Call(ctx, serviceName, "PUT", path, data)
}

// PutJSON 发送 PUT 请求并解析 JSON 响应
func (c *ServiceClient) PutJSON(ctx context.Context, serviceName, path string, requestData, responseData interface{}) error {
	return c.CallJSON(ctx, serviceName, "PUT", path, requestData, responseData)
}

// Delete 发送 DELETE 请求
func (c *ServiceClient) Delete(ctx context.Context, serviceName, path string) ([]byte, error) {
	return c.Call(ctx, serviceName, "DELETE", path, nil)
}

// DeleteJSON 发送 DELETE 请求并解析 JSON 响应
func (c *ServiceClient) DeleteJSON(ctx context.Context, serviceName, path string, responseData interface{}) error {
	return c.CallJSON(ctx, serviceName, "DELETE", path, nil, responseData)
}

// CircuitBreaker 熔断器接口
type CircuitBreaker interface {
	// Execute 执行操作
	Execute(ctx context.Context, operation func() error) error

	// IsOpen 检查熔断器是否开启
	IsOpen() bool

	// Reset 重置熔断器
	Reset()
}

// SimpleCircuitBreaker 简单熔断器实现
type SimpleCircuitBreaker struct {
	failureThreshold int
	failureCount     int
	lastFailureTime  time.Time
	timeout          time.Duration
	state            CircuitBreakerState
	mutex            sync.RWMutex
}

// CircuitBreakerState 熔断器状态
type CircuitBreakerState string

const (
	CircuitBreakerClosed CircuitBreakerState = "closed"
	CircuitBreakerOpen   CircuitBreakerState = "open"
	CircuitBreakerHalf   CircuitBreakerState = "half-open"
)

// NewSimpleCircuitBreaker 创建简单熔断器
func NewSimpleCircuitBreaker(failureThreshold int, timeout time.Duration) *SimpleCircuitBreaker {
	return &SimpleCircuitBreaker{
		failureThreshold: failureThreshold,
		timeout:          timeout,
		state:            CircuitBreakerClosed,
	}
}

// Execute 执行操作
func (cb *SimpleCircuitBreaker) Execute(ctx context.Context, operation func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case CircuitBreakerOpen:
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = CircuitBreakerHalf
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	case CircuitBreakerHalf:
		// 允许一次尝试
		cb.state = CircuitBreakerClosed
	}

	// 执行操作
	err := operation()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		if cb.failureCount >= cb.failureThreshold {
			cb.state = CircuitBreakerOpen
		}
	} else {
		// 成功时重置
		cb.failureCount = 0
		cb.state = CircuitBreakerClosed
	}

	return err
}

// IsOpen 检查熔断器是否开启
func (cb *SimpleCircuitBreaker) IsOpen() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state == CircuitBreakerOpen
}

// Reset 重置熔断器
func (cb *SimpleCircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failureCount = 0
	cb.state = CircuitBreakerClosed
}
