package internal

import (
	"fmt"
	"github.com/streadway/amqp"
)

const rabbitConnectionString = "amqp://guest:guest@messaging:5672/"

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

func NewRabbitMQ(local bool) RabbitMQ {
	fmt.Println("Attempting to connect to rabbitmq")
	var connectionString string
	if local {
		fmt.Printf("Running in LOCAL mode, connecting to localhost...\n")
		connectionString = "amqp://guest:guest@localhost:5672/"
	} else {
		fmt.Printf("Running in PRODUCTION mode, connecting to messaging...\n")
		connectionString = rabbitConnectionString
	}

	conn, err := amqp.Dial(connectionString)
	FailOnError(err, "Failed to connect to RabbitMQ")
	fmt.Println("Successfully connected!")
	return RabbitMQ{
		Connection: conn,
		Channel:    nil,
		Queue: amqp.Queue{
			Name:      "",
			Messages:  0,
			Consumers: 0,
		},
	}
}

func (rmq *RabbitMQ) CreateSubmissionChannel() {
	channel, err := rmq.Connection.Channel()
	FailOnError(err, "Failed to open a channel")
	rmq.Channel = channel
}

func (rmq *RabbitMQ) DeclareQueue() {
	var err error
	rmq.Queue, err = rmq.Channel.QueueDeclare(
		"submissions", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	FailOnError(err, "Failed to declare a queue")
}
