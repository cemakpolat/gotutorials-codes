package main

import "fmt"

func add(a, b int) int {
    return a + b
}

func subtract(a, b int) int {
    return a - b
}

func main() {
    var operation func(int, int) int  // Declare a variable of function type
    operation = add // Assign add function to operation
    fmt.Println("Addition:", operation(5, 3)) // Output: Addition: 8

operation = subtract // Assign subtract function to operation
    fmt.Println("Subtraction:", operation(5, 3)) // Output: Subtraction: 2
}

