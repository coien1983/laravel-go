package database

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// 测试模型结构体
type User struct {
	Model
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

// 自定义表名
func (u *User) TableName() string {
	return "users"
}

// 带钩子的用户模型
type UserWithHooks struct {
	Model
	Name     string `db:"name"`
	Email    string `db:"email"`
	HookFlag bool   `db:"hook_flag"`
}

func (u *UserWithHooks) TableName() string {
	return "users_with_hooks"
}

func (u *UserWithHooks) BeforeSave(conn Connection) error {
	u.HookFlag = true
	return nil
}

func (u *UserWithHooks) AfterSave(conn Connection) error {
	// 可以在这里执行保存后的逻辑
	return nil
}

func (u *UserWithHooks) BeforeDelete(conn Connection) error {
	// 可以在这里执行删除前的逻辑
	return nil
}

func (u *UserWithHooks) AfterDelete(conn Connection) error {
	// 可以在这里执行删除后的逻辑
	return nil
}

// 关联模型
type Profile struct {
	Model
	UserID int64  `db:"user_id"`
	Bio    string `db:"bio"`
	Avatar string `db:"avatar"`
}

func (p *Profile) TableName() string {
	return "profiles"
}

type Post struct {
	Model
	UserID  int64  `db:"user_id"`
	Title   string `db:"title"`
	Content string `db:"content"`
	Status  string `db:"status"`
}

func (p *Post) TableName() string {
	return "posts"
}

type Comment struct {
	Model
	PostID  int64  `db:"post_id"`
	UserID  int64  `db:"user_id"`
	Content string `db:"content"`
}

func (c *Comment) TableName() string {
	return "comments"
}

// 多对多关联的中间表
type UserRole struct {
	Model
	UserID int64 `db:"user_id"`
	RoleID int64 `db:"role_id"`
}

func (ur *UserRole) TableName() string {
	return "user_roles"
}

type Role struct {
	Model
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *Role) TableName() string {
	return "roles"
}

// 测试基础CRUD操作
func TestModelBasicCRUD(t *testing.T) {
	// 创建测试数据库连接
	config := &ConnectionConfig{
		Driver: SQLite,
		Host:   "test_basic_crud.db",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// 测试插入
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	model := &Model{}
	err = model.Save(conn, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID should be set after save")
	}

	if user.CreatedAt == nil {
		t.Error("CreatedAt should be set after save")
	}

	if user.UpdatedAt == nil {
		t.Error("UpdatedAt should be set after save")
	}

	// 测试查询
	var foundUser User
	err = model.Find(conn, user.ID, &foundUser)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if foundUser.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, foundUser.Name)
	}

	// 测试更新
	user.Name = "Jane Doe"
	err = model.Save(conn, user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	var updatedUser User
	err = model.Find(conn, user.ID, &updatedUser)
	if err != nil {
		t.Fatalf("Failed to find updated user: %v", err)
	}

	if updatedUser.Name != "Jane Doe" {
		t.Errorf("Expected updated name Jane Doe, got %s", updatedUser.Name)
	}

	// 测试条件查询
	var users []User
	err = model.Where(conn, &users, "age", ">=", 25)
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}

	if len(users) == 0 {
		t.Error("Should find at least one user")
	}

	// 测试统计
	count, err := model.Count(conn, &User{})
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	if count < 1 {
		t.Error("Should have at least one user")
	}

	// 测试存在性检查
	exists, err := model.Exists(conn, &User{}, "email", "john@example.com")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}

	if !exists {
		t.Error("User should exist")
	}

	// 测试软删除
	err = model.Delete(conn, user)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	if user.DeletedAt == nil {
		t.Error("DeletedAt should be set after soft delete")
	}

	// 验证软删除后无法通过普通查询找到
	var deletedUser User
	err = model.Find(conn, user.ID, &deletedUser)
	if err == nil {
		t.Error("Should not find soft deleted user")
	}
}

