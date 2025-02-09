Let’s move to the next project: **A Chat Application**. This will be a real-time chat app where users can join rooms and send messages to each other.  

---

### **Chat Application Overview**

#### **Basic Features (MVP)**  
1. Users can join a chat room.  
2. Messages sent in a room are broadcast to all users in that room.  
3. Built using WebSocket for real-time communication.

#### **Enhancements (Future)**  
1. User authentication with unique usernames.  
2. Persistent chat history using a database.  
3. Multiple rooms with user creation capabilities.  
4. Typing indicators and timestamps.  
5. Add a front-end using HTML/CSS/JavaScript.  

---

### **Step 1: Set Up the Project**

#### 1. **Create a Directory**  
```bash
mkdir chat-app
cd chat-app
go mod init chat-app
```

#### 2. **Directory Structure**  
```
chat-app/
├── main.go        // Entry point
├── hub.go         // Manages rooms and message broadcasting
├── client.go      // Handles individual client connections
├── README.md      // Documentation
```

---

### **Step 2: Code Implementation**

#### **1. `main.go` (Entry Point)**  
This file initializes the WebSocket server and starts the application.  

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(hub, w, r)
	})

	log.Println("Chat server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
```

---

#### **2. `hub.go` (The Hub)**  
This file manages chat rooms and broadcasting messages.  

```go
package main

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
```

---

#### **3. `client.go` (Handling Clients)**  
This file handles individual WebSocket connections.  

```go
package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait  = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
```

---

### **Step 3: Test the Application**  

#### **1. Install Dependencies**  
Install the Gorilla WebSocket package:  
```bash
go get github.com/gorilla/websocket
```

#### **2. Start the Server**  
Run the server:  
```bash
go run main.go
```

#### **3. Test with a WebSocket Client**  
Use a WebSocket testing tool like [Postman](https://www.postman.com/) or [websocat](https://github.com/vi/websocat).  

**Connect to WebSocket**:  
```bash
ws://localhost:8080/ws
```

**Send a Message**:  
Send text messages, and you'll see them broadcast to all connected clients.  

---

### **Step 4: Next Steps**  
1. Add support for multiple chat rooms.  
2. Implement persistent storage for chat history using SQLite or MongoDB.  
3. Add authentication (e.g., JWT tokens) for user management.  
4. Build a simple front-end for better user experience.  

Would you like to expand this project or move to the next one?