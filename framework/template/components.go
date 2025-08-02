package template

import (
	"fmt"
	"regexp"
	"strings"
)

// Component 组件定义
type Component struct {
	Name    string
	Content string
	Props   map[string]interface{}
}

// Layout 布局模板
type Layout struct {
	Name     string
	Content  string
	Sections map[string]string
}

// View 部分视图
type View struct {
	Name    string
	Content string
	Data    Data
}

// ComponentManager 组件管理器
type ComponentManager struct {
	components map[string]*Component
	layouts    map[string]*Layout
	views      map[string]*View
}

// NewComponentManager 创建新的组件管理器
func NewComponentManager() *ComponentManager {
	return &ComponentManager{
		components: make(map[string]*Component),
		layouts:    make(map[string]*Layout),
		views:      make(map[string]*View),
	}
}

// RegisterComponent 注册组件
func (cm *ComponentManager) RegisterComponent(name, content string) {
	cm.components[name] = &Component{
		Name:    name,
		Content: content,
		Props:   make(map[string]interface{}),
	}
}

// RegisterLayout 注册布局
func (cm *ComponentManager) RegisterLayout(name, content string) {
	cm.layouts[name] = &Layout{
		Name:     name,
		Content:  content,
		Sections: make(map[string]string),
	}
}

// RegisterView 注册部分视图
func (cm *ComponentManager) RegisterView(name, content string) {
	cm.views[name] = &View{
		Name:    name,
		Content: content,
		Data:    make(Data),
	}
}

// GetComponent 获取组件
func (cm *ComponentManager) GetComponent(name string) (*Component, bool) {
	component, exists := cm.components[name]
	return component, exists
}

// GetLayout 获取布局
func (cm *ComponentManager) GetLayout(name string) (*Layout, bool) {
	layout, exists := cm.layouts[name]
	return layout, exists
}

// GetView 获取部分视图
func (cm *ComponentManager) GetView(name string) (*View, bool) {
	view, exists := cm.views[name]
	return view, exists
}

// parseComponents 解析组件
func (e *Engine) parseComponents(content string) string {
	// 解析 @component 指令
	re := regexp.MustCompile(`@component\s*\(\s*['"]([^'"]+)['"]\s*\)\s*(.*?)\s*@endcomponent`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 2 {
			componentName := parts[1]
			componentContent := strings.TrimSpace(parts[2])

			// 解析组件属性
			props := e.parseComponentProps(componentContent)

			// 获取组件定义
			if component, exists := e.componentManager.GetComponent(componentName); exists {
				return e.renderComponent(component, props)
			}
		}
		return match
	})
}

// parseComponentProps 解析组件属性
func (e *Engine) parseComponentProps(content string) map[string]interface{} {
	props := make(map[string]interface{})

	// 解析 @slot 指令
	slotRe := regexp.MustCompile(`@slot\s*\(\s*['"]([^'"]+)['"]\s*\)\s*(.*?)\s*@endslot`)
	matches := slotRe.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 2 {
			slotName := match[1]
			slotContent := strings.TrimSpace(match[2])
			props[slotName] = slotContent
		}
	}

	return props
}

// renderComponent 渲染组件
func (e *Engine) renderComponent(component *Component, props map[string]interface{}) string {
	content := component.Content

	// 替换插槽
	for slotName, slotContent := range props {
		placeholder := fmt.Sprintf("@slot('%s')", slotName)
		content = strings.Replace(content, placeholder, slotContent.(string), -1)
	}

	return content
}

// parseLayouts 解析布局
func (e *Engine) parseLayouts(content string) string {
	// 解析 @extends 指令
	extendsRe := regexp.MustCompile(`@extends\s*\(\s*['"]([^'"]+)['"]\s*\)`)
	matches := extendsRe.FindStringSubmatch(content)
	if len(matches) > 1 {
		layoutName := matches[1]

		// 获取布局定义
		if layout, exists := e.componentManager.GetLayout(layoutName); exists {
			return e.renderLayout(layout, content)
		}
	}

	return content
}

// renderLayout 渲染布局
func (e *Engine) renderLayout(layout *Layout, content string) string {
	layoutContent := layout.Content

	// 解析子模板的 @section 指令
	sections := e.parseSections(content)

	// 将子模板的块内容注入到布局中
	for sectionName, sectionContent := range sections {
		placeholder := fmt.Sprintf("@yield('%s')", sectionName)
		layoutContent = strings.Replace(layoutContent, placeholder, sectionContent, -1)
	}

	return layoutContent
}

// parseYields 解析 @yield 指令
func (e *Engine) parseYields(content string) string {
	// @yield('section_name') 指令
	re := regexp.MustCompile(`@yield\s*\(\s*['"]([^'"]+)['"]\s*\)`)
	return re.ReplaceAllString(content, "{{template \"$1\" .}}")
}

