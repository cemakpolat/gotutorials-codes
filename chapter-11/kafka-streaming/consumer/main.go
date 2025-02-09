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
	time.Sleep(10 * time.Second) // wait for kafka to be up
	brokers := []string{getEnv("KAFKA_BROKERS", "172.25.0.3:9092")}
	topic := getEnv("KAFKA_TOPIC", "my-topic")
	groupId := getEnv("KAFKA_GROUPID", "my-group")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupId,
		StartOffset: kafka.LastOffset,
	})

	defer r.Close()
	log.Println("Consumer waits for ... ")
	ctx := context.Background()
	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			log.Println("error fetching message ", err)
			time.Sleep(time.Second)
			continue
		}
		fmt.Printf("Message at offset %v: Key: %s, Value: %s\n", m.Offset, string(m.Key), string(m.Value))
		r.CommitMessages(ctx, m)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
