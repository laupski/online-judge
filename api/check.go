package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// route the question key
// check the answer with the correlation ID in the payload

type CheckRequest struct {
	SubmissionID string `json:"submission_id"`
}

type CheckResponse struct {
	StdOut  string `json:"stdout"`
	StdErr  string `json:"stderr"`
	Correct bool   `json:"correct"`
	Error   string `json:"error"`
}

func checkSubmission(c *gin.Context) {
	//_ = c.Param("question")
	var check CheckRequest
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(bodyBytes, &check)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewCheckResponse("", "", "bad request"))
		return
	}
	if check.SubmissionID == "" {
		c.JSON(http.StatusBadRequest, NewCheckResponse("", "", "submission_id cannot be empty"))
		return
	}

	submission, err := Redis.Get(check.SubmissionID).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, NewCheckResponse("", "", "no submission found"))
		return
	}

	var response CheckResponse
	err = json.Unmarshal([]byte(submission), &response)
	c.JSON(http.StatusOK, response)
	//c.Res(http.StatusOK, val)
}

func NewCheckResponse(stdout string, stderr string, err string) CheckResponse {
	return CheckResponse{
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
