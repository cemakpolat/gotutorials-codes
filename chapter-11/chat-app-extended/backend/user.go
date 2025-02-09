package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

// RegisterHandler registers a new user.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// LoginHandler authenticates a user and returns a JWT token.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user User
	err = userCollection.FindOne(context.Background(), bson.M{"username": credentials.Username}).Decode(&user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := GenerateJWT(user.Username)
	w.Write([]byte(token))
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")

	// Retrieve messages from the database for a specific room
	cursor, err := messageCollection.Find(context.Background(), bson.M{"room": room})
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var messages []Message
	for cursor.Next(context.Background()) {
		var msg Message
		if err := cursor.Decode(&msg); err != nil {
			log.Println("Error decoding message:", err)
			continue
		}
		messages = append(messages, msg)
	}

	// Return the messages as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Println("Error encoding messages:", err)
	}
}
