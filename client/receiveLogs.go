package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError2(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError2(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError2(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	failOnError2(err, "Failed to declare a queue")

	err = ch.QueueBind(q.Name, "", "logs", false, nil)
	failOnError2(err, "Failed to bind a queue")

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

func failOnError2(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
