# 缓存系统指南

## 📖 概述

Laravel-Go Framework 提供了强大的缓存系统，支持多种缓存驱动、缓存标签、缓存策略和性能优化，帮助提升应用程序的响应速度和性能。

> 📚 **相关文档**: 如需查看详细的 API 接口说明，请参考 [缓存系统 API 参考](../api/cache.md)

## 🚀 快速开始

### 1. 基本使用

```go
// 设置缓存
cache.Put("user:1", user, time.Hour)

// 获取缓存
if user, exists := cache.Get("user:1"); exists {
    return user.(*User)
}

// 删除缓存
cache.Forget("user:1")

// 检查缓存是否存在
if cache.Has("user:1") {
    // 缓存存在
}

// 获取或设置缓存
user := cache.Remember("user:1", time.Hour, func() interface{} {
    return db.First(&User{}, 1)
})
```

### 2. 缓存标签

```go
// 使用标签
cache.Tags("users", "profiles").Put("user:1", user, time.Hour)

// 通过标签清除缓存
cache.Tags("users").Flush()

// 多标签操作
cache.Tags("users", "posts").Put("user:1:posts", posts, time.Hour)
```

## 🔧 缓存驱动

### 1. Redis 驱动

```go
// 配置 Redis
config.Set("cache.driver", "redis")
config.Set("cache.redis.host", "localhost")
config.Set("cache.redis.port", 6379)
config.Set("cache.redis.password", "")
config.Set("cache.redis.database", 0)

// 使用 Redis 缓存
redisCache := cache.NewRedisDriver(config.Get("cache.redis"))
redisCache.Put("key", "value", time.Hour)
```

### 2. 内存驱动

```go
// 配置内存缓存
config.Set("cache.driver", "memory")

// 使用内存缓存
memoryCache := cache.NewMemoryDriver()
memoryCache.Put("key", "value", time.Hour)
```

### 3. 文件驱动

```go
// 配置文件缓存
config.Set("cache.driver", "file")
config.Set("cache.file.path", "storage/cache")

// 使用文件缓存
fileCache := cache.NewFileDriver(config.Get("cache.file.path"))
fileCache.Put("key", "value", time.Hour)
```

### 4. 数据库驱动

```go
// 配置数据库缓存
config.Set("cache.driver", "database")
config.Set("cache.database.table", "cache")

// 使用数据库缓存
dbCache := cache.NewDatabaseDriver(db, "cache")
dbCache.Put("key", "value", time.Hour)
```

### 5. MongoDB 驱动

```go
// 配置 MongoDB 缓存
config.Set("cache.driver", "mongodb")
config.Set("cache.mongodb.uri", "mongodb://localhost:27017")
config.Set("cache.mongodb.database", "laravel_go")
config.Set("cache.mongodb.collection", "cache")

// 使用 MongoDB 缓存
mongoClient, _ := mongo.Connect(context.Background(), "mongodb://localhost:27017")
mongoCache := cache.NewMongoStore(mongoClient, "laravel_go", "cache")
cache.Cache.Extend("mongodb", mongoCache)
```

### 6. Memcached 驱动

```go
// 配置 Memcached 缓存
config.Set("cache.driver", "memcached")
config.Set("cache.memcached.host", "127.0.0.1")
config.Set("cache.memcached.port", 11211)

// 简单使用
memcachedCache := cache.NewMemcachedStore("127.0.0.1:11211")
cache.Cache.Extend("memcached", memcachedCache)

// 使用配置
config := map[string]interface{}{
    "host": "127.0.0.1",
    "port": "11211",
}
memcachedCache := cache.NewMemcachedStoreWithConfig(config)
cache.Cache.Extend("memcached", memcachedCache)
```

## 📊 缓存策略

### 1. 查询缓存

