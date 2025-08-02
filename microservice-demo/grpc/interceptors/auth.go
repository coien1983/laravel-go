package interceptors

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor 认证拦截器
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 跳过认证的方法
	skipAuthMethods := map[string]bool{
		"/user.UserService/GetUser": true,
		"/user.UserService/ListUsers": true,
	}
	
	if skipAuthMethods[info.FullMethod] {
		return handler(ctx, req)
	}
	
	// 从元数据中获取token
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	
	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	
	token := authHeader[0]
	if !strings.HasPrefix(token, "Bearer ") {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
	}
	
	// TODO: 验证token
	tokenValue := strings.TrimPrefix(token, "Bearer ")
	if tokenValue == "" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	
	// 将用户信息添加到上下文
	userID := "123" // TODO: 从token中解析用户ID
	newCtx := context.WithValue(ctx, "user_id", userID)
	
	return handler(newCtx, req)
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in context")
	}
	return userID, nil
}