package main

import "fmt"

// Task represents a single to-do item
type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// TaskList manages a list of tasks
type TaskList struct {
	Tasks []Task `json:"tasks"`
}

// Add adds a new task to the list
func (t *TaskList) Add(text string) {
	id := len(t.Tasks)
	task := Task{ID: id, Text: text}
	t.Tasks = append(t.Tasks, task)
}

// List displays all tasks
func (t *TaskList) List() {
	if len(t.Tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for _, task := range t.Tasks {
		fmt.Printf("[%d] %s\n", task.ID, task.Text)
	}
}

// Delete removes a task by ID
func (t *TaskList) Delete(id int) error {
	if id < 0 || id >= len(t.Tasks) {
		return fmt.Errorf("invalid task ID: %d", id)
	}

	t.Tasks = append(t.Tasks[:id], t.Tasks[id+1:]...)
	return nil
}
