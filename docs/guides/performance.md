# 性能优化指南

## 📖 概述

Laravel-Go Framework 提供了全面的性能优化功能，包括缓存优化、数据库优化、内存管理、并发控制、监控分析等，帮助构建高性能的应用程序。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [性能系统 API 参考](../api/performance.md)

## 🚀 快速开始

### 1. 性能监控

```go
// 创建性能监控器
monitor := performance.NewMonitor()

// 监控 HTTP 请求
type PerformanceMiddleware struct {
    monitor *performance.Monitor
}

func (m *PerformanceMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    start := time.Now()

    // 记录请求开始
    m.monitor.RecordRequestStart(request.Path, request.Method)

    // 处理请求
    response := next(request)

    // 记录请求完成
    duration := time.Since(start)
    m.monitor.RecordRequestEnd(request.Path, request.Method, duration, response.StatusCode)

    return response
}

// 监控数据库查询
func (s *UserService) GetUsers() ([]*Models.User, error) {
    start := time.Now()

    var users []*Models.User
    err := s.db.Find(&users).Error

    duration := time.Since(start)
    s.monitor.RecordDatabaseQuery("SELECT * FROM users", duration, len(users), err)

    return users, err
}
```

### 2. 缓存优化

```go
// 查询缓存
func (s *UserService) GetUserWithCache(id uint) (*Models.User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 尝试从缓存获取
    if cached, exists := cache.Get(cacheKey); exists {
        return cached.(*Models.User), nil
    }

    // 从数据库获取
    var user Models.User
    err := s.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }

    // 缓存结果
    cache.Set(cacheKey, &user, time.Hour)

    return &user, nil
}

// 批量缓存
func (s *UserService) GetUsersBatch(ids []uint) ([]*Models.User, error) {
    users := make([]*Models.User, 0)
    missingIDs := make([]uint, 0)

    // 批量从缓存获取
    for _, id := range ids {
        cacheKey := fmt.Sprintf("user:%d", id)
        if cached, exists := cache.Get(cacheKey); exists {
            users = append(users, cached.(*Models.User))
        } else {
            missingIDs = append(missingIDs, id)
        }
    }

    // 从数据库获取缺失的用户
    if len(missingIDs) > 0 {
        var missingUsers []*Models.User
        err := s.db.Where("id IN ?", missingIDs).Find(&missingUsers).Error
        if err != nil {
            return nil, err
        }

        // 缓存新获取的用户
        for _, user := range missingUsers {
            cacheKey := fmt.Sprintf("user:%d", user.ID)
            cache.Set(cacheKey, user, time.Hour)
            users = append(users, user)
        }
    }

    return users, nil
}

// 缓存预热
func (s *UserService) WarmupCache() error {
    // 获取热门用户
    var popularUsers []*Models.User
    err := s.db.Order("login_count DESC").Limit(100).Find(&popularUsers).Error
    if err != nil {
        return err
    }

    // 预热缓存
    for _, user := range popularUsers {
        cacheKey := fmt.Sprintf("user:%d", user.ID)
        cache.Set(cacheKey, user, time.Hour)
    }

    return nil
}
```

### 3. 数据库优化

```go
// 查询优化
func (s *UserService) GetUsersOptimized(page, limit int) ([]*Models.User, int64, error) {
    var total int64
    var users []*Models.User

    // 使用索引优化查询
    query := s.db.Model(&Models.User{}).
        Select("id, name, email, created_at").
        Where("status = ?", "active").
        Order("created_at DESC")

    // 获取总数
    err := query.Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    // 分页查询
    err = query.Offset((page - 1) * limit).Limit(limit).Find(&users).Error
    if err != nil {
        return nil, 0, err
    }

    return users, total, nil
}

// 批量操作
func (s *UserService) BatchCreateUsers(users []*Models.User) error {
    // 使用事务批量插入
    return s.db.Transaction(func(tx *database.Connection) error {
        batchSize := 100
        for i := 0; i < len(users); i += batchSize {
            end := i + batchSize
            if end > len(users) {
                end = len(users)
            }

            batch := users[i:end]
            if err := tx.CreateInBatches(batch, len(batch)).Error; err != nil {
                return err
            }
        }
        return nil
    })
}

// 连接池优化
func ConfigureDatabasePool() {
    config.Set("database.connections.mysql.pool.max_open_conns", 100)
    config.Set("database.connections.mysql.pool.max_idle_conns", 10)
    config.Set("database.connections.mysql.pool.conn_max_lifetime", time.Hour)
    config.Set("database.connections.mysql.pool.conn_max_idle_time", time.Minute*30)
}
```

### 4. 内存优化

