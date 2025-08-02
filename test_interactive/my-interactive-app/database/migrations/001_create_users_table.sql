-- Migration: Create Users Table
-- Description: 创建用户表
-- Version: 1.0

-- UP Migration
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入示例数据
INSERT INTO users (name, email, password) VALUES 
('John Doe', 'john@example.com', 'hashed_password_1'),
('Jane Smith', 'jane@example.com', 'hashed_password_2');

-- DOWN Migration (如果需要回滚)
-- DROP TABLE IF EXISTS users;