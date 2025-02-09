package converter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func ProcessCSV(csvPath string) ([]Person, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("error opening csv file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read the header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	var people []Person
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading csv record: %w", err)
		}

		if len(record) != 3 {
			return nil, fmt.Errorf("invalid csv record, expected 3 values, got %d, record: %v", len(record), record)
		}

		var person Person
		person.Name = record[0]

		var age int
		if _, err := fmt.Sscan(record[1], &age); err != nil {
			return nil, fmt.Errorf("error parsing age: %w", err)
		}

		person.Age = age
		person.City = record[2]

		people = append(people, person)
	}
	return people, nil
}

func WriteJSON(people []Person, jsonPath string) error {
	file, err := os.Create(jsonPath)
	if err != nil {
		return fmt.Errorf("error creating json file: %w", err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(people, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling json: %w", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to json file %w", err)
	}

	return nil
}
