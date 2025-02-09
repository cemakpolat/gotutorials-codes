## Project: Log Analysis Tool with Tests

This project will create a command-line tool that analyzes log files, allowing users to filter log entries based on a given keyword and output the count of log entries based on their severity. The tool will also include unit tests to ensure the core analysis logic functions correctly.

**Project Structure:**

```
loganalyzer/
├── go.mod
├── main.go
└── analyzer/
    ├── analyzer.go
    └── analyzer_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir loganalyzer
    cd loganalyzer
    ```

2.  **Initialize the Go Module:**
    Inside the `loganalyzer` directory, run:

    ```bash
    go mod init loganalyzer
    ```

3.  **Create the `analyzer` Package Directory:**

    ```bash
    mkdir analyzer
    ```

4.  **Create the `analyzer.go` file (inside the `analyzer` directory):**

    ```bash
    touch analyzer/analyzer.go
    ```

5.  **Create the `analyzer_test.go` file (inside the `analyzer` directory):**

    ```bash
    touch analyzer/analyzer_test.go
    ```

6.  **Create the `main.go` file (inside the `loganalyzer` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
loganalyzer/
├── go.mod
├── main.go
└── analyzer/
    ├── analyzer.go
    └── analyzer_test.go
```

**Now, paste the following code into their corresponding files:**

**1. `go.mod` File:**

```
module loganalyzer

go 1.21
```

**2. `analyzer/analyzer.go` (Log Analysis Logic):**

```go
// analyzer/analyzer.go
package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func AnalyzeLogs(filename string, filter string) (map[string]int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    severityCounts := make(map[string]int)
	
    scanner := bufio.NewScanner(file);
	regex := regexp.MustCompile(filter);

    for scanner.Scan() {
        line := scanner.Text();
		if filter != "" && !regex.MatchString(line) {
			continue
		}
        
        parts := strings.SplitN(line, " ", 2);
		if len(parts) < 2 {
			continue
		}
        
		firstPart := parts[0]
		if strings.Contains(firstPart, "[") && strings.Contains(firstPart, "]") {
			severity := strings.ReplaceAll(firstPart, "[", "");
			severity = strings.ReplaceAll(severity, "]", "");
			severityCounts[severity]++;
		}
    }
	if err = scanner.Err(); err != nil {
		return nil, err
	}

    return severityCounts, nil
}
```
*   This file contains all the core logic of the log analysis.
*   `AnalyzeLogs` reads and parses the file, counts the errors based on the log severity, and filters the lines based on a regular expression.

**3. `analyzer/analyzer_test.go` (Log Analysis Logic Tests):**

```go
// analyzer/analyzer_test.go
package analyzer

import (
	"os"
	"reflect"
	"testing"
)

func TestAnalyzeLogs(t *testing.T) {
	testLog := `[info] this is an info log message
[error] this is an error message
[warning] this is a warning message
[info] another info message
[error] another error message with more details
[debug] some debug information
`

	filename := "test.log";

	err := os.WriteFile(filename, []byte(testLog), 0644);
	if err != nil {
		t.Fatalf("failed to create file: %v", err);
	}
	defer os.Remove(filename);

	t.Run("Test without filter", func(t *testing.T) {
		severityCounts, err := AnalyzeLogs(filename, "")
		if err != nil {
			t.Errorf("AnalyzeLogs failed with: %v", err)
		}
		
		expected := map[string]int{
			"info": 2,
			"error": 2,
			"warning": 1,
			"debug": 1,
		}
		if !reflect.DeepEqual(severityCounts, expected) {
			t.Errorf("AnalyzeLogs failed: expected: %v, got: %v", expected, severityCounts);
		}
	})

	t.Run("Test with filter", func(t *testing.T) {
		severityCounts, err := AnalyzeLogs(filename, "error");
		if err != nil {
			t.Errorf("AnalyzeLogs failed with filter, error: %v", err);
		}
		expected := map[string]int{
			"error": 2,
		}

		if !reflect.DeepEqual(severityCounts, expected) {
			t.Errorf("AnalyzeLogs failed with filter, expected: %v, got: %v", expected, severityCounts);
		}
	})

	t.Run("Test with invalid logs", func(t *testing.T) {
		testLog := `invalid log`
		err = os.WriteFile(filename, []byte(testLog), 0644);
		if err != nil {
			t.Fatalf("failed to write test file: %v", err);
		}
		severityCounts, err := AnalyzeLogs(filename, "")
		if err != nil {
			t.Fatalf("AnalyzeLogs should not have failed")
		}
		
		expected := map[string]int{};
		if !reflect.DeepEqual(severityCounts, expected) {
			t.Errorf("AnalyzeLogs failed with invalid logs, expected %v got %v", expected, severityCounts);
		}
	})
}
```
*   Tests the `AnalyzeLogs` function with different log files and filters, verifying if the result is as expected.
* Includes a scenario to test for invalid logs, and verifies that the results are correct.
* Uses `os.WriteFile` to write some test log files.

**4. `main.go` (Main Application Logic):**

```go
// main.go
package main

import (
	"fmt"
    "log"
    "os"
	"loganalyzer/analyzer"
)


func main() {
    if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run main.go <filename> [filter]");
		os.Exit(1);
    }
    
    filename := os.Args[1];
	var filter string;
	if len(os.Args) == 3 {
		filter = os.Args[2];
	}
	
    severityCounts, err := analyzer.AnalyzeLogs(filename, filter);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing logs %v\n", err);
		os.Exit(1);
	}
	
	fmt.Println("Log Analysis Results");
	for severity, count := range severityCounts {
		fmt.Printf("%s: %d\n", severity, count);
	}
}
```

*  Imports the `analyzer` package, gets the command line arguments and calls the `AnalyzeLogs` function.
*  It will output to standard error if an error happens while analyzing the logs, and print to standard output all the analysis results.

**How to Run the Project and Tests:**

1.  **Run the Application:**
    Open a terminal, navigate to the `loganalyzer` directory, and run:

    ```bash
    go run . app.log
    ```
	Create a file named `app.log` with some test log information.
	You can also use a filter as a second argument:
    ```bash
	go run . app.log error
	```

2.  **Run the Tests:**
    Open a terminal, navigate to the `loganalyzer` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok      loganalyzer/analyzer        0.003s
```

**Output (If Tests Fail):**
If any test fails, the output will give additional information to help with the debugging process.
```
--- FAIL: TestAnalyzeLogs/Test_with_filter (0.00s)
    analyzer_test.go:44: AnalyzeLogs failed with filter, expected: map[error:2 debug:1] , got: map[error:2]
FAIL
exit status 1
FAIL	loganalyzer/analyzer	0.003s
```

**Key Features of This Project:**

*   **Log Analysis Tool:** Implements the core functionality of a log analysis tool that is capable of filtering and grouping the log lines.
*   **Modularity:** Encapsulates the analysis logic in a separate package called `analyzer`.
*   **Testing:** Includes unit tests for all key components of the `analyzeLogs` function.
*	**Filtering**: The tool has the ability to filter logs based on a regex.
*	**Error Handling**: The tool has a basic error handling for errors happening during the log parsing.

This project demonstrates how to create a tool to analyze log files, and also shows how to create different test scenarios and handle errors.
