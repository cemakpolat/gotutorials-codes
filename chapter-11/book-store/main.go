package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:"id"`
	BookID string `json:"bookId"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// Struct for the json response of the Google books API
type GoogleBooksResponse struct {
	Items []struct {
		Id         string `json:"id"`
		VolumeInfo struct {
			Title   string   `json:"title"`
			Authors []string `json:"authors"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

// Struct for the json response of the open library API
type OpenLibraryResponse struct {
	Docs []struct {
		Key        string   `json:"key"`
		Title      string   `json:"title"`
		AuthorName []string `json:"author_name"`
	} `json:"docs"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest) // Respond with error to the client
		return
	}
	for _, book := range books {
		if book.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.NotFound(w, r)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = len(books) + 1
	books = append(books, book)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	for index, book := range books {
		if book.ID == id {
			books = append(books[:index], books[index+1:]...)
			w.WriteHeader(http.StatusNoContent) // 204 - Successful deletion
			return
		}
	}
	http.NotFound(w, r)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var updatedBook Book
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for index, book := range books {
		if book.ID == id {
			updatedBook.ID = id
			books[index] = updatedBook
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}
	http.NotFound(w, r)
}

func fetchBooksFromGoogleBooks(bookCount int) {
	// Use the search endpoint with a query for book titles
	url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=lord+of+the+rings&maxResults=%d", bookCount)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching data from google books: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Decode the JSON response
	var googleBooksResp GoogleBooksResponse
	err = json.Unmarshal(body, &googleBooksResp)
	if err != nil {
		log.Fatalf("Error unmarshalling json from google books: %v", err)
	}

	// Loop over each returned book and create a new book in our storage
	for _, item := range googleBooksResp.Items {
		author := "Unknown Author" // Default if no author
		if item.VolumeInfo.Authors != nil && len(item.VolumeInfo.Authors) > 0 {
			author = item.VolumeInfo.Authors[0]
		}

		books = append(books, Book{
			ID:     len(books) + 1,
			BookID: item.Id, // Using google id as id
			Title:  item.VolumeInfo.Title,
			Author: author,
		})
	}
}

func fetchBooksFromOpenLibrary(bookCount int) {
	url := fmt.Sprintf("https://openlibrary.org/search.json?q=the+lord+of+the+rings&limit=%d", bookCount)
	// Use the search endpoint with a query for book titles
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching data from open library: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Decode the JSON response
	var openLibraryResp OpenLibraryResponse
	err = json.Unmarshal(body, &openLibraryResp)
	if err != nil {
		log.Fatalf("Error unmarshalling json from open library: %v", err)
	}

	// Loop over each returned book and create a new book in our storage
	for index, doc := range openLibraryResp.Docs {
		author := "Unknown Author"
		if doc.AuthorName != nil && len(doc.AuthorName) > 0 {
			author = doc.AuthorName[0]
		}
		books = append(books, Book{
			ID:     len(books) + 1,
			BookID: strconv.Itoa(index + 1), // Using index as id, you should use the proper id
			Title:  doc.Title,
			Author: author,
		})
	}
}

func main() {
	// Initialize book data
	bookCount := 20
	fetchBooksFromGoogleBooks(bookCount)
	fetchBooksFromOpenLibrary(bookCount)

	r := mux.NewRouter()
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", removeBook).Methods("DELETE")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
