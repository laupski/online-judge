package judge

import (
	"fmt"
	"github.com/laupski/online-judge/internal"
	"log"
	"os"
)

const sandboxDirectory = "./sandbox/"

var RabbitMQConnection internal.RabbitMQ

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
	RabbitMQConnection = internal.NewRabbitMQ(local)
	defer RabbitMQConnection.Connection.Close()
	RabbitMQConnection.CreateSubmissionChannel()
	defer RabbitMQConnection.Channel.Close()
	RabbitMQConnection.DeclareQueue()

	msgs, err := RabbitMQConnection.Channel.Consume(
		RabbitMQConnection.Queue.Name, // queue
		"",                            // consumer
		true,                          // auto-ack
		false,                         // exclusive
		false,                         // no-local
		false,                         // no-wait
		nil,                           // args
	)
	internal.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
