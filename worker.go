package webhookcicd

import (
	"errors"
	"log"
	"strconv"
)

type BranchUpdate struct {
	Name string
	SHA  string
}
type pipeline struct {
	branch chan BranchUpdate
	db     *DB
}

func newPipeline(db *DB) (error, *pipeline) {

	b := make(chan BranchUpdate)
	return nil, &pipeline{branch: b, db: db}

}

func (pp *pipeline) StartWorker() {
	for {
		select {
		case bra := <-pp.branch:

			ver := pp.db.GetBuildNo()
			repoName := pp.db.GetRepoName()
			if err := pp.db.BuildFinish(); err != nil {
				log.Println(err)
			}
			log.Println(bra.Name, ver)

			if err := prepareTheSource(); err != nil {
				cleanTheSource()
				log.Println(err)
				return
			}

			if err := buildTheImage(repoName); err != nil {
				cleanTheSource()
				log.Println(err)
				return
			}

			if err := pushTheImage(repoName, ver); err != nil {
				cleanTheSource()
				log.Println(err)
				return
			}

			cleanTheSource()

		}
	}
}

func prepareTheSource() error {
	log.Println(" 🚀  Cloning the repositiory")
	if err := runCmd("rm -rf *"); err != nil {
		return errors.New(" 👾  repository clone failed, temp directory error")
	}

	if err := runCmd("git clone git@github.com:grapetechadmin/dewa-test.git"); err != nil {
		return errors.New(" 👾  clone from github failed with error :" + err.Error())
	}

	return nil

}

func buildTheImage(repoName string) error {
	log.Println(" 🚀  Building the repository to " + repoName + " latest")
	if err := runCmdSetDir("docker build -t "+repoName+":latest .", WorkDir+"/"+repoName); err != nil {
		handleError(err)
		return errors.New(" 👾  docker build failed")
	}
	return nil
}

func pushTheImage(repoName string, buildNo int) error {
	log.Println(" 🚀  Push " + repoName + ":latest")

	if err := runCmd("docker tag " + repoName + ":latest 670907057868.dkr.ecr.us-east-2.amazonaws.com/" + repoName + ":" + strconv.Itoa(buildNo)); err != nil {
		handleError(err)
		return errors.New(" 👾  docker tag failed, ")
	}

	log.Println(" 🚀  Tagged " + repoName + ":latest as  670907057868.dkr.ecr.us-east-2.amazonaws.com/" + repoName + ":" + strconv.Itoa(buildNo))

	if err := runCmd("docker push 670907057868.dkr.ecr.us-east-2.amazonaws.com/" + repoName + ":" + strconv.Itoa(buildNo)); err != nil {
		handleError(err)
		return errors.New(" 👾  docker push failed")
	}

	log.Println(" 🚀  Pushed to registry 670907057868.dkr.ecr.us-east-2.amazonaws.com/" + repoName + ":" + strconv.Itoa(buildNo))
	return nil
}

func cleanTheSource() {
	log.Println(" 🚀  Clear the work dir")

	if err := runCmd("rm -rf *"); err != nil {
		handleError(err)
	}

	if err := runCmd("docker rmi -f $(docker images -aq)"); err != nil {
		handleError(err)
	}
}
