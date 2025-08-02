package errors

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// AppError 应用错误类型
type AppError struct {
	Code    int
	Message string
	Err     error
	Stack   []string
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithStack 添加堆栈信息
func (e *AppError) WithStack() *AppError {
	e.Stack = getStackTrace()
	return e
}

// WithError 包装原始错误
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// WithMessage 设置错误消息
func (e *AppError) WithMessage(message string) *AppError {
	e.Message = message
	return e
}

// WithCode 设置错误代码
func (e *AppError) WithCode(code int) *AppError {
	e.Code = code
	return e
}

// 预定义错误类型
var (
	ErrNotFound           = &AppError{Code: 404, Message: "Resource not found"}
	ErrUnauthorized       = &AppError{Code: 401, Message: "Unauthorized"}
	ErrForbidden          = &AppError{Code: 403, Message: "Forbidden"}
	ErrValidation         = &AppError{Code: 422, Message: "Validation failed"}
	ErrInternalServer     = &AppError{Code: 500, Message: "Internal server error"}
	ErrBadRequest         = &AppError{Code: 400, Message: "Bad request"}
	ErrMethodNotAllowed   = &AppError{Code: 405, Message: "Method not allowed"}
	ErrConflict           = &AppError{Code: 409, Message: "Conflict"}
	ErrTooManyRequests    = &AppError{Code: 429, Message: "Too many requests"}
	ErrServiceUnavailable = &AppError{Code: 503, Message: "Service unavailable"}
)

// New 创建新的应用错误
func New(message string) *AppError {
	return &AppError{
		Code:    500,
		Message: message,
		Stack:   getStackTrace(),
	}
}

// NewWithCode 创建带错误代码的应用错误
func NewWithCode(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Stack:   getStackTrace(),
	}
}

// Wrap 包装错误
func Wrap(err error, message string) *AppError {
	return &AppError{
		Code:    500,
		Message: message,
		Err:     err,
		Stack:   getStackTrace(),
	}
}

// WrapWithCode 包装错误并设置代码
func WrapWithCode(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Stack:   getStackTrace(),
	}
}

// IsAppError 检查是否为应用错误
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError 获取应用错误
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}

// getStackTrace 获取堆栈跟踪
func getStackTrace() []string {
	var stack []string

	// 跳过前3个调用帧（getStackTrace, New/Wrap, 调用者）
	for i := 3; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		funcName := runtime.FuncForPC(pc).Name()

		// 跳过标准库和框架内部函数
		if strings.Contains(funcName, "runtime.") ||
			strings.Contains(funcName, "framework/errors.") {
			continue
		}

		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, funcName))

		// 限制堆栈深度
		if len(stack) >= 10 {
			break
		}
	}

	return stack
}

// ErrorHandler 错误处理器接口
type ErrorHandler interface {
	Handle(err error) error
	Log(err error)
	Report(err error)
}

// DefaultErrorHandler 默认错误处理器
type DefaultErrorHandler struct {
	logger Logger
}

// Logger 日志接口
type Logger interface {
	Error(message string, context map[string]interface{})
	Warning(message string, context map[string]interface{})
	Info(message string, context map[string]interface{})
	Debug(message string, context map[string]interface{})
}

// NewDefaultErrorHandler 创建默认错误处理器
func NewDefaultErrorHandler(logger Logger) *DefaultErrorHandler {
	return &DefaultErrorHandler{
		logger: logger,
	}
}

// Handle 处理错误
func (h *DefaultErrorHandler) Handle(err error) error {
	defer func() {
		if r := recover(); r != nil {
			// 记录错误处理器本身的panic
			if h.logger != nil {
				h.logger.Error("Error handler panic", map[string]interface{}{
					"panic": r,
					"stack": getStackTrace(),
				})
			}
		}
	}()

	if err == nil {
		return nil
	}

	// 记录错误
	h.Log(err)

	// 如果是应用错误，直接返回
	if IsAppError(err) {
		return err
	}

	// 包装为应用错误
	return Wrap(err, "An unexpected error occurred")
}

// Log 记录错误
func (h *DefaultErrorHandler) Log(err error) {
	if h.logger == nil {
		return
	}

	context := map[string]interface{}{
		"error":       err.Error(),
		"timestamp":   time.Now(),
		"error_type":  reflect.TypeOf(err).String(),
	}

	if appErr := GetAppError(err); appErr != nil {
		context["code"] = appErr.Code
		context["stack"] = appErr.Stack
		if appErr.Err != nil {
			context["cause"] = appErr.Err.Error()
		}
	}

	h.logger.Error("Application error", context)
}

// Report 报告错误（用于外部错误报告服务）
func (h *DefaultErrorHandler) Report(err error) {
	defer func() {
		if r := recover(); r != nil {
			// 记录报告错误时的panic
			if h.logger != nil {
				h.logger.Error("Error report panic", map[string]interface{}{
					"panic": r,
					"stack": getStackTrace(),
				})
			}
		}
	}()
	
	// 这里可以集成外部错误报告服务，如 Sentry
	// 暂时只记录日志
	h.Log(err)
}

// ErrorMiddleware 错误处理中间件
type ErrorMiddleware struct {
	handler ErrorHandler
}

// NewErrorMiddleware 创建错误处理中间件
func NewErrorMiddleware(handler ErrorHandler) *ErrorMiddleware {
	return &ErrorMiddleware{
		handler: handler,
	}
}

// Handle 处理请求中的错误
func (m *ErrorMiddleware) Handle(request interface{}, next func() interface{}) interface{} {
	defer func() {
		if r := recover(); r != nil {
			// 将 panic 转换为错误
			var err error
			if errVal, ok := r.(error); ok {
				err = errVal
			} else {
				err = New(fmt.Sprintf("Panic: %v", r))
			}

			// 处理错误
			m.handler.Handle(err)
		}
	}()

	// 执行下一个处理器
	result := next()

	// 检查结果是否为错误
	if err, ok := result.(error); ok {
		return m.handler.Handle(err)
	}

	return result
}

// SafeExecute 安全执行函数
func SafeExecute(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if errVal, ok := r.(error); ok {
				err = errVal
			} else {
				err = New(fmt.Sprintf("Panic: %v", r))
			}
		}
	}()
	
	return fn()
}

// SafeExecuteWithContext 带上下文的安全执行函数
func SafeExecuteWithContext(ctx context.Context, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if errVal, ok := r.(error); ok {
				err = errVal
			} else {
				err = New(fmt.Sprintf("Panic: %v", r))
			}
		}
	}()
	
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	
	return fn()
}

// 注意：ValidationError 和 ValidationErrors 已在 error_types.go 中定义
