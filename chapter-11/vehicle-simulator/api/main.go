// api/main.go
package main

import (
	_ "api/docs" // import the generated docs
	"fmt"
	"io"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger" // Import the swagger middleware
)

//	@title Swagger Example API
//	@version 1.0
//	@description This is a sample server for Autonomous Vehicles.
//	@termsOfService http://swagger.io/terms/
//  @BasePath /

type Task struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	AssignedTo string `json:"assigned_to"`
	Status     string `json:"status"`
}

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

// @Summary Get all vehicles
// @Description Get the information of all the vehicles that have been connected to the controller.
// @Produce json
// @Success 200 {array} map[string]VehicleData
// @Router /vehicles [get]
func handleVehicles(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://controller:8081/vehicles")
	if err != nil {
		http.Error(w, "Error making request to controller", http.StatusInternalServerError)
		log.Println("Error making request to controller: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error receiving data from controller", resp.StatusCode)
		log.Println("Received invalid status code: ", resp.StatusCode)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading body from controller", http.StatusInternalServerError)
		log.Println("Error reading body from controller: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// @Summary Get all tasks
// @Description Get the information of all tasks that are being managed by the controller.
// @Produce json
// @Success 200 {array} Task
// @Router /tasks [get]
func handleTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://controller:8081/tasks")
	if err != nil {
		http.Error(w, "Error making request to controller", http.StatusInternalServerError)
		log.Println("Error making request to controller: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error receiving data from controller", resp.StatusCode)
		log.Println("Received invalid status code: ", resp.StatusCode)
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading body from controller", http.StatusInternalServerError)
		log.Println("Error reading body from controller: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/vehicles", handleVehicles)
	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler) // handler for swagger.

	fmt.Println("API is running at port 8080 and swagger at /swagger")
	log.Fatal(http.ListenAndServe(":8080", nil))
}