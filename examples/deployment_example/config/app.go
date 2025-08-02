package config

import (
	"laravel-go/framework/config"
	"os"
)

var Config *config.Config

func init() {
	Config = config.NewConfig()
	Config.Set("app", map[string]interface{}{
		"name":     "Laravel-Go 部署示例",
		"version":  "1.0.0",
		"port":     getEnv("APP_PORT", "8080"),
		"debug":    getEnv("APP_DEBUG", "true") == "true",
		"timezone": getEnv("APP_TIMEZONE", "Asia/Shanghai"),
		"url":      getEnv("APP_URL", "http://localhost:8080"),
		"env":      getEnv("APP_ENV", "development"),
	})
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 