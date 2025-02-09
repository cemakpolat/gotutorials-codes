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
		t.Errorf("AddTask failed, tasks: %v", tasks)
	}

	AddTask(&tasks, "Go to gym")
	if len(tasks) != 2 || tasks[1].Description != "Go to gym" {
		t.Errorf("AddTask failed, tasks: %v", tasks)
	}
}

func TestCompleteTask(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Buy groceries", Completed: false},
		{ID: 2, Description: "Go to gym", Completed: false},
	}

	CompleteTask(&tasks, 1)

	if tasks[0].Completed != true || tasks[1].Completed != false {
		t.Errorf("CompleteTask failed, tasks: %v", tasks)
	}
	CompleteTask(&tasks, 3) // should do nothing
	if tasks[0].Completed != true || tasks[1].Completed != false {
		t.Errorf("CompleteTask failed, tasks: %v", tasks)
	}
}

func TestLoadAndSaveTasks(t *testing.T) {
	tasks := []Task{
		{ID: 1, Description: "Buy groceries", Completed: false},
		{ID: 2, Description: "Go to gym", Completed: true},
	}

	filename := "test_tasks.txt"
	err := SaveTasksToFile(tasks, filename)
	if err != nil {
		t.Fatalf("SaveTasksToFile failed: %v", err)
	}

	loadedTasks, err := LoadTasksFromFile(filename)
	if err != nil {
		t.Fatalf("LoadTasksFromFile failed: %v", err)
	}

	if !reflect.DeepEqual(tasks, loadedTasks) {
		t.Errorf("LoadTasksFromFile failed: expected %v, got %v", tasks, loadedTasks)
	}

	os.Remove(filename) // cleanup
}
