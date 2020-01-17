package main

import (
	"flag"
	"github.com/sivsivsree/webhookcicd"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var pipelineFile string
	flag.StringVar(&pipelineFile, "script", "bash.sh", "file to run when event happens")
	flag.Parse()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	_, srv := webhookcicd.NewServer()

	srv.SetPipeline(pipelineFile)
	srv.Start()

	<-done

	srv.Stop()
}
