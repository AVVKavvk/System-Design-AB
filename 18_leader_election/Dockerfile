FROM golang:1.24-alpine

WORKDIR /app

# Copy go files
COPY go.mod ./
COPY go.sum ./

COPY *.go ./

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o leader-election .

# Run the application
CMD ["./leader-election"]