package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"laravel-go/framework/cache"
	"laravel-go/framework/config"
	"laravel-go/framework/container"
	"laravel-go/framework/database"
	"laravel-go/framework/queue"
)

// CacheQueueIntegrationTestSuite 缓存和队列集成测试套件
type CacheQueueIntegrationTestSuite struct {
	suite.Suite
	app   *container.Container
	db    database.Connection
	cache cache.Cache
	queue queue.Queue
}

// SetupSuite 设置测试套件
func (suite *CacheQueueIntegrationTestSuite) SetupSuite() {
	// 初始化应用容器
	suite.app = container.NewContainer()

	// 注册配置
	suite.app.Singleton("config", func() interface{} {
		return config.NewConfig()
	})

	// 注册数据库连接
	suite.app.Singleton("database", func() interface{} {
		db, err := database.NewConnection(&database.Config{
			Driver:   "sqlite",
			Database: ":memory:",
		})
		if err != nil {
			suite.T().Fatalf("Failed to create database connection: %v", err)
		}
		return db
	})

	// 注册缓存
	suite.app.Singleton("cache", func() interface{} {
		return cache.NewCache("memory", &cache.Config{})
	})

	// 注册队列
	suite.app.Singleton("queue", func() interface{} {
		return queue.NewQueue("memory", &queue.Config{})
	})

	// 获取服务实例
	suite.db = suite.app.Make("database").(database.Connection)
	suite.cache = suite.app.Make("cache").(cache.Cache)
	suite.queue = suite.app.Make("queue").(queue.Queue)

	// 设置数据库表
	suite.setupDatabase()
}

// TearDownSuite 清理测试套件
func (suite *CacheQueueIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest 设置每个测试
func (suite *CacheQueueIntegrationTestSuite) SetupTest() {
	// 清理数据库
	suite.cleanupDatabase()
	// 清理缓存
	suite.cache.Flush()
	// 清理队列
	suite.queue.Clear()
}

// setupDatabase 设置数据库表
func (suite *CacheQueueIntegrationTestSuite) setupDatabase() {
	// 创建用户表
	_, err := suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		suite.T().Fatalf("Failed to create users table: %v", err)
	}

	// 创建缓存统计表
	_, err = suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS cache_stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key_name VARCHAR(255) NOT NULL,
			hits INTEGER DEFAULT 0,
			misses INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		suite.T().Fatalf("Failed to create cache_stats table: %v", err)
	}
}

// cleanupDatabase 清理数据库
func (suite *CacheQueueIntegrationTestSuite) cleanupDatabase() {
	suite.db.Exec("DELETE FROM cache_stats")
	suite.db.Exec("DELETE FROM users")
}

// TestCacheBasicOperations 测试缓存基础操作
func (suite *CacheQueueIntegrationTestSuite) TestCacheBasicOperations() {
	// 测试设置缓存
	err := suite.cache.Set("test_key", "test_value", 60*time.Second)
	suite.NoError(err)

	// 测试获取缓存
	value, err := suite.cache.Get("test_key")
	suite.NoError(err)
	suite.Equal("test_value", value)

	// 测试检查缓存是否存在
	exists, err := suite.cache.Has("test_key")
	suite.NoError(err)
	suite.True(exists)

	// 测试删除缓存
	err = suite.cache.Delete("test_key")
	suite.NoError(err)

	// 验证删除
	exists, err = suite.cache.Has("test_key")
	suite.NoError(err)
	suite.False(exists)
}

// TestCacheExpiration 测试缓存过期
func (suite *CacheQueueIntegrationTestSuite) TestCacheExpiration() {
	// 设置短期缓存
	err := suite.cache.Set("expire_key", "expire_value", 100*time.Millisecond)
	suite.NoError(err)

	// 立即获取应该存在
	value, err := suite.cache.Get("expire_key")
	suite.NoError(err)
	suite.Equal("expire_value", value)

	// 等待过期
	time.Sleep(200 * time.Millisecond)

	// 过期后应该不存在
	value, err = suite.cache.Get("expire_key")
	suite.Error(err) // 应该返回错误
	suite.Empty(value)
}

// TestCacheTags 测试缓存标签
func (suite *CacheQueueIntegrationTestSuite) TestCacheTags() {
	// 设置带标签的缓存
	err := suite.cache.Tags("users", "profile").Set("user:1", "user_data", 60*time.Second)
	suite.NoError(err)

	err = suite.cache.Tags("users").Set("user:2", "user_data_2", 60*time.Second)
	suite.NoError(err)

	// 获取带标签的缓存
	value, err := suite.cache.Tags("users", "profile").Get("user:1")
	suite.NoError(err)
	suite.Equal("user_data", value)

	// 刷新标签
	err = suite.cache.Tags("users").Flush()
	suite.NoError(err)

	// 验证标签内的缓存被清除
	value, err = suite.cache.Tags("users", "profile").Get("user:1")
	suite.Error(err)
	suite.Empty(value)

	value, err = suite.cache.Tags("users").Get("user:2")
	suite.Error(err)
	suite.Empty(value)
}

