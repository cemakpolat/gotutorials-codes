package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// HandleShorten handles requests to shorten URLs
func HandleShorten(w http.ResponseWriter, r *http.Request, store *MemoryStore) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the input JSON
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Generate a short code and store it
	shortCode := generateShortCode()
	store.Save(shortCode, req.URL)

	// Respond with the short URL
	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"short_url": shortURL})
}

// HandleRedirect handles requests to redirect to the original URL
func HandleRedirect(w http.ResponseWriter, r *http.Request, store *MemoryStore) {
	shortCode := r.URL.Path[1:] // Get the short code from the path
	originalURL, found := store.Get(shortCode)
	if !found {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

// generateShortCode creates a random 6-character string
func generateShortCode() string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}
