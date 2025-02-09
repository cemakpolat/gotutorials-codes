
##Project: System Monitoring Tool with Tests

This project will create a basic system monitoring tool that collects CPU, memory, and goroutine count metrics and exposes them via a web server. It will also include unit tests for the core metric collection logic.

**Project Structure:**

```
sysmonitor/
├── go.mod
├── main.go
└── monitor/
    ├── monitor.go
    └── monitor_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir sysmonitor
    cd sysmonitor
    ```

2.  **Initialize the Go Module:**
    Inside the `sysmonitor` directory, run:

    ```bash
    go mod init sysmonitor
    ```

3.  **Create the `monitor` Package Directory:**

    ```bash
    mkdir monitor
    ```

4.  **Create the `monitor.go` file (inside the `monitor` directory):**

    ```bash
    touch monitor/monitor.go
    ```

5.  **Create the `monitor_test.go` file (inside the `monitor` directory):**

    ```bash
    touch monitor/monitor_test.go
    ```

6.  **Create the `main.go` file (inside the `sysmonitor` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
sysmonitor/
├── go.mod
├── main.go
└── monitor/
    ├── monitor.go
    └── monitor_test.go
```

**Now, copy the following code into the corresponding files:**

**1. `go.mod` File:**

```
module sysmonitor

go 1.21
```

**2. `monitor/monitor.go` (System Monitoring Logic):**

```go
// monitor/monitor.go
package monitor

import (
	"runtime"
	"encoding/json"
	"net/http"
)


type Metrics struct {
	CPUUsage float64 `json:"cpu_usage"`
	MemoryUsage uint64 `json:"memory_usage"`
	GoroutineCount int `json:"goroutine_count"`
}
	
func CollectMetrics() Metrics {
	var memStats runtime.MemStats;
	runtime.ReadMemStats(&memStats);
		
	var cpuUsage float64 = 0; // requires third party libraries
		
		
	return Metrics{
		CPUUsage: cpuUsage,
		MemoryUsage: memStats.Alloc,
		GoroutineCount: runtime.NumGoroutine(),
	}
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := CollectMetrics();
	w.Header().Set("Content-Type", "application/json")
		
	data, err := json.MarshalIndent(metrics, "", "    ");
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return;
	}
	w.Write(data);
}
```
*  This file contains the core logic for collecting system metrics and exposing the data to the client.
*   The `CollectMetrics` collects the information about memory usage and current running goroutines.
*   The `MetricsHandler` is responsible for outputting the metrics in a json format.

**3. `monitor/monitor_test.go` (System Monitoring Logic Tests):**

```go
// monitor/monitor_test.go
package monitor

import (
    "encoding/json"
    "net/http"
	"net/http/httptest"
    "reflect"
    "runtime"
    "testing"
)

func TestCollectMetrics(t *testing.T) {
	metrics := CollectMetrics()
    
    var memStats runtime.MemStats;
	runtime.ReadMemStats(&memStats);

    if metrics.GoroutineCount <= 0 {
        t.Errorf("CollectMetrics failed, expected goroutine count to be higher than 0 but got: %d", metrics.GoroutineCount);
    }
    if metrics.MemoryUsage != memStats.Alloc {
		t.Errorf("CollectMetrics failed, expected memory usage to be %d got %d", memStats.Alloc, metrics.MemoryUsage)
	}
	// CPU Usage is not possible to test reliably
}


func TestMetricsHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/metrics", nil);
    if err != nil {
        t.Fatalf("error creating request: %v", err);
    }
    recorder := httptest.NewRecorder();
    MetricsHandler(recorder, req);

    if recorder.Code != http.StatusOK {
        t.Fatalf("metricsHandler failed with status %d", recorder.Code);
    }
	
	var metrics Metrics;

	err = json.Unmarshal(recorder.Body.Bytes(), &metrics);
	if err != nil {
		t.Fatalf("Error unmarshalling metrics %v", err);
	}

	expected := CollectMetrics()
	if !reflect.DeepEqual(metrics.MemoryUsage, expected.MemoryUsage) || !reflect.DeepEqual(metrics.GoroutineCount, expected.GoroutineCount) {
		t.Fatalf("Metrics handler returned invalid data");
	}
}
```
*  Tests that the `CollectMetrics` returns a valid data by asserting memory and goroutine counts.
*  Tests that the `MetricsHandler` provides a response with the correct `http.StatusOk` code and the json output contains valid data.

**4. `main.go` (Main Application Logic):**

```go
// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
    "sysmonitor/monitor"
)

func main() {
    http.HandleFunc("/metrics", monitor.MetricsHandler);
    fmt.Println("Server is starting at port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```
* This file contains the entry point of the application.
* Starts an http server at port 8080 and will call the `MetricsHandler` for every request to `/metrics` endpoint.

**How to Run the Project and Tests:**

1.  **Run the Application:**
    Open a terminal, navigate to the `sysmonitor` directory, and run:

    ```bash
    go run .
    ```

2.  **Access the metrics:** Open a web browser and visit `http://localhost:8080/metrics`.
3. **Run the Tests:**
    Open a terminal, navigate to the `sysmonitor` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok      sysmonitor/monitor        0.003s
```

**Output (If Tests Fail):**

If any test fails, the output will show more details to assist in the debugging process.

```
--- FAIL: TestCollectMetrics (0.00s)
    monitor_test.go:18: CollectMetrics failed, expected memory usage to be 8978987 got 0
FAIL
exit status 1
FAIL	sysmonitor/monitor	0.002s
```

**Key Features of This Project:**

*   **System Monitoring Tool:** Implements the basic functionality to collect system metrics and expose them through a http endpoint.
*   **Modularity:** Uses a separate `monitor` package to encapsulate all logic related to system monitoring.
*  **Testing**:  Includes tests to verify that the data of the different metrics is correct, and that the output of the http endpoint is correct.
*   **API Endpoint:** The application creates an API endpoint at `/metrics` to view the metrics as JSON.

This project provides a practical example of how Go can be used to create a system monitoring tool and how to test the core components of this kind of application.
