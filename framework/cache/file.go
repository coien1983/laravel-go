package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// FileItem 文件缓存项
type FileItem struct {
	Value      interface{} `json:"value"`
	Expiration int64       `json:"expiration"`
}

// IsExpired 检查是否过期
func (item *FileItem) IsExpired() bool {
	return item.Expiration > 0 && time.Now().Unix() >= item.Expiration
}

// FileStore 文件缓存存储
type FileStore struct {
	directory string
	prefix    string
}

// NewFileStore 创建新的文件缓存存储
func NewFileStore(directory string) *FileStore {
	// 确保目录存在
	os.MkdirAll(directory, 0755)

	return &FileStore{
		directory: directory,
		prefix:    "",
	}
}

// getFilePath 获取缓存文件路径
func (store *FileStore) getFilePath(key string) string {
	return filepath.Join(store.directory, store.prefix+key+".cache")
}

// Get 获取缓存值
func (store *FileStore) Get(key string) (interface{}, error) {
	filePath := store.getFilePath(key)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cache key not found: %s", key)
	}

	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	// 解析JSON
	var item FileItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	// 检查是否过期
	if item.IsExpired() {
		// 删除过期文件
		os.Remove(filePath)
		return nil, fmt.Errorf("cache key expired: %s", key)
	}

	return item.Value, nil
}

// GetString 获取字符串缓存值
func (store *FileStore) GetString(key string) (string, error) {
	value, err := store.Get(key)
	if err != nil {
		return "", err
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// GetInt 获取整数缓存值
func (store *FileStore) GetInt(key string) (int, error) {
	value, err := store.Get(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %v to int", value)
	}
}

// GetFloat 获取浮点数缓存值
func (store *FileStore) GetFloat(key string) (float64, error) {
	value, err := store.Get(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to float64", value)
	}
}

// GetBool 获取布尔值缓存值
func (store *FileStore) GetBool(key string) (bool, error) {
	value, err := store.Get(key)
	if err != nil {
		return false, err
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("cannot convert %v to bool", value)
	}
}

// GetBytes 获取字节数组缓存值
func (store *FileStore) GetBytes(key string) ([]byte, error) {
	value, err := store.Get(key)
	if err != nil {
		return nil, err
	}

	switch v := value.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		// 尝试JSON序列化
		return json.Marshal(v)
	}
}

// Set 设置缓存值
func (store *FileStore) Set(key string, value interface{}, ttl time.Duration) error {
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).Unix()
	}

	item := FileItem{
		Value:      value,
		Expiration: expiration,
	}

	// 序列化为JSON
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	// 写入文件
	filePath := store.getFilePath(key)
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// SetString 设置字符串缓存值
func (store *FileStore) SetString(key string, value string, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetInt 设置整数缓存值
func (store *FileStore) SetInt(key string, value int, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetFloat 设置浮点数缓存值
func (store *FileStore) SetFloat(key string, value float64, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetBool 设置布尔值缓存值
func (store *FileStore) SetBool(key string, value bool, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetBytes 设置字节数组缓存值
func (store *FileStore) SetBytes(key string, value []byte, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// Delete 删除缓存
func (store *FileStore) Delete(key string) error {
	filePath := store.getFilePath(key)
	return os.Remove(filePath)
}

// DeleteMultiple 批量删除缓存
func (store *FileStore) DeleteMultiple(keys []string) error {
	for _, key := range keys {
		if err := store.Delete(key); err != nil {
			// 继续删除其他键，不中断
			continue
		}
	}
	return nil
}

// Clear 清空所有缓存
func (store *FileStore) Clear() error {
	// 删除目录下的所有.cache文件
	files, err := filepath.Glob(filepath.Join(store.directory, "*.cache"))
	if err != nil {
		return fmt.Errorf("failed to list cache files: %w", err)
	}

	for _, file := range files {
		os.Remove(file)
	}

	return nil
}

// Has 检查缓存是否存在
func (store *FileStore) Has(key string) bool {
	filePath := store.getFilePath(key)
	
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	
	// 读取文件内容检查是否过期
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false
	}
	
	// 解析JSON
	var item FileItem
	if err := json.Unmarshal(data, &item); err != nil {
		return false
	}
	
	// 检查是否过期
	return !item.IsExpired()
}

// Missing 检查缓存是否不存在
func (store *FileStore) Missing(key string) bool {
	return !store.Has(key)
}

// Increment 递增缓存值
func (store *FileStore) Increment(key string, value int) (int, error) {
	current, err := store.GetInt(key)
	if err != nil {
		current = 0
	}

	newValue := current + value
	if err := store.SetInt(key, newValue, 0); err != nil {
		return 0, err
	}

	return newValue, nil
}

// Decrement 递减缓存值
func (store *FileStore) Decrement(key string, value int) (int, error) {
	return store.Increment(key, -value)
}

// Remember 记住缓存值
func (store *FileStore) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// 先尝试获取缓存
	if value, err := store.Get(key); err == nil {
		return value, nil
	}

	// 缓存不存在，执行回调函数
	value, err := callback()
	if err != nil {
		return nil, err
	}

	// 设置缓存
	if err := store.Set(key, value, ttl); err != nil {
		return nil, err
	}

	return value, nil
}

// RememberForever 永久记住缓存值
func (store *FileStore) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return store.Remember(key, 0, callback)
}

// Tags 获取标签管理器
func (store *FileStore) Tags(names ...string) TaggedStore {
	return NewMemoryTaggedStore(store, names...)
}

// Flush 刷新缓存
func (store *FileStore) Flush() error {
	return store.Clear()
}

// GetPrefix 获取缓存键前缀
func (store *FileStore) GetPrefix() string {
	return store.prefix
}

// SetPrefix 设置缓存键前缀
func (store *FileStore) SetPrefix(prefix string) {
	store.prefix = prefix
}