```go
// 缓存查询结果
func (s *UserService) GetUser(id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 尝试从缓存获取
    if cached, exists := cache.Get(cacheKey); exists {
        return cached.(*User), nil
    }

    // 查询数据库
    var user User
    err := s.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }

    // 缓存结果
    cache.Put(cacheKey, &user, time.Hour)

    return &user, nil
}

// 使用 Remember 方法简化
func (s *UserService) GetUser(id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    user := cache.Remember(cacheKey, time.Hour, func() interface{} {
        var user User
        s.db.First(&user, id)
        return &user
    })

    return user.(*User), nil
}
```

### 2. 列表缓存

```go
// 缓存分页列表
func (s *PostService) GetPosts(page, limit int, filters map[string]interface{}) ([]*Post, int64) {
    // 生成缓存键
    cacheKey := fmt.Sprintf("posts:page:%d:limit:%d:filters:%v", page, limit, filters)

    result := cache.Remember(cacheKey, time.Minute*30, func() interface{} {
        query := s.db.Model(&Post{})

        // 应用过滤条件
        for key, value := range filters {
            query = query.Where(key+" = ?", value)
        }

        var total int64
        query.Count(&total)

        var posts []*Post
        offset := (page - 1) * limit
        query.Offset(offset).Limit(limit).Find(&posts)

        return map[string]interface{}{
            "posts": posts,
            "total": total,
        }
    })

    data := result.(map[string]interface{})
    return data["posts"].([]*Post), data["total"].(int64)
}
```

### 3. 关联数据缓存

```go
// 缓存用户及其关联数据
func (s *UserService) GetUserWithRelations(id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d:with_relations", id)

    user := cache.Remember(cacheKey, time.Hour, func() interface{} {
        var user User
        s.db.Preload("Posts").Preload("Comments").First(&user, id)
        return &user
    })

    return user.(*User), nil
}
```

## 🏷️ 缓存标签

### 1. 标签管理

```go
// 使用标签组织缓存
func (s *UserService) GetUser(id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    user := cache.Tags("users").Remember(cacheKey, time.Hour, func() interface{} {
        var user User
        s.db.First(&user, id)
        return &user
    })

    return user.(*User), nil
}

// 清除用户相关缓存
func (s *UserService) ClearUserCache(id uint) {
    cache.Tags("users").Flush()
}

// 清除特定用户缓存
func (s *UserService) ClearUserCache(id uint) {
    cache.Tags("users").Forget(fmt.Sprintf("user:%d", id))
}
```

### 2. 多标签策略

```go
// 使用多个标签
func (s *PostService) GetPost(id uint) (*Post, error) {
    cacheKey := fmt.Sprintf("post:%d", id)

    post := cache.Tags("posts", "users", "categories").Remember(cacheKey, time.Hour, func() interface{} {
        var post Post
        s.db.Preload("User").Preload("Category").First(&post, id)
        return &post
    })

    return post.(*Post), nil
}

// 更新文章时清除相关缓存
func (s *PostService) UpdatePost(id uint, data map[string]interface{}) error {
    // 更新数据库
    err := s.db.Model(&Post{}).Where("id = ?", id).Updates(data).Error
    if err != nil {
        return err
    }

    // 清除相关缓存
    cache.Tags("posts").Forget(fmt.Sprintf("post:%d", id))
    cache.Tags("users").Flush() // 清除用户相关缓存
    cache.Tags("categories").Flush() // 清除分类相关缓存

    return nil
}
```

## ⚡ 性能优化

### 1. 缓存预热

```go
// 应用启动时预热缓存
func (s *CacheService) WarmUp() {
    // 预热热门文章
    s.warmUpPopularPosts()

    // 预热用户统计
    s.warmUpUserStats()

    // 预热系统配置
    s.warmUpSystemConfig()
}

func (s *CacheService) warmUpPopularPosts() {
    var posts []Post
    s.db.Where("status = ?", "published").
        Order("view_count desc").
        Limit(100).
        Find(&posts)

    for _, post := range posts {
        cacheKey := fmt.Sprintf("post:%d", post.ID)
        cache.Put(cacheKey, &post, time.Hour)
    }
}
```

