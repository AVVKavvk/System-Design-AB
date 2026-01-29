package main

import (
	"log"
	"net"

	pb "github.com/AVVKavvk/grpc/proto"
	"google.golang.org/grpc"
)

type userServer struct {
	pb.UnimplementedUserServiceServer
	// dummy users map
	Users map[string]*pb.User
}

var ()

func NewUserServer() *userServer {
	return &userServer{
		Users: make(map[string]*pb.User),
	}
}

func main() {

	var port = ":8080"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	us := NewUserServer()

	pb.RegisterUserServiceServer(grpcServer, us)

	log.Printf("Server started at %v", lis.Addr())

	//list is the port, the grpc server needs to start there
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

}

func init() {

}
