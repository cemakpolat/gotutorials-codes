package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Parse command-line arguments
	sourceDir := flag.String("source", ".", "Directory to organize (default is current directory)")
	flag.Parse()

	// Check if the source directory exists
	if _, err := os.Stat(*sourceDir); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", *sourceDir)
	}

	fmt.Printf("Organizing files in: %s\n", *sourceDir)

	// Call the organizer logic
	if err := OrganizeFiles(*sourceDir); err != nil {
		log.Fatalf("Error organizing files: %v", err)
	}

	fmt.Println("Files organized successfully!")
}
