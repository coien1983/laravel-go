package config

func init() {
	Config.Set("logging", map[string]interface{}{
		"default": "stack",
		"deprecations": map[string]interface{}{
			"channel": "null",
			"trace":   false,
		},
		"channels": map[string]interface{}{
			"stack": map[string]interface{}{
				"driver":   "stack",
				"channels": []string{"single"},
				"ignore_exceptions": false,
			},
			"single": map[string]interface{}{
				"driver": "single",
				"path":   "storage/logs/laravel.log",
				"level":  getEnv("LOG_LEVEL", "debug"),
			},
			"daily": map[string]interface{}{
				"driver": "daily",
				"path":   "storage/logs/laravel.log",
				"level":  getEnv("LOG_LEVEL", "debug"),
				"days":   14,
			},
			"slack": map[string]interface{}{
				"driver":   "slack",
				"url":      getEnv("LOG_SLACK_WEBHOOK_URL", ""),
				"username": "Laravel Log",
				"emoji":    ":boom:",
				"level":    getEnv("LOG_LEVEL", "critical"),
			},
			"papertrail": map[string]interface{}{
				"driver":       "monolog",
				"level":        getEnv("LOG_LEVEL", "debug"),
				"handler":      "Monolog\\Handler\\SyslogUdpHandler",
				"handler_with": map[string]interface{}{
					"host": getEnv("PAPERTRAIL_URL", ""),
					"port": getEnv("PAPERTRAIL_PORT", ""),
				},
			},
			"stderr": map[string]interface{}{
				"driver":    "monolog",
				"level":     getEnv("LOG_LEVEL", "debug"),
				"handler":   "Monolog\\Handler\\StreamHandler",
				"formatter": "Monolog\\Formatter\\JsonFormatter",
				"with": map[string]interface{}{
					"stream": "php://stderr",
				},
			},
			"syslog": map[string]interface{}{
				"driver": "syslog",
				"level":  getEnv("LOG_LEVEL", "debug"),
			},
			"errorlog": map[string]interface{}{
				"driver": "errorlog",
				"level":  getEnv("LOG_LEVEL", "debug"),
			},
			"null": map[string]interface{}{
				"driver":  "monolog",
				"handler": "Monolog\\Handler\\NullHandler",
			},
			"emergency": map[string]interface{}{
				"path": "storage/logs/laravel.log",
			},
		},
	})
} 