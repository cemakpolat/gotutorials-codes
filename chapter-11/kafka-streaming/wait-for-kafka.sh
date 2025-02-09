#!/bin/sh
while ! nc -z cp-kafka 29092; do
  echo "Waiting for Kafka..."
  sleep 1
done
echo "Kafka is up. Starting consumer..."
exec "$@"

