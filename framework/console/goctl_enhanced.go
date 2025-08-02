package console

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// GoZeroGenerator go-zero代码生成器
type GoZeroGenerator struct {
	output Output
}

// NewGoZeroGenerator 创建新的go-zero生成器
func NewGoZeroGenerator(output Output) *GoZeroGenerator {
	return &GoZeroGenerator{
		output: output,
	}
}

// GenerateFromProto 从proto文件生成完整的go-zero服务
func (g *GoZeroGenerator) GenerateFromProto(protoFile, outputDir string) error {
	g.output.Info("🚀 开始从proto文件生成go-zero服务...")

	// 检查proto文件
	if _, err := os.Stat(protoFile); os.IsNotExist(err) {
		return fmt.Errorf("proto file not found: %s", protoFile)
	}

	// 解析proto文件
	protoInfo, err := g.parseProtoFile(protoFile)
	if err != nil {
		return fmt.Errorf("failed to parse proto file: %w", err)
	}

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 生成go-zero项目结构
	if err := g.generateGoZeroStructure(protoInfo, outputDir); err != nil {
		return fmt.Errorf("failed to generate go-zero structure: %w", err)
	}

	g.output.Success("✅ go-zero服务生成完成")
	return nil
}

// GenerateFromApi 从.api文件生成完整的go-zero API服务
func (g *GoZeroGenerator) GenerateFromApi(apiFile, outputDir string) error {
	g.output.Info("🚀 开始从.api文件生成go-zero API服务...")

	// 检查api文件
	if _, err := os.Stat(apiFile); os.IsNotExist(err) {
		return fmt.Errorf("api file not found: %s", apiFile)
	}

	// 解析api文件
	apiInfo, err := g.parseApiFile(apiFile)
	if err != nil {
		return fmt.Errorf("failed to parse api file: %w", err)
	}

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 生成go-zero API项目结构
	if err := g.generateGoZeroApiStructure(apiInfo, outputDir); err != nil {
		return fmt.Errorf("failed to generate go-zero API structure: %w", err)
	}

	g.output.Success("✅ go-zero API服务生成完成")
	return nil
}

// GenerateMicroservice 生成完整的微服务
func (g *GoZeroGenerator) GenerateMicroservice(name, protoFile, apiFile, outputDir string) error {
	g.output.Info("🚀 开始生成微服务...")

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 生成RPC服务（如果有proto文件）
	if protoFile != "" {
		if err := g.GenerateFromProto(protoFile, filepath.Join(outputDir, "rpc")); err != nil {
			return fmt.Errorf("failed to generate RPC service: %w", err)
		}
	}

	// 生成API服务（如果有api文件）
	if apiFile != "" {
		if err := g.GenerateFromApi(apiFile, filepath.Join(outputDir, "api")); err != nil {
			return fmt.Errorf("failed to generate API service: %w", err)
		}
	}

	// 生成网关配置
	if err := g.generateGatewayConfig(name, outputDir); err != nil {
		return fmt.Errorf("failed to generate gateway config: %w", err)
	}

	g.output.Success("✅ 微服务生成完成")
	return nil
}

// GenerateLogic 生成logic层
func (g *GoZeroGenerator) GenerateLogic(methodName, serviceName, outputDir string) error {
	g.output.Info(fmt.Sprintf("🚀 开始生成logic层: %s", methodName))

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 生成logic文件
	logicFile := filepath.Join(outputDir, strings.ToLower(methodName)+"logic.go")
	logicTemplate := `package logic

import (
	"context"
	"{{ .ProjectName }}/internal/svc"
	"{{ .ProjectName }}/internal/types"
)

type {{ .MethodName }}Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func New{{ .MethodName }}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *{{ .MethodName }}Logic {
	return &{{ .MethodName }}Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *{{ .MethodName }}Logic) {{ .MethodName }}(req *types.{{ .MethodName }}Req) (resp *types.{{ .MethodName }}Resp, err error) {
	// TODO: 在这里实现你的业务逻辑
	return &types.{{ .MethodName }}Resp{}, nil
}
`

	data := map[string]interface{}{
		"ProjectName": g.getProjectName(),
		"MethodName":  g.toPascalCase(methodName),
	}

	if err := g.writeTemplateToFile(logicFile, logicTemplate, data); err != nil {
		return fmt.Errorf("failed to write logic file: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ Logic层生成完成: %s", logicFile))
	return nil
}

