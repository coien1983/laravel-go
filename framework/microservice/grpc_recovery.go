package microservice

import (
	"context"
	"fmt"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"laravel-go/framework/errors"
)

// RecoveryInterceptor gRPC panic恢复拦截器
type RecoveryInterceptor struct {
	errorHandler errors.ErrorHandler
	logger       Logger
}

// Logger 日志接口
type Logger interface {
	Error(message string, context map[string]interface{})
	Warning(message string, context map[string]interface{})
	Info(message string, context map[string]interface{})
	Debug(message string, context map[string]interface{})
}

// NewRecoveryInterceptor 创建panic恢复拦截器
func NewRecoveryInterceptor(errorHandler errors.ErrorHandler, logger Logger) *RecoveryInterceptor {
	return &RecoveryInterceptor{
		errorHandler: errorHandler,
		logger:       logger,
	}
}

// UnaryServerInterceptor 一元服务拦截器
func (ri *RecoveryInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if panicVal := recover(); panicVal != nil {
				// 记录panic信息
				ri.logPanic(panicVal, info.FullMethod, ctx)

				// 创建错误
				err = status.Errorf(codes.Internal, "Internal server error: %v", panicVal)

				// 处理错误
				if ri.errorHandler != nil {
					appErr := errors.New(fmt.Sprintf("gRPC panic: %v", panicVal))
					_ = ri.errorHandler.Handle(appErr)
				}
			}
		}()

		return handler(ctx, req)
	}
}

// StreamServerInterceptor 流服务拦截器
func (ri *RecoveryInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if panicVal := recover(); panicVal != nil {
				// 记录panic信息
				ri.logPanic(panicVal, info.FullMethod, stream.Context())

				// 创建错误
				err = status.Errorf(codes.Internal, "Internal server error: %v", panicVal)

				// 处理错误
				if ri.errorHandler != nil {
					appErr := errors.New(fmt.Sprintf("gRPC stream panic: %v", panicVal))
					_ = ri.errorHandler.Handle(appErr)
				}
			}
		}()

		return handler(srv, stream)
	}
}

// logPanic 记录panic信息
func (ri *RecoveryInterceptor) logPanic(panic interface{}, method string, ctx context.Context) {
	if ri.logger == nil {
		return
	}

	stack := debug.Stack()

	context := map[string]interface{}{
		"panic":   panic,
		"stack":   string(stack),
		"method":  method,
		"context": ctx,
	}

	ri.logger.Error("gRPC panic recovered", context)
}

// SafeUnaryHandler 安全一元处理器包装器
func SafeUnaryHandler(handler grpc.UnaryHandler, errorHandler errors.ErrorHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		defer func() {
			if panicVal := recover(); panicVal != nil {
				// 记录panic
				if errorHandler != nil {
					appErr := errors.New(fmt.Sprintf("Unary handler panic: %v", panicVal))
					_ = errorHandler.Handle(appErr)
				}

				// 返回错误
				err = status.Errorf(codes.Internal, "Internal server error: %v", panicVal)
			}
		}()

		return handler(ctx, req)
	}
}

// SafeStreamHandler 安全流处理器包装器
func SafeStreamHandler(handler grpc.StreamHandler, errorHandler errors.ErrorHandler) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) (err error) {
		defer func() {
			if panicVal := recover(); panicVal != nil {
				// 记录panic
				if errorHandler != nil {
					appErr := errors.New(fmt.Sprintf("Stream handler panic: %v", panicVal))
					_ = errorHandler.Handle(appErr)
				}

				// 返回错误
				err = status.Errorf(codes.Internal, "Internal server error: %v", panicVal)
			}
		}()

		return handler(srv, stream)
	}
}
