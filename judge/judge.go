package judge

import (
	"encoding/json"
	"fmt"
	"github.com/laupski/online-judge/internal"
	"github.com/streadway/amqp"
	"log"
	"os"
)

const sandboxDirectory = "./sandbox/"

var RabbitMQ internal.RabbitMQ
var Deliveries <-chan amqp.Delivery

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
	RabbitMQ.DeclareAndBindQueue("responses")
	Deliveries = RabbitMQ.SetConsumer("requests")

	//TODO figure out why the channel keeps closing after processing
	forever := make(chan bool)
	go func() {
		for d := range Deliveries {
			log.Printf("Received a message: %s", d.Body)
			response := compileSubmission(d.Body)
			bodyBytes, _ := json.Marshal(response)

			err := RabbitMQ.Channel.Publish(
				"submissions", // exchange
				"responses",   // routing key
				false,         // mandatory
				false,         // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          bodyBytes,
				})
			internal.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
