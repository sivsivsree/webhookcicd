package webhookcicd

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strconv"
)

const (
	repoName   = "repo"
	branchName = "branch"
	aws_ecr    = "aws_ecr"
)

type DB struct {
	*leveldb.DB
}

func NewDB() (error, *DB) {
	db, err := leveldb.OpenFile(".siv", nil)
	if err != nil {
		log.Fatal("opening db -c", err)
	}

	return nil, &DB{db}
}

func (db *DB) LogON() string {
	return "ON"
}

func (db DB) GetBuildNo() int {
	buildNo, _ := db.Get([]byte("buildNO"), nil)
	// handleErrorMsg("GetBuildNo",err)
	ver := 0
	if buildNo != nil {
		newBuildNo, _ := strconv.Atoi(string(buildNo))
		newBuildNo = newBuildNo + 1
		buildNo = []byte(strconv.Itoa(newBuildNo))
		ver = newBuildNo
	}

	return ver
}

func (db *DB) BuildFinish() error {
	buildNo, err := db.Get([]byte("buildNO"), nil)
	if err != nil {
		buildNo = []byte("0")
	} else {
		newBuildNo, _ := strconv.Atoi(string(buildNo))
		newBuildNo = newBuildNo + 1
		buildNo = []byte(strconv.Itoa(newBuildNo))
	}

	return db.Put([]byte("buildNO"), buildNo, nil)
}

func (db *DB) SetRepo(repo string) error {
	return db.Put([]byte(repoName), []byte(repo), nil)
}

func (db *DB) SetBranch(branch string) error {
	return db.Put([]byte(branchName), []byte(branch), nil)
}

func (db *DB) SetECR(ecr string) error {
	return db.Put([]byte(aws_ecr), []byte(ecr), nil)
}

func (db *DB) GetBranch() string {
	branch, err := db.Get([]byte(branchName), nil)
	handleErrorMsg("GetBranch", err)
	return string(branch)
}

func (db *DB) GetRepoName() string {
	branch, err := db.Get([]byte(repoName), nil)
	handleErrorMsg("GetRepoName", err)
	return string(branch)
}

func (db *DB) GetECR() string {
	branch, err := db.Get([]byte(aws_ecr), nil)
	handleErrorMsg("GetECR", err)
	return string(branch)
}
