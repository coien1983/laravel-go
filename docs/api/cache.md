# 缓存 API 参考

## 🗄️ 缓存系统概览

Laravel-Go Framework 提供了统一的缓存接口，支持多种缓存驱动，包括内存、文件、Redis、MongoDB 等。

## 🚀 快速开始

### 基本使用

```go
import "laravel-go/framework/cache"

// 获取缓存实例
cache := cache.Driver("default")

// 存储数据
cache.Put("key", "value", time.Hour)

// 获取数据
value := cache.Get("key")

// 检查是否存在
exists := cache.Has("key")

// 删除数据
cache.Forget("key")
```

## 📋 API 参考

### 核心方法

#### Put - 存储数据

```go
// 存储数据到缓存
func (c *Cache) Put(key string, value interface{}, ttl time.Duration) error

// 示例
err := cache.Put("user:1", user, time.Hour)
err := cache.Put("config:app", config, 24*time.Hour)
err := cache.Put("temp:data", data, time.Minute*5)
```

#### Get - 获取数据

```go
// 获取缓存数据
func (c *Cache) Get(key string) (interface{}, bool)

// 示例
value, exists := cache.Get("user:1")
if exists {
    user := value.(*User)
    // 使用 user
}

// 获取带默认值的数据
func (c *Cache) Get(key string, defaultValue interface{}) interface{}

// 示例
user := cache.Get("user:1", &User{}).(*User)
```

#### Has - 检查存在

```go
// 检查键是否存在
func (c *Cache) Has(key string) bool

// 示例
if cache.Has("user:1") {
    // 键存在
}
```

#### Forget - 删除数据

```go
// 删除指定的键
func (c *Cache) Forget(key string) error

// 示例
err := cache.Forget("user:1")
err := cache.Forget("temp:data")
```

#### Flush - 清空缓存

```go
// 清空所有缓存
func (c *Cache) Flush() error

// 示例
err := cache.Flush()
```

#### Remember - 记住数据

```go
// 如果键不存在，则执行回调函数并缓存结果
func (c *Cache) Remember(key string, ttl time.Duration, callback func() interface{}) interface{}

// 示例
user := cache.Remember("user:1", time.Hour, func() interface{} {
    return userService.GetUser(1)
}).(*User)
```

#### RememberForever - 永久记住

```go
// 永久缓存数据
func (c *Cache) RememberForever(key string, callback func() interface{}) interface{}

// 示例
config := cache.RememberForever("config:app", func() interface{} {
    return loadAppConfig()
}).(*Config)
```

### 高级方法

#### Increment - 递增

```go
// 递增数值
func (c *Cache) Increment(key string, value int64) (int64, error)

// 示例
count, err := cache.Increment("visits:page:1", 1)
count, err := cache.Increment("score:user:1", 10)
```

#### Decrement - 递减

```go
// 递减数值
func (c *Cache) Decrement(key string, value int64) (int64, error)

// 示例
count, err := cache.Decrement("stock:product:1", 1)
count, err := cache.Decrement("lives:player:1", 1)
```

#### Tags - 标签管理

```go
// 使用标签
func (c *Cache) Tags(names ...string) *TaggedCache

// 示例
taggedCache := cache.Tags("users", "profiles")
taggedCache.Put("user:1", user, time.Hour)
taggedCache.Flush() // 清除所有带这些标签的缓存
```

#### Multiple - 批量操作

```go
// 批量获取
func (c *Cache) Many(keys []string) map[string]interface{}

// 示例
keys := []string{"user:1", "user:2", "user:3"}
values := cache.Many(keys)

// 批量存储
func (c *Cache) PutMany(values map[string]interface{}, ttl time.Duration) error

// 示例
values := map[string]interface{}{
    "user:1": user1,
    "user:2": user2,
    "user:3": user3,
}
err := cache.PutMany(values, time.Hour)
```

## 🔧 缓存驱动

### 内存驱动

```go
// 使用内存缓存
cache := cache.Driver("memory")

// 配置
type MemoryConfig struct {
    DefaultTTL time.Duration `env:"CACHE_MEMORY_TTL" default:"1h"`
    MaxSize    int           `env:"CACHE_MEMORY_MAX_SIZE" default:"1000"`
}
```

### 文件驱动

```go
// 使用文件缓存
cache := cache.Driver("file")

// 配置
type FileConfig struct {
    Path      string        `env:"CACHE_FILE_PATH" default:"storage/cache"`
    DefaultTTL time.Duration `env:"CACHE_FILE_TTL" default:"1h"`
}
```

### Redis 驱动

```go
// 使用 Redis 缓存
cache := cache.Driver("redis")

// 配置
type RedisConfig struct {
    Host      string `env:"CACHE_REDIS_HOST" default:"localhost"`
    Port      int    `env:"CACHE_REDIS_PORT" default:"6379"`
    Password  string `env:"CACHE_REDIS_PASSWORD"`
    Database  int    `env:"CACHE_REDIS_DB" default:"0"`
    Prefix    string `env:"CACHE_REDIS_PREFIX" default:"laravel_go:"`
}
```

### MongoDB 驱动

