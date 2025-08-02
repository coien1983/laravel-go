package routing

import (
	"strings"
)

// RadixNode Radix树节点
type RadixNode struct {
	path     string
	children map[string]*RadixNode
	handlers map[string]interface{} // method -> handler
	params   []string               // 参数名列表
	isParam  bool                   // 是否为参数节点
}

// RadixTree Radix树路由匹配器
type RadixTree struct {
	root *RadixNode
}

// NewRadixTree 创建新的Radix树
func NewRadixTree() *RadixTree {
	return &RadixTree{
		root: &RadixNode{
			children: make(map[string]*RadixNode),
			handlers: make(map[string]interface{}),
		},
	}
}

// Insert 插入路由
func (rt *RadixTree) Insert(method, path string, handler interface{}) {
	parts := rt.splitPath(path)
	rt.insertNode(rt.root, parts, method, handler, 0)
}

// insertNode 递归插入节点
func (rt *RadixTree) insertNode(node *RadixNode, parts []string, method string, handler interface{}, depth int) {
	if depth == len(parts) {
		// 到达叶子节点，设置处理器
		node.handlers[method] = handler
		return
	}

	part := parts[depth]
	isParam := strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}")
	paramName := ""

	if isParam {
		paramName = strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")
		part = ":" // 使用冒号表示参数节点
	}

	// 查找或创建子节点
	child, exists := node.children[part]
	if !exists {
		child = &RadixNode{
			path:     part,
			children: make(map[string]*RadixNode),
			handlers: make(map[string]interface{}),
			isParam:  isParam,
		}
		if isParam {
			child.params = append(child.params, paramName)
		}
		node.children[part] = child
	}

	// 递归插入
	rt.insertNode(child, parts, method, handler, depth+1)
}

// Match 匹配路由
func (rt *RadixTree) Match(method, path string) (interface{}, map[string]string, bool) {
	parts := rt.splitPath(path)
	params := make(map[string]string)

	handler, found := rt.matchNode(rt.root, parts, method, params, 0)
	return handler, params, found
}

// matchNode 递归匹配节点
func (rt *RadixTree) matchNode(node *RadixNode, parts []string, method string, params map[string]string, depth int) (interface{}, bool) {
	if depth == len(parts) {
		// 到达叶子节点，查找处理器
		if handler, exists := node.handlers[method]; exists {
			return handler, true
		}
		return nil, false
	}

	part := parts[depth]

	// 尝试精确匹配
	if child, exists := node.children[part]; exists {
		if handler, found := rt.matchNode(child, parts, method, params, depth+1); found {
			return handler, true
		}
	}

	// 尝试参数匹配
	for _, paramChild := range node.children {
		if paramChild.isParam {
			// 设置参数值
			if len(paramChild.params) > 0 {
				paramName := paramChild.params[0]
				params[paramName] = part
			}

			if handler, found := rt.matchNode(paramChild, parts, method, params, depth+1); found {
				return handler, true
			}
		}
	}

	return nil, false
}

// splitPath 分割路径
func (rt *RadixTree) splitPath(path string) []string {
	if path == "/" {
		return []string{""}
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 1 && parts[0] == "" {
		return []string{""}
	}

	return parts
}

// GetAllRoutes 获取所有路由
func (rt *RadixTree) GetAllRoutes() []Route {
	var routes []Route
	rt.collectRoutes(rt.root, "", &routes)
	return routes
}

// collectRoutes 收集所有路由
func (rt *RadixTree) collectRoutes(node *RadixNode, currentPath string, routes *[]Route) {
	if len(node.handlers) > 0 {
		for method, handler := range node.handlers {
			route := Route{
				Method:  method,
				Path:    currentPath,
				Handler: handler,
			}
			*routes = append(*routes, route)
		}
	}

	for path, child := range node.children {
		childPath := currentPath
		if path == ":" {
			childPath += "/{param}"
		} else {
			childPath += "/" + path
		}
		rt.collectRoutes(child, childPath, routes)
	}
}
