package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	recorder := httptest.NewRecorder()
	HomeHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("HomeHandler failed, expected %d, got %d", http.StatusOK, recorder.Code)
	}

	expected := "Welcome to the Home page!"
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("HomeHandler failed, expected %q, got %q", expected, actual)
	}
}

func TestAboutHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/about", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	recorder := httptest.NewRecorder()
	AboutHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("AboutHandler failed, expected %d, got %d", http.StatusOK, recorder.Code)
	}

	expected := "This is the about page!"
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("AboutHandler failed, expected %q, got %q", expected, actual)
	}
}

func TestPostsHandler(t *testing.T) {
	t.Run("Test GET request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/posts", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		recorder := httptest.NewRecorder()
		PostsHandler(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("PostsHandler failed, expected %d, got %d", http.StatusOK, recorder.Code)
		}

		expected := "This is a list of posts!"
		actual := recorder.Body.String()
		if actual != expected {
			t.Errorf("PostsHandler GET failed, expected %q, got %q", expected, actual)
		}
	})

	t.Run("Test POST request", func(t *testing.T) {
		body := `{"message": "Test message", "author": "test"}`
		req, err := http.NewRequest(http.MethodPost, "/posts", strings.NewReader(body))
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		PostsHandler(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("PostsHandler failed, expected %d, got %d", http.StatusOK, recorder.Code)
		}

		actual := recorder.Body.String()

		expected := "Post created:\nMessage: Test message\nAuthor: test\n"
		if actual != expected {
			t.Errorf("PostsHandler POST failed, expected %q, got %q", expected, actual)
		}
	})
}
