package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	out, err := os.Create("./uploads/" + header.Filename)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "File uploaded successfully: %s", header.Filename)

}
func main() {
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		log.Fatal("Unable to create uploads directory")
	}
	http.HandleFunc("/upload", uploadFile)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
