package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	toml "github.com/pelletier/go-toml"
)

type Config struct {
	Target []TargetConfig
}

type TargetConfig struct {
	URL string
}

type Worker struct {
	URL     string
	Timeout int
}

func (w *Worker) run(ctx context.Context) {
	// jitter
	jitter := rand.Intn(10)
	time.Sleep(time.Duration(jitter) * time.Second)

	fmt.Printf("start worker: %s\n", w.URL)

	c := time.Tick(1 * time.Minute)
	for {
		fmt.Printf("check: %s\n", w.URL)
		time.Sleep(time.Second * 10)

		select {
		case <-c:
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	config, err := toml.LoadFile("config.toml")
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

	for _, url := range targetURLs {
		worker := &Worker{URL: url, Timeout: 10}
		go worker.run(ctx)
	}

	select {
	case <-ctx.Done():
		stop()
		fmt.Println("Interrupt")
	}
}
