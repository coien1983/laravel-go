package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container 接口定义
type Container interface {
	// 注册服务
	Bind(abstract interface{}, concrete interface{})
	BindSingleton(abstract interface{}, concrete interface{})
	BindCallback(abstract interface{}, callback func(Container) interface{})

	// 解析服务
	Make(abstract interface{}) interface{}
	Resolve(abstract interface{}) error

	// 检查服务是否存在
	Has(abstract interface{}) bool

	// 调用方法并注入依赖
	Call(callback interface{}, parameters ...interface{}) ([]interface{}, error)

	// 获取容器实例
	Instance(abstract interface{}, instance interface{})
}

// ServiceProvider 服务提供者接口
type ServiceProvider interface {
	Register(container Container)
	Boot(container Container)
}

// container 容器实现
type container struct {
	bindings   map[reflect.Type]*binding
	singletons map[reflect.Type]interface{}
	instances  map[reflect.Type]interface{}
	resolving  map[reflect.Type]bool
	mutex      sync.RWMutex
}

// binding 绑定信息
type binding struct {
	concrete interface{}
	shared   bool
	callback func(Container) interface{}
}

// NewContainer 创建新的容器实例
func NewContainer() Container {
	return &container{
		bindings:   make(map[reflect.Type]*binding),
		singletons: make(map[reflect.Type]interface{}),
		instances:  make(map[reflect.Type]interface{}),
		resolving:  make(map[reflect.Type]bool),
	}
}

// Bind 注册服务
func (c *container) Bind(abstract interface{}, concrete interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	abstractType := reflect.TypeOf(abstract)
	if abstractType.Kind() == reflect.Ptr && abstractType.Elem().Kind() == reflect.Interface {
		// 如果 abstract 是指向接口的指针，使用接口类型
		abstractType = abstractType.Elem()
	}
	c.bindings[abstractType] = &binding{
		concrete: concrete,
		shared:   false,
	}
}

// BindSingleton 注册单例服务
func (c *container) BindSingleton(abstract interface{}, concrete interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	abstractType := reflect.TypeOf(abstract)
	if abstractType.Kind() == reflect.Ptr && abstractType.Elem().Kind() == reflect.Interface {
		// 如果 abstract 是指向接口的指针，使用接口类型
		abstractType = abstractType.Elem()
	}
	c.bindings[abstractType] = &binding{
		concrete: concrete,
		shared:   true,
	}
}

// BindCallback 注册回调函数
func (c *container) BindCallback(abstract interface{}, callback func(Container) interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	abstractType := reflect.TypeOf(abstract)
	if abstractType.Kind() == reflect.Ptr && abstractType.Elem().Kind() == reflect.Interface {
		// 如果 abstract 是指向接口的指针，使用接口类型
		abstractType = abstractType.Elem()
	}
	c.bindings[abstractType] = &binding{
		callback: callback,
		shared:   false,
	}
}

// Make 解析服务
func (c *container) Make(abstract interface{}) interface{} {
	abstractType := reflect.TypeOf(abstract)
	if abstractType.Kind() == reflect.Ptr && abstractType.Elem().Kind() == reflect.Interface {
		// 如果 abstract 是指向接口的指针，使用接口类型
		abstractType = abstractType.Elem()
	}

	// 检查是否已经解析过
	if instance, exists := c.instances[abstractType]; exists {
		return instance
	}

	// 检查单例
	if singleton, exists := c.singletons[abstractType]; exists {
		return singleton
	}

	// 解析服务
	instance, err := c.resolve(abstractType)
	if err != nil {
		panic(fmt.Sprintf("Unable to resolve %v: %v", abstractType, err))
	}

	// 如果是单例，保存实例
	c.mutex.Lock()
	if binding, exists := c.bindings[abstractType]; exists && binding.shared {
		c.singletons[abstractType] = instance
	}
	c.mutex.Unlock()

	return instance
}

// Resolve 解析服务并返回错误
func (c *container) Resolve(abstract interface{}) error {
	abstractType := reflect.TypeOf(abstract)
	_, err := c.resolve(abstractType)
	return err
}

// Has 检查服务是否存在
func (c *container) Has(abstract interface{}) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	abstractType := reflect.TypeOf(abstract)
	if abstractType.Kind() == reflect.Ptr && abstractType.Elem().Kind() == reflect.Interface {
		// 如果 abstract 是指向接口的指针，使用接口类型
		abstractType = abstractType.Elem()
	}
	_, exists := c.bindings[abstractType]
	return exists
}

