package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cases := []struct {
		name   string
		config []byte
		want   *Config
	}{
		{
			name: "standard config",
			config: []byte(`[notification.slack]
token = "dummytoken"
channel = "#general"

[[monitor]]
name = "example.com check"
url = "https://example.com/check"
`),
			want: &Config{
				Notification: &Notification{
					Slack: &Slack{Token: "dummytoken", Channel: "#general"},
				},
				Monitors: []*Monitor{
					{
						Name:   "example.com check",
						Method: "GET",
						URL:    parseURL(t, "https://example.com/check"),
						Follow: false,
					},
				},
			},
		},
		{
			name: "multiple monitors",
			config: []byte(`[notification.slack]
token = "dummytoken"
channel = "#general"

[[monitor]]
name = "example.com check"
url = "https://example.com/check"

[[monitor]]
name = "example.com check 2"
url = "https://example.com/check2"
`),
			want: &Config{
				Notification: &Notification{
					Slack: &Slack{Token: "dummytoken", Channel: "#general"},
				},
				Monitors: []*Monitor{
					{
						Name:   "example.com check",
						Method: "GET",
						URL:    parseURL(t, "https://example.com/check"),
						Follow: false,
					},
					{
						Name:   "example.com check 2",
						Method: "GET",
						URL:    parseURL(t, "https://example.com/check2"),
						Follow: false,
					},
				},
			},
		},
		{
			name: "envvar",
			config: []byte(`[notification.slack]
token = '{{ env "TEST_SLACK_TOKEN" }}'
channel = "#general"

[[monitor]]
name = "example.com check"
url = "https://example.com/check"
`),
			want: &Config{
				Notification: &Notification{
					Slack: &Slack{Token: "envtoken", Channel: "#general"},
				},
				Monitors: []*Monitor{
					{
						Name:   "example.com check",
						Method: "GET",
						URL:    parseURL(t, "https://example.com/check"),
						Follow: false,
					},
				},
			},
		},
		{
			name: "method and follow",
			config: []byte(`[notification.slack]
token = "dummytoken"
channel = "#general"

[[monitor]]
name = "example.com post check"
method = "POST"
url = "https://example.com/check"

[[monitor]]
name = "example.com follow check"
url = "https://example.com/check"
follow = true
`),
			want: &Config{
				Notification: &Notification{
					Slack: &Slack{Token: "dummytoken", Channel: "#general"},
				},
				Monitors: []*Monitor{
					{
						Name:   "example.com post check",
						Method: "POST",
						URL:    parseURL(t, "https://example.com/check"),
						Follow: false,
					},
					{
						Name:   "example.com follow check",
						Method: "GET",
						URL:    parseURL(t, "https://example.com/check"),
						Follow: true,
					},
				},
			},
		},
	}

	if err := os.Setenv("TEST_SLACK_TOKEN", "envtoken"); err != nil {
		t.Fatal("set env failed")
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal("create temporary directory failed")
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
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
		})
	}
}
