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
	branch       chan BranchUpdate
	db           *DB
	notification *slack
}

func newPipeline(db *DB) (error, *pipeline) {

	noti := NewSlack()
	b := make(chan BranchUpdate)
	return nil, &pipeline{branch: b, db: db, notification: noti}

}

func (pp *pipeline) StartWorker() {
	for {
		select {
		case bra := <-pp.branch:

			ver := pp.db.GetBuildNo()
			repoName := pp.db.GetRepoName()
			awsRegistry := pp.db.GetECR()

			go func() {
				pp.notification.msg <- Msg{
					Text:    " ⏳ Build # " + strconv.Itoa(ver) + " Started for " + bra.Name,
					BuildNo: ver,
				}
			}()

			if err := pp.db.BuildFinish(); err != nil {
				log.Println(err)
			}

			if err := prepareTheSource(repoName, bra.Name); err != nil {
				go func() {
					pp.notification.msg <- Msg{
						Text:    " 💥  Build # " + strconv.Itoa(ver) + " failed for " + bra.Name + "\n" + err.Error(),
						BuildNo: ver,
					}
				}()
				cleanTheSource()
				log.Println(err)
				return
			}

			if err := buildTheImage(repoName); err != nil {
				go func() {
					pp.notification.msg <- Msg{
						Text:    " 💥  Build # " + strconv.Itoa(ver) + " failed for " + bra.Name + "\n" + err.Error(),
						BuildNo: ver,
					}
				}()
				cleanTheSource()
				log.Println(err)
				return
			}

			buildVer := bra.Name + "-" + strconv.Itoa(ver)
			if err := pushTheImage(repoName, awsRegistry, buildVer); err != nil {
				go func() {
					pp.notification.msg <- Msg{
						Text:    " 💥  Build # " + strconv.Itoa(ver) + " failed for " + bra.Name + "\n" + err.Error(),
						BuildNo: ver,
					}

				}()
				cleanTheSource()
				log.Println(err)
				return
			}

			go func() {

				pp.notification.msg <- Msg{
					Text:    "📦 Container pushed to " + awsRegistry + ":" + buildVer + "  🏷 tagged '" + buildVer + "' ready to  ship 🛳 ",
					BuildNo: ver,
				}

				pp.notification.msg <- Msg{
					Text:    " 🍻 Build # " + strconv.Itoa(ver) + " Successful for " + bra.Name,
					BuildNo: ver,
				}

			}()
			cleanTheSource()

		}
	}
}

func prepareTheSource(repoName, branch string) error {
	log.Println(" 🚀  Cloning the repositiory  git@github.com:grapetechadmin/" + repoName + ".git")
	if err := runCmd("rm -rf *"); err != nil {
		return errors.New(" 👾  repository clone failed, temp directory error")
	}

	if err := runCmd("git clone git@github.com:grapetechadmin/" + repoName + ".git"); err != nil {
		return errors.New(" 👾  clone from github failed with error :" + err.Error())
	}

	if err := runCmdSetDir("git checkout "+branch, WorkDir+"/"+repoName); err != nil {
		return errors.New(" 👾  git checkout to '" + branch + "' branch failed :" + err.Error())
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

func pushTheImage(repoName, awsRegistry, buildNo string) error {
	log.Println(" 🚀  Push " + repoName + ":latest")

	if err := runCmd("docker tag " + repoName + ":latest " + awsRegistry + ":" + buildNo); err != nil {
		handleError(err)
		return errors.New(" 👾  docker tag failed, ")
	}

	log.Println(" 🚀  Tagged " + repoName + ":latest as " + awsRegistry + ":" + buildNo)

	if err := runCmd("docker push " + awsRegistry + ":" + buildNo); err != nil {
		handleError(err)
		return errors.New(" 👾  docker push failed")
	}

	log.Println(" 🚀  Pushed to registry 670907057868.dkr.ecr.us-east-2.amazonaws.com/" + repoName + ":" + buildNo)
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
