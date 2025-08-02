package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"laravel-go/framework/errors"
)

// Migration 迁移接口
type Migration interface {
	// 获取迁移信息
	GetName() string
	GetVersion() string
	GetDescription() string

	// 执行迁移
	Up(conn Connection) error
	Down(conn Connection) error
}

// MigrationFile 迁移文件结构
type MigrationFile struct {
	Version     string
	Name        string
	Description string
	UpSQL       string
	DownSQL     string
	CreatedAt   time.Time
}

// MigrationManager 迁移管理器
type MigrationManager struct {
	conn          Connection
	migrations    map[string]Migration
	migrationsDir string
}

// NewMigrationManager 创建新的迁移管理器
func NewMigrationManager(conn Connection, migrationsDir string) *MigrationManager {
	return &MigrationManager{
		conn:          conn,
		migrations:    make(map[string]Migration),
		migrationsDir: migrationsDir,
	}
}

// RegisterMigration 注册迁移
func (mm *MigrationManager) RegisterMigration(migration Migration) {
	mm.migrations[migration.GetVersion()] = migration
}

// LoadMigrationsFromDirectory 从目录加载迁移文件
func (mm *MigrationManager) LoadMigrationsFromDirectory() error {
	// 确保目录存在
	if err := os.MkdirAll(mm.migrationsDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create migrations directory")
	}

	// 读取目录中的所有 .sql 文件
	files, err := os.ReadDir(mm.migrationsDir)
	if err != nil {
		return errors.Wrap(err, "failed to read migrations directory")
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		migration, err := mm.loadMigrationFromFile(file.Name())
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to load migration %s", file.Name()))
		}

		mm.RegisterMigration(migration)
	}

	return nil
}

// loadMigrationFromFile 从文件加载迁移
func (mm *MigrationManager) loadMigrationFromFile(filename string) (Migration, error) {
	filePath := filepath.Join(mm.migrationsDir, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read migration file")
	}

	// 解析文件名获取版本和名称
	// 格式: 20240101000000_create_users_table.sql
	parts := strings.Split(strings.TrimSuffix(filename, ".sql"), "_")
	if len(parts) < 2 {
		return nil, errors.New("invalid migration filename format")
	}

	version := parts[0]
	name := strings.Join(parts[1:], "_")

	// 解析SQL内容
	upSQL, downSQL, err := mm.parseMigrationSQL(string(content))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse migration SQL")
	}

	return &fileMigration{
		version:     version,
		name:        name,
		description: fmt.Sprintf("Migration %s", name),
		upSQL:       upSQL,
		downSQL:     downSQL,
		createdAt:   time.Now(),
	}, nil
}

// parseMigrationSQL 解析迁移SQL文件
func (mm *MigrationManager) parseMigrationSQL(content string) (string, string, error) {
	// 简单的SQL解析，查找 -- UP 和 -- DOWN 标记
	lines := strings.Split(content, "\n")
	var upLines, downLines []string
	var currentSection string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "-- UP") {
			currentSection = "up"
			continue
		}

		if strings.HasPrefix(line, "-- DOWN") {
			currentSection = "down"
			continue
		}

		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		switch currentSection {
		case "up":
			upLines = append(upLines, line)
		case "down":
			downLines = append(downLines, line)
		}
	}

	upSQL := strings.Join(upLines, "\n")
	downSQL := strings.Join(downLines, "\n")

	return upSQL, downSQL, nil
}

// CreateMigrationTable 创建迁移表
func (mm *MigrationManager) CreateMigrationTable() error {
	query := `CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		description TEXT,
		executed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		batch INTEGER NOT NULL
	)`

	_, err := mm.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	return nil
}

// GetExecutedMigrations 获取已执行的迁移
func (mm *MigrationManager) GetExecutedMigrations() (map[string]bool, error) {
	query := "SELECT version FROM migrations ORDER BY version"
	rows, err := mm.conn.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query executed migrations")
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, errors.Wrap(err, "failed to scan migration version")
		}
		executed[version] = true
	}

	return executed, nil
}

// GetNextBatchNumber 获取下一个批次号
func (mm *MigrationManager) GetNextBatchNumber() (int, error) {
	query := "SELECT COALESCE(MAX(batch), 0) + 1 FROM migrations"
	row := mm.conn.QueryRow(query)

	var batch int
	err := row.Scan(&batch)
	if err != nil && err != sql.ErrNoRows {
		return 0, errors.Wrap(err, "failed to get next batch number")
	}

	return batch, nil
}

// RunMigrations 运行迁移
func (mm *MigrationManager) RunMigrations() error {
	// 确保迁移表存在
	if err := mm.CreateMigrationTable(); err != nil {
		return err
	}

	// 获取已执行的迁移
	executed, err := mm.GetExecutedMigrations()
	if err != nil {
		return err
	}

	// 获取下一个批次号
	batch, err := mm.GetNextBatchNumber()
	if err != nil {
		return err
	}

	// 获取所有迁移版本并排序
	var versions []string
	for version := range mm.migrations {
		versions = append(versions, version)
	}
	sort.Strings(versions)

	// 执行未执行的迁移
	for _, version := range versions {
		if executed[version] {
			continue
		}

		migration := mm.migrations[version]
		fmt.Printf("Running migration: %s (%s)\n", migration.GetName(), version)

		// 开始事务
		tx, err := mm.conn.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// 执行迁移
		if err := migration.Up(mm.conn); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to run migration %s: %w", version, err)
		}

		// 记录迁移执行
		insertQuery := `
			INSERT INTO migrations (version, name, description, batch)
			VALUES (?, ?, ?, ?)
		`
		_, err = mm.conn.Exec(insertQuery, version, migration.GetName(), migration.GetDescription(), batch)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration: %w", err)
		}

		// 提交事务
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		fmt.Printf("✓ Migration completed: %s\n", migration.GetName())
	}

	return nil
}

