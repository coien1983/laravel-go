package main

import (
	"fmt"
	"laravel-go/framework/routing"
)

// 示例处理器
func homeHandler() string {
	return "Welcome to Laravel-Go Framework!"
}

func userListHandler() string {
	return "User List"
}

func userDetailHandler(id string) string {
	return fmt.Sprintf("User Detail: %s", id)
}

func createUserHandler() string {
	return "Create User"
}

func updateUserHandler(id string) string {
	return fmt.Sprintf("Update User: %s", id)
}

func deleteUserHandler(id string) string {
	return fmt.Sprintf("Delete User: %s", id)
}

func postListHandler() string {
	return "Post List"
}

func postDetailHandler(id string) string {
	return fmt.Sprintf("Post Detail: %s", id)
}

func postCommentsHandler(id string) string {
	return fmt.Sprintf("Post Comments: %s", id)
}

func adminDashboardHandler() string {
	return "Admin Dashboard"
}

func main() {
	// 创建路由器
	router := routing.NewRouter()

	// 添加全局中间件
	router.Use("auth", "cors")

	// 基础路由
	router.Get("/", homeHandler)

	// API 路由分组
	router.Group("/api/v1", func(r routing.Router) {
		r.Use("api")

		// 用户相关路由
		r.Group("/users", func(r routing.Router) {
			r.Get("", userListHandler)           // GET /api/v1/users
			r.Post("", createUserHandler)        // POST /api/v1/users
			r.Get("/{id}", userDetailHandler)    // GET /api/v1/users/{id}
			r.Put("/{id}", updateUserHandler)    // PUT /api/v1/users/{id}
			r.Delete("/{id}", deleteUserHandler) // DELETE /api/v1/users/{id}
		})

		// 文章相关路由
		r.Group("/posts", func(r routing.Router) {
			r.Get("", postListHandler)                   // GET /api/v1/posts
			r.Post("", createUserHandler)                // POST /api/v1/posts
			r.Get("/{id}", postDetailHandler)            // GET /api/v1/posts/{id}
			r.Get("/{id}/comments", postCommentsHandler) // GET /api/v1/posts/{id}/comments
			r.Post("/{id}/comments", createUserHandler)  // POST /api/v1/posts/{id}/comments
		})
	})

	// 管理后台路由分组
	router.Group("/admin", func(r routing.Router) {
		r.Use("admin")
		r.Get("/dashboard", adminDashboardHandler) // GET /admin/dashboard
		r.Get("/users", userListHandler)           // GET /admin/users
	})

	// 添加路由约束
	router.Where("id", "[0-9]+")

	// 创建路由命令
	command := routing.NewRouteCommand(router)

	// 显示所有路由
	fmt.Println("=== 所有路由 ===")
	fmt.Println(command.List())

	// 显示GET方法的路由
	fmt.Println("\n=== GET方法路由 ===")
	fmt.Println(command.ListByMethod("GET"))

	// 显示API分组的路由
	fmt.Println("\n=== API分组路由 ===")
	fmt.Println(command.ListByGroup("/api/v1"))

	// 显示特定路由详情
	fmt.Println("\n=== 路由详情 ===")
	fmt.Println(command.Show("/api/v1/users/{id}"))

	// 测试路由匹配
	fmt.Println("\n=== 路由匹配测试 ===")

	testCases := []struct {
		method string
		path   string
	}{
		{"GET", "/"},
		{"GET", "/api/v1/users"},
		{"GET", "/api/v1/users/123"},
		{"POST", "/api/v1/users"},
		{"PUT", "/api/v1/users/456"},
		{"DELETE", "/api/v1/users/789"},
		{"GET", "/api/v1/posts/123/comments"},
		{"GET", "/admin/dashboard"},
		{"GET", "/invalid/path"},
	}

	for _, tc := range testCases {
		route, found := router.Match(tc.method, tc.path)
		if found {
			fmt.Printf("✓ %s %s -> Handler: %T, Parameters: %v\n",
				tc.method, tc.path, route.Handler, route.Parameters)
		} else {
			fmt.Printf("✗ %s %s -> Not found\n", tc.method, tc.path)
		}
	}
}
