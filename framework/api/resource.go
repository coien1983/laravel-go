package api

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

// Resource 资源转换器接口
type Resource interface {
	// ToArray 转换为数组格式
	ToArray() map[string]interface{}
	// ToJSON 转换为 JSON 格式
	ToJSON() ([]byte, error)
	// With 添加额外的字段
	With(fields ...string) Resource
	// Without 移除指定字段
	Without(fields ...string) Resource
	// When 条件性包含字段
	When(condition bool, fields ...string) Resource
	// Merge 合并其他资源
	Merge(resource Resource) Resource
}

// Collection 集合转换器接口
type Collection interface {
	// ToArray 转换为数组格式
	ToArray() []map[string]interface{}
	// ToJSON 转换为 JSON 格式
	ToJSON() ([]byte, error)
	// With 为所有资源添加额外字段
	With(fields ...string) Collection
	// Without 为所有资源移除指定字段
	Without(fields ...string) Collection
	// When 条件性包含字段
	When(condition bool, fields ...string) Collection
	// Map 映射集合中的每个资源
	Map(fn func(Resource) Resource) Collection
	// Filter 过滤集合
	Filter(fn func(Resource) bool) Collection
	// Paginate 分页
	Paginate(page, perPage int) Collection
}

// BaseResource 基础资源转换器
type BaseResource struct {
	data       interface{}
	fields     []string
	hidden     []string
	conditions map[string]bool
	additional map[string]interface{}
}

// NewResource 创建新的资源转换器
func NewResource(data interface{}) *BaseResource {
	return &BaseResource{
		data:       data,
		fields:     []string{},
		hidden:     []string{},
		conditions: make(map[string]bool),
		additional: make(map[string]interface{}),
	}
}

// ToArray 转换为数组格式
func (r *BaseResource) ToArray() map[string]interface{} {
	result := make(map[string]interface{})
	
	// 获取结构体的所有字段
	v := reflect.ValueOf(r.data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	
	if v.Kind() != reflect.Struct {
		// 如果不是结构体，直接返回
		return result
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
		
		// 检查是否应该隐藏此字段
		if r.shouldHide(fieldName) {
			continue
		}
		
		// 检查是否应该包含此字段
		if !r.shouldInclude(fieldName) {
			continue
		}
		
		// 获取字段值
		fieldValue := r.getFieldValue(value)
		result[fieldName] = fieldValue
	}
	
	// 添加额外字段
	for key, value := range r.additional {
		result[key] = value
	}
	
	return result
}

// ToJSON 转换为 JSON 格式
func (r *BaseResource) ToJSON() ([]byte, error) {
	return json.Marshal(r.ToArray())
}

// With 添加额外的字段
func (r *BaseResource) With(fields ...string) Resource {
	newResource := &BaseResource{
		data:       r.data,
		fields:     append([]string{}, r.fields...),
		hidden:     append([]string{}, r.hidden...),
		conditions: make(map[string]bool),
		additional: make(map[string]interface{}),
	}
	
	// 复制条件
	for k, v := range r.conditions {
		newResource.conditions[k] = v
	}
	
	// 复制额外字段
	for k, v := range r.additional {
		newResource.additional[k] = v
	}
	
	// 添加新字段
	newResource.fields = append(newResource.fields, fields...)
	
	return newResource
}

// Without 移除指定字段
func (r *BaseResource) Without(fields ...string) Resource {
	newResource := &BaseResource{
		data:       r.data,
		fields:     append([]string{}, r.fields...),
		hidden:     append([]string{}, r.hidden...),
		conditions: make(map[string]bool),
		additional: make(map[string]interface{}),
	}
	
	// 复制条件
	for k, v := range r.conditions {
		newResource.conditions[k] = v
	}
	
	// 复制额外字段
	for k, v := range r.additional {
		newResource.additional[k] = v
	}
	
	// 添加隐藏字段
	newResource.hidden = append(newResource.hidden, fields...)
	
	return newResource
}

// When 条件性包含字段
func (r *BaseResource) When(condition bool, fields ...string) Resource {
	newResource := &BaseResource{
		data:       r.data,
		fields:     append([]string{}, r.fields...),
		hidden:     append([]string{}, r.hidden...),
		conditions: make(map[string]bool),
		additional: make(map[string]interface{}),
	}
	
	// 复制条件
	for k, v := range r.conditions {
		newResource.conditions[k] = v
	}
	
	// 复制额外字段
	for k, v := range r.additional {
		newResource.additional[k] = v
	}
	
	// 添加条件字段
	for _, field := range fields {
		newResource.conditions[field] = condition
	}
	
	return newResource
}

// Merge 合并其他资源
func (r *BaseResource) Merge(resource Resource) Resource {
	if other, ok := resource.(*BaseResource); ok {
		for key, value := range other.additional {
			r.additional[key] = value
		}
	}
	return r
}

// Add 添加额外字段
func (r *BaseResource) Add(key string, value interface{}) Resource {
	r.additional[key] = value
	return r
}

// shouldHide 检查是否应该隐藏字段
func (r *BaseResource) shouldHide(fieldName string) bool {
	for _, hidden := range r.hidden {
		if hidden == fieldName {
			return true
		}
	}
	return false
}

// shouldInclude 检查是否应该包含字段
func (r *BaseResource) shouldInclude(fieldName string) bool {
	// 如果有指定字段，只包含指定的字段
	if len(r.fields) > 0 {
		for _, field := range r.fields {
			if field == fieldName {
				return true
			}
		}
		return false
	}
	
	// 检查条件字段
	if condition, exists := r.conditions[fieldName]; exists {
		return condition
	}
	
	return true
}

// getFieldValue 获取字段值
func (r *BaseResource) getFieldValue(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint()
	case reflect.Float32, reflect.Float64:
		return value.Float()
	case reflect.Bool:
		return value.Bool()
	case reflect.Struct:
		// 处理时间类型
		if value.Type() == reflect.TypeOf(time.Time{}) {
			return value.Interface().(time.Time).Format(time.RFC3339)
		}
		// 递归处理嵌套结构体
		return NewResource(value.Interface()).ToArray()
	case reflect.Ptr:
		if value.IsNil() {
			return nil
		}
		return r.getFieldValue(value.Elem())
	case reflect.Slice, reflect.Array:
		// 处理切片
		result := make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			result[i] = r.getFieldValue(value.Index(i))
		}
		return result
	case reflect.Map:
		// 处理映射
		result := make(map[string]interface{})
		for _, key := range value.MapKeys() {
			result[key.String()] = r.getFieldValue(value.MapIndex(key))
		}
		return result
	default:
		return value.Interface()
	}
}

