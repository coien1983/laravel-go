package config

import (
	"laravel-go/framework/config"
	"os"
)

var Config *config.Config

func init() {
	Config = config.NewConfig()
	Config.Set("database", map[string]interface{}{
		"default": "postgres",
		"connections": map[string]interface{}{
			"postgres": map[string]interface{}{
				"driver":   "postgres",
				"host":     getEnv("DB_HOST", "localhost"),
				"port":     getEnv("DB_PORT", "5432"),
				"database": getEnv("DB_DATABASE", "laravel_go_blog"),
				"username": getEnv("DB_USERNAME", "laravel_go"),
				"password": getEnv("DB_PASSWORD", "password"),
				"charset":  "utf8mb4",
				"collation": "utf8mb4_unicode_ci",
				"prefix":   "",
				"strict":   true,
				"engine":   "",
				"options": map[string]interface{}{
					"sslmode": "disable",
				},
			},
			"sqlite": map[string]interface{}{
				"driver":   "sqlite",
				"database": getEnv("DB_DATABASE", "database/blog.sqlite"),
				"prefix":   "",
			},
		},
		"migrations": "database/migrations",
		"redis": map[string]interface{}{
			"client": "predis",
			"default": map[string]interface{}{
				"host":     getEnv("REDIS_HOST", "127.0.0.1"),
				"password": getEnv("REDIS_PASSWORD", ""),
				"port":     getEnv("REDIS_PORT", "6379"),
				"database": "0",
			},
			"cache": map[string]interface{}{
				"host":     getEnv("REDIS_HOST", "127.0.0.1"),
				"password": getEnv("REDIS_PASSWORD", ""),
				"port":     getEnv("REDIS_PORT", "6379"),
				"database": "1",
			},
		},
	})
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 