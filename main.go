package main

import (
	"context"
	"fmt"
	"math/rand"
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
	Status Status

	Probe *Probe

	MessageCh chan<- Message

	Logger *Logger
}

func (w *Worker) run(ctx context.Context) {
	// jitter
	rand.Seed(time.Now().UnixNano())
	jitter := rand.Intn(10)
	time.Sleep(time.Duration(jitter) * time.Second)

	w.Logger.Info(w.ID, w.Probe.Monitor.URL.String(), "start worker")

	c := time.Tick(1 * time.Minute)
	for {
		w.Logger.Info(w.ID, w.Probe.Monitor.URL.String(), "check")

		ok, reason, err := w.Probe.Check(ctx)
		if err == nil {
			if ok && (!w.Status.Is(OK)) {
				w.Status.Recovery()
				w.MessageCh <- Message{
					Text: fmt.Sprintf("%s: %s\n%s - %s",
						w.Status.String(),
						w.Probe.Monitor.Name,
						w.Probe.Monitor.URL.String(),
						reason,
					),
					StatusType: OK,
				}
				w.Logger.Info(
					w.ID,
					w.Probe.Monitor.URL.String(),
					fmt.Sprintf("status canged: %s", w.Status.String()),
				)

			} else if !ok && (w.Status.Is(OK)) {
				w.Status.Trigger()
				w.MessageCh <- Message{
					Text: fmt.Sprintf("%s: %s\n%s - %s",
						w.Status.String(),
						w.Probe.Monitor.Name,
						w.Probe.Monitor.URL.String(),
						reason,
					),
					StatusType: Critical,
				}
				w.Logger.Info(
					w.ID,
					w.Probe.Monitor.URL.String(),
					fmt.Sprintf("status canged: %s", w.Status.String()),
				)
			}
		} else {
			if !w.Status.Is(Unknown) {
				w.Status.Unknown()
				w.MessageCh <- Message{
					Text: fmt.Sprintf("%s: %s\n%s - %s",
						w.Status.String(),
						w.Probe.Monitor.Name,
						w.Probe.Monitor.URL.String(),
						reason,
					),
					StatusType: Unknown,
				}
				w.Logger.Info(
					w.ID,
					w.Probe.Monitor.URL.String(),
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

	for i, m := range config.Monitors {
		id := i + 1

		p := &Probe{
			Monitor: m,
		}

		worker := &Worker{
			ID:     id,
			Status: OK,

			Probe: p,

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
