package main

import (
	"net/url"
	"os"

	toml "github.com/pelletier/go-toml"
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
	URL URL `toml:"url"`
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

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
