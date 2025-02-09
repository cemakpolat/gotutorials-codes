// vehicle/utils/utils.go
package utils

import (
	"math/rand"
	"encoding/json"
	"fmt"
	"time"
	"github.com/google/uuid"
)


type Location struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}


type VehicleData struct {
    ID string `json:"id"`
    Location Location `json:"location"`
    Battery int `json:"battery"`
    Status string `json:"status"`
}

func GenerateRandomID() string {
    return uuid.New().String();
}

func GenerateRandomLocation() Location {
    return Location {
        Latitude:  rand.Float64() * 180 - 90,
        Longitude: rand.Float64() * 360 - 180,
    }
}

func StructToJson(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling to json %w", err);
	}
	return jsonData, nil
}

func CurrentTime() string {
	now := time.Now();
	return now.Format(time.RFC3339);
}
