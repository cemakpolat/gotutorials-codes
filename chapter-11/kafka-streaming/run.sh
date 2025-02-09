#!/bin/bash

# Default topic name and number of messages to send
export KAFKA_TOPIC="${KAFKA_TOPIC:-my-topic}"
export MESSAGE_COUNT="${MESSAGE_COUNT:-10}"

# Launch Docker Compose
docker compose up --build -d


echo "Kafka is ready. Creating topic..."
docker-compose exec cp-kafka kafka-topics \
  --create \
  --bootstrap-server cp-kafka:9092 \
  --replication-factor 1 \
  --partitions 1 \
  --topic my-topic

  
echo "Kafka Cluster, producer and consumer started. Using topic: $KAFKA_TOPIC and producer will send $MESSAGE_COUNT messages"

# Wait for the producer and consumer to start
sleep 5

# See the logs
# docker compose logs -f consumer producer