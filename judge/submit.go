package judge

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const timeout = 5

var supportedLanguages = []string{"go"}

type outputs struct {
	stdout string
	stderr string
}

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

	fmt.Println("Running code...")
	stdout, stderr, err := compileSubmission(submission)

	c.JSON(http.StatusOK, gin.H{
		"stdout": stdout,
		"stderr": stderr,
		"error":  err,
	})

}

func compileSubmission(s submissionRequest) (string, string, error) {
	found := false
	for _, v := range supportedLanguages {
		if v == s.Language {
			found = true
		}
	}
	if found == false {
		return "", "", errors.New("unsupported language requested")
	}

	tmpfile, err := ioutil.TempFile(sandboxDirectory, "judge*.go")
	if err != nil {
		return "", "", err
	}
	defer os.Remove(tmpfile.Name())
	_, err = tmpfile.Write([]byte(s.Code))
	if err != nil {
		return "", "", err
	}

	// Create a timeout to prevent hanging resources
	c1 := make(chan outputs, 1)
	go func() {
		run := exec.Command("go", "run", tmpfile.Name())
		stdout, _ := run.StdoutPipe()
		stderr, _ := run.StderrPipe()
		err = run.Start()
		standardOutput, _ := ioutil.ReadAll(stdout)
		errorOutput, _ := ioutil.ReadAll(stderr)
		c1 <- outputs{
			stdout: string(standardOutput),
			stderr: string(errorOutput),
		}
	}()
	select {
	case res := <-c1:
		tmpfile.Close()
		return res.stdout, res.stderr, err
	case <-time.After(timeout * time.Second):
		tmpfile.Close()
		return "", "", errors.New(fmt.Sprintf("timeout after %v seconds", timeout))
	}
}
