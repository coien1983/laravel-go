package config

func init() {
	Config.Set("api", map[string]interface{}{
		"version": "v1",
		"prefix":  "/api",
		"middleware": []string{
			"cors",
			"logging",
		},
		"rate_limit": map[string]interface{}{
			"enabled": true,
			"requests_per_minute": 60,
		},
		"cors": map[string]interface{}{
			"allowed_origins": []string{"*"},
			"allowed_methods": []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			"allowed_headers": []string{"*"},
			"exposed_headers": []string{},
			"allow_credentials": true,
			"max_age": 86400,
		},
		"documentation": map[string]interface{}{
			"enabled": true,
			"path":    "/docs",
			"title":   "Laravel-Go API Documentation",
			"version": "1.0.0",
		},
		"response": map[string]interface{}{
			"format": "json",
			"include_meta": true,
			"include_links": true,
		},
	})
} 