package judge

import (
	"fmt"
	"github.com/laupski/online-judge/internal"
	"github.com/streadway/amqp"
	"log"
	"os"
)

const sandboxDirectory = "./sandbox/"

var RabbitMQ internal.RabbitMQ

func init() {
	if _, err := os.Stat(sandboxDirectory); os.IsNotExist(err) {
		fmt.Printf("Creating sandbox directory in %v\n", sandboxDirectory)
		err = os.Mkdir(sandboxDirectory, 0777)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Sandbox directory already created in %v\n", sandboxDirectory)
	}
}

func StartJudge(local bool) {
	RabbitMQ = internal.NewRabbitMQ(local)
	defer RabbitMQ.Connection.Close()
	RabbitMQ.CreateSubmissionChannel()
	defer RabbitMQ.Channel.Close()
	RabbitMQ.DeclareQueue()

	msgs, err := RabbitMQ.Channel.Consume(
		RabbitMQ.Queue.Name, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	internal.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			compileSubmission(d.Body)

			err = RabbitMQ.Channel.Publish(
				"",                  // exchange
				RabbitMQ.Queue.Name, // routing key
				false,               // mandatory
				false,               // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte("success"),
				})
			internal.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
