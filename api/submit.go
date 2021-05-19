package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func postSubmission(c *gin.Context) {
	fmt.Println("Running submission...")
	c.Param("question")
}

func checkAnswer(key, output string) bool {
	return false
}
