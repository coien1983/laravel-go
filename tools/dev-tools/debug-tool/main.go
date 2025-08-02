package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"time"
)

type DebugInfo struct {
	Timestamp     time.Time              `json:"timestamp"`
	RuntimeInfo   RuntimeInfo            `json:"runtime_info"`
	MemoryInfo    MemoryInfo             `json:"memory_info"`
	GoroutineInfo []GoroutineInfo        `json:"goroutine_info"`
	StackTraces   map[string]interface{} `json:"stack_traces"`
	Environment   map[string]string      `json:"environment"`
}

type RuntimeInfo struct {
	GoVersion    string `json:"go_version"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
	NumCgoCall   int64  `json:"num_cgo_call"`
}

type MemoryInfo struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"num_gc"`
	HeapAlloc  uint64 `json:"heap_alloc"`
	HeapSys    uint64 `json:"heap_sys"`
	HeapIdle   uint64 `json:"heap_idle"`
	HeapInuse  uint64 `json:"heap_inuse"`
	HeapReleased uint64 `json:"heap_released"`
	HeapObjects uint64 `json:"heap_objects"`
}

type GoroutineInfo struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type DebugTool struct {
	port     int
	profiles map[string]*os.File
}

func NewDebugTool(port int) *DebugTool {
	return &DebugTool{
		port:     port,
		profiles: make(map[string]*os.File),
	}
}

func (dt *DebugTool) Start() {
	http.HandleFunc("/debug", dt.handleDebug)
	http.HandleFunc("/debug/memory", dt.handleMemoryProfile)
	http.HandleFunc("/debug/cpu", dt.handleCPUProfile)
	http.HandleFunc("/debug/goroutines", dt.handleGoroutines)
	http.HandleFunc("/debug/stack", dt.handleStack)
	http.HandleFunc("/debug/gc", dt.handleGC)
	http.HandleFunc("/debug/pprof/", dt.handlePProf)

	fmt.Printf("Debug tool started on port %d\n", dt.port)
	fmt.Printf("Available endpoints:\n")
	fmt.Printf("  GET /debug - General debug information\n")
	fmt.Printf("  GET /debug/memory - Memory profile\n")
	fmt.Printf("  GET /debug/cpu - CPU profile\n")
	fmt.Printf("  GET /debug/goroutines - Goroutine information\n")
	fmt.Printf("  GET /debug/stack - Stack traces\n")
	fmt.Printf("  GET /debug/gc - Force garbage collection\n")
	fmt.Printf("  GET /debug/pprof/* - Go pprof endpoints\n")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", dt.port), nil))
}

func (dt *DebugTool) handleDebug(w http.ResponseWriter, r *http.Request) {
	info := dt.collectDebugInfo()
	
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(info)
}

func (dt *DebugTool) handleMemoryProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=memory.prof")
	
	pprof.WriteHeapProfile(w)
}

func (dt *DebugTool) handleCPUProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=cpu.prof")
	
	pprof.StartCPUProfile(w)
	defer pprof.StopCPUProfile()
	
	// Collect CPU profile for 30 seconds
	time.Sleep(30 * time.Second)
}

func (dt *DebugTool) handleGoroutines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	goroutines := dt.collectGoroutineInfo()
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(goroutines)
}

func (dt *DebugTool) handleStack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	
	debug.WriteHeapDump(1)
	fmt.Fprintf(w, "Stack dump written to heap dump file\n")
}

func (dt *DebugTool) handleGC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	before := dt.getMemoryStats()
	runtime.GC()
	after := dt.getMemoryStats()
	
	result := map[string]interface{}{
		"before": before,
		"after":  after,
		"freed":  before.Alloc - after.Alloc,
	}
	
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(result)
}

func (dt *DebugTool) handlePProf(w http.ResponseWriter, r *http.Request) {
	// Redirect to Go's built-in pprof handler
	http.Redirect(w, r, "/debug/pprof/"+r.URL.Path[len("/debug/pprof/"):], http.StatusFound)
}

func (dt *DebugTool) collectDebugInfo() DebugInfo {
	return DebugInfo{
		Timestamp:   time.Now(),
		RuntimeInfo: dt.getRuntimeInfo(),
		MemoryInfo:  dt.getMemoryStats(),
		GoroutineInfo: dt.collectGoroutineInfo(),
		StackTraces: dt.getStackTraces(),
		Environment: dt.getEnvironment(),
	}
}

func (dt *DebugTool) getRuntimeInfo() RuntimeInfo {
	return RuntimeInfo{
		GoVersion:    runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		NumCgoCall:   runtime.NumCgoCall(),
	}
}

func (dt *DebugTool) getMemoryStats() MemoryInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return MemoryInfo{
		Alloc:        m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		Sys:          m.Sys,
		NumGC:        m.NumGC,
		HeapAlloc:    m.HeapAlloc,
		HeapSys:      m.HeapSys,
		HeapIdle:     m.HeapIdle,
		HeapInuse:    m.HeapInuse,
		HeapReleased: m.HeapReleased,
		HeapObjects:  m.HeapObjects,
	}
}

func (dt *DebugTool) collectGoroutineInfo() []GoroutineInfo {
	// This is a simplified version. In a real implementation,
	// you would use runtime.Stack() to get detailed goroutine information
	var info []GoroutineInfo
	
	// For demonstration, we'll just return basic info
	info = append(info, GoroutineInfo{
		ID:       1,
		Status:   "running",
		Function: "main.main",
		File:     "main.go",
		Line:     1,
	})
	
	return info
}

func (dt *DebugTool) getStackTraces() map[string]interface{} {
	// Collect stack traces for all goroutines
	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	
	return map[string]interface{}{
		"stack_trace": string(buf[:n]),
		"size":        n,
	}
}

func (dt *DebugTool) getEnvironment() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			env[pair[0]] = pair[1]
		}
	}
	return env
}

func (dt *DebugTool) StartProfiling() {
	// Start CPU profiling
	cpuFile, err := os.Create("debug_cpu.prof")
	if err != nil {
		log.Printf("Failed to create CPU profile: %v", err)
		return
	}
	dt.profiles["cpu"] = cpuFile
	pprof.StartCPUProfile(cpuFile)
	
	// Start memory profiling
	memFile, err := os.Create("debug_memory.prof")
	if err != nil {
		log.Printf("Failed to create memory profile: %v", err)
		return
	}
	dt.profiles["memory"] = memFile
	
	fmt.Println("Profiling started. CPU profile: debug_cpu.prof, Memory profile: debug_memory.prof")
}

func (dt *DebugTool) StopProfiling() {
	// Stop CPU profiling
	if cpuFile, exists := dt.profiles["cpu"]; exists {
		pprof.StopCPUProfile()
		cpuFile.Close()
		delete(dt.profiles, "cpu")
	}
	
	// Write memory profile
	if memFile, exists := dt.profiles["memory"]; exists {
		pprof.WriteHeapProfile(memFile)
		memFile.Close()
		delete(dt.profiles, "memory")
	}
	
	fmt.Println("Profiling stopped.")
}

func main() {
	var (
		port     = flag.Int("port", 6060, "Port for debug server")
		profile  = flag.Bool("profile", false, "Start profiling on startup")
	)
	flag.Parse()
	
	debugTool := NewDebugTool(*port)
	
	if *profile {
		debugTool.StartProfiling()
		defer debugTool.StopProfiling()
	}
	
	debugTool.Start()
} 