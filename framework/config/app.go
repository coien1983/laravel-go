package config

import (
	"os"
	"strconv"
)

// AppConfig 应用配置结构
type AppConfig struct {
	Name      string   `json:"name"`
	Version   string   `json:"version"`
	Env       string   `json:"env"`
	Debug     bool     `json:"debug"`
	URL       string   `json:"url"`
	Port      string   `json:"port"`
	Timezone  string   `json:"timezone"`
	Locale    string   `json:"locale"`
	Key       string   `json:"key"`
	Providers []string `json:"providers"`
}

// LoadAppConfig 加载应用配置
func LoadAppConfig() *AppConfig {
	return &AppConfig{
		Name:     getEnv("APP_NAME", "Laravel-Go"),
		Version:  getEnv("APP_VERSION", "1.0.0"),
		Env:      getEnv("APP_ENV", "production"),
		Debug:    getEnvBool("APP_DEBUG", false),
		URL:      getEnv("APP_URL", "http://localhost"),
		Port:     getEnv("APP_PORT", "8080"),
		Timezone: getEnv("APP_TIMEZONE", "UTC"),
		Locale:   getEnv("APP_LOCALE", "en"),
		Key:      getEnv("APP_KEY", ""),
		Providers: []string{
			"laravel-go/framework/providers/AppServiceProvider",
			"laravel-go/framework/providers/RouteServiceProvider",
			"laravel-go/framework/providers/DatabaseServiceProvider",
			"laravel-go/framework/providers/CacheServiceProvider",
			"laravel-go/framework/providers/QueueServiceProvider",
			"laravel-go/framework/providers/EventServiceProvider",
			"laravel-go/framework/providers/LogServiceProvider",
		},
	}
}

// DatabaseConfig 数据库配置结构
type DatabaseConfig struct {
	Default     string                      `json:"default"`
	Connections map[string]ConnectionConfig `json:"connections"`
	Migrations  string                      `json:"migrations"`
	Redis       RedisConfig                 `json:"redis"`
}

