package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"laravel-go/framework/microservice"
)

// 导入多注册中心演示函数
func demonstrateMultiRegistry() {
	fmt.Println("\n=== 多注册中心演示 ===")

	// 1. 内存注册中心演示
	fmt.Println("\n1. 内存注册中心演示")
	demonstrateMemoryRegistry()

	// 2. 构建器模式演示
	fmt.Println("\n2. 构建器模式演示")
	demonstrateRegistryBuilder()

	// 3. 配置模式演示
	fmt.Println("\n3. 配置模式演示")
	demonstrateConfigMode()
}

// 演示内存注册中心
func demonstrateMemoryRegistry() {
	// 创建内存注册中心
	registry, err := microservice.NewServiceRegistry(&microservice.RegistryConfig{
		Type: microservice.RegistryTypeMemory,
	})
	if err != nil {
		log.Printf("创建内存注册中心失败: %v", err)
		return
	}

	// 创建服务发现
	discovery := microservice.NewServiceDiscovery(registry, microservice.NewRoundRobinLoadBalancer())

	// 注册服务
	service := &microservice.ServiceInfo{
		ID:       "user-service-1",
		Name:     "user-service",
		Version:  "v1.0.0",
		Address:  "localhost",
		Port:     8081,
		Protocol: "http",
		Health:   "healthy",
		Metadata: map[string]string{
			"environment": "development",
			"region":      "us-east-1",
		},
		Tags: []string{"api", "user"},
	}

	ctx := context.Background()
	err = registry.Register(ctx, service)
	if err != nil {
		log.Printf("注册服务失败: %v", err)
		return
	}

	fmt.Printf("✅ 服务已注册: %s (%s:%d)\n", service.Name, service.Address, service.Port)

	// 发现服务
	services, err := discovery.Discover(ctx, "user-service")
	if err != nil {
		log.Printf("发现服务失败: %v", err)
		return
	}

	fmt.Printf("✅ 发现 %d 个服务实例\n", len(services))
	for _, s := range services {
		fmt.Printf("  - %s (%s:%d) [%s]\n", s.ID, s.Address, s.Port, s.Health)
	}

	// 关闭注册中心
	registry.Close()
}

// 演示构建器模式
func demonstrateRegistryBuilder() {
	// 使用构建器创建 etcd 注册中心（仅演示配置，不实际连接）
	registry, err := microservice.NewRegistryBuilder().
		WithType(microservice.RegistryTypeEtcd).
		WithEtcd([]string{"localhost:2379"}, "/laravel-go/services").
		Build()

	if err != nil {
		log.Printf("创建 etcd 注册中心失败: %v", err)
		return
	}

	fmt.Println("✅ 使用构建器模式创建了 etcd 注册中心")

	// 使用构建器创建 Consul 注册中心
	consulRegistry, err := microservice.NewRegistryBuilder().
		WithType(microservice.RegistryTypeConsul).
		WithConsul("localhost:8500", "laravel-go/services").
		Build()

	if err != nil {
		log.Printf("创建 Consul 注册中心失败: %v", err)
		return
	}

	fmt.Println("✅ 使用构建器模式创建了 Consul 注册中心")

	// 使用构建器创建 Nacos 注册中心
	nacosRegistry, err := microservice.NewRegistryBuilder().
		WithType(microservice.RegistryTypeNacos).
		WithNacos("localhost:8848", "public", "DEFAULT_GROUP").
		Build()

	if err != nil {
		log.Printf("创建 Nacos 注册中心失败: %v", err)
		return
	}

	fmt.Println("✅ 使用构建器模式创建了 Nacos 注册中心")

	// 使用构建器创建 Zookeeper 注册中心
	zkRegistry, err := microservice.NewRegistryBuilder().
		WithType(microservice.RegistryTypeZookeeper).
		WithZookeeper([]string{"localhost:2181"}, "/laravel-go/services").
		Build()

	if err != nil {
		log.Printf("创建 Zookeeper 注册中心失败: %v", err)
		return
	}

	fmt.Println("✅ 使用构建器模式创建了 Zookeeper 注册中心")

	// 关闭所有注册中心
	registry.Close()
	consulRegistry.Close()
	nacosRegistry.Close()
	zkRegistry.Close()
}

