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

	"laravel-go/framework/database"
	"laravel-go/framework/routing"
)

// SimpleIntegrationTestSuite 简化集成测试套件
type SimpleIntegrationTestSuite struct {
	suite.Suite
	router routing.Router
	server *httptest.Server
	db     database.Connection
}

// SetupSuite 设置测试套件
func (suite *SimpleIntegrationTestSuite) SetupSuite() {
	// 初始化路由
	suite.router = routing.NewRouter()

	// 创建测试服务器
	suite.server = httptest.NewServer(suite.router)

	// 设置数据库
	suite.setupDatabase()
}

// TearDownSuite 清理测试套件
func (suite *SimpleIntegrationTestSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest 设置每个测试
func (suite *SimpleIntegrationTestSuite) SetupTest() {
	// 清理数据库
	suite.cleanupDatabase()
}

// setupDatabase 设置数据库
func (suite *SimpleIntegrationTestSuite) setupDatabase() {
	var err error
	suite.db, err = database.NewConnection(&database.ConnectionConfig{
		Driver:   "sqlite",
		Database: ":memory:",
	})
	if err != nil {
		suite.T().Fatalf("Failed to create database connection: %v", err)
	}

	// 创建测试表
	_, err = suite.db.Exec(`
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
}

// cleanupDatabase 清理数据库
func (suite *SimpleIntegrationTestSuite) cleanupDatabase() {
	if suite.db != nil {
		suite.db.Exec("DELETE FROM users")
	}
}

// TestBasicRouting 测试基础路由
func (suite *SimpleIntegrationTestSuite) TestBasicRouting() {
	// 注册测试路由
	suite.router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// 测试GET请求
	req, _ := http.NewRequest("GET", suite.server.URL+"/test", nil)
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 读取响应体
	var body bytes.Buffer
	body.ReadFrom(resp.Body)
	suite.Equal("test response", body.String())
}

// TestRouteParameters 测试路由参数
func (suite *SimpleIntegrationTestSuite) TestRouteParameters() {
	// 注册带参数的路由
	suite.router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := routing.GetRouteParameters(r)
		id := params["id"]

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("user id: %s", id)))
	})

	// 测试参数路由
	req, _ := http.NewRequest("GET", suite.server.URL+"/users/123", nil)
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var body bytes.Buffer
	body.ReadFrom(resp.Body)
	suite.Equal("user id: 123", body.String())
}

// TestRouteGroup 测试路由分组
func (suite *SimpleIntegrationTestSuite) TestRouteGroup() {
	// 创建API分组
	suite.router.Group("/api", func(r routing.Router) {
		r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("api users"))
		})

		r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("user created"))
		})
	})

	// 测试分组路由
	req, _ := http.NewRequest("GET", suite.server.URL+"/api/users", nil)
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var body bytes.Buffer
	body.ReadFrom(resp.Body)
	suite.Equal("api users", body.String())

	// 测试POST请求
	req, _ = http.NewRequest("POST", suite.server.URL+"/api/users", nil)
	resp, err = http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)
}

// TestDatabaseOperations 测试数据库操作
func (suite *SimpleIntegrationTestSuite) TestDatabaseOperations() {
	// 插入测试数据
	_, err := suite.db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		"John Doe", "john@example.com", "password123")
	suite.NoError(err)

	// 查询数据
	var count int
	err = suite.db.Query("SELECT COUNT(*) FROM users").Scan(&count)
	suite.NoError(err)
	suite.Equal(1, count)

	// 查询用户信息
	var name, email string
	err = suite.db.Query("SELECT name, email FROM users WHERE email = ?", "john@example.com").Scan(&name, &email)
	suite.NoError(err)
	suite.Equal("John Doe", name)
	suite.Equal("john@example.com", email)
}

// TestJSONResponse 测试JSON响应
func (suite *SimpleIntegrationTestSuite) TestJSONResponse() {
	// 注册JSON响应路由
	suite.router.Get("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Hello World",
			"status":  "success",
		})
	})

	// 测试JSON响应
	req, _ := http.NewRequest("GET", suite.server.URL+"/json", nil)
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("application/json", resp.Header.Get("Content-Type"))

	// 解析JSON响应
	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)
	suite.Equal("Hello World", response["message"])
	suite.Equal("success", response["status"])
}

// TestErrorHandling 测试错误处理
func (suite *SimpleIntegrationTestSuite) TestErrorHandling() {
	// 注册错误处理路由
	suite.router.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	suite.router.Get("/not-found", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	// 测试500错误
	req, _ := http.NewRequest("GET", suite.server.URL+"/error", nil)
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)

	// 测试404错误
	req, _ = http.NewRequest("GET", suite.server.URL+"/not-found", nil)
	resp, err = http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

// TestRequestHeaders 测试请求头
func (suite *SimpleIntegrationTestSuite) TestRequestHeaders() {
	// 注册请求头处理路由
	suite.router.Get("/headers", func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		accept := r.Header.Get("Accept")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"user_agent": userAgent,
			"accept":     accept,
		})
	})

	// 测试请求头
	req, _ := http.NewRequest("GET", suite.server.URL+"/headers", nil)
	req.Header.Set("User-Agent", "TestAgent/1.0")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	suite.Equal("TestAgent/1.0", response["user_agent"])
	suite.Equal("application/json", response["accept"])
}

// TestDatabaseCRUD 测试数据库CRUD操作
func (suite *SimpleIntegrationTestSuite) TestDatabaseCRUD() {
	// 创建用户
	_, err := suite.db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		"Jane Doe", "jane@example.com", "password123")
	suite.NoError(err)

	// 查询用户
	var id int64
	var name, email string
	err = suite.db.Query("SELECT id, name, email FROM users WHERE email = ?", "jane@example.com").Scan(&id, &name, &email)
	suite.NoError(err)
	suite.Equal("Jane Doe", name)
	suite.Equal("jane@example.com", email)

	// 更新用户
	_, err = suite.db.Exec("UPDATE users SET name = ? WHERE id = ?", "Jane Smith", id)
	suite.NoError(err)

	// 验证更新
	err = suite.db.Query("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	suite.NoError(err)
	suite.Equal("Jane Smith", name)

	// 删除用户
	_, err = suite.db.Exec("DELETE FROM users WHERE id = ?", id)
	suite.NoError(err)

	// 验证删除
	var count int
	err = suite.db.Query("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&count)
	suite.NoError(err)
	suite.Equal(0, count)
}

// TestConcurrentRequests 测试并发请求
func (suite *SimpleIntegrationTestSuite) TestConcurrentRequests() {
	// 注册计数器路由
	counter := 0
	suite.router.Get("/counter", func(w http.ResponseWriter, r *http.Request) {
		counter++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("count: %d", counter)))
	})

	// 发送多个并发请求
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			req, _ := http.NewRequest("GET", suite.server.URL+"/counter", nil)
			resp, err := http.DefaultClient.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				done <- true
			} else {
				done <- false
			}
		}()
	}

	// 等待所有请求完成
	successCount := 0
	for i := 0; i < 5; i++ {
		if <-done {
			successCount++
		}
	}

	suite.Equal(5, successCount)
}

// TestTimeoutHandling 测试超时处理
func (suite *SimpleIntegrationTestSuite) TestTimeoutHandling() {
	// 注册慢速路由
	suite.router.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("slow response"))
	})

	// 创建带超时的客户端
	client := &http.Client{
		Timeout: 50 * time.Millisecond,
	}

	// 测试超时
	req, _ := http.NewRequest("GET", suite.server.URL+"/slow", nil)
	_, err := client.Do(req)
	suite.Error(err) // 应该超时
}

// 运行简化集成测试套件
func TestSimpleIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleIntegrationTestSuite))
}
