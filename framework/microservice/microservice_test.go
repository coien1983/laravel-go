package microservice

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestServiceInfo(t *testing.T) {
	service := &ServiceInfo{
		ID:       "test-service-1",
		Name:     "test-service",
		Version:  "1.0.0",
		Address:  "localhost",
		Port:     8080,
		Protocol: "http",
		Health:   "healthy",
		Metadata: map[string]string{
			"environment": "test",
			"region":      "us-west-1",
		},
		Tags: []string{"api", "web"},
	}

	if service.ID != "test-service-1" {
		t.Errorf("Expected service ID to be 'test-service-1', got %s", service.ID)
	}

	if service.Name != "test-service" {
		t.Errorf("Expected service name to be 'test-service', got %s", service.Name)
	}

	if len(service.Metadata) != 2 {
		t.Errorf("Expected 2 metadata items, got %d", len(service.Metadata))
	}

	if len(service.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(service.Tags))
	}
}

func TestMemoryServiceRegistry(t *testing.T) {
	registry := NewMemoryServiceRegistry()
	ctx := context.Background()

	// 测试注册服务
	service := &ServiceInfo{
		ID:       "test-service-1",
		Name:     "test-service",
		Version:  "1.0.0",
		Address:  "localhost",
		Port:     8080,
		Protocol: "http",
		Health:   "healthy",
	}

	err := registry.Register(ctx, service)
	if err != nil {
		t.Errorf("Failed to register service: %v", err)
	}

	// 测试获取服务
	retrieved, err := registry.GetService(ctx, "test-service-1")
	if err != nil {
		t.Errorf("Failed to get service: %v", err)
	}

	if retrieved.ID != service.ID {
		t.Errorf("Expected service ID %s, got %s", service.ID, retrieved.ID)
	}

	// 测试列出所有服务
	services, err := registry.ListServices(ctx)
	if err != nil {
		t.Errorf("Failed to list services: %v", err)
	}

	if len(services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(services))
	}

	// 测试更新服务
	service.Health = "unhealthy"
	err = registry.Update(ctx, service)
	if err != nil {
		t.Errorf("Failed to update service: %v", err)
	}

	updated, err := registry.GetService(ctx, "test-service-1")
	if err != nil {
		t.Errorf("Failed to get updated service: %v", err)
	}

	if updated.Health != "unhealthy" {
		t.Errorf("Expected health to be 'unhealthy', got %s", updated.Health)
	}

	// 测试注销服务
	err = registry.Deregister(ctx, "test-service-1")
	if err != nil {
		t.Errorf("Failed to deregister service: %v", err)
	}

	_, err = registry.GetService(ctx, "test-service-1")
	if err == nil {
		t.Error("Expected error when getting deregistered service")
	}

	// 测试关闭注册中心
	err = registry.Close()
	if err != nil {
		t.Errorf("Failed to close registry: %v", err)
	}
}

func TestMemoryServiceDiscovery(t *testing.T) {
	registry := NewMemoryServiceRegistry()
	discovery := NewMemoryServiceDiscovery(registry, NewRoundRobinLoadBalancer())
	ctx := context.Background()

	// 注册多个服务实例
	services := []*ServiceInfo{
		{
			ID:       "service-1",
			Name:     "user-service",
			Version:  "1.0.0",
			Address:  "localhost",
			Port:     8081,
			Protocol: "http",
			Health:   "healthy",
		},
		{
			ID:       "service-2",
			Name:     "user-service",
			Version:  "1.0.0",
			Address:  "localhost",
			Port:     8082,
			Protocol: "http",
			Health:   "healthy",
		},
		{
			ID:       "service-3",
			Name:     "order-service",
			Version:  "1.0.0",
			Address:  "localhost",
			Port:     8083,
			Protocol: "http",
			Health:   "healthy",
		},
	}

	for _, service := range services {
		err := registry.Register(ctx, service)
		if err != nil {
			t.Errorf("Failed to register service %s: %v", service.ID, err)
		}
	}

	// 测试发现服务
	userServices, err := discovery.Discover(ctx, "user-service")
	if err != nil {
		t.Errorf("Failed to discover user-service: %v", err)
	}

	if len(userServices) != 2 {
		t.Errorf("Expected 2 user services, got %d", len(userServices))
	}

	// 测试发现单个服务（负载均衡）
	selected, err := discovery.DiscoverOne(ctx, "user-service")
	if err != nil {
		t.Errorf("Failed to discover one user service: %v", err)
	}

	if selected.Name != "user-service" {
		t.Errorf("Expected service name 'user-service', got %s", selected.Name)
	}

	// 测试缓存统计
	stats := discovery.GetCacheStats()
	if stats["user-service"] != 2 {
		t.Errorf("Expected 2 cached user services, got %d", stats["user-service"])
	}

	// 测试清除缓存
	discovery.ClearCache()
	stats = discovery.GetCacheStats()
	if len(stats) != 0 {
		t.Errorf("Expected empty cache stats, got %d items", len(stats))
	}

	// 测试关闭发现服务
	err = discovery.Close()
	if err != nil {
		t.Errorf("Failed to close discovery: %v", err)
	}
}

