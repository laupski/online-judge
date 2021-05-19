package api

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func StartAPI() {
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./frontend/public", true)))
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

func getQuestion(c *gin.Context) {
	fmt.Println("Getting question")
}

func postSubmission(c *gin.Context) {
	fmt.Println("Posting submission")
}

func checkAnswer(key, output string) bool {
	return false
}
