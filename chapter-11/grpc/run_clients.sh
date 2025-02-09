#!/bin/bash

#Number of clients to launch
NUM_CLIENTS=${1:-3}

for (( i=0; i<$NUM_CLIENTS; i++ ))
do
    go run client/client.go & # Run client in background
done

wait # Wait for background processes to finish. (in our case it will wait indefinitely)