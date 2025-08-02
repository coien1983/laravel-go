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

type ProductService struct{}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}

func main() {
	fmt.Println("=== 产品微服务启动 ===")

	// 创建产品服务
	productService := &ProductService{}

	// 设置路由
	http.HandleFunc("/health", productService.HealthCheck)
	http.HandleFunc("/products", productService.GetProducts)
	http.HandleFunc("/products/", productService.GetProduct)

	// 启动服务器
	port := 8083
	go func() {
		fmt.Printf("产品服务启动在 http://localhost:%d\n", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatal("启动服务器失败:", err)
		}
	}()

	// 等待中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("正在关闭产品服务...")
	fmt.Println("产品服务已关闭")
}

func (s *ProductService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"service": "product-service",
		"version": "1.0.0",
	})
}

func (s *ProductService) GetProducts(w http.ResponseWriter, r *http.Request) {
	// 模拟产品列表
	products := []Product{
		{
			ID:          1,
			Name:        "iPhone 15",
			Description: "最新款iPhone手机",
			Price:       6999.00,
			Stock:       100,
			Category:    "电子产品",
		},
		{
			ID:          2,
			Name:        "MacBook Pro",
			Description: "专业级笔记本电脑",
			Price:       12999.00,
			Stock:       50,
			Category:    "电子产品",
		},
		{
			ID:          3,
			Name:        "AirPods Pro",
			Description: "无线降噪耳机",
			Price:       1999.00,
			Stock:       200,
			Category:    "配件",
		},
	}

	response := map[string]interface{}{
		"success": true,
		"data":    products,
		"total":   len(products),
		"service": "product-service",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *ProductService) GetProduct(w http.ResponseWriter, r *http.Request) {
	// 从路径中提取产品ID
	productID := r.URL.Path[len("/products/"):]

	// 模拟获取产品
	product := Product{
		ID:          1,
		Name:        "产品" + productID,
		Description: "这是产品" + productID + "的描述",
		Price:       100.00,
		Stock:       100,
		Category:    "默认分类",
	}

	response := map[string]interface{}{
		"success": true,
		"data":    product,
		"service": "product-service",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 