// GenerateHandler 生成handler层
func (g *GoZeroGenerator) GenerateHandler(endpointName, method, path, outputDir string) error {
	g.output.Info(fmt.Sprintf("🚀 开始生成handler层: %s", endpointName))

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 生成handler文件
	handlerFile := filepath.Join(outputDir, strings.ToLower(endpointName)+"handler.go")
	handlerTemplate := `package handler

import (
	"net/http"
	"{{ .ProjectName }}/internal/logic"
	"{{ .ProjectName }}/internal/svc"
	"{{ .ProjectName }}/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func {{ .EndpointName }}Handler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.{{ .EndpointName }}Req
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.New{{ .EndpointName }}Logic(r.Context(), serverCtx)
		resp, err := l.{{ .EndpointName }}(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
`

	data := map[string]interface{}{
		"ProjectName":  g.getProjectName(),
		"EndpointName": g.toPascalCase(endpointName),
	}

	if err := g.writeTemplateToFile(handlerFile, handlerTemplate, data); err != nil {
		return fmt.Errorf("failed to write handler file: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ Handler层生成完成: %s", handlerFile))
	return nil
}

// generateGatewayConfig 生成网关配置
func (g *GoZeroGenerator) generateGatewayConfig(name, outputDir string) error {
	// 生成网关配置文件
	gatewayFile := filepath.Join(outputDir, "gateway.yaml")
	gatewayTemplate := `Name: {{ .Name }}-gateway
Host: 0.0.0.0
Port: 8888
Mode: dev

Routes:
  - Path: /api
    Target: http://localhost:8080
    StripPrefix: true
  - Path: /rpc
    Target: grpc://localhost:9090
    StripPrefix: true
`

	data := map[string]interface{}{
		"Name": name,
	}

	if err := g.writeTemplateToFile(gatewayFile, gatewayTemplate, data); err != nil {
		return fmt.Errorf("failed to write gateway config: %w", err)
	}

	return nil
}

// ProtoInfo proto文件信息
type ProtoInfo struct {
	Package  string
	Services []ProtoServiceInfo
	Messages []ProtoMessageInfo
	Imports  []string
}

// ProtoServiceInfo 服务信息
type ProtoServiceInfo struct {
	Name    string
	Methods []ProtoMethodInfo
}

// ProtoMethodInfo 方法信息
type ProtoMethodInfo struct {
	Name   string
	Input  string
	Output string
}

// ProtoMessageInfo 消息信息
type ProtoMessageInfo struct {
	Name   string
	Fields []ProtoFieldInfo
}

// ProtoFieldInfo 字段信息
type ProtoFieldInfo struct {
	Number int
	Name   string
	Type   string
	Rule   string // repeated, optional, required
}

