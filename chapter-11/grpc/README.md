## gRPC example 

gRPC is a high-performance RPC framework developed by Google, often used in microservices. Go has excellent support for gRPC.

Core Functionality:

- Define a gRPC service using protocol buffers.
- Implement a gRPC server.
- Create a gRPC client.

First, install the needed libraries:
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc@latest
```

Then generate the go file from the proto using the following command
`protoc --go_out=. --go-grpc_out=. proto/service.proto`


This file implements a gRPC server listening on port 50051, and creates a handler function for the defined function in the proto file.
This program defines a gRPC client, and calls the SayHello function in the gRPC server.
Requires the dependencies google.golang.org/grpc and google.golang.org/protobuf.


Install protoc

`brew install protobuf`

`protoc --version`


` protoc --go_out=. --go-grpc_out=. grpc/proto/service.proto`


Create a Script grpc/run_clients.sh (Shell Script):
Create a file in the grpc directory called run_clients.sh, and include the following content:
```
#!/bin/bash

#Number of clients to launch
NUM_CLIENTS=${1:-3}

for (( i=0; i<$NUM_CLIENTS; i++ ))
do
    go run client/client.go & # Run client in background
done

wait # Wait for background processes to finish. (in our case it will wait indefinitely)
```

This script will receive the number of clients to launch as the first argument. The default is 3.
It uses a for loop to launch each client in background using & at the end of the line.
It uses wait to wait for all background process, which in our case will wait indefinitely.
Make the script executable: Run the following command to make the shell script executable:

`chmod +x run_clients.sh`


Create the run_clients.sh file in the grpc directory and paste the content above.
Make the script executable: run chmod +x run_clients.sh to make the shell script executable.
Run the gRPC Server:
Open a terminal, navigate to the grpc/server directory, and run the gRPC server:
cd grpc/server
go run server.go
cd ../..


Run Multiple Clients:
Open a new terminal, navigate to the grpc directory, and run the run_clients.sh script using one of the commands below:
To launch 5 clients:
`./run_clients.sh 5`

Key Points:
run_clients.sh: This script simplifies launching multiple client instances, simulating real-world scenarios.
Background Execution: The & symbol runs each client instance in the background, allowing all the clients to be executed concurrently.
Concurrency: By running the clients concurrently, you can see the server processing multiple requests from different clients at once.
Flexibility: You can easily modify the run_clients.sh to launch a different number of clients by passing a different value as the first argument.

**Updated Instructions:**
1.  **Install `protoc` and go plugins**: Install the `protoc` compiler and the go plugins.

 ```bash
 go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
 go install google.golang.org/protobuf/cmd/protoc@latest
 ```
2. **Create the `proto` directory**: Create the `proto` directory inside your project
 ```bash
 mkdir proto
 ```
3. **Create the `service.proto` file:** Create a new file called `service.proto` inside the `proto` directory and paste the code for the proto file inside of it.
4.  **Generate Go code:** You must generate the go code based on the proto file by using the following command. Make sure to run this command in the root of your project, where the `proto` folder is located:
 ```bash
  protoc --go_out=. --go-grpc_out=. proto/service.proto
 ```
 This will generate two files: `service.pb.go` and `service_grpc.pb.go` inside the directory where you executed the command, as specified by the `option go_package = "./;service";`. This file contains the Go code for the server and client.
5.  **Create Directories for the client and server:** Create two folders `client` and `server` in the root of your project.
6.  **Paste Code:** Paste the code for the `client.go` and `server.go` in their respective folders.
7.  **Run the Server:**
 Open a terminal, navigate to the `server` directory, and run:

     ```bash
     go run server.go
     ```

**Complete Folder Structure for the gRPC Example**

Here's the complete structure:

```
grpc/
├── client/
│   └── client.go
├── proto/
│   └── service.proto
├── server/
│   └── server.go
├── go.mod
└── go.sum
```

**Description:**

*   **`grpc/` (Root Directory):**
    *   This is the main directory that encompasses the entire gRPC example project.
    *   It contains the `go.mod` and `go.sum` files for dependency management.
*   **`grpc/client/` (Client Directory):**
    *   Contains the Go source file for the gRPC client.
    *   **`client.go`**: The Go source file for the client that contains the code to connect to the gRPC server and make calls.
*   **`grpc/proto/` (Proto Directory):**
    *   Contains the `.proto` file that defines the gRPC service and messages.
    *   **`service.proto`**: The protocol buffer definition file, containing the service definition.
*   **`grpc/server/` (Server Directory):**
    *   Contains the Go source file for the gRPC server.
    *   **`server.go`**: The Go source file for the gRPC server that implements the service defined in the `.proto` file.

*  **`go.mod` and `go.sum`**: The go module definition and sum files.

**File Contents (Recap)**

1.  **`grpc/proto/service.proto` (Protocol Buffer Definition):**

    ```protobuf
    syntax = "proto3";

    package service;

    option go_package = "grpc/service";

    service MyService {
      rpc SayHello (HelloRequest) returns (HelloResponse);
    }

    message HelloRequest {
      string name = 1;
    }

    message HelloResponse {
      string message = 1;
    }
    ```

2.  **`grpc/server/server.go` (gRPC Server Implementation):**

    ```go
	// grpc/server/server.go
	package main

	import (
		"context"
		"fmt"
		"log"
		"net"
		"google.golang.org/grpc"
		pb "grpc/service"
	)
	
	type server struct {
		pb.UnimplementedMyServiceServer
	}
	
	func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
		fmt.Printf("Received message with name: %v\n", req.Name);
		message := fmt.Sprintf("Hello %v, from gRPC server", req.Name);
		return &pb.HelloResponse{Message: message}, nil;
	}
	
	func main() {
		listener, err := net.Listen("tcp", ":50051");
		if err != nil {
			log.Fatalf("Error creating the listener: %v", err);
		}
		
		s := grpc.NewServer();
		pb.RegisterMyServiceServer(s, &server{})
		fmt.Println("gRPC server listening at 50051")
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Error serving the grpc: %v", err);
		}
	}
    ```

3.  **`grpc/client/client.go` (gRPC Client Implementation):**

    ```go
	// grpc/client/client.go
	package main

	import (
		"context"
		"fmt"
		"log"
		"time"
		"google.golang.org/grpc"
		pb "grpc/service"
	)
	
	func main() {
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock());
		if err != nil {
			log.Fatalf("did not connect %v", err);
		}
		defer conn.Close();
		
		c := pb.NewMyServiceClient(conn);
		ctx, cancel := context.WithTimeout(context.Background(), time.Second);
		defer cancel()
		
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "John"})
		if err != nil {
			log.Fatalf("could not greet: %v", err);
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}
    ```

**How to Create the Folder Structure:**

You can use these commands in your terminal to create the folder structure:

```bash
mkdir grpc
mkdir grpc/client
mkdir grpc/proto
mkdir grpc/server
```

Then create the `go.mod` and `go.sum` files using `go mod init grpc`, and `go mod tidy`, if you have any external dependency.

And finally, create the `client.go` inside the `grpc/client` folder, the `server.go` inside the `grpc/server`, and the `service.proto` inside the `grpc/proto` folder.

**Important Note:**
Make sure to run the `protoc` command from the root folder where the `grpc` directory exists:

```bash
protoc --go_out=. --go-grpc_out=. grpc/proto/service.proto
```

After running the `protoc` command, this will generate two new files in the `grpc` directory: `service.pb.go` and `service_grpc.pb.go` and a `service` folder in `grpc`.

This file structure will help you to organize your project.