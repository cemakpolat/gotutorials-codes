// controller/main.go
package main

import (
	"context"
	"controller/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
)

type Task struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	AssignedTo string `json:"assigned_to"`
	Status     string `json:"status"`
}

type VehicleData struct {
	ID       string         `json:"id"`
	Location utils.Location `json:"location"`
	Battery  int            `json:"battery"`
	Status   string         `json:"status"`
}

func main() {
	ctx := context.Background()
	// Redis setup
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis service name from docker compose
		Password: "",           // No password in this example
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to ping redis, %v", err)
	}
	fmt.Println("Successfully connected to redis")

	// MQTT setup
	var client mqtt.Client
	var token mqtt.Token
	retryInterval := 5 * time.Second
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		opts := mqtt.NewClientOptions()
		opts.AddBroker("tcp://emqx:1883") // use the service name to connect to emqx
		opts.SetClientID("go-mqtt-controller")
		client = mqtt.NewClient(opts)

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
	defer client.Disconnect(250)

	fmt.Println("Controller connected to mqtt")
	dataTopic := "vehicle/data"
	taskCompleteTopic := "vehicle/task-completed"

	// MQTT subscriber for data
	vehicleData := make(map[string]utils.VehicleData)
	var mu sync.RWMutex // Mutex to protect the shared vehicleData variable

	client.Subscribe(dataTopic, 1, func(client mqtt.Client, msg mqtt.Message) {
		var vehicleDataFromMqtt utils.VehicleData
		err = json.Unmarshal(msg.Payload(), &vehicleDataFromMqtt)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			return
		}
		mu.Lock()
		vehicleData[vehicleDataFromMqtt.ID] = vehicleDataFromMqtt
		mu.Unlock()
		log.Println("Received vehicle data and saved it to memory, vehicle id:", vehicleDataFromMqtt.ID)
	})

	// MQTT Subsciber for completed tasks
	client.Subscribe(taskCompleteTopic, 1, func(client mqtt.Client, msg mqtt.Message) {
		var task Task
		err = json.Unmarshal(msg.Payload(), &task)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			return
		}
		log.Println("Received task completed for task id: ", task.ID, ", from vehicle id: ", task.AssignedTo)

		err = rdb.Del(ctx, task.ID).Err()
		if err != nil {
			log.Println("Error deleting task from redis: ", err)
		}

	})

	// Task generator
	go func() {
		for {
			time.Sleep(20 * time.Second) // create a new task every 5 seconds.
			mu.RLock()
			keys := make([]string, 0, len(vehicleData))
			for k := range vehicleData {
				keys = append(keys, k)
			}
			mu.RUnlock()

			if len(keys) == 0 {
				log.Println("There are no available vehicles")
				continue
			}

			vehicleId := keys[rand.Intn(len(keys))] // choose a random vehicle for the task.
			taskId := utils.GenerateRandomID()
			task := Task{
				ID:         taskId,
				Type:       "delivery",
				AssignedTo: vehicleId,
				Status:     "pending",
			}
			jsonTask, err := utils.StructToJson(task)
			if err != nil {
				log.Println("Error converting task to json:", err)
				continue
			}

			err = rdb.Set(ctx, taskId, jsonTask, 0).Err()
			if err != nil {
				log.Println("Error saving the task to redis:", err)
				continue
			}
			topic := fmt.Sprintf("vehicle/task/%s", vehicleId) // send the task to a dedicated topic.
			client.Publish(topic, 1, false, jsonTask)
			log.Println("Generated task:", taskId, "assigned to vehicle id", vehicleId)
			mu.Lock()
			if vehicle, exists := vehicleData[vehicleId]; exists {
				vehicle.Status = "working"
				vehicleData[vehicleId] = vehicle
			}
			mu.Unlock()
		}
	}()

	http.HandleFunc("/vehicles", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json, err := json.MarshalIndent(vehicleData, "", "    ")
		if err != nil {
			http.Error(w, "Error while marshalling vehicle data", http.StatusInternalServerError)
			return
		}
		w.Write(json)
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		var tasks []Task

		keys, err := rdb.Keys(ctx, "*").Result()
		if err != nil {
			http.Error(w, "Error reading keys from redis", http.StatusInternalServerError)
			log.Println("Error reading keys from redis:", err)
			return
		}

		for _, key := range keys {
			if key == "latest_sensor_data" { // skip this one if it exists
				continue
			}
			taskJson, err := rdb.Get(ctx, key).Result()
			if err != nil {
				log.Println("Error reading tasks from redis", err)
				continue
			}
			var task Task
			err = json.Unmarshal([]byte(taskJson), &task)
			if err != nil {
				log.Println("Error unmarshalling task", err)
				continue
			}
			tasks = append(tasks, task)
		}
		w.Header().Set("Content-Type", "application/json")
		json, err := json.MarshalIndent(tasks, "", "    ")
		if err != nil {
			http.Error(w, "Error marshalling tasks to json", http.StatusInternalServerError)
			return
		}
		w.Write(json)
	})

	fmt.Println("Controller is running, using mqtt at port 1883 and http at 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
