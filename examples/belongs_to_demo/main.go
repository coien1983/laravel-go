package main

import (
	"fmt"
	"log"

	"laravel-go/framework/database"
)

// User 用户模型
type User struct {
	database.Model
	Name  string `db:"name"`
	Email string `db:"email"`
}

func (u *User) TableName() string {
	return "users"
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

func main() {
	// 创建数据库连接
	config := &database.ConnectionConfig{
		Driver: database.SQLite,
		Host:   "belongs_to_demo.db",
	}

	conn, err := database.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// 创建表
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT,
			status TEXT DEFAULT 'draft',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create posts table: %v", err)
	}

	// 创建用户
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	model := &database.Model{}
	err = model.Save(conn, user)
	if err != nil {
		log.Fatalf("Failed to save user: %v", err)
	}

	fmt.Printf("Created user: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)

	// 创建文章
	post := &Post{
		UserID:  user.ID,
		Title:   "My First Post",
		Content: "This is the content of my first post.",
		Status:  "published",
	}

	err = model.Save(conn, post)
	if err != nil {
		log.Fatalf("Failed to save post: %v", err)
	}

	fmt.Printf("Created post: ID=%d, Title=%s, UserID=%d\n", post.ID, post.Title, post.UserID)

	// 使用 BelongsToModel 获取文章的作者
	var postAuthor User
	err = model.BelongsToModel(conn, post, &postAuthor, "ID", "UserID")
	if err != nil {
		log.Fatalf("Failed to load post author: %v", err)
	}

	fmt.Printf("Post author: ID=%d, Name=%s, Email=%s\n", postAuthor.ID, postAuthor.Name, postAuthor.Email)

	// 创建另一个文章
	post2 := &Post{
		UserID:  user.ID,
		Title:   "My Second Post",
		Content: "This is the content of my second post.",
		Status:  "draft",
	}

	err = model.Save(conn, post2)
	if err != nil {
		log.Fatalf("Failed to save post2: %v", err)
	}

	// 获取第二个文章的作者
	var post2Author User
	err = model.BelongsToModel(conn, post2, &post2Author, "ID", "UserID")
	if err != nil {
		log.Fatalf("Failed to load post2 author: %v", err)
	}

	fmt.Printf("Post2 author: ID=%d, Name=%s, Email=%s\n", post2Author.ID, post2Author.Name, post2Author.Email)

	fmt.Println("\nBelongsToModel demo completed successfully!")
}
