package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type PerformanceMetrics struct {
	Timestamp     time.Time     `json:"timestamp"`
	MemoryUsage   MemoryStats   `json:"memory_usage"`
	CPUUsage      CPUStats      `json:"cpu_usage"`
	HTTPMetrics   HTTPStats     `json:"http_metrics"`
	GoroutineInfo GoroutineInfo `json:"goroutine_info"`
}

type MemoryStats struct {
	Alloc      uint64  `json:"alloc"`
	TotalAlloc uint64  `json:"total_alloc"`
	Sys        uint64  `json:"sys"`
	NumGC      uint32  `json:"num_gc"`
	HeapAlloc  uint64  `json:"heap_alloc"`
	HeapSys    uint64  `json:"heap_sys"`
	HeapIdle   uint64  `json:"heap_idle"`
	HeapInuse  uint64  `json:"heap_inuse"`
}

type CPUStats struct {
	NumCPU       int     `json:"num_cpu"`
	NumGoroutine int     `json:"num_goroutine"`
	CPUPercent   float64 `json:"cpu_percent"`
}

type HTTPStats struct {
	TotalRequests    int64         `json:"total_requests"`
	ActiveRequests   int64         `json:"active_requests"`
	AverageResponse  time.Duration `json:"average_response"`
	SlowestRequest   time.Duration `json:"slowest_request"`
	FastestRequest   time.Duration `json:"fastest_request"`
	ErrorCount       int64         `json:"error_count"`
	SuccessCount     int64         `json:"success_count"`
}

type GoroutineInfo struct {
	Count int `json:"count"`
}

type PerformanceAnalyzer struct {
	metrics     []PerformanceMetrics
	outputFile  string
	interval    time.Duration
	httpMetrics *HTTPMetricsCollector
}

type HTTPMetricsCollector struct {
	requests    []time.Duration
	totalCount  int64
	errorCount  int64
	successCount int64
}

func NewPerformanceAnalyzer(outputFile string, interval time.Duration) *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		outputFile: outputFile,
		interval:   interval,
		httpMetrics: &HTTPMetricsCollector{
			requests: make([]time.Duration, 0),
		},
	}
}

func (pa *PerformanceAnalyzer) Start() {
	ticker := time.NewTicker(pa.interval)
	defer ticker.Stop()

	fmt.Printf("Performance analyzer started. Collecting metrics every %v\n", pa.interval)
	fmt.Printf("Output will be saved to: %s\n", pa.outputFile)

	for {
		select {
		case <-ticker.C:
			metrics := pa.collectMetrics()
			pa.metrics = append(pa.metrics, metrics)
			pa.printMetrics(metrics)
		}
	}
}

func (pa *PerformanceAnalyzer) collectMetrics() PerformanceMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return PerformanceMetrics{
		Timestamp: time.Now(),
		MemoryUsage: MemoryStats{
			Alloc:      m.Alloc,
			TotalAlloc: m.TotalAlloc,
			Sys:        m.Sys,
			NumGC:      m.NumGC,
			HeapAlloc:  m.HeapAlloc,
			HeapSys:    m.HeapSys,
			HeapIdle:   m.HeapIdle,
			HeapInuse:  m.HeapInuse,
		},
		CPUUsage: CPUStats{
			NumCPU:       runtime.NumCPU(),
			NumGoroutine: runtime.NumGoroutine(),
		},
		HTTPMetrics: pa.httpMetrics.getStats(),
		GoroutineInfo: GoroutineInfo{
			Count: runtime.NumGoroutine(),
		},
	}
}

