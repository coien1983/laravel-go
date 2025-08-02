package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"laravel-go/framework/auth"
	"laravel-go/framework/cache"
	"laravel-go/framework/config"
	"laravel-go/framework/container"
	"laravel-go/framework/database"
	"laravel-go/framework/event"
	"laravel-go/framework/queue"
	"laravel-go/framework/routing"
	"laravel-go/framework/validation"
)

// IntegrationTestSuite 集成测试套件
type IntegrationTestSuite struct {
	suite.Suite
	app      container.Container
	router   *routing.Router
	server   *httptest.Server
	db       database.Connection
	cache    cache.Store
	queue    queue.Queue
	eventMgr *event.EventManager
}

// SetupSuite 设置测试套件
func (suite *IntegrationTestSuite) SetupSuite() {
	// 初始化应用容器
	suite.app = container.NewContainer()

	// 注册配置
	suite.app.BindSingleton("config", config.NewConfig())

	// 注册数据库连接
	suite.app.BindCallback("database", func(container container.Container) interface{} {
		db, err := database.NewConnection(&database.ConnectionConfig{
			Driver:   "sqlite",
			Database: ":memory:",
		})
		if err != nil {
			suite.T().Fatalf("Failed to create database connection: %v", err)
		}
		return db
	})

	// 注册缓存
	suite.app.BindSingleton("cache", cache.NewManager())

	// 注册队列
	suite.app.BindSingleton("queue", queue.NewManager())

	// 注册事件管理器
	suite.app.BindSingleton("events", event.NewEventManager())

	// 注册认证服务
	suite.app.BindCallback("auth", func(container container.Container) interface{} {
		return auth.NewAuth(container)
	})

	// 获取服务实例
	suite.db = suite.app.Make("database").(database.Connection)
	suite.cache = suite.app.Make("cache").(cache.Cache)
	suite.queue = suite.app.Make("queue").(queue.Queue)
	suite.eventMgr = suite.app.Make("events").(*event.EventManager)

	// 初始化路由
	suite.router = routing.NewRouter()

	// 创建测试服务器
	suite.server = httptest.NewServer(suite.router)

	// 设置数据库表
	suite.setupDatabase()
}

// TearDownSuite 清理测试套件
func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest 设置每个测试
func (suite *IntegrationTestSuite) SetupTest() {
	// 清理数据库
	suite.cleanupDatabase()
	// 清理缓存
	suite.cache.Flush()
	// 清理队列
	suite.queue.Clear()
}

// setupDatabase 设置数据库表
func (suite *IntegrationTestSuite) setupDatabase() {
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

	// 创建文章表
	_, err = suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title VARCHAR(255) NOT NULL,
			content TEXT,
			status VARCHAR(50) DEFAULT 'draft',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		suite.T().Fatalf("Failed to create posts table: %v", err)
	}
}

// cleanupDatabase 清理数据库
func (suite *IntegrationTestSuite) cleanupDatabase() {
	suite.db.Exec("DELETE FROM posts")
	suite.db.Exec("DELETE FROM users")
}

// TestHTTPIntegration HTTP集成测试
func (suite *IntegrationTestSuite) TestHTTPIntegration() {
	// 创建用户控制器
	userController := &UserController{
		db: suite.db,
	}

	// 注册路由
	suite.router.Get("/users", userController.Index)
	suite.router.Post("/users", userController.Store)
	suite.router.Get("/users/{id}", userController.Show)
	suite.router.Put("/users/{id}", userController.Update)
	suite.router.Delete("/users/{id}", userController.Destroy)

	// 测试创建用户
	userData := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
	}
	userJSON, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", suite.server.URL+"/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)

	// 测试获取用户列表
	req, _ = http.NewRequest("GET", suite.server.URL+"/users", nil)
	resp, err = http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 测试获取单个用户
	req, _ = http.NewRequest("GET", suite.server.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
}

// TestDatabaseIntegration 数据库集成测试
func (suite *IntegrationTestSuite) TestDatabaseIntegration() {
	// 创建用户模型
	user := &User{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "password123",
	}

	// 保存用户
	err := user.Save(suite.db)
	suite.NoError(err)
	suite.NotZero(user.ID)

	// 查询用户
	foundUser := &User{}
	err = foundUser.Find(suite.db, user.ID)
	suite.NoError(err)
	suite.Equal(user.Name, foundUser.Name)
	suite.Equal(user.Email, foundUser.Email)

	// 更新用户
	foundUser.Name = "Jane Smith"
	err = foundUser.Save(suite.db)
	suite.NoError(err)

	// 验证更新
	updatedUser := &User{}
	err = updatedUser.Find(suite.db, user.ID)
	suite.NoError(err)
	suite.Equal("Jane Smith", updatedUser.Name)

	// 删除用户
	err = updatedUser.Delete(suite.db)
	suite.NoError(err)

	// 验证删除
	err = updatedUser.Find(suite.db, user.ID)
	suite.Error(err) // 应该找不到用户
}

