package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	conn   net.Conn
	writer *bufio.Writer
}

var clients = make(map[net.Conn]Client)
var clientChannel = make(chan Client)
var messageChannel = make(chan string)
var quitChannel = make(chan net.Conn)

func HandleConnection(conn net.Conn) {
	client := Client{
		conn:   conn,
		writer: bufio.NewWriter(conn),
	}
	clientChannel <- client
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			quitChannel <- conn
			return
		}
		messageChannel <- message
	}
}

func BroadcastMessage(message string, sender net.Conn) {
	for conn, client := range clients {
		if conn != sender {
			_, err := client.writer.WriteString(message)
			if err != nil {
				log.Println("Error sending message:", err)
				quitChannel <- conn
			}
			err = client.writer.Flush() // flush the messages immediately.
			if err != nil {
				log.Println("Error flushing the messages:", err)
				quitChannel <- conn
			}
		}
	}
}

func ManageConnections() {
	for {
		select {
		case client := <-clientChannel:
			clients[client.conn] = client
			fmt.Println("New client connected: ", client.conn.RemoteAddr())

		case message := <-messageChannel:
			if len(clients) == 0 {
				continue // skip if no clients are available
			}
			// to test functionality we will use the first element of the client map
			for conn := range clients {
				BroadcastMessage(message, conn)
				break
			}
		case conn := <-quitChannel:
			fmt.Println("Client disconnected: ", conn.RemoteAddr())
			delete(clients, conn)
			conn.Close()
		}
	}
}
