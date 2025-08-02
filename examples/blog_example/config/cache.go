package config

func init() {
	Config.Set("cache", map[string]interface{}{
		"default": "redis",
		"stores": map[string]interface{}{
			"apc": map[string]interface{}{
				"driver": "apc",
			},
			"array": map[string]interface{}{
				"driver": "array",
				"serialize": false,
			},
			"file": map[string]interface{}{
				"driver": "file",
				"path":   "storage/framework/cache/data",
			},
			"redis": map[string]interface{}{
				"driver":     "redis",
				"connection": "cache",
				"lock_connection": "default",
			},
			"database": map[string]interface{}{
				"driver": "database",
				"table":  "cache",
				"connection": "",
				"lock_connection": "",
			},
			"dynamodb": map[string]interface{}{
				"driver": "dynamodb",
				"key":    getEnv("AWS_ACCESS_KEY_ID", ""),
				"secret": getEnv("AWS_SECRET_ACCESS_KEY", ""),
				"region": getEnv("AWS_DEFAULT_REGION", "us-east-1"),
				"table":  getEnv("DYNAMODB_CACHE_TABLE", "cache"),
				"endpoint": getEnv("DYNAMODB_ENDPOINT", ""),
			},
			"memcached": map[string]interface{}{
				"driver": "memcached",
				"persistent_id": getEnv("MEMCACHED_PERSISTENT_ID", ""),
				"sasl": []string{
					getEnv("MEMCACHED_USERNAME", ""),
					getEnv("MEMCACHED_PASSWORD", ""),
				},
				"options": []interface{}{
					// Memcached::OPT_CONNECT_TIMEOUT => 2000,
				},
				"servers": []map[string]interface{}{
					{
						"host":   getEnv("MEMCACHED_HOST", "127.0.0.1"),
						"port":   getEnv("MEMCACHED_PORT", "11211"),
						"weight": 100,
					},
				},
			},
		},
		"prefix": "laravel_go_cache",
	})
} 