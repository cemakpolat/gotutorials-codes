package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client
var userCollection *mongo.Collection
var messageCollection *mongo.Collection

func initMongoDB() {
	// Connect to MongoDB
	// clientOptions := options.Client().ApplyURI("mongodb://MONGO_URL:27017")
	mongoURI := os.Getenv("MONGO_URL")
	if mongoURI == "" {
		log.Fatal("MONGO_URL environment variable is required")
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	db, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping the MongoDB server to ensure the connection is established
	if err := db.Ping(context.Background(), nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	// Set the user collection (MongoDB creates it if it doesn't exist)
	userCollection = db.Database("chatapp").Collection("users")
	messageCollection = db.Database("chatapp").Collection("messages")
	log.Println("MongoDB connection established and collection set.")
}

// Close the MongoDB connection when the app shuts down
func closeMongoDB() {
	if err := db.Disconnect(context.Background()); err != nil {
		log.Fatal("Error disconnecting from MongoDB:", err)
	}
	log.Println("MongoDB connection closed.")
}
