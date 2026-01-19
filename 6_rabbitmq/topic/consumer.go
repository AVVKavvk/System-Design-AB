package topic

import (
	"fmt"

	"github.com/AVVKavvk/rabbitmq/rabbitmq"
)

func EmailConsumer() {

	rbmqc, err := rabbitmq.GetRabbitMQConnection()
	if err != nil {
		panic(err)
	}

	ch, err := rbmqc.Channel()
	if err != nil {
		panic(err)
	}
	err = ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("[email *], waiting for producer to produce email message")

	emailQueue, err := ch.QueueDeclare(EmailQueueName, true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(emailQueue.Name, EmailQueueRoutingKey, ExchangeName, false, nil)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(emailQueue.Name, "", false, false, false, false, nil)
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
	rbmqc, err := rabbitmq.GetRabbitMQConnection()
	if err != nil {
		panic(err)
	}

	ch, err := rbmqc.Channel()
	if err != nil {
		panic(err)
	}
	err = ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	whatsappQueue, err := ch.QueueDeclare(WhatsappQueueName, true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(whatsappQueue.Name, WhatsappQueueRoutingKey, ExchangeName, false, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("[whatsapp *], waiting for producer to produce whatsapp message")

	msgs, err := ch.Consume(whatsappQueue.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		fmt.Println("Body", msg.Body)
		fmt.Println("Exchange", msg.Exchange)
		fmt.Println("Routing Key", msg.RoutingKey)
	}
}
