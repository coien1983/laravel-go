-- Migration: create_users_table
-- Description: Create users table
-- Version: 20250802134250

-- UP Migration
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- DOWN Migration
DROP TABLE IF EXISTS users;
