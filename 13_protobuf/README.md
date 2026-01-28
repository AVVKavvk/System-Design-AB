# Protobuf

**Protocol Buffers** (Protobuf) is a free and open-source cross-platform data format used to serialize structured data. Developed by Google, it is similar to XML or JSON but **smaller, faster, and simpler**.

Think of it as a way to package data so it can be sent over a network or stored in a file, but instead of being human-readable text (like JSON), it is converted into a **compact binary format**.

## How It Works

The Protobuf workflow involves three main steps:

1. **Define the Schema**: You write a .proto file to define how your data is structured.

2. **Compile**: You use the Protobuf compiler (protoc) to generate code in your preferred language (Go, Python, Java, C++, etc.).

3. **Serialize/Deserialize**: You use the generated code to convert your data objects into binary strings and back again.

## Why Use Protobuf?

1. **Performance and Efficiency**

   Because Protobuf is a binary format, it is significantly more compact than JSON. This leads to:
   - **Reduced Bandwidth**: Less data is sent over the wire.

   - **Faster Parsing**: Computers can read binary much faster than they can parse text-based formats like JSON or XML.

2. **Strong Typing and Validation**

   In JSON, you can accidentally send a string where a number is expected. In Protobuf, the schema is strictly enforced. If your `.proto` file says a field is an `int32`, the generated code ensures it stays that way.

3. **Backward and Forward Compatibility**

   Protobuf is designed to handle evolving data structures. You can add new fields to your message format without breaking old services that aren’t updated to use them yet.

4. **Code Generation**

   Instead of manually writing boilerplate code to parse data, Protobuf generates the classes/structs and methods for you across multiple languages. This makes it a backbone for **gRPC**.

## Comparison: Protobuf vs. JSON

| Feature    | Protocol Buffers             | JSON                      |
| ---------- | ---------------------------- | ------------------------- |
| **Format** | Binary (not human-readable)  | Text (human-readable)     |
| **Size**   | Very Small                   | Larger                    |
| **Speed**  | Extremely Fast               | Slower                    |
| **Schema** | Required (.proto)            | Optional                  |
| **Usage**  | Internal APIs, Microservices | Public APIs, Web Browsers |

## A Simple Example

If you were defining a user profile, your `.proto` file would look like this:

```protobuf
syntax = "proto3";

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
}
```

The numbers (= 1, = 2) are "tags" used to identify fields in the binary format, which is part of why it stays so small—it doesn't need to send the full field name "email" every single time.

<br/>

# Go Protobuf Code

This repo included a golang code with `protoc` as a compiler through docker

## Docker

1. Up the container

   ```bash
   docker compose up -d
   ```

2. Check container status

   ```bash
   docker ps
   ```

## Compile `.proto` file

### User

Compile `user.proto` using `protoc` using docker

```docker
docker exec -it proto-gen protoc --go_out=. --go_opt=paths=source_relative proto/user.proto
```

### Post

Compile `post.proto` using `protoc` using docker

```docker
docker exec -it proto-gen protoc --go_out=. --go_opt=paths=source_relative proto/post.proto
```

## Run the program

1. Using `air`

   ```bash
   air
   ```

2. Using `go run`

   ```bash
   go run main.go
   ```

## Test Program

In new terminal run the test script

```bash
 go run tests/test_client.go
```

### Test script

```go
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
```
