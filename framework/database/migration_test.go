package database

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMigration 测试迁移实现
type TestMigration struct {
	version     string
	name        string
	description string
	upSQL       string
	downSQL     string
}

func NewTestMigration(version, name, description, upSQL, downSQL string) *TestMigration {
	return &TestMigration{
		version:     version,
		name:        name,
		description: description,
		upSQL:       upSQL,
		downSQL:     downSQL,
	}
}

func (tm *TestMigration) GetName() string {
	return tm.name
}

func (tm *TestMigration) GetVersion() string {
	return tm.version
}

func (tm *TestMigration) GetDescription() string {
	return tm.description
}

func (tm *TestMigration) Up(conn Connection) error {
	if tm.upSQL == "" {
		return nil
	}
	_, err := conn.Exec(tm.upSQL)
	return err
}

func (tm *TestMigration) Down(conn Connection) error {
	if tm.downSQL == "" {
		return nil
	}
	_, err := conn.Exec(tm.downSQL)
	return err
}

// TestMigrationManager 测试迁移管理器
func TestMigrationManager(t *testing.T) {
	// 创建测试数据库连接
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: ":memory:",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建临时迁移目录
	tempDir := t.TempDir()

	// 创建迁移管理器
	manager := NewMigrationManager(conn, tempDir)

	// 测试创建迁移表
	err = manager.CreateMigrationTable()
	if err != nil {
		t.Fatalf("Failed to create migration table: %v", err)
	}

	// 验证表是否创建成功
	rows, err := conn.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='migrations'")
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Error("Migrations table was not created")
	}

	// 测试获取已执行的迁移
	executed, err := manager.GetExecutedMigrations()
	if err != nil {
		t.Fatalf("Failed to get executed migrations: %v", err)
	}

	if len(executed) != 0 {
		t.Error("Expected no executed migrations initially")
	}

	// 测试获取下一个批次号
	batch, err := manager.GetNextBatchNumber()
	if err != nil {
		t.Fatalf("Failed to get next batch number: %v", err)
	}

	if batch != 1 {
		t.Errorf("Expected batch number 1, got %d", batch)
	}
}

// TestMigrationRegistration 测试迁移注册
func TestMigrationRegistration(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: ":memory:",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	tempDir := t.TempDir()
	manager := NewMigrationManager(conn, tempDir)

	// 创建测试迁移
	migration1 := NewTestMigration("20240101000001", "create_users_table", "Create users table",
		"CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)",
		"DROP TABLE users")

	migration2 := NewTestMigration("20240101000002", "create_posts_table", "Create posts table",
		"CREATE TABLE posts (id INTEGER PRIMARY KEY, title TEXT, user_id INTEGER)",
		"DROP TABLE posts")

	// 注册迁移
	manager.RegisterMigration(migration1)
	manager.RegisterMigration(migration2)

	// 验证迁移已注册
	if len(manager.migrations) != 2 {
		t.Errorf("Expected 2 migrations, got %d", len(manager.migrations))
	}

	// 验证迁移内容
	if manager.migrations["20240101000001"] != migration1 {
		t.Error("Migration 1 not found")
	}

	if manager.migrations["20240101000002"] != migration2 {
		t.Error("Migration 2 not found")
	}
}

// TestMigrationExecution 测试迁移执行
func TestMigrationExecution(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: ":memory:",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	tempDir := t.TempDir()
	manager := NewMigrationManager(conn, tempDir)

	// 创建迁移表
	err = manager.CreateMigrationTable()
	if err != nil {
		t.Fatalf("Failed to create migration table: %v", err)
	}

	// 创建测试迁移
	migration := NewTestMigration("20240101000001", "create_test_table", "Create test table",
		"CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)",
		"DROP TABLE test_table")

	manager.RegisterMigration(migration)

	// 执行迁移
	err = manager.RunMigrations()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// 验证表已创建
	rows, err := conn.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='test_table'")
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Error("Test table was not created")
	}

	// 验证迁移记录
	executed, err := manager.GetExecutedMigrations()
	if err != nil {
		t.Fatalf("Failed to get executed migrations: %v", err)
	}

	if !executed["20240101000001"] {
		t.Error("Migration was not recorded as executed")
	}
}

// TestMigrationRollback 测试迁移回滚
func TestMigrationRollback(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: ":memory:",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	tempDir := t.TempDir()
	manager := NewMigrationManager(conn, tempDir)

	// 创建迁移表
	err = manager.CreateMigrationTable()
	if err != nil {
		t.Fatalf("Failed to create migration table: %v", err)
	}

	// 创建测试迁移
	migration := NewTestMigration("20240101000001", "create_test_table", "Create test table",
		"CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)",
		"DROP TABLE test_table")

	manager.RegisterMigration(migration)

	// 执行迁移
	err = manager.RunMigrations()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// 验证表已创建
	rows, err := conn.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='test_table'")
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Error("Test table was not created")
	}

	// 回滚迁移
	err = manager.RollbackMigrations(1)
	if err != nil {
		t.Fatalf("Failed to rollback migrations: %v", err)
	}

	// 验证表已删除
	rows, err = conn.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='test_table'")
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		t.Error("Test table was not dropped")
	}

	// 验证迁移记录已删除
	executed, err := manager.GetExecutedMigrations()
	if err != nil {
		t.Fatalf("Failed to get executed migrations: %v", err)
	}

	if executed["20240101000001"] {
		t.Error("Migration record was not deleted")
	}
}

