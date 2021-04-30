package main

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Probe struct {
	Target *Target
}

func (p *Probe) Check(ctx context.Context) (bool, string, error) {
	client := http.DefaultClient
	if p.Target.Follow {
		client.CheckRedirect = nil
	} else {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	client.Timeout = 15 * time.Second

	req, err := http.NewRequestWithContext(ctx, p.Target.Method, p.Target.URL.String(), nil)
	if err != nil {
		return false, "error", err
	}

	resp, err := client.Do(req)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return false, "timeout", nil
		}
		return false, "error", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return false, resp.Status, nil
	}

	return true, resp.Status, nil
}
