package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type UserAPIController struct{}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func (c *UserAPIController) Index(w http.ResponseWriter, r *http.Request) {
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
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *UserAPIController) Show(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// 尝试从路径参数获取
		idStr = r.URL.Path[len("/api/v1/users/"):]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的用户ID",
		})
		return
	}

	// 模拟获取用户
	user := User{
		ID:    id,
		Name:  "用户" + idStr,
		Email: "user" + idStr + "@example.com",
		Age:   25 + id,
	}

	response := map[string]interface{}{
		"success": true,
		"data":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *UserAPIController) Store(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的请求数据",
		})
		return
	}

	// 模拟创建用户
	user := User{
		ID:    4,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	response := map[string]interface{}{
		"success": true,
		"message": "用户创建成功",
		"data":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *UserAPIController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// 尝试从路径参数获取
		idStr = r.URL.Path[len("/api/v1/users/"):]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的用户ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的请求数据",
		})
		return
	}

	// 模拟更新用户
	user := User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	response := map[string]interface{}{
		"success": true,
		"message": "用户更新成功",
		"data":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *UserAPIController) Destroy(w http.ResponseWriter, r *http.Request) {
	// 从URL参数获取ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// 尝试从路径参数获取
		idStr = r.URL.Path[len("/api/v1/users/"):]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "无效的用户ID",
		})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "用户删除成功",
		"id":      id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 