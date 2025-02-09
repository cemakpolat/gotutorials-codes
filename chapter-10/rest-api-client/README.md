Implement the second example from Chapter 10, the "REST API Client," as a complete project with tests. This project will focus on making HTTP requests to an external API, handling the JSON response, and printing the data.

**Project: REST API Client with Tests**

This project will create a command-line tool that fetches data from a specified REST API endpoint and displays the relevant information, all while including tests to ensure the reliability of the core logic.

**Project Structure:**

```
apiclient/
├── go.mod
├── main.go
└── client/
    ├── client.go
    └── client_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir apiclient
    cd apiclient
    ```

2.  **Initialize the Go Module:**
    Inside the `apiclient` directory, run:

    ```bash
    go mod init apiclient
    ```

3.  **Create the `client` Package Directory:**

    ```bash
    mkdir client
    ```

4.  **Create the `client.go` file (inside the `client` directory):**

    ```bash
    touch client/client.go
    ```

5.  **Create the `client_test.go` file (inside the `client` directory):**

    ```bash
    touch client/client_test.go
    ```

6.  **Create the `main.go` file (inside the `apiclient` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
apiclient/
├── go.mod
├── main.go
└── client/
    ├── client.go
    └── client_test.go
```

**Now, paste the following code into the corresponding files:**

**1. `go.mod` File:**

```
module apiclient

go 1.21
```

**2. `client/client.go` (API Client Logic):**

```go
// client/client.go
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func FetchAndPrintData(url string) (*Todo, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close();

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 status code: %d", resp.StatusCode)
	}

    data, err := io.ReadAll(resp.Body);
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err);
	}


	var todo Todo;
    if err := json.Unmarshal(data, &todo); err != nil {
		return nil, fmt.Errorf("error unmarshalling json %w", err);
	}
	
    fmt.Printf("Todo ID: %d\n", todo.ID)
	fmt.Printf("Title: %s\n", todo.Title)
	fmt.Printf("Completed: %t\n", todo.Completed)
	return &todo, nil
}
```
*   Defines the struct that is used to unmarshall the json payload.
*   `FetchAndPrintData` performs an http get request, handles any potential error, and prints the result to standard output.
*  The functions are exported, allowing them to be called from `main.go` and `client_test.go`.

**3. `client/client_test.go` (API Client Logic Tests):**

```go
// client/client_test.go
package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchAndPrintDataSuccess(t *testing.T) {
    // Create a mock server
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintln(w, `{"userId": 1, "id": 1, "title": "Test Todo", "completed": false}`)
    }))
    defer testServer.Close()

    // Call the function with the mock server URL
	todo, err := FetchAndPrintData(testServer.URL);
    if err != nil {
        t.Errorf("FetchAndPrintData() failed with error: %v", err)
		return;
    }
	
	if todo.Title != "Test Todo" {
		t.Fatalf("Returned incorrect todo item, expected title Test Todo, got %s", todo.Title)
	}
}

func TestFetchAndPrintDataFailure(t *testing.T) {
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer testServer.Close()
    
	_, err := FetchAndPrintData(testServer.URL)
	if err == nil {
		t.Fatalf("FetchAndPrintData should have failed with 500");
	}
}
```
*   Tests both successful and failed HTTP requests, using a mock server.
*   Asserts that the output is valid for the success scenario and the test fails when an error was expected.

**4. `main.go` (Main Application Logic):**

```go
// main.go
package main

import (
	"fmt"
	"log"
    "os"
	"apiclient/client"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <url>");
		os.Exit(1);
	}
    url := os.Args[1]

	_, err := client.FetchAndPrintData(url);
    if err != nil {
		log.Fatal("error fetching data", err);
    }
}
```

*   Imports the `client` package.
*   Gets the url from the command line arguments.
*   Calls the `FetchAndPrintData` function to perform the API call, log any errors, and return to the user.

**How to Run the Project and Tests:**

1.  **Run the Application:**
    Open a terminal, navigate to the `apiclient` directory, and run:

    ```bash
    go run . https://jsonplaceholder.typicode.com/todos/1
    ```
Replace the url for your desired endpoint.

2.  **Run the Tests:**
    Open a terminal, navigate to the `apiclient` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok      apiclient/client        0.004s
```

**Output (If Tests Fail):**
If any test fails, the output will show details of the error to help in the debugging process.
```
--- FAIL: TestFetchAndPrintDataFailure (0.00s)
    client_test.go:46: FetchAndPrintData should have failed with 500
FAIL
exit status 1
FAIL	apiclient/client	0.004s
```

**Key Features of This Project:**

*   **API Client:** Implements the core logic for making requests to a REST API.
*   **JSON Parsing:** Successfully parses the data returned by the API using `encoding/json`.
*   **Modularity:** Encapsulates the logic for API calls and output in the `client` package.
*   **Testing:** Includes unit tests that cover the success and failure cases when making an api request.
* **Arguments**: Reads arguments from the command line to get the URL of the desired endpoint.

This project provides a practical example of how to build a REST API client in Go, and demonstrates the modularity, testing and error handling concepts.
