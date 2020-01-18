package main

import (
	"context"
	"fmt"
	"github.com/sivsivsree/webhookcicd/internal"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	cc := internal.NewConfigServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	in := &internal.Config{
		Repo:   "dewa-test",
		Branch: "master",
		ECR:    "<>",
	}
	data, err := cc.ChangeConfig(ctx, in)

	fmt.Println(data, err)

}
