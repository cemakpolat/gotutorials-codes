package main

import (
	"fmt"
	"xproject/calculator"
	"xproject/greeting"
)

func main() {
	greeting.SayHello("Alice")
	sum := calculator.Add(10, 5)
	fmt.Println("Sum:", sum)
	product := calculator.Multiply(4, 3)
	fmt.Println("Product:", product)
}
