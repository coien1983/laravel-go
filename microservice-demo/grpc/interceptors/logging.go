package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor 日志拦截器
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	
	// 调用实际的RPC方法
	resp, err := handler(ctx, req)
	
	// 记录日志
	duration := time.Since(start)
	statusCode := codes.OK
	if err != nil {
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}
		log.Printf("gRPC: %s | %s | %v | %s", info.FullMethod, statusCode, duration, err)
	} else {
		log.Printf("gRPC: %s | %s | %v", info.FullMethod, statusCode, duration)
	}
	
	return resp, err
}

// RecoveryInterceptor 恢复拦截器
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("gRPC panic: %v", r)
			err = status.Errorf(codes.Internal, "Internal server error")
		}
	}()
	
	return handler(ctx, req)
}