package config

import (
	"os"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Connection string
	Host       string
	Port       string
	Database   string
	Username   string
	Password   string
}

// LoadDatabaseConfig 加载数据库配置
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Connection: getEnv("DB_CONNECTION", "sqlite"),
		Host:       getEnv("DB_HOST", "127.0.0.1"),
		Port:       getEnv("DB_PORT", "3306"),
		Database:   getEnv("DB_DATABASE", "app.db"),
		Username:   getEnv("DB_USERNAME", ""),
		Password:   getEnv("DB_PASSWORD", ""),
	}
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}