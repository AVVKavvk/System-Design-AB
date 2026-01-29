package main

import (
	"io"
	"log"

	"google.golang.org/grpc"

	pb "github.com/AVVKavvk/grpc/proto"
)

func (us *userServer) UserChat(stream grpc.BidiStreamingServer[pb.ChatMessage, pb.ChatMessage]) error {

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		log.Printf("Client Mess: %s\n", req.Message)

		res := pb.ChatMessage{
			UserId:  req.UserId,
			Message: "Acknowledged Message: " + req.Message,
		}

		err = stream.Send(&res)
		if err != nil {
			return err
		}
	}
}
