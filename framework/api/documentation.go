package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// OpenAPISpec OpenAPI 规范
type OpenAPISpec struct {
	OpenAPI    string                 `json:"openapi"`
	Info       *Info                  `json:"info"`
	Servers    []*Server              `json:"servers,omitempty"`
	Paths      map[string]*PathItem   `json:"paths"`
	Components *Components            `json:"components,omitempty"`
	Tags       []*Tag                 `json:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// Info API 信息
type Info struct {
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	Version        string   `json:"version"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

// Contact 联系信息
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License 许可证信息
type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Server 服务器信息
type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty"`
}

// ServerVariable 服务器变量
type ServerVariable struct {
	Default     string   `json:"default"`
	Enum        []string `json:"enum,omitempty"`
	Description string   `json:"description,omitempty"`
}

// PathItem 路径项
type PathItem struct {
	Ref         string     `json:"$ref,omitempty"`
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	GET         *Operation `json:"get,omitempty"`
	PUT         *Operation `json:"put,omitempty"`
	POST        *Operation `json:"post,omitempty"`
	DELETE      *Operation `json:"delete,omitempty"`
	OPTIONS     *Operation `json:"options,omitempty"`
	HEAD        *Operation `json:"head,omitempty"`
	PATCH       *Operation `json:"patch,omitempty"`
	TRACE       *Operation `json:"trace,omitempty"`
	Servers     []*Server  `json:"servers,omitempty"`
	Parameters  []*Parameter `json:"parameters,omitempty"`
}

// Operation 操作
type Operation struct {
	Tags         []string               `json:"tags,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	Description  string                 `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	OperationID  string                 `json:"operationId,omitempty"`
	Parameters   []*Parameter           `json:"parameters,omitempty"`
	RequestBody  *RequestBody           `json:"requestBody,omitempty"`
	Responses    map[string]*Response   `json:"responses"`
	Callbacks    map[string]*Callback   `json:"callbacks,omitempty"`
	Deprecated   bool                   `json:"deprecated,omitempty"`
	Security     []*SecurityRequirement `json:"security,omitempty"`
	Servers      []*Server              `json:"servers,omitempty"`
}

// Parameter 参数
type Parameter struct {
	Name            string      `json:"name"`
	In              string      `json:"in"`
	Description     string      `json:"description,omitempty"`
	Required        bool        `json:"required,omitempty"`
	Deprecated      bool        `json:"deprecated,omitempty"`
	AllowEmptyValue bool        `json:"allowEmptyValue,omitempty"`
	Style           string      `json:"style,omitempty"`
	Explode         bool        `json:"explode,omitempty"`
	AllowReserved   bool        `json:"allowReserved,omitempty"`
	Schema          *Schema     `json:"schema,omitempty"`
	Example         interface{} `json:"example,omitempty"`
	Examples        map[string]*Example `json:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty"`
}

// RequestBody 请求体
type RequestBody struct {
	Description string                    `json:"description,omitempty"`
	Content     map[string]*MediaType    `json:"content"`
	Required    bool                     `json:"required,omitempty"`
}

// Response 响应
type Response struct {
	Description string                    `json:"description"`
	Headers     map[string]*Header       `json:"headers,omitempty"`
	Content     map[string]*MediaType    `json:"content,omitempty"`
	Links       map[string]*Link         `json:"links,omitempty"`
}

// Header 头部
type Header struct {
	Description     string      `json:"description,omitempty"`
	Required        bool        `json:"required,omitempty"`
	Deprecated      bool        `json:"deprecated,omitempty"`
	AllowEmptyValue bool        `json:"allowEmptyValue,omitempty"`
	Style           string      `json:"style,omitempty"`
	Explode         bool        `json:"explode,omitempty"`
	AllowReserved   bool        `json:"allowReserved,omitempty"`
	Schema          *Schema     `json:"schema,omitempty"`
	Example         interface{} `json:"example,omitempty"`
	Examples        map[string]*Example `json:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty"`
}

// MediaType 媒体类型
type MediaType struct {
	Schema   *Schema                `json:"schema,omitempty"`
	Example  interface{}            `json:"example,omitempty"`
	Examples map[string]*Example    `json:"examples,omitempty"`
	Encoding map[string]*Encoding   `json:"encoding,omitempty"`
}

