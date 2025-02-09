package monitor

import (
	"encoding/json"
	"net/http"
	"runtime"
)

type Metrics struct {
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryUsage    uint64  `json:"memory_usage"`
	GoroutineCount int     `json:"goroutine_count"`
}

func CollectMetrics() Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var cpuUsage float64 = 0 //TODO: requires third party libraries

	return Metrics{
		CPUUsage:       cpuUsage,
		MemoryUsage:    memStats.Alloc,
		GoroutineCount: runtime.NumGoroutine(),
	}
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := CollectMetrics()
	w.Header().Set("Content-Type", "application/json")

	data, err := json.MarshalIndent(metrics, "", "    ")
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