### 2. 缓存穿透防护

```go
// 使用布隆过滤器防止缓存穿透
type BloomFilter struct {
    bitset []bool
    size   int
    hashCount int
}

func (bf *BloomFilter) Add(key string) {
    for i := 0; i < bf.hashCount; i++ {
        hash := bf.hash(key, i)
        bf.bitset[hash%bf.size] = true
    }
}

func (bf *BloomFilter) Contains(key string) bool {
    for i := 0; i < bf.hashCount; i++ {
        hash := bf.hash(key, i)
        if !bf.bitset[hash%bf.size] {
            return false
        }
    }
    return true
}

// 在缓存中使用
func (s *UserService) GetUser(id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 检查布隆过滤器
    if !bloomFilter.Contains(cacheKey) {
        return nil, errors.New("user not found")
    }

    user := cache.Remember(cacheKey, time.Hour, func() interface{} {
        var user User
        s.db.First(&user, id)
        return &user
    })

    return user.(*User), nil
}
```

### 3. 缓存雪崩防护

```go
// 使用随机过期时间防止缓存雪崩
func (s *UserService) GetUser(id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 随机过期时间（基础时间 ± 10%）
    baseExpiration := time.Hour
    randomOffset := time.Duration(rand.Intn(12)) * time.Minute // ±10分钟
    expiration := baseExpiration + randomOffset

    user := cache.Remember(cacheKey, expiration, func() interface{} {
        var user User
        s.db.First(&user, id)
        return &user
    })

    return user.(*User), nil
}
```

## 🔄 缓存更新策略

### 1. 写入时更新

```go
// 更新数据时同步更新缓存
func (s *UserService) UpdateUser(id uint, data map[string]interface{}) error {
    // 更新数据库
    err := s.db.Model(&User{}).Where("id = ?", id).Updates(data).Error
    if err != nil {
        return err
    }

    // 更新缓存
    cacheKey := fmt.Sprintf("user:%d", id)
    if user, exists := cache.Get(cacheKey); exists {
        // 更新缓存中的用户数据
        cachedUser := user.(*User)
        for key, value := range data {
            // 使用反射更新字段
            reflect.ValueOf(cachedUser).Elem().FieldByName(key).Set(reflect.ValueOf(value))
        }
        cache.Put(cacheKey, cachedUser, time.Hour)
    }

    return nil
}
```

### 2. 写入时删除

```go
// 更新数据时删除缓存（Cache Aside 模式）
func (s *UserService) UpdateUser(id uint, data map[string]interface{}) error {
    // 更新数据库
    err := s.db.Model(&User{}).Where("id = ?", id).Updates(data).Error
    if err != nil {
        return err
    }

    // 删除缓存，下次读取时重新加载
    cacheKey := fmt.Sprintf("user:%d", id)
    cache.Forget(cacheKey)

    return nil
}
```

### 3. 延迟双删

```go
// 延迟双删策略
func (s *UserService) UpdateUser(id uint, data map[string]interface{}) error {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 第一次删除缓存
    cache.Forget(cacheKey)

    // 更新数据库
    err := s.db.Model(&User{}).Where("id = ?", id).Updates(data).Error
    if err != nil {
        return err
    }

    // 延迟删除缓存（防止并发问题）
    go func() {
        time.Sleep(500 * time.Millisecond)
        cache.Forget(cacheKey)
    }()

    return nil
}
```

## 📈 缓存监控

### 1. 缓存统计

```go
// 缓存统计信息
type CacheStats struct {
    Hits   int64 `json:"hits"`
    Misses int64 `json:"misses"`
    Keys   int64 `json:"keys"`
    Memory int64 `json:"memory"`
}

// 缓存监控中间件
type CacheMonitorMiddleware struct {
    http.Middleware
    stats *CacheStats
}

func (m *CacheMonitorMiddleware) Handle(request http.Request, next http.HandlerFunc) http.Response {
    start := time.Now()

    response := next(request)

    // 记录缓存统计
    duration := time.Since(start)
    if duration > time.Millisecond*100 {
        // 记录慢查询
        log.Printf("Slow cache operation: %v", duration)
    }

    return response
}
```

