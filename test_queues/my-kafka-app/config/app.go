package config

import (
	"os"
	"strconv"
)

// AppConfig 应用配置
type AppConfig struct {
	Name   string
	Env    string
	Debug  bool
	URL    string
	Port   string
}

// LoadAppConfig 加载应用配置
func LoadAppConfig() *AppConfig {
	debug, _ := strconv.ParseBool(getEnv("APP_DEBUG", "true"))
	
	return &AppConfig{
		Name:  getEnv("APP_NAME", "Laravel-Go App"),
		Env:   getEnv("APP_ENV", "local"),
		Debug: debug,
		URL:   getEnv("APP_URL", "http://localhost:8080"),
		Port:  getEnv("PORT", "8080"),
	}
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}