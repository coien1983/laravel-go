package config

func init() {
	Config.Set("deployment", map[string]interface{}{
		"environment": getEnv("DEPLOYMENT_ENV", "development"),
		"platform":    getEnv("DEPLOYMENT_PLATFORM", "docker"),
		"docker": map[string]interface{}{
			"image": map[string]interface{}{
				"name":     getEnv("DOCKER_IMAGE_NAME", "laravel-go-app"),
				"tag":      getEnv("DOCKER_IMAGE_TAG", "latest"),
				"registry": getEnv("DOCKER_REGISTRY", ""),
			},
			"compose": map[string]interface{}{
				"file": getEnv("DOCKER_COMPOSE_FILE", "docker-compose.yml"),
			},
		},
		"kubernetes": map[string]interface{}{
			"namespace": getEnv("K8S_NAMESPACE", "laravel-go"),
			"replicas":  getEnv("K8S_REPLICAS", "3"),
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    getEnv("K8S_CPU_REQUEST", "100m"),
					"memory": getEnv("K8S_MEMORY_REQUEST", "128Mi"),
				},
				"limits": map[string]interface{}{
					"cpu":    getEnv("K8S_CPU_LIMIT", "500m"),
					"memory": getEnv("K8S_MEMORY_LIMIT", "512Mi"),
				},
			},
		},
		"monitoring": map[string]interface{}{
			"enabled": getEnv("MONITORING_ENABLED", "true") == "true",
			"prometheus": map[string]interface{}{
				"enabled": getEnv("PROMETHEUS_ENABLED", "true") == "true",
				"port":    getEnv("PROMETHEUS_PORT", "9090"),
			},
			"grafana": map[string]interface{}{
				"enabled": getEnv("GRAFANA_ENABLED", "true") == "true",
				"port":    getEnv("GRAFANA_PORT", "3000"),
			},
		},
		"database": map[string]interface{}{
			"type": getEnv("DB_TYPE", "postgres"),
			"postgres": map[string]interface{}{
				"host":     getEnv("DB_HOST", "localhost"),
				"port":     getEnv("DB_PORT", "5432"),
				"database": getEnv("DB_DATABASE", "laravel_go"),
				"username": getEnv("DB_USERNAME", "laravel_go"),
				"password": getEnv("DB_PASSWORD", "password"),
			},
		},
		"cache": map[string]interface{}{
			"type": getEnv("CACHE_TYPE", "redis"),
			"redis": map[string]interface{}{
				"host":     getEnv("REDIS_HOST", "localhost"),
				"port":     getEnv("REDIS_PORT", "6379"),
				"password": getEnv("REDIS_PASSWORD", ""),
				"database": getEnv("REDIS_DATABASE", "0"),
			},
		},
		"nginx": map[string]interface{}{
			"enabled": getEnv("NGINX_ENABLED", "true") == "true",
			"port":    getEnv("NGINX_PORT", "80"),
			"ssl": map[string]interface{}{
				"enabled": getEnv("NGINX_SSL_ENABLED", "false") == "true",
				"port":    getEnv("NGINX_SSL_PORT", "443"),
			},
		},
	})
} 