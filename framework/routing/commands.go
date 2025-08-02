package routing

import (
	"fmt"
	"sort"
	"strings"
)

// RouteCommand 路由命令接口
type RouteCommand interface {
	List() string
	ListByMethod(method string) string
	ListByGroup(group string) string
	Show(path string) string
}

// routeCommand 路由命令实现
type routeCommand struct {
	router Router
}

// NewRouteCommand 创建新的路由命令
func NewRouteCommand(router Router) RouteCommand {
	return &routeCommand{
		router: router,
	}
}

// List 列出所有路由
func (rc *routeCommand) List() string {
	routes := rc.router.GetRoutes()

	if len(routes) == 0 {
		return "No routes found."
	}

	// 按方法排序
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Method != routes[j].Method {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})

	var output strings.Builder
	output.WriteString("Route List:\n")
	output.WriteString("===========\n\n")

	output.WriteString(fmt.Sprintf("%-8s %-40s %-20s %s\n", "Method", "Path", "Group", "Handler"))
	output.WriteString(strings.Repeat("-", 80) + "\n")

	for _, route := range routes {
		handlerName := rc.getHandlerName(route.Handler)
		group := route.Group
		if group == "" {
			group = "-"
		}
		output.WriteString(fmt.Sprintf("%-8s %-40s %-20s %s\n",
			route.Method,
			route.Path,
			group,
			handlerName))
	}

	return output.String()
}

// ListByMethod 按方法列出路由
func (rc *routeCommand) ListByMethod(method string) string {
	routes := rc.router.GetRoutes()

	var filteredRoutes []Route
	for _, route := range routes {
		if strings.ToUpper(route.Method) == strings.ToUpper(method) {
			filteredRoutes = append(filteredRoutes, route)
		}
	}

	if len(filteredRoutes) == 0 {
		return fmt.Sprintf("No routes found for method: %s", method)
	}

	// 按路径排序
	sort.Slice(filteredRoutes, func(i, j int) bool {
		return filteredRoutes[i].Path < filteredRoutes[j].Path
	})

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Routes for %s method:\n", strings.ToUpper(method)))
	output.WriteString("========================\n\n")

	output.WriteString(fmt.Sprintf("%-40s %-20s %s\n", "Path", "Group", "Handler"))
	output.WriteString(strings.Repeat("-", 70) + "\n")

	for _, route := range filteredRoutes {
		handlerName := rc.getHandlerName(route.Handler)
		group := route.Group
		if group == "" {
			group = "-"
		}
		output.WriteString(fmt.Sprintf("%-40s %-20s %s\n",
			route.Path,
			group,
			handlerName))
	}

	return output.String()
}

// ListByGroup 按分组列出路由
func (rc *routeCommand) ListByGroup(group string) string {
	routes := rc.router.GetRoutes()

	var filteredRoutes []Route
	for _, route := range routes {
		if route.Group == group {
			filteredRoutes = append(filteredRoutes, route)
		}
	}

	if len(filteredRoutes) == 0 {
		return fmt.Sprintf("No routes found for group: %s", group)
	}

	// 按方法排序
	sort.Slice(filteredRoutes, func(i, j int) bool {
		if filteredRoutes[i].Method != filteredRoutes[j].Method {
			return filteredRoutes[i].Method < filteredRoutes[j].Method
		}
		return filteredRoutes[i].Path < filteredRoutes[j].Path
	})

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Routes for group: %s\n", group))
	output.WriteString("=====================\n\n")

	output.WriteString(fmt.Sprintf("%-8s %-40s %s\n", "Method", "Path", "Handler"))
	output.WriteString(strings.Repeat("-", 60) + "\n")

	for _, route := range filteredRoutes {
		handlerName := rc.getHandlerName(route.Handler)
		output.WriteString(fmt.Sprintf("%-8s %-40s %s\n",
			route.Method,
			route.Path,
			handlerName))
	}

	return output.String()
}

// Show 显示特定路由详情
func (rc *routeCommand) Show(path string) string {
	routes := rc.router.GetRoutes()

	var foundRoutes []Route
	for _, route := range routes {
		if route.Path == path {
			foundRoutes = append(foundRoutes, route)
		}
	}

	if len(foundRoutes) == 0 {
		return fmt.Sprintf("No routes found for path: %s", path)
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Route details for: %s\n", path))
	output.WriteString("=====================\n\n")

	for _, route := range foundRoutes {
		output.WriteString(fmt.Sprintf("Method: %s\n", route.Method))
		output.WriteString(fmt.Sprintf("Path: %s\n", route.Path))
		output.WriteString(fmt.Sprintf("Handler: %s\n", rc.getHandlerName(route.Handler)))

		if route.Group != "" {
			output.WriteString(fmt.Sprintf("Group: %s\n", route.Group))
		}

		if len(route.Parameters) > 0 {
			output.WriteString("Parameters:\n")
			for name, value := range route.Parameters {
				output.WriteString(fmt.Sprintf("  %s: %s\n", name, value))
			}
		}

		if len(route.Constraints) > 0 {
			output.WriteString("Constraints:\n")
			for name, pattern := range route.Constraints {
				output.WriteString(fmt.Sprintf("  %s: %s\n", name, pattern))
			}
		}

		if route.CacheTTL > 0 {
			output.WriteString(fmt.Sprintf("Cache TTL: %d seconds\n", route.CacheTTL))
		}

		if len(route.Middleware) > 0 {
			output.WriteString("Middleware:\n")
			for _, middleware := range route.Middleware {
				output.WriteString(fmt.Sprintf("  - %s\n", rc.getHandlerName(middleware)))
			}
		}

		output.WriteString("\n")
	}

	return output.String()
}

// getHandlerName 获取处理器名称
func (rc *routeCommand) getHandlerName(handler interface{}) string {
	if handler == nil {
		return "nil"
	}

	// 尝试获取函数名
	if _, ok := handler.(func()); ok {
		return "anonymous function"
	}

	// 尝试获取结构体名称
	handlerType := fmt.Sprintf("%T", handler)
	if strings.Contains(handlerType, ".") {
		parts := strings.Split(handlerType, ".")
		return parts[len(parts)-1]
	}

	return handlerType
}
