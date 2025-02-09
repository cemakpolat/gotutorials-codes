## Project: Autonomous Vehicle System Simulation

This project will simulate a system where:

1.  **Vehicle Simulators (Go):** Multiple autonomous vehicle simulators generate data (location, battery, status) and push it to a central controller via MQTT.
2.  **Controller (Go):** A central controller subscribes to vehicle data via MQTT, generates random tasks, assigns them to vehicles based on available resources (simulated), and stores vehicle and task status in Redis.
3.  **REST API (Go):** The controller exposes a REST API to query vehicle statuses, task assignments, and overall system state.
4.  **Redis:** A Redis database stores the vehicle and task data.
5.  **Docker:** Each component (vehicle simulator, controller, Redis) will have its own Dockerfile, and all will be orchestrated with Docker Compose.

**Project Structure:**

```
autonomous-vehicles/
├── docker-compose.yml
├── vehicle/
│   ├── Dockerfile
│   ├── main.go
│   └── utils/
│       └── utils.go
├── controller/
│   ├── Dockerfile
│   ├── main.go
│   └── utils/
│        └── utils.go
└── api/
    ├── Dockerfile
    └── main.go
```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir autonomous-vehicles
    cd autonomous-vehicles
    ```

2.  **Create the `vehicle`, `controller`, and `api` Directories:**

    ```bash
    mkdir vehicle controller api
    ```

3.  **Create Subdirectories**

    ```bash
     mkdir vehicle/utils controller/utils
    ```
4.  **Create the files inside their respective directories:**

    ```bash
    touch vehicle/Dockerfile vehicle/main.go vehicle/utils/utils.go
    touch controller/Dockerfile controller/main.go controller/utils/utils.go
    touch api/Dockerfile api/main.go
    touch docker-compose.yml
    ```

Now, you should have the following project structure:

```
autonomous-vehicles/
├── docker-compose.yml
├── vehicle/
│   ├── Dockerfile
│   ├── main.go
│   └── utils/
│       └── utils.go
├── controller/
│   ├── Dockerfile
│   ├── main.go
│   └── utils/
│        └── utils.go
└── api/
    ├── Dockerfile
    └── main.go
```

**Now, paste the following code into the corresponding files:**

**1.  `docker-compose.yml` (Docker Compose Orchestration):**

```yaml
version: "3.8"
services:
  redis:
    image: redis:latest
    ports:
      - "6379:6379"
  controller:
    build: ./controller
    ports:
      - "8081:8081"
    depends_on:
      - redis
  vehicle:
    build: ./vehicle
    depends_on:
        - controller
  api:
    build: ./api
    ports:
      - "8080:8080"
    depends_on:
      - controller
```

*   Defines three services: `redis`, `controller`, `vehicle` and `api`.
*   Specifies the build contexts, ports, and dependencies between the services.

**2.  `vehicle/Dockerfile` (Vehicle Simulator Dockerfile):**

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /vehicle-sim main.go

CMD ["/vehicle-sim"]
```

*   Specifies the base image, working directory, and commands to download dependencies, build the executable, and run it.

**3. `vehicle/main.go` (Vehicle Simulator):**

```go
// vehicle/main.go
package main

import (
    "fmt"
    "log"
    "math/rand"
	"time"
    "vehicle/utils"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
    vehicleID := utils.GenerateRandomID();

	opts := mqtt.NewClientOptions();
	opts.AddBroker("tcp://controller:1883");
	opts.SetClientID(fmt.Sprintf("vehicle-%s", vehicleID));
	client := mqtt.NewClient(opts);
	
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT: %v", token.Error());
	}

	fmt.Println("Vehicle ", vehicleID, " connected to the controller via MQTT");
	defer client.Disconnect(250);

	topic := "vehicle/data";
	
	for {
        location := utils.GenerateRandomLocation();
        battery := rand.Intn(100);

        data := utils.VehicleData{
            ID: vehicleID,
            Location: location,
            Battery: battery,
            Status: "idle",
        };
		json, err := utils.StructToJson(data);
		if err != nil {
			log.Println("Error converting to json: ", err);
			continue;
		}

        client.Publish(topic, 0, false, json)
		time.Sleep(2 * time.Second); // simulate a vehicle sending data every 2 seconds.
    }
}
```