// BaseCollection 基础集合转换器
type BaseCollection struct {
	resources []Resource
	fields    []string
	hidden    []string
	conditions map[string]bool
}

// NewCollection 创建新的集合转换器
func NewCollection(resources []Resource) *BaseCollection {
	return &BaseCollection{
		resources:  resources,
		fields:     []string{},
		hidden:     []string{},
		conditions: make(map[string]bool),
	}
}

// ToArray 转换为数组格式
func (c *BaseCollection) ToArray() []map[string]interface{} {
	result := make([]map[string]interface{}, len(c.resources))
	for i, resource := range c.resources {
		if baseResource, ok := resource.(*BaseResource); ok {
			// 应用集合级别的字段过滤
			baseResource.fields = append(baseResource.fields, c.fields...)
			baseResource.hidden = append(baseResource.hidden, c.hidden...)
			for key, value := range c.conditions {
				baseResource.conditions[key] = value
			}
		}
		result[i] = resource.ToArray()
	}
	return result
}

// ToJSON 转换为 JSON 格式
func (c *BaseCollection) ToJSON() ([]byte, error) {
	return json.Marshal(c.ToArray())
}

// With 为所有资源添加额外字段
func (c *BaseCollection) With(fields ...string) Collection {
	c.fields = append(c.fields, fields...)
	return c
}

// Without 为所有资源移除指定字段
func (c *BaseCollection) Without(fields ...string) Collection {
	c.hidden = append(c.hidden, fields...)
	return c
}

// When 条件性包含字段
func (c *BaseCollection) When(condition bool, fields ...string) Collection {
	for _, field := range fields {
		c.conditions[field] = condition
	}
	return c
}

// Map 映射集合中的每个资源
func (c *BaseCollection) Map(fn func(Resource) Resource) Collection {
	newResources := make([]Resource, len(c.resources))
	for i, resource := range c.resources {
		newResources[i] = fn(resource)
	}
	return &BaseCollection{
		resources:  newResources,
		fields:     c.fields,
		hidden:     c.hidden,
		conditions: c.conditions,
	}
}

// Filter 过滤集合
func (c *BaseCollection) Filter(fn func(Resource) bool) Collection {
	var filteredResources []Resource
	for _, resource := range c.resources {
		if fn(resource) {
			filteredResources = append(filteredResources, resource)
		}
	}
	return &BaseCollection{
		resources:  filteredResources,
		fields:     c.fields,
		hidden:     c.hidden,
		conditions: c.conditions,
	}
}

// Paginate 分页
func (c *BaseCollection) Paginate(page, perPage int) Collection {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	
	start := (page - 1) * perPage
	end := start + perPage
	
	if start >= len(c.resources) {
		return &BaseCollection{
			resources:  []Resource{},
			fields:     c.fields,
			hidden:     c.hidden,
			conditions: c.conditions,
		}
	}
	
	if end > len(c.resources) {
		end = len(c.resources)
	}
	
	return &BaseCollection{
		resources:  c.resources[start:end],
		fields:     c.fields,
		hidden:     c.hidden,
		conditions: c.conditions,
	}
}

// ResourceCollection 资源集合，用于从数据切片创建集合
type ResourceCollection struct {
	*BaseCollection
}

// NewResourceCollection 从数据切片创建资源集合
func NewResourceCollection(data interface{}) *ResourceCollection {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return &ResourceCollection{
			BaseCollection: NewCollection([]Resource{}),
		}
	}
	
	resources := make([]Resource, v.Len())
	for i := 0; i < v.Len(); i++ {
		resources[i] = NewResource(v.Index(i).Interface())
	}
	
	return &ResourceCollection{
		BaseCollection: NewCollection(resources),
	}
}

// 便捷函数

// NewResourceFromData 创建单个资源
func NewResourceFromData(data interface{}) Resource {
	return NewResource(data)
}

// NewCollectionFromData 创建资源集合
func NewCollectionFromData(data interface{}) Collection {
	return NewResourceCollection(data)
}

// NewResourceFromSlice 从切片创建资源集合
func NewResourceFromSlice(data interface{}) Collection {
	return NewResourceCollection(data)
} 