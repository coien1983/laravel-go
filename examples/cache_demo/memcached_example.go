package main

import (
	"fmt"
	"log"
	"time"

	"laravel-go/framework/cache"
)

func memcachedExample() {
	fmt.Println("=== Memcached 缓存驱动示例 ===")

	// 初始化缓存管理器
	cache.Init()

	// 创建 Memcached 存储
	memcachedStore := cache.NewMemcachedStore("127.0.0.1:11211")
	cache.Cache.Extend("memcached", memcachedStore)

	// 设置默认存储为 memcached
	cache.Cache.SetDefaultStore("memcached")

	fmt.Println("✅ Memcached 缓存驱动初始化完成")

	// 基本操作示例
	memcachedBasicOperations()

	// 类型化操作示例
	memcachedTypedOperations()

	// 高级功能示例
	memcachedAdvancedFeatures()

	fmt.Println("✅ Memcached 缓存驱动示例完成")
}

// 基本操作示例
func memcachedBasicOperations() {
	fmt.Println("\n--- 基本操作示例 ---")

	// 设置缓存
	err := cache.Cache.Set("user:1", map[string]interface{}{
		"id":   1,
		"name": "张三",
		"age":  25,
	}, 5*time.Minute)
	if err != nil {
		log.Printf("设置缓存失败: %v", err)
		return
	}
	fmt.Println("✅ 设置缓存: user:1")

	// 获取缓存
	value, err := cache.Cache.Get("user:1")
	if err != nil {
		log.Printf("获取缓存失败: %v", err)
		return
	}
	fmt.Printf("✅ 获取缓存: user:1 = %v\n", value)

	// 检查缓存是否存在
	if cache.Cache.Has("user:1") {
		fmt.Println("✅ 缓存存在: user:1")
	}

	// 删除缓存
	err = cache.Cache.Delete("user:1")
	if err != nil {
		log.Printf("删除缓存失败: %v", err)
		return
	}
	fmt.Println("✅ 删除缓存: user:1")

	// 检查缓存是否不存在
	if !cache.Cache.Has("user:1") {
		fmt.Println("✅ 缓存不存在: user:1")
	}
}

// 类型化操作示例
func memcachedTypedOperations() {
	fmt.Println("\n--- 类型化操作示例 ---")

	// 字符串操作
	err := cache.Cache.SetString("greeting", "Hello, Memcached!", 2*time.Minute)
	if err != nil {
		log.Printf("设置字符串缓存失败: %v", err)
		return
	}
	greeting, err := cache.Cache.GetString("greeting")
	if err != nil {
		log.Printf("获取字符串缓存失败: %v", err)
		return
	}
	fmt.Printf("✅ 字符串缓存: %s\n", greeting)

	// 整数操作
	err = cache.Cache.SetInt("counter", 100, 3*time.Minute)
	if err != nil {
		log.Printf("设置整数缓存失败: %v", err)
		return
	}
	counter, err := cache.Cache.GetInt("counter")
	if err != nil {
		log.Printf("获取整数缓存失败: %v", err)
		return
	}
	fmt.Printf("✅ 整数缓存: %d\n", counter)

	// 浮点数操作
	err = cache.Cache.SetFloat("price", 99.99, 1*time.Minute)
	if err != nil {
		log.Printf("设置浮点数缓存失败: %v", err)
		return
	}
	price, err := cache.Cache.GetFloat("price")
	if err != nil {
		log.Printf("获取浮点数缓存失败: %v", err)
		return
	}
	fmt.Printf("✅ 浮点数缓存: %.2f\n", price)

	// 布尔值操作
	err = cache.Cache.SetBool("is_active", true, 5*time.Minute)
	if err != nil {
		log.Printf("设置布尔值缓存失败: %v", err)
		return
	}
	isActive, err := cache.Cache.GetBool("is_active")
	if err != nil {
		log.Printf("获取布尔值缓存失败: %v", err)
		return
	}
	fmt.Printf("✅ 布尔值缓存: %t\n", isActive)

	// 字节数组操作
	data := []byte("Hello, World!")
	err = cache.Cache.SetBytes("binary_data", data, 1*time.Minute)
	if err != nil {
		log.Printf("设置字节数组缓存失败: %v", err)
		return
	}
	binaryData, err := cache.Cache.GetBytes("binary_data")
	if err != nil {
		log.Printf("获取字节数组缓存失败: %v", err)
		return
	}
	fmt.Printf("✅ 字节数组缓存: %s\n", string(binaryData))
}

