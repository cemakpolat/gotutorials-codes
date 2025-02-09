package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Post struct {
	Message string `json:"message"`
	Author  string `json:"author"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Home page!")
}
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is the about page!")
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Query parameters:", r.URL.Query())
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "This is a list of posts!")
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading the body", http.StatusBadRequest)
			return // Important to return here!
		}
		var post Post
		err = json.Unmarshal(body, &post)
		if err != nil {
			http.Error(w, "Error unmarshalling json", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Post created:\n")
		fmt.Fprintf(w, "Message: %s\nAuthor: %s\n", post.Message, post.Author)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
