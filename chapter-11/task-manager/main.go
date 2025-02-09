package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	// Define flags
	add := flag.String("add", "", "Add a new task")
	list := flag.Bool("list", false, "List all tasks")
	delete := flag.Int("delete", -1, "Delete a task by ID")
	flag.Parse()

	// Load tasks from storage
	tasks, err := LoadTasks("tasks.json")
	if err != nil {
		log.Fatalf("Error loading tasks: %v", err)
	}

	// Handle commands
	switch {
	case *add != "":
		tasks.Add(*add)
		fmt.Println("Task added!")
	case *list:
		tasks.List()
	case *delete >= 0:
		if err := tasks.Delete(*delete); err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Task deleted!")
		}
	default:
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	// Save tasks back to storage
	if err := tasks.Save("tasks.json"); err != nil {
		log.Fatalf("Error saving tasks: %v", err)
	}
}
