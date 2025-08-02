package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite 驱动

	"laravel-go/framework/config"
	"laravel-go/framework/errors"
)

// Driver 数据库驱动类型
type Driver string

const (
	MySQL      Driver = "mysql"
	PostgreSQL Driver = "postgres"
	SQLite     Driver = "sqlite"
	SQLServer  Driver = "sqlserver"
)

// ConnectionConfig 连接配置
type ConnectionConfig struct {
	Driver          Driver
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	Charset         string
	Timezone        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	SSLMode         string
	Options         map[string]string
}

// Connection 数据库连接接口
type Connection interface {
	// 获取原始数据库连接
	DB() *sql.DB

	// 执行查询
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// 执行命令
	Exec(query string, args ...interface{}) (sql.Result, error)

	// 开始事务
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	// 关闭连接
	Close() error

	// 检查连接状态
	Ping() error
	PingContext(ctx context.Context) error

	// 获取连接统计信息
	Stats() sql.DBStats
}

// connection 数据库连接实现
type connection struct {
	db     *sql.DB
	config *ConnectionConfig
	mutex  sync.RWMutex
}

// NewConnection 创建新的数据库连接
func NewConnection(config *ConnectionConfig) (Connection, error) {
	dsn, err := buildDSN(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build DSN")
	}

	// SQLite 驱动名称是 "sqlite3"，不是 "sqlite"
	driverName := string(config.Driver)
	if config.Driver == SQLite {
		driverName = "sqlite3"
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database connection")
	}

	// 设置连接池参数
	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, errors.Wrap(err, fmt.Sprintf("failed to ping database with DSN: %s", dsn))
	}

	return &connection{
		db:     db,
		config: config,
	}, nil
}

// DB 获取原始数据库连接
func (c *connection) DB() *sql.DB {
	return c.db
}

// Query 执行查询
func (c *connection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

// QueryRow 执行单行查询
func (c *connection) QueryRow(query string, args ...interface{}) *sql.Row {
	return c.db.QueryRow(query, args...)
}

// QueryContext 执行查询（带上下文）
func (c *connection) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

// QueryRowContext 执行单行查询（带上下文）
func (c *connection) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return c.db.QueryRowContext(ctx, query, args...)
}

// Exec 执行命令
func (c *connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

// Begin 开始事务
func (c *connection) Begin() (*sql.Tx, error) {
	return c.db.Begin()
}

// BeginTx 开始事务（带上下文）
func (c *connection) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}

// Close 关闭连接
func (c *connection) Close() error {
	return c.db.Close()
}

// Ping 检查连接状态
func (c *connection) Ping() error {
	return c.db.Ping()
}

// PingContext 检查连接状态（带上下文）
func (c *connection) PingContext(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// Stats 获取连接统计信息
func (c *connection) Stats() sql.DBStats {
	return c.db.Stats()
}

// buildDSN 构建数据库连接字符串
func buildDSN(config *ConnectionConfig) (string, error) {
	switch config.Driver {
	case MySQL:
		return buildMySQLDSN(config), nil
	case PostgreSQL:
		return buildPostgreSQLDSN(config), nil
	case SQLite:
		return buildSQLiteDSN(config), nil
	case SQLServer:
		return buildSQLServerDSN(config), nil
	default:
		return "", errors.New("unsupported database driver: " + string(config.Driver))
	}
}

// buildMySQLDSN 构建 MySQL 连接字符串
func buildMySQLDSN(config *ConnectionConfig) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	if config.Charset != "" {
		dsn += "?charset=" + config.Charset
	}
	if config.Timezone != "" {
		if config.Charset != "" {
			dsn += "&time_zone=" + config.Timezone
		} else {
			dsn += "?time_zone=" + config.Timezone
		}
	}

	// 添加其他选项
	if len(config.Options) > 0 {
		separator := "?"
		if config.Charset != "" || config.Timezone != "" {
			separator = "&"
		}
		for key, value := range config.Options {
			dsn += separator + key + "=" + value
			separator = "&"
		}
	}

	return dsn
}

// buildPostgreSQLDSN 构建 PostgreSQL 连接字符串
func buildPostgreSQLDSN(config *ConnectionConfig) string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database)

	if config.SSLMode != "" {
		dsn += " sslmode=" + config.SSLMode
	}

	// 添加其他选项
	for key, value := range config.Options {
		dsn += " " + key + "=" + value
	}

	return dsn
}

