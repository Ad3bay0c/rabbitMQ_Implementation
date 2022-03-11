package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func failOnError1(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError1(err, "Failed to connect to RabbitMQ")
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %s", err)
		}
	}()
	ch, err := conn.Channel()
	failOnError1(err, "Failed to open a channel")
	defer func() {
		if err := ch.Close(); err != nil {
			log.Fatalf("Failed to close channel: %s", err)
		}
	}()

	q, err := ch.QueueDeclare(
		"new_task", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError1(err, "Failed to publish a message")

	body := bodyFrom(os.Args)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError1(err, "Failed to publish a message")

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

func bodyFrom(args []string) string {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
