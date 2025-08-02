package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestUser 测试用户结构体
type TestUser struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Password  string    `json:"password,omitempty"`
}

// TestPost 测试文章结构体
type TestPost struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	UserID  int       `json:"user_id"`
	User    *TestUser `json:"user,omitempty"`
	Tags    []string  `json:"tags,omitempty"`
}

// TestResource 测试资源转换器
func TestResource(t *testing.T) {
	user := &TestUser{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		Password:  "secret",
	}

	// 测试基础资源转换
	resource := NewResource(user)
	result := resource.ToArray()

	// 验证基本字段
	if result["id"] != int64(1) {
		t.Errorf("Expected id to be 1, got %v", result["id"])
	}
	if result["name"] != "John Doe" {
		t.Errorf("Expected name to be 'John Doe', got %v", result["name"])
	}
	if result["email"] != "john@example.com" {
		t.Errorf("Expected email to be 'john@example.com', got %v", result["email"])
	}

	// 验证时间格式
	if result["created_at"] != "2023-01-01T00:00:00Z" {
		t.Errorf("Expected created_at to be formatted time, got %v", result["created_at"])
	}

	// 测试 Without 方法
	resourceWithoutPassword := resource.Without("password")
	resultWithoutPassword := resourceWithoutPassword.ToArray()

	if _, exists := resultWithoutPassword["password"]; exists {
		t.Error("Password should be hidden")
	}

	// 测试 With 方法
	resourceWithExtra := resource.With("extra_field")
	if baseResource, ok := resourceWithExtra.(*BaseResource); ok {
		resourceWithExtra = baseResource.Add("extra_field", "extra_value")
		resultWithExtra := resourceWithExtra.ToArray()

		if resultWithExtra["extra_field"] != "extra_value" {
			t.Errorf("Expected extra_field to be 'extra_value', got %v", resultWithExtra["extra_field"])
		}
	}

	// 测试 When 方法
	resourceWhen := resource.When(true, "conditional_field")
	if baseResource, ok := resourceWhen.(*BaseResource); ok {
		resourceWhen = baseResource.Add("conditional_field", "conditional_value")
		resultWhen := resourceWhen.ToArray()

		if resultWhen["conditional_field"] != "conditional_value" {
			t.Errorf("Expected conditional_field to be 'conditional_value', got %v", resultWhen["conditional_field"])
		}
	}

	// 测试 JSON 转换
	jsonData, err := resource.ToJSON()
	if err != nil {
		t.Errorf("Failed to convert to JSON: %v", err)
	}

	var jsonResult map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonResult); err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if jsonResult["id"] != float64(1) {
		t.Errorf("Expected JSON id to be 1, got %v", jsonResult["id"])
	}
}

// TestCollection 测试集合转换器
func TestCollection(t *testing.T) {
	users := []*TestUser{
		{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		},
		{
			ID:    2,
			Name:  "Jane Smith",
			Email: "jane@example.com",
		},
	}

	// 创建资源集合
	resources := make([]Resource, len(users))
	for i, user := range users {
		resources[i] = NewResource(user)
	}

	collection := NewCollection(resources)
	result := collection.ToArray()

	if len(result) != 2 {
		t.Errorf("Expected 2 items in collection, got %d", len(result))
	}

	// 测试集合级别的字段过滤
	collectionWithoutEmail := collection.Without("email")
	resultWithoutEmail := collectionWithoutEmail.ToArray()

	for _, item := range resultWithoutEmail {
		if _, exists := item["email"]; exists {
			t.Error("Email should be hidden in all collection items")
		}
	}

	// 测试分页
	paginated := collection.Paginate(1, 1)
	paginatedResult := paginated.ToArray()

	if len(paginatedResult) != 1 {
		t.Errorf("Expected 1 item after pagination, got %d", len(paginatedResult))
	}

	// 测试过滤
	filtered := collection.Filter(func(resource Resource) bool {
		data := resource.ToArray()
		return data["id"] == int64(1)
	})
	filteredResult := filtered.ToArray()

	if len(filteredResult) != 1 {
		t.Errorf("Expected 1 item after filtering, got %d", len(filteredResult))
		return
	}

	if filteredResult[0]["id"] != int64(1) {
		t.Errorf("Expected filtered item to have id 1, got %v", filteredResult[0]["id"])
	}

	// 测试映射
	mapped := collection.Map(func(resource Resource) Resource {
		if baseResource, ok := resource.(*BaseResource); ok {
			return baseResource.Add("mapped", true)
		}
		return resource
	})
	mappedResult := mapped.ToArray()

	for _, item := range mappedResult {
		if item["mapped"] != true {
			t.Error("Expected mapped field to be true")
		}
	}
}

