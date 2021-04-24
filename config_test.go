package main

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseURL(t *testing.T, u string) URL {
	t.Helper()

	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatal(err)
	}

	return URL(*parsed)
}

func TestLoadConfig(t *testing.T) {
	cases := []struct {
		config []byte
		want   *Config
	}{
		{
			config: []byte(`[notification.slack]
token = "dummytoken"
channel = "#general"

[[target]]
url = "https://example.com/check"
`),
			want: &Config{
				Notification: &Notification{
					Slack: &Slack{Token: "dummytoken", Channel: "#general"},
				},
				Target: []*Target{
					{URL: parseURL(t, "https://example.com/check")},
				},
			},
		},
		{
			config: []byte(`[notification.slack]
token = "dummytoken"
channel = "#general"

[[target]]
url = "https://example.com/check"

[[target]]
url = "https://example.com/check2"
`),
			want: &Config{
				Notification: &Notification{
					Slack: &Slack{Token: "dummytoken", Channel: "#general"},
				},
				Target: []*Target{
					{URL: parseURL(t, "https://example.com/check")},
					{URL: parseURL(t, "https://example.com/check2")},
				},
			},
		},
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal("create temporary directory failed")
	}
	for _, c := range cases {
		f, err := os.CreateTemp(tmpDir, "")
		if err != nil {
			t.Fatal("create temporary file failed")
		}

		if err := os.WriteFile(f.Name(), c.config, os.ModeTemporary); err != nil {
			t.Fatal("write file failed")
		}

		config, err := LoadConfig(f.Name())
		assert.Nil(t, err)
		assert.Equal(t, c.want, config)
	}
}
