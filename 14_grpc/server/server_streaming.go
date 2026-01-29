package main

import (
	"log"
	"time"

	pb "github.com/AVVKavvk/grpc/proto"
	"google.golang.org/grpc"
)

// Server Streaming: Server calls stream.Send() multiple times
func (us *userServer) GetAllUsers(req *pb.GetAllUserRequest, stream grpc.ServerStreamingServer[pb.User]) error {

	log.Println("\n .................... Sending User data..................")
	for _, user := range us.Users {
		log.Printf("Sending User: %v \n", user)
		err := stream.Send(user)
		time.Sleep(1 * time.Second)

		if err != nil {
			return err
		}
	}
	log.Println("\n ....................Done With Sending User data..................")
	return nil
}