// buildSQLiteDSN 构建 SQLite 连接字符串
func buildSQLiteDSN(config *ConnectionConfig) string {
	// SQLite 使用 file: 前缀，内存数据库使用 :memory:
	if config.Database == ":memory:" {
		return "file::memory:?cache=shared"
	}
	return config.Database
}

// buildSQLServerDSN 构建 SQL Server 连接字符串
func buildSQLServerDSN(config *ConnectionConfig) string {
	dsn := fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database)

	// 添加其他选项
	for key, value := range config.Options {
		dsn += ";" + key + "=" + value
	}

	return dsn
}

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]Connection
	configs     map[string]*ConnectionConfig
	mutex       sync.RWMutex
	// 添加健康检查和清理机制
	healthTicker *time.Ticker
	stopChan     chan struct{}
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager() *ConnectionManager {
	cm := &ConnectionManager{
		connections:  make(map[string]Connection),
		configs:      make(map[string]*ConnectionConfig),
		healthTicker: time.NewTicker(60 * time.Second), // 每分钟检查一次
		stopChan:     make(chan struct{}),
	}

	// 启动健康检查协程
	go cm.healthCheckRoutine()

	return cm
}

// healthCheckRoutine 健康检查协程
func (cm *ConnectionManager) healthCheckRoutine() {
	for {
		select {
		case <-cm.healthTicker.C:
			cm.checkConnectionsHealth()
		case <-cm.stopChan:
			return
		}
	}
}

// checkConnectionsHealth 检查连接健康状态
func (cm *ConnectionManager) checkConnectionsHealth() {
	cm.mutex.RLock()
	connections := make(map[string]Connection)
	for name, conn := range cm.connections {
		connections[name] = conn
	}
	cm.mutex.RUnlock()

	for name, conn := range connections {
		// 使用超时上下文进行健康检查
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := conn.PingContext(ctx)
		cancel()

		if err != nil {
			// 连接不健康，尝试重新连接
			go cm.reconnectConnection(name)
		}
	}
}

// reconnectConnection 重新连接
func (cm *ConnectionManager) reconnectConnection(name string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 关闭旧连接
	if conn, exists := cm.connections[name]; exists {
		conn.Close()
		delete(cm.connections, name)
	}

	// 创建新连接
	config, exists := cm.configs[name]
	if !exists {
		return
	}

	conn, err := NewConnection(config)
	if err != nil {
		// 记录错误，但不阻塞其他操作
		return
	}

	cm.connections[name] = conn
}

// AddConnection 添加连接配置
func (cm *ConnectionManager) AddConnection(name string, config *ConnectionConfig) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.configs[name] = config
}

// GetConnection 获取连接
func (cm *ConnectionManager) GetConnection(name string) (Connection, error) {
	return cm.GetConnectionWithTimeout(name, 10*time.Second)
}

