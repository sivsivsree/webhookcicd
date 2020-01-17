package webhookcicd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type pipeline struct {
	file string
}

func newPipeline(file string) (error, *pipeline) {

	if len(file) < 0 {
		return errors.New("no pipeline script/file provided"), nil
	}


	return nil, &pipeline{file: file,}

}

func (pp pipeline) Run() {
	fmt.Println("Running Service")
	cmd := exec.Command("/bin/sh", pp.file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	handleError(err)
}
