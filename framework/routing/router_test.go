package routing

import (
	"fmt"
	"testing"
)

// 测试处理器
func testHandler() string {
	return "test response"
}

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Error("NewRouter() should not return nil")
	}
}

func TestRouterBasicMethods(t *testing.T) {
	router := NewRouter()

	// 测试所有HTTP方法
	router.Get("/test", testHandler)
	router.Post("/test", testHandler)
	router.Put("/test", testHandler)
	router.Delete("/test", testHandler)
	router.Patch("/test", testHandler)
	router.Options("/test", testHandler)
	router.Head("/test", testHandler)

	routes := router.GetRoutes()
	if len(routes) != 7 {
		t.Errorf("Expected 7 routes, got %d", len(routes))
	}

	// 验证方法
	expectedMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	for i, method := range expectedMethods {
		if routes[i].Method != method {
			t.Errorf("Expected method %s, got %s", method, routes[i].Method)
		}
	}
}

func TestRouterGroup(t *testing.T) {
	router := NewRouter()

	router.Group("/api", func(r Router) {
		r.Get("/users", testHandler)
		r.Post("/users", testHandler)
	})

	router.Group("/admin", func(r Router) {
		r.Get("/dashboard", testHandler)
	})

	routes := router.GetRoutes()
	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}

	// 验证分组路径
	expectedPaths := []string{"/api/users", "/api/users", "/admin/dashboard"}
	for i, path := range expectedPaths {
		if routes[i].Path != path {
			t.Errorf("Expected path %s, got %s", path, routes[i].Path)
		}
	}
}

func TestRouterParameters(t *testing.T) {
	router := NewRouter()

	router.Get("/users/{id}", testHandler)
	router.Get("/posts/{id}/comments/{commentId}", testHandler)

	// 测试参数匹配
	route, found := router.Match("GET", "/users/123")
	if !found {
		t.Error("Route should be found")
	}
	if route.Parameters["id"] != "123" {
		t.Errorf("Expected parameter id=123, got %s", route.Parameters["id"])
	}

	route, found = router.Match("GET", "/posts/456/comments/789")
	if !found {
		t.Error("Route should be found")
	}
	if route.Parameters["id"] != "456" {
		t.Errorf("Expected parameter id=456, got %s", route.Parameters["id"])
	}
	if route.Parameters["commentId"] != "789" {
		t.Errorf("Expected parameter commentId=789, got %s", route.Parameters["commentId"])
	}
}

func TestRouterMiddleware(t *testing.T) {
	router := NewRouter()

	// 添加全局中间件
	router.Use("auth", "cors")

	router.Get("/test", testHandler)

	routes := router.GetRoutes()
	if len(routes) != 1 {
		t.Error("Expected 1 route")
	}

	if len(routes[0].Middleware) != 2 {
		t.Errorf("Expected 2 middleware, got %d", len(routes[0].Middleware))
	}
}

func TestRouterConstraints(t *testing.T) {
	router := NewRouter()

	router.Where("id", "[0-9]+")
	router.Get("/users/{id}", testHandler)

	// 测试约束匹配
	_, found := router.Match("GET", "/users/123")
	if !found {
		t.Error("Route should be found")
	}

	_, found = router.Match("GET", "/users/abc")
	if found {
		t.Error("Route should not be found with invalid id")
	}
}

func TestRouterCache(t *testing.T) {
	router := NewRouter()

	router.Get("/test", testHandler).Cache(300)

	routes := router.GetRoutes()
	if len(routes) != 1 {
		t.Error("Expected 1 route")
	}

	if routes[0].CacheTTL != 300 {
		t.Errorf("Expected cache TTL 300, got %d", routes[0].CacheTTL)
	}
}

