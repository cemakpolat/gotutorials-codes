FROM golang:1.21-alpine

# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .
# RUN go build -o /consumer-app main.go

# CMD ["/consumer-app"]



WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Copy consumer code
COPY . .

# Add wait-for-kafka script
COPY wait-for-kafka.sh /usr/local/bin/wait-for-kafka.sh
RUN chmod +x /usr/local/bin/wait-for-kafka.sh

# Build the consumer application
RUN go build -o consumer .

# Use the wait-for-kafka script as the entrypoint
ENTRYPOINT ["wait-for-kafka.sh", "cp-kafka", "29092", "./consumer"]
