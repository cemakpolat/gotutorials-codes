package printer

import (
	"io"
	"os"
	"testing"
)

type Address struct {
	Street string
	City   string
}

type Person struct {
	Name    string
	Age     int
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
		City:   "Anytown",
	}
	person := Person{
		Name:    "Alice",
		Age:     30,
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
	})

	if actualOutput != expectedOutput {
		t.Errorf("PrintStructFields failed, expected: \n%q, got: \n%q", expectedOutput, actualOutput)
	}

	t.Run("Test with invalid parameter", func(t *testing.T) {
		actualOutput := captureOutput(func() {
			PrintStructFields(10)
		})

		expectedOutput := "Provided input is not a struct\n"

		if actualOutput != expectedOutput {
			t.Errorf("PrintStructFields failed with a non struct, expected %q, got %q", expectedOutput, actualOutput)
		}
	})

}
