package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// InitConfig 初始化配置
func InitConfig() error {
	// 创建配置目录
	configDirs := []string{
		"config",
		"storage/framework/cache/data",
		"storage/framework/sessions",
		"storage/logs",
		"database",
		"database/migrations",
	}

	for _, dir := range configDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败 %s: %v", dir, err)
		}
	}

	// 复制环境变量示例文件
	if err := copyEnvExample(); err != nil {
		return fmt.Errorf("复制环境变量文件失败: %v", err)
	}

	// 创建默认配置文件
	if err := createDefaultConfigs(); err != nil {
		return fmt.Errorf("创建默认配置文件失败: %v", err)
	}

	return nil
}

// copyEnvExample 复制环境变量示例文件
func copyEnvExample() error {
	// 检查是否已存在 .env 文件
	if _, err := os.Stat(".env"); err == nil {
		return nil // 文件已存在，跳过
	}

	// 读取示例文件
	exampleContent := `# Laravel-Go Framework Environment Configuration
# Copy this file to .env and modify as needed

# Application Configuration
APP_NAME="Laravel-Go"
APP_VERSION=1.0.0
APP_ENV=production
APP_DEBUG=false
APP_URL=http://localhost:8080
APP_PORT=8080
APP_TIMEZONE=UTC
APP_LOCALE=en
APP_KEY=

# Database Configuration
DB_CONNECTION=sqlite
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel_go
DB_USERNAME=root
DB_PASSWORD=
DB_MIGRATIONS=database/migrations

# Redis Configuration
REDIS_HOST=127.0.0.1
REDIS_PASSWORD=null
REDIS_PORT=6379
REDIS_DB=0
REDIS_CACHE_DB=1
REDIS_QUEUE_DB=2

# Cache Configuration
CACHE_DRIVER=file
CACHE_PATH=storage/framework/cache/data
CACHE_PREFIX=laravel_go_cache

# Queue Configuration
QUEUE_CONNECTION=sync
REDIS_QUEUE=default

# Session Configuration
SESSION_DRIVER=file
SESSION_LIFETIME=120
SESSION_COOKIE=laravel_go_session
SESSION_DOMAIN=
SESSION_SECURE_COOKIE=false
SESSION_FILES=storage/framework/sessions

# Logging Configuration
LOG_CHANNEL=single
LOG_LEVEL=debug

# Mail Configuration (Optional)
MAIL_MAILER=smtp
MAIL_HOST=smtp.mailtrap.io
MAIL_PORT=2525
MAIL_USERNAME=null
MAIL_PASSWORD=null
MAIL_ENCRYPTION=null
MAIL_FROM_ADDRESS="hello@example.com"
MAIL_FROM_NAME="${APP_NAME}"

# AWS Configuration (Optional)
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_DEFAULT_REGION=us-east-1
AWS_BUCKET=

# SQS Configuration (Optional)
SQS_PREFIX=https://sqs.us-east-1.amazonaws.com/your-account-id
SQS_QUEUE=default
SQS_SUFFIX=

# Memcached Configuration (Optional)
MEMCACHED_HOST=127.0.0.1
MEMCACHED_PORT=11211
MEMCACHED_USERNAME=
MEMCACHED_PASSWORD=
MEMCACHED_PERSISTENT_ID=

# Beanstalkd Configuration (Optional)
BEANSTALKD_HOST=127.0.0.1
BEANSTALKD_QUEUE=default

# DynamoDB Configuration (Optional)
DYNAMODB_CACHE_TABLE=cache
DYNAMODB_ENDPOINT=

# Logging Services (Optional)
LOG_SLACK_WEBHOOK_URL=
PAPERTRAIL_URL=
PAPERTRAIL_PORT=

# Security Configuration
SESSION_SECURE_COOKIE=false
`

	// 写入 .env 文件
	if err := os.WriteFile(".env", []byte(exampleContent), 0644); err != nil {
		return err
	}

	fmt.Println("✅ 已创建 .env 文件")
	return nil
}

