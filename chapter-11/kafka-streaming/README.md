Dcd delve deeper into the Kafka example and transform it into a more complete project with enhanced functionality and a clearer structure. We'll focus on creating a producer and consumer that can be easily configured and run, providing a better demonstration of how Go interacts with Kafka.

**Project: Kafka Producer and Consumer with Configuration**

This project will create two separate Go applications:

1.  **Kafka Producer:**
    *   Reads configuration from environment variables.
    *   Connects to a Kafka cluster.
    *   Allows users to specify the topic and the number of messages to send via CLI arguments.
    *   Sends messages to Kafka with keys and values.
2.  **Kafka Consumer:**
    *   Reads configuration from environment variables.
    *   Connects to a Kafka cluster.
    *   Subscribes to a specific topic.
    *   Receives and prints messages with their keys and values.
    *   Commits the messages after reading them.

**Project Structure:**

```
kafka-streaming/
├── producer/
│   ├── main.go
│   └── go.mod
└── consumer/
    ├── main.go
    └── go.mod

```

**Setup Instructions:**

1.  **Create the Project Directory:**
    Open your terminal and run:

    ```bash
    mkdir kafka-streaming
    cd kafka-streaming
    ```

2.  **Create the `producer` and `consumer` Directories:**

    ```bash
    mkdir producer consumer
    ```

3.  **Create the `producer/main.go` and `consumer/main.go` files:**

    ```bash
    touch producer/main.go consumer/main.go
    ```

4.  **Create `go.mod` files:**
	```bash
	touch producer/go.mod
	touch consumer/go.mod
	```

Now you should have the following structure:
```
kafka-streaming/
├── producer/
│   ├── main.go
│   └── go.mod
└── consumer/
    ├── main.go
    └── go.mod
```

**Now, paste the code into the corresponding files:**

**1.  `producer/go.mod` (Producer Module):**
    ```
    module producer

    go 1.21

    require github.com/segmentio/kafka-go v0.4.47
    ```

**2.  `producer/main.go` (Kafka Producer):**

    ```go
	// producer/main.go
    package main

    import (
        "context"
        "fmt"
        "log"
		"os"
		"strconv"
		"time"
		"github.com/segmentio/kafka-go"
    )
	
    func main() {
		brokers := []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
		topic := getEnv("KAFKA_TOPIC", "my-topic")
		messageCount, err := strconv.Atoi(getEnv("MESSAGE_COUNT", "10"));
		if err != nil {
			log.Fatal("invalid message count");
		}

        w := &kafka.Writer{
            Addr:     kafka.TCP(brokers...),
            Topic:    topic,
            Balancer: &kafka.LeastBytes{},
        }

		defer w.Close();

		for i := 0; i < messageCount; i++ {
			message := fmt.Sprintf("message %d", i);
			ctx := context.Background();
			err := w.WriteMessages(ctx, kafka.Message{
				Key: []byte(fmt.Sprintf("key %d", i)),
				Value: []byte(message),
			});
			if err != nil {
				log.Fatalf("error producing message %v", err);
			}
			log.Println("Produced message: ", message);
			time.Sleep(time.Second)
		}
		
		log.Println("All messages were sent")
    }

	func getEnv(key, defaultValue string) string {
		value := os.Getenv(key);
		if value == "" {
			return defaultValue;
		}
		return value
	}
    ```

    *   Reads Kafka broker addresses, topic name, and message count from environment variables or uses defaults.
    *  Sends messages to Kafka with incremental keys and values.
	* Uses the `getEnv` helper function to read environment variables or return a default value, when the variable is not found.

**3.  `consumer/go.mod` (Consumer Module):**
    ```
    module consumer

    go 1.21

    require github.com/segmentio/kafka-go v0.4.47
    ```

**4.  `consumer/main.go` (Kafka Consumer):**

    ```go
	// consumer/main.go
	package main

	import (
		"context"
		"fmt"
		"log"
		"os"
		"time"
		"github.com/segmentio/kafka-go"
	)
	
	func main() {
		brokers := []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
		topic := getEnv("KAFKA_TOPIC", "my-topic")
		groupId := getEnv("KAFKA_GROUPID", "my-group");
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupId,
			StartOffset: kafka.LastOffset,
		})

		defer r.Close();
		ctx := context.Background();
		for {
			m, err := r.FetchMessage(ctx);
			if err != nil {
				log.Println("error fetching message ", err);
				time.Sleep(time.Second);
				continue
			}
			fmt.Printf("Message at offset %v: Key: %s, Value: %s\n", m.Offset, string(m.Key), string(m.Value));
			r.CommitMessages(ctx, m)
		}
	}

	func getEnv(key, defaultValue string) string {
		value := os.Getenv(key);
		if value == "" {
			return defaultValue;
		}
		return value
	}
    ```
    *   Reads Kafka broker addresses, topic name, and group ID from environment variables or uses defaults.
    *   Receives and prints all messages in the topic.
	*   Uses the `getEnv` helper function to read environment variables or return a default value, when the variable is not found.

