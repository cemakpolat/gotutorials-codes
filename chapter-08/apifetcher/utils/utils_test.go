package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchTodo(t *testing.T) {
	// Create a mock server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"userId": 1, "id": 1, "title": "Test Todo", "completed": false}`))
	}))
	defer testServer.Close()

	// Call the function with the mock server URL
	todo, err := FetchTodo(testServer.URL)
	if err != nil {
		t.Errorf("FetchTodo() failed with error: %v", err)
		return
	}

	// Assert the expected value
	if todo.Title != "Test Todo" {
		t.Errorf("FetchTodo() returned incorrect title: got %v, want %v", todo.Title, "Test Todo")
	}

	if todo.ID != 1 {
		t.Errorf("FetchTodo() returned incorrect id: got %d, want %d", todo.ID, 1)
	}
}

func TestFetchTodoFailure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	_, err := FetchTodo(testServer.URL)
	if err == nil {
		t.Errorf("FetchTodo() should have failed with 500 error")
	}
}
