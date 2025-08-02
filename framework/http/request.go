package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"laravel-go/framework/errors"
)

// Request HTTP请求接口
type Request interface {
	// 获取请求方法
	Method() string
	// 获取请求路径
	Path() string
	// 获取请求参数
	Param(key string) string
	// 获取查询参数
	Query(key string) string
	// 获取请求头
	Header(key string) string
	// 获取客户端IP
	IP() string
	// 获取请求体
	Body() []byte
	// 解析JSON请求体
	Json(v interface{}) error
	// 获取表单数据
	Form(key string) string
	// 获取文件
	File(key string) ([]byte, string, error)
	// 获取Cookie
	Cookie(key string) string
	// 获取原始请求
	Raw() *http.Request
}

// request HTTP请求实现
type request struct {
	req  *http.Request
	body []byte
}

// NewRequest 创建新的请求
func NewRequest(req *http.Request) Request {
	return &request{
		req: req,
	}
}

// Method 获取请求方法
func (r *request) Method() string {
	return r.req.Method
}

// Path 获取请求路径
func (r *request) Path() string {
	return r.req.URL.Path
}

// Param 获取路径参数 (需要从路由中获取，这里暂时返回空)
func (r *request) Param(key string) string {
	// 这里需要从路由参数中获取，暂时返回空
	// 在实际使用中，需要从路由匹配结果中获取
	return ""
}

// Query 获取查询参数
func (r *request) Query(key string) string {
	return r.req.URL.Query().Get(key)
}

// Header 获取请求头
func (r *request) Header(key string) string {
	return r.req.Header.Get(key)
}

// IP 获取客户端IP
func (r *request) IP() string {
	// 尝试从各种头部获取真实IP
	if ip := r.req.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.req.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.req.Header.Get("X-Client-IP"); ip != "" {
		return ip
	}
	return r.req.RemoteAddr
}

// Body 获取请求体
func (r *request) Body() []byte {
	if r.body != nil {
		return r.body
	}

	body, err := io.ReadAll(r.req.Body)
	if err != nil {
		return nil
	}
	r.body = body
	return body
}

// Json 解析JSON请求体
func (r *request) Json(v interface{}) error {
	body := r.Body()
	if body == nil {
		return errors.New("failed to read request body")
	}
	return json.Unmarshal(body, v)
}

// Form 获取表单数据
func (r *request) Form(key string) string {
	return r.req.FormValue(key)
}

// File 获取文件
func (r *request) File(key string) ([]byte, string, error) {
	file, header, err := r.req.FormFile(key)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, "", err
	}

	return data, header.Filename, nil
}

// Cookie 获取Cookie
func (r *request) Cookie(key string) string {
	cookie, err := r.req.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// Raw 获取原始请求
func (r *request) Raw() *http.Request {
	return r.req
}

// RequestValidator 请求验证器
type RequestValidator interface {
	Validate(request Request) error
}

// ValidationRequest 带验证的请求
type ValidationRequest struct {
	Request
	validator interface{}
}

// NewValidationRequest 创建带验证的请求
func NewValidationRequest(req Request, validator interface{}) *ValidationRequest {
	return &ValidationRequest{
		Request:   req,
		validator: validator,
	}
}

// Validate 验证请求
func (vr *ValidationRequest) Validate() error {
	// 暂时返回nil，后续实现验证逻辑
	return nil
}

// RequestBuilder 请求构建器
type RequestBuilder struct {
	req *http.Request
}

// NewRequestBuilder 创建请求构建器
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{}
}

// Method 设置请求方法
func (rb *RequestBuilder) Method(method string) *RequestBuilder {
	if rb.req == nil {
		rb.req = &http.Request{}
	}
	rb.req.Method = method
	return rb
}

// URL 设置请求URL
func (rb *RequestBuilder) URL(url string) *RequestBuilder {
	if rb.req == nil {
		rb.req = &http.Request{}
	}
	// 暂时跳过URL解析，后续实现
	return rb
}

// Header 设置请求头
func (rb *RequestBuilder) Header(key, value string) *RequestBuilder {
	if rb.req == nil {
		rb.req = &http.Request{}
	}
	if rb.req.Header == nil {
		rb.req.Header = make(http.Header)
	}
	rb.req.Header.Set(key, value)
	return rb
}

// Body 设置请求体
func (rb *RequestBuilder) Body(body io.Reader) *RequestBuilder {
	if rb.req == nil {
		rb.req = &http.Request{}
	}
	rb.req.Body = io.NopCloser(body)
	return rb
}

// Build 构建请求
func (rb *RequestBuilder) Build() *http.Request {
	return rb.req
}

// RequestMiddleware 请求中间件接口
type RequestMiddleware interface {
	Process(request Request, next func(Request) Response) Response
}

// RequestMiddlewareFunc 请求中间件函数
type RequestMiddlewareFunc func(request Request, next func(Request) Response) Response

// Process 处理请求
func (rmf RequestMiddlewareFunc) Process(request Request, next func(Request) Response) Response {
	return rmf(request, next)
}

// RequestPipeline 请求管道
type RequestPipeline struct {
	middlewares []RequestMiddleware
}

// NewRequestPipeline 创建请求管道
func NewRequestPipeline() *RequestPipeline {
	return &RequestPipeline{
		middlewares: make([]RequestMiddleware, 0),
	}
}

// Use 添加中间件
func (rp *RequestPipeline) Use(middleware RequestMiddleware) *RequestPipeline {
	rp.middlewares = append(rp.middlewares, middleware)
	return rp
}

// Process 处理请求
func (rp *RequestPipeline) Process(request Request, handler func(Request) Response) Response {
	// 创建处理器函数
	handlerFunc := func(req Request) Response {
		return handler(req)
	}

	// 构建中间件链
	next := handlerFunc
	for i := len(rp.middlewares) - 1; i >= 0; i-- {
		currentMiddleware := rp.middlewares[i]
		currentNext := next
		next = func(req Request) Response {
			return currentMiddleware.Process(req, currentNext)
		}
	}

	// 执行中间件链
	return next(request)
}
