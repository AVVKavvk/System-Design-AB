package main

import (
	"context"
	"io"
	"log"

	pb "github.com/AVVKavvk/grpc/proto"
)

func getAllUsers(client pb.UserServiceClient) error {

	stream, err := client.GetAllUsers(context.Background(), &pb.GetAllUserRequest{})

	if err != nil {
		return err
	}

	log.Println("\n .........................Getting user from server.....................")
	for {
		user, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Failed to get user: %v", err)
		}

		log.Printf("Get User: %v", user)
	}

	log.Println("\n .........................Done with getting user from server...................")
	return nil
}
