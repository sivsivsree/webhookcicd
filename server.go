package webhookcicd

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"log"
	"net/http"
	"os"
	"time"
)

var WorkDir = os.TempDir() + "cicd"

type service struct {
	git      *http.Server
	pipeline *pipeline
	secret   []byte
}

func NewServer() (error, *service) {

	server := &http.Server{Addr: ":8080", Handler: nil}
	return nil, &service{git: server, secret: []byte("")}

}
func (s *service) SetPipeline(db *DB) *service {
	err, pline := newPipeline(db)
	handleError(err)
	s.pipeline = pline
	return s
}

func (s *service) SetSecret(secret string) *service {
	if len(secret) < 0 {
		log.Println("no secret for github hook provided")
	}
	s.secret = []byte(secret)
	return s
}

func (s *service) SetWorkDir() {
	log.Println("current working dir", WorkDir)
	handleErrorMsg("[NewServer]", os.MkdirAll(WorkDir, os.ModePerm))
}

func (s *service) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", s.handleWebhook)
	s.git.Handler = mux

	go func() {
		err := s.git.ListenAndServe()
		log.Fatal(err)
	}()
	log.Println("git listener started")
	go s.pipeline.StartWorker()
	log.Println("worker started")
}

func (s *service) Stop() {
	log.Println("Shutting Down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	_ = s.git.Shutdown(ctx)
}

func (s service) handleWebhook(w http.ResponseWriter, r *http.Request) {

	payload, err := github.ValidatePayload(r, s.secret)
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.PushEvent:

		ref := *e.Ref
		branch := ref[len("refs/heads/"):]
		if branch == "master" { // todo: can set as arg branch set master
			s.pipeline.branch <- BranchUpdate{
				Name: branch,
				SHA:  *e.After,
			}
		}

	case *github.PingEvent:

	case *github.PullRequestEvent:
		// this is a pull request, do something with it
	case *github.WatchEvent:
		// https://developer.github.com/v3/activity/events/types/#watchevent
		// someone starred our repository
		if e.Action != nil && *e.Action == "starred" {
			fmt.Printf("%s starred repository %s\n",
				*e.Sender.Login, *e.Repo.FullName)
		}
	default:
		log.Printf("unknown WebHookType: %s, webhook-id: %s skipping\n", github.WebHookType(r), r.Header.Get("X-GitHub-Delivery"))

		return
	}
}