**How to Run:**

1.  **Start Kafka:** You need to have a Kafka cluster running. If you don't have one, you can quickly set up a single-node Kafka cluster using Docker:

    ```bash
    docker run -d --name kafka -p 9092:9092 -p 2181:2181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092  confluentinc/cp-kafka:latest
    ```

2.  **Set Environment Variables:** If you want to customize the broker addresses or topic names, you can do that using environment variables.

3.  **Run the Producer:** Open a terminal, navigate to the `kafka-streaming/producer` directory and run the producer. For example, if you want to use a different topic name, you can do:
    ```bash
	export KAFKA_TOPIC=my-other-topic
    go run main.go
    ```
    Or to send 20 messages instead of the default 10, you can do:
     ```bash
	export MESSAGE_COUNT=20
    go run main.go
    ```

4.  **Run the Consumer:** Open a new terminal, navigate to the `kafka-streaming/consumer` directory, and run the consumer:
	```bash
	export KAFKA_TOPIC=my-other-topic
    go run main.go
    ```
	Remember to use the same topic in both the producer and the consumer.

Now you should have a Kafka producer sending messages and the consumer receiving them and printing their information in the console.

**Key Features of This Project:**

*   **Clear separation:** The project is divided into the `producer` and `consumer`, which is a good approach for clarity and maintainability.
*  **Configuration:**  The applications use environment variables to configure the kafka brokers and topic name, so that they can be modified easily.
*   **Message Details:** The consumer now prints both key and value for each message, for more clarity of how kafka works.
*   **Commit**: The consumer now commits messages after reading them.
*   **Flexibility:** You can use different command line arguments to control the behavior of the producer application, such as the number of messages sent.
*   **Testability:** The new configuration makes the programs more flexible and easier to test.

This updated example shows a more complete and flexible implementation of a Kafka producer and consumer, demonstrating how Go can be used to implement real-world data streaming scenarios.

Let me know if you have any further questions or if you want to explore any other variations!




Okay, let's break down the provided `docker-compose.yml` file, explaining each part and its purpose in setting up a Kafka cluster with Zookeeper, along with the producer and consumer services.

**`docker-compose.yml` Explanation**

```yaml
version: "3.8"
services:
  cp-zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    hostname: zookeeper
    container_name: zookeeper
    ports:
      - "2181:2181"
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  cp-kafka:
    image: confluentinc/cp-kafka:7.5.0
    hostname: kafka
    container_name: kafka
    depends_on:
    - cp-zookeeper
    ports:
    - "29092:29092"
    - "9092:9092"
    - "9997:9997"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://cp-kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:29092,PLAINTEXT_HOST://0.0.0.0:9092
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0
      CONFLUENT_METRICS_REPORTER_BOOTSTRAP_SERVERS: cp-kafka:29092
      CONFLUENT_METRICS_REPORTER_ZOOKEEPER_CONNECT: cp-zookeeper:2181
      CONFLUENT_METRICS_REPORTER_TOPIC_REPLICAS: 1
      CONFLUENT_METRICS_ENABLE: 'false'
      KAFKA_JMX_PORT: 9997
      KAFKA_MESSAGE_MAX_BYTES: 10000000
      KAFKA_SOCKET_REQUEST_MAX_BYTES: 100001200
      KAFKA_JMX_OPTS: -Dcom.sun.management.jmxremote -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.ssl=false -Djava.rmi.server.hostname=kafka0 -Dcom.sun.management.jmxremote.rmi.port=9997
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: 'true'

  
  producer:
    build: ./producer
    depends_on:
      - cp-kafka
    environment:
      KAFKA_BROKERS: "cp-kafka:29092"
      KAFKA_TOPIC: my-topic
      KAFKA_GROUPID: my-group
  
  consumer:
    build: ./consumer
    depends_on:
      - cp-kafka
    environment:
      KAFKA_BROKERS: cp-kafka:29092
      KAFKA_TOPIC: my-topic
      KAFKA_GROUPID: my-group
```

**Explanation:**

1.  **`version: "3.8"`**:
    *   Specifies the version of the Docker Compose file format.

2.  **`services:`**:
    *   Defines the different services that will be launched as containers.

3.  **`cp-zookeeper` Service**:

    *   **`image: confluentinc/cp-zookeeper:7.5.0`**: Specifies the Docker image to use for Zookeeper, which is a central service for managing the kafka cluster.
    *   **`hostname: zookeeper`**: Sets the hostname of the container to `zookeeper`.
    *  **`container_name: zookeeper`**: Sets the container name to `zookeeper`.
    *   **`ports: - "2181:2181"`**: Maps port 2181 on the host to port 2181 on the container, so that you can access zookeeper from outside the docker network.
    *   **`environment:`**: Defines environment variables:
        *   `ZOOKEEPER_CLIENT_PORT: 2181`: Sets the Zookeeper client port.
        *   `ZOOKEEPER_TICK_TIME: 2000`:  Sets the tick time for Zookeeper, defining the timeout.

