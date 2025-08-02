package tests

import (
	"os"
	"strconv"
	"time"

	"laravel-go/framework/config"
)

// TestConfig 测试配置
type TestConfig struct {
	// 数据库配置
	Database struct {
		Driver   string `json:"driver"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Database string `json:"database"`
		Username string `json:"username"`
		Password string `json:"password"`
		Charset  string `json:"charset"`
	} `json:"database"`

	// 缓存配置
	Cache struct {
		Driver string `json:"driver"`
		Host   string `json:"host"`
		Port   int    `json:"port"`
		DB     int    `json:"db"`
	} `json:"cache"`

	// 队列配置
	Queue struct {
		Driver string `json:"driver"`
		Host   string `json:"host"`
		Port   int    `json:"port"`
		DB     int    `json:"db"`
	} `json:"queue"`

	// HTTP配置
	HTTP struct {
		Port    int           `json:"port"`
		Timeout time.Duration `json:"timeout"`
	} `json:"http"`

	// 测试配置
	Test struct {
		Parallel    bool          `json:"parallel"`
		Timeout     time.Duration `json:"timeout"`
		RetryCount  int           `json:"retry_count"`
		CleanupData bool          `json:"cleanup_data"`
	} `json:"test"`

	// 日志配置
	Log struct {
		Level  string `json:"level"`
		Output string `json:"output"`
	} `json:"log"`
}

// LoadTestConfig 加载测试配置
func LoadTestConfig() *TestConfig {
	cfg := &TestConfig{}

	// 设置默认值
	cfg.setDefaults()

	// 从环境变量加载配置
	cfg.loadFromEnv()

	// 从配置文件加载（如果存在）
	cfg.loadFromFile()

	return cfg
}

// setDefaults 设置默认配置
func (cfg *TestConfig) setDefaults() {
	// 数据库默认配置
	cfg.Database.Driver = "sqlite"
	cfg.Database.Database = ":memory:"
	cfg.Database.Charset = "utf8mb4"

	// 缓存默认配置
	cfg.Cache.Driver = "memory"
	cfg.Cache.Host = "localhost"
	cfg.Cache.Port = 6379
	cfg.Cache.DB = 0

	// 队列默认配置
	cfg.Queue.Driver = "memory"
	cfg.Queue.Host = "localhost"
	cfg.Queue.Port = 6379
	cfg.Queue.DB = 1

	// HTTP默认配置
	cfg.HTTP.Port = 8080
	cfg.HTTP.Timeout = 30 * time.Second

	// 测试默认配置
	cfg.Test.Parallel = false
	cfg.Test.Timeout = 60 * time.Second
	cfg.Test.RetryCount = 3
	cfg.Test.CleanupData = true

	// 日志默认配置
	cfg.Log.Level = "info"
	cfg.Log.Output = "stdout"
}

// loadFromEnv 从环境变量加载配置
func (cfg *TestConfig) loadFromEnv() {
	// 数据库配置
	if driver := os.Getenv("TEST_DB_DRIVER"); driver != "" {
		cfg.Database.Driver = driver
	}
	if host := os.Getenv("TEST_DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("TEST_DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.Port = p
		}
	}
	if database := os.Getenv("TEST_DB_DATABASE"); database != "" {
		cfg.Database.Database = database
	}
	if username := os.Getenv("TEST_DB_USERNAME"); username != "" {
		cfg.Database.Username = username
	}
	if password := os.Getenv("TEST_DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}

	// 缓存配置
	if driver := os.Getenv("TEST_CACHE_DRIVER"); driver != "" {
		cfg.Cache.Driver = driver
	}
	if host := os.Getenv("TEST_CACHE_HOST"); host != "" {
		cfg.Cache.Host = host
	}
	if port := os.Getenv("TEST_CACHE_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Cache.Port = p
		}
	}
	if db := os.Getenv("TEST_CACHE_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			cfg.Cache.DB = d
		}
	}

	// 队列配置
	if driver := os.Getenv("TEST_QUEUE_DRIVER"); driver != "" {
		cfg.Queue.Driver = driver
	}
	if host := os.Getenv("TEST_QUEUE_HOST"); host != "" {
		cfg.Queue.Host = host
	}
	if port := os.Getenv("TEST_QUEUE_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Queue.Port = p
		}
	}
	if db := os.Getenv("TEST_QUEUE_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			cfg.Queue.DB = d
		}
	}

	// HTTP配置
	if port := os.Getenv("TEST_HTTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.HTTP.Port = p
		}
	}
	if timeout := os.Getenv("TEST_HTTP_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			cfg.HTTP.Timeout = t
		}
	}

	// 测试配置
	if parallel := os.Getenv("TEST_PARALLEL"); parallel != "" {
		cfg.Test.Parallel = parallel == "true"
	}
	if timeout := os.Getenv("TEST_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			cfg.Test.Timeout = t
		}
	}
	if retryCount := os.Getenv("TEST_RETRY_COUNT"); retryCount != "" {
		if r, err := strconv.Atoi(retryCount); err == nil {
			cfg.Test.RetryCount = r
		}
	}
	if cleanupData := os.Getenv("TEST_CLEANUP_DATA"); cleanupData != "" {
		cfg.Test.CleanupData = cleanupData == "true"
	}

	// 日志配置
	if level := os.Getenv("TEST_LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}
	if output := os.Getenv("TEST_LOG_OUTPUT"); output != "" {
		cfg.Log.Output = output
	}
}

// loadFromFile 从配置文件加载配置
func (cfg *TestConfig) loadFromFile() {
	// 检查是否存在测试配置文件
	configFile := "tests/test_config.json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return
	}

	// 加载配置文件
	config := config.NewConfig()
	if err := config.LoadFile(configFile); err != nil {
		return
	}

	// 从配置文件读取值（这里简化处理，实际应该解析JSON）
	// 在实际实现中，应该使用JSON解析器来读取配置文件
}

// GetDatabaseDSN 获取数据库连接字符串
func (cfg *TestConfig) GetDatabaseDSN() string {
	switch cfg.Database.Driver {
	case "sqlite":
		return cfg.Database.Database
	case "mysql":
		return cfg.Database.Username + ":" + cfg.Database.Password + "@tcp(" + cfg.Database.Host + ":" + strconv.Itoa(cfg.Database.Port) + ")/" + cfg.Database.Database + "?charset=" + cfg.Database.Charset
	case "postgres":
		return "host=" + cfg.Database.Host + " port=" + strconv.Itoa(cfg.Database.Port) + " user=" + cfg.Database.Username + " password=" + cfg.Database.Password + " dbname=" + cfg.Database.Database + " sslmode=disable"
	default:
		return cfg.Database.Database
	}
}

// GetCacheDSN 获取缓存连接字符串
func (cfg *TestConfig) GetCacheDSN() string {
	switch cfg.Cache.Driver {
	case "redis":
		return cfg.Cache.Host + ":" + strconv.Itoa(cfg.Cache.Port)
	case "memory":
		return "memory"
	default:
		return "memory"
	}
}

// GetQueueDSN 获取队列连接字符串
func (cfg *TestConfig) GetQueueDSN() string {
	switch cfg.Queue.Driver {
	case "redis":
		return cfg.Queue.Host + ":" + strconv.Itoa(cfg.Queue.Port)
	case "memory":
		return "memory"
	default:
		return "memory"
	}
}

// IsParallel 是否并行运行测试
func (cfg *TestConfig) IsParallel() bool {
	return cfg.Test.Parallel
}

// GetTimeout 获取测试超时时间
func (cfg *TestConfig) GetTimeout() time.Duration {
	return cfg.Test.Timeout
}

// GetRetryCount 获取重试次数
func (cfg *TestConfig) GetRetryCount() int {
	return cfg.Test.RetryCount
}

// ShouldCleanupData 是否清理测试数据
func (cfg *TestConfig) ShouldCleanupData() bool {
	return cfg.Test.CleanupData
}

// GetLogLevel 获取日志级别
func (cfg *TestConfig) GetLogLevel() string {
	return cfg.Log.Level
}

// GetLogOutput 获取日志输出
func (cfg *TestConfig) GetLogOutput() string {
	return cfg.Log.Output
}

// 全局测试配置实例
var GlobalTestConfig *TestConfig

// init 初始化全局测试配置
func init() {
	GlobalTestConfig = LoadTestConfig()
}
