package main

import (
	"concurrentpipeline/data"
	"concurrentpipeline/utils"
	"sync"
)

func main() {
	const dataCount = 100
	const numWorkers = 5

	dataChannel := make(chan data.Data, 100)
	processedDataChannel := make(chan int, 100)
	var results []int

	var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish
	var m sync.Mutex      // Mutex to protect the shared slice.

	// Data generation goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		data.GenerateData(dataCount, dataChannel)
		close(dataChannel) // Close the data channel when data generation is done
	}()

	// Worker goroutines to process data
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			utils.ProcessData(workerID, dataChannel, processedDataChannel)
		}(i)
	}

	// Aggregator goroutine
	go func() {
		wg.Wait()                   // Wait for data generation and workers to finish
		close(processedDataChannel) // Close the processed data channel
	}()

	// Process aggregated data
	utils.AggregateData(processedDataChannel, &results, &m)

	// Print final results
	utils.PrintResults(results)
}
