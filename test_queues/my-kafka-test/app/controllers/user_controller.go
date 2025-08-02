package controllers

import (
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/gorilla/mux"
)

// User 用户模型
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserController 用户控制器
type UserController struct {
	users []User
}

// NewUserController 创建新的用户控制器
func NewUserController() *UserController {
	// 初始化一些示例数据
	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	
	return &UserController{
		users: users,
	}
}

// Index 获取用户列表
func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c.users)
}

// Show 获取单个用户
func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	for _, user := range c.users {
		if user.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	
	http.Error(w, "User not found", http.StatusNotFound)
}

// Store 创建用户
func (c *UserController) Store(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// 简单的 ID 生成
	user.ID = len(c.users) + 1
	c.users = append(c.users, user)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}