package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type OrderAPIController struct{}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type UpdateOrderRequest struct {
	Status string `json:"status"`
}

func (c *OrderAPIController) Index(w http.ResponseWriter, r *http.Request) {
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
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *OrderAPIController) Show(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// 尝试从路径参数获取
		idStr = r.URL.Path[len("/api/v1/orders/"):]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的订单ID",
		})
		return
	}

	// 模拟获取订单
	order := Order{
		ID:        id,
		UserID:    id,
		ProductID: id,
		Quantity:  1,
		Total:     100.00 + float64(id)*50,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	response := map[string]interface{}{
		"success": true,
		"data":    order,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *OrderAPIController) Store(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的请求数据",
		})
		return
	}

	// 模拟创建订单
	order := Order{
		ID:        4,
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Total:     float64(req.Quantity) * 100.00,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	response := map[string]interface{}{
		"success": true,
		"message": "订单创建成功",
		"data":    order,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *OrderAPIController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// 尝试从路径参数获取
		idStr = r.URL.Path[len("/api/v1/orders/"):]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的订单ID",
		})
		return
	}

	var req UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的请求数据",
		})
		return
	}

	// 模拟更新订单
	order := Order{
		ID:        id,
		UserID:    id,
		ProductID: id,
		Quantity:  1,
		Total:     100.00 + float64(id)*50,
		Status:    req.Status,
		CreatedAt: time.Now(),
	}

	response := map[string]interface{}{
		"success": true,
		"message": "订单更新成功",
		"data":    order,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *OrderAPIController) Destroy(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// 尝试从路径参数获取
		idStr = r.URL.Path[len("/api/v1/orders/"):]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的订单ID",
		})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "订单删除成功",
		"id":      id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 