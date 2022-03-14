package main

import (
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErrorRpc(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnErrorRpc(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,
	)
	failOnErrorRpc(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnErrorRpc(err, "Failed to set QoS")

	msg, err := ch.Consume(
		q.Name,
		"",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	failOnErrorRpc(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msg {
			n, err := strconv.Atoi(string(d.Body))
			failOnErrorRpc(err, "Failed to convert body to integer")

			log.Printf(" [.] fib(%d)", n)
			response := fib(n)

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(response)),
				})
			failOnErrorRpc(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<- forever
}

func failOnErrorRpc(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
