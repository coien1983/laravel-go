package main

import (
	"fmt"
	"log"
	"time"

	"laravel-go/framework/cache"
)

func main() {
	fmt.Println("=== Laravel-Go 缓存系统演示 ===\n")

	// 初始化缓存系统
	cache.Init()

	// 1. 基本缓存操作
	fmt.Println("1. 基本缓存操作:")
	demoBasicCache()

	// 2. 缓存类型操作
	fmt.Println("\n2. 缓存类型操作:")
	demoCacheTypes()

	// 3. 缓存过期
	fmt.Println("\n3. 缓存过期:")
	demoCacheExpiration()

	// 4. 递增递减操作
	fmt.Println("\n4. 递增递减操作:")
	demoIncrementDecrement()

	// 5. Remember 操作
	fmt.Println("\n5. Remember 操作:")
	demoRemember()

	// 6. 标签缓存
	fmt.Println("\n6. 标签缓存:")
	demoTaggedCache()

	// 7. 缓存优化功能
	fmt.Println("\n7. 缓存优化功能:")
	demoCacheOptimization()

	// 8. 文件缓存
	fmt.Println("\n8. 文件缓存:")
	demoFileCache()

	// 9. 多驱动支持
	fmt.Println("\n9. 多驱动支持:")
	demoMultiDriver()

	fmt.Println("\n=== 演示完成 ===")
}

func demoBasicCache() {
	// 设置缓存
	err := cache.Set("user:1", map[string]interface{}{
		"id":    1,
		"name":  "张三",
		"email": "zhangsan@example.com",
	}, time.Minute*5)
	if err != nil {
		log.Printf("设置缓存失败: %v", err)
		return
	}

	// 获取缓存
	value, err := cache.Get("user:1")
	if err != nil {
		log.Printf("获取缓存失败: %v", err)
		return
	}

	fmt.Printf("  获取用户缓存: %v\n", value)

	// 检查缓存是否存在
	if cache.Has("user:1") {
		fmt.Println("  用户缓存存在")
	}

	// 删除缓存
	err = cache.Delete("user:1")
	if err != nil {
		log.Printf("删除缓存失败: %v", err)
		return
	}

	if !cache.Has("user:1") {
		fmt.Println("  用户缓存已删除")
	}
}

func demoCacheTypes() {
	// 字符串缓存
	cache.Cache.SetString("app_name", "Laravel-Go", time.Hour)
	if value, err := cache.Cache.GetString("app_name"); err == nil {
		fmt.Printf("  字符串缓存: %s\n", value)
	}

	// 整数缓存
	cache.Cache.SetInt("user_count", 1000, time.Hour)
	if value, err := cache.Cache.GetInt("user_count"); err == nil {
		fmt.Printf("  整数缓存: %d\n", value)
	}

	// 浮点数缓存
	cache.Cache.SetFloat("pi", 3.14159, time.Hour)
	if value, err := cache.Cache.GetFloat("pi"); err == nil {
		fmt.Printf("  浮点数缓存: %.5f\n", value)
	}

	// 布尔值缓存
	cache.Cache.SetBool("debug_mode", true, time.Hour)
	if value, err := cache.Cache.GetBool("debug_mode"); err == nil {
		fmt.Printf("  布尔值缓存: %t\n", value)
	}

	// 字节数组缓存
	cache.Cache.SetBytes("binary_data", []byte("Hello World"), time.Hour)
	if value, err := cache.Cache.GetBytes("binary_data"); err == nil {
		fmt.Printf("  字节数组缓存: %s\n", string(value))
	}
}

func demoCacheExpiration() {
	// 设置短期过期的缓存
	err := cache.Set("temp_data", "临时数据", time.Second*2)
	if err != nil {
		log.Printf("设置临时缓存失败: %v", err)
		return
	}

	fmt.Println("  设置2秒过期的缓存")

	// 立即检查
	if cache.Has("temp_data") {
		fmt.Println("  缓存存在")
	}

	// 等待3秒
	time.Sleep(time.Second * 3)

	// 再次检查
	if !cache.Has("temp_data") {
		fmt.Println("  缓存已过期")
	}
}

func demoIncrementDecrement() {
	// 递增操作
	value, err := cache.Cache.Increment("counter", 5)
	if err != nil {
		log.Printf("递增失败: %v", err)
		return
	}
	fmt.Printf("  递增5: %d\n", value)

	value, err = cache.Cache.Increment("counter", 3)
	if err != nil {
		log.Printf("递增失败: %v", err)
		return
	}
	fmt.Printf("  再递增3: %d\n", value)

	// 递减操作
	value, err = cache.Cache.Decrement("counter", 2)
	if err != nil {
		log.Printf("递减失败: %v", err)
		return
	}
	fmt.Printf("  递减2: %d\n", value)
}

func demoRemember() {
	// 使用Remember获取缓存，如果不存在则执行回调
	value, err := cache.Remember("expensive_data", time.Minute*5, func() (interface{}, error) {
		fmt.Println("  执行昂贵的计算...")
		time.Sleep(time.Millisecond * 100) // 模拟耗时操作
		return map[string]interface{}{
			"result": "计算结果",
			"time":   time.Now().Unix(),
		}, nil
	})
	if err != nil {
		log.Printf("Remember失败: %v", err)
		return
	}

	fmt.Printf("  获取数据: %v\n", value)

	// 再次调用，应该直接返回缓存
	value2, err := cache.Remember("expensive_data", time.Minute*5, func() (interface{}, error) {
		fmt.Println("  这不应该执行")
		return "新数据", nil
	})
	if err != nil {
		log.Printf("Remember失败: %v", err)
		return
	}

	fmt.Printf("  再次获取数据: %v\n", value2)
}

