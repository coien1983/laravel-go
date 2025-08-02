package main

import (
	"encoding/json"
	"laravel-go/framework/http"
	"laravel-go/framework/validation"
	"net/http"
	"strconv"
)

type PostController struct{}

type PostRequest struct {
	Title   string `json:"title" validate:"required,min=5,max=200"`
	Content string `json:"content" validate:"required,min=10"`
	Status  string `json:"status" validate:"required,oneof=draft published"`
}

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
	UserID  int    `json:"user_id"`
}

func (c *PostController) Index(ctx *http.Context) {
	// 模拟文章列表
	posts := []Post{
		{
			ID:      1,
			Title:   "第一篇博客文章",
			Content: "这是第一篇博客文章的内容...",
			Status:  "published",
			UserID:  1,
		},
		{
			ID:      2,
			Title:   "第二篇博客文章",
			Content: "这是第二篇博客文章的内容...",
			Status:  "published",
			UserID:  1,
		},
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"posts": posts,
		"total": len(posts),
	})
}

func (c *PostController) Show(ctx *http.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "无效的文章ID",
		})
		return
	}

	// 模拟获取文章
	post := Post{
		ID:      id,
		Title:   "示例文章",
		Content: "这是示例文章的内容...",
		Status:  "published",
		UserID:  1,
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"post": post,
	})
}

func (c *PostController) Store(ctx *http.Context) {
	var req PostRequest
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

	// 模拟创建文章
	post := Post{
		ID:      3,
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
		UserID:  1,
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "文章创建成功",
		"post":    post,
	})
}

func (c *PostController) Update(ctx *http.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "无效的文章ID",
		})
		return
	}

	var req PostRequest
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

	// 模拟更新文章
	post := Post{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
		UserID:  1,
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "文章更新成功",
		"post":    post,
	})
}

func (c *PostController) Destroy(ctx *http.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "无效的文章ID",
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "文章删除成功",
		"id":      id,
	})
} 