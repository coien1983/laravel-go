package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Version 版本信息
type Version struct {
	Version     string    `json:"version"`
	Status      string    `json:"status"` // stable, deprecated, beta
	DeprecatedAt *time.Time `json:"deprecated_at,omitempty"`
	SunsetAt    *time.Time `json:"sunset_at,omitempty"`
	Message     string    `json:"message,omitempty"`
}

// VersionManager 版本管理器
type VersionManager struct {
	versions map[string]*Version
	defaultVersion string
	deprecationWarnings map[string]bool
}

// NewVersionManager 创建版本管理器
func NewVersionManager() *VersionManager {
	return &VersionManager{
		versions: make(map[string]*Version),
		defaultVersion: "v1",
		deprecationWarnings: make(map[string]bool),
	}
}

// RegisterVersion 注册版本
func (vm *VersionManager) RegisterVersion(version string, status string) *Version {
	v := &Version{
		Version: version,
		Status:  status,
	}
	vm.versions[version] = v
	return v
}

// SetDefaultVersion 设置默认版本
func (vm *VersionManager) SetDefaultVersion(version string) {
	vm.defaultVersion = version
}

// GetDefaultVersion 获取默认版本
func (vm *VersionManager) GetDefaultVersion() string {
	return vm.defaultVersion
}

// GetVersion 获取版本信息
func (vm *VersionManager) GetVersion(version string) (*Version, bool) {
	v, exists := vm.versions[version]
	return v, exists
}

// IsVersionDeprecated 检查版本是否已弃用
func (vm *VersionManager) IsVersionDeprecated(version string) bool {
	if v, exists := vm.versions[version]; exists {
		return v.Status == "deprecated"
	}
	return false
}

// IsVersionSunset 检查版本是否已停止支持
func (vm *VersionManager) IsVersionSunset(version string) bool {
	if v, exists := vm.versions[version]; exists {
		if v.SunsetAt != nil {
			return time.Now().After(*v.SunsetAt)
		}
	}
	return false
}

// DeprecateVersion 弃用版本
func (vm *VersionManager) DeprecateVersion(version string, message string, sunsetAt time.Time) error {
	if v, exists := vm.versions[version]; exists {
		v.Status = "deprecated"
		v.Message = message
		v.DeprecatedAt = &time.Time{}
		*v.DeprecatedAt = time.Now()
		v.SunsetAt = &sunsetAt
		return nil
	}
	return fmt.Errorf("version %s not found", version)
}

// GetSupportedVersions 获取支持的版本列表
func (vm *VersionManager) GetSupportedVersions() []*Version {
	var supported []*Version
	for _, version := range vm.versions {
		if version.Status != "sunset" {
			supported = append(supported, version)
		}
	}
	return supported
}

// VersionMiddleware 版本中间件
type VersionMiddleware struct {
	versionManager *VersionManager
	headerName     string
	paramName      string
	required       bool
}

// NewVersionMiddleware 创建版本中间件
func NewVersionMiddleware(versionManager *VersionManager) *VersionMiddleware {
	return &VersionMiddleware{
		versionManager: versionManager,
		headerName:     "Accept-Version",
		paramName:      "version",
		required:       false,
	}
}

// SetHeaderName 设置版本头名称
func (vm *VersionMiddleware) SetHeaderName(name string) *VersionMiddleware {
	vm.headerName = name
	return vm
}

// SetParamName 设置版本参数名称
func (vm *VersionMiddleware) SetParamName(name string) *VersionMiddleware {
	vm.paramName = name
	return vm
}

// SetRequired 设置是否必需版本
func (vm *VersionMiddleware) SetRequired(required bool) *VersionMiddleware {
	vm.required = required
	return vm
}

// Handle 处理版本检查
func (vm *VersionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求中获取版本
		version := vm.extractVersion(r)
		
		// 如果没有指定版本，使用默认版本
		if version == "" {
			version = vm.versionManager.GetDefaultVersion()
		}
		
		// 检查版本是否存在
		if _, exists := vm.versionManager.GetVersion(version); !exists {
			if vm.required {
				http.Error(w, fmt.Sprintf("Unsupported API version: %s", version), http.StatusBadRequest)
				return
			}
			version = vm.versionManager.GetDefaultVersion()
		}
		
		// 检查版本是否已弃用
		if vm.versionManager.IsVersionDeprecated(version) {
			vm.addDeprecationWarning(w, version)
		}
		
		// 检查版本是否已停止支持
		if vm.versionManager.IsVersionSunset(version) {
			http.Error(w, fmt.Sprintf("API version %s is no longer supported", version), http.StatusGone)
			return
		}
		
		// 将版本信息添加到请求上下文
		ctx := r.Context()
		ctx = contextWithVersion(ctx, version)
		r = r.WithContext(ctx)
		
		next(w, r)
	}
}

// extractVersion 从请求中提取版本
func (vm *VersionMiddleware) extractVersion(r *http.Request) string {
	// 从头部获取版本
	if version := r.Header.Get(vm.headerName); version != "" {
		return version
	}
	
	// 从查询参数获取版本
	if version := r.URL.Query().Get(vm.paramName); version != "" {
		return version
	}
	
	// 从路径获取版本 (例如 /api/v1/users)
	path := r.URL.Path
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, "v") && i > 0 {
			return part
		}
	}
	
	return ""
}

// addDeprecationWarning 添加弃用警告
func (vm *VersionMiddleware) addDeprecationWarning(w http.ResponseWriter, version string) {
	if v, exists := vm.versionManager.GetVersion(version); exists {
		warning := fmt.Sprintf("299 - \"API version %s is deprecated", version)
		if v.Message != "" {
			warning += fmt.Sprintf(": %s", v.Message)
		}
		if v.SunsetAt != nil {
			warning += fmt.Sprintf(". Sunset date: %s", v.SunsetAt.Format("2006-01-02"))
		}
		warning += "\""
		w.Header().Add("Warning", warning)
	}
}

