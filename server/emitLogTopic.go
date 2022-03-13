package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErrorTopic(err, "Failed to connect to RabbitMQ")
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()
	ch, err := conn.Channel()
	failOnErrorTopic(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"logs_topic",   // name
		"topic", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnErrorTopic(err, "Failed to declare an exchange")

	body := bodyFromTopic(os.Args)

	err = ch.Publish("logs_topic", severityTopicFrom(os.Args), false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	failOnErrorTopic(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)

}

func failOnErrorTopic(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyFromTopic(args []string) string {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func severityTopicFrom(args []string) string {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}