// createDefaultConfigs 创建默认配置文件
func createDefaultConfigs() error {
	configs := map[string]interface{}{
		"config/app.json": map[string]interface{}{
			"name":     "Laravel-Go",
			"version":  "1.0.0",
			"env":      "production",
			"debug":    false,
			"url":      "http://localhost:8080",
			"port":     "8080",
			"timezone": "UTC",
			"locale":   "en",
			"key":      "",
		},
		"config/database.json": map[string]interface{}{
			"default": "sqlite",
			"connections": map[string]interface{}{
				"sqlite": map[string]interface{}{
					"driver":   "sqlite",
					"database": "database/laravel-go.sqlite",
					"prefix":   "",
				},
				"mysql": map[string]interface{}{
					"driver":   "mysql",
					"host":     "127.0.0.1",
					"port":     "3306",
					"database": "laravel_go",
					"username": "root",
					"password": "",
					"charset":  "utf8mb4",
					"prefix":   "",
					"options": map[string]interface{}{
						"parseTime": "true",
						"loc":       "Local",
					},
				},
				"postgres": map[string]interface{}{
					"driver":   "postgres",
					"host":     "127.0.0.1",
					"port":     "5432",
					"database": "laravel_go",
					"username": "postgres",
					"password": "",
					"charset":  "utf8",
					"prefix":   "",
					"options": map[string]interface{}{
						"sslmode": "disable",
					},
				},
			},
			"migrations": "database/migrations",
			"redis": map[string]interface{}{
				"client": "predis",
				"default": map[string]interface{}{
					"host":     "127.0.0.1",
					"password": "",
					"port":     "6379",
					"database": "0",
				},
				"cache": map[string]interface{}{
					"host":     "127.0.0.1",
					"password": "",
					"port":     "6379",
					"database": "1",
				},
				"queue": map[string]interface{}{
					"host":     "127.0.0.1",
					"password": "",
					"port":     "6379",
					"database": "2",
				},
			},
		},
		"config/cache.json": map[string]interface{}{
			"default": "file",
			"stores": map[string]interface{}{
				"apc": map[string]interface{}{
					"driver": "apc",
				},
				"array": map[string]interface{}{
					"driver":    "array",
					"serialize": false,
				},
				"file": map[string]interface{}{
					"driver": "file",
					"path":   "storage/framework/cache/data",
				},
				"redis": map[string]interface{}{
					"driver":     "redis",
					"connection": "cache",
				},
				"database": map[string]interface{}{
					"driver": "database",
					"table":  "cache",
				},
			},
			"prefix": "laravel_go_cache",
		},
		"config/queue.json": map[string]interface{}{
			"default": "sync",
			"connections": map[string]interface{}{
				"sync": map[string]interface{}{
					"driver": "sync",
				},
				"database": map[string]interface{}{
					"driver":       "database",
					"table":        "jobs",
					"queue":        "default",
					"retry_after":  90,
					"after_commit": false,
				},
				"redis": map[string]interface{}{
					"driver":       "redis",
					"connection":   "default",
					"queue":        "default",
					"retry_after":  90,
					"after_commit": false,
				},
			},
			"failed": map[string]interface{}{
				"driver":   "database-uuids",
				"database": "pgsql",
				"table":    "failed_jobs",
			},
		},
		"config/session.json": map[string]interface{}{
			"driver":          "file",
			"lifetime":        120,
			"expire_on_close": false,
			"encrypt":         false,
			"files":           "storage/framework/sessions",
			"connection":      "",
			"table":           "sessions",
			"store":           "",
			"lottery":         []int{2, 100},
			"cookie":          "laravel_go_session",
			"path":            "/",
			"domain":          "",
			"secure":          false,
			"http_only":       true,
			"same_site":       "lax",
		},
		"config/logging.json": map[string]interface{}{
			"default": "stack",
			"deprecations": map[string]interface{}{
				"channel": "null",
				"trace":   false,
			},
			"channels": map[string]interface{}{
				"stack": map[string]interface{}{
					"driver":            "stack",
					"channels":          []string{"single"},
					"ignore_exceptions": false,
				},
				"single": map[string]interface{}{
					"driver": "single",
					"path":   "storage/logs/laravel.log",
					"level":  "debug",
				},
				"daily": map[string]interface{}{
					"driver": "daily",
					"path":   "storage/logs/laravel.log",
					"level":  "debug",
					"days":   14,
				},
				"stderr": map[string]interface{}{
					"driver":    "monolog",
					"level":     "debug",
					"handler":   "Monolog\\Handler\\StreamHandler",
					"formatter": "Monolog\\Formatter\\JsonFormatter",
					"with": map[string]interface{}{
						"stream": "php://stderr",
					},
				},
				"syslog": map[string]interface{}{
					"driver": "syslog",
					"level":  "debug",
				},
				"errorlog": map[string]interface{}{
					"driver": "errorlog",
					"level":  "debug",
				},
				"null": map[string]interface{}{
					"driver":  "monolog",
					"handler": "Monolog\\Handler\\NullHandler",
				},
				"emergency": map[string]interface{}{
					"path": "storage/logs/laravel.log",
				},
			},
		},
	}

	for filePath, config := range configs {
		// 检查文件是否已存在
		if _, err := os.Stat(filePath); err == nil {
			continue // 文件已存在，跳过
		}

		// 创建目录
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败 %s: %v", dir, err)
		}

		// 序列化为JSON
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化配置失败 %s: %v", filePath, err)
		}

		// 写入文件
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("写入配置文件失败 %s: %v", filePath, err)
		}

		fmt.Printf("✅ 已创建配置文件: %s\n", filePath)
	}

	return nil
}

// GenerateAppKey 生成应用密钥
func GenerateAppKey() string {
	// 这里应该实现一个安全的密钥生成算法
	// 暂时返回一个简单的示例
	return "base64:your-32-character-app-key-here"
}
