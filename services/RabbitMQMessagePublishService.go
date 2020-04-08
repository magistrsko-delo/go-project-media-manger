package services

import (
	"github.com/streadway/amqp"
	"go-project-media-manger/Models"
	"log"
)

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
	q    amqp.Queue
}

func InitRabbitMQ() *RabbitMQ  {

	env := Models.GetEnvStruct()
	conn, err := amqp.Dial("amqp://" + env.RabbitUser + ":" + env.RabbitPassword + "@" + env.RabbitHost + ":5672")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		env.RabbitQueue, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return &RabbitMQ{
		Conn: conn,
		Ch:   ch,
		q:    q,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}