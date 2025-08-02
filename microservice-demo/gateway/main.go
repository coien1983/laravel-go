package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "microservice-demo/proto/user"
)

// Gateway API网关
type Gateway struct {
	userClient pb.UserServiceClient
	router     *mux.Router
}

// NewGateway 创建网关实例
func NewGateway() (*Gateway, error) {
	// 连接gRPC服务
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	userClient := pb.NewUserServiceClient(conn)

	router := mux.NewRouter()
	gateway := &Gateway{
		userClient: userClient,
		router:     router,
	}

	// 注册路由
	gateway.registerRoutes()

	return gateway, nil
}

// registerRoutes 注册路由
func (gateway *Gateway) registerRoutes() {
	// 中间件
	gateway.router.Use(gateway.loggingMiddleware)
	gateway.router.Use(gateway.corsMiddleware)

	// API路由
	api := gateway.router.PathPrefix("/api/v1").Subrouter()
	
	// 用户相关路由
	api.HandleFunc("/users", gateway.getUsers).Methods("GET")
	api.HandleFunc("/users/{id}", gateway.getUser).Methods("GET")
	api.HandleFunc("/users", gateway.createUser).Methods("POST")
	api.HandleFunc("/users/{id}", gateway.updateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", gateway.deleteUser).Methods("DELETE")

	// 健康检查
	gateway.router.HandleFunc("/health", gateway.healthCheck).Methods("GET")
}

// loggingMiddleware 日志中间件
func (gateway *Gateway) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("API Gateway: %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// corsMiddleware CORS中间件
func (gateway *Gateway) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// healthCheck 健康检查
func (gateway *Gateway) healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
		"service": "api-gateway",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getUsers 获取用户列表
func (gateway *Gateway) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 从查询参数获取分页信息
	page := int32(1)
	pageSize := int32(10)
	search := r.URL.Query().Get("search")

	resp, err := gateway.userClient.ListUsers(ctx, &pb.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// getUser 获取单个用户
func (gateway *Gateway) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	
	// TODO: 解析用户ID
	id := int64(1) // 示例

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := gateway.userClient.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// createUser 创建用户
func (gateway *Gateway) createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := gateway.userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// updateUser 更新用户
func (gateway *Gateway) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	
	// TODO: 解析用户ID
	id := int64(1) // 示例

	var req struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Avatar string `json:"avatar"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := gateway.userClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:     id,
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Avatar: req.Avatar,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// deleteUser 删除用户
func (gateway *Gateway) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	
	// TODO: 解析用户ID
	id := int64(1) // 示例

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := gateway.userClient.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	gateway, err := NewGateway()
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	port := ":8080"
	if envPort := os.Getenv("GATEWAY_PORT"); envPort != "" {
		port = ":" + envPort
	}

	server := &http.Server{
		Addr:    port,
		Handler: gateway.router,
	}

	// 启动服务器
	go func() {
		fmt.Printf("🚀 API Gateway starting on http://localhost%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gateway error: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	fmt.Println("\n🛑 Shutting down API Gateway...")
	fmt.Println("✅ API Gateway stopped gracefully")
}