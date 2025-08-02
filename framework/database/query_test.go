package database

import (
	"strings"
	"testing"
)

func TestNewQueryBuilder(t *testing.T) {
	// 创建测试连接
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建查询构建器
	qb := NewQueryBuilder(conn)
	if qb == nil {
		t.Fatal("QueryBuilder should not be nil")
	}

	if qb.connection != conn {
		t.Error("Connection should be set correctly")
	}
}

func TestQueryBuilderTable(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	qb := NewQueryBuilder(conn)
	qb.Table("users")

	if qb.table != "users" {
		t.Errorf("Expected table 'users', got '%s'", qb.table)
	}
}

func TestQueryBuilderSelect(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	qb := NewQueryBuilder(conn)
	qb.Select("id", "name", "email")

	if len(qb.selects) != 3 {
		t.Errorf("Expected 3 selects, got %d", len(qb.selects))
	}

	expected := []string{"id", "name", "email"}
	for i, select_ := range qb.selects {
		if select_ != expected[i] {
			t.Errorf("Expected select '%s', got '%s'", expected[i], select_)
		}
	}
}

func TestQueryBuilderWhere(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	qb := NewQueryBuilder(conn)
	qb.Where("age", ">", 18)

	if len(qb.wheres) != 1 {
		t.Errorf("Expected 1 where condition, got %d", len(qb.wheres))
	}

	where := qb.wheres[0]
	if where.Column != "age" {
		t.Errorf("Expected column 'age', got '%s'", where.Column)
	}
	if where.Operator != ">" {
		t.Errorf("Expected operator '>', got '%s'", where.Operator)
	}
	if where.Value != 18 {
		t.Errorf("Expected value 18, got %v", where.Value)
	}
}

func TestQueryBuilderWhereEq(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	qb := NewQueryBuilder(conn)
	qb.WhereEq("name", "John")

	if len(qb.wheres) != 1 {
		t.Errorf("Expected 1 where condition, got %d", len(qb.wheres))
	}

	where := qb.wheres[0]
	if where.Column != "name" {
		t.Errorf("Expected column 'name', got '%s'", where.Column)
	}
	if where.Operator != "=" {
		t.Errorf("Expected operator '=', got '%s'", where.Operator)
	}
	if where.Value != "John" {
		t.Errorf("Expected value 'John', got %v", where.Value)
	}
}

func TestQueryBuilderWhereIn(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	qb := NewQueryBuilder(conn)
	qb.WhereIn("id", []interface{}{1, 2, 3})

	if len(qb.wheres) != 1 {
		t.Errorf("Expected 1 where condition, got %d", len(qb.wheres))
	}

	where := qb.wheres[0]
	if !where.Raw {
		t.Error("Expected raw SQL for WHERE IN")
	}
	if !strings.Contains(where.RawSQL, "id IN (?,?,?)") {
		t.Errorf("Expected SQL to contain 'id IN (?,?,?)', got '%s'", where.RawSQL)
	}
}

func TestQueryBuilderOrderBy(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	qb := NewQueryBuilder(conn)
	qb.OrderBy("created_at", "DESC")

	if len(qb.orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(qb.orders))
	}

	order := qb.orders[0]
	if order.Column != "created_at" {
		t.Errorf("Expected column 'created_at', got '%s'", order.Column)
	}
	if order.Direction != "DESC" {
		t.Errorf("Expected direction 'DESC', got '%s'", order.Direction)
	}
}

func TestQueryBuilderLimitOffset(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	qb := NewQueryBuilder(conn)
	qb.Limit(10).Offset(20)

	if qb.limit != 10 {
		t.Errorf("Expected limit 10, got %d", qb.limit)
	}
	if qb.offset != 20 {
		t.Errorf("Expected offset 20, got %d", qb.offset)
	}
}

func TestQueryBuilderBuildSelectQuery(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	qb := NewQueryBuilder(conn)
	qb.Table("users").
		Select("id", "name").
		WhereEq("active", true).
		OrderBy("name", "ASC").
		Limit(10)

	query, args := qb.buildSelectQuery()

	expectedQuery := "SELECT id, name FROM users WHERE active = ? ORDER BY name ASC LIMIT 10"
	if query != expectedQuery {
		t.Errorf("Expected query '%s', got '%s'", expectedQuery, query)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(args))
	}
	if args[0] != true {
		t.Errorf("Expected argument true, got %v", args[0])
	}
}

func TestQueryBuilderComplexQuery(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	qb := NewQueryBuilder(conn)
	qb.Table("users").
		Select("users.id", "users.name", "profiles.bio").
		LeftJoin("profiles", "users.id = profiles.user_id").
		WhereEq("users.active", true).
		WhereGt("users.age", 18).
		WhereLike("users.name", "John").
		OrWhere("users.email", "LIKE", "%@example.com").
		GroupBy("users.id").
		Having("COUNT(*)", ">", 1).
		OrderBy("users.created_at", "DESC").
		Limit(20).
		Offset(40)

	query, args := qb.buildSelectQuery()

	// 验证查询包含所有必要的子句
	if !strings.Contains(query, "SELECT users.id, users.name, profiles.bio") {
		t.Error("Query should contain SELECT clause")
	}
	if !strings.Contains(query, "FROM users") {
		t.Error("Query should contain FROM clause")
	}
	if !strings.Contains(query, "LEFT JOIN profiles") {
		t.Error("Query should contain LEFT JOIN clause")
	}
	if !strings.Contains(query, "WHERE users.active = ?") {
		t.Error("Query should contain WHERE clause")
	}
	if !strings.Contains(query, "GROUP BY users.id") {
		t.Error("Query should contain GROUP BY clause")
	}
	if !strings.Contains(query, "HAVING COUNT(*) > ?") {
		t.Error("Query should contain HAVING clause")
	}
	if !strings.Contains(query, "ORDER BY users.created_at DESC") {
		t.Error("Query should contain ORDER BY clause")
	}
	if !strings.Contains(query, "LIMIT 20") {
		t.Error("Query should contain LIMIT clause")
	}
	if !strings.Contains(query, "OFFSET 40") {
		t.Error("Query should contain OFFSET clause")
	}

	// 验证参数数量
	expectedArgs := 5 // active, age, name, email, having
	if len(args) != expectedArgs {
		t.Errorf("Expected %d arguments, got %d", expectedArgs, len(args))
	}
}

