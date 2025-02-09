package main

import (
	"fmt"
	"log"
	"net/http"
	"webserver/server"
)

func main() {
	http.HandleFunc("/", server.HomeHandler)
	http.HandleFunc("/about", server.AboutHandler)
	http.HandleFunc("/posts", server.PostsHandler)
	fmt.Println("Server is starting at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
