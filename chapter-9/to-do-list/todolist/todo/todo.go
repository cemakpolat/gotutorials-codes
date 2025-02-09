package todo

import (
	"bufio"
	"fmt"
	"os"
)

type Task struct {
	ID          int
	Description string
	Completed   bool
}

func AddTask(tasks *[]Task, description string) {
	newTask := Task{
		ID:          len(*tasks) + 1,
		Description: description,
		Completed:   false,
	}
	*tasks = append(*tasks, newTask)
	fmt.Println("Task added!")
}
func ListTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No task found")
		return
	}
	for _, task := range tasks {
		status := "[ ]"
		if task.Completed {
			status = "[x]"
		}
		fmt.Printf("%d %s %s\n", task.ID, status, task.Description)
	}
}

func CompleteTask(tasks *[]Task, id int) {
	for i, task := range *tasks {
		if task.ID == id {
			(*tasks)[i].Completed = true
			fmt.Println("Task marked as complete")
			return
		}
	}
}

func LoadTasksFromFile(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var tasks []Task
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var task Task
		var completedStr string
		_, err := fmt.Sscan(line, &task.ID)
		if err != nil {
			continue
		}

		parts := []rune(line)
		index := 0
		for i := 0; i < len(parts); i++ {
			if parts[i] == ' ' {
				index = i + 1
				break
			}
		}

		indexCompleted := 0
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] == ' ' {
				indexCompleted = i + 1
				break
			}
		}

		if index < indexCompleted {
			task.Description = string(parts[index : indexCompleted-1])
		}

		completedStr = string(parts[indexCompleted:])

		if completedStr == "true" {
			task.Completed = true
		} else {
			task.Completed = false
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}
func SaveTasksToFile(tasks []Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, task := range tasks {
		_, err = fmt.Fprintf(file, "%d %s %t\n", task.ID, task.Description, task.Completed)
		if err != nil {
			return err
		}
	}
	return nil
}
