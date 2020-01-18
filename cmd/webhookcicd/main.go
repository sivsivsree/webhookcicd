package main

import (
	"github.com/sivsivsree/webhookcicd"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	_, db := webhookcicd.NewDB()

	_, srv := webhookcicd.NewServer()

	srv.SetWorkDir()
	srv.SetPipeline(db)
	srv.SetSecret(os.Getenv("SECRET"))
	srv.Start()

	<-done

	_ = db.Close()
	srv.Stop()
}
