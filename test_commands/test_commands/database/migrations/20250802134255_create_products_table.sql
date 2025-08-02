-- Migration: create_products_table
-- Description: Create products table
-- Version: 20250802134255

-- UP Migration
CREATE TABLE IF NOT EXISTS products (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- DOWN Migration
DROP TABLE IF EXISTS products;
