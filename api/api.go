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

func StartAPI() {
	time.Sleep(5 * time.Second)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

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

func getQuestion(c *gin.Context, db *sql.DB) {
	fmt.Println("Getting question...")
	var question string
	err := db.QueryRow("SELECT question FROM public.questions WHERE id = ?1", c.Request.Header.Get("id")).Scan(&question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching list",
		})
		return
	}
	c.JSONP(http.StatusOK, question)
}

func getQuestionList(c *gin.Context, db *sql.DB) {
	fmt.Println("Getting question list...")
	rows, err := db.Query("SELECT key, number, title FROM public.questions ORDER BY number")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching list",
		})
		return
	}

	for rows.Next() {
		var key string
		var number int
		var title string
		err = rows.Scan(&key, &number, &title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error parsing rows",
			})
			return
		}
		fmt.Printf("key: %v, number: %v, title: %v", key, number, title)
	}
	c.JSON(http.StatusOK, rows)
}

func postSubmission(c *gin.Context) {
	fmt.Println("Posting submission...")
}

func checkAnswer(key, output string) bool {
	return false
}
