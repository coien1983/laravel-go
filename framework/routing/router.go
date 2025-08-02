package routing

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Router 路由接口
type Router interface {
	// 基础路由方法
	Get(path string, handler interface{}) Router
	Post(path string, handler interface{}) Router
	Put(path string, handler interface{}) Router
	Delete(path string, handler interface{}) Router
	Patch(path string, handler interface{}) Router
	Options(path string, handler interface{}) Router
	Head(path string, handler interface{}) Router

	// 路由分组
	Group(prefix string, callback func(Router))

	// 中间件
	Use(middleware ...interface{})

	// 路由参数
	Where(name, pattern string)

	// 路由缓存
	Cache(ttl int) Router

	// 获取路由列表
	GetRoutes() []Route

	// 匹配路由
	Match(method, path string) (Route, bool)
}

// Route 路由信息
type Route struct {
	Method      string
	Path        string
	Handler     interface{}
	Middleware  []interface{}
	Parameters  map[string]string
	Constraints map[string]string
	CacheTTL    int
	Group       string
}

// router 路由实现
type router struct {
	routes      []Route
	groups      map[string]*routeGroup
	middleware  []interface{}
	constraints map[string]string
	cache       *routeCache
	radixTree   *RadixTree
	mutex       sync.RWMutex
}

// routeGroup 路由分组
type routeGroup struct {
	prefix      string
	middleware  []interface{}
	constraints map[string]string
	routes      []Route
}

// routeCache 路由缓存
type routeCache struct {
	cache map[string]Route
	mutex sync.RWMutex
}

// NewRouter 创建新的路由器
func NewRouter() Router {
	return &router{
		routes:      make([]Route, 0),
		groups:      make(map[string]*routeGroup),
		middleware:  make([]interface{}, 0),
		constraints: make(map[string]string),
		cache: &routeCache{
			cache: make(map[string]Route),
		},
		radixTree: NewRadixTree(),
	}
}

// Get 注册 GET 路由
func (r *router) Get(path string, handler interface{}) Router {
	r.addRoute(http.MethodGet, path, handler)
	return r
}

// Post 注册 POST 路由
func (r *router) Post(path string, handler interface{}) Router {
	r.addRoute(http.MethodPost, path, handler)
	return r
}

// Put 注册 PUT 路由
func (r *router) Put(path string, handler interface{}) Router {
	r.addRoute(http.MethodPut, path, handler)
	return r
}

// Delete 注册 DELETE 路由
func (r *router) Delete(path string, handler interface{}) Router {
	r.addRoute(http.MethodDelete, path, handler)
	return r
}

// Patch 注册 PATCH 路由
func (r *router) Patch(path string, handler interface{}) Router {
	r.addRoute(http.MethodPatch, path, handler)
	return r
}

// Options 注册 OPTIONS 路由
func (r *router) Options(path string, handler interface{}) Router {
	r.addRoute(http.MethodOptions, path, handler)
	return r
}

// Head 注册 HEAD 路由
func (r *router) Head(path string, handler interface{}) Router {
	r.addRoute(http.MethodHead, path, handler)
	return r
}

// Group 创建路由分组
func (r *router) Group(prefix string, callback func(Router)) {
	group := &routeGroup{
		prefix:      prefix,
		middleware:  make([]interface{}, 0),
		constraints: make(map[string]string),
		routes:      make([]Route, 0),
	}

	// 创建分组路由器
	groupRouter := &router{
		routes:      make([]Route, 0),
		groups:      make(map[string]*routeGroup),
		middleware:  group.middleware,
		constraints: group.constraints,
		cache:       r.cache,
		radixTree:   NewRadixTree(),
	}

	// 执行回调函数
	callback(groupRouter)

	// 将分组路由添加到主路由器
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, route := range groupRouter.routes {
		route.Path = prefix + route.Path
		route.Group = prefix
		r.routes = append(r.routes, route)

		// 同时添加到Radix Tree
		r.radixTree.Insert(route.Method, route.Path, route.Handler)
	}
}

// Use 添加中间件
func (r *router) Use(middleware ...interface{}) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.middleware = append(r.middleware, middleware...)
}

// Where 添加路由参数约束
func (r *router) Where(name, pattern string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.constraints[name] = pattern
}

// Cache 设置路由缓存
func (r *router) Cache(ttl int) Router {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// 为当前路由设置缓存TTL
	if len(r.routes) > 0 {
		r.routes[len(r.routes)-1].CacheTTL = ttl
	}
	return r
}

// GetRoutes 获取所有路由
func (r *router) GetRoutes() []Route {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.routes
}

// Match 匹配路由
func (r *router) Match(method, path string) (Route, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 检查缓存
	cacheKey := fmt.Sprintf("%s:%s", method, path)
	if route, exists := r.cache.get(cacheKey); exists {
		return route, true
	}

	// 使用Radix Tree进行匹配
	if _, params, found := r.radixTree.Match(method, path); found {
		// 查找对应的路由信息
		for _, route := range r.routes {
			if route.Method == method && r.matchPath(route.Path, path) {
				// 创建路由副本并设置参数
				matchedRoute := route
				matchedRoute.Parameters = params
				// 缓存结果
				r.cache.set(cacheKey, matchedRoute)
				return matchedRoute, true
			}
		}
	}

	return Route{}, false
}

// addRoute 添加路由
func (r *router) addRoute(method, path string, handler interface{}) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	route := Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Middleware:  make([]interface{}, len(r.middleware)),
		Parameters:  make(map[string]string),
		Constraints: make(map[string]string),
		CacheTTL:    0,
	}

	// 复制中间件
	copy(route.Middleware, r.middleware)

	// 复制约束
	for k, v := range r.constraints {
		route.Constraints[k] = v
	}

	r.routes = append(r.routes, route)

	// 同时添加到Radix Tree
	r.radixTree.Insert(method, path, handler)
}

// matchPath 匹配路径
func (r *router) matchPath(pattern, path string) bool {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, patternPart := range patternParts {
		if strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
			// 参数匹配，检查约束
			paramName := strings.TrimSuffix(strings.TrimPrefix(patternPart, "{"), "}")
			if constraint, exists := r.constraints[paramName]; exists {
				if !r.validateConstraint(pathParts[i], constraint) {
					return false
				}
			}
			continue
		}
		if patternPart != pathParts[i] {
			return false
		}
	}

	return true
}

// validateConstraint 验证参数约束
func (r *router) validateConstraint(value, pattern string) bool {
	// 简单的正则表达式匹配实现
	// 这里可以扩展为更复杂的正则表达式引擎
	if pattern == "[0-9]+" {
		for _, char := range value {
			if char < '0' || char > '9' {
				return false
			}
		}
		return true
	}
	return true
}

// extractParameters 提取参数
func (r *router) extractParameters(pattern, path string) map[string]string {
	params := make(map[string]string)
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	for i, patternPart := range patternParts {
		if strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
			paramName := strings.TrimSuffix(strings.TrimPrefix(patternPart, "{"), "}")
			if i < len(pathParts) {
				params[paramName] = pathParts[i]
			}
		}
	}

	return params
}

// routeCache 方法
func (rc *routeCache) get(key string) (Route, bool) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()
	route, exists := rc.cache[key]
	return route, exists
}

func (rc *routeCache) set(key string, route Route) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.cache[key] = route
}

func (rc *routeCache) clear() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.cache = make(map[string]Route)
}
