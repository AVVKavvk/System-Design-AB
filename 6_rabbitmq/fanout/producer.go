package fanout

import (
	"encoding/json"
	"time"

	"github.com/AVVKavvk/rabbitmq/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
)

func Producer() {

	var message map[string]interface{}

	message = map[string]interface{}{
		"message": "New product has been added",
	}

	messageStr, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	rbmqc, err := rabbitmq.GetRabbitMQConnection()
	if err != nil {
		panic(err)
	}

	ch, err := rbmqc.Channel()
	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(ExchangeName, "fanout", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	err = ch.Publish(ExchangeName, "", false, false, amqp091.Publishing{
		ContentType: "text/plan",
		Body:        []byte(messageStr),
	})

	if err != nil {
		panic(err)
	}

	defer ch.Close()
}
