package main

import (
	"net/url"

	gc "github.com/kayac/go-config"
)

type Config struct {
	Notification *Notification `toml:"notification"`
	Target       []*Target     `toml:"target"`
}

type Notification struct {
	Slack *Slack `toml:"slack"`
}

type Slack struct {
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
}

type Target struct {
	Name   string `toml:"name"`
	Method string `toml:"method"`
	URL    URL    `toml:"url"`
	Follow bool   `toml:"follow"`
}

type URL url.URL

func (u *URL) UnmarshalText(text []byte) error {
	parsedURL, err := url.Parse(string(text))
	if err != nil {
		return err
	}

	*u = URL(*parsedURL)

	return nil
}

func (u *URL) String() string {
	return (*url.URL)(u).String()
}

func LoadConfig(filename string) (*Config, error) {
	config := Config{}

	if err := gc.LoadWithEnvTOML(&config, filename); err != nil {
		return nil, err
	}

	for _, t := range config.Target {
		if t.Method == "" {
			t.Method = "GET"
		}
	}

	return &config, nil
}
