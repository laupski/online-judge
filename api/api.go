package api

import (
	"context"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"
)

const (
	postgres = "postgres://postgres:postgres@database:5432/online-judge"
	rabbit   = "amqp://guest:guest@messaging:5672/"
)

var PostgresConnection *pgx.Conn
var RabbitMQConnection *amqp.Connection

func StartAPI(local bool) {
	PostgresConnection = connectToPostgres(local)
	defer PostgresConnection.Close(context.Background())
	RabbitMQConnection = connectToRabbitMQ(local)
	defer RabbitMQConnection.Close()

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./frontend/public", true)))
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, &gin.H{
			"health": "good",
		})
	})
	router.GET("/api/questions", func(c *gin.Context) {
		getQuestionList(c)
	})
	router.GET("/api/question/:key", func(c *gin.Context) {
		getQuestion(c)
	})
	router.POST("/api/submit/:key", func(c *gin.Context) {
		postSubmission(c)
	})

	api := &http.Server{
		Handler: router,
		Addr:    ":1337",
	}

	fmt.Println("Now serving the online-judge API server on port 1337")
	log.Fatal(api.ListenAndServe())
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

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")*/
}

func connectToPostgres(local bool) *pgx.Conn {
	fmt.Println("Attempting to connect to postgres")
	var connectionString string
	if local {
		fmt.Printf("Running in LOCAL mode, connecting to localhost...\n")
		connectionString = "postgres://postgres:postgres@localhost:5432/online-judge"
	} else {
		time.Sleep(5 * time.Second)
		fmt.Printf("Running in PRODUCTION mode, connecting to database...\n")
		connectionString = postgres
	}

	conn, err := pgx.Connect(context.Background(), connectionString)
	failOnError(err, "Failed to connect to postgres")
	fmt.Println("Successfully connected!")
	return conn
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
