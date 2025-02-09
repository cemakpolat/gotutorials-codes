package chat

import (
	"bufio"
	"net"
	"sync"
	"testing"
	"time"
)

func TestBroadcastMessage(t *testing.T) {
	// Simulate two clients
	listener, err := net.Listen("tcp", "127.0.0.1:0") // selecting random port
	if err != nil {
		t.Fatalf("Error creating listener: %v", err)
	}
	defer listener.Close()

	// Create 2 clients to simulate chat
	var conns []net.Conn
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			conn, err := net.Dial("tcp", listener.Addr().String())
			if err != nil {
				t.Fatalf("Error connecting to server: %v", err)
			}
			conns = append(conns, conn)
		}()
	}

	wg.Wait() // wait for all connections to be established

	for _, conn := range conns {
		go HandleConnection(conn)
	}

	go ManageConnections()

	// Send a message from the first client
	message := "Test message from client 1\n"

	client := Client{conn: conns[0], writer: bufio.NewWriter(conns[0])}
	clients[conns[0]] = client // manually adding the client as we did not receive it in the client channel

	messageChannel <- message

	// Give the server some time to broadcast the message
	time.Sleep(time.Second)

	// Read the message received by the other client
	reader := bufio.NewReader(conns[1])
	receivedMessage, err := reader.ReadString('\n')

	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}
	if receivedMessage != message {
		t.Fatalf("Incorrect message received by client, expected %q got %q", message, receivedMessage)
	}

	//close the connections.
	for _, conn := range conns {
		conn.Close()
	}
}
