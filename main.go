package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Message struct {
	Text       string
	StatusType Status
}

type Worker struct {
	ID         int
	URL        string
	Status     Status

	Client *http.Client

	MessageCh chan<- Message

	Logger *Logger
}

func (w *Worker) run(ctx context.Context) {
	// jitter
	rand.Seed(time.Now().UnixNano())
	jitter := rand.Intn(10)
	time.Sleep(time.Duration(jitter) * time.Second)

	w.Logger.Info(w.ID, w.URL, "start worker")

	c := time.Tick(1 * time.Minute)
	for {
		w.Logger.Info(w.ID, w.URL, "check")

		ok, err := w.check(ctx)
		if err == nil {
			if ok && (!w.Status.Is(Initial) && !w.Status.Is(OK)) {
				w.Status.Recovery()
				w.MessageCh <- Message{
					Text:       fmt.Sprintf("%s: %s", w.Status.String(), w.URL),
					StatusType: OK,
				}
				w.Logger.Info(w.ID, w.URL, fmt.Sprintf("status canged: %s", w.Status.String()))
			} else if !ok && (w.Status.Is(Initial) || w.Status.Is(OK)) {
				w.Status.Trigger()
				w.MessageCh <- Message{
					Text:       fmt.Sprintf("%s: %s", w.Status.String(), w.URL),
					StatusType: Alert,
				}
				w.Logger.Info(w.ID, w.URL, fmt.Sprintf("status canged: %s", w.Status.String()))
			}
		} else {
			if !w.Status.Is(Unknown) {
				w.Status.Unknown()
				w.MessageCh <- Message{
					Text:       fmt.Sprintf("%s: %s", w.Status.String(), w.URL),
					StatusType: Unknown,
				}
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
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return false, nil
		}
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

	messageCh := make(chan Message)
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
		client.Timeout = 15 * time.Second

		worker := &Worker{
			ID:     id,
			URL:    url,
			Status: Initial,

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