// TestCacheBatchOperations 测试缓存批量操作
func (suite *CacheQueueIntegrationTestSuite) TestCacheBatchOperations() {
	// 批量设置
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	err := suite.cache.SetMany(data, 60*time.Second)
	suite.NoError(err)

	// 批量获取
	keys := []string{"key1", "key2", "key3"}
	values, err := suite.cache.GetMany(keys)
	suite.NoError(err)
	suite.Len(values, 3)
	suite.Equal("value1", values["key1"])
	suite.Equal("value2", values["key2"])
	suite.Equal("value3", values["key3"])

	// 批量删除
	err = suite.cache.DeleteMany(keys)
	suite.NoError(err)

	// 验证删除
	for _, key := range keys {
		exists, err := suite.cache.Has(key)
		suite.NoError(err)
		suite.False(exists)
	}
}

// TestCacheWithDatabase 测试缓存与数据库集成
func (suite *CacheQueueIntegrationTestSuite) TestCacheWithDatabase() {
	// 创建用户
	user := &User{
		Name:     "Cache User",
		Email:    "cache@example.com",
		Password: "password",
	}
	err := user.Save(suite.db)
	suite.NoError(err)

	// 缓存用户数据
	userData := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}
	userJSON, _ := json.Marshal(userData)

	err = suite.cache.Set(fmt.Sprintf("user:%d", user.ID), string(userJSON), 60*time.Second)
	suite.NoError(err)

	// 从缓存获取用户数据
	cachedData, err := suite.cache.Get(fmt.Sprintf("user:%d", user.ID))
	suite.NoError(err)

	var cachedUser map[string]interface{}
	err = json.Unmarshal([]byte(cachedData), &cachedUser)
	suite.NoError(err)
	suite.Equal(user.Name, cachedUser["name"])
	suite.Equal(user.Email, cachedUser["email"])

	// 更新数据库中的用户
	user.Name = "Updated Cache User"
	err = user.Save(suite.db)
	suite.NoError(err)

	// 清除缓存
	err = suite.cache.Delete(fmt.Sprintf("user:%d", user.ID))
	suite.NoError(err)

	// 重新从数据库加载并缓存
	foundUser := &User{}
	err = foundUser.Find(suite.db, user.ID)
	suite.NoError(err)
	suite.Equal("Updated Cache User", foundUser.Name)
}

// TestQueueBasicOperations 测试队列基础操作
func (suite *CacheQueueIntegrationTestSuite) TestQueueBasicOperations() {
	// 创建测试任务
	task := &TestJob{
		Data: "test queue data",
	}

	// 推送任务到队列
	err := suite.queue.Push(task)
	suite.NoError(err)

	// 检查队列长度
	length, err := suite.queue.Size()
	suite.NoError(err)
	suite.Equal(1, length)

	// 处理任务
	job, err := suite.queue.Pop()
	suite.NoError(err)
	suite.NotNil(job)

	// 执行任务
	err = job.Handle()
	suite.NoError(err)

	// 验证队列为空
	length, err = suite.queue.Size()
	suite.NoError(err)
	suite.Equal(0, length)
}

// TestQueueDelayedJobs 测试延迟任务
func (suite *CacheQueueIntegrationTestSuite) TestQueueDelayedJobs() {
	// 创建延迟任务
	task := &TestJob{
		Data: "delayed job",
	}

	// 推送延迟任务（延迟100毫秒）
	err := suite.queue.Later(task, 100*time.Millisecond)
	suite.NoError(err)

	// 立即检查队列长度
	length, err := suite.queue.Size()
	suite.NoError(err)
	suite.Equal(1, length)

	// 立即尝试获取任务应该失败
	job, err := suite.queue.Pop()
	suite.Error(err) // 应该返回错误，因为任务还在延迟中
	suite.Nil(job)

	// 等待延迟时间
	time.Sleep(150 * time.Millisecond)

	// 现在应该可以获取任务
	job, err = suite.queue.Pop()
	suite.NoError(err)
	suite.NotNil(job)

	// 执行任务
	err = job.Handle()
	suite.NoError(err)
}

// TestQueueBatchOperations 测试队列批量操作
func (suite *CacheQueueIntegrationTestSuite) TestQueueBatchOperations() {
	// 创建多个任务
	tasks := []queue.Job{
		&TestJob{Data: "batch job 1"},
		&TestJob{Data: "batch job 2"},
		&TestJob{Data: "batch job 3"},
	}

	// 批量推送任务
	err := suite.queue.PushMany(tasks)
	suite.NoError(err)

	// 检查队列长度
	length, err := suite.queue.Size()
	suite.NoError(err)
	suite.Equal(3, length)

	// 批量处理任务
	for i := 0; i < 3; i++ {
		job, err := suite.queue.Pop()
		suite.NoError(err)
		suite.NotNil(job)

		err = job.Handle()
		suite.NoError(err)
	}

	// 验证队列为空
	length, err = suite.queue.Size()
	suite.NoError(err)
	suite.Equal(0, length)
}

