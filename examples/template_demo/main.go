package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"laravel-go/framework/template"
)

func main() {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "template_demo")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	// 创建模板引擎
	engine := template.NewEngine(tempDir, tempDir, false)

	// 添加助手函数
	engine.AddHelper("uppercase", func(s string) string {
		return strings.ToUpper(s)
	})

	// 注册布局
	layoutContent := `
<!DOCTYPE html>
<html>
<head>
	<title>{{ .title }}</title>
</head>
<body>
	<header>
		<h1>{{ .title }}</h1>
	</header>
	<main>
		@yield('content')
	</main>
	<footer>
		<p>&copy; 2024 Laravel-Go Framework</p>
	</footer>
</body>
</html>
`
	engine.GetComponentManager().RegisterLayout("app", layoutContent)

	// 注册组件
	alertComponent := `
<div class="alert alert-{{ .type }}">
	<h4>{{ .title }}</h4>
	<p>{{ .message }}</p>
</div>
`
	engine.GetComponentManager().RegisterComponent("alert", alertComponent)

	// 创建主页面模板
	mainTemplate := `
@extends('app')

@section('content')
	<div class="container">
		<h2>Welcome to Laravel-Go Framework</h2>
		
		@if (showAlert)
			@component('alert')
				@slot('title')
					Welcome!
				@endslot
				This is a welcome message from the template engine.
			@endcomponent
		@endif
		
		<div class="features">
			<h3>Features:</h3>
			@foreach (features as feature)
				<li>{{ .feature }}</li>
			@endforeach
		</div>
		
		<p>Current user: {{ .user }}</p>
		<p>Uppercase message: {{ uppercase .message }}</p>
	</div>
@endsection
`

	// 写入模板文件
	templatePath := filepath.Join(tempDir, "main.blade.php")
	err = os.WriteFile(templatePath, []byte(mainTemplate), 0644)
	if err != nil {
		panic(err)
	}

	// 准备数据
	data := template.Data{
		"title":     "Laravel-Go Template Demo",
		"showAlert": true,
		"user":      "John Doe",
		"message":   "Hello from template engine!",
		"features": []map[string]string{
			{"feature": "Template Inheritance"},
			{"feature": "Component System"},
			{"feature": "Control Structures"},
			{"feature": "Helper Functions"},
		},
	}

	// 渲染模板
	result, err := engine.Render("main", data)
	if err != nil {
		panic(err)
	}

	// 输出结果
	fmt.Println("=== Template Engine Demo ===")
	fmt.Println(result)
	fmt.Println("=== End Demo ===")
}
