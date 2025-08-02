package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Gateway struct {
	services map[string]string
}

func main() {
	fmt.Println("=== API 网关启动 ===")

	// 创建网关
	gateway := &Gateway{
		services: map[string]string{
			"user-service":    "http://localhost:8082",
			"product-service": "http://localhost:8083",
			"order-service":   "http://localhost:8084",
		},
	}

	// 设置路由
	http.HandleFunc("/", gateway.Route)
	http.HandleFunc("/health", gateway.HealthCheck)
	http.HandleFunc("/services", gateway.ListServices)

	// 启动服务器
	port := 8080
	go func() {
		fmt.Printf("API 网关启动在 http://localhost:%d\n", port)
		fmt.Printf("用户服务: http://localhost:%d/users\n", port)
		fmt.Printf("产品服务: http://localhost:%d/products\n", port)
		fmt.Printf("订单服务: http://localhost:%d/orders\n", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatal("启动服务器失败:", err)
		}
	}()

	// 等待中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("正在关闭API网关...")
	fmt.Println("API网关已关闭")
}

func (g *Gateway) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"service": "api-gateway",
		"version": "1.0.0",
	})
}

func (g *Gateway) ListServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    g.services,
		"service": "api-gateway",
	})
}

func (g *Gateway) Route(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// 添加CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 路由到相应的服务
	if strings.HasPrefix(path, "/users") {
		g.proxyToService(w, r, "user-service", path)
	} else if strings.HasPrefix(path, "/products") {
		g.proxyToService(w, r, "product-service", path)
	} else if strings.HasPrefix(path, "/orders") {
		g.proxyToService(w, r, "order-service", path)
	} else {
		// 默认响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "欢迎使用 Laravel-Go 微服务 API 网关",
			"version": "1.0.0",
			"services": map[string]string{
				"users":    "/users",
				"products": "/products",
				"orders":   "/orders",
			},
		})
	}
}

func (g *Gateway) proxyToService(w http.ResponseWriter, r *http.Request, serviceName, path string) {
	serviceURL, exists := g.services[serviceName]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "服务不存在: " + serviceName,
		})
		return
	}

	// 创建目标URL
	targetURL, err := url.Parse(serviceURL + path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "无效的服务URL",
		})
		return
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 修改请求
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL.Host

	// 代理请求
	proxy.ServeHTTP(w, r)
}
