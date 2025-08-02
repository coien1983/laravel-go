package microservice

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// StreamManager gRPC 流管理器
type StreamManager struct {
	streams map[string]*StreamInfo
	mutex   sync.RWMutex
	// 添加清理机制
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

// StreamInfo 流信息
type StreamInfo struct {
	ID        string
	Type      StreamType
	StartTime time.Time
	Metadata  metadata.MD
	Context   context.Context
	Cancel    context.CancelFunc
	// 添加超时时间
	Timeout time.Duration
}

// StreamType 流类型
type StreamType string

const (
	StreamTypeUnary         StreamType = "unary"
	StreamTypeClient        StreamType = "client"
	StreamTypeServer        StreamType = "server"
	StreamTypeBidirectional StreamType = "bidirectional"
)

// NewStreamManager 创建流管理器
func NewStreamManager() *StreamManager {
	sm := &StreamManager{
		streams:       make(map[string]*StreamInfo),
		cleanupTicker: time.NewTicker(30 * time.Second), // 每30秒清理一次
		stopChan:      make(chan struct{}),
	}
	
	// 启动清理协程
	go sm.cleanupRoutine()
	
	return sm
}

// cleanupRoutine 清理协程
func (sm *StreamManager) cleanupRoutine() {
	for {
		select {
		case <-sm.cleanupTicker.C:
			sm.cleanupExpiredStreams()
		case <-sm.stopChan:
			return
		}
	}
}

// cleanupExpiredStreams 清理过期的流
func (sm *StreamManager) cleanupExpiredStreams() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	now := time.Now()
	var expiredStreams []string
	
	for id, stream := range sm.streams {
		// 检查流是否超时（默认5分钟）
		timeout := stream.Timeout
		if timeout == 0 {
			timeout = 5 * time.Minute
		}
		
		if now.Sub(stream.StartTime) > timeout {
			expiredStreams = append(expiredStreams, id)
		}
	}
	
	// 清理过期的流
	for _, id := range expiredStreams {
		stream := sm.streams[id]
		stream.Cancel()
		delete(sm.streams, id)
	}
}

// RegisterStream 注册流
func (sm *StreamManager) RegisterStream(id string, streamType StreamType, ctx context.Context) {
	sm.RegisterStreamWithTimeout(id, streamType, ctx, 5*time.Minute)
}

// RegisterStreamWithTimeout 注册流（带超时）
func (sm *StreamManager) RegisterStreamWithTimeout(id string, streamType StreamType, ctx context.Context, timeout time.Duration) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	md, _ := metadata.FromIncomingContext(ctx)
	ctx, cancel := context.WithCancel(ctx)

	sm.streams[id] = &StreamInfo{
		ID:        id,
		Type:      streamType,
		StartTime: time.Now(),
		Metadata:  md,
		Context:   ctx,
		Cancel:    cancel,
		Timeout:   timeout,
	}
}

// UnregisterStream 注销流
func (sm *StreamManager) UnregisterStream(id string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if stream, exists := sm.streams[id]; exists {
		stream.Cancel()
		delete(sm.streams, id)
	}
}

// GetStream 获取流信息
func (sm *StreamManager) GetStream(id string) (*StreamInfo, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	stream, exists := sm.streams[id]
	return stream, exists
}

// ListStreams 列出所有流
func (sm *StreamManager) ListStreams() []*StreamInfo {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	streams := make([]*StreamInfo, 0, len(sm.streams))
	for _, stream := range sm.streams {
		streams = append(streams, stream)
	}
	return streams
}

// CloseAllStreams 关闭所有流
func (sm *StreamManager) CloseAllStreams() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for _, stream := range sm.streams {
		stream.Cancel()
	}
	sm.streams = make(map[string]*StreamInfo)
}

// Close 关闭流管理器
func (sm *StreamManager) Close() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// 停止清理协程
	close(sm.stopChan)
	sm.cleanupTicker.Stop()
	
	// 关闭所有流
	for _, stream := range sm.streams {
		stream.Cancel()
	}
	
	// 清空流映射
	sm.streams = nil
}

// StreamHandler 流处理器接口
type StreamHandler interface {
	HandleStream(ctx context.Context, stream interface{}) error
}

// BidirectionalStreamHandler 双向流处理器
type BidirectionalStreamHandler struct {
	manager *StreamManager
}