// 演示配置模式
func demonstrateConfigMode() {
	// 创建 etcd 配置
	etcdConfig := &microservice.RegistryConfig{
		Type: microservice.RegistryTypeEtcd,
		Etcd: &microservice.EtcdConfig{
			Endpoints: []string{"localhost:2379", "localhost:2380"},
			Username:  "admin",
			Password:  "password",
			Prefix:    "/laravel-go/services",
			TTL:       30 * time.Second,
		},
	}

	// 创建 Consul 配置
	consulConfig := &microservice.RegistryConfig{
		Type: microservice.RegistryTypeConsul,
		Consul: &microservice.ConsulConfig{
			Address:    "localhost:8500",
			Token:      "consul-token",
			Datacenter: "dc1",
			Prefix:     "laravel-go/services",
			TTL:        30 * time.Second,
		},
	}

	// 创建 Nacos 配置
	nacosConfig := &microservice.RegistryConfig{
		Type: microservice.RegistryTypeNacos,
		Nacos: &microservice.NacosConfig{
			ServerAddr: "localhost:8848",
			Namespace:  "public",
			Group:      "DEFAULT_GROUP",
			Username:   "nacos",
			Password:   "nacos",
			TTL:        30 * time.Second,
		},
	}

	// 创建 Zookeeper 配置
	zkConfig := &microservice.RegistryConfig{
		Type: microservice.RegistryTypeZookeeper,
		Zookeeper: &microservice.ZookeeperConfig{
			Servers:        []string{"localhost:2181", "localhost:2182"},
			Prefix:         "/laravel-go/services",
			TTL:            30 * time.Second,
			SessionTimeout: 10 * time.Second,
		},
	}

	// 创建注册中心（仅演示配置，不实际连接）
	registries := []struct {
		name   string
		config *microservice.RegistryConfig
	}{
		{"Etcd", etcdConfig},
		{"Consul", consulConfig},
		{"Nacos", nacosConfig},
		{"Zookeeper", zkConfig},
	}

	for _, r := range registries {
		registry, err := microservice.NewServiceRegistry(r.config)
		if err != nil {
			log.Printf("创建 %s 注册中心失败: %v", r.name, err)
			continue
		}

		fmt.Printf("✅ 使用配置模式创建了 %s 注册中心\n", r.name)
		registry.Close()
	}
}

// 演示服务客户端选项
func demonstrateServiceClientOptions() {
	fmt.Println("\n4. 服务客户端选项演示")

	// 创建内存注册中心和服务发现
	registry, _ := microservice.NewServiceRegistry(&microservice.RegistryConfig{
		Type: microservice.RegistryTypeMemory,
	})
	discovery := microservice.NewServiceDiscovery(registry, microservice.NewRoundRobinLoadBalancer())

	// 使用选项创建服务客户端
	_ = microservice.NewServiceClient(
		discovery,
		microservice.WithTimeout(10*time.Second),
		microservice.WithRetry(5, 2*time.Second),
	)

	fmt.Printf("✅ 创建了服务客户端:\n")
	fmt.Printf("  - 超时时间: 10s\n")
	fmt.Printf("  - 重试次数: 5\n")
	fmt.Printf("  - 重试延迟: 2s\n")

	// 关闭资源
	registry.Close()
}

// 演示负载均衡器
func demonstrateLoadBalancers() {
	fmt.Println("\n5. 负载均衡器演示")

	// 创建内存注册中心
	registry, _ := microservice.NewServiceRegistry(&microservice.RegistryConfig{
		Type: microservice.RegistryTypeMemory,
	})

	// 注册多个服务实例
	services := []*microservice.ServiceInfo{
		{
			ID:       "user-service-1",
			Name:     "user-service",
			Address:  "localhost",
			Port:     8081,
			Protocol: "http",
			Health:   "healthy",
		},
		{
			ID:       "user-service-2",
			Name:     "user-service",
			Address:  "localhost",
			Port:     8082,
			Protocol: "http",
			Health:   "healthy",
		},
		{
			ID:       "user-service-3",
			Name:     "user-service",
			Address:  "localhost",
			Port:     8083,
			Protocol: "http",
			Health:   "healthy",
		},
	}

	ctx := context.Background()
	for _, service := range services {
		registry.Register(ctx, service)
	}

	// 测试轮询负载均衡器
	fmt.Println("\n轮询负载均衡器:")
	roundRobinDiscovery := microservice.NewServiceDiscovery(registry, microservice.NewRoundRobinLoadBalancer())
	for i := 0; i < 6; i++ {
		service, _ := roundRobinDiscovery.DiscoverOne(ctx, "user-service")
		if service != nil {
			fmt.Printf("  选择服务: %s:%d\n", service.Address, service.Port)
		}
	}

	// 测试随机负载均衡器
	fmt.Println("\n随机负载均衡器:")
	randomDiscovery := microservice.NewServiceDiscovery(registry, microservice.NewRandomLoadBalancer())
	for i := 0; i < 6; i++ {
		service, _ := randomDiscovery.DiscoverOne(ctx, "user-service")
		if service != nil {
			fmt.Printf("  选择服务: %s:%d\n", service.Address, service.Port)
		}
	}

	registry.Close()
}

