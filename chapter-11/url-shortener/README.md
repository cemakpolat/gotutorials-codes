Let’s move to the next project: **A URL Shortener**. This tool allows users to input long URLs and get a shortened version that redirects to the original link when accessed.  

---

### **URL Shortener Overview**  

#### **Basic Features (MVP)**:  
1. Input a long URL and generate a short URL.  
2. Store the mapping of short URLs to long URLs in memory.  
3. Redirect users to the original URL when they access the short URL.  

#### **Enhancements (Future)**:  
1. Store URL mappings in a database (e.g., SQLite or Redis).  
2. Allow custom aliases for shortened URLs.  
3. Add expiration times for URLs.  
4. Build a simple web interface for the tool.  

---

### **Step 1: Set Up the Project**  

#### 1. **Create the Directory**  
```bash
mkdir url-shortener
cd url-shortener
go mod init url-shortener
```

#### 2. **Directory Structure**  
```
url-shortener/
├── main.go        // Entry point
├── handler.go     // HTTP handlers for shortening and redirecting
├── storage.go     // Logic for storing and retrieving URLs
└── README.md      // Documentation
```

---

### **Step 2: Code Implementation**

#### **1. `main.go` (Entry Point)**  
This file sets up the HTTP server.  

```go
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
```

---

#### **2. `handler.go` (HTTP Handlers)**  
This file contains the logic for shortening and redirecting URLs.  

```go
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
```

---

#### **3. `storage.go` (In-Memory Storage)**  
This file provides a simple way to store URL mappings.  

```go
package main

import "sync"

// MemoryStore is an in-memory store for URL mappings
type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewMemoryStore initializes a new MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]string),
	}
}

// Save stores a short URL and its corresponding long URL
func (s *MemoryStore) Save(shortCode, longURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[shortCode] = longURL
}

// Get retrieves the original URL for a given short code
func (s *MemoryStore) Get(shortCode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	longURL, found := s.store[shortCode]
	return longURL, found
}
```

---

### **Step 3: Test the Tool**  

#### **1. Start the Server**  
```bash
go run .
```

#### **2. Shorten a URL**  
Use a tool like `curl` or Postman to send a POST request:  

```bash
curl -X POST -H "Content-Type: application/json" -d '{"url": "https://example.com"}' http://localhost:8080/shorten
```

You should receive a response like:  
```json
{"short_url":"http://localhost:8080/abc123"}
```

#### **3. Redirect to the Original URL**  
Access the short URL in your browser or using `curl`:  

```bash
curl -v http://localhost:8080/abc123
```

This should redirect you to `https://example.com`.  