```go
// 使用 MongoDB 缓存
cache := cache.Driver("mongodb")

// 配置
type MongoConfig struct {
    URI       string `env:"CACHE_MONGODB_URI" default:"mongodb://localhost:27017"`
    Database  string `env:"CACHE_MONGODB_DB" default:"laravel_go"`
    Collection string `env:"CACHE_MONGODB_COLLECTION" default:"cache"`
}
```

### Memcached 驱动

```go
// 使用 Memcached 缓存
cache := cache.Driver("memcached")

// 配置
type MemcachedConfig struct {
    Host      string `env:"CACHE_MEMCACHED_HOST" default:"127.0.0.1"`
    Port      int    `env:"CACHE_MEMCACHED_PORT" default:"11211"`
    Username  string `env:"CACHE_MEMCACHED_USERNAME"`
    Password  string `env:"CACHE_MEMCACHED_PASSWORD"`
    PersistentID string `env:"CACHE_MEMCACHED_PERSISTENT_ID"`
}

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

## 🎯 使用示例

### 用户数据缓存

```go
type UserService struct {
    cache cache.Cache
    db    *database.Connection
}

func (s *UserService) GetUser(id int) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 尝试从缓存获取
    if cached, exists := s.cache.Get(cacheKey); exists {
        return cached.(*User), nil
    }

    // 从数据库获取
    var user User
    err := s.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }

    // 缓存用户数据
    s.cache.Put(cacheKey, &user, time.Hour)

    return &user, nil
}

func (s *UserService) UpdateUser(id int, data map[string]interface{}) error {
    // 更新数据库
    err := s.db.Model(&User{}).Where("id = ?", id).Updates(data).Error
    if err != nil {
        return err
    }

    // 清除缓存
    cacheKey := fmt.Sprintf("user:%d", id)
    s.cache.Forget(cacheKey)

    return nil
}
```

### 配置缓存

```go
type ConfigService struct {
    cache cache.Cache
}

func (s *ConfigService) GetConfig(key string) interface{} {
    cacheKey := fmt.Sprintf("config:%s", key)

    return s.cache.Remember(cacheKey, 24*time.Hour, func() interface{} {
        return s.loadConfigFromDatabase(key)
    })
}

func (s *ConfigService) SetConfig(key string, value interface{}) error {
    // 更新数据库
    err := s.updateConfigInDatabase(key, value)
    if err != nil {
        return err
    }

    // 更新缓存
    cacheKey := fmt.Sprintf("config:%s", key)
    return s.cache.Put(cacheKey, value, 24*time.Hour)
}
```

### 会话缓存

```go
type SessionService struct {
    cache cache.Cache
}

func (s *SessionService) GetSession(sessionID string) (*Session, error) {
    cacheKey := fmt.Sprintf("session:%s", sessionID)

    if cached, exists := s.cache.Get(cacheKey); exists {
        return cached.(*Session), nil
    }

    return nil, errors.New("session not found")
}

func (s *SessionService) SetSession(session *Session) error {
    cacheKey := fmt.Sprintf("session:%s", session.ID)
    return s.cache.Put(cacheKey, session, 30*time.Minute)
}

func (s *SessionService) DeleteSession(sessionID string) error {
    cacheKey := fmt.Sprintf("session:%s", sessionID)
    return s.cache.Forget(cacheKey)
}
```

### 页面缓存

```go
type PageCache struct {
    cache cache.Cache
}

func (p *PageCache) GetPage(url string) (string, bool) {
    cacheKey := fmt.Sprintf("page:%s", url)
    content, exists := p.cache.Get(cacheKey)
    if exists {
        return content.(string), true
    }
    return "", false
}

func (p *PageCache) SetPage(url, content string) error {
    cacheKey := fmt.Sprintf("page:%s", url)
    return p.cache.Put(cacheKey, content, time.Hour)
}

func (p *PageCache) ClearPage(url string) error {
    cacheKey := fmt.Sprintf("page:%s", url)
    return p.cache.Forget(cacheKey)
}
```

## 🔄 缓存标签

### 标签使用

```go
// 使用标签
usersCache := cache.Tags("users")
profilesCache := cache.Tags("profiles")

// 存储带标签的数据
usersCache.Put("user:1", user1, time.Hour)
usersCache.Put("user:2", user2, time.Hour)
profilesCache.Put("profile:1", profile1, time.Hour)

// 清除特定标签的缓存
usersCache.Flush() // 清除所有 users 标签的缓存

// 清除多个标签的缓存
cache.Tags("users", "profiles").Flush()
```

### 标签示例

```go
type UserService struct {
    cache cache.Cache
}

func (s *UserService) GetUser(id int) (*User, error) {
    taggedCache := s.cache.Tags("users")
    cacheKey := fmt.Sprintf("user:%d", id)

    return taggedCache.Remember(cacheKey, time.Hour, func() interface{} {
        var user User
        s.db.First(&user, id)
        return &user
    }).(*User), nil
}

