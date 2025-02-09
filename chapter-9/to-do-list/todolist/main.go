package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todolist/todo"
)

func main() {
	tasks, err := todo.LoadTasksFromFile("tasks.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading tasks from file: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command (add/list/complete/exit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.SplitN(input, " ", 2)
		command := parts[0]

		switch command {
		case "add":
			if len(parts) < 2 {
				fmt.Println("Description required!")
				continue
			}
			todo.AddTask(&tasks, parts[1])
			err = todo.SaveTasksToFile(tasks, "tasks.txt")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving tasks to file: %v\n", err)
			}
		case "list":
			todo.ListTasks(tasks)
		case "complete":
			if len(parts) < 2 {
				fmt.Println("Task ID required!")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid task id!")
				continue
			}
			todo.CompleteTask(&tasks, id)
			err = todo.SaveTasksToFile(tasks, "tasks.txt")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving tasks to file: %v\n", err)
			}
		case "exit":
			fmt.Println("Exiting the application")
			return
		default:
			fmt.Println("Invalid command. Please enter add/list/complete or exit")
		}
	}
}
