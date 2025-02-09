Implement the first example from Chapter 9, the "Simple To-Do List Manager," as a complete project with tests. We'll follow the structure we've been using and will include all the necessary setup instructions.

**Project: Simple To-Do List Manager with Tests**

This project will implement a command-line to-do list manager with features for adding, listing, marking tasks as complete, and storing the data persistently in a file. It will include unit tests to ensure its functionality.

**Project Structure:**

```
todolist/
├── go.mod
├── main.go
└── todo/
    ├── todo.go
    └── todo_test.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir todolist
    cd todolist
    ```

2.  **Initialize the Go Module:**
    Inside the `todolist` directory, run:

    ```bash
    go mod init todolist
    ```

3.  **Create the `todo` Package Directory:**

    ```bash
    mkdir todo
    ```

4.  **Create the `todo.go` file (inside the `todo` directory):**

    ```bash
    touch todo/todo.go
    ```

5.  **Create the `todo_test.go` file (inside the `todo` directory):**

    ```bash
    touch todo/todo_test.go
    ```

6.  **Create the `main.go` file (inside the `todolist` directory):**

    ```bash
    touch main.go
    ```

Now, you should have the following project structure:

```
todolist/
├── go.mod
├── main.go
└── todo/
    ├── todo.go
    └── todo_test.go
```

**Now, copy the following code to their corresponding files:**

**1. `go.mod` File:**

```
module todolist

go 1.21
```

**2. `todo/todo.go` (To-Do Logic):**

```go
// todo/todo.go
package todo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	ID int
	Description string
	Completed bool
}

func AddTask(tasks *[]Task, description string) {
	newTask := Task {
		ID: len(*tasks) + 1,
		Description: description,
		Completed: false,
	}
	*tasks = append(*tasks, newTask);
	fmt.Println("Task added!");
}

func ListTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
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
	fmt.Println("Task not found")
}

func LoadTasksFromFile(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil;
		}
		return nil, err;
	}

	defer file.Close();
	var tasks []Task;

	scanner := bufio.NewScanner(file);
	for scanner.Scan() {
		line := scanner.Text();
		var task Task;
		_, err = fmt.Sscan(line, &task.ID, &task.Description, &task.Completed);
		if err != nil {
			// skip invalid lines
			continue
		}
		tasks = append(tasks, task);
	}
	return tasks, nil
}

func SaveTasksToFile(tasks []Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err;
	}
	defer file.Close();
	for _, task := range tasks {
		_, err = fmt.Fprintf(file, "%d %s %t\n", task.ID, task.Description, task.Completed)
		if err != nil {
			return err;
		}
	}
	return nil;
}
```

*   Implements the core logic for to-do list operations, file loading and saving.
*   Functions are exported to be used in `main.go` and `todo_test.go`.

**3. `todo/todo_test.go` (To-Do Logic Tests):**

```go
// todo/todo_test.go
package todo

import (
	"os"
	"reflect"
	"testing"
)

func TestAddTask(t *testing.T) {
	tasks := []Task{}
	AddTask(&tasks, "Buy groceries")
	if len(tasks) != 1 || tasks[0].Description != "Buy groceries" || tasks[0].Completed != false {
		t.Errorf("AddTask failed, tasks: %v", tasks);
	}

	AddTask(&tasks, "Go to gym");
	if len(tasks) != 2 || tasks[1].Description != "Go to gym" {
		t.Errorf("AddTask failed, tasks: %v", tasks);
	}
}

func TestCompleteTask(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Buy groceries", Completed: false},
		{ID: 2, Description: "Go to gym", Completed: false},
	}

	CompleteTask(&tasks, 1);

	if tasks[0].Completed != true || tasks[1].Completed != false {
		t.Errorf("CompleteTask failed, tasks: %v", tasks);
	}
	CompleteTask(&tasks, 3); // should do nothing
	if tasks[0].Completed != true || tasks[1].Completed != false {
		t.Errorf("CompleteTask failed, tasks: %v", tasks);
	}
}

func TestLoadAndSaveTasks(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Buy groceries", Completed: false},
		{ID: 2, Description: "Go to gym", Completed: true},
	}

	filename := "test_tasks.txt";
	err := SaveTasksToFile(tasks, filename);
	if err != nil {
		t.Fatalf("SaveTasksToFile failed: %v", err);
	}

	loadedTasks, err := LoadTasksFromFile(filename);
	if err != nil {
		t.Fatalf("LoadTasksFromFile failed: %v", err);
	}

	if !reflect.DeepEqual(tasks, loadedTasks) {
		t.Errorf("LoadTasksFromFile failed: expected %v, got %v", tasks, loadedTasks);
	}

	os.Remove(filename); // cleanup
}
```

*   Tests the logic of adding tasks, completing tasks, and saving and loading the tasks from a file.
*   Uses a variety of assertion techniques to check different aspects of the functionality.

**4. `main.go` (Main Application Logic):**

```go
// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
	"todolist/todo"
)


func main() {
	tasks, err := todo.LoadTasksFromFile("tasks.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading tasks from file: %v\n", err);
		os.Exit(1)
	}
	
	reader := bufio.NewReader(os.Stdin)

	for {
	fmt.Print("Enter command (add/list/complete/exit): ");
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	parts := strings.SplitN(input, " ", 2);
	command := parts[0];

	switch command {
	case "add":
		if len(parts) < 2 {
			fmt.Println("Description required!");
			continue;
		}
		todo.AddTask(&tasks, parts[1])
		err = todo.SaveTasksToFile(tasks, "tasks.txt");
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving tasks to file: %v\n", err);
		}
	case "list":
		todo.ListTasks(tasks)
	case "complete":
		if len(parts) < 2 {
			fmt.Println("Task ID required!");
			continue;
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Invalid task id!");
			continue;
		}
		todo.CompleteTask(&tasks, id)
		err = todo.SaveTasksToFile(tasks, "tasks.txt");
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving tasks to file: %v\n", err);
		}
	case "exit":
		fmt.Println("Exiting the application");
		return
	default:
		fmt.Println("Invalid command. Please enter add/list/complete or exit");
		}
	}
}
```
*   Imports the `todo` package and uses its exported functions.
*   Provides command-line user interface to interact with the To-do logic.

**How to Run the Project and Tests:**

1.  **Run the Application:**
    Open a terminal, navigate to the `todolist` directory, and run:

    ```bash
    go run .
    ```

    Follow the prompts to add, list, and complete tasks.

2.  **Run the Tests:**
    Open a terminal, navigate to the `todolist` directory, and run:

    ```bash
    go test ./...
    ```

**Output (If Tests Pass):**

```
ok      todolist/todo   0.003s
```

**Output (If Tests Fail):**
If any tests fail, the output will provide a description of the failure, the expected output and the actual output.

```
--- FAIL: TestAddTask (0.00s)
    todo_test.go:14: AddTask failed, tasks: [{1 Buy groceries false} {2 Go to the gym true}]
FAIL
exit status 1
FAIL    todolist/todo   0.003s
```

**Key Features of This Project:**

*   **Complete CLI Application:** Implements the full functionality of the to-do list manager.
*   **Modularity:** Uses a `todo` package to encapsulate the core logic.
*   **Testing:** Includes unit tests to verify the correctness of the logic.
*   **Persistent Storage:** Saves and loads data to a text file to persist between executions.

This project provides a practical example of how to build a CLI tool in Go, combining user interaction with core logic, and incorporating tests.