// VersionRouter 版本路由器
type VersionRouter struct {
	versionManager *VersionManager
	routes         map[string]map[string]http.HandlerFunc
	middleware     map[string][]func(http.HandlerFunc) http.HandlerFunc
}

// NewVersionRouter 创建版本路由器
func NewVersionRouter(versionManager *VersionManager) *VersionRouter {
	return &VersionRouter{
		versionManager: versionManager,
		routes:         make(map[string]map[string]http.HandlerFunc),
		middleware:     make(map[string][]func(http.HandlerFunc) http.HandlerFunc),
	}
}

// Route 添加路由
func (vr *VersionRouter) Route(version, method, path string, handler http.HandlerFunc) {
	if vr.routes[version] == nil {
		vr.routes[version] = make(map[string]http.HandlerFunc)
	}
	key := fmt.Sprintf("%s:%s", method, path)
	vr.routes[version][key] = handler
}

// GET 添加 GET 路由
func (vr *VersionRouter) GET(version, path string, handler http.HandlerFunc) {
	vr.Route(version, "GET", path, handler)
}

// POST 添加 POST 路由
func (vr *VersionRouter) POST(version, path string, handler http.HandlerFunc) {
	vr.Route(version, "POST", path, handler)
}

// PUT 添加 PUT 路由
func (vr *VersionRouter) PUT(version, path string, handler http.HandlerFunc) {
	vr.Route(version, "PUT", path, handler)
}

// DELETE 添加 DELETE 路由
func (vr *VersionRouter) DELETE(version, path string, handler http.HandlerFunc) {
	vr.Route(version, "DELETE", path, handler)
}

// Use 添加中间件
func (vr *VersionRouter) Use(version string, middleware ...func(http.HandlerFunc) http.HandlerFunc) {
	vr.middleware[version] = append(vr.middleware[version], middleware...)
}

// ServeHTTP 处理 HTTP 请求
func (vr *VersionRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 从请求中提取版本
	version := vr.extractVersion(r)
	
	// 检查版本是否存在
	if _, exists := vr.versionManager.GetVersion(version); !exists {
		version = vr.versionManager.GetDefaultVersion()
	}
	
	// 查找路由 - 移除版本前缀
	pathWithoutVersion := r.URL.Path
	parts := strings.Split(pathWithoutVersion, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, "v") && i > 0 {
			// 重建路径，移除版本部分
			pathWithoutVersion = "/" + strings.Join(parts[i+1:], "/")
			break
		}
	}
	
	key := fmt.Sprintf("%s:%s", r.Method, pathWithoutVersion)
	if handler, exists := vr.routes[version][key]; exists {
		// 应用中间件
		finalHandler := handler
		for i := len(vr.middleware[version]) - 1; i >= 0; i-- {
			finalHandler = vr.middleware[version][i](finalHandler)
		}
		
		// 将版本信息添加到请求上下文
		ctx := r.Context()
		ctx = contextWithVersion(ctx, version)
		r = r.WithContext(ctx)
		
		// 检查版本状态
		if vr.versionManager.IsVersionSunset(version) {
			http.Error(w, fmt.Sprintf("API version %s is no longer supported", version), http.StatusGone)
			return
		}
		
		if vr.versionManager.IsVersionDeprecated(version) {
			vr.addDeprecationWarning(w, version)
		}
		
		finalHandler(w, r)
		return
	}
	
	// 路由未找到
	http.NotFound(w, r)
}

// extractVersion 从请求中提取版本
func (vr *VersionRouter) extractVersion(r *http.Request) string {
	// 从路径获取版本 (例如 /api/v1/users)
	path := r.URL.Path
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, "v") && i > 0 {
			return part
		}
	}
	
	// 从头部获取版本
	if version := r.Header.Get("Accept-Version"); version != "" {
		return version
	}
	
	// 从查询参数获取版本
	if version := r.URL.Query().Get("version"); version != "" {
		return version
	}
	
	return vr.versionManager.GetDefaultVersion()
}

// addDeprecationWarning 添加弃用警告
func (vr *VersionRouter) addDeprecationWarning(w http.ResponseWriter, version string) {
	if v, exists := vr.versionManager.GetVersion(version); exists {
		warning := fmt.Sprintf("299 - \"API version %s is deprecated", version)
		if v.Message != "" {
			warning += fmt.Sprintf(": %s", v.Message)
		}
		if v.SunsetAt != nil {
			warning += fmt.Sprintf(". Sunset date: %s", v.SunsetAt.Format("2006-01-02"))
		}
		warning += "\""
		w.Header().Add("Warning", warning)
	}
}

// 上下文相关函数（需要实现 context 包）

// contextKey 上下文键类型
type contextKey string

const versionContextKey contextKey = "api_version"

// contextWithVersion 创建包含版本信息的上下文
func contextWithVersion(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, versionContextKey, version)
}

// VersionFromContext 从上下文获取版本
func VersionFromContext(ctx context.Context) string {
	if version, ok := ctx.Value(versionContextKey).(string); ok {
		return version
	}
	return ""
}

// 便捷函数

// NewAPIVersion 创建新的 API 版本
func NewAPIVersion(version string, status string) *Version {
	return &Version{
		Version: version,
		Status:  status,
	}
}

// NewVersionedAPI 创建版本化的 API
func NewVersionedAPI() (*VersionManager, *VersionRouter) {
	vm := NewVersionManager()
	vr := NewVersionRouter(vm)
	return vm, vr
} 