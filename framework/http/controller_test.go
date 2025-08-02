package http

import (
	"net/http"
	"testing"

	"laravel-go/framework/container"
)

// TestController 测试控制器
type TestController struct {
	*BaseController
}

// NewTestController 创建测试控制器
func NewTestController() *TestController {
	return &TestController{
		BaseController: NewBaseController(),
	}
}

// Index 首页方法
func (c *TestController) Index(request Request) Response {
	return c.Success(map[string]interface{}{
		"message": "Hello from TestController",
		"method":  request.Method(),
		"path":    request.Path(),
	})
}

// Show 显示方法
func (c *TestController) Show(request Request) Response {
	id := request.Param("id")
	if id == "" {
		return c.NotFound("ID parameter is required")
	}
	
	return c.Success(map[string]interface{}{
		"id":      id,
		"message": "Show method called",
	})
}

// Store 存储方法
func (c *TestController) Store(request Request) Response {
	// 解析JSON请求体
	var data map[string]interface{}
	if err := request.Json(&data); err != nil {
		return c.Error("Invalid JSON data")
	}
	
	return c.Success(data, "Data stored successfully")
}

// Update 更新方法
func (c *TestController) Update(request Request) Response {
	id := request.Param("id")
	if id == "" {
		return c.NotFound("ID parameter is required")
	}
	
	var data map[string]interface{}
	if err := request.Json(&data); err != nil {
		return c.Error("Invalid JSON data")
	}
	
	data["id"] = id
	return c.Success(data, "Data updated successfully")
}

// Delete 删除方法
func (c *TestController) Delete(request Request) Response {
	id := request.Param("id")
	if id == "" {
		return c.NotFound("ID parameter is required")
	}
	
	return c.Success(map[string]interface{}{
		"id":      id,
		"message": "Data deleted successfully",
	})
}

// TestBaseController 测试基础控制器
func TestBaseController(t *testing.T) {
	controller := NewBaseController()
	
	// 测试设置和获取请求
	req, _ := http.NewRequest("GET", "/test", nil)
	request := NewRequest(req)
	controller.SetRequest(request)
	
	if controller.GetRequest() != request {
		t.Error("Request not set correctly")
	}
	
	// 测试设置和获取容器
	container := container.NewContainer()
	controller.SetContainer(container)
	
	if controller.GetContainer() != container {
		t.Error("Container not set correctly")
	}
}

// TestControllerResponses 测试控制器响应方法
func TestControllerResponses(t *testing.T) {
	controller := NewBaseController()
	
	// 测试JSON响应
	response := controller.Json(map[string]string{"key": "value"})
	if response.Status() != 200 {
		t.Error("JSON response status should be 200")
	}
	
	// 测试文本响应
	response = controller.Text("Hello World")
	if response.Status() != 200 {
		t.Error("Text response status should be 200")
	}
	
	// 测试重定向响应
	response = controller.Redirect("/new-location")
	if response.Status() != 302 {
		t.Error("Redirect response status should be 302")
	}
	
	// 测试成功响应
	response = controller.Success(map[string]string{"data": "test"})
	if response.Status() != 200 {
		t.Error("Success response status should be 200")
	}
	
	// 测试错误响应
	response = controller.Error("Something went wrong")
	if response.Status() != 400 {
		t.Error("Error response status should be 400")
	}
	
	// 测试404响应
	response = controller.NotFound("Resource not found")
	if response.Status() != 404 {
		t.Error("NotFound response status should be 404")
	}
	
	// 测试401响应
	response = controller.Unauthorized("Unauthorized access")
	if response.Status() != 401 {
		t.Error("Unauthorized response status should be 401")
	}
	
	// 测试403响应
	response = controller.Forbidden("Access forbidden")
	if response.Status() != 403 {
		t.Error("Forbidden response status should be 403")
	}
	
	// 测试500响应
	response = controller.InternalServerError("Server error")
	if response.Status() != 500 {
		t.Error("InternalServerError response status should be 500")
	}
}

// TestControllerResolver 测试控制器解析器
func TestControllerResolver(t *testing.T) {
	container := container.NewContainer()
	resolver := NewControllerResolver(container)
	
	// 测试函数处理器
	handler := ControllerHandlerFunc(func(request Request) Response {
		return NewJsonResponse(200, map[string]string{
			"message": "Handler function called",
		})
	})
	
	// 解析函数处理器
	resolvedHandler, err := resolver.Resolve(handler)
	if err != nil {
		t.Fatalf("Failed to resolve handler: %v", err)
	}
	
	// 创建测试请求
	req, _ := http.NewRequest("GET", "/test", nil)
	request := NewRequest(req)
	
	// 执行处理器
	response := resolvedHandler.Handle(request)
	if response.Status() != 200 {
		t.Error("Handler should return 200 status")
	}
}

