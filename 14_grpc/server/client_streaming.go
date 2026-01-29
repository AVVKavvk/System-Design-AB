package main

import (
	"io"
	"log"

	pb "github.com/AVVKavvk/grpc/proto"
	"google.golang.org/grpc"
)

func (us *userServer) GetUsersByIds(stream grpc.ClientStreamingServer[pb.GetUserByIdRequest, pb.BulkUserResponse]) error {

	ids := []string{}
	log.Println("************************ Getting request from client for users by ids *****************************")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
		userId := req.Id

		log.Printf("Get request for Id: %s\n", userId)
		ids = append(ids, userId)
	}

	var res pb.BulkUserResponse

	for _, id := range ids {
		user, exists := us.Users[id]

		if exists {
			res.Users = append(res.Users, user)
		}
	}

	log.Println("************************ Sending Response at once *****************************")

	return stream.SendAndClose(&res)
}
