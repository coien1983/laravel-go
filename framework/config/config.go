package config

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Config 配置管理器
type Config struct {
	data  map[string]interface{}
	env   *Env
	mutex sync.RWMutex
}

// Env 环境变量管理器
type Env struct {
	loaded bool
	mutex  sync.RWMutex
}

// NewConfig 创建新的配置管理器
func NewConfig() *Config {
	return &Config{
		data: make(map[string]interface{}),
		env:  &Env{},
	}
}

// LoadEnv 加载环境变量文件
func (c *Config) LoadEnv(filenames ...string) error {
	c.env.mutex.Lock()
	defer c.env.mutex.Unlock()

	if c.env.loaded {
		return nil
	}

	filename := ".env"
	if len(filenames) > 0 {
		filename = filenames[0]
	}

	file, err := os.Open(filename)
	if err != nil {
		// 如果文件不存在，不报错
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析键值对
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// 移除引号
			if len(value) >= 2 && (value[0] == '"' && value[len(value)-1] == '"') {
				value = value[1 : len(value)-1]
			}

			os.Setenv(key, value)
		}
	}

	c.env.loaded = true
	return scanner.Err()
}

// Get 获取配置值
func (c *Config) Get(key string, defaultValue ...interface{}) interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// 支持点号分隔的嵌套键
	keys := strings.Split(key, ".")
	var value interface{} = c.data

	for _, k := range keys {
		if currentMap, ok := value.(map[string]interface{}); ok {
			if v, exists := currentMap[k]; exists {
				value = v
			} else {
				// 尝试从环境变量获取
				if envValue := c.getEnvValue(key); envValue != nil {
					return envValue
				}

				if len(defaultValue) > 0 {
					return defaultValue[0]
				}
				return nil
			}
		} else {
			// 如果不是map类型，尝试从环境变量获取
			if envValue := c.getEnvValue(key); envValue != nil {
				return envValue
			}

			if len(defaultValue) > 0 {
				return defaultValue[0]
			}
			return nil
		}
	}

	return value
}

// Set 设置配置值
func (c *Config) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	keys := strings.Split(key, ".")
	current := c.data

	// 遍历到最后一个键
	for _, k := range keys[:len(keys)-1] {
		if _, exists := current[k]; !exists {
			current[k] = make(map[string]interface{})
		}
		current = current[k].(map[string]interface{})
	}

	// 设置最后一个键的值
	current[keys[len(keys)-1]] = value
}

// Has 检查配置是否存在
func (c *Config) Has(key string) bool {
	return c.Get(key) != nil
}

// All 获取所有配置
func (c *Config) All() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	result := make(map[string]interface{})
	for k, v := range c.data {
		result[k] = v
	}
	return result
}

// LoadFromFile 从文件加载配置
func (c *Config) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	configMap := make(map[string]interface{})

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析键值对
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// 移除引号
			if len(value) >= 2 && (value[0] == '"' && value[len(value)-1] == '"') {
				value = value[1 : len(value)-1]
			}

			configMap[key] = interface{}(value)
		}
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for k, v := range configMap {
		if mapValue, ok := v.(map[string]interface{}); ok {
			c.data[k] = mapValue
		} else {
			c.data[k] = v
		}
	}

	return scanner.Err()
}

// LoadFromStruct 从结构体加载配置
func (c *Config) LoadFromStruct(config interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	configValue := reflect.ValueOf(config)
	configType := configValue.Type()

	// 如果是指针，获取元素
	if configType.Kind() == reflect.Ptr {
		configValue = configValue.Elem()
		configType = configType.Elem()
	}

	// 遍历结构体字段
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		fieldValue := configValue.Field(i)

		// 获取字段标签
		envTag := field.Tag.Get("env")
		defaultTag := field.Tag.Get("default")

		var value interface{}

		// 优先从环境变量获取
		if envTag != "" {
			if envValue := os.Getenv(envTag); envValue != "" {
				value = c.convertValue(envValue, field.Type)
			}
		}

		// 如果环境变量没有值，使用默认值
		if value == nil && defaultTag != "" {
			value = c.convertValue(defaultTag, field.Type)
		}

		// 如果都没有，使用字段的零值
		if value == nil {
			value = fieldValue.Interface()
		}

		// 设置配置值
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			fieldName = strings.Split(jsonTag, ",")[0]
		}

		c.data[fieldName] = value
	}

	return nil
}

// getEnvValue 从环境变量获取值
func (c *Config) getEnvValue(key string) interface{} {
	// 将配置键转换为环境变量格式
	envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))

	if value := os.Getenv(envKey); value != "" {
		return value
	}

	return nil
}

// convertValue 转换值类型
func (c *Config) convertValue(value string, targetType reflect.Type) interface{} {
	switch targetType.Kind() {
	case reflect.String:
		return value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			return uintValue
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}

	return value
}

// GetString 获取字符串配置
func (c *Config) GetString(key string, defaultValue ...string) string {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}

	if str, ok := value.(string); ok {
		return str
	}

	return fmt.Sprintf("%v", value)
}

// GetInt 获取整数配置
func (c *Config) GetInt(key string, defaultValue ...int) int {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if intValue, err := strconv.Atoi(v); err == nil {
			return intValue
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetBool 获取布尔配置
func (c *Config) GetBool(key string, defaultValue ...bool) bool {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		if boolValue, err := strconv.ParseBool(v); err == nil {
			return boolValue
		}
	case int:
		return v != 0
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// GetFloat 获取浮点数配置
func (c *Config) GetFloat(key string, defaultValue ...float64) float64 {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0.0
	}

	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		if floatValue, err := strconv.ParseFloat(v, 64); err == nil {
			return floatValue
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0.0
}

// GetStringSlice 获取字符串切片配置
func (c *Config) GetStringSlice(key string, defaultValue ...[]string) []string {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []string{}
	}

	switch v := value.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result
	case string:
		return strings.Split(v, ",")
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return []string{}
}

// GetMap 获取映射配置
func (c *Config) GetMap(key string, defaultValue ...map[string]interface{}) map[string]interface{} {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return make(map[string]interface{})
	}

	if m, ok := value.(map[string]interface{}); ok {
		return m
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return make(map[string]interface{})
}

// Validate 验证配置
func (c *Config) Validate(rules map[string]string) error {
	for key, rule := range rules {
		value := c.Get(key)

		// 解析验证规则
		ruleParts := strings.Split(rule, "|")
		for _, part := range ruleParts {
			if err := c.validateRule(key, value, part); err != nil {
				return fmt.Errorf("validation failed for %s: %v", key, err)
			}
		}
	}

	return nil
}

// validateRule 验证单个规则
func (c *Config) validateRule(key string, value interface{}, rule string) error {
	switch rule {
	case "required":
		if value == nil || value == "" {
			return fmt.Errorf("%s is required", key)
		}
	case "string":
		if value != nil && reflect.TypeOf(value).Kind() != reflect.String {
			return fmt.Errorf("%s must be a string", key)
		}
	case "int":
		if value != nil {
			switch value.(type) {
			case int, int64, float64:
				// 有效
			default:
				return fmt.Errorf("%s must be an integer", key)
			}
		}
	case "bool":
		if value != nil && reflect.TypeOf(value).Kind() != reflect.Bool {
			return fmt.Errorf("%s must be a boolean", key)
		}
	}

	return nil
}