// User 用户模型
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Order 订单模型
type Order struct {
	ID     int     `json:"id"`
	UserID int     `json:"user_id"`
	Status string  `json:"status"`
	Total  float64 `json:"total"`
}

func main() {
	fmt.Println("=== Laravel-Go 微服务系统演示 ===\n")

	// 1. 演示服务注册和发现
	demoServiceRegistry()

	// 2. 演示负载均衡
	demoLoadBalancing()

	// 3. 演示服务通信
	demoServiceCommunication()

	// 4. 演示熔断器
	demoCircuitBreaker()

	// 5. 多注册中心演示
	demonstrateMultiRegistry()
	demonstrateServiceClientOptions()
	demonstrateLoadBalancers()

	// 6. 启动微服务演示服务器
	startMicroserviceDemo()
}

// demoServiceRegistry 演示服务注册和发现
func demoServiceRegistry() {
	fmt.Println("1. 服务注册和发现演示")
	fmt.Println("======================")

	// 创建服务注册中心
	registry := microservice.NewMemoryServiceRegistry()
	discovery := microservice.NewMemoryServiceDiscovery(registry, microservice.NewRoundRobinLoadBalancer())
	ctx := context.Background()

	// 注册用户服务实例
	userServices := []*microservice.ServiceInfo{
		{
			ID:       "user-service-1",
			Name:     "user-service",
			Version:  "1.0.0",
			Address:  "localhost",
			Port:     8081,
			Protocol: "http",
			Health:   "healthy",
			Metadata: map[string]string{
				"environment": "production",
				"region":      "us-west-1",
			},
			Tags: []string{"api", "user"},
		},
		{
			ID:       "user-service-2",
			Name:     "user-service",
			Version:  "1.0.0",
			Address:  "localhost",
			Port:     8082,
			Protocol: "http",
			Health:   "healthy",
			Metadata: map[string]string{
				"environment": "production",
				"region":      "us-west-1",
			},
			Tags: []string{"api", "user"},
		},
		{
			ID:       "user-service-3",
			Name:     "user-service",
			Version:  "1.0.0",
			Address:  "localhost",
			Port:     8083,
			Protocol: "http",
			Health:   "unhealthy", // 不健康的实例
			Metadata: map[string]string{
				"environment": "production",
				"region":      "us-west-1",
			},
			Tags: []string{"api", "user"},
		},
	}

	// 注册订单服务实例
	orderServices := []*microservice.ServiceInfo{
		{
			ID:       "order-service-1",
			Name:     "order-service",
			Version:  "1.0.0",
			Address:  "localhost",
			Port:     8084,
			Protocol: "http",
			Health:   "healthy",
			Metadata: map[string]string{
				"environment": "production",
				"region":      "us-west-1",
			},
			Tags: []string{"api", "order"},
		},
		{
			ID:       "order-service-2",
			Name:     "order-service",
			Version:  "1.0.0",
			Address:  "localhost",
			Port:     8085,
			Protocol: "http",
			Health:   "healthy",
			Metadata: map[string]string{
				"environment": "production",
				"region":      "us-west-1",
			},
			Tags: []string{"api", "order"},
		},
	}

	// 注册所有服务
	fmt.Println("\n注册服务实例:")
	for _, service := range append(userServices, orderServices...) {
		err := registry.Register(ctx, service)
		if err != nil {
			fmt.Printf("注册服务 %s 失败: %v\n", service.ID, err)
		} else {
			fmt.Printf("✓ 注册服务: %s (%s:%d)\n", service.ID, service.Address, service.Port)
		}
	}

	// 列出所有服务
	fmt.Println("\n所有注册的服务:")
	allServices, err := registry.ListServices(ctx)
	if err != nil {
		fmt.Printf("获取服务列表失败: %v\n", err)
	} else {
		for _, service := range allServices {
			fmt.Printf("- %s (%s) - %s\n", service.Name, service.ID, service.Health)
		}
	}

	// 发现用户服务
	fmt.Println("\n发现用户服务:")
	userServiceInstances, err := discovery.Discover(ctx, "user-service")
	if err != nil {
		fmt.Printf("发现用户服务失败: %v\n", err)
	} else {
		fmt.Printf("找到 %d 个用户服务实例:\n", len(userServiceInstances))
		for _, instance := range userServiceInstances {
			fmt.Printf("- %s (%s:%d) - %s\n", instance.ID, instance.Address, instance.Port, instance.Health)
		}
	}

	// 发现订单服务
	fmt.Println("\n发现订单服务:")
	orderServiceInstances, err := discovery.Discover(ctx, "order-service")
	if err != nil {
		fmt.Printf("发现订单服务失败: %v\n", err)
	} else {
		fmt.Printf("找到 %d 个订单服务实例:\n", len(orderServiceInstances))
		for _, instance := range orderServiceInstances {
			fmt.Printf("- %s (%s:%d) - %s\n", instance.ID, instance.Address, instance.Port, instance.Health)
		}
	}

	// 测试负载均衡选择
	fmt.Println("\n负载均衡选择:")
	for i := 0; i < 5; i++ {
		selected, err := discovery.DiscoverOne(ctx, "user-service")
		if err != nil {
			fmt.Printf("选择用户服务失败: %v\n", err)
		} else {
			fmt.Printf("第 %d 次选择: %s (%s:%d)\n", i+1, selected.ID, selected.Address, selected.Port)
		}
	}

	// 测试缓存统计
	fmt.Println("\n缓存统计:")
	stats := discovery.GetCacheStats()
	for serviceName, count := range stats {
		fmt.Printf("- %s: %d 个实例\n", serviceName, count)
	}

	fmt.Println()
}

