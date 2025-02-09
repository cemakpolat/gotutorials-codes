package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"
)

var CollectMetricsFunc = CollectMetrics

func TestCollectMetrics(t *testing.T) {
	metrics := CollectMetricsFunc()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	if metrics.GoroutineCount <= 0 {
		t.Errorf("CollectMetrics failed, expected goroutine count to be higher than 0 but got: %d", metrics.GoroutineCount)
	}
	if metrics.MemoryUsage != memStats.Alloc {
		t.Errorf("CollectMetrics failed, expected memory usage to be %d got %d", memStats.Alloc, metrics.MemoryUsage)
	}
	// CPU Usage is not possible to test reliably
}
func TestMetricsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}

	// Create a mock CollectMetrics to control the values.
	var memStats runtime.MemStats
	var readMemStatsFunc = runtime.ReadMemStats
	readMemStatsFunc(&memStats)
	expectedMetrics := Metrics{
		CPUUsage:       0,
		MemoryUsage:    memStats.Alloc,
		GoroutineCount: runtime.NumGoroutine(),
	}
	originalCollectMetrics := CollectMetricsFunc
	originalReadMemStats := readMemStatsFunc
	CollectMetricsFunc = func() Metrics {
		return expectedMetrics
	}
	readMemStatsFunc = func(m *runtime.MemStats) {
		*m = memStats
	}
	defer func() {
		CollectMetricsFunc = originalCollectMetrics
		readMemStatsFunc = originalReadMemStats
	}()

	recorder := httptest.NewRecorder()
	MetricsHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("metricsHandler failed with status %d", recorder.Code)
	}

	var metrics Metrics

	err = json.Unmarshal(recorder.Body.Bytes(), &metrics)
	if err != nil {
		t.Fatalf("Error unmarshalling metrics %v", err)
	}

	if !reflect.DeepEqual(metrics, expectedMetrics) {
		t.Fatalf("Metrics handler returned invalid data, Expected: %v, Got: %v", expectedMetrics, metrics)
	}
}
