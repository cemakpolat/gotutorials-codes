package utils

import (
	"concurrentpipeline/data"
	"testing"
)

func TestProcessData(t *testing.T) {
	jobs := make(chan data.Data, 10)
	results := make(chan int, 10)

	// Populate the jobs channel
	for i := 0; i < 10; i++ {
		jobs <- data.Data{Value: i}
	}
	close(jobs)

	// Start the worker
	go ProcessData(1, jobs, results)

	// Collect results from the results channel
	var receivedResults []int
	for result := range results {
		receivedResults = append(receivedResults, result)
	}

	// Define the expected results
	expectedResults := []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}

	// Validate the results
	if len(expectedResults) != len(receivedResults) {
		t.Fatalf("expected %d results, but got %d", len(expectedResults), len(receivedResults))
	}

	for i := 0; i < len(expectedResults); i++ {
		if expectedResults[i] != receivedResults[i] {
			t.Fatalf("expected %d but got %d at index %d", expectedResults[i], receivedResults[i], i)
		}
	}
}
