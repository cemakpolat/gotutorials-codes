#!/bin/sh
while ! nc -z controller 8081; do
    echo "Waiting for controller..."
    sleep 1
done
echo "Controller is ready!"
exec "$@"
