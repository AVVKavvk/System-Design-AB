FROM golang:1.24-alpine

# Install networking + shell tools
RUN apk add --no-cache \
    bash \
    iptables \
    iproute2 \
    curl

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /server/bin

EXPOSE 8080

CMD ["/server/bin"]
