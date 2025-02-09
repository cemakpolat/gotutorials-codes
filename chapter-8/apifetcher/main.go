package main

import (
	"apifetcher/utils"
	"fmt"
	"log"
	"net/http"
)

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://jsonplaceholder.typicode.com/todos/1"

	todo, err := utils.FetchTodo(url)
	if err != nil {
		http.Error(w, "Failed to fetch the data", http.StatusInternalServerError)
		log.Println("Error fetching data:", err)
		return
	}

	fmt.Fprintf(w, "Todo ID: %d\n", todo.ID)
	fmt.Fprintf(w, "Title: %s\n", todo.Title)
	fmt.Fprintf(w, "Completed: %t", todo.Completed)

}

func main() {
	http.HandleFunc("/fetch", fetchHandler)
	fmt.Println("Server is starting at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
