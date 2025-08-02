package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoStore MongoDB缓存存储
type MongoStore struct {
	client     *mongo.Client
	database   string
	collection string
	prefix     string
}

// MongoItem MongoDB缓存项
type MongoItem struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Key        string             `bson:"key"`
	Value      interface{}        `bson:"value"`
	Expiration *time.Time         `bson:"expiration,omitempty"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

// NewMongoStore 创建新的MongoDB缓存存储
func NewMongoStore(client *mongo.Client, database, collection string) *MongoStore {
	store := &MongoStore{
		client:     client,
		database:   database,
		collection: collection,
		prefix:     "",
	}

	// 创建索引
	store.createIndexes()

	return store
}

// createIndexes 创建索引
func (store *MongoStore) createIndexes() error {
	ctx := context.Background()
	coll := store.client.Database(store.database).Collection(store.collection)

	// 创建key的唯一索引
	keyIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "key", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// 创建过期时间的TTL索引
	expirationIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "expiration", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := coll.Indexes().CreateMany(ctx, []mongo.IndexModel{keyIndex, expirationIndex})
	return err
}

// Get 获取缓存值
func (store *MongoStore) Get(key string) (interface{}, error) {
	ctx := context.Background()
	coll := store.client.Database(store.database).Collection(store.collection)

	var item MongoItem
	err := coll.FindOne(ctx, bson.M{
		"key": store.prefix + key,
		"$or": []bson.M{
			{"expiration": bson.M{"$exists": false}},
			{"expiration": bson.M{"$gt": time.Now()}},
		},
	}).Decode(&item)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("cache key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	// 检查是否过期
	if item.Expiration != nil && time.Now().After(*item.Expiration) {
		// 删除过期项
		store.Delete(key)
		return nil, fmt.Errorf("cache key expired: %s", key)
	}

	return item.Value, nil
}

// GetString 获取字符串缓存值
func (store *MongoStore) GetString(key string) (string, error) {
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
func (store *MongoStore) GetInt(key string) (int, error) {
	value, err := store.Get(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
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
func (store *MongoStore) GetFloat(key string) (float64, error) {
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
func (store *MongoStore) GetBool(key string) (bool, error) {
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
func (store *MongoStore) GetBytes(key string) ([]byte, error) {
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
func (store *MongoStore) Set(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	coll := store.client.Database(store.database).Collection(store.collection)

	now := time.Now()
	var expiration *time.Time
	if ttl > 0 {
		exp := now.Add(ttl)
		expiration = &exp
	}

	item := MongoItem{
		Key:        store.prefix + key,
		Value:      value,
		Expiration: expiration,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 使用upsert操作
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"key": store.prefix + key}
	_, err := coll.ReplaceOne(ctx, filter, item, opts)

	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// SetString 设置字符串缓存值
func (store *MongoStore) SetString(key string, value string, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetInt 设置整数缓存值
func (store *MongoStore) SetInt(key string, value int, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetFloat 设置浮点数缓存值
func (store *MongoStore) SetFloat(key string, value float64, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetBool 设置布尔值缓存值
func (store *MongoStore) SetBool(key string, value bool, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// SetBytes 设置字节数组缓存值
func (store *MongoStore) SetBytes(key string, value []byte, ttl time.Duration) error {
	return store.Set(key, value, ttl)
}

// Delete 删除缓存
func (store *MongoStore) Delete(key string) error {
	ctx := context.Background()
	coll := store.client.Database(store.database).Collection(store.collection)

	_, err := coll.DeleteOne(ctx, bson.M{"key": store.prefix + key})
	return err
}

// DeleteMultiple 批量删除缓存
func (store *MongoStore) DeleteMultiple(keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	ctx := context.Background()
	coll := store.client.Database(store.database).Collection(store.collection)

	// 构建key列表
	keyList := make([]string, len(keys))
	for i, key := range keys {
		keyList[i] = store.prefix + key
	}

	filter := bson.M{"key": bson.M{"$in": keyList}}
	_, err := coll.DeleteMany(ctx, filter)
	return err
}

// Clear 清空所有缓存
func (store *MongoStore) Clear() error {
	ctx := context.Background()
	coll := store.client.Database(store.database).Collection(store.collection)

	// 删除所有以prefix开头的key
	filter := bson.M{"key": bson.M{"$regex": "^" + store.prefix}}
	_, err := coll.DeleteMany(ctx, filter)
	return err
}

// Has 检查缓存是否存在
func (store *MongoStore) Has(key string) bool {
	ctx := context.Background()
	coll := store.client.Database(store.database).Collection(store.collection)

	count, err := coll.CountDocuments(ctx, bson.M{
		"key": store.prefix + key,
		"$or": []bson.M{
			{"expiration": bson.M{"$exists": false}},
			{"expiration": bson.M{"$gt": time.Now()}},
		},
	})

	return err == nil && count > 0
}

// Missing 检查缓存是否不存在
func (store *MongoStore) Missing(key string) bool {
	return !store.Has(key)
}

// Increment 递增缓存值
func (store *MongoStore) Increment(key string, value int) (int, error) {
	ctx := context.Background()
	coll := store.client.Database(store.database).Collection(store.collection)

	// 获取当前值
	var item MongoItem
	err := coll.FindOne(ctx, bson.M{"key": store.prefix + key}).Decode(&item)

	var current int
	if err == mongo.ErrNoDocuments {
		current = 0
	} else if err != nil {
		return 0, err
	} else {
		// 解析当前值
		switch v := item.Value.(type) {
		case int:
			current = v
		case int32:
			current = int(v)
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

	newValue := current + value

	// 设置新值
	now := time.Now()
	upsertItem := MongoItem{
		Key:       store.prefix + key,
		Value:     newValue,
		CreatedAt: now,
		UpdatedAt: now,
	}

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"key": store.prefix + key}
	_, err = coll.ReplaceOne(ctx, filter, upsertItem, opts)
	if err != nil {
		return 0, err
	}

	return newValue, nil
}

// Decrement 递减缓存值
func (store *MongoStore) Decrement(key string, value int) (int, error) {
	return store.Increment(key, -value)
}

// Remember 记住缓存值
func (store *MongoStore) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
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
func (store *MongoStore) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return store.Remember(key, 0, callback)
}

// Tags 获取标签管理器
func (store *MongoStore) Tags(names ...string) TaggedStore {
	return NewMemoryTaggedStore(store, names...)
}

// Flush 刷新缓存
func (store *MongoStore) Flush() error {
	return store.Clear()
}

// GetPrefix 获取缓存键前缀
func (store *MongoStore) GetPrefix() string {
	return store.prefix
}

// SetPrefix 设置缓存键前缀
func (store *MongoStore) SetPrefix(prefix string) {
	store.prefix = prefix
}

// CleanupExpired 清理过期缓存（MongoDB会自动清理TTL索引）
func (store *MongoStore) CleanupExpired() error {
	// MongoDB的TTL索引会自动清理过期文档
	// 这里可以添加额外的清理逻辑
	return nil
}

// GetStats 获取缓存统计信息
func (store *MongoStore) GetStats() (map[string]interface{}, error) {
	ctx := context.Background()
	coll := store.client.Database(store.database).Collection(store.collection)

	// 总文档数
	total, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// 未过期文档数
	valid, err := coll.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"expiration": bson.M{"$exists": false}},
			{"expiration": bson.M{"$gt": time.Now()}},
		},
	})
	if err != nil {
		return nil, err
	}

	// 过期文档数
	expired, err := coll.CountDocuments(ctx, bson.M{
		"expiration": bson.M{"$lte": time.Now()},
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":   total,
		"valid":   valid,
		"expired": expired,
	}, nil
}
