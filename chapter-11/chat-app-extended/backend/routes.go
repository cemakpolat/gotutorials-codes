package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(hub *Hub) *mux.Router {
	r := mux.NewRouter()

	// User routes for registration and login
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/messages", GetMessagesHandler).Methods("GET")

	// WebSocket endpoint
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	}).Methods("GET")

	return r
}
