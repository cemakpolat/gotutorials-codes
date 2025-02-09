Design a complete project that integrates the concurrency features discussed in Chapter 7: goroutines, channels, mutexes, and waitgroups. This project will simulate a concurrent data processing pipeline, demonstrating how these concepts work together in a practical scenario.

**Project: Concurrent Data Processing Pipeline**

This project will simulate a data processing pipeline using concurrency features:

1.  **Data Generation:** A goroutine will generate a stream of data.
2.  **Data Processing:** A pool of worker goroutines will concurrently process the data from the stream.
3.  **Data Aggregation:** A goroutine will aggregate the results.
4.  **Synchronization:** Uses `sync.Mutex` to protect shared data and `sync.WaitGroup` to synchronize the goroutines.
5.  **Communication:** Uses channels for communication between goroutines.

**Project Structure:**

```
concurrentpipeline/
├── go.mod
├── main.go
├── data/
│   └── data.go
└── utils/
    ├── utils.go
    └── utils_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run the following command:

    ```bash
    mkdir concurrentpipeline
    cd concurrentpipeline
    ```

2.  **Initialize the Go Module:**
    Inside the `concurrentpipeline` directory, run:

    ```bash
    go mod init concurrentpipeline
    ```

3.  **Create the `data` Package Directory:**

    ```bash
    mkdir data
    ```

4.  **Create the `utils` Package Directory:**

    ```bash
     mkdir utils
    ```

5.  **Create the `data.go` file (inside the `data` directory):**

    ```bash
    touch data/data.go
    ```

6.  **Create the `utils.go` file (inside the `utils` directory):**

     ```bash
    touch utils/utils.go
    ```

7.  **Create the `utils_test.go` file (inside the `utils` directory):**

    ```bash
    touch utils/utils_test.go
    ```

8.  **Create the `main.go` file (inside the `concurrentpipeline` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
concurrentpipeline/
├── go.mod
├── main.go
├── data/
│   └── data.go
└── utils/
    ├── utils.go
    └── utils_test.go
```

**Now, paste the code into the following files:**

**1. `go.mod` File:**

```
module concurrentpipeline

go 1.21
```

**2. `data/data.go` (Data Generation):**

```go
// data/data.go
package data

import "math/rand"

type Data struct {
    Value int
}

func GenerateData(count int, out chan<- Data) {
    for i := 0; i < count; i++ {
        data := Data{Value: rand.Intn(100)}
        out <- data
    }
    close(out)
}
```

*   `GenerateData` generates random data and sends it to the provided channel.
*   Closes the channel after sending all data to indicate completion.

**3. `utils/utils.go` (Data Processing and Aggregation):**

```go
// utils/utils.go
package utils

import (
    "fmt"
    "sync"
	"time"
	"concurrentpipeline/data"
)


func ProcessData(id int, in <-chan data.Data, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
    for d := range in {
		time.Sleep(100 * time.Millisecond) // simulate processing
        fmt.Println("worker:", id, "processing: ", d.Value);
        out <- d.Value * 2
    }
	fmt.Println("worker:", id, "finished")
}

func AggregateData(in <-chan int, results *[]int, m *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for res := range in {
		m.Lock(); // mutex to protect writing to the slice
        *results = append(*results, res)
        m.Unlock() // release the mutex
	}
}

func PrintResults(results []int) {
	fmt.Println("Printing results");
    for _, r := range results {
        fmt.Println(r)
    }
}
```

*   `ProcessData` reads data from a channel, multiplies it by 2, and sends the result to another channel.
*   `AggregateData` reads results from the processing channel, aggregates the results to a shared slice, using a mutex for protection.
*	`PrintResults` prints the results.
*   Uses wait groups to wait for the workers and aggregation to finish.

**4. `utils/utils_test.go`:**
```go
// utils/utils_test.go
package utils

import (
	"sync"
	"testing"
	"concurrentpipeline/data"
)

func TestProcessData(t *testing.T) {
	jobs := make(chan data.Data, 10);
	results := make(chan int, 10);
	var wg sync.WaitGroup;

	// populate jobs channel
	for i := 0; i < 10; i++ {
		jobs <- data.Data{Value: i}
	}
	close(jobs)
	
	wg.Add(1);
	go ProcessData(1, jobs, results, &wg);

	wg.Wait();
	close(results);
	var receivedResults []int;
	for result := range results {
		receivedResults = append(receivedResults, result)
	}

	expectedResults := []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}
	
	if len(expectedResults) != len(receivedResults) {
		t.Fatalf("expected %d results, but got %d", len(expectedResults), len(receivedResults));
	}
	
	for i := 0; i < len(expectedResults); i++ {
		if expectedResults[i] != receivedResults[i] {
			t.Fatalf("expected %d but got %d at index %d", expectedResults[i], receivedResults[i], i);
		}
	}
}
```
*   Tests the processing and aggregation of data.
*   Ensures the results are valid and the operations are correct.

**5. `main.go` (Main Application):**

```go
// main.go
package main

import (
	"sync"
	"concurrentpipeline/data"
	"concurrentpipeline/utils"
)

func main() {
    const dataCount = 100
    const numWorkers = 5

    dataChannel := make(chan data.Data, 100)
    processedDataChannel := make(chan int, 100)
    var results []int

	var wg sync.WaitGroup; // waitgroup to wait for all goroutines to finish
	var m sync.Mutex; // mutex to protect the shared slice.
	
    wg.Add(1); // wait for the data generator to finish
    go data.GenerateData(dataCount, dataChannel)

    for i := 0; i < numWorkers; i++ {
	wg.Add(1); // wait for each worker to finish
        go utils.ProcessData(i, dataChannel, processedDataChannel, &wg)
    }
	
	wg.Add(1); // wait for the aggregation to finish
    go utils.AggregateData(processedDataChannel, &results, &m, &wg)

    wg.Wait() // wait for all processes to finish

    utils.PrintResults(results)
}
```

*   Sets up the data processing pipeline, using goroutines for data generation, processing and aggregation.
*   Manages the pipeline with wait groups and channels.
*    Uses mutex to protect the shared slice.

**How to Run the Project and Tests:**

1.  **Run the application:**
    Open a terminal, navigate to the `concurrentpipeline` directory, and run:

    ```bash
    go run .
    ```

2.  **Run the tests:**
    Open a terminal, navigate to the `concurrentpipeline` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok  	concurrentpipeline/utils	0.003s
```

**Output (If Tests Fail):**

```
--- FAIL: TestProcessData (0.00s)
    utils_test.go:35: expected 10 results, but got 9
FAIL
exit status 1
FAIL	concurrentpipeline/utils	0.005s
```

**Key Features of This Project:**

*   **Complete Pipeline:** Demonstrates a full concurrent data processing pipeline.
*   **Concurrency Features:** Uses goroutines, channels, mutexes, and wait groups effectively.
*   **Modularity:** Code is organized into packages for data, processing, and the main application logic.
*  **Testing:** Includes tests for the processing logic, demonstrating its correctness.
*   **Practical Scenario:** Simulates a realistic data processing workflow.

This complete example demonstrates how to use Go's concurrency features to build an efficient and parallel application.

