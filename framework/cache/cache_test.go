package cache

import (
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()

	if manager.stores == nil {
		t.Error("Expected stores map to be initialized")
	}

	if manager.defaultStore != "memory" {
		t.Errorf("Expected default store 'memory', got %s", manager.defaultStore)
	}

	if manager.config == nil {
		t.Error("Expected config map to be initialized")
	}
}

func TestManagerExtend(t *testing.T) {
	manager := NewManager()
	memoryStore := NewMemoryStore()

	manager.Extend("test", memoryStore)

	if _, exists := manager.stores["test"]; !exists {
		t.Error("Store should be extended")
	}
}

func TestManagerSetDefaultStore(t *testing.T) {
	manager := NewManager()
	memoryStore := NewMemoryStore()

	manager.Extend("test", memoryStore)
	manager.SetDefaultStore("test")

	if manager.defaultStore != "test" {
		t.Errorf("Expected default store 'test', got %s", manager.defaultStore)
	}
}

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore()

	// 测试设置和获取
	err := store.Set("test_key", "test_value", time.Minute)
	if err != nil {
		t.Errorf("Set should not return error: %v", err)
	}

	value, err := store.Get("test_key")
	if err != nil {
		t.Errorf("Get should not return error: %v", err)
	}

	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %v", value)
	}

	// 测试检查存在
	if !store.Has("test_key") {
		t.Error("Key should exist")
	}

	if store.Missing("test_key") {
		t.Error("Key should not be missing")
	}

	// 测试删除
	err = store.Delete("test_key")
	if err != nil {
		t.Errorf("Delete should not return error: %v", err)
	}

	if store.Has("test_key") {
		t.Error("Key should not exist after deletion")
	}
}

func TestMemoryStoreExpiration(t *testing.T) {
	store := NewMemoryStore()

	// 设置短期过期的缓存
	err := store.Set("expire_key", "expire_value", time.Millisecond*10)
	if err != nil {
		t.Errorf("Set should not return error: %v", err)
	}

	// 等待过期
	time.Sleep(time.Millisecond * 20)

	// 检查是否过期
	if store.Has("expire_key") {
		t.Error("Key should be expired")
	}

	_, err = store.Get("expire_key")
	if err == nil {
		t.Error("Get should return error for expired key")
	}
}

func TestMemoryStoreTypes(t *testing.T) {
	store := NewMemoryStore()

	// 测试字符串
	err := store.SetString("string_key", "string_value", time.Minute)
	if err != nil {
		t.Errorf("SetString should not return error: %v", err)
	}

	value, err := store.GetString("string_key")
	if err != nil {
		t.Errorf("GetString should not return error: %v", err)
	}

	if value != "string_value" {
		t.Errorf("Expected 'string_value', got %s", value)
	}

	// 测试整数
	err = store.SetInt("int_key", 42, time.Minute)
	if err != nil {
		t.Errorf("SetInt should not return error: %v", err)
	}

	intValue, err := store.GetInt("int_key")
	if err != nil {
		t.Errorf("GetInt should not return error: %v", err)
	}

	if intValue != 42 {
		t.Errorf("Expected 42, got %d", intValue)
	}

	// 测试浮点数
	err = store.SetFloat("float_key", 3.14, time.Minute)
	if err != nil {
		t.Errorf("SetFloat should not return error: %v", err)
	}

	floatValue, err := store.GetFloat("float_key")
	if err != nil {
		t.Errorf("GetFloat should not return error: %v", err)
	}

	if floatValue != 3.14 {
		t.Errorf("Expected 3.14, got %f", floatValue)
	}

	// 测试布尔值
	err = store.SetBool("bool_key", true, time.Minute)
	if err != nil {
		t.Errorf("SetBool should not return error: %v", err)
	}

	boolValue, err := store.GetBool("bool_key")
	if err != nil {
		t.Errorf("GetBool should not return error: %v", err)
	}

	if !boolValue {
		t.Error("Expected true, got false")
	}
}

func TestMemoryStoreIncrementDecrement(t *testing.T) {
	store := NewMemoryStore()

	// 测试递增
	value, err := store.Increment("counter", 5)
	if err != nil {
		t.Errorf("Increment should not return error: %v", err)
	}

	if value != 5 {
		t.Errorf("Expected 5, got %d", value)
	}

	// 再次递增
	value, err = store.Increment("counter", 3)
	if err != nil {
		t.Errorf("Increment should not return error: %v", err)
	}

	if value != 8 {
		t.Errorf("Expected 8, got %d", value)
	}

	// 测试递减
	value, err = store.Decrement("counter", 2)
	if err != nil {
		t.Errorf("Decrement should not return error: %v", err)
	}

	if value != 6 {
		t.Errorf("Expected 6, got %d", value)
	}
}

