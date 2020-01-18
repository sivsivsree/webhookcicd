package main

import (
	"fmt"
	"github.com/sivsivsree/webhookcicd"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	_, db := webhookcicd.NewDB()
	_ = db.SetRepo("dewa-test")
	_ = db.SetBranch("master")
	_ = db.SetECR("<>")

	_, srv := webhookcicd.NewServer()

	srv.SetWorkDir()
	srv.SetPipeline(db)
	srv.SetSecret(os.Getenv("SECRET"))
	srv.Start()

	grpcServer := webhookcicd.NewApiServer(webhookcicd.GrpcServer{DB: db})

	fmt.Println("waiting for clicd connection")

	<-done

	_ = db.Close()
	grpcServer.GracefulStop()
	srv.Stop()
}
