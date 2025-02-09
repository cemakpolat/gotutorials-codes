package main

import (
	"fmt"
	"os"
	"strconv"
	"wordcounter/counter"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <filename> <topN>")
		os.Exit(1)
	}
	filename := os.Args[1]
	topN, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid topN value!")
		os.Exit(1)
	}

	wordCounts, err := counter.CountWords(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting words: %v\n", err)
		os.Exit(1)
	}

	counter.PrintTopWords(wordCounts, topN)
}
