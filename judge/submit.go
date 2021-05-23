package judge

import (
	"fmt"
	"github.com/laupski/online-judge/api"
	"github.com/laupski/online-judge/internal"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func compileSubmission(body []byte) api.SubmissionResponse {
	fmt.Println("Received message from Queue")
	submission, err := api.CheckJSONSubmissionRequest(body)
	if err != nil {
		return api.NewSubmissionResponse("", "", err.Error())
	}

	tmpfile, err := ioutil.TempFile(sandboxDirectory, "judge*.go")
	if err != nil {
		return api.NewSubmissionResponse("", "", err.Error())
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	_, err = tmpfile.Write([]byte(submission.Code))
	if err != nil {
		return api.NewSubmissionResponse("", "", err.Error())
	}

	// Create a timeout to prevent hanging resources
	c1 := make(chan api.SubmissionResponse, 1)
	go func() {
		run := exec.Command("go", "run", tmpfile.Name())
		stdout, _ := run.StdoutPipe()
		stderr, _ := run.StderrPipe()
		err = run.Start()
		standardOutput, _ := ioutil.ReadAll(stdout)
		errorOutput, _ := ioutil.ReadAll(stderr)
		c1 <- api.NewSubmissionResponse(string(standardOutput), string(errorOutput), "")
	}()
	select {
	case res := <-c1:
		return res
	case <-time.After(internal.Timeout * time.Second):
		return api.NewSubmissionResponse("", "", fmt.Sprintf("timeout after %v seconds", internal.Timeout))
	}
}
