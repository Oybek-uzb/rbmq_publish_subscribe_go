package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err, "failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	failOnError(err, "failed to declare queue")

	err = ch.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil)
	failOnError(err, "failed to bind queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err, "failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			log.Printf("Consumed message: %s", msg.Body)
		}
	}()

	log.Print(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
