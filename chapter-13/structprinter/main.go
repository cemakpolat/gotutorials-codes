package main

import (
	"fmt"
	"structprinter/printer"
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

func main() {
	person := Person{
		Name: "Alice",
		Age:  30,
		Address: Address{
			Street: "Main Street 123",
			City:   "New York",
		},
	}
	fmt.Println("Printing the fields of a struct using reflection")
	printer.PrintStructFields(person)
}