// TestQueueWithDatabase 测试队列与数据库集成
func (suite *CacheQueueIntegrationTestSuite) TestQueueWithDatabase() {
	// 创建用户
	user := &User{
		Name:     "Queue User",
		Email:    "queue@example.com",
		Password: "password",
	}
	err := user.Save(suite.db)
	suite.NoError(err)

	// 创建数据库任务
	dbTask := &DatabaseJob{
		UserID: user.ID,
		Action: "update_profile",
		Data:   map[string]interface{}{"status": "processed"},
	}

	// 推送任务到队列
	err = suite.queue.Push(dbTask)
	suite.NoError(err)

	// 处理任务
	job, err := suite.queue.Pop()
	suite.NoError(err)
	suite.NotNil(job)

	// 执行任务
	err = job.Handle()
	suite.NoError(err)

	// 验证任务执行结果（这里可以检查数据库中的变化）
	// 在实际应用中，任务可能会更新数据库记录
}

// TestCacheQueueIntegration 测试缓存和队列集成
func (suite *CacheQueueIntegrationTestSuite) TestCacheQueueIntegration() {
	// 创建用户
	user := &User{
		Name:     "Integration User",
		Email:    "integration@example.com",
		Password: "password",
	}
	err := user.Save(suite.db)
	suite.NoError(err)

	// 缓存用户数据
	err = suite.cache.Set(fmt.Sprintf("user:%d", user.ID), user.Name, 60*time.Second)
	suite.NoError(err)

	// 创建更新用户的任务
	updateTask := &UserUpdateJob{
		UserID:  user.ID,
		NewName: "Updated Integration User",
	}

	// 推送任务到队列
	err = suite.queue.Push(updateTask)
	suite.NoError(err)

	// 处理任务
	job, err := suite.queue.Pop()
	suite.NoError(err)
	suite.NotNil(job)

	// 执行任务
	err = job.Handle()
	suite.NoError(err)

	// 验证数据库更新
	foundUser := &User{}
	err = foundUser.Find(suite.db, user.ID)
	suite.NoError(err)
	suite.Equal("Updated Integration User", foundUser.Name)

	// 验证缓存被清除（任务执行后应该清除相关缓存）
	cachedName, err := suite.cache.Get(fmt.Sprintf("user:%d", user.ID))
	suite.Error(err) // 缓存应该被清除
	suite.Empty(cachedName)
}

// TestCacheStatistics 测试缓存统计
func (suite *CacheQueueIntegrationTestSuite) TestCacheStatistics() {
	// 设置一些缓存
	suite.cache.Set("stat1", "value1", 60*time.Second)
	suite.cache.Set("stat2", "value2", 60*time.Second)
	suite.cache.Set("stat3", "value3", 60*time.Second)

	// 获取缓存
	suite.cache.Get("stat1")       // hit
	suite.cache.Get("stat2")       // hit
	suite.cache.Get("nonexistent") // miss

	// 获取统计信息
	stats := suite.cache.GetStats()
	suite.NotNil(stats)
	suite.GreaterOrEqual(stats.Hits, 2)
	suite.GreaterOrEqual(stats.Misses, 1)
}

// TestQueueStatistics 测试队列统计
func (suite *CacheQueueIntegrationTestSuite) TestQueueStatistics() {
	// 推送一些任务
	suite.queue.Push(&TestJob{Data: "stat job 1"})
	suite.queue.Push(&TestJob{Data: "stat job 2"})
	suite.queue.Push(&TestJob{Data: "stat job 3"})

	// 处理一些任务
	job1, _ := suite.queue.Pop()
	job1.Handle()
	job2, _ := suite.queue.Pop()
	job2.Handle()

	// 获取统计信息
	stats := suite.queue.GetStats()
	suite.NotNil(stats)
	suite.GreaterOrEqual(stats.Pushed, 3)
	suite.GreaterOrEqual(stats.Processed, 2)
	suite.Equal(1, stats.Size) // 还剩1个任务
}

// 运行缓存队列集成测试套件
func TestCacheQueueIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CacheQueueIntegrationTestSuite))
}

// DatabaseJob 数据库任务
type DatabaseJob struct {
	UserID int64                  `json:"user_id"`
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

func (j *DatabaseJob) Handle() error {
	// 模拟数据库操作
	fmt.Printf("Processing database job: %s for user %d\n", j.Action, j.UserID)
	return nil
}

// UserUpdateJob 用户更新任务
type UserUpdateJob struct {
	UserID  int64  `json:"user_id"`
	NewName string `json:"new_name"`
}

func (j *UserUpdateJob) Handle() error {
	// 这里应该更新数据库中的用户
	fmt.Printf("Updating user %d name to: %s\n", j.UserID, j.NewName)
	return nil
}
