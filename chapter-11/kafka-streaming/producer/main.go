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
	time.Sleep(10 * time.Second) // wait for kafka to be up
	brokers := []string{getEnv("KAFKA_BROKERS", "172.25.0.3:9092")}
	topic := getEnv("KAFKA_TOPIC", "my-topic")
	messageCount, err := strconv.Atoi(getEnv("MESSAGE_COUNT", "100"))
	if err != nil {
		log.Fatal("invalid message count")
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer w.Close()

	for i := 0; i < messageCount; i++ {
		message := fmt.Sprintf("message %d", i)
		ctx := context.Background()
		err := w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(fmt.Sprintf("key %d", i)),
			Value: []byte(message),
		})
		if err != nil {
			log.Fatalf("error producing message %v", err)
		}
		log.Println("Produced message: ", message)
		time.Sleep(time.Second)
	}

	log.Println("All messages were sent")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
