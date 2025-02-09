package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Hourly struct {
	Time        []string  `json:"time"`
	Temperature []float64 `json:"temperature_2m"`
}

type WeatherResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Hourly    `json:"hourly"`
}

func main() {
	url := "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&hourly=temperature_2m"
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error fetching weather data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	fmt.Println("Weather Forecast for Latitude:", weather.Latitude, "Longitude:", weather.Longitude)
	for i, time := range weather.Hourly.Time {
		fmt.Printf("Time: %s, Temperature: %.2f°C\n", time, weather.Hourly.Temperature[i])
	}
}