func TestLoadBalancers(t *testing.T) {
	services := []*ServiceInfo{
		{ID: "1", Name: "service", Health: "healthy"},
		{ID: "2", Name: "service", Health: "healthy"},
		{ID: "3", Name: "service", Health: "unhealthy"},
		{ID: "4", Name: "service", Health: "healthy"},
	}

	// 测试轮询负载均衡器
	rr := NewRoundRobinLoadBalancer()

	// 应该只选择健康服务
	selected1 := rr.Select(services)
	if selected1 == nil {
		t.Error("Expected to select a healthy service")
	}

	selected2 := rr.Select(services)
	if selected2 == nil {
		t.Error("Expected to select a healthy service")
	}

	selected3 := rr.Select(services)
	if selected3 == nil {
		t.Error("Expected to select a healthy service")
	}

	// 验证轮询行为（应该选择不同的服务）
	selectedIDs := []string{selected1.ID, selected2.ID, selected3.ID}
	if selectedIDs[0] == selectedIDs[1] && selectedIDs[1] == selectedIDs[2] {
		t.Error("Expected round-robin to select different services")
	}

	// 测试随机负载均衡器
	random := NewRandomLoadBalancer()
	selected := random.Select(services)
	if selected == nil {
		t.Error("Expected to select a healthy service")
	}

	if selected.Health != "healthy" {
		t.Errorf("Expected healthy service, got %s", selected.Health)
	}

	// 测试空服务列表
	emptyServices := []*ServiceInfo{}
	selected = rr.Select(emptyServices)
	if selected != nil {
		t.Error("Expected nil when no services available")
	}

	// 测试只有不健康服务
	unhealthyServices := []*ServiceInfo{
		{ID: "1", Name: "service", Health: "unhealthy"},
		{ID: "2", Name: "service", Health: "unhealthy"},
	}

	selected = rr.Select(unhealthyServices)
	if selected != nil {
		t.Error("Expected nil when no healthy services available")
	}
}

func TestHealthChecker(t *testing.T) {
	checker := NewHTTPHealthChecker(5 * time.Second)
	ctx := context.Background()

	service := &ServiceInfo{
		ID:       "test-service",
		Name:     "test-service",
		Address:  "localhost",
		Port:     8080,
		Protocol: "http",
		Health:   "unknown",
	}

	// 测试健康检查
	err := checker.Check(ctx, service)
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	if service.Health != "healthy" {
		t.Errorf("Expected health to be 'healthy', got %s", service.Health)
	}

	// 测试健康状态判断
	isHealthy := checker.IsHealthy(service)
	if !isHealthy {
		t.Error("Expected service to be healthy")
	}

	service.Health = "unhealthy"
	isHealthy = checker.IsHealthy(service)
	if isHealthy {
		t.Error("Expected service to be unhealthy")
	}
}

