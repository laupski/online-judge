package api

import (
	"context"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/laupski/online-judge/internal"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const postgres = "postgres://postgres:postgres@database:5432/online-judge"

var PostgresConnection *pgx.Conn
var RabbitMQ internal.RabbitMQ
var Msgs <-chan amqp.Delivery

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func StartAPI(local bool) {
	var err error
	PostgresConnection = connectToPostgres(local)
	defer PostgresConnection.Close(context.Background())
	RabbitMQ = internal.NewRabbitMQ(local)
	defer RabbitMQ.Connection.Close()
	RabbitMQ.CreateSubmissionChannel()
	defer RabbitMQ.Channel.Close()
	RabbitMQ.DeclareQueue()
	Msgs, err = RabbitMQ.Channel.Consume(
		RabbitMQ.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	internal.FailOnError(err, "Failed to register a consumer")

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
	internal.FailOnError(err, "Failed to connect to postgres")
	fmt.Println("Successfully connected!")
	return conn
}
