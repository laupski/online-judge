package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laupski/online-judge/internal"
	"github.com/streadway/amqp"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type SubmissionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type SubmissionResponse struct {
	StdOut  string `json:"stdout"`
	StdErr  string `json:"stderr"`
	Correct bool   `json:"correct"`
	Error   string `json:"error"`
}

var supportedLanguages = []string{"go"}

func postSubmission(c *gin.Context) {
	fmt.Println("Routing submission...")
	_ = c.Param("question")
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	_, err := CheckJSONSubmissionRequest(bodyBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewSubmissionResponse("", "", err.Error()))
		return
	}

	corrId := randomString(32)
	err = RabbitMQ.Channel.Publish(
		"submissions", // exchange
		"requests",    // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			//ReplyTo:       "requests",
			Body: bodyBytes,
		})
	internal.FailOnError(err, "Failed to publish a message")

	// Create a timeout to prevent hanging resources
	c1 := make(chan bool, 1)
	go func() {
		for d := range CompiledResults {
			if corrId == d.CorrelationId {
				fmt.Println("Successfully received an answer from the Judge server!")
				//fmt.Println(d.Body)
				var response SubmissionResponse
				_ = json.Unmarshal(d.Body, &response)
				c.JSON(http.StatusOK, response)
				c1 <- true
			}
		}
	}()
	select {
	case <-c1:
		return
	case <-time.After(internal.Timeout * time.Second):
		c.JSON(http.StatusInternalServerError, SubmissionResponse{
			StdOut:  "",
			StdErr:  "",
			Correct: false,
			Error:   "timeout"})
		return
	}
}

func CheckJSONSubmissionRequest(bodyBytes []byte) (SubmissionRequest, error) {
	var submission SubmissionRequest
	err := json.Unmarshal(bodyBytes, &submission)
	if err != nil {
		return submission, errors.New("bad submission")
	}
	if submission.Language == "" {
		return submission, errors.New("language cannot be empty")
	}
	if submission.Code == "" {
		return submission, errors.New("submission cannot be empty")
	}

	found := false
	for _, v := range supportedLanguages {
		if v == submission.Language {
			found = true
		}
	}
	if found == false {
		return submission, errors.New("unsupported language requested")
	}

	return submission, nil
}

func NewSubmissionResponse(stdout string, stderr string, err string) SubmissionResponse {
	return SubmissionResponse{
		stdout,
		stderr,
		false,
		err,
	}
}

// Check redis then postgres for the answer
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