// TestControllerHandlerFunc 测试控制器处理器函数
func TestControllerHandlerFunc(t *testing.T) {
	// 创建处理器函数
	handler := ControllerHandlerFunc(func(request Request) Response {
		return NewJsonResponse(200, map[string]string{
			"message": "Handler function called",
		})
	})
	
	// 创建测试请求
	req, _ := http.NewRequest("GET", "/test", nil)
	request := NewRequest(req)
	
	// 执行处理器
	response := handler.Handle(request)
	if response.Status() != 200 {
		t.Error("Handler function should return 200 status")
	}
}

// TestControllerBuilder 测试控制器构建器
func TestControllerBuilder(t *testing.T) {
	container := container.NewContainer()
	builder := NewControllerBuilder(container)
	
	// 注册控制器
	testController := NewTestController()
	err := builder.Register(testController)
	if err != nil {
		t.Fatalf("Failed to register controller: %v", err)
	}
	
	// 测试函数处理器
	handler := ControllerHandlerFunc(func(request Request) Response {
		return NewJsonResponse(200, map[string]string{
			"message": "Builder handler called",
		})
	})
	
	// 解析函数处理器
	resolvedHandler, err := builder.Resolve(handler)
	if err != nil {
		t.Fatalf("Failed to resolve handler: %v", err)
	}
	
	// 创建测试请求
	req, _ := http.NewRequest("GET", "/test", nil)
	request := NewRequest(req)
	
	// 执行处理器
	response := resolvedHandler.Handle(request)
	if response.Status() != 200 {
		t.Error("Builder should create valid handler")
	}
}

// TestControllerIntegration 测试控制器集成
func TestControllerIntegration(t *testing.T) {
	container := container.NewContainer()
	builder := NewControllerBuilder(container)
	
	// 注册控制器
	testController := NewTestController()
	builder.Register(testController)
	
	// 测试不同的控制器方法
	testCases := []struct {
		name     string
		method   string
		path     string
		expected int
	}{
		{"Index", "GET", "/", 200},
		{"Show", "GET", "/1", 200},
		{"Store", "POST", "/", 200},
		{"Update", "PUT", "/1", 200},
		{"Delete", "DELETE", "/1", 200},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建请求
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			request := NewRequest(req)
			
			// 创建处理器函数
			handler := ControllerHandlerFunc(func(req Request) Response {
				return NewJsonResponse(200, map[string]string{
					"method": tc.method,
					"path":   tc.path,
				})
			})
			
			// 解析处理器
			resolvedHandler, err := builder.Resolve(handler)
			if err != nil {
				t.Fatalf("Failed to resolve handler: %v", err)
			}
			
			// 执行处理器
			response := resolvedHandler.Handle(request)
			if response.Status() != tc.expected {
				t.Errorf("Expected status %d, got %d", tc.expected, response.Status())
			}
		})
	}
}

// TestControllerMiddleware 测试控制器中间件
func TestControllerMiddleware(t *testing.T) {
	container := container.NewContainer()
	resolver := NewControllerResolver(container)
	middleware := NewControllerMiddleware(resolver)
	
	// 创建测试请求
	req, _ := http.NewRequest("GET", "/test", nil)
	request := NewRequest(req)
	
	// 创建下一个处理器
	next := func(req Request) Response {
		return NewJsonResponse(200, map[string]string{
			"message": "Next handler called",
		})
	}
	
	// 执行中间件
	response := middleware.Process(request, next)
	if response.Status() != 200 {
		t.Error("Middleware should pass through to next handler")
	}
}

// BenchmarkControllerResponse 基准测试控制器响应
func BenchmarkControllerResponse(b *testing.B) {
	controller := NewBaseController()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		controller.Success(map[string]string{"key": "value"})
	}
}

// BenchmarkControllerResolver 基准测试控制器解析器
func BenchmarkControllerResolver(b *testing.B) {
	container := container.NewContainer()
	resolver := NewControllerResolver(container)
	
	testController := NewTestController()
	container.Bind("TestController", testController)
	
	req, _ := http.NewRequest("GET", "/test", nil)
	request := NewRequest(req)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler, _ := resolver.Resolve("TestController@Index")
		handler.Handle(request)
	}
} 