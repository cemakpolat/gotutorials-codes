Okay, let's implement the final example from Chapter 9, the "Command-Line File Downloader," as a complete project with tests. We'll maintain our established structure, setup instructions, and testing approach.

**Project: Command-Line File Downloader with Tests**

This project will create a command-line tool that downloads a file from a given URL and saves it to a specified local file path. It will include unit tests to ensure its functionality.

**Project Structure:**

```
downloader/
├── go.mod
├── main.go
└── download/
    ├── download.go
    └── download_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir downloader
    cd downloader
    ```

2.  **Initialize the Go Module:**
    Inside the `downloader` directory, run:

    ```bash
    go mod init downloader
    ```

3.  **Create the `download` Package Directory:**

    ```bash
    mkdir download
    ```

4.  **Create the `download.go` file (inside the `download` directory):**

    ```bash
    touch download/download.go
    ```

5.  **Create the `download_test.go` file (inside the `download` directory):**

    ```bash
    touch download/download_test.go
    ```

6.  **Create the `main.go` file (inside the `downloader` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
downloader/
├── go.mod
├── main.go
└── download/
    ├── download.go
    └── download_test.go
```

**Now, copy the following code into the corresponding files:**

**1.  `go.mod` File:**

```
module downloader

go 1.21
```

**2.  `download/download.go` (Download Logic):**

```go
// download/download.go
package download

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

func DownloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making get request %w", err);
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file %w", err);
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body);
	if err != nil {
		return fmt.Errorf("error copying data %w", err);
	}
	return nil;
}
```
* This file contains the logic to download a file from an URL to a local path.
* Handles potential errors such as incorrect response code, creation error and copy error.
* Exports the `DownloadFile` function, to be used in `main.go` and `download_test.go`.

**3.  `download/download_test.go` (Download Logic Tests):**

```go
// download/download_test.go
package download

import (
	"fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
)

func TestDownloadFileSuccess(t *testing.T) {
    // Create a mock server
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain");
        fmt.Fprintln(w, "Test file content")
    }))
    defer testServer.Close()

    // Call DownloadFile with the mock server URL
	filepath := "testfile.txt";
    err := DownloadFile(testServer.URL, filepath);
    if err != nil {
        t.Fatalf("DownloadFile() failed with error: %v", err)
    }

	defer os.Remove(filepath)
    // read the file content
	content, err := os.ReadFile(filepath);
	if err != nil {
		t.Fatalf("Error reading test file %v", err);
	}

	if string(content) != "Test file content\n" {
		t.Fatalf("Content is incorrect: %v", string(content));
	}
}


func TestDownloadFileFailure(t *testing.T) {
    testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer testServer.Close()

	filepath := "testfile.txt";
    err := DownloadFile(testServer.URL, filepath)
    if err == nil {
        t.Errorf("DownloadFile should have failed with 500 error")
    }

	_, err = os.Stat(filepath);
	if !os.IsNotExist(err) {
		t.Fatalf("File should not exist");
	}
}
```

*   Tests the `DownloadFile` with both valid and invalid URLs.
* Uses httptest to create a server and test the correct use cases.
* Assert the existence of the file, it's contents and the http response code.

**4. `main.go` (Main Application Logic):**

```go
// main.go
package main

import (
	"fmt"
    "os"
	"downloader/download"
)

func main() {
    if len(os.Args) != 3 {
	    fmt.Println("Usage: go run main.go <url> <filepath>");
	    os.Exit(1)
    }

	url := os.Args[1];
	filepath := os.Args[2];
	
	err := download.DownloadFile(url, filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading file: %v\n", err);
		os.Exit(1)
	}
	fmt.Println("File downloaded successfully!");
}
```

*   Imports the `download` package and uses the `DownloadFile` function to perform the downloading logic.
*   Uses the command-line arguments to get the URL and file path.
*   Outputs to standard error if an issue was found while downloading the file.

**How to Run the Project and Tests:**

1.  **Run the Application:**
    Open a terminal, navigate to the `downloader` directory, and run:

    ```bash
    go run . <URL> <filePath>
    ```

    Replace `<URL>` with the URL of the file you want to download and `<filePath>` with the desired local file path (including the file name). For example:
    ```bash
	go run . https://www.example.com/myfile.txt myfile.txt
    ```

2.  **Run the Tests:**
    Open a terminal, navigate to the `downloader` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok      downloader/download	0.004s
```

**Output (If Tests Fail):**

If any tests fail, the output will show details about the error and where the test failed.

```
--- FAIL: TestDownloadFileFailure (0.00s)
    download_test.go:44: DownloadFile should have failed with 500 error
FAIL
exit status 1
FAIL    downloader/download	0.003s
```

**Key Features of This Project:**

*   **Complete CLI Tool:** Implements the full file download functionality.
*   **Modularity:** Uses a `download` package to encapsulate the core logic.
*   **Testing:** Includes unit tests to verify the success and error cases of the download logic.
*   **Error Handling:** The program handles the most common errors like file creation error, network error and invalid HTTP status codes.
*   **Arguments:** Reads arguments from the command line to receive the URL and local path.

This completes our collection of practical CLI tool projects with tests, and demonstrates how to apply core go concepts for common CLI functionalities, combining testing, modularity and clear APIs.
