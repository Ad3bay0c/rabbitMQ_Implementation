package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErrorTopic(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnErrorTopic(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnErrorTopic(err, "Failed to declare an exchange")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	failOnErrorTopic(err, "Failed to declare a queue")

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", s, "logs_topic", s)
		err = ch.QueueBind(
			q.Name,        // queue name
			s,             // routing key
			"logs_topic", // exchange
			false,
			nil,
		)
		failOnErrorTopic(err, "Failed to bind a queue")
	}

	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func failOnErrorTopic(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
