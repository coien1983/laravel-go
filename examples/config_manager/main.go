package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"laravel-go/framework/config"
	"os"
	"strings"
)

func main() {
	var (
		configFile = flag.String("config", "", "配置文件路径")
		key        = flag.String("key", "", "配置键")
		value      = flag.String("value", "", "配置值")
		action     = flag.String("action", "get", "操作类型: get, set, list, validate")
		format     = flag.String("format", "json", "输出格式: json, yaml, env")
	)
	flag.Parse()

	// 创建配置管理器
	cfg := config.NewConfig()

	// 加载配置文件
	if *configFile != "" {
		if err := cfg.LoadFromFile(*configFile); err != nil {
			fmt.Printf("加载配置文件失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 加载环境变量
	if err := cfg.LoadEnv(); err != nil {
		fmt.Printf("加载环境变量失败: %v\n", err)
		os.Exit(1)
	}

	switch *action {
	case "get":
		if *key == "" {
			fmt.Println("请指定配置键")
			os.Exit(1)
		}
		value := cfg.Get(*key)
		outputValue(value, *format)

	case "set":
		if *key == "" || *value == "" {
			fmt.Println("请指定配置键和值")
			os.Exit(1)
		}
		cfg.Set(*key, *value)
		fmt.Printf("配置已设置: %s = %s\n", *key, *value)

	case "list":
		allConfig := cfg.All()
		outputValue(allConfig, *format)

	case "validate":
		rules := map[string]string{
			"app.name":      "required",
			"app.version":   "required",
			"app.port":      "required|numeric",
			"app.debug":     "required|boolean",
			"database.host": "required",
			"database.port": "required|numeric",
		}
		if err := cfg.Validate(rules); err != nil {
			fmt.Printf("配置验证失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("配置验证通过")

	default:
		fmt.Printf("未知操作: %s\n", *action)
		flag.Usage()
		os.Exit(1)
	}
}

func outputValue(value interface{}, format string) {
	switch format {
	case "json":
		data, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			fmt.Printf("JSON序列化失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))

	case "yaml":
		// 简单的YAML输出
		outputYAML(value, "")

	case "env":
		outputEnv(value, "")

	default:
		fmt.Printf("未知格式: %s\n", format)
		os.Exit(1)
	}
}

func outputYAML(value interface{}, prefix string) {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			fullKey := key
			if prefix != "" {
				fullKey = prefix + "." + key
			}
			outputYAML(val, fullKey)
		}
	case []interface{}:
		for i, val := range v {
			fullKey := fmt.Sprintf("%s[%d]", prefix, i)
			outputYAML(val, fullKey)
		}
	default:
		fmt.Printf("%s: %v\n", prefix, v)
	}
}

func outputEnv(value interface{}, prefix string) {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			fullKey := strings.ToUpper(key)
			if prefix != "" {
				fullKey = strings.ToUpper(prefix + "_" + key)
			}
			outputEnv(val, fullKey)
		}
	case []interface{}:
		for i, val := range v {
			fullKey := fmt.Sprintf("%s_%d", prefix, i)
			outputEnv(val, fullKey)
		}
	default:
		fmt.Printf("%s=%v\n", prefix, v)
	}
}