// TestResourceCollection 测试资源集合
func TestResourceCollection(t *testing.T) {
	users := []*TestUser{
		{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		},
		{
			ID:    2,
			Name:  "Jane Smith",
			Email: "jane@example.com",
		},
	}

	collection := NewResourceCollection(users)
	result := collection.ToArray()

	if len(result) != 2 {
		t.Errorf("Expected 2 items in collection, got %d", len(result))
	}

	// 测试 JSON 转换
	jsonData, err := collection.ToJSON()
	if err != nil {
		t.Errorf("Failed to convert collection to JSON: %v", err)
	}

	var jsonResult []map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonResult); err != nil {
		t.Errorf("Failed to unmarshal collection JSON: %v", err)
	}

	if len(jsonResult) != 2 {
		t.Errorf("Expected 2 items in JSON result, got %d", len(jsonResult))
	}
}

// TestVersionManager 测试版本管理器
func TestVersionManager(t *testing.T) {
	vm := NewVersionManager()

	// 测试注册版本
	v1 := vm.RegisterVersion("v1", "stable")
	if v1.Version != "v1" {
		t.Errorf("Expected version to be 'v1', got %s", v1.Version)
	}
	if v1.Status != "stable" {
		t.Errorf("Expected status to be 'stable', got %s", v1.Status)
	}

	v2 := vm.RegisterVersion("v2", "beta")
	if v2.Version != "v2" {
		t.Errorf("Expected version to be 'v2', got %s", v2.Version)
	}

	// 测试获取版本
	if version, exists := vm.GetVersion("v1"); !exists {
		t.Error("Version v1 should exist")
	} else if version.Version != "v1" {
		t.Errorf("Expected version to be 'v1', got %s", version.Version)
	}

	// 测试默认版本
	if vm.GetDefaultVersion() != "v1" {
		t.Errorf("Expected default version to be 'v1', got %s", vm.GetDefaultVersion())
	}

	vm.SetDefaultVersion("v2")
	if vm.GetDefaultVersion() != "v2" {
		t.Errorf("Expected default version to be 'v2', got %s", vm.GetDefaultVersion())
	}

	// 测试版本弃用
	sunsetTime := time.Now().Add(24 * time.Hour)
	err := vm.DeprecateVersion("v1", "v1 is deprecated", sunsetTime)
	if err != nil {
		t.Errorf("Failed to deprecate version: %v", err)
	}

	if !vm.IsVersionDeprecated("v1") {
		t.Error("Version v1 should be deprecated")
	}

	// 测试获取支持的版本
	supported := vm.GetSupportedVersions()
	if len(supported) != 2 {
		t.Errorf("Expected 2 supported versions, got %d", len(supported))
	}
}

// TestVersionMiddleware 测试版本中间件
func TestVersionMiddleware(t *testing.T) {
	vm := NewVersionManager()
	vm.RegisterVersion("v1", "stable")
	vm.RegisterVersion("v2", "beta")

	middleware := NewVersionMiddleware(vm)

	// 测试从头部获取版本
	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Accept-Version", "v2")
	w := httptest.NewRecorder()

	var capturedVersion string
	handler := middleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		capturedVersion = VersionFromContext(r.Context())
	})

	handler(w, req)

	if capturedVersion != "v2" {
		t.Errorf("Expected version to be 'v2', got %s", capturedVersion)
	}

	// 测试从查询参数获取版本
	req = httptest.NewRequest("GET", "/api/users?version=v1", nil)
	w = httptest.NewRecorder()

	handler(w, req)

	if capturedVersion != "v1" {
		t.Errorf("Expected version to be 'v1', got %s", capturedVersion)
	}

	// 测试从路径获取版本
	req = httptest.NewRequest("GET", "/api/v2/users", nil)
	w = httptest.NewRecorder()

	handler(w, req)

	if capturedVersion != "v2" {
		t.Errorf("Expected version to be 'v2', got %s", capturedVersion)
	}

	// 测试不支持的版本
	req = httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Accept-Version", "v3")
	w = httptest.NewRecorder()

	handler(w, req)

	if capturedVersion != "v1" { // 应该使用默认版本
		t.Errorf("Expected version to be 'v1' (default), got %s", capturedVersion)
	}
}

