package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laupski/online-judge/internal"
	"github.com/streadway/amqp"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type submissionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type submissionResponse struct {
	Response string `json:"submission"`
}

func postSubmission(c *gin.Context) {
	fmt.Println("Routing submission...")
	c.Param("question")
	var submission submissionRequest
	body := c.Request.Body
	bodyBytes, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(bodyBytes, &submission)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Bad submission",
		})
		return
	}

	if submission.Language == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Language cannot be empty",
		})
		return
	}
	if submission.Code == "" {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Submission cannot be empty",
		})
		return
	}

	corrId := randomString(32)

	err = RabbitMQ.Channel.Publish(
		"",                  // exchange
		RabbitMQ.Queue.Name, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       RabbitMQ.Queue.Name,
			Body:          bodyBytes,
		})
	internal.FailOnError(err, "Failed to publish a message")

	for d := range Msgs {
		if corrId == d.CorrelationId {
			fmt.Println("Successfully received an answer from the Judge server!")
			break
		}
	}

	c.JSON(http.StatusOK, submissionResponse{Response: "Successfully submitted"})
}

func checkAnswer(key, output string) bool {
	return false
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
