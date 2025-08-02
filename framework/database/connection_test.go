package database

import (
	"testing"
	"time"

	"laravel-go/framework/config"
)

func TestNewConnection(t *testing.T) {
	// 测试 SQLite 连接（不需要真实数据库）
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db", // 使用文件数据库
	}

	// 打印 DSN 用于调试
	dsn, _ := buildDSN(config)
	t.Logf("DSN: %s", dsn)

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 测试 Ping
	if err := conn.Ping(); err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	// 测试查询
	rows, err := conn.Query("SELECT 1")
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}
	defer rows.Close()

	// 测试统计信息
	stats := conn.Stats()
	// SQLite 的默认 MaxOpenConnections 是 0，这是正常的
	if stats.MaxOpenConnections < 0 {
		t.Error("MaxOpenConnections should not be negative")
	}
}

func TestConnectionManager(t *testing.T) {
	manager := NewConnectionManager()

	// 添加连接配置
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	manager.AddConnection("test", config)

	// 获取连接
	conn, err := manager.GetConnection("test")
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}

	// 测试连接
	if err := conn.Ping(); err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	// 测试重复获取（应该返回同一个连接）
	conn2, err := manager.GetConnection("test")
	if err != nil {
		t.Fatalf("Failed to get connection again: %v", err)
	}

	if conn != conn2 {
		t.Error("Should return the same connection instance")
	}

	// 测试获取不存在的连接
	_, err = manager.GetConnection("nonexistent")
	if err == nil {
		t.Error("Should return error for nonexistent connection")
	}

	// 测试统计信息
	stats := manager.GetStats()
	if len(stats) != 1 {
		t.Errorf("Expected 1 connection, got %d", len(stats))
	}

	// 测试关闭连接
	if err := manager.CloseConnection("test"); err != nil {
		t.Errorf("Failed to close connection: %v", err)
	}

	// 测试关闭所有连接
	manager.AddConnection("test2", config)
	conn3, _ := manager.GetConnection("test2")
	if err := manager.CloseAll(); err != nil {
		t.Errorf("Failed to close all connections: %v", err)
	}

	// 验证连接已关闭
	if err := conn3.Ping(); err == nil {
		t.Error("Connection should be closed")
	}
}

func TestBuildDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   *ConnectionConfig
		expected string
	}{
		{
			name: "MySQL DSN",
			config: &ConnectionConfig{
				Driver:   MySQL,
				Host:     "localhost",
				Port:     3306,
				Database: "testdb",
				Username: "user",
				Password: "pass",
				Charset:  "utf8mb4",
			},
			expected: "user:pass@tcp(localhost:3306)/testdb?charset=utf8mb4",
		},
		{
			name: "PostgreSQL DSN",
			config: &ConnectionConfig{
				Driver:   PostgreSQL,
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				Username: "user",
				Password: "pass",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=user password=pass dbname=testdb sslmode=disable",
		},
		{
			name: "SQLite DSN",
			config: &ConnectionConfig{
				Driver:   SQLite,
				Database: "test.db",
			},
			expected: "test.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn, err := buildDSN(tt.config)
			if err != nil {
				t.Fatalf("buildDSN failed: %v", err)
			}
			if dsn != tt.expected {
				t.Errorf("Expected DSN %s, got %s", tt.expected, dsn)
			}
		})
	}
}