// parseProtoFile 解析proto文件
func (g *GoZeroGenerator) parseProtoFile(protoFile string) (*ProtoInfo, error) {
	file, err := os.Open(protoFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	protoInfo := &ProtoInfo{
		Services: []ProtoServiceInfo{},
		Messages: []ProtoMessageInfo{},
		Imports:  []string{},
	}

	scanner := bufio.NewScanner(file)
	var currentService *ProtoServiceInfo
	var currentMessage *ProtoMessageInfo
	var inService, inMessage bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// 解析package
		if strings.HasPrefix(line, "package ") {
			protoInfo.Package = strings.TrimSpace(strings.TrimPrefix(line, "package"))
			protoInfo.Package = strings.TrimSuffix(protoInfo.Package, ";")
			continue
		}

		// 解析import
		if strings.HasPrefix(line, "import ") {
			importStr := strings.TrimSpace(strings.TrimPrefix(line, "import"))
			importStr = strings.TrimSuffix(importStr, ";")
			importStr = strings.Trim(importStr, "\"")
			protoInfo.Imports = append(protoInfo.Imports, importStr)
			continue
		}

		// 解析service
		if strings.HasPrefix(line, "service ") {
			if currentService != nil {
				protoInfo.Services = append(protoInfo.Services, *currentService)
			}
			serviceName := strings.TrimSpace(strings.TrimPrefix(line, "service"))
			serviceName = strings.TrimSpace(strings.TrimSuffix(serviceName, "{"))
			currentService = &ProtoServiceInfo{Name: serviceName, Methods: []ProtoMethodInfo{}}
			inService = true
			inMessage = false
			continue
		}

		// 解析message
		if strings.HasPrefix(line, "message ") {
			if currentMessage != nil {
				protoInfo.Messages = append(protoInfo.Messages, *currentMessage)
			}
			messageName := strings.TrimSpace(strings.TrimPrefix(line, "message"))
			messageName = strings.TrimSpace(strings.TrimSuffix(messageName, "{"))
			currentMessage = &ProtoMessageInfo{Name: messageName, Fields: []ProtoFieldInfo{}}
			inMessage = true
			inService = false
			continue
		}

		// 解析service方法
		if inService && currentService != nil && strings.Contains(line, "(") && strings.Contains(line, ")") {
			methodRegex := regexp.MustCompile(`rpc\s+(\w+)\s*\(\s*(\w+)\s*\)\s*returns\s*\(\s*(\w+)\s*\)`)
			matches := methodRegex.FindStringSubmatch(line)
			if len(matches) >= 4 {
				currentService.Methods = append(currentService.Methods, ProtoMethodInfo{
					Name:   matches[1],
					Input:  matches[2],
					Output: matches[3],
				})
			}
		}

		// 解析message字段
		if inMessage && currentMessage != nil && strings.Contains(line, " ") {
			fieldRegex := regexp.MustCompile(`(\w+)\s+(\w+)\s+(\w+)\s*=\s*(\d+)`)
			matches := fieldRegex.FindStringSubmatch(line)
			if len(matches) >= 5 {
				fieldNumber := 0
				fmt.Sscanf(matches[4], "%d", &fieldNumber)
				currentMessage.Fields = append(currentMessage.Fields, ProtoFieldInfo{
					Number: fieldNumber,
					Type:   matches[1],
					Name:   matches[2],
					Rule:   matches[3],
				})
			}
		}

		// 处理结束括号
		if line == "}" {
			if inService && currentService != nil {
				protoInfo.Services = append(protoInfo.Services, *currentService)
				currentService = nil
				inService = false
			}
			if inMessage && currentMessage != nil {
				protoInfo.Messages = append(protoInfo.Messages, *currentMessage)
				currentMessage = nil
				inMessage = false
			}
		}
	}

	// 添加最后一个service和message
	if currentService != nil {
		protoInfo.Services = append(protoInfo.Services, *currentService)
	}
	if currentMessage != nil {
		protoInfo.Messages = append(protoInfo.Messages, *currentMessage)
	}

	return protoInfo, nil
}

// ApiFileInfo .api文件信息
type ApiFileInfo struct {
	Info     map[string]string
	Types    []ApiTypeInfo
	Services []ApiServiceInfo
}

// ApiTypeInfo API类型信息
type ApiTypeInfo struct {
	Name   string
	Fields []ApiFieldInfo
}

// ApiFieldInfo API字段信息
type ApiFieldInfo struct {
	Name string
	Type string
	Tag  string
}

// ApiServiceInfo API服务信息
type ApiServiceInfo struct {
	Name    string
	Methods []ApiMethodInfo
}

