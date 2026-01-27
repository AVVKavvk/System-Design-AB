package topic

import (
	"fmt"
	"time"
)

var (
	ExchangeName            = "notification_exchange_topic"
	EmailQueueName          = "email_queue_topic"
	WhatsappQueueName       = "whatsapp_queue_topic"
	EmailQueueRoutingKey    = "email.*"
	WhatsappQueueRoutingKey = "whatsapp.*"
)

func safeRun(name string, task func()) {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("%s crashed: %v. Restarting in 5s...\n", name, r)
					time.Sleep(5 * time.Second) // Cooldown to prevent CPU spikes
				}
			}()
			task()
		}()
	}
}

func main() {
	go safeRun("Producer", Producer)
	go safeRun("EmailConsumer", EmailConsumer)
	go safeRun("WhatsappConsumer", WhatsappConsumer)

	// Block main from exiting
	select {}
}
func Main() {

	main()

}