// NewBidirectionalStreamHandler 创建双向流处理器
func NewBidirectionalStreamHandler(manager *StreamManager) *BidirectionalStreamHandler {
	return &BidirectionalStreamHandler{
		manager: manager,
	}
}

// HandleBidirectionalStream 处理双向流
func (h *BidirectionalStreamHandler) HandleBidirectionalStream(stream interface{}) error {
	// 这里需要根据具体的 gRPC 流类型来实现
	// 例如：pb.UserService_StreamChatServer
	return nil
}

// StreamClient gRPC 流客户端
type StreamClient struct {
	conn   *grpc.ClientConn
	client interface{}
}

// NewStreamClient 创建流客户端
func NewStreamClient(address string, opts ...grpc.DialOption) (*StreamClient, error) {
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &StreamClient{
		conn: conn,
	}, nil
}

// Close 关闭连接
func (c *StreamClient) Close() error {
	return c.conn.Close()
}

// GetConnection 获取连接
func (c *StreamClient) GetConnection() *grpc.ClientConn {
	return c.conn
}

// StreamServer gRPC 流服务器
type StreamServer struct {
	server   *grpc.Server
	manager  *StreamManager
	handlers map[string]StreamHandler
	mutex    sync.RWMutex
}

// NewStreamServer 创建流服务器
func NewStreamServer(opts ...grpc.ServerOption) *StreamServer {
	return &StreamServer{
		server:   grpc.NewServer(opts...),
		manager:  NewStreamManager(),
		handlers: make(map[string]StreamHandler),
	}
}

// RegisterHandler 注册流处理器
func (s *StreamServer) RegisterHandler(method string, handler StreamHandler) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.handlers[method] = handler
}

// GetHandler 获取流处理器
func (s *StreamServer) GetHandler(method string) (StreamHandler, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	handler, exists := s.handlers[method]
	return handler, exists
}

// GetServer 获取底层服务器
func (s *StreamServer) GetServer() *grpc.Server {
	return s.server
}

// GetStreamManager 获取流管理器
func (s *StreamServer) GetStreamManager() *StreamManager {
	return s.manager
}

// StreamWrapper 流包装器
type StreamWrapper struct {
	stream   interface{}
	manager  *StreamManager
	streamID string
}

// NewStreamWrapper 创建流包装器
func NewStreamWrapper(stream interface{}, manager *StreamManager, streamID string) *StreamWrapper {
	return &StreamWrapper{
		stream:   stream,
		manager:  manager,
		streamID: streamID,
	}
}

// GetStream 获取底层流
func (sw *StreamWrapper) GetStream() interface{} {
	return sw.stream
}

// GetStreamID 获取流ID
func (sw *StreamWrapper) GetStreamID() string {
	return sw.streamID
}

// Close 关闭流
func (sw *StreamWrapper) Close() {
	sw.manager.UnregisterStream(sw.streamID)
}

// StreamMessage 流消息
type StreamMessage struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Time    time.Time   `json:"time"`
	Headers metadata.MD `json:"headers,omitempty"`
}

// StreamProcessor 流处理器
type StreamProcessor struct {
	manager *StreamManager
}

// NewStreamProcessor 创建流处理器
func NewStreamProcessor(manager *StreamManager) *StreamProcessor {
	return &StreamProcessor{
		manager: manager,
	}
}

// ProcessClientStream 处理客户端流
func (sp *StreamProcessor) ProcessClientStream(ctx context.Context, stream interface{}) error {
	// 这里需要根据具体的 gRPC 流类型来实现
	// 例如：pb.UserService_UploadFileClient
	return nil
}

// ProcessServerStream 处理服务器流
func (sp *StreamProcessor) ProcessServerStream(ctx context.Context, stream interface{}) error {
	// 这里需要根据具体的 gRPC 流类型来实现
	// 例如：pb.UserService_DownloadFileServer
	return nil
}

// ProcessBidirectionalStream 处理双向流
func (sp *StreamProcessor) ProcessBidirectionalStream(ctx context.Context, stream interface{}) error {
	// 这里需要根据具体的 gRPC 流类型来实现
	// 例如：pb.UserService_ChatServer
	return nil
}

// StreamTransformer 流转换器
type StreamTransformer struct {
	transformers map[string]TransformFunc
}

// TransformFunc 转换函数
type TransformFunc func(interface{}) (interface{}, error)

