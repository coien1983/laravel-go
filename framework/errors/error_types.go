package errors

import (
	"fmt"
	"strings"
	"time"
)

// ErrorCode 错误码类型
type ErrorCode string

// 系统级错误码
const (
	// 成功
	ErrorCodeSuccess ErrorCode = "SUCCESS"

	// 客户端错误 (4xx)
	ErrorCodeBadRequest           ErrorCode = "BAD_REQUEST"
	ErrorCodeUnauthorized         ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden            ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound             ErrorCode = "NOT_FOUND"
	ErrorCodeMethodNotAllowed     ErrorCode = "METHOD_NOT_ALLOWED"
	ErrorCodeConflict             ErrorCode = "CONFLICT"
	ErrorCodeValidationFailed     ErrorCode = "VALIDATION_FAILED"
	ErrorCodeTooManyRequests      ErrorCode = "TOO_MANY_REQUESTS"
	ErrorCodeRequestTimeout       ErrorCode = "REQUEST_TIMEOUT"
	ErrorCodeUnsupportedMediaType ErrorCode = "UNSUPPORTED_MEDIA_TYPE"

	// 服务器错误 (5xx)
	ErrorCodeInternalServer       ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeNotImplemented       ErrorCode = "NOT_IMPLEMENTED"
	ErrorCodeServiceUnavailable   ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeGatewayTimeout       ErrorCode = "GATEWAY_TIMEOUT"
	ErrorCodeDatabaseError        ErrorCode = "DATABASE_ERROR"
	ErrorCodeCacheError           ErrorCode = "CACHE_ERROR"
	ErrorCodeExternalServiceError ErrorCode = "EXTERNAL_SERVICE_ERROR"

	// 业务错误
	ErrorCodeBusinessLogic     ErrorCode = "BUSINESS_LOGIC_ERROR"
	ErrorCodeResourceExhausted ErrorCode = "RESOURCE_EXHAUSTED"
	ErrorCodeQuotaExceeded     ErrorCode = "QUOTA_EXCEEDED"
	ErrorCodeRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"
)

// ErrorSeverity 错误严重程度
type ErrorSeverity string

const (
	ErrorSeverityLow      ErrorSeverity = "LOW"
	ErrorSeverityMedium   ErrorSeverity = "MEDIUM"
	ErrorSeverityHigh     ErrorSeverity = "HIGH"
	ErrorSeverityCritical ErrorSeverity = "CRITICAL"
)

// ErrorCategory 错误分类
type ErrorCategory string

const (
	ErrorCategorySystem     ErrorCategory = "SYSTEM"
	ErrorCategoryBusiness   ErrorCategory = "BUSINESS"
	ErrorCategoryValidation ErrorCategory = "VALIDATION"
	ErrorCategorySecurity   ErrorCategory = "SECURITY"
	ErrorCategoryNetwork    ErrorCategory = "NETWORK"
	ErrorCategoryDatabase   ErrorCategory = "DATABASE"
	ErrorCategoryCache      ErrorCategory = "CACHE"
	ErrorCategoryExternal   ErrorCategory = "EXTERNAL"
)

// BusinessError 业务错误
type BusinessError struct {
	Code      ErrorCode     `json:"code"`
	Message   string        `json:"message"`
	Details   interface{}   `json:"details,omitempty"`
	Severity  ErrorSeverity `json:"severity"`
	Category  ErrorCategory `json:"category"`
	Timestamp time.Time     `json:"timestamp"`
	RequestID string        `json:"request_id,omitempty"`
	UserID    string        `json:"user_id,omitempty"`
	Stack     []string      `json:"stack,omitempty"`
	Cause     error         `json:"-"`
}

