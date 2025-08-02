package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"laravel-go/framework/config"
	"laravel-go/framework/container"
	"laravel-go/framework/log"
	"laravel-go/framework/routing"
)

// Server HTTP服务器接口
type Server interface {
	// 启动服务器
	Start() error
	// 停止服务器
	Stop() error
	// 获取路由器
	Router() routing.Router
	// 设置中间件
	Use(middleware ...string) Server
	// 设置静态文件
	Static(path, dir string) Server
}

// server HTTP服务器实现
type server struct {
	config     config.Config
	container  container.Container
	router     routing.Router
	httpServer *http.Server
	middleware []string
	static     map[string]string
}

// NewServer 创建新的HTTP服务器
func NewServer(config config.Config, container container.Container) Server {
	return &server{
		config:     config,
		container:  container,
		router:     routing.NewRouter(),
		middleware: make([]string, 0),
		static:     make(map[string]string),
	}
}

// Start 启动服务器
func (s *server) Start() error {
	// 获取配置
	host := s.config.Get("app.host", "localhost").(string)
	port := s.config.Get("app.port", "8080").(string)
	addr := fmt.Sprintf("%s:%s", host, port)

	// 创建HTTP服务器
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.createHandler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 记录启动日志
	log.Info("HTTP server starting on "+addr, nil)

	// 启动服务器
	return s.httpServer.ListenAndServe()
}

// Stop 停止服务器
func (s *server) Stop() error {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// Router 获取路由器
func (s *server) Router() routing.Router {
	return s.router
}

// Use 设置中间件
func (s *server) Use(middleware ...string) Server {
	s.middleware = append(s.middleware, middleware...)
	return s
}

// Static 设置静态文件
func (s *server) Static(path, dir string) Server {
	s.static[path] = dir
	return s
}

// createHandler 创建HTTP处理器
func (s *server) createHandler() http.Handler {
	// 创建多路复用器
	mux := http.NewServeMux()

	// 添加静态文件处理
	for path, dir := range s.static {
		mux.Handle(path+"/", http.StripPrefix(path, http.FileServer(http.Dir(dir))))
	}

	// 添加路由处理
	mux.HandleFunc("/", s.handleRequest)

	return mux
}

// handleRequest 处理HTTP请求
func (s *server) handleRequest(w http.ResponseWriter, r *http.Request) {
	// 匹配路由
	route, found := s.router.Match(r.Method, r.URL.Path)
	if !found {
		// 404 处理
		s.handleNotFound(w, r)
		return
	}

	// 执行中间件链
	s.executeMiddlewareChain(w, r, route)
}

// handleNotFound 处理404错误
func (s *server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	response := map[string]interface{}{
		"success": false,
		"message": "Not Found",
		"path":    r.URL.Path,
	}

	json.NewEncoder(w).Encode(response)
}

// executeMiddlewareChain 执行中间件链
func (s *server) executeMiddlewareChain(w http.ResponseWriter, r *http.Request, route routing.Route) {
	// 创建处理器
	handler := &routeHandler{
		route: route,
	}

	// 创建中间件管道
	pipeline := NewPipeline()

	// 添加全局中间件
	for _, middlewareName := range s.middleware {
		if middleware := s.getMiddleware(middlewareName); middleware != nil {
			pipeline.Use(middleware)
		}
	}

	// 添加路由中间件
	for _, middlewareName := range route.Middleware {
		if name, ok := middlewareName.(string); ok {
			if middleware := s.getMiddleware(name); middleware != nil {
				pipeline.Use(middleware)
			}
		}
	}

	// 执行中间件链
	response := pipeline.Process(NewRequest(r), handler)

	// 发送响应
	response.Send(w)
}

// routeHandler 路由处理器
type routeHandler struct {
	route routing.Route
}

// Handle 实现 Handler 接口
func (rh *routeHandler) Handle(request Request) Response {
	// 调用路由处理器
	if rh.route.Handler != nil {
		// 这里需要调用实际的处理器
		// 暂时返回简单的响应
		return NewJsonResponse(http.StatusOK, map[string]interface{}{
			"message": "Handler executed",
			"path":    request.Path(),
		})
	}
	return NewJsonResponse(http.StatusOK, map[string]interface{}{
		"message": "No handler found",
	})
}

// getMiddleware 获取中间件
func (s *server) getMiddleware(name string) Middleware {
	// 从容器中获取中间件
	if s.container != nil {
		if middleware := s.container.Make("middleware." + name); middleware != nil {
			if mw, ok := middleware.(Middleware); ok {
				return mw
			}
		}
	}

	// 返回默认中间件
	return s.getDefaultMiddleware(name)
}

// getDefaultMiddleware 获取默认中间件
func (s *server) getDefaultMiddleware(name string) Middleware {
	switch name {
	case "cors":
		return &CORSMiddleware{}
	case "auth":
		return &AuthMiddleware{}
	case "logging":
		return &LoggingMiddleware{}
	default:
		return nil
	}
}

// responseWriter 响应写入器
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 写入状态码
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write 写入数据
func (rw *responseWriter) Write(data []byte) (int, error) {
	return rw.ResponseWriter.Write(data)
}

// ServerBuilder 服务器构建器
type ServerBuilder struct {
	config     config.Config
	container  container.Container
	middleware []string
	static     map[string]string
}

// NewServerBuilder 创建服务器构建器
func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		middleware: make([]string, 0),
		static:     make(map[string]string),
	}
}

// Config 设置配置
func (sb *ServerBuilder) Config(config config.Config) *ServerBuilder {
	sb.config = config
	return sb
}

// Container 设置容器
func (sb *ServerBuilder) Container(container container.Container) *ServerBuilder {
	sb.container = container
	return sb
}

// Middleware 添加中间件
func (sb *ServerBuilder) Middleware(middleware ...string) *ServerBuilder {
	sb.middleware = append(sb.middleware, middleware...)
	return sb
}

// Static 设置静态文件
func (sb *ServerBuilder) Static(path, dir string) *ServerBuilder {
	sb.static[path] = dir
	return sb
}

// Build 构建服务器
func (sb *ServerBuilder) Build() Server {
	server := NewServer(sb.config, sb.container)

	// 添加中间件
	for _, middleware := range sb.middleware {
		server.Use(middleware)
	}

	// 添加静态文件
	for path, dir := range sb.static {
		server.Static(path, dir)
	}

	return server
}
