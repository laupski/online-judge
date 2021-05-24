package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type SubmissionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type SubmissionResponse struct {
	Submitted    bool   `json:"submitted"`
	SubmissionID string `json:"submission_id"`
	Error        string `json:"error"`
}

var supportedLanguages = []string{"go"}

func postSubmission(c *gin.Context) {
	fmt.Println("Submitting code...")
	//_ = c.Param("question")
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	_, err := CheckJSONSubmissionRequest(bodyBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewSubmitResponse(false, "", err.Error()))
		return
	}

	submissionId := randomString(32)
	err = RabbitMQ.PublishMessage(submissionId, bodyBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewSubmitResponse(false, "", "could not submit response"))
		return
	}

	c.JSON(http.StatusOK, NewSubmitResponse(true, submissionId, ""))
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

func NewSubmitResponse(submitted bool, submission_id string, error string) SubmissionResponse {
	return SubmissionResponse{
		submitted,
		submission_id,
		error,
	}
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