func TestRadixTree(t *testing.T) {
	tree := NewRadixTree()

	// 插入路由
	tree.Insert("GET", "/users", testHandler)
	tree.Insert("GET", "/users/{id}", testHandler)
	tree.Insert("POST", "/users", testHandler)

	// 测试匹配
	handler, params, found := tree.Match("GET", "/users")
	if !found {
		t.Error("Route should be found")
	}
	if handler == nil {
		t.Error("Handler should not be nil")
	}

	handler, params, found = tree.Match("GET", "/users/123")
	if !found {
		t.Error("Route should be found")
	}
	if params["id"] != "123" {
		t.Errorf("Expected parameter id=123, got %s", params["id"])
	}

	// 测试不匹配
	_, _, found = tree.Match("GET", "/invalid")
	if found {
		t.Error("Route should not be found")
	}
}

func TestRouteCommand(t *testing.T) {
	router := NewRouter()

	router.Get("/users", testHandler)
	router.Post("/users", testHandler)
	router.Get("/posts", testHandler)

	command := NewRouteCommand(router)

	// 测试列出所有路由
	list := command.List()
	if list == "" {
		t.Error("List should not be empty")
	}

	// 测试按方法列出路由
	getList := command.ListByMethod("GET")
	if getList == "" {
		t.Error("ListByMethod should not be empty")
	}

	// 测试显示路由详情
	show := command.Show("/users")
	if show == "" {
		t.Error("Show should not be empty")
	}
}

func TestRouterComplexScenario(t *testing.T) {
	router := NewRouter()

	// 添加全局中间件
	router.Use("auth")

	// 创建API分组
	router.Group("/api/v1", func(r Router) {
		r.Use("api")

		// 用户相关路由
		r.Group("/users", func(r Router) {
			r.Get("", testHandler)         // GET /api/v1/users
			r.Post("", testHandler)        // POST /api/v1/users
			r.Get("/{id}", testHandler)    // GET /api/v1/users/{id}
			r.Put("/{id}", testHandler)    // PUT /api/v1/users/{id}
			r.Delete("/{id}", testHandler) // DELETE /api/v1/users/{id}
		})

		// 文章相关路由
		r.Group("/posts", func(r Router) {
			r.Get("", testHandler)                // GET /api/v1/posts
			r.Post("", testHandler)               // POST /api/v1/posts
			r.Get("/{id}", testHandler)           // GET /api/v1/posts/{id}
			r.Get("/{id}/comments", testHandler)  // GET /api/v1/posts/{id}/comments
			r.Post("/{id}/comments", testHandler) // POST /api/v1/posts/{id}/comments
		})
	})

	// 创建管理后台分组
	router.Group("/admin", func(r Router) {
		r.Use("admin")
		r.Get("/dashboard", testHandler) // GET /admin/dashboard
		r.Get("/users", testHandler)     // GET /admin/users
	})

	routes := router.GetRoutes()
	if len(routes) != 12 {
		t.Errorf("Expected 12 routes, got %d", len(routes))
	}

	// 测试API路由匹配
	route, found := router.Match("GET", "/api/v1/users")
	if !found {
		t.Error("API users route should be found")
	}

	route, found = router.Match("GET", "/api/v1/users/123")
	if !found {
		t.Error("API user detail route should be found")
	}
	if route.Parameters["id"] != "123" {
		t.Errorf("Expected parameter id=123, got %s", route.Parameters["id"])
	}

	// 测试管理后台路由匹配
	route, found = router.Match("GET", "/admin/dashboard")
	if !found {
		t.Error("Admin dashboard route should be found")
	}
}

func BenchmarkRouterMatch(b *testing.B) {
	router := NewRouter()

	// 添加大量路由
	for i := 0; i < 1000; i++ {
		router.Get(fmt.Sprintf("/users/%d", i), testHandler)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.Match("GET", "/users/500")
	}
}

func BenchmarkRadixTreeMatch(b *testing.B) {
	tree := NewRadixTree()

	// 添加大量路由
	for i := 0; i < 1000; i++ {
		tree.Insert("GET", fmt.Sprintf("/users/%d", i), testHandler)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Match("GET", "/users/500")
	}
}