// NewStreamTransformer 创建流转换器
func NewStreamTransformer() *StreamTransformer {
	return &StreamTransformer{
		transformers: make(map[string]TransformFunc),
	}
}

// RegisterTransformer 注册转换器
func (st *StreamTransformer) RegisterTransformer(name string, transformer TransformFunc) {
	st.transformers[name] = transformer
}

// Transform 转换数据
func (st *StreamTransformer) Transform(name string, data interface{}) (interface{}, error) {
	if transformer, exists := st.transformers[name]; exists {
		return transformer(data)
	}
	return data, nil
}

// StreamFilter 流过滤器
type StreamFilter struct {
	filters map[string]FilterFunc
}

// FilterFunc 过滤函数
type FilterFunc func(interface{}) (bool, error)

// NewStreamFilter 创建流过滤器
func NewStreamFilter() *StreamFilter {
	return &StreamFilter{
		filters: make(map[string]FilterFunc),
	}
}

// RegisterFilter 注册过滤器
func (sf *StreamFilter) RegisterFilter(name string, filter FilterFunc) {
	sf.filters[name] = filter
}

// Filter 过滤数据
func (sf *StreamFilter) Filter(name string, data interface{}) (bool, error) {
	if filter, exists := sf.filters[name]; exists {
		return filter(data)
	}
	return true, nil
}

// StreamAggregator 流聚合器
type StreamAggregator struct {
	aggregators map[string]AggregateFunc
	buffers     map[string][]interface{}
	mutex       sync.RWMutex
}

// AggregateFunc 聚合函数
type AggregateFunc func([]interface{}) (interface{}, error)

// NewStreamAggregator 创建流聚合器
func NewStreamAggregator() *StreamAggregator {
	return &StreamAggregator{
		aggregators: make(map[string]AggregateFunc),
		buffers:     make(map[string][]interface{}),
	}
}

// RegisterAggregator 注册聚合器
func (sa *StreamAggregator) RegisterAggregator(name string, aggregator AggregateFunc) {
	sa.aggregators[name] = aggregator
}

// AddData 添加数据
func (sa *StreamAggregator) AddData(name string, data interface{}) {
	sa.mutex.Lock()
	defer sa.mutex.Unlock()

	if _, exists := sa.buffers[name]; !exists {
		sa.buffers[name] = make([]interface{}, 0)
	}
	sa.buffers[name] = append(sa.buffers[name], data)
}

// Aggregate 聚合数据
func (sa *StreamAggregator) Aggregate(name string) (interface{}, error) {
	sa.mutex.Lock()
	defer sa.mutex.Unlock()

	if aggregator, exists := sa.aggregators[name]; exists {
		if buffer, exists := sa.buffers[name]; exists {
			result, err := aggregator(buffer)
			if err == nil {
				// 清空缓冲区
				sa.buffers[name] = make([]interface{}, 0)
			}
			return result, err
		}
	}
	return nil, fmt.Errorf("aggregator not found: %s", name)
}

// StreamRateLimiter 流限流器
type StreamRateLimiter struct {
	limiters map[string]*TokenBucket
	mutex    sync.RWMutex
}

// TokenBucket 令牌桶
type TokenBucket struct {
	tokens     int
	capacity   int
	rate       float64
	lastRefill time.Time
	mutex      sync.Mutex
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(capacity int, rate float64) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		rate:       rate,
		lastRefill: time.Now(),
	}
}

// Take 获取令牌
func (tb *TokenBucket) Take() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	// 补充令牌
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tokensToAdd := int(elapsed * tb.rate)

	if tokensToAdd > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// NewStreamRateLimiter 创建流限流器
func NewStreamRateLimiter() *StreamRateLimiter {
	return &StreamRateLimiter{
		limiters: make(map[string]*TokenBucket),
	}
}

// AddLimiter 添加限流器
func (srl *StreamRateLimiter) AddLimiter(name string, capacity int, rate float64) {
	srl.mutex.Lock()
	defer srl.mutex.Unlock()
	srl.limiters[name] = NewTokenBucket(capacity, rate)
}

// Allow 检查是否允许
func (srl *StreamRateLimiter) Allow(name string) bool {
	srl.mutex.RLock()
	limiter, exists := srl.limiters[name]
	srl.mutex.RUnlock()

	if !exists {
		return true
	}

	return limiter.Take()
}

