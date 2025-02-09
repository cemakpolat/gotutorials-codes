Create a project that integrates the HTTP and API concepts we've covered in Chapter 8. This project will involve creating a simple HTTP server that fetches data from an external API, parses the JSON response, and displays the information to the client.

**Project: Simple API Fetcher**

This project will:

1.  Create an HTTP server that listens on a specific port.
2.  Define a handler for a specific path (`/fetch`).
3.  Inside the handler, fetch data from the JSONPlaceholder API (or any other public API).
4.  Parse the JSON response.
5.  Display the parsed data in the HTTP response.
6.  Include a `utils` package that has all the logic to perform http requests and parse json.

**Project Structure:**

```
apifetcher/
├── go.mod
├── main.go
└── utils/
    ├── utils.go
	└── utils_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run the following commands:

    ```bash
    mkdir apifetcher
    cd apifetcher
    ```

2.  **Initialize the Go Module:**
    Inside the `apifetcher` directory, run:

    ```bash
    go mod init apifetcher
    ```

3.  **Create the `utils` Package Directory:**

    ```bash
     mkdir utils
    ```

4. **Create the `utils.go` file (inside the `utils` directory):**
	```bash
	 touch utils/utils.go
	```

5. **Create the `utils_test.go` file (inside the `utils` directory):**
	```bash
	 touch utils/utils_test.go
	```

6.  **Create the `main.go` file (inside the `apifetcher` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
apifetcher/
├── go.mod
├── main.go
└── utils/
    ├── utils.go
	└── utils_test.go
```

**Now, paste the code into the respective files:**

**1. `go.mod` File:**

```
module apifetcher

go 1.21
```

**2. `utils/utils.go` (Utility functions):**
```go
// utils/utils.go
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
func FetchTodo(url string) (*Todo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making http request %w", err);
	}
	defer resp.Body.Close();
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 status code %d", resp.StatusCode)
	}
	
	data, err := io.ReadAll(resp.Body);
	if err != nil {
		return nil, fmt.Errorf("error reading the body %w", err);
	}
	
	var todo Todo;
	if err := json.Unmarshal(data, &todo); err != nil {
		return nil, fmt.Errorf("error unmarshalling json %w", err);
	}

	return &todo, nil;
}
```
* This package will contain utility functions for making http requests, and parsing data.

**3. `utils/utils_test.go`:**
```go
// utils/utils_test.go
package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchTodo(t *testing.T) {
    // Create a mock server
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"userId": 1, "id": 1, "title": "Test Todo", "completed": false}`))
    }))
    defer testServer.Close()

    // Call the function with the mock server URL
    todo, err := FetchTodo(testServer.URL)
    if err != nil {
        t.Errorf("FetchTodo() failed with error: %v", err)
        return
    }

    // Assert the expected value
    if todo.Title != "Test Todo" {
        t.Errorf("FetchTodo() returned incorrect title: got %v, want %v", todo.Title, "Test Todo")
    }
	
	if todo.ID != 1 {
		t.Errorf("FetchTodo() returned incorrect id: got %d, want %d", todo.ID, 1);
	}
}


func TestFetchTodoFailure(t *testing.T) {
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer testServer.Close()

    _, err := FetchTodo(testServer.URL)
    if err == nil {
        t.Errorf("FetchTodo() should have failed with 500 error")
    }
}
```
* This test will test both failure and success of the FetchTodo function.
*   The test creates a mock server to make sure that the function does not depend on an external api.

**4. `main.go` (Main Application):**

```go
// main.go
package main

import (
    "fmt"
    "net/http"
	"apifetcher/utils"
	"log"
)

func fetchHandler(w http.ResponseWriter, r *http.Request) {
    url := "https://jsonplaceholder.typicode.com/todos/1"
	
	todo, err := utils.FetchTodo(url)
	if err != nil {
		http.Error(w, "Failed to fetch the data", http.StatusInternalServerError);
		log.Println("Error fetching data:", err);
		return;
	}


    fmt.Fprintf(w, "Todo ID: %d\n", todo.ID)
    fmt.Fprintf(w, "Title: %s\n", todo.Title)
	fmt.Fprintf(w, "Completed: %t", todo.Completed)

}

func main() {
    http.HandleFunc("/fetch", fetchHandler)
    fmt.Println("Server is starting at port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

*   Defines a handler for the `/fetch` path that fetches data from the JSONPlaceholder API and returns it to the client.
*   Includes error handling for API requests.
*  Starts the HTTP server at port 8080.

**How to Run the Project and Tests:**

1.  **Run the application:**
    Open a terminal, navigate to the `apifetcher` directory, and run:

    ```bash
    go run .
    ```

2.  **Access the endpoint:** Open a web browser and visit `http://localhost:8080/fetch`. You should see the data coming from JSONPlaceholder.
3.  **Run the tests:**
    Open a terminal, navigate to the `apifetcher` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok      apifetcher/utils	0.003s
```

**Output (If Tests Fail):**

```
--- FAIL: TestFetchTodo (0.00s)
    utils_test.go:24: FetchTodo() returned incorrect title: got Test, want Test Todo
FAIL
exit status 1
FAIL    apifetcher/utils        0.003s
```

**Key Features of This Project:**

*   **HTTP Server:** Creates a basic HTTP server.
*   **API Request:** Makes requests to an external API.
*   **JSON Parsing:** Parses the JSON response using Go's standard library.
*   **Modularity:** Code is organized by packages using `utils` and `main`.
*   **Error Handling:** Includes error handling for API requests and server operations.
*   **Testing:** Tests for the `fetchTodo` function are included in `utils_test.go`.

This example provides a complete demonstration of using Go's HTTP and API capabilities and gives the reader a practical example to apply what they have learned in Chapter 8.
