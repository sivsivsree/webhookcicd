package webhookcicd

import (
	"context"
	"fmt"
	"github.com/sivsivsree/webhookcicd/internal"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewApiServer(srv GrpcServer) *grpc.Server {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	internal.RegisterConfigServiceServer(grpcServer, srv)
	go grpcServer.Serve(lis)
	return grpcServer
}

// Server represents the gRPC server
type GrpcServer struct {
	DB *DB
}

// ChangeConfig generates response to request
func (s GrpcServer) ChangeConfig(ctx context.Context, in *internal.Config) (*internal.Config, error) {
	log.Println("Receive ChangeConfig", in)

	ecr := s.DB.GetECR()
	branch := s.DB.GetBranch()
	repoName := s.DB.GetRepoName()

	var err error
	if in.ECR != "" && in.Branch != "" && in.Repo != "" {
		err = s.DB.SetECR(in.ECR)
		ecr = in.ECR

		err = s.DB.SetBranch(in.Branch)
		branch = in.Branch

		err = s.DB.SetBranch(in.Branch)
		repoName = in.Repo
	} else {
		errors.New("All are critical configurations")
	}

	return &internal.Config{ECR: ecr, Branch: branch, Repo: repoName}, err
}

func (s GrpcServer) GetConfig(ctx context.Context, in *internal.Config) (*internal.Config, error) {

	log.Println("GetConfig called", in)

	ecr := s.DB.GetECR()
	branch := s.DB.GetBranch()
	repoName := s.DB.GetRepoName()
	return &internal.Config{ECR: ecr, Branch: branch, Repo: repoName}, nil
}