// StreamBuffer 流缓冲区
type StreamBuffer struct {
	buffers map[string]*RingBuffer
	mutex   sync.RWMutex
}

// RingBuffer 环形缓冲区
type RingBuffer struct {
	buffer []interface{}
	size   int
	head   int
	tail   int
	count  int
	mutex  sync.Mutex
}

// NewRingBuffer 创建环形缓冲区
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]interface{}, size),
		size:   size,
	}
}

// Push 推入数据
func (rb *RingBuffer) Push(data interface{}) bool {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()

	if rb.count >= rb.size {
		return false // 缓冲区已满
	}

	rb.buffer[rb.tail] = data
	rb.tail = (rb.tail + 1) % rb.size
	rb.count++
	return true
}

// Pop 弹出数据
func (rb *RingBuffer) Pop() (interface{}, bool) {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()

	if rb.count == 0 {
		return nil, false // 缓冲区为空
	}

	data := rb.buffer[rb.head]
	rb.head = (rb.head + 1) % rb.size
	rb.count--
	return data, true
}

// Size 获取大小
func (rb *RingBuffer) Size() int {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()
	return rb.count
}

// NewStreamBuffer 创建流缓冲区
func NewStreamBuffer() *StreamBuffer {
	return &StreamBuffer{
		buffers: make(map[string]*RingBuffer),
	}
}

// AddBuffer 添加缓冲区
func (sb *StreamBuffer) AddBuffer(name string, size int) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()
	sb.buffers[name] = NewRingBuffer(size)
}

// Push 推入数据
func (sb *StreamBuffer) Push(name string, data interface{}) bool {
	sb.mutex.RLock()
	buffer, exists := sb.buffers[name]
	sb.mutex.RUnlock()

	if !exists {
		return false
	}

	return buffer.Push(data)
}

// Pop 弹出数据
func (sb *StreamBuffer) Pop(name string) (interface{}, bool) {
	sb.mutex.RLock()
	buffer, exists := sb.buffers[name]
	sb.mutex.RUnlock()

	if !exists {
		return nil, false
	}

	return buffer.Pop()
}

// StreamPartitioner 流分区器
type StreamPartitioner struct {
	partitions map[string]PartitionFunc
}

// PartitionFunc 分区函数
type PartitionFunc func(interface{}) (string, error)

// NewStreamPartitioner 创建流分区器
func NewStreamPartitioner() *StreamPartitioner {
	return &StreamPartitioner{
		partitions: make(map[string]PartitionFunc),
	}
}

// RegisterPartition 注册分区函数
func (sp *StreamPartitioner) RegisterPartition(name string, partition PartitionFunc) {
	sp.partitions[name] = partition
}

// Partition 分区数据
func (sp *StreamPartitioner) Partition(name string, data interface{}) (string, error) {
	if partition, exists := sp.partitions[name]; exists {
		return partition(data)
	}
	return "default", nil
}

// StreamScheduler 流调度器
type StreamScheduler struct {
	schedulers map[string]SchedulerFunc
	queues     map[string]chan interface{}
	mutex      sync.RWMutex
}

// SchedulerFunc 调度函数
type SchedulerFunc func(interface{}) (int, error)

// NewStreamScheduler 创建流调度器
func NewStreamScheduler() *StreamScheduler {
	return &StreamScheduler{
		schedulers: make(map[string]SchedulerFunc),
		queues:     make(map[string]chan interface{}),
	}
}

// RegisterScheduler 注册调度器
func (ss *StreamScheduler) RegisterScheduler(name string, scheduler SchedulerFunc, queueSize int) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	ss.schedulers[name] = scheduler
	ss.queues[name] = make(chan interface{}, queueSize)
}

// Schedule 调度数据
func (ss *StreamScheduler) Schedule(name string, data interface{}) error {
	ss.mutex.RLock()
	scheduler, exists := ss.schedulers[name]
	queue, queueExists := ss.queues[name]
	ss.mutex.RUnlock()

	if !exists || !queueExists {
		return fmt.Errorf("scheduler not found: %s", name)
	}

	_, err := scheduler(data)
	if err != nil {
		return err
	}

	// 这里可以实现基于优先级的调度逻辑
	select {
	case queue <- data:
		return nil
	default:
		return fmt.Errorf("queue is full: %s", name)
	}
}

