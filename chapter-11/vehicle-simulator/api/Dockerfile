FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /api-app main.go
COPY ./docs ./docs
EXPOSE 8080
CMD ["/api-app"]
