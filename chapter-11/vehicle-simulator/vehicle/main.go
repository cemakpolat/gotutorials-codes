// vehicle/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
	"vehicle/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Task struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	AssignedTo string `json:"assigned_to"`
	Status     string `json:"status"`
}

func main() {
	vehicleID := utils.GenerateRandomID()

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://emqx:1883")
	opts.SetClientID(fmt.Sprintf("vehicle-%s", vehicleID))
	client := mqtt.NewClient(opts)
	var token mqtt.Token

	retryInterval := 5 * time.Second
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		token = client.Connect()
		if token.Wait() && token.Error() == nil {
			break
		}
		log.Println("Error connecting to MQTT, retrying in ", retryInterval)
		time.Sleep(retryInterval)
	}

	if token.Error() != nil {
		log.Fatalf("Error connecting to MQTT after max retries: %v", token.Error())
	}

	fmt.Println("Vehicle ", vehicleID, " connected to the controller via MQTT")
	defer client.Disconnect(250)

	dataTopic := "vehicle/data"
	taskTopic := fmt.Sprintf("vehicle/task/%s", vehicleID) // dedicated topic for this vehicle.

	// Subscribe to task topic
	client.Subscribe(taskTopic, 1, func(client mqtt.Client, msg mqtt.Message) {
		var task Task
		err := json.Unmarshal(msg.Payload(), &task)
		if err != nil {
			log.Println("Error unmarshalling task: ", err)
			return
		}
		log.Println("Vehicle ", vehicleID, " received task id:", task.ID)

		// Simulate some task processing.
		time.Sleep(time.Duration(rand.Intn(20)+5) * time.Second)
		task.Status = "completed"

		jsonTask, err := utils.StructToJson(task)
		if err != nil {
			log.Println("Error encoding task to json:", err)
			return
		}

		// Send task completed message back to controller
		client.Publish("vehicle/task-completed", 1, false, jsonTask)
		log.Println("Vehicle", vehicleID, "completed task id: ", task.ID)
	})

	// Simulate data generation and send data every 2 seconds.
	for {
		location := utils.GenerateRandomLocation()
		battery := rand.Intn(100)

		data := utils.VehicleData{
			ID:       vehicleID,
			Location: location,
			Battery:  battery,
			Status:   "idle",
		}
		json, err := utils.StructToJson(data)
		if err != nil {
			log.Println("Error converting to json: ", err)
			continue
		}

		client.Publish(dataTopic, 0, false, json)
		time.Sleep(2 * time.Second)
	}
}
