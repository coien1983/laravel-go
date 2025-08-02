package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Engine 模板引擎
type Engine struct {
	viewsPath        string
	cachePath        string
	useCache         bool
	cache            map[string]*Template
	cacheMutex       sync.RWMutex
	helpers          map[string]interface{}
	helpersMutex     sync.RWMutex
	layouts          map[string]*Template
	layoutsMutex     sync.RWMutex
	componentManager *ComponentManager
}

// Template 编译后的模板
type Template struct {
	Name     string
	Content  string
	Compiled *template.Template
	Blocks   map[string]string
	Extends  string
}

// Data 模板数据
type Data map[string]interface{}

// NewEngine 创建新的模板引擎
func NewEngine(viewsPath, cachePath string, useCache bool) *Engine {
	return &Engine{
		viewsPath:        viewsPath,
		cachePath:        cachePath,
		useCache:         useCache,
		cache:            make(map[string]*Template),
		helpers:          make(map[string]interface{}),
		layouts:          make(map[string]*Template),
		componentManager: NewComponentManager(),
	}
}

// AddHelper 添加助手函数
func (e *Engine) AddHelper(name string, fn interface{}) {
	e.helpersMutex.Lock()
	defer e.helpersMutex.Unlock()
	e.helpers[name] = fn
}

