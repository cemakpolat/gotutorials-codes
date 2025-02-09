package queue

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	jobs := make(chan Task, 10)
	var wg sync.WaitGroup

	wg.Add(1)
	go Worker(1, jobs, &wg)

	for i := 0; i < 10; i++ {
		jobs <- Task{
			ID:       i + 1,
			Priority: i,
			Data:     fmt.Sprintf("data %d", i),
		}
	}

	close(jobs)
	wg.Wait() // waits until all jobs are completed.
	// if it completed without errors, the test is a success.
}

func TestWorkerMultiple(t *testing.T) {
	jobs := make(chan Task, 10)
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go Worker(i, jobs, &wg)
	}

	for i := 0; i < 10; i++ {
		jobs <- Task{
			ID:       i + 1,
			Priority: i,
			Data:     fmt.Sprintf("data %d", i),
		}
	}
	close(jobs)
	wg.Wait()
	// if it completed without errors, the test is a success.
}

func TestWorkerMultipleLongRunning(t *testing.T) {
	jobs := make(chan Task, 10)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				time.Sleep(10 * time.Millisecond) // simulate work
				if j.ID == 5 {
					time.Sleep(500 * time.Millisecond)
				}
			}
		}()
	}

	for i := 0; i < 10; i++ {
		jobs <- Task{
			ID:       i + 1,
			Priority: i,
			Data:     fmt.Sprintf("data %d", i),
		}
	}
	close(jobs)
	wg.Wait()
	// if it completed without errors, the test is a success.
}
