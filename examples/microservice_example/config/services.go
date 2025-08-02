package config

func init() {
	Config.Set("services", map[string]interface{}{
		"gateway": map[string]interface{}{
			"name": "api-gateway",
			"port": getEnv("GATEWAY_PORT", "8080"),
			"host": getEnv("GATEWAY_HOST", "localhost"),
		},
		"user_service": map[string]interface{}{
			"name": "user-service",
			"port": getEnv("USER_SERVICE_PORT", "8082"),
			"host": getEnv("USER_SERVICE_HOST", "localhost"),
		},
		"product_service": map[string]interface{}{
			"name": "product-service",
			"port": getEnv("PRODUCT_SERVICE_PORT", "8083"),
			"host": getEnv("PRODUCT_SERVICE_HOST", "localhost"),
		},
		"order_service": map[string]interface{}{
			"name": "order-service",
			"port": getEnv("ORDER_SERVICE_PORT", "8084"),
			"host": getEnv("ORDER_SERVICE_HOST", "localhost"),
		},
		"registry": map[string]interface{}{
			"type": getEnv("REGISTRY_TYPE", "memory"),
			"consul": map[string]interface{}{
				"host": getEnv("CONSUL_HOST", "localhost"),
				"port": getEnv("CONSUL_PORT", "8500"),
			},
			"etcd": map[string]interface{}{
				"endpoints": []string{getEnv("ETCD_ENDPOINTS", "localhost:2379")},
			},
			"zookeeper": map[string]interface{}{
				"servers": []string{getEnv("ZOOKEEPER_SERVERS", "localhost:2181")},
			},
		},
		"load_balancer": map[string]interface{}{
			"type": getEnv("LOAD_BALANCER_TYPE", "round_robin"),
			"health_check": map[string]interface{}{
				"enabled": true,
				"interval": 30,
				"timeout": 5,
			},
		},
		"circuit_breaker": map[string]interface{}{
			"enabled": true,
			"failure_threshold": 5,
			"recovery_timeout": 60,
			"half_open_requests": 3,
		},
	})
} 