```go
// 内存池
type UserPool struct {
    pool sync.Pool
}

func NewUserPool() *UserPool {
    return &UserPool{
        pool: sync.Pool{
            New: func() interface{} {
                return &Models.User{}
            },
        },
    }
}

func (p *UserPool) Get() *Models.User {
    return p.pool.Get().(*Models.User)
}

func (p *UserPool) Put(user *Models.User) {
    // 重置用户对象
    user.ID = 0
    user.Name = ""
    user.Email = ""
    user.Password = ""
    user.CreatedAt = time.Time{}
    user.UpdatedAt = time.Time{}

    p.pool.Put(user)
}

// 使用内存池
func (s *UserService) ProcessUsers(ids []uint) error {
    pool := NewUserPool()

    for _, id := range ids {
        user := pool.Get()
        defer pool.Put(user)

        // 处理用户
        if err := s.processUser(user, id); err != nil {
            return err
        }
    }

    return nil
}

// 对象复用
type BufferPool struct {
    pool sync.Pool
}

func NewBufferPool() *BufferPool {
    return &BufferPool{
        pool: sync.Pool{
            New: func() interface{} {
                return make([]byte, 0, 1024)
            },
        },
    }
}

func (p *BufferPool) Get() []byte {
    return p.pool.Get().([]byte)
}

func (p *BufferPool) Put(buf []byte) {
    // 重置缓冲区
    buf = buf[:0]
    p.pool.Put(buf)
}
```

### 5. 并发控制

```go
// 并发限制
type ConcurrencyLimiter struct {
    semaphore chan struct{}
}

func NewConcurrencyLimiter(maxConcurrency int) *ConcurrencyLimiter {
    return &ConcurrencyLimiter{
        semaphore: make(chan struct{}, maxConcurrency),
    }
}

func (l *ConcurrencyLimiter) Execute(task func() error) error {
    l.semaphore <- struct{}{}
    defer func() { <-l.semaphore }()

    return task()
}

// 使用并发限制
func (s *UserService) ProcessUsersConcurrently(users []*Models.User) error {
    limiter := NewConcurrencyLimiter(10) // 最多10个并发

    var wg sync.WaitGroup
    errors := make(chan error, len(users))

    for _, user := range users {
        wg.Add(1)
        go func(u *Models.User) {
            defer wg.Done()

            err := limiter.Execute(func() error {
                return s.processUser(u)
            })

            if err != nil {
                errors <- err
            }
        }(user)
    }

    wg.Wait()
    close(errors)

    // 检查错误
    for err := range errors {
        if err != nil {
            return err
        }
    }

    return nil
}

// 工作池
type WorkerPool struct {
    workers    int
    jobQueue   chan func()
    workerPool chan chan func()
    quit       chan bool
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers:    workers,
        jobQueue:   make(chan func(), 1000),
        workerPool: make(chan chan func(), workers),
        quit:       make(chan bool),
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        worker := NewWorker(p.workerPool)
        worker.Start()
    }

    go p.dispatch()
}

func (p *WorkerPool) dispatch() {
    for {
        select {
        case job := <-p.jobQueue:
            worker := <-p.workerPool
            worker <- job
        case <-p.quit:
            return
        }
    }
}

func (p *WorkerPool) Submit(job func()) {
    p.jobQueue <- job
}

func (p *WorkerPool) Stop() {
    p.quit <- true
}
```

### 6. 异步处理

```go
// 异步任务处理
func (s *UserService) ProcessUserAsync(userID uint) {
    // 提交到队列异步处理
    queue.Push(&ProcessUserJob{
        UserID: userID,
        Time:   time.Now(),
    })
}

// 异步任务
type ProcessUserJob struct {
    UserID uint      `json:"user_id"`
    Time   time.Time `json:"time"`
}

func (j *ProcessUserJob) Handle() error {
    // 异步处理用户数据
    userService := Services.NewUserService(database.NewConnection())

    user, err := userService.GetUser(j.UserID)
    if err != nil {
        return err
    }

    // 执行耗时操作
    return userService.ProcessUserData(user)
}

// 异步邮件发送
func (s *EmailService) SendWelcomeEmailAsync(user *Models.User) {
    queue.Push(&SendWelcomeEmailJob{
        UserID:   user.ID,
        Email:    user.Email,
        Name:     user.Name,
    })
}

type SendWelcomeEmailJob struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
}

func (j *SendWelcomeEmailJob) Handle() error {
    // 异步发送邮件
    emailService := Services.NewEmailService()

    return emailService.SendWelcomeEmail(j.Email, j.Name)
}
```

### 7. 资源管理

