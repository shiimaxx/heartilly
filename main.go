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

	"github.com/jessevdk/go-flags"
)

type Message struct {
	Text       string
	StatusType Status
}

type Worker struct {
	ID     int
	Target *Target
	Status Status

	Client *http.Client

	MessageCh chan<- Message

	Logger *Logger
}

func (w *Worker) run(ctx context.Context) {
	// jitter
	rand.Seed(time.Now().UnixNano())
	jitter := rand.Intn(10)
	time.Sleep(time.Duration(jitter) * time.Second)

	w.Logger.Info(w.ID, w.Target.URL.String(), "start worker")

	c := time.Tick(1 * time.Minute)
	for {
		w.Logger.Info(w.ID, w.Target.URL.String(), "check")

		ok, err := w.check(ctx)
		if err == nil {
			if ok && (!w.Status.Is(OK)) {
				w.Status.Recovery()
				w.MessageCh <- Message{
					Text: fmt.Sprintf("%s: %s\n%s - %s",
						w.Status.String(),
						w.Target.Name,
						w.Target.URL.String(),
						"200 OK", // TODO: reasone
					),
					StatusType: OK,
				}
				w.Logger.Info(
					w.ID,
					w.Target.URL.String(),
					fmt.Sprintf("status canged: %s", w.Status.String()),
				)

			} else if !ok && (w.Status.Is(OK)) {
				w.Status.Trigger()
				w.MessageCh <- Message{
					Text: fmt.Sprintf("%s: %s\n%s - %s",
						w.Status.String(),
						w.Target.Name,
						w.Target.URL.String(),
						"500 Internal Server Error", // TODO: reasone
					),
					StatusType: Critical,
				}
				w.Logger.Info(
					w.ID,
					w.Target.URL.String(),
					fmt.Sprintf("status canged: %s", w.Status.String()),
				)
			}
		} else {
			if !w.Status.Is(Unknown) {
				w.Status.Unknown()
				w.MessageCh <- Message{
					Text: fmt.Sprintf("[%s]\n%s: %s",
						w.Target.Name,
						w.Status.String(),
						w.Target.URL.String(),
					),
					StatusType: Unknown,
				}
				w.Logger.Info(
					w.ID,
					w.Target.URL.String(),
					fmt.Sprintf("status canged: %s", w.Status.String()),
				)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.Target.URL.String(), nil)
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

type Options struct {
	Config string `short:"c" long:"config" default:"config.toml" description:"configuration file"`
}

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	logger, err := NewLogger()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	config, err := LoadConfig(opts.Config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for i, t := range config.Target {
		id := i + 1

		client := http.DefaultClient
		client.Timeout = 15 * time.Second

		worker := &Worker{
			ID:     id,
			Target: t,
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
			os.Exit(0)
		}
	}
}