// GetConnectionWithTimeout 获取连接（带超时）
func (cm *ConnectionManager) GetConnectionWithTimeout(name string, timeout time.Duration) (Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 使用 channel 避免死锁
	connChan := make(chan Connection, 1)
	errChan := make(chan error, 1)

	go func() {
		conn, err := cm.getConnectionInternal(name)
		if err != nil {
			errChan <- err
			return
		}
		connChan <- conn
	}()

	select {
	case conn := <-connChan:
		return conn, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// getConnectionInternal 内部获取连接方法
func (cm *ConnectionManager) getConnectionInternal(name string) (Connection, error) {
	cm.mutex.RLock()
	conn, exists := cm.connections[name]
	cm.mutex.RUnlock()

	if exists {
		// 检查连接是否健康
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := conn.PingContext(ctx)
		cancel()

		if err == nil {
			return conn, nil
		}

		// 连接不健康，需要重新创建
		cm.mutex.Lock()
		conn.Close()
		delete(cm.connections, name)
		cm.mutex.Unlock()
	}

	// 创建新连接
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 双重检查
	if conn, exists = cm.connections[name]; exists {
		return conn, nil
	}

	config, exists := cm.configs[name]
	if !exists {
		return nil, errors.New("connection config not found: " + name)
	}

	conn, err := NewConnection(config)
	if err != nil {
		return nil, err
	}

	cm.connections[name] = conn
	return conn, nil
}

// CloseConnection 关闭连接
func (cm *ConnectionManager) CloseConnection(name string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if conn, exists := cm.connections[name]; exists {
		err := conn.Close()
		delete(cm.connections, name)
		return err
	}

	return nil
}

// CloseAll 关闭所有连接
func (cm *ConnectionManager) CloseAll() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 停止健康检查
	close(cm.stopChan)
	cm.healthTicker.Stop()

	var lastError error
	for name, conn := range cm.connections {
		if err := conn.Close(); err != nil {
			lastError = err
		}
		delete(cm.connections, name)
	}

	return lastError
}

// GetStats 获取所有连接的统计信息
func (cm *ConnectionManager) GetStats() map[string]sql.DBStats {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	stats := make(map[string]sql.DBStats)
	for name, conn := range cm.connections {
		stats[name] = conn.Stats()
	}

	return stats
}

// LoadFromConfig 从配置加载连接
func (cm *ConnectionManager) LoadFromConfig(cfg *config.Config) error {
	connections := cfg.Get("database.connections", map[string]interface{}{}).(map[string]interface{})

	for name, connConfig := range connections {
		config, err := parseConnectionConfig(connConfig.(map[string]interface{}))
		if err != nil {
			return errors.Wrap(err, "failed to parse connection config: "+name)
		}
		cm.AddConnection(name, config)
	}

	return nil
}

// parseConnectionConfig 解析连接配置
func parseConnectionConfig(data map[string]interface{}) (*ConnectionConfig, error) {
	config := &ConnectionConfig{
		Options: make(map[string]string),
	}

	// 解析基本配置
	if driver, ok := data["driver"].(string); ok {
		config.Driver = Driver(driver)
	} else {
		return nil, errors.New("driver is required")
	}

	if host, ok := data["host"].(string); ok {
		config.Host = host
	}

	if port, ok := data["port"].(int); ok {
		config.Port = port
	} else {
		// 设置默认端口
		switch config.Driver {
		case MySQL:
			config.Port = 3306
		case PostgreSQL:
			config.Port = 5432
		case SQLServer:
			config.Port = 1433
		}
	}

	if database, ok := data["database"].(string); ok {
		config.Database = database
	}

	if username, ok := data["username"].(string); ok {
		config.Username = username
	}

	if password, ok := data["password"].(string); ok {
		config.Password = password
	}

	if charset, ok := data["charset"].(string); ok {
		config.Charset = charset
	}

	if timezone, ok := data["timezone"].(string); ok {
		config.Timezone = timezone
	}

	if maxOpenConns, ok := data["max_open_conns"].(int); ok {
		config.MaxOpenConns = maxOpenConns
	}

	if maxIdleConns, ok := data["max_idle_conns"].(int); ok {
		config.MaxIdleConns = maxIdleConns
	}

	if connMaxLifetime, ok := data["conn_max_lifetime"].(string); ok {
		if duration, err := time.ParseDuration(connMaxLifetime); err == nil {
			config.ConnMaxLifetime = duration
		}
	}

	if connMaxIdleTime, ok := data["conn_max_idle_time"].(string); ok {
		if duration, err := time.ParseDuration(connMaxIdleTime); err == nil {
			config.ConnMaxIdleTime = duration
		}
	}

	if sslMode, ok := data["ssl_mode"].(string); ok {
		config.SSLMode = sslMode
	}

	// 解析其他选项
	if options, ok := data["options"].(map[string]interface{}); ok {
		for key, value := range options {
			if strValue, ok := value.(string); ok {
				config.Options[key] = strValue
			}
		}
	}

	return config, nil
}

// Validate 验证连接配置
func (cfg *ConnectionConfig) Validate() error {
	if cfg.Driver == "" {
		return errors.New("database driver is required")
	}

	if cfg.Host == "" && cfg.Driver != SQLite {
		return errors.New("database host is required")
	}

	if cfg.Database == "" {
		return errors.New("database name is required")
	}

	if cfg.MaxOpenConns <= 0 {
		return errors.New("MaxOpenConns must be positive")
	}

	if cfg.MaxIdleConns <= 0 {
		return errors.New("MaxIdleConns must be positive")
	}

	if cfg.ConnMaxLifetime <= 0 {
		return errors.New("ConnMaxLifetime must be positive")
	}

	if cfg.ConnMaxIdleTime <= 0 {
		return errors.New("ConnMaxIdleTime must be positive")
	}

	return nil
}
