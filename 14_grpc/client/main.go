package main

import (
	"log"

	pb "github.com/AVVKavvk/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ()

func main() {

	var (
		port          = ":8080"
		connectionStr = "localhost" + port
	)

	conn, err := grpc.NewClient(connectionStr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// 1. Unary
	user1 := pb.CreateUserRequest{
		Name:  "Vipin Kumawat",
		Email: "kumawatvipin066@gmail.com",
		Age:   22,
	}
	user2 := pb.CreateUserRequest{
		Name:  "Avvk",
		Email: "avvk@gmail",
		Age:   22,
	}

	if err := createUser(client, &user1); err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	if err := createUser(client, &user2); err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// 2. Server Streaming
	if err := getAllUsers(client); err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}

	// 3. Client Streaming
	err = getUsersByIds(client, []string{"812", "063", "134", "284"})

	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}

	// 4. Bidi Streaming
	err = chat(client)

	if err != nil {
		log.Fatalf("Failed to chat : %v", err)
	}

}