// ApiMethodInfo API方法信息
type ApiMethodInfo struct {
	Name   string
	Path   string
	Method string
	Req    string
	Resp   string
}

// parseApiFile 解析.api文件
func (g *GoZeroGenerator) parseApiFile(apiFile string) (*ApiFileInfo, error) {
	file, err := os.Open(apiFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	apiInfo := &ApiFileInfo{
		Info:     make(map[string]string),
		Types:    []ApiTypeInfo{},
		Services: []ApiServiceInfo{},
	}

	scanner := bufio.NewScanner(file)
	var currentType *ApiTypeInfo
	var currentService *ApiServiceInfo

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// 解析info块
		if strings.HasPrefix(line, "info") {
			// 解析info信息
			continue
		}

		// 解析type块
		if strings.HasPrefix(line, "type") {
			if currentType != nil {
				apiInfo.Types = append(apiInfo.Types, *currentType)
			}
			typeName := strings.TrimSpace(strings.TrimPrefix(line, "type"))
			typeName = strings.TrimSpace(strings.TrimSuffix(typeName, "{"))
			currentType = &ApiTypeInfo{Name: typeName, Fields: []ApiFieldInfo{}}
			continue
		}

		// 解析service块
		if strings.HasPrefix(line, "service") {
			if currentService != nil {
				apiInfo.Services = append(apiInfo.Services, *currentService)
			}
			serviceName := strings.TrimSpace(strings.TrimPrefix(line, "service"))
			serviceName = strings.TrimSpace(strings.TrimSuffix(serviceName, "{"))
			currentService = &ApiServiceInfo{Name: serviceName, Methods: []ApiMethodInfo{}}
			continue
		}

		// 解析字段
		if currentType != nil && strings.Contains(line, " ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				fieldName := parts[0]
				fieldType := parts[1]
				tag := ""
				if len(parts) > 2 {
					tag = strings.Join(parts[2:], " ")
				}
				currentType.Fields = append(currentType.Fields, ApiFieldInfo{
					Name: fieldName,
					Type: fieldType,
					Tag:  tag,
				})
			}
		}

		// 解析HTTP方法
		if currentService != nil && (strings.Contains(line, "get") || strings.Contains(line, "post") ||
			strings.Contains(line, "put") || strings.Contains(line, "delete")) {
			// 简单的HTTP方法解析
			methodRegex := regexp.MustCompile(`(get|post|put|delete)\s+([^\s]+)\s*\(([^)]*)\)\s*returns\s*\(([^)]*)\)`)
			matches := methodRegex.FindStringSubmatch(line)
			if len(matches) >= 5 {
				currentService.Methods = append(currentService.Methods, ApiMethodInfo{
					Method: strings.ToUpper(matches[1]),
					Path:   matches[2],
					Req:    strings.TrimSpace(matches[3]),
					Resp:   strings.TrimSpace(matches[4]),
				})
			}
		}

		// 处理结束括号
		if line == "}" {
			if currentType != nil {
				apiInfo.Types = append(apiInfo.Types, *currentType)
				currentType = nil
			}
			if currentService != nil {
				apiInfo.Services = append(apiInfo.Services, *currentService)
				currentService = nil
			}
		}
	}

	// 添加最后一个type和service
	if currentType != nil {
		apiInfo.Types = append(apiInfo.Types, *currentType)
	}
	if currentService != nil {
		apiInfo.Services = append(apiInfo.Services, *currentService)
	}

	return apiInfo, nil
}

