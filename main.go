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

		result := &Result{
			CheckedAt:   time.Now().UTC(),
			Reason:    reason,
			MonitorID: w.Probe.Monitor.ID,
		}

		if err == nil {
			if ok && (!w.Status.Is(OK)) {
				w.Status.Recovery()

				result.Status = w.Status.String()
				if err := CreateResult(result); err != nil {
					w.Logger.Error(
						w.ID,
						w.Probe.Monitor.URL.String(),
						fmt.Sprintf("save result failed: %s", err.Error()),
					)
				}

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

				result.Status = w.Status.String()
				if err := CreateResult(result); err != nil {
					w.Logger.Error(
						w.ID,
						w.Probe.Monitor.URL.String(),
						fmt.Sprintf("save result failed: %s", err.Error()),
					)
				}

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

				result.Status = w.Status.String()
				if err := CreateResult(result); err != nil {
					w.Logger.Error(
						w.ID,
						w.Probe.Monitor.URL.String(),
						fmt.Sprintf("save result failed: %s", err.Error()),
					)
				}

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

	if err := OpenDB(config.DBFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	logger.Info(0, "", fmt.Sprint("open dbfile: ", config.DBFile))

	monitors, err := InitSyncMonitor(config.Monitors)
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

	for i, m := range monitors {
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

	httpSrv := NewHTTPServer()
	go func() {
		if err := httpSrv.Start(":8000"); err != nil {
			errCh <- err
		}
	}()

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
