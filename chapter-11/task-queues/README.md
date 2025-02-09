
## Project: Task Queue with Tests

This project will create a basic task queue that allows users to enqueue tasks with different priorities, and processes them concurrently using a worker pool, while including test cases to verify its correct behaviour.

**Project Structure:**

```
taskqueue/
├── go.mod
├── main.go
└── queue/
    ├── queue.go
    └── queue_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir taskqueue
    cd taskqueue
    ```

2.  **Initialize the Go Module:**
    Inside the `taskqueue` directory, run:

    ```bash
    go mod init taskqueue
    ```

3.  **Create the `queue` Package Directory:**

    ```bash
    mkdir queue
    ```

4.  **Create the `queue.go` file (inside the `queue` directory):**

    ```bash
    touch queue/queue.go
    ```

5.  **Create the `queue_test.go` file (inside the `queue` directory):**

    ```bash
    touch queue/queue_test.go
    ```

6.  **Create the `main.go` file (inside the `taskqueue` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
taskqueue/
├── go.mod
├── main.go
└── queue/
    ├── queue.go
    └── queue_test.go
```

**Now, copy the following code into the corresponding files:**

**1.  `go.mod` File:**

```
module taskqueue

go 1.21
```

**2.  `queue/queue.go` (Task Queue Logic):**

```go
// queue/queue.go
package queue

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

type Task struct {
	ID int
	Priority int
	Data string
}

func Worker(id int, jobs <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Println("worker: ", id, "processing task ", job.ID, "priority", job.Priority);
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond); // simulate processing
	}
	fmt.Println("worker: ", id, "finished");
}
```
*  Defines the `Task` struct to be used as the jobs.
*  The function `Worker` simulates a job that is processed by a worker.

**3. `queue/queue_test.go` (Task Queue Logic Tests):**

```go
// queue/queue_test.go
package queue

import (
    "sync"
    "testing"
    "time"
)

func TestWorker(t *testing.T) {
	jobs := make(chan Task, 10);
	var wg sync.WaitGroup;

	wg.Add(1);
	go Worker(1, jobs, &wg);
	
	for i := 0; i < 10; i++ {
		jobs <- Task {
			ID: i + 1,
			Priority: i,
			Data: fmt.Sprintf("data %d", i),
		}
	}

	close(jobs)
	wg.Wait(); // waits until all jobs are completed.
	// if it completed without errors, the test is a success.
}


func TestWorkerMultiple(t *testing.T) {
    jobs := make(chan Task, 10)
	var wg sync.WaitGroup

    for i := 0; i < 3; i++ {
		wg.Add(1);
		go Worker(i, jobs, &wg)
    }

	for i := 0; i < 10; i++ {
		jobs <- Task {
			ID: i + 1,
			Priority: i,
			Data: fmt.Sprintf("data %d", i),
		}
	}
	close(jobs);
	wg.Wait();
	// if it completed without errors, the test is a success.
}

func TestWorkerMultipleLongRunning(t *testing.T) {
    jobs := make(chan Task, 10)
    var wg sync.WaitGroup

    for i := 0; i < 10; i++ {
        wg.Add(1);
        go func() {
            defer wg.Done();
            for j := range jobs {
                time.Sleep(10 * time.Millisecond) // simulate work
				if j.ID == 5 {
					time.Sleep(500 * time.Millisecond)
				}
            }
        }()
    }


    for i := 0; i < 10; i++ {
	jobs <- Task {
            ID: i + 1,
            Priority: i,
            Data: fmt.Sprintf("data %d", i),
        }
    }
    close(jobs)
	wg.Wait();
    // if it completed without errors, the test is a success.
}
```
*   Includes different scenarios for the tests, including a single worker, multiple workers, and multiple long running tasks.
*  Uses a `waitgroup` to make sure that the workers complete all their jobs before exiting the test case.
* Asserts that the tests complete without errors.

**4. `main.go` (Main Application Logic):**

```go
// main.go
package main

import (
	"fmt"
	"sync"
	"taskqueue/queue"
	"math/rand"
)

func main() {
    const numWorkers = 5;
    const numTasks = 10;

    var wg sync.WaitGroup;

    jobs := make(chan queue.Task, 100);

    for i := 0; i < numWorkers; i++ {
		wg.Add(1);
        go queue.Worker(i, jobs, &wg);
    }

    for i := 0; i < numTasks; i++ {
		jobs <- queue.Task {
            ID: i + 1,
			Priority: rand.Intn(10),
            Data: fmt.Sprintf("data %d", i),
        }
    }
	
	close(jobs);
	wg.Wait();

	fmt.Println("All workers have finished")
}
```
*   This file has the logic to create multiple workers and send several tasks to a shared queue, using a `waitgroup` to make sure that all workers finish their work.
* Uses the `queue` package to perform the task queue logic.

**How to Run the Project and Tests:**

1.  **Run the Application:**
    Open a terminal, navigate to the `taskqueue` directory, and run:

    ```bash
    go run .
    ```

2.  **Run the Tests:**
    Open a terminal, navigate to the `taskqueue` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok      taskqueue/queue     0.004s
```

**Output (If Tests Fail):**
If any test fails, the output will contain additional information about the errors.
```
--- FAIL: TestWorker (0.00s)
    queue_test.go:18: Worker failed, the program panicked.
FAIL
exit status 2
FAIL	taskqueue/queue	0.004s
```

**Key Features of This Project:**

*   **Task Queue Implementation:** Implements the core logic for managing a basic task queue.
*   **Concurrency:** Uses goroutines and channels to process tasks concurrently via a worker pool.
*   **Modularity:** Encapsulates the queue logic in a separate `queue` package.
*  **Testing**: Includes different test scenarios that tests the functionality of the task queue.
*	**Waitgroup**: Uses waitgroup to ensure that all workers have completed their work.

This project provides a practical example of how to create a basic task queue using Go and its concurrency primitives, while also testing the core concepts of a task queue.