// 测试钩子功能
func TestModelHooks(t *testing.T) {
	// 创建测试数据库连接
	config := &ConnectionConfig{
		Driver: SQLite,
		Host:   "test_hooks.db",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS users_with_hooks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			hook_flag BOOLEAN DEFAULT FALSE,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// 测试保存钩子
	user := &UserWithHooks{
		Name:  "Hook User",
		Email: "hook@example.com",
	}

	model := &Model{}
	err = model.Save(conn, user)
	if err != nil {
		t.Fatalf("Failed to save user with hooks: %v", err)
	}

	if !user.HookFlag {
		t.Error("HookFlag should be set by BeforeSave hook")
	}

	// 测试删除钩子
	err = model.Delete(conn, user)
	if err != nil {
		t.Fatalf("Failed to delete user with hooks: %v", err)
	}
}

// 测试关联查询
func TestModelAssociations(t *testing.T) {
	// 创建测试数据库连接
	config := &ConnectionConfig{
		Driver: SQLite,
		Host:   "test_associations.db",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			bio TEXT,
			avatar TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create profiles table: %v", err)
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT,
			status TEXT DEFAULT 'published',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create posts table: %v", err)
	}

	// 插入测试数据
	user := &User{
		Name:  "Association User",
		Email: "assoc@example.com",
	}

	model := &Model{}
	err = model.Save(conn, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	profile := &Profile{
		UserID: user.ID,
		Bio:    "This is my bio",
		Avatar: "avatar.jpg",
	}
	err = model.Save(conn, profile)
	if err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	post1 := &Post{
		UserID:  user.ID,
		Title:   "First Post",
		Content: "Content of first post",
		Status:  "published",
	}
	err = model.Save(conn, post1)
	if err != nil {
		t.Fatalf("Failed to save post1: %v", err)
	}

	post2 := &Post{
		UserID:  user.ID,
		Title:   "Second Post",
		Content: "Content of second post",
		Status:  "draft",
	}
	err = model.Save(conn, post2)
	if err != nil {
		t.Fatalf("Failed to save post2: %v", err)
	}

	// 测试 HasOne 关联
	var userProfile Profile
	err = user.HasOne(conn, &userProfile, "user_id", "ID")
	if err != nil {
		t.Fatalf("Failed to load user profile: %v", err)
	}

	if userProfile.Bio != "This is my bio" {
		t.Errorf("Expected bio 'This is my bio', got %s", userProfile.Bio)
	}

	// 测试 HasMany 关联
	var userPosts []Post
	err = user.HasMany(conn, &userPosts, "user_id", "ID")
	if err != nil {
		t.Fatalf("Failed to load user posts: %v", err)
	}

	if len(userPosts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(userPosts))
	}

	// 测试 BelongsTo 关联
	var postUser User
	err = model.BelongsToModel(conn, post1, &postUser, "ID", "UserID")
	if err != nil {
		t.Fatalf("Failed to load post user: %v", err)
	}

	if postUser.Name != "Association User" {
		t.Errorf("Expected user name 'Association User', got %s", postUser.Name)
	}
}

// 测试多对多关联
func TestModelManyToMany(t *testing.T) {
	// 创建测试数据库连接
	config := &ConnectionConfig{
		Driver: SQLite,
		Host:   "test_many_to_many.db",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create roles table: %v", err)
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS user_roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			role_id INTEGER NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (role_id) REFERENCES roles (id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create user_roles table: %v", err)
	}

	// 插入测试数据
	user := &User{
		Name:  "ManyToMany User",
		Email: "manytomany@example.com",
	}

	model := &Model{}
	err = model.Save(conn, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	role1 := &Role{
		Name:        "Admin",
		Description: "Administrator role",
	}
	err = model.Save(conn, role1)
	if err != nil {
		t.Fatalf("Failed to save role1: %v", err)
	}

	role2 := &Role{
		Name:        "Editor",
		Description: "Editor role",
	}
	err = model.Save(conn, role2)
	if err != nil {
		t.Fatalf("Failed to save role2: %v", err)
	}

	// 创建关联
	userRole1 := &UserRole{
		UserID: user.ID,
		RoleID: role1.ID,
	}
	err = model.Save(conn, userRole1)
	if err != nil {
		t.Fatalf("Failed to save user_role1: %v", err)
	}

	userRole2 := &UserRole{
		UserID: user.ID,
		RoleID: role2.ID,
	}
	err = model.Save(conn, userRole2)
	if err != nil {
		t.Fatalf("Failed to save user_role2: %v", err)
	}

	// 测试 ManyToMany 关联
	var userRoles []Role
	err = user.ManyToMany(conn, &userRoles, "user_roles", "ID", "user_id", "role_id", "id")
	if err != nil {
		t.Fatalf("Failed to load user roles: %v", err)
	}

	if len(userRoles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(userRoles))
	}

	// 验证角色名称
	roleNames := make(map[string]bool)
	for _, role := range userRoles {
		roleNames[role.Name] = true
	}

	if !roleNames["Admin"] {
		t.Error("Expected Admin role not found")
	}

	if !roleNames["Editor"] {
		t.Error("Expected Editor role not found")
	}
}

// 测试自动时间戳
func TestModelAutoTimestamps(t *testing.T) {
	// 创建测试数据库连接
	config := &ConnectionConfig{
		Driver: SQLite,
		Host:   "test_timestamps.db",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// 测试插入时自动设置时间戳
	user := &User{
		Name:  "Timestamp User",
		Email: "timestamp@example.com",
	}

	model := &Model{}
	err = model.Save(conn, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	if user.CreatedAt == nil {
		t.Error("CreatedAt should be automatically set")
	}

	if user.UpdatedAt == nil {
		t.Error("UpdatedAt should be automatically set")
	}

	originalCreatedAt := user.CreatedAt
	originalUpdatedAt := user.UpdatedAt

	// 等待一秒以确保时间戳不同
	time.Sleep(1 * time.Second)

	// 测试更新时自动更新时间戳
	user.Name = "Updated Timestamp User"
	err = model.Save(conn, user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if user.CreatedAt != originalCreatedAt {
		t.Error("CreatedAt should not change on update")
	}

	if user.UpdatedAt == originalUpdatedAt {
		t.Error("UpdatedAt should be updated on save")
	}

	if !user.UpdatedAt.After(*originalUpdatedAt) {
		t.Error("UpdatedAt should be later than original UpdatedAt")
	}
}

// 测试软删除
func TestModelSoftDelete(t *testing.T) {
	// 创建测试数据库连接
	config := &ConnectionConfig{
		Driver: SQLite,
		Host:   "test_soft_delete.db",
	}

	conn, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// 插入测试数据
	user := &User{
		Name:  "Soft Delete User",
		Email: "softdelete@example.com",
	}

	model := &Model{}
	err = model.Save(conn, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// 验证记录存在
	exists, err := model.Exists(conn, &User{}, "id", user.ID)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}

	if !exists {
		t.Error("User should exist before soft delete")
	}

	// 执行软删除
	err = model.Delete(conn, user)
	if err != nil {
		t.Fatalf("Failed to soft delete user: %v", err)
	}

	if user.DeletedAt == nil {
		t.Error("DeletedAt should be set after soft delete")
	}

	// 验证软删除后无法通过普通查询找到
	exists, err = model.Exists(conn, &User{}, "id", user.ID)
	if err != nil {
		t.Fatalf("Failed to check existence after soft delete: %v", err)
	}

	if exists {
		t.Error("User should not exist after soft delete")
	}

	// 验证数据库中确实有 deleted_at 值
	var deletedAt *time.Time
	err = conn.QueryRow("SELECT deleted_at FROM users WHERE id = ?", user.ID).Scan(&deletedAt)
	if err != nil {
		t.Fatalf("Failed to query deleted_at: %v", err)
	}

	if deletedAt == nil {
		t.Error("deleted_at should be set in database")
	}
}

// 测试辅助函数
func TestModelHelperFunctions(t *testing.T) {
	// 测试 getTableName
	user := &User{}
	tableName := getTableName(user)
	if tableName != "users" {
		t.Errorf("Expected table name 'users', got %s", tableName)
	}

	// 测试 getPrimaryKey
	pk := getPrimaryKey(user)
	if pk != "id" {
		t.Errorf("Expected primary key 'id', got %s", pk)
	}

	// 测试 structToMap
	user.Name = "Test User"
	user.Email = "test@example.com"
	user.Age = 25

	data := structToMap(user)
	if data["name"] != "Test User" {
		t.Errorf("Expected name 'Test User', got %v", data["name"])
	}

	if data["email"] != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %v", data["email"])
	}

	if data["age"] != 25 {
		t.Errorf("Expected age 25, got %v", data["age"])
	}

	// 测试 isZero
	if !isZero(reflect.ValueOf(0)) {
		t.Error("0 should be zero value")
	}

	if !isZero(reflect.ValueOf("")) {
		t.Error("Empty string should be zero value")
	}

	if isZero(reflect.ValueOf("test")) {
		t.Error("'test' should not be zero value")
	}
}

// 性能测试
func BenchmarkModelSave(b *testing.B) {
	// 创建测试数据库连接
	config := &ConnectionConfig{
		Driver: SQLite,
		Host:   "benchmark_save.db",
	}

	conn, err := NewConnection(config)
	if err != nil {
		b.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		b.Fatalf("Failed to create table: %v", err)
	}

	model := &Model{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &User{
			Name:  fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   i % 100,
		}

		err := model.Save(conn, user)
		if err != nil {
			b.Fatalf("Failed to save user: %v", err)
		}
	}
}

func BenchmarkModelFind(b *testing.B) {
	// 创建测试数据库连接
	config := &ConnectionConfig{
		Driver: SQLite,
		Host:   "benchmark_find.db",
	}

	conn, err := NewConnection(config)
	if err != nil {
		b.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建测试表并插入数据
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		b.Fatalf("Failed to create table: %v", err)
	}

	// 插入测试数据
	model := &Model{}
	userIDs := make([]int64, 100)
	for i := 0; i < 100; i++ {
		user := &User{
			Name:  fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   i % 100,
		}

		err := model.Save(conn, user)
		if err != nil {
			b.Fatalf("Failed to save user: %v", err)
		}
		userIDs[i] = user.ID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user User
		userID := userIDs[i%len(userIDs)]
		err := model.Find(conn, userID, &user)
		if err != nil {
			b.Fatalf("Failed to find user: %v", err)
		}
	}
}