// TestVersionRouter 测试版本路由器
func TestVersionRouter(t *testing.T) {
	vm := NewVersionManager()
	vm.RegisterVersion("v1", "stable")
	vm.RegisterVersion("v2", "beta")

	router := NewVersionRouter(vm)

	// 添加路由
	var capturedVersion string
	router.GET("v1", "/users", func(w http.ResponseWriter, r *http.Request) {
		capturedVersion = VersionFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	router.GET("v2", "/users", func(w http.ResponseWriter, r *http.Request) {
		capturedVersion = VersionFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	// 测试 v1 路由
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if capturedVersion != "v1" {
		t.Errorf("Expected version to be 'v1', got %s", capturedVersion)
	}

	// 测试 v2 路由
	req = httptest.NewRequest("GET", "/api/v2/users", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if capturedVersion != "v2" {
		t.Errorf("Expected version to be 'v2', got %s", capturedVersion)
	}

	// 测试不存在的路由
	req = httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// TestAPIDocumentation 测试 API 文档生成
func TestAPIDocumentation(t *testing.T) {
	doc := NewAPIDocumentation("Test API", "1.0.0", "Test API Documentation")

	// 添加服务器
	doc.AddServer("https://api.example.com", "Production server")
	doc.AddServer("https://staging-api.example.com", "Staging server")

	// 添加标签
	doc.AddTag("users", "User management operations")
	doc.AddTag("posts", "Post management operations")

	// 生成用户模式
	userSchema := doc.GenerateSchemaFromStruct("User", &TestUser{})
	doc.AddSchema("User", userSchema)

	// 生成文章模式
	postSchema := doc.GenerateSchemaFromStruct("Post", &TestPost{})
	doc.AddSchema("Post", postSchema)

	// 添加路径
	operation := NewOperation("Get users", "Retrieve a list of users")
	operation.Responses["200"] = NewResponse("Successful response")
	operation.Responses["400"] = NewResponse("Bad request")

	doc.AddPath("/users", "GET", operation)

	// 测试 JSON 生成
	jsonData, err := doc.ToJSON()
	if err != nil {
		t.Errorf("Failed to generate JSON documentation: %v", err)
	}

	var spec OpenAPISpec
	if err := json.Unmarshal(jsonData, &spec); err != nil {
		t.Errorf("Failed to unmarshal JSON documentation: %v", err)
	}

	// 验证基本信息
	if spec.OpenAPI != "3.0.0" {
		t.Errorf("Expected OpenAPI version to be '3.0.0', got %s", spec.OpenAPI)
	}
	if spec.Info.Title != "Test API" {
		t.Errorf("Expected title to be 'Test API', got %s", spec.Info.Title)
	}
	if spec.Info.Version != "1.0.0" {
		t.Errorf("Expected version to be '1.0.0', got %s", spec.Info.Version)
	}

	// 验证服务器
	if len(spec.Servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(spec.Servers))
	}

	// 验证标签
	if len(spec.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(spec.Tags))
	}

	// 验证模式
	if len(spec.Components.Schemas) != 2 {
		t.Errorf("Expected 2 schemas, got %d", len(spec.Components.Schemas))
	}

	// 验证路径
	if len(spec.Paths) != 1 {
		t.Errorf("Expected 1 path, got %d", len(spec.Paths))
	}

	// 测试 HTML 生成
	html := doc.GenerateHTML()
	if html == "" {
		t.Error("Generated HTML should not be empty")
	}

	if !strings.Contains(html, "Test API") {
		t.Error("Generated HTML should contain API title")
	}
}

// TestSchemaGeneration 测试模式生成
func TestSchemaGeneration(t *testing.T) {
	doc := NewAPIDocumentation("Test API", "1.0.0", "Test API Documentation")

	// 测试基本类型
	schema := doc.GenerateSchemaFromStruct("TestUser", &TestUser{})
	
	if schema.Type != "object" {
		t.Errorf("Expected schema type to be 'object', got %s", schema.Type)
	}

	// 验证必需字段
	requiredFields := []string{"id", "name", "email", "created_at", "updated_at"}
	for _, field := range requiredFields {
		found := false
		for _, required := range schema.Required {
			if required == field {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Field %s should be required", field)
		}
	}

	// 验证属性
	if schema.Properties["id"].Type != "integer" {
		t.Errorf("Expected id type to be 'integer', got %s", schema.Properties["id"].Type)
	}
	if schema.Properties["name"].Type != "string" {
		t.Errorf("Expected name type to be 'string', got %s", schema.Properties["name"].Type)
	}
	if schema.Properties["email"].Type != "string" {
		t.Errorf("Expected email type to be 'string', got %s", schema.Properties["email"].Type)
	}
	if schema.Properties["created_at"].Type != "string" {
		t.Errorf("Expected created_at type to be 'string', got %s", schema.Properties["created_at"].Type)
	}
	if schema.Properties["created_at"].Format != "date-time" {
		t.Errorf("Expected created_at format to be 'date-time', got %s", schema.Properties["created_at"].Format)
	}

	// 验证 omitempty 字段
	if _, exists := schema.Properties["password"]; exists {
		t.Error("Password field should not be in schema due to omitempty")
	}
}

// TestNestedResource 测试嵌套资源
func TestNestedResource(t *testing.T) {
	user := &TestUser{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	post := &TestPost{
		ID:      1,
		Title:   "Test Post",
		Content: "Test content",
		UserID:  1,
		User:    user,
		Tags:    []string{"test", "example"},
	}

	resource := NewResource(post)
	result := resource.ToArray()

	// 验证嵌套用户
	if userData, exists := result["user"]; !exists {
		t.Error("User field should exist")
	} else {
		userMap := userData.(map[string]interface{})
		if userMap["id"] != int64(1) {
			t.Errorf("Expected user id to be 1, got %v", userMap["id"])
		}
		if userMap["name"] != "John Doe" {
			t.Errorf("Expected user name to be 'John Doe', got %v", userMap["name"])
		}
	}

	// 验证切片
	if tags, exists := result["tags"]; !exists {
		t.Error("Tags field should exist")
	} else {
		tagsSlice := tags.([]interface{})
		if len(tagsSlice) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(tagsSlice))
		}
		if tagsSlice[0] != "test" {
			t.Errorf("Expected first tag to be 'test', got %v", tagsSlice[0])
		}
	}
}

// TestConvenienceFunctions 测试便捷函数
func TestConvenienceFunctions(t *testing.T) {
	user := &TestUser{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// 测试 NewResourceFromData
	resource := NewResourceFromData(user)
	if resource == nil {
		t.Error("NewResourceFromData should not return nil")
	}

	// 测试 NewCollectionFromData
	users := []*TestUser{user}
	collection := NewCollectionFromData(users)
	if collection == nil {
		t.Error("NewCollectionFromData should not return nil")
	}

	// 测试 NewResourceFromSlice
	sliceCollection := NewResourceFromSlice(users)
	if sliceCollection == nil {
		t.Error("NewResourceFromSlice should not return nil")
	}

	// 测试 NewAPIVersion
	version := NewAPIVersion("v1", "stable")
	if version.Version != "v1" {
		t.Errorf("Expected version to be 'v1', got %s", version.Version)
	}
	if version.Status != "stable" {
		t.Errorf("Expected status to be 'stable', got %s", version.Status)
	}

	// 测试 NewVersionedAPI
	vm, vr := NewVersionedAPI()
	if vm == nil {
		t.Error("VersionManager should not be nil")
	}
	if vr == nil {
		t.Error("VersionRouter should not be nil")
	}
} 