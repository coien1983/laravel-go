package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"laravel-go/framework/errors"
)

// ConfigManager 配置管理器
type ConfigManager struct {
	configs map[string]interface{}
	mutex   sync.RWMutex
	// 配置文件路径
	configPath string
	// 配置变更监听器
	listeners     map[string][]ConfigChangeListener
	listenerMutex sync.RWMutex
	// 配置缓存
	cache      map[string]*ConfigCache
	cacheMutex sync.RWMutex
}

// ConfigChangeListener 配置变更监听器
type ConfigChangeListener func(key string, oldValue, newValue interface{})

// ConfigCache 配置缓存
type ConfigCache struct {
	Value      interface{}
	Expiration time.Time
}

// IsExpired 检查缓存是否过期
func (cc *ConfigCache) IsExpired() bool {
	return !cc.Expiration.IsZero() && time.Now().After(cc.Expiration)
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configPath string) *ConfigManager {
	cm := &ConfigManager{
		configs:    make(map[string]interface{}),
		configPath: configPath,
		listeners:  make(map[string][]ConfigChangeListener),
		cache:      make(map[string]*ConfigCache),
	}

	// 加载配置文件
	if err := cm.loadConfigs(); err != nil {
		// 记录错误但不中断启动
		fmt.Printf("Warning: Failed to load configs: %v\n", err)
	}

	return cm
}

// loadConfigs 加载配置文件
func (cm *ConfigManager) loadConfigs() error {
	if cm.configPath == "" {
		return nil
	}

	// 遍历配置文件目录
	err := filepath.Walk(cm.configPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		// 读取配置文件
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var config map[string]interface{}
		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}

		// 获取配置名称（文件名去掉扩展名）
		configName := filepath.Base(path[:len(path)-len(filepath.Ext(path))])

		cm.mutex.Lock()
		cm.configs[configName] = config
		cm.mutex.Unlock()

		return nil
	})

	return err
}

// Get 获取配置值
func (cm *ConfigManager) Get(key string, defaultValue interface{}) interface{} {
	// 先检查缓存
	cm.cacheMutex.RLock()
	if cache, exists := cm.cache[key]; exists && !cache.IsExpired() {
		cm.cacheMutex.RUnlock()
		return cache.Value
	}
	cm.cacheMutex.RUnlock()

	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// 解析配置键路径
	value := cm.getNestedValue(cm.configs, key)
	if value == nil {
		return defaultValue
	}

	// 缓存结果（5分钟）
	cm.cacheMutex.Lock()
	cm.cache[key] = &ConfigCache{
		Value:      value,
		Expiration: time.Now().Add(5 * time.Minute),
	}
	cm.cacheMutex.Unlock()

	return value
}

// GetString 获取字符串配置
func (cm *ConfigManager) GetString(key string, defaultValue string) string {
	value := cm.Get(key, defaultValue)
	if str, ok := value.(string); ok {
		return str
	}
	return defaultValue
}

// GetInt 获取整数配置
func (cm *ConfigManager) GetInt(key string, defaultValue int) int {
	value := cm.Get(key, defaultValue)
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

// GetBool 获取布尔配置
func (cm *ConfigManager) GetBool(key string, defaultValue bool) bool {
	value := cm.Get(key, defaultValue)
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "1"
	case int:
		return v != 0
	}
	return defaultValue
}

// Set 设置配置值
func (cm *ConfigManager) Set(key string, value interface{}) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 获取旧值
	oldValue := cm.getNestedValue(cm.configs, key)

	// 设置新值
	if err := cm.setNestedValue(cm.configs, key, value); err != nil {
		return err
	}

	// 清除缓存
	cm.cacheMutex.Lock()
	delete(cm.cache, key)
	cm.cacheMutex.Unlock()

	// 通知监听器
	cm.notifyListeners(key, oldValue, value)

	return nil
}

