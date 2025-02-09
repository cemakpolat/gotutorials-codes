
**Complete Example: `greeting` and `calculator` packages with Tests (Including Setup Instructions)**

Here's the complete example, including commands to create the folder structure and files:

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run the following command to create the main project directory:

    ```bash
    mkdir xproject
    cd xproject
    ```

2.  **Initialize the Go Module:**
    Inside the `xproject` directory, run:

    ```bash
    go mod init xproject
    ```

    This creates the `go.mod` file for managing dependencies.

3.  **Create the `greeting` Package Directory:**

    ```bash
    mkdir greeting
    ```

4.  **Create the `calculator` Package Directory:**

    ```bash
    mkdir calculator
    ```

5.  **Create the `greeting.go` file (inside `greeting` directory):**

    ```bash
    touch greeting/greeting.go
    ```

6.  **Create the `doc.go` file (inside `greeting` directory):**

    ```bash
    touch greeting/doc.go
    ```

7.  **Create the `greeting_test.go` file (inside `greeting` directory):**

    ```bash
     touch greeting/greeting_test.go
    ```

8.  **Create the `calculator.go` file (inside `calculator` directory):**

    ```bash
    touch calculator/calculator.go
    ```

9. **Create the `doc.go` file (inside `calculator` directory):**

    ```bash
    touch calculator/doc.go
    ```

10. **Create the `calculator_test.go` file (inside `calculator` directory):**

    ```bash
    touch calculator/calculator_test.go
    ```

11. **Create the `main.go` file (inside `xproject` directory):**

    ```bash
    touch main.go
    ```

Now, you have the following folder structure:

```
xproject/
├── go.mod
├── main.go
├── greeting/
│   ├── greeting.go
│   ├── greeting_test.go
│   └── doc.go
└── calculator/
    ├── calculator.go
    ├── calculator_test.go
    └── doc.go
```

**Now, you can paste the code in their respective files:**

**1. `go.mod` File:**

```
module xproject

go 1.21
```

**2. `main.go` (Main Application):**

```go
// main.go
package main

import (
    "fmt"
    "xproject/greeting"
    "xproject/calculator"
)

func main() {
    greeting.SayHello("Alice")
    sum := calculator.Add(10, 5)
    fmt.Println("Sum:", sum)
    product := calculator.Multiply(4,3)
    fmt.Println("Product:", product)
}
```

**3. `greeting/greeting.go` (Greeting Package):**

```go
// greeting/greeting.go
package greeting

import (
	"fmt"
	"os"
	"io"
)

// SayHello prints a greeting message to the console
func SayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

func captureOutput(f func()) string {
    old := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    defer func() {
        os.Stdout = old
    }()
    f()
    w.Close()
    out, _ := io.ReadAll(r)
    return string(out)
}
```

**4. `greeting/doc.go` (Greeting Package Documentation):**

```go
// Package greeting provides functions for generating greeting messages.
package greeting
```

**5.  `greeting/greeting_test.go` (Greeting Package Tests):**

    ```go
    // greeting/greeting_test.go
    package greeting

    import (
        "testing"
	"os"
        "io"
    )

    func TestSayHello(t *testing.T) {
        testName := "Test User"
        expected := "Hello, Test User!\n"
        
        // This will capture any output that is printed in the console
        output := captureOutput(func() {
            SayHello(testName)
        })
	
	if output != expected {
		t.Errorf("SayHello() failed, expected %q, got %q", expected, output)
	}
    }
    ```

**6. `calculator/calculator.go` (Calculator Package):**

```go
// calculator/calculator.go
package calculator

// Add performs addition of two numbers
func Add(a int, b int) int {
    return a + b
}

// multiply performs multiplication of two numbers
func Multiply(a int, b int) int {
	return a * b;
}
```

**7. `calculator/doc.go` (Calculator Package Documentation):**

```go
// Package calculator provides basic mathematical operations.
package calculator
```

**8.  `calculator/calculator_test.go` (Calculator Package Tests):**

    ```go
    // calculator/calculator_test.go
    package calculator

    import "testing"

    func TestAdd(t *testing.T) {
	testCases := []struct{
		a int
		b int
		expected int
		}{
			{10, 5, 15},
			{0, 0, 0},
			{-5, 5, 0},
			{-10, -5, -15},
		}
	for _, tc := range testCases {
		if actual := Add(tc.a, tc.b); actual != tc.expected {
			t.Errorf("Add(%d,%d) failed, expected %d, got %d", tc.a, tc.b, tc.expected, actual)
		}
	}
    }

    func TestMultiply(t *testing.T) {
	testCases := []struct{
		a int
		b int
		expected int
		}{
			{10, 5, 50},
			{0, 0, 0},
			{-5, 5, -25},
			{-10, -5, 50},
		}
	for _, tc := range testCases {
		if actual := Multiply(tc.a, tc.b); actual != tc.expected {
			t.Errorf("Multiply(%d,%d) failed, expected %d, got %d", tc.a, tc.b, tc.expected, actual)
		}
	}
    }
    ```

**How to Run the Code and Tests:**

1.  **Run the Application:**
    Inside the `xproject` directory:

    ```bash
    go run .
    ```

2.  **Run the Tests:**
    Inside the `xproject` directory:

    ```bash
    go test ./...
    ```

**Explanation:**

*   The setup instructions are now clearly defined before the actual code.
*   We are using `touch` command in the setup instructions to create the files, but any text editor or tool will suffice.
*   Users can easily copy these instructions, which will allow them to fully set up the project.

This comprehensive example should provide a complete and self-contained guide for the readers to understand and implement the concepts we have been discussing.

Let me know if you have any other questions!
