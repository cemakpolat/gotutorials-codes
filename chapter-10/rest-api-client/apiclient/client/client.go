package client

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

func FetchAndPrintData(url string) (*Todo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non 200 status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var todo Todo
	if err := json.Unmarshal(data, &todo); err != nil {
		return nil, fmt.Errorf("error unmarshalling json %w", err)
	}

	fmt.Printf("Todo ID: %d\n", todo.ID)
	fmt.Printf("Title: %s\n", todo.Title)
	fmt.Printf("Completed: %t\n", todo.Completed)
	return &todo, nil
}
