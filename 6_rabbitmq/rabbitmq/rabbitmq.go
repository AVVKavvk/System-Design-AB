package rabbitmq

import (
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RabbitMQConnection amqp.Connection
	once               sync.Once
)

func GetRabbitMQConnection() amqp.Connection {
	return RabbitMQConnection
}

func init() {
	once.Do(func() {
		var err error
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			panic(err)
		}
		RabbitMQConnection = *conn
		fmt.Println("Rabbitmq connected")
	})
}