// GetQueue 获取队列
func (ss *StreamScheduler) GetQueue(name string) (<-chan interface{}, error) {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	if queue, exists := ss.queues[name]; exists {
		return queue, nil
	}

	return nil, fmt.Errorf("queue not found: %s", name)
}

// StreamValidator 流验证器
type StreamValidator struct {
	validators map[string]ValidatorFunc
}

// ValidatorFunc 验证函数
type ValidatorFunc func(interface{}) error

// NewStreamValidator 创建流验证器
func NewStreamValidator() *StreamValidator {
	return &StreamValidator{
		validators: make(map[string]ValidatorFunc),
	}
}

// RegisterValidator 注册验证器
func (sv *StreamValidator) RegisterValidator(name string, validator ValidatorFunc) {
	sv.validators[name] = validator
}

// Validate 验证数据
func (sv *StreamValidator) Validate(name string, data interface{}) error {
	if validator, exists := sv.validators[name]; exists {
		return validator(data)
	}
	return nil
}

// StreamEnricher 流丰富器
type StreamEnricher struct {
	enrichers map[string]EnrichFunc
}

// EnrichFunc 丰富函数
type EnrichFunc func(interface{}) (interface{}, error)

// NewStreamEnricher 创建流丰富器
func NewStreamEnricher() *StreamEnricher {
	return &StreamEnricher{
		enrichers: make(map[string]EnrichFunc),
	}
}

// RegisterEnricher 注册丰富器
func (se *StreamEnricher) RegisterEnricher(name string, enricher EnrichFunc) {
	se.enrichers[name] = enricher
}

// Enrich 丰富数据
func (se *StreamEnricher) Enrich(name string, data interface{}) (interface{}, error) {
	if enricher, exists := se.enrichers[name]; exists {
		return enricher(data)
	}
	return data, nil
}

// StreamRouter 流路由器
type StreamRouter struct {
	routes map[string]RouteFunc
}

// RouteFunc 路由函数
type RouteFunc func(interface{}) (string, error)

// NewStreamRouter 创建流路由器
func NewStreamRouter() *StreamRouter {
	return &StreamRouter{
		routes: make(map[string]RouteFunc),
	}
}

// RegisterRoute 注册路由
func (sr *StreamRouter) RegisterRoute(name string, route RouteFunc) {
	sr.routes[name] = route
}

// Route 路由数据
func (sr *StreamRouter) Route(name string, data interface{}) (string, error) {
	if route, exists := sr.routes[name]; exists {
		return route(data)
	}
	return "default", nil
}

// 辅助函数

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 内置转换器

// JSONTransformer JSON转换器
func JSONTransformer() TransformFunc {
	return func(data interface{}) (interface{}, error) {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		return string(jsonData), nil
	}
}

// Base64Transformer Base64转换器
func Base64Transformer() TransformFunc {
	return func(data interface{}) (interface{}, error) {
		switch v := data.(type) {
		case string:
			return base64.StdEncoding.EncodeToString([]byte(v)), nil
		case []byte:
			return base64.StdEncoding.EncodeToString(v), nil
		default:
			return nil, fmt.Errorf("unsupported type for base64 encoding")
		}
	}
}

// 内置过滤器

// SizeFilter 大小过滤器
func SizeFilter(maxSize int) FilterFunc {
	return func(data interface{}) (bool, error) {
		switch v := data.(type) {
		case string:
			return len(v) <= maxSize, nil
		case []byte:
			return len(v) <= maxSize, nil
		default:
			return true, nil
		}
	}
}

// TypeFilter 类型过滤器
func TypeFilter(allowedTypes ...string) FilterFunc {
	return func(data interface{}) (bool, error) {
		dataType := fmt.Sprintf("%T", data)
		for _, allowedType := range allowedTypes {
			if dataType == allowedType {
				return true, nil
			}
		}
		return false, nil
	}
}

// 内置聚合器

// CountAggregator 计数聚合器
func CountAggregator() AggregateFunc {
	return func(data []interface{}) (interface{}, error) {
		return len(data), nil
	}
}

// SumAggregator 求和聚合器
func SumAggregator() AggregateFunc {
	return func(data []interface{}) (interface{}, error) {
		var sum float64
		for _, item := range data {
			switch v := item.(type) {
			case int:
				sum += float64(v)
			case int64:
				sum += float64(v)
			case float64:
				sum += v
			case float32:
				sum += float64(v)
			default:
				return nil, fmt.Errorf("unsupported type for sum aggregation")
			}
		}
		return sum, nil
	}
}

