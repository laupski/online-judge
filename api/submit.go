package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func postSubmission(c *gin.Context) {
	fmt.Println("Posting submission...")
}

func checkAnswer(key, output string) bool {
	return false
}