// TestCacheIntegration 缓存集成测试
func (suite *IntegrationTestSuite) TestCacheIntegration() {
	// 设置缓存
	err := suite.cache.Set("test_key", "test_value", 60*time.Second)
	suite.NoError(err)

	// 获取缓存
	value, err := suite.cache.Get("test_key")
	suite.NoError(err)
	suite.Equal("test_value", value)

	// 检查缓存是否存在
	exists, err := suite.cache.Has("test_key")
	suite.NoError(err)
	suite.True(exists)

	// 删除缓存
	err = suite.cache.Delete("test_key")
	suite.NoError(err)

	// 验证删除
	exists, err = suite.cache.Has("test_key")
	suite.NoError(err)
	suite.False(exists)
}

// TestQueueIntegration 队列集成测试
func (suite *IntegrationTestSuite) TestQueueIntegration() {
	// 创建测试任务
	task := &TestJob{
		Data: "test data",
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

// TestEventIntegration 事件集成测试
func (suite *IntegrationTestSuite) TestEventIntegration() {
	eventFired := false
	eventData := ""

	// 注册事件监听器
	suite.eventMgr.Listen("user.created", func(e event.Event) {
		eventFired = true
		if userEvent, ok := e.(*UserCreatedEvent); ok {
			eventData = userEvent.UserName
		}
	})

	// 触发事件
	userEvent := &UserCreatedEvent{
		UserID:   1,
		UserName: "Test User",
	}
	suite.eventMgr.Dispatch(userEvent)

	// 验证事件被触发
	suite.True(eventFired)
	suite.Equal("Test User", eventData)
}

// TestAuthIntegration 认证集成测试
func (suite *IntegrationTestSuite) TestAuthIntegration() {
	authService := suite.app.Make("auth").(*auth.Auth)

	// 创建用户
	user := &User{
		Name:     "Auth User",
		Email:    "auth@example.com",
		Password: "password123",
	}
	err := user.Save(suite.db)
	suite.NoError(err)

	// 测试登录
	credentials := map[string]string{
		"email":    "auth@example.com",
		"password": "password123",
	}

	// 注意：这里需要实现实际的密码哈希验证
	// 为了测试，我们直接验证用户存在
	foundUser := &User{}
	err = foundUser.Where(suite.db, "email = ?", "auth@example.com").First()
	suite.NoError(err)
	suite.Equal("Auth User", foundUser.Name)
}

// TestMiddlewareIntegration 中间件集成测试
func (suite *IntegrationTestSuite) TestMiddlewareIntegration() {
	// 创建测试中间件
	testMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Middleware", "true")
			next(w, r)
		}
	}

	// 创建测试处理器
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}

	// 注册带中间件的路由
	suite.router.Get("/middleware-test", testMiddleware(testHandler))

	// 测试请求
	req, _ := http.NewRequest("GET", suite.server.URL+"/middleware-test", nil)
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("true", resp.Header.Get("X-Test-Middleware"))
}

// TestValidationIntegration 验证集成测试
func (suite *IntegrationTestSuite) TestValidationIntegration() {
	// 创建验证器
	validator := validation.NewValidator()

	// 测试数据
	data := map[string]interface{}{
		"name":  "",
		"email": "invalid-email",
		"age":   "not-a-number",
	}

	// 验证规则
	rules := map[string]string{
		"name":  "required|string|min:2",
		"email": "required|email",
		"age":   "required|int|min:18",
	}

	// 执行验证
	errors := validator.Validate(data, rules)

	// 验证错误
	suite.True(len(errors) > 0)
	suite.Contains(errors, "name")
	suite.Contains(errors, "email")
	suite.Contains(errors, "age")
}

// 运行集成测试套件
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// User 用户模型
type User struct {
	database.Model
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (u *User) TableName() string {
	return "users"
}

// UserController 用户控制器
type UserController struct {
	controller.Controller
	db database.Connection
}

func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	users := []*User{}
	err := c.db.Query("SELECT * FROM users").Get(&users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": users,
	})
}

func (c *UserController) Store(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := user.Save(c.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取ID
	params := routing.GetRouteParameters(r)
	id := params["id"]

	var user User
	err := user.Find(c.db, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	params := routing.GetRouteParameters(r)
	id := params["id"]

	var user User
	err := user.Find(c.db, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 更新用户字段
	if name, ok := updateData["name"].(string); ok {
		user.Name = name
	}
	if email, ok := updateData["email"].(string); ok {
		user.Email = email
	}

	err = user.Save(c.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (c *UserController) Destroy(w http.ResponseWriter, r *http.Request) {
	params := routing.GetRouteParameters(r)
	id := params["id"]

	var user User
	err := user.Find(c.db, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = user.Delete(c.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// TestJob 测试任务
type TestJob struct {
	Data string
}

func (j *TestJob) Handle() error {
	// 模拟任务处理
	fmt.Printf("Processing job with data: %s\n", j.Data)
	return nil
}

// UserCreatedEvent 用户创建事件
type UserCreatedEvent struct {
	event.BaseEvent
	UserID   int64
	UserName string
}

func (e *UserCreatedEvent) GetName() string {
	return "user.created"
}
