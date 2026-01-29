package main

import (
	"context"
	"log"
	"time"

	pb "github.com/AVVKavvk/grpc/proto"
)

func getUsersByIds(client pb.UserServiceClient, ids []string) error {

	stream, err := client.GetUsersByIds(context.Background())

	log.Println("*********************** Getting users from server by ids ************************")

	if err != nil {
		log.Fatalf("Error while getting user: %v", err)
		return err
	}

	for _, id := range ids {

		var req pb.GetUserByIdRequest

		req.Id = id

		log.Printf("Sending request for user with id: %v\n", req.Id)
		err := stream.Send(&req)

		time.Sleep(1 * time.Second) // sleep for 1 second for simulation

		if err != nil {
			log.Fatalf("Error while getting user: %v", err)
			return err
		}
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while getting user: %v", err)
		return err
	}

	log.Printf("Get Users from server : %v", res)

	return nil
}
