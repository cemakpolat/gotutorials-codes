package main

import (
	"fmt"
	"math/rand"
	"sync"
	"taskqueue/queue"
)

func main() {
	const numWorkers = 5
	const numTasks = 10

	var wg sync.WaitGroup

	jobs := make(chan queue.Task, 100)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go queue.Worker(i, jobs, &wg)
	}

	for i := 0; i < numTasks; i++ {
		jobs <- queue.Task{
			ID:       i + 1,
			Priority: rand.Intn(10),
			Data:     fmt.Sprintf("data %d", i),
		}
	}

	close(jobs)
	wg.Wait()

	fmt.Println("All workers have finished")
}
