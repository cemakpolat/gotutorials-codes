package main

import (
	"downloader/download"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <url> <filepath>")
		os.Exit(1)
	}

	url := os.Args[1]
	filepath := os.Args[2]

	err := download.DownloadFile(url, filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("File downloaded successfully!")
}