// parseIncludes 解析 @include 指令（增强版）
func (e *Engine) parseIncludes(content string) string {
	// @include('view_name', ['key' => 'value']) 指令
	re := regexp.MustCompile(`@include\s*\(\s*['"]([^'"]+)['"]\s*(?:,\s*(\[[^\]]+\]))?\s*\)`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 1 {
			viewName := parts[1]

			// 解析传递的数据
			var data Data
			if len(parts) > 2 && parts[2] != "" {
				data = e.parseIncludeData(parts[2])
			}

			// 获取部分视图
			if view, exists := e.componentManager.GetView(viewName); exists {
				return e.renderView(view, data)
			}

			// 尝试读取模板文件
			if content, err := e.readTemplateFile(viewName); err == nil {
				return content
			}
		}
		return match
	})
}

// parseIncludeData 解析包含的数据
func (e *Engine) parseIncludeData(dataStr string) Data {
	data := make(Data)

	// 解析 ['key' => 'value'] 格式
	re := regexp.MustCompile(`['"]([^'"]+)['"]\s*=>\s*['"]([^'"]+)['"]`)
	matches := re.FindAllStringSubmatch(dataStr, -1)

	for _, match := range matches {
		if len(match) > 2 {
			key := match[1]
			value := match[2]
			data[key] = value
		}
	}

	return data
}

// renderView 渲染部分视图
func (e *Engine) renderView(view *View, data Data) string {
	// 合并视图数据
	mergedData := make(Data)
	for k, v := range view.Data {
		mergedData[k] = v
	}
	for k, v := range data {
		mergedData[k] = v
	}

	// 渲染视图内容
	content := view.Content

	// 替换变量
	for key, value := range mergedData {
		placeholder := fmt.Sprintf("{{ .%s }}", key)
		content = strings.Replace(content, placeholder, fmt.Sprintf("%v", value), -1)
	}

	return content
}

// parseStacks 解析 @stack 指令
func (e *Engine) parseStacks(content string) string {
	// @stack('name') 指令
	re := regexp.MustCompile(`@stack\s*\(\s*['"]([^'"]+)['"]\s*\)`)
	return re.ReplaceAllString(content, "{{template \"stack_$1\" .}}")
}

// parsePush 解析 @push 指令
func (e *Engine) parsePush(content string) string {
	// @push('name') ... @endpush 指令
	re := regexp.MustCompile(`@push\s*\(\s*['"]([^'"]+)['"]\s*\)\s*(.*?)\s*@endpush`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 2 {
			stackContent := strings.TrimSpace(parts[2])
			return stackContent
		}
		return match
	})
}

// parsePrepend 解析 @prepend 指令
func (e *Engine) parsePrepend(content string) string {
	// @prepend('name') ... @endprepend 指令
	re := regexp.MustCompile(`@prepend\s*\(\s*['"]([^'"]+)['"]\s*\)\s*(.*?)\s*@endprepend`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 2 {
			stackName := parts[1]
			stackContent := strings.TrimSpace(parts[2])
			return fmt.Sprintf("{{define \"stack_%s\"}}%s{{end}}", stackName, stackContent)
		}
		return match
	})
}

// parseOnce 解析 @once 指令
func (e *Engine) parseOnce(content string) string {
	// @once ... @endonce 指令
	re := regexp.MustCompile(`@once\s*(.*?)\s*@endonce`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 1 {
			onceContent := strings.TrimSpace(parts[1])
			// 生成唯一标识符
			uniqueID := fmt.Sprintf("once_%d", len(match))
			return fmt.Sprintf("{{if not .%s}}{{$%s := true}}%s{{end}}", uniqueID, uniqueID, onceContent)
		}
		return match
	})
}

// parseError 解析 @error 指令
func (e *Engine) parseError(content string) string {
	// @error('field') 指令
	re := regexp.MustCompile(`@error\s*\(\s*['"]([^'"]+)['"]\s*\)`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 1 {
			fieldName := parts[1]
			return fmt.Sprintf("{{if .errors.%s}}<span class=\"error\">{{.errors.%s}}</span>{{end}}", fieldName, fieldName)
		}
		return match
	})
}

// parseOld 解析 @old 指令
func (e *Engine) parseOld(content string) string {
	// @old('field', 'default') 指令
	re := regexp.MustCompile(`@old\s*\(\s*['"]([^'"]+)['"]\s*(?:,\s*['"]([^'"]*)['"])?\s*\)`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 1 {
			fieldName := parts[1]
			defaultValue := ""
			if len(parts) > 2 {
				defaultValue = parts[2]
			}
			return fmt.Sprintf("{{if .old.%s}}{{.old.%s}}{{else}}%s{{end}}", fieldName, fieldName, defaultValue)
		}
		return match
	})
}
