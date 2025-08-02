package performance

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemMonitor 系统监控器
type SystemMonitor struct {
	monitor     Monitor
	running     bool
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
	collectors  []Collector
}

// Collector 指标收集器接口
type Collector interface {
	// Collect 收集指标
	Collect(monitor Monitor) error
	// Name 收集器名称
	Name() string
}

// CPUMonitor CPU监控器
type CPUMonitor struct {
	interval time.Duration
}

// NewCPUMonitor 创建CPU监控器
func NewCPUMonitor(interval time.Duration) *CPUMonitor {
	return &CPUMonitor{
		interval: interval,
	}
}

func (cm *CPUMonitor) Name() string {
	return "cpu"
}

func (cm *CPUMonitor) Collect(monitor Monitor) error {
	// 获取CPU使用率
	percentages, err := cpu.Percent(cm.interval, false)
	if err != nil {
		return fmt.Errorf("failed to get CPU usage: %w", err)
	}

	var cpuUsage float64
	if len(percentages) > 0 {
		cpuUsage = percentages[0]
	}

	// 获取或创建CPU使用率指标
	cpuMetric := monitor.GetMetric("cpu_usage")
	if cpuMetric == nil {
		cpuMetric = NewGauge("cpu_usage", map[string]string{"type": "percentage"})
		monitor.RegisterMetric(cpuMetric)
	}

	if gauge, ok := cpuMetric.(*Gauge); ok {
		gauge.Set(cpuUsage)
	}

	// 获取CPU核心数
	numCPU := runtime.NumCPU()
	cpuCoresMetric := monitor.GetMetric("cpu_cores")
	if cpuCoresMetric == nil {
		cpuCoresMetric = NewGauge("cpu_cores", map[string]string{"type": "count"})
		monitor.RegisterMetric(cpuCoresMetric)
	}

	if gauge, ok := cpuCoresMetric.(*Gauge); ok {
		gauge.Set(float64(numCPU))
	}

	return nil
}

// MemoryMonitor 内存监控器
type MemoryMonitor struct {
	interval time.Duration
}

// NewMemoryMonitor 创建内存监控器
func NewMemoryMonitor(interval time.Duration) *MemoryMonitor {
	return &MemoryMonitor{
		interval: interval,
	}
}

func (mm *MemoryMonitor) Name() string {
	return "memory"
}

func (mm *MemoryMonitor) Collect(monitor Monitor) error {
	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get memory info: %w", err)
	}

	// 总内存
	totalMemoryMetric := monitor.GetMetric("memory_total")
	if totalMemoryMetric == nil {
		totalMemoryMetric = NewGauge("memory_total", map[string]string{"unit": "bytes"})
		monitor.RegisterMetric(totalMemoryMetric)
	}
	if gauge, ok := totalMemoryMetric.(*Gauge); ok {
		gauge.Set(float64(memInfo.Total))
	}

	// 已用内存
	usedMemoryMetric := monitor.GetMetric("memory_used")
	if usedMemoryMetric == nil {
		usedMemoryMetric = NewGauge("memory_used", map[string]string{"unit": "bytes"})
		monitor.RegisterMetric(usedMemoryMetric)
	}
	if gauge, ok := usedMemoryMetric.(*Gauge); ok {
		gauge.Set(float64(memInfo.Used))
	}

	// 可用内存
	availableMemoryMetric := monitor.GetMetric("memory_available")
	if availableMemoryMetric == nil {
		availableMemoryMetric = NewGauge("memory_available", map[string]string{"unit": "bytes"})
		monitor.RegisterMetric(availableMemoryMetric)
	}
	if gauge, ok := availableMemoryMetric.(*Gauge); ok {
		gauge.Set(float64(memInfo.Available))
	}

	// 内存使用率
	memoryUsageMetric := monitor.GetMetric("memory_usage_percent")
	if memoryUsageMetric == nil {
		memoryUsageMetric = NewGauge("memory_usage_percent", map[string]string{"unit": "percentage"})
		monitor.RegisterMetric(memoryUsageMetric)
	}
	if gauge, ok := memoryUsageMetric.(*Gauge); ok {
		gauge.Set(memInfo.UsedPercent)
	}

	// Go运行时内存统计
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Go堆内存
	heapAllocMetric := monitor.GetMetric("go_heap_alloc")
	if heapAllocMetric == nil {
		heapAllocMetric = NewGauge("go_heap_alloc", map[string]string{"unit": "bytes"})
		monitor.RegisterMetric(heapAllocMetric)
	}
	if gauge, ok := heapAllocMetric.(*Gauge); ok {
		gauge.Set(float64(m.HeapAlloc))
	}

	// Go堆内存系统
	heapSysMetric := monitor.GetMetric("go_heap_sys")
	if heapSysMetric == nil {
		heapSysMetric = NewGauge("go_heap_sys", map[string]string{"unit": "bytes"})
		monitor.RegisterMetric(heapSysMetric)
	}
	if gauge, ok := heapSysMetric.(*Gauge); ok {
		gauge.Set(float64(m.HeapSys))
	}

	// Go协程数量
	goroutinesMetric := monitor.GetMetric("go_goroutines")
	if goroutinesMetric == nil {
		goroutinesMetric = NewGauge("go_goroutines", map[string]string{"type": "count"})
		monitor.RegisterMetric(goroutinesMetric)
	}
	if gauge, ok := goroutinesMetric.(*Gauge); ok {
		gauge.Set(float64(runtime.NumGoroutine()))
	}

	return nil
}

// DiskMonitor 磁盘监控器
type DiskMonitor struct {
	interval time.Duration
	paths    []string
}