func demoTaggedCache() {
	// 创建带标签的缓存
	taggedCache := cache.Cache.Tags("users", "profiles")

	// 设置带标签的缓存
	err := taggedCache.Set("user_profile:1", map[string]interface{}{
		"user_id": 1,
		"profile": map[string]interface{}{
			"avatar": "avatar1.jpg",
			"bio":    "用户简介",
		},
	}, time.Hour)
	if err != nil {
		log.Printf("设置标签缓存失败: %v", err)
		return
	}

	fmt.Println("  设置带标签的缓存")

	// 获取标签列表
	tags := taggedCache.GetTags()
	fmt.Printf("  标签列表: %v\n", tags)

	// 刷新标签（使所有相关缓存失效）
	err = taggedCache.Flush()
	if err != nil {
		log.Printf("刷新标签失败: %v", err)
		return
	}

	fmt.Println("  标签缓存已刷新")
}

func demoCacheOptimization() {
	// 创建缓存优化器
	optimizer := cache.NewOptimizer(cache.Cache.DefaultStore())

	// 缓存预热
	items := map[string]interface{}{
		"config:app": map[string]interface{}{
			"name":    "Laravel-Go",
			"version": "1.0.0",
		},
		"config:database": map[string]interface{}{
			"driver": "mysql",
			"host":   "localhost",
		},
		"config:cache": map[string]interface{}{
			"driver": "memory",
			"ttl":    3600,
		},
	}

	err := optimizer.WarmUp(items, time.Hour)
	if err != nil {
		log.Printf("缓存预热失败: %v", err)
		return
	}

	fmt.Println("  缓存预热完成")

	// 批量获取
	keys := []string{"config:app", "config:database", "config:cache"}
	results, err := optimizer.BatchGet(keys)
	if err != nil {
		log.Printf("批量获取失败: %v", err)
		return
	}

	fmt.Printf("  批量获取结果: %d 个项目\n", len(results))

	// 创建带统计的缓存
	cacheWithStats := cache.NewCacheWithStats(cache.Cache.DefaultStore())

	// 进行一些操作
	cacheWithStats.Set("test_key", "test_value", time.Minute)
	cacheWithStats.Get("test_key")     // 命中
	cacheWithStats.Get("non_existent") // 未命中

	// 获取统计信息
	stats := cacheWithStats.GetStats()
	fmt.Printf("  缓存统计 - 命中: %d, 未命中: %d, 命中率: %.2f%%\n",
		stats.Hits, stats.Misses, stats.HitRate)
}

func demoFileCache() {
	// 创建文件缓存存储
	fileStore := cache.NewFileStore("./storage/cache")

	// 设置文件缓存
	err := fileStore.Set("file_data", map[string]interface{}{
		"message":   "这是文件缓存数据",
		"timestamp": time.Now().Unix(),
	}, time.Minute*10)
	if err != nil {
		log.Printf("设置文件缓存失败: %v", err)
		return
	}

	fmt.Println("  文件缓存已设置")

	// 获取文件缓存
	value, err := fileStore.Get("file_data")
	if err != nil {
		log.Printf("获取文件缓存失败: %v", err)
		return
	}

	fmt.Printf("  文件缓存数据: %v\n", value)

	// 检查文件缓存是否存在
	if fileStore.Has("file_data") {
		fmt.Println("  文件缓存存在")
	}

	// 删除文件缓存
	err = fileStore.Delete("file_data")
	if err != nil {
		log.Printf("删除文件缓存失败: %v", err)
		return
	}

	fmt.Println("  文件缓存已删除")
}

func demoMultiDriver() {
	fmt.Println("  当前支持的缓存驱动:")
	fmt.Println("    - 内存驱动 (MemoryStore): 高性能内存缓存")
	fmt.Println("    - 文件驱动 (FileStore): 持久化文件缓存")
	fmt.Println("    - 数据库驱动 (DatabaseStore): 支持MySQL、PostgreSQL、SQLite")
	fmt.Println("    - Redis驱动 (RedisStore): 高性能分布式缓存")
	fmt.Println("    - MongoDB驱动 (MongoStore): 文档型数据库缓存")

	fmt.Println("\n  扩展新驱动的方法:")
	fmt.Println("    1. 实现 Store 接口")
	fmt.Println("    2. 注册到缓存管理器")
	fmt.Println("    3. 设置默认驱动")

	fmt.Println("\n  Redis驱动示例:")
	fmt.Println("    // 创建Redis客户端")
	fmt.Println("    redisClient := redis.NewClient(&redis.Options{")
	fmt.Println("        Addr: \"localhost:6379\",")
	fmt.Println("    })")
	fmt.Println("    redisStore := cache.NewRedisStore(redisClient)")
	fmt.Println("    cache.Cache.Extend(\"redis\", redisStore)")

	fmt.Println("\n  MongoDB驱动示例:")
	fmt.Println("    // 创建MongoDB客户端")
	fmt.Println("    mongoClient, _ := mongo.Connect(context.Background(), \"mongodb://localhost:27017\")")
	fmt.Println("    mongoStore := cache.NewMongoStore(mongoClient, \"cache_db\", \"cache_collection\")")
	fmt.Println("    cache.Cache.Extend(\"mongodb\", mongoStore)")
}
