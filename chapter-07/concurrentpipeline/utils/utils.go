package utils

import (
	"concurrentpipeline/data"
	"fmt"
	"sync"
	"time"
)

// ProcessData processes items from the dataChannel, doubles their value,
// and sends the results to the processedDataChannel.
func ProcessData(workerID int, dataChannel chan data.Data, processedDataChannel chan int) {
	for item := range dataChannel {
		fmt.Printf("worker: %d processing: %d\n", workerID, item.Value)
		// Simulate some work
		time.Sleep(time.Millisecond * 100)
		processedDataChannel <- item.Value * 2
	}
	fmt.Printf("worker: %d finished\n", workerID)
}

// AggregateData collects results from processedDataChannel and appends them to the results slice.
func AggregateData(processedDataChannel chan int, results *[]int, m *sync.Mutex) {
	for value := range processedDataChannel {
		m.Lock()
		*results = append(*results, value)
		m.Unlock()
	}
	fmt.Println("Finished aggregating")
}

// PrintResults outputs the final results to the console.
func PrintResults(results []int) {
	fmt.Println("Results:")
	for _, result := range results {
		fmt.Println(result)
	}
}