4.  **`cp-kafka` Service**:

    *   **`image: confluentinc/cp-kafka:7.5.0`**: Specifies the Docker image for Kafka.
    *   **`hostname: kafka`**: Sets the hostname of the container to `kafka`.
    *   **`container_name: kafka`**: Sets the name of the container to `kafka`.
    *   **`depends_on:`**: Specifies that this service depends on the `cp-zookeeper` service, ensuring that Zookeeper starts before Kafka.
    *   **`ports:`**: Exposes several ports on the host to ports on the container:
        *   `29092:29092`: The main port for Kafka brokers, using an internal port for communication inside the docker network.
		* `9092:9092`: A secondary port for communicating from outside the docker network to the kafka cluster.
        *   `9997:9997`: A port to access Kafka JMX Metrics.
    *   **`environment:`**: Defines the Kafka-specific environment variables:
        *   `KAFKA_BROKER_ID: 1`: Sets the broker ID for the Kafka instance.
        *   `KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'`**: Specifies the Zookeeper connection string, using the `zookeeper` service name to connect.
        *   `KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://cp-kafka:29092,PLAINTEXT_HOST://localhost:9092`:
          *   Specifies that `cp-kafka:29092` is for inter-broker communication, from other containers.
          *   Specifies that `localhost:9092` is for external communication to the kafka broker.
        *   `KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:29092,PLAINTEXT_HOST://0.0.0.0:9092`: Specifies the listeners to use for all interfaces, which will make it available from within the docker network.
        *   `KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT`: Defines the security protocol for the listener.
        *   `KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT`: The protocol name used for inter-broker communication
        *   `KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1`: Configures how to replicate the internal topics, which contains information about the consumer offset.
        *   `KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0`: How much time a consumer group needs to wait before rebalancing.
        *   `CONFLUENT_METRICS_REPORTER_BOOTSTRAP_SERVERS: cp-kafka:29092`: Specifies bootstrap servers for metrics reporting.
        *   `CONFLUENT_METRICS_REPORTER_ZOOKEEPER_CONNECT: cp-zookeeper:2181`: Specifies zookeeper connection for metrics reporting.
        *  `CONFLUENT_METRICS_REPORTER_TOPIC_REPLICAS: 1`: Specifies the replication factor for internal topics of metrics.
        *   `CONFLUENT_METRICS_ENABLE: 'false'`: Disables the metrics.
        *   `KAFKA_JMX_PORT: 9997`: Exposes the JMX metrics port.
        *  `KAFKA_MESSAGE_MAX_BYTES`: The maximum size in bytes for the messages.
		*  `KAFKA_SOCKET_REQUEST_MAX_BYTES`: Maximum size of the requests for each socket.
        *   `KAFKA_JMX_OPTS`: Configures the JMX options for the Kafka.
        * `KAFKA_AUTO_CREATE_TOPICS_ENABLE: 'true'`: Enables automatic topic creation when using the kafka producer and consumer.

5.  **`producer` Service**:

    *   `build: ./producer`: Specifies the build context to be the `producer` directory
    *   **`depends_on:`**: Specifies that this service depends on the `cp-kafka` service, ensuring that Kafka is running before launching the producer.
    *   **`environment:`**: Defines the environment variables for producer application:
        *   `KAFKA_BROKERS: "cp-kafka:29092"`: Specifies the broker address to `cp-kafka:29092`, as the producer will be running inside docker.
        *   `KAFKA_TOPIC: my-topic`: Specifies the name of the topic to be used for publishing messages.
        * `KAFKA_GROUPID: my-group`: Specifies the name of the consumer group to be used to commit the offsets.

6.  **`consumer` Service**:

    *  `build: ./consumer`: Specifies the build context to be the `consumer` directory.
    *   **`depends_on:`**: Specifies that this service depends on the `cp-kafka` service, ensuring that Kafka is running before launching the consumer.
    *   **`environment:`**: Defines the environment variables for the consumer application:
        *   `KAFKA_BROKERS: cp-kafka:29092`: Specifies the broker address to `cp-kafka:29092` to be used by the consumer to connect to the kafka service, as it will be running inside docker.
		* `KAFKA_TOPIC: my-topic`: Specifies the topic name that the consumer should read the messages from.
		* `KAFKA_GROUPID: my-group`: Specifies the consumer group to connect to, for the consumer to use.

This detailed explanation provides a complete understanding of all configurations, ports, and environment variables being used in this setup, allowing you to debug problems and adapt it for other applications. This provides a clear description of the docker compose file, so you can have a clear picture of the entire process.