*   Creates a mock vehicle that generates random data and publishes it to a specific mqtt topic.
*	Uses utils function to generate random data, and json encoding.

**4.  `vehicle/utils/utils.go` (Vehicle Simulator Utils):**

```go
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
```
*  This file includes all utility functions such as id generation, location generation and json encoding.

**5.  `controller/Dockerfile` (Controller Dockerfile):**

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /controller-app main.go

CMD ["/controller-app"]
```
* This Dockerfile is similar to the vehicle Dockerfile, but it builds the controller application.

**6.  `controller/main.go` (Controller Logic):**

```go
// controller/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
    "controller/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
)

type Task struct {
	ID string `json:"id"`
	Type string `json:"type"`
	AssignedTo string `json:"assigned_to"`
	Status string `json:"status"`
}

func main() {
    ctx := context.Background()
    // Redis setup
    rdb := redis.NewClient(&redis.Options{
        Addr:     "redis:6379", // Redis service name from docker compose
        Password: "",          // No password in this example
        DB:       0,
    })

	_, err := rdb.Ping(ctx).Result();
	if err != nil {
		log.Fatalf("Failed to ping redis, %v", err);
	}
	fmt.Println("Successfully connected to redis");


	// MQTT setup
	opts := mqtt.NewClientOptions();
	opts.AddBroker("tcp://0.0.0.0:1883"); // accept connection from all interfaces.
	opts.SetClientID("go-mqtt-controller");

	client := mqtt.NewClient(opts);
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT: %v", token.Error());
	}
	defer client.Disconnect(250)

    fmt.Println("Controller connected to mqtt");
	topic := "vehicle/data"
	
    // MQTT subscriber
	vehicleData := make(map[string]utils.VehicleData);
	var mu sync.RWMutex; // Mutex to protect the shared vehicleData variable
	
    client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		var vehicleDataFromMqtt utils.VehicleData;
		err = json.Unmarshal(msg.Payload(), &vehicleDataFromMqtt);
		if err != nil {
			log.Println("Error unmarshalling message:", err);
			return;
		}
		mu.Lock();
		vehicleData[vehicleDataFromMqtt.ID] = vehicleDataFromMqtt;
		mu.Unlock();
		log.Println("Received vehicle data and saved it to memory, vehicle id:", vehicleDataFromMqtt.ID)
	});

    // Task generator
    go func() {
		for {
			time.Sleep(5 * time.Second); // create a new task every 5 seconds.
			mu.RLock();
			keys := make([]string, 0, len(vehicleData));
			for k := range vehicleData {
				keys = append(keys, k)
			}
			mu.RUnlock();

			if len(keys) == 0 {
				log.Println("There are no available vehicles")
				continue;
			}

			vehicleId := keys[rand.Intn(len(keys))]; // choose a random vehicle for the task.
			taskId := utils.GenerateRandomID();
			task := Task {
				ID: taskId,
				Type: "delivery",
				AssignedTo: vehicleId,
				Status: "pending",
			}
			jsonTask, err := utils.StructToJson(task);
			if err != nil {
				log.Println("Error converting task to json:", err);
				continue;
			}

			err = rdb.Set(ctx, taskId, jsonTask, 0).Err();
			if err != nil {
				log.Println("Error saving the task to redis:", err)
				continue;
			}
			log.Println("Generated task:", taskId, "assigned to vehicle id", vehicleId);

			mu.Lock()
			if vehicle, exists := vehicleData[vehicleId]; exists {
				vehicle.Status = "working"
				vehicleData[vehicleId] = vehicle
			}
			mu.Unlock()
		}
	}()
    
	
	http.HandleFunc("/vehicles", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock();
		defer mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json, err := json.MarshalIndent(vehicleData, "", "    ");
		if err != nil {
			http.Error(w, "Error while marshalling vehicle data", http.StatusInternalServerError)
			return;
		}
		w.Write(json)
	});
	
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		var tasks []Task;

		keys, err := rdb.Keys(ctx, "*").Result();
		if err != nil {
			http.Error(w, "Error reading keys from redis", http.StatusInternalServerError);
			log.Println("Error reading keys from redis:", err);
			return
		}

		for _, key := range keys {
			if key == "latest_sensor_data" { // skip this one if it exists
				continue;
			}
			taskJson, err := rdb.Get(ctx, key).Result();
			if err != nil {
				log.Println("Error reading tasks from redis", err);
				continue
			}
			var task Task;
			err = json.Unmarshal([]byte(taskJson), &task);
			if err != nil {
				log.Println("Error unmarshalling task", err)
				continue
			}
			tasks = append(tasks, task);
		}
		w.Header().Set("Content-Type", "application/json");
		json, err := json.MarshalIndent(tasks, "", "    ");
		if err != nil {
			http.Error(w, "Error marshalling tasks to json", http.StatusInternalServerError);
			return
		}
		w.Write(json)
	})
	
	fmt.Println("Controller is running, using mqtt at port 1883 and http at 8081")
	log.Fatal(http.ListenAndServe(":8081", nil));
}
```
*   This component acts as the central controller of the application.
*   It connects to redis and listens to the `vehicle/data` mqtt topic, storing that data in a local map in memory, using a mutex to protect it.
*   It generates random tasks and stores them in redis.
*   It updates the vehicle status when a task is assigned.
*  It exposes two different http endpoints: `/vehicles` which returns all vehicle information, and `/tasks` which returns all tasks.
*  It uses the `utils` package to perform json encoding and id generation.
*  It requires libraries `github.com/eclipse/paho.mqtt.golang` and `github.com/go-redis/redis/v8`.

**7. `controller/utils/utils.go` (Controller Utils):**

```go
// controller/utils/utils.go
package utils

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)
type VehicleData struct {
    ID string `json:"id"`
    Location Location `json:"location"`
    Battery int `json:"battery"`
    Status string `json:"status"`
}

