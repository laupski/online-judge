package judge

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func StartJudge(local bool) {
	router := gin.Default()

	router.POST("/judge", func(c *gin.Context) {
		postSubmission(c)
	})

	judge := &http.Server{
		Handler: router,
		Addr:    ":1338",
	}

	fmt.Println("Now serving the online-judge Judge server on port 1338")
	log.Fatal(judge.ListenAndServe())
}
