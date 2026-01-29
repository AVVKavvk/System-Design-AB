package main

import (
	"context"
	"log"

	pb "github.com/AVVKavvk/grpc/proto"
)

func createUser(client pb.UserServiceClient, user *pb.CreateUserRequest) error {

	result, err := client.CreateUser(context.Background(), user)

	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
		return err
	}

	log.Printf("User created: %v", result)

	return nil
}