type Location struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}


func GenerateRandomID() string {
    return uuid.New().String();
}

func StructToJson(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling to json %w", err);
	}
	return jsonData, nil
}
```
*    Contains utility functions for the controller, to generate ids and json encoding.

**8. `api/Dockerfile` (API Dockerfile):**

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /api-app main.go

CMD ["/api-app"]
```

* This Dockerfile builds the `api` application.

**9. `api/main.go` (REST API):**

```go
// api/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"io"
)
type Task struct {
    ID string `json:"id"`
    Type string `json:"type"`
    AssignedTo string `json:"assigned_to"`
    Status string `json:"status"`
}

type VehicleData struct {
    ID string `json:"id"`
    Location Location `json:"location"`
    Battery int `json:"battery"`
    Status string `json:"status"`
}
type Location struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
    http.HandleFunc("/vehicles", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://controller:8081/vehicles");
		if err != nil {
			http.Error(w, "Error making request to controller", http.StatusInternalServerError);
			log.Println("Error making request to controller: ", err);
			return;
		}
		defer resp.Body.Close();

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Error receiving data from controller", resp.StatusCode);
			log.Println("Received invalid status code: ", resp.StatusCode);
			return;
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading body from controller", http.StatusInternalServerError)
			log.Println("Error reading body from controller: ", err);
			return;
		}

        w.Header().Set("Content-Type", "application/json")
        w.Write(data);
    })
	
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://controller:8081/tasks");
		if err != nil {
			http.Error(w, "Error making request to controller", http.StatusInternalServerError);
			log.Println("Error making request to controller: ", err);
			return;
		}
		defer resp.Body.Close();

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Error receiving data from controller", resp.StatusCode);
			log.Println("Received invalid status code: ", resp.StatusCode);
			return;
		}
		data, err := io.ReadAll(resp.Body);
		if err != nil {
			http.Error(w, "Error reading body from controller", http.StatusInternalServerError)
			log.Println("Error reading body from controller: ", err);
			return;
		}
		
        w.Header().Set("Content-Type", "application/json")
        w.Write(data);
    })

    fmt.Println("API is running at port 8080");
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

*   This component exposes the data that is handled by the controller.
*   It has an endpoint at `/vehicles` and another one at `/tasks` which retrieves the vehicle and task data from the controller using the corresponding http endpoints.

**How to Run the Project:**

1.  Save all the code into their respective files and directories.
2.  Open a terminal, navigate to the `autonomous-vehicles` directory, and run:

    ```bash
    docker compose up
    ```
This will start all containers and you should be able to:
-  See the mqtt messages being printed in the controller logs
-   Access the http endpoints `/vehicles` and `/tasks` of the api application at http://localhost:8080/vehicles and http://localhost:8080/tasks.
 - Check the redis database using a redis client to see the assigned tasks.

**Key Features of This Project:**

*   **Complex System Simulation:** Simulates a complete system with multiple components.
*  **MQTT and Redis Integration:** Demonstrates communication between components using MQTT, and data persistence using Redis.
*   **REST API:** Exposes a REST API for querying vehicle statuses and assigned tasks.
*   **Dockerized:** All components are containerized, simplifying deployment and management.
* **Scalability:** The system is scalable by running multiple instances of the vehicle simulator.
* **Modularity**: The system is split into different containers, each with a single responsability.

This complex project showcases Go's power in building distributed systems, and demonstrates how to use multiple technologies to perform a realistic scenario of autonomous vehicles and their management.

Let me know when you're ready for the other projects. This was a big one, so let's take a break.


## SWAGGER

Okay, let's add OpenAPI (Swagger) documentation to your `api` service. This will make it easier for developers to understand and interact with your API. We'll use the `swaggo/swag` library to generate the OpenAPI documentation.

**Project Modifications:**

1.  **Install `swaggo/swag`:**
    You'll need to install the `swaggo/swag` tool in your local environment:
    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

2.  **Add API Documentation Comments to `api/main.go`:**
    We'll add comments in the Go source code following the Swaggo conventions.

3.  **Generate Swagger Documentation:**
    We'll use the `swag init` command inside the `api` directory to generate the necessary files.

4.  **Serve Swagger UI:**
    We'll modify our API to serve Swagger UI using the generated files.

**Modified Files:**

1.  **`api/main.go` (Modified with API Documentation):**

    ```go
    // api/main.go
    package main

    import (
    	"fmt"
    	"log"
    	"net/http"
    	"io"
    	"github.com/swaggo/http-swagger" // import the swagger middleware
    	_ "api/docs" // import the generated docs
    )

    //	@title Swagger Example API
    //	@version 1.0
    //	@description This is a sample server for Autonomous Vehicles.
    //	@termsOfService http://swagger.io/terms/
	//  @BasePath /
	
    type Task struct {
        ID string `json:"id"`
        Type string `json:"type"`
        AssignedTo string `json:"assigned_to"`
        Status string `json:"status"`
    }

    type VehicleData struct {
        ID string `json:"id"`
        Location Location `json:"location"`
        Battery int `json:"battery"`
        Status string `json:"status"`
    }
    type Location struct {
        Latitude float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
    }
    
    //	@Summary Get all vehicles
    //	@Description Get the information of all the vehicles that have been connected to the controller.
    //	@Produce json
    //	@Success 200 {array} map[string]VehicleData
    //	@Router /vehicles [get]
    func handleVehicles(w http.ResponseWriter, r *http.Request) {
    	resp, err := http.Get("http://controller:8081/vehicles");
    	if err != nil {
    		http.Error(w, "Error making request to controller", http.StatusInternalServerError);
    		log.Println("Error making request to controller: ", err);
    		return;
    	}
    	defer resp.Body.Close();

    	if resp.StatusCode != http.StatusOK {
    		http.Error(w, "Error receiving data from controller", resp.StatusCode);
    		log.Println("Received invalid status code: ", resp.StatusCode);
    		return;
    	}
    
    	data, err := io.ReadAll(resp.Body)
    	if err != nil {
    		http.Error(w, "Error reading body from controller", http.StatusInternalServerError)
    		log.Println("Error reading body from controller: ", err);
    		return;
    	}
    
        w.Header().Set("Content-Type", "application/json")
        w.Write(data);
    }
    
    //	@Summary Get all tasks
    //	@Description Get the information of all tasks that are being managed by the controller.
    //	@Produce json
    //	@Success 200 {array} Task
    //	@Router /tasks [get]
    func handleTasks(w http.ResponseWriter, r *http.Request) {
    	resp, err := http.Get("http://controller:8081/tasks");
    	if err != nil {
    		http.Error(w, "Error making request to controller", http.StatusInternalServerError);
    		log.Println("Error making request to controller: ", err);
    		return;
    	}
    	defer resp.Body.Close();
    
    	if resp.StatusCode != http.StatusOK {
    		http.Error(w, "Error receiving data from controller", resp.StatusCode);
    		log.Println("Received invalid status code: ", resp.StatusCode);
    		return;
    	}
    	data, err := io.ReadAll(resp.Body);
    	if err != nil {
    		http.Error(w, "Error reading body from controller", http.StatusInternalServerError)
    		log.Println("Error reading body from controller: ", err);
    		return;
    	}
    	
        w.Header().Set("Content-Type", "application/json")
        w.Write(data);
    }
    
    func main() {
    	http.HandleFunc("/vehicles", handleVehicles)
    	http.HandleFunc("/tasks", handleTasks)
		http.HandleFunc("/swagger/*", httpSwagger.WrapHandler) // handler for swagger.
    
        fmt.Println("API is running at port 8080 and swagger at /swagger");
    	log.Fatal(http.ListenAndServe(":8080", nil))
    }
    ```
	*  The `httpSwagger.WrapHandler` will expose the swagger documentation at the `/swagger` path, and it will use the generated files in the `/docs` directory.
	*   The comments above the functions are all used by the `swag` tool to generate the OpenAPI definition of this API.
	*  The `import _ "api/docs"` is needed to make sure that the `docs` package is initialized.

2. **`api/Dockerfile` (Modified):**
```dockerfile
FROM golang:1.24-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /api-app main.go
EXPOSE 8080
CMD ["/api-app"]
```
   * No changes required in this dockerfile, only adding the `EXPOSE` instruction.

3.  **Create `api/docs/`**: We need to create a folder to store the generated files:
    ```bash
	mkdir api/docs
    ```

**How to Apply the Changes:**

1.  **Update `api/main.go` and `api/Dockerfile`**: Replace the content with the code above.
2.  **Navigate to the `api` directory:**
    ```bash
    cd api
    ```
3.  **Generate Swagger Docs:** Use the `swag init` command to generate the necessary Swagger documentation files:
    ```bash
    swag init
    ```
    This command will generate a `docs` folder in your `api` directory containing the OpenAPI specification.
4.  **Navigate back to the `autonomous-vehicles` directory:**

    ```bash
	cd ..
    ```
4. Docker update
```   
   FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /api-app main.go
COPY ./docs ./docs
EXPOSE 8080

CMD ["/api-app"]
```
Test 
`docker build -t api-app .`
`docker run -p 8080:8080 api-app`

5.  **Rebuild and Run:** Run the docker compose command again:

    ```bash
    docker compose up --build
    ```

**Accessing Swagger UI:**

After running the application, you can access the Swagger UI in your web browser by visiting `http://localhost:8080/swagger/index.html`. This provides an interactive UI to explore and test your API.

**Key Points:**

*   **`swaggo/swag`:** This library uses annotations in your code to generate OpenAPI specifications automatically.
*   **API Comments:** The comments added in `api/main.go` are crucial for generating the documentation. The library can interpret comments to generate the information about the path, method, parameters, types and response.
*   **Swagger UI:** The `/swagger/index.html` endpoint serves the interactive Swagger UI, which helps developers interact with your API endpoints in an interactive manner.

By adding Swagger documentation, you will make your API more accessible and easier to use for both yourself and other developers that might be using it.
