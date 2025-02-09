 FROM golang:1.21-alpine

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .
# RUN go build -o /producer-app main.go

# CMD ["/producer-app"]


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Copy producer code
COPY . .

# Add wait-for-kafka script
COPY wait-for-kafka.sh /usr/local/bin/wait-for-kafka.sh
RUN chmod +x /usr/local/bin/wait-for-kafka.sh

# Build the producer application
RUN go build -o producer .

# Use the wait-for-kafka script as the entrypoint
ENTRYPOINT ["wait-for-kafka.sh", "cp-kafka", "29092", "./producer"]
