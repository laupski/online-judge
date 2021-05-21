package judge

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

const (
	sandboxDirectory = "./sandbox/"
	rabbit           = "amqp://guest:guest@messaging:5672/"
)

var RabbitMQConnection *amqp.Connection

func StartJudge(local bool) {
	RabbitMQConnection = connectToRabbitMQ(local)
	defer RabbitMQConnection.Close()

	if _, err := os.Stat(sandboxDirectory); os.IsNotExist(err) {
		fmt.Printf("Creating sandbox directory in %v\n", sandboxDirectory)
		err = os.Mkdir(sandboxDirectory, 0777)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Sandbox directory already created in %v\n", sandboxDirectory)
	}

	ch, err := RabbitMQConnection.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
}

func connectToRabbitMQ(local bool) *amqp.Connection {
	fmt.Println("Attempting to connect to rabbitmq")
	var connectionString string
	if local {
		fmt.Printf("Running in LOCAL mode, connecting to localhost...\n")
		connectionString = "amqp://guest:guest@localhost:5672/"
	} else {
		fmt.Printf("Running in PRODUCTION mode, connecting to messaging...\n")
		connectionString = rabbit
	}

	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	fmt.Println("Successfully connected!")
	return conn

	/* TODO implement messaging
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever*/
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
