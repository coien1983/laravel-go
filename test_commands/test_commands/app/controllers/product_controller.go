package controllers

import (
	"github.com/coien1983/laravel-go/framework/http"
)

// Product 控制器
type Product struct {
	http.BaseController
}

// NewProduct 创建新的控制器实例
func NewProduct() *Product {
	return &Product{}
}

// Index 显示资源列表
func (c *Product) Index() http.Response {
	return c.Json(map[string]interface{}{
		"message": "Product Index",
	})
}

// Show 显示指定资源
func (c *Product) Show(id string) http.Response {
	return c.Json(map[string]interface{}{
		"message": "Product Show",
		"id":      id,
	})
}

// Store 存储新创建的资源
func (c *Product) Store() http.Response {
	return c.Json(map[string]interface{}{
		"message": "Product Store",
	})
}

// Update 更新指定资源
func (c *Product) Update(id string) http.Response {
	return c.Json(map[string]interface{}{
		"message": "Product Update",
		"id":      id,
	})
}

// Delete 删除指定资源
func (c *Product) Delete(id string) http.Response {
	return c.Json(map[string]interface{}{
		"message": "Product Delete",
		"id":      id,
	})
}
