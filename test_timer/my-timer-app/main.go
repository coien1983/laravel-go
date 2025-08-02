package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 设置服务器
	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	// 创建 HTTP 服务器
	mux := http.NewServeMux()
	
	// 注册路由
	registerRoutes(mux)
	
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// 启动服务器
	go func() {
		fmt.Printf("🚀 Server starting on http://localhost%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	fmt.Println("\n🛑 Shutting down server...")
	fmt.Println("✅ Server stopped gracefully")
}

// registerRoutes 注册路由
func registerRoutes(mux *http.ServeMux) {
	// 导入路由包
	// 这里会在运行时动态加载路由
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message": "Welcome to Laravel-Go!",
			"version": "1.0.0",
			"status":  "running",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status": "ok",
			"time":   "2024-01-01T00:00:00Z",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}