// Error 实现 error 接口
func (e *BusinessError) Error() string {
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口
func (e *BusinessError) Unwrap() error {
	return e.Cause
}

// WithDetails 添加详细信息
func (e *BusinessError) WithDetails(details interface{}) *BusinessError {
	e.Details = details
	return e
}

// WithSeverity 设置严重程度
func (e *BusinessError) WithSeverity(severity ErrorSeverity) *BusinessError {
	e.Severity = severity
	return e
}

// WithCategory 设置错误分类
func (e *BusinessError) WithCategory(category ErrorCategory) *BusinessError {
	e.Category = category
	return e
}

// WithRequestID 设置请求ID
func (e *BusinessError) WithRequestID(requestID string) *BusinessError {
	e.RequestID = requestID
	return e
}

// WithUserID 设置用户ID
func (e *BusinessError) WithUserID(userID string) *BusinessError {
	e.UserID = userID
	return e
}

// WithStack 添加堆栈信息
func (e *BusinessError) WithStack() *BusinessError {
	e.Stack = getStackTrace()
	return e
}

// WithCause 设置原始错误
func (e *BusinessError) WithCause(cause error) *BusinessError {
	e.Cause = cause
	return e
}

// NewBusinessError 创建业务错误
func NewBusinessError(code ErrorCode, message string) *BusinessError {
	return &BusinessError{
		Code:      code,
		Message:   message,
		Severity:  ErrorSeverityMedium,
		Category:  ErrorCategoryBusiness,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// WrapBusinessError 包装错误为业务错误
func WrapBusinessError(err error, code ErrorCode, message string) *BusinessError {
	return &BusinessError{
		Code:      code,
		Message:   message,
		Severity:  ErrorSeverityMedium,
		Category:  ErrorCategoryBusiness,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
		Cause:     err,
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Field     string      `json:"field"`
	Message   string      `json:"message"`
	Value     interface{} `json:"value,omitempty"`
	Rule      string      `json:"rule,omitempty"`
	Expected  interface{} `json:"expected,omitempty"`
	Actual    interface{} `json:"actual,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// Error 实现 error 接口
func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation failed for field '%s': %s", e.Field, e.Message)
}

// WithValue 设置字段值
func (e *ValidationError) WithValue(value interface{}) *ValidationError {
	e.Value = value
	return e
}

// WithRule 设置验证规则
func (e *ValidationError) WithRule(rule string) *ValidationError {
	e.Rule = rule
	return e
}

// WithExpected 设置期望值
func (e *ValidationError) WithExpected(expected interface{}) *ValidationError {
	e.Expected = expected
	return e
}

// WithActual 设置实际值
func (e *ValidationError) WithActual(actual interface{}) *ValidationError {
	e.Actual = actual
	return e
}

// WithRequestID 设置请求ID
func (e *ValidationError) WithRequestID(requestID string) *ValidationError {
	e.RequestID = requestID
	return e
}

// NewValidationError 创建验证错误
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:     field,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// ValidationErrors 验证错误集合
type ValidationErrors []*ValidationError

// Error 实现 error 接口
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "No validation errors"
	}

	messages := make([]string, len(e))
	for i, err := range e {
		messages[i] = err.Error()
	}

	return strings.Join(messages, "; ")
}

// Add 添加验证错误
func (e *ValidationErrors) Add(field, message string) {
	*e = append(*e, NewValidationError(field, message))
}

// AddWithValue 添加带值的验证错误
func (e *ValidationErrors) AddWithValue(field, message string, value interface{}) {
	err := NewValidationError(field, message).WithValue(value)
	*e = append(*e, err)
}

// AddWithRule 添加带规则的验证错误
func (e *ValidationErrors) AddWithRule(field, message, rule string) {
	err := NewValidationError(field, message).WithRule(rule)
	*e = append(*e, err)
}

// HasErrors 检查是否有错误
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// GetErrors 获取所有错误
func (e ValidationErrors) GetErrors() []*ValidationError {
	return e
}

// GetErrorsByField 根据字段获取错误
func (e ValidationErrors) GetErrorsByField(field string) []*ValidationError {
	var errors []*ValidationError
	for _, err := range e {
		if err.Field == field {
			errors = append(errors, err)
		}
	}
	return errors
}

// ToMap 转换为字段错误映射
func (e ValidationErrors) ToMap() map[string][]string {
	result := make(map[string][]string)
	for _, err := range e {
		result[err.Field] = append(result[err.Field], err.Message)
	}
	return result
}

// SecurityError 安全错误
type SecurityError struct {
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	IP        string    `json:"ip,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	Action    string    `json:"action,omitempty"`
	Resource  string    `json:"resource,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id,omitempty"`
	Stack     []string  `json:"stack,omitempty"`
	Cause     error     `json:"-"`
}

// Error 实现 error 接口
func (e *SecurityError) Error() string {
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口
func (e *SecurityError) Unwrap() error {
	return e.Cause
}

// WithIP 设置IP地址
func (e *SecurityError) WithIP(ip string) *SecurityError {
	e.IP = ip
	return e
}

// WithUserAgent 设置用户代理
func (e *SecurityError) WithUserAgent(userAgent string) *SecurityError {
	e.UserAgent = userAgent
	return e
}

// WithUserID 设置用户ID
func (e *SecurityError) WithUserID(userID string) *SecurityError {
	e.UserID = userID
	return e
}

// WithAction 设置操作
func (e *SecurityError) WithAction(action string) *SecurityError {
	e.Action = action
	return e
}

// WithResource 设置资源
func (e *SecurityError) WithResource(resource string) *SecurityError {
	e.Resource = resource
	return e
}

// WithRequestID 设置请求ID
func (e *SecurityError) WithRequestID(requestID string) *SecurityError {
	e.RequestID = requestID
	return e
}

// WithStack 添加堆栈信息
func (e *SecurityError) WithStack() *SecurityError {
	e.Stack = getStackTrace()
	return e
}

// WithCause 设置原始错误
func (e *SecurityError) WithCause(cause error) *SecurityError {
	e.Cause = cause
	return e
}

// NewSecurityError 创建安全错误
func NewSecurityError(code ErrorCode, message string) *SecurityError {
	return &SecurityError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// DatabaseError 数据库错误
type DatabaseError struct {
	Code       ErrorCode   `json:"code"`
	Message    string      `json:"message"`
	Operation  string      `json:"operation,omitempty"`
	Table      string      `json:"table,omitempty"`
	Query      string      `json:"query,omitempty"`
	Parameters interface{} `json:"parameters,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  string      `json:"request_id,omitempty"`
	Stack      []string    `json:"stack,omitempty"`
	Cause      error       `json:"-"`
}

// Error 实现 error 接口
func (e *DatabaseError) Error() string {
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口
func (e *DatabaseError) Unwrap() error {
	return e.Cause
}

// WithOperation 设置操作
func (e *DatabaseError) WithOperation(operation string) *DatabaseError {
	e.Operation = operation
	return e
}

// WithTable 设置表名
func (e *DatabaseError) WithTable(table string) *DatabaseError {
	e.Table = table
	return e
}

// WithQuery 设置查询
func (e *DatabaseError) WithQuery(query string) *DatabaseError {
	e.Query = query
	return e
}

// WithParameters 设置参数
func (e *DatabaseError) WithParameters(parameters interface{}) *DatabaseError {
	e.Parameters = parameters
	return e
}

// WithRequestID 设置请求ID
func (e *DatabaseError) WithRequestID(requestID string) *DatabaseError {
	e.RequestID = requestID
	return e
}

// WithStack 添加堆栈信息
func (e *DatabaseError) WithStack() *DatabaseError {
	e.Stack = getStackTrace()
	return e
}

// WithCause 设置原始错误
func (e *DatabaseError) WithCause(cause error) *DatabaseError {
	e.Cause = cause
	return e
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(code ErrorCode, message string) *DatabaseError {
	return &DatabaseError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Stack:     getStackTrace(),
	}
}

// ExternalServiceError 外部服务错误
type ExternalServiceError struct {
	Code         ErrorCode     `json:"code"`
	Message      string        `json:"message"`
	ServiceName  string        `json:"service_name"`
	Endpoint     string        `json:"endpoint,omitempty"`
	Method       string        `json:"method,omitempty"`
	StatusCode   int           `json:"status_code,omitempty"`
	ResponseBody string        `json:"response_body,omitempty"`
	Timeout      time.Duration `json:"timeout,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
	RequestID    string        `json:"request_id,omitempty"`
	Stack        []string      `json:"stack,omitempty"`
	Cause        error         `json:"-"`
}

// Error 实现 error 接口
func (e *ExternalServiceError) Error() string {
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口
func (e *ExternalServiceError) Unwrap() error {
	return e.Cause
}

// WithEndpoint 设置端点
func (e *ExternalServiceError) WithEndpoint(endpoint string) *ExternalServiceError {
	e.Endpoint = endpoint
	return e
}

// WithMethod 设置方法
func (e *ExternalServiceError) WithMethod(method string) *ExternalServiceError {
	e.Method = method
	return e
}

// WithStatusCode 设置状态码
func (e *ExternalServiceError) WithStatusCode(statusCode int) *ExternalServiceError {
	e.StatusCode = statusCode
	return e
}

// WithResponseBody 设置响应体
func (e *ExternalServiceError) WithResponseBody(responseBody string) *ExternalServiceError {
	e.ResponseBody = responseBody
	return e
}

// WithTimeout 设置超时时间
func (e *ExternalServiceError) WithTimeout(timeout time.Duration) *ExternalServiceError {
	e.Timeout = timeout
	return e
}

// WithRequestID 设置请求ID
func (e *ExternalServiceError) WithRequestID(requestID string) *ExternalServiceError {
	e.RequestID = requestID
	return e
}

// WithStack 添加堆栈信息
func (e *ExternalServiceError) WithStack() *ExternalServiceError {
	e.Stack = getStackTrace()
	return e
}

// WithCause 设置原始错误
func (e *ExternalServiceError) WithCause(cause error) *ExternalServiceError {
	e.Cause = cause
	return e
}

// NewExternalServiceError 创建外部服务错误
func NewExternalServiceError(serviceName, message string) *ExternalServiceError {
	return &ExternalServiceError{
		Code:        ErrorCodeExternalServiceError,
		Message:     message,
		ServiceName: serviceName,
		Timestamp:   time.Now(),
		Stack:       getStackTrace(),
	}
}

// 预定义错误实例
var (
	// 业务错误
	ErrUserNotFound            = NewBusinessError(ErrorCodeNotFound, "User not found")
	ErrUserAlreadyExists       = NewBusinessError(ErrorCodeConflict, "User already exists")
	ErrInvalidCredentials      = NewBusinessError(ErrorCodeUnauthorized, "Invalid credentials")
	ErrInsufficientPermissions = NewBusinessError(ErrorCodeForbidden, "Insufficient permissions")
	ErrResourceNotFound        = NewBusinessError(ErrorCodeNotFound, "Resource not found")
	ErrResourceConflict        = NewBusinessError(ErrorCodeConflict, "Resource conflict")
	ErrRateLimitExceeded       = NewBusinessError(ErrorCodeRateLimitExceeded, "Rate limit exceeded")
	ErrQuotaExceeded           = NewBusinessError(ErrorCodeQuotaExceeded, "Quota exceeded")

	// 安全错误
	ErrAuthenticationFailed = NewSecurityError(ErrorCodeUnauthorized, "Authentication failed")
	ErrAuthorizationFailed  = NewSecurityError(ErrorCodeForbidden, "Authorization failed")
	ErrInvalidToken         = NewSecurityError(ErrorCodeUnauthorized, "Invalid token")
	ErrTokenExpired         = NewSecurityError(ErrorCodeUnauthorized, "Token expired")
	ErrAccessDenied         = NewSecurityError(ErrorCodeForbidden, "Access denied")

	// 数据库错误
	ErrDatabaseConnection  = NewDatabaseError(ErrorCodeDatabaseError, "Database connection failed")
	ErrDatabaseQuery       = NewDatabaseError(ErrorCodeDatabaseError, "Database query failed")
	ErrDatabaseTransaction = NewDatabaseError(ErrorCodeDatabaseError, "Database transaction failed")
	ErrDatabaseConstraint  = NewDatabaseError(ErrorCodeDatabaseError, "Database constraint violation")

	// 外部服务错误
	ErrExternalServiceUnavailable = NewExternalServiceError("external", "External service unavailable")
	ErrExternalServiceTimeout     = NewExternalServiceError("external", "External service timeout")
	ErrExternalServiceError       = NewExternalServiceError("external", "External service error")
)