// TestMigrationGenerator 测试迁移生成器
func TestMigrationGenerator(t *testing.T) {
	tempDir := t.TempDir()
	generator := NewMigrationGenerator(tempDir)

	// 生成迁移文件
	err := generator.GenerateMigration("create_users_table", "Create users table for authentication")
	if err != nil {
		t.Fatalf("Failed to generate migration: %v", err)
	}

	// 验证文件已创建
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	// 验证文件名格式
	filename := files[0].Name()
	if filepath.Ext(filename) != ".sql" {
		t.Error("Generated file is not a .sql file")
	}

	// 验证文件内容
	filePath := filepath.Join(tempDir, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "create_users_table") {
		t.Error("Generated file does not contain migration name")
	}

	if !strings.Contains(contentStr, "Create users table for authentication") {
		t.Error("Generated file does not contain description")
	}
}

// TestMigrationFileParsing 测试迁移文件解析
func TestMigrationFileParsing(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: ":memory:",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	tempDir := t.TempDir()
	manager := NewMigrationManager(conn, tempDir)

	// 创建测试迁移文件
	migrationContent := `-- Migration: create_test_table
-- Description: Create test table
-- Version: 20240101000001

-- UP Migration
CREATE TABLE test_table (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- DOWN Migration
DROP TABLE test_table;
`

	filename := "20240101000001_create_test_table.sql"
	filePath := filepath.Join(tempDir, filename)

	err = os.WriteFile(filePath, []byte(migrationContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	// 加载迁移
	err = manager.LoadMigrationsFromDirectory()
	if err != nil {
		t.Fatalf("Failed to load migrations: %v", err)
	}

	// 验证迁移已加载
	if len(manager.migrations) != 1 {
		t.Errorf("Expected 1 migration, got %d", len(manager.migrations))
	}

	migration := manager.migrations["20240101000001"]
	if migration == nil {
		t.Error("Migration not found")
	}

	if migration.GetName() != "create_test_table" {
		t.Errorf("Expected name 'create_test_table', got '%s'", migration.GetName())
	}

	// 执行迁移
	err = manager.CreateMigrationTable()
	if err != nil {
		t.Fatalf("Failed to create migration table: %v", err)
	}

	err = manager.RunMigrations()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// 验证表已创建
	rows, err := conn.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='test_table'")
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Error("Test table was not created")
	}
}

// TestMigrationStatus 测试迁移状态
func TestMigrationStatus(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: ":memory:",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	tempDir := t.TempDir()
	manager := NewMigrationManager(conn, tempDir)

	// 创建迁移表
	err = manager.CreateMigrationTable()
	if err != nil {
		t.Fatalf("Failed to create migration table: %v", err)
	}

	// 创建测试迁移
	migration := NewTestMigration("20240101000001", "create_test_table", "Create test table",
		"CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)",
		"DROP TABLE test_table")

	manager.RegisterMigration(migration)

	// 执行迁移
	err = manager.RunMigrations()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// 获取迁移状态
	status, err := manager.GetMigrationStatus()
	if err != nil {
		t.Fatalf("Failed to get migration status: %v", err)
	}

	if len(status) != 1 {
		t.Errorf("Expected 1 migration status, got %d", len(status))
	}

	ms := status[0]
	if ms.Version != "20240101000001" {
		t.Errorf("Expected version '20240101000001', got '%s'", ms.Version)
	}

	if ms.Name != "create_test_table" {
		t.Errorf("Expected name 'create_test_table', got '%s'", ms.Name)
	}

	if ms.Status != "executed" {
		t.Errorf("Expected status 'executed', got '%s'", ms.Status)
	}

	if ms.Batch != 1 {
		t.Errorf("Expected batch 1, got %d", ms.Batch)
	}
}

// BenchmarkMigrationExecution 基准测试迁移执行
func BenchmarkMigrationExecution(b *testing.B) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: ":memory:",
	}

	conn, err := NewConnection(config)
	if err != nil {
		b.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	tempDir := b.TempDir()
	manager := NewMigrationManager(conn, tempDir)

	// 创建迁移表
	err = manager.CreateMigrationTable()
	if err != nil {
		b.Fatalf("Failed to create migration table: %v", err)
	}

	// 创建测试迁移
	migration := NewTestMigration("20240101000001", "create_test_table", "Create test table",
		"CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)",
		"DROP TABLE test_table")

	manager.RegisterMigration(migration)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 执行迁移
		err = manager.RunMigrations()
		if err != nil {
			b.Fatalf("Failed to run migrations: %v", err)
		}

		// 回滚迁移
		err = manager.RollbackMigrations(1)
		if err != nil {
			b.Fatalf("Failed to rollback migrations: %v", err)
		}
	}
}