// Call 调用方法并注入依赖
func (c *container) Call(callback interface{}, parameters ...interface{}) ([]interface{}, error) {
	callbackType := reflect.TypeOf(callback)
	if callbackType.Kind() != reflect.Func {
		return nil, fmt.Errorf("callback must be a function")
	}

	// 准备参数
	args := make([]reflect.Value, callbackType.NumIn())

	for i := 0; i < callbackType.NumIn(); i++ {
		paramType := callbackType.In(i)

		// 检查是否有提供的参数
		if i < len(parameters) && parameters[i] != nil {
			args[i] = reflect.ValueOf(parameters[i])
		} else {
			// 从容器解析依赖
			instance, err := c.resolve(paramType)
			if err != nil {
				return nil, fmt.Errorf("unable to resolve dependency %v: %v", paramType, err)
			}
			args[i] = reflect.ValueOf(instance)
		}
	}

	// 调用函数
	results := reflect.ValueOf(callback).Call(args)

	// 转换返回值
	outputs := make([]interface{}, len(results))
	for i, result := range results {
		outputs[i] = result.Interface()
	}

	return outputs, nil
}

// Instance 获取容器实例
func (c *container) Instance(abstract interface{}, instance interface{}) {
	abstractType := reflect.TypeOf(abstract)
	if abstractType.Kind() == reflect.Ptr && abstractType.Elem().Kind() == reflect.Interface {
		// 如果 abstract 是指向接口的指针，使用接口类型
		abstractType = abstractType.Elem()
	}

	c.mutex.RLock()
	if existing, exists := c.instances[abstractType]; exists {
		c.mutex.RUnlock()
		reflect.ValueOf(instance).Elem().Set(reflect.ValueOf(existing))
		return
	}

	if singleton, exists := c.singletons[abstractType]; exists {
		c.mutex.RUnlock()
		reflect.ValueOf(instance).Elem().Set(reflect.ValueOf(singleton))
		return
	}
	c.mutex.RUnlock()

	// 解析新实例
	resolved, err := c.resolve(abstractType)
	if err != nil {
		panic(fmt.Sprintf("Unable to resolve %v: %v", abstractType, err))
	}

	// 将解析的实例设置到传入的指针中
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}
	instanceValue.Set(reflect.ValueOf(resolved))
}

// resolve 解析服务实现
func (c *container) resolve(abstractType reflect.Type) (interface{}, error) {
	c.mutex.Lock()

	// 检查循环依赖
	if c.resolving[abstractType] {
		c.mutex.Unlock()
		return nil, fmt.Errorf("circular dependency detected for %v", abstractType)
	}

	c.resolving[abstractType] = true
	c.mutex.Unlock()

	defer func() {
		c.mutex.Lock()
		delete(c.resolving, abstractType)
		c.mutex.Unlock()
	}()

	// 查找绑定
	binding, exists := c.bindings[abstractType]
	if !exists {
		// 尝试直接实例化
		return c.build(abstractType)
	}

	// 如果有回调函数，执行回调
	if binding.callback != nil {
		return binding.callback(c), nil
	}

	// 如果 concrete 已经是实例，直接返回
	if binding.concrete != nil {
		concreteType := reflect.TypeOf(binding.concrete)
		if concreteType.Kind() != reflect.Func {
			return binding.concrete, nil
		}

		// 如果是函数，调用它
		results := reflect.ValueOf(binding.concrete).Call(nil)
		if len(results) > 0 {
			return results[0].Interface(), nil
		}
		return nil, nil
	}

	// 尝试直接实例化
	return c.build(abstractType)
}

// build 构建实例
func (c *container) build(concreteType reflect.Type) (interface{}, error) {
	// 检查是否是接口
	if concreteType.Kind() == reflect.Interface {
		return nil, fmt.Errorf("cannot instantiate interface %v", concreteType)
	}

	// 检查是否是指针
	if concreteType.Kind() == reflect.Ptr {
		concreteType = concreteType.Elem()
	}

	// 创建实例
	instance := reflect.New(concreteType).Interface()

	// 注入依赖
	if err := c.injectDependencies(instance); err != nil {
		return nil, err
	}

	return instance, nil
}

// injectDependencies 注入依赖
func (c *container) injectDependencies(instance interface{}) error {
	instanceValue := reflect.ValueOf(instance)
	instanceType := instanceValue.Type()

	// 如果是指针，获取元素
	if instanceType.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
		instanceType = instanceType.Elem()
	}

	// 遍历字段，查找需要注入的依赖
	for i := 0; i < instanceType.NumField(); i++ {
		field := instanceType.Field(i)
		fieldValue := instanceValue.Field(i)

		// 检查是否有注入标签
		if tag := field.Tag.Get("inject"); tag != "" {
			var dependency interface{}

			if tag == "container" {
				// 注入容器本身
				dependency = c
			} else {
				// 根据字段类型从容器解析依赖
				fieldType := field.Type
				if fieldType.Kind() == reflect.Interface {
					// 对于接口类型，尝试从容器解析
					var err error
					dependency, err = c.resolve(fieldType)
					if err != nil {
						// 如果解析失败，跳过这个字段
						continue
					}
				} else {
					// 对于其他类型，暂时跳过
					continue
				}
			}

			if fieldValue.CanSet() {
				fieldValue.Set(reflect.ValueOf(dependency))
			}
		}
	}

	return nil
}