// RollbackMigrations 回滚迁移
func (mm *MigrationManager) RollbackMigrations(steps int) error {
	// 获取最后一批次的迁移
	query := `
		SELECT version, name FROM migrations 
		WHERE batch = (SELECT MAX(batch) FROM migrations)
		ORDER BY version DESC
		LIMIT ?
	`
	rows, err := mm.conn.Query(query, steps)
	if err != nil {
		return errors.Wrap(err, "failed to query migrations to rollback")
	}
	defer rows.Close()

	var migrationsToRollback []struct {
		version string
		name    string
	}

	for rows.Next() {
		var version, name string
		if err := rows.Scan(&version, &name); err != nil {
			return errors.Wrap(err, "failed to scan migration")
		}
		migrationsToRollback = append(migrationsToRollback, struct {
			version string
			name    string
		}{version, name})
	}

	// 回滚迁移
	for _, migrationInfo := range migrationsToRollback {
		migration, exists := mm.migrations[migrationInfo.version]
		if !exists {
			fmt.Printf("Warning: Migration %s not found, skipping\n", migrationInfo.version)
			continue
		}

		fmt.Printf("Rolling back migration: %s (%s)\n", migrationInfo.name, migrationInfo.version)

		// 先删除迁移记录
		deleteQuery := "DELETE FROM migrations WHERE version = ?"
		_, err = mm.conn.Exec(deleteQuery, migrationInfo.version)
		if err != nil {
			return fmt.Errorf("failed to delete migration record: %w", err)
		}

		// 执行回滚
		if err := migration.Down(mm.conn); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", migrationInfo.version, err)
		}

		fmt.Printf("✓ Migration rolled back: %s\n", migrationInfo.name)
	}

	return nil
}

// GetMigrationStatus 获取迁移状态
func (mm *MigrationManager) GetMigrationStatus() ([]MigrationStatus, error) {
	query := `
		SELECT m.version, m.name, m.description, m.executed_at, m.batch,
		       CASE WHEN m.version IS NOT NULL THEN 'executed' ELSE 'pending' END as status
		FROM (
			SELECT version, name, description, executed_at, batch FROM migrations
		) m
		ORDER BY m.version
	`
	rows, err := mm.conn.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query migration status")
	}
	defer rows.Close()

	var status []MigrationStatus
	for rows.Next() {
		var ms MigrationStatus
		err := rows.Scan(&ms.Version, &ms.Name, &ms.Description, &ms.ExecutedAt, &ms.Batch, &ms.Status)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan migration status")
		}
		status = append(status, ms)
	}

	return status, nil
}

// MigrationStatus 迁移状态
type MigrationStatus struct {
	Version     string
	Name        string
	Description string
	ExecutedAt  *time.Time
	Batch       int
	Status      string
}

// fileMigration 文件迁移实现
type fileMigration struct {
	version     string
	name        string
	description string
	upSQL       string
	downSQL     string
	createdAt   time.Time
}

func (fm *fileMigration) GetName() string {
	return fm.name
}

func (fm *fileMigration) GetVersion() string {
	return fm.version
}

func (fm *fileMigration) GetDescription() string {
	return fm.description
}

func (fm *fileMigration) Up(conn Connection) error {
	if fm.upSQL == "" {
		return errors.New("no UP SQL found in migration")
	}

	_, err := conn.Exec(fm.upSQL)
	if err != nil {
		return fmt.Errorf("failed to execute UP migration: %w", err)
	}
	return nil
}

func (fm *fileMigration) Down(conn Connection) error {
	if fm.downSQL == "" {
		return errors.New("no DOWN SQL found in migration")
	}

	_, err := conn.Exec(fm.downSQL)
	if err != nil {
		return fmt.Errorf("failed to execute DOWN migration: %w", err)
	}
	return nil
}

// MigrationGenerator 迁移生成器
type MigrationGenerator struct {
	migrationsDir string
}

// NewMigrationGenerator 创建新的迁移生成器
func NewMigrationGenerator(migrationsDir string) *MigrationGenerator {
	return &MigrationGenerator{
		migrationsDir: migrationsDir,
	}
}

// GenerateMigration 生成迁移文件
func (mg *MigrationGenerator) GenerateMigration(name, description string) error {
	// 生成版本号（时间戳格式）
	version := time.Now().Format("20060102150405")

	// 生成文件名
	filename := fmt.Sprintf("%s_%s.sql", version, strings.ToLower(strings.ReplaceAll(name, " ", "_")))
	filePath := filepath.Join(mg.migrationsDir, filename)

	// 生成迁移内容
	content := fmt.Sprintf(`-- Migration: %s
-- Description: %s
-- Version: %s

-- UP Migration
-- Add your UP migration SQL here
-- Example:
-- CREATE TABLE users (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     name VARCHAR(255) NOT NULL,
--     email VARCHAR(255) UNIQUE NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- DOWN Migration
-- Add your DOWN migration SQL here
-- Example:
-- DROP TABLE users;
`, name, description, version)

	// 确保目录存在
	if err := os.MkdirAll(mg.migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	fmt.Printf("Migration file created: %s\n", filePath)
	return nil
}
