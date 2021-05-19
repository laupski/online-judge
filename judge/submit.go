package judge

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type submissionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

func postSubmission(c *gin.Context) {
	var submission submissionRequest
	body := c.Request.Body
	x, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(x, &submission)
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

	fmt.Println("Compiling code...")
	err = compileSubmission(submission)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
		return
	}

	fmt.Println("Running code...")
	results, err := runSubmission(submission)
	c.JSON(http.StatusOK, gin.H{
		"error":   err,
		"results": results,
	})
}

func compileSubmission(s submissionRequest) error {
	return nil
}

func runSubmission(s submissionRequest) (string, error) {
	return "", nil
}