// AverageAggregator 平均值聚合器
func AverageAggregator() AggregateFunc {
	return func(data []interface{}) (interface{}, error) {
		if len(data) == 0 {
			return 0.0, nil
		}

		sum, err := SumAggregator()(data)
		if err != nil {
			return nil, err
		}

		return sum.(float64) / float64(len(data)), nil
	}
}

// 内置验证器

// RequiredFieldValidator 必填字段验证器
func RequiredFieldValidator(fields ...string) ValidatorFunc {
	return func(data interface{}) error {
		// 这里需要根据具体的数据结构来实现
		// 例如：检查结构体字段是否为空
		return nil
	}
}

// RangeValidator 范围验证器
func RangeValidator(min, max float64) ValidatorFunc {
	return func(data interface{}) error {
		switch v := data.(type) {
		case int:
			val := float64(v)
			if val < min || val > max {
				return fmt.Errorf("value %f is out of range [%f, %f]", val, min, max)
			}
		case float64:
			if v < min || v > max {
				return fmt.Errorf("value %f is out of range [%f, %f]", v, min, max)
			}
		default:
			return fmt.Errorf("unsupported type for range validation")
		}
		return nil
	}
}

// StreamInterceptor 流拦截器
func StreamInterceptor(manager *StreamManager) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 生成流ID
		streamID := fmt.Sprintf("%s_%d", info.FullMethod, time.Now().UnixNano())

		// 注册流
		manager.RegisterStream(streamID, StreamTypeBidirectional, ss.Context())
		defer manager.UnregisterStream(streamID)

		// 添加流ID到上下文
		md := metadata.New(map[string]string{"stream-id": streamID})
		ctx := metadata.NewIncomingContext(ss.Context(), md)

		// 创建包装的流
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// 调用处理器
		return handler(srv, wrappedStream)
	}
}

// wrappedServerStream 包装的服务器流
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context 获取上下文
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// StreamMetrics 流指标
type StreamMetrics struct {
	ActiveStreams    int64
	TotalStreams     int64
	StreamDuration   time.Duration
	MessagesSent     int64
	MessagesReceived int64
	Errors           int64
}

// StreamMetricsCollector 流指标收集器
type StreamMetricsCollector struct {
	metrics map[string]*StreamMetrics
	mutex   sync.RWMutex
}

// NewStreamMetricsCollector 创建流指标收集器
func NewStreamMetricsCollector() *StreamMetricsCollector {
	return &StreamMetricsCollector{
		metrics: make(map[string]*StreamMetrics),
	}
}

// RecordStreamStart 记录流开始
func (smc *StreamMetricsCollector) RecordStreamStart(method string) {
	smc.mutex.Lock()
	defer smc.mutex.Unlock()

	if _, exists := smc.metrics[method]; !exists {
		smc.metrics[method] = &StreamMetrics{}
	}

	smc.metrics[method].ActiveStreams++
	smc.metrics[method].TotalStreams++
}

// RecordStreamEnd 记录流结束
func (smc *StreamMetricsCollector) RecordStreamEnd(method string, duration time.Duration) {
	smc.mutex.Lock()
	defer smc.mutex.Unlock()

	if metrics, exists := smc.metrics[method]; exists {
		metrics.ActiveStreams--
		metrics.StreamDuration = duration
	}
}

// RecordMessageSent 记录发送消息
func (smc *StreamMetricsCollector) RecordMessageSent(method string) {
	smc.mutex.Lock()
	defer smc.mutex.Unlock()

	if metrics, exists := smc.metrics[method]; exists {
		metrics.MessagesSent++
	}
}

// RecordMessageReceived 记录接收消息
func (smc *StreamMetricsCollector) RecordMessageReceived(method string) {
	smc.mutex.Lock()
	defer smc.mutex.Unlock()

	if metrics, exists := smc.metrics[method]; exists {
		metrics.MessagesReceived++
	}
}

// RecordError 记录错误
func (smc *StreamMetricsCollector) RecordError(method string) {
	smc.mutex.Lock()
	defer smc.mutex.Unlock()

	if metrics, exists := smc.metrics[method]; exists {
		metrics.Errors++
	}
}

