package webhookcicd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type BranchUpdate struct {
	Name string
	SHA  string
}
type pipeline struct {
	file   string
	branch chan BranchUpdate
	db     *DB
}

func newPipeline(db *DB, file string) (error, *pipeline) {

	if len(file) < 0 {
		return errors.New("no pipeline script/file provided"), nil
	}
	b := make(chan BranchUpdate)
	return nil, &pipeline{file: file, branch: b, db: db}

}

func (pp *pipeline) Run() {
	ver := pp.db.GetBuildNo()

	//handleError(pp.db.Put([]byte("buildNO"), buildNo, nil))

	log.Println("Build Started")

	src, err := ioutil.ReadFile(pp.file)
	if err != nil {
		log.Fatal(err)
	}

	src = bytes.ReplaceAll(src, []byte("${BUILD_NO}"), []byte(ver))

	if err = ioutil.WriteFile("process.lock", src, 0666); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("/bin/sh", "process.lock")
	cmd.Stdout = log.Writer()
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	// deleteFile("process.lock")
	_ = pp.db.BuildFinish()
	handleError(err)

}

func (pp *pipeline) Monit() {
	for {
		select {
		case bra := <-pp.branch:

			ver := pp.db.GetBuildNo()
			_ = pp.db.BuildFinish()
			log.Println(bra.Name, ver)

			if err := prepareTheSource(); err != nil {
				cleanTheSource()
				log.Println(err)
			}

			buildTheImage()
			pushTheImage()

		}
	}
}

func prepareTheSource() error {
	log.Println(" Cloning the repositiory")
	if err := runCmd("rm -rf *"); err != nil {
		handleError(err)
		return errors.New("repository clone failed, temp directory error")
	}

	if err := runCmd("git clone git@github.com:grapetechadmin/dewa-test.git"); err != nil {
		handleError(err)
		return errors.New("clone from github failed with error :" + err.Error())
	}

	return nil

}

func buildTheImage() {
	log.Println("Build the Image")
	fmt.Println("docker build -t [dewa-ev]:latest .")
}

func pushTheImage() {
	log.Println("Push the Image")
	fmt.Println("docker tag [dewa-ev]:latest 670907057868.dkr.ecr.us-east-2.amazonaws.com/dewa-ev:${BUILD_NO}")
	fmt.Println("docker push 670907057868.dkr.ecr.us-east-2.amazonaws.com/[dewa-ev]:${BUILD_NO}")
}

func cleanTheSource() {
	log.Println("Clear the work dir")
	err := runCmd("rm -rf *")
	if err != nil {
		handleError(err)
	}
}

func runCmd(command string) error {
	return runCmdSetDir(command, WorkDir)
}

func runCmdSetDir(command string, dir string) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = dir
	cmd.Stdout = log.Writer()
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
