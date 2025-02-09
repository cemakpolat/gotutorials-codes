package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan Message
	Room string
}

// ReadPump handles incoming messages from the WebSocket connection.
// func (c *Client) ReadPump() {
// 	defer func() {
// 		c.Hub.unregister <- c
// 		c.Conn.Close()
// 	}()

// 	for {
// 		var msg Message
// 		err := c.Conn.ReadJSON(&msg)
// 		print("read pump")
// 		if err != nil {
// 			break
// 		}

//			msg.Timestamp = time.Now()
//			c.Hub.broadcast <- msg
//		}
//	}
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Set the room and timestamp in the message
		msg.Room = c.Room
		msg.Timestamp = time.Now().UTC()

		// Save the message to MongoDB
		_, err = messageCollection.InsertOne(context.Background(), msg)
		if err != nil {
			log.Println("Error saving message to MongoDB:", err)
		}

		// Send the message to the room
		c.Hub.broadcast <- msg
	}
}

// WritePump handles outgoing messages to the WebSocket connection.
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		msg, ok := <-c.Send
		if !ok {
			return
		}

		c.Conn.WriteJSON(msg)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}

	// Get the chat room from the URL query parameter
	room := r.URL.Query().Get("room")
	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan Message),
		Room: room,
	}

	// Register the new client to the hub
	client.Hub.register <- client

	// Start reading and writing in separate goroutines
	go client.ReadPump()
	go client.WritePump()
}
