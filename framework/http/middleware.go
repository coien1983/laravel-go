package http

import (
	"sync"
)

// Handler 处理器接口
type Handler interface {
	// Handle 处理请求
	Handle(request Request) Response
}

// Middleware 中间件接口
type Middleware interface {
	// Handle 处理请求
	Handle(request Request, next Next) Response
}

// MiddlewareFunc 中间件函数类型
type MiddlewareFunc func(request Request, next Next) Response

// Handle 实现 Middleware 接口
func (f MiddlewareFunc) Handle(request Request, next Next) Response {
	return f(request, next)
}

// Next 下一个中间件或处理器的函数类型
type Next func(request Request) Response

// Pipeline 中间件管道
type Pipeline struct {
	middlewares []Middleware
	mutex       sync.RWMutex
}

// NewPipeline 创建新的中间件管道
func NewPipeline() *Pipeline {
	return &Pipeline{
		middlewares: make([]Middleware, 0),
	}
}

// Use 添加中间件
func (p *Pipeline) Use(middleware ...Middleware) *Pipeline {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.middlewares = append(p.middlewares, middleware...)
	return p
}

// UseFunc 添加中间件函数
func (p *Pipeline) UseFunc(middlewareFuncs ...MiddlewareFunc) *Pipeline {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, mf := range middlewareFuncs {
		p.middlewares = append(p.middlewares, mf)
	}
	return p
}

// Process 处理请求
func (p *Pipeline) Process(request Request, handler Handler) Response {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// 创建处理器函数
	handlerFunc := func(req Request) Response {
		return handler.Handle(req)
	}

	// 构建中间件链
	next := handlerFunc
	for i := len(p.middlewares) - 1; i >= 0; i-- {
		currentMiddleware := p.middlewares[i]
		currentNext := next
		next = func(req Request) Response {
			return currentMiddleware.Handle(req, currentNext)
		}
	}

	// 执行中间件链
	return next(request)
}

// GetMiddlewares 获取所有中间件
func (p *Pipeline) GetMiddlewares() []Middleware {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	middlewares := make([]Middleware, len(p.middlewares))
	copy(middlewares, p.middlewares)
	return middlewares
}

// Clear 清空中间件
func (p *Pipeline) Clear() *Pipeline {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.middlewares = make([]Middleware, 0)
	return p
}

// MiddlewareManager 中间件管理器
type MiddlewareManager struct {
	globalMiddlewares []Middleware
	routeMiddlewares  map[string][]Middleware
	groupMiddlewares  map[string][]Middleware
	mutex             sync.RWMutex
}

// NewMiddlewareManager 创建新的中间件管理器
func NewMiddlewareManager() *MiddlewareManager {
	return &MiddlewareManager{
		globalMiddlewares: make([]Middleware, 0),
		routeMiddlewares:  make(map[string][]Middleware),
		groupMiddlewares:  make(map[string][]Middleware),
	}
}

// UseGlobal 添加全局中间件
func (mm *MiddlewareManager) UseGlobal(middleware ...Middleware) *MiddlewareManager {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	mm.globalMiddlewares = append(mm.globalMiddlewares, middleware...)
	return mm
}

// UseGlobalFunc 添加全局中间件函数
func (mm *MiddlewareManager) UseGlobalFunc(middlewareFuncs ...MiddlewareFunc) *MiddlewareManager {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	for _, mf := range middlewareFuncs {
		mm.globalMiddlewares = append(mm.globalMiddlewares, mf)
	}
	return mm
}

// UseForRoute 为特定路由添加中间件
func (mm *MiddlewareManager) UseForRoute(routeName string, middleware ...Middleware) *MiddlewareManager {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	mm.routeMiddlewares[routeName] = append(mm.routeMiddlewares[routeName], middleware...)
	return mm
}

// UseForGroup 为路由分组添加中间件
func (mm *MiddlewareManager) UseForGroup(groupName string, middleware ...Middleware) *MiddlewareManager {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	mm.groupMiddlewares[groupName] = append(mm.groupMiddlewares[groupName], middleware...)
	return mm
}