// generateGoZeroStructure 生成go-zero项目结构
func (g *GoZeroGenerator) generateGoZeroStructure(protoInfo *ProtoInfo, outputDir string) error {
	// 创建目录结构
	dirs := []string{
		"internal/logic",
		"internal/svc",
		"internal/config",
		"internal/server",
		"etc",
	}

	for _, dir := range dirs {
		fullDir := filepath.Join(outputDir, dir)
		if err := os.MkdirAll(fullDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullDir, err)
		}
	}

	// 生成main.go
	mainFile := filepath.Join(outputDir, "main.go")
	mainTemplate := `package main

import (
	"flag"
	"fmt"
	"log"

	"{{ .ProjectName }}/internal/config"
	"{{ .ProjectName }}/internal/server"
	"{{ .ProjectName }}/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/{{ .ServiceName }}.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	srv := server.New{{ .ServiceName }}Server(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		{{ range .Services }}
		{{ .Name }}.Register{{ .Name }}Server(grpcServer, srv)
		{{ end }}

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
`

	data := map[string]interface{}{
		"ProjectName": g.getProjectName(),
		"ServiceName": g.toPascalCase(protoInfo.Services[0].Name),
		"Services":    protoInfo.Services,
	}

	if err := g.writeTemplateToFile(mainFile, mainTemplate, data); err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	// 生成配置文件
	configFile := filepath.Join(outputDir, "etc", protoInfo.Services[0].Name+".yaml")
	configTemplate := `Name: {{ .ServiceName }}
Host: 0.0.0.0
Port: 8080
Mode: dev

RpcServerConf:
  Endpoints:
    - 0.0.0.0:8080
  Timeout: 30000
`

	if err := g.writeTemplateToFile(configFile, configTemplate, data); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// 生成配置结构
	configStructFile := filepath.Join(outputDir, "internal/config/config.go")
	configStructTemplate := `package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mode string
}
`

	if err := g.writeTemplateToFile(configStructFile, configStructTemplate, data); err != nil {
		return fmt.Errorf("failed to write config struct file: %w", err)
	}

	// 生成服务上下文
	svcFile := filepath.Join(outputDir, "internal/svc/servicecontext.go")
	svcTemplate := `package svc

import (
	"{{ .ProjectName }}/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
`

	if err := g.writeTemplateToFile(svcFile, svcTemplate, data); err != nil {
		return fmt.Errorf("failed to write service context file: %w", err)
	}

	// 生成服务器实现
	for _, service := range protoInfo.Services {
		serverFile := filepath.Join(outputDir, "internal/server", service.Name+"server.go")
		serverTemplate := `package server

import (
	"context"
	"{{ .ProjectName }}/internal/logic"
	"{{ .ProjectName }}/internal/svc"
	"{{ .ProjectName }}/types"
)

type {{ .ServiceName }}Server struct {
	svcCtx *svc.ServiceContext
	types.Unimplemented{{ .ServiceName }}Server
}

func New{{ .ServiceName }}Server(svcCtx *svc.ServiceContext) *{{ .ServiceName }}Server {
	return &{{ .ServiceName }}Server{
		svcCtx: svcCtx,
	}
}

{{ range .Methods }}
func (s *{{ $.ServiceName }}Server) {{ .Name }}(ctx context.Context, req *types.{{ .Input }}) (*types.{{ .Output }}, error) {
	l := logic.New{{ .Name }}Logic(ctx, s.svcCtx)
	return l.{{ .Name }}(req)
}
{{ end }}
`

		serviceData := map[string]interface{}{
			"ProjectName": g.getProjectName(),
			"ServiceName": g.toPascalCase(service.Name),
			"Methods":     service.Methods,
		}

		if err := g.writeTemplateToFile(serverFile, serverTemplate, serviceData); err != nil {
			return fmt.Errorf("failed to write server file: %w", err)
		}
	}

	g.output.Success(fmt.Sprintf("✅ go-zero项目结构生成完成: %s", outputDir))
	return nil
}

