package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/AVVKavvk/grpc/proto"
)

func chat(client pb.UserServiceClient) error {

	stream, err := client.UserChat(context.Background())

	if err != nil {
		panic(err)
	}

	messages := []string{"Hi", "Is this working?", "Hello", "How are you?", "Vipin", "Kumawat", "Bye!"}

	// Use a channel to wait for the server to finish
	wait := make(chan struct{})

	go func() {

		for {
			res, err := stream.Recv()

			if err == io.EOF {
				close(wait)
				return
			}
			log.Printf("Server Says: %s\n", res.Message)
		}
	}()

	for _, msg := range messages {

		req := pb.ChatMessage{
			UserId:  "user_123",
			Message: msg,
		}

		err := stream.Send(&req)

		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second) // sleep for 1 second for simulation

	}

	err = stream.CloseSend()
	if err != nil {
		panic(err)
	}

	// Since after CloseSend(), the client can't send any more messages, we can now wait for the server to finish the stream
	<-wait

	return nil
}
