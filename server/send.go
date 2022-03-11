package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %s", err)
		}
	}()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func() {
		if err := ch.Close(); err != nil {
			log.Fatalf("Failed to close channel: %s", err)
		}
	}()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to publish a message")
	for i := 0; i < 10; i++ {
		body := fmt.Sprintf("Message %d", i)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
	}

	//err = ch.Publish(
	//	"",
	//	q.Name,
	//	false,
	//	false,
	//	amqp.Publishing{
	//		ContentType: "text/plain",
	//		Body:        []byte("Hello World!"),
	//	})
}