// Schema 模式
type Schema struct {
	Type                 string                 `json:"type,omitempty"`
	Format               string                 `json:"format,omitempty"`
	Description          string                 `json:"description,omitempty"`
	Title                string                 `json:"title,omitempty"`
	Default              interface{}            `json:"default,omitempty"`
	Example              interface{}            `json:"example,omitempty"`
	Deprecated           bool                   `json:"deprecated,omitempty"`
	ReadOnly             bool                   `json:"readOnly,omitempty"`
	WriteOnly            bool                   `json:"writeOnly,omitempty"`
	XML                  *XML                   `json:"xml,omitempty"`
	ExternalDocs         *ExternalDocumentation `json:"externalDocs,omitempty"`
	Properties           map[string]*Schema     `json:"properties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	Items                *Schema                `json:"items,omitempty"`
	AllOf                []*Schema              `json:"allOf,omitempty"`
	OneOf                []*Schema              `json:"oneOf,omitempty"`
	AnyOf                []*Schema              `json:"anyOf,omitempty"`
	Not                  *Schema                `json:"not,omitempty"`
	AdditionalProperties *Schema                `json:"additionalProperties,omitempty"`
	Discriminator        *Discriminator         `json:"discriminator,omitempty"`
}

// XML XML 信息
type XML struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}

// ExternalDocumentation 外部文档
type ExternalDocumentation struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

// Example 示例
type Example struct {
	Summary       string      `json:"summary,omitempty"`
	Description   string      `json:"description,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty"`
}

// Encoding 编码
type Encoding struct {
	ContentType   string                    `json:"contentType,omitempty"`
	Headers       map[string]*Header        `json:"headers,omitempty"`
	Style         string                    `json:"style,omitempty"`
	Explode       bool                      `json:"explode,omitempty"`
	AllowReserved bool                      `json:"allowReserved,omitempty"`
}

// Discriminator 鉴别器
type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

// Link 链接
type Link struct {
	OperationRef string                 `json:"operationRef,omitempty"`
	OperationID  string                 `json:"operationId,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	RequestBody  interface{}            `json:"requestBody,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Server       *Server                `json:"server,omitempty"`
}

// Callback 回调
type Callback struct {
	Expression string `json:"expression"`
}

// SecurityRequirement 安全要求
type SecurityRequirement struct {
	Name string `json:"name"`
}

// Tag 标签
type Tag struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// Components 组件
type Components struct {
	Schemas         map[string]*Schema         `json:"schemas,omitempty"`
	Responses       map[string]*Response       `json:"responses,omitempty"`
	Parameters      map[string]*Parameter      `json:"parameters,omitempty"`
	Examples        map[string]*Example        `json:"examples,omitempty"`
	RequestBodies   map[string]*RequestBody    `json:"requestBodies,omitempty"`
	Headers         map[string]*Header         `json:"headers,omitempty"`
	SecuritySchemes map[string]*SecurityScheme `json:"securitySchemes,omitempty"`
	Links           map[string]*Link           `json:"links,omitempty"`
	Callbacks       map[string]*Callback       `json:"callbacks,omitempty"`
}

// SecurityScheme 安全方案
type SecurityScheme struct {
	Type             string      `json:"type"`
	Description      string      `json:"description,omitempty"`
	Name             string      `json:"name,omitempty"`
	In               string      `json:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty"`
	OpenIDConnectURL string      `json:"openIdConnectUrl,omitempty"`
}

// OAuthFlows OAuth 流程
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

// OAuthFlow OAuth 流程
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"`
}

// APIDocumentation API 文档生成器
type APIDocumentation struct {
	spec     *OpenAPISpec
	version  string
	basePath string
}

// NewAPIDocumentation 创建 API 文档生成器
func NewAPIDocumentation(title, version, description string) *APIDocumentation {
	return &APIDocumentation{
		spec: &OpenAPISpec{
			OpenAPI: "3.0.0",
			Info: &Info{
				Title:       title,
				Version:     version,
				Description: description,
			},
			Paths:      make(map[string]*PathItem),
			Components: &Components{},
		},
		version:  version,
		basePath: "/api",
	}
}

// SetBasePath 设置基础路径
func (ad *APIDocumentation) SetBasePath(basePath string) *APIDocumentation {
	ad.basePath = basePath
	return ad
}

// AddServer 添加服务器
func (ad *APIDocumentation) AddServer(url, description string) *APIDocumentation {
	server := &Server{
		URL:         url,
		Description: description,
	}
	ad.spec.Servers = append(ad.spec.Servers, server)
	return ad
}

// AddTag 添加标签
func (ad *APIDocumentation) AddTag(name, description string) *APIDocumentation {
	tag := &Tag{
		Name:        name,
		Description: description,
	}
	ad.spec.Tags = append(ad.spec.Tags, tag)
	return ad
}

