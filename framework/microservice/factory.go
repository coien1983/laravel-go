package microservice

import (
	"fmt"
	"time"
)

// RegistryType 注册中心类型
type RegistryType string

const (
	RegistryTypeMemory     RegistryType = "memory"
	RegistryTypeEtcd       RegistryType = "etcd"
	RegistryTypeConsul     RegistryType = "consul"
	RegistryTypeNacos      RegistryType = "nacos"
	RegistryTypeZookeeper  RegistryType = "zookeeper"
)

// RegistryConfig 注册中心配置
type RegistryConfig struct {
	Type     RegistryType              `json:"type"`
	Memory   *MemoryRegistryConfig     `json:"memory,omitempty"`
	Etcd     *EtcdConfig               `json:"etcd,omitempty"`
	Consul   *ConsulConfig             `json:"consul,omitempty"`
	Nacos    *NacosConfig              `json:"nacos,omitempty"`
	Zookeeper *ZookeeperConfig         `json:"zookeeper,omitempty"`
}

// MemoryRegistryConfig 内存注册中心配置
type MemoryRegistryConfig struct {
	CleanupInterval time.Duration `json:"cleanup_interval"`
}

// NewServiceRegistry 创建服务注册中心
func NewServiceRegistry(config *RegistryConfig) (ServiceRegistry, error) {
	if config == nil {
		// 默认使用内存注册中心
		return NewMemoryServiceRegistry(), nil
	}

	switch config.Type {
	case RegistryTypeMemory:
		return NewMemoryServiceRegistry(), nil

	case RegistryTypeEtcd:
		if config.Etcd == nil {
			config.Etcd = &EtcdConfig{
				Endpoints: []string{"localhost:2379"},
				Prefix:    "/laravel-go/services",
				TTL:       30 * time.Second,
			}
		}
		return NewEtcdServiceRegistry(config.Etcd)

	case RegistryTypeConsul:
		if config.Consul == nil {
			config.Consul = &ConsulConfig{
				Address: "localhost:8500",
				Prefix:  "laravel-go/services",
				TTL:     30 * time.Second,
			}
		}
		return NewConsulServiceRegistry(config.Consul)

	case RegistryTypeNacos:
		if config.Nacos == nil {
			config.Nacos = &NacosConfig{
				ServerAddr: "localhost:8848",
				Namespace:  "public",
				Group:      "DEFAULT_GROUP",
				TTL:        30 * time.Second,
			}
		}
		return NewNacosServiceRegistry(config.Nacos)

	case RegistryTypeZookeeper:
		if config.Zookeeper == nil {
			config.Zookeeper = &ZookeeperConfig{
				Servers:        []string{"localhost:2181"},
				Prefix:         "/laravel-go/services",
				TTL:            30 * time.Second,
				SessionTimeout: 10 * time.Second,
			}
		}
		return NewZookeeperServiceRegistry(config.Zookeeper)

	default:
		return nil, fmt.Errorf("unsupported registry type: %s", config.Type)
	}
}

// NewServiceDiscovery 创建服务发现
func NewServiceDiscovery(registry ServiceRegistry, loadBalancer LoadBalancer) ServiceDiscovery {
	if loadBalancer == nil {
		loadBalancer = NewRoundRobinLoadBalancer()
	}

	return &MemoryServiceDiscovery{
		registry:     registry,
		loadBalancer: loadBalancer,
		cache:        make(map[string][]*ServiceInfo),
		watchers:     make(map[string]chan ServiceEvent),
	}
}



// RegistryBuilder 注册中心构建器
type RegistryBuilder struct {
	config *RegistryConfig
}

// NewRegistryBuilder 创建注册中心构建器
func NewRegistryBuilder() *RegistryBuilder {
	return &RegistryBuilder{
		config: &RegistryConfig{},
	}
}

// WithType 设置注册中心类型
func (b *RegistryBuilder) WithType(registryType RegistryType) *RegistryBuilder {
	b.config.Type = registryType
	return b
}

// WithEtcd 配置 etcd
func (b *RegistryBuilder) WithEtcd(endpoints []string, prefix string) *RegistryBuilder {
	b.config.Etcd = &EtcdConfig{
		Endpoints: endpoints,
		Prefix:    prefix,
		TTL:       30 * time.Second,
	}
	return b
}

// WithConsul 配置 Consul
func (b *RegistryBuilder) WithConsul(address, prefix string) *RegistryBuilder {
	b.config.Consul = &ConsulConfig{
		Address: address,
		Prefix:  prefix,
		TTL:     30 * time.Second,
	}
	return b
}

// WithNacos 配置 Nacos
func (b *RegistryBuilder) WithNacos(serverAddr, namespace, group string) *RegistryBuilder {
	b.config.Nacos = &NacosConfig{
		ServerAddr: serverAddr,
		Namespace:  namespace,
		Group:      group,
		TTL:        30 * time.Second,
	}
	return b
}

// WithZookeeper 配置 Zookeeper
func (b *RegistryBuilder) WithZookeeper(servers []string, prefix string) *RegistryBuilder {
	b.config.Zookeeper = &ZookeeperConfig{
		Servers:        servers,
		Prefix:         prefix,
		TTL:            30 * time.Second,
		SessionTimeout: 10 * time.Second,
	}
	return b
}

// Build 构建注册中心
func (b *RegistryBuilder) Build() (ServiceRegistry, error) {
	return NewServiceRegistry(b.config)
} 