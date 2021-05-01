package main

import (
	gc "github.com/kayac/go-config"
)

type Config struct {
	Notification *Notification `toml:"notification"`
	Monitors     []*Monitor    `toml:"monitor"`
}

type Notification struct {
	Slack *Slack `toml:"slack"`
}

type Slack struct {
	Token   string `toml:"token"`
	Channel string `toml:"channel"`
}

func LoadConfig(filename string) (*Config, error) {
	config := Config{}

	if err := gc.LoadWithEnvTOML(&config, filename); err != nil {
		return nil, err
	}

	for _, m := range config.Monitors {
		if m.Method == "" {
			m.Method = "GET"
		}
	}

	return &config, nil
}
