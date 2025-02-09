#!/bin/sh

host=$1
port=$2
shift 2
cmd="$@"

while ! nc -z "$host" "$port"; do
  echo "Waiting for Kafka at $host:$port..."
  sleep 2
done

echo "Kafka is up at $host:$port. Starting producer..."
exec $cmd
