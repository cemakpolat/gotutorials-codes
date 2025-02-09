package queue

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID       int
	Priority int
	Data     string
}

func Worker(id int, jobs <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Println("worker: ", id, "processing task ", job.ID, "priority", job.Priority)
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // simulate processing
	}
	fmt.Println("worker: ", id, "finished")
}