// Render 渲染模板
func (e *Engine) Render(name string, data Data) (string, error) {
	tmpl, err := e.compile(name)
	if err != nil {
		return "", fmt.Errorf("failed to compile template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Compiled.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.String(), nil
}

// RenderToWriter 渲染模板到Writer
func (e *Engine) RenderToWriter(name string, data Data, w io.Writer) error {
	tmpl, err := e.compile(name)
	if err != nil {
		return fmt.Errorf("failed to compile template %s: %w", name, err)
	}

	return tmpl.Compiled.Execute(w, data)
}

// compile 编译模板
func (e *Engine) compile(name string) (*Template, error) {
	// 检查缓存
	if e.useCache {
		e.cacheMutex.RLock()
		if cached, exists := e.cache[name]; exists {
			e.cacheMutex.RUnlock()
			return cached, nil
		}
		e.cacheMutex.RUnlock()
	}

	// 读取模板文件
	content, err := e.readTemplateFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// 解析模板
	tmpl, err := e.parseTemplate(name, content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// 编译模板
	err = e.compileTemplate(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to compile template: %w", err)
	}

	// 缓存模板
	if e.useCache {
		e.cacheMutex.Lock()
		e.cache[name] = tmpl
		e.cacheMutex.Unlock()
	}

	return tmpl, nil
}

// readTemplateFile 读取模板文件
func (e *Engine) readTemplateFile(name string) (string, error) {
	filePath := filepath.Join(e.viewsPath, name+".blade.php")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// parseTemplate 解析模板
func (e *Engine) parseTemplate(name, content string) (*Template, error) {
	tmpl := &Template{
		Name:    name,
		Content: content,
		Blocks:  make(map[string]string),
	}

	// 解析 @extends 指令
	if extends := e.parseExtends(content); extends != "" {
		tmpl.Extends = extends
	}

	// 解析 @section 指令
	tmpl.Blocks = e.parseSections(content)

	// 解析 @include 指令
	content = e.parseIncludes(content)

	// 解析 @component 指令
	content = e.parseComponents(content)

	// 解析 @extends 和 @yield 指令
	content = e.parseLayouts(content)
	content = e.parseYields(content)

	// 解析 @stack, @push, @prepend 指令
	content = e.parseStacks(content)
	content = e.parsePush(content)
	content = e.parsePrepend(content)

	// 解析 @once 指令
	content = e.parseOnce(content)

	// 解析 @error 和 @old 指令
	content = e.parseError(content)
	content = e.parseOld(content)

	// 解析 @if, @foreach, @for 等控制结构
	content = e.parseControlStructures(content)

	// 解析 {{ }} 变量输出
	content = e.parseVariables(content)

	// 解析 @{{ }} 原始输出
	content = e.parseRawOutput(content)

	tmpl.Content = content
	return tmpl, nil
}

// parseExtends 解析 @extends 指令
func (e *Engine) parseExtends(content string) string {
	re := regexp.MustCompile(`@extends\s*\(\s*['"]([^'"]+)['"]\s*\)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		layoutName := matches[1]

		// 检查是否是注册的布局
		if _, exists := e.componentManager.GetLayout(layoutName); exists {
			return layoutName
		}

		// 否则尝试从文件读取
		return layoutName
	}
	return ""
}

// parseSections 解析 @section 指令
func (e *Engine) parseSections(content string) map[string]string {
	sections := make(map[string]string)

	// 匹配 @section('name') ... @endsection 或 @section('name') ... @stop
	re := regexp.MustCompile(`@section\s*\(\s*['"]([^'"]+)['"]\s*\)\s*(.*?)\s*(?:@endsection|@stop)`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 2 {
			sections[match[1]] = strings.TrimSpace(match[2])
		}
	}

	return sections
}

// parseControlStructures 解析控制结构
func (e *Engine) parseControlStructures(content string) string {
	// 解析 @if
	content = e.parseIfStatements(content)

	// 解析 @foreach
	content = e.parseForeachStatements(content)

	// 解析 @for
	content = e.parseForStatements(content)

	// 解析 @while
	content = e.parseWhileStatements(content)

	return content
}

// parseIfStatements 解析 @if 语句
func (e *Engine) parseIfStatements(content string) string {
	// @if (condition) ... @endif
	re := regexp.MustCompile(`@if\s*\((.*?)\)\s*(.*?)\s*@endif`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		// 提取条件和内容
		parts := re.FindStringSubmatch(match)
		if len(parts) > 2 {
			condition := strings.TrimSpace(parts[1])
			body := strings.TrimSpace(parts[2])
			// 将条件中的变量引用转换为Go模板语法
			condition = strings.ReplaceAll(condition, ".", "")
			return fmt.Sprintf("{{if %s}}%s{{end}}", condition, body)
		}
		return match
	})
}

// parseForeachStatements 解析 @foreach 语句
func (e *Engine) parseForeachStatements(content string) string {
	// @foreach ($items as $item) ... @endforeach
	re := regexp.MustCompile(`@foreach\s*\(\s*\$([^)]+)\s+as\s+\$([^)]+)\s*\)\s*(.*?)\s*@endforeach`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 3 {
			collection := strings.TrimSpace(parts[1])
			item := strings.TrimSpace(parts[2])
			body := strings.TrimSpace(parts[3])
			return fmt.Sprintf("{{range $%s := .%s}}%s{{end}}", item, collection, body)
		}
		return match
	})
}

// parseForStatements 解析 @for 语句
func (e *Engine) parseForStatements(content string) string {
	// @for ($i = 0; $i < 10; $i++) ... @endfor
	re := regexp.MustCompile(`@for\s*\(\s*\$([^=]+)\s*=\s*([^;]+);\s*\$([^<]+)\s*([^;]+);\s*\$([^)]+)\s*\)\s*(.*?)\s*@endfor`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 6 {
			init := strings.TrimSpace(parts[1]) + " := " + strings.TrimSpace(parts[2])
			condition := "$" + strings.TrimSpace(parts[3]) + " " + strings.TrimSpace(parts[4])
			increment := "$" + strings.TrimSpace(parts[5])
			body := strings.TrimSpace(parts[6])
			return fmt.Sprintf("{{%s; for %s; %s {}}%s{{}}", init, condition, increment, body)
		}
		return match
	})
}

// parseWhileStatements 解析 @while 语句
func (e *Engine) parseWhileStatements(content string) string {
	// @while (condition) ... @endwhile
	re := regexp.MustCompile(`@while\s*\((.*?)\)\s*(.*?)\s*@endwhile`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 2 {
			condition := strings.TrimSpace(parts[1])
			body := strings.TrimSpace(parts[2])
			return fmt.Sprintf("{{for %s {}}%s{{}}", condition, body)
		}
		return match
	})
}

// parseVariables 解析变量输出
func (e *Engine) parseVariables(content string) string {
	// 将 {{ $variable }} 转换为 {{ .variable }}
	re := regexp.MustCompile(`\{\{\s*\$([^}]+)\s*\}\}`)
	return re.ReplaceAllString(content, "{{ .$1 }}")
}

// parseRawOutput 解析原始输出
func (e *Engine) parseRawOutput(content string) string {
	// 将 @{{ }} 转换为 {{ }}
	re := regexp.MustCompile(`@\{\{\s*([^}]+)\s*\}\}`)
	return re.ReplaceAllString(content, "{{ $1 }}")
}

// compileTemplate 编译模板
func (e *Engine) compileTemplate(tmpl *Template) error {
	// 如果有继承，先编译父模板
	if tmpl.Extends != "" {
		// 检查是否是注册的布局
		if layout, exists := e.componentManager.GetLayout(tmpl.Extends); exists {
			// 使用注册的布局
			content := layout.Content
			for name, blockContent := range tmpl.Blocks {
				placeholder := fmt.Sprintf("@yield('%s')", name)
				content = strings.Replace(content, placeholder, blockContent, 1)
			}
			tmpl.Content = content
		} else {
			// 尝试从文件编译父模板
			parent, err := e.compile(tmpl.Extends)
			if err != nil {
				return fmt.Errorf("failed to compile parent template: %w", err)
			}

			// 将子模板的块内容注入到父模板中
			content := parent.Content
			for name, blockContent := range tmpl.Blocks {
				placeholder := fmt.Sprintf("{{template \"%s\" .}}", name)
				content = strings.Replace(content, placeholder, blockContent, 1)
			}
			tmpl.Content = content
		}
	}

	// 创建Go模板
	goTemplate := template.New(tmpl.Name)

	// 添加助手函数
	e.helpersMutex.RLock()
	for name, fn := range e.helpers {
		goTemplate = goTemplate.Funcs(template.FuncMap{name: fn})
	}
	e.helpersMutex.RUnlock()

	// 解析模板
	parsed, err := goTemplate.Parse(tmpl.Content)
	if err != nil {
		return fmt.Errorf("failed to parse Go template: %w", err)
	}

	tmpl.Compiled = parsed
	return nil
}

// ClearCache 清除缓存
func (e *Engine) ClearCache() {
	e.cacheMutex.Lock()
	defer e.cacheMutex.Unlock()
	e.cache = make(map[string]*Template)
}

// GetCachedTemplate 获取缓存的模板
func (e *Engine) GetCachedTemplate(name string) (*Template, bool) {
	e.cacheMutex.RLock()
	defer e.cacheMutex.RUnlock()
	tmpl, exists := e.cache[name]
	return tmpl, exists
}

// GetComponentManager 获取组件管理器
func (e *Engine) GetComponentManager() *ComponentManager {
	return e.componentManager
}