func (s *UserService) UpdateUser(id int, data map[string]interface{}) error {
    // 更新数据库
    err := s.db.Model(&User{}).Where("id = ?", id).Updates(data).Error
    if err != nil {
        return err
    }

    // 清除用户相关的所有缓存
    s.cache.Tags("users").Flush()

    return nil
}
```

## 📊 缓存统计

### 统计信息

```go
// 获取缓存统计信息
func (c *Cache) Stats() *CacheStats

type CacheStats struct {
    Hits   int64 `json:"hits"`
    Misses int64 `json:"misses"`
    Keys   int64 `json:"keys"`
    Size   int64 `json:"size"`
}

// 示例
stats := cache.Stats()
fmt.Printf("Cache hits: %d, misses: %d, keys: %d\n",
    stats.Hits, stats.Misses, stats.Keys)
```

### 监控缓存

```go
type CacheMonitor struct {
    cache cache.Cache
}

func (m *CacheMonitor) GetHitRate() float64 {
    stats := m.cache.Stats()
    total := stats.Hits + stats.Misses
    if total == 0 {
        return 0
    }
    return float64(stats.Hits) / float64(total)
}

func (m *CacheMonitor) GetCacheSize() int64 {
    stats := m.cache.Stats()
    return stats.Size
}
```

## 🛡️ 错误处理

### 错误类型

```go
// 缓存错误类型
type CacheError struct {
    Message string
    Key     string
    Err     error
}

func (e *CacheError) Error() string {
    return fmt.Sprintf("cache error for key %s: %s", e.Key, e.Message)
}

// 处理缓存错误
func handleCacheError(err error, key string) {
    if cacheErr, ok := err.(*CacheError); ok {
        log.Printf("Cache error for key %s: %v", cacheErr.Key, cacheErr.Err)
    } else {
        log.Printf("Unknown cache error: %v", err)
    }
}
```

### 错误处理示例

```go
func (s *UserService) GetUser(id int) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 尝试从缓存获取
    if cached, exists := s.cache.Get(cacheKey); exists {
        return cached.(*User), nil
    }

    // 从数据库获取
    var user User
    err := s.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }

    // 缓存用户数据，处理错误
    if err := s.cache.Put(cacheKey, &user, time.Hour); err != nil {
        log.Printf("Failed to cache user %d: %v", id, err)
        // 不返回错误，因为数据已经获取成功
    }

    return &user, nil
}
```

## 📝 最佳实践

### 1. 缓存键命名

```go
// 使用一致的命名规范
const (
    UserCachePrefix     = "user:"
    ConfigCachePrefix   = "config:"
    SessionCachePrefix  = "session:"
    PageCachePrefix     = "page:"
)

// 生成缓存键
func getUserCacheKey(id int) string {
    return fmt.Sprintf("%s%d", UserCachePrefix, id)
}

func getConfigCacheKey(key string) string {
    return fmt.Sprintf("%s%s", ConfigCachePrefix, key)
}
```

### 2. TTL 设置

```go
// 根据数据类型设置合适的 TTL
const (
    UserCacheTTL      = time.Hour
    ConfigCacheTTL    = 24 * time.Hour
    SessionCacheTTL   = 30 * time.Minute
    PageCacheTTL      = time.Hour
    TempDataTTL       = 5 * time.Minute
)
```

### 3. 缓存穿透防护

```go
func (s *UserService) GetUser(id int) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 尝试从缓存获取
    if cached, exists := s.cache.Get(cacheKey); exists {
        if cached == nil {
            return nil, errors.New("user not found")
        }
        return cached.(*User), nil
    }

    // 从数据库获取
    var user User
    err := s.db.First(&user, id).Error
    if err != nil {
        // 缓存空值，防止缓存穿透
        s.cache.Put(cacheKey, nil, time.Minute*5)
        return nil, err
    }

    // 缓存用户数据
    s.cache.Put(cacheKey, &user, time.Hour)

    return &user, nil
}
```

### 4. 缓存更新策略

```go
// 写入时更新缓存
func (s *UserService) UpdateUser(id int, data map[string]interface{}) error {
    // 更新数据库
    err := s.db.Model(&User{}).Where("id = ?", id).Updates(data).Error
    if err != nil {
        return err
    }

    // 更新缓存
    cacheKey := fmt.Sprintf("user:%d", id)
    var user User
    s.db.First(&user, id)
    s.cache.Put(cacheKey, &user, time.Hour)

    return nil
}

// 或者清除缓存，让下次读取时重新加载
func (s *UserService) UpdateUser(id int, data map[string]interface{}) error {
    // 更新数据库
    err := s.db.Model(&User{}).Where("id = ?", id).Updates(data).Error
    if err != nil {
        return err
    }

    // 清除缓存
    cacheKey := fmt.Sprintf("user:%d", id)
    s.cache.Forget(cacheKey)

    return nil
}
```

## 📚 总结

Laravel-Go Framework 的缓存 API 提供了：

1. **统一接口**: 支持多种缓存驱动
2. **丰富功能**: 支持标签、批量操作、统计等
3. **高性能**: 内存缓存和分布式缓存支持
4. **易用性**: 简洁的 API 设计
5. **可扩展性**: 易于添加新的缓存驱动

通过合理使用缓存 API，可以显著提升应用程序的性能和用户体验。
