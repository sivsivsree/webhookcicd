package webhookcicd

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type slack struct {
	endpoint string
	msg      chan Msg
}

func NewSlack() *slack {
	msgCh := make(chan Msg)
	return &slack{endpoint: SlackEndpoint, msg: msgCh}
}

type Msg struct {
	Text    string
	BuildNo int
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

func (s *slack) Start() {
	log.Println("notification worker started")
	for {
		select {
		case msg := <-s.msg:
			err := s.SendNotification(msg.Text)
			handleErrorMsg("[SendNotification]", err)
		}
	}
}

func (s *slack) SendNotification(msg string) error {

	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, s.endpoint, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
