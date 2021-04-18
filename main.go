package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Worker struct {
	ID     int
	URL    string
	Status Status

	Client *http.Client

	MessageCh chan<- string

	Logger *Logger
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
			if w.Status != OK {
				w.Status.Recovery()
				w.MessageCh <- "recovery"
				w.Logger.Info(w.ID, w.URL, fmt.Sprintf("status canged: %s", w.Status.String()))
			}
		} else {
			if w.Status != ALERT {
				w.Status.Trigger()
				w.MessageCh <- "alert"
				w.Logger.Info(w.ID, w.URL, fmt.Sprintf("status canged: %s", w.Status.String()))
			}
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
	errCh := make(chan error)
	alertSender := &AlertSender{
		MessageCh: messageCh,
		ErrCh:     errCh,
	}
	if slackConf := config.Notification.Slack; slackConf != nil {
		token := slackConf.Token
		channel := slackConf.Channel

		alertSender.SetNotifier(
			NewSlackNotifier(token, channel),
		)
	}
	go alertSender.Run()

	var targetURLs []string
	for _, t := range config.Target {
		targetURLs = append(targetURLs, t.URL)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for i, url := range targetURLs {
		id := i + 1

		client := http.DefaultClient
		client.Timeout = 10 * time.Second

		worker := &Worker{
			ID:     id,
			URL:    url,
			Status: OK,

			Client: client,

			MessageCh: messageCh,

			Logger: logger,
		}
		go worker.run(ctx)
	}

	for {
		select {
		case err := <-errCh:
			logger.Error(0, "", err.Error())
		case <-ctx.Done():
			stop()
			logger.Info(0, "", "Interrupt")
			return
		}
	}
}
