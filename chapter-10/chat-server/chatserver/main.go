package main

import (
	"chatserver/chat"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error creating listener: ", err)
	}
	defer listener.Close()

	fmt.Println("Server is listening at port 8080")
	go chat.ManageConnections()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error accepting connection: ", err)
			continue
		}
		go chat.HandleConnection(conn)
	}
}
