package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"laravel-go/framework/auth"
	"laravel-go/framework/config"
	"laravel-go/framework/container"
	"laravel-go/framework/database"
	"laravel-go/framework/routing"
	"laravel-go/framework/validation"
)

// HTTPIntegrationTestSuite HTTP集成测试套件
type HTTPIntegrationTestSuite struct {
	suite.Suite
	app    *container.Container
	router *routing.Router
	server *httptest.Server
	db     database.Connection
}

// SetupSuite 设置测试套件
func (suite *HTTPIntegrationTestSuite) SetupSuite() {
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

	// 注册认证服务
	suite.app.Singleton("auth", func() interface{} {
		return auth.NewAuth(suite.app)
	})

	// 获取服务实例
	suite.db = suite.app.Make("database").(database.Connection)

	// 初始化路由
	suite.router = routing.NewRouter()

	// 创建测试服务器
	suite.server = httptest.NewServer(suite.router)

	// 设置数据库表
	suite.setupDatabase()
}

// TearDownSuite 清理测试套件
func (suite *HTTPIntegrationTestSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest 设置每个测试
func (suite *HTTPIntegrationTestSuite) SetupTest() {
	// 清理数据库
	suite.cleanupDatabase()
}

// setupDatabase 设置数据库表
func (suite *HTTPIntegrationTestSuite) setupDatabase() {
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
func (suite *HTTPIntegrationTestSuite) cleanupDatabase() {
	suite.db.Exec("DELETE FROM posts")
	suite.db.Exec("DELETE FROM users")
}

// TestBasicRouting 测试基础路由
func (suite *HTTPIntegrationTestSuite) TestBasicRouting() {
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
func (suite *HTTPIntegrationTestSuite) TestRouteParameters() {
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
func (suite *HTTPIntegrationTestSuite) TestRouteGroup() {
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

// TestMiddleware 测试中间件
func (suite *HTTPIntegrationTestSuite) TestMiddleware() {
	// 创建测试中间件
	authMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next(w, r)
		}
	}

	loggingMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-ID", "test-123")
			next(w, r)
		}
	}

	// 注册带中间件的路由
	suite.router.Get("/protected", authMiddleware(loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("protected content"))
	})))

	// 测试未授权访问
	req, _ := http.NewRequest("GET", suite.server.URL+"/protected", nil)
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)

	// 测试授权访问
	req, _ = http.NewRequest("GET", suite.server.URL+"/protected", nil)
	req.Header.Set("Authorization", "Bearer token123")
	resp, err = http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("test-123", resp.Header.Get("X-Request-ID"))
}

// TestControllerCRUD 测试控制器CRUD操作
func (suite *HTTPIntegrationTestSuite) TestControllerCRUD() {
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

	// 测试更新用户
	updateData := map[string]interface{}{
		"name": "John Smith",
	}
	updateJSON, _ := json.Marshal(updateData)

	req, _ = http.NewRequest("PUT", suite.server.URL+"/users/1", bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 测试删除用户
	req, _ = http.NewRequest("DELETE", suite.server.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusNoContent, resp.StatusCode)
}

// TestValidationMiddleware 测试验证中间件
func (suite *HTTPIntegrationTestSuite) TestValidationMiddleware() {
	// 创建验证中间件
	validationMiddleware := func(rules map[string]string) func(http.HandlerFunc) http.HandlerFunc {
		return func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				// 解析请求体
				var data map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				// 验证数据
				validator := validation.NewValidator()
				errors := validator.Validate(data, rules)

				if len(errors) > 0 {
					w.WriteHeader(http.StatusUnprocessableEntity)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"errors": errors,
					})
					return
				}

				next(w, r)
			}
		}
	}

	// 定义验证规则
	rules := map[string]string{
		"name":  "required|string|min:2",
		"email": "required|email",
		"age":   "required|int|min:18",
	}

	// 注册带验证的路由
	suite.router.Post("/validate", validationMiddleware(rules)(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("validation passed"))
	}))

	// 测试验证失败
	invalidData := map[string]interface{}{
		"name":  "",
		"email": "invalid-email",
		"age":   "not-a-number",
	}
	invalidJSON, _ := json.Marshal(invalidData)

	req, _ := http.NewRequest("POST", suite.server.URL+"/validate", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusUnprocessableEntity, resp.StatusCode)

	// 测试验证成功
	validData := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   25,
	}
	validJSON, _ := json.Marshal(validData)

	req, _ = http.NewRequest("POST", suite.server.URL+"/validate", bytes.NewBuffer(validJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
}

// TestErrorHandling 测试错误处理
func (suite *HTTPIntegrationTestSuite) TestErrorHandling() {
	// 注册测试路由
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

// TestContentTypes 测试内容类型
func (suite *HTTPIntegrationTestSuite) TestContentTypes() {
	// 注册测试路由
	suite.router.Post("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "JSON response"})
	})

	suite.router.Post("/xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<message>XML response</message>"))
	})

	// 测试JSON响应
	req, _ := http.NewRequest("POST", suite.server.URL+"/json", nil)
	resp, err := http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("application/json", resp.Header.Get("Content-Type"))

	// 测试XML响应
	req, _ = http.NewRequest("POST", suite.server.URL+"/xml", nil)
	resp, err = http.DefaultClient.Do(req)

	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("application/xml", resp.Header.Get("Content-Type"))
}

// TestRequestHeaders 测试请求头
func (suite *HTTPIntegrationTestSuite) TestRequestHeaders() {
	// 注册测试路由
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

// 运行HTTP集成测试套件
func TestHTTPIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPIntegrationTestSuite))
}
