package webhookcicd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func handleError(err error) {
	if err != nil {
		log.Printf("[error] err=%s\n", err, err)
		return
	}
}

func handleErrorMsg(tag string, err error) {
	if err != nil {
		log.Printf("[%s] err=%s\n", tag, err)
		return
	}
}

func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	handleErrorMsg("Delete FIle: ", err)

	fmt.Println("File Deleted")
}

func runCmd(command string) error {
	return runCmdSetDir(command, WorkDir)
}

func runCmdSetDir(command string, dir string) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = dir
	handleErrorMsg("[runCmdSetDir]", os.MkdirAll(os.TempDir()+"cicd-logs", os.ModePerm))
	outfile, err := os.Create(os.TempDir() + "cicd-logs/log.log")
	if err != nil {
		return err
	}
	defer outfile.Close()
	errFile, err := os.Create(os.TempDir() + "cicd-logs/error.log")
	if err != nil {
		return err
	}
	defer errFile.Close()
	cmd.Stdout = outfile
	//cmd.Stdout = log.Writer()
	cmd.Stderr = errFile
	return cmd.Run()
}
