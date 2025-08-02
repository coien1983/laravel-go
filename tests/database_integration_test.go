package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"laravel-go/framework/config"
	"laravel-go/framework/container"
	"laravel-go/framework/database"
)

// DatabaseIntegrationTestSuite 数据库集成测试套件
type DatabaseIntegrationTestSuite struct {
	suite.Suite
	app *container.Container
	db  database.Connection
}

// SetupSuite 设置测试套件
func (suite *DatabaseIntegrationTestSuite) SetupSuite() {
	// 初始化应用容器
	suite.app = container.NewContainer()

	// 注册配置
	suite.app.Singleton("config", func() interface{} {
		return config.NewConfig()
	})

	// 注册数据库连接
	suite.app.Singleton("database", func() interface{} {
		db, err := database.NewConnection(&database.Config{
			Driver:   "sqlite",
			Database: ":memory:",
		})
		if err != nil {
			suite.T().Fatalf("Failed to create database connection: %v", err)
		}
		return db
	})

	// 获取数据库连接
	suite.db = suite.app.Make("database").(database.Connection)

	// 设置数据库表
	suite.setupDatabase()
}

// TearDownSuite 清理测试套件
func (suite *DatabaseIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest 设置每个测试
func (suite *DatabaseIntegrationTestSuite) SetupTest() {
	// 清理数据库
	suite.cleanupDatabase()
}

// setupDatabase 设置数据库表
func (suite *DatabaseIntegrationTestSuite) setupDatabase() {
	// 创建用户表
	_, err := suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			age INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME NULL
		)
	`)
	if err != nil {
		suite.T().Fatalf("Failed to create users table: %v", err)
	}

	// 创建文章表
	_, err = suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title VARCHAR(255) NOT NULL,
			content TEXT,
			status VARCHAR(50) DEFAULT 'draft',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		suite.T().Fatalf("Failed to create posts table: %v", err)
	}

	// 创建评论表
	_, err = suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		suite.T().Fatalf("Failed to create comments table: %v", err)
	}

	// 创建标签表
	_, err = suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		suite.T().Fatalf("Failed to create tags table: %v", err)
	}

	// 创建文章标签关联表
	_, err = suite.db.Exec(`
		CREATE TABLE IF NOT EXISTS post_tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id),
			UNIQUE(post_id, tag_id)
		)
	`)
	if err != nil {
		suite.T().Fatalf("Failed to create post_tags table: %v", err)
	}
}

// cleanupDatabase 清理数据库
func (suite *DatabaseIntegrationTestSuite) cleanupDatabase() {
	suite.db.Exec("DELETE FROM post_tags")
	suite.db.Exec("DELETE FROM tags")
	suite.db.Exec("DELETE FROM comments")
	suite.db.Exec("DELETE FROM posts")
	suite.db.Exec("DELETE FROM users")
}

// TestBasicCRUD 测试基础CRUD操作
func (suite *DatabaseIntegrationTestSuite) TestBasicCRUD() {
	// 创建用户
	user := &User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
		Age:      25,
	}

	// 保存用户
	err := user.Save(suite.db)
	suite.NoError(err)
	suite.NotZero(user.ID)

	// 查询用户
	foundUser := &User{}
	err = foundUser.Find(suite.db, user.ID)
	suite.NoError(err)
	suite.Equal(user.Name, foundUser.Name)
	suite.Equal(user.Email, foundUser.Email)
	suite.Equal(user.Age, foundUser.Age)

	// 更新用户
	foundUser.Name = "John Smith"
	foundUser.Age = 30
	err = foundUser.Save(suite.db)
	suite.NoError(err)

	// 验证更新
	updatedUser := &User{}
	err = updatedUser.Find(suite.db, user.ID)
	suite.NoError(err)
	suite.Equal("John Smith", updatedUser.Name)
	suite.Equal(30, updatedUser.Age)

	// 删除用户
	err = updatedUser.Delete(suite.db)
	suite.NoError(err)

	// 验证删除
	err = updatedUser.Find(suite.db, user.ID)
	suite.Error(err) // 应该找不到用户
}

// TestQueryBuilder 测试查询构建器
func (suite *DatabaseIntegrationTestSuite) TestQueryBuilder() {
	// 创建测试数据
	users := []*User{
		{Name: "Alice", Email: "alice@example.com", Password: "pass1", Age: 25},
		{Name: "Bob", Email: "bob@example.com", Password: "pass2", Age: 30},
		{Name: "Charlie", Email: "charlie@example.com", Password: "pass3", Age: 35},
	}

	for _, user := range users {
		err := user.Save(suite.db)
		suite.NoError(err)
	}

	// 测试WHERE查询
	var youngUsers []*User
	err := suite.db.Query("SELECT * FROM users WHERE age < ?", 30).Get(&youngUsers)
	suite.NoError(err)
	suite.Len(youngUsers, 2)

	// 测试ORDER BY
	var orderedUsers []*User
	err = suite.db.Query("SELECT * FROM users ORDER BY age DESC").Get(&orderedUsers)
	suite.NoError(err)
	suite.Len(orderedUsers, 3)
	suite.Equal("Charlie", orderedUsers[0].Name)

	// 测试LIMIT
	var limitedUsers []*User
	err = suite.db.Query("SELECT * FROM users LIMIT 2").Get(&limitedUsers)
	suite.NoError(err)
	suite.Len(limitedUsers, 2)

	// 测试聚合函数
	var count int
	err = suite.db.Query("SELECT COUNT(*) FROM users").Scan(&count)
	suite.NoError(err)
	suite.Equal(3, count)

	var avgAge float64
	err = suite.db.Query("SELECT AVG(age) FROM users").Scan(&avgAge)
	suite.NoError(err)
	suite.Equal(30.0, avgAge)
}

// TestModelAssociations 测试模型关联
func (suite *DatabaseIntegrationTestSuite) TestModelAssociations() {
	// 创建用户
	user := &User{
		Name:     "Author",
		Email:    "author@example.com",
		Password: "password",
		Age:      28,
	}
	err := user.Save(suite.db)
	suite.NoError(err)

	// 创建文章
	post := &Post{
		UserID:  user.ID,
		Title:   "Test Post",
		Content: "This is a test post content",
		Status:  "published",
	}
	err = post.Save(suite.db)
	suite.NoError(err)

	// 创建评论
	comment := &Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "Great post!",
	}
	err = comment.Save(suite.db)
	suite.NoError(err)

	// 测试BelongsTo关联
	foundPost := &Post{}
	err = foundPost.Find(suite.db, post.ID)
	suite.NoError(err)

	// 获取文章的作者
	author, err := foundPost.BelongsTo(suite.db, "users", "user_id", "id")
	suite.NoError(err)
	suite.NotNil(author)
	suite.Equal(user.Name, author["name"])

	// 测试HasMany关联
	// 获取用户的所有文章
	userPosts, err := user.HasMany(suite.db, "posts", "user_id", "id")
	suite.NoError(err)
	suite.Len(userPosts, 1)
	suite.Equal(post.Title, userPosts[0]["title"])

	// 获取文章的所有评论
	postComments, err := post.HasMany(suite.db, "comments", "post_id", "id")
	suite.NoError(err)
	suite.Len(postComments, 1)
	suite.Equal(comment.Content, postComments[0]["content"])
}

// TestManyToMany 测试多对多关联
func (suite *DatabaseIntegrationTestSuite) TestManyToMany() {
	// 创建用户
	user := &User{
		Name:     "Blogger",
		Email:    "blogger@example.com",
		Password: "password",
		Age:      25,
	}
	err := user.Save(suite.db)
	suite.NoError(err)

	// 创建文章
	post := &Post{
		UserID:  user.ID,
		Title:   "Multi-tag Post",
		Content: "This post has multiple tags",
		Status:  "published",
	}
	err = post.Save(suite.db)
	suite.NoError(err)

	// 创建标签
	tags := []*Tag{
		{Name: "Technology"},
		{Name: "Programming"},
		{Name: "Go"},
	}

	for _, tag := range tags {
		err := tag.Save(suite.db)
		suite.NoError(err)
	}

	// 关联文章和标签
	for _, tag := range tags {
		_, err := suite.db.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", post.ID, tag.ID)
		suite.NoError(err)
	}

	// 测试获取文章的所有标签
	postTags, err := post.BelongsToMany(suite.db, "tags", "post_tags", "post_id", "tag_id", "id", "id")
	suite.NoError(err)
	suite.Len(postTags, 3)

	// 测试获取标签的所有文章
	tagPosts, err := tags[0].BelongsToMany(suite.db, "posts", "post_tags", "tag_id", "post_id", "id", "id")
	suite.NoError(err)
	suite.Len(tagPosts, 1)
	suite.Equal(post.Title, tagPosts[0]["title"])
}

// TestSoftDelete 测试软删除
func (suite *DatabaseIntegrationTestSuite) TestSoftDelete() {
	// 创建用户
	user := &User{
		Name:     "Soft Delete User",
		Email:    "soft@example.com",
		Password: "password",
		Age:      30,
	}
	err := user.Save(suite.db)
	suite.NoError(err)

	// 软删除用户
	err = user.SoftDelete(suite.db)
	suite.NoError(err)

	// 验证软删除
	suite.NotNil(user.DeletedAt)

	// 查询应该找不到用户（默认不包含软删除）
	foundUser := &User{}
	err = foundUser.Find(suite.db, user.ID)
	suite.Error(err)

	// 强制查询包含软删除的记录
	err = suite.db.Query("SELECT * FROM users WHERE id = ?", user.ID).Get(foundUser)
	suite.NoError(err)
	suite.Equal(user.Name, foundUser.Name)
	suite.NotNil(foundUser.DeletedAt)

	// 恢复用户
	err = user.Restore(suite.db)
	suite.NoError(err)

	// 验证恢复
	err = foundUser.Find(suite.db, user.ID)
	suite.NoError(err)
	suite.Nil(foundUser.DeletedAt)
}

// TestModelHooks 测试模型钩子
func (suite *DatabaseIntegrationTestSuite) TestModelHooks() {
	// 创建带钩子的用户
	user := &UserWithHooks{
		Name:     "Hook User",
		Email:    "hook@example.com",
		Password: "password",
		Age:      25,
	}

	// 保存用户（应该触发BeforeSave和AfterSave钩子）
	err := user.Save(suite.db)
	suite.NoError(err)
	suite.True(user.HookFlag) // BeforeSave钩子设置的标志

	// 删除用户（应该触发BeforeDelete和AfterDelete钩子）
	err = user.Delete(suite.db)
	suite.NoError(err)
}

// TestTransactions 测试事务
func (suite *DatabaseIntegrationTestSuite) TestTransactions() {
	// 开始事务
	tx, err := suite.db.Begin()
	suite.NoError(err)

	// 在事务中创建用户
	user := &User{
		Name:     "Transaction User",
		Email:    "tx@example.com",
		Password: "password",
		Age:      25,
	}
	err = user.Save(tx)
	suite.NoError(err)

	// 提交事务
	err = tx.Commit()
	suite.NoError(err)

	// 验证用户已保存
	foundUser := &User{}
	err = foundUser.Find(suite.db, user.ID)
	suite.NoError(err)
	suite.Equal(user.Name, foundUser.Name)

	// 测试回滚
	tx, err = suite.db.Begin()
	suite.NoError(err)

	rollbackUser := &User{
		Name:     "Rollback User",
		Email:    "rollback@example.com",
		Password: "password",
		Age:      30,
	}
	err = rollbackUser.Save(tx)
	suite.NoError(err)

	// 回滚事务
	err = tx.Rollback()
	suite.NoError(err)

	// 验证用户未保存
	err = foundUser.Find(suite.db, rollbackUser.ID)
	suite.Error(err)
}

// TestRawQueries 测试原始查询
func (suite *DatabaseIntegrationTestSuite) TestRawQueries() {
	// 创建测试数据
	users := []*User{
		{Name: "Raw1", Email: "raw1@example.com", Password: "pass", Age: 20},
		{Name: "Raw2", Email: "raw2@example.com", Password: "pass", Age: 25},
		{Name: "Raw3", Email: "raw3@example.com", Password: "pass", Age: 30},
	}

	for _, user := range users {
		err := user.Save(suite.db)
		suite.NoError(err)
	}

	// 执行原始查询
	rows, err := suite.db.Query("SELECT name, age FROM users WHERE age >= ? ORDER BY age", 25)
	suite.NoError(err)
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var name string
		var age int
		err := rows.Scan(&name, &age)
		suite.NoError(err)
		results = append(results, map[string]interface{}{
			"name": name,
			"age":  age,
		})
	}

	suite.Len(results, 2)
	suite.Equal("Raw2", results[0]["name"])
	suite.Equal("Raw3", results[1]["name"])
}

// TestModelValidation 测试模型验证
func (suite *DatabaseIntegrationTestSuite) TestModelValidation() {
	// 测试空名称
	user := &User{
		Name:     "",
		Email:    "test@example.com",
		Password: "password",
		Age:      25,
	}

	err := user.Validate()
	suite.Error(err)

	// 测试无效邮箱
	user.Name = "Test User"
	user.Email = "invalid-email"
	err = user.Validate()
	suite.Error(err)

	// 测试有效数据
	user.Email = "valid@example.com"
	err = user.Validate()
	suite.NoError(err)
}

// 运行数据库集成测试套件
func TestDatabaseIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationTestSuite))
}

// UserWithHooks 带钩子的用户模型
type UserWithHooks struct {
	database.Model
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Age      int    `db:"age"`
	HookFlag bool   `db:"hook_flag"`
}

func (u *UserWithHooks) TableName() string {
	return "users"
}

func (u *UserWithHooks) BeforeSave(conn database.Connection) error {
	u.HookFlag = true
	return nil
}

func (u *UserWithHooks) AfterSave(conn database.Connection) error {
	// 可以在这里执行保存后的逻辑
	return nil
}

func (u *UserWithHooks) BeforeDelete(conn database.Connection) error {
	// 可以在这里执行删除前的逻辑
	return nil
}

func (u *UserWithHooks) AfterDelete(conn database.Connection) error {
	// 可以在这里执行删除后的逻辑
	return nil
}

// Post 文章模型
type Post struct {
	database.Model
	UserID  int64  `db:"user_id"`
	Title   string `db:"title"`
	Content string `db:"content"`
	Status  string `db:"status"`
}

func (p *Post) TableName() string {
	return "posts"
}

// Comment 评论模型
type Comment struct {
	database.Model
	PostID  int64  `db:"post_id"`
	UserID  int64  `db:"user_id"`
	Content string `db:"content"`
}

func (c *Comment) TableName() string {
	return "comments"
}

// Tag 标签模型
type Tag struct {
	database.Model
	Name string `db:"name"`
}

func (t *Tag) TableName() string {
	return "tags"
}
