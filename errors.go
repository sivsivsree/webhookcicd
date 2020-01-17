package webhookcicd

import "log"

func handleError(err error) {
	if err != nil {
		log.Printf("[error] err=%s\n", err)
		return
	}
}
