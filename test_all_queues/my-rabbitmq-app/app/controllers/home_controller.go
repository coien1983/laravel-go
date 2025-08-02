package controllers

import (
	"net/http"
	"encoding/json"
)

// HomeController 首页控制器
type HomeController struct{}

// NewHomeController 创建新的首页控制器
func NewHomeController() *HomeController {
	return &HomeController{}
}

// Index 首页
func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Welcome to Laravel-Go!",
		"version": "1.0.0",
		"status":  "running",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Health 健康检查
func (c *HomeController) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "ok",
		"time":   "2024-01-01T00:00:00Z",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}