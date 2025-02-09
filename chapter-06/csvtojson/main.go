package main

import (
	"bufio"
	"csvtojson/converter"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the path to the CSV file: ")
	csvPath, _ := reader.ReadString('\n')
	csvPath = csvPath[:len(csvPath)-1] // remove the last character, the newline character

	fmt.Print("Enter the path to the output JSON file: ")
	jsonPath, _ := reader.ReadString('\n')
	jsonPath = jsonPath[:len(jsonPath)-1] // remove the last character, the newline character

	people, err := converter.ProcessCSV(csvPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing CSV: %v\n", err)
		return
	}

	err = converter.WriteJSON(people, jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSON: %v\n", err)
		return
	}

	fmt.Println("CSV data converted to JSON successfully.")
}
