package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/AVVKavvk/protobuf/proto"
	proto_lib "google.golang.org/protobuf/proto"
)

func main() {
	// 1. Create a valid User object
	user := &proto.User{
		Id:    1,
		Name:  "vipin",
		Email: "v@v.com",
		Age:   25,
	}

	// 2. Marshal to binary
	data, err := proto_lib.Marshal(user)
	if err != nil {
		panic(err)
	}

	// 3. Send POST request
	resp, err := http.Post("http://localhost:8080/users", "application/protobuf", bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 4. Read response
	body, _ := io.ReadAll(resp.Body)
	resUser := &proto.User{}
	proto_lib.Unmarshal(body, resUser)

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Received User: ID=%d, Name=%s, Age=%d\n", resUser.Id, resUser.Name, resUser.Age)
}
