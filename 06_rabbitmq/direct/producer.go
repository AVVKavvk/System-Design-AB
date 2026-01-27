package direct

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AVVKavvk/rabbitmq/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
)

func Producer() {

	var emailMessage map[string]interface{}
	var whatsappMessage map[string]interface{}

	emailMessage = map[string]interface{}{
		"type":    "email",
		"from":    "v@v.com",
		"to":      "kumawatvipin066@gmail.com",
		"subject": "Order Confirmed",
		"body":    "Your order has been confirmed",
	}

	whatsappMessage = map[string]interface{}{
		"type":    "whatsapp",
		"from":    "+91 1234567890",
		"to":      "+91 9876543210",
		"message": "Your order has been confirmed",
	}

	emailMessageStr, err := json.Marshal(emailMessage)
	if err != nil {
		panic(err)
	}

	whatsappMessageStr, err := json.Marshal(whatsappMessage)
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

	err = ch.ExchangeDeclare(ExchangeName, "direct", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	emailQueue, err := ch.QueueDeclare(EmailQueueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	whatsappQueue, err := ch.QueueDeclare(WhatsappQueueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(emailQueue.Name, EmailQueueRoutingKey, ExchangeName, false, nil)
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(whatsappQueue.Name, WhatsappQueueRoutingKey, ExchangeName, false, nil)
	if err != nil {
		panic(err)
	}

	time.Sleep(30 * time.Second)

	err = ch.Publish(ExchangeName, EmailQueueRoutingKey, false, false, amqp091.Publishing{
		ContentType: "text/plan",
		Body:        []byte(emailMessageStr),
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Sent Email Message %s \n", emailMessageStr)

	time.Sleep(30 * time.Second)

	err = ch.Publish(ExchangeName, WhatsappQueueRoutingKey, false, false, amqp091.Publishing{
		ContentType: "text/plan",
		Body:        []byte(whatsappMessageStr),
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Sent Whatsapp Message %s \n", whatsappMessageStr)

	defer ch.Close()
}
