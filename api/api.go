package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"
)

const (
	host     = "database"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "online-judge"
)

func StartAPI(local bool) {
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
	failOnError(err, "Failed to publish a message")

	var psqlInfo string
	if local == true {
		fmt.Printf("Running in LOCAL mode, connecting to localhost...\n")
		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			"localhost", port, user, password, dbname)

	} else {
		time.Sleep(5 * time.Second)
		fmt.Printf("Running in PRODUCTION mode, connecting to %v...\n", host)
		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Trying to connect to the database...")
	err = db.Ping()
	if err != nil {
		fmt.Println("Could not connect!")
		panic(err)
	} else {
		fmt.Println("Successfully connected!")
	}

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./frontend/public", true)))
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, &gin.H{
			"health": "good",
		})
	})
	router.GET("/api/questions", func(c *gin.Context) {
		getQuestionList(c, db)
	})
	router.GET("/api/question/:key", func(c *gin.Context) {
		getQuestion(c, db)
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
