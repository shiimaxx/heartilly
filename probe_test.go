package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/ok")
		w.WriteHeader(http.StatusMovedPermanently)
	})

	ts := httptest.NewServer(mux)
	return ts
}

func TestProbe_Check(t *testing.T) {
	cases := []struct {
		name   string
		method string
		path   string
		follow bool

		wantResult bool
		wantReason string
		wantErr    error
	}{
		{
			name:   "ok",
			method: "GET",
			path:   "/ok",
			follow: false,

			wantResult: true,
			wantReason: "200 OK",
			wantErr:    nil,
		},
		{
			name:   "error",
			method: "GET",
			path:   "/error",
			follow: false,

			wantResult: false,
			wantReason: "500 Internal Server Error",
			wantErr:    nil,
		},
		{
			name:   "redirect",
			method: "GET",
			path:   "/redirect",
			follow: false,

			wantResult: true,
			wantReason: "301 Moved Permanently",
			wantErr:    nil,
		},
		{
			name:   "follow redirect",
			method: "GET",
			path:   "/redirect",
			follow: true,

			wantResult: true,
			wantReason: "200 OK",
			wantErr:    nil,
		},
	}

	ts := newTestServer()
	defer ts.Close()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			target := &Target{
				Method: c.method,
				URL:    parseURL(t, fmt.Sprintf("%s%s", ts.URL, c.path)),
				Follow: c.follow,
			}
			probe := &Probe{Target: target}
			result, reason, err := probe.Check(context.TODO())

			assert.Equal(t, c.wantResult, result)
			assert.Equal(t, c.wantReason, reason)
			assert.Equal(t, c.wantErr, err)
		})
	}
}
