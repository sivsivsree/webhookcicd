package webhookcicd

import (
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
)

type DB struct {
	*leveldb.DB
}

func (db *DB) GetBuildNo() string {
	buildNo, _ := db.Get([]byte("buildNO"), nil)
	// handleErrorMsg("GetBuildNo",err)
	ver := "build-0"
	if buildNo != nil {
		newBuildNo, _ := strconv.Atoi(string(buildNo))
		newBuildNo = newBuildNo + 1
		buildNo = []byte(strconv.Itoa(newBuildNo))
		ver = "build-" + strconv.Itoa(newBuildNo)
	}

	return ver
}

func (db *DB) BuildFinish() error {
	buildNo, err := db.Get([]byte("buildNO"), nil)

	if err != nil {
		return err
	}

	if buildNo != nil {
		newBuildNo, _ := strconv.Atoi(string(buildNo))
		newBuildNo = newBuildNo + 1
		buildNo = []byte(strconv.Itoa(newBuildNo))
	} else {
		buildNo = []byte("0")
	}

	return db.Put([]byte("buildNO"), buildNo, nil)
}
