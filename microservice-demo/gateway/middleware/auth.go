package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 跳过认证的路径
		skipAuthPaths := map[string]bool{
			"/health": true,
			"/api/v1/users": true, // GET请求
		}
		
		if skipAuthPaths[r.URL.Path] && r.Method == "GET" {
			next.ServeHTTP(w, r)
			return
		}
		
		// 获取Authorization头
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		// 验证Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		
		// TODO: 验证token
		// 这里应该调用认证服务验证token
		
		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: 实现限流逻辑
		// 这里可以使用Redis或其他存储来实现限流
		
		next.ServeHTTP(w, r)
	})
}