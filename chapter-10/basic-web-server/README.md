Okay, let's implement the first example from Chapter 10, the "Basic Web Server with Routing," as a complete project with tests. We will focus on providing a solid foundation for a web server, incorporating basic routing and testing.

**Project: Basic Web Server with Routing and Tests**

This project will create a web server that handles different routes and responds with varying content, including handling query parameters, and handling POST requests. We'll also include tests to ensure that our routing and response logic are correct.

**Project Structure:**

```
webserver/
├── go.mod
├── main.go
└── server/
    ├── server.go
    └── server_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir webserver
    cd webserver
    ```

2.  **Initialize the Go Module:**
    Inside the `webserver` directory, run:

    ```bash
    go mod init webserver
    ```

3.  **Create the `server` Package Directory:**

    ```bash
    mkdir server
    ```

4.  **Create the `server.go` file (inside the `server` directory):**

    ```bash
    touch server/server.go
    ```

5.  **Create the `server_test.go` file (inside the `server` directory):**

    ```bash
    touch server/server_test.go
    ```

6.  **Create the `main.go` file (inside the `webserver` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
webserver/
├── go.mod
├── main.go
└── server/
    ├── server.go
    └── server_test.go
```

**Now, paste the following code into their corresponding files:**

**1. `go.mod` File:**

```
module webserver

go 1.21
```

**2. `server/server.go` (Server Logic):**

```go
// server/server.go
package server

import (
    "fmt"
	"net/http"
	"encoding/json"
	"io"
)

type Post struct {
	Message string `json:"message"`
	Author string `json:"author"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the Home page!")
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "This is the about page!")
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Query parameters:", r.URL.Query());

	switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "This is a list of posts!");
		case http.MethodPost:
			body, err := io.ReadAll(r.Body);
			if err != nil {
				http.Error(w, "Error reading the body", http.StatusBadRequest);
				return;
			}
			var post Post;

			err = json.Unmarshal(body, &post);

			if err != nil {
				http.Error(w, "Error unmarshalling json", http.StatusBadRequest);
				return;
			}
			fmt.Fprintf(w, "Post created:\n");
			fmt.Fprintf(w, "Message: %s\nAuthor: %s\n", post.Message, post.Author);
		default:
		    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed);
		}
}
```
* This file contains all handlers for the http server.
*  `HomeHandler` is the handler for the `/` route.
* `AboutHandler` is the handler for the `/about` route.
* `PostsHandler` is the handler for the `/posts` route, and it handles both `GET` and `POST` requests.
*  It is also responsible for unmarshalling the json in the body of the `POST` requests.

**3. `server/server_test.go` (Server Logic Tests):**

```go
// server/server_test.go
package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil);
	if err != nil {
		t.Fatalf("Error creating request: %v", err);
	}
	
	recorder := httptest.NewRecorder();
	HomeHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("HomeHandler failed, expected %d, got %d", http.StatusOK, recorder.Code)
	}

	expected := "Welcome to the Home page!"
	actual := recorder.Body.String();
	
	if actual != expected {
		t.Errorf("HomeHandler failed, expected %q, got %q", expected, actual);
	}
}

func TestAboutHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/about", nil);
	if err != nil {
		t.Fatalf("Error creating request: %v", err);
	}

	recorder := httptest.NewRecorder();
	AboutHandler(recorder, req);

	if recorder.Code != http.StatusOK {
		t.Fatalf("AboutHandler failed, expected %d, got %d", http.StatusOK, recorder.Code)
	}

	expected := "This is the about page!"
	actual := recorder.Body.String();
	
	if actual != expected {
		t.Errorf("AboutHandler failed, expected %q, got %q", expected, actual);
	}
}

func TestPostsHandler(t *testing.T) {
    t.Run("Test GET request", func(t *testing.T) {
	    req, err := http.NewRequest(http.MethodGet, "/posts", nil)
	    if err != nil {
		    t.Fatalf("Error creating request: %v", err)
	    }

	    recorder := httptest.NewRecorder()
	    PostsHandler(recorder, req)

	    if recorder.Code != http.StatusOK {
		    t.Fatalf("PostsHandler failed, expected %d, got %d", http.StatusOK, recorder.Code)
	    }

	    expected := "This is a list of posts!"
	    actual := recorder.Body.String()
	    if actual != expected {
		t.Errorf("PostsHandler GET failed, expected %q, got %q", expected, actual);
	    }
    })

    t.Run("Test POST request", func(t *testing.T) {
	    body := `{"message": "Test message", "author": "test"}`
	    req, err := http.NewRequest(http.MethodPost, "/posts", strings.NewReader(body));
	    if err != nil {
		t.Fatalf("Error creating request: %v", err);
	    }

		req.Header.Set("Content-Type", "application/json");

	    recorder := httptest.NewRecorder()
	    PostsHandler(recorder, req);
	    
	    if recorder.Code != http.StatusOK {
		t.Fatalf("PostsHandler failed, expected %d, got %d", http.StatusOK, recorder.Code);
	    }

	    actual := recorder.Body.String()

	    expected := "Post created:\nMessage: Test message\nAuthor: test\n"
	    if actual != expected {
		t.Errorf("PostsHandler POST failed, expected %q, got %q", expected, actual);
	    }
    })
}
```
*   Tests all the handlers for their expected behaviours.
*   Uses httptest package for creating http recorder for testing http handlers.
*  Tests both the success scenarios, and verifies that the response is as expected.
*  The `TestPostsHandler` tests both `GET` and `POST` methods.

**4. `main.go` (Main Application Logic):**

```go
// main.go
package main

import (
	"fmt"
    "log"
    "net/http"
    "webserver/server"
)


func main() {
    http.HandleFunc("/", server.HomeHandler)
	http.HandleFunc("/about", server.AboutHandler)
	http.HandleFunc("/posts", server.PostsHandler)
    fmt.Println("Server is starting at port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

*   Imports the `server` package and uses its exported functions.
*   Starts the http server at port 8080, and logs any errors if they happen during the execution.

**How to Run the Project and Tests:**

1.  **Run the Application:**
    Open a terminal, navigate to the `webserver` directory, and run:

    ```bash
    go run .
    ```
2. **Test the application**: Open a web browser and visit `http://localhost:8080`, `http://localhost:8080/about` and `http://localhost:8080/posts`, and also use tools such as `curl` or `Postman` to send `POST` requests.

{
    "message":"Welcome to World",
    "author": "Cem Akpolat"
}

3.  **Run the Tests:**
    Open a terminal, navigate to the `webserver` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok      webserver/server        0.004s
```

**Output (If Tests Fail):**

If any of the tests fail, you will be given the corresponding output, which should help you debug what has gone wrong:

```
--- FAIL: TestHomeHandler (0.00s)
    server_test.go:24: HomeHandler failed, expected "welcome to the Home page!", got "Welcome to the Home page!"
FAIL
exit status 1
FAIL    webserver/server        0.003s
```

**Key Features of This Project:**

*   **Basic Web Server:** Creates a functional web server using Go's `net/http` package.
*   **Routing:** Handles requests on different routes using custom handlers.
*   **Query Parameters:** The `/posts` endpoint will log all query parameters passed to it.
*	**Post Methods:** The `/posts` endpoint can receive json data via a post request.
*   **Modularity:** The application logic is organized by packages.
*   **Testing:** Includes unit tests that verify the correctness of each handler.

This project provides a foundation for building more complex web servers in Go and combines routing, request handling, testing, and modularity.
