package main

import (
	"fmt"
	"laravel-go/framework/core"
	"laravel-go/framework/http"
	"laravel-go/framework/validation"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("=== Laravel-Go API 示例 ===")

	// 初始化应用
	app := core.NewApplication()

	// 创建HTTP服务器
	server := http.NewServer(app.Config, app.Container)

	// 注册中间件
	server.Use(LoggingMiddleware)
	server.Use(CORSMiddleware)

	// 注册路由
	registerAPIRoutes(server)

	// 启动服务器
	go func() {
		port := "8081"
		fmt.Printf("API 服务器启动在 http://localhost:%s\n", port)
		fmt.Printf("API 文档: http://localhost:%s/docs\n", port)
		if err := server.Start(":" + port); err != nil {
			log.Fatal("启动服务器失败:", err)
		}
	}()

	// 等待中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("正在关闭API服务器...")
	fmt.Println("API服务器已关闭")
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

func registerAPIRoutes(server *http.Server) {
	// API 版本控制
	v1 := server.Group("/api/v1")
	
	// 健康检查
	v1.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","version":"1.0.0"}`))
	})

	// 用户API
	userController := &UserAPIController{}
	v1.Get("/users", userController.Index)
	v1.Get("/users/:id", userController.Show)
	v1.Post("/users", userController.Store)
	v1.Put("/users/:id", userController.Update)
	v1.Delete("/users/:id", userController.Destroy)

	// 产品API
	productController := &ProductAPIController{}
	v1.Get("/products", productController.Index)
	v1.Get("/products/:id", productController.Show)
	v1.Post("/products", productController.Store)
	v1.Put("/products/:id", productController.Update)
	v1.Delete("/products/:id", productController.Destroy)

	// 订单API
	orderController := &OrderAPIController{}
	v1.Get("/orders", orderController.Index)
	v1.Get("/orders/:id", orderController.Show)
	v1.Post("/orders", orderController.Store)
	v1.Put("/orders/:id", orderController.Update)
	v1.Delete("/orders/:id", orderController.Destroy)

	// API文档
	server.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Laravel-Go API 文档</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .method { font-weight: bold; color: #007bff; }
        .url { font-family: monospace; background: #f8f9fa; padding: 5px; }
        .description { color: #666; margin-top: 10px; }
    </style>
</head>
<body>
    <h1>Laravel-Go API 文档</h1>
    
    <h2>健康检查</h2>
    <div class="endpoint">
        <div class="method">GET</div>
        <div class="url">/api/v1/health</div>
        <div class="description">检查API服务状态</div>
    </div>

    <h2>用户管理</h2>
    <div class="endpoint">
        <div class="method">GET</div>
        <div class="url">/api/v1/users</div>
        <div class="description">获取用户列表</div>
    </div>
    <div class="endpoint">
        <div class="method">GET</div>
        <div class="url">/api/v1/users/:id</div>
        <div class="description">获取单个用户</div>
    </div>
    <div class="endpoint">
        <div class="method">POST</div>
        <div class="url">/api/v1/users</div>
        <div class="description">创建新用户</div>
    </div>
    <div class="endpoint">
        <div class="method">PUT</div>
        <div class="url">/api/v1/users/:id</div>
        <div class="description">更新用户信息</div>
    </div>
    <div class="endpoint">
        <div class="method">DELETE</div>
        <div class="url">/api/v1/users/:id</div>
        <div class="description">删除用户</div>
    </div>

    <h2>产品管理</h2>
    <div class="endpoint">
        <div class="method">GET</div>
        <div class="url">/api/v1/products</div>
        <div class="description">获取产品列表</div>
    </div>
    <div class="endpoint">
        <div class="method">GET</div>
        <div class="url">/api/v1/products/:id</div>
        <div class="description">获取单个产品</div>
    </div>
    <div class="endpoint">
        <div class="method">POST</div>
        <div class="url">/api/v1/products</div>
        <div class="description">创建新产品</div>
    </div>
    <div class="endpoint">
        <div class="method">PUT</div>
        <div class="url">/api/v1/products/:id</div>
        <div class="description">更新产品信息</div>
    </div>
    <div class="endpoint">
        <div class="method">DELETE</div>
        <div class="url">/api/v1/products/:id</div>
        <div class="description">删除产品</div>
    </div>

    <h2>订单管理</h2>
    <div class="endpoint">
        <div class="method">GET</div>
        <div class="url">/api/v1/orders</div>
        <div class="description">获取订单列表</div>
    </div>
    <div class="endpoint">
        <div class="method">GET</div>
        <div class="url">/api/v1/orders/:id</div>
        <div class="description">获取单个订单</div>
    </div>
    <div class="endpoint">
        <div class="method">POST</div>
        <div class="url">/api/v1/orders</div>
        <div class="description">创建新订单</div>
    </div>
    <div class="endpoint">
        <div class="method">PUT</div>
        <div class="url">/api/v1/orders/:id</div>
        <div class="description">更新订单信息</div>
    </div>
    <div class="endpoint">
        <div class="method">DELETE</div>
        <div class="url">/api/v1/orders/:id</div>
        <div class="description">删除订单</div>
    </div>
</body>
</html>`
		w.Write([]byte(html))
	})
} 