// demoLoadBalancing 演示负载均衡
func demoLoadBalancing() {
	fmt.Println("2. 负载均衡演示")
	fmt.Println("===============")

	services := []*microservice.ServiceInfo{
		{ID: "service-1", Name: "test-service", Health: "healthy"},
		{ID: "service-2", Name: "test-service", Health: "healthy"},
		{ID: "service-3", Name: "test-service", Health: "unhealthy"},
		{ID: "service-4", Name: "test-service", Health: "healthy"},
	}

	// 测试轮询负载均衡器
	fmt.Println("\n轮询负载均衡器:")
	rr := microservice.NewRoundRobinLoadBalancer()
	for i := 0; i < 6; i++ {
		selected := rr.Select(services)
		if selected != nil {
			fmt.Printf("第 %d 次选择: %s\n", i+1, selected.ID)
		} else {
			fmt.Printf("第 %d 次选择: 无可用服务\n", i+1)
		}
	}

	// 测试随机负载均衡器
	fmt.Println("\n随机负载均衡器:")
	random := microservice.NewRandomLoadBalancer()
	for i := 0; i < 5; i++ {
		selected := random.Select(services)
		if selected != nil {
			fmt.Printf("第 %d 次选择: %s\n", i+1, selected.ID)
		} else {
			fmt.Printf("第 %d 次选择: 无可用服务\n", i+1)
		}
	}

	fmt.Println()
}

// demoServiceCommunication 演示服务通信
func demoServiceCommunication() {
	fmt.Println("3. 服务通信演示")
	fmt.Println("===============")

	// 创建服务注册中心和发现服务
	registry := microservice.NewMemoryServiceRegistry()
	discovery := microservice.NewMemoryServiceDiscovery(registry, microservice.NewRoundRobinLoadBalancer())
	ctx := context.Background()

	// 注册模拟服务
	service := &microservice.ServiceInfo{
		ID:       "api-service",
		Name:     "api-service",
		Address:  "localhost",
		Port:     8086,
		Protocol: "http",
		Health:   "healthy",
		Metadata: map[string]string{
			"version": "1.0.0",
		},
	}

	err := registry.Register(ctx, service)
	if err != nil {
		fmt.Printf("注册服务失败: %v\n", err)
		return
	}

	// 创建服务客户端
	client := microservice.NewServiceClient(
		discovery,
		microservice.WithTimeout(5*time.Second),
		microservice.WithRetry(3, 1*time.Second),
	)

	// 模拟服务调用
	fmt.Println("\n模拟服务调用:")

	// GET 请求
	fmt.Println("发送 GET 请求...")
	response, err := client.Get(ctx, "api-service", "/users")
	if err != nil {
		fmt.Printf("GET 请求失败: %v\n", err)
	} else {
		fmt.Printf("GET 响应: %s\n", string(response))
	}

	// POST 请求
	fmt.Println("发送 POST 请求...")
	user := User{ID: 1, Name: "张三", Email: "zhangsan@example.com"}
	response, err = client.Post(ctx, "api-service", "/users", user)
	if err != nil {
		fmt.Printf("POST 请求失败: %v\n", err)
	} else {
		fmt.Printf("POST 响应: %s\n", string(response))
	}

	// JSON 请求
	fmt.Println("发送 JSON 请求...")
	var responseUser User
	err = client.PostJSON(ctx, "api-service", "/users", user, &responseUser)
	if err != nil {
		fmt.Printf("JSON 请求失败: %v\n", err)
	} else {
		fmt.Printf("JSON 响应: %+v\n", responseUser)
	}

	fmt.Println()
}

