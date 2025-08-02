package config

import "laravel-go/framework/config"

func init() {
	config.Set("app", map[string]interface{}{
		"name":     "Laravel-Go 博客系统",
		"version":  "1.0.0",
		"port":     "8080",
		"debug":    true,
		"timezone": "Asia/Shanghai",
	})

	config.Set("database", map[string]interface{}{
		"default": "sqlite",
		"connections": map[string]interface{}{
			"sqlite": map[string]interface{}{
				"driver":   "sqlite",
				"database": "blog.db",
			},
		},
	})

	config.Set("cache", map[string]interface{}{
		"default": "memory",
		"stores": map[string]interface{}{
			"memory": map[string]interface{}{
				"driver": "memory",
			},
		},
	})

	config.Set("queue", map[string]interface{}{
		"default": "memory",
		"connections": map[string]interface{}{
			"memory": map[string]interface{}{
				"driver": "memory",
			},
		},
	})
}
