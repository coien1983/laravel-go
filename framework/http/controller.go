package http

import (
	"reflect"
	"strings"

	"laravel-go/framework/container"
	"laravel-go/framework/errors"
)

// ControllerInterface 基础控制器接口
type ControllerInterface interface {
	// 设置请求和响应
	SetRequest(request Request)
	SetResponse(response Response)
	
	// 获取请求和响应
	GetRequest() Request
	GetResponse() Response
	
	// 设置容器
	SetContainer(container container.Container)
	GetContainer() container.Container
	
	// 响应方法
	Json(data interface{}, status ...int) Response
	Text(text string, status ...int) Response
	Redirect(location string, status ...int) Response
	File(filename string) Response
	
	// 标准响应
	Success(data interface{}, message ...string) Response
	Error(message string, status ...int) Response
	ValidationError(errors interface{}) Response
	NotFound(message ...string) Response
	Unauthorized(message ...string) Response
	Forbidden(message ...string) Response
	InternalServerError(message ...string) Response
	
	// 验证方法
	Validate(request Request, rules interface{}) error
	ValidateRequest(request Request, validator interface{}) error
}

// BaseController 基础控制器实现
type BaseController struct {
	request   Request
	response  Response
	container container.Container
}

// NewBaseController 创建新的基础控制器
func NewBaseController() *BaseController {
	return &BaseController{}
}

// SetRequest 设置请求
func (c *BaseController) SetRequest(request Request) {
	c.request = request
}

// SetResponse 设置响应
func (c *BaseController) SetResponse(response Response) {
	c.response = response
}

// GetRequest 获取请求
func (c *BaseController) GetRequest() Request {
	return c.request
}

// GetResponse 获取响应
func (c *BaseController) GetResponse() Response {
	return c.response
}

// SetContainer 设置容器
func (c *BaseController) SetContainer(container container.Container) {
	c.container = container
}

// GetContainer 获取容器
func (c *BaseController) GetContainer() container.Container {
	return c.container
}

// Json 返回JSON响应
func (c *BaseController) Json(data interface{}, status ...int) Response {
	statusCode := 200
	if len(status) > 0 {
		statusCode = status[0]
	}
	return NewJsonResponse(statusCode, data)
}

// Text 返回文本响应
func (c *BaseController) Text(text string, status ...int) Response {
	statusCode := 200
	if len(status) > 0 {
		statusCode = status[0]
	}
	return NewTextResponse(statusCode, text)
}

// Redirect 返回重定向响应
func (c *BaseController) Redirect(location string, status ...int) Response {
	statusCode := 302
	if len(status) > 0 {
		statusCode = status[0]
	}
	return NewRedirectResponse(statusCode, location)
}

// File 返回文件响应
func (c *BaseController) File(filename string) Response {
	return NewFileResponse(filename)
}

// Success 返回成功响应
func (c *BaseController) Success(data interface{}, message ...string) Response {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}
	
	response := map[string]interface{}{
		"success": true,
		"message": msg,
		"data":    data,
	}
	
	return c.Json(response, 200)
}

// Error 返回错误响应
func (c *BaseController) Error(message string, status ...int) Response {
	statusCode := 400
	if len(status) > 0 {
		statusCode = status[0]
	}
	
	response := map[string]interface{}{
		"success": false,
		"message": message,
		"data":    nil,
	}
	
	return c.Json(response, statusCode)
}

// ValidationError 返回验证错误响应
func (c *BaseController) ValidationError(errors interface{}) Response {
	response := map[string]interface{}{
		"success": false,
		"message": "Validation failed",
		"errors":  errors,
	}
	
	return c.Json(response, 422)
}

// NotFound 返回404响应
func (c *BaseController) NotFound(message ...string) Response {
	msg := "Resource not found"
	if len(message) > 0 {
		msg = message[0]
	}
	return c.Error(msg, 404)
}

// Unauthorized 返回401响应
func (c *BaseController) Unauthorized(message ...string) Response {
	msg := "Unauthorized"
	if len(message) > 0 {
		msg = message[0]
	}
	return c.Error(msg, 401)
}

// Forbidden 返回403响应
func (c *BaseController) Forbidden(message ...string) Response {
	msg := "Forbidden"
	if len(message) > 0 {
		msg = message[0]
	}
	return c.Error(msg, 403)
}

// InternalServerError 返回500响应
func (c *BaseController) InternalServerError(message ...string) Response {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}
	return c.Error(msg, 500)
}

