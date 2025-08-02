package config

func init() {
	Config.Set("queue", map[string]interface{}{
		"default": "redis",
		"connections": map[string]interface{}{
			"sync": map[string]interface{}{
				"driver": "sync",
			},
			"database": map[string]interface{}{
				"driver": "database",
				"table":  "jobs",
				"queue":  "default",
				"retry_after": 90,
				"after_commit": false,
			},
			"redis": map[string]interface{}{
				"driver":     "redis",
				"connection": "default",
				"queue":      getEnv("REDIS_QUEUE", "default"),
				"retry_after": 90,
				"block_for":   nil,
				"after_commit": false,
			},
			"sqs": map[string]interface{}{
				"driver": "sqs",
				"key":    getEnv("AWS_ACCESS_KEY_ID", ""),
				"secret": getEnv("AWS_SECRET_ACCESS_KEY", ""),
				"prefix": getEnv("SQS_PREFIX", "https://sqs.us-east-1.amazonaws.com/your-account-id"),
				"queue":  getEnv("SQS_QUEUE", "default"),
				"suffix": getEnv("SQS_SUFFIX", ""),
				"region": getEnv("AWS_DEFAULT_REGION", "us-east-1"),
				"group":  "default",
				"deduplicator": "sqs",
				"allow_delay_failures": false,
				"after_commit": false,
			},
			"beanstalkd": map[string]interface{}{
				"driver": "beanstalkd",
				"host":   getEnv("BEANSTALKD_HOST", "127.0.0.1"),
				"queue":  getEnv("BEANSTALKD_QUEUE", "default"),
				"retry_after": 90,
				"block_for": 0,
				"after_commit": false,
			},
		},
		"failed": map[string]interface{}{
			"driver": "database-uuids",
			"database": "pgsql",
			"table":    "failed_jobs",
		},
	})
} 