// generateGoZeroApiStructure 生成go-zero API项目结构
func (g *GoZeroGenerator) generateGoZeroApiStructure(apiInfo *ApiFileInfo, outputDir string) error {
	// 创建目录结构
	dirs := []string{
		"internal/handler",
		"internal/logic",
		"internal/svc",
		"internal/config",
		"internal/types",
		"etc",
	}

	for _, dir := range dirs {
		fullDir := filepath.Join(outputDir, dir)
		if err := os.MkdirAll(fullDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullDir, err)
		}
	}

	// 生成main.go
	mainFile := filepath.Join(outputDir, "main.go")
	mainTemplate := `package main

import (
	"flag"
	"fmt"

	"{{ .ProjectName }}/internal/config"
	"{{ .ProjectName }}/internal/handler"
	"{{ .ProjectName }}/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
`

	data := map[string]interface{}{
		"ProjectName": g.getProjectName(),
	}

	if err := g.writeTemplateToFile(mainFile, mainTemplate, data); err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	// 生成配置文件
	configFile := filepath.Join(outputDir, "etc", "api.yaml")
	configTemplate := `Name: api
Host: 0.0.0.0
Port: 8888
Mode: dev
`

	if err := g.writeTemplateToFile(configFile, configTemplate, data); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// 生成配置结构
	configStructFile := filepath.Join(outputDir, "internal/config/config.go")
	configStructTemplate := `package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
}
`

	if err := g.writeTemplateToFile(configStructFile, configStructTemplate, data); err != nil {
		return fmt.Errorf("failed to write config struct file: %w", err)
	}

	// 生成服务上下文
	svcFile := filepath.Join(outputDir, "internal/svc/servicecontext.go")
	svcTemplate := `package svc

import (
	"{{ .ProjectName }}/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
`

	if err := g.writeTemplateToFile(svcFile, svcTemplate, data); err != nil {
		return fmt.Errorf("failed to write service context file: %w", err)
	}

	// 生成类型定义
	if len(apiInfo.Types) > 0 {
		typesFile := filepath.Join(outputDir, "internal/types/types.go")
		typesTemplate := `package types

{{ range .Types }}
type {{ .Name }} struct {
{{ range .Fields }}	{{ .Name }} {{ .Type }} {{ .Tag }}
{{ end }}}
{{ end }}
`

		if err := g.writeTemplateToFile(typesFile, typesTemplate, apiInfo); err != nil {
			return fmt.Errorf("failed to write types file: %w", err)
		}
	}

	// 生成处理器
	handlerFile := filepath.Join(outputDir, "internal/handler/handlers.go")
	handlerTemplate := `package handler

import (
	"net/http"

	"{{ .ProjectName }}/internal/logic"
	"{{ .ProjectName }}/internal/svc"
	"{{ .ProjectName }}/internal/types"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
{{ range .Services }}
	// {{ .Name }} 服务路由
	server.AddRoutes(
		[]rest.Route{
{{ range .Methods }}			{
				Method:  http.Method{{ .Method }},
				Path:    "{{ .Path }}",
				Handler: {{ .Name }}Handler(serverCtx),
			},
{{ end }}		},
	)
{{ end }}
}

{{ range .Services }}
{{ range .Methods }}
func {{ .Name }}Handler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.{{ .Req }}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.New{{ .Name }}Logic(r.Context(), serverCtx)
		resp, err := l.{{ .Name }}(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
{{ end }}
{{ end }}
`

	handlerData := map[string]interface{}{
		"ProjectName": g.getProjectName(),
		"Services":    apiInfo.Services,
	}

	if err := g.writeTemplateToFile(handlerFile, handlerTemplate, handlerData); err != nil {
		return fmt.Errorf("failed to write handler file: %w", err)
	}

	g.output.Success(fmt.Sprintf("✅ go-zero API项目结构生成完成: %s", outputDir))
	return nil
}

// writeTemplateToFile 将模板写入文件
func (g *GoZeroGenerator) writeTemplateToFile(filePath, templateStr string, data interface{}) error {
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// getProjectName 获取项目名称
func (g *GoZeroGenerator) getProjectName() string {
	// 这里可以从 go.mod 文件读取项目名称
	// 暂时返回默认值
	return "laravel-go"
}

// toPascalCase 转换为PascalCase
func (g *GoZeroGenerator) toPascalCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