func TestQueryBuilderIntegration(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test_integration.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS test_users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE,
			age INTEGER,
			active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// 插入测试数据
	_, err = conn.Exec(`
		INSERT OR REPLACE INTO test_users (id, name, email, age, active) VALUES 
		(1, 'John Doe', 'john@example.com', 25, 1),
		(2, 'Jane Smith', 'jane@example.com', 30, 1),
		(3, 'Bob Johnson', 'bob@example.com', 35, 0)
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// 测试查询构建器
	qb := NewQueryBuilder(conn)
	results, err := qb.Table("test_users").
		Select("id", "name", "email", "age").
		WhereEq("active", true).
		OrderBy("name", "ASC").
		Get()

	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// 验证第一个结果
	first := results[0]
	if first["name"] != "Jane Smith" {
		t.Errorf("Expected first name 'Jane Smith', got '%s'", first["name"])
	}

	// 测试计数
	count, err := qb.Table("test_users").WhereEq("active", true).Count()
	if err != nil {
		t.Fatalf("Failed to count: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// 测试分页
	paginated, err := qb.Table("test_users").
		WhereEq("active", true).
		Paginate(1, 1)

	if err != nil {
		t.Fatalf("Failed to paginate: %v", err)
	}

	if paginated["total"] != int64(2) {
		t.Errorf("Expected total 2, got %v", paginated["total"])
	}
	if paginated["per_page"] != 1 {
		t.Errorf("Expected per_page 1, got %v", paginated["per_page"])
	}
	if paginated["current_page"] != 1 {
		t.Errorf("Expected current_page 1, got %v", paginated["current_page"])
	}
}

func TestQueryBuilderAggregates(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test_aggregates.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS test_orders (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			amount DECIMAL(10,2),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// 插入测试数据
	_, err = conn.Exec(`
		INSERT OR REPLACE INTO test_orders (id, user_id, amount) VALUES 
		(1, 1, 100.50),
		(2, 1, 200.75),
		(3, 2, 150.25),
		(4, 2, 300.00)
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	qb := NewQueryBuilder(conn)

	// 测试 SUM
	sum, err := qb.Table("test_orders").Sum("amount")
	if err != nil {
		t.Fatalf("Failed to sum: %v", err)
	}
	expectedSum := 751.5
	if sum != expectedSum {
		t.Errorf("Expected sum %.2f, got %.2f", expectedSum, sum)
	}

	// 测试 AVG
	avg, err := qb.Table("test_orders").Avg("amount")
	if err != nil {
		t.Fatalf("Failed to avg: %v", err)
	}
	expectedAvg := 187.875
	if avg != expectedAvg {
		t.Errorf("Expected avg %.3f, got %.3f", expectedAvg, avg)
	}

	// 测试 MAX
	max, err := qb.Table("test_orders").Max("amount")
	if err != nil {
		t.Fatalf("Failed to max: %v", err)
	}
	expectedMax := 300.0
	if max != expectedMax {
		t.Errorf("Expected max %.1f, got %.1f", expectedMax, max)
	}

	// 测试 MIN
	min, err := qb.Table("test_orders").Min("amount")
	if err != nil {
		t.Fatalf("Failed to min: %v", err)
	}
	expectedMin := 100.5
	if min != expectedMin {
		t.Errorf("Expected min %.1f, got %.1f", expectedMin, min)
	}
}

func TestQueryBuilderExists(t *testing.T) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test_exists.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS test_users_exists (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// 插入测试数据
	_, err = conn.Exec(`
		INSERT OR REPLACE INTO test_users_exists (id, name, email) VALUES 
		(1, 'John Doe', 'john@example.com')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	qb := NewQueryBuilder(conn)

	// 测试存在
	exists, err := qb.Table("test_users_exists").WhereEq("name", "John Doe").Exists()
	if err != nil {
		t.Fatalf("Failed to check exists: %v", err)
	}
	if !exists {
		t.Error("Expected user to exist")
	}

	// 测试不存在
	notExists, err := qb.Table("test_users_exists").WhereEq("name", "Non Existent").Exists()
	if err != nil {
		t.Fatalf("Failed to check exists: %v", err)
	}
	if notExists {
		t.Error("Expected user to not exist")
	}
}

func BenchmarkQueryBuilder(b *testing.B) {
	config := &ConnectionConfig{
		Driver:   SQLite,
		Database: "test.db",
	}
	conn, err := NewConnection(config)
	if err != nil {
		b.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb := NewQueryBuilder(conn)
		qb.Table("test_users").
			Select("id", "name", "email").
			WhereEq("active", true).
			OrderBy("name", "ASC").
			Limit(10)
		
		_, _ = qb.buildSelectQuery()
	}
} 