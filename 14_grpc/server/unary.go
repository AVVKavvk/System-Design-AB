package main

import (
	"context"
	"fmt"
	"math/rand"

	pb "github.com/AVVKavvk/grpc/proto"
)

// Unary: Standard Request-Response
func (us *userServer) CreateUser(c context.Context, userReq *pb.CreateUserRequest) (*pb.User, error) {
	userId := generateThreeDigitUserId()

	user := &pb.User{
		Id:    userId,
		Name:  userReq.Name,
		Email: userReq.Email,
		Age:   userReq.Age,
	}

	_, exists := us.Users[userId]

	if exists {
		return nil, fmt.Errorf("user with id %s already exists", userId)
	}

	us.Users[userId] = user

	return user, nil
}

func generateThreeDigitUserId() string {
	return fmt.Sprintf("%03d", rand.Intn(1000))
}