### 2. 缓存健康检查

```go
// 缓存健康检查
func (s *CacheService) HealthCheck() error {
    testKey := "health_check"
    testValue := "ok"

    // 测试写入
    err := cache.Put(testKey, testValue, time.Minute)
    if err != nil {
        return fmt.Errorf("cache write failed: %v", err)
    }

    // 测试读取
    if value, exists := cache.Get(testKey); !exists || value != testValue {
        return fmt.Errorf("cache read failed")
    }

    // 测试删除
    cache.Forget(testKey)
    if cache.Has(testKey) {
        return fmt.Errorf("cache delete failed")
    }

    return nil
}
```

## 🛠️ 高级功能

### 1. 缓存锁

```go
// 分布式锁
type CacheLock struct {
    cache cache.Cache
    key   string
    ttl   time.Duration
}

func (l *CacheLock) Acquire() bool {
    return l.cache.Add(l.key+":lock", "locked", l.ttl)
}

func (l *CacheLock) Release() {
    l.cache.Forget(l.key + ":lock")
}

// 使用锁防止缓存击穿
func (s *UserService) GetUser(id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    lock := &CacheLock{cache: cache, key: cacheKey, ttl: time.Second * 10}

    // 尝试获取缓存
    if user, exists := cache.Get(cacheKey); exists {
        return user.(*User), nil
    }

    // 获取锁
    if !lock.Acquire() {
        // 等待其他进程加载数据
        time.Sleep(time.Millisecond * 100)
        if user, exists := cache.Get(cacheKey); exists {
            return user.(*User), nil
        }
    }

    defer lock.Release()

    // 加载数据
    var user User
    err := s.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }

    // 缓存数据
    cache.Put(cacheKey, &user, time.Hour)

    return &user, nil
}
```

### 2. 缓存版本控制

```go
// 缓存版本控制
type CacheVersion struct {
    cache cache.Cache
}

func (cv *CacheVersion) GetVersion(key string) int64 {
    versionKey := key + ":version"
    if version, exists := cv.cache.Get(versionKey); exists {
        return version.(int64)
    }
    return 0
}

func (cv *CacheVersion) IncrementVersion(key string) {
    versionKey := key + ":version"
    currentVersion := cv.GetVersion(key)
    cv.cache.Put(versionKey, currentVersion+1, time.Hour*24)
}

// 使用版本控制
func (s *UserService) GetUser(id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    version := cacheVersion.GetVersion(cacheKey)
    versionedKey := fmt.Sprintf("%s:v%d", cacheKey, version)

    user := cache.Remember(versionedKey, time.Hour, func() interface{} {
        var user User
        s.db.First(&user, id)
        return &user
    })

    return user.(*User), nil
}

func (s *UserService) UpdateUser(id uint, data map[string]interface{}) error {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 更新数据库
    err := s.db.Model(&User{}).Where("id = ?", id).Updates(data).Error
    if err != nil {
        return err
    }

    // 增加版本号，使旧缓存失效
    cacheVersion.IncrementVersion(cacheKey)

    return nil
}
```

## 📚 总结

Laravel-Go Framework 的缓存系统提供了：

1. **多种驱动**: Redis、内存、文件、数据库
2. **标签管理**: 灵活的缓存组织和清理
3. **性能优化**: 预热、穿透防护、雪崩防护
4. **更新策略**: 写入更新、写入删除、延迟双删
5. **监控功能**: 统计信息、健康检查
6. **高级功能**: 分布式锁、版本控制

通过合理使用缓存系统，可以显著提升应用性能和用户体验。
