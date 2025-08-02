# Laravel-Go 缓存系统

Laravel-Go 缓存系统提供了统一的缓存接口和多驱动支持，支持多种缓存后端。

## 支持的驱动

### 1. 内存驱动 (MemoryStore) ✅ 已实现

- **特点**: 高性能内存缓存，支持自动过期清理
- **适用场景**: 单机应用，临时数据缓存
- **优势**: 速度最快，无需外部依赖

```go
// 使用内存驱动
memoryStore := cache.NewMemoryStore()
cache.Cache.Extend("memory", memoryStore)
```

### 2. 文件驱动 (FileStore) ✅ 已实现

- **特点**: 基于文件的持久化缓存
- **适用场景**: 单机应用，需要持久化的缓存
- **优势**: 数据持久化，无需数据库

```go
// 使用文件驱动
fileStore := cache.NewFileStore("./storage/cache")
cache.Cache.Extend("file", fileStore)
```

### 3. 数据库驱动 (DatabaseStore) ✅ 已实现

- **特点**: 基于关系型数据库的缓存
- **支持数据库**: MySQL、PostgreSQL、SQLite
- **适用场景**: 多实例应用，需要共享缓存

```go
// 使用数据库驱动
db, _ := sql.Open("mysql", "dsn")
dbStore := cache.NewDatabaseStore(db, "cache_table")
cache.Cache.Extend("database", dbStore)
```

### 4. Redis 驱动 (RedisStore) ✅ 已实现

- **特点**: 高性能分布式缓存
- **适用场景**: 分布式应用，高并发场景
- **优势**: 高性能，支持集群，丰富的数据结构

```go
// 使用Redis驱动（需要安装依赖）
// go get github.com/redis/go-redis/v9
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})
redisStore := cache.NewRedisStore(redisClient)
cache.Cache.Extend("redis", redisStore)
```

### 5. MongoDB 驱动 (MongoStore) ✅ 已实现

- **特点**: 基于 MongoDB 的文档型缓存
- **适用场景**: 需要复杂数据结构缓存的场景
- **优势**: 支持复杂查询，自动 TTL 索引

```go
// 使用MongoDB驱动（需要安装依赖）
// go get go.mongodb.org/mongo-driver/mongo
mongoClient, _ := mongo.Connect(context.Background(), "mongodb://localhost:27017")
mongoStore := cache.NewMongoStore(mongoClient, "cache_db", "cache_collection")
cache.Cache.Extend("mongodb", mongoStore)
```

### 6. Memcached 驱动 (MemcachedStore) ✅ 已实现

- **特点**: 高性能分布式内存缓存
- **适用场景**: 高并发、高吞吐量的缓存场景
- **优势**: 高性能、低延迟、支持集群

```go
// 使用Memcached驱动（需要安装依赖）
// go get github.com/bradfitz/gomemcache/memcache

// 简单使用
memcachedStore := cache.NewMemcachedStore("127.0.0.1:11211")
cache.Cache.Extend("memcached", memcachedStore)

// 使用配置
config := map[string]interface{}{
    "host": "127.0.0.1",
    "port": "11211",
}
memcachedStore := cache.NewMemcachedStoreWithConfig(config)
cache.Cache.Extend("memcached", memcachedStore)
```

## 核心接口

### Store 接口

所有缓存驱动都必须实现 `Store` 接口：

```go
type Store interface {
    // 基本操作
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    Clear() error

    // 类型化操作
    GetString(key string) (string, error)
    GetInt(key string) (int, error)
    GetFloat(key string) (float64, error)
    GetBool(key string) (bool, error)
    GetBytes(key string) ([]byte, error)

    SetString(key string, value string, ttl time.Duration) error
    SetInt(key string, value int, ttl time.Duration) error
    SetFloat(key string, value float64, ttl time.Duration) error
    SetBool(key string, value bool, ttl time.Duration) error
    SetBytes(key string, value []byte, ttl time.Duration) error

    // 批量操作
    DeleteMultiple(keys []string) error

    // 检查操作
    Has(key string) bool
    Missing(key string) bool

    // 数值操作
    Increment(key string, value int) (int, error)
    Decrement(key string, value int) (int, error)

    // 高级操作
    Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error)
    RememberForever(key string, callback func() (interface{}, error)) (interface{}, error)

    // 标签支持
    Tags(names ...string) TaggedStore

    // 前缀管理
    GetPrefix() string
    SetPrefix(prefix string)

    // 刷新
    Flush() error
}
```

## 扩展自定义驱动

### 1. 实现 Store 接口

```go
type CustomStore struct {
    // 你的存储实现
}

func (s *CustomStore) Get(key string) (interface{}, error) {
    // 实现获取逻辑
}

func (s *CustomStore) Set(key string, value interface{}, ttl time.Duration) error {
    // 实现设置逻辑
}

// ... 实现其他接口方法
```

### 2. 注册驱动

```go
// 创建自定义驱动
customStore := NewCustomStore()

// 注册到缓存管理器
cache.Cache.Extend("custom", customStore)

// 设置为默认驱动
cache.Cache.SetDefaultStore("custom")
```

### 3. 使用驱动

```go
// 使用特定驱动
customCache := cache.Cache.Store("custom")
customCache.Set("key", "value", time.Hour)

// 使用默认驱动
cache.Set("key", "value", time.Hour)
```

## 高级功能

### 缓存标签

支持缓存标签，可以批量管理相关缓存：

