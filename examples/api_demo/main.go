package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"laravel-go/framework/api"
)

// User 用户模型
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Password  string    `json:"password,omitempty"`
}

// Post 文章模型
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	User      *User     `json:"user,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Comment 评论模型
type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	fmt.Println("=== Laravel-Go API 系统演示 ===\n")

	// 1. 演示资源转换器
	demoResourceTransformer()

	// 2. 演示版本控制
	demoVersionControl()

	// 3. 演示 API 文档生成
	demoAPIDocumentation()

	// 4. 启动演示服务器
	startDemoServer()
}

// demoResourceTransformer 演示资源转换器
func demoResourceTransformer() {
	fmt.Println("1. 资源转换器演示")
	fmt.Println("==================")

	// 创建测试数据
	user := &User{
		ID:        1,
		Name:      "张三",
		Email:     "zhangsan@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  "secret123",
	}

	post := &Post{
		ID:        1,
		Title:     "我的第一篇文章",
		Content:   "这是文章内容...",
		UserID:    1,
		User:      user,
		Tags:      []string{"技术", "Go语言"},
		CreatedAt: time.Now(),
	}

	// 基础资源转换
	fmt.Println("\n基础资源转换:")
	userResource := api.NewResource(user)
	result := userResource.ToArray()
	printJSON("用户资源", result)

	// 隐藏敏感字段
	fmt.Println("\n隐藏敏感字段:")
	resourceWithoutPassword := userResource.Without("password")
	resultWithoutPassword := resourceWithoutPassword.ToArray()
	printJSON("用户资源（隐藏密码）", resultWithoutPassword)

	// 添加额外字段
	fmt.Println("\n添加额外字段:")
	resourceWithExtra := userResource.Add("extra_field", "额外信息")
	resultWithExtra := resourceWithExtra.ToArray()
	printJSON("用户资源（添加额外字段）", resultWithExtra)

	// 条件字段
	fmt.Println("\n条件字段:")
	isAdmin := true
	resourceWhen := userResource.When(isAdmin, "admin_field")
	if baseResource, ok := resourceWhen.(*api.BaseResource); ok {
		resourceWhen = baseResource.Add("admin_field", "管理员信息")
	}
	resultWhen := resourceWhen.ToArray()
	printJSON("用户资源（条件字段）", resultWhen)

	// 嵌套资源
	fmt.Println("\n嵌套资源:")
	postResource := api.NewResource(post)
	postResult := postResource.ToArray()
	printJSON("文章资源（包含用户）", postResult)

	// 集合转换
	fmt.Println("\n集合转换:")
	users := []*User{
		{
			ID:    1,
			Name:  "张三",
			Email: "zhangsan@example.com",
		},
		{
			ID:    2,
			Name:  "李四",
			Email: "lisi@example.com",
		},
		{
			ID:    3,
			Name:  "王五",
			Email: "wangwu@example.com",
		},
	}

	collection := api.NewResourceCollection(users)
	collectionResult := collection.ToArray()
	printJSON("用户集合", collectionResult)

	// 集合操作
	fmt.Println("\n集合操作:")

	// 分页
	paginated := collection.Paginate(1, 2)
	paginatedResult := paginated.ToArray()
	printJSON("用户集合（分页 1-2）", paginatedResult)

	// 过滤
	filtered := collection.Filter(func(resource api.Resource) bool {
		data := resource.ToArray()
		return data["id"] == float64(1)
	})
	filteredResult := filtered.ToArray()
	printJSON("用户集合（过滤 ID=1）", filteredResult)

	// 映射
	mapped := collection.Map(func(resource api.Resource) api.Resource {
		if baseResource, ok := resource.(*api.BaseResource); ok {
			return baseResource.Add("processed", true)
		}
		return resource
	})
	mappedResult := mapped.ToArray()
	printJSON("用户集合（映射处理）", mappedResult)

	fmt.Println()
}

// demoVersionControl 演示版本控制
func demoVersionControl() {
	fmt.Println("2. 版本控制演示")
	fmt.Println("===============")

	// 创建版本管理器
	vm := api.NewVersionManager()

	// 注册版本
	fmt.Println("\n注册版本:")
	v1 := vm.RegisterVersion("v1", "stable")
	v2 := vm.RegisterVersion("v2", "beta")
	v3 := vm.RegisterVersion("v3", "alpha")

	fmt.Printf("v1: %s (%s)\n", v1.Version, v1.Status)
	fmt.Printf("v2: %s (%s)\n", v2.Version, v2.Status)
	fmt.Printf("v3: %s (%s)\n", v3.Version, v3.Status)

	// 设置默认版本
	vm.SetDefaultVersion("v2")
	fmt.Printf("\n默认版本: %s\n", vm.GetDefaultVersion())

	// 弃用版本
	sunsetTime := time.Now().Add(30 * 24 * time.Hour) // 30天后
	err := vm.DeprecateVersion("v1", "v1 版本将在30天后停止支持", sunsetTime)
	if err != nil {
		fmt.Printf("弃用版本失败: %v\n", err)
	}

	fmt.Printf("\nv1 是否已弃用: %t\n", vm.IsVersionDeprecated("v1"))

	// 获取支持的版本
	supported := vm.GetSupportedVersions()
	fmt.Printf("\n支持的版本数量: %d\n", len(supported))

	// 创建版本路由器
	router := api.NewVersionRouter(vm)

	// 添加不同版本的路由
	router.GET("v1", "/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"version": "v1",
			"message": "这是 v1 版本的用户接口",
			"data":    []map[string]interface{}{},
		}
		json.NewEncoder(w).Encode(response)
	})

	router.GET("v2", "/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"version":  "v2",
			"message":  "这是 v2 版本的用户接口",
			"data":     []map[string]interface{}{},
			"features": []string{"分页", "过滤", "排序"},
		}
		json.NewEncoder(w).Encode(response)
	})

	router.GET("v3", "/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"version":  "v3",
			"message":  "这是 v3 版本的用户接口（实验性）",
			"data":     []map[string]interface{}{},
			"features": []string{"分页", "过滤", "排序", "实时更新"},
		}
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("\n版本路由器已配置完成")
	fmt.Println()
}

// demoAPIDocumentation 演示 API 文档生成
func demoAPIDocumentation() {
	fmt.Println("3. API 文档生成演示")
	fmt.Println("===================")

	// 创建 API 文档
	doc := api.NewAPIDocumentation(
		"Laravel-Go API",
		"1.0.0",
		"基于 Laravel-Go 框架的 RESTful API 文档",
	)

	// 设置基础路径
	doc.SetBasePath("/api")

	// 添加服务器
	doc.AddServer("https://api.example.com", "生产环境")
	doc.AddServer("https://staging-api.example.com", "预发布环境")
	doc.AddServer("http://localhost:8080", "本地开发环境")

	// 添加标签
	doc.AddTag("users", "用户管理相关接口")
	doc.AddTag("posts", "文章管理相关接口")
	doc.AddTag("comments", "评论管理相关接口")
	doc.AddTag("auth", "认证相关接口")

	// 生成模式
	userSchema := doc.GenerateSchemaFromStruct("User", &User{})
	postSchema := doc.GenerateSchemaFromStruct("Post", &Post{})
	commentSchema := doc.GenerateSchemaFromStruct("Comment", &Comment{})

	doc.AddSchema("User", userSchema)
	doc.AddSchema("Post", postSchema)
	doc.AddSchema("Comment", commentSchema)

	// 添加路径和操作
	fmt.Println("\n添加 API 路径:")

	// 用户相关接口
	userListOp := api.NewOperation("获取用户列表", "获取系统中的所有用户")
	userListOp.Tags = []string{"users"}
	userListOp.Parameters = []*api.Parameter{
		api.NewParameter("page", "query", "页码", false),
		api.NewParameter("per_page", "query", "每页数量", false),
		api.NewParameter("search", "query", "搜索关键词", false),
	}
	userListOp.Responses["200"] = api.NewResponse("成功获取用户列表")
	userListOp.Responses["400"] = api.NewResponse("请求参数错误")
	userListOp.Responses["401"] = api.NewResponse("未授权")
	doc.AddPath("/users", "GET", userListOp)

	userCreateOp := api.NewOperation("创建用户", "创建新用户")
	userCreateOp.Tags = []string{"users"}
	userCreateOp.RequestBody = &api.RequestBody{
		Description: "用户信息",
		Required:    true,
		Content: map[string]*api.MediaType{
			"application/json": {
				Schema: userSchema,
			},
		},
	}
	userCreateOp.Responses["201"] = api.NewResponse("用户创建成功")
	userCreateOp.Responses["400"] = api.NewResponse("请求参数错误")
	userCreateOp.Responses["409"] = api.NewResponse("用户已存在")
	doc.AddPath("/users", "POST", userCreateOp)

	userGetOp := api.NewOperation("获取用户详情", "根据用户ID获取用户详细信息")
	userGetOp.Tags = []string{"users"}
	userGetOp.Parameters = []*api.Parameter{
		api.NewParameter("id", "path", "用户ID", true),
	}
	userGetOp.Responses["200"] = api.NewResponse("成功获取用户详情")
	userGetOp.Responses["404"] = api.NewResponse("用户不存在")
	doc.AddPath("/users/{id}", "GET", userGetOp)

	// 文章相关接口
	postListOp := api.NewOperation("获取文章列表", "获取系统中的所有文章")
	postListOp.Tags = []string{"posts"}
	postListOp.Parameters = []*api.Parameter{
		api.NewParameter("page", "query", "页码", false),
		api.NewParameter("per_page", "query", "每页数量", false),
		api.NewParameter("user_id", "query", "作者ID", false),
		api.NewParameter("tag", "query", "标签", false),
	}
	postListOp.Responses["200"] = api.NewResponse("成功获取文章列表")
	postListOp.Responses["400"] = api.NewResponse("请求参数错误")
	doc.AddPath("/posts", "GET", postListOp)

	postCreateOp := api.NewOperation("创建文章", "创建新文章")
	postCreateOp.Tags = []string{"posts"}
	postCreateOp.RequestBody = &api.RequestBody{
		Description: "文章信息",
		Required:    true,
		Content: map[string]*api.MediaType{
			"application/json": {
				Schema: postSchema,
			},
		},
	}
	postCreateOp.Responses["201"] = api.NewResponse("文章创建成功")
	postCreateOp.Responses["400"] = api.NewResponse("请求参数错误")
	postCreateOp.Responses["401"] = api.NewResponse("未授权")
	doc.AddPath("/posts", "POST", postCreateOp)

	// 评论相关接口
	commentListOp := api.NewOperation("获取评论列表", "获取文章的评论列表")
	commentListOp.Tags = []string{"comments"}
	commentListOp.Parameters = []*api.Parameter{
		api.NewParameter("post_id", "path", "文章ID", true),
		api.NewParameter("page", "query", "页码", false),
		api.NewParameter("per_page", "query", "每页数量", false),
	}
	commentListOp.Responses["200"] = api.NewResponse("成功获取评论列表")
	commentListOp.Responses["404"] = api.NewResponse("文章不存在")
	doc.AddPath("/posts/{post_id}/comments", "GET", commentListOp)

	commentCreateOp := api.NewOperation("创建评论", "为文章创建新评论")
	commentCreateOp.Tags = []string{"comments"}
	commentCreateOp.Parameters = []*api.Parameter{
		api.NewParameter("post_id", "path", "文章ID", true),
	}
	commentCreateOp.RequestBody = &api.RequestBody{
		Description: "评论信息",
		Required:    true,
		Content: map[string]*api.MediaType{
			"application/json": {
				Schema: commentSchema,
			},
		},
	}
	commentCreateOp.Responses["201"] = api.NewResponse("评论创建成功")
	commentCreateOp.Responses["400"] = api.NewResponse("请求参数错误")
	commentCreateOp.Responses["401"] = api.NewResponse("未授权")
	commentCreateOp.Responses["404"] = api.NewResponse("文章不存在")
	doc.AddPath("/posts/{post_id}/comments", "POST", commentCreateOp)

	// 生成文档
	fmt.Println("\n生成 API 文档:")

	// JSON 格式
	jsonData, err := doc.ToJSON()
	if err != nil {
		fmt.Printf("生成 JSON 文档失败: %v\n", err)
	} else {
		fmt.Printf("JSON 文档大小: %d 字节\n", len(jsonData))
	}

	// HTML 格式
	html := doc.GenerateHTML()
	fmt.Printf("HTML 文档大小: %d 字节\n", len(html))

	// 显示文档统计信息
	fmt.Printf("\n文档统计:\n")
	fmt.Printf("- 文档生成完成\n")
	fmt.Printf("- JSON 格式: %d 字节\n", len(jsonData))
	fmt.Printf("- HTML 格式: %d 字节\n", len(html))

	fmt.Println()
}

// startDemoServer 启动演示服务器
func startDemoServer() {
	fmt.Println("4. 启动演示服务器")
	fmt.Println("==================")

	// 创建版本管理器
	vm := api.NewVersionManager()
	vm.RegisterVersion("v1", "stable")
	vm.RegisterVersion("v2", "beta")

	// 创建版本路由器
	router := api.NewVersionRouter(vm)

	// 添加路由处理函数
	router.GET("v1", "/users", handleV1Users)
	router.GET("v2", "/users", handleV2Users)
	router.GET("v1", "/posts", handleV1Posts)
	router.GET("v2", "/posts", handleV2Posts)

	// 创建 API 文档处理器
	doc := createAPIDocumentation()

	// 设置路由
	http.HandleFunc("/api/", router.ServeHTTP)
	http.HandleFunc("/api-docs.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonData, _ := doc.ToJSON()
		w.Write(jsonData)
	})
	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(doc.GenerateHTML()))
	})

	fmt.Println("服务器启动在 http://localhost:8080")
	fmt.Println("API 文档: http://localhost:8080/docs")
	fmt.Println("OpenAPI JSON: http://localhost:8080/api-docs.json")
	fmt.Println("\n测试端点:")
	fmt.Println("- GET http://localhost:8080/api/v1/users")
	fmt.Println("- GET http://localhost:8080/api/v2/users")
	fmt.Println("- GET http://localhost:8080/api/v1/posts")
	fmt.Println("- GET http://localhost:8080/api/v2/posts")
	fmt.Println("\n按 Ctrl+C 停止服务器")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 路由处理函数
func handleV1Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	users := []*User{
		{ID: 1, Name: "张三", Email: "zhangsan@example.com"},
		{ID: 2, Name: "李四", Email: "lisi@example.com"},
	}

	collection := api.NewResourceCollection(users)
	jsonData, _ := collection.ToJSON()

	response := map[string]interface{}{
		"version": "v1",
		"message": "用户列表 (v1)",
		"data":    json.RawMessage(jsonData),
	}

	json.NewEncoder(w).Encode(response)
}

func handleV2Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	users := []*User{
		{ID: 1, Name: "张三", Email: "zhangsan@example.com"},
		{ID: 2, Name: "李四", Email: "lisi@example.com"},
		{ID: 3, Name: "王五", Email: "wangwu@example.com"},
	}

	collection := api.NewResourceCollection(users)
	collectionWithFields := collection.With("created_at", "updated_at")
	jsonData, _ := collectionWithFields.ToJSON()

	response := map[string]interface{}{
		"version": "v2",
		"message": "用户列表 (v2) - 包含更多字段",
		"data":    json.RawMessage(jsonData),
		"pagination": map[string]interface{}{
			"current_page": 1,
			"per_page":     10,
			"total":        3,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func handleV1Posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	posts := []*Post{
		{ID: 1, Title: "第一篇文章", Content: "内容...", UserID: 1},
		{ID: 2, Title: "第二篇文章", Content: "内容...", UserID: 2},
	}

	collection := api.NewResourceCollection(posts)
	jsonData, _ := collection.ToJSON()

	response := map[string]interface{}{
		"version": "v1",
		"message": "文章列表 (v1)",
		"data":    json.RawMessage(jsonData),
	}

	json.NewEncoder(w).Encode(response)
}

func handleV2Posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	user := &User{ID: 1, Name: "张三", Email: "zhangsan@example.com"}
	posts := []*Post{
		{ID: 1, Title: "第一篇文章", Content: "内容...", UserID: 1, User: user, Tags: []string{"技术"}},
		{ID: 2, Title: "第二篇文章", Content: "内容...", UserID: 2, Tags: []string{"生活"}},
	}

	collection := api.NewResourceCollection(posts)
	collectionWithFields := collection.With("user", "tags", "created_at")
	jsonData, _ := collectionWithFields.ToJSON()

	response := map[string]interface{}{
		"version": "v2",
		"message": "文章列表 (v2) - 包含用户和标签信息",
		"data":    json.RawMessage(jsonData),
		"meta": map[string]interface{}{
			"total_posts": 2,
			"total_users": 2,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// createAPIDocumentation 创建 API 文档
func createAPIDocumentation() *api.APIDocumentation {
	doc := api.NewAPIDocumentation(
		"Laravel-Go API 演示",
		"1.0.0",
		"API 系统演示文档",
	)

	doc.SetBasePath("/api")
	doc.AddServer("http://localhost:8080", "本地演示服务器")

	// 添加标签
	doc.AddTag("users", "用户管理")
	doc.AddTag("posts", "文章管理")

	// 生成模式
	userSchema := doc.GenerateSchemaFromStruct("User", &User{})
	postSchema := doc.GenerateSchemaFromStruct("Post", &Post{})

	doc.AddSchema("User", userSchema)
	doc.AddSchema("Post", postSchema)

	// 添加路径
	userListOp := api.NewOperation("获取用户列表", "获取所有用户")
	userListOp.Tags = []string{"users"}
	userListOp.Responses["200"] = api.NewResponse("成功")
	doc.AddPath("/v1/users", "GET", userListOp)
	doc.AddPath("/v2/users", "GET", userListOp)

	postListOp := api.NewOperation("获取文章列表", "获取所有文章")
	postListOp.Tags = []string{"posts"}
	postListOp.Responses["200"] = api.NewResponse("成功")
	doc.AddPath("/v1/posts", "GET", postListOp)
	doc.AddPath("/v2/posts", "GET", postListOp)

	return doc
}

// printJSON 打印 JSON 数据
func printJSON(title string, data interface{}) {
	fmt.Printf("\n%s:\n", title)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("JSON 序列化失败: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}
