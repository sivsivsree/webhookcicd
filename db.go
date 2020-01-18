package webhookcicd

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"strconv"
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

func (db *DB) GetRepoName() string {
	return "dewa-test"
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
