// grpc/client/client.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "grpc/service"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMyServiceClient(conn)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		name := fmt.Sprintf("John %s", time.Now().Format(time.RFC3339))
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			log.Printf("could not greet: %v", err)
		} else {
			log.Printf("Greeting: %s", r.GetMessage())
		}

		time.Sleep(5 * time.Second)
	}
}