// ConnectionConfig 连接配置结构
type ConnectionConfig struct {
	Driver   string            `json:"driver"`
	Host     string            `json:"host"`
	Port     string            `json:"port"`
	Database string            `json:"database"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Charset  string            `json:"charset"`
	Prefix   string            `json:"prefix"`
	Options  map[string]string `json:"options"`
}

// RedisConfig Redis配置结构
type RedisConfig struct {
	Client  string                `json:"client"`
	Default RedisConnectionConfig `json:"default"`
	Cache   RedisConnectionConfig `json:"cache"`
	Queue   RedisConnectionConfig `json:"queue"`
}

// RedisConnectionConfig Redis连接配置结构
type RedisConnectionConfig struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

// LoadDatabaseConfig 加载数据库配置
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Default: getEnv("DB_CONNECTION", "sqlite"),
		Connections: map[string]ConnectionConfig{
			"sqlite": {
				Driver:   "sqlite",
				Database: getEnv("DB_DATABASE", "database/laravel-go.sqlite"),
				Prefix:   "",
			},
			"mysql": {
				Driver:   "mysql",
				Host:     getEnv("DB_HOST", "127.0.0.1"),
				Port:     getEnv("DB_PORT", "3306"),
				Database: getEnv("DB_DATABASE", "laravel_go"),
				Username: getEnv("DB_USERNAME", "root"),
				Password: getEnv("DB_PASSWORD", ""),
				Charset:  "utf8mb4",
				Prefix:   "",
				Options: map[string]string{
					"parseTime": "true",
					"loc":       "Local",
				},
			},
			"postgres": {
				Driver:   "postgres",
				Host:     getEnv("DB_HOST", "127.0.0.1"),
				Port:     getEnv("DB_PORT", "5432"),
				Database: getEnv("DB_DATABASE", "laravel_go"),
				Username: getEnv("DB_USERNAME", "postgres"),
				Password: getEnv("DB_PASSWORD", ""),
				Charset:  "utf8",
				Prefix:   "",
				Options: map[string]string{
					"sslmode": "disable",
				},
			},
		},
		Migrations: getEnv("DB_MIGRATIONS", "database/migrations"),
		Redis: RedisConfig{
			Client: "predis",
			Default: RedisConnectionConfig{
				Host:     getEnv("REDIS_HOST", "127.0.0.1"),
				Password: getEnv("REDIS_PASSWORD", ""),
				Port:     getEnv("REDIS_PORT", "6379"),
				Database: getEnv("REDIS_DB", "0"),
			},
			Cache: RedisConnectionConfig{
				Host:     getEnv("REDIS_HOST", "127.0.0.1"),
				Password: getEnv("REDIS_PASSWORD", ""),
				Port:     getEnv("REDIS_PORT", "6379"),
				Database: getEnv("REDIS_CACHE_DB", "1"),
			},
			Queue: RedisConnectionConfig{
				Host:     getEnv("REDIS_HOST", "127.0.0.1"),
				Password: getEnv("REDIS_PASSWORD", ""),
				Port:     getEnv("REDIS_PORT", "6379"),
				Database: getEnv("REDIS_QUEUE_DB", "2"),
			},
		},
	}
}

// CacheConfig 缓存配置结构
type CacheConfig struct {
	Default string                 `json:"default"`
	Stores  map[string]StoreConfig `json:"stores"`
	Prefix  string                 `json:"prefix"`
}

// StoreConfig 存储配置结构
type StoreConfig struct {
	Driver     string                 `json:"driver"`
	Path       string                 `json:"path,omitempty"`
	Connection string                 `json:"connection,omitempty"`
	Table      string                 `json:"table,omitempty"`
	Options    map[string]interface{} `json:"options,omitempty"`
}

// LoadCacheConfig 加载缓存配置
func LoadCacheConfig() *CacheConfig {
	return &CacheConfig{
		Default: getEnv("CACHE_DRIVER", "file"),
		Stores: map[string]StoreConfig{
			"apc": {
				Driver: "apc",
			},
			"array": {
				Driver: "array",
				Options: map[string]interface{}{
					"serialize": false,
				},
			},
			"file": {
				Driver: "file",
				Path:   getEnv("CACHE_PATH", "storage/framework/cache/data"),
			},
			"redis": {
				Driver:     "redis",
				Connection: "cache",
			},
			"database": {
				Driver: "database",
				Table:  "cache",
			},
		},
		Prefix: getEnv("CACHE_PREFIX", "laravel_go_cache"),
	}
}

// QueueConfig 队列配置结构
type QueueConfig struct {
	Default     string                     `json:"default"`
	Connections map[string]QueueConnection `json:"connections"`
	Failed      FailedJobConfig            `json:"failed"`
}

// QueueConnection 队列连接配置
type QueueConnection struct {
	Driver       string        `json:"driver"`
	Table        string        `json:"table,omitempty"`
	Queue        string        `json:"queue,omitempty"`
	RetryAfter   int           `json:"retry_after,omitempty"`
	AfterCommit  bool          `json:"after_commit,omitempty"`
	Connection   string        `json:"connection,omitempty"`
	BlockFor     *int          `json:"block_for,omitempty"`
	Host         string        `json:"host,omitempty"`
	Key          string        `json:"key,omitempty"`
	Secret       string        `json:"secret,omitempty"`
	Prefix       string        `json:"prefix,omitempty"`
	Suffix       string        `json:"suffix,omitempty"`
	Region       string        `json:"region,omitempty"`
	Group        string        `json:"group,omitempty"`
	PersistentID string        `json:"persistent_id,omitempty"`
	SASL         []string      `json:"sasl,omitempty"`
	Servers      []QueueServer `json:"servers,omitempty"`
}

// QueueServer 队列服务器配置
type QueueServer struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	Weight int    `json:"weight"`
}

// FailedJobConfig 失败任务配置
type FailedJobConfig struct {
	Driver   string `json:"driver"`
	Database string `json:"database,omitempty"`
	Table    string `json:"table,omitempty"`
}

// LoadQueueConfig 加载队列配置
func LoadQueueConfig() *QueueConfig {
	return &QueueConfig{
		Default: getEnv("QUEUE_CONNECTION", "sync"),
		Connections: map[string]QueueConnection{
			"sync": {
				Driver: "sync",
			},
			"database": {
				Driver:      "database",
				Table:       "jobs",
				Queue:       "default",
				RetryAfter:  90,
				AfterCommit: false,
			},
			"redis": {
				Driver:      "redis",
				Connection:  "default",
				Queue:       getEnv("REDIS_QUEUE", "default"),
				RetryAfter:  90,
				AfterCommit: false,
			},
		},
		Failed: FailedJobConfig{
			Driver:   "database-uuids",
			Database: "pgsql",
			Table:    "failed_jobs",
		},
	}
}

// SessionConfig 会话配置结构
type SessionConfig struct {
	Driver        string `json:"driver"`
	Lifetime      int    `json:"lifetime"`
	ExpireOnClose bool   `json:"expire_on_close"`
	Encrypt       bool   `json:"encrypt"`
	Files         string `json:"files"`
	Connection    string `json:"connection"`
	Table         string `json:"table"`
	Store         string `json:"store"`
	Lottery       []int  `json:"lottery"`
	Cookie        string `json:"cookie"`
	Path          string `json:"path"`
	Domain        string `json:"domain"`
	Secure        bool   `json:"secure"`
	HTTPOnly      bool   `json:"http_only"`
	SameSite      string `json:"same_site"`
}

// LoadSessionConfig 加载会话配置
func LoadSessionConfig() *SessionConfig {
	return &SessionConfig{
		Driver:        getEnv("SESSION_DRIVER", "file"),
		Lifetime:      getEnvInt("SESSION_LIFETIME", 120),
		ExpireOnClose: false,
		Encrypt:       false,
		Files:         getEnv("SESSION_FILES", "storage/framework/sessions"),
		Connection:    getEnv("SESSION_CONNECTION", ""),
		Table:         "sessions",
		Store:         getEnv("SESSION_STORE", ""),
		Lottery:       []int{2, 100},
		Cookie:        getEnv("SESSION_COOKIE", "laravel_go_session"),
		Path:          "/",
		Domain:        getEnv("SESSION_DOMAIN", ""),
		Secure:        getEnvBool("SESSION_SECURE_COOKIE", false),
		HTTPOnly:      true,
		SameSite:      "lax",
	}
}

// LoggingConfig 日志配置结构
type LoggingConfig struct {
	Default      string                   `json:"default"`
	Deprecations DeprecationConfig        `json:"deprecations"`
	Channels     map[string]ChannelConfig `json:"channels"`
}

// DeprecationConfig 弃用配置
type DeprecationConfig struct {
	Channel string `json:"channel"`
	Trace   bool   `json:"trace"`
}

// ChannelConfig 通道配置
type ChannelConfig struct {
	Driver           string                 `json:"driver"`
	Path             string                 `json:"path,omitempty"`
	Level            string                 `json:"level,omitempty"`
	Days             int                    `json:"days,omitempty"`
	URL              string                 `json:"url,omitempty"`
	Username         string                 `json:"username,omitempty"`
	Emoji            string                 `json:"emoji,omitempty"`
	Handler          string                 `json:"handler,omitempty"`
	HandlerWith      map[string]interface{} `json:"handler_with,omitempty"`
	Formatter        string                 `json:"formatter,omitempty"`
	With             map[string]interface{} `json:"with,omitempty"`
	Channels         []string               `json:"channels,omitempty"`
	IgnoreExceptions bool                   `json:"ignore_exceptions,omitempty"`
}

// LoadLoggingConfig 加载日志配置
func LoadLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		Default: getEnv("LOG_CHANNEL", "stack"),
		Deprecations: DeprecationConfig{
			Channel: "null",
			Trace:   false,
		},
		Channels: map[string]ChannelConfig{
			"stack": {
				Driver:           "stack",
				Channels:         []string{"single"},
				IgnoreExceptions: false,
			},
			"single": {
				Driver: "single",
				Path:   "storage/logs/laravel.log",
				Level:  getEnv("LOG_LEVEL", "debug"),
			},
			"daily": {
				Driver: "daily",
				Path:   "storage/logs/laravel.log",
				Level:  getEnv("LOG_LEVEL", "debug"),
				Days:   14,
			},
			"slack": {
				Driver:   "slack",
				URL:      getEnv("LOG_SLACK_WEBHOOK_URL", ""),
				Username: "Laravel Log",
				Emoji:    ":boom:",
				Level:    getEnv("LOG_LEVEL", "critical"),
			},
			"stderr": {
				Driver:    "monolog",
				Level:     getEnv("LOG_LEVEL", "debug"),
				Handler:   "Monolog\\Handler\\StreamHandler",
				Formatter: "Monolog\\Formatter\\JsonFormatter",
				With: map[string]interface{}{
					"stream": "php://stderr",
				},
			},
			"syslog": {
				Driver: "syslog",
				Level:  getEnv("LOG_LEVEL", "debug"),
			},
			"errorlog": {
				Driver: "errorlog",
				Level:  getEnv("LOG_LEVEL", "debug"),
			},
			"null": {
				Driver:  "monolog",
				Handler: "Monolog\\Handler\\NullHandler",
			},
			"emergency": {
				Path: "storage/logs/laravel.log",
			},
		},
	}
}

// 辅助函数
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value == "true" || value == "1" {
			return true
		}
		return false
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
