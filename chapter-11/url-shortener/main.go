package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize in-memory storage
	store := NewMemoryStore()

	// Define HTTP routes
	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		HandleShorten(w, r, store)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleRedirect(w, r, store)
	})

	// Start the server
	log.Println("Starting URL Shortener on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
