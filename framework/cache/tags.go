package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// TagSet 标签集合
type TagSet struct {
	store     Store
	names     []string
	namespace string
	mutex     sync.RWMutex
}

// NewMemoryTaggedStore 创建新的内存标签存储
func NewMemoryTaggedStore(store Store, names ...string) TaggedStore {
	// 生成命名空间
	namespace := generateTagNamespace(names)

	return &TagSet{
		store:     store,
		names:     names,
		namespace: namespace,
	}
}

// generateTagNamespace 生成标签命名空间
func generateTagNamespace(names []string) string {
	sort.Strings(names)
	combined := strings.Join(names, "|")
	hash := md5.Sum([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// Get 获取缓存值
func (ts *TagSet) Get(key string) (interface{}, error) {
	return ts.store.Get(ts.namespace + ":" + key)
}

// GetString 获取字符串缓存值
func (ts *TagSet) GetString(key string) (string, error) {
	return ts.store.GetString(ts.namespace + ":" + key)
}

// GetInt 获取整数缓存值
func (ts *TagSet) GetInt(key string) (int, error) {
	return ts.store.GetInt(ts.namespace + ":" + key)
}

// GetFloat 获取浮点数缓存值
func (ts *TagSet) GetFloat(key string) (float64, error) {
	return ts.store.GetFloat(ts.namespace + ":" + key)
}

// GetBool 获取布尔值缓存值
func (ts *TagSet) GetBool(key string) (bool, error) {
	return ts.store.GetBool(ts.namespace + ":" + key)
}

// GetBytes 获取字节数组缓存值
func (ts *TagSet) GetBytes(key string) ([]byte, error) {
	return ts.store.GetBytes(ts.namespace + ":" + key)
}

// Set 设置缓存值
func (ts *TagSet) Set(key string, value interface{}, ttl time.Duration) error {
	return ts.store.Set(ts.namespace+":"+key, value, ttl)
}

// SetString 设置字符串缓存值
func (ts *TagSet) SetString(key string, value string, ttl time.Duration) error {
	return ts.store.SetString(ts.namespace+":"+key, value, ttl)
}

// SetInt 设置整数缓存值
func (ts *TagSet) SetInt(key string, value int, ttl time.Duration) error {
	return ts.store.SetInt(ts.namespace+":"+key, value, ttl)
}

// SetFloat 设置浮点数缓存值
func (ts *TagSet) SetFloat(key string, value float64, ttl time.Duration) error {
	return ts.store.SetFloat(ts.namespace+":"+key, value, ttl)
}

// SetBool 设置布尔值缓存值
func (ts *TagSet) SetBool(key string, value bool, ttl time.Duration) error {
	return ts.store.SetBool(ts.namespace+":"+key, value, ttl)
}

// SetBytes 设置字节数组缓存值
func (ts *TagSet) SetBytes(key string, value []byte, ttl time.Duration) error {
	return ts.store.SetBytes(ts.namespace+":"+key, value, ttl)
}

// Delete 删除缓存
func (ts *TagSet) Delete(key string) error {
	return ts.store.Delete(ts.namespace + ":" + key)
}

// DeleteMultiple 批量删除缓存
func (ts *TagSet) DeleteMultiple(keys []string) error {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = ts.namespace + ":" + key
	}
	return ts.store.DeleteMultiple(prefixedKeys)
}

// Clear 清空所有缓存
func (ts *TagSet) Clear() error {
	// 对于标签存储，清空意味着刷新标签
	return ts.Flush()
}

// Has 检查缓存是否存在
func (ts *TagSet) Has(key string) bool {
	return ts.store.Has(ts.namespace + ":" + key)
}

// Missing 检查缓存是否不存在
func (ts *TagSet) Missing(key string) bool {
	return ts.store.Missing(ts.namespace + ":" + key)
}

// Increment 递增缓存值
func (ts *TagSet) Increment(key string, value int) (int, error) {
	return ts.store.Increment(ts.namespace+":"+key, value)
}

// Decrement 递减缓存值
func (ts *TagSet) Decrement(key string, value int) (int, error) {
	return ts.store.Decrement(ts.namespace+":"+key, value)
}

// Remember 记住缓存值
func (ts *TagSet) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	return ts.store.Remember(ts.namespace+":"+key, ttl, callback)
}

// RememberForever 永久记住缓存值
func (ts *TagSet) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return ts.store.RememberForever(ts.namespace+":"+key, callback)
}

// AddTags 添加标签
func (ts *TagSet) AddTags(names ...string) TaggedStore {
	allNames := append(ts.names, names...)
	return NewMemoryTaggedStore(ts.store, allNames...)
}

// Flush 刷新标签下的所有缓存
func (ts *TagSet) Flush() error {
	// 更新标签版本号来使所有相关缓存失效
	versionKey := fmt.Sprintf("tag_version:%s", strings.Join(ts.names, "|"))
	currentVersion, _ := ts.store.GetInt(versionKey)
	ts.store.SetInt(versionKey, currentVersion+1, 0)
	return nil
}

// GetTags 获取标签列表
func (ts *TagSet) GetTags() []string {
	return ts.names
}

// Tags 获取标签管理器（兼容Store接口）
func (ts *TagSet) Tags(names ...string) TaggedStore {
	return ts.AddTags(names...)
}

// GetPrefix 获取缓存键前缀
func (ts *TagSet) GetPrefix() string {
	return ts.namespace + ":"
}

// SetPrefix 设置缓存键前缀
func (ts *TagSet) SetPrefix(prefix string) {
	// 标签存储不支持直接设置前缀
}

// TagManager 标签管理器
type TagManager struct {
	store Store
}

// NewTagManager 创建新的标签管理器
func NewTagManager(store Store) *TagManager {
	return &TagManager{
		store: store,
	}
}

// Flush 刷新指定标签的所有缓存
func (tm *TagManager) Flush(names ...string) error {
	if len(names) == 0 {
		return fmt.Errorf("no tags specified")
	}

	sort.Strings(names)
	versionKey := fmt.Sprintf("tag_version:%s", strings.Join(names, "|"))
	currentVersion, _ := tm.store.GetInt(versionKey)
	return tm.store.SetInt(versionKey, currentVersion+1, 0)
}

// GetVersion 获取标签版本号
func (tm *TagManager) GetVersion(names ...string) int {
	if len(names) == 0 {
		return 0
	}

	sort.Strings(names)
	versionKey := fmt.Sprintf("tag_version:%s", strings.Join(names, "|"))
	version, _ := tm.store.GetInt(versionKey)
	return version
}