func TestMemoryStoreRemember(t *testing.T) {
	store := NewMemoryStore()

	// 测试Remember
	value, err := store.Remember("remember_key", time.Minute, func() (interface{}, error) {
		return "remembered_value", nil
	})
	if err != nil {
		t.Errorf("Remember should not return error: %v", err)
	}

	if value != "remembered_value" {
		t.Errorf("Expected 'remembered_value', got %v", value)
	}

	// 再次调用应该返回缓存的值
	value, err = store.Remember("remember_key", time.Minute, func() (interface{}, error) {
		return "new_value", nil
	})
	if err != nil {
		t.Errorf("Remember should not return error: %v", err)
	}

	if value != "remembered_value" {
		t.Errorf("Expected 'remembered_value', got %v", value)
	}
}

func TestMemoryStoreBatchOperations(t *testing.T) {
	store := NewMemoryStore()

	// 设置多个值
	items := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for key, value := range items {
		err := store.Set(key, value, time.Minute)
		if err != nil {
			t.Errorf("Set should not return error: %v", err)
		}
	}

	// 测试批量删除
	err := store.DeleteMultiple([]string{"key1", "key2"})
	if err != nil {
		t.Errorf("DeleteMultiple should not return error: %v", err)
	}

	// 检查删除结果
	if store.Has("key1") || store.Has("key2") {
		t.Error("Keys should be deleted")
	}

	if !store.Has("key3") {
		t.Error("Key3 should still exist")
	}

	// 测试清空
	err = store.Clear()
	if err != nil {
		t.Errorf("Clear should not return error: %v", err)
	}

	if store.Has("key3") {
		t.Error("All keys should be cleared")
	}
}

func TestFileStore(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	store := NewFileStore(tempDir)

	// 测试设置和获取
	err := store.Set("test_key", "test_value", time.Minute)
	if err != nil {
		t.Errorf("Set should not return error: %v", err)
	}

	value, err := store.Get("test_key")
	if err != nil {
		t.Errorf("Get should not return error: %v", err)
	}

	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %v", value)
	}

	// 测试检查存在
	if !store.Has("test_key") {
		t.Error("Key should exist")
	}

	// 测试删除
	err = store.Delete("test_key")
	if err != nil {
		t.Errorf("Delete should not return error: %v", err)
	}

	if store.Has("test_key") {
		t.Error("Key should not exist after deletion")
	}
}

func TestFileStoreExpiration(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	store := NewFileStore(tempDir)

	// 设置短期过期的缓存
	err := store.Set("expire_key", "expire_value", time.Millisecond*10)
	if err != nil {
		t.Errorf("Set should not return error: %v", err)
	}

	// 等待过期
	time.Sleep(time.Millisecond * 20)

	// 检查是否过期
	if store.Has("expire_key") {
		t.Error("Key should be expired")
	}

	_, err = store.Get("expire_key")
	if err == nil {
		t.Error("Get should return error for expired key")
	}
}

func TestTaggedStore(t *testing.T) {
	store := NewMemoryStore()
	taggedStore := store.Tags("users", "profiles")

	// 设置带标签的缓存
	err := taggedStore.Set("user_1", "user_data", time.Minute)
	if err != nil {
		t.Errorf("Set should not return error: %v", err)
	}

	// 获取带标签的缓存
	value, err := taggedStore.Get("user_1")
	if err != nil {
		t.Errorf("Get should not return error: %v", err)
	}

	if value != "user_data" {
		t.Errorf("Expected 'user_data', got %v", value)
	}

	// 测试标签列表
	tags := taggedStore.GetTags()
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}

	// 测试刷新标签
	err = taggedStore.Flush()
	if err != nil {
		t.Errorf("Flush should not return error: %v", err)
	}
}

