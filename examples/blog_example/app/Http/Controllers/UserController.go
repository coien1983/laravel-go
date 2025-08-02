package main

import (
	"encoding/json"
	"laravel-go/framework/http"
	"laravel-go/framework/validation"
	"net/http"
)

type UserController struct{}

type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Email string `json:"email" validate:"required,email"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c *UserController) Profile(ctx *http.Context) {
	// 模拟获取用户信息
	user := User{
		ID:    1,
		Name:  "测试用户",
		Email: "test@example.com",
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

func (c *UserController) UpdateProfile(ctx *http.Context) {
	var req UpdateProfileRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "无效的请求数据",
		})
		return
	}

	// 验证请求数据
	validator := validation.NewValidator()
	if err := validator.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// 模拟更新用户信息
	user := User{
		ID:    1,
		Name:  req.Name,
		Email: req.Email,
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "用户信息更新成功",
		"user":    user,
	})
} 