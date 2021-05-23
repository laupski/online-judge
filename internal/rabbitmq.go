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
	err = channel.ExchangeDeclare(
		"submissions",
		"direct",
		false,
		false,
		false,
		false,
		nil)
	FailOnError(err, "Failed to declare an exchange on the channel")
	rmq.Channel = channel
}

func (rmq *RabbitMQ) DeclareAndBindQueue(n string) {
	var err error
	rmq.Queue, err = rmq.Channel.QueueDeclare(
		n,     // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare a queue")
	err = rmq.Channel.QueueBind(
		rmq.Queue.Name,
		n,
		"submissions",
		false,
		nil)
	FailOnError(err, "Failed to bind the queue")
}

func (rmq *RabbitMQ) SetConsumer(k string) <-chan amqp.Delivery {
	channel, err := rmq.Channel.Consume(
		k,
		"",
		true,
		false,
		false,
		false,
		nil)
	FailOnError(err, "Failed to register a consumer")
	return channel
}
