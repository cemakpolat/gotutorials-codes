// grpc/server/server.go
package main

import (
	"context"
	"fmt"
	pb "grpc/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMyServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("Received message with name: %v\n", req.Name)
	message := fmt.Sprintf("Hello %v, from gRPC server", req.Name)
	return &pb.HelloResponse{Message: message}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error creating the listener: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, &server{})
	fmt.Println("gRPC server listening at 50051")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Error serving the grpc: %v", err)
	}
}