```go
// 资源清理
type ResourceManager struct {
    resources []io.Closer
    mu        sync.Mutex
}

func (rm *ResourceManager) AddResource(resource io.Closer) {
    rm.mu.Lock()
    defer rm.mu.Unlock()

    rm.resources = append(rm.resources, resource)
}

func (rm *ResourceManager) Cleanup() error {
    rm.mu.Lock()
    defer rm.mu.Unlock()

    var errors []error
    for _, resource := range rm.resources {
        if err := resource.Close(); err != nil {
            errors = append(errors, err)
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("cleanup errors: %v", errors)
    }

    return nil
}

// 使用资源管理器
func (s *UserService) ProcessLargeFile(filepath string) error {
    rm := &ResourceManager{}
    defer rm.Cleanup()

    file, err := os.Open(filepath)
    if err != nil {
        return err
    }
    rm.AddResource(file)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        // 处理每一行
        if err := s.processLine(line); err != nil {
            return err
        }
    }

    return scanner.Err()
}
```

### 8. 性能分析

```go
// 性能分析器
type PerformanceProfiler struct {
    metrics map[string]*Metric
    mu      sync.RWMutex
}

type Metric struct {
    Count   int64         `json:"count"`
    Total   time.Duration `json:"total"`
    Average time.Duration `json:"average"`
    Min     time.Duration `json:"min"`
    Max     time.Duration `json:"max"`
}

func (p *PerformanceProfiler) Record(name string, duration time.Duration) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.metrics[name] == nil {
        p.metrics[name] = &Metric{}
    }

    metric := p.metrics[name]
    metric.Count++
    metric.Total += duration

    if metric.Min == 0 || duration < metric.Min {
        metric.Min = duration
    }

    if duration > metric.Max {
        metric.Max = duration
    }

    metric.Average = metric.Total / time.Duration(metric.Count)
}

func (p *PerformanceProfiler) GetMetrics() map[string]*Metric {
    p.mu.RLock()
    defer p.mu.RUnlock()

    metrics := make(map[string]*Metric)
    for k, v := range p.metrics {
        metrics[k] = v
    }

    return metrics
}

// 使用性能分析器
func (s *UserService) GetUsersWithProfiling() ([]*Models.User, error) {
    profiler := performance.GetProfiler()

    start := time.Now()
    defer func() {
        profiler.Record("get_users", time.Since(start))
    }()

    var users []*Models.User
    err := s.db.Find(&users).Error

    return users, err
}
```

### 9. 负载均衡

```go
// 负载均衡器
type LoadBalancer struct {
    servers []*Server
    current int
    mu      sync.Mutex
}

type Server struct {
    Address string
    Weight  int
    Healthy bool
}

func (lb *LoadBalancer) Next() *Server {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    if len(lb.servers) == 0 {
        return nil
    }

    // 轮询算法
    server := lb.servers[lb.current]
    lb.current = (lb.current + 1) % len(lb.servers)

    return server
}

func (lb *LoadBalancer) AddServer(server *Server) {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    lb.servers = append(lb.servers, server)
}

// 使用负载均衡器
func (s *UserService) GetUserFromBalancedServer(userID uint) (*Models.User, error) {
    lb := s.loadBalancer

    for i := 0; i < len(lb.servers); i++ {
        server := lb.Next()
        if server == nil || !server.Healthy {
            continue
        }

        // 尝试从服务器获取用户
        user, err := s.getUserFromServer(server, userID)
        if err == nil {
            return user, nil
        }
    }

    return nil, errors.New("no available server")
}
```

### 10. 性能测试

```go
// 性能测试
func BenchmarkUserService_GetUser(b *testing.B) {
    service := setupTestUserService()

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _, err := service.GetUser(uint(i % 100))
        if err != nil {
            b.Fatalf("Failed to get user: %v", err)
        }
    }
}

func BenchmarkUserService_GetUserWithCache(b *testing.B) {
    service := setupTestUserService()

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _, err := service.GetUserWithCache(uint(i % 100))
        if err != nil {
            b.Fatalf("Failed to get user: %v", err)
        }
    }
}

// 压力测试
func TestUserService_ConcurrentAccess(t *testing.T) {
    service := setupTestUserService()

    var wg sync.WaitGroup
    concurrency := 100
    requests := 1000

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            for j := 0; j < requests/concurrency; j++ {
                _, err := service.GetUser(uint(j % 100))
                if err != nil {
                    t.Errorf("Failed to get user: %v", err)
                }
            }
        }()
    }

    wg.Wait()
}
```

## 📚 总结

Laravel-Go Framework 的性能优化系统提供了：

1. **性能监控**: 请求、数据库、缓存性能监控
2. **缓存优化**: 查询缓存、批量缓存、缓存预热
3. **数据库优化**: 查询优化、批量操作、连接池
4. **内存优化**: 内存池、对象复用
5. **并发控制**: 并发限制、工作池
6. **异步处理**: 队列任务、异步操作
7. **资源管理**: 资源清理、内存管理
8. **性能分析**: 性能指标收集和分析
9. **负载均衡**: 服务器负载均衡
10. **性能测试**: 基准测试、压力测试

通过合理使用这些性能优化功能，可以构建高性能、可扩展的应用程序。
