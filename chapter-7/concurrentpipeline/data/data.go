package data

import "fmt"

type Data struct {
	Value int
}

func GenerateData(count int, dataChannel chan Data) {
	for i := 0; i < count; i++ {
		dataChannel <- Data{Value: i}
		fmt.Println("Generating data: ", i)
	}
	fmt.Println("Finished Generating data")
}
