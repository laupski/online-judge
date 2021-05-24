package judge

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/laupski/online-judge/internal"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

const sandboxDirectory = "./sandbox/"

var RabbitMQ internal.RabbitMQ
var Deliveries <-chan amqp.Delivery
var Redis *redis.Client

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
	RabbitMQ.DeclareAndBindQueue("requests")
	Deliveries = RabbitMQ.SetConsumer("requests")
	Redis = internal.NewRedis(local)
	defer Redis.Close()

	forever := make(chan bool)
	go func() {
		for d := range Deliveries {
			log.Printf("Received a message: %s", d.Body)
			response := compileSubmission(d.Body)
			jsonResponse, _ := json.Marshal(response)
			err := d.Ack(false)
			if err != nil {
				fmt.Println(err)
			}

			err = Redis.Set(d.CorrelationId, jsonResponse, 2*internal.Timeout*time.Second).Err()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
