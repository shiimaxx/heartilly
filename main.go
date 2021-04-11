package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	toml "github.com/pelletier/go-toml"
	"go.uber.org/zap"
)

var logger *zap.Logger

var (
	StatusOK    = 0
	StatusAlert = 1
)

type Config struct {
	Target []TargetConfig
}

type TargetConfig struct {
	URL string
}

type Worker struct {
	ID     int
	URL    string
	Status int

	Client *http.Client
}

func (w *Worker) run(ctx context.Context) {
	// jitter
	jitter := rand.Intn(10)
	time.Sleep(time.Duration(jitter) * time.Second)

	logger.Info("start worker", zap.Int("id", w.ID), zap.String("target", w.URL))

	c := time.Tick(1 * time.Minute)
	for {
		logger.Info("check", zap.Int("id", w.ID), zap.String("target", w.URL))

		if ok, err := w.check(ctx); ok && err == nil {
			if w.Status != StatusOK {
				w.Status = StatusOK
				logger.Info("recovery")
			}
			logger.Debug("ok")
		} else {
			if w.Status != StatusAlert {
				w.Status = StatusAlert
				logger.Info("alert")
			}
			logger.Debug("failed")
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
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logger = l
	defer logger.Sync()

	config, err := toml.LoadFile("conf/config.toml")
	if err != nil {
		panic(err)
	}

	targetConfig := config.Get("target").([]*toml.Tree)

	var targetURLs []string
	for _, tc := range targetConfig {
		targetURLs = append(targetURLs, tc.Get("url").(string))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for i, url := range targetURLs {
		id := i

		client := http.DefaultClient
		client.Timeout = 10 * time.Second

		worker := &Worker{
			ID:     id,
			URL:    url,
			Status: StatusOK,

			Client: client,
		}
		go worker.run(ctx)
	}

	select {
	case <-ctx.Done():
		stop()
		logger.Info("Interrupt")
	}
}
