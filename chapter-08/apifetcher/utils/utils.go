package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func FetchTodo(url string) (*Todo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making http request %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the body %w", err)
	}

	var todo Todo
	if err := json.Unmarshal(data, &todo); err != nil {
		return nil, fmt.Errorf("error unmarshalling json %w", err)
	}

	return &todo, nil
}
