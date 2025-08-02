package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type OrderService struct{}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	fmt.Println("=== 订单微服务启动 ===")

	// 创建订单服务
	orderService := &OrderService{}

	// 设置路由
	http.HandleFunc("/health", orderService.HealthCheck)
	http.HandleFunc("/orders", orderService.GetOrders)
	http.HandleFunc("/orders/", orderService.GetOrder)

	// 启动服务器
	port := 8084
	go func() {
		fmt.Printf("订单服务启动在 http://localhost:%d\n", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatal("启动服务器失败:", err)
		}
	}()

	// 等待中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("正在关闭订单服务...")
	fmt.Println("订单服务已关闭")
}

func (s *OrderService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"service": "order-service",
		"version": "1.0.0",
	})
}

func (s *OrderService) GetOrders(w http.ResponseWriter, r *http.Request) {
	// 模拟订单列表
	orders := []Order{
		{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  2,
			Total:     13998.00,
			Status:    "pending",
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        2,
			UserID:    2,
			ProductID: 2,
			Quantity:  1,
			Total:     12999.00,
			Status:    "completed",
			CreatedAt: time.Now().Add(-48 * time.Hour),
		},
		{
			ID:        3,
			UserID:    3,
			ProductID: 3,
			Quantity:  3,
			Total:     5997.00,
			Status:    "shipped",
			CreatedAt: time.Now().Add(-72 * time.Hour),
		},
	}

	response := map[string]interface{}{
		"success": true,
		"data":    orders,
		"total":   len(orders),
		"service": "order-service",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *OrderService) GetOrder(w http.ResponseWriter, r *http.Request) {
	// 从路径中提取订单ID
	orderID := r.URL.Path[len("/orders/"):]

	// 模拟获取订单
	order := Order{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  1,
		Total:     100.00,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	response := map[string]interface{}{
		"success": true,
		"data":    order,
		"service": "order-service",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