func TestCircuitBreaker(t *testing.T) {
	cb := NewSimpleCircuitBreaker(3, 5*time.Second)
	ctx := context.Background()

	// 测试正常操作
	callCount := 0
	operation := func() error {
		callCount++
		return nil
	}

	err := cb.Execute(ctx, operation)
	if err != nil {
		t.Errorf("Expected no error for successful operation: %v", err)
	}

	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// 测试失败操作
	failureCount := 0
	failingOperation := func() error {
		failureCount++
		return fmt.Errorf("operation failed")
	}

	// 执行失败操作直到熔断器开启
	for i := 0; i < 4; i++ {
		err = cb.Execute(ctx, failingOperation)
		if i < 3 {
			// 前3次应该返回错误但熔断器未开启
			if err == nil {
				t.Error("Expected error for failing operation")
			}
		} else {
			// 第4次应该因为熔断器开启而失败
			if err == nil {
				t.Error("Expected circuit breaker to be open")
			}
		}
	}

	if failureCount != 3 {
		t.Errorf("Expected 3 failures, got %d", failureCount)
	}

	// 测试熔断器状态
	if !cb.IsOpen() {
		t.Error("Expected circuit breaker to be open")
	}

	// 测试重置熔断器
	cb.Reset()
	if cb.IsOpen() {
		t.Error("Expected circuit breaker to be closed after reset")
	}
}

func TestServiceRegistryCleanup(t *testing.T) {
	registry := NewMemoryServiceRegistry()
	ctx := context.Background()

	// 注册一个短期服务
	service := &ServiceInfo{
		ID:       "short-lived-service",
		Name:     "test-service",
		Address:  "localhost",
		Port:     8080,
		Protocol: "http",
		Health:   "healthy",
		TTL:      1 * time.Second, // 1秒后过期
	}

	err := registry.Register(ctx, service)
	if err != nil {
		t.Errorf("Failed to register service: %v", err)
	}

	// 验证服务存在
	_, err = registry.GetService(ctx, "short-lived-service")
	if err != nil {
		t.Errorf("Service should exist: %v", err)
	}

	// 等待服务过期
	time.Sleep(2 * time.Second)

	// 手动清理过期服务
	registry.CleanupExpiredServices()

	// 验证服务已被清理
	_, err = registry.GetService(ctx, "short-lived-service")
	if err == nil {
		t.Error("Service should have been cleaned up")
	}
}

func TestServiceDiscoveryWithUnhealthyServices(t *testing.T) {
	registry := NewMemoryServiceRegistry()
	discovery := NewMemoryServiceDiscovery(registry, NewRoundRobinLoadBalancer())
	ctx := context.Background()

	// 注册健康和不健康的服务
	services := []*ServiceInfo{
		{ID: "healthy-1", Name: "test-service", Health: "healthy"},
		{ID: "healthy-2", Name: "test-service", Health: "healthy"},
		{ID: "unhealthy-1", Name: "test-service", Health: "unhealthy"},
		{ID: "unhealthy-2", Name: "test-service", Health: "unhealthy"},
	}

	for _, service := range services {
		err := registry.Register(ctx, service)
		if err != nil {
			t.Errorf("Failed to register service %s: %v", service.ID, err)
		}
	}

	// 发现所有服务
	allServices, err := discovery.Discover(ctx, "test-service")
	if err != nil {
		t.Errorf("Failed to discover services: %v", err)
	}

	if len(allServices) != 4 {
		t.Errorf("Expected 4 services, got %d", len(allServices))
	}

	// 发现单个服务（应该只选择健康的）
	selected, err := discovery.DiscoverOne(ctx, "test-service")
	if err != nil {
		t.Errorf("Failed to discover one service: %v", err)
	}

	if selected.Health != "healthy" {
		t.Errorf("Expected healthy service, got %s", selected.Health)
	}
}

func TestServiceRegistryConcurrency(t *testing.T) {
	registry := NewMemoryServiceRegistry()
	ctx := context.Background()

	// 并发注册服务
	const numGoroutines = 10
	const servicesPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < servicesPerGoroutine; j++ {
				serviceID := fmt.Sprintf("service-%d-%d", goroutineID, j)
				service := &ServiceInfo{
					ID:       serviceID,
					Name:     "concurrent-service",
					Address:  "localhost",
					Port:     8080 + goroutineID*100 + j,
					Protocol: "http",
					Health:   "healthy",
				}

				err := registry.Register(ctx, service)
				if err != nil {
					t.Errorf("Failed to register service %s: %v", serviceID, err)
				}
			}
		}(i)
	}

	wg.Wait()

	// 验证所有服务都已注册
	services, err := registry.ListServices(ctx)
	if err != nil {
		t.Errorf("Failed to list services: %v", err)
	}

	expectedCount := numGoroutines * servicesPerGoroutine
	if len(services) != expectedCount {
		t.Errorf("Expected %d services, got %d", expectedCount, len(services))
	}
}
