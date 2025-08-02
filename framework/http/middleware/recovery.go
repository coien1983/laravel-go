package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"laravel-go/framework/errors"
)

// RecoveryMiddleware HTTP panic恢复中间件
type RecoveryMiddleware struct {
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

// NewRecoveryMiddleware 创建panic恢复中间件
func NewRecoveryMiddleware(errorHandler errors.ErrorHandler, logger Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		errorHandler: errorHandler,
		logger:       logger,
	}
}

// Handle 处理HTTP请求
func (m *RecoveryMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if panicVal := recover(); panicVal != nil {
				// 记录panic信息
				m.logPanic(panicVal, r)
				
				// 创建错误响应
				err := errors.New(fmt.Sprintf("Internal server error: %v", panicVal))
				if m.errorHandler != nil {
					_ = m.errorHandler.Handle(err)
				}
				
				// 返回500错误
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// logPanic 记录panic信息
func (m *RecoveryMiddleware) logPanic(panic interface{}, r *http.Request) {
	if m.logger == nil {
		return
	}
	
	stack := debug.Stack()
	
	context := map[string]interface{}{
		"panic":     panic,
		"stack":     string(stack),
		"timestamp": time.Now(),
		"method":    r.Method,
		"url":       r.URL.String(),
		"remote_addr": r.RemoteAddr,
		"user_agent": r.UserAgent(),
	}
	
	m.logger.Error("HTTP panic recovered", context)
}

// SafeHandler 安全处理器包装器
func SafeHandler(handler http.HandlerFunc, errorHandler errors.ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				// 记录panic
				if errorHandler != nil {
					err := errors.New(fmt.Sprintf("Handler panic: %v", r))
					errorHandler.Handle(err)
				}
				
				// 返回错误响应
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		
		handler(w, r)
	}
}

// SafeHandlerWithContext 带上下文的安全处理器
func SafeHandlerWithContext(handler func(context.Context, http.ResponseWriter, *http.Request) error, errorHandler errors.ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				// 记录panic
				if errorHandler != nil {
					err := errors.New(fmt.Sprintf("Handler panic: %v", r))
					errorHandler.Handle(err)
				}
				
				// 返回错误响应
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		
		// 检查上下文是否已取消
		select {
		case <-r.Context().Done():
			http.Error(w, "Request cancelled", http.StatusRequestTimeout)
			return
		default:
		}
		
		// 执行处理器
		if err := handler(r.Context(), w, r); err != nil {
			if errorHandler != nil {
				err = errorHandler.Handle(err)
			}
			
			// 根据错误类型返回相应的HTTP状态码
			if appErr := errors.GetAppError(err); appErr != nil {
				http.Error(w, appErr.Message, appErr.Code)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	}
} 