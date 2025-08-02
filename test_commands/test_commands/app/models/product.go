package models

import (
	"github.com/coien1983/laravel-go/framework/database"
)

// Product 模型
type Product struct {
	database.Model
	
}

// TableName 获取表名
func (m *Product) TableName() string {
	return "products"
}

// NewProduct 创建新的模型实例
func NewProduct() *Product {
	return &Product{}
}
