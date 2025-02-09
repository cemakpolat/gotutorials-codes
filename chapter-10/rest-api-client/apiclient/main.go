package main

import (
	"apiclient/client"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <url>")
		os.Exit(1)
	}
	url := os.Args[1]

	_, err := client.FetchAndPrintData(url)
	if err != nil {
		log.Fatal("error fetching data", err)
	}
}