func TestCacheOptimizer(t *testing.T) {
	store := NewMemoryStore()
	optimizer := NewOptimizer(store)

	// 测试缓存预热
	items := map[string]interface{}{
		"warm_key1": "warm_value1",
		"warm_key2": "warm_value2",
	}

	err := optimizer.WarmUp(items, time.Minute)
	if err != nil {
		t.Errorf("WarmUp should not return error: %v", err)
	}

	// 验证预热结果
	for key, expectedValue := range items {
		value, err := store.Get(key)
		if err != nil {
			t.Errorf("Get should not return error: %v", err)
		}

		if value != expectedValue {
			t.Errorf("Expected %v, got %v", expectedValue, value)
		}
	}

	// 测试批量获取
	keys := []string{"warm_key1", "warm_key2"}
	results, err := optimizer.BatchGet(keys)
	if err != nil {
		t.Errorf("BatchGet should not return error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestCacheWithStats(t *testing.T) {
	store := NewMemoryStore()
	cacheWithStats := NewCacheWithStats(store)

	// 设置缓存
	err := cacheWithStats.Set("stats_key", "stats_value", time.Minute)
	if err != nil {
		t.Errorf("Set should not return error: %v", err)
	}

	// 获取缓存（命中）
	value, err := cacheWithStats.Get("stats_key")
	if err != nil {
		t.Errorf("Get should not return error: %v", err)
	}

	if value != "stats_value" {
		t.Errorf("Expected 'stats_value', got %v", value)
	}

	// 获取不存在的缓存（未命中）
	_, err = cacheWithStats.Get("non_existent_key")
	if err == nil {
		t.Error("Get should return error for non-existent key")
	}

	// 检查统计信息
	stats := cacheWithStats.GetStats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}

	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}

	if stats.HitRate != 50.0 {
		t.Errorf("Expected 50.0%% hit rate, got %f", stats.HitRate)
	}
}

func TestGlobalCache(t *testing.T) {
	// 初始化全局缓存
	Init()

	// 测试全局函数
	err := Set("global_key", "global_value", time.Minute)
	if err != nil {
		t.Errorf("Set should not return error: %v", err)
	}

	value, err := Get("global_key")
	if err != nil {
		t.Errorf("Get should not return error: %v", err)
	}

	if value != "global_value" {
		t.Errorf("Expected 'global_value', got %v", value)
	}

	if !Has("global_key") {
		t.Error("Key should exist")
	}

	err = Delete("global_key")
	if err != nil {
		t.Errorf("Delete should not return error: %v", err)
	}

	if Has("global_key") {
		t.Error("Key should not exist after deletion")
	}
}

func TestCachePrefix(t *testing.T) {
	store := NewMemoryStore()

	// 设置前缀
	store.SetPrefix("app:")

	// 设置缓存
	err := store.Set("key", "value", time.Minute)
	if err != nil {
		t.Errorf("Set should not return error: %v", err)
	}

	// 检查前缀
	if store.GetPrefix() != "app:" {
		t.Errorf("Expected prefix 'app:', got %s", store.GetPrefix())
	}

	// 获取缓存
	value, err := store.Get("key")
	if err != nil {
		t.Errorf("Get should not return error: %v", err)
	}

	if value != "value" {
		t.Errorf("Expected 'value', got %v", value)
	}
}

func TestCacheRememberForever(t *testing.T) {
	store := NewMemoryStore()

	// 测试永久记住
	value, err := store.RememberForever("forever_key", func() (interface{}, error) {
		return "forever_value", nil
	})
	if err != nil {
		t.Errorf("RememberForever should not return error: %v", err)
	}

	if value != "forever_value" {
		t.Errorf("Expected 'forever_value', got %v", value)
	}

	// 验证缓存存在
	if !store.Has("forever_key") {
		t.Error("Key should exist")
	}
}

func TestCacheGetBytes(t *testing.T) {
	store := NewMemoryStore()

	// 测试字节数组
	testBytes := []byte("test bytes")
	err := store.SetBytes("bytes_key", testBytes, time.Minute)
	if err != nil {
		t.Errorf("SetBytes should not return error: %v", err)
	}

	bytes, err := store.GetBytes("bytes_key")
	if err != nil {
		t.Errorf("GetBytes should not return error: %v", err)
	}

	if string(bytes) != "test bytes" {
		t.Errorf("Expected 'test bytes', got %s", string(bytes))
	}

	// 测试复杂对象的JSON序列化
	complexObj := map[string]interface{}{
		"name": "test",
		"age":  25,
	}

	err = store.Set("complex_key", complexObj, time.Minute)
	if err != nil {
		t.Errorf("Set should not return error: %v", err)
	}

	bytes, err = store.GetBytes("complex_key")
	if err != nil {
		t.Errorf("GetBytes should not return error: %v", err)
	}

	if len(bytes) == 0 {
		t.Error("Bytes should not be empty")
	}
}
