package config

func init() {
	Config.Set("session", map[string]interface{}{
		"driver": getEnv("SESSION_DRIVER", "file"),
		"lifetime": getEnv("SESSION_LIFETIME", "120"),
		"expire_on_close": false,
		"encrypt": false,
		"files": "storage/framework/sessions",
		"connection": getEnv("SESSION_CONNECTION", ""),
		"table": "sessions",
		"store": getEnv("SESSION_STORE", ""),
		"lottery": []int{2, 100},
		"cookie": getEnv("SESSION_COOKIE", "laravel_go_session"),
		"path": "/",
		"domain": getEnv("SESSION_DOMAIN", ""),
		"secure": getEnv("SESSION_SECURE_COOKIE", "false"),
		"http_only": true,
		"same_site": "lax",
	})
} 