// Validate 验证请求数据
func (c *BaseController) Validate(request Request, rules interface{}) error {
	// 这里可以集成验证器
	// 暂时返回nil，实际实现中需要调用验证器
	return nil
}

// ValidateRequest 验证请求
func (c *BaseController) ValidateRequest(request Request, validator interface{}) error {
	// 创建验证请求
	validationRequest := NewValidationRequest(request, validator)
	return validationRequest.Validate()
}

// ControllerHandler 控制器处理器接口
type ControllerHandler interface {
	Handle(request Request) Response
}

// ControllerHandlerFunc 控制器处理器函数类型
type ControllerHandlerFunc func(request Request) Response

// Handle 实现ControllerHandler接口
func (f ControllerHandlerFunc) Handle(request Request) Response {
	return f(request)
}

// ControllerResolver 控制器解析器
type ControllerResolver struct {
	container container.Container
}

// NewControllerResolver 创建新的控制器解析器
func NewControllerResolver(container container.Container) *ControllerResolver {
	return &ControllerResolver{
		container: container,
	}
}

// Resolve 解析控制器
func (r *ControllerResolver) Resolve(controller interface{}) (ControllerHandler, error) {
	// 如果controller是函数，直接返回
	if handler, ok := controller.(ControllerHandlerFunc); ok {
		return handler, nil
	}
	
	// 如果controller是字符串，从容器中解析
	if controllerName, ok := controller.(string); ok {
		return r.resolveFromContainer(controllerName)
	}
	
	// 如果controller是结构体，创建实例
	if reflect.TypeOf(controller).Kind() == reflect.Struct {
		return r.resolveStruct(controller)
	}
	
	return nil, errors.New("unsupported controller type")
}

// resolveFromContainer 从容器中解析控制器
func (r *ControllerResolver) resolveFromContainer(controllerName string) (ControllerHandler, error) {
	// 解析控制器名称，格式: Controller@method
	parts := strings.Split(controllerName, "@")
	if len(parts) != 2 {
		return nil, errors.New("invalid controller format, expected: Controller@method")
	}
	
	// 从容器中获取控制器实例
	// 注意：这里需要修改容器的使用方式，因为字符串不能直接作为抽象类型
	// 暂时返回错误，实际实现中需要更好的容器管理
	return nil, errors.New("container resolution not implemented for string keys")
}

// resolveStruct 解析结构体控制器
func (r *ControllerResolver) resolveStruct(controller interface{}) (ControllerHandler, error) {
	// 创建处理器
	handler := func(request Request) Response {
		// 设置请求到控制器
		if controller, ok := controller.(ControllerInterface); ok {
			controller.SetRequest(request)
			controller.SetContainer(r.container)
		}
		
		// 这里可以添加更多的处理逻辑
		return NewJsonResponse(500, map[string]interface{}{
			"success": false,
			"message": "Struct controller not implemented",
		})
	}
	
	return ControllerHandlerFunc(handler), nil
}

// ControllerMiddleware 控制器中间件
type ControllerMiddleware struct {
	resolver *ControllerResolver
}

// NewControllerMiddleware 创建新的控制器中间件
func NewControllerMiddleware(resolver *ControllerResolver) *ControllerMiddleware {
	return &ControllerMiddleware{
		resolver: resolver,
	}
}

// Process 处理请求
func (m *ControllerMiddleware) Process(request Request, next func(Request) Response) Response {
	// 这里可以添加控制器相关的中间件逻辑
	// 比如依赖注入、请求验证等
	
	return next(request)
}

// ControllerBuilder 控制器构建器
type ControllerBuilder struct {
	container container.Container
	resolver  *ControllerResolver
}

// NewControllerBuilder 创建新的控制器构建器
func NewControllerBuilder(container container.Container) *ControllerBuilder {
	resolver := NewControllerResolver(container)
	return &ControllerBuilder{
		container: container,
		resolver:  resolver,
	}
}

// Register 注册控制器
func (b *ControllerBuilder) Register(controller interface{}) error {
	// 注册控制器到容器
	b.container.Bind(controller, controller)
	return nil
}

// RegisterSingleton 注册单例控制器
func (b *ControllerBuilder) RegisterSingleton(controller interface{}) error {
	// 注册单例控制器到容器
	b.container.BindSingleton(controller, controller)
	return nil
}

// Resolve 解析控制器
func (b *ControllerBuilder) Resolve(controller interface{}) (ControllerHandler, error) {
	return b.resolver.Resolve(controller)
}

// CreateHandler 创建处理器
func (b *ControllerBuilder) CreateHandler(controller interface{}) (ControllerHandler, error) {
	return b.Resolve(controller)
} 