// GetMiddlewaresForRoute 获取路由的中间件
func (mm *MiddlewareManager) GetMiddlewaresForRoute(routeName, groupName string) []Middleware {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	var middlewares []Middleware

	// 添加全局中间件
	middlewares = append(middlewares, mm.globalMiddlewares...)

	// 添加分组中间件
	if groupMiddlewares, exists := mm.groupMiddlewares[groupName]; exists {
		middlewares = append(middlewares, groupMiddlewares...)
	}

	// 添加路由中间件
	if routeMiddlewares, exists := mm.routeMiddlewares[routeName]; exists {
		middlewares = append(middlewares, routeMiddlewares...)
	}

	return middlewares
}

// ConditionalMiddleware 条件中间件
type ConditionalMiddleware struct {
	condition  func(Request) bool
	middleware Middleware
}

// NewConditionalMiddleware 创建条件中间件
func NewConditionalMiddleware(condition func(Request) bool, middleware Middleware) *ConditionalMiddleware {
	return &ConditionalMiddleware{
		condition:  condition,
		middleware: middleware,
	}
}

// Handle 实现 Middleware 接口
func (cm *ConditionalMiddleware) Handle(request Request, next Next) Response {
	if cm.condition(request) {
		return cm.middleware.Handle(request, next)
	}
	return next(request)
}

// PriorityMiddleware 优先级中间件
type PriorityMiddleware struct {
	priority   int
	middleware Middleware
}

// NewPriorityMiddleware 创建优先级中间件
func NewPriorityMiddleware(priority int, middleware Middleware) *PriorityMiddleware {
	return &PriorityMiddleware{
		priority:   priority,
		middleware: middleware,
	}
}

// Handle 实现 Middleware 接口
func (pm *PriorityMiddleware) Handle(request Request, next Next) Response {
	return pm.middleware.Handle(request, next)
}

// GetPriority 获取优先级
func (pm *PriorityMiddleware) GetPriority() int {
	return pm.priority
}

// PriorityPipeline 优先级中间件管道
type PriorityPipeline struct {
	middlewares []*PriorityMiddleware
	mutex       sync.RWMutex
}

// NewPriorityPipeline 创建优先级中间件管道
func NewPriorityPipeline() *PriorityPipeline {
	return &PriorityPipeline{
		middlewares: make([]*PriorityMiddleware, 0),
	}
}

// Use 添加优先级中间件
func (pp *PriorityPipeline) Use(priority int, middleware Middleware) *PriorityPipeline {
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	priorityMiddleware := NewPriorityMiddleware(priority, middleware)
	pp.middlewares = append(pp.middlewares, priorityMiddleware)

	// 按优先级排序
	pp.sortByPriority()
	return pp
}

// Process 处理请求
func (pp *PriorityPipeline) Process(request Request, handler Handler) Response {
	pp.mutex.RLock()
	defer pp.mutex.RUnlock()

	// 创建处理器函数
	handlerFunc := func(req Request) Response {
		return handler.Handle(req)
	}

	// 构建中间件链
	next := handlerFunc
	for i := len(pp.middlewares) - 1; i >= 0; i-- {
		currentMiddleware := pp.middlewares[i]
		currentNext := next
		next = func(req Request) Response {
			return currentMiddleware.Handle(req, currentNext)
		}
	}

	// 执行中间件链
	return next(request)
}

// sortByPriority 按优先级排序
func (pp *PriorityPipeline) sortByPriority() {
	// 简单的冒泡排序，按优先级从高到低排序
	for i := 0; i < len(pp.middlewares)-1; i++ {
		for j := 0; j < len(pp.middlewares)-i-1; j++ {
			if pp.middlewares[j].GetPriority() < pp.middlewares[j+1].GetPriority() {
				pp.middlewares[j], pp.middlewares[j+1] = pp.middlewares[j+1], pp.middlewares[j]
			}
		}
	}
}
