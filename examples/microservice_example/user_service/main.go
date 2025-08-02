package main

import (
	"encoding/json"
	"fmt"
	"laravel-go/framework/microservice"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type UserService struct {
	registry microservice.Registry
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {
	fmt.Println("=== 用户微服务启动 ===")

	// 创建服务注册中心
	registry, err := microservice.NewMemoryRegistry()
	if err != nil {
		log.Fatal("创建注册中心失败:", err)
	}

	// 创建用户服务
	userService := &UserService{
		registry: registry,
	}

	// 注册服务
	serviceInfo := &microservice.ServiceInfo{
		Name:    "user-service",
		Version: "1.0.0",
		Address: "localhost:8082",
		Port:    8082,
		Tags:    []string{"user", "api"},
		Metadata: map[string]string{
			"protocol": "http",
			"health":   "/health",
		},
	}

	if err := registry.Register(serviceInfo); err != nil {
		log.Fatal("注册服务失败:", err)
	}

	// 设置路由
	http.HandleFunc("/health", userService.HealthCheck)
	http.HandleFunc("/users", userService.GetUsers)
	http.HandleFunc("/users/", userService.GetUser)

	// 启动服务器
	go func() {
		fmt.Printf("用户服务启动在 http://localhost:%d\n", serviceInfo.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", serviceInfo.Port), nil); err != nil {
			log.Fatal("启动服务器失败:", err)
		}
	}()

	// 等待中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("正在关闭用户服务...")

	// 注销服务
	if err := registry.Deregister(serviceInfo.Name); err != nil {
		log.Printf("注销服务失败: %v", err)
	}

	fmt.Println("用户服务已关闭")
}

func (s *UserService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"service": "user-service",
		"version": "1.0.0",
	})
}

func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	// 模拟用户列表
	users := []User{
		{ID: 1, Name: "张三", Email: "zhangsan@example.com", Age: 25},
		{ID: 2, Name: "李四", Email: "lisi@example.com", Age: 30},
		{ID: 3, Name: "王五", Email: "wangwu@example.com", Age: 28},
	}

	response := map[string]interface{}{
		"success": true,
		"data":    users,
		"total":   len(users),
		"service": "user-service",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *UserService) GetUser(w http.ResponseWriter, r *http.Request) {
	// 从路径中提取用户ID
	userID := r.URL.Path[len("/users/"):]

	// 模拟获取用户
	user := User{
		ID:    1,
		Name:  "用户" + userID,
		Email: "user" + userID + "@example.com",
		Age:   25,
	}

	response := map[string]interface{}{
		"success": true,
		"data":    user,
		"service": "user-service",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 