package main

import (
	"fmt"
	"loganalyzer/analyzer"
	"os"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run main.go <filename> [filter]")
		os.Exit(1)
	}

	filename := os.Args[1]
	var filter string
	if len(os.Args) == 3 {
		filter = os.Args[2]
	}

	severityCounts, err := analyzer.AnalyzeLogs(filename, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing logs %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Log Analysis Results")
	for severity, count := range severityCounts {
		fmt.Printf("%s: %d\n", severity, count)
	}
}
