package webhookcicd

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type service struct {
	git      *http.Server
	pipeline *pipeline
}

func NewServer() (error, *service) {
	server := &http.Server{Addr: ":8080", Handler: nil}
	return nil, &service{git: server,}

}
func (s *service) SetPipeline(file string) *service {
	err, pline := newPipeline(file)
	handleError(err)
	s.pipeline = pline
	return s
}
func (s *service) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", s.handleWebhook)
	s.git.Handler = mux
	log.Println("git listener started")
	go s.git.ListenAndServe()
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
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: err=%s\n", err)
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

		s.pipeline.Run()
		// this is a commit push, do something with it
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

		log.Printf("unknown event type %s\n", github.WebHookType(r))
		s.pipeline.Run()
		return
	}
}