// demoCircuitBreaker 演示熔断器
func demoCircuitBreaker() {
	fmt.Println("4. 熔断器演示")
	fmt.Println("=============")

	// 创建熔断器
	cb := microservice.NewSimpleCircuitBreaker(3, 5*time.Second)
	ctx := context.Background()

	// 模拟成功操作
	fmt.Println("\n测试成功操作:")
	successCount := 0
	successOperation := func() error {
		successCount++
		fmt.Printf("执行成功操作 #%d\n", successCount)
		return nil
	}

	for i := 0; i < 3; i++ {
		err := cb.Execute(ctx, successOperation)
		if err != nil {
			fmt.Printf("操作失败: %v\n", err)
		} else {
			fmt.Printf("操作成功\n")
		}
	}

	// 模拟失败操作
	fmt.Println("\n测试失败操作:")
	failureCount := 0
	failingOperation := func() error {
		failureCount++
		fmt.Printf("执行失败操作 #%d\n", failureCount)
		return fmt.Errorf("操作失败")
	}

	for i := 0; i < 5; i++ {
		err := cb.Execute(ctx, failingOperation)
		if err != nil {
			fmt.Printf("操作失败: %v\n", err)
		} else {
			fmt.Printf("操作成功\n")
		}
	}

	// 检查熔断器状态
	fmt.Printf("\n熔断器状态: %s\n", func() string {
		if cb.IsOpen() {
			return "开启"
		}
		return "关闭"
	}())

	// 重置熔断器
	fmt.Println("重置熔断器...")
	cb.Reset()
	fmt.Printf("重置后熔断器状态: %s\n", func() string {
		if cb.IsOpen() {
			return "开启"
		}
		return "关闭"
	}())

	fmt.Println()
}

// startMicroserviceDemo 启动微服务演示服务器
func startMicroserviceDemo() {
	fmt.Println("5. 启动微服务演示服务器")
	fmt.Println("=======================")

	// 创建服务注册中心
	registry := microservice.NewMemoryServiceRegistry()
	ctx := context.Background()

	// 启动清理工作协程
	registry.StartCleanupWorker(10 * time.Second)

	// 注册演示服务
	demoService := &microservice.ServiceInfo{
		ID:       "demo-service",
		Name:     "demo-service",
		Address:  "localhost",
		Port:     8087,
		Protocol: "http",
		Health:   "healthy",
		TTL:      30 * time.Second,
	}

	err := registry.Register(ctx, demoService)
	if err != nil {
		fmt.Printf("注册演示服务失败: %v\n", err)
		return
	}

	// 设置 HTTP 路由
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			users := []User{
				{ID: 1, Name: "张三", Email: "zhangsan@example.com"},
				{ID: 2, Name: "李四", Email: "lisi@example.com"},
			}
			json.NewEncoder(w).Encode(users)
		case "POST":
			var user User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			user.ID = 999 // 模拟分配 ID
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(user)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			orders := []Order{
				{ID: 1, UserID: 1, Status: "pending", Total: 100.50},
				{ID: 2, UserID: 2, Status: "completed", Total: 200.75},
			}
			json.NewEncoder(w).Encode(orders)
		case "POST":
			var order Order
			if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			order.ID = 999 // 模拟分配 ID
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(order)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		services, err := registry.ListServices(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(services)
	})

	fmt.Println("服务器启动在 http://localhost:8087")
	fmt.Println("健康检查: http://localhost:8087/health")
	fmt.Println("用户服务: http://localhost:8087/users")
	fmt.Println("订单服务: http://localhost:8087/orders")
	fmt.Println("服务列表: http://localhost:8087/services")
	fmt.Println("\n按 Ctrl+C 停止服务器")

	log.Fatal(http.ListenAndServe(":8087", nil))
}
