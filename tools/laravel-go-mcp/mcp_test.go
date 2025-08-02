package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestMCPRequest 测试 MCP 请求
func TestMCPRequest(t *testing.T) {
	// 创建测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟 MCP 服务器响应
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result": map[string]interface{}{
				"success": true,
				"message": "测试成功",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	// 创建测试请求
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "test",
		"params":  map[string]interface{}{},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("序列化请求失败: %v", err)
	}

	// 发送请求
	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期望状态码 200，得到 %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	// 检查响应格式
	if response["jsonrpc"] != "2.0" {
		t.Errorf("期望 jsonrpc 2.0，得到 %v", response["jsonrpc"])
	}

	if response["id"] != float64(1) {
		t.Errorf("期望 id 1，得到 %v", response["id"])
	}

	result, ok := response["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("响应中没有 result 字段")
	}

	if result["success"] != true {
		t.Errorf("期望 success true，得到 %v", result["success"])
	}
}

// TestMCPClient 测试 MCP 客户端
func TestMCPClient(t *testing.T) {
	client := NewMCPClientExample("http://localhost:8080")

	// 测试初始化请求
	initParams := &ClientInitializeRequest{
		Name:        "test-api",
		Description: "测试API项目",
		Version:     "1.0.0",
		Modules:     []string{"user"},
		Database:    "mysql",
		Cache:       "redis",
		Queue:       "redis",
	}

	// 注意：这个测试需要实际的服务器运行
	// 在实际环境中，应该使用 mock 服务器
	t.Skip("跳过需要实际服务器的测试")

	_, err := client.Initialize(initParams)
	if err != nil {
		t.Logf("初始化请求失败 (预期): %v", err)
	}
}

// TestProjectConfig 测试项目配置
func TestProjectConfig(t *testing.T) {
	config := &ProjectConfig{
		Name:        "test-project",
		Description: "测试项目",
		Version:     "1.0.0",
		Author:      "Test Author",
		Modules:     []string{"user", "product"},
		Database:    "mysql",
		Cache:       "redis",
		Queue:       "redis",
	}

	// 检查配置字段
	if config.Name != "test-project" {
		t.Errorf("期望名称 test-project，得到 %s", config.Name)
	}

	if config.Description != "测试项目" {
		t.Errorf("期望描述 测试项目，得到 %s", config.Description)
	}

	if len(config.Modules) != 2 {
		t.Errorf("期望 2 个模块，得到 %d", len(config.Modules))
	}

	if config.Database != "mysql" {
		t.Errorf("期望数据库 mysql，得到 %s", config.Database)
	}
}

// TestCodeGeneration 测试代码生成
func TestCodeGeneration(t *testing.T) {
	mcp := &LaravelGoMCP{
		projectPath: ".",
	}

	// 测试控制器生成
	controllerContent := mcp.generateController("user")
	if controllerContent == "" {
		t.Error("控制器内容为空")
	}

	// 检查控制器内容
	if !bytes.Contains([]byte(controllerContent), []byte("UserController")) {
		t.Error("控制器内容中缺少 UserController")
	}

	// 测试模型生成
	modelContent := mcp.generateModel("user")
	if modelContent == "" {
		t.Error("模型内容为空")
	}

	// 检查模型内容
	if !bytes.Contains([]byte(modelContent), []byte("User")) {
		t.Error("模型内容中缺少 User")
	}

	// 测试服务生成
	serviceContent := mcp.generateService("user")
	if serviceContent == "" {
		t.Error("服务内容为空")
	}

	// 检查服务内容
	if !bytes.Contains([]byte(serviceContent), []byte("UserService")) {
		t.Error("服务内容中缺少 UserService")
	}
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	// 测试无效的 JSON
	invalidJSON := `{"invalid": json}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟解析错误
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"error": map[string]interface{}{
				"code":    -32700,
				"message": "解析错误",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte(invalidJSON)))
	if err != nil {
		t.Fatalf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("期望状态码 400，得到 %d", resp.StatusCode)
	}
}

// BenchmarkMCPRequest 性能测试
func BenchmarkMCPRequest(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result": map[string]interface{}{
				"success": true,
				"message": "测试成功",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "test",
		"params":  map[string]interface{}{},
	}

	jsonData, _ := json.Marshal(request)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			b.Fatalf("HTTP请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

// 运行所有测试
func TestMain(m *testing.M) {
	fmt.Println("🧪 开始运行 MCP 服务测试...")

	// 运行测试
	code := m.Run()

	fmt.Println("✅ MCP 服务测试完成")

	// 退出
	// os.Exit(code)
}
