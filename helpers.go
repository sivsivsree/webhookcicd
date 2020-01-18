package webhookcicd

import (
	"fmt"
	"log"
	"os"
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
