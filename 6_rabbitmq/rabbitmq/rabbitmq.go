package rabbitmq

import (
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Store the POINTER
var (
	instance *amqp.Connection
	mu       sync.Mutex
)

func GetRabbitMQConnection() (*amqp.Connection, error) {
	mu.Lock()
	defer mu.Unlock()

	// If connection doesn't exist or is closed, reconnect
	if instance == nil || instance.IsClosed() {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			return nil, err
		}
		instance = conn
		fmt.Println("RabbitMQ connected successfully")
	}
	return instance, nil
}
