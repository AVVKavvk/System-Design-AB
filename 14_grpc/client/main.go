package main

import (
	"log"

	pb "github.com/AVVKavvk/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port          = ":8080"
	connectionStr = "localhost" + port
)

func main() {

	conn, err := grpc.NewClient(connectionStr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

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
}
