package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type ProductAPIController struct{}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}

func (c *ProductAPIController) Index(w http.ResponseWriter, r *http.Request) {
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
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *ProductAPIController) Show(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// 尝试从路径参数获取
		idStr = r.URL.Path[len("/api/v1/products/"):]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的产品ID",
		})
		return
	}

	// 模拟获取产品
	product := Product{
		ID:          id,
		Name:        "产品" + idStr,
		Description: "这是产品" + idStr + "的描述",
		Price:       100.00 + float64(id)*50,
		Stock:       100 + id*10,
		Category:    "默认分类",
	}

	response := map[string]interface{}{
		"success": true,
		"data":    product,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *ProductAPIController) Store(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的请求数据",
		})
		return
	}

	// 模拟创建产品
	product := Product{
		ID:          4,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}

	response := map[string]interface{}{
		"success": true,
		"message": "产品创建成功",
		"data":    product,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *ProductAPIController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// 尝试从路径参数获取
		idStr = r.URL.Path[len("/api/v1/products/"):]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的产品ID",
		})
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的请求数据",
		})
		return
	}

	// 模拟更新产品
	product := Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}

	response := map[string]interface{}{
		"success": true,
		"message": "产品更新成功",
		"data":    product,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *ProductAPIController) Destroy(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// 尝试从路径参数获取
		idStr = r.URL.Path[len("/api/v1/products/"):]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的产品ID",
		})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "产品删除成功",
		"id":      id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 