// NewDiskMonitor 创建磁盘监控器
func NewDiskMonitor(interval time.Duration, paths []string) *DiskMonitor {
	if len(paths) == 0 {
		paths = []string{"/"}
	}
	return &DiskMonitor{
		interval: interval,
		paths:    paths,
	}
}

func (dm *DiskMonitor) Name() string {
	return "disk"
}

func (dm *DiskMonitor) Collect(monitor Monitor) error {
	for _, path := range dm.paths {
		// 获取磁盘使用情况
		usage, err := disk.Usage(path)
		if err != nil {
			continue // 跳过无法访问的路径
		}

		// 总空间
		totalMetric := monitor.GetMetric(fmt.Sprintf("disk_total_%s", path))
		if totalMetric == nil {
			totalMetric = NewGauge(fmt.Sprintf("disk_total_%s", path), map[string]string{
				"path": path,
				"unit": "bytes",
			})
			monitor.RegisterMetric(totalMetric)
		}
		if gauge, ok := totalMetric.(*Gauge); ok {
			gauge.Set(float64(usage.Total))
		}

		// 已用空间
		usedMetric := monitor.GetMetric(fmt.Sprintf("disk_used_%s", path))
		if usedMetric == nil {
			usedMetric = NewGauge(fmt.Sprintf("disk_used_%s", path), map[string]string{
				"path": path,
				"unit": "bytes",
			})
			monitor.RegisterMetric(usedMetric)
		}
		if gauge, ok := usedMetric.(*Gauge); ok {
			gauge.Set(float64(usage.Used))
		}

		// 可用空间
		freeMetric := monitor.GetMetric(fmt.Sprintf("disk_free_%s", path))
		if freeMetric == nil {
			freeMetric = NewGauge(fmt.Sprintf("disk_free_%s", path), map[string]string{
				"path": path,
				"unit": "bytes",
			})
			monitor.RegisterMetric(freeMetric)
		}
		if gauge, ok := freeMetric.(*Gauge); ok {
			gauge.Set(float64(usage.Free))
		}

		// 使用率
		usagePercentMetric := monitor.GetMetric(fmt.Sprintf("disk_usage_percent_%s", path))
		if usagePercentMetric == nil {
			usagePercentMetric = NewGauge(fmt.Sprintf("disk_usage_percent_%s", path), map[string]string{
				"path": path,
				"unit": "percentage",
			})
			monitor.RegisterMetric(usagePercentMetric)
		}
		if gauge, ok := usagePercentMetric.(*Gauge); ok {
			gauge.Set(usage.UsedPercent)
		}
	}

	return nil
}

// NetworkMonitor 网络监控器
type NetworkMonitor struct {
	interval time.Duration
}

// NewNetworkMonitor 创建网络监控器
func NewNetworkMonitor(interval time.Duration) *NetworkMonitor {
	return &NetworkMonitor{
		interval: interval,
	}
}

func (nm *NetworkMonitor) Name() string {
	return "network"
}

func (nm *NetworkMonitor) Collect(monitor Monitor) error {
	// 这里可以添加网络监控逻辑
	// 由于gopsutil的网络监控比较复杂，这里先预留接口
	return nil
}

// NewSystemMonitor 创建系统监控器
func NewSystemMonitor(monitor Monitor) *SystemMonitor {
	sm := &SystemMonitor{
		monitor: monitor,
	}

	// 添加默认收集器
	sm.collectors = []Collector{
		NewCPUMonitor(5 * time.Second),
		NewMemoryMonitor(5 * time.Second),
		NewDiskMonitor(10 * time.Second, []string{"/"}),
		NewNetworkMonitor(10 * time.Second),
	}

	return sm
}

// AddCollector 添加收集器
func (sm *SystemMonitor) AddCollector(collector Collector) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.collectors = append(sm.collectors, collector)
}

// Start 启动系统监控
func (sm *SystemMonitor) Start(ctx context.Context) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.running {
		return nil
	}

	sm.ctx, sm.cancel = context.WithCancel(ctx)
	sm.running = true

	// 启动监控循环
	go sm.monitorLoop()

	return nil
}

// Stop 停止系统监控
func (sm *SystemMonitor) Stop() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.running {
		return nil
	}

	if sm.cancel != nil {
		sm.cancel()
	}
	sm.running = false

	return nil
}

// monitorLoop 监控循环
func (sm *SystemMonitor) monitorLoop() {
	ticker := time.NewTicker(5 * time.Second) // 每5秒收集一次
	defer ticker.Stop()

	for {
		select {
		case <-sm.ctx.Done():
			return
		case <-ticker.C:
			sm.collectAll()
		}
	}
}

// collectAll 收集所有指标
func (sm *SystemMonitor) collectAll() {
	sm.mu.RLock()
	collectors := make([]Collector, len(sm.collectors))
	copy(collectors, sm.collectors)
	sm.mu.RUnlock()

	for _, collector := range collectors {
		if err := collector.Collect(sm.monitor); err != nil {
			// 记录错误但不中断其他收集器
			continue
		}
	}
}

// GetSystemMetrics 获取系统指标
func (sm *SystemMonitor) GetSystemMetrics() map[string]interface{} {
	metrics := sm.monitor.GetAllMetrics()
	result := make(map[string]interface{})

	for name, metric := range metrics {
		result[name] = map[string]interface{}{
			"type":      metric.Type(),
			"value":     metric.Value(),
			"labels":    metric.Labels(),
			"timestamp": metric.Timestamp(),
		}
	}

	return result
} 