// AddListener 添加配置变更监听器
func (cm *ConfigManager) AddListener(key string, listener ConfigChangeListener) {
	cm.listenerMutex.Lock()
	defer cm.listenerMutex.Unlock()

	cm.listeners[key] = append(cm.listeners[key], listener)
}

// RemoveListener 移除配置变更监听器
func (cm *ConfigManager) RemoveListener(key string, listener ConfigChangeListener) {
	cm.listenerMutex.Lock()
	defer cm.listenerMutex.Unlock()

	if listeners, exists := cm.listeners[key]; exists {
		for i, l := range listeners {
			if fmt.Sprintf("%p", l) == fmt.Sprintf("%p", listener) {
				cm.listeners[key] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
}

// notifyListeners 通知监听器
func (cm *ConfigManager) notifyListeners(key string, oldValue, newValue interface{}) {
	cm.listenerMutex.RLock()
	listeners := cm.listeners[key]
	cm.listenerMutex.RUnlock()

	for _, listener := range listeners {
		go func(l ConfigChangeListener) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Config listener panic: %v\n", r)
				}
			}()
			l(key, oldValue, newValue)
		}(listener)
	}
}

// getNestedValue 获取嵌套值
func (cm *ConfigManager) getNestedValue(data interface{}, key string) interface{} {
	keys := splitKey(key)
	current := data

	for _, k := range keys {
		switch v := current.(type) {
		case map[string]interface{}:
			if val, exists := v[k]; exists {
				current = val
			} else {
				return nil
			}
		case map[interface{}]interface{}:
			if val, exists := v[k]; exists {
				current = val
			} else {
				return nil
			}
		default:
			return nil
		}
	}

	return current
}

// setNestedValue 设置嵌套值
func (cm *ConfigManager) setNestedValue(data interface{}, key string, value interface{}) error {
	keys := splitKey(key)

	if len(keys) == 0 {
		return errors.New("empty key")
	}

	// 确保根是map
	if cm.configs == nil {
		cm.configs = make(map[string]interface{})
	}

	current := cm.configs
	for i, k := range keys[:len(keys)-1] {
		if _, exists := current[k]; !exists {
			current[k] = make(map[string]interface{})
		}

		if next, ok := current[k].(map[string]interface{}); ok {
			current = next
		} else {
			return errors.New(fmt.Sprintf("key %s is not a map", k))
		}
	}

	// 设置最终值
	current[keys[len(keys)-1]] = value
	return nil
}

// splitKey 分割配置键
func splitKey(key string) []string {
	var keys []string
	var current string
	inBracket := false

	for _, char := range key {
		switch char {
		case '.':
			if !inBracket {
				if current != "" {
					keys = append(keys, current)
					current = ""
				}
			} else {
				current += string(char)
			}
		case '[':
			if !inBracket {
				if current != "" {
					keys = append(keys, current)
					current = ""
				}
				inBracket = true
			} else {
				current += string(char)
			}
		case ']':
			if inBracket {
				keys = append(keys, current)
				current = ""
				inBracket = false
			} else {
				current += string(char)
			}
		default:
			current += string(char)
		}
	}

	if current != "" {
		keys = append(keys, current)
	}

	return keys
}

// ClearCache 清除配置缓存
func (cm *ConfigManager) ClearCache() {
	cm.cacheMutex.Lock()
	defer cm.cacheMutex.Unlock()

	cm.cache = make(map[string]*ConfigCache)
}

// GetStats 获取配置管理器统计信息
func (cm *ConfigManager) GetStats() map[string]interface{} {
	cm.cacheMutex.RLock()
	cacheSize := len(cm.cache)
	cm.cacheMutex.RUnlock()

	cm.listenerMutex.RLock()
	listenerCount := 0
	for _, listeners := range cm.listeners {
		listenerCount += len(listeners)
	}
	cm.listenerMutex.RUnlock()

	return map[string]interface{}{
		"config_count":   len(cm.configs),
		"cache_size":     cacheSize,
		"listener_count": listenerCount,
	}
}
