package judge

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

const sandboxDirectory = "./sandbox/"

func StartJudge(local bool) {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, &gin.H{
			"health": "good",
		})
	})
	router.POST("/judge", func(c *gin.Context) {
		postSubmission(c)
	})

	if _, err := os.Stat(sandboxDirectory); os.IsNotExist(err) {
		fmt.Printf("Creating sandbox directory in %v\n", sandboxDirectory)
		err = os.Mkdir(sandboxDirectory, 0777)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Sandbox directory already created in %v\n", sandboxDirectory)
	}

	judge := &http.Server{
		Handler: router,
		Addr:    ":1338",
	}

	fmt.Println("Now serving the online-judge Judge server on port 1338")
	log.Fatal(judge.ListenAndServe())
}
