package direct

import (
	"fmt"

	"github.com/AVVKavvk/rabbitmq/rabbitmq"
)

func EmailConsumer() {

	rbmqc := rabbitmq.GetRabbitMQConnection()

	ch, err := rbmqc.Channel()
	if err != nil {
		panic(err)
	}

	fmt.Println("[email *], waiting for producer to produce email message")

	msgs, err := ch.Consume(EmailQueueName, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		fmt.Println("Body", msg.Body)
		fmt.Println("Exchange", msg.Exchange)
		fmt.Println("Routing Key", msg.RoutingKey)
	}
}

func WhatsappConsumer() {
	rbmqc := rabbitmq.GetRabbitMQConnection()

	ch, err := rbmqc.Channel()
	if err != nil {
		panic(err)
	}
	fmt.Println("[whatsapp *], waiting for producer to produce whatsapp message")

	msgs, err := ch.Consume(WhatsappQueueName, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		fmt.Println("Body", msg.Body)
		fmt.Println("Exchange", msg.Exchange)
		fmt.Println("Routing Key", msg.RoutingKey)
	}
}
