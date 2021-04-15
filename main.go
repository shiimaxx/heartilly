package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/slack-go/slack"
)

var (
	StatusOK    = 0
	StatusAlert = 1
)

type Worker struct {
	ID     int
	URL    string
	Status int

	Client *http.Client

	MessageCh chan<- string

	Logger *Logger
}

type AlertSender struct {
	Notifiers []Notifier

	MessageCh <-chan string
}

func (as *AlertSender) Run() error {
	for {
		msg := <-as.MessageCh

		for _, notifier := range as.Notifiers {
			if err := notifier.Notify(msg); err != nil {
				return err
			}
		}
	}
}

type Notifier interface {
	Notify(string) error
}

type SlackNotifier struct {
	Channel string
	Client  *slack.Client
}

func (s *SlackNotifier) Notify(msg string) error {
	_, _, err := s.Client.PostMessage(
		s.Channel,
		slack.MsgOptionText(msg, false),
	)
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) run(ctx context.Context) {
	// jitter
	jitter := rand.Intn(10)
	time.Sleep(time.Duration(jitter) * time.Second)

	w.Logger.Info(w.ID, w.URL, "start worker")

	c := time.Tick(1 * time.Minute)
	for {
		w.Logger.Info(w.ID, w.URL, "check")

		if ok, err := w.check(ctx); ok && err == nil {
			if w.Status != StatusOK {
				w.Status = StatusOK
				w.MessageCh <- "recovery"
				w.Logger.Info(w.ID, w.URL, "recovery")
			}
			w.Logger.Debug(w.ID, w.URL, "ok")
		} else {
			if w.Status != StatusAlert {
				w.Status = StatusAlert
				w.MessageCh <- "alert"
				w.Logger.Info(w.ID, w.URL, "alert")
			}
			w.Logger.Debug(w.ID, w.URL, "failed")
		}

		select {
		case <-c:
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) check(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.URL, nil)
	if err != nil {
		return false, err
	}

	resp, err := w.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func main() {
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}

	config, err := LoadConfig("conf/config.toml")
	if err != nil {
		panic(err)
	}

	messageCh := make(chan string)
	alertSender := &AlertSender{
		Notifiers: []Notifier{
			&SlackNotifier{
				Channel: config.Notification.Slack.Channel,
				Client:  slack.New(config.Notification.Slack.Token),
			},
		},
		MessageCh: messageCh,
	}
	go alertSender.Run()

	var targetURLs []string
	for _, t := range config.Target {
		targetURLs = append(targetURLs, t.URL)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for i, url := range targetURLs {
		id := i+1

		client := http.DefaultClient
		client.Timeout = 10 * time.Second

		worker := &Worker{
			ID:     id,
			URL:    url,
			Status: StatusOK,

			Client: client,

			MessageCh: messageCh,

			Logger: logger,
		}
		go worker.run(ctx)
	}

	select {
	case <-ctx.Done():
		stop()
		logger.Info(0, "", "Interrupt")
	}
}