func (pa *PerformanceAnalyzer) printMetrics(metrics PerformanceMetrics) {
	fmt.Printf("\n=== Performance Metrics (%s) ===\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Memory Usage:\n")
	fmt.Printf("  Alloc: %s\n", formatBytes(metrics.MemoryUsage.Alloc))
	fmt.Printf("  Total Alloc: %s\n", formatBytes(metrics.MemoryUsage.TotalAlloc))
	fmt.Printf("  Sys: %s\n", formatBytes(metrics.MemoryUsage.Sys))
	fmt.Printf("  Heap Alloc: %s\n", formatBytes(metrics.MemoryUsage.HeapAlloc))
	fmt.Printf("  Heap Sys: %s\n", formatBytes(metrics.MemoryUsage.HeapSys))
	fmt.Printf("  Num GC: %d\n", metrics.MemoryUsage.NumGC)
	
	fmt.Printf("CPU Usage:\n")
	fmt.Printf("  Num CPU: %d\n", metrics.CPUUsage.NumCPU)
	fmt.Printf("  Num Goroutines: %d\n", metrics.CPUUsage.NumGoroutine)
	
	fmt.Printf("HTTP Metrics:\n")
	fmt.Printf("  Total Requests: %d\n", metrics.HTTPMetrics.TotalRequests)
	fmt.Printf("  Active Requests: %d\n", metrics.HTTPMetrics.ActiveRequests)
	fmt.Printf("  Average Response: %v\n", metrics.HTTPMetrics.AverageResponse)
	fmt.Printf("  Slowest Request: %v\n", metrics.HTTPMetrics.SlowestRequest)
	fmt.Printf("  Fastest Request: %v\n", metrics.HTTPMetrics.FastestRequest)
	fmt.Printf("  Success Count: %d\n", metrics.HTTPMetrics.SuccessCount)
	fmt.Printf("  Error Count: %d\n", metrics.HTTPMetrics.ErrorCount)
}

func (pa *PerformanceAnalyzer) SaveMetrics() error {
	file, err := os.Create(pa.outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(pa.metrics)
}

func (hmc *HTTPMetricsCollector) RecordRequest(duration time.Duration, success bool) {
	hmc.requests = append(hmc.requests, duration)
	hmc.totalCount++
	
	if success {
		hmc.successCount++
	} else {
		hmc.errorCount++
	}
}

func (hmc *HTTPMetricsCollector) getStats() HTTPStats {
	if len(hmc.requests) == 0 {
		return HTTPStats{}
	}

	var total time.Duration
	slowest := hmc.requests[0]
	fastest := hmc.requests[0]

	for _, req := range hmc.requests {
		total += req
		if req > slowest {
			slowest = req
		}
		if req < fastest {
			fastest = req
		}
	}

	average := total / time.Duration(len(hmc.requests))

	return HTTPStats{
		TotalRequests:   hmc.totalCount,
		ActiveRequests:  int64(len(hmc.requests)),
		AverageResponse: average,
		SlowestRequest:  slowest,
		FastestRequest:  fastest,
		SuccessCount:    hmc.successCount,
		ErrorCount:      hmc.errorCount,
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func startProfiling() {
	// Start CPU profiling
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	// Start memory profiling
	memFile, err := os.Create("memory.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		pprof.WriteHeapProfile(memFile)
		memFile.Close()
	}()
}

func main() {
	var (
		outputFile = flag.String("output", "performance_metrics.json", "Output file for metrics")
		interval   = flag.Duration("interval", 5*time.Second, "Metrics collection interval")
		profile    = flag.Bool("profile", false, "Enable CPU and memory profiling")
		port       = flag.Int("port", 8080, "Port for HTTP server")
	)
	flag.Parse()

	if *profile {
		startProfiling()
	}

	analyzer := NewPerformanceAnalyzer(*outputFile, *interval)

	// Start metrics collection in a goroutine
	go analyzer.Start()

	// Start HTTP server for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Simulate some work
		time.Sleep(10 * time.Millisecond)
		
		w.Write([]byte("Hello, Laravel-Go!"))
		
		duration := time.Since(start)
		analyzer.httpMetrics.RecordRequest(duration, true)
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Simulate error
		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred"))
		
		duration := time.Since(start)
		analyzer.httpMetrics.RecordRequest(duration, false)
	})

	fmt.Printf("HTTP server starting on port %d\n", *port)
	fmt.Printf("Test endpoints:\n")
	fmt.Printf("  GET / - Success response\n")
	fmt.Printf("  GET /error - Error response\n")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
} 