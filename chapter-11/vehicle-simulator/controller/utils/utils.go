// controller/utils/utils.go
package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type VehicleData struct {
	ID       string   `json:"id"`
	Location Location `json:"location"`
	Battery  int      `json:"battery"`
	Status   string   `json:"status"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Task struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	AssignedTo string `json:"assigned_to"`
	Status     string `json:"status"`
}

func GenerateRandomID() string {
	return uuid.New().String()
}

func StructToJson(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling to json %w", err)
	}
	return jsonData, nil
}

func CurrentTime() string {
	now := time.Now()
	return now.Format(time.RFC3339)
}
