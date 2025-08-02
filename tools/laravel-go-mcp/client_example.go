package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// MCPClientExample MCP客户端示例
type MCPClientExample struct {
	baseURL string
	client  *http.Client
}

// NewMCPClientExample 创建新的MCP客户端示例
func NewMCPClientExample(baseURL string) *MCPClientExample {
	return &MCPClientExample{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ClientInitializeRequest 初始化请求参数
type ClientInitializeRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Modules     []string `json:"modules"`
	Database    string   `json:"database"`
	Cache       string   `json:"cache"`
	Queue       string   `json:"queue"`
}

// ClientGenerateRequest 生成请求参数
type ClientGenerateRequest struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// ClientDeployRequest 部署请求参数
type ClientDeployRequest struct {
	Environment string `json:"environment"`
}

// Call 调用MCP方法
func (c *MCPClientExample) Call(method string, params interface{}) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      time.Now().Unix(),
		"method":  method,
		"params":  params,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	resp, err := c.client.Post(c.baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return response, nil
}

// Initialize 初始化项目
func (c *MCPClientExample) Initialize(params *ClientInitializeRequest) (map[string]interface{}, error) {
	return c.Call("initialize", params)
}

// Generate 生成模块
func (c *MCPClientExample) Generate(params *ClientGenerateRequest) (map[string]interface{}, error) {
	return c.Call("generate", params)
}

// Build 构建项目
func (c *MCPClientExample) Build() (map[string]interface{}, error) {
	return c.Call("build", nil)
}

// Test 运行测试
func (c *MCPClientExample) Test() (map[string]interface{}, error) {
	return c.Call("test", nil)
}

// Deploy 部署项目
func (c *MCPClientExample) Deploy(params *ClientDeployRequest) (map[string]interface{}, error) {
	return c.Call("deploy", params)
}

// Monitor 获取性能监控数据
func (c *MCPClientExample) Monitor() (map[string]interface{}, error) {
	return c.Call("monitor", nil)
}

// Analyze 代码分析
func (c *MCPClientExample) Analyze() (map[string]interface{}, error) {
	return c.Call("analyze", nil)
}

// Optimize 性能优化
func (c *MCPClientExample) Optimize() (map[string]interface{}, error) {
	return c.Call("optimize", nil)
}

// Info 获取项目信息
func (c *MCPClientExample) Info() (map[string]interface{}, error) {
	return c.Call("info", nil)
}

// RunDemo 演示MCP客户端使用
func RunDemo() {
	client := NewMCPClientExample("http://localhost:8080")

	fmt.Println("🚀 Laravel-Go MCP 客户端演示")
	fmt.Println(strings.Repeat("=", 50))

	// 1. 初始化项目
	fmt.Println("\n1. 初始化项目...")
	initParams := &ClientInitializeRequest{
		Name:        "demo-api",
		Description: "演示API项目",
		Version:     "1.0.0",
		Modules:     []string{"user", "product", "order"},
		Database:    "mysql",
		Cache:       "redis",
		Queue:       "redis",
	}

	resp, err := client.Initialize(initParams)
	if err != nil {
		fmt.Printf("❌ 初始化失败: %v\n", err)
		return
	}

	if resp["error"] != nil {
		fmt.Printf("❌ 初始化错误: %v\n", resp["error"])
		return
	}

	fmt.Printf("✅ 初始化成功: %v\n", resp["result"])

	// 2. 生成新模块
	fmt.Println("\n2. 生成新模块...")
	generateParams := &ClientGenerateRequest{
		Type: "api",
		Name: "category",
	}

	resp, err = client.Generate(generateParams)
	if err != nil {
		fmt.Printf("❌ 生成模块失败: %v\n", err)
		return
	}

	if resp["error"] != nil {
		fmt.Printf("❌ 生成模块错误: %v\n", resp["error"])
		return
	}

	fmt.Printf("✅ 生成模块成功: %v\n", resp["result"])

	// 3. 构建项目
	fmt.Println("\n3. 构建项目...")
	resp, err = client.Build()
	if err != nil {
		fmt.Printf("❌ 构建失败: %v\n", err)
		return
	}

	if resp["error"] != nil {
		fmt.Printf("❌ 构建错误: %v\n", resp["error"])
		return
	}

	fmt.Printf("✅ 构建成功: %v\n", resp["result"])

	// 4. 运行测试
	fmt.Println("\n4. 运行测试...")
	resp, err = client.Test()
	if err != nil {
		fmt.Printf("❌ 测试失败: %v\n", err)
		return
	}

	if resp["error"] != nil {
		fmt.Printf("❌ 测试错误: %v\n", resp["error"])
		return
	}

	fmt.Printf("✅ 测试完成: %v\n", resp["result"])

	// 5. 获取性能监控数据
	fmt.Println("\n5. 获取性能监控数据...")
	resp, err = client.Monitor()
	if err != nil {
		fmt.Printf("❌ 获取监控数据失败: %v\n", err)
		return
	}

	if resp["error"] != nil {
		fmt.Printf("❌ 获取监控数据错误: %v\n", resp["error"])
		return
	}

	fmt.Printf("✅ 监控数据: %v\n", resp["result"])

	// 6. 代码分析
	fmt.Println("\n6. 代码分析...")
	resp, err = client.Analyze()
	if err != nil {
		fmt.Printf("❌ 代码分析失败: %v\n", err)
		return
	}

	if resp["error"] != nil {
		fmt.Printf("❌ 代码分析错误: %v\n", resp["error"])
		return
	}

	fmt.Printf("✅ 代码分析完成: %v\n", resp["result"])

	// 7. 性能优化
	fmt.Println("\n7. 性能优化...")
	resp, err = client.Optimize()
	if err != nil {
		fmt.Printf("❌ 性能优化失败: %v\n", err)
		return
	}

	if resp["error"] != nil {
		fmt.Printf("❌ 性能优化错误: %v\n", resp["error"])
		return
	}

	fmt.Printf("✅ 性能优化完成: %v\n", resp["result"])

	// 8. 获取项目信息
	fmt.Println("\n8. 获取项目信息...")
	resp, err = client.Info()
	if err != nil {
		fmt.Printf("❌ 获取项目信息失败: %v\n", err)
		return
	}

	if resp["error"] != nil {
		fmt.Printf("❌ 获取项目信息错误: %v\n", resp["error"])
		return
	}

	fmt.Printf("✅ 项目信息: %v\n", resp["result"])

	fmt.Println("\n�� MCP 客户端演示完成!")
}
