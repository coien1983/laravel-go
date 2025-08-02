package config

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	if config == nil {
		t.Fatal("NewConfig() should not return nil")
	}
}

func TestSetAndGet(t *testing.T) {
	config := NewConfig()

	// 设置配置
	config.Set("app.name", "Laravel-Go")
	config.Set("app.debug", true)
	config.Set("database.host", "localhost")

	// 获取配置
	if name := config.GetString("app.name"); name != "Laravel-Go" {
		t.Fatalf("Expected 'Laravel-Go', got '%s'", name)
	}

	if debug := config.GetBool("app.debug"); !debug {
		t.Fatal("Expected true, got false")
	}

	if host := config.GetString("database.host"); host != "localhost" {
		t.Fatalf("Expected 'localhost', got '%s'", host)
	}
}

func TestGetWithDefault(t *testing.T) {
	config := NewConfig()

	// 获取不存在的配置，使用默认值
	if name := config.GetString("app.name", "default"); name != "default" {
		t.Fatalf("Expected 'default', got '%s'", name)
	}

	if port := config.GetInt("app.port", 8080); port != 8080 {
		t.Fatalf("Expected 8080, got %d", port)
	}

	if debug := config.GetBool("app.debug", true); !debug {
		t.Fatal("Expected true, got false")
	}
}

func TestHas(t *testing.T) {
	config := NewConfig()

	// 检查不存在的配置
	if config.Has("app.name") {
		t.Fatal("Has() should return false for non-existent key")
	}

	// 设置配置
	config.Set("app.name", "test")

	// 检查存在的配置
	if !config.Has("app.name") {
		t.Fatal("Has() should return true for existing key")
	}
}

func TestLoadFromStruct(t *testing.T) {
	config := NewConfig()

	// 设置环境变量
	os.Setenv("APP_NAME", "TestApp")
	os.Setenv("APP_DEBUG", "true")

	// 定义配置结构体
	type AppConfig struct {
		Name  string `env:"APP_NAME" default:"Laravel-Go"`
		Debug bool   `env:"APP_DEBUG" default:"false"`
		Port  int    `env:"APP_PORT" default:"8080"`
	}

	appConfig := &AppConfig{}

	// 从结构体加载配置
	if err := config.LoadFromStruct(appConfig); err != nil {
		t.Fatalf("LoadFromStruct() should not return error: %v", err)
	}

	// 验证配置
	if name := config.GetString("Name"); name != "TestApp" {
		t.Fatalf("Expected 'TestApp', got '%s'", name)
	}

	if debug := config.GetBool("Debug"); !debug {
		t.Fatal("Expected true, got false")
	}

	if port := config.GetInt("Port"); port != 8080 {
		t.Fatalf("Expected 8080, got %d", port)
	}
}

func TestLoadEnv(t *testing.T) {
	config := NewConfig()

	// 创建临时环境文件
	envContent := `APP_NAME=Laravel-Go
APP_DEBUG=true
DB_HOST=localhost
DB_PORT=3306`

	tmpFile, err := os.CreateTemp("", "test.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(envContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// 加载环境文件
	if err := config.LoadEnv(tmpFile.Name()); err != nil {
		t.Fatalf("LoadEnv() should not return error: %v", err)
	}

	// 验证环境变量
	if name := os.Getenv("APP_NAME"); name != "Laravel-Go" {
		t.Fatalf("Expected 'Laravel-Go', got '%s'", name)
	}

	if debug := os.Getenv("APP_DEBUG"); debug != "true" {
		t.Fatalf("Expected 'true', got '%s'", debug)
	}
}

func TestGetStringSlice(t *testing.T) {
	config := NewConfig()

	// 设置字符串切片
	config.Set("app.middleware", []string{"auth", "cors", "log"})

	// 获取字符串切片
	middleware := config.GetStringSlice("app.middleware")
	if len(middleware) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(middleware))
	}

	expected := []string{"auth", "cors", "log"}
	for i, item := range middleware {
		if item != expected[i] {
			t.Fatalf("Expected '%s', got '%s'", expected[i], item)
		}
	}
}

func TestGetMap(t *testing.T) {
	config := NewConfig()

	// 设置映射
	config.Set("database", map[string]interface{}{
		"host": "localhost",
		"port": 3306,
		"name": "laravel_go",
	})

	// 获取映射
	dbConfig := config.GetMap("database")
	if len(dbConfig) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(dbConfig))
	}

	if host := dbConfig["host"]; host != "localhost" {
		t.Fatalf("Expected 'localhost', got '%v'", host)
	}

	if port := dbConfig["port"]; port != 3306 {
		t.Fatalf("Expected 3306, got '%v'", port)
	}
}

func TestValidate(t *testing.T) {
	config := NewConfig()

	// 设置配置
	config.Set("app.name", "Laravel-Go")
	config.Set("app.port", 8080)
	config.Set("app.debug", true)

	// 定义验证规则
	rules := map[string]string{
		"app.name":  "required|string",
		"app.port":  "required|int",
		"app.debug": "required|bool",
		"app.host":  "required", // 这个不存在，应该失败
	}

	// 验证配置
	if err := config.Validate(rules); err == nil {
		t.Fatal("Validate() should return error for missing required field")
	}

	// 验证存在的字段
	validRules := map[string]string{
		"app.name":  "required|string",
		"app.port":  "required|int",
		"app.debug": "required|bool",
	}

	if err := config.Validate(validRules); err != nil {
		t.Fatalf("Validate() should not return error: %v", err)
	}
}