// AddPath 添加路径
func (ad *APIDocumentation) AddPath(path, method string, operation *Operation) *APIDocumentation {
	fullPath := ad.basePath + path
	
	if ad.spec.Paths[fullPath] == nil {
		ad.spec.Paths[fullPath] = &PathItem{}
	}
	
	pathItem := ad.spec.Paths[fullPath]
	
	switch strings.ToUpper(method) {
	case "GET":
		pathItem.GET = operation
	case "POST":
		pathItem.POST = operation
	case "PUT":
		pathItem.PUT = operation
	case "DELETE":
		pathItem.DELETE = operation
	case "PATCH":
		pathItem.PATCH = operation
	case "OPTIONS":
		pathItem.OPTIONS = operation
	case "HEAD":
		pathItem.HEAD = operation
	case "TRACE":
		pathItem.TRACE = operation
	}
	
	return ad
}

// AddSchema 添加模式
func (ad *APIDocumentation) AddSchema(name string, schema *Schema) *APIDocumentation {
	if ad.spec.Components.Schemas == nil {
		ad.spec.Components.Schemas = make(map[string]*Schema)
	}
	ad.spec.Components.Schemas[name] = schema
	return ad
}

// GenerateSchemaFromStruct 从结构体生成模式
func (ad *APIDocumentation) GenerateSchemaFromStruct(name string, data interface{}) *Schema {
	schema := &Schema{
		Type:       "object",
		Properties: make(map[string]*Schema),
	}
	
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	
	if v.Kind() != reflect.Struct {
		return schema
	}
	
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		
		// 获取字段名
		fieldName := field.Name
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}
		
		// 检查是否应该跳过 omitempty 字段
		if strings.Contains(jsonTag, "omitempty") {
			// 跳过 omitempty 字段
			continue
		}
		
		// 生成字段模式
		fieldSchema := ad.generateFieldSchema(value.Type())
		
		// 添加到必需字段
		schema.Required = append(schema.Required, fieldName)
		
		schema.Properties[fieldName] = fieldSchema
	}
	
	return schema
}

// generateFieldSchema 生成字段模式
func (ad *APIDocumentation) generateFieldSchema(t reflect.Type) *Schema {
	switch t.Kind() {
	case reflect.String:
		return &Schema{Type: "string"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &Schema{Type: "integer", Format: "int64"}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &Schema{Type: "integer", Format: "uint64"}
	case reflect.Float32, reflect.Float64:
		return &Schema{Type: "number", Format: "double"}
	case reflect.Bool:
		return &Schema{Type: "boolean"}
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			return &Schema{Type: "string", Format: "date-time"}
		}
		// 递归处理嵌套结构体
		return ad.GenerateSchemaFromStruct(t.Name(), reflect.New(t).Interface())
	case reflect.Ptr:
		return ad.generateFieldSchema(t.Elem())
	case reflect.Slice, reflect.Array:
		return &Schema{
			Type:  "array",
			Items: ad.generateFieldSchema(t.Elem()),
		}
	case reflect.Map:
		return &Schema{
			Type: "object",
			AdditionalProperties: ad.generateFieldSchema(t.Elem()),
		}
	default:
		return &Schema{Type: "string"}
	}
}

// GenerateExample 生成示例
func (ad *APIDocumentation) GenerateExample(data interface{}) interface{} {
	return data
}

// ToJSON 生成 JSON 格式的文档
func (ad *APIDocumentation) ToJSON() ([]byte, error) {
	return json.MarshalIndent(ad.spec, "", "  ")
}

// ToYAML 生成 YAML 格式的文档（需要 yaml 包）
func (ad *APIDocumentation) ToYAML() ([]byte, error) {
	// 这里需要导入 yaml 包
	// return yaml.Marshal(ad.spec)
	return nil, fmt.Errorf("YAML support not implemented")
}

// GenerateHTML 生成 HTML 格式的文档
func (ad *APIDocumentation) GenerateHTML() string {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>` + ad.spec.Info.Title + ` - API Documentation</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
        .swagger-ui .topbar { display: none; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/api-docs.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
	return html
}

// 便捷函数

// NewOperation 创建新的操作
func NewOperation(summary, description string) *Operation {
	return &Operation{
		Summary:     summary,
		Description: description,
		Responses:   make(map[string]*Response),
	}
}

// NewResponse 创建新的响应
func NewResponse(description string) *Response {
	return &Response{
		Description: description,
	}
}

// NewParameter 创建新的参数
func NewParameter(name, in, description string, required bool) *Parameter {
	return &Parameter{
		Name:        name,
		In:          in,
		Description: description,
		Required:    required,
	}
}

// NewSchema 创建新的模式
func NewSchema(schemaType, format string) *Schema {
	return &Schema{
		Type:   schemaType,
		Format: format,
	}
} 