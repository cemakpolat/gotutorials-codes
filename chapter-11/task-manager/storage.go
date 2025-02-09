package main

import (
	"encoding/json"
	"errors"
	"os"
)

// LoadTasks loads tasks from a JSON file
func LoadTasks(filename string) (TaskList, error) {
	var tasks TaskList
	file, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		return tasks, nil // Return an empty task list if the file doesn't exist
	}
	if err != nil {
		return tasks, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		return tasks, err
	}
	return tasks, nil
}

// SaveTasks saves tasks to a JSON file
func (t *TaskList) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(t)
}
