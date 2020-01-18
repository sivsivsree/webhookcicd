package main

import (
	"flag"
	"github.com/sivsivsree/webhookcicd"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var pipelineFile string
	flag.StringVar(&pipelineFile, "script", "bash.siv", "file to run when event happens")
	flag.Parse()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	db, err := leveldb.OpenFile("siv", nil)
	if err != nil {
		log.Fatal("opening db -c", err)
	}
	_, srv := webhookcicd.NewServer()

	srv.SetWorkDir()
	srv.SetPipeline(&webhookcicd.DB{DB: db}, pipelineFile)
	srv.SetSecret("my-hook")
	srv.Start()

	<-done

	db.Close()
	srv.Stop()
}
