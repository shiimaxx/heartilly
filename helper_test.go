package main

import (
	"net/url"
	"testing"
)

func parseURL(t *testing.T, u string) URL {
	t.Helper()

	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatal(err)
	}

	return URL(*parsed)
}
