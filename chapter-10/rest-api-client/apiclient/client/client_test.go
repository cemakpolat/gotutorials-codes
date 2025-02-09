package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchAndPrintDataSuccess(t *testing.T) {
	// Create a mock server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"userId": 1, "id": 1, "title": "Test Todo", "completed": false}`)
	}))
	defer testServer.Close()

	// Call the function with the mock server URL
	todo, err := FetchAndPrintData(testServer.URL)
	if err != nil {
		t.Errorf("FetchAndPrintData() failed with error: %v", err)
		return
	}

	if todo.Title != "Test Todo" {
		t.Fatalf("Returned incorrect todo item, expected title Test Todo, got %s", todo.Title)
	}
}

func TestFetchAndPrintDataFailure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	_, err := FetchAndPrintData(testServer.URL)
	if err == nil {
		t.Fatalf("FetchAndPrintData should have failed with 500")
	}
}
