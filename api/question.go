package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getQuestionResponse struct {
	Question string `json:"question"`
}

func getQuestion(c *gin.Context) {
	fmt.Println("Getting question...")
	var question string
	err := PostgresConnection.QueryRow(context.Background(), "SELECT question FROM public.questions WHERE key = $1", c.Param("key")).Scan(&question)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No question found",
		})
		return
	}
	fmt.Println(question)
	response := &getQuestionResponse{
		Question: question,
	}
	c.JSON(http.StatusOK, response)
}