```go
// 创建带标签的缓存
taggedCache := cache.Cache.Tags("users", "profiles")
taggedCache.Set("user:1", userData, time.Hour)

// 刷新标签下的所有缓存
taggedCache.Flush()
```

### 缓存优化

提供缓存预热、批量操作、统计等功能：

```go
// 创建优化器
optimizer := cache.NewOptimizer(cache.Cache.DefaultStore())

// 缓存预热
items := map[string]interface{}{
    "config:app": appConfig,
    "config:db":  dbConfig,
}
optimizer.WarmUp(items, time.Hour)

// 批量获取
results, _ := optimizer.BatchGet([]string{"key1", "key2"})

// 获取统计信息
stats := optimizer.GetStats()
```

### 缓存统计

带统计功能的缓存包装器：

```go
// 创建带统计的缓存
cacheWithStats := cache.NewCacheWithStats(cache.Cache.DefaultStore())

// 进行缓存操作
cacheWithStats.Set("key", "value", time.Hour)
cacheWithStats.Get("key")

// 获取统计信息
stats := cacheWithStats.GetStats()
fmt.Printf("命中率: %.2f%%\n", stats.HitRate)
```

## 配置示例

### 基本配置

```go
// 初始化缓存系统
cache.Init()

// 注册多个驱动
memoryStore := cache.NewMemoryStore()
fileStore := cache.NewFileStore("./storage/cache")

cache.Cache.Extend("memory", memoryStore)
cache.Cache.Extend("file", fileStore)

// 设置默认驱动
cache.Cache.SetDefaultStore("memory")
```

### 生产环境配置

```go
// 使用Redis作为主缓存
redisClient := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})
redisStore := cache.NewRedisStore(redisClient)

// 使用文件缓存作为备用
fileStore := cache.NewFileStore("./storage/cache")

cache.Cache.Extend("redis", redisStore)
cache.Cache.Extend("file", fileStore)
cache.Cache.SetDefaultStore("redis")
```

## 性能考虑

1. **内存驱动**: 最快，但数据不持久
2. **文件驱动**: 中等性能，数据持久
3. **数据库驱动**: 较慢，但支持复杂查询
4. **Redis 驱动**: 高性能，支持分布式
5. **MongoDB 驱动**: 中等性能，支持复杂数据结构

## 最佳实践

1. **选择合适的驱动**: 根据应用场景选择最适合的驱动
2. **使用缓存标签**: 合理使用标签管理相关缓存
3. **设置合理的 TTL**: 避免缓存过期时间过长或过短
4. **监控缓存性能**: 使用统计功能监控缓存命中率
5. **实现缓存预热**: 在应用启动时预热重要缓存
6. **处理缓存穿透**: 使用 Remember 方法避免缓存穿透

## 依赖管理

### Redis 驱动依赖

```bash
go get github.com/redis/go-redis/v9
```

### MongoDB 驱动依赖

```bash
go get go.mongodb.org/mongo-driver/mongo
```

### 快速安装所有依赖

```bash
# 安装Redis和MongoDB驱动
go get github.com/redis/go-redis/v9
go get go.mongodb.org/mongo-driver/mongo

# 或者一次性安装
go mod tidy
```

## 测试

所有驱动都有完整的单元测试：

```bash
go test ./framework/cache -v
```

## 使用示例

### 基本使用

```go
package main

import (
    "time"
    "laravel-go/framework/cache"
)

func main() {
    // 初始化缓存系统
    cache.Init()

    // 设置缓存
    cache.Set("key", "value", time.Hour)

    // 获取缓存
    value, err := cache.Get("key")
    if err == nil {
        fmt.Printf("缓存值: %v\n", value)
    }
}
```

### 多驱动配置

```go
package main

import (
    "context"
    "time"
    "github.com/redis/go-redis/v9"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "laravel-go/framework/cache"
)

func main() {
    // 初始化缓存系统
    cache.Init()

    // 注册内存驱动
    memoryStore := cache.NewMemoryStore()
    cache.Cache.Extend("memory", memoryStore)

    // 注册文件驱动
    fileStore := cache.NewFileStore("./storage/cache")
    cache.Cache.Extend("file", fileStore)

    // 注册Redis驱动
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    redisStore := cache.NewRedisStore(redisClient)
    cache.Cache.Extend("redis", redisStore)

    // 注册MongoDB驱动
    mongoClient, _ := mongo.Connect(context.Background(),
        options.Client().ApplyURI("mongodb://localhost:27017"))
    mongoStore := cache.NewMongoStore(mongoClient, "cache_db", "cache_collection")
    cache.Cache.Extend("mongodb", mongoStore)

    // 设置默认驱动
    cache.Cache.SetDefaultStore("redis")

    // 使用缓存
    cache.Set("app_name", "Laravel-Go", time.Hour)
    value, _ := cache.Get("app_name")
    fmt.Printf("应用名称: %v\n", value)
}
```

### 驱动切换

```go
// 使用特定驱动
redisCache := cache.Cache.Store("redis")
redisCache.Set("redis_key", "redis_value", time.Hour)

mongoCache := cache.Cache.Store("mongodb")
mongoCache.Set("mongo_key", "mongo_value", time.Hour)

// 切换默认驱动
cache.Cache.SetDefaultStore("memory")
cache.Set("default_key", "default_value", time.Hour)
```

## 示例程序

查看 `examples/cache_demo/main.go` 了解完整的使用示例。
