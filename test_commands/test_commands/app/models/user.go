package models

import (
	"github.com/coien1983/laravel-go/framework/database"
)

// User 模型
type User struct {
	database.Model
	
}

// TableName 获取表名
func (m *User) TableName() string {
	return "users"
}

// NewUser 创建新的模型实例
func NewUser() *User {
	return &User{}
}