// GetMetrics 获取指标
func (smc *StreamMetricsCollector) GetMetrics() map[string]*StreamMetrics {
	smc.mutex.RLock()
	defer smc.mutex.RUnlock()

	metrics := make(map[string]*StreamMetrics)
	for k, v := range smc.metrics {
		metrics[k] = v
	}
	return metrics
}

// StreamError 流错误
type StreamError struct {
	Code    codes.Code
	Message string
	Details []interface{}
}

// NewStreamError 创建流错误
func NewStreamError(code codes.Code, message string, details ...interface{}) *StreamError {
	return &StreamError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Error 实现 error 接口
func (se *StreamError) Error() string {
	return se.Message
}

// ToStatus 转换为 gRPC 状态
func (se *StreamError) ToStatus() *status.Status {
	st := status.New(se.Code, se.Message)
	// 注意：这里需要将 interface{} 转换为 proto.Message
	// 在实际使用中，Details 应该是 proto.Message 类型
	return st
}

// StreamContext 流上下文
type StreamContext struct {
	context.Context
	streamID  string
	method    string
	metadata  metadata.MD
	startTime time.Time
	metrics   *StreamMetricsCollector
}

// NewStreamContext 创建流上下文
func NewStreamContext(ctx context.Context, streamID, method string, md metadata.MD, metrics *StreamMetricsCollector) *StreamContext {
	return &StreamContext{
		Context:   ctx,
		streamID:  streamID,
		method:    method,
		metadata:  md,
		startTime: time.Now(),
		metrics:   metrics,
	}
}

// GetStreamID 获取流ID
func (sc *StreamContext) GetStreamID() string {
	return sc.streamID
}

// GetMethod 获取方法名
func (sc *StreamContext) GetMethod() string {
	return sc.method
}

// GetMetadata 获取元数据
func (sc *StreamContext) GetMetadata() metadata.MD {
	return sc.metadata
}

// GetStartTime 获取开始时间
func (sc *StreamContext) GetStartTime() time.Time {
	return sc.startTime
}

// GetDuration 获取持续时间
func (sc *StreamContext) GetDuration() time.Duration {
	return time.Since(sc.startTime)
}

// RecordMessageSent 记录发送消息
func (sc *StreamContext) RecordMessageSent() {
	if sc.metrics != nil {
		sc.metrics.RecordMessageSent(sc.method)
	}
}

// RecordMessageReceived 记录接收消息
func (sc *StreamContext) RecordMessageReceived() {
	if sc.metrics != nil {
		sc.metrics.RecordMessageReceived(sc.method)
	}
}

// RecordError 记录错误
func (sc *StreamContext) RecordError() {
	if sc.metrics != nil {
		sc.metrics.RecordError(sc.method)
	}
}

// StreamReader 流读取器
type StreamReader struct {
	stream interface{}
	ctx    *StreamContext
}

// NewStreamReader 创建流读取器
func NewStreamReader(stream interface{}, ctx *StreamContext) *StreamReader {
	return &StreamReader{
		stream: stream,
		ctx:    ctx,
	}
}

// Read 读取消息
func (sr *StreamReader) Read(msg interface{}) error {
	// 这里需要根据具体的 gRPC 流类型来实现
	// 例如：stream.RecvMsg(msg)
	return nil
}

// ReadAll 读取所有消息
func (sr *StreamReader) ReadAll() ([]interface{}, error) {
	var messages []interface{}

	for {
		var msg interface{}
		if err := sr.Read(&msg); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		messages = append(messages, msg)
		sr.ctx.RecordMessageReceived()
	}

	return messages, nil
}

// StreamWriter 流写入器
type StreamWriter struct {
	stream interface{}
	ctx    *StreamContext
}

// NewStreamWriter 创建流写入器
func NewStreamWriter(stream interface{}, ctx *StreamContext) *StreamWriter {
	return &StreamWriter{
		stream: stream,
		ctx:    ctx,
	}
}

// Write 写入消息
func (sw *StreamWriter) Write(msg interface{}) error {
	// 这里需要根据具体的 gRPC 流类型来实现
	// 例如：stream.SendMsg(msg)
	sw.ctx.RecordMessageSent()
	return nil
}

// WriteAll 写入所有消息
func (sw *StreamWriter) WriteAll(messages []interface{}) error {
	for _, msg := range messages {
		if err := sw.Write(msg); err != nil {
			return err
		}
	}
	return nil
}
