package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MainExchange  = "notification_exchange"
	MainQueue     = "email_queue"
	RetryExchange = "retry_exchange"
	RetryQueue    = "email_retry_queue"
	RoutingKey    = "email_routing_key"
)

func DLX() {
	// 1. Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	// 2. Setup Infrastructure (Exchanges and Queues)
	setupInfrastructure(ch)

	// 3. Start Consumer in a goroutine
	go Consumer(ch)

	time.Sleep(5 * time.Second)

	// 4. Run Producer to send one test message
	Producer(ch)

	// Keep main alive to watch the retries happen
	select {}
}

func Producer(ch *amqp.Channel) {
	body := "Test Email Content"
	err := ch.Publish(
		MainExchange,
		RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Printf("Failed to publish: %v", err)
	}
	fmt.Println(" [→] Producer sent message")
}

func Consumer(ch *amqp.Channel) {
	// Important: Don't take more than 1 message at a time
	ch.Qos(1, 0, false)

	msgs, _ := ch.Consume(MainQueue, "", false, false, false, false, nil)

	for d := range msgs {
		retryCount := getRetryCount(d.Headers)
		fmt.Printf("\n [←] Received. Body: %s | Failure Count: %d\n", string(d.Body), retryCount)

		// SIMULATION: Fail until retry count is 3
		if retryCount < 3 {
			fmt.Printf(" [!] Simulating Error... Sending to Delay Queue for 5s\n")
			// requeue: false triggers the Dead Letter Exchange
			d.Nack(false, false)
		} else {
			fmt.Println(" [✔] Success after 3 retries! Acknowledging.")
			d.Ack(false)
		}
	}
}

func setupInfrastructure(ch *amqp.Channel) {
	// Declare Exchanges
	ch.ExchangeDeclare(MainExchange, "direct", true, false, false, false, nil)
	ch.ExchangeDeclare(RetryExchange, "direct", true, false, false, false, nil)

	ch.QueueDeclare(MainQueue, true, false, false, false, nil)
	ch.QueueBind(MainQueue, RoutingKey, MainExchange, false, nil)

	ch.QueueDeclare(RetryQueue, true, false, false, false, nil)
	ch.QueueBind(RetryQueue, RoutingKey, RetryExchange, false, nil)

	err := UpsertPolicy("main_retry_policy", "^email_queue$", PolicyDefinition{
		"dead-letter-exchange": RetryExchange,
	})
	if err != nil {
		fmt.Printf("Error setting main policy: %v\n", err)
	}

	// 2. Sync Policy for the Retry Queue (TTL and return to Main)
	err = UpsertPolicy("retry_delay_policy", "^email_retry_queue$", PolicyDefinition{
		"dead-letter-exchange":    MainExchange,
		"dead-letter-routing-key": RoutingKey,
		"message-ttl":             5000, // 5 seconds
	})
	if err != nil {
		fmt.Printf("Error setting retry policy: %v\n", err)
	}
}

func getRetryCount(headers amqp.Table) int64 {
	if val, ok := headers["x-death"]; ok {
		fmt.Println("val:", val)
		if slice, ok := val.([]interface{}); ok && len(slice) > 0 {
			if table, ok := slice[0].(amqp.Table); ok {
				return table["count"].(int64)
			}
		}
	}
	return 0
}
