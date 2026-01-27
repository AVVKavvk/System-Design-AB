package direct

import (
	"fmt"
	"time"
)

var (
	ExchangeName            = "notification_exchange"
	EmailQueueName          = "email_queue"
	WhatsappQueueName       = "whatsapp_queue"
	EmailQueueRoutingKey    = "email_routing_key"
	WhatsappQueueRoutingKey = "whatsapp_routing_key"
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
