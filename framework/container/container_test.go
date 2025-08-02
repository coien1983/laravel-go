package container

import (
	"testing"
)

// TestService 测试服务
type TestService struct {
	Name string
}

// TestServiceWithDependency 带依赖的测试服务
type TestServiceWithDependency struct {
	Service *TestService `inject:"container"`
}

func TestNewContainer(t *testing.T) {
	container := NewContainer()
	if container == nil {
		t.Fatal("NewContainer() should not return nil")
	}
}

func TestBindAndMake(t *testing.T) {
	container := NewContainer()

	// 注册服务
	service := &TestService{Name: "test"}
	container.Bind((*TestService)(nil), service)

	// 解析服务
	result := container.Make((*TestService)(nil))
	if result == nil {
		t.Fatal("Make() should not return nil")
	}

	if resultService, ok := result.(*TestService); !ok {
		t.Fatal("Make() should return *TestService")
	} else if resultService.Name != "test" {
		t.Fatalf("Expected name 'test', got '%s'", resultService.Name)
	}
}

func TestBindSingleton(t *testing.T) {
	container := NewContainer()

	// 注册单例服务
	service := &TestService{Name: "singleton"}
	container.BindSingleton((*TestService)(nil), service)

	// 解析服务两次
	result1 := container.Make((*TestService)(nil))
	result2 := container.Make((*TestService)(nil))

	// 应该返回同一个实例
	if result1 != result2 {
		t.Fatal("Singleton should return the same instance")
	}
}

func TestHas(t *testing.T) {
	container := NewContainer()

	// 检查不存在的服务
	if container.Has((*TestService)(nil)) {
		t.Fatal("Has() should return false for unregistered service")
	}

	// 注册服务
	service := &TestService{Name: "test"}
	container.Bind((*TestService)(nil), service)

	// 检查存在的服务
	if !container.Has((*TestService)(nil)) {
		t.Fatal("Has() should return true for registered service")
	}
}

func TestCall(t *testing.T) {
	container := NewContainer()

	// 注册服务
	service := &TestService{Name: "test"}
	container.Bind((*TestService)(nil), service)

	// 定义回调函数
	callback := func(service *TestService) string {
		return service.Name
	}

	// 调用函数
	results, err := container.Call(callback)
	if err != nil {
		t.Fatalf("Call() should not return error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if name, ok := results[0].(string); !ok {
		t.Fatal("Result should be string")
	} else if name != "test" {
		t.Fatalf("Expected 'test', got '%s'", name)
	}
}

func TestBindCallback(t *testing.T) {
	container := NewContainer()

	// 注册回调函数
	container.BindCallback((*TestService)(nil), func(c Container) interface{} {
		return &TestService{Name: "callback"}
	})

	// 解析服务
	result := container.Make((*TestService)(nil))
	if result == nil {
		t.Fatal("Make() should not return nil")
	}

	if resultService, ok := result.(*TestService); !ok {
		t.Fatal("Make() should return *TestService")
	} else if resultService.Name != "callback" {
		t.Fatalf("Expected name 'callback', got '%s'", resultService.Name)
	}
}
