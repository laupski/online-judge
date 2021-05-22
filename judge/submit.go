package judge

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	Language string
	Code     string
}

func compileSubmission(body []byte) (string, string, error) {
	fmt.Println("Received message from Queue")
	var submission submissionRequest
	err := json.Unmarshal(body, &submission)
	if err != nil {
		return "", "", errors.New("Bad submission")
	}

	if submission.Language == "" {
		return "", "", errors.New("Language cannot be empty")
	}
	if submission.Code == "" {
		return "", "", errors.New("Submission cannot be empty")
	}

	found := false
	for _, v := range supportedLanguages {
		if v == submission.Language {
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
	defer tmpfile.Close()

	_, err = tmpfile.Write([]byte(submission.Code))
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
		return res.stdout, res.stderr, err
	case <-time.After(timeout * time.Second):
		return "", "", errors.New(fmt.Sprintf("timeout after %v seconds", timeout))
	}
}
