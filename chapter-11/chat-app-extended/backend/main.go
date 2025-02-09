package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	// Initialize MongoDB connection
	initMongoDB()
	defer closeMongoDB() // Ensure MongoDB connection is closed when the application shuts down

	// Initialize the hub
	hub := NewHub()
	go hub.Run()

	// Setup routes using the SetupRoutes function
	r := SetupRoutes(hub)

	// Initialize CORS settings
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},         // Allow frontend (Angular) origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},  // Allowed HTTP methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Allowed request headers
		AllowCredentials: true,                                      // Allow credentials such as cookies
	})

	// Apply the CORS middleware to the routes
	handler := corsOptions.Handler(r)

	// Start the server
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
