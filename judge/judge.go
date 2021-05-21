package judge

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

const sandboxDirectory = "./sandbox/"

func StartJudge(local bool) {
	if _, err := os.Stat(sandboxDirectory); os.IsNotExist(err) {
		fmt.Printf("Creating sandbox directory in %v\n", sandboxDirectory)
		err = os.Mkdir(sandboxDirectory, 0777)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Sandbox directory already created in %v\n", sandboxDirectory)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

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
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
