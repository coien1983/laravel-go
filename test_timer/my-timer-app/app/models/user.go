package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 获取表名
func (u *User) TableName() string {
	return "users"
}

// NewUser 创建新用户
func NewUser() *User {
	return &User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Fillable 可填充字段
func (u *User) Fillable() []string {
	return []string{"name", "email", "password"}
}

// Hidden 隐藏字段
func (u *User) Hidden() []string {
	return []string{"password"}
}