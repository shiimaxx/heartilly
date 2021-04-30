package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProbe_Check(t *testing.T) {
	cases := []struct {
		f http.HandlerFunc

		name   string
		method string
		follow bool

		wantResult bool
		wantReason string
		wantErr    error
	}{
		{
			f: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "ok")
			}),

			name:   "ok",
			method: "GET",
			follow: false,

			wantResult: true,
			wantReason: "200 OK",
			wantErr:    nil,
		},
		{
			f: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "error")
			}),

			name:   "error",
			method: "GET",
			follow: false,

			wantResult: false,
			wantReason: "500 Internal Server Error",
			wantErr:    nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ts := httptest.NewServer(c.f)
			defer ts.Close()

			target := &Target{
				Method: c.method,
				URL:    parseURL(t, ts.URL),
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
