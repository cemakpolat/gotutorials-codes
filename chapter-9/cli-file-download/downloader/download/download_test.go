package download

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDownloadFileSuccess(t *testing.T) {
	// Create a mock server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Test file content")
	}))
	defer testServer.Close()

	// Call DownloadFile with the mock server URL
	filepath := "testfile.txt"
	err := DownloadFile(testServer.URL, filepath)
	if err != nil {
		t.Fatalf("DownloadFile() failed with error: %v", err)
	}

	defer os.Remove(filepath)
	// read the file content
	content, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Error reading test file %v", err)
	}

	if string(content) != "Test file content\n" {
		t.Fatalf("Content is incorrect: %v", string(content))
	}
}

func TestDownloadFileFailure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	filepath := "testfile.txt"
	err := DownloadFile(testServer.URL, filepath)
	if err == nil {
		t.Errorf("DownloadFile should have failed with 500 error")
	}

	_, err = os.Stat(filepath)
	if !os.IsNotExist(err) {
		t.Fatalf("File should not exist")
	}
}
