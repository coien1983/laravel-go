package main

import (
	"encoding/json"
	"laravel-go/framework/http"
	"laravel-go/framework/validation"
	"net/http"
)

type AuthController struct{}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (c *AuthController) Register(ctx *http.Context) {
	var req RegisterRequest
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

	// 模拟用户注册
	user := map[string]interface{}{
		"id":    1,
		"name":  req.Name,
		"email": req.Email,
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "用户注册成功",
		"user":    user,
	})
}

func (c *AuthController) Login(ctx *http.Context) {
	var req LoginRequest
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

	// 模拟用户登录
	user := map[string]interface{}{
		"id":    1,
		"name":  "测试用户",
		"email": req.Email,
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "登录成功",
		"user":    user,
		"token":   "mock-jwt-token",
	})
}

func (c *AuthController) Logout(ctx *http.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "退出登录成功",
	})
} 