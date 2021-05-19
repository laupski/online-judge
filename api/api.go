package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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
	router.GET("api/questions", func(c *gin.Context) {
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
