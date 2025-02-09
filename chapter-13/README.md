
## Project: Dynamic Struct Printer with Reflection

This project will implement a simple data printer that can dynamically print the fields of any struct, including nested structs, using reflection. This will showcase how to inspect struct types and values at runtime. This project will create a command-line tool that:

1.  Accepts a struct as input (defined in the code).
2.  Uses reflection to iterate through the fields of the struct, including nested structs.
3.  Prints the name, type, and value of each field.

**Project Structure:**

```
structprinter/
├── go.mod
├── main.go
└── printer/
    ├── printer.go
    └── printer_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir structprinter
    cd structprinter
    ```

2.  **Initialize the Go Module:**
    Inside the `structprinter` directory, run:

    ```bash
    go mod init structprinter
    ```

3.  **Create the `printer` Package Directory:**

    ```bash
    mkdir printer
    ```

4.  **Create the `printer.go` file (inside the `printer` directory):**

    ```bash
    touch printer/printer.go
    ```

5.  **Create the `printer_test.go` file (inside the `printer` directory):**

    ```bash
    touch printer/printer_test.go
    ```

6.  **Create the `main.go` file (inside the `structprinter` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
structprinter/
├── go.mod
├── main.go
└── printer/
    ├── printer.go
    └── printer_test.go
```

**Now, paste the following code into the corresponding files:**

**1. `go.mod` File:**

```
module structprinter

go 1.21
```

**2. `printer/printer.go` (Struct Printer Logic):**

```go
// printer/printer.go
package printer

import (
	"fmt"
	"reflect"
)

func PrintStructFields(data interface{}) {
    val := reflect.ValueOf(data)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    typ := val.Type()
    if typ.Kind() != reflect.Struct {
        fmt.Println("Provided input is not a struct")
        return
    }
    for i := 0; i < val.NumField(); i++ {
        field := typ.Field(i);
        value := val.Field(i);
        fmt.Printf("Field: %s, Type: %v, Value: %v\n", field.Name, field.Type, value)
        if value.Kind() == reflect.Struct {
		    printNestedFields(value, "\t")
        }
    }
}

func printNestedFields(val reflect.Value, prefix string) {
    typ := val.Type()
    for i := 0; i < val.NumField(); i++ {
        field := typ.Field(i);
        value := val.Field(i);
        fmt.Printf("%sNested Field: %s, Nested Type: %v, Nested Value: %v\n", prefix, field.Name, field.Type, value)
        if value.Kind() == reflect.Struct {
		   printNestedFields(value, prefix + "\t");
        }
    }
}
```
*   This file contains the core logic for printing the fields of a struct using reflection.
*   The function `PrintStructFields` uses reflection to iterate through a struct's fields and prints their name, type, and value.
* The function `printNestedFields` has a recursive logic to print the fields of nested structs.

**3. `printer/printer_test.go` (Struct Printer Logic Tests):**

```go
// printer/printer_test.go
package printer

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
)


type Address struct {
	Street string
	City string
}

type Person struct {
	Name string
	Age int
	Address Address
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

func TestPrintStructFields(t *testing.T) {
	address := Address{
		Street: "123 Main Street",
		City: "Anytown",
	}
	person := Person{
		Name: "Alice",
		Age: 30,
		Address: address,
	}

    expectedOutput := `Field: Name, Type: string, Value: Alice
Field: Age, Type: int, Value: 30
Field: Address, Type: printer.Address, Value: {123 Main Street Anytown}
	Nested Field: Street, Nested Type: string, Nested Value: 123 Main Street
	Nested Field: City, Nested Type: string, Nested Value: Anytown
`
    actualOutput := captureOutput(func() {
		PrintStructFields(person)
	});

	if actualOutput != expectedOutput {
		t.Errorf("PrintStructFields failed, expected: \n%q, got: \n%q", expectedOutput, actualOutput)
	}

	t.Run("Test with invalid parameter", func(t *testing.T) {
		actualOutput := captureOutput(func() {
			PrintStructFields(10);
		})

		expectedOutput := "Provided input is not a struct\n";

		if actualOutput != expectedOutput {
			t.Errorf("PrintStructFields failed with a non struct, expected %q, got %q", expectedOutput, actualOutput);
		}
	})

}
```
*   Includes different scenarios for testing, a nested struct and also a scenario where a non struct variable is passed to the `PrintStructFields` method.
*   Asserts if the output is correct.

**4. `main.go` (Main Application Logic):**

```go
// main.go
package main

import (
	"fmt"
	"structprinter/printer"
)
type Address struct {
	Street string
	City string
}

type Person struct {
	Name string
	Age int
	Address Address
}

func main() {
    person := Person{
        Name: "Alice",
        Age: 30,
		Address: Address {
			Street: "Main Street 123",
			City: "New York",
		},
    }
	fmt.Println("Printing the fields of a struct using reflection");
    printer.PrintStructFields(person);
}
```
*   Defines the structs `Address` and `Person`.
*  Creates a `Person` variable and prints its fields using the `printer.PrintStructFields` method.

**How to Run the Project and Tests:**

1.  **Run the Application:**
    Open a terminal, navigate to the `structprinter` directory, and run:

    ```bash
    go run .
    ```

2.  **Run the Tests:**
    Open a terminal, navigate to the `structprinter` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok      structprinter/printer       0.003s
```

**Output (If Tests Fail):**
If any of the tests fail, you'll be given detailed information on the failure.

```
--- FAIL: TestPrintStructFields (0.00s)
    printer_test.go:37: PrintStructFields failed, expected:
"Field: Name, Type: string, Value: Alice
Field: Age, Type: int, Value: 30
Field: Address, Type: printer.Address, Value: {Main Street 123 New York}
Nested Field: Street, Nested Type: string, Nested Value: Main Street 123
Nested Field: City, Nested Type: string, Nested Value: New York
", got:
"Field: Name, Type: string, Value: Alice
Field: Age, Type: int, Value: 30
Field: Address, Type: structprinter/printer.Address, Value: {Main Street 123 New York}
	Nested Field: Street, Nested Type: string, Nested Value: Main Street 123
	Nested Field: City, Nested Type: string, Nested Value: New York
"
FAIL
exit status 1
FAIL    structprinter/printer       0.003s
```

**Key Features of This Project:**

*   **Dynamic Struct Printing:** Uses reflection to dynamically inspect and print struct fields.
*   **Modularity:** Encapsulates the reflection logic in a separate `printer` package.
*   **Testing:** Includes unit tests that verify the correctness of the struct printing logic.
*   **Nested Structs:** Implements a method that also prints the fields of nested structs.

This project provides a good example of how to use reflection to achieve dynamic behavior in Go, while also following good practices for modularity and testing.