func TestParseConnectionConfig(t *testing.T) {
	data := map[string]interface{}{
		"driver":             "mysql",
		"host":               "localhost",
		"port":               3306,
		"database":           "testdb",
		"username":           "user",
		"password":           "pass",
		"charset":            "utf8mb4",
		"timezone":           "+08:00",
		"max_open_conns":     10,
		"max_idle_conns":     5,
		"conn_max_lifetime":  "1h",
		"conn_max_idle_time": "30m",
		"ssl_mode":           "disable",
		"options": map[string]interface{}{
			"parseTime": "true",
			"loc":       "Local",
		},
	}

	config, err := parseConnectionConfig(data)
	if err != nil {
		t.Fatalf("parseConnectionConfig failed: %v", err)
	}

	if config.Driver != MySQL {
		t.Errorf("Expected driver MySQL, got %s", config.Driver)
	}
	if config.Host != "localhost" {
		t.Errorf("Expected host localhost, got %s", config.Host)
	}
	if config.Port != 3306 {
		t.Errorf("Expected port 3306, got %d", config.Port)
	}
	if config.Database != "testdb" {
		t.Errorf("Expected database testdb, got %s", config.Database)
	}
	if config.Username != "user" {
		t.Errorf("Expected username user, got %s", config.Username)
	}
	if config.Password != "pass" {
		t.Errorf("Expected password pass, got %s", config.Password)
	}
	if config.Charset != "utf8mb4" {
		t.Errorf("Expected charset utf8mb4, got %s", config.Charset)
	}
	if config.Timezone != "+08:00" {
		t.Errorf("Expected timezone +08:00, got %s", config.Timezone)
	}
	if config.MaxOpenConns != 10 {
		t.Errorf("Expected max_open_conns 10, got %d", config.MaxOpenConns)
	}
	if config.MaxIdleConns != 5 {
		t.Errorf("Expected max_idle_conns 5, got %d", config.MaxIdleConns)
	}
	if config.ConnMaxLifetime != time.Hour {
		t.Errorf("Expected conn_max_lifetime 1h, got %v", config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime != 30*time.Minute {
		t.Errorf("Expected conn_max_idle_time 30m, got %v", config.ConnMaxIdleTime)
	}
	if config.SSLMode != "disable" {
		t.Errorf("Expected ssl_mode disable, got %s", config.SSLMode)
	}
	if config.Options["parseTime"] != "true" {
		t.Errorf("Expected option parseTime true, got %s", config.Options["parseTime"])
	}
	if config.Options["loc"] != "Local" {
		t.Errorf("Expected option loc Local, got %s", config.Options["loc"])
	}
}

func TestLoadFromConfig(t *testing.T) {
	// 创建测试配置
	cfg := config.NewConfig()
	cfg.Set("database.connections.mysql", map[string]interface{}{
		"driver":   "mysql",
		"host":     "localhost",
		"port":     3306,
		"database": "testdb",
		"username": "user",
		"password": "pass",
	})
	cfg.Set("database.connections.sqlite", map[string]interface{}{
		"driver":   "sqlite",
		"database": "test.db",
	})

	manager := NewConnectionManager()
	err := manager.LoadFromConfig(cfg)
	if err != nil {
		t.Fatalf("LoadFromConfig failed: %v", err)
	}

	// 测试获取连接
	conn, err := manager.GetConnection("sqlite")
	if err != nil {
		t.Fatalf("Failed to get sqlite connection: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		t.Errorf("Ping failed: %v", err)
	}
}

func TestConnectionPoolSettings(t *testing.T) {
	config := &ConnectionConfig{
		Driver:          SQLite,
		Database:        "test.db",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 验证连接池设置
	db := conn.DB()
	stats := db.Stats()

	if stats.MaxOpenConnections != 10 {
		t.Errorf("Expected MaxOpenConnections 10, got %d", stats.MaxOpenConnections)
	}
}

func TestTransaction(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test_transaction.db",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// 测试事务
	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// 插入数据
	_, err = tx.Exec("INSERT OR REPLACE INTO test (id, name) VALUES (?, ?)", 1, "test")
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	// 验证数据
	row := conn.QueryRow("SELECT name FROM test WHERE id = ?", 1)
	var name string
	if err := row.Scan(&name); err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}

	if name != "test" {
		t.Errorf("Expected name 'test', got '%s'", name)
	}
}

func BenchmarkConnectionCreation(b *testing.B) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := NewConnection(config)
		if err != nil {
			b.Fatalf("Failed to create connection: %v", err)
		}
		conn.Close()
	}
}

func BenchmarkConnectionManager(b *testing.B) {
	manager := NewConnectionManager()
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	manager.AddConnection("test", config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := manager.GetConnection("test")
		if err != nil {
			b.Fatalf("Failed to get connection: %v", err)
		}
		// 使用连接
		conn.Ping()
	}
}
