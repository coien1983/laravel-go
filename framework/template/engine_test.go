package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine("/views", "/cache", true)

	if engine.viewsPath != "/views" {
		t.Errorf("Expected viewsPath /views, got %s", engine.viewsPath)
	}

	if engine.cachePath != "/cache" {
		t.Errorf("Expected cachePath /cache, got %s", engine.cachePath)
	}

	if !engine.useCache {
		t.Error("Expected useCache to be true")
	}

	if engine.componentManager == nil {
		t.Error("Expected componentManager to be initialized")
	}
}

func TestAddHelper(t *testing.T) {
	engine := NewEngine("/views", "/cache", false)

	// 添加助手函数
	engine.AddHelper("uppercase", func(s string) string {
		return strings.ToUpper(s)
	})

	// 验证助手函数已添加
	engine.helpersMutex.RLock()
	_, exists := engine.helpers["uppercase"]
	engine.helpersMutex.RUnlock()

	if !exists {
		t.Error("Helper function was not added")
	}
}

func TestRenderSimpleTemplate(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	engine := NewEngine(tempDir, tempDir, false)

	// 创建测试模板文件
	templateContent := `Hello {{ .name }}!`
	templatePath := filepath.Join(tempDir, "test.blade.php")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// 渲染模板
	data := Data{"name": "World"}
	result, err := engine.Render("test", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	expected := "Hello World!"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestRenderWithControlStructures(t *testing.T) {
	tempDir := t.TempDir()
	engine := NewEngine(tempDir, tempDir, false)

	// 创建包含控制结构的模板
	templateContent := `
@if (show)
	Hello {{ .name }}!
@endif
@foreach (items as item)
	- {{ .item }}
@endforeach
`
	templatePath := filepath.Join(tempDir, "control.blade.php")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// 渲染模板
	data := Data{
		"show":  true,
		"name":  "World",
		"items": []string{"item1", "item2", "item3"},
	}
	result, err := engine.Render("control", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 验证结果包含期望的内容
	if !strings.Contains(result, "Hello World!") {
		t.Error("Expected result to contain 'Hello World!'")
	}

	if !strings.Contains(result, "item1") || !strings.Contains(result, "item2") || !strings.Contains(result, "item3") {
		t.Error("Expected result to contain all items")
	}
}

func TestRenderWithLayout(t *testing.T) {
	tempDir := t.TempDir()
	engine := NewEngine(tempDir, tempDir, false)

	// 注册布局
	layoutContent := `
<!DOCTYPE html>
<html>
<head>
	<title>{{ .title }}</title>
</head>
<body>
	@yield('content')
</body>
</html>
`
	engine.componentManager.RegisterLayout("app", layoutContent)

	// 创建子模板
	childContent := `
@extends('app')

@section('content')
	<h1>{{ .title }}</h1>
	<p>{{ .message }}</p>
@endsection
`
	templatePath := filepath.Join(tempDir, "child.blade.php")
	err := os.WriteFile(templatePath, []byte(childContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// 渲染模板
	data := Data{
		"title":   "My Page",
		"message": "Hello World!",
	}
	result, err := engine.Render("child", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 验证结果
	if !strings.Contains(result, "<title>My Page</title>") {
		t.Error("Expected result to contain title")
	}

	if !strings.Contains(result, "<h1>My Page</h1>") {
		t.Error("Expected result to contain h1")
	}

	if !strings.Contains(result, "<p>Hello World!</p>") {
		t.Error("Expected result to contain message")
	}
}

func TestRenderWithComponent(t *testing.T) {
	tempDir := t.TempDir()
	engine := NewEngine(tempDir, tempDir, false)

	// 注册组件
	componentContent := `
<div class="alert alert-{{ .type }}">
	@slot('title')
		Default Title
	@endslot
	{{ .message }}
</div>
`
	engine.componentManager.RegisterComponent("alert", componentContent)

	// 创建使用组件的模板
	templateContent := `
@component('alert')
	@slot('title')
		Warning!
	@endslot
	This is a warning message.
@endcomponent
`
	templatePath := filepath.Join(tempDir, "component.blade.php")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// 渲染模板
	data := Data{"type": "warning", "message": "This is a warning message."}
	result, err := engine.Render("component", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 验证结果
	if !strings.Contains(result, "alert-warning") {
		t.Error("Expected result to contain alert-warning class")
	}

	if !strings.Contains(result, "Warning!") {
		t.Error("Expected result to contain slot content")
	}

	if !strings.Contains(result, "This is a warning message.") {
		t.Error("Expected result to contain message")
	}
}

func TestRenderWithInclude(t *testing.T) {
	tempDir := t.TempDir()
	engine := NewEngine(tempDir, tempDir, false)

	// 创建被包含的模板
	includeContent := `<p>{{ .message }}</p>`
	includePath := filepath.Join(tempDir, "partial.blade.php")
	err := os.WriteFile(includePath, []byte(includeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write include file: %v", err)
	}

	// 创建主模板
	mainContent := `
<div>
	<h1>{{ .title }}</h1>
	@include('partial')
</div>
`
	mainPath := filepath.Join(tempDir, "main.blade.php")
	err = os.WriteFile(mainPath, []byte(mainContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write main file: %v", err)
	}

	// 渲染模板
	data := Data{
		"title":   "My Page",
		"message": "Hello from partial!",
	}
	result, err := engine.Render("main", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 验证结果
	if !strings.Contains(result, "<h1>My Page</h1>") {
		t.Error("Expected result to contain title")
	}

	if !strings.Contains(result, "<p>Hello from partial!</p>") {
		t.Error("Expected result to contain partial content")
	}
}

func TestRenderWithStack(t *testing.T) {
	tempDir := t.TempDir()
	engine := NewEngine(tempDir, tempDir, false)

	// 创建使用stack的模板
	templateContent := `
<!DOCTYPE html>
<html>
<head>
	@stack('scripts')
</head>
<body>
	@push('scripts')
		<script src="app.js"></script>
	@endpush
	@push('scripts')
		<script src="utils.js"></script>
	@endpush
</body>
</html>
`
	templatePath := filepath.Join(tempDir, "stack.blade.php")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// 渲染模板
	data := Data{}
	result, err := engine.Render("stack", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 验证结果
	if !strings.Contains(result, "app.js") {
		t.Error("Expected result to contain app.js")
	}

	if !strings.Contains(result, "utils.js") {
		t.Error("Expected result to contain utils.js")
	}
}

func TestRenderWithErrorAndOld(t *testing.T) {
	tempDir := t.TempDir()
	engine := NewEngine(tempDir, tempDir, false)

	// 创建包含错误和旧值的模板
	templateContent := `
<form>
	<input type="text" name="name" value="@old('name', 'Default')">
	@error('name')
		<span class="error">{{ .errors.name }}</span>
	@enderror
</form>
`
	templatePath := filepath.Join(tempDir, "form.blade.php")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// 渲染模板（无错误）
	data := Data{
		"old": map[string]string{
			"name": "John",
		},
		"errors": map[string]string{},
	}
	result, err := engine.Render("form", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 验证结果
	if !strings.Contains(result, `value="John"`) {
		t.Error("Expected result to contain old value")
	}

	// 渲染模板（有错误）
	dataWithError := Data{
		"old": map[string]string{},
		"errors": map[string]string{
			"name": "Name is required",
		},
	}
	resultWithError, err := engine.Render("form", dataWithError)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 验证结果
	if !strings.Contains(resultWithError, `value="Default"`) {
		t.Error("Expected result to contain default value")
	}

	if !strings.Contains(resultWithError, "Name is required") {
		t.Error("Expected result to contain error message")
	}
}

func TestCacheFunctionality(t *testing.T) {
	tempDir := t.TempDir()
	engine := NewEngine(tempDir, tempDir, true)

	// 创建模板文件
	templateContent := `Hello {{ .name }}!`
	templatePath := filepath.Join(tempDir, "cache.blade.php")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// 第一次渲染（应该编译并缓存）
	data := Data{"name": "World"}
	result1, err := engine.Render("cache", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 验证缓存
	cached, exists := engine.GetCachedTemplate("cache")
	if !exists {
		t.Error("Template should be cached")
	}

	if cached.Name != "cache" {
		t.Errorf("Expected cached template name 'cache', got %s", cached.Name)
	}

	// 第二次渲染（应该使用缓存）
	result2, err := engine.Render("cache", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	if result1 != result2 {
		t.Error("Cached and non-cached results should be identical")
	}

	// 清除缓存
	engine.ClearCache()

	// 验证缓存已清除
	_, exists = engine.GetCachedTemplate("cache")
	if exists {
		t.Error("Template should not be cached after clearing")
	}
}

func TestRenderToWriter(t *testing.T) {
	tempDir := t.TempDir()
	engine := NewEngine(tempDir, tempDir, false)

	// 创建模板文件
	templateContent := `Hello {{ .name }}!`
	templatePath := filepath.Join(tempDir, "writer.blade.php")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// 渲染到Writer
	data := Data{"name": "World"}
	var buf strings.Builder
	err = engine.RenderToWriter("writer", data, &buf)
	if err != nil {
		t.Fatalf("Failed to render template to writer: %v", err)
	}

	result := buf.String()
	expected := "Hello World!"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
