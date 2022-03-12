package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError3(err, "Failed to connect to RabbitMQ")
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()
	ch, err := conn.Channel()
	failOnError3(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"logs_direct",   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError3(err, "Failed to declare an exchange")

	body := bodyFrom3(os.Args)

	err = ch.Publish("logs_direct", severityFrom(os.Args), false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	failOnError3(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)

}

func failOnError3(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyFrom3(args []string) string {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}