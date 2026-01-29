FROM golang:1.23-alpine

RUN apk add --no-cache protobuf git

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2 \
 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

ENV PATH="/go/bin:$PATH"

WORKDIR /workspace