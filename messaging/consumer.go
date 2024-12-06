package messaging

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handleReplicaAdded(body interface{}) {
	data, ok := body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for replica-added: %v", body)
		return
	}

	replicaName, _ := data["name"].(string)
	replicaURL, _ := data["url"].(string)

	log.Printf("Replica added: Name=%s, URL=%s", replicaName, replicaURL)

}

func handleReplicaRemoved(body interface{}) {
	data, ok := body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid message body for replica-removed: %v", body)
		return
	}

	replicaName, _ := data["name"].(string)
	replicaURL, _ := data["url"].(string)
	
	log.Printf("Replica removed: Name=%s, URL=%s", replicaName, replicaURL)

}

func SetupConsumer() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"reverseproxy-to-admin", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			processMessage(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
