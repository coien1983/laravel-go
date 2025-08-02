package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DatabaseStore 数据库缓存存储
type DatabaseStore struct {
	db     *sql.DB
	table  string
	prefix string
}

// DatabaseItem 数据库缓存项
type DatabaseItem struct {
	Key        string    `db:"key"`
	Value      string    `db:"value"`
	Expiration time.Time `db:"expiration"`
}

// NewDatabaseStore 创建新的数据库缓存存储
func NewDatabaseStore(db *sql.DB, table string) *DatabaseStore {
	store := &DatabaseStore{
		db:     db,
		table:  table,
		prefix: "",
	}

	// 确保缓存表存在
	store.createTable()

	return store
}

// createTable 创建缓存表
func (store *DatabaseStore) createTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			key VARCHAR(255) PRIMARY KEY,
			value TEXT,
			expiration TIMESTAMP NULL
		)
	`, store.table)

	_, err := store.db.Exec(query)
	return err
}

// Get 获取缓存值
func (store *DatabaseStore) Get(key string) (interface{}, error) {
	query := fmt.Sprintf(`
		SELECT value, expiration 
		FROM %s 
		WHERE key = ? AND (expiration IS NULL OR expiration > ?)
	`, store.table)

	var value string
	var expiration sql.NullTime

	err := store.db.QueryRow(query, store.prefix+key, time.Now()).Scan(&value, &expiration)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cache key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	// 检查是否过期
	if expiration.Valid && time.Now().After(expiration.Time) {
		// 删除过期项
		store.Delete(key)
		return nil, fmt.Errorf("cache key expired: %s", key)
	}

	// 尝试解析JSON
	var result interface{}
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		// 如果不是JSON，返回原始字符串
		return value, nil
	}

	return result, nil
}

// GetString 获取字符串缓存值
func (store *DatabaseStore) GetString(key string) (string, error) {
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
func (store *DatabaseStore) GetInt(key string) (int, error) {
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
func (store *DatabaseStore) GetFloat(key string) (float64, error) {
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
func (store *DatabaseStore) GetBool(key string) (bool, error) {
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
func (store *DatabaseStore) GetBytes(key string) ([]byte, error) {
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
func (store *DatabaseStore) Set(key string, value interface{}, ttl time.Duration) error {
	// 序列化为JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	var expiration sql.NullTime
	if ttl > 0 {
		expiration.Time = time.Now().Add(ttl)
		expiration.Valid = true
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (key, value, expiration) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE 
		value = VALUES(value), 
		expiration = VALUES(expiration)
	`, store.table)

	_, err = store.db.Exec(query, store.prefix+key, string(data), expiration)
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// SetString 设置字符串缓存值
func (store *DatabaseStore) SetString(key string, value string, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetInt 设置整数缓存值
func (store *DatabaseStore) SetInt(key string, value int, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetFloat 设置浮点数缓存值
func (store *DatabaseStore) SetFloat(key string, value float64, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetBool 设置布尔值缓存值
func (store *DatabaseStore) SetBool(key string, value bool, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetBytes 设置字节数组缓存值
func (store *DatabaseStore) SetBytes(key string, value []byte, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// Delete 删除缓存
func (store *DatabaseStore) Delete(key string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE key = ?", store.table)
	_, err := store.db.Exec(query, store.prefix+key)
	return err
}

// DeleteMultiple 批量删除缓存
func (store *DatabaseStore) DeleteMultiple(keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	// 构建批量删除查询
	placeholders := make([]string, len(keys))
	args := make([]interface{}, len(keys))

	for i, key := range keys {
		placeholders[i] = "?"
		args[i] = store.prefix + key
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE key IN (%s)", store.table,
		fmt.Sprintf(strings.Join(placeholders, ",")))

	_, err := store.db.Exec(query, args...)
	return err
}

// Clear 清空所有缓存
func (store *DatabaseStore) Clear() error {
	query := fmt.Sprintf("DELETE FROM %s WHERE key LIKE ?", store.table)
	_, err := store.db.Exec(query, store.prefix+"%")
	return err
}

// Has 检查缓存是否存在
func (store *DatabaseStore) Has(key string) bool {
	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE key = ? AND (expiration IS NULL OR expiration > ?)
	`, store.table)

	var count int
	err := store.db.QueryRow(query, store.prefix+key, time.Now()).Scan(&count)
	return err == nil && count > 0
}

// Missing 检查缓存是否不存在
func (store *DatabaseStore) Missing(key string) bool {
	return !store.Has(key)
}

// Increment 递增缓存值
func (store *DatabaseStore) Increment(key string, value int) (int, error) {
	// 使用事务确保原子性
	tx, err := store.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// 获取当前值
	var currentValue string
	query := fmt.Sprintf("SELECT value FROM %s WHERE key = ?", store.table)
	err = tx.QueryRow(query, store.prefix+key).Scan(&currentValue)

	var current int
	if err == sql.ErrNoRows {
		current = 0
	} else if err != nil {
		return 0, err
	} else {
		// 解析当前值
		var parsed interface{}
		if err := json.Unmarshal([]byte(currentValue), &parsed); err != nil {
			current = 0
		} else {
			switch v := parsed.(type) {
			case int:
				current = v
			case int64:
				current = int(v)
			case float64:
				current = int(v)
			case string:
				if parsed, err := strconv.Atoi(v); err == nil {
					current = parsed
				}
			default:
				current = 0
			}
		}
	}

	newValue := current + value

	// 设置新值
	data, err := json.Marshal(newValue)
	if err != nil {
		return 0, err
	}

	upsertQuery := fmt.Sprintf(`
		INSERT INTO %s (key, value) VALUES (?, ?) 
		ON DUPLICATE KEY UPDATE value = VALUES(value)
	`, store.table)

	_, err = tx.Exec(upsertQuery, store.prefix+key, string(data))
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newValue, nil
}

// Decrement 递减缓存值
func (store *DatabaseStore) Decrement(key string, value int) (int, error) {
	return store.Increment(key, -value)
}

// Remember 记住缓存值
func (store *DatabaseStore) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
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
func (store *DatabaseStore) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return store.Remember(key, 0, callback)
}

// Tags 获取标签管理器
func (store *DatabaseStore) Tags(names ...string) TaggedStore {
	return NewMemoryTaggedStore(store, names...)
}

// Flush 刷新缓存
func (store *DatabaseStore) Flush() error {
	return store.Clear()
}

// GetPrefix 获取缓存键前缀
func (store *DatabaseStore) GetPrefix() string {
	return store.prefix
}

// SetPrefix 设置缓存键前缀
func (store *DatabaseStore) SetPrefix(prefix string) {
	store.prefix = prefix
}

// CleanupExpired 清理过期缓存
func (store *DatabaseStore) CleanupExpired() error {
	query := fmt.Sprintf("DELETE FROM %s WHERE expiration IS NOT NULL AND expiration <= ?", store.table)
	_, err := store.db.Exec(query, time.Now())
	return err
}
