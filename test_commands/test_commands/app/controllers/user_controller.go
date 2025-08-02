package controllers

import (
	"github.com/coien1983/laravel-go/framework/http"
)

// User 控制器
type User struct {
	http.BaseController
}

// NewUser 创建新的控制器实例
func NewUser() *User {
	return &User{}
}

// Index 显示资源列表
func (c *User) Index() http.Response {
	return c.Json(map[string]interface{}{
		"message": "User Index",
	})
}

// Show 显示指定资源
func (c *User) Show(id string) http.Response {
	return c.Json(map[string]interface{}{
		"message": "User Show",
		"id":      id,
	})
}

// Store 存储新创建的资源
func (c *User) Store() http.Response {
	return c.Json(map[string]interface{}{
		"message": "User Store",
	})
}

// Update 更新指定资源
func (c *User) Update(id string) http.Response {
	return c.Json(map[string]interface{}{
		"message": "User Update",
		"id":      id,
	})
}

// Delete 删除指定资源
func (c *User) Delete(id string) http.Response {
	return c.Json(map[string]interface{}{
		"message": "User Delete",
		"id":      id,
	})
}
