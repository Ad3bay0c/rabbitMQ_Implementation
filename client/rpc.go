package main

import (
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	n := bodyFrom(os.Args)

	log.Printf(" [x] Requesting fib(%d)", n)
	res, err := fibonacciRPC(n)
	failOnErrorRpc(err, "Failed to handle RPC request")

	log.Printf(" [.] Got %d", res)
}

func bodyFrom(args []string) int {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	failOnErrorRpc(err, "Failed to convert arg to integer")
	return n
}

func fibonacciRPC(n int) (res int, err error) {
	conn,err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErrorRpc(err, "Failed to connect to RabbitMQ")
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	ch, err := conn.Channel()
	failOnErrorRpc(err, "Failed to open a channel")
	defer func() {
		if err != nil {
			ch.Close()
		}
	}()

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnErrorRpc(err, "Failed to declare a queue")

	msq, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnErrorRpc(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = ch.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(strconv.Itoa(n)),
			CorrelationId: corrId,
			ReplyTo:       q.Name,
		},
	)
	failOnErrorRpc(err, "Failed to publish a message")

	for d := range msq {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			failOnErrorRpc(err, "Failed to convert body to integer")
			break
		}
	}
	return
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func failOnErrorRpc(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