// 高级功能示例
func memcachedAdvancedFeatures() {
	fmt.Println("\n--- 高级功能示例 ---")

	// Remember 功能
	value, err := cache.Cache.Remember("expensive_operation", 10*time.Minute, func() (interface{}, error) {
		fmt.Println("🔄 执行昂贵的操作...")
		time.Sleep(100 * time.Millisecond) // 模拟耗时操作
		return "计算结果", nil
	})
	if err != nil {
		log.Printf("Remember 操作失败: %v", err)
		return
	}
	fmt.Printf("✅ Remember 结果: %v\n", value)

	// 再次调用，应该从缓存获取
	value2, err := cache.Cache.Remember("expensive_operation", 10*time.Minute, func() (interface{}, error) {
		fmt.Println("🔄 执行昂贵的操作...") // 这行不应该执行
		return "计算结果", nil
	})
	if err != nil {
		log.Printf("Remember 操作失败: %v", err)
		return
	}
	fmt.Printf("✅ Remember 缓存结果: %v\n", value2)

	// 递增和递减操作
	err = cache.Cache.SetInt("visit_count", 0, 1*time.Hour)
	if err != nil {
		log.Printf("设置访问计数失败: %v", err)
		return
	}

	// 递增
	newCount, err := cache.Cache.Increment("visit_count", 1)
	if err != nil {
		log.Printf("递增失败: %v", err)
		return
	}
	fmt.Printf("✅ 访问计数递增: %d\n", newCount)

	// 再次递增
	newCount, err = cache.Cache.Increment("visit_count", 5)
	if err != nil {
		log.Printf("递增失败: %v", err)
		return
	}
	fmt.Printf("✅ 访问计数递增: %d\n", newCount)

	// 递减
	newCount, err = cache.Cache.Decrement("visit_count", 2)
	if err != nil {
		log.Printf("递减失败: %v", err)
		return
	}
	fmt.Printf("✅ 访问计数递减: %d\n", newCount)

	// 批量操作
	keys := []string{"key1", "key2", "key3"}
	for i, key := range keys {
		err := cache.Cache.SetString(key, fmt.Sprintf("value%d", i+1), 5*time.Minute)
		if err != nil {
			log.Printf("设置缓存失败: %v", err)
			continue
		}
		fmt.Printf("✅ 设置缓存: %s\n", key)
	}

	// 批量删除
	err = cache.Cache.DeleteMultiple(keys)
	if err != nil {
		log.Printf("批量删除失败: %v", err)
		return
	}
	fmt.Println("✅ 批量删除完成")

	// 检查是否都已删除
	for _, key := range keys {
		if !cache.Cache.Has(key) {
			fmt.Printf("✅ 缓存已删除: %s\n", key)
		}
	}
}

// 配置示例
func configExample() {
	fmt.Println("\n--- 配置示例 ---")

	// 使用配置创建 Memcached 存储
	config := map[string]interface{}{
		"host": "127.0.0.1",
		"port": "11211",
	}

	memcachedStore := cache.NewMemcachedStoreWithConfig(config)
	cache.Cache.Extend("memcached_config", memcachedStore)

	// 使用配置的存储
	err := cache.Cache.Store("memcached_config").SetString("config_test", "配置测试", 1*time.Minute)
	if err != nil {
		log.Printf("配置测试失败: %v", err)
		return
	}

	value, err := cache.Cache.Store("memcached_config").GetString("config_test")
	if err != nil {
		log.Printf("获取配置测试失败: %v", err)
		return
	}

	fmt.Printf("✅ 配置测试: